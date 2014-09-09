package sql

import (
	"fmt"
	"os/exec"
	"time"
	"io/ioutil"
)

type File struct {
	Schema
	path string
}

func NewFileSchema(resource string) (*File, error){
	var sc *File = &File{path: resource}

	return sc, nil
}

func (file *File) GetSchema() error {
	tmpDbName := "tmp_" + time.Now().Format("20060102150405")

	// read sql file
	contents, err := ioutil.ReadFile(file.path)
	if err != nil {
		return err
	}

	// drop tempoarary database
	err = exec.Command("mysqladmin", "-uroot", "create", tmpDbName).Run()
	if err != nil {
		return err
	}

	// create tables written in sql file
	err = exec.Command("mysql", "-uroot", tmpDbName, "-e", string(contents)).Run()
	if err != nil {
		return err
	}

	out, err := exec.Command("mysqldump",
		"-uroot",
		"--no-data=true",
		tmpDbName).Output()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// drop tempoarary database
	err = exec.Command("mysql", "-uroot", tmpDbName, "-e", "DROP DATABASE " + tmpDbName).Run()
	if err != nil {
		return err
	}

	err = file.ParseSQL(out)

	return nil
}
