package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func GenerateMigration() (err error) {

	filename := fmt.Sprintf("./migrations/%d.migration.go", time.Now().Unix())

	if err = ioutil.WriteFile(filename, make([]byte, 0), os.FileMode(0744)); err != nil {
		return
	}

	fmt.Printf("created %s\n", filename)

	return
}
