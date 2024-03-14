package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

const (
	errArg = "check for argument correctnes"
)

func main() {
	uri := flag.String("s", "", "Uniform Resource Identifier, specify base url for successful download of the site")
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println(errArg)
		return
	}

	if ok, err := regexp.MatchString("^(http|https)://", *uri); !ok || err != nil {
		fmt.Println("invalid url")
		return
	}

	if err := wget(*uri); err != nil {
		log.Fatalf("[Error] %v\n", err)
	}

}

// wget func - write resualts of a GET request to a file.
func wget(url string) error {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create("index.html")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
