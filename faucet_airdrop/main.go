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
	"github.com/magicmarkets/solana-go-tools/generated/escrow_token_mint"
)

func main() {
	ctx := context.Background()

	var payerFile, mintPubkey, clusterRPC, clusterWS string
	var amount float64

	flag.StringVar(&payerFile, "payer", "", "payer private key from solana-keygen file")
	flag.StringVar(&mintPubkey, "mint", "", "mint address")
	flag.Float64Var(&amount, "amount", 1, "amount of SOL to exchange for tokens")
	flag.StringVar(&clusterRPC, "url", rpc.LocalNet_RPC, "solana cluster rpc url")
	flag.StringVar(&clusterWS, "ws", rpc.LocalNet_WS, "solana cluster websocket url")
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

	mint, err := solana.PublicKeyFromBase58(mintPubkey)
	if err != nil {
		fmt.Println("solana.PrivateKeyFromSolanaKeygenFile failed:", err)
		os.Exit(1)
	}

	fmt.Println("mint: ", mint)

	faucet, _, err := solana.FindProgramAddress([][]byte{mint.Bytes(), []byte("faucet_vault")}, escrow_token_mint.ProgramID)
	if err != nil {
		fmt.Println("solana.FindProgramAddress failed:", err)
		os.Exit(1)
	}

	authority, _, err := solana.FindProgramAddress([][]byte{[]byte("faucet_authority")}, escrow_token_mint.ProgramID)
	if err != nil {
		fmt.Println("solana.FindProgramAddress failed:", err)
		os.Exit(1)
	}

	ata, _, err := solana.FindAssociatedTokenAddress(payer.PublicKey(), mint)
	if err != nil {
		fmt.Println("solana.FindAssociatedTokenAddress failed:", err)
		os.Exit(1)
	}

	_, _ = sgo.CreateAssociatedTokenAccount(ctx, rpcClient, wsClient, payer.PublicKey(), mint, payer)

	lamports := uint64(amount * math.Pow10(9))

	builder := escrow_token_mint.NewSwapInstruction(
		lamports,
		payer.PublicKey(),
		ata,
		mint,
		faucet,
		authority,
		solana.SystemProgramID,
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

	fmt.Println("airdropped using lamports", lamports)
	fmt.Println(sig)
}
