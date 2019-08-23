package api

import (
	"encoding/json"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"log"
	"os"
	"reflect"
	"strings"
)

func VaultClientInit() (c *vault.Client, err error) {
	var token string
	var vaultAddress string

	// initialize a new Vault client
	vaultAddress, ok := os.LookupEnv(EnvVaultAddr)
	if !ok {
		vaultAddress = vaultDefaultAddr
	}

	c, err = vault.NewClient(&vault.Config{
		Address: vaultAddress,
	})
	if err != nil {
		return
	}

	authMethod, ok := os.LookupEnv(EnvVaultAuthMethod)
	if !ok {
		authMethod = vaultDefaultAuthenticationMethod
	}

	switch authMethod {
	case "kubernetes":
		log.Println("Kubernetes authentication method selected")
		token, _, err = AuthKubernetes()
	case "token":
		log.Println("Token authentication method selected")
		token, err = AuthToken()
	default:
		err = fmt.Errorf("%s is not supported as an authentication method, choose between kubernetes,token", authMethod)
	}
	if err != nil {
		return
	}
	c.SetToken(token)
	c.SetLimiter(10.0, 30)
	return
}

func VaultSecretDataIsDifferent(newData map[string]interface{}, vaultSecret *vault.Secret, kvVersion int) (b bool, err error) {
	if vaultSecret == nil {
		return true, nil
	}

	// Use the same Json Marshaller instead of the Vault API Marshaller (integers are not of the same types)
	var vaultSecretDataBytes []byte
	var decodedVaultSecret interface{}

	switch kvVersion {
	case 1:
		vaultSecretDataBytes, err = json.Marshal(vaultSecret.Data)
		if err != nil {
			return false, err
		}
		err = json.Unmarshal(vaultSecretDataBytes, &decodedVaultSecret)
		if err != nil {
			return
		}
		switch decodedVaultSecret {
		case nil:
		default:
			if !reflect.DeepEqual(newData, decodedVaultSecret) {
				b = true
			}
		}
	case 2:
		_, exists := vaultSecret.Data["data"]
		switch exists {
		case true:
			vaultSecretDataBytes, err = json.Marshal(vaultSecret.Data["data"])
			if err != nil {
				return false, err
			}
		case false:
			return true, nil
		}
		err = json.Unmarshal(vaultSecretDataBytes, &decodedVaultSecret)
		if err != nil {
			return
		}
		switch decodedVaultSecret {
		case nil:
		default:
			m := decodedVaultSecret.(map[string]interface{})
			if !reflect.DeepEqual(newData, m) {
				b = true
			}
		}
	}
	return
}

func VaultKVIdempotentWrite(secret interface{}, kvMount string, kvVersion int, kvPath string, c *vault.Client) (err error) {
	var b strings.Builder
	var inputSecret map[string]interface{}

	b.WriteString(kvMount)
	b.WriteString("/")

	switch kvVersion {
	case 1:
		b.WriteString(kvPath)
	case 2:
		b.WriteString("data/")
		b.WriteString(kvPath)
	default:
		return fmt.Errorf("ERROR: unsupported kv version: %d", kvVersion)
	}
	kvCompletePath := b.String()

	vaultSecret, err := c.Logical().Read(kvCompletePath)
	if err != nil {
		return fmt.Errorf("could not read secret from path %s/%s: %s", kvMount, kvPath, err)
	}

	// TODO: ensure we are always passing a map[string]interface instead of testing
	// Convert secret to a map[string]interface{} value
	switch secret.(type) {
	case string:
		b := []byte(secret.(string))
		var f interface{}
		err = json.Unmarshal(b, &f)
		if err != nil {
			return err
		}
		inputSecret = f.(map[string]interface{})
	case map[string]interface{}:
		inputSecret = secret.(map[string]interface{})
	default:
		return fmt.Errorf("unsupported secret type: %s", reflect.TypeOf(secret))
	}

	// Testing for both new secret and strict equality
	ok, err := VaultSecretDataIsDifferent(inputSecret, vaultSecret, kvVersion)
	if err != nil {
		return
	}

	if !ok {
		log.Printf("Secret did not change, %s/%s", kvMount, kvPath)
		return
	}

	switch kvVersion {
	case 1:
		_, err = c.Logical().Write(kvCompletePath, inputSecret)
	case 2:
		_, err = c.Logical().Write(kvCompletePath, map[string]interface{}{
			"data": inputSecret,
		})
	}

	if err != nil {
		return fmt.Errorf("could not write secret to path %s/%s: %s", kvMount, kvPath, err)
	}
	log.Printf("Successfully updated secret: %s/%s", kvMount, kvPath)
	return
}
