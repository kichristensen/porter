package mongodb

import (
	"context"
	"testing"

	"get.porter.sh/porter/pkg/portercontext"
	"github.com/stretchr/testify/assert"
)

func TestParseDatabase(t *testing.T) {
	ctx := context.Background()

	tc := portercontext.NewTestContext(t)
	t.Run("db specified", func(t *testing.T) {
		mongo := NewStore(tc.Context, PluginConfig{URL: "mongodb://localhost:27017/test/"})
		if err := mongo.Connect(ctx); err != nil {
			t.Fatal(err)
		}
		defer mongo.Close()
		assert.Equal(t, "test", mongo.database)
	})

	t.Run("default db", func(t *testing.T) {
		mongo := NewStore(tc.Context, PluginConfig{URL: "mongodb://localhost:27017"})
		if err := mongo.Connect(ctx); err != nil {
			t.Fatal(err)
		}
		defer mongo.Close()
		assert.Equal(t, "porter", mongo.database)
	})
}
