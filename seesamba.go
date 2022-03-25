package main

import (
	"fmt"
	"os"

	"github.com/massarakhsh/lik"
)

var optServ = ""
var optBase = ""
var optUser = ""
var optPass = ""

func main() {
	getArgs()

	if !OpenDB(optServ, optBase, optUser, optPass) {
		fmt.Println("Database NOT opened")
		return
	}
	UpdateSamba()
	CloseDB()
}

func getArgs() bool {
	args, ok := lik.GetArgs(os.Args[1:])
	if val := args.GetString("serv"); val != "" {
		optServ = val
	}
	if val := args.GetString("base"); val != "" {
		optBase = val
	}
	if val := args.GetString("user"); val != "" {
		optUser = val
	}
	if val := args.GetString("pass"); val != "" {
		optPass = val
	}
	return ok
}
