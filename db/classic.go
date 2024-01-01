package db

import (
	"github.com/hkdb/app/env"

	"os"
)

func RecordClassic(pm, pkg string, classic bool) error {

	// Write new entry
	werr := writeClassicEntry(pm, pkg, classic)
	if werr != nil {
		return werr
	}
	
	return nil

}

func GetClassic(pm, name string) (bool, error) {

	classic, err := readClassicEntry(pm, name)
	if err != nil {
		return false, err
	}

	return classic.Classic, nil

}

func RemoveClassic(pm, p string) error {

	cRec := env.DBDir + "/packages/repo/local/" + pm + "/" + p + ".json"
	if _, err := os.Stat(cRec); err == nil {
		if err := os.Remove(cRec); err != nil {
			return err
		}
	}

	return nil

}

func writeClassicEntry(pm, n string, classic bool) error {

	// Init DB to write
	pdb, err := initDB("packages/repo/local")
	if err != nil {
		return err
	}

	// Concat new string
	c := Classics{}
	c.Classic = classic

	// Write to DB
	if err := pdb.Write(pm, n, c); err != nil {
		return err
	}	

	return nil

}

func readClassicEntry(pm, name string) (Classics, error) {

	classic := Classics{}

	pdb, err := initDB("packages/repo/local")
	if err != nil {
		return classic, err
	}
	if err := pdb.Read(pm, name, &classic); err != nil {
		return classic, err		
	}

	return classic, nil

}

