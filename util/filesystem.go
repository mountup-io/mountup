package util

import (
	"github.com/mountup-io/mountup/api"
	"io/ioutil"
	"os"
	"path/filepath"
)

func SavePrivateKeyToFS(key *api.PrivateKey) error {
	homeDir, _ := os.UserHomeDir()
	pkeyDir := filepath.Join(homeDir, ".mountup/keys", key.KeyName+".pem")

	keyBytes := []byte(key.KeyMaterial)
	err := ioutil.WriteFile(pkeyDir, keyBytes, 0644)
	if err != nil {
		return err
	}

	// Set read-only
	err = os.Chmod(pkeyDir, 0444)
	if err != nil {
		return err
	}

	return nil
}
