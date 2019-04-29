## Intro

This repository contains example of using [Monetha Payment Go SDK](https://github.com/monetha/payment-go-sdk). Examples show main capabilities of SDK in order to achieve a decentralized payment via Monetha Payment Gateway

## Build & Run

### Prerequisites

1. Make sure you have [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) installed.
1. Install [Go 1.12](https://golang.org/dl/)
1. Setup `$GOPATH` environment variable as described [here](https://github.com/golang/go/wiki/SettingGOPATH).
1. Clone the repository:
    ```bash
    mkdir -p $GOPATH/src/github.com/monetha
    cd $GOPATH/src/github.com/monetha
    git clone git@github.com:monetha/payment-example.git
    cd reputation-go-sdk
    ```

**Note**: You can skip steps 2-3 on Linux and use the official docker image for Go after step 4 to build the project:

```bash
docker run -it --rm \
  -v "$PWD":/go/src/github.com/monetha/payment-example \
  -w /go/src/github.com/monetha/payment-example \
  golang:1.12 \
  /bin/bash
```

### Build

Install dependencies:

    make dependencies

Once the dependencies are installed, run 

    make cmd

to build the full suite of utilities. After the executable files are built, they can be found in the directory `./bin/`.

### Running examples