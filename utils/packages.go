package utils

import (
	"github.com/hkdb/app/db"

	"fmt"
	"strings"
)

func History(pm, p string) {

	pkgs, err := db.ReadPkgSlice("", "packages", pm)
	if err != nil {
		PrintErrorExit("Read History Error:", err)
	}

	if p != "" {
		fmt.Println("Packages manually installed by app that matches \"" + p + "\":\n")
		if len(pkgs) == 1 && pkgs[0] == "" {
			fmt.Println(p + " was not installed with app...")
		} else {
			for i := 0; i < len(pkgs); i++ {
				pkg := strings.ToLower(pkgs[i])
				search := strings.ToLower(p)
				if strings.Contains(pkg, search) {
					fmt.Println("\n" + pkgs[i])
				}
			}
		}
	} else {
		fmt.Println("Packages manually installed by app:\n")
		if len(pkgs) == 1 && pkgs[0] == "" {
			fmt.Println("No packages have been installed with app yet...")
		} else {
			for i := 0; i < len(pkgs); i++ {
				fmt.Println(pkgs[i])
			}
		}
	}

	fmt.Println("\n")

}
