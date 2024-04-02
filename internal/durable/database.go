package durable

import (
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDB(c string) error {
	var err error

	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  c,
		PreferSimpleProtocol: true,
		//DriverName:           "cloudsqlpostgres",
	}), &gorm.Config{SkipDefaultTransaction: false})

	if err != nil {
		return err
	}
	return nil
}

func Connection() *gorm.DB {
	return db
}
