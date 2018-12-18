package main

import (
	"fmt"
	"os"

	"github.com/pi9min/spoon"
	"github.com/pi9min/spoon/_example/readme/user"
)

func main() {
	cli, err := spoon.New()
	if err != nil {
		panic(err)
	}

	schema, err := cli.GenerateCreateTable(&user.User{})
	if err != nil {
		panic(err)
	}

	// output to stdout
	fmt.Fprint(os.Stdout, schema)
}
