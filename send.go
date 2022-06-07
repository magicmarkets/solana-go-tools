package sgo

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

func SendTx(ctx context.Context, rpcClient *rpc.Client, wsClient *ws.Client, instructions []solana.Instruction, signers []solana.PrivateKey, payer solana.PrivateKey, synchronous bool) (*solana.Signature, error) {
	recent, err := rpcClient.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return nil, fmt.Errorf("sendTx: %w", err)
	}

	tx, err := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(payer.PublicKey()),
	)
	if err != nil {
		return nil, fmt.Errorf("sendTx: %w", err)
	}

	signers = append(signers, payer)

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			for _, signer := range signers {
				if signer.PublicKey().Equals(key) {
					return &signer
				}
			}

			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("sendTx: %w", err)
	}

	var sig solana.Signature
	if synchronous {
		sig, err = confirm.SendAndConfirmTransaction(
			ctx,
			rpcClient,
			wsClient,
			tx,
		)
	} else {
		sig, err = rpcClient.SendTransactionWithOpts(
			ctx,
			tx,
			false,
			rpc.CommitmentFinalized,
		)
	}

	if err != nil {
		return &sig, fmt.Errorf("sendTx: %w", err)
	}

	return &sig, nil
}
