package config

import "fmt"

const (
	// connStringFormat is our default DB connection string format to use.
	//
	// Examples showing text and url based formats in case we need them later.
	//
	// - host=host port=port dbname=dbname user=theusername password=pass
	// - postgresql://user:pass@host:port/dbname
	connStringFormat = "%s://%s:%s@%s:%v/%v"
)

// RDSSecret represents a specific format DB connection details required by
// RDS for password rotation.
type RDSSecret struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Engine   string `json:"engine"`
	DBName   string `json:"dbname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// ConnectionString returns a usable connection string from the data.
func (r *RDSSecret) ConnectionString() string {
	return fmt.Sprintf(connStringFormat, r.Engine, r.Username, r.Password, r.Host, r.Port, r.DBName)
}
