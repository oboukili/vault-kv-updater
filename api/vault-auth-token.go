package api

import (
	"fmt"
	"os"
)

func AuthToken() (token string, err error) {
	vaultAddr = os.Getenv(EnvVaultAddr)
	if vaultAddr == "" {
		vaultAddr = vaultDefaultAddr
	}

	token, ok := os.LookupEnv(EnvVaultToken)
	if !ok {
		err = fmt.Errorf("missing %s environment variable", EnvVaultToken)
	}
	return
}
