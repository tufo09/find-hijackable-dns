package main

import (
	"flag"
	"fmt"
	"github.com/bwesterb/go-zonefile"
	"os"
)

func main() {
	zonefileFilename := flag.String("zonefile", "gov.txt", "file name of the zonefile you want to analyze")
	flag.Parse()

	file, err := os.ReadFile(*zonefileFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open Zonefile: %s", *zonefileFilename))
	}
	zf, err := zonefile.Load(file)
	if err != nil {
		exit(fmt.Sprintf("Failed to load Zonefile: %s\n Error: %s", *zonefileFilename, err))
	}

	domains := make(map[string][]string)

	for _, entrie := range zf.Entries() {
		domain := string(entrie.Domain())
		ns := string(entrie.Values()[0])

		if string(entrie.Type()) == "NS" && len(entrie.Values()) > 0 {
			domains[domain] = append(domains[domain], ns)
		}
	}
	for domain, ns := range domains {
		fmt.Printf("%s, %s\n", domain, ns)
	}
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
