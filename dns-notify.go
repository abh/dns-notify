package main

import (
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	domainFlag = flag.String("domain", "", "Domain to notify, required")
	verbose    = flag.Bool("verbose", false, "Be extra verbose")
	quiet      = flag.Bool("quiet", false, "Only output on errors")
	timeout    = flag.Int64("timeout", 2000, "Timeout for response (in milliseconds)")
)

func main() {

	flag.Parse()

	servers := flag.Args()

	if len(*domainFlag) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	sendNotify(servers, *domainFlag)

}

func sendNotify(servers []string, domain string) {

	if !strings.HasSuffix(domain, ".") {
		domain = domain + "."
	}

	if len(servers) == 0 {
		fmt.Println("No servers")
	}

	c := new(dns.Client)
	c.ReadTimeout = time.Duration(*timeout) * time.Millisecond

	m := new(dns.Msg)
	m.SetNotify(domain)

	wg := new(sync.WaitGroup)

	for _, server := range servers {

		serverPort, err := fixupHost(server)
		if err != nil {
			fmt.Printf("%s: %s\n", server, err)
			continue
		}

		wg.Add(1)

		go func(target string) {
			defer wg.Done()

			if *verbose {
				fmt.Println("Sending notify to", target)
			}

			resp, rtt, err := c.Exchange(m, target)

			if err != nil {
				fmt.Printf("%s: %s\n", target, err.Error())
				return
			}

			ok := "ok"
			if !resp.Authoritative {
				ok = fmt.Sprintf("not ok (%s)", dns.RcodeToString[resp.Rcode])
			}

			if !*quiet {
				fmt.Printf("%s: %s (%s)\n",
					target, ok, rtt.String())
			}
		}(serverPort)

	}

	wg.Wait()

}

func fixupHost(s string) (string, error) {

	_, _, err := net.SplitHostPort(s)
	if err != nil && strings.HasPrefix(err.Error(), "missing port in address") {
		return s + ":53", nil
	}
	if err != nil {
		return "", err
	}

	// input was ok ...
	return s, nil

}
