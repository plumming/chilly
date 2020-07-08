module github.com/plumming/chilly

go 1.13

require (
	github.com/fatih/color v1.9.0
	github.com/ghodss/yaml v1.0.0
	github.com/hashicorp/go-version v1.2.1
	github.com/jenkins-x/jx-logging v0.0.11
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/maxbrunsfeld/counterfeiter/v6 v6.2.3
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	github.com/xeonx/timeago v1.0.0-rc4
	golang.org/x/tools v0.0.0-20200515010526-7d3b6ebf133d // indirect
	gopkg.in/AlecAivazis/survey.v1 v1.8.8
	k8s.io/client-go v0.18.5
	k8s.io/utils v0.0.0-20200327001022-6496210b90e8 // indirect
)

replace k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190718092204-1918f780cd63
