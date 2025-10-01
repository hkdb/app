package arch

import (
	"github.com/hkdb/app/db"
	"github.com/hkdb/app/env"
	"github.com/hkdb/app/utils"

	"fmt"
	"os"
	"os/exec"
)

func ParuAddRepo(s, g string, restore bool) {

	repo := s
	if sExt := utils.GetFileExtension(s); sExt == "sh" {
		repo = utils.GetFileName(s)
	} else {
		utils.PrintErrorMsgExit("Input Error:", "Arch based distros will only take bash scripts as the add-repo arg...")
	}

	if restore == true {
		exists, err := db.RepoExists("paru", repo)
		if err != nil {
			utils.PrintErrorExit("Error:", err)
		}
		if exists == true {
			utils.PrintErrorMsgExit("Error:", "This repo has already been added to app before...")
		}
	}

	if ws := utils.HasWhiteSpace(s); ws == true {
		utils.PrintErrorMsgExit("Input Error:", "Can't add more than one repo at a time...")
	}

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
	utils.CreateDirIfNotExist(env.DBDir + "/packages/repo/local/paru")
	utils.Copy(sFull, env.DBDir+"/packages/repo/local/paru/"+s)
	name := utils.GetFileName(s)

	fmt.Println("\n" + name + " has been added...\n")
	// Record added repo
	if err := db.RecordRepo("paru", name); err != nil {
		utils.PrintErrorExit("Repo Record Error:", err)
	}

}

func ParuRemoveRepo(s string) {

	repo := s
	if sExt := utils.GetFileExtension(s); sExt == "sh" {
		repo = utils.StripExtension(s, "."+sExt)
	} else {
		utils.PrintErrorMsgExit("Input Error:", "Arch based distros will only take bash scripts as the rm-repo arg...")
	}
	exists, err := db.RepoExists("paru", repo)
	if err != nil {
		utils.PrintErrorExit("Error:", err)
	}
	if exists == false {
		utils.PrintErrorMsgExit("Error:", "This repo does not exist or was not added by app...")
	}

	if ws := utils.HasWhiteSpace(s); ws == true {
		utils.PrintErrorMsgExit("Input Error:", "Can't remove more than one repo at a time...")
	}

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

	// Record removed repo
	if err = db.RemoveRepo("paru", "sh", repo); err != nil {
		utils.PrintErrorExit("Repo Remove Error:", err)
	}

	utils.DeleteDirIfEmpty(env.DBDir + "/packages/repo/local/paru")

	fmt.Println(repo + " has been removed...\n")

}
