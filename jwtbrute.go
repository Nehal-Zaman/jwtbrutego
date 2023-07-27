package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt"
)

func main() {

	wordlist_file, token, numberOfWorkers := getCliArgs()

	bytes, err := os.ReadFile(wordlist_file)
	checkError(err)
	wordlist := strings.Split(string(bytes), "\n")

	wordlistChannel := make(chan string)
	var wg sync.WaitGroup
	wg.Add(numberOfWorkers)

	for i := 0; i < numberOfWorkers; i++ {
		go func() {
			defer wg.Done()

			for word := range wordlistChannel {
				if checkJwtKey(word, token) {
					fmt.Printf("Found valid key: %v\n", word)
					os.Exit(0)
				}
			}
		}()
	}

	for _, word := range wordlist {
		wordlistChannel <- word
	}

	close(wordlistChannel)

	wg.Wait()

}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func getCliArgs() (string, string, int) {
	printBanner()

	wordlistPtr := flag.String("wordlist", "", "specify the wordlist file to use")
	tokenPtr := flag.String("token", "", "specify the JWT token for the attack")
	numOfWorkersPtr := flag.Int("threads", 20, "specify number of threads (default: 20)")
	flag.Parse()

	wordlist := *wordlistPtr
	token := *tokenPtr
	numOfWorkers := *numOfWorkersPtr

	if wordlist == "" || token == "" {
		flag.Usage()
		panic("Wordlist/token must be given")
	}

	return wordlist, token, numOfWorkers
}

func checkJwtKey(key, tokenString string) bool {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return false
	}
	return true
}

func printBanner() {
	fmt.Println("---------------------------------------------------------")
	fmt.Println("|  jwtbrute: multithreaded jwt secret finder in golang  |")
	fmt.Println("|    Author: n3hal_ (github.com/Nehal-Zaman)            |")
	fmt.Println("---------------------------------------------------------")
	fmt.Println("")
}
