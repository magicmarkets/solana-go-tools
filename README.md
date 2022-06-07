# Solana Go Tools

Tools we used in early development while building
[magic.markets](https://magic.markets) to create and manage a stable coin mint
substitute on localnet and devnet. Such mint is used by our escrow program and
market participants. This heavily depends on the fantastic Go package
[gagliardetto/solana-go](https://github.com/gagliardetto/solana-go).

Custom rpc and websocket URLs are supported as well as monikers (and by their
first letter), mainnet-beta, testnet, devnet, localhost.

Create a new token keypair using the [solana toolchain](https://github.com/solana-labs/solana/releases):

    solana-keygen new --outfile mytoken.json

Create the mint using the keypair:

    go run create_escrow_mint/main.go -decimals 6 -payer ~/.config/solana/id.json -mint mytoken.json -url devnet

Mint 100k tokens to an acccount:

    go run mint_escrow_tokens/main.go -amount 100000 -mint mytoken.json -payer ~/.config/solana/id.json -receiver <PUBKEY> -url d
