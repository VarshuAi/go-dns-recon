package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	domain := flag.String("d", "", "Target domain (e.g. google.com)")
	wordlist := flag.String("w", "subdomains.txt", "Path to subdomain wordlist file")
	threads := flag.Int("t", 20, "Number of concurrent lookup threads")
	flag.Parse()

	if *domain == "" {
		fmt.Println("Usage: dns-recon -d <domain> -w <wordlist> -t <threads>")
		return
	}

	file, err := os.Open(*wordlist)
	if err != nil {
		fmt.Printf("[-] Unable to open wordlist file: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Printf("[*] Starting DNS Subdomain Recon on: %s\n", *domain)
	
	jobs := make(chan string, 100)
	var wg sync.WaitGroup

	for i := 0; i < *threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for sub := range jobs {
				target := fmt.Sprintf("%s.%s", sub, *domain)
				ips, err := net.LookupIP(target)
				if err == nil {
					fmt.Printf("[+] Discovered: %s -> %v\n", target, ips)
				}
			}
		}()
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sub := strings.TrimSpace(scanner.Text())
		if sub != "" {
			jobs <- sub
		}
	}
	close(jobs)
	wg.Wait()
	fmt.Println("[+] Subdomain discovery complete.")
}