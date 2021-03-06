package cliconfig

import (
	"testing"
	"thy/errors"

	"github.com/stretchr/testify/assert"
)

var (
	tn = "mocktenant"
)

func TestGetSecureSetting(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		profile       string
		expectedError error
		expectedVal   string
	}{
		{"missing-key", "", "some-profile", errors.NewS("key cannot be empty"), ""},
		{"empty-value", "hello", "some-profile", nil, ""},
		{"missing-profile", "hello", "", errors.NewS("profile cannot be empty"), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := GetSecureSettingForProfile(tt.key, tt.profile)

			if tt.expectedError != nil {
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			}
			if err == nil {
				assert.Equal(t, tt.expectedVal, val)
			}
		})

	}
}

func TestInstallCommand(t *testing.T) {
	validInstallCmd := IsInstallCmd([]string{"--install"})
	assert.Equal(t, true, validInstallCmd)
	validInstallCmd = IsInstallCmd([]string{"-install"})
	assert.Equal(t, true, validInstallCmd)
}
