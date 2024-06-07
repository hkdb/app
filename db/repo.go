package db

import (
	"github.com/hkdb/app/env"

	"os"
	"strings"

	"golang.org/x/exp/slices"
)

func RepoExists(pm, r string) (bool, error) {

	aRepos := strings.Split(r, " ")
	repos, err := ReadPkgSlice("packages", "repo", pm)
	if err != nil {
		return false, err
	}

	for i := 0; i < len(aRepos); i++ {
		installed := slices.Contains(repos, aRepos[i])
		if installed == true {
			return true, nil
		}
	}

	return false, nil

}

func RecordRepo(pm, r string) error {

	current := ""

	if _, err := os.Stat(env.DBDir + "/packages/repo/" + pm + ".json"); err == nil {
		// Read current package string from package manager
		current, err = ReadPkgs("packages", "repo", pm)
		if err != nil {
			return err
		}

		// Delete current Package entry for package manager
		derr := delPkgsEntry("packages", "repo", pm)
		if derr != nil {
			return derr
		}
	}

	repos := current + " " + r

	// Write new entry
	werr := writePkgsEntry("packages", "repo", pm, repos)
	if werr != nil {
		return werr
	}

	return nil

}

func RecordSetup(pm, name, r, g string) error {

	// Write new entry
	werr := writeRepoEntry(pm, name, r, g)
	if werr != nil {
		return werr
	}

	return nil

}

func GetSetup(pm, name string) (string, string, error) {

	setup, err := readRepoEntry(pm, name)
	if err != nil {
		return "", "", err
	}

	return setup.Repo, setup.Gpg, nil

}

func RemoveRepo(pm, rType, r string) error {

	rmRepos := strings.Split(r, " ")

	repos, err := ReadPkgSlice("packages", "repo", pm)
	if err != nil {
		return err
	}

	// Remove matching elements from slice
	for i := 0; i < len(rmRepos); i++ {
		for x := 0; x < len(repos); x++ {
			if rmRepos[i] == repos[x] {
				if rType != "ppa" {
					// Remove script
					urlRec := env.DBDir + "/packages/repo/local/" + pm + "/" + r + ".json"
					if _, err = os.Stat(urlRec); err == nil {
						if err := os.Remove(urlRec); err != nil {
							return err
						}
					}
					script := env.DBDir + "/packages/repo/local/" + pm + "/" + r + ".sh"
					if _, err = os.Stat(script); err == nil {
						if err := os.Remove(script); err != nil {
							return err
						}
					}
				}
				repos = removeIndex(repos, x)
			}
		}
	}

	// Join the new slice into string
	newRepos := strings.Join(repos, " ")

	// Delete current entry
	derr := delPkgsEntry("packages", "repo", pm)
	if derr != nil {
		return derr
	}

	// Write new entry
	writePkgsEntry("packages", "repo", pm, newRepos)

	return nil

}

func writeRepoEntry(pm, n, r, g string) error {

	// Init DB to write
	pdb, err := initDB("packages/repo/local")
	if err != nil {
		return err
	}

	// Concat new string
	repo := Repos{}
	repo.Repo = r
	repo.Gpg = g

	// Write to DB
	if err := pdb.Write(pm, n, repo); err != nil {
		return err
	}

	return nil

}

func readRepoEntry(pm, name string) (Repos, error) {

	repo := Repos{}

	pdb, err := initDB("packages/repo/local")
	if err != nil {
		return repo, err
	}
	if err := pdb.Read(pm, name, &repo); err != nil {
		return repo, err
	}

	return repo, nil

}
