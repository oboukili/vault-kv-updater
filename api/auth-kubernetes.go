// sourced from https://github.com/sethvargo/vault-kubernetes-authenticator
// Thanks Seth! :)

package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"golang.org/x/net/http2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	vaultAddr         string
	vaultCaPem        string
	vaultCaCert       string
	vaultCaPath       string
	vaultNamespace    string
	vaultSkipVerify   bool
	vaultServerName   string
	vaultK8SMountPath string
)

func AuthKubernetes() (token string, accessor string, err error) {
	vaultAddr = os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		vaultAddr = "http://127.0.0.1:8200"
	}

	role := os.Getenv("VAULT_ROLE")
	if role == "" {
		log.Fatal("missing VAULT_ROLE")
	}

	vaultCaPem = os.Getenv("VAULT_CAPEM")
	vaultCaCert = os.Getenv("VAULT_CACERT")
	vaultCaPath = os.Getenv("VAULT_CAPATH")
	vaultNamespace = os.Getenv("VAULT_NAMESPACE")
	vaultServerName = os.Getenv("VAULT_TLS_SERVER_NAME")

	if s := os.Getenv("VAULT_SKIP_VERIFY"); s != "" {
		b, err := strconv.ParseBool(s)
		if err != nil {
			log.Fatal(err)
		}
		vaultSkipVerify = b
	}

	vaultK8SMountPath = os.Getenv("VAULT_K8S_MOUNT_PATH")
	if vaultK8SMountPath == "" {
		vaultK8SMountPath = "kubernetes"
	}

	saPath := os.Getenv("SERVICE_ACCOUNT_PATH")
	if saPath == "" {
		saPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	}

	// Read the JWT token from disk
	jwt, err := ReadJwtToken(saPath)
	if err != nil {
		err = fmt.Errorf("could not read Kubernetes SA JWT token from disk: %s", err)
		return
	}

	// Authenticate to vault using the jwt token
	token, accessor, err = authenticate(role, jwt)
	if err != nil {
		err = fmt.Errorf("could not authenticate to Vault: %s", err)
	}
	return
}

func ReadJwtToken(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read jwt token: %s", err)
	}

	return string(bytes.TrimSpace(data)), nil
}

func authenticate(role, jwt string) (string, string, error) {
	// Setup the TLS (especially required for custom CAs)
	rootCAs, err := rootCAs()
	if err != nil {
		return "", "", err
	}

	tlsClientConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    rootCAs,
	}

	if vaultSkipVerify {
		tlsClientConfig.InsecureSkipVerify = true
	}

	if vaultServerName != "" {
		tlsClientConfig.ServerName = vaultServerName
	}

	transport := &http.Transport{
		TLSClientConfig: tlsClientConfig,
	}

	if err := http2.ConfigureTransport(transport); err != nil {
		return "", "", fmt.Errorf("failed to configure http2")
	}

	client := &http.Client{
		Transport: transport,
	}

	transport.Proxy = http.ProxyFromEnvironment

	addr := vaultAddr + "/v1/auth/" + vaultK8SMountPath + "/login"
	body := fmt.Sprintf(`{"role": "%s", "jwt": "%s"}`, role, jwt)

	req, err := http.NewRequest(http.MethodPost, addr, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %s", err)
	}
	if vaultNamespace != "" {
		req.Header.Set("X-Vault-Namespace", vaultNamespace)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to login %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var b bytes.Buffer
		if _, err := io.Copy(&b, resp.Body); err != nil {
			log.Printf("failed to copy response body: %s", err)
		}
		return "", "", fmt.Errorf("failed to get successful response: %#v, %s",
			resp, b.String())
	}

	var s struct {
		Auth struct {
			ClientToken    string `json:"client_token"`
			ClientAccessor string `json:"accessor"`
		} `json:"auth"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return "", "", fmt.Errorf("failed to read body: %s", err)
	}

	return s.Auth.ClientToken, s.Auth.ClientAccessor, nil
}

// rootCAs returns the list of trusted root CAs based off the provided
// configuration. If no CAs were specified, the system roots are used.
func rootCAs() (*x509.CertPool, error) {
	switch {
	case vaultCaPem != "":
		pool := x509.NewCertPool()
		if err := loadCert(pool, []byte(vaultCaPem)); err != nil {
			return nil, err
		}
		return pool, nil
	case vaultCaCert != "":
		pool := x509.NewCertPool()
		if err := loadCertFile(pool, vaultCaCert); err != nil {
			return nil, err
		}
		return pool, nil
	case vaultCaPath != "":
		pool := x509.NewCertPool()
		if err := loadCertFolder(pool, vaultCaPath); err != nil {
			return nil, err
		}
		return pool, nil
	default:
		pool, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("failed to load system certs %s", err)
		}
		return pool, err
	}
}

// loadCert loads a single pem-encoded certificate into the given pool.
func loadCert(pool *x509.CertPool, pem []byte) error {
	if ok := pool.AppendCertsFromPEM(pem); !ok {
		return fmt.Errorf("failed to parse PEM")
	}
	return nil
}

// loadCertFile loads the certificate at the given path into the given pool.
func loadCertFile(pool *x509.CertPool, p string) error {
	pem, err := ioutil.ReadFile(p)
	if err != nil {
		return fmt.Errorf("failed to read CA file from disk %s", err)
	}

	if err := loadCert(pool, pem); err != nil {
		return fmt.Errorf("failed to load CA at %s: %s", p, err)
	}

	return nil
}

// loadCertFolder iterates exactly one level below the given directory path and
// loads all certificates in that path. It does not recurse
func loadCertFolder(pool *x509.CertPool, p string) error {
	if err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return loadCertFile(pool, path)
	}); err != nil {
		return fmt.Errorf("failed to load CAs at %s: %s", p, err)
	}
	return nil
}
