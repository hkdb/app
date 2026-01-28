
package cli

import (
	"github.com/hkdb/app/env"
	"github.com/hkdb/app/utils"
	
	"flag"
	"fmt"
	"os"
)

func basicReminders(flag *flag.FlagSet, f env.Flags) {

	// For app update
	if f.M == "app" && f.A != "update" {
		utils.PrintErrorMsgExit("Input Error:", "app can only be specified to update itself...")
	}

	// Adding reminder for multiple packages without quoting
	if flag.Arg(2) != "" || flag.Arg(3) != "" {
		utils.PrintErrorMsgExit("Input Error:", "If you are trying to specify multiple packages, wrap the pacakges with single quotes....")
	}

	// Adding enable reminder
	if f.A == "enable" || f.A == "disable" {
		if f.P != "" {
			utils.PrintErrorMsgExit("Syntax Error:", "The proper command to " + f.A + " " + f.P + " is:\n\n\t\tapp -m " + f.P + " " + f.A)
		}
	}

}

func restoreCheck(f env.Flags) {

	if f.M != "" {
		fmt.Println(utils.ColorRed,
			"-m (package manager) must not be specified when you are restoring the system with -r... Try",
			utils.ColorYellow, "./app -h", utils.ColorRed, "to learn more...\n", utils.ColorReset)
		os.Exit(1)
	}
	if f.A != "" {
		fmt.Println(utils.ColorRed,
			"-a (action) must not be specified when you are restoring the system with -r... Try",
			utils.ColorYellow, "./app -h", utils.ColorRed, "to learn more...\n", utils.ColorReset)
		os.Exit(1)
	}
	if f.P != "" {
		fmt.Println(utils.ColorRed,
			"package(s) must not be specified when you are restoring the system with -r... Try",
			utils.ColorYellow, "./app -h", utils.ColorRed, "to learn more...\n", utils.ColorReset)
		os.Exit(1)
	}
	if f.C != "" {
		fmt.Println(utils.ColorRed,
			"-c (channel) must not be specified when you are restoring the system with -r... Try",
			utils.ColorYellow, "./app -h", utils.ColorRed, "to learn more...\n", utils.ColorReset)
		os.Exit(1)
	}
	if f.G != "" {
		fmt.Println(utils.ColorRed,
			"-gpg (gpg key url) must not be specified when you are restoring the system with -r... Try",
			utils.ColorYellow, "./app -h", utils.ColorRed, "to learn more...\n", utils.ColorReset)
		os.Exit(1)
	}

}

func actionCheck(f env.Flags) env.Flags {

	if f.A == "upgrade" && f.P == "all" && f.M != "" {
		fmt.Println(utils.ColorRed,
			"-m must not be specified when executing \"app upgrade all\"... Try",
			utils.ColorYellow, "./app -h", utils.ColorRed, "to learn more...\n", utils.ColorReset)
		os.Exit(1)
	}
	if f.M == "" {
		f.M = defaultPkgMgr()
	}
	if f.A == "" {
		fmt.Println(utils.ColorRed,
			"Action must be specified unless you are restoring the system with -r... Try",
			utils.ColorYellow, "./app -h", utils.ColorRed, "to learn more...\n", utils.ColorReset)
		os.Exit(1)
	}
	if f.P == "" && f.A != "list" && f.A != "autoremove" && f.A != "update" && f.A != "history" && f.A != "upgrade" && f.A != "dist-upgrade" && f.A != "enable" && f.A != "disable" && f.A != "ls-repo" && f.A != "settings" {
		fmt.Println(utils.ColorRed,
			"Package(s) must be specified... Try",
			utils.ColorYellow, "./app -h", utils.ColorRed, "to learn more...\n", utils.ColorReset)
		os.Exit(1)
	}
	if f.P != "" {
		if f.A == "autoremove" || f.A == "update" || f.A == "settings" {
			fmt.Println(utils.ColorRed,
				"Package(s) must not be specified for these actions... Try",
				utils.ColorYellow, "./app -h", utils.ColorRed, "to learn more...\n", utils.ColorReset)
			os.Exit(1)
		}
	}
	if f.A == "search" || f.A == "history" || f.A == "list" || f.A == "ls-repo" {
		if whitespace := utils.HasWhiteSpace(f.P); whitespace == true {
			utils.PrintErrorMsgExit("Input Error:", "Cannot specify more than one package in this action...\n")
		}
	}

	return f

}

func nativePkgMgrCheck(m string) {

	// Check if it's the right native package manager
	switch m {
	case "apt":
		if env.Base != "debian" {
			utils.PrintErrorMsgExit("Error: wrong native package manager specified...\n", "")
		}
	case "dnf":
		if env.Base != "redhat" {
			utils.PrintErrorMsgExit("Error: wrong native package manager specified...\n", "")
		}
	case "pacman":
		if env.Base != "arch" {
			utils.PrintErrorMsgExit("Error: wrong native package manager specified...\n", "")
		}
	case "zypper":
		if env.Base != "suse" {
			utils.PrintErrorMsgExit("Error: wrong native package manager specified...\n", "")
		}
	case "pkg":
		if env.Base != "freebsd" {
			utils.PrintErrorMsgExit("Error: wrong native package manager specified...\n", "")
		}
	}

}

func subFlagsCheck(f env.Flags) {

	if f.Y == true {
		if f.M != "apt" && f.M != "dnf" && f.M != "zypper" {
			if f.A != "install" && f.A != "remove" && f.A != "purge" && f.A != "upgrade" && f.A != "dist-upgrade" {
				utils.PrintErrorMsgExit("Error: -y can only be used to install, remove, purge or upgrade...", "")
			}
			utils.PrintErrorMsgExit("Error: -y can only be used with native package managers...", "")
		}
		env.AutoYes = true
	}

	if f.G != "" {
		if env.Base != "redhat" && env.Base != "suse" {
			utils.PrintErrorMsgExit("Error: -gpg should only be used with adding repos for Redhat & Suse based systems...", "")
		}
		isUrl := utils.IsUrl(f.G)
		if isUrl == false {
			utils.PrintErrorMsgExit("Error: -gpg can only a url as argument...", "")
		}
	}

	classicCheck(f)

	if f.A == "settings" {
		if f.P != "" && f.R != "" && f.G != "" && f.Y != false && f.Classic != false && f.C != "" {
			utils.PrintErrorMsgExit("Error: The \"settings\" action can't be executed with any flags or options...", "")
		}
	}

	if f.Tag != "" {
		if f.M != "cargo" || f.A != "install" {
			utils.PrintErrorMsgExit("Error: This flag is only available for installing git url with cargo...", "")
		}
	}

}

func classicCheck(f env.Flags) {

	if f.C != "" {
		if f.M != "snap" || f.A != "install" {
			utils.PrintErrorMsgExit("Error: -c should only be used with installing packages with snap...", "")
		}
	}

	if f.Classic != false {
		if f.M != "snap" || f.A != "install" {
			utils.PrintErrorMsgExit("Error: -classic should only be used with installing packages with snap...", "")
		}
	}

}

func sortCheck(raw bool, f env.Flags) env.Flags {

	sort := true
	if raw == true {
		sort = false
	}
	if sort != true && f.A != "history" {
		utils.PrintErrorMsgExit("Error: This flag is only available for the history action...", "")
	}

	f.Sort = sort

	return f

}
