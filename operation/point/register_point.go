package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-point/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	RegisterPointFactHint = hint.MustNewHint("mitum-point-register-point-operation-fact-v0.0.1")
	RegisterPointHint     = hint.MustNewHint("mitum-point-register-point-operation-v0.0.1")
)

type RegisterPointFact struct {
	PointFact
	symbol        currencytypes.CurrencyID
	name          string
	initialSupply common.Big
}

func NewRegisterPointFact(
	token []byte,
	sender, contract base.Address,
	currency currencytypes.CurrencyID,
	symbol currencytypes.CurrencyID,
	name string,
	initialSupply common.Big,
) RegisterPointFact {
	fact := RegisterPointFact{
		PointFact: NewPointFact(
			base.NewBaseFact(RegisterPointFactHint, token), sender, contract, currency,
		),
		symbol:        symbol,
		name:          name,
		initialSupply: initialSupply,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact RegisterPointFact) IsValid(b []byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(fact))

	if err := util.CheckIsValiders(nil, false, fact.PointFact, fact.symbol); err != nil {
		return e.Wrap(err)
	}

	if fact.name == "" {
		return e.Wrap(errors.Errorf("empty symbol"))
	}

	if !fact.initialSupply.OverNil() {
		return e.Wrap(errors.Errorf("nil big"))
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}
	return nil
}

func (fact RegisterPointFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact RegisterPointFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.PointFact.Bytes(),
		fact.symbol.Bytes(),
		[]byte(fact.name),
		fact.initialSupply.Bytes(),
	)
}

func (fact RegisterPointFact) Name() string {
	return fact.name
}

func (fact RegisterPointFact) Symbol() currencytypes.CurrencyID {
	return fact.symbol
}

func (fact RegisterPointFact) InitialSupply() common.Big {
	return fact.initialSupply
}

type RegisterPoint struct {
	common.BaseOperation
}

func NewRegisterPoint(fact RegisterPointFact) RegisterPoint {
	return RegisterPoint{BaseOperation: common.NewBaseOperation(RegisterPointHint, fact)}
}
