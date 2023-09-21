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
	RegisterTokenFactHint = hint.MustNewHint("mitum-token-register-token-operation-fact-v0.0.1")
	RegisterTokenHint     = hint.MustNewHint("mitum-token-register-token-operation-v0.0.1")
)

type RegisterTokenFact struct {
	TokenFact
	symbol      currencytypes.CurrencyID
	name        string
	totalSupply common.Big
}

func NewRegisterTokenFact(
	token []byte,
	sender, contract base.Address,
	currency currencytypes.CurrencyID,
	symbol currencytypes.CurrencyID,
	name string,
	totalSupply common.Big,
) RegisterTokenFact {
	fact := RegisterTokenFact{
		TokenFact: NewTokenFact(
			base.NewBaseFact(RegisterTokenFactHint, token), sender, contract, currency,
		),
		symbol:      symbol,
		name:        name,
		totalSupply: totalSupply,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact RegisterTokenFact) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(fact))

	if err := util.CheckIsValiders(nil, false, fact.TokenFact, fact.symbol); err != nil {
		return e.Wrap(err)
	}

	if fact.name == "" {
		return e.Wrap(errors.Errorf("empty symbol"))
	}

	if !fact.totalSupply.OverNil() {
		return e.Wrap(errors.Errorf("nil big"))
	}

	return nil
}

func (fact RegisterTokenFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.TokenFact.Bytes(),
		fact.symbol.Bytes(),
		[]byte(fact.name),
		fact.totalSupply.Bytes(),
	)
}

func (fact RegisterTokenFact) Name() string {
	return fact.name
}

func (fact RegisterTokenFact) Symbol() currencytypes.CurrencyID {
	return fact.symbol
}

func (fact RegisterTokenFact) TotalSupply() common.Big {
	return fact.totalSupply
}

type RegisterToken struct {
	common.BaseOperation
}

func NewRegisterToken(fact RegisterTokenFact) RegisterToken {
	return RegisterToken{BaseOperation: common.NewBaseOperation(RegisterTokenHint, fact)}
}
