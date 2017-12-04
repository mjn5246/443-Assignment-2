// This file records some helper functions that I have implemented

package main

import (
	"math/rand"
	"time"
	"fmt"
	"crypto/sha1"
	"strconv"
	"unocide"  
	"string"
)

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


// This function generates a crazy 16 byte salt
func saltGenerator() []byte {
	rand.Seed(time.Now().UTC().UnixNano())
	str1 := strconv.Itoa(rand.Intn(10000))
	str2 := time.Now().String()

	hash := sha1.New()
	hash.Write([]byte( str1+str2 ))
	bs := hash.Sum(nil)
	salt := []byte(bs[0:16])
	return salt
}

func main{
	// You can try to test any of these functions

}
