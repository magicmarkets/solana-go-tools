package sgo

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

func NewAccountInstruction(ctx context.Context, client *rpc.Client, account, owner, payer solana.PublicKey, size uint64) (solana.Instruction, error) {
	lamports, err := client.GetMinimumBalanceForRentExemption(
		ctx,
		size,
		rpc.CommitmentFinalized,
	)

	if err != nil {
		return nil, fmt.Errorf("new account: %w", err)
	}

	instBuilder := system.NewCreateAccountInstruction(lamports, size, owner, payer, account)

	inst, err := instBuilder.ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("new account: %w", err)
	}

	return inst, nil
}

func CreateAssociatedTokenAccount(ctx context.Context, client *rpc.Client, ws *ws.Client, wallet, mint solana.PublicKey, payer solana.PrivateKey) (*solana.PublicKey, error) {
	createInst := associatedtokenaccount.NewCreateInstruction(payer.PublicKey(), wallet, mint)

	inst, err := createInst.ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("create associated token account: %w", err)
	}

	_, err = SendTx(ctx, client, ws, []solana.Instruction{inst}, []solana.PrivateKey{}, payer, true)
	if err != nil {
		return nil, fmt.Errorf("send tx: %w", err)
	}

	ata, _, err := solana.FindAssociatedTokenAddress(wallet, mint)
	if err != nil {
		return nil, fmt.Errorf("find associated token account: %w", err)
	}

	return ata.ToPointer(), nil
}
