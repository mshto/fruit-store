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
	User     string `json:"User"      envconfig:"DB_USER"     validate:"required"`
	Password string `json:"Password"  envconfig:"DB_PASSWORD" validate:"required"`
	Host     string `json:"Host"      envconfig:"DB_HOST"     validate:"required"`
	Port     int    `json:"Port"      envconfig:"DB_PORT"     validate:"required"`
	DBName   string `json:"DBName"    envconfig:"DB_NAME"     validate:"required"`
	DBType   string `json:"DBType"    envconfig:"DB_TYPE"     validate:"required"`
}

const (
	pssqlDsnFormat = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"

	pingTimeoutSec = 5
)

// New init database client
func New(db Database) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(pssqlDsnFormat, db.Host, db.Port, db.User, db.Password, db.DBName)
	conn, err := sql.Open(db.DBType, psqlInfo)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*pingTimeoutSec)
	defer cancel()

	return conn, conn.PingContext(ctx)
}
