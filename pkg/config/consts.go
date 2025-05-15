package config

const (
	// DbMySQL Gorm Drivers mysql || postgres || sqlite || sqlserver
	DbMySQL          DBType = "mysql"
	DbPostgres       DBType = "postgres"
	DbSQLite         DBType = "sqlite"
	DbSQLServer      DBType = "sqlserver"
	DbClickHouse     DBType = "clickhouse"
	DefaultQueryPath        = "./dao/query"
	// default dsn
	defaultMysqlDSN      = "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	defaultPostgresDSN   = "host=localhost port=5432 user=postgres password=<PASSWORD> dbname=test sslmode=disable"
	defaultSQLiteDSN     = "file:./test.db?cache=shared&mode=rwc"
	defaultSQLServerDSN  = "sqlserver://sa:123456@localhost:1433?database=test"
	defaultClickHouseDSN = "tcp://127.0.0.1:9000?username=&database=&read_timeout=10&write_timeout=20&alt_hosts=127.0.0.2:9000,127.0.0.3:9000"
	version              = `v1.0.9`
)
