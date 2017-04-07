package util

import (
	"os"
	"os/user"
)

func GetHomeDir() string {
	// Get the OS specific home directory via the Go standard lib.
	var homeDir string
	usr, err := user.Current()
	if err == nil {
		homeDir = usr.HomeDir
	}

	// Fall back to standard HOME environment variable that works
	// for most POSIX OSes if the directory from the Go standard
	// lib failed.
	if err != nil || homeDir == "" {
		homeDir = os.Getenv("HOME")
	}

	return homeDir + "/"
}

type ApiError struct {
	LogError  error
	UserError error
}

func NewApiError(logError error, userError error) *ApiError {
	apiError := new(ApiError)
	apiError.LogError = logError
	apiError.UserError = userError
	return apiError
}
