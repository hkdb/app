package db

import (
	"github.com/hkdb/app/env"

	"fmt"
	"strings"
	"os"
  "io"

	"golang.org/x/exp/slices"
	scribble "github.com/nanobox-io/golang-scribble"
)
	
func IsInstalled(root, col, pm, p string) (bool, error) {

	iPkgs := strings.Split(p, " ")
	pkgs, err := ReadPkgSlice(root, col, pm)
	if err != nil {
		return false, err
	}

	for i := 0; i < len(iPkgs); i++ {
		installed := slices.Contains(pkgs, iPkgs[i])
		if installed == true {
			return true, nil
		}
	}

	return false, nil 

}

// Add newly isntalled packages to the end of the string of packages
func RecordPkg(root, col, pm, p string) error {
	
  current := ""

	if _, err := os.Stat(env.DBDir + "/" + col + "/" + pm + ".json"); err == nil {
		// Read current package string from package manager
		current, err = ReadPkgs(root, col, pm)
		if err != nil {
			return err
		}

		// Delete current Package entry for package manager
		derr := delPkgsEntry(root, col, pm)
		if derr != nil {
			return derr
		}
	}

	
	// Assemble new string	
	pkgs := current + " " + p

	// Write new entry
	werr := writePkgsEntry(root, col, pm, pkgs)
	if werr != nil {
		return werr
	}
	
	return nil

}

func RemovePkg(root, col, pm, p string) error {
	
	rmPkgs := strings.Split(p, " ")
  
	pkgs, err := ReadPkgSlice(root, col, pm)
	if err != nil {
		return err
	}

	// Remove matching elements from slice
	for i := 0; i < len(rmPkgs); i++ {
		for x := 0; x < len(pkgs); x++ {
      if rmPkgs[i] == pkgs[x] {
				pkgs = removeIndex(pkgs, x)
			}
		}	
	}

	// Join the new slice into string
	newPkgs := strings.Join(pkgs, " ")

	// Delete current entry
	derr := delPkgsEntry(root, col, pm)
	if derr != nil {
		return derr
	}

	// Write new entry
	writePkgsEntry(root, col, pm, newPkgs)

  // Check to see if there are associated local packages and if so, delete
  rlerr := removeLocalPkgs(root, col, pm, rmPkgs)	
  if rlerr != nil {
    return rlerr
  }

	return nil

}

func removeLocalPkgs(root string, col string, pm string, rmPkgs []string) error {

	var file = ""

	switch pm {
	case "apt":
		file = "deb"
	case "dnf":
		file = "rpm"
  default:
    return nil
	}

	// ie. ~/.config/app/packages/deb.json
  if _, oerr := os.Stat(env.DBDir + "/packages/" + file + ".json"); oerr != nil {
    return nil
  } 

	if file != "" {
		lpkgs, lerr := ReadPkgSlice("", col, file)
		if lerr != nil {
			return lerr
		}

		// Remove matching elements from local slice
		for j := 0; j < len(rmPkgs); j++ {
			for y := 0; y < len(lpkgs); y++ {
				if rmPkgs[j] == lpkgs[y] {
          fmt.Println(" Removing " + rmPkgs[j] + " from " + file + " file history...\n")
					// Remove package from local 
          lpkgs = removeIndex(lpkgs, y)
          
          rmFileRecord, rerr := ReadPkgs("packages/local", file, rmPkgs[j])
          if rerr == nil {
	          rmFile := strings.TrimSpace(rmFileRecord)

            // Remove package file
            fmt.Println(" Removing:", rmFile + "\n")
            reerr := os.Remove(env.DBDir + "/packages/local/" + file + "/" + rmFile)
            if reerr != nil {
              return reerr
            }
            // Remove package.json
            derr := os.Remove(env.DBDir + "/packages/local/" + file + "/" + rmPkgs[j] + ".json")
            if derr != nil {
              return derr
            }
          }
				}
			}	
		}

    // Clean up empty folders in local
    ldir := env.DBDir + "/packages/local/" + file
    empty, _ := dirIsEmpty(ldir)
    if empty == true {
      rmerr := os.Remove(ldir)
      if rmerr != nil {
        return rmerr
      }
    }

		// Join the new local slice into string
		newLPkgs := strings.Join(lpkgs, " ")

		// Delete current entry
		derr := delPkgsEntry(root, col, file)
		if derr != nil {
			return derr
		}

		// Write new local entry
		writePkgsEntry(root, col, file, newLPkgs)

	}

  return nil

}

func ReadPkgs(root, col, pm string) (string, error) {

	pdb, err := initDB(root)
	if err != nil {
		return "", err
	}
	pkgs := Packages{}
	if err := pdb.Read(col, pm, &pkgs); err != nil {
		return "", err		
	}

	return pkgs.Packages, nil

}

func ReadPkgSlice(root, col, pm string) ([]string, error) {
	
	pkgs, err := ReadPkgs(root, col, pm)
	if err != nil {
		errStr := err.Error()
		errMsg := errStr[len(errStr)-25:]
		if errMsg == "no such file or directory" {
			return []string{""}, nil
		}
		return []string{""}, err
	}

	pkgSlice := strings.Split(pkgs, " ")

	return pkgSlice, nil

}

func delPkg(root, col, pm string) error {
	
	pdb, err := initDB(root)
	if err != nil {
		return err
	}
	if err := pdb.Delete(col, pm); err != nil {
		return err
	}

	return nil
}

func delPkgsEntry(root, col, pm string) error {

	pdb, err := initDB(root)
	if err != nil {
		return err
	}
	if err := pdb.Delete(col , pm); err != nil {
		return err
	}

	return nil

}

func writePkgsEntry(root, col, pm, p string) error {

	// Init DB to write
	pdb, err := initDB(root)
	if err != nil {
		return err
	}

	// Concat new string
	pkgs := Packages{}
	pkgs.Packages = p 

	// Write to DB
	if err := pdb.Write(col, pm, pkgs); err != nil {
		return err
	}	

	return nil

}

func initDB(root string) (*scribble.Driver, error) {

	pdb, err := scribble.New(env.DBDir + "/" + root, nil)
	if err != nil {
		return pdb, err
	}

	return pdb, nil

}

func removeIndex(s []string, index int) []string {

	return append(s[:index], s[index+1:]...)

}

func dirIsEmpty(name string) (bool, error) {
  f, err := os.Open(name)
  if err != nil {
    return false, err
  }
  defer f.Close()

  _, err = f.Readdirnames(1) // Or f.Readdir(1)
  if err == io.EOF {
    return true, nil
  }
  return false, err
}
