package debian

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
	exists, err := db.RepoExists("apt", pkg)
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
	case "ppa":
		cmd := exec.Command(sudo[0], sudo[1], sudo[2], "add-apt-repository "+s)
		utils.RunCmd(cmd, "Add Repo Error:")
		name = s
	case "sh":
		sFull := utils.GetWorkPath() + "/" + s
		s1 := s[0:1]
		if s1 == "/" {
			sFull = s
			s = utils.GetFileFromFullPath(s)
		}
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
		utils.CreateDirIfNotExist(env.DBDir + "/packages/repo/local/apt")
		utils.Copy(sFull, env.DBDir+"/packages/repo/local/apt/"+s)
		name = utils.GetFileName(s)

		fmt.Println("\nRepo added... Let's run update...\n")
		Update()

	default:
		utils.PrintErrorMsgExit("Repo Source Error:", "This is not a supported type of repo source...")
	}

	fmt.Println("\n" + name + " has been added...\n")
	// Record added repo
	if err := db.RecordRepo("apt", name); err != nil {
		utils.PrintErrorExit("Repo Record Error:", err)
	}

}

func RemoveRepo(s string) {

	exists, err := db.RepoExists("apt", s)
	if err != nil {
		utils.PrintErrorExit("Error:", err)
	}
	if exists == false {
		utils.PrintErrorMsgExit("Error:", "This repo does not exist or was not added by app...")
	}

	// Check if it's PPA
	sType := utils.GetRepoType(s)

	switch sType {
	case "ppa":
		cmd := exec.Command(sudo[0], sudo[1], sudo[2], "add-apt-repository --remove "+s)
		utils.RunCmd(cmd, "Remove Repo Error:")
	default:
		sources := strings.Split(s, " ")
		for i := 0; i < len(sources); i++ {
			fmt.Println("Removing /etc/apt/sources.list.d/" + sources[i] + ".list with sudo:")
			rmRepo := exec.Command("/usr/bin/sudo", "/bin/bash", "-c", "rm /etc/apt/sources.list.d/"+sources[i]+".list")
			utils.RunCmd(rmRepo, "Repo Remove Error:")
		}
	}

	// Record added repo
	if err = db.RemoveRepo("apt", sType, s); err != nil {
		utils.PrintErrorExit("Repo Remove Error:", err)
	}

	utils.DeleteDirIfEmpty(env.DBDir + "/packages/repo/local/apt")

	fmt.Println(s + " has been removed...\n")

}
