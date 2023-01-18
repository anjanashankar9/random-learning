package client_test

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"github.com/anjanashankar9/random-learning/go-sqlite/database/client"
)

func TestNewPostgresClientReturnsClient(t *testing.T) {
	unit, err := client.New("sqlite3", ":memory:")
	require.NoError(t, err)
	require.NotNil(t, unit)
}
