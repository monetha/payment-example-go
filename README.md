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
    cd payment-example
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

After dependencies are installed you can start the examples

### Running examples

All examples are run with Test Customer and Merchant addresses. We've provided private keys to those address for simplicity of onboarding and showing how SDK works. We call everyone who are executing examples not to withdraw Ether funds from those addresses.

| Name | Ropsten Address |
| ---- | --------------- |
| PaymentProcessor Contract Address | [0x35D6708FD36DCb902adce5D9d6ABeB4838318554](https://ropsten.etherscan.io/address/0x35D6708FD36DCb902adce5D9d6ABeB4838318554) |
| Merchant Address | [0x8c77F5BA864718f098F83114efEC7180649afB85](https://ropsten.etherscan.io/address/0x8c77F5BA864718f098F83114efEC7180649afB85) |
| Customer Address | [0xdF8c3E2c8506F67705acB0a4dCa28Cf44934B511](https://ropsten.etherscan.io/address/0xdF8c3E2c8506F67705acB0a4dCa28Cf44934B511) |

We are using Infura.io JSON RPC in the provided examples in order to execute transaction on a chain. Feel free to change `backendURL` variable to any JSON RPC node you want.

```golang
backendURL := "https://ropsten.infura.io/v3/7c0fc2888a824c62a3651fd446c8f989"
```

**Note**: Currently all PaymentProcessor contract instances are provided by Monetha. Contact [team@monetha.io](mailto:team@monetha.io) in case if you would like to have your own instance deployed. 

#### Example flow where Customer initiates the purchase

```bash
go run order_initiated_by_customer/main.go
```

Short description of what is going on

- Customer initiates the order
- Customer pays for the Order
- Merchant verifies if Order was paid
- Merchant processes the payment after service was provided
- Merchant withdraws funds from MerchantWallet to his Address

### Example flow where Merchant initiates the purchase

```bash
go run order_initiated_by_merchant/main.go
```

Short description of what is going on

- Merchant initiates the order and provides Customer with information to pay
- Customer pays for the Order
- Merchant verifies if Order was paid
- Merchant processes the payment after service was provided
- Merchant withdraws funds from MerchantWallet to his Address