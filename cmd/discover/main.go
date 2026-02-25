// Command discover captures network traffic from YouTube and extracts ad-related
// domains that can be added to the blocklist. It can also generate the AdGuard/uBlock
// filter list from the blocklist.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/miguelmartens/opgecanceld-blocklist/internal/blocklist"
	"github.com/miguelmartens/opgecanceld-blocklist/internal/discover"
)

// Version is set at build time via -ldflags.
var Version string

const (
	defaultBlocklistPath = "opgecanceld-blocklist.txt"
	defaultFiltersPath   = "opgecanceld-filters.txt"
)

func main() {
	duration := flag.Duration("duration", 2*time.Minute, "How long to capture traffic")
	output := flag.String("output", "", "Output file for new domains (default: stdout)")
	doAppend := flag.Bool("append", false, "Append new domains to blocklist")
	buildFilters := flag.Bool("build-filters", false, "Generate AdGuard/uBlock filter list from blocklist (no discovery)")
	flag.Parse()

	if *buildFilters {
		runBuildFilters()
		return
	}

	existing, err := blocklist.LoadDomainSet(defaultBlocklistPath)
	if err != nil {
		log.Fatal(err)
	}

	client := discover.NewClient(discover.Config{
		Duration:  *duration,
		Blocklist: defaultBlocklistPath,
	}, existing)

	log.Println("Starting browser and capturing network traffic...")
	log.Printf("Will capture for %v. Browse YouTube or wait for ads to load.\n", *duration)

	newDomains, err := client.Run(context.Background())
	if err != nil {
		log.Fatalf("discovery failed: %v", err)
	}

	if len(newDomains) == 0 {
		log.Println("No new ad-related domains discovered.")
		return
	}

	log.Printf("Discovered %d new potential ad domains.\n", len(newDomains))

	if *doAppend {
		if err := blocklist.AppendDomains(defaultBlocklistPath, newDomains); err != nil {
			log.Fatal(err)
		}
		log.Printf("Appended %d domains to %s\n", len(newDomains), defaultBlocklistPath)
		if n, err := blocklist.GenerateFilters(defaultBlocklistPath, defaultFiltersPath); err != nil {
			log.Fatalf("generate filters: %v", err)
		} else {
			log.Printf("Generated %s with %d filter rules\n", defaultFiltersPath, n)
		}
	} else if *output != "" {
		if err := blocklist.WriteDomains(*output, newDomains); err != nil {
			log.Fatal(err)
		}
		log.Printf("Wrote %d domains to %s\n", len(newDomains), *output)
	} else {
		fmt.Println("# New domains to add:")
		for _, d := range newDomains {
			fmt.Println(d)
		}
	}
}

func runBuildFilters() {
	n, err := blocklist.GenerateFilters(defaultBlocklistPath, defaultFiltersPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Generated %s with %d filter rules\n", defaultFiltersPath, n)
}
