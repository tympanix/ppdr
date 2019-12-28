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
    ?       github.com/tympanix/master-2019 [no test files]
    ?       github.com/tympanix/master-2019/debug   [no test files]
    ok      github.com/tympanix/master-2019/ltl     (cached)
    ok      github.com/tympanix/master-2019/ltl/parser      (cached)
    ok      github.com/tympanix/master-2019/ltl/scanner     (cached)
    ?       github.com/tympanix/master-2019/ltl/scanner/token       [no test files]
    ok      github.com/tympanix/master-2019/repo    (cached)
    ok      github.com/tympanix/master-2019/systems (cached)
    ?       github.com/tympanix/master-2019/systems/ba      [no test files]
    ok      github.com/tympanix/master-2019/systems/gnba    (cached)
    ?       github.com/tympanix/master-2019/systems/mock/ts [no test files]
    ok      github.com/tympanix/master-2019/systems/nba     (cached)
    ok      github.com/tympanix/master-2019/systems/product (cached)
    ?       github.com/tympanix/master-2019/systems/ts      [no test files]
```