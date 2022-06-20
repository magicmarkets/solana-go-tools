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
	"github.com/magicmarkets/solana-go-tools/generated/escrow_token_mint"
)

func main() {
	ctx := context.Background()

	var authority, mintFile, receiverAddress, clusterRPC, clusterWS string

	flag.StringVar(&authority, "authority", "", "payer private key from solana-keygen file that becomes the sweep authority")
	flag.StringVar(&mintFile, "mint", "", "mint key from solana-keygen file")
	flag.StringVar(&clusterRPC, "url", rpc.LocalNet_RPC, "solana cluster rpc url")
	flag.StringVar(&clusterWS, "ws", rpc.LocalNet_WS, "solana cluster websocket url")
	flag.Parse()

	payer, err := solana.PrivateKeyFromSolanaKeygenFile(authority)
	if err != nil {
		fmt.Println("solana.PrivateKeyFromSolanaKeygenFile failed:", err)
		os.Exit(1)
	}

	fmt.Println("payer pubkey:", payer.PublicKey())

	if receiverAddress == "" {
		receiverAddress = payer.PublicKey().String()
	}

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

	mint, err := solana.PrivateKeyFromSolanaKeygenFile(mintFile)
	if err != nil {
		fmt.Println("solana.PrivateKeyFromSolanaKeygenFile failed:", err)
		os.Exit(1)
	}

	fmt.Println("mint: ", mint.PublicKey())

	faucet, _, err := solana.FindProgramAddress([][]byte{mint.PublicKey().Bytes(), []byte("faucet_vault")}, escrow_token_mint.ProgramID)
	if err != nil {
		fmt.Println("solana.FindProgramAddress failed:", err)
		os.Exit(1)
	}

	builder := escrow_token_mint.NewInitializeInstruction(
		payer.PublicKey(),
		mint.PublicKey(),
		faucet,
		solana.SystemProgramID,
		solana.SysVarRentPubkey,
		solana.TokenProgramID,
	)

	inst, err := builder.ValidateAndBuild()
	if err != nil {
		fmt.Println("escrow_token_mint.NewInitializeInstruction failed:", err)
		os.Exit(1)
	}

	sig, err := sgo.SendTx(ctx, rpcClient, wsClient, []solana.Instruction{inst}, []solana.PrivateKey{}, payer, false)
	if err != nil {
		fmt.Println("sgo.SendTx failed:", err)
		os.Exit(1)
	}

	fmt.Println("faucet created:", faucet.String())
	fmt.Println(sig)
}
