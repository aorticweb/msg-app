package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	db := testDB(t)
	srv := testServer(t, db)
	defer clean(t, db, srv)

	resp, err := http.Get(url(srv.URL, "/health"))

	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)
	checkContentType(t, resp)
}
