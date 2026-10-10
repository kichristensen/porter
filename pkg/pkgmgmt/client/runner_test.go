package client

import (
	"context"
	"os"
	"testing"

	"get.porter.sh/porter/pkg/pkgmgmt"
	"get.porter.sh/porter/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunner_Validate(t *testing.T) {
	r := NewTestRunner(t, "lucky-charms", "cereals", true)

	err := r.Validate()
	require.NoError(t, err)
}

func TestRunner_Validate_MissingName(t *testing.T) {
	// Setup failure: empty package name
	r := NewTestRunner(t, "", "candy", true)

	err := r.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "package name to execute not specified")
}

func TestRunner_Validate_MissingExecutable(t *testing.T) {
	r := NewTestRunner(t, "mypackage", "packages", true)

	// Setup failure: Don't copy the package binary into the test context
	err := r.FileSystem.Remove(r.getExecutablePath())
	require.NoError(t, err)

	err = r.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "package not found")
}

func TestRunner_Run_CensorsError(t *testing.T) {
	r := NewTestRunner(t, "mypackage", "mixins", true)
	r.TestContext.Setenv(test.ExpectedCommandExitCodeEnv, "1")
	r.TestContext.Setenv(test.ExpectedCommandErrorEnv, "couldn't run command fail open_door topsecret")
	r.SetSensitiveValues([]string{"topsecret"})

	err := r.Run(context.Background(), pkgmgmt.CommandOptions{Command: "install", Runtime: true})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "couldn't run command fail open_door *******")
	assert.NotContains(t, err.Error(), "topsecret", "expected the sensitive value to be masked in the error")
	assert.NotContains(t, r.TestContext.GetError(), "topsecret", "expected the sensitive value to be masked in the output of the command")
}

func TestRunner_Run_ErrorIncludesRunnableCommand(t *testing.T) {
	r := NewTestRunner(t, "mypackage", "mixins", true)
	r.TestContext.Setenv(test.ExpectedCommandExitCodeEnv, "1")

	err := r.Run(context.Background(), pkgmgmt.CommandOptions{Command: "install", Runtime: true})
	require.Error(t, err)
	// The test command is executed with the test binary
	assert.Contains(t, err.Error(), "package command failed "+os.Args[0]+" ", "expected the error to include a runnable command")
	assert.NotContains(t, err.Error(), r.Getwd()+os.Args[0], "expected the command to not be prefixed with the working directory")
}

func TestRunner_Run_ErrorQuotesCommandArguments(t *testing.T) {
	r := NewTestRunner(t, "mypackage", "mixins", true)
	r.TestContext.Setenv(test.ExpectedCommandExitCodeEnv, "1")

	err := r.Run(context.Background(), pkgmgmt.CommandOptions{Command: "install", File: "my file.yaml", Runtime: true})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "install -f 'my file.yaml'", "expected arguments with spaces to be quoted so that the command is runnable")
}
