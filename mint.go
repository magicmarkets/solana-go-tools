package sgo

import (
	"context"
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

const MintAccountSize = 82

func NewMintInstruction(ctx context.Context, client *rpc.Client, decimals uint8, mint, owner, payer solana.PublicKey) ([]solana.Instruction, error) {
	createAccountInst, err := NewAccountInstruction(ctx, client, mint, solana.TokenProgramID, payer, MintAccountSize)
	if err != nil {
		return nil, fmt.Errorf("create token account: %w", err)
	}

	initMintInst, err := InitMintInstruction(ctx, decimals, mint, owner)
	if err != nil {
		return nil, fmt.Errorf("create token account: %w", err)
	}

	return []solana.Instruction{createAccountInst, initMintInst}, nil
}

func InitMintInstruction(ctx context.Context, decimals uint8, mint, owner solana.PublicKey) (solana.Instruction, error) {
	instBuilder := token.NewInitializeMintInstructionBuilder().
		SetDecimals(decimals).
		SetMintAuthority(owner).
		SetMintAccount(mint).
		SetSysVarRentPubkeyAccount(solana.SysVarRentPubkey)

	inst, err := instBuilder.ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("initialize account: %w", err)
	}

	return inst, nil
}

func MintToInstruction(ctx context.Context, mintAccount, mintAuthority, tokenAccount solana.PublicKey, amount uint64) (solana.Instruction, error) {
	instBuilder := token.NewMintToInstructionBuilder().
		SetMintAccount(mintAccount).
		SetDestinationAccount(tokenAccount).
		SetAuthorityAccount(mintAuthority).
		SetAmount(amount)

	inst, err := instBuilder.ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("mint to: %w", err)
	}

	return inst, nil
}

func GetMint(ctx context.Context, client *rpc.Client, address solana.PublicKey) (*token.Mint, error) {
	acct, err := client.GetAccountInfo(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("getMint: %w", err)
	}

	var mint token.Mint

	if err = bin.NewBinDecoder(acct.Value.Data.GetBinary()).Decode(&mint); err != nil {
		return nil, fmt.Errorf("getMint: %w", err)
	}

	return &mint, nil
}

func MintTo(ctx context.Context, client *rpc.Client, ws *ws.Client, mintAccount, tokenAccount solana.PublicKey, amount uint64, mintAuthority, payer solana.PrivateKey) (*solana.Signature, error) {
	mintToInst, err := MintToInstruction(ctx, mintAccount, mintAuthority.PublicKey(), tokenAccount, amount)
	if err != nil {
		return nil, fmt.Errorf("create token account: %w", err)
	}

	return SendTx(ctx, client, ws, []solana.Instruction{mintToInst}, []solana.PrivateKey{mintAuthority, payer}, payer, false)
}
