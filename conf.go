package main

// Conf is the program configuration.
type Conf struct {
	Domains []string `conf:"domains"`
	Servers []DomainConf
}

// DomainConf is a domain configuration.
type DomainConf struct {
	Listen string `conf:"listen"`
	Name   string `conf:"name"`

	Certs CertsConf  `conf:"certs"`
	Repos []RepoConf `conf:"repo"`
}

// CertsConf contains the certificate configuration for SSL.
type CertsConf struct {
	Public  string `conf:"public"`
	Private string `conf:"private"`
}

// RepoConf is an individual repo configuration.
//
// The domain name and module name are combined as:
//	domain '/' module
//
// For example the url: go.company.corp/find/this can be split into:
//	domain= go.company.corp
//	module= find/this
type RepoConf struct {
	Module []string `conf:"module"`
	Repo   string   `conf:"repo"`
	VCS    string   `conf:"vcs"`
}
