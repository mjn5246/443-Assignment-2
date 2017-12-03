////////////////////////////////////////////////////////////////////////////////
//
//  File           : swallet443.go
//  Description    : This is the implementaiton file for the swallet password
//                   wallet program program.  See assignment details.
//
//  Collaborators  : **TODO**: FILL ME IN
//  Last Modified  : **TODO**: FILL ME IN
//

// Package statement
package main

// Imports
import ( 
	"fmt"
	"os"
	"bufio"
	"time"
	"strings"
	"math/rand"
	"github.com/pborman/getopt"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	// There will likely be several mode APIs you need
)

// Type definition  ** YOU WILL NEED TO ADD TO THESE **

// A single password
type walletEntry struct {
	password []byte    // Should be exactly 32 bytes with zero right padding
	salt []byte        // Should be exactly 16 bytes 
	comment []byte     // Should be exactly 128 bytes with zero right padding
}

// The wallet as a whole
type wallet struct {
	filename string
	masterPassword []byte   // Should be exactly 32 bytes with zero right padding
	passwords []walletEntry
}

// Global data
var usageText string = `USAGE: swallet443 [-h] [-v] <wallet-file> [create|add|del|show|chpw|reset|list]

where:
    -h - help mode (display this message)
    -v - enable verbose output

    <wallet-file> - wallet file to manage
    [create|add|del|show|chpw] - is a command to execute, where

     create - create a new wallet file
     add - adds a password to the wallet
     del - deletes a password from the wallet
     show - show a password in the wallet
     chpw - changes the password for an entry in the wallet
     reset - changes the password for the wallet
     list - list the entries in the wallet (without passwords)`

var verbose bool = true

var path = ""

// You may want to create more global variables

//
// Functions

// Up to you to decide which functions you want to add

////////////////////////////////////////////////////////////////////////////////
//
// Function     : walletUsage
// Description  : This function prints out the wallet help
//
// Inputs       : none
// Outputs      : none

func walletUsage() {
	fmt.Fprintf(os.Stderr, "%s\n\n", usageText)
}

////////////////////////////////////////////////////////////////////////////////
//
// Function     : createWallet
// Description  : This function creates a wallet if it does not exist
//
// Inputs       : filename - the name of the wallet file
// Outputs      : the wallet if created, nil otherwise

func createWallet(filename string) *wallet {

	// Setup the wallet
	var wal443 wallet 
	wal443.filename = filename
	wal443.masterPassword = make([]byte, 32, 32) // You need to take it from here

	var newPath = path + filename
	var input, input2 []byte

	var _, err = os.Stat(newPath)
	if !os.IsNotExist(err) {				// checks if the file already exists
		fmt.Println("This file already exists" )
		return nil
	}

	fmt.Print("Enter Master Password (no longer than 32bytes): ")	// asking for the master password from the user		

	fmt.Scanln(&input)
	if cap(input) > 32 {	// check if the password is too long
		fmt.Print("Master Password must be no longer than 32 length\n")
		os.Exit(0)
	}
//	if cap(input) < 8 {	// check if the password is too short
//		fmt.Pringln("Please enter a longer Master Password\n")
//		os.Exit(0)
//	}

	fmt.Print("Re-enter Master Password: ")
	fmt.Scanln(&input2)

	// check if master passwords match
	if string(input) != string(input2) {
		fmt.Print("Master passwords do not match\n")
		return nil
	}

	// create the file
	f, err := os.Create(newPath)
	if err != nil {
		fmt.Println("Can't create wallet")
		os.Exit(0)
	}
	// close the file at the end
	defer f.Close()

	if err != nil {
		os.Exit(0)
	} else {
		wal443.masterPassword = input
		key := input

		// write the top line to the file
		topline := time.Now().String() + "\t" + "|| generation: 1 ||\n"
		f.WriteString(topline)

		// create hmac of the topline using master password as the key
		hash := hmac.New(sha1.New, key)
		hash.Write([]byte(topline))

		// encode the result of hmac into base64
		encode := base64.StdEncoding.EncodeToString(hash.Sum(nil))
		f.WriteString(encode)

		fmt.Println("Wallet created")
	}

	return &wal443
}

////////////////////////////////////////////////////////////////////////////////
//
// Function     : loadWallet
// Description  : This function loads an existing wallet
//
// Inputs       : filename - the name of the wallet file
// Outputs      : the wallet if created, nil otherwise

func loadWallet(filename string) *wallet {

	// Setup the wallet
	var wal443 wallet 
	// DO THE LOADING HERE
	// Open the file
	var newPath = path + filename
	f,_ := os.Open(newPath)
	defer f.Close()

	// ask for master password
	var key []byte
	fmt.Print("Please enter the Master Password: ")
	fmt.Scanln(&key)

	// The following code checks if the entered password is correct
	var content []string		// store all the content of the file

	scanner := bufio.NewScanner(f)	// read the file line by line
	for scanner.Scan() {
		str := scanner.Text()
		content = append(content, str)
	}

	// convert each line (except for the last line) to byte and write to hash
	hash := hmac.New(sha1.New, key)
	for i := 0; i < len(content) - 1; i++ {
		line := []byte(content[i] + "\n")	// For some readon + "\n" is needed
		hash.Write(line)
	}

	HMAC_value := content[ len(content)-1 ]		// The last line is the HMAC value

	// check hmac with the given password
	encode := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	if (encode == HMAC_value) {
		fmt.Println("Right password")
	}

	// Return the wall
	return &wal443
}

////////////////////////////////////////////////////////////////////////////////
//
// Function     : saveWallet
// Description  : This function save a wallet to the file specified
//
// Inputs       : walletFile - the name of the wallet file
// Outputs      : true if successful test, false if failure

func (wal443 wallet) saveWallet() bool {

	// Setup the wallet

	// Return successfully
	return true
}


func (wal443 wallet) addPassword(password string) bool {

	var newPath = path + wal443.filename
	f, _ := os.Open(newPath)
	defer f.Close()

	return true
}
////////////////////////////////////////////////////////////////////////////////
//
// Function     : processWalletCommand
// Description  : This is the main processing function for the wallet
//
// Inputs       : walletFile - the name of the wallet file
//                command - the command to execute
// Outputs      : true if successful test, false if failure

func (wal443 wallet) processWalletCommand(command string) bool {

	// Process the command 
	switch command {
	case "add":
		// DO SOMETHING HERE, e.g., wal443.addPassword(...)

	case "del":
		// DO SOMETHING HERE
		
	case "show":
		// DO SOMETHING HERE
		
	case "chpw":
		// DO SOMETHING HERE
		
	case "reset":
		// DO SOMETHING HERE
		
	case "list":
		// DO SOMETHING HERE
		
	default:
		// Handle error, return failure
		fmt.Fprintf(os.Stderr, "Bad/unknown command for wallet [%s], aborting.\n", command)
		return false
	}

	// Return sucessfull
	return true
}

////////////////////////////////////////////////////////////////////////////////
//
// Function     : main
// Description  : The main function for the password generator program
//
// Inputs       : none
// Outputs      : 0 if successful test, -1 if failure

func main() {

	// Setup options for the program content
	getopt.SetUsage(walletUsage)
	rand.Seed(time.Now().UTC().UnixNano())
	helpflag := getopt.Bool('h', "", "help (this menu)")
	verboseflag := getopt.Bool('v', "", "enable verbose output")

	// Now parse the command line arguments
	err := getopt.Getopt(nil)
	if err != nil {
		// Handle error
		fmt.Fprintln(os.Stderr, err)
		getopt.Usage()
		os.Exit(-1)
	}

	// Process the flags
	fmt.Printf("help flag [%t]\n", *helpflag)
	fmt.Printf("verbose flag [%t]\n", *verboseflag)
	verbose = *verboseflag
	if *helpflag == true {
		getopt.Usage()
		os.Exit(-1)
	}

	// Check the arguments to make sure we have enough, process if OK
	if getopt.NArgs() < 2 {
		fmt.Printf("Not enough arguments for wallet operation.\n")
		getopt.Usage()
		os.Exit(-1)
	}
	fmt.Printf("wallet file [%t]\n", getopt.Arg(0))
	filename := getopt.Arg(0)
	fmt.Printf("command [%t]\n", getopt.Arg(1))
	command := strings.ToLower(getopt.Arg(1))

	// Now check if we are creating a wallet
	if command == "create" {

		// Create and save the wallet as needed
		wal443 := createWallet(filename)
		if wal443 != nil {
			wal443.saveWallet()
		}

	} else {

		// Load the wallet, then process the command
		wal443 := loadWallet(filename)
		if wal443 != nil && wal443.processWalletCommand(command) {
			wal443.saveWallet()
		}

	}

	// Return (no return code)
	return
}
