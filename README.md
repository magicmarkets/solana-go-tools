[![Go Report Card](https://goreportcard.com/badge/github.com/magicmarkets/solana-go-tools)](https://goreportcard.com/report/github.com/magicmarkets/solana-go-tools)
[![GoDoc](https://godoc.org/github.com/magicmarkets/solana-go-tools?status.svg)](https://godoc.org/github.com/magicmarkets/solana-go-tools)

# Solana Go Tools

Solana Go cli tools we used in early development while building
[magic.markets](https://magic.markets) to create and manage a stable coin mint
substitute on local and devnet. Such mint is used by our escrow program and is
the medium of exchange for market participants. This heavily depends on the
fantastic Go package
[gagliardetto/solana-go](https://github.com/gagliardetto/solana-go).

Custom rpc and websocket URLs are supported as well as monikers (and by their
first letter), mainnet-beta, testnet, devnet, localhost.

Create a new token keypair using the [solana toolchain](https://github.com/solana-labs/solana/releases):

    solana-keygen new --outfile mytoken.json

## Create and Mint Escrow Tokens

Create the mint using the keypair on devnet:

    go run create_escrow_mint/main.go -decimals 6 -payer ~/.config/solana/id.json -mint mytoken.json -url devnet

Mint 100k tokens to acccount `<PUBKEY>` on devnet:

    go run mint_escrow_tokens/main.go -amount 100000 -mint mytoken.json -payer ~/.config/solana/id.json -receiver <PUBKEY> -url d

## On-Chain Escrow Token Faucet

Transfer the mint authority to an on-chain faucet on devnet:

    go run init_escrow_faucet/main.go -mint mytoken.json -authority ~/.config/solana/id.json -url d

Run an airdrop against the faucet. In order to avoid spam the faucet gives out
tokens for devent SOL at a rate of 100x by default. Here we swap 1 SOL for 100
escrow tokens:

    go run faucet_airdrop/main.go -amount 1 -payer ~/.config/solana/id.json -mint <MINT PUBKEY> -url d

If you hold the faucet sweep authority, you can regularly claim the devnet SOL
deposited for tokens. On devnet for example:

    go run sweep_faucet/main.go -amount 1.2 -authority ~/.config/solana/id.json -mint <MINT PUBKEY> -url d
