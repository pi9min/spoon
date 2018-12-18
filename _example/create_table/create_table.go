package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/pi9min/spoon"
	ex "github.com/pi9min/spoon/_example"
)

func main() {
	var (
		outFilePath string
	)
	flag.StringVar(&outFilePath, "o", "./_example/sql/create_table.sql", "set ddl output file path")
	flag.StringVar(&outFilePath, "outfile", "./_example/sql/create_table.sql", "set ddl output file path")
	flag.Parse()

	if outFilePath == "" {
		log.Println("Please set outFilePath. -o or -outfile")
		return
	}

	cli, err := spoon.New()
	if err != nil {
		log.Println(err.Error())
		return
	}

	ebs := []spoon.EntityBehavior{
		&ex.User{},
		ex.Entry{},
		ex.PlayerComment{},
		ex.Bookmark{},
		ex.Balance{},
		ex.NestParent{
			NestChild1: &ex.NestChild1{},
			NestChild2: &ex.NestChild2{},
		},
	}

	schemas, err := cli.GenerateCreateTables(ebs)
	if err != nil {
		log.Println(err.Error())
		return
	}

	f, err := os.Create(outFilePath)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer f.Close()

	if _, err := f.WriteString(strings.Join(schemas, "\n\n")); err != nil {
		log.Println(err.Error())
		return
	}
}
