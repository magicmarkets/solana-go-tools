package sgo

import (
	metadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	"github.com/gagliardetto/solana-go"
)

/*
Example:
{
  "name": "USD Coin",
  "symbol": "USDC",
  "description": "Fully reserved fiat-backed stablecoin created by Circle.",
  "image": "https://www.circle.com/hs-fs/hubfs/sundaes/USDC.png?width=540&height=540&name=USDC.png"
}
*/

type TokenMetadata struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func CreateMetadataAccountV2(name, symbol string, mint, authority solana.PublicKey) (solana.Instruction, error) {
	// Metadata key (pda of ['metadata', program id, mint id])
	pda, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("metadata"),
			solana.TokenMetadataProgramID.Bytes(),
			mint.Bytes(),
		}, solana.TokenMetadataProgramID,
	)
	if err != nil {
		return nil, err
	}

	builder := metadata.NewCreateMetadataAccountV2InstructionBuilder().
		SetArgs(metadata.CreateMetadataAccountArgsV2{
			Data: metadata.DataV2{
				Name:   name,
				Symbol: symbol,
			},
			IsMutable: true,
		}).
		SetMetadataKeyPDAAccount(pda).
		SetMintOfTokenAccount(mint).
		SetMintAuthorityAccount(authority).
		SetPayerAccount(authority).
		SetUpdateAuthorityInfoAccount(authority).
		SetSystemAccount(solana.SystemProgramID).
		SetRentAccount(solana.SysVarRentPubkey)

	inst, err := builder.ValidateAndBuild()
	if err != nil {
		return nil, err
	}

	return inst, nil
}
