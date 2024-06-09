package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type data struct {
	hosts  []string
	paths  []string
	params []string
}

func main() {
	var filePath string
	var printHelp bool

	flag.StringVar(&filePath, "f", "", "File path to read CSV")
	flag.BoolVar(&printHelp, "h", false, "Print help")

	flag.Parse()

	if printHelp {
		flag.PrintDefaults()
		return
	}

	if filePath == "" {
		fmt.Println("Please provide a file path using the -f flag.")
		return
	}

	data := readCsv(filePath)
	if data == nil {
		return
	}

	writeArrayToFile("hosts.txt", data.hosts)
	writeArrayToFile("paths.txt", data.paths)
	writeArrayToFile("params.txt", data.params)
}

func writeArrayToFile(fileName string, lines []string) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func readCsv(fileName string) *data {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("File not exists")
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Can't read file")
		return nil
	}

	if len(records) == 0 {
		fmt.Println("No records found")
		return nil
	}

	header := records[0]
	hostIndex := -1
	pathIndex := -1
	queryIndex := -1
	for i, col := range header {
		if col == "host" {
			hostIndex = i
		} else if col == "path" {
			pathIndex = i
		} else if col == "query" {
			queryIndex = i
		}
	}

	if hostIndex == -1 || pathIndex == -1 || queryIndex == -1 {
		fmt.Println("Host or path or query column not found")
		return nil
	}

	data := &data{
		hosts:  []string{},
		paths:  []string{},
		params: []string{},
	}

	uniqueHosts := make(map[string]bool)
	uniquePaths := make(map[string]bool)
	uniqueParams := make(map[string]bool)

	for _, record := range records[1:] {
		host := processHost(record[hostIndex])
		path := record[pathIndex]
		param := record[queryIndex]

		if host != "" {
			appendUnique(&data.hosts, host, &uniqueHosts)
		}
		if path != "/" {
			appendUnique(&data.paths, path, &uniquePaths)
		}
		if param != "" {
			paramNames := processParams(param)
			for _, paramName := range paramNames {
				appendUnique(&data.params, paramName, &uniqueParams)
			}
		}
	}
	return data
}

func appendUnique(slice *[]string, item string, uniqueMap *map[string]bool) {
	if _, exists := (*uniqueMap)[item]; !exists {
		(*uniqueMap)[item] = true
		*slice = append(*slice, item)
	}
}

func processHost(rawHost string) string {
	parts := strings.Split(rawHost, ".")
	if len(parts) > 2 {
		return strings.Join(parts[:len(parts)-2], ".")
	}
	return ""
}

func processParams(query string) []string {
	values, err := url.ParseQuery(query)
	if err != nil {
		return nil
	}

	var paramNames []string
	for key := range values {
		paramNames = append(paramNames, key)
	}
	return paramNames
}
