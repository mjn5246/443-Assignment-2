package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"math/rand"
	"time"
	"unicode"
)
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

func main(){
	// /usr/share/dict/words
	f,_ := os.Open("/tmp/File")
	var countWords [][]string	// define a hash table of a set of strings
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
	//for i, v := range countWords{
	//	fmt.Printf("%d words have length %d\n", len(v), i)
	//}

	// fmt.Println("Words of length 5:", countWords[5])
	// test if split works
	str := "w2d122"
	fmt.Println(split(str))

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
