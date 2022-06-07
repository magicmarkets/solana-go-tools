package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	sgo "github.com/magicmarkets/solana-go-tools"
)

const defaultDecimals = 6

func main() {
	ctx := context.Background()

	var payerFile, mintFile, clusterRPC, clusterWS string

	var decimals int

	flag.StringVar(&payerFile, "payer", "", "payer private key from solana-keygen file")
	flag.StringVar(&mintFile, "mint", "", "mint key from solana-keygen file")
	flag.StringVar(&clusterRPC, "url", rpc.LocalNet_RPC, "solana cluster rpc url")
	flag.StringVar(&clusterWS, "ws", rpc.LocalNet_WS, "solana cluster websocket url")
	flag.IntVar(&decimals, "decimals", defaultDecimals, "mint decimals")
	flag.Parse()

	payer, err := solana.PrivateKeyFromSolanaKeygenFile(payerFile)
	if err != nil {
		fmt.Println("solana.PrivateKeyFromSolanaKeygenFile failed:", err)
		os.Exit(1)
	}

	fmt.Println("payer pubkey:", payer.PublicKey())

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

	escrowMintInst, err := sgo.NewMintInstruction(ctx, rpcClient, uint8(decimals), escrowMint.PublicKey(), payer.PublicKey(), payer.PublicKey())
	if err != nil {
		fmt.Println("sgo.NewMintInstruction failed:", err)
		os.Exit(1)
	}

	fmt.Println("creating mint...")

	sig, err := sgo.SendTx(ctx, rpcClient, wsClient, escrowMintInst, []solana.PrivateKey{escrowMint, payer}, payer, true)
	if err != nil {
		fmt.Println("sgo.SendTx failed:", err)
		os.Exit(1)
	}

	fmt.Println(sig)
}
