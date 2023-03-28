package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

const (
	fp  = "list.txt"
	rfp = "result.md"
)

func main() {
	urls, err := readLines(fp)
	if err != nil {
		panic(err)
	}

	var rslt []string
	for _, url := range urls {
		title, err := getTitle(url)
		if err != nil {
			rslt = append(rslt, fmt.Sprintf("[%s](%s)", "", url))
			continue
		}
		rslt = append(rslt, fmt.Sprintf("[%s](%s)", title, url))
	}
	if err = writeLines(rfp, rslt); err != nil {
		panic(err)
	}
}

func readLines(fp string) ([]string, error) {
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	s := bufio.NewScanner(f)
	var ls []string
	for s.Scan() {
		ls = append(ls, s.Text())
	}
	if err := s.Err(); err != nil {

		return nil, err
	}
	return ls, nil
}

var re = regexp.MustCompile(`<title>([\s\S]*?)</title>`)

func getTitle(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	match := re.FindStringSubmatch(string(body))
	if len(match) >= 2 {
		return match[1], nil
	}
	return "", nil
}

func writeLines(fp string, ls []string) error {
	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range ls {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
