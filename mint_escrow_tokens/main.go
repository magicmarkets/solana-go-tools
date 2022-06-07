package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	sgo "github.com/magicmarkets/solana-go-tools"
)

func main() {
	ctx := context.Background()

	var payerFile, mintFile, receiverAddress, clusterRPC, clusterWS string

	var amount int

	flag.StringVar(&payerFile, "payer", "", "payer private key from solana-keygen file")
	flag.StringVar(&mintFile, "mint", "", "mint key from solana-keygen file")
	flag.StringVar(&receiverAddress, "receiver", "", "address of the receiving account, if empty to payer")
	flag.StringVar(&clusterRPC, "url", rpc.LocalNet_RPC, "solana cluster rpc url")
	flag.StringVar(&clusterWS, "ws", rpc.LocalNet_WS, "solana cluster websocket url")
	flag.IntVar(&amount, "amount", 100, "mint amount")
	flag.Parse()

	payer, err := solana.PrivateKeyFromSolanaKeygenFile(payerFile)
	if err != nil {
		fmt.Println("solana.PrivateKeyFromSolanaKeygenFile failed:", err)
		os.Exit(1)
	}

	fmt.Println("payer pubkey:", payer.PublicKey())

	if receiverAddress == "" {
		receiverAddress = payer.PublicKey().String()
	}

	rpcURL, err := sgo.RPCFromMoniker(clusterRPC)
	if err == nil {
		clusterRPC = rpcURL
	}

	wsURL, err := sgo.WSFromMoniker(clusterWS)
	if err == nil {
		clusterWS = wsURL
	}

	rpcClient := rpc.New(clusterRPC)

	wsClient, err := ws.Connect(context.Background(), clusterWS)
	if err != nil {
		fmt.Println("ws.Connect failed:", err)
		os.Exit(1)
	}

	escrowMint, err := solana.PrivateKeyFromSolanaKeygenFile(mintFile)
	if err != nil {
		fmt.Println("solana.PrivateKeyFromSolanaKeygenFile failed:", err)
		os.Exit(1)
	}

	fmt.Println("mint pubkey: ", escrowMint.PublicKey())

	mint, err := sgo.GetMint(ctx, rpcClient, escrowMint.PublicKey())
	if err != nil {
		fmt.Println("sgo.GetMint failed:", err)
		os.Exit(1)
	}

	mintAmount := uint64(math.Pow10(int(mint.Decimals)) * float64(amount))

	receiver := solana.MustPublicKeyFromBase58(receiverAddress)

	ata, _, err := solana.FindAssociatedTokenAddress(receiver, escrowMint.PublicKey())
	if err != nil {
		fmt.Println("solana.FindAssociatedTokenAddress failed:", err)
		os.Exit(1)
	}

	// fmt.Printf("associated token account create: %s for %s\n", ata, receiver)

	_, _ = sgo.CreateAssociatedTokenAccount(ctx, rpcClient, wsClient, receiver, escrowMint.PublicKey(), payer)

	sig, err := sgo.MintTo(ctx, rpcClient, wsClient, escrowMint.PublicKey(), ata, mintAmount, payer, payer)
	if err != nil {
		fmt.Println("sgo.MintTo failed:", err)
		os.Exit(1)
	}

	fmt.Printf("minted %v escrow tokens\n", amount)
	fmt.Println(sig)
}
