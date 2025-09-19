package mgr

import (
	"github.com/hkdb/app/appimage"
	"github.com/hkdb/app/arch"
	"github.com/hkdb/app/cargo"
	"github.com/hkdb/app/debian"
	"github.com/hkdb/app/env"
	"github.com/hkdb/app/flatpak"
	"github.com/hkdb/app/freebsd"
	"github.com/hkdb/app/golang"
	brew "github.com/hkdb/app/macos"
	"github.com/hkdb/app/pip"
	"github.com/hkdb/app/redhat"
	"github.com/hkdb/app/suse"
	"github.com/hkdb/app/restore"
	"github.com/hkdb/app/snap"
	"github.com/hkdb/app/utils"

	"fmt"
	"os"
)

func Process(flag env.Flags) {

	a := flag.A
	p := flag.P
	m := flag.M
	r := flag.R
	g := flag.G
	c := flag.C
	classic := flag.Classic
	tag := flag.Tag
	sort := flag.Sort

	if r != "" {
		if r == "all" {
			restoreAll()
		} else {
			restoreOne(r)
		}
	} else {
		if a != "enable" && a != "disable" {
			enabled := isEnabled(m)
			if enabled == false {
				utils.PrintErrorMsgExit(m+" has been disabled. You must first enable it before using it in app...\n", "")
			}
		}
		execute(m, a, p, g, c, classic, tag, sort)
	}

}

func restoreAll() {

	switch env.OSType {
	case "Linux":
		pm := utils.GetNativePkgMgr()
		nable := false
		switch pm {
		case "apt":
			if env.Apt == true {
				nable = true
			} 
		case "dnf":
			if env.Dnf == true {
				nable = true
			}
		case "pacman":
			if env.Pacman == true {
				nable = true
			}
		case "zypper":
			if env.Zypper == true {
				nable = true
			}
		}
		if nable == true {
			fmt.Print("Restore all " + pm + " repos and packages? (Y/n) ")
			native := utils.Confirm()
			if native == true {
				switch env.Base {
				case "debian":
					restore.RestoreAllRepos("apt")
					debian.InstallAll()
					//fmt.Println("Debain Install All")
				case "redhat":
					restore.RestoreAllRepos("dnf")
					redhat.InstallAll()
					//fmt.Println("Redhat Install All")
				case "arch":
					restore.RestoreAllRepos("pacman")
					arch.InstallAll()
					//fmt.Println("Arch Install All")
				case "suse":
					restore.RestoreAllRepos("suse")
					arch.InstallAll()
				case "nixos":
					fmt.Println("Skipping nix-env... It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
				}
			}
		}
		
		if env.Base == "arch" && env.Yay == true {
			fmt.Print("Restore all yay repos and packages? (Y/n) ")
			fp := utils.Confirm()
			if fp == true {
				restore.RestoreAllRepos("yay")
				arch.YayInstallAll()
			}
		}

		if env.Base == "arch" && env.Paru == true {
			fmt.Print("Restore all paru repos and packages? (Y/n) ")
			fp := utils.Confirm()
			if fp == true {
				restore.RestoreAllRepos("paru")
				arch.ParuInstallAll()
			}
		}

		if env.Flatpak == true {
			fmt.Print("Restore all Flatpak repos and apps? (Y/n) ")
			fp := utils.Confirm()
			if fp == true {
				restore.RestoreAllRepos("flatpak")
				flatpak.InstallAll()
			}
		}

		if env.Snap == true {
			fmt.Print("Restore all Snap apps? (Y/n) ")
			s := utils.Confirm()
			if s == true {
				snap.InstallAll()
			}
		}

		if env.Brew == true {
			fmt.Print("Restore all Homebrew apps? (Y/n) ")
			s := utils.Confirm()
			if s == true {
				brew.InstallAll()
			}
		}

		if env.AppImage == true {
			fmt.Print("Restore all AppImage apps? (Y/n) ")
			s := utils.Confirm()
			if s == true {
				appimage.InstallAll()
			}
		}
	case "Mac":
		fmt.Print("Restore all Homebrew apps? (Y/n) ")
		s := utils.Confirm()
		if s == true {
			brew.InstallAll()
		}
	case "FreeBSD":
		fmt.Print("Restore all PKG apps? (Y/n) ")
		s := utils.Confirm()
		if s == true {
			freebsd.InstallAll()
		}
	case "Windows":
		fmt.Println("Not implemented yet... Coming Soon!")
	default:
		utils.PrintErrorMsgExit("OS not supported...", "")
	}

	if env.Go == true {
		fmt.Print("Restore all Go apps? (Y/n) ")
		s := utils.Confirm()
		if s == true {
			golang.InstallAll()
		}
	}

	if env.Pip == true {
		fmt.Print("Restore all Python(pip) apps? (Y/n) ")
		s := utils.Confirm()
		if s == true {
			pip.InstallAll()
		}
	}

	if env.Cargo == true {
		fmt.Print("Restore all Rust(cargo) apps? (Y/n) ")
		s := utils.Confirm()
		if s == true {
			cargo.InstallAll()
		}
	}

	fmt.Println("")

}

func restoreOne(r string) {
	switch env.OSType {
	case "Linux":
		switch r {
		case "apt":
			npm := utils.GetNativePkgMgr()
			if npm == "apt" {
				restore.RestoreAllRepos("apt")
				debian.InstallAll()
			} else {
				utils.PrintErrorMsgExit("Error:", "You are trying to restore with apt on a non-apt system...\n")
			}
		case "dnf":
			npm := utils.GetNativePkgMgr()
			if npm == "dnf" {
				restore.RestoreAllRepos("dnf")
				redhat.InstallAll()
			} else {
				utils.PrintErrorMsgExit("Error:", "You are trying to restore with dnf on a non-dnf system...\n")
			}
		case "pacman":
			if env.Pacman == false {
				utils.PrintErrorMsgExit(r+" is disabled... To enable it, run", "app -m "+r+" enable\n")
			}
			npm := utils.GetNativePkgMgr()
			if npm == "pacman" {
				restore.RestoreAllRepos("pacman")
				arch.InstallAll()
			} else {
				utils.PrintErrorMsgExit("Error:", "You are trying to restore with pacman on a non-Arch system...\n")
			}
		case "yay":
			if env.Yay == false {
				utils.PrintErrorMsgExit(r+" is disabled... To enable it, run", "app -m "+r+" enable\n")
			}
			npm := utils.GetNativePkgMgr()
			if npm == "pacman" {
				arch.YayInstallAll()
			} else {
				utils.PrintErrorMsgExit("Error:", "You are trying to restore with yay on a non-Arch system...\n")
			}
		case "paru":
			if env.Paru == false {
				utils.PrintErrorMsgExit(r+" is disabled... To enable it, run", "app -m "+r+" enable\n")
			}
			npm := utils.GetNativePkgMgr()
			if npm == "pacman" {
				arch.ParuInstallAll()
			} else {
				utils.PrintErrorMsgExit("Error:", "You are trying to restore with paru on a non-Arch system...\n")
			}
		case "zypper":
			npm := utils.GetNativePkgMgr()
			if npm == "zypper" {
				restore.RestoreAllRepos("zypper")
				redhat.InstallAll()
			} else {
				utils.PrintErrorMsgExit("Error:", "You are trying to restore with dnf on a non-zypper system...\n")
			}
		case "nix-env":
			npm := utils.GetNativePkgMgr()
			if npm == "nix-env" {
				utils.PrintErrorMsgExit("Error: It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
			}
		case "flatpak":
			if env.Flatpak == false {
				utils.PrintErrorMsgExit(r+" is disabled... To enable it, run", "app -m "+r+" enable\n")
			}
			restore.RestoreAllRepos("flatpak")
			flatpak.InstallAll()
		case "snap":
			if env.Snap == false {
				utils.PrintErrorMsgExit(r+" is disabled... To enable it, run ", "app -m "+r+" enable\n")
			}
			snap.InstallAll()
		case "appimage":
			if env.AppImage == false {
				utils.PrintErrorMsgExit(r+" is disabled... To enable it, run ", "app -m "+r+" enable\n")
			}
			appimage.InstallAll()
		case "brew":
			if env.Brew == false {
				utils.PrintErrorMsgExit(r+" is disabled... To enable it, run ", "app -m "+r+" enable\n")
			}
			brew.InstallAll()
		case "go":
			if env.Go == false {
				utils.PrintErrorMsgExit(r+" is disabled... To enable it, run ", "app -m "+r+" enable\n")
			}
			brew.InstallAll()
		case "pip":
			if env.Pip == false {
				utils.PrintErrorMsgExit(r+" is disabled... To enable it, run ", "app -m "+r+" enable\n")
			}
			pip.InstallAll()
		case "cargo":
			if env.Cargo == false {
				utils.PrintErrorMsgExit(r+" is disabled... To enable it, run ", "app -m "+r+" enable\n")
			}
			cargo.InstallAll()
		case "scoop":
			utils.PrintErrorMsgExit("Error:", "Not implemented yet... Coming Soon!")
		default:
			utils.PrintErrorMsgExit("Error:", "Package manager not supported...")
		}
	case "Mac":
		if r != "brew" {
			utils.PrintErrorMsgExit("Error:", "Package Manager not supported...")
		}
	case "FreeBSD":
		if r != "pkg" {
			utils.PrintErrorMsgExit("Error:", "Package Manager not supported...")
		}
	case "Windows":
		utils.PrintErrorMsgExit("Error:", "Windows support is not implemented yet...")
	default:
		utils.PrintErrorMsgExit("Error:", "Unspported operating system...")
	}

}

func execute(m, a, p, g, c string, classic bool, tag string, sort bool) {

	if env.OSType == "Mac" && m != "brew" && m != "app" {
		utils.PrintErrorMsgExit("Error:", "macOS currently only supports Homebrew...")
	}

	switch a {
	case "install":
		switch m {
		case "apt":
			debian.Install(p)
		case "dnf":
			redhat.Install(p)
		case "pacman":
			arch.Install(p)
		case "yay":
			arch.YayInstall(p)
		case "paru":
			arch.ParuInstall(p)
		case "pkg":
			freebsd.Install(p)
		case "zypper":
			suse.Install(p)
		case "nix-env":
			utils.PrintErrorMsgExit("Error: It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
		case "flatpak":
			flatpak.Install(p)
		case "snap":
			snap.Install(p, c, classic)
		case "brew":
			brew.Install(p)
		case "go":
			golang.Install(p)
		case "pip":
			pip.Install(p)
		case "cargo":
			cargo.Install(p, tag)
		case "appimage":
			appimage.Install(p)
		default:
			fmt.Println("Unsupported package manager... Exiting...\n")
			os.Exit(1)
		}
	case "remove":
		switch m {
		case "apt":
			debian.Remove(p)
		case "dnf":
			redhat.Remove(p)
		case "pacman":
			arch.Remove(p)
		case "yay":
			arch.YayRemove(p)
		case "paru":
			arch.ParuRemove(p)
		case "pkg":
			freebsd.Remove(p)
		case "zypper":
			suse.Remove(p)
		case "nix-env":
			utils.PrintErrorMsgExit("Error: It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
		case "flatpak":
			flatpak.Remove(p)
		case "snap":
			snap.Remove(p)
		case "brew":
			brew.Remove(p)
		case "go":
			golang.Remove(p)
		case "pip":
			pip.Remove(p)
		case "cargo":
			cargo.Remove(p)
		case "appimage":
			appimage.Remove(p)
		default:
			fmt.Println("Unsupported package manager... Exiting...\n")
			os.Exit(1)
		}
	case "purge":
		switch m {
		case "apt":
			debian.Purge(p)
		case "dnf":
			redhat.Purge(p)
		case "pacman":
			arch.Purge(p)
		case "yay":
			arch.YayPurge(p)
		case "paru":
			arch.ParuPurge(p)
		case "pkg":
			freebsd.Purge(p)
		case "zypper":
			suse.Purge(p)
		case "nix-env":
			utils.PrintErrorMsgExit("Error: It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
		case "flatpak":
			flatpak.Purge(p)
		case "snap":
			snap.Purge(p)
		case "brew":
			brew.Purge(p)
		case "go":
			golang.Purge(p)
		case "pip":
			pip.Purge(p)
		case "cargo":
			cargo.Purge(p)
		default:
			fmt.Println("Unsupported package manager... Exiting...\n")
			os.Exit(1)
		}
	case "search":
		switch m {
		case "apt":
			debian.Search(p)
		case "dnf":
			redhat.Search(p)
		case "pacman":
			arch.Search(p)
		case "yay":
			arch.YaySearch(p)
		case "paru":
			arch.ParuSearch(p)
		case "pkg":
			freebsd.Search(p)
		case "zypper":
			suse.Search(p)
		case "nix-env":
			utils.PrintErrorMsgExit("Error: It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
		case "flatpak":
			flatpak.Search(p)
		case "snap":
			snap.Search(p)
		case "brew":
			brew.Search(p)
		case "go":
			golang.Search(p)
		case "pip":
			pip.Search(p)
		case "cargo":
			cargo.Search(p)
		default:
			fmt.Println("Unsupported package manager... Exiting...\n")
			os.Exit(1)
		}
	case "update":
		switch m {
		case "apt":
			debian.Update()
		case "dnf":
			redhat.Update()
		case "pacman":
			arch.Update()
		case "yay":
			arch.YayUpdate()
		case "paru":
			arch.ParuUpdate()
		case "pkg":
			freebsd.Update()
		case "zypper":
			suse.Update()
		case "nix-env":
			utils.PrintErrorMsgExit("Error: It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
		case "flatpak":
			flatpak.Update()
		case "snap":
			snap.Update()
		case "brew":
			brew.Update()
		case "go":
			golang.Update()
		case "pip":
			pip.Update()
		case "cargo":
			cargo.Update()
		case "app":
			utils.UpdateApp()
		default:
			fmt.Println("Unsupported package manager... Exiting...\n")
			os.Exit(1)
		}
	case "upgrade":
		switch p {
		case "":
			switch m {
			case "apt":
				debian.Upgrade()
			case "dnf":
				redhat.Upgrade()
			case "pacman":
				arch.Upgrade()
			case "yay":
				arch.YayUpgrade()
			case "paru":
				arch.ParuUpgrade()
			case "pkg":
				freebsd.Upgrade()
			case "zypper":
				suse.Upgrade()
			case "nix-env":
				utils.PrintErrorMsgExit("Error: It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
			case "flatpak":
				flatpak.Upgrade()
			case "snap":
				snap.Upgrade()
			case "brew":
				brew.Upgrade()
			case "go":
				golang.Upgrade()
			case "pip":
				pip.Upgrade()
			case "cargo":
				cargo.Upgrade()
			default:
				fmt.Println("Unsupported package manager... Exiting...\n")
				os.Exit(1)
			}
		case "all":
			nabled := utils.IsNativeEnabled()
			if nabled == true {
				switch m {
				case "apt":
					fmt.Println("Updating APT repos:\n")
					debian.Update()
					fmt.Println("\nUpgrading with APT:\n")
					debian.Upgrade()
				case "dnf":
					fmt.Println("\nUpgrading with DNF:\n")
					redhat.Update()
					redhat.Upgrade()
				case "pacman":
					fmt.Println("\nUpgrading with PACMAN:\n")
					arch.Upgrade()
				case "pkg":
					fmt.Println("\nUpgrading with PKG:\n")
					freebsd.Update()
					freebsd.Upgrade()
				case "zypper":
					fmt.Println("\nUpgrading with ZYPPER:\n")
					suse.Update()
					suse.Upgrade()
				case "nix-env":
					utils.PrintErrorMsgExit("Error: It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
				case "brew":
					if env.OSType != "Mac" {
						utils.PrintErrorMsgExit("Unsupported OS/Distro....\n", "")
					}
				default:
					utils.PrintErrorMsgExit("Unsupported OS/Distro....\n", "")
				}
			}

			if env.OSType == "Linux" {
				npm := utils.GetNativePkgMgr()
				if npm == "pacman" && env.Yay == true {
					fmt.Println("\nUpgrade with YAY:\n")
					arch.YayUpgrade()
				}
				if npm == "pacman" && env.Yay == false {
					fmt.Println("Yay is disabled... Skipping...\n")
				}
				if npm == "pacman" && env.Paru == true {
					fmt.Println("\nUpgrade with PARU:\n")
					arch.ParuUpgrade()
				}
				if npm == "pacman" && env.Paru == false {
					fmt.Println("Paru is disabled... Skipping...\n")
				}
				if env.Flatpak != false {
					fmt.Println("\nUpgrading with FLATPAK:\n")
					flatpak.Upgrade()
				} else {
					fmt.Println("Flatpak is disabled... Skipping...\n")
				}
				if env.Snap != false {
					fmt.Println("\nUpgrading with SNAP:\n")
					snap.Upgrade()
				} else {
					fmt.Println("Snap is disabled... Skipping...\n")
				}
			}
			if env.Brew != false {
				fmt.Println("\nUpgrading with HOMEBREW:\n")
				brew.Upgrade()
			} else {
				fmt.Println("Brew is disabled... Skipping...\n")
			}
			if env.Go != false {
				fmt.Println("\nUpgrading GO binaries:\n")
				golang.Upgrade()
			} else {
				fmt.Println("Go is disabled... Skipping...\n")
			}
			if env.Pip != false {
				fmt.Println("\nUpgrading with PIP installed packages:\n")
				pip.Upgrade()
			} else {
				fmt.Println("PIP is disabled... Skipping...\n")
			}
			if env.Cargo != false {
				fmt.Println("\nUpgrading with CARGO:\n")
				cargo.Upgrade()
			} else {
				fmt.Println("Cargo is disabled... Skipping...\n")
			}
		default:
			utils.PrintErrorMsgExit("Not a recognized value for the upgrade action...\n", "")
		}
	case "dist-upgrade":
		switch p {
		case "":
			switch m {
			case "apt":
				debian.DistUpgrade()
			case "dnf":
				redhat.DistUpgrade()
			case "pacman":
				arch.DistUpgrade()
			case "yay":
				arch.YayDistUpgrade()
			case "paru":
				arch.ParuDistUpgrade()
			case "pkg":
				freebsd.DistUpgrade()
			case "zypper":
				suse.DistUpgrade()
			case "nix-env":
				utils.PrintErrorMsgExit("Error: It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
			case "flatpak":
				flatpak.DistUpgrade()
			case "snap":
				snap.DistUpgrade()
			case "brew":
				brew.DistUpgrade()
			case "go":
				golang.DistUpgrade()
			case "pip":
				pip.DistUpgrade()
			case "cargo":
				cargo.DistUpgrade()
			default:
				fmt.Println("Unsupported package manager...\n")
				os.Exit(1)
			}
		case "all":
			switch m {
			case "apt":
				fmt.Println("Updating APT repos:\n")
				debian.Update()
				fmt.Println("\nUpgrading with APT:\n")
				debian.DistUpgrade()
			case "dnf":
				fmt.Println("\nUpgrading with DNF:\n")
				redhat.DistUpgrade()
			case "pacman":
				fmt.Println("\nUpgrading with PACMAN:\n")
				arch.DistUpgrade()
			case "pkg":
				fmt.Println("\nUpgrading with PKG:\n")
				freebsd.DistUpgrade()
			case "zypper":
				fmt.Println("\nUpgrading with ZYPPER:\n")
				suse.DistUpgrade()
			case "nix-env":
				fmt.Println("Skipping nix-env: It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
			}
			npm := utils.GetNativePkgMgr()
			if npm == "pacman" && env.Yay != false {
				fmt.Println("\nUpgrade with YAY:\n")
				arch.YayDistUpgrade()
			}
			if npm == "pacman" && env.Paru != false {
				fmt.Println("\nUpgrade with PARU:\n")
				arch.ParuDistUpgrade()
			}
			if env.Flatpak != false {
				fmt.Println("\nUpgrading with FLATPAK:\n")
				flatpak.Upgrade()
			}
			if env.Snap != false {
				fmt.Println("\nUpgrading with SNAP:\n")
				snap.Upgrade()
			}
			if env.Brew != false {
				fmt.Println("\nUpgrading with HOMEBREW:\n")
				brew.Upgrade()
			}
			if env.Go != false {
				fmt.Println("\nUpgrading GO binaries:\n")
				golang.Upgrade()
			} else {
				fmt.Println("Go is disabled... Skipping...\n")
			}
			if env.Pip != false {
				fmt.Println("\nUpgrading with PIP installed packages:\n")
				pip.Upgrade()
			} else {
				fmt.Println("PIP is disabled... Skipping...\n")
			}
			if env.Cargo != false {
				fmt.Println("\nUpgrading with CARGO:\n")
				cargo.Upgrade()
			} else {
				fmt.Println("Cargo is disabled... Skipping...\n")
			}
		default:
			utils.PrintErrorMsgExit("Not a recognized value for the dist-upgrade action...\n", "")
		}
	case "autoremove":
		switch m {
		case "apt":
			debian.AutoRemove()
		case "dnf":
			redhat.AutoRemove()
		case "pacman":
			arch.AutoRemove()
		case "yay":
			arch.YayAutoRemove()
		case "paru":
			arch.ParuAutoRemove()
		case "pkg":
			freebsd.AutoRemove()
		case "zypper":
			suse.AutoRemove()
		case "nix-env":
			utils.PrintErrorMsgExit("Error: It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
		case "flatpak":
			flatpak.AutoRemove()
		case "snap":
			snap.AutoRemove()
		case "brew":
			brew.AutoRemove()
		case "go":
			golang.AutoRemove()
		case "pip":
			pip.AutoRemove()
		case "cargo":
			cargo.AutoRemove()
		default:
			fmt.Println("Unsupported package manager... Exiting...\n")
			os.Exit(1)
		}
	case "list":
		switch m {
		case "apt":
			if p == "" {
				debian.ListSystem()
			} else {
				debian.ListSystemSearch(p)
			}
		case "dnf":
			if p == "" {
				redhat.ListSystem()
			} else {
				redhat.ListSystemSearch(p)
			}
		case "pacman":
			if p == "" {
				arch.ListSystem()
			} else {
				arch.ListSystemSearch(p)
			}
		case "yay":
			if p == "" {
				arch.YayListSystem()
			} else {
				arch.YayListSystemSearch(p)
			}
		case "paru":
			if p == "" {
				arch.ParuListSystem()
			} else {
				arch.ParuListSystemSearch(p)
			}
		case "pkg":
			if p == "" {
				freebsd.ListSystem()
			} else {
				freebsd.ListSystemSearch(p)
			}
		case "zypper":
			if p == "" {
				suse.ListSystem()
			} else {
				suse.ListSystemSearch(p)
			}
		case "nix-env":
			utils.PrintErrorMsgExit("Error: It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
		case "flatpak":
			if p == "" {
				flatpak.ListSystem()
			} else {
				flatpak.ListSystemSearch(p)
			}
		case "snap":
			if p == "" {
				snap.ListSystem()
			} else {
				snap.ListSystemSearch(p)
			}
		case "brew":
			if p == "" {
				brew.ListSystem()
			} else {
				brew.ListSystemSearch(p)
			}
		case "go":
			if p == "" {
				golang.ListSystem()
			} else {
				golang.ListSystemSearch(p)
			}
		case "pip":
			if p == "" {
				pip.ListSystem()
			} else {
				pip.ListSystemSearch(p)
			}
		case "cargo":
			if p == "" {
				cargo.ListSystem()
			} else {
				cargo.ListSystemSearch(p)
			}
		default:
			fmt.Println("Unsupported package manager... Exiting...\n")
			os.Exit(1)
		}
	case "history":
		switch m {
		case "apt", "dnf", "pacman", "yay", "paru", "pkg", "zypper", "flatpak", "snap", "brew", "go", "pip", "cargo", "appimage":
			utils.History(m, p, sort)
		case "nix-env":
			utils.PrintErrorMsgExit("Error: It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
		default:
			fmt.Println("Unsupported package manager... Exiting...\n")
			os.Exit(1)
		}
	case "enable":
		switch m {
		case "apt":
			env.Apt = true
			utils.EditSettings("APT = ", "y")
			fmt.Println("Apt has been enabled...\n")
		case "dnf":
			env.Dnf = true
			utils.EditSettings("DNF = ", "y")
			fmt.Println("Dnf has been enabled...\n")
		case "pacman":
			env.Yay = true
			utils.EditSettings("PACMAN = ", "y")
			fmt.Println("Pacman has been enabled...\n")
		case "zypper":
			env.Yay = true
			utils.EditSettings("ZYPPER = ", "y")
			fmt.Println("Zypper has been enabled...\n")
		case "yay":
			env.Yay = true
			utils.EditSettings("YAY = ", "y")
			fmt.Println("Yay has been enabled...\n")
		case "paru":
			env.Paru = true
			utils.EditSettings("PARU = ", "y")
			fmt.Println("Paru has been enabled...\n")
		case "nix-env":
			utils.PrintErrorMsgExit("Error: It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
		case "flatpak":
			env.Flatpak = true
			utils.EditSettings("FLATPAK = ", "y")
			fmt.Println("Flatpak has been enabled...\n")
		case "snap":
			env.Snap = true
			utils.EditSettings("SNAP = ", "y")
			fmt.Println("Snap has been enabled...\n")
		case "brew":
			if env.OSType == "Mac" {
				utils.PrintErrorMsgExit("", "Homebrew is the default for macOS. It's already enabled...")
			}
			env.Brew = true
			utils.EditSettings("BREW = ", "y")
			fmt.Println("Homebrew has been enabled...\n")
		case "go":
			env.Go = true
			utils.EditSettings("GOLANG = ", "y")
			fmt.Println("Go has been enabled...\n")
		case "pip":
			env.Pip = true
			utils.EditSettings("PIP = ", "y")
			fmt.Println("Pip has been enabled...\n")
		case "cargo":
			env.Cargo = true
			utils.EditSettings("CARGO = ", "y")
			fmt.Println("Cargo has been enabled...\n")
		case "appimage":
			env.AppImage = true
			utils.EditSettings("APPIMAGE = ", "y")
			fmt.Println("AppImage has been enabled...\n")
		default:
			fmt.Println("Can't enable... Unsupported package manager...\n")
			os.Exit(1)
		}
	case "disable":
		switch m {
		case "apt":
			env.Apt = false
			utils.EditSettings("APT = ", "n")
			fmt.Println("Apt has been disabled...\n")
		case "dnf":
			env.Dnf = false
			utils.EditSettings("DNF = ", "n")
			fmt.Println("Dnf has been disabled...\n")
		case "pacman":
			env.Yay = false
			utils.EditSettings("PACMAN = ", "n")
			fmt.Println("Pacman has been disabled...\n")
		case "zypper":
			env.Yay = false
			utils.EditSettings("ZYPPER = ", "n")
			fmt.Println("Zypper has been disabled...\n")
		case "yay":
			env.Yay = false
			utils.EditSettings("YAY = ", "n")
			fmt.Println("Yay has been disabled...\n")
		case "paru":
			env.Paru = false
			utils.EditSettings("PARU = ", "n")
			fmt.Println("Paru has been disabled...\n")
		case "nix-env":
			utils.PrintErrorMsgExit("Error: It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
		case "flatpak":
			env.Flatpak = false
			utils.EditSettings("FLATPAK = ", "n")
			fmt.Println("Flatpak has been disabled...\n")
		case "snap":
			env.Snap = false
			utils.EditSettings("SNAP = ", "n")
			fmt.Println("Snap has been disabled...\n")
		case "brew":
			if env.OSType == "Mac" {
				utils.PrintErrorMsgExit("Error:", "Homebrew is the default for macOS. It can't be disabled...")
			}
			env.Brew = false
			utils.EditSettings("BREW = ", "n")
			fmt.Println("Homebrew has been disabled...\n")
		case "go":
			env.Go = false
			utils.EditSettings("GOLANG = ", "n")
			fmt.Println("Go has been disabled...\n")
		case "pip":
			env.Pip = false
			utils.EditSettings("PIP = ", "n")
			fmt.Println("Pip has been disabled...\n")
		case "cargo":
			env.Cargo = false
			utils.EditSettings("CARGO = ", "n")
			fmt.Println("Cargo has been disabled...\n")
		case "appimage":
			env.AppImage = false
			utils.EditSettings("APPIMAGE = ", "n")
			fmt.Println("AppImage has been disabled...\n")
		default:
			fmt.Println("Can't disable... Unsupported package manager...\n")
			os.Exit(1)
		}
	case "add-repo":
		switch m {
		case "apt":
			debian.AddRepo(p, g)
		case "dnf":
			redhat.AddRepo(p, g)
		case "pacman":
			arch.AddRepo(p, g)
		case "yay":
			arch.YayAddRepo(p, g)
		case "paru":
			arch.ParuAddRepo(p, g)
		case "pkg":
			freebsd.AddRepo(p, g)
		case "zypper":
			suse.AddRepo(p, g)
		case "nix-env":
			utils.PrintErrorMsgExit("Error: It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
		case "flatpak":
			flatpak.AddRepo(p, g)
		default:
			fmt.Println("Unsupported pacakge manager...\n")
			os.Exit(1)
		}
	case "rm-repo":
		switch m {
		case "apt":
			debian.RemoveRepo(p)
		case "dnf":
			redhat.RemoveRepo(p)
		case "pacman":
			arch.RemoveRepo(p)
		case "yay":
			arch.YayRemoveRepo(p)
		case "paru":
			arch.ParuRemoveRepo(p)
		case "pkg":
			freebsd.RemoveRepo(p)
		case "zypper":
			suse.RemoveRepo(p)
		case "nix-env":
			utils.PrintErrorMsgExit("Error: It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
		case "flatpak":
			flatpak.RemoveRepo(p)
		default:
			fmt.Println("Unsupported pacakge manager...\n")
			os.Exit(1)
		}
	case "ls-repo":
		switch m {
		case "apt", "dnf", "pacman", "yay", "paru", "zypper", "flatpak", "snap":
			utils.ListRepo(m, p)
		case "nix-env":
			utils.PrintErrorMsgExit("Error: It's redundant to use app for nix-env packages when you can just install packages by adding packages to /etc/nixos/configurations.nix. app is only meant for managing other package managers that are used on nixos...", "")
		default:
			fmt.Println("Unsupported pacakge manager...\n")
			os.Exit(1)
		}
	case "settings":
		fmt.Println("Package Managers:\n")
		if env.OSType == "Linux" {
			npm := utils.GetNativePkgMgr()
			fmt.Print(npm + ": ")
			nabled := false
			switch npm {
			case "apt":
				if env.Apt == true {
					nabled = true
				}
			case "dnf":
				if env.Dnf == true {
					nabled = true
				}
			case "pacman":
				if env.Pacman == true {
					nabled = true
				}
			case "zypper":
				if env.Zypper == true {
					nabled = true
				}
			}
			if nabled == true {
				fmt.Println(utils.ColorGreen, "Enabled", utils.ColorReset)
			} else {
				fmt.Println(utils.ColorRed, "Disabled", utils.ColorReset)
			}
			fmt.Print("yay: ")
			if env.Yay == true {
				fmt.Println(utils.ColorGreen, "Enabled", utils.ColorReset)
			} else {
				fmt.Println(utils.ColorRed, "Disabled", utils.ColorReset)
			}
			fmt.Print("paru: ")
			if env.Paru == true {
				fmt.Println(utils.ColorGreen, "Enabled", utils.ColorReset)
			} else {
				fmt.Println(utils.ColorRed, "Disabled", utils.ColorReset)
			}
			fmt.Print("flatpak: ")
			if env.Flatpak == true {
				fmt.Println(utils.ColorGreen, "Enabled", utils.ColorReset)
			} else {
				fmt.Println(utils.ColorRed, "Disabled", utils.ColorReset)
			}
			fmt.Print("snap: ")
			if env.Snap == true {
				fmt.Println(utils.ColorGreen, "Enabled", utils.ColorReset)
			} else {
				fmt.Println(utils.ColorRed, "Disabled", utils.ColorReset)
			}
			fmt.Print("brew: ")
			if env.Brew == true {
				fmt.Println(utils.ColorGreen, "Enabled", utils.ColorReset)
			} else {
				fmt.Println(utils.ColorRed, "Disabled", utils.ColorReset)
			}
			fmt.Print("appimage: ")
			if env.AppImage == true {
				fmt.Println(utils.ColorGreen, "Enabled", utils.ColorReset)
			} else {
				fmt.Println(utils.ColorRed, "Disabled", utils.ColorReset)
			}
			fmt.Print("go: ")
			if env.Go == true {
				fmt.Println(utils.ColorGreen, "Enabled", utils.ColorReset)
			} else {
				fmt.Println(utils.ColorRed, "Disabled", utils.ColorReset)
			}
			fmt.Print("pip: ")
			if env.Pip == true {
				fmt.Println(utils.ColorGreen, "Enabled", utils.ColorReset)
			} else {
				fmt.Println(utils.ColorRed, "Disabled", utils.ColorReset)
			}
			fmt.Print("cargo: ")
			if env.Cargo == true {
				fmt.Println(utils.ColorGreen, "Enabled", utils.ColorReset)
			} else {
				fmt.Println(utils.ColorRed, "Disabled", utils.ColorReset)
			}
		}

		if env.OSType == "FreeBSD" {
			fmt.Print("go: ")
			if env.Go == true {
				fmt.Println(utils.ColorGreen, "Enabled", utils.ColorReset)
			} else {
				fmt.Println(utils.ColorRed, "Disabled", utils.ColorReset)
			}
			fmt.Print("pip: ")
			if env.Pip == true {
				fmt.Println(utils.ColorGreen, "Enabled", utils.ColorReset)
			} else {
				fmt.Println(utils.ColorRed, "Disabled", utils.ColorReset)
			}
			fmt.Print("cargo: ")
			if env.Cargo == true {
				fmt.Println(utils.ColorGreen, "Enabled", utils.ColorReset)
			} else {
				fmt.Println(utils.ColorRed, "Disabled", utils.ColorReset)
			}
		}

		if env.OSType == "Mac" {
			fmt.Print("brew: ")
			if env.Brew == true {
				fmt.Println(utils.ColorGreen, "Enabled", utils.ColorReset)
			} else {
				fmt.Println(utils.ColorRed, "Disabled", utils.ColorReset)
			}
			fmt.Print("go: ")
			if env.Go == true {
				fmt.Println(utils.ColorGreen, "Enabled", utils.ColorReset)
			} else {
				fmt.Println(utils.ColorRed, "Disabled", utils.ColorReset)
			}
			fmt.Print("pip: ")
			if env.Pip == true {
				fmt.Println(utils.ColorGreen, "Enabled", utils.ColorReset)
			} else {
				fmt.Println(utils.ColorRed, "Disabled", utils.ColorReset)
			}
			fmt.Print("cargo: ")
			if env.Cargo == true {
				fmt.Println(utils.ColorGreen, "Enabled", utils.ColorReset)
			} else {
				fmt.Println(utils.ColorRed, "Disabled", utils.ColorReset)
			}
		}
		fmt.Println("")
	default:
		fmt.Println("Unsupported action... Exiting...\n")
		os.Exit(1)
	}

}
