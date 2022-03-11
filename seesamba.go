package main

import (
	"fmt"
	"os"

	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
)

var optServ = ""
var optBase = ""
var optUser = ""
var optPass = ""

var DB likbase.DBaser

func main() {
	getArgs()
	logon := optUser + ":" + optPass
	addr := "tcp(" + optServ + ":3306)"
	if DB = likbase.OpenDBase("mysql", logon, addr, optBase); DB == nil {
		lik.SayError(fmt.Sprint("DB not opened"))
		return
	}

	UpdateSamba()

	DB.Close()
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

func GetElm(part string, id lik.IDB) lik.Seter {
	return DB.GetOneById(part, id)
}

func GetList(part string) lik.Lister {
	return DB.GetListElm("*", part, "", "SysNum")
}
