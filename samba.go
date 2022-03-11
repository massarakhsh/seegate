package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/massarakhsh/lik"
)

type sysResource struct {
	Namely   string
	Server   string
	Path     string
	Disk     string
	Comments string
	Access   string
	Roles    int
}

var sysMapResource map[string]*sysResource
var sysListResource []*sysResource

func UpdateSamba() {
	sysLoadResourses()
	if sysUpdateSamba("root", "/etc/samba/public.conf2") {
		//sysExecute("/etc/init.d/smbd restart")
	}
}

func sysLoadResourses() {
	sysMapResource = make(map[string]*sysResource)
	sysListResource = []*sysResource{}
	if list := GetList("Resource"); list != nil {
		for ne := 0; ne < list.Count(); ne++ {
			if elm := list.GetSet(ne); elm != nil {
				name := confNameSymbols(elm.GetString("Namely"))
				resource := &sysResource{Namely: name}
				resource.Server = elm.GetString("Server")
				resource.Path = elm.GetString("Path")
				resource.Disk = elm.GetString("Disk")
				resource.Roles = elm.GetInt("Roles")
				resource.Comments = elm.GetString("Comments")
				resource.Access = elm.GetString("Access")
				sysListResource = append(sysListResource, resource)
				sysMapResource[name] = resource
			}
		}
	}
	sort.SliceStable(sysListResource, func(i, j int) bool {
		return sysListResource[i].Namely < sysListResource[j].Namely
	})
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

func sysUpdateSamba(server string, namefile string) bool {
	code := "# Samba list access\n"
	for _, resource := range sysListResource {
		if server := strings.ToLower(resource.Server); server == "root" || true {
			code += fmt.Sprintf("[%s]\n", resource.Namely)
			if val := resource.Comments; val != "" {
				code += fmt.Sprintf("\tcomment = %s\n", val)
			}
			if val := resource.Path; val != "" {
				code += fmt.Sprintf("\tpath = %s\n", val)
			}
			if val := resource.Roles; val >= 0 {
				if (val & 0x1) != 0 {
					code += "\tbrowseable = yes\n"
				} else {
					code += "\tbrowseable = no\n"
				}
				if (val & 0x2) != 0 {
					code += "\twriteable = yes\n"
				} else {
					code += "\twriteable = no\n"
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

func confBuildListIP() []*ElmIP {
	var ips []*ElmIP
	for _, elm := range IPMapSys {
		ips = append(ips, elm)
	}
	sort.SliceStable(ips, func(i, j int) bool {
		return ips[i].IP < ips[j].IP
	})
	return ips
}

func confWrite(namefile string, code string) bool {
	if HostName != "root" {
		if match := lik.RegExParse(namefile, "public"); match != nil {
			namefile = "/mnt/filedata/etc/samba/public.conf2"
		} else if match := lik.RegExParse(namefile, "([^/]+)$"); match != nil {
			namefile = "./root/" + match[1]
		}
	}
	fmt.Printf("configurate %s\n", namefile)
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
		fmt.Printf("updated\n")
		return true
	}
	return false
}
