package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

func main() {
	endpoint := "http://localhost:8080/"

	fmt.Println("Enter your URL")

	reader := bufio.NewReader(os.Stdin)

	url, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	url = strings.TrimSuffix(url, "\n")

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "text/plain; charset=utf-8").
		SetBody(url).
		Post(endpoint)

	if err != nil {
		panic(err)
	}

	fmt.Println("Status code ", resp.Status())
	fmt.Println(string(resp.Body()))
}
