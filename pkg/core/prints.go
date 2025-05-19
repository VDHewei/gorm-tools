package core

import (
	"fmt"
	"github.com/liushuochen/gotable"
	"gorm.io/gorm"
	"log"
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
	for _, t := range values {
		ty, _ := migrator.TableType(t)
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
