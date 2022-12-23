package gcp_auth

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var defaultCredentialsPath string

// Authenticate checks if user has credentials to use Google SDK and if not authenticates through a gcloud shell out
func Authenticate() error {
	if !hasServiceAccount() {
		err := createDefaultCredentials()
		if err != nil {
			return err
		}
	}
	return nil
}

func IsAuthenticated() (bool, error) {
	credentials, errCred := hasDefaultCredentials()
	if errCred != nil {
		return false, fmt.Errorf("hasDefaultCredentials check failed: %w", errCred)
	}
	return credentials, nil
}

// hasDefaultCredentials checks the default credentials used by google SDK, which by default are stored in ~/.config/gcloud/application_default_credentials.json
func hasDefaultCredentials() (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, fmt.Errorf("retrieving homedir failed: %w", err)
	}
	defaultCredentialsPath = filepath.Join(homeDir, ".config", "gcloud", "application_default_credentials.json")
	if _, err := os.Stat(defaultCredentialsPath); os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}

// hasServiceAccount checks if service account credentials are used (use case for ci usage)
func hasServiceAccount() bool {
	if _, exists := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS"); exists {
		return true
	}
	return false
}

// createDefaultCredentials redirects to browser to retrieve the default application credentials
func createDefaultCredentials() error {
	command := "gcloud auth application-default login"
	cmd := exec.Command("bash", "-c", command)

	// Windows usually does not ship bash, and execution in a subshell is not needed
	if runtime.GOOS == "windows" {
		commandParts := strings.Split(command, " ")
		cmd = exec.Command(commandParts[0], commandParts[1:]...)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	stdErr := stderr.String()
	if err != nil {
		return fmt.Errorf("'gcloud auth application-default login' failed: '%s'. "+
			"Please make sure gcloud is installed and if you are not having a bash shell, excecute 'gcloud auth application-default login' in another shell", stdErr)
	}
	return nil
}
