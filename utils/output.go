package utils

import (
	"fmt"
	"log"
	"os"
)

// Toggle on and off timestamp in logging
func TimeStamp(ts bool) {

	if ts == true {
		log.SetFlags(3)
	}
	if ts == false {
		log.SetFlags(0)
	}

}

// Insert newline in log
func NewLine() {

	log.SetFlags(0)
	log.Println("")
	log.SetFlags(3)

}

// Display header banner
func LogLaunchBanner(version1 string) {

	TimeStamp(false)
	logo := `
_____  ______ ______  
\__  \ \____ \\____ \ 
 / __ \|  |_> >  |_> >
(____  /   __/|   __/ 
     \/|__|   |__|    
`
	log.Println(string(ColorCyan), logo, string(ColorReset))
	log.Println(string(ColorGreen), "ϟ app - a package management assistant with super powers", string(ColorReset), version1) // can add + " & " + version2
	TimeStamp(true)

	NewLine()

}

// Print error and exit
func PrintErrorExit(eType string, e error) {

	NewLine()
	fmt.Println(ColorRed, eType, ColorReset, e)
	NewLine()
	os.Exit(1)

}

// Print error
func PrintError(eType string, e error) {

	NewLine()
	fmt.Println(ColorRed, eType, ColorReset, e)
	NewLine()

}

// Print error msg and exit
func PrintErrorMsgExit(eType string, e string) {

	NewLine()
	fmt.Println(ColorRed, eType, ColorReset, e)
	NewLine()
	os.Exit(1)

}

// Print error msg and exit
func PrintErrorMsg(eType string, e string) {

	NewLine()
	fmt.Println(ColorRed, eType, ColorReset, e)
	NewLine()

}
