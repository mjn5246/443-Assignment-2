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
	"io/ioutil"
	"log"

	"time"
	"strings"
	"strconv"

	"math/rand"
	"github.com/pborman/getopt"

	"crypto/hmac"
	"crypto/sha1"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"unicode"
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
	generation int
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

// This function finds the generation number from the topline of a wallet
func findGeneration(topline string) int {

	parts := strings.Split(topline, "||")
	// parts look like {time, generation: #, \n}

	str := parts[1]

	// find the starting digit
	i := 0
	for i < len(str) && !unicode.IsDigit( []rune(str)[i] ) {
		i++
	}

	// find the ending digit
	j := i + 1
	for j < len(str) && unicode.IsDigit( []rune(str)[j] ) {
		j++
	}

	// define the generation number
	generation, _ := strconv.Atoi(str[i:j])

	return generation
}

// This function generates a 16 byte long salt
// Salt is encoded using base64 at the end
func saltGenerator() []byte {
	rand.Seed(time.Now().UTC().UnixNano())
	str := ""
	digits := "0123456789"
	letters := "qazwsxedcrfvtgbyhnujmikolp" + "QAZWSXEDCRFVTGBYHNUJMIKOLP"
	specials := "~!@#$%^&*()-_=+{}[]:;/?<>,.|"

	for i := 0; i < 16; i++ {
		number := rand.Intn(3)
		switch number {
		case 0:
			index := rand.Intn(10)
			str += digits[index: index+1]
		case 1:
			index := rand.Intn(52)
			str += letters[index: index+1]
		default:
			index := rand.Intn(28)
			str += specials[index: index+1]
		}
	}
	// encode the string into base64
	encode := base64.StdEncoding.EncodeToString([]byte(str))

	// return a 16-byte long salt
	return []byte(encode[0:16])
}

// This function generates a key by taking the top 16 bytes of sha1 hash
// of the master password
func keyGenerator( masterPassword []byte ) []byte {
	hash := sha1.New()
	hash.Write(masterPassword)
	result := hash.Sum(nil)

	return result[0:16]
}

// This function pads a password to 32 bytes long
func padding(password []byte) []byte {
	pad := 16 - len(password)
	password_long := append(password, make([]byte, pad, pad)...)

	return password_long
}

// This function encrypts the password and encodes the output using base64
func AES_encrypt(key []byte, salt []byte, password []byte) []byte {
	
	plaintext := append(salt, padding(password)...)
	block, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(block)

	nonce := make([]byte, 12)
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	// encode using base 64
	encode := base64.StdEncoding.EncodeToString(ciphertext)
	return []byte(encode)
}

func AES_decrypt(key []byte, encode []byte) []byte {

	// decode the encoded ciphertext
	ciphertext, _ := base64.StdEncoding.DecodeString(string(encode))

	nonce := make([]byte, 12)

	block, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(block)
	plaintext, _ := aesgcm.Open(nil, nonce, ciphertext, nil)

	// igore the salt
	return plaintext[16:]
}

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
	for cap(input) > 32 {	// check if the password is too long
		fmt.Print("Master Password must be no longer than 32 length\n")
		fmt.Print("Try a different Master Password: ")
		fmt.Scanln(&input)
	}
//	if cap(input) < 8 {	// check if the password is too short
//		fmt.Pringln("Please enter a longer Master Password\n")
//		os.Exit(0)
//	}

	fmt.Print("Confirm Master Password: ")
	fmt.Scanln(&input2)

	// check if master passwords match
	if string(input) != string(input2) {
		fmt.Print("Master passwords do not match\n")
		return nil
	}

	// create the file
	f, err := os.Create(newPath)
	defer f.Close()					// close file at the end
	if err != nil {
		fmt.Println("Can't create wallet")
		os.Exit(0)
	} else {
		// store info into wal443
		wal443.filename = filename
		wal443.masterPassword = input
		wal443.generation = 0

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

	var newPath = path + filename

	// ask for master password
	var key []byte
	fmt.Print("Please enter the Master Password: ")
	fmt.Scanln(&key)

	// ##### The following code checks if the entered password is correct #####
	// read everything in the file
        input, err := ioutil.ReadFile(newPath)
        if err != nil {
                log.Fatalln(err)
        }

        content := strings.Split(string(input), "\n")	// content has type []string

	// convert each line (except for the last line) to byte and write to hash
	hash := hmac.New(sha1.New, key)
	for i := 0; i < len(content) - 1; i++ {
		line := []byte(content[i] + "\n")	// For some readon + "\n" is needed
		hash.Write(line)
	}

	HMAC_value := content[ len(content)-1 ]		// The last line is the HMAC value

	// check hmac with the given password
	encode := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	if (encode != HMAC_value) {
		fmt.Println("Wrong password")
		return nil
	}

 	// ##### The following code loads data into wal443 #####
	// record the masterPassword and find the genertion number from the top line
	wal443.filename = filename
	wal443.masterPassword = key
	wal443.generation = findGeneration(content[0])

	// try to load the enties of the wallet
	for i := 1; i < len(content) - 1; i++ {
		parts := strings.Split(content[i], " || ")
		// parts look like {entry, salt, password, comment}

		var entry walletEntry
		entry.salt = []byte(parts[1])
		entry.password = []byte(parts[2])
		entry.comment = []byte(parts[3])

		// append the entry to the field passwords
		wal443.passwords = append(wal443.passwords, entry)
	}
	fmt.Println("Wallet loaded successfully\n")

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
	// generation number is added by 1 every time we save the wallet
	// Setup the wallet
	var newPath = path + wal443.filename
	f,_ := os.Create(newPath)			// The purpose is to rewrite the file
	hash := hmac.New(sha1.New, wal443.masterPassword)		// Define hash

	// Write the topline
	gen := strconv.Itoa(wal443.generation + 1)
	topline := time.Now().String() + " || generation: " + gen + " ||\n"
	f.WriteString(topline)
	hash.Write([]byte(topline))

	// Write the entries
	for i := 0; i < len(wal443.passwords); i++ {
		entry := strconv.Itoa(i+1)
		line := entry + " || " + string(wal443.passwords[i].salt) + " || " + string(wal443.passwords[i].password)
		line = line + " || " + string(wal443.passwords[i].comment) + "\n"
		f.Write([]byte(line))
		hash.Write([]byte(line))
	}

	// encode the result of hmac into base64 and write the last line (HMAC)
	encode := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	f.WriteString(encode)

	// Return successfully
	return true
}

// This function adds a password to the wallet
func (wal443 *wallet) addPassword() bool {

	var input, input2, comment []byte
	// ask for the password to be added
	fmt.Print("Please enter a password (no longer than 16 bytes): ")

	fmt.Scanln(&input)
	for cap(input) > 16 {	// check if the password is too long
		fmt.Print("Password is too long, try a different password: ")
		fmt.Scanln(&input)
	}

	fmt.Print("Confirm the password: ")
	fmt.Scanln(&input2)

	// check if the passwords match
	if string(input) != string(input2) {
		fmt.Print("Two passwords do not match\n")
		return false
	}

	// prompt for comment
	fmt.Print("Enter any comment for the password (maxsize is 128 bytes): ")
	fmt.Scanln(&comment)
	for len(comment) > 128 {
		fmt.Print("Comment is too long, re-enter the comment: ")
		fmt.Scanln(&comment)
	}

	salt := saltGenerator()				// salt is []byte; base64
	key := keyGenerator(wal443.masterPassword)	// key is []byte

	// Generate the encryption of password using aes, the output is base64
	password := AES_encrypt(key, salt, input)	// password is []byte, base64

	// Add the new password to wal443, note that the wal443 is modified after this function returns
	var entry = walletEntry{password, salt, comment}
	wal443.passwords = append(wal443.passwords, entry)
	
	fmt.Println("Password added successfully")

	return true
}

func (wal443 *wallet) delPassword() bool {
	return true
}

func (wal443 *wallet) showPassword() bool {
	var entry_number int

	if len(wal443.passwords) == 0 {
		fmt.Println("No password is in the wallet")
		return false
	}

	fmt.Print("Please enter an entry number: (from 1 to ", len(wal443.passwords), " ): ")
	fmt.Scanln(&entry_number)

//	for i:=0; i<len(wal443.passwords); i++ {
//		fmt.Printf("%s\n", wal443.passwords[i].password)
//	}

	ciphertext := wal443.passwords[entry_number-1].password
	key := keyGenerator(wal443.masterPassword)

	// decrypt the encrypted password
	plaintext := AES_decrypt(key, ciphertext)

	// display password and comment
	line := "Password: " + string(plaintext) + " || comment: "
	line += string(wal443.passwords[entry_number-1].comment)
	fmt.Println(line)

	return true
}

func (wal443 *wallet) changePassword() bool {


	return true
}

// This function resets the master password
func (wal443 *wallet) reset() bool {

	var input, input2 []byte

	// Prompt for a new master password
	fmt.Print("Please enter a new Master Password (no longer than 32bytes): ")
	fmt.Scanln(&input)

	for cap(input) > 32 {	// check if the password is too long
		fmt.Print("Master Password must be no longer than 32 length\n")
		fmt.Print("Try a different Master Password: ")
		fmt.Scanln(&input)
	}

	fmt.Print("Confirm Master Password: ")
	fmt.Scanln(&input2)

	// check if master passwords match
	if string(input) != string(input2) {
		fmt.Print("Master passwords do not match\n")
		return false
	}
	
	// generate old_key and new_key
	old_key := keyGenerator(wal443.masterPassword)
	new_key := keyGenerator(input)

	// reset master password
	wal443.masterPassword = input

	// we need to re-encrypt all passwords using the new master password
	for i := 0; i < len(wal443.passwords); i++ {
		old_ciphertext := wal443.passwords[i].password
		salt := wal443.passwords[i].salt

		// plaintext of the password
		plaintext := AES_decrypt(old_key, old_ciphertext)

		// re-encrypt the password
		new_ciphertext := AES_encrypt(new_key, salt, plaintext)
		wal443.passwords[i].password = new_ciphertext
	}

	fmt.Println("Master Password is reset successfully")
	// return succeccfully
	return true
}

// This function lists all entries in the wallet
func (wal443 *wallet) list() bool {

	// if no password is in the wallet
	if len(wal443.passwords) == 0 {
		fmt.Println("No password is in the wallet")
		return false
	}
	
	// Write the entries
	fmt.Println("entry || comment")
	for i := 0; i < len(wal443.passwords); i++ {
		entry := strconv.Itoa(i+1)
		line := entry + " || " + string(wal443.passwords[i].comment)
		fmt.Println(line)
	}

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

func (wal443 *wallet) processWalletCommand(command string) bool {

	// Process the command
	switch command {
	case "add":
		wal443.addPassword()
		// DO SOMETHING HERE, e.g., wal443.addPassword(...)

	case "del":
		// DO SOMETHING HERE

	case "show":
		wal443.showPassword()

		// The return is set to false because there is no need to call savewallet()
		return false

	case "chpw":
		// DO SOMETHING HERE

	case "reset":
		wal443.reset()

	case "list":
		wal443.list()

		// The return is set to false because there is no need to call savewallet()
		return false

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
