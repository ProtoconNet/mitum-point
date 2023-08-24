package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	MintFactHint = hint.MustNewHint("mitum-token-mint-operation-fact-v0.0.1")
	MintHint     = hint.MustNewHint("mitum-token-mint-operation-v0.0.1")
)

type MintFact struct {
	TokenFact
	receiver base.Address
	amount   common.Big
}

func NewMintFact(
	token []byte,
	sender, contract base.Address,
	tokenID, currency currencytypes.CurrencyID,
	receiver base.Address,
	amount common.Big,
) MintFact {
	fact := MintFact{
		TokenFact: NewTokenFact(
			base.NewBaseFact(MintFactHint, token), sender, contract, tokenID, currency,
		),
		receiver: receiver,
		amount:   amount,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact MintFact) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(fact))

	if err := fact.TokenFact.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if err := fact.receiver.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if fact.sender.Equal(fact.receiver) {
		return e.Wrap(errors.Errorf("sender is same with receiver, %s", fact.receiver))
	}

	if !fact.amount.OverZero() {
		return e.Wrap(errors.Errorf("zero amount"))
	}

	return nil
}

func (fact MintFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.TokenFact.Bytes(),
		fact.receiver.Bytes(),
		fact.amount.Bytes(),
	)
}

func (fact MintFact) Receiver() base.Address {
	return fact.receiver
}

func (fact MintFact) Amount() common.Big {
	return fact.amount
}

type Mint struct {
	TokenOperation
}

func NewMint(fact MintFact) Mint {
	return Mint{TokenOperation: NewTokenOperation(MintHint, fact)}
}
