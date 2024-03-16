package main

import (
	"fmt"
	"lang/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("[%s] START REPL SESSION\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)

}
