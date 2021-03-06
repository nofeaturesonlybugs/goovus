`goovus` serves vanity URLs to Go tools.

## What's In A Name?

* `go`  
  Made for Go.
* `o`  
  Open as in open source.
* `vus`  
  `v`anity `u`rl `s`erver.

`go + o + vus` gives `goovus`.

## Quick Note
We'll be using `go.company.corp` a lot in this document.  However you can use any internal, private, or public domain as long as you can install `goovus` on the relevant machine.

## What Does It Do?

It serves vanity URLs to go tools.  Let's say you import:
```go
import (
    "go.company.corp/libA"
    "go.company.corp/libB"
    "go.company.corp/libB/pkg"
)
```

Commands such as `go get` will attempt to communicate with a web server at `go.company.corp`.  That web server can tell the Go tools where the source code repos are for each of those modules.

Essentially `"go.company.corp/libA"` can be mapped to an internal source code repo; likewise with `"go.company.corp/libB"`.  Or any other vanity URL you can conjure up.

## Go Modules and Private Repos
Let's assume that `go.company.corp/libA` and `go.company.corp/libB` are private Go modules.  Without a vanity URL server you have to import them as:
```go
import (
    "go.company.corp/libA.git"
    "go.company.corp/libB.git"
    "go.company.corp/libB.git/pkg"
)
```
For private repos a vanity URL server can remove `.git` from appearing in import paths.

I think it's worth it.

## Not a Go Module Proxy
`goovus` is **NOT** a Go module proxy.  It does not build a module cache.  It does not fetch things from the internet.  It does not fetch your private repos and serve them as if it was a proxy.

In fact it's ideal when you have private repos and you DON'T want to configure a Go module proxy.

## Building goovus

```bash
$ go build
```

But if you want `goovus -v` to print useful information build it with the build script.
```bash
$ ./build.sh
```

## Configuration
Add domains used for vanity URLs to GOPRIVATE.

Continuing with `libA` and `libB` you'd set `GOPRIVATE=go.company.corp` in your environment.  This stops `go get` or `go mod tidy` from searching the internet for any modules beginning with `go.company.corp`.

`goovus` uses a main configuration file and then one additional configuration file per vanity domain.

The default configuration file directory is `$EXEHOME/conf` but you can set it to any directory with the `-c` or `-conf` flags.

```ini
# $EXEHOME/conf/conf.ini ~or~ $EXEHOME/conf/hostname.ini
#
# Main configuration file.
# Create a domains= line for each vanity domain to serve.  The value is the name of a domain.ini
# file in the same directory as this file.
#
# We are serving vanity URLs for go.company.corp and decided to name the domain ini file
# go-company-corp.ini and that file exists in the same directory as this one.
domains = go-company-corp.ini
```

```ini
# $EXEHOME/conf/go-company-corp.ini
# This is the ini file for go.company.corp.

# Set our bind network address and the vanity domain we are serving.
listen = 0.0.0.0:443
name = go.company.corp

[certs]
# If public and private are set to paths on disk they are considered
# public/private keys and a TLS listener is created at the "host:port"
# value above.
public = /ssl/certs/go.company.corp.pem
private = /ssl/certs/go.company.corp.pem

# Next we need a [repo] section for EACH go module we are serving vanity URLs.
# In our example these are libA and libB.
#
# Note that the value for "name =" is prefixed to each repo.module and they are
# separated with a slash.
#
# In other words "go.company.corp/" is prefixed to both "libA" and "libB" to create the
# complete name the go tools are searching.
#

[repo]
module = libA
repo = ssh://gitserver/libA
vcs = git

[repo]
# As a bonus example let's say libB has v2, v3, and v4 all served out of the same repo.
module = libB
module = libB/v2
module = libB/v3
module = libB/v4
repo = ssh://gitserver/libB
vcs = git
```

Run `goovus` with:
```bash
$ goovus -s
```
or
```bash
$ goovus -serve
```
