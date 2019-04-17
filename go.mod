module github.com/manifoldco/healthz

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/alecthomas/gometalinter v3.0.0+incompatible
	github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf // indirect
	github.com/client9/misspell v0.3.4
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/shlex v0.0.0-20181106134648-c34317bd91bf // indirect
	github.com/gordonklaus/ineffassign v0.0.0-20180909121442-1003c8bd00dc
	github.com/kr/pretty v0.1.0 // indirect
	github.com/nicksnyder/go-i18n v1.10.0 // indirect
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/stretchr/testify v1.3.0 // indirect
	github.com/tsenart/deadcode v0.0.0-20160724212837-210d2dc333e9
	golang.org/x/lint v0.0.0-20181217174547-8f45f776aaf1
	golang.org/x/tools v0.0.0-20190206175659-4892ae6946ab // indirect
	gopkg.in/alecthomas/kingpin.v3-unstable v3.0.0-20171010053543-63abe20a23e2 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)

// This version of kingpin is incompatible with the released version of
// gometalinter until the next release of gometalinter, and possibly until it
// has go module support, we'll need this exclude, and perhaps some more.
//
// After that point, we should be able to remove it.
exclude gopkg.in/alecthomas/kingpin.v3-unstable v3.0.0-20180810215634-df19058c872c
