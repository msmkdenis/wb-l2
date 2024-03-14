package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
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
	rec := flag.Int("r", 1, "depth of recursion")
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println(errArg)
		return
	}

	err := os.Mkdir("html", 0755)
	if err != nil {
		return
	}

	if ok, err := regexp.MatchString("^(http|https)://", *uri); !ok || err != nil {
		fmt.Println("invalid url")
		return
	}

	if err := wget(*uri, *rec); err != nil {
		log.Fatalf("wget error: %s", err)
	}

}

func generateString(n int, generator rand.Source) string {
	result := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		randomNumber := generator.Int63()
		result = append(result, byte(randomNumber%26+97))
	}
	return string(result)
}

func wget(u string, n int) error {

	if n == 0 {
		return nil
	}

	rnd := rand.NewSource(time.Now().Unix())

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading page: %v\n", err)
		return nil
	}

	ur, err := url.Parse(u)
	if err != nil {
		return err
	}

	fmt.Println(ur)
	file, err := os.Create("html/" + generateString(4, rnd) + ".html")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(b)
	if err != nil {
		return err
	}

	links, err := parseLinks(ur, string(b))
	if err != nil {
		return err
	}

	for _, link := range links {
		if strings.HasPrefix(link, u) {
			err := wget(link, n-1)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return nil
}

func parseLinks(urlToGet *url.URL, content string) ([]string, error) {
	var (
		err       error
		links     []string = make([]string, 0)
		matches   [][]string
		findLinks = regexp.MustCompile("<a.*?href=\"(.*?)\"")
	)

	// Retrieve all anchor tag URLs from string
	matches = findLinks.FindAllStringSubmatch(content, -1)

	for _, val := range matches {
		var linkUrl *url.URL

		// Parse the anchr tag URL
		if linkUrl, err = url.Parse(val[1]); err != nil {
			return links, err
		}

		// If the URL is absolute, add it to the slice
		// If the URL is relative, build an absolute URL
		if linkUrl.IsAbs() {
			links = append(links, linkUrl.String())
		} else {
			links = append(links, urlToGet.Scheme+"://"+urlToGet.Host+linkUrl.String())
		}
	}

	return links, err
}
