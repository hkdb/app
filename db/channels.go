package db

import (
	"github.com/hkdb/app/env"

	"os"
)

func ChannelPreferred(pm, p string) bool {

	if _, err := os.Stat(env.DBDir + "/packages/repo/channel/" + pm + "/" + p + ".json"); err == nil {
		return true
	}
	return false 

}

func RecordChan(pm, name, c string) error {

	// Write new entry
	werr := writeChanEntry(pm, name, c)
	if werr != nil {
		return werr
	}
	
	return nil

}

func GetChan(pm, name string) (string, error) {

	c, err := readChanEntry(pm, name)
	if err != nil {
		return "", err
	}

	return c.Channel, nil

}

func RemoveChan(pm, name string) error {

	err := deleteChanFile(pm, name)
	if err != nil {
		return err
	}
	
	return nil

}

func writeChanEntry(pm, n, c string) error {

	// Init DB to write
	pdb, err := initDB("packages/repo/channel")
	if err != nil {
		return err
	}

	// Concat new string
	ch := Channels{}
	ch.Channel = c

	// Write to DB
	if err := pdb.Write(pm, n, ch); err != nil {
		return err
	}	

	return nil

}

func readChanEntry(pm, name string) (Channels, error) {

	ch := Channels{}

	pdb, err := initDB("packages/repo/channel")
	if err != nil {
		return ch, err
	}
	if err := pdb.Read(pm, name, &ch); err != nil {
		return ch, err		
	}

	return ch, nil

}

func deleteChanFile(pm, name string) error {

	file := env.DBDir + "/packages/repo/channel/" + pm + "/" + name + ".json"
	if _, err := os.Stat(file); err == nil {
		derr := os.Remove(file)
		if derr != nil {
			return derr
		}
	}

	return nil

}
