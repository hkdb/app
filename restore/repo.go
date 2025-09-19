package restore

import (
	"github.com/hkdb/app/arch"
	"github.com/hkdb/app/db"
	"github.com/hkdb/app/debian"
	"github.com/hkdb/app/env"
	"github.com/hkdb/app/flatpak"
	"github.com/hkdb/app/redhat"
	"github.com/hkdb/app/suse"
	"github.com/hkdb/app/utils"

	"fmt"
	"os"
)

func RestoreAllRepos(pm string) {

	pmrDir := env.DBDir + "/packages/repo/" + pm + ".json"
	if _, err := os.Stat(pmrDir); os.IsNotExist(err) {
		fmt.Println("There are no repos added by app... Moving on...")
		return
	}

	repos, err := db.ReadPkgSlice("packages", "repo", pm)
	if err != nil {
		utils.PrintErrorExit("Read Repo List Error:", err)
	}

	for i := 1; i < len(repos); i++ {
		name := repos[i]
		rType := utils.GetRepoRestoreType(pm, repos[i])
		switch pm {
		case "apt":
			switch rType {
			case "ppa":
				debian.AddRepo(name, "")
			case "sh":
				debian.AddRepo(env.DBDir+"/packages/repo/local/"+pm+"/"+name+".sh", "")
			case "json":
				utils.PrintErrorMsgExit("Url/Gpg repo records are not supported for Debian based distros...", "")
			default:
				utils.PrintErrorMsgExit("Repo Records Error:", "Unrecognized record type for Debian based distros...")
			}
		case "dnf":
			switch rType {
			case "sh":
				redhat.AddRepo(env.DBDir+"/packages/repo/local/"+pm+"/"+name+".sh", "")
			case "json":
				r, g, err := db.GetSetup(pm, repos[i])
				if err != nil {
					utils.PrintErrorExit("Read Repo Error:", err)
				}
				redhat.AddRepo(r, g)
			default:
				utils.PrintErrorMsgExit("Repo Records Error:", "Unrecognized record type for Redhat based distros...")
			}
		case "zypper":
			switch rType {
			case "sh":
				suse.AddRepo(env.DBDir+"/packages/repo/local/"+pm+"/"+name+".sh", "")
			case "json":
				r, g, err := db.GetSetup(pm, repos[i])
				if err != nil {
					utils.PrintErrorExit("Read Repo Error:", err)
				}
				suse.AddRepo(r, g)
			default:
				utils.PrintErrorMsgExit("Repo Records Error:", "Unrecognized record type for Suse based distros...")
			}
		case "pacman":
			switch rType {
			case "sh":
				arch.AddRepo(env.DBDir+"/packages/repo/local/"+pm+"/"+name+".sh", "")
			default:
				utils.PrintErrorMsgExit("Repo Records Error:", "Unrecognized record type for Arch based distros...")
			}
		case "yay":
			switch rType {
			case "sh":
				arch.YayAddRepo(env.DBDir+"/packages/repo/local/"+pm+"/"+name+".sh", "")
			default:
				utils.PrintErrorMsgExit("Repo Records Error:", "Unrecognized record type for Arch based distros...")
			}
		case "paru":
			switch rType {
			case "sh":
				arch.ParuAddRepo(env.DBDir+"/packages/repo/local/"+pm+"/"+name+".sh", "")
			default:
				utils.PrintErrorMsgExit("Repo Records Error:", "Unrecognized record type for Arch based distros...")
			}
		case "flatpak":
			switch rType {
			case "json":
				r, g, err := db.GetSetup(pm, repos[i])
				if err != nil {
					utils.PrintErrorExit("Read Repo Error:", err)
				}
				flatpak.AddRepo(r, g)
			default:
				utils.PrintErrorMsgExit("Repo Records Error:", "Unrecognized record type for Flatpak...")
			}
		case "snap":
			utils.PrintErrorMsgExit("Error:", "Adding third party repos for snap is not supported...")
		default:
			utils.PrintErrorMsg("Error:", "Unsupported package manager")
		}
	}

}
