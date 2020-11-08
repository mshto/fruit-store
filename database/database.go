package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // postgres driver
)

//Database struct stores system state database configuration
type Database struct {
	User     string `validate:"required"`
	Password string `validate:"required"`
	Host     string `validate:"required"`
	Port     int    `validate:"required"`
	DBName   string `validate:"required"`
}

const (
	pssqlType      = "postgres"
	pssqlDsnFormat = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"

	pingTimeoutSec = 5
)

// New New
func New(db Database) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(pssqlDsnFormat, db.Host, db.Port, db.User, db.Password, db.DBName)
	conn, err := sql.Open(pssqlType, psqlInfo)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*pingTimeoutSec)
	defer cancel()

	return conn, conn.PingContext(ctx)
}
