package utils

import (
	"github.com/hkdb/app/db"
	"github.com/hkdb/app/env"

	"fmt"
	"os"
	"strings"
)

func ListRepo(pm, r string) {

	repos, err := db.ReadPkgSlice("packages", "repo", pm)
	if err != nil {
		PrintErrorExit("Read Repo List Error:", err)
	}

	if r != "" {
		fmt.Println("Repos added by app that matches \"" + r + "\":")
		if len(repos) == 1 && repos[0] == "" {
			fmt.Println(r + " was not added with app...")
		} else {
			for i := 0; i < len(repos); i++ {
				repo := strings.ToLower(repos[i])
				search := strings.ToLower(r)
				if strings.Contains(repo, search) {
					fmt.Println("\n" + repos[i])
				}
			}
		}
	} else {
		fmt.Println("Repos manually installed by app:")
		if len(repos) == 1 && repos[0] == "" {
			fmt.Println("\nNo repos have been added with app yet...")
		} else {
			for i := 0; i < len(repos); i++ {
				fmt.Println(repos[i])
			}
		}
	}

	fmt.Println("\n")

}

func IsPPA(p string) bool {

	pkg := p[0:4]
	if pkg == "ppa:" {
		return true
	}

	return false

}

func IsScript(p string) bool {

	ext := GetFileExtension(p)
	if ext == "sh" {
		return true
	}

	return false

}

func IsUrl(p string) bool {

	http_last := 7
	https_last := 8
	if len(p) <= 8 {
		http_last = len(p) - 1
		https_last = len(p) - 1
	}
	http := p[0:http_last]
	https := p[0:https_last]

	if http == "http://" || https == "https://" {
		return true
	}

	return false

}

func GetRepoType(s string) string {

	// Check if it's PPA
	isPPA := IsPPA(s)

	if isPPA == true {
		return "ppa"
	}

	// Check if it's Deb Source
	isScript := IsScript(s)
	if isScript == true {
		return "sh"
	}

	// Check if it's URL
	isUrl := IsUrl(s)
	if isUrl == true {
		return "url"
	}

	//PrintErrorMsgExit("Error:", "This is not a supported type of repo source")
	return ""

}

// Returns string of extension

func GetNameFromUrl(g string) (string, string) {

	gArg := IsUrl(g)
	if gArg == false {
		PrintErrorMsgExit("Input Error:", "The gpg key should be in the form of a url...")
	}
	gFile := GetFileFromUrl(g)
	gExt := GetFileExtension(gFile)

	return gExt, gFile

}

func GetRepoName(p, pType, g string) string {

	switch pType {
	case "ppa":
		ws := HasWhiteSpace(p)
		if ws == true {
			PrintErrorMsgExit("Error:", "PPAs cannot contain spaces...")
		}
		return p
	case "sh":
		name := GetFileName(p)
		return name
	case "url":
		ws := HasWhiteSpace(p)
		if ws == true {
			PrintErrorMsgExit("Error:", "URLs should not contain spaces...")
		}
		file := GetFileFromUrl(p)
		return file
	default:
		PrintErrorMsgExit("Error:", "Source format error")
	}

	return ""

}

func GetRepoRestoreType(pm, repo string) string {

	rTag := IsPPA(repo)
	if rTag == true {
		return "ppa"
	}

	rType := ""
	prFile := env.DBDir + "/packages/repo/local/" + pm + "/" + repo + ".json"
	if _, err := os.Stat(prFile); err == nil {
		rType = "json"
	}
	psFile := env.DBDir + "/packages/repo/local/" + pm + "/" + repo + ".sh"
	if _, err := os.Stat(psFile); err == nil {
		if rType == "json" {
			PrintErrorMsgExit("Record Error:", "Something is wrong with your app records. There should not be both a .json and .sh file referencing "+repo+"...")
		}
		rType = "sh"
	}

	return rType

}
