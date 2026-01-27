package xyDbMysql

import (
	"fmt"
	"github.com/doobcontrol/gDb/xyDb"
	"testing"
)

func TestDbMysqlAccess_SetConnect(t *testing.T) {
	mysqlAccess := &DbMysqlAccess{}
    err := mysqlAccess.SetConnect("")
	var expectErr = "Error 1045 (28000): Access denied for user ''@'localhost' (using password: NO)"
    if err == nil {
        t.Error(fmt.Printf("DbAccess.SetConnect err expected an error, but got nil"))
    } else {
		if err.Error() != expectErr {
        	t.Error(fmt.Printf("DbAccess.SetConnect err expected an error: %s, but got: %s", 
			expectErr, err.Error()))
    	} 
	}

	// Ensure the connection string is valid and can connect to the test database
    err = mysqlAccess.SetConnect(testConnect_admin)
    if err != nil {
        t.Error(fmt.Printf("DbAccess.SetConnect err expected nil, but got an error: %s", err.Error()))
    }
}
func TestDbMysqlAccess_makeConnectString(t *testing.T) {
	pars := &map[string]string{
		"adminName": "root",
		"adminPass" : "D00.control",
		"newName": "testDbadmin",
		"newPass" : "ynyhcWh@7admin",
		"address" : "127.0.0.1:3306",
		"params" : "multiStatements=true",
	}
    adminConnectString, newDbConnectString := makeConnectString(pars, "testDb")
    if adminConnectString != testConnect_admin {
        t.Error(fmt.Printf("DbAccess.makeConnectString adminConnectString expected: %s, but got: %s", 
		adminConnectString, testConnect_admin))
    }
    if newDbConnectString != testConnect_user {
        t.Error(fmt.Printf("DbAccess.makeConnectString newDbConnectString expected: %s, but got: %s", 
		testConnect_user, newDbConnectString))
    }
}
func TestDbMysqlAccess_InitDb(t *testing.T) {
	pars := &map[string]string{
		"adminName": "root",
		"adminPass" : "D00.control",
		"newName": "testDbadmin",
		"newPass" : "ynyhcWh@7admin",
		"address" : "127.0.0.1:3306",
		"params" : "multiStatements=true",
	}
	db := xyDb.DbStructure{
		DbName: "testDb",
		Tables: []xyDb.DbTable{
			{
				TableName: "table1",
				Fields: []xyDb.DbField{
					{
						FieldName: "F1",
						DataType: "string",
						Length: 10,
						IsKey: true,
					},
					{
						FieldName: "F2",
						DataType: "string",
						Length: 10,
						IsKey: false,
					},
				},
			},
			{
				TableName: "table2",
				Fields: []xyDb.DbField{
					{
						FieldName: "F1",
						DataType: "string",
						Length: 10,
						IsKey: true,
					},
					{
						FieldName: "F2",
						DataType: "string",
						Length: 10,
						IsKey: false,
					},
				},
			},
		},
	}

	mysqlAccess := &DbMysqlAccess{}

	//clean db
	cleanDb(mysqlAccess)

    newConnectString, err := mysqlAccess.InitDb(pars, db)
    if err != nil {
        t.Error(fmt.Printf("DbAccess.InitDb err expected nil, but got an error: %s", err.Error()))
    }
    if newConnectString != testConnect_user {
        t.Error(fmt.Printf("DbAccess.InitDb newConnectString expected: %s, but got: %s", 
		testConnect_user, newConnectString))
    }

	//clean db
	cleanDb(mysqlAccess)
}
func TestDbMysqlAccess_sql(t *testing.T) {
	mysqlAccess := &DbMysqlAccess{}

	//clean db
	cleanDb(mysqlAccess)

    //init db
	initDb(mysqlAccess)
	mysqlAccess.Close()

	//connect
	mysqlAccess.SetConnect(testConnect_user)

	err := mysqlAccess.ExSql(fmt.Sprintf("insert into %s(%s) values(%s)","table1","F1,F2","'abc','efg'"))
    if err != nil {
        t.Error(fmt.Printf("DbAccess.ExSql.insert err expected nil, but got an error: %s", err.Error()))
    } else {
		t.Log("DbAccess.ExSql.insert succeed")
	}
	err = mysqlAccess.ExSql(fmt.Sprintf("insert into %s(%s) values(%s)","table1","F1,F2","'abc1','efg2'"))
    if err != nil {
        t.Error(fmt.Printf("DbAccess.ExSql.insert err expected nil, but got an error: %s", err.Error()))
	} else {
		t.Log("DbAccess.ExSql.insert succeed")
	}

	err = mysqlAccess.ExSql(fmt.Sprintf("update %s set %s where %s='%s'","table1","F2='zzz'","F1", "abc"))
    if err != nil {
        t.Error(fmt.Printf("DbAccess.ExSql.update err expected nil, but got an error: %s", err.Error()))
    } else {
		t.Log("DbAccess.ExSql.update succeed")
	}

	record, err := mysqlAccess.Query(fmt.Sprintf("select * from %s","table1"))
    if err != nil {
        t.Error(fmt.Printf("DbAccess.ExSql.select err expected nil, but got an error: %s", err.Error()))
    } else {
		if len(*record) != 2 {
			t.Error(fmt.Printf("DbAccess.ExSql.select expected: %d records, but got: %d records", 2, len(*record)))
		} else {
			t.Log("DbAccess.ExSql.select succeed")
		}
	}

	err = mysqlAccess.ExSql(fmt.Sprintf("delete from %s where %s='%s'","table1","F1", "abc"))
    if err != nil {
        t.Error(fmt.Printf("DbAccess.ExSql.delete err expected nil, but got an error: %s", err.Error()))
    } else {
		t.Log("DbAccess.ExSql.delete succeed")
	}

	//clean db
	cleanDb(mysqlAccess)
}

var testConnect_admin = "root:D00.control@tcp(127.0.0.1:3306)/mysql?multiStatements=true"
var testConnect_user = "testDbadmin:ynyhcWh@7admin@tcp(127.0.0.1:3306)/testDb?multiStatements=true"
func cleanDb(mysqlAccess *DbMysqlAccess){
	mysqlAccess.SetConnect(testConnect_admin)
	mysqlAccess.ExSql("DROP DATABASE testDb;")
	mysqlAccess.ExSql("DROP USER 'testDbadmin'@'localhost';")
	mysqlAccess.Close()
}
func initDb(mysqlAccess *DbMysqlAccess){
	pars := &map[string]string{
		"adminName": "root",
		"adminPass" : "D00.control",
		"newName": "testDbadmin",
		"newPass" : "ynyhcWh@7admin",
		"address" : "127.0.0.1:3306",
		"params" : "multiStatements=true",
	}
	db := xyDb.DbStructure{
		DbName: "testDb",
		Tables: []xyDb.DbTable{
			{
				TableName: "table1",
				Fields: []xyDb.DbField{
					{
						FieldName: "F1",
						DataType: "string",
						Length: 10,
						IsKey: true,
					},
					{
						FieldName: "F2",
						DataType: "string",
						Length: 10,
						IsKey: false,
					},
				},
			},
			{
				TableName: "table2",
				Fields: []xyDb.DbField{
					{
						FieldName: "F1",
						DataType: "string",
						Length: 10,
						IsKey: true,
					},
					{
						FieldName: "F2",
						DataType: "string",
						Length: 10,
						IsKey: false,
					},
				},
			},
		},
	}
	mysqlAccess.InitDb(pars, db)
}