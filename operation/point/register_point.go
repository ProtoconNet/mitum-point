package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-point/types"
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
	symbol        types.PointID
	name          string
	initialSupply common.Big
}

func NewRegisterPointFact(
	token []byte,
	sender, contract base.Address,
	currency currencytypes.CurrencyID,
	symbol types.PointID,
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
	if err := util.CheckIsValiders(nil, false, fact.PointFact, fact.symbol); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if fact.name == "" {
		return common.ErrFactInvalid.Wrap(common.ErrValueInvalid.Wrap(errors.Errorf("empty symbol")))
	}

	if !fact.initialSupply.OverNil() {
		return common.ErrFactInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("zero initial supply")))
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
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

func (fact RegisterPointFact) Symbol() types.PointID {
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
