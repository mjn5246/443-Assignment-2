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
	"time"
	"strconv"
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
	if pattern != "" {
		var s string 
		for m := 0; m <= len(pattern)-1; m++ {
			s = string(pattern[m])
			l, err := strconv.ParseInt(s , 10, 64)
			if err != nil{
				l = 1
			} else {
				s = string(pattern[m - 1])
				l -= 1
			}
			for l > 0 {
				if s == "d"{			//for digit
					pwd += digit()
				} else if s == "c" {		//for either lower or upper case letter
					n := rand.Intn(2)	
					if n == 0{
						pwd += lower()
					} else {
						pwd += upper()
					}
				} else if s == "l" {		//for lower case letter
					pwd += lower()
				} else if s == "u" {		//for upper case letter
					pwd += upper()
				} else {
					//for word
				}
				l--
			}
		}
		return pwd
	}	
	var j int8
	for j = 1; j<=length; j++{
		i := rand.Intn(3) + 1
		if i == 1 {		//if 1 digit
			pwd += digit()
		} else if i == 2{	//if 2 lower case letter
			pwd += lower()
		} else {		//if 3 special character
			pwd += special()
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
		if _, err := strconv.Atoi(*length); err != nil {
			fmt.Printf("Bad length passed in [%s]\n", *length)
			fmt.Fprintln(os.Stderr, err)
			getopt.Usage()
			os.Exit(-1)
		}
		if plength <= 0 || plength > 64 {
			plength = 16
		}
		l, err := strconv.ParseInt(*length , 10, 8)
		if err != nil{
			os.Exit(-1)
		}
		plength = int8(l)
	}


	// Now generate the password and print it out
	pwd := generatePasword(plength, *pattern, *webflag)
	fmt.Printf("Generated password:  %s\n", pwd)

	// Return (no return code)
	return
}

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

// random English letter generagor
func character() string {
	rand.Seed(time.Now().UTC().UnixNano())
	if rand.Intn(2) % 2 == 1 {
		return lower()
	} else {
		return upper()
	}
}

// special characters: ~!@#$%^&*()-_=+{}[]:;/?<>,.|\

func special() string {
	characters := "~!@#$%^&*()-_=+{}[]:;/?<>,.|"
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(28)

	return characters[index:index+1]
}

