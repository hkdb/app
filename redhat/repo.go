package redhat

import (
	"github.com/hkdb/app/db"
	"github.com/hkdb/app/env"
	"github.com/hkdb/app/utils"

	"fmt"
	"os"
	"os/exec"
	"strings"
)

func AddRepo(s, g string) {

	pkg := s
	if sExt := utils.GetFileExtension(s); sExt == "sh" {
		pkg = utils.GetFileName(s)
	}
	exists, err := db.RepoExists("dnf", pkg)
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
		if g == "" {
			utils.PrintErrorMsgExit("Input Error:", "No gpg url was specified...")
		}
		_, sFile := utils.GetNameFromUrl(s)
		gpg := exec.Command(sudo[0], sudo[1], sudo[2], "/usr/bin/rpm --import "+g)
		utils.RunCmd(gpg, "Import PGP Key Error:")
		cmd := exec.Command(sudo[0], sudo[1], sudo[2], "dnf config-manager --add-repo "+s)
		utils.RunCmd(cmd, "Add Repo Error:")
		name = sFile
		db.RecordSetup("dnf", name, s, g)
	case "sh":
		sFull := utils.GetWorkPath() + "/" + s
		fLine, err := utils.ReadFileLine(sFull, 1)
		if err != nil {
			utils.PrintErrorExit("File Error:", err)
		}
		if fLine != "#!/bin/bash" {
			utils.PrintErrorMsgExit("File Syntax Error:", "Did you add #!/bin/bash to the top of the script?")
		}
		os.Chmod(sFull, 0755)
		runScript := exec.Command(sudo[0], sudo[1], sudo[2], sFull)
		utils.RunCmd(runScript, "Script Error:")
		utils.CreateDirIfNotExist(env.DBDir + "/packages/repo/local/dnf")
		utils.Copy(sFull, env.DBDir+"/packages/repo/local/dnf/"+s)
		name = utils.GetFileName(s)
	default:
		utils.PrintErrorMsgExit("Repo Source Error:", "This is not a supported type of repo source...")
	}

	fmt.Println("\n" + name + " has been added...\n")
	// Record added repo
	if err := db.RecordRepo("dnf", name); err != nil {
		utils.PrintErrorExit("Repo Record Error:", err)
	}

}

func RemoveRepo(s string) {

	exists, err := db.RepoExists("dnf", s)
	if err != nil {
		utils.PrintErrorExit("Error:", err)
	}
	if exists == false {
		utils.PrintErrorMsgExit("Error:", "This repo does not exist or was not added by app...")
	}

	// Check if it's url
	sType := utils.GetRepoType(s)

	sources := strings.Split(s, " ")
	for i := 0; i < len(sources); i++ {
		fmt.Println("Removing /etc/yum.repos.d/" + sources[i] + " with sudo:")
		rmRepo := exec.Command("/usr/bin/sudo", "/bin/bash", "-c", "rm /etc/yum.repos.d/"+sources[i])
		utils.RunCmd(rmRepo, "Repo Remove Error:")
	}

	// Record added repo
	if err = db.RemoveRepo("dnf", sType, s); err != nil {
		utils.PrintErrorExit("Repo Remove Error:", err)
	}

	utils.DeleteDirIfEmpty(env.DBDir + "/packages/repo/local/dnf")

	fmt.Println(s + " has been removed...\n")

}
