package cli

import (
	"github.com/hkdb/app/env"
	"github.com/hkdb/app/utils"

	"flag"
	"fmt"
	"os"
)

var m = flag.String("m", "", 
	"Package Manager\n   usage: app -m <package manager> install neovim\n   default: auto-detect of native pkg manager <apt/dnf/pacman>\n   example: app install neovim\n   options:\n\t- apt\n\t- dnf\n\t- pacman\n\t- yay\n\t- flatpak\n\t- snap\n\t- brew\n\t- appimage\n")
var r = flag.String("r", "", 
	"Restore / Install all on new system\n   usage: app -r <type>\n   option:\n\t- apt\n\t- dnf\n\t- pacman\n\t- yay\n\t- flatpak\n\t- snap\n\t- brew\n\t- appimage\n\t- all\n")
var y = flag.Bool("y", false, 
	"Auto Yes - Skips the package manager confirmation (APT & DNF Only)\n   usage: app -y install neovim\n")
var gpg = flag.String("gpg", "", 
	"PGP Key URL - Used in combination with a url arg for add-repo (DNF Only)\n   usage: app -gpg <url> add-repo <url>\n")
var c = flag.String("c", "",
	"Channel - Used in combination with installing snap packages (SNAP Only)\n   usage: app -m snap -c <channel> install vlc\n   options:\n\t- beta\n\t- candidate\n\t- edge\n\t- stable\n")
var classic = flag.Bool("classic", false,
	"Classic Confinement for Snaps (SNAP Only)\n   usage: app -m snap -classic install flow\n")
var tag = flag.String("tag", "",
	"Tag (version) for cargo\n   usage: app -m cargo -tag <version> install <git url>\n   example: app -m cargo -tag 0.2.0 install https://github.com/donovanglover/hyprnome\n")

func ParseFlags() env.Flags {

	flag.Usage = usage
	flag.Parse()

	a := flag.Arg(0)
	p := flag.Arg(1)

	if flag.Arg(2) != "" || flag.Arg(3) != "" {
		utils.PrintErrorMsgExit("Input Error:", "If you are trying to specify multiple packages, wrap the pacakges with single quotes....")
	}


	// If -r flag is empty, then it's not a restore and therefore, we are either 
	// installing, upgrading, dist-upgrading, uninstalling, purging,  searching, or autoremoving
	if *r == "" {
		if a == "upgrade" && p == "all" && *m != "" {
			fmt.Println(utils.ColorRed, 
				"-m must not be specified when executing \"app upgrade all\"... Try", 
				utils.ColorYellow, "./app -h", utils.ColorRed, "to learn more...\n", utils.ColorReset)
			os.Exit(1)
		}
		if *m == "" {
			*m = defaultPkgMgr()
		}
		if a == "" {
			fmt.Println(utils.ColorRed, 
				"Action must be specified unless you are restoring the system with -r... Try", 
				utils.ColorYellow, "./app -h", utils.ColorRed, "to learn more...\n", utils.ColorReset)
			os.Exit(1)
		}
		if p == "" && a != "list" && a != "autoremove" && a != "update" && a != "history" && a != "upgrade" && a != "dist-upgrade" && a != "enable" && a != "disable" && a != "ls-repo"  && a != "settings"{
			fmt.Println(utils.ColorRed, 
				"Package(s) must be specified... Try", 
				utils.ColorYellow, "./app -h", utils.ColorRed, "to learn more...\n", utils.ColorReset)
			os.Exit(1)
		}
		if p != "" {
			if a == "autoremove" || a == "update" || a == "settings" {
				fmt.Println(utils.ColorRed, 
					"Package(s) must not be specified for these actions... Try", 
					utils.ColorYellow, "./app -h", utils.ColorRed, "to learn more...\n", utils.ColorReset)
				os.Exit(1)
			}
		}
		if a == "search" || a == "history" || a == "list" || a == "ls-repo" {
			if whitespace := utils.HasWhiteSpace(p); whitespace == true {
				utils.PrintErrorMsgExit("Input Error:", "Cannot specify more than one package in this action...\n")
			}
		}
		if a == "enable" || a == "disable" {
			if *m == "apt" || *m == "dnf" || *m == "pacman" {
				utils.PrintErrorMsgExit("Native package managers cannot be enabled/disabled...\n", "")
			}
		}
	} else {
		if *m != "" {
			fmt.Println(utils.ColorRed, 
				"-m (package manager) must not be specified when you are restoring the system with -r... Try", 
				utils.ColorYellow, "./app -h", utils.ColorRed, "to learn more...\n", utils.ColorReset)
			os.Exit(1)
		}
		if a != "" {
			fmt.Println(utils.ColorRed, 
				"-a (action) must not be specified when you are restoring the system with -r... Try", 
				utils.ColorYellow, "./app -h", utils.ColorRed, "to learn more...\n", utils.ColorReset)
			os.Exit(1)
		}
		if p != "" {
			fmt.Println(utils.ColorRed, 
				"-p (package(s)) must not be specified when you are restoring the system with -r... Try", 
				utils.ColorYellow, "./app -h", utils.ColorRed, "to learn more...\n", utils.ColorReset)
			os.Exit(1)
		}
		if *c != "" {
			fmt.Println(utils.ColorRed, 
				"-c (channel) must not be specified when you are restoring the system with -r... Try", 
				utils.ColorYellow, "./app -h", utils.ColorRed, "to learn more...\n", utils.ColorReset)
			os.Exit(1)
		}
		if *gpg != "" {
			fmt.Println(utils.ColorRed, 
				"-gpg (gpg key url) must not be specified when you are restoring the system with -r... Try", 
				utils.ColorYellow, "./app -h", utils.ColorRed, "to learn more...\n", utils.ColorReset)
			os.Exit(1)
		}
	}

	if *y == true {
		if *m != "apt" && *m != "dnf" {
			if a != "install" && a != "remove" && a != "purge" && a != "upgrade" && a != "dist-upgrade" {
				utils.PrintErrorMsgExit("Error: -y can only be used to install, remove, purge or upgrade...", "")
			} 
			utils.PrintErrorMsgExit("Error: -y can only be used with native package managers...", "")
		}
		env.AutoYes = true
	}

	if *gpg != "" {
		if env.Base != "redhat" {
			utils.PrintErrorMsgExit("Error: -gpg should only be used with adding repos for Redhat based systems...", "")
		}
		isUrl := utils.IsUrl(*gpg)
		if isUrl == false {
			utils.PrintErrorMsgExit("Error: -gpg can only a url as argument...", "")
		} 
	}

	if *c != "" {
		if *m != "snap" || a != "install" {
			utils.PrintErrorMsgExit("Error: -c should only be used with installing packages with snap...", "")
		}
	}

	if *classic != false {
		if *m != "snap" || a != "install" {
			utils.PrintErrorMsgExit("Error: -c should only be used with installing packages with snap...", "")
		}
	}

	if a == "settings" {
		if p != "" && *r != "" && *gpg != "" && *y != false && *classic != false && *c != "" {
			utils.PrintErrorMsgExit("Error: The \"settings\" action can't be executed with any flags or options...", "")
		}
	}
	
	if *tag != "" {
		if *m != "cargo" || a != "install" {
			utils.PrintErrorMsgExit("Error: This flag is only available for installing git url with cargo...", "")
		}
	}

	f := env.Flags{}
	f.A = a
	f.P = p
	f.M = *m
	f.R = *r
	f.G = *gpg
	f.C = *c
	f.Classic = *classic
	f.Tag = *tag
	return f

}

func usage() {
	fmt.Printf("USAGE:\n\tapp [OPTIONS] <ACTION> <PACKAGE>")
	fmt.Printf("\n\tEXAMPLE:")
	fmt.Printf("\n\t\tapp install neovim")
	fmt.Printf("\n\t\tapp -m flatpak install Geary")
	fmt.Printf("\n\n")

	fmt.Printf("ACTIONS:\n\t- install ~ Install package. Takes package name as argument\n\t- remove ~ Uninstall package. Takes package name as argument\n\t- update ~ Refreshes repos\n\t- upgrade ~ Upgrade packages. Takes \"all\" as a value to upgrade with all package managers\n\t- dist-upgrade ~ A more advanced upgrade that can add or remove packages during upgrade (APT Only)\n\t- autoremove ~ Remove dependency packages that are no longer required\n\t- purge ~ Same as remove but removes configs too (APT only)\n\t- search ~ Search for packages in repos\n\t- list ~ List packages installed on system. Greps for package if argument is provided.\n\t- history ~ List pacakges installed by app. Takes package name as argument to search.\n\t- enable ~ Enable Package Manager (yay, flatpak, snap, brew, appimage)\n\t- disable ~ Disable Package Manager (yay, flatpak, snap, brew, appimage)\n\t- add-repo ~ Add package manager repo. Takes a .sh, ppa, or url as argument.\n\t- rm-repo ~ Remove package manager repo. Takes repo identifier as argument\n\t- ls-repo ~ List package manager repos\n\t- settings ~ List settings including the status of packages managers (enabled/disabled)\n\n")
	fmt.Printf("PACKAGE:\n\tPackage name(s). For multiple packages, wrap the argument with quotes.")
	fmt.Printf("\n\tEXAMPLE:\n\t\tapp install 'neovim whois nmap'")
	fmt.Printf("\n\n")
	fmt.Printf("OPTIONS:\n")
	flag.PrintDefaults()
}

func defaultPkgMgr() string {
	
	switch env.OSType {
	case "Linux":
		switch env.Base {
		case "debian":
			return "apt"
		case "redhat":
			return "dnf"
		case "arch":
			return "pacman"
		default:
			fmt.Println(utils.ColorRed, "\nDistro not supported...\n", utils.ColorReset)
			os.Exit(1)
		}
	case "Mac":
		return "brew"
	case "Windows":
		return "scoop"
	default:
		fmt.Println(utils.ColorRed, "\nUnsupported OS...\n", utils.ColorReset)
		os.Exit(1)
	}
	return ""

}

