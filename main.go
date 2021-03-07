package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/nofeaturesonlybugs/routines"
)

var (
	// Version is the application version.
	Version = ""
	// Date is the application build date.
	Date = ""
	// GoVersion is the Go version.
	GoVersion = ""
)

// Flags is a container for command line flags.
var Flags struct {
	Files struct {
		Conf string
	}
	Paths struct {
		Home string
		Conf string
	}
	Help    bool
	Serve   bool
	Test    bool
	Version bool
}

// fatal tests if err is not nil and exits the app if it is.
func fatal(err error) {
	if err != nil {
		fmt.Printf("fatal error: %v\n", err.Error())
		os.Exit(0x01)
	}
}

func main() {
	var conf *Conf
	var rtns routines.Routines
	var err error
	//
	LoadFlags()
	if Flags.Serve || Flags.Test {
		err = LoadPaths()
		fatal(err)
		conf, err = LoadConfig(Flags.Files.Conf)
		fatal(err)
		if Flags.Test {
			fmt.Println("config ok")
			os.Exit(0)
		}
	} else if Flags.Version {
		exe, _ := os.Executable()
		exe = filepath.Base(exe)
		if Version == "" {
			Version = "n/a"
		}
		if Date == "" {
			Date = "n/a"
		}
		if GoVersion == "" {
			GoVersion = "n/a"
		}
		fmt.Printf("%v %v\n", exe, Version)
		fmt.Printf("\tbuilt %v\n\twith %v\n", Date, GoVersion)
		os.Exit(0)
	}
	//
	rtns = routines.NewRoutines()
	defer fmt.Println("Clean stop.") //TODO log better
	defer rtns.Wait()
	defer rtns.Stop()
	//
	if Flags.Serve {
		var vanity *VanityServer
		for _, server := range conf.Servers {
			if vanity, err = NewVanityServer(server); err != nil {
				fmt.Printf("creating vanity server: %v", err.Error()) // TODO log better
				return
			} else if err = vanity.Start(rtns); err != nil {
				fmt.Printf("starting vanity server: %v", err.Error()) // TODO log better
				return
			}
		}
	}
	fmt.Println("Running.") //TODO log better
	//
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh)
	<-sigCh
}
