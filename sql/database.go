package sql

import (
	"strconv"
	"strings"
	"os/exec"
	"fmt"
)

type Database struct {
	Schema
	host string
	port int
	user string
	pass string
	db   string
}

func NewDatabaseSchema(resource string) (*Database, error){
	argv := strings.Fields(resource)
	var database *Database = &Database{host: "localhost", port: 3306, user: "", pass: "", db: "test"}
	for _, arg := range argv {
		if strings.Index(arg, "-u") == 0 {
			database.user = arg[2:]
		} else if strings.Index(arg, "-p") == 0 {
			database.pass = arg[2:]
		} else if strings.Index(arg, "-P") == 0 {
			port, _ := strconv.Atoi(arg[2:])
			database.port = port
		} else if strings.Index(arg, "-h") == 0 {
			database.host = arg[2:]
		} else if strings.Index(arg, "-") == -1 {
			database.db = arg
		}
	}
	return database, nil
}

func (database *Database) GetSchema() error {
	out, err := exec.Command("mysqldump",
		"-h"+database.host,
		"-P"+strconv.Itoa(database.port),
		"-u"+database.user,
		"-p"+database.pass,
		"--no-data=true",
		database.db).Output()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = database.ParseSQL(out)
	return nil
}
