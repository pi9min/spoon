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
	flag.StringVar(&outFilePath, "o", "./_example/sql/create_index.sql", "set ddl output file path")
	flag.StringVar(&outFilePath, "outfile", "./_example/sql/create_index.sql", "set ddl output file path")
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

	indexes := make([]string, 0, len(ebs))
	for i := range ebs {
		schemas, err := cli.GenerateCreateIndexes(ebs[i])
		if err != nil {
			log.Println(err.Error())
			return
		}
		indexes = append(indexes, schemas...)
	}

	f, err := os.Create(outFilePath)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer f.Close()

	if _, err := f.WriteString(strings.Join(indexes, "\n\n")); err != nil {
		log.Println(err.Error())
		return
	}
}
