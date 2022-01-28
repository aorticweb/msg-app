package tests

import (
	"testing"
	"time"

	"github.com/aorticweb/msg-app/app/crud"
	"github.com/stretchr/testify/require"
)

// TestDBConnection ... Live db connection to assert content of database
func TestDBConnection(t *testing.T) {
	_, err := crud.WaitForDB(3 * time.Second)
	require.NoError(t, err)
}
