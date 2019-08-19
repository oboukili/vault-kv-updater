package main

import (
	"encoding/json"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"log"
	"os"
	"reflect"
)

const vaultDefaultAddr = "http://127.0.0.1:8200"

func vaultClientInit() (c *vault.Client) {
	var token string
	var vaultAddress string

	// initialize a new Vault client
	vaultAddress, ok := os.LookupEnv("VAULT_ADDR")
	if !ok {
		vaultAddress = vaultDefaultAddr
	}

	c, err := vault.NewClient(&vault.Config{
		Address: vaultAddress,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// TODO: implement a proper authentication method cli choice
	// Use Kubernetes Vault authentication
	token, _, err = authKubernetes()
	if err != nil {
		log.Fatalln(err)
	}
	c.SetToken(token)

	return
}

func vaultKVIdempotentWrite(secret interface{}, path string, c *vault.Client) (err error) {
	var inputSecret map[string]interface{}
	vaultSecret, err := c.Logical().Read(path)
	if err != nil {
		return
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
	if vaultSecret != nil && reflect.DeepEqual(inputSecret, vaultSecret.Data) {
		log.Printf("Secret did not change, %s", path)
		return

	}

	_, err = c.Logical().Write(path, inputSecret)
	if err != nil {
		return
	}
	log.Printf("Successfully updated secret: %s", path)
	return
}
