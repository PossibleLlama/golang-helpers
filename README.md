# GoLang Helpers

Go module for common helpers/bits of code that will be shared between projects.

Use via the following:

``` bash
go get github.com/PossibleLlama/golang-helpers
```

Testing via the following:

``` bash
gotestsum --packages="./..." -- -coverprofile="./coverage.out" -race
go tool cover -html="./coverage.out"
```
