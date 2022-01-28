package tests

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/aorticweb/msg-app/app/crud"
	api "github.com/aorticweb/msg-app/app/handlers"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func url(base string, route string) string {
	return fmt.Sprintf("%s%s", base, route)
}

// testDB ... fixture for live db connection ... do not forget to rollback transaction
func testDB(t *testing.T) *gorm.DB {
	db, err := crud.WaitForDB(time.Second)
	require.NoError(t, err)
	db.DisableNestedTransaction = false
	log.Println(db.DisableNestedTransaction)
	tx := db.Begin()
	tx.SavePoint("beforeTest")
	return tx
}

// testServer ... fixture for test server ... do not forget to close server
func testServer(t *testing.T, db *gorm.DB) *httptest.Server {
	logger := log.New(os.Stdout, "msg-app: ", log.LstdFlags|log.Llongfile)
	srv := httptest.NewServer(api.NewAPI(db, logger))
	return srv
}

func clean(t *testing.T, db *gorm.DB, srv *httptest.Server) {
	if db != nil {
		db.RollbackTo("beforeTest")
	}
	if srv != nil {
		srv.Close()
	}
}

func checkContentType(t *testing.T, resp *http.Response) {
	require.Equal(t, resp.Header.Get("Content-Type"), "application/json")
}
