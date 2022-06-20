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
	var name, symbol string

	var decimals int

	flag.StringVar(&payerFile, "payer", "", "payer private key from solana-keygen file")
	flag.StringVar(&mintFile, "mint", "", "mint key from solana-keygen file")
	flag.StringVar(&clusterRPC, "url", rpc.LocalNet_RPC, "solana cluster rpc url")
	flag.StringVar(&clusterWS, "ws", "", "solana cluster websocket url")
	flag.IntVar(&decimals, "decimals", defaultDecimals, "mint decimals")
	flag.StringVar(&name, "name", "", "optional name for the metaplex token metadata")
	flag.StringVar(&symbol, "symbol", "", "optional symbol for the metaplex token metadata")
	flag.Parse()

	payer, err := solana.PrivateKeyFromSolanaKeygenFile(payerFile)
	if err != nil {
		fmt.Println("solana.PrivateKeyFromSolanaKeygenFile failed:", err)
		os.Exit(1)
	}

	fmt.Println("payer pubkey:", payer.PublicKey())

	wsURL, err := sgo.WSFromMoniker(clusterRPC)
	if err == nil {
		clusterWS = wsURL
	}

	rpcURL, err := sgo.RPCFromMoniker(clusterRPC)
	if err == nil {
		clusterRPC = rpcURL
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

	var instructions []solana.Instruction

	escrowMintInst, err := sgo.NewMintInstruction(ctx, rpcClient, uint8(decimals), escrowMint.PublicKey(), payer.PublicKey(), payer.PublicKey())
	if err != nil {
		fmt.Println("sgo.NewMintInstruction failed:", err)
		os.Exit(1)
	}

	instructions = append(instructions, escrowMintInst...)

	if name != "" && symbol != "" {
		fmt.Println("adding metaplex token metadata")

		metadataInst, err := sgo.CreateMetadataAccountV2(name, symbol, escrowMint.PublicKey(), payer.PublicKey())
		if err != nil {
			fmt.Println("sgo.CreateMetadataAccountV2 failed:", err)
			os.Exit(1)
		}

		instructions = append(instructions, metadataInst)
	}

	fmt.Println("creating mint...")

	sig, err := sgo.SendTx(ctx, rpcClient, wsClient, instructions, []solana.PrivateKey{escrowMint, payer}, payer, false)
	if err != nil {
		fmt.Println("sgo.SendTx failed:", err)
		os.Exit(1)
	}

	fmt.Println(sig)
}
