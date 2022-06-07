package sgo

import (
	"errors"

	"github.com/gagliardetto/solana-go/rpc"
)

var ErrUnknownMoniker = errors.New("unknown moniker")

func RPCFromMoniker(moniker string) (string, error) {
	switch moniker {
	case "localnet", "localhost", "l":
		return rpc.LocalNet_RPC, nil
	case "testnet", "t":
		return rpc.TestNet_RPC, nil
	case "devnet", "d":
		return rpc.DevNet_RPC, nil
	case "mainnet", "m":
		return rpc.MainNetBeta_RPC, nil
	}

	return "", ErrUnknownMoniker
}

func WSFromMoniker(moniker string) (string, error) {
	switch moniker {
	case "localnet", "localhost", "l":
		return rpc.LocalNet_WS, nil
	case "testnet", "t":
		return rpc.TestNet_WS, nil
	case "devnet", "d":
		return rpc.DevNet_WS, nil
	case "mainnet", "m":
		return rpc.MainNetBeta_WS, nil
	}

	return "", ErrUnknownMoniker
}
