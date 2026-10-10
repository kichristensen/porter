package client

import (
	"context"
	"testing"

	"get.porter.sh/porter/pkg/config"
	"get.porter.sh/porter/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileSystem_List(t *testing.T) {
	c := config.NewTestConfig(t)

	p := NewFileSystem(c.Config, "mixins")
	mixins, err := p.List()

	require.Nil(t, err)
	require.Len(t, mixins, 2)
	assert.Equal(t, mixins[0], "exec")
	assert.Equal(t, mixins[1], "testmixin")
}

func TestFileSystem_GetMetadata_CommandFails(t *testing.T) {
	c := config.NewTestConfig(t)
	c.TestContext.Setenv(test.ExpectedCommandExitCodeEnv, "1")
	c.TestContext.Setenv(test.ExpectedCommandErrorEnv, "unknown command version")

	p := NewFileSystem(c.Config, "mixins")
	_, err := p.GetMetadata(context.Background(), "exec")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown command version", "expected the error to include the output of the failed command")
	assert.NotContains(t, c.TestContext.GetError(), "unknown command version", "expected the output of the failed command to not be printed when verbosity is not debug")
}
