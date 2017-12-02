////////////////////////////////////////////////////////////////////////////////
//
//  File           : spwgen443.go
//  Description    : This is the implementaiton file for the spwgen443 password
//                   generator program.  See assignment details.
//
//  Collaborators  : **TODO**: Patrick Colville, Mauro Notaro, Weiyu Luo
//  Last Modified  : **TODO**: 2017/12/1 19:35
//

// Package statement
package main

// Imports
import (
	"fmt"
	"os"
	"math/rand"
	"bufio"
	"strconv"
	"time"
	"github.com/pborman/getopt"
	"strings"
	"unicode"
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

var countWords [][]string	// define a hash table of a set of strings

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

// characters which are allowed in a web password
func web() string {
	rand.Seed(time.Now().UTC().UnixNano())
	if rand.Intn(2) == 0 {
		return character()
	} else {
		return digit()
	}
}

// random English letter generator
func character() string {
	rand.Seed(time.Now().UTC().UnixNano())
	if rand.Intn(2) == 1 {
		return lower()
	} else {
		return upper()
	}
}

// uniform random generator
func uniform(webflag bool) string {
	rand.Seed(time.Now().UTC().UnixNano())
	if webflag == true {
		return web()
	}

	i := rand.Intn(3)
	switch i {
  case 0:
		return character()
	case 1:
		return digit()
	default:
		return special()
	}
}

// This functions selects a random word from the file
func Word(dict [][]string) string {	// dict is a hash table of set of strings
	rand.Seed(time.Now().UTC().UnixNano())
	
	sum := 0		// sum is the total number of words in the dict
	for _, s := range dict {
		sum += len(s)
	}

	index := rand.Intn(sum)
	i := 0
	for index > len(dict[i]) - 1 {
		index -= len(dict[i])
		i++
	}

	return dict[i][index]
}

// This function selects a random word with a given length from the file
func WordLength(dict [][]string, length int) string {
	if length > len(dict) - 1 || len( dict[length] ) == 0 {		
	// if the word length exceeds the size of the hash table or the bucket is empty
		fmt.Printf("Generated password: None\n")
		fmt.Printf("word of length %d does not exist\n", length)
		os.Exit(-1)
		return "-1"
	} else {
		index := rand.Intn( len(dict[length]) )
		return dict[length][index]
	}
}

// The function split splits a string into a slice of characters
func split(str string) []string {
	var mySlice []string
	i := 0
	for i < len(str) {
		if strings.ContainsAny(str[i:i+1], "dclus") {
			mySlice = append(mySlice, str[i:i+1])
			i++
		} else if str[i:i+1] == "w" {		// detect if the length of the word is given
			j := i+1
			for j < len(str) && unicode.IsDigit( []rune(str)[j] ) {	// if this position is a digit
				j++
			}
			mySlice = append(mySlice, str[i:j])
			i = j
		} else {
			fmt.Println("Invalid pattern")
			os.Exit(-1)
		}
	}
	return mySlice
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
	if pattern != ""{ 	// if a pattern is provided
		mySlice := split(pattern)
		for _, str := range mySlice {

			switch str {		// differentiate different cases
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
			default:			// if the pattern asks to select a word
				var word string
				if len(str) == 1 {	// if it is a random word
					word = Word(countWords)
				} else {		// random word with a given length
					length, _ := strconv.Atoi(str[1:])
					word = WordLength(countWords, length)
				}

				if len(pwd) + len(word) >= 64 {
					fmt.Println("Maximum size of a password is reached.")
					return pwd
				} else {
					pwd += word
				}
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

  // Read the file /usr/share/dict/words and place words of different lengths
	// into different buckets
	f,_ := os.Open("/usr/share/dict/words")
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		str := scanner.Text() 			// read in a new word
		if !strings.Contains(str, "'") {	// if the word does not contain "'"
			if len(countWords) <= len(str) {
				countWords = append(countWords, make([][]string, len(str)-len(countWords)+1)...)
				countWords[len(str)] = append(countWords[len(str)], str)
			} else {
				countWords[len(str)] = append(countWords[len(str)], str)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	// Now generate the password and print it out
	pwd := generatePasword(plength, *pattern, *webflag)
	fmt.Printf("Generated password:  %s\n", pwd)

	// Return (no return code)
	return
}
