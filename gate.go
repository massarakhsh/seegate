package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/massarakhsh/lik"
)

type itResource struct {
	SysNum   int
	Namely   string
	Access   string
	Server   string
	Path     string
	Disk     string
	Comments string
	Roles    int64
}

type itAccess struct {
	SysNum      int
	SysResource int
	Namely      string
	IP          string
	SysOperator int
	SysDepart   int
	Comments    string
	Access      string
	Roles       int64
}

const RAT_BRAWSABLE = 0x1
const RAT_WRITEABLE = 0x2
const RAT_PRIVATE = 0x4

var sysMapResource map[int]*itResource
var sysMapAccess map[int]*itAccess

func UpdateGate() {
	for sysLoadResources() {
		if sysUpdateSamba("/etc/samba/public.conf") {
			if strings.ToLower(hostName) == "master" {
				sysExecute("/etc/init.d/smbd restart")
			}
		}
		time.Sleep(time.Second * 60)
	}
}

func sysLoadResources() bool {
	if !sysTableResource() {
		return false
	}
	if !sysTableAccess() {
		return false
	}
	if !sysSynchroTables() {
		return false
	}
	return true
}

func sysTableResource() bool {
	sysMapResource = make(map[int]*itResource)
	if list := GetList("Resource"); list != nil {
		for ne := 0; ne < list.Count(); ne++ {
			if elm := list.GetSet(ne); elm != nil {
				sys := int(elm.GetInt("SysNum"))
				name := confNameSymbols(elm.GetString("Namely"))
				resource := &itResource{Namely: name, SysNum: sys}
				resource.Server = elm.GetString("Server")
				resource.Path = elm.GetString("Path")
				resource.Disk = elm.GetString("Disk")
				resource.Roles = elm.GetInt("Roles")
				resource.Comments = elm.GetString("Comments")
				resource.Access = elm.GetString("Access")
				sysMapResource[sys] = resource
			}
		}
	}
	return true
}

func sysTableAccess() bool {
	sysMapAccess = make(map[int]*itAccess)
	if list := GetList("Access"); list != nil {
		for ne := 0; ne < list.Count(); ne++ {
			if elm := list.GetSet(ne); elm != nil {
				sys := int(elm.GetInt("SysNum"))
				access := &itAccess{SysNum: sys}
				access.SysResource = int(elm.GetInt("SysResource"))
				access.SysOperator = int(elm.GetInt("SysOperator"))
				access.SysDepart = int(elm.GetInt("SysDepart"))
				access.Namely = elm.GetString("Namely")
				access.IP = elm.GetString("IP")
				access.Roles = elm.GetInt("Roles")
				access.Comments = elm.GetString("Comments")
				sysMapAccess[sys] = access
			}
		}
	}
	return true
}

func sysSynchroTables() bool {
	for sys, acc := range sysMapAccess {
		sysRes := acc.SysResource
		if sysMapResource[sysRes] == nil {
			delete(sysMapAccess, sys)
		}
	}
	for _, res := range sysMapResource {
		sysSynchroResource(res)
	}
	return true
}

func sysSynchroResource(resource *itResource) {
	list := []*itAccess{}
	for _, acc := range sysMapAccess {
		if acc.SysResource == resource.SysNum {
			list = append(list, acc)
		}
	}
	if len(list) == 0 && resource.Access != "" {
		words := strings.Split(resource.Access, " ")
		for _, word := range words {
			if addr := strings.Trim(word, "\r\n\t "); addr != "" {
				set := lik.BuildSet("SysResource", resource.SysNum, "Namely", addr)
				if sys := int(InsertElm("Access", set)); sys > 0 {
					access := &itAccess{SysNum: sys}
					access.SysResource = resource.SysNum
					access.Namely = addr
					sysMapAccess[sys] = access
				}
			}
		}
	}
}

func sysExecute(cmd string) {
	lik.SayInfo(cmd)
	cmds := strings.Split(cmd, " ")
	cmdc := cmds[0]
	cmds = cmds[1:]
	if exe := exec.Command(cmdc, cmds...); exe != nil {
		exe.Run()
	}
}

func sysUpdateSamba(namefile string) bool {
	resources := []*itResource{}
	for _, resource := range sysMapResource {
		resources = append(resources, resource)
	}
	sort.SliceStable(resources, func(i, j int) bool {
		return resources[i].Namely < resources[j].Namely
	})
	code := "# Samba list access. DO NOT EDIT this file!\n"
	for _, resource := range resources {
		if server := strings.ToLower(resource.Server); server == "master" {
			code += fmt.Sprintf("[%s]\n", resource.Namely)
			if val := resource.Comments; val != "" {
				code += fmt.Sprintf("\tcomment = %s\n", val)
			}
			if val := resource.Path; val != "" {
				code += fmt.Sprintf("\tpath = %s\n", val)
			}
			if val := resource.Roles; val >= 0 {
				if (val & RAT_BRAWSABLE) != 0 {
					code += "\tbrowseable = yes\n"
				} else {
					code += "\tbrowseable = no\n"
				}
				if (val & RAT_WRITEABLE) != 0 {
					code += "\twriteable = yes\n"
				} else {
					code += "\twriteable = no\n"
				}
				if (val & RAT_PRIVATE) != 0 {
					code += "\tforce user = root\n"
					code += "\tforce group = root\n"
				}
			}
			if val := resource.Access; val != "" {
				val = regexp.MustCompile("[^0-9\\./\\ ]").ReplaceAllString(val, "")
				code += fmt.Sprintf("\thosts allow = %s\n", val)
			}
			code += "\n"
		}
	}
	return confWrite(namefile, code)
}

func confNameSymbols(name string) string {
	//name = strings.ToLower(name)
	name = lik.Transliterate(name)
	name = regexp.MustCompile("[^0-9a-zA-Z\\-\\_]").ReplaceAllString(name, "-")
	return name
}

func confWrite(namefile string, code string) bool {
	if file, err := os.Open(namefile); err == nil {
		oldcode := ""
		buf := make([]byte, 4096)
		for {
			if n, err := file.Read(buf); err == nil {
				oldcode += string(buf[:n])
			} else {
				break
			}
		}
		file.Close()
		if oldcode == code {
			fmt.Printf("equal\n")
			return false
		}
	}
	if file, err := os.Create(namefile); err == nil {
		file.WriteString(code)
		file.Close()
		fmt.Printf("Configuration file %s was updates\n", namefile)
		return true
	}
	return false
}