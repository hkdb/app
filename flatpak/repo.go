package flatpak

import (
	"github.com/hkdb/app/env"
	"github.com/hkdb/app/utils"
	"github.com/hkdb/app/db"
	
	"os/exec"
	"strings"
	"fmt"
)

func AddRepo(s, g string) {

	pkg := s
	if sExt := utils.GetFileExtension(s); sExt == "sh" {
		pkg = utils.GetFileName(s)
	}
	exists, err := db.RepoExists("flatpak", pkg)
	if err != nil {
		utils.PrintErrorExit("Error:", err)
	}
	if exists == true {
		utils.PrintErrorMsgExit("Error:", "This repo has already been added to app before...")
	}

	if ws := utils.HasWhiteSpace(s); ws == true {
		utils.PrintErrorMsgExit("Input Error:", "Can't add more than one repo at a time...")
	}

	// Check if it's PPA
	sType := utils.GetRepoType(s)
	name := ""

	switch sType {
	case "url":
		_, sFile := utils.GetNameFromUrl(s)
		name = utils.StripExtension(sFile, ".flatpakrepo")
		cmd := exec.Command("/usr/bin/flatpak", "remote-add", "--if-not-exists", name, s)
		utils.RunCmd(cmd, "Add Repo Error:")
		db.RecordSetup("flatpak", name, s, g)
	default:
		utils.PrintErrorMsgExit("Repo Source Error:", "This is not a supported type of repo source...")
	}

	fmt.Println("\n" + name + " has been added...\n")
	// Record added repo
	if err := db.RecordRepo("flatpak", name); err != nil {
		utils.PrintErrorExit("Repo Record Error:", err)
	}

}

func RemoveRepo(s string) {

	exists, err := db.RepoExists("flatpak", s)
	if err != nil {
		utils.PrintErrorExit("Error:", err)
	}
	if exists == false {
		utils.PrintErrorMsgExit("Error:", "This repo does not exist or was not added by app...")
	}

	sources := strings.Split(s, " ")
	for i := 0; i < len(sources); i++ {
		cmd := exec.Command("/usr/bin/flatpak", "remote-delete", s)
		utils.RunCmd(cmd, "Remove Repo Error:")
	}

	// Record added repo
	if err = db.RemoveRepo("flatpak", "url", s); err != nil {
		utils.PrintErrorExit("Repo Remove Error:", err)
	}

	utils.DeleteDirIfEmpty(env.DBDir + "/packages/repo/local/flatpak")

	fmt.Println(s + " has been removed...\n")

}

