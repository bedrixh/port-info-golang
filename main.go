package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	portNumber := -1
	portProtocol := "*"

	if portFromArg, err := strconv.ParseInt(os.Args[1], 10, 64); err == nil {
		portNumber = int(portFromArg)
	} else {
		splited := strings.Split(os.Args[1], "/")
		if len(splited) == 2 {
			portFromArg, err := strconv.ParseInt(splited[0], 10, 64)
			if err != nil {
				cliUsage()
			}
			portNumber = int(portFromArg)
			portProtocol = splited[1]
		}
	}

	if portNumber == -1 {
		cliUsage()
	}

	CSVPortData := getCSVString()

	r := csv.NewReader(strings.NewReader(CSVPortData))

	records, err := r.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	for i := range len(records) {
		if strconv.FormatInt(int64(portNumber), 10) == records[i][1] {
			if portProtocol == "*" || portProtocol == records[i][2] {
				for j := range len(records[0]) {
					if records[i][j] != "" {
						fmt.Printf("%s: %s\n", records[0][j], records[i][j])
					}
				}
				fmt.Println("------------")
			}
		}
	}

}

func cliUsage() {
	println(os.Args[0] + " ([number]) or ([number]/[protocol])")
	os.Exit(0)
}

func getCSVString() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	cacheDir = filepath.Join(cacheDir, "port-info-golang")
	cacheCSVPath := filepath.Join(cacheDir, "service-names-port-numbers.csv")

	data, err := os.ReadFile(cacheCSVPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			client := http.Client{}
			resp, _ := client.Get("https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.csv")
			responseCSV, err := io.ReadAll(resp.Body)

			if err != nil {
				panic("oh shit")
			}
			os.MkdirAll(cacheDir, 0775)
			os.WriteFile(cacheCSVPath, responseCSV, 0664)
			return string(responseCSV)
		}
		panic(err)
	}

	return string(data)
}
