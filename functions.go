package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/nofeaturesonlybugs/conf"
	"github.com/nofeaturesonlybugs/errors"
)

// LoadFlags loads command line arguments.
//
// As a special case if no flags are provided a help message is printed and the
// program exits.
func LoadFlags() {
	usage := struct {
		Conf    string
		Help    string
		Serve   string
		Test    string
		Version string
	}{
		Conf:    "Set the configuration directory.",
		Help:    "Print help information.",
		Serve:   "Run the server to serve requests from Go tools.",
		Test:    "Test configuration.",
		Version: "Print version information.",
	}
	flag.StringVar(&Flags.Paths.Conf, "c", "", usage.Conf)
	flag.StringVar(&Flags.Paths.Conf, "conf", "", usage.Conf)
	flag.BoolVar(&Flags.Help, "h", false, usage.Help)
	flag.BoolVar(&Flags.Help, "help", false, usage.Help)
	flag.BoolVar(&Flags.Serve, "s", false, usage.Serve)
	flag.BoolVar(&Flags.Serve, "serve", false, usage.Serve)
	flag.BoolVar(&Flags.Test, "t", false, usage.Test)
	flag.BoolVar(&Flags.Test, "test", false, usage.Test)
	flag.BoolVar(&Flags.Version, "v", false, usage.Version)
	flag.BoolVar(&Flags.Version, "version", false, usage.Version)
	flag.Parse()
	if Flags.Help || (flag.NArg() == 0 && flag.NFlag() == 0) {
		flag.PrintDefaults()
		os.Exit(0)
	}
}

// LoadPaths loads remaining paths and files into Flags.
func LoadPaths() error {
	var err error
	//
	// Application home.
	if Flags.Paths.Home, err = os.Executable(); err != nil {
		return errors.Go(err)
	}
	Flags.Paths.Home = filepath.Dir(Flags.Paths.Home)
	//
	// Configuration path.
	if Flags.Paths.Conf == "" {
		Flags.Paths.Conf = filepath.Join(Flags.Paths.Home, "conf")
	}
	if stat, err := os.Stat(Flags.Paths.Conf); err != nil {
		return errors.Go(err)
	} else if !stat.IsDir() {
		errors.Errorf("Configuration directory does not exist @ %v", Flags.Paths.Conf)
	}
	//
	// Hostname for host-specific configurations.
	host, err := os.Hostname()
	if err != nil {
		return errors.Go(err)
	}
	//
	// Look for hostname.ini, fallback to conf.ini within conf directory.
	files := []string{strings.ToLower(host) + ".ini", "conf.ini"}
	for _, file := range files {
		file = filepath.Join(Flags.Paths.Conf, file)
		if stat, err := os.Stat(file); err == nil {
			if !stat.IsDir() {
				Flags.Files.Conf = file
				break
			}
		}
	}
	if Flags.Files.Conf == "" {
		return errors.Errorf("Config file not found; searched both %v in %v", strings.Join(files, ", "), Flags.Paths.Conf)
	}
	//
	return nil
}

// LoadConfig creates and loads a configuration object from a given config file name.
func LoadConfig(file string) (*Conf, error) {
	rv := &Conf{}
	if c, err := conf.File(file); err != nil {
		return nil, errors.Go(err)
	} else if err = c.FillByTag("conf", rv); err != nil {
		return nil, errors.Go(err)
	}
	//
	// Domain configurations are in same path as main configuration.
	path := filepath.Dir(file)
	//
	// Load each domain.
	for _, domain := range rv.Domains {
		d := DomainConf{}
		if c, err := conf.File(filepath.Join(path, domain)); err != nil {
			return nil, errors.Go(err)
		} else if err = c.FillByTag("conf", &d); err != nil {
			return nil, errors.Go(err)
		}
		rv.Servers = append(rv.Servers, d)
	}
	return rv, nil
}
