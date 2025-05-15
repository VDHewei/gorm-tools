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
		tables = all
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
		ty, err := migrator.TableType(t)
		if err != nil {
			log.Fatalln("query table("+t+")info failed: ", err.Error())
			return true
		}
		comment, _ := ty.Comment()
		printTable.AddRow([]string{ty.Name(), comment})
	}
	fmt.Println(printTable)
	return true
}
