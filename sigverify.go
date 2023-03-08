package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	ctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc"
)

func main() {
	// create the sign doc json obj. and marshal it
	Msg := createMsgSignData("juno1n9e6zfv956xn2m36q3qzjq2gdpa5zqnzxqlrtn", "dGVzdA==")
	enc, err := json.Marshal(Msg)
	if err != nil {
		panic(err)
	}

	// utf-8 encoding of string, sorted by key
	fmt.Println(string(enc))
	// register interfaces for unpacking protobuf data
	reg := ctypes.NewInterfaceRegistry()
	authtypes.RegisterInterfaces(reg)
	cryptocodec.RegisterInterfaces(reg)
	pCodec := codec.NewProtoCodec(reg)

	// create grpc conn to juno uni-5 node
	grpcConn, err := grpc.Dial(
		"194.61.28.217:9090",
		grpc.WithInsecure(),
	)
	if err != nil {
		panic(err)
	}
	defer grpcConn.Close()
	authClient := authtypes.NewQueryClient(grpcConn)
	// retrieve pubKey for account
	authResp, err := authClient.Account(context.Background(), &authtypes.QueryAccountRequest{
		Address: "juno1n9e6zfv956xn2m36q3qzjq2gdpa5zqnzxqlrtn",
	})
	var valOperAcct authtypes.AccountI
	// cast account response to pubKey
	if err = pCodec.UnpackAny(authResp.Account, &valOperAcct); err != nil {
		panic(err)
	}
	// retrieve signature bytes from base64 encoded signature
	sig, err := base64.StdEncoding.DecodeString("fTfH8K1B1D/GyOCT7AjjSx35qsgDdonSl5yJnFbtT6YDfuwzB50PRIfsxCz8oHVqMpsyq7iLT/acFaeVMzw2nw==")
	if err != nil {
		panic(err)
	}

	if !valOperAcct.GetPubKey().VerifySignature(enc, sig) {
		panic(fmt.Errorf("WRONG SIGNATURE"))
	}
	fmt.Println("SIGNATURE IS CORRECT!!!!")
}

type MsgSignData struct {
	Account_number string `json:"account_number"`
	Chain_id       string `json:"chain_id"`
	Fee            struct {
		Amount []string `json:"amount"`
		Gas    string   `json:"gas"`
	} `json:"fee"`
	Memo string `json:"memo"`
	Msgs []struct {
		Typ   string `json:"type"`
		Value struct {
			Data   string `json:"data"`
			Signer string `json:"signer"`
		} `json:"value"`
	} `json:"msgs"`
	Sequence string `json:"sequence"`
}

// MsgSignData signDoc where signer is bech32 address of signer
// data is the base64 encoded payload to sign
func createMsgSignData(signer, data string) *MsgSignData {
	return &MsgSignData{
		Chain_id:       "",
		Account_number: "0",
		Sequence:       "0",
		Fee: struct {
			Amount []string "json:\"amount\""
			Gas    string   "json:\"gas\""
		}{
			Gas:    "0",
			Amount: []string{},
		},
		Msgs: []struct {
			Typ   string "json:\"type\""
			Value struct {
				Data   string "json:\"data\""
				Signer string "json:\"signer\""
			} "json:\"value\""
		}{{
			Typ: "sign/MsgSignData",
			Value: struct {
				Data   string "json:\"data\""
				Signer string "json:\"signer\""
			}{
				Signer: signer,
				Data:   data,
			},
		},
		},
		Memo: "",
	}
}
