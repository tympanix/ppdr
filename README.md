# Setting up

## Prerequisites

Go is installed.
* https://golang.org/doc/install

`GOPATH` and `PATH` is set correctly. 
* Can be tested by following the steps on https://golang.org/doc/install#testing

## Steps

The following instructions are intended for unix/linux systems.

Navigate to the root of the project.

Install the stringer library.
```bash
    go get golang.org/x/tools/cmd/stringer
```

Use it to generate required methods from root of the project.
```bash
    go generate ./...
```

# Test
From root of the project to execute all tests.
```bash
    go test ./...
    ?       github.com/tympanix/ppdr [no test files]
    ?       github.com/tympanix/ppdr/debug   [no test files]
    ok      github.com/tympanix/ppdr/ltl     (cached)
    ok      github.com/tympanix/ppdr/ltl/parser      (cached)
    ok      github.com/tympanix/ppdr/ltl/scanner     (cached)
    ?       github.com/tympanix/ppdr/ltl/scanner/token       [no test files]
    ok      github.com/tympanix/ppdr/repo    (cached)
    ok      github.com/tympanix/ppdr/systems (cached)
    ?       github.com/tympanix/ppdr/systems/ba      [no test files]
    ok      github.com/tympanix/ppdr/systems/gnba    (cached)
    ?       github.com/tympanix/ppdr/systems/mock/ts [no test files]
    ok      github.com/tympanix/ppdr/systems/nba     (cached)
    ok      github.com/tympanix/ppdr/systems/product (cached)
    ?       github.com/tympanix/ppdr/systems/ts      [no test files]
```