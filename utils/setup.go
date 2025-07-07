package utils

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
	"github.com/sajad-dev/gingo-helpers/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDB() (*gorm.DB, *dockertest.Resource, *dockertest.Pool) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	if err := pool.Client.Ping(); err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.Run("mysql", "latest", []string{
		"MYSQL_ROOT_PASSWORD=secret",
		"MYSQL_DATABASE=testdb",
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	var db *gorm.DB
	if err := pool.Retry(func() error {
		dsn := fmt.Sprintf("root:secret@(127.0.0.1:%s)/testdb?parseTime=true", resource.GetPort("3306/tcp"))
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			time.Sleep(1 * time.Second)
		}
		return err
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	if err := db.AutoMigrate(config.ConfigStUtils.DATABASE...); err != nil {
		log.Fatalf("Migration failed: %s", err)
	}

	return db, resource, pool
}
