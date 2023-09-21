# GoLang Helpers

Go module for common helpers/bits of code that will be shared
between projects.

Use via the following:

``` bash
go get github.com/PossibleLlama/golang-helpers
```

Testing via the following:

``` bash
gotestsum --packages="./..." -- -coverprofile="./coverage.out" -count=10 -race
go tool cover -html="./coverage.out"
```

## Logging

This is a wrapper around zap that sets up a series of fields
to logged consistently between projects.

### Intended use

During initialization, if setting the `scmLinkToRepo` parameter
to a link (in this repo's case `https://github.com/PossibleLlama/golang-helpers`)
all logs will have a direct link to the repository.
You can use this in combination with setting the `version`
parameter to a git hash to allow those links to always go to the
commit that the application was built from, so you can travel to
the correct point in time as well as location.

The caller is expected to expose the version and scmLink as
build time variables, which then accept these being passed in
during the build.
If the caller is using Github actions, you can use the following
variables.

- `${{ github.repository }}` - Will give you `PossibleLlama/golang-helpers`.
  You'll need to prepend `https://github.com/` to that.
- `${{ github.sha }}` - Will give you the full git sha.
