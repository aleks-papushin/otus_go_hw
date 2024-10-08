package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	dir, _ := os.MkdirTemp("", "testdir")
	defer os.RemoveAll(dir)

	testCases := []struct {
		name          string
		varName       string
		varValue      string
		expectedValue string
	}{
		{
			name:          "Test simple",
			varName:       "VAR1",
			varValue:      "value1",
			expectedValue: "value1",
		},
		{
			name:          "Test spaces",
			varName:       "VAR2",
			varValue:      "value2 with spaces",
			expectedValue: "value2 with spaces",
		},
		{
			name:          "Test tab at the end",
			varName:       "VAR3",
			varValue:      "value3\t",
			expectedValue: "value3",
		},
		{
			name:          "Test terminal 0",
			varName:       "VAR4",
			varValue:      "value4\x00with\x00newlines",
			expectedValue: "value4\nwith\nnewlines",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.WriteFile(filepath.Join(dir, tc.varName), []byte(tc.varValue), 0600) //nolint:gofumpt
			env, _ := ReadDir(dir)
			envValue, ok := env[tc.varName]

			require.True(t, ok, fmt.Sprintf("Variable %s is missing", tc.varName))
			require.Equal(
				t,
				envValue.Value,
				tc.expectedValue,
				fmt.Sprintf("Variable %s has wrong value, got %s, expected %s", tc.varName, envValue.Value, tc.varValue))
		})
	}
}

func TestErrorIfEqualSignInFileName(t *testing.T) {
	dir, _ := os.MkdirTemp("", "testdir")
	defer os.RemoveAll(dir)

	fileName := "invalid=name.txt"
	filePath := filepath.Join(dir, fileName)
	os.WriteFile(filePath, []byte("some content"), 0600) //nolint:gofumpt
	_, err := ReadDir(dir)

	require.Error(t, err, "Expected error is not nil")
	require.Contains(t, err.Error(), "=", "Error message should indicate issue with '=' in file name")
}
