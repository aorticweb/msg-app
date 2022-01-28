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

// func TestDBTrans(t *testing.T) {
// 	db := testDB(t)
// 	defer clean(t, db, nil)
// 	exist, err := crud.UserExist(db, "brisco")
// 	require.NoError(t, err)
// 	require.False(t, exist)
// 	require.NoError(t, db.Error)
// 	db = db.Commit()
// 	require.NoError(t, db.Error)
// 	// db = db.Rollback()
// 	db = db.Begin()
// 	require.NoError(t, db.Error)
// 	// err = crud.CreateUser(db, crud.User{Username: "brisco"})
// 	// require.NoError(t, err)
// }
