package main

import (
	"flag"
	"fmt"
	"github.com/bwesterb/go-zonefile"
	"github.com/openrdap/rdap"
	"golang.org/x/net/publicsuffix"
	"os"
	"strings"
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

	for _, entry := range zf.Entries() {
		domain := string(entry.Domain())
		ns := string(entry.Values()[0])

		if string(entry.Type()) == "NS" && len(entry.Values()) > 0 {
			domains[domain] = append(domains[domain], ns)
		}
	}

	domainsNameservers := make(map[string][]string)

	for domain, ns := range domains {
		exists := false
		for _, ns := range ns {
			notFQDN, err := publicsuffix.EffectiveTLDPlusOne(strings.TrimSuffix(ns, "."))
			if err != nil {
				exit(fmt.Sprintf("Failed to get public suffix: %s, Error %s", ns, err))
			}
			for _, ns := range domainsNameservers[domain] {
				if notFQDN == ns {
					exists = true
					break
				}
			}
			if !exists {
				domainsNameservers[domain] = append(domainsNameservers[domain], notFQDN)
			}
		}
	}
	//	for domain, ns := range domainsNameservers {
	//		fmt.Printf("%s: %s\n", domain, ns)
	//	}

	rdapData := searchForDomains(domainsNameservers)

	for domain, info := range rdapData {
		fmt.Printf("Domain: %s, Registrar: %s\n", domain, info.Remarks)
	}
}

func searchForDomains(domainsNameservers map[string][]string) map[string]*rdap.Domain {
	client := &rdap.Client{}
	results := make(map[string]*rdap.Domain)
	for _, nsList := range domainsNameservers {
		for _, domain := range nsList {
			fmt.Println("Searching for", domain)
			data, err := client.QueryDomain(domain)
			if err != nil {
				fmt.Println("Error searching for", domain, err.Error())
			} else {
				fmt.Println(data.Events, "\n", data.LDHName, "\n")
				results[domain] = data
			}
		}
	}

	return results
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
