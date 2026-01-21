package xyDbMysql

import (
	"database/sql"
	"fmt"
	"github.com/doobcontrol/gDb/xyDb"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

// Define DbDbMysqlAccess
const S_mysqlDriverName = "mysql"
type DbMysqlAccess struct{
   xyDb.DbAccess
}
func (dba *DbMysqlAccess) InitDb(initPars map[string]string, dbStructure xyDb.DbStructure) (string, error) {
   dba.SetDriverName(S_mysqlDriverName)

   adminConnectString, newDbConnectString := makeConnectString(initPars, dbStructure.DbName)
   db, err := sql.Open(dba.DbDriverName, adminConnectString)
   if err != nil {
	   return "", err
   }
   defer db.Close() // Defer closing the database connection until the main function finishes

   //Create tables script
   dScript, err := dba.DbScript(dbStructure)
   if err != nil {
	   return "", err
   }

   //Create Db script
   var dbBuilder strings.Builder
   dbBuilder.WriteString(fmt.Sprintf("CREATE DATABASE %s;", dbStructure.DbName))
   dbBuilder.WriteString(fmt.Sprintf(
	"CREATE USER '%s'@'localhost' IDENTIFIED BY '%s';", 
	initPars[S_newName],
	initPars[S_newPass],
	))
   dbBuilder.WriteString(fmt.Sprintf(
	"GRANT ALL PRIVILEGES ON %s.* TO '%s'@'localhost';", 
	dbStructure.DbName,
	initPars[S_newName],
	))
	dbBuilder.WriteString(dScript)

	if _, err = db.Exec(dbBuilder.String()); err != nil {
		return "", err
	}

   db, err = sql.Open(dba.DbDriverName, newDbConnectString)
   if err != nil {
	   return "", err
   }

   dba.Db = db
   return newDbConnectString, nil
}

const S_adminName = "adminName"
const S_adminPass = "adminPass"
const S_newName = "newName" //the new user for new dbname
const S_newPass = "newPass" //the new user's password for new dbname
const S_protocol = "protocol"
const S_address = "address"
//const S_dbname = "dbname"
const S_params = "params"
func makeConnectString(initPars map[string]string, dbName string) (string, string) {
	var username = initPars[S_adminName]
	var password = initPars[S_adminPass]
	var protocol = "tcp"
	var address = initPars[S_address]
	var dbname = ""
	var params = ""

	if value, ok := initPars[S_protocol]; ok {
		protocol = value
	}
	if value, ok := initPars[S_params]; ok {
		params = "?" + value
	}

	adminConnectString := fmt.Sprintf(
		"%s:%s@%s(%s)/%s%s", 
		username, 
		password,
		protocol, 
		address,
		dbname, 
		params,
	)

	username = initPars[S_newName]
	password = initPars[S_newPass]
	dbname = dbName
	newDbConnectString := fmt.Sprintf(
		"%s:%s@%s(%s)/%s%s", 
		username, 
		password,
		protocol, 
		address,
		dbname, 
		params,
	)

	return adminConnectString, newDbConnectString
}
func (dba *DbMysqlAccess) SetConnect(connectString string) error {
   dba.SetDriverName(S_mysqlDriverName)
   return dba.DbAccess.SetConnect(connectString)
}
