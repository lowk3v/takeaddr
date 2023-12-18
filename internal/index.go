package internal

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/lowk3v/takeaddr/internal/enum"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Clone from https://github.com/003random/getJS

type GetJS struct {
	Options Options
}
type Options struct {
	Action     enum.ACTION
	Version    bool
	Output     string
	Verbose    bool
	Url        string
	Method     string
	OutputFile string
	InputFile  string
	Resolve    bool
	Complete   bool
	NoColors   bool
	Headers    []string
	Insecure   bool
	Timeout    int
}

func Run(opt Options) ([]string, error) {
	var allSources []string
	var urls []string

	if opt.InputFile != "" {
		f, err := os.Open(opt.InputFile)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Couldn't open input file: %v", err)
			return nil, err
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Couldn't close input file: %v", err)
			}
		}(f)

		lines, err := readLines(f)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Couldn't read from input file: %v", err)
		}
		_, _ = fmt.Fprintf(os.Stdout, "Read urls from input file: %s\n", opt.InputFile)
		urls = append(urls, lines...)
	}

	if opt.Url != "" {
		urls = append(urls, opt.Url)
	}

	if len(urls) == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "No urls supplied\n")
		os.Exit(3)
	}

	if opt.Resolve && !opt.Complete {
		_, _ = fmt.Fprintf(os.Stderr, "Resolve can only be used in combination with -complete\n")
		os.Exit(3)
	}

	for _, e := range urls {
		var sourcesBak []string
		var completedSuccessfully = true
		_, _ = fmt.Fprintf(os.Stdout, "Getting sources from url: %s\n", e)
		sources, err := getScriptSrc(e, opt.Method, opt.Headers, opt.Insecure, opt.Timeout)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Couldn't get sources from url: %v\n", err)
		}

		if opt.Complete {
			_, _ = fmt.Fprintf(os.Stdout, "Completing URLs\n")
			sourcesBak = sources
			sources, err = completeUrls(sources, e)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Couldn't complete URLs: %v\n", err)
				sources = sourcesBak
				completedSuccessfully = false
			}
		}

		if opt.Resolve && opt.Complete {
			if completedSuccessfully {
				sourcesBak = sources
				sources, err = resolveUrlsAndGrepScAddr(sources)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Couldn't resolve URLs: %v\n", err)
					sources = sourcesBak
				}
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "Couldn't resolve URLs\n")
			}
		} else if opt.Resolve {
			_, _ = fmt.Fprintf(os.Stderr, "Resolve can only be used in combination with -complete\n")
		}

		allSources = append(allSources, sources...)

		if opt.OutputFile != "" {
			allSources = append(allSources, sources...)
		}

	}

	// Save to file
	if opt.OutputFile != "" {
		_, _ = fmt.Fprintf(os.Stdout, "Saving output to file: %s\n", opt.OutputFile)
		err := saveToFile(allSources, opt.OutputFile)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Couldn't save to output file: %v\n", err)
		}
	}

	return allSources, nil
}

func saveToFile(sources []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Couldn't close output file: %v\n", err)
		}
	}(file)

	w := bufio.NewWriter(file)
	for _, line := range sources {
		_, _ = fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func getScriptSrc(url string, method string, headers []string, insecure bool, timeout int) ([]string, error) {
	// Request the HTML page.
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return []string{}, err
	}

	for _, d := range headers {
		values := strings.Split(d, ":")
		if len(values) == 2 {
			_, _ = fmt.Fprintf(os.Stderr, "Setting header: %s: %s\n", values[0], values[1])
			req.Header.Set(values[0], values[1])
		}
	}

	tr := &http.Transport{
		ResponseHeaderTimeout: time.Duration(timeout) * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: insecure},
	}

	var client = &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: tr,
	}

	res, err := client.Do(req)
	if err != nil {
		return []string{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Couldn't close response body: %v\n", err)
		}
	}(res.Body)
	if res.StatusCode != 200 {
		_, _ = fmt.Fprintf(os.Stderr, "Url returned an %d instead of 200\n", res.StatusCode)
		return nil, nil
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var sources []string

	// Find the script tags, and get the src
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		dsrc, _ := s.Attr("data-src")
		if src != "" {
			sources = append(sources, src)
		}
		if dsrc != "" {
			sources = append(sources, dsrc)
		}
	})

	return sources, nil
}

func readLines(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func resolveUrlsAndGrepScAddr(s []string) ([]string, error) {
	patternRegex := `0x[a-fA-F\d]{40}`
	smartContractAddresses := make(map[string]bool)

	for i := len(s) - 1; i >= 0; i-- {
		resp, err := http.Get(s[i])
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 && resp.StatusCode != 304 {
			s = append(s[:i], s[i+1:]...)
		}

		// find all sc addresses in the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		scAddresses := regexp.MustCompile(patternRegex).FindAllString(string(body), -1)
		for _, scAddress := range scAddresses {
			smartContractAddresses[scAddress] = true
		}
	}

	// convert map to slice, and ignore 0x00 and 0xff
	addresses := make([]string, 0, len(smartContractAddresses))
	for k := range smartContractAddresses {
		if k != "0x0000000000000000000000000000000000000000" && k != "0xffffffffffffffffffffffffffffffffffffffff" {
			addresses = append(addresses, k)
		}
	}

	// sort it
	sort.Strings(addresses)

	return addresses, nil
}

func completeUrls(s []string, mainUrl string) ([]string, error) {
	u, err := url.Parse(mainUrl)
	if err != nil {
		return nil, err
	}

	for i := range s {
		if strings.HasPrefix(s[i], "//") {
			s[i] = u.Scheme + ":" + s[i]
		} else if strings.HasPrefix(s[i], "/") && string(s[i][1]) != "/" {
			s[i] = u.Scheme + "://" + u.Host + s[i]
		} else if !strings.HasPrefix(s[i], "http://") && !strings.HasPrefix(s[i], "https://") {
			s[i] = u.Scheme + "://" + u.Host + u.Path + "/" + s[i]
		}
	}
	return s, nil
}
