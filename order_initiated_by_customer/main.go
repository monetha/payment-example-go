package main

import (
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	eth "github.com/monetha/go-ethereum"
	"github.com/monetha/payment-example-go/utils"
	"github.com/monetha/payment-go-sdk/processor"
	"github.com/monetha/payment-go-sdk/wallet"
)

func main() {
	var p *processor.Processor
	var customerSession, merchantSession *eth.Session
	var orderID, price *big.Int
	var txHash common.Hash

	// paymentProcessorContractAddress address of the Merchant's PaymnetProcessor smart contract deployed and provided by Monetha
	var paymentProcessorContractAddress = "0x35D6708FD36DCb902adce5D9d6ABeB4838318554"

	// backendURL JSON RPC url for communication with Ethereum blockchain
	backendURL := "https://ropsten.infura.io/v3/7c0fc2888a824c62a3651fd446c8f989"

	// customerPrivateKey private key used to control Customer's funds
	// Note: this is  written down for demonstration purpose only
	customerPrivateKey := "d4266a34672534b063d4d7911a7450d93f35075ef8583b3c96a408285fc21a38" // customerAccount public address: 0x8c77F5BA864718f098F83114efEC7180649afB85

	// merchantPrivateKey private key used to control Customer's funds
	// Note: this is  written down for demonstration purpose only
	merchantPrivateKey := "ad6c05a5d77f993cf6c23eab52c6b9db7894dabf5647df63155d3e66280c2dc3" // merchantAccount public address: 0xdF8c3E2c8506F67705acB0a4dCa28Cf44934B511

	processorContractAddress := common.HexToAddress(paymentProcessorContractAddress)

	// Initialize Customer account
	customerAccount, err := crypto.HexToECDSA(customerPrivateKey)
	if err != nil {
		log.Printf("error: %v", err.Error())
	}
	log.Printf("info: customer address - %v", crypto.PubkeyToAddress(customerAccount.PublicKey).String())
	customerAddress := crypto.PubkeyToAddress(customerAccount.PublicKey)

	// Initialize Merchant account
	merchantAccount, err := crypto.HexToECDSA(merchantPrivateKey)
	if err != nil {
		log.Printf("error: %v", err.Error())
	}
	log.Printf("info: merchant address - %v", crypto.PubkeyToAddress(merchantAccount.PublicKey).String())
	merchantAddress := crypto.PubkeyToAddress(merchantAccount.PublicKey)

	log.SetOutput(os.Stderr)
	ctx := utils.CreateCtrlCContext()

	client, err := ethclient.Dial(backendURL)
	if err != nil {
		log.Printf("error: %v", err.Error())
		return
	}

	e := eth.New(client, log.Printf)

	customerBalance, _ := e.Backend.BalanceAt(ctx, customerAddress, nil)
	if customerBalance.Cmp(math.Exp(big.NewInt(5), big.NewInt(15))) < 0 {
		log.Printf("error: Example Customer's balance insufficient to show the example")
		return
	}

	merchantBalance, _ := e.Backend.BalanceAt(ctx, merchantAddress, nil)
	if merchantBalance.Cmp(math.Exp(big.NewInt(10), big.NewInt(15))) < 0 {
		log.Printf("error: Example Merchant's balance insufficient to show the example")
		return
	}

	log.Print("info: starting example flow: customer initializes the order")

	customerSession = e.NewSession(customerAccount)
	p = processor.NewProcessor(customerSession, processorContractAddress)

	// Customer initializes an Order this way showing his intent to purchase Goods or Service

	// Order metadata
	orderID = big.NewInt(time.Now().Unix())
	price = math.Exp(big.NewInt(1*10), big.NewInt(15)) // price = 0.001 ETH

	log.Print("info: Customer adds order")
	txHash, err = p.AddOrder(ctx, orderID, price, customerAddress, common.Address{}, big.NewInt(0))
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

	// Customer transfers the order information to Merchant
	// Merchant receives information from Customer that order is created

	// Merchant validates that Order is paid
	order, err := p.ContractHandler.Orders(nil, orderID)

	if order.State != processor.OrderStatePaid {
		log.Println("error: Customer didn't pay for the Order yet")
	} else {

		// Initialize a new session with a Merchant's private key
		merchantSession = e.NewSession(merchantAccount)

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

		txHash, err = wallet.WithdrawAllTo(ctx, merchantAddress, big.NewInt(0))
		if err != nil {
			log.Printf("error: %v", err)
		}

		receipt, err = e.WaitForTxReceipt(ctx, txHash)
		if err != nil {
			log.Printf("error: %v", err)
		}
	}
	log.Print("info: done!")
}
