package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/miekg/dns"
)

type apiNotifyResponse struct {
	Error  string
	Result []NotifyResponse
}

// NotifyResponse is returned by the API (as a list, since there can
// be multiple responses)
type NotifyResponse struct {
	Server string
	Result string
	Error  bool
}

var (
	domainFlag = flag.String("domain", "", "Domain to notify")
	verbose    = flag.Bool("verbose", false, "Be extra verbose")
	quiet      = flag.Bool("quiet", false, "Only output on errors")
	timeout    = flag.Int64("timeout", 2000, "Timeout for response (in milliseconds)")
	listen     = flag.String("listen", "", "Listen on this ip:port for the HTTP API")
)

var servers []string

func main() {

	flag.Parse()

	servers = flag.Args()

	if len(*domainFlag) == 0 && len(*listen) == 0 {
		fmt.Println("-listen or -domain parameter required\n")
		flag.Usage()
		os.Exit(2)
	}

	if len(*listen) == 0 {
		sendNotify(servers, *domainFlag)
		return
	}

	startHTTP(*listen)
}

func buildMux() *http.ServeMux {

	mux := http.NewServeMux()

	restHandler := rest.ResourceHandler{}
	restHandler.EnableGzip = true
	restHandler.EnableLogAsJson = true
	restHandler.EnableResponseStackTrace = true
	//restHandler.EnableStatusService = true

	restHandler.SetRoutes(
		&rest.Route{"POST", "/api/v1/notify/*domain", notifyHandler},
	)

	mux.Handle("/api/v1/", &restHandler)

	return mux

}

func startHTTP(listen string) {
	fmt.Printf("Listening on http://%s\n", listen)
	err := http.ListenAndServe(listen, buildMux())
	fmt.Printf("Could not listen to %s: %s", listen, err)
}

func notifyHandler(w rest.ResponseWriter, r *rest.Request) {

	domain := r.PathParam("domain")

	resp := new(apiNotifyResponse)

	resp.Result = sendNotify(servers, domain)

	for _, r := range resp.Result {
		if r.Error {
			resp.Error = r.Result
		}
	}

	w.WriteJson(resp)

}

func sendNotify(servers []string, domain string) []NotifyResponse {

	if !strings.HasSuffix(domain, ".") {
		domain = domain + "."
	}

	if len(servers) == 0 {
		fmt.Println("No servers")
		resp := NotifyResponse{Result: "No servers", Error: true}
		fmt.Println("No servers")
		return []NotifyResponse{resp}
	}

	c := new(dns.Client)

	c.ReadTimeout = time.Duration(*timeout) * time.Millisecond

	m := new(dns.Msg)
	m.SetNotify(domain)

	responseChannel := make(chan NotifyResponse, len(servers))

	for _, server := range servers {

		go func(server string) {

			result := NotifyResponse{Server: server}

			defer func() {
				if result.Error || !*quiet {
					fmt.Printf("%s: %s\n", result.Server, result.Result)
				}
				responseChannel <- result
			}()

			target, err := fixupHost(server)
			if err != nil {
				result.Result = fmt.Sprintf("%s: %s", server, err)
				fmt.Println(result.Result)
				result.Error = true
				return
			}

			result.Server = target

			if *verbose {
				fmt.Println("Sending notify to", target)
			}

			resp, rtt, err := c.Exchange(m, target)

			if err != nil {
				result.Error = true
				result.Result = err.Error()
				return
			}

			ok := "ok"
			if !resp.Authoritative {
				ok = fmt.Sprintf("not ok (%s)", dns.RcodeToString[resp.Rcode])
			}

			result.Result = fmt.Sprintf("%s: %s (%s)",
				target, ok, rtt.String())

		}(server)

	}

	responses := make([]NotifyResponse, len(servers))

	for i := 0; i < len(servers); i++ {
		responses[i] = <-responseChannel
	}

	return responses

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
