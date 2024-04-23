module github.com/alexrios/endpoints

go 1.21

require (
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.4
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/afero v1.3.2
	github.com/stretchr/testify v1.4.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20190422165155-953cdadca894 // indirect
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)

// This version (and below) does not support "go install" command. Try with v0.5.0 of greater.
retract v0.4.1
