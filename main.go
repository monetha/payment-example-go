package main

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/chequebook"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/monetha/payment-go-sdk/paymenthandler"
	"github.com/monetha/reputation-go-sdk/eth/backend"
)

func main() {
	log.SetOutput(os.Stderr)
	ctx := createCtrlCContext()

	var paymentBackend chequebook.Backend

	backendURL := "https://ropsten.infura.io/v3/9341cea07e634c21be9d5a5ccb892db5"
	paymentBackend, err := ethclient.Dial(backendURL)
	if err != nil {
		log.Printf("error: %v", err.Error())
		return
	}

	// merchantKey public address: 0xdF8c3E2c8506F67705acB0a4dCa28Cf44934B511
	merchantKey, err := crypto.HexToECDSA("ad6c05a5d77f993cf6c23eab52c6b9db7894dabf5647df63155d3e66280c2dc3")
	if err != nil {
		log.Printf("error: %v", err.Error())
	}
	log.Printf("info: merchant address - %v", crypto.PubkeyToAddress(merchantKey.PublicKey))

	backend.NewHandleNonceBackend(paymentBackend, []common.Address{crypto.PubkeyToAddress(merchantKey)})

	// Create an account
	key, err := crypto.GenerateKey()

	// Get the address
	address := crypto.PubkeyToAddress(key.PublicKey)
	// 0x8ee3333cDE801ceE9471ADf23370c48b011f82a6

	// Get the private key
	keyStr := hex.EncodeToString(key.D.Bytes())
	// 05b14254a1d0c77a49eae3bdf080f926a2df17d8e2ebdf7af941ea001481e57f

	balance, err := paymentBackend.BalanceAt(ctx, address, nil)

	zeroBalance := big.NewInt(0)

	if balance.CmpAbs(zeroBalance) > 0 {

	}

	// processorContractAddress address of the Merchant's smart contract
	processorContractAddress := common.HexToAddress("0x35D6708FD36DCb902adce5D9d6ABeB4838318554")

	operationsAuth := bind.NewKeyedTransactor(merchantKey)

	processor, err := paymenthandler.New("ropsten.infura.io/v3/9341cea07e634c21be9d5a5ccb892db5", keyStr)
	if err != nil {
		log.Printf("error: %v", err.Error())
	}

	// 0.2 ETH == 200000000000000000 == 2 * 10^17
	//processor.AddOrder(ctx, processorContractAddress, suggestedGasprice, orderRef, 200000000000000000, )


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
