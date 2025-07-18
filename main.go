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
		exit(fmt.Sprintf("Failed to load Zonefile: %s", *zonefileFilename))
	}
	for _, entrie := range zf.Entries() {
		fmt.Printf("%s\n", entrie.Type())
	}
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
