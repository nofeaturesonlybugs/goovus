package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"path"
	"strings"

	"github.com/nofeaturesonlybugs/errors"
	"github.com/nofeaturesonlybugs/routines"
)

// VanityServer is an http server that serves vanity urls for a domain.
type VanityServer struct {
	routines.Service
	//
	conf  DomainConf
	httpd *http.Server
}

// NewVanityServer creates a new Server type from the configuration.
func NewVanityServer(conf DomainConf) (*VanityServer, error) {
	var err error
	var cert tls.Certificate
	var tlsConfig *tls.Config
	//
	// Validate host:port configuration.
	_, _, err = net.SplitHostPort(conf.Listen)
	if err != nil {
		return nil, errors.Errorf("Invalid host:port value @ %v", conf.Listen)
	}
	//
	// Check for TLS
	if conf.Certs.Private != "" && conf.Certs.Public != "" {
		if cert, err = tls.LoadX509KeyPair(conf.Certs.Public, conf.Certs.Private); err != nil {
			return nil, errors.Go(err)
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}
	//
	// Map our modules configuration for faster lookups in the http handler.
	modules := map[string]VanityTemplateData{}
	for _, repo := range conf.Repos {
		for _, module := range repo.Module {
			modules["/"+module] = VanityTemplateData{
				Root:    path.Join(conf.Name, module),
				RepoURL: repo.Repo,
				VCS:     repo.VCS,
			}
		}
	}
	//
	// Our http handler.
	handler := func(w http.ResponseWriter, r *http.Request) {
		var module VanityTemplateData
		var ok bool
		var err error
		//
		// Make a local copy.
		modules := modules
		if !strings.Contains(r.URL.RawQuery, "go-get=1") {
			http.NotFound(w, r)
			return
		} else if module, ok = modules[r.URL.Path]; !ok {
			http.NotFound(w, r)
			return
		}
		if err = VanityTemplate.Execute(w, module); err != nil {
			fmt.Printf("error during render: %v", err.Error()) // TODO Logging facility would be better.
		}
	}
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handler))
	//
	rv := &VanityServer{
		conf: conf,
		httpd: &http.Server{
			Addr:      conf.Listen,
			Handler:   mux,
			TLSConfig: tlsConfig,
		},
	}
	rv.Service = routines.NewService(rv.start)
	//
	return rv, nil
}

// start starts the httpd listener.
func (me *VanityServer) start(rtns routines.Routines) error {
	var listener net.Listener
	var err error
	if me.httpd.TLSConfig == nil {
		if listener, err = net.Listen("tcp", me.conf.Listen); err != nil {
			return errors.Go(err)
		}
	} else {
		if listener, err = tls.Listen("tcp", me.conf.Listen, me.httpd.TLSConfig); err != nil {
			return errors.Go(err)
		}
	}
	fmt.Printf("Listening on %v for %v\n", me.conf.Listen, me.conf.Name) //TODO log better
	//
	rtns.Go(func() {
		defer fmt.Printf("Stopped %v for %v\n", me.conf.Listen, me.conf.Name) // TODO log better
		me.httpd.Serve(listener)
	})
	rtns.Go(func() {
		<-rtns.Done()
		me.httpd.Shutdown(context.Background())
	})
	//
	return nil
}
