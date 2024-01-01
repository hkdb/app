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
		default:
			PrintErrorMsgExit("Error: Distro not supported...", "")
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
