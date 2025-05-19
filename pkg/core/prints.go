package core

import (
	"fmt"
	"github.com/liushuochen/gotable"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/migrator"
	"log"
	"reflect"
	"sort"
	"strings"
)

func PrintTables(db *gorm.DB, tables []string, excludes []string) bool {
	var migrator = db.Migrator()
	all, err := migrator.GetTables()
	if err != nil {
		log.Fatalln("query database tables failed, error:", err)
		return true
	}
	var (
		filterMap  = make(map[string]struct{})
		excludeMap = make(map[string]struct{})
	)
	for _, v := range excludes {
		v = strings.TrimSpace(v)
		excludeMap[v] = struct{}{}
	}
	if len(tables) == 0 {
		for _, v := range all {
			v = strings.TrimSpace(v)
			if _, ok := excludeMap[v]; ok {
				continue
			}
			filterMap[v] = struct{}{}
		}
	} else {
		for _, v := range tables {
			v = strings.TrimSpace(v)
			if _, ok := excludeMap[v]; ok {
				continue
			}
			if _, ok := filterMap[v]; ok {
				continue
			}
			filterMap[v] = struct{}{}
		}
	}
	var values []string
	for k, _ := range filterMap {
		values = append(values, k)
	}
	if len(values) < 0 {
		return true
	}
	sort.Strings(values)
	printTable, err := gotable.Create("table_name", "comment")
	if err != nil {
		log.Fatalln("Create table failed: ", err.Error())
		return true
	}
	var m = newExtrasMigrate(db, migrator)
	for _, t := range values {
		ty, _ := m.TableType(t)
		if ty == nil {
			printTable.AddRow([]string{t, ""})
		} else {
			comment, _ := ty.Comment()
			printTable.AddRow([]string{ty.Name(), comment})
		}
	}
	fmt.Println(printTable)
	return true
}

func PrintTable(db *gorm.DB, tableName string) bool {
	var migrator = db.Migrator()
	all, err := migrator.GetTables()
	if err != nil {
		log.Fatalln("query database tables failed, error:", err)
		return true
	}
	var (
		filterMap = make(map[string]struct{})
	)
	for _, v := range all {
		v = strings.TrimSpace(v)
		filterMap[v] = struct{}{}
	}
	if _, ok := filterMap[tableName]; !ok {
		tableName = strings.TrimSpace(tableName)
		if _, ok = filterMap[tableName]; !ok {
			log.Fatalln("table" + tableName + " not found in database")
			return true
		}
	}
	printTable, err := gotable.Create("field", "type", "null", "pk/uk", "default", "comment")
	if err != nil {
		log.Fatalln("Create table failed: ", err.Error())
		return true
	}
	ty, _ := migrator.TableType(tableName)
	types, _ := migrator.ColumnTypes(tableName)
	if types != nil {
		for _, v := range types {
			var keyValue = ""
			uk, _ := v.Unique()
			pk, _ := v.PrimaryKey()
			comment, _ := v.Comment()
			nullable, _ := v.Nullable()
			defaultValue, _ := v.DefaultValue()
			nullableStr := fmt.Sprintf("%v", nullable)
			if pk {
				keyValue = "pk"
			} else if uk {
				keyValue = "uk"
			}
			typeVal, _ := v.ColumnType()
			printTable.AddRow([]string{v.Name(), typeVal, nullableStr, keyValue, defaultValue, comment})
		}
	}
	var tableComment string
	if ty != nil {
		comment, _ := ty.Comment()
		if comment != "" {
			tableComment = comment
		}
	}
	if tableComment == "" {
		fmt.Println("<" + tableName + ">")
	} else {
		fmt.Println("<" + tableName + "> -- " + tableComment)
	}
	fmt.Println(printTable)
	return true
}

type migratorImpl struct {
	db *gorm.DB
	m  gorm.Migrator
}

func newExtrasMigrate(db *gorm.DB, migrator gorm.Migrator) *migratorImpl {
	return &migratorImpl{
		db: db,
		m:  migrator,
	}
}

func (m migratorImpl) RunWithValue(value interface{}, fc func(*gorm.Statement) error) error {
	if m.m != nil {
		if v, ok := m.m.(postgres.Migrator); ok {
			return v.RunWithValue(value, fc)
		}
		v := reflect.ValueOf(m.m)
		if caller := v.MethodByName(`RunWithValue`); caller.IsValid() {
			res := caller.Call([]reflect.Value{reflect.ValueOf(value), reflect.ValueOf(fc)})
			return res[0].Interface().(error)
		}
	}
	stmt := &gorm.Statement{DB: m.db}
	if m.db.Statement != nil {
		stmt.Table = m.db.Statement.Table
		stmt.TableExpr = m.db.Statement.TableExpr
	}

	if table, ok := value.(string); ok {
		stmt.Table = table
	} else if err := stmt.ParseWithSpecialTableName(value, stmt.Table); err != nil {
		return err
	}

	return fc(stmt)
}

func (m migratorImpl) CurrentSchema(stmt *gorm.Statement, table string) (interface{}, interface{}) {
	if m.m != nil {
		if _, ok := m.m.(postgres.Migrator); ok {
			if strings.Contains(table, ".") {
				if tables := strings.Split(table, `.`); len(tables) == 2 {
					return tables[0], tables[1]
				}
			}
			if stmt.TableExpr != nil {
				if tables := strings.Split(stmt.TableExpr.SQL, `"."`); len(tables) == 2 {
					return strings.TrimPrefix(tables[0], `"`), table
				}
			}
			return gorm.Expr("CURRENT_SCHEMA()"), table
		}
		v := reflect.ValueOf(m.m)
		if caller := v.MethodByName(`CurrentSchema`); caller.IsValid() {
			res := caller.Call([]reflect.Value{reflect.ValueOf(stmt), reflect.ValueOf(table)})
			return res[0].Interface(), res[1].Interface()
		}
	}
	if tables := strings.Split(table, `.`); len(tables) == 2 {
		return tables[0], tables[1]
	}
	m.db = m.db.Table(table)
	return m.m.CurrentDatabase(), table
}

// TableType table type return tableType,error
func (m migratorImpl) TableType(value interface{}) (tableType gorm.TableType, err error) {
	var table migrator.TableType
	err = m.RunWithValue(value, func(stmt *gorm.Statement) error {
		var (
			values = []interface{}{
				&table.SchemaValue, &table.NameValue, &table.TypeValue, // &table.CommentValue,
			}
			currentDatabase, tableName = m.CurrentSchema(stmt, stmt.Table)
			tableTypeSQL               = `SELECT schemaname, tablename, 
                           CASE WHEN schemaname LIKE 'pg_%' THEN 'SYSTEM TABLE'
                                ELSE 'BASE TABLE' END
                           FROM pg_catalog.pg_tables 
                           WHERE schemaname = ? AND tablename = ?`
			tableCommentSQL = `SELECT obj_description(pg_class.oid, 'pg_class') AS table_comment 
FROM information_schema.tables as tables
JOIN pg_class ON pg_class.relname = tables.table_name
WHERE table_schema = ? and table_name= ?;
`
		)
		row := m.db.Table(tableName.(string)).Raw(tableTypeSQL, currentDatabase, tableName).Row()

		if scanErr := row.Scan(values...); scanErr != nil {
			return scanErr
		}

		row = m.db.Table(tableName.(string)).Raw(tableCommentSQL, currentDatabase, tableName).Row()
		if scanErr := row.Scan(&table.CommentValue); scanErr != nil {
			return scanErr
		}
		return nil
	})
	return table, err
}
