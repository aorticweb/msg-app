package crud

import (
	"errors"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func dbConnection(url string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(url), &gorm.Config{
		SkipDefaultTransaction: true,
	})
}

func dbUrl() (string, error) {
	url, exist := os.LookupEnv("POSTGRES_URL")
	if !exist {
		return "", errors.New("Env varialbe POSTGRES_URL is not set")
	}
	return url, nil
}

// WaitForDB ... wait up to timeout for database to be up then return connection
func WaitForDB(timeout time.Duration) (*gorm.DB, error) {
	url, err := dbUrl()
	if err != nil {
		return nil, err
	}
	db, err := dbConnection(url)
	for start := time.Now(); time.Since(start) < timeout; {
		if err == nil {
			break
		}
		db, err = dbConnection(url)
		time.Sleep(1 * time.Second)
	}
	return db, err
}
