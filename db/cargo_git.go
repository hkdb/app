package db

import (
	"github.com/hkdb/app/env"

	"os"
	"strings"

	"golang.org/x/exp/slices"
)

func GitExists(pm, r string) (bool, error) {

	aRepos := strings.Fields(r)
	repos, err := ReadPkgSlice("packages", "repo", pm)
	if err != nil {
		return false, err
	}

	for i := 0; i < len(aRepos); i++ {
		installed := slices.Contains(repos, aRepos[i])
		if installed == true {
			if len(aRepos[i]) < 4 {
				return false, nil
			}
			if aRepos[i][:4] == "http" { 
				return true, nil
			}
		}
	}

	return false, nil

}

func RecordGit(pm, r string) error {

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

func RecordGitSetup(pm, name, u, t string) error {

	// Write new entry
	werr := writeGitEntry(pm, name, u, t)
	if werr != nil {
		return werr
	}

	return nil

}

func GetGitSetup(pm, name string) (string, string, error) {

	setup, err := readGitEntry(pm, name)
	if err != nil {
		return "", "", err
	}

	return setup.Url, setup.Tag, nil

}

func RemoveGit(pm, p string) error {

	rmRepos := strings.Split(p, " ")

	repos, err := ReadPkgSlice("packages", "repo", pm)
	if err != nil {
		return err
	}

	// Remove matching elements from slice
	for i := 0; i < len(rmRepos); i++ {
		for x := 0; x < len(repos); x++ {
			if rmRepos[i] == repos[x] {
				// Remove package git data
				urlRec := env.DBDir + "/packages/repo/local/" + pm + "/" + p + ".json"
				if _, err = os.Stat(urlRec); err == nil {
					if err := os.Remove(urlRec); err != nil {
						return err
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

func writeGitEntry(pm, n, u, t string) error {

	// Init DB to write
	pdb, err := initDB("packages/repo/local")
	if err != nil {
		return err
	}

	// Concat new string
	repo := Git{}
	repo.Url = u
	repo.Tag = t

	// Write to DB
	if err := pdb.Write(pm, n, repo); err != nil {
		return err
	}

	return nil

}

func readGitEntry(pm, name string) (Git, error) {

	repo := Git{}

	pdb, err := initDB("packages/repo/local")
	if err != nil {
		return repo, err
	}
	if err := pdb.Read(pm, name, &repo); err != nil {
		return repo, err
	}

	return repo, nil

}
