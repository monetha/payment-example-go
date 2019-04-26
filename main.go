package main

import (
	"context"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/monetha/payment-go-sdk/wallet"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	eth "github.com/monetha/go-ethereum"
	"github.com/monetha/payment-go-sdk/processor"
)

func main() {
	log.SetOutput(os.Stderr)
	ctx := createCtrlCContext()

	backendURL := "https://ropsten.infura.io/v3/9341cea07e634c21be9d5a5ccb892db5"

	client, err := ethclient.Dial(backendURL)
	if err != nil {
		log.Printf("error: %v", err.Error())
		return
	}

	e := eth.New(client, log.Printf)

	// processorContractAddress address of the Merchant's smart contract
	processorContractAddress := common.HexToAddress("0x35D6708FD36DCb902adce5D9d6ABeB4838318554")

	// Initialize Customer
	customerAccount, err := crypto.HexToECDSA("d4266a34672534b063d4d7911a7450d93f35075ef8583b3c96a408285fc21a38") // customerAccount public address: 0x8c77F5BA864718f098F83114efEC7180649afB85
	if err != nil {
		log.Printf("error: %v", err.Error())
	}
	log.Printf("info: customer address - %v", crypto.PubkeyToAddress(customerAccount.PublicKey).String())

	customerAddress := crypto.PubkeyToAddress(customerAccount.PublicKey)

	// Initialize Merchant
	merchantAccount, err := crypto.HexToECDSA("ad6c05a5d77f993cf6c23eab52c6b9db7894dabf5647df63155d3e66280c2dc3") // merchantKey public address: 0xdF8c3E2c8506F67705acB0a4dCa28Cf44934B511
	if err != nil {
		log.Printf("error: %v", err.Error())
	}
	log.Printf("info: merchant address - %v", crypto.PubkeyToAddress(merchantAccount.PublicKey).String())

	// Flow 1.
	// Customer initializes an Order

	log.Print("info: starting example flow: customer initializes the order")

	customerSession := e.NewSession(customerAccount)
	p := processor.NewProcessor(customerSession, processorContractAddress)

	// Order metadata
	orderID := big.NewInt(time.Now().Unix())
	price := math.Exp(big.NewInt(1*10), big.NewInt(15)) // price = 0.001 ETH
	tokenAddress := common.Address{} // Payment is done in ETH

	log.Print("info: Customer adds order")
	txHash, err := p.AddOrder(ctx, orderID, price, customerAddress, tokenAddress, big.NewInt(0))
	receipt, err := e.WaitForTxReceipt(ctx, txHash)
	if err != nil {
		log.Printf("error: could not create order")
	}

	if receipt.Status == types.ReceiptStatusSuccessful {

		log.Print("info: Customer pays for the order")
		txHash, err = p.SecurePay(ctx, orderID, price)
		if err != nil {
			log.Printf("error: customer could not submit order confirmation %v", err)
			return
		}
		_, err := e.WaitForTxReceipt(ctx, txHash)
		if err != nil {
			log.Printf("error: customer could not pay for the orders")
			return
		}
	}

	// Merchant received information from Customer that order is created

	// Merchant validates that Order is paid


	// TODO: ask user to sign the text message to proof the public address is his address

	// Initialize a new session with a Merchant's private key
	merchantSession := e.NewSession(merchantAccount)

	// Initialize the payment Processor
	p = processor.NewProcessor(merchantSession, processorContractAddress)

	dealContent := `
	{ 
		"description": "Service Agreement No. 201901-987",
		"customer_full_name": "John Smith",
		"phone_number": "+1-202-555-0116"
		
	}`
	dealHash := new(big.Int).SetBytes(crypto.Keccak256Hash([]byte(dealContent)).Bytes())

	log.Print("info: Merchant processes the payment")
	// Execute ProcessPayment to transfer order funds to Merchant's Wallet
	txHash, err = p.ProcessPayment(ctx, orderID, dealHash)
	if err != nil {
		log.Printf("error: %v", err)
	}

	receipt, err = e.WaitForTxReceipt(ctx, txHash)
	if err != nil {
		log.Printf("error: could not process order")
	}

	merchantWalletAddress, err := p.ContractHandler.MerchantWallet(nil)
	if err != nil {
		log.Printf("error: could not retrieve merchant wallet address from PaymentProcessor contract")
	}

	log.Print("info: Merchant transfers funds from his wallet contract to Merchant's Address")
	wallet := wallet.NewWallet(merchantSession, merchantWalletAddress)

	txHash, err = wallet.WithdrawAllTo(ctx, crypto.PubkeyToAddress(merchantAccount.PublicKey), big.NewInt(0))
	if err != nil {
		log.Printf("error: %v", err)
	}

	receipt, err = e.WaitForTxReceipt(ctx, txHash)
	if err != nil {
		log.Printf("error: %v", err)
	}

	log.Print("info: ending example flow: customer initializes the order")

	// Flow 2.
	// Merchant initializes an Order

	log.Print("info: starting example flow: customer initializes the order")

	// TODO: implement

	log.Print("info: ending example flow: customer initializes the order")

	defer log.Print("info: done.")
}

func createCtrlCContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
		<-sigChan
		log.Println("got interrupt signal")
	}()

	return ctx
}
