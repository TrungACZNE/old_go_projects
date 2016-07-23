package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"net/http"
	"net/url"

	"github.com/abh/geoip"
	"github.com/codegangsta/cli"
	"github.com/kr/pretty"
)

var (
	geodb   *geoip.GeoIP
	hostMap map[string]string

	backendHosts       map[string][]string
	defaultBackendHost []string
	transportList      = map[string]http.RoundTripper{}

	blacklist = []string{}

	hostRotation = map[string]int{}
)

func tunnel(w http.ResponseWriter, req *http.Request, backend string) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Can't CONNECT - webserver doest not support hijacking", http.StatusInternalServerError)
		return
	}

	conn, _, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() { _ = conn.Close() }()

	backConn, err := net.Dial("tcp", backend)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() { _ = backConn.Close() }()

	additionalHeaders := ""
	for key, vallist := range req.Header {
		for _, val := range vallist {
			additionalHeaders += fmt.Sprintf("%s: %s\n", key, val)
		}
	}

	host := stripPort(req.URL.Host)
	message := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n%s\r\n\r\n", req.URL.Host, host, additionalHeaders)

	_, err = backConn.Write([]byte(message))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go func() { _, _ = io.Copy(backConn, conn) }()
	_, _ = io.Copy(conn, backConn)
}

func stripPort(host string) string {
	if p := strings.Index(host, ":"); p != -1 {
		return host[:p]
	}
	return host
}

func selectHost(code string) string {
	var hostList []string
	if code == "*" {
		hostList = defaultBackendHost
	} else {
		hostList = backendHosts[code]
	}
	hostRotation[code] = (hostRotation[code] + 1) % len(hostList)
	return hostList[hostRotation[code]]
}

func chooseBackendForHost(host string) string {
	host = stripPort(host)
	tld := strings.ToLower(host[strings.LastIndex(host, ".")+1:])

	for mappedHost, code := range hostMap {
		if strings.HasPrefix(host, mappedHost) {
			tld = code
			break
		}
	}

	for code, _ := range backendHosts {
		if code == tld {
			return selectHost(code)
		}
	}

	ips, err := net.LookupHost(host)
	if err != nil {
		log.Println("Could not look up", host)
		return selectHost("*")
	}

	if len(ips) > 0 {
		ip := ips[0]
		countrycode, _ := geodb.GetCountry(ip)
		if _, ok := backendHosts[countrycode]; ok {
			return selectHost(countrycode)
		}

	}
	return selectHost("*")
}

func initTransport(backend string) error {
	backendurl, err := url.Parse("http://" + backend)
	if err != nil {
		return err
	}

	transportList[backend] = &http.Transport{
		Proxy: http.ProxyURL(backendurl),
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	return nil
}

func Proxy(w http.ResponseWriter, req *http.Request) {
	if len(blacklist) > 0 {
		fullurl := req.URL.String()
		for _, piece := range blacklist {
			if len(piece) > 2 && strings.Index(fullurl, piece) != -1 {
				log.Println("Blocked ", fullurl)
				http.NotFoundHandler().ServeHTTP(w, req)
				return
			}
		}
	}

	backend := chooseBackendForHost(req.Host)
	if req.Method == "CONNECT" {
		tunnel(w, req, backend)
		return
	}

	tp := transportList[backend]

	backendRequest := new(http.Request)
	*backendRequest = *req
	/*
		backendRequest.URL.Path = backendRequest.URL.String()
		backendRequest.URL.Host = backend*/

	res, err := tp.RoundTrip(backendRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() { _ = res.Body.Close() }()

	wHeader := w.Header()
	for key, vallist := range res.Header {
		for _, val := range vallist {
			wHeader.Add(key, val)
		}
	}

	w.WriteHeader(res.StatusCode)
	_, _ = io.Copy(w, res.Body)
}

func parseHostList(text string) ([]string, error) {
	// no array
	if !strings.HasPrefix(text, "[") && !strings.HasSuffix(text, "]") {
		return []string{text}, nil
	}

	// semicolon separated array
	if strings.HasPrefix(text, "[") && strings.HasSuffix(text, "]") {
		return strings.Split(text[1:len(text)-1], ";"), nil
	}

	return nil, fmt.Errorf("Uneven brackets %s", text)
}

// parseBackendHosts from command line flag to map[string]string
// requires at least one "*" key, does not check fo duplication
// text format: "key1=value2,key2=value2,..."
func parseBackendHosts(text string) (map[string][]string, []string, error) {
	defaultBackendHost := []string{}
	result := map[string][]string{}
	for _, block := range strings.Split(text, ",") {
		tokens := strings.Split(block, "=")
		if len(tokens) != 2 {
			return nil, nil, fmt.Errorf("Fail to parse %s", block)
		}
		countryCode := strings.ToLower(tokens[0])
		hostList, err := parseHostList(tokens[1])
		if err != nil {
			return nil, nil, err
		}

		if countryCode == "*" {
			defaultBackendHost = hostList
		} else {
			result[countryCode] = hostList
		}
	}
	if len(defaultBackendHost) == 0 {
		return result, nil, fmt.Errorf("No default host (*) found")
	}
	return result, defaultBackendHost, nil
}

func parseHostMap(text string) (map[string]string, error) {
	result := map[string]string{}
	for _, block := range strings.Split(text, ",") {
		tokens := strings.Split(block, "=")
		if len(tokens) != 2 {
			return nil, fmt.Errorf("Fail to parse %s", block)
		}
		host := strings.ToLower(tokens[0])
		countryCode := strings.ToLower(tokens[1])
		result[host] = countryCode
	}
	return result, nil
}

func loadBlacklist(file string) ([]string, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	text := string(bytes)
	return strings.Split(text, "\n"), nil
}

func start(backend, bind, geodat, hostmap, blacklistFile string) {
	var err error
	if backend == "" || bind == "" {
		log.Fatal("No backend or bind specified")
	}

	if geodat == "" {
		log.Fatal("Need to specify geodat")
	}

	geodb, err = geoip.Open(geodat)
	if err != nil {
		log.Fatal(err)
	}

	backendHosts, defaultBackendHost, err = parseBackendHosts(backend)
	if err != nil {
		log.Fatal(err)
	}
	for _, hostList := range backendHosts {
		for _, host := range hostList {
			err := initTransport(host)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	for _, host := range defaultBackendHost {
		err := initTransport(host)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Main list %# v\n", pretty.Formatter(backendHosts))
	fmt.Printf("Default list %# v\n", pretty.Formatter(defaultBackendHost))

	hostMap, err = parseHostMap(hostmap)
	if err != nil {
		log.Fatal(err)
	}

	if blacklistFile != "" {
		blacklist, err = loadBlacklist(blacklistFile)
	}
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(bind, http.HandlerFunc(Proxy))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	log.SetOutput(os.Stdout)
	runtime.GOMAXPROCS(runtime.NumCPU())

	app := cli.NewApp()
	app.Name = "Name"
	app.Usage = "Do stuff"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "backend",
			Value: "*=[0.0.0.0:3002]",
			Usage: "specifies backend: use a comma separated list, each element being a key=value pair, of which key is the country code and the value is a semicolon separated list of the backend hosts matching that key; country code be compared in lower case",
		},
		cli.StringFlag{
			Name:  "bind",
			Value: "0.0.0.0:3001",
			Usage: "specifies bind point",
		},
		cli.StringFlag{
			Name:  "geodat",
			Value: "/root/GeoIP.dat",
			Usage: "specifies geo IP data file",
		},
		cli.StringFlag{
			Name:  "map",
			Value: "vod-hds-uk-live=uk",
			Usage: "maps hostnames to country code manually",
		},
		cli.StringFlag{
			Name:  "blacklist",
			Value: "",
			Usage: "specifies file name containing black list pieces of a URL",
		},
	}
	app.Action = func(c *cli.Context) {
		start(c.String("backend"), c.String("bind"), c.String("geodat"), c.String("map"), c.String("blacklist"))
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Println("app.Run() error:", err)
	}
}
