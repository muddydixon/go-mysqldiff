package table

import (
	"strings"
	"regexp"
	"reflect"
)

/******************************
 *
 * Table
 *
 */
type Table struct {
	name        string
	sql         string
	columns     []Column
	primaryKeys map[string]PrimaryKey
	uniqueKeys  map[string]UniqueKey
	keys        map[string]Key
}
func (t *Table) GetSQL() string { return t.sql}

func NewTable(name, sql string) Table {
	t := *&Table{name: name, sql: sql}
	t.columns     = make([]Column, 0)
	t.primaryKeys = make(map[string]PrimaryKey)
	t.uniqueKeys  = make(map[string]UniqueKey)
	t.keys        = make(map[string]Key)

	rePrimaryKey  := regexp.MustCompile("(?i)^ *PRIMARY KEY +\\((.*)\\)")
	reUniqueKey   := regexp.MustCompile("(?i)^ *UNIQUE KEY +`(.*)` +\\((.*)\\)")
	reKey         := regexp.MustCompile("(?i)^ *KEY +`(.*)` +\\((.*)\\)")
	reColumn      := regexp.MustCompile("(?i)^ *`(.*?)` +(.+?)[\\n,]?$")

	for _, line := range strings.Split(sql, "\n") {
		if m, _ := regexp.MatchString("(?i)^CREATE", line); m {
			continue
		} else if m, _ := regexp.MatchString(`^\)`, line); m {
			continue
		}

		if m := rePrimaryKey.FindStringSubmatch(line); len(m) > 1 {
			t.primaryKeys[m[1]] = *&PrimaryKey{column: m[1]}
		} else if m := reUniqueKey.FindStringSubmatch(line); len(m) > 1 {
			t.uniqueKeys[m[1]] = *&UniqueKey{name: m[1], column: m[2], }
		} else if m := reKey.FindStringSubmatch(line); len(m) > 1 {
			t.keys[m[1]] = *&Key{name: m[1], column: m[2]}
		} else if m := reColumn.FindStringSubmatch(line); len(m) > 1 {
			t.columns = append(t.columns, *&Column{column: m[1], definition: m[2]})
		}
	}
	return t
}

func DiffTable(src, dst Table) string {
	diff := make([]string, 0)
	// pkDiff  := DiffPrimaryKey(src.primaryKeys, dst.primaryKeys)
	ukDiff  := DiffUniqueKey(src.uniqueKeys, dst.uniqueKeys)
	kDiff   := DiffKey(src.keys, dst.keys)
	colDiff := DiffColumn(src.columns, dst.columns)

	if len(colDiff) > 0 {
		diff = append(diff,  "ALTER TABLE `" + src.name + "` " + strings.Join(colDiff, ", ") + ";")
	}
	if len(ukDiff) > 0 {
		diff = append(diff,  strings.Join(ukDiff, "\n"))
	}
	if len(kDiff) > 0 {
		diff = append(diff,  strings.Join(kDiff, "\n"))
	}
	if len(diff) > 0 {
		return strings.Join(diff, "\n")
	} else {
		return ""
	}
}

/******************************
 *
 * PrimaryKey
 *
 */
type PrimaryKey struct {
	column string
}
func DiffPrimaryKey(src, dst map[string]PrimaryKey) {
}
/******************************
 *
 * UniqueKey
 *
 */
type UniqueKey struct {
	name string
	column string
}
func DiffUniqueKey(src, dst map[string]UniqueKey) []string{
	diff := make([]string, 0)
	nilKey := *&Key{}

	for name, k := range dst {
		if !reflect.DeepEqual(src[name], nilKey) {
			continue
		}
		diff = append(diff, "ADD UNIQUE INDEX `" + name + "` " + k.column + ";")
	}

	for name, _ := range src {
		if reflect.DeepEqual(dst[name], nilKey) {
			diff = append(diff, "DROP INDEX `" + name + "`;")
		}
	}
	return diff
}


/******************************
 *
 * Key
 *
 */
type Key struct {
	name string
	column string
}
func DiffKey(src, dst map[string]Key) []string {
	diff := make([]string, 0)
	nilKey := *&Key{}
	for name, k := range dst {
		if !reflect.DeepEqual(src[name], nilKey) {
			continue
		}
		diff = append(diff, "ADD INDEX `" + name + "` " + k.column + ";")
	}

	for name, _ := range src {
		if reflect.DeepEqual(dst[name], nilKey) {
			diff = append(diff, "DROP INDEX `" + name + "`;")
		}
	}
	return diff
}
/******************************
 *
 * Colmun
 *
 */
type Column struct {
	column     string
	definition string
}
func DiffColumn(src, dst []Column) []string {
	diff := make([]string, 0)
	srcColMap := make(map[string]Column)
	dstColMap := make(map[string]Column)
	for _, col := range src {
		srcColMap[col.column] = col
	}
	for _, col := range dst {
		dstColMap[col.column] = col
	}
	nilColumn := *&Column{}

	for idx, column := range dst {
		if reflect.DeepEqual(srcColMap[column.column], nilColumn) {
			if idx != 0 {
				diff = append(diff, "ADD `" + column.column + "` " + column.definition +
					" AFTER `" + dst[idx - 1].column + "`")
			} else {
				diff = append(diff, "ADD `" + column.column + "` " + column.definition)
			}
		} else if srcColMap[column.column].definition != column.definition {
			diff = append(diff, "MODIFY `" + column.column + "` "+column.definition)
		}
	}
	for _, column := range src {
		if reflect.DeepEqual(dstColMap[column.column], nilColumn) {
			diff = append(diff, "DROP `" + column.column + "`")
		}
	}
	return diff
}
