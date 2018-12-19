module github.com/manifoldco/healthz

require (
	github.com/fatih/color v1.7.0 // indirect
	github.com/gogo/protobuf v1.2.0 // indirect
	github.com/golang/mock v1.2.0 // indirect
	github.com/golangci/gocyclo v0.0.0-20180528144436-0a533e8fa43d // indirect
	github.com/golangci/golangci-lint v1.12.3
	github.com/golangci/revgrep v0.0.0-20180812185044-276a5c0a1039 // indirect
	github.com/golangci/tools v0.0.0-20181110070903-2cefd77fef9b // indirect
	github.com/golangci/unparam v0.0.0-20180902115109-7ad9dbcccc16 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/nbutton23/zxcvbn-go v0.0.0-20180912185939-ae427f1e4c1d // indirect
	github.com/onsi/ginkgo v1.7.0 // indirect
	github.com/onsi/gomega v1.4.3 // indirect
	github.com/shurcooL/go v0.0.0-20181215222900-0143a8f55f04 // indirect
	github.com/sirupsen/logrus v1.2.0 // indirect
	github.com/spf13/afero v1.2.0 // indirect
	github.com/spf13/cobra v0.0.3 // indirect
	github.com/spf13/viper v1.3.1 // indirect
	golang.org/x/net v0.0.0-20181217023233-e147a9138326 // indirect
	golang.org/x/sync v0.0.0-20181108010431-42b317875d0f // indirect
	golang.org/x/sys v0.0.0-20181218192612-074acd46bca6 // indirect
	golang.org/x/tools v0.0.0-20181218204010-d4971274fe38 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	sourcegraph.com/sourcegraph/go-diff v0.5.0 // indirect
	sourcegraph.com/sqs/pbtypes v1.0.0 // indirect
)

// This version of kingpin is incompatible with the released version of
// gometalinter until the next release of gometalinter, and possibly until it
// has go module support, we'll need this exclude, and perhaps some more.
//
// After that point, we should be able to remove it.
exclude gopkg.in/alecthomas/kingpin.v3-unstable v3.0.0-20180810215634-df19058c872c
