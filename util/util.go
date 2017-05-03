package util

import (
	"os"
	"os/user"
	"regexp"
	"strings"
	"unicode"
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

func (a *ApiError) Error() string {
	return a.LogError.Error()
}

func NewAPIError(logError error, userError error) *ApiError {
	apiError := new(ApiError)
	apiError.LogError = logError
	apiError.UserError = userError
	return apiError
}

func NewAPIErrorFromOne(err error) *ApiError {
	apiError := new(ApiError)
	apiError.LogError = err
	apiError.UserError = err
	return apiError
}

func DashDelimiterToCamelCase(input string) string {
	input = strings.TrimSpace(input)
	if len(input) > 0 {
		r, _ := regexp.Compile("-(.)")
		out := []rune(input)
		out[0] = unicode.ToUpper(out[0])
		return r.ReplaceAllStringFunc(string(out), func(w string) string {
			return string(strings.ToUpper(w)[1])
		})
	}
	return input
}
