package providers

import (
	"os"
	"os/user"

	"github.com/pacific/output"
	"gopkg.in/yaml.v3"
)

type Config struct {
	IncludeSubdomains bool
	Providers         []string
	Blacklist         []string
	ApiKeys           API
}

type API struct {
	Alienvault     []string `yaml:"alienvault"`
	Commoncrawl    []string `yaml:"commoncrawl"`
	Grayhatwarfare []string `yaml:"grayhatwarfare"`
	Grepapp        []string `yaml:"grepapp"`
	Hybridanalysis []string `yaml:"hybridanalysis"`
	Packettotal    []string `yaml:"packettotal"`
	Virustotal     []string `yaml:"virustotal"`
	Wayback        []string `yaml:"wayback"`
	Urlscan        []string `yaml:"urlscan"`
}

type ApiConfig struct {
	Api API
}

type Provider interface {
	Fetch(string, chan<- string)
}

func getFilename() string {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	username := currentUser.Username
	return "/home/" + username + "/.config/pacific/config.yaml"

}

func (C *ApiConfig) ReadApiConfig(confile string) {
	var fname string = getFilename()

	if len(confile) > 0 {
		fname = confile
	}

	data, err := os.ReadFile(fname)
	if err != nil {
		output.Info("Could not found configuration file in " + fname)
	}

	var config API
	yaml.Unmarshal(data, &config)
	C.Api = config
}
