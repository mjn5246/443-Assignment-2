////////////////////////////////////////////////////////////////////////////////
//
//  File           : spwgen443.go
//  Description    : This is the implementaiton file for the spwgen443 password
//                   generator program.  See assignment details.
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
	"math/rand"
	"strconv"
	"time"
	"github.com/pborman/getopt"
	// There will likely be several mode APIs you need
)

// Global data
var patternval string = `pattern (set of symbols defining password)

        A pattern consists of a string of characters "xxxxx",
        where the x pattern characters include:

          d - digit
          c - upper or lower case character
          l - lower case character
          u - upper case character
          w - random word from /usr/share/dict/words (or /usr/dict/words)
              note that w# will identify a word of length #, if possible
          s - special character in ~!@#$%^&*()-_=+{}[]:;/?<>,.|\

        Note: the pattern overrides other flags, e.g., -w`

// You may want to create more global variables

//
// Functions

// random digit generator
func digit() string {
	digits := "0123456789"
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(10)

	return digits[index:index+1]
}

// random lower case letter generator
func lower() string {
	letters := "qazwsxedcrfvtgbyhnujmikolp"
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(26)

	return letters[index:index+1]
}

// random upper case letter generator
func upper() string {
	letters := "QAZWSXEDCRFVTGBYHNUJMIKOLP"
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(26)

	return letters[index:index+1]
}

// special characters: ~!@#$%^&*()-_=+{}[]:;/?<>,.|
// special character random generator
func special() string {
	characters := "~!@#$%^&*()-_=+{}[]:;/?<>,.|"
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(28)

	return characters[index:index+1]
}

// NOTE! This function needs to be modified because I am not sure what
// characters are allowed in a web password
func web() string {
	characters := "~!@#$%^&"
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(8)

	return characters[index:index+1]
}

// random English letter generator
func character() string {
	rand.Seed(time.Now().UTC().UnixNano())
	if rand.Intn(2) % 2 == 1 {
		return lower()
	} else {
		return upper()
	}
}

// uniform random generator
func uniform(webflag bool) string {
	rand.Seed(time.Now().UTC().UnixNano())
	i := rand.Intn(3)
	switch i {
  case 0:
		return character()
	case 1:
		return digit()
	default:
		if webflag == true {	// if it is a web password
			return web()
		} else {
			return special()
		}
	}
}

// Up to you to decide which functions you want to add

////////////////////////////////////////////////////////////////////////////////
//
// Function     : generatePasword
// Description  : This is the function to generate the password.
//
// Inputs       : length - length of password
//                pattern - pattern of the file ("" if no pattern)
//                webflag - is this a web password?
// Outputs      : 0 if successful test, -1 if failure

func generatePasword(length int8, pattern string, webflag bool) string {

	pwd := "" // Start with nothing and add code
	if pattern != ""{ // if a pattern is provided
		for _, char := range pattern {
			c := string(char)
			switch c {
			case "d":
				pwd += digit()
			case "c":
				pwd += character()
			case "l":
				pwd += lower()
			case "u":
				pwd += upper()
			case "s":
				pwd += special()
			}
		}
	} else {	// if no pattern is given
		for i := int8(0); i<length; i++ {
			pwd += uniform(webflag)
		}
	}

	// Now return the password
	return pwd
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
	rand.Seed(time.Now().UTC().UnixNano())
	helpflag := getopt.Bool('h', "", "help (this menu)")
	webflag := getopt.Bool('w', "", "web flag (no symbol characters, e.g., no &*...)")
	length := getopt.String('l', "", "length of password (in characters)")
	pattern := getopt.String('p', "", patternval)

	// Now parse the command line arguments
	err := getopt.Getopt(nil)
	if err != nil {
		// Handle error
		fmt.Fprintln(os.Stderr, err)
		getopt.Usage()
		os.Exit(-1)
	}

	// Get the flags
	fmt.Printf("helpflag [%t]\n", *helpflag)
	fmt.Printf("webflag [%t]\n", *webflag)
	fmt.Printf("length [%s]\n", *length)
	fmt.Printf("pattern [%s]\n", *pattern)
	// Normally, we we use getopt.Arg{#) to get the non-flag paramters
	// Safety check length parameter
	var plength int8 = 16
	if *length != "" {
		if temp, err := strconv.Atoi(*length); err != nil {
			fmt.Printf("Bad length passed in [%s]\n", *length)
			fmt.Fprintln(os.Stderr, err)
			getopt.Usage()
			os.Exit(-1)
		} else {
			plength = int8(temp)
		}
		if plength <= 0 || plength > 64 {
			plength = 16
		}
	}


	// Now generate the password and print it out
	pwd := generatePasword(plength, *pattern, *webflag)
	fmt.Printf("Generated password:  %s\n", pwd)

	// Return (no return code)
	return
}
