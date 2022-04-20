package main

import (
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
)

var DB likbase.DBaser

func OpenDB(serv string, name string, user string, pass string) bool {
	likbase.FId = "SysNum"
	logon := optUser + ":" + optPass
	addr := "tcp(" + optServ + ":3306)"
	if DB = likbase.OpenDBase("mysql", logon, addr, optBase); DB == nil {
		lik.SayError("DB not opened")
		return false
	}
	return true
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		DB = nil
	}
}

func GetElm(part string, id lik.IDB) lik.Seter {
	return DB.GetOneById(part, id)
}

func InsertElm(part string, sets lik.Seter) lik.IDB {
	return DB.InsertElm(part, sets)
}

func UpdateElm(part string, id lik.IDB, sets lik.Seter) bool {
	return DB.UpdateElm(part, id, sets)
}

func DeleteElm(part string, id lik.IDB) bool {
	return DB.DeleteElm(part, id)
}

func GetList(part string) lik.Lister {
	return DB.GetListElm("*", part, "", "SysNum")
}

func CalculateString(sql string) string {
	val := ""
	if one := DB.GetOneBySql(sql); one != nil {
		for _, set := range one.Values() {
			if set.Val != nil {
				val = set.Val.ToString()
				break
			}
		}
	}
	return val
}
