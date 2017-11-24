package main

// !!!Notice!!! 
// I change the package "rand" to "math/rand" because my go program can not find the package "rand"

import (
	"fmt"
	"math/rand"
	"time"
)

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

func character() string {
	rand.Seed(time.Now().UTC().UnixNano())
	if rand.Intn(2) {
		return lower()
	} else {
		return upper()
	}
}

// This function generates a password with a given pattern
// d := digit (is defined to be)
// c := lower or upper case letter
// l := lower case letter
// u := upper case letter

func main() {
	MyString := "ddlulddcc"
	output := ""
	for _, char := range MyString {
		c := string(char)
		switch c {
		case "d":
			output += digit()
		case "c":
			output += character()				
		case "l":
			output += lower()
		case "u":
			output += upper()
		}
	}
	fmt.Println(output)
}
