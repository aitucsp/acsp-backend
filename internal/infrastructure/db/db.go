package db

import (
	"database/sql"
	"fmt"
	"net/url"

	"acsp/infrastructure/config"
)

// InitDB initializes DB
func InitDB(driver string, d *config.DatabaseConfig) (*sql.DB, error) {

	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s%s",
		d.Username,
		d.UserPwd,
		d.Host,
		d.Port,
		d.Name,
		fmt.Sprintf("?loc=UTC&parseTime=true&time_zone=%v", url.QueryEscape("'+00:00'")),
	)

	dbConn, err := sql.Open(driver, connectionString)
	if err != nil {
		return nil, err
	}

	err = dbConn.Ping()
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}
