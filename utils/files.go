package utils

import (
	"github.com/hkdb/app/db"
	"github.com/hkdb/app/env"

	"os"
	"os/exec"
	"strings"
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
	"bufio"
)

func IsLocalInstall(filename string) (bool, string, string) {

	// Get the extension of the file
	ext := GetFileExtension(filename)
	if ext == "NONE" {
      return false, "", ""  //if no extension is present print failure
  }

	// Check extension against current environment
	switch env.Base {
	case "debian":
		if ext != "deb" {
			return false, "", ""
		}
	case "redhat":
		if ext != "rpm" {
			PrintErrorMsgExit("Error:", "This is not a supported file type.")
		}
	default:
		PrintErrorMsgExit("File Error:", "This is not a supported file type...")
	}

	// Get working path
	path := GetWorkPath() 
  // fmt.Println("Working Path:", path)

	// Check if config dir for local package exists. If it doesn't, mkdir
	dir := env.DBDir + "/packages/local/" + ext
	if _, derr := os.Stat(dir); os.IsNotExist(derr) {
    merr := os.MkdirAll(dir, os.ModePerm)
		if merr != nil {
			PrintErrorExit("Error:", merr)
		}
	}

  // fmt.Println("DB Dir:", dir)

	// Set source package file
	srcpkg := path + "/" + filename
  destpkg := dir + "/" + filename

  //fmt.Println("Copying package from " + srcpkg + " to " + destpkg + "...")

	// Copy local package to config dir
	cerr := Copy(srcpkg, destpkg)
	if cerr != nil {
		PrintErrorExit("Copy Error:", cerr)
	}

	// Get package name
	var pNameCmd []byte
	var perr error

	switch ext {
	case "deb":
		pNameCmd, perr = exec.Command("dpkg-deb", "-f", filename, "Package").Output()
	case "rpm":
		pNameCmd, perr = exec.Command("rpm", "-q", "--qf", "\"%{NAME}\n\"", "-p", filename).Output()
	default:
		PrintErrorMsgExit("File Error:", "This is not a supported file type...")
	}
	if perr != nil {
		PrintErrorExit("Read Package Name Error:", perr)
	}

	pname := strings.TrimSuffix(string(pNameCmd), "\n")

  // fmt.Println("Package Name:", pname)

	// Mark package as local
	rerr := db.RecordPkg("", "packages" , ext, pname)
	if rerr != nil {
		PrintErrorExit("Record Local Package Error: ", rerr)
	}
	
	// Record package and file name association
	rferr := db.RecordPkg("packages/local", ext, pname, filename)
	if rferr != nil {
		PrintErrorExit("Record Local Package Association Error: ", rferr)
	}

	return true, pname, srcpkg
	
}

func Copy(src, dst string) error {
	
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		regular := errors.New("This is not a regular file...")
		return regular
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, cerr := io.Copy(destination, source)
	return cerr

}

func CreateDirIfNotExist(dir string) {

	if _, err := os.Stat(dir); os.IsNotExist(err) {
    err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			PrintErrorExit("Error:", err)
		}
	}

}

func DeleteDirIfEmpty(dir string) {
	
	// Clean-up dir
	if empty, _ := DirIsEmpty(dir); empty == true {
		if err := os.Remove(dir); err != nil {
			PrintErrorExit("Remove Error:", err)
		}
	}

}

func DirIsEmpty(dir string) (bool, error) {
  
	f, err := os.Open(dir)
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

func GetWorkPath() string {

	// Get working path
	path, err := os.Getwd()
	if err != nil {
		PrintErrorExit("Path Read Error:", err)
	}
	
	return path

}

func GetFileExtension(filename string) string {

	// Get the extension of the file
	extIndex := strings.LastIndex(filename, ".")
	if extIndex == -1 {
      return "NONE" 
  }

	ext := filename[extIndex+1:]
	
	return ext

}


func GetFileName(filename string) string {

	var ext = filepath.Ext(filename)
	var name = filename[0:len(filename)-len(ext)]

	return name

}

func StripExtension(filename, ext string) string {

	name := strings.TrimSuffix(filename, ext)

	return name

}

func WriteToFile(line string, file string) error {

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(line + "\n"); err != nil {
		return err
	}
	return nil

}

func AppendToFile(line string, file string) error {

	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(line + "\n"); err != nil {
		return err
	}
	return nil

}

func WriteDotDesktop(pkg, icon, file, dest, share string) {

	input, err := ioutil.ReadFile(file)
	if err != nil {
		PrintErrorExit("Read .desktop Error:", err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "Icon=") {
			lines[i] = "Icon=" + icon
		}
		if strings.Contains(line, "Exec=") {
			lines[i] = "Exec=" + pkg
		}
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(dest, []byte(output), 0644)
	if err != nil {
		PrintErrorExit("Write .desktop Error:", err)
	}

	err = ioutil.WriteFile(share, []byte(output), 0644)
	if err != nil {
		PrintErrorExit("Write .desktop Error:", err)
	}

}

func EditSettings(lType, value string) {

	input, err := ioutil.ReadFile(env.DBDir + "/settings.conf" )
	if err != nil {
		PrintErrorExit("Read settings Error:", err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, lType) {
			lines[i] = lType + value
		}
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(env.DBDir + "/settings.conf", []byte(output), 0644)
	if err != nil {
		PrintErrorExit("Write settings Error:", err)
	}

}

func ReadDotDesktop(file string) (*DotDesktop, error) {

  dd := &DotDesktop{}
	f, err := os.Open(file)
	if err != nil {
		return dd, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	
	for scanner.Scan() {
    line := scanner.Text()
    entry := strings.Split(line, "=")
    if len(entry) > 2 {
      PrintErrorMsgExit("Read .desktop Error:", "Malformed .desktop file...")
    }
    if entry[0] == "Name" {
      dd.Name = entry[1]
      if dd.Icon != "" {
        break
      }
    }
    if entry[0] == "Icon" {
      dd.Icon = entry[1]
      if dd.Name != "" {
        break
      }
    }
	}

	if err := scanner.Err(); err != nil {
		return dd, err
	}

	return dd, nil

}

func ReadFileLine(file string, line int) (string, error) {
	
	line = line-1 

	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var current int
	for scanner.Scan() {
		if current == line {
			return scanner.Text(), nil
		}
		current++
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil

}

func RemoveFileLine(file string, ln int) {

	ln = ln - 1

	input, err := ioutil.ReadFile(file)
	if err != nil {
		PrintErrorExit("Read .desktop Error:", err)
	}

	lines := strings.Split(string(input), "\n")

	for i := 0; i < len(lines); i++ {
		if i == ln {
			lines = append(lines[:ln], lines[ln+1:]...)
			break
		}
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(file, []byte(output), 0644)
	if err != nil {
		PrintErrorExit("Write .desktop Error:", err)
	}

}

func RemoveDirRecursive(dir string) {

	if err := os.RemoveAll(dir); err != nil {
		PrintErrorExit("Remove Dir Error:", err)
	}

}


func RemoveFile(file string) {

	if err := os.Remove(file); err != nil {
		PrintErrorExit("Remove File Error:", err)
	}

}

func GetFileFromFullPath(file string) string {

	return filepath.Base(file)

}

func CheckIfExists(fileOrDir string) (bool, error) {

	if _, err := os.Stat(fileOrDir); os.IsNotExist(err) {
		return false, err
	}

	return true, nil

}
