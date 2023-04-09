package main

import (
	"bufio"
	"flag"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pacific/output"
	"github.com/pacific/providers"
)

func main() {
	log.SetFlags(0)
	var domains []string = []string{}
	results := make(chan string)
	var f *os.File = os.Stdout

	const logo string = `
    ____  ___   __________________________
   / __ \/   | / ____/  _/ ____/  _/ ____/
  / /_/ / /| |/ /    / // /_   / // /     
 / ____/ ___ / /____/ // __/ _/ // /___   
/_/   /_/  |_\____/___/_/   /___/\____/   
                                          
`

	target := flag.String("d", "", "Target domain")
	confile := flag.String("cf", "", "Config file path")
	outfile := flag.String("o", "", "Filename to write results to")
	blackList := flag.String("b", "", "extensions to skip, ex: ttf,woff,svg,png,jpg")
	provider := flag.String("providers", "alienvault,commoncrawl,grayhatwarfare,hybridanalysis,wayback,urlscan", "Active providers")
	flag.Parse()
	println(logo)

	uparse, _ := url.Parse(*target)
	if len(uparse.Scheme) > 0 {
		target = &uparse.Host
	}
	if len(*target) <= 0 {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			domains = append(domains, s.Text())
		}
	}
	domains = append(domains, *target)
	t1 := time.Now()
	println("Target: " + *target + "\nTime: " + t1.Format(time.RFC850))

	var apiConfig providers.ApiConfig
	apiConfig.ReadApiConfig(*confile)

	var config providers.Config

	config.Providers = strings.Split(*provider, ",")
	config.Blacklist = strings.Split(*blackList, ",")
	config.ApiKeys = apiConfig.Api

	if *outfile != "" {
		f, _ = os.OpenFile(*outfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	var providerList []providers.Provider
	for _, p := range config.Providers {
		switch p {
		case "alienvault":
			otx := providers.NewAlienvault(&config)
			providerList = append(providerList, otx)
		case "commoncrawl":
			ccw := providers.NewCommonCrawl(&config)
			providerList = append(providerList, ccw)
		case "grayhatwarfare":
			ghat := providers.NewGrayHatWarfare(&config)
			providerList = append(providerList, ghat)
		case "hybridanalysis":
			ha := providers.NewHybridAnalysis(&config)
			providerList = append(providerList, ha)
		case "wayback":
			wayback := providers.NewWayback(&config)
			providerList = append(providerList, wayback)
		case "urlscan":
			urlscan := providers.NewUrlscan(&config)
			providerList = append(providerList, urlscan)
		default:
			output.Err(p + " is not a valid provider.")
		}
	}

	outwg := &sync.WaitGroup{}
	outwg.Add(1)
	go func() {
		output.Out(f, results, config.Blacklist)
		outwg.Done()
	}()

	wg := &sync.WaitGroup{}
	var i = 0
	for _, domain := range domains {
		domain := domain
		wg.Add(len(providerList))
		i++
		for _, provider := range providerList {
			go func(provider providers.Provider) {
				defer wg.Done()
				provider.Fetch(domain, results)
			}(provider)
		}
		if i >= 40 {
			i = 0
			wg.Wait()
		}
	}

	wg.Wait()
	close(results)
	outwg.Wait()

}
