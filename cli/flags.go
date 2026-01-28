package cli

import (
	"github.com/hkdb/app/env"
	"github.com/hkdb/app/utils"

	"flag"
	"fmt"
	"os"
)

var m = flag.String("m", "",
	"Package Manager\n   usage: app -m <package manager> install neovim\n   default: auto-detect of native pkg manager <apt/dnf/pacman/pkg>\n   example: app install neovim\n   options:\n\t- apt\n\t- dnf\n\t- pacman\n\t- yay\n\t- paru\n\t- pkg\n\t- zypper\n\t- flatpak\n\t- snap\n\t- brew\n\t- go\n\t- pip\n\t- pipx\n\t- cargo\n\t- appimage\n\t- app (ONLY FOR UPDATING app ITSELF)\n")
var r = flag.String("r", "",
	"Restore / Install all on new system\n   usage: app -r <type>\n   option:\n\t- apt\n\t- dnf\n\t- pacman\n\t- yay\n\t- paru\n\t- pkg\n\t- zypper\n\t- flatpak\n\t- snap\n\t- brew\n\t- go\n\t- pip\n\t- pipx\n\t- cargo\n\t- appimage\n\t- all\n")
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
var raw = flag.Bool("raw", false,
	"Raw history output instead of alphabetically sorted output\n   usage: app -raw history\n   example: app -raw history\n")


func ParseFlags() env.Flags {

	flag.Usage = usage
	flag.Parse()

	f := env.Flags{}
	f.A = flag.Arg(0)
	f.P = flag.Arg(1)
	f.M = *m
	f.R = *r
	f.G = *gpg
	f.C = *c
	f.Classic = *classic
	f.Tag = *tag

	basicReminders(flag.CommandLine, f)

	// If -r flag is empty, then it's not a restore and therefore, we are either
	// installing, upgrading, dist-upgrading, uninstalling, purging, searching, or autoremoving
	if *r == "" {
		f = actionCheck(f)
	} else {
		restoreCheck(f)
	}

	nativePkgMgrCheck(*m)

	subFlagsCheck(f)

	f = sortCheck(*raw, f)

	return f

}

func usage() {
	fmt.Printf("USAGE:\n\tapp [OPTIONS] <ACTION> <PACKAGE>")
	fmt.Printf("\n\tEXAMPLE:")
	fmt.Printf("\n\t\tapp install neovim")
	fmt.Printf("\n\t\tapp -m flatpak install Geary")
	fmt.Printf("\n\n")

	fmt.Printf("ACTIONS:\n\t- install ~ Install package. Takes package name as argument\n\t- remove ~ Uninstall package. Takes package name as argument\n\t- update ~ Refreshes repos\n\t- upgrade ~ Upgrade packages. Takes \"all\" as a value to upgrade with all package managers\n\t- dist-upgrade ~ A more advanced upgrade that can add or remove packages during upgrade (APT Only)\n\t- autoremove ~ Remove dependency packages that are no longer required\n\t- purge ~ Same as remove but removes configs too (APT only)\n\t- search ~ Search for packages in repos\n\t- list ~ List packages installed on system. Greps for package if argument is provided.\n\t- history ~ List pacakges installed by app. Takes package name as argument to search.\n\t- enable ~ Enable Package Manager (yay, paru, flatpak, snap, brew, appimage)\n\t- disable ~ Disable Package Manager (pacman, yay, paru, flatpak, snap, brew, appimage)\n\t- add-repo ~ Add package manager repo. Takes a .sh, ppa, or url as argument.\n\t- rm-repo ~ Remove package manager repo. Takes repo identifier as argument\n\t- ls-repo ~ List package manager repos\n\t- settings ~ List settings including the status of packages managers (enabled/disabled)\n\n")
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
		case "suse":
			return "zypper"
		case "arch":
			if env.Pacman == false {
				if env.Yay == true {
					return "yay"
				}
				if env.Paru == true {
					return "paru"
				}
			}
			return "pacman"
		case "nixos":
			return "nix-env"
		default:
			fmt.Println(utils.ColorRed, "\nDistro not supported...\n", utils.ColorReset)
			os.Exit(1)
		}
	case "Mac":
		return "brew"
	case "FreeBSD":
		return "pkg"
	case "Windows":
		return "scoop"
	default:
		fmt.Println(utils.ColorRed, "\nUnsupported OS...\n", utils.ColorReset)
		os.Exit(1)
	}
	return ""

}
