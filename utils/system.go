package utils

import (
	"github.com/hkdb/app/env"
)

func GetNativePkgMgr() string {
	switch env.OSType {
	case "Linux":
		switch env.Base {
		case "debian":
			return "apt"
		case "redhat":
			return "dnf"
		case "arch":
			return "pacman"
		case "suse":
			return "zypper"
		case "nixos":
			return "nix-env"
		default:
			PrintErrorMsgExit("Get Native PM Error: Distro not supported...", "")
		}
	case "macos":
		PrintErrorMsgExit("macos is not implemented yet...", "")
	case "windows":
		PrintErrorMsgExit("windows is not implemented yet...", "")
	default:
		PrintErrorMsgExit("Operating system not supported...", "")
	}

	return ""

}

func IsNativeEnabled() bool {
	pm := GetNativePkgMgr()

	switch pm {
	case "apt":
		return env.Apt
	case "dnf":
		return env.Dnf
	case "pacman":
		return env.Pacman
	case "zypper":
		return env.Zypper
	case "nix-env":
		return env.NixEnv
	default:
		PrintErrorMsgExit("Native PM Check Error:", "This native package amanger is not supported...")
	}

	return false
}
