package sql

import (
	"strings"
	"regexp"
	"reflect"
	"os"
	"../table"
)

/******************************
 *
 * Schema interface
 *
 */
type SqlSchema interface {
	GetSchema() error
	ParseSQL(contents []byte) error
	GetTables() map[string]table.Table
}

func DiffSchema (src, dst SqlSchema) string {
	diff := make([]string, 0)
	nilTable := *&table.Table{}
	srcTable := src.GetTables()

	for name, t := range dst.GetTables() {
		// table sql on both
		if !reflect.DeepEqual(srcTable[name], nilTable) {
			diff = append(diff, table.DiffTable(srcTable[name], t))
		}else
		// table sql only destination
		{
			diff = append(diff, t.GetSQL() + ";")
		}
	}
	validDiff := make([]string, 0)
	for _, d := range diff {
		if d != "" {
			validDiff = append(validDiff, d)
		}
	}
	return strings.Join(validDiff, "\n")
}

/******************************
 *
 * Schema
 *
 */
type Schema struct {
	tables map[string]table.Table
}

func (schema *Schema) GetTables () (map[string]table.Table) {
	return schema.tables
}

func (schema *Schema) ParseSQL (contents []byte) error {
	schema.tables = make(map[string]table.Table)

	reCreateTable := regexp.MustCompile("(?msi)(CREATE TABLE `(?P<name>[^`]+)`.*? ENGINE[^;]*)")
	tableSqls := reCreateTable.FindAllStringSubmatch(string(contents), -1)

	for _, tableSql := range tableSqls {
		t := table.NewTable(tableSql[2], tableSql[0])
		schema.tables[tableSql[2]] = t
	}

	return nil
}

func GetSchema(resource string) (SqlSchema, error) {
	_, err := os.Stat(resource)
	if err != nil {
		sc, err := NewDatabaseSchema(resource)
		return sc, err
	} else {
		sc, err := NewFileSchema(resource)
		return sc, err
	}
}
