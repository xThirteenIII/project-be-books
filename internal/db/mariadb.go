// package db handles connection to MariaDB
package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	defaultMaxLifeTime = 3 * time.Second
	defaultMaxOpenConn = 10
	defaultConnTimeout = 5 * time.Second
)

var connAttempts = 10

func Connect() error {
	// parseTime=true changes the output type of DATE and DATETIME values to time.Time
	// instead of []byte / string
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// Use default localhost
		dsn = "user:password@tcp(localhost:3306)/books?parseTime=true"
	}

	// Open may just validate its arguments without creating a connection
	// to the database. To verify that the data source name is valid, call
	// [DB.Ping].

	// The returned [DB] is safe for concurrent use by multiple goroutines
	// and maintains its own pool of idle connections. Thus, the Open
	// function should be called just once. It is rarely necessary to
	// close a [DB].
	var err error
	var db *sql.DB
	for connAttempts > 0 {
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		// Ping ensures the connection has happened.
		if err := db.Ping(); err == nil {
			break
		}

		log.Printf("Mysql is trying to connect, attempt left: %d\n", connAttempts)
		time.Sleep(defaultConnTimeout)
		connAttempts--
	}

	if err != nil {
		return err
	}

	/* db.SetConnMaxLifetime() is required to ensure connections are closed by the driver safely before
	 * connection is closed by MySQL server, OS, or other middlewares.
	 * Since some middlewares close idle connections by 5 minutes, we recommend timeout shorter than 5 minutes.
	 * This setting helps load balancing and changing system variables too.
	 */
	db.SetConnMaxLifetime(defaultMaxLifeTime)
	db.SetMaxOpenConns(defaultMaxOpenConn) // Should not be impactful for this project.
	db.SetMaxIdleConns(10)

	log.Printf("Connected to MariaDB\n")

	return nil
}
