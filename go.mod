module github.com/doobcontrol/gDbMysql

go 1.25.6

require (
	github.com/doobcontrol/gDb v1.0.1
	github.com/go-sql-driver/mysql v1.9.3
)

require filippo.io/edwards25519 v1.1.0 // indirect

replace github.com/doobcontrol/gDb => ../gDb //only for local develop
