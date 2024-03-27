package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	endpoint := "http://localhost:3000/"

	fmt.Println("Enter long URL")

	reader := bufio.NewReader(os.Stdin)

	url, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	url = strings.TrimSuffix(url, "\n")

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(url))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "text/plain; charset=utf-8")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println("Status code ", res.Status)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}
