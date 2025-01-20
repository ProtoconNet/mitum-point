package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-point/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	RegisterModelFactHint = hint.MustNewHint("mitum-point-register-model-operation-fact-v0.0.1")
	RegisterModelHint     = hint.MustNewHint("mitum-point-register-model-operation-v0.0.1")
)

type RegisterModelFact struct {
	PointFact
	symbol        types.PointSymbol
	name          string
	decimal       common.Big
	initialSupply common.Big
}

func NewRegisterModelFact(
	token []byte,
	sender, contract base.Address,
	currency ctypes.CurrencyID,
	symbol types.PointSymbol,
	name string,
	decimal common.Big,
	initialSupply common.Big,
) RegisterModelFact {
	fact := RegisterModelFact{
		PointFact: NewPointFact(
			base.NewBaseFact(RegisterModelFactHint, token), sender, contract, currency,
		),
		symbol:        symbol,
		name:          name,
		decimal:       decimal,
		initialSupply: initialSupply,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact RegisterModelFact) IsValid(b []byte) error {
	if err := util.CheckIsValiders(nil, false, fact.PointFact, fact.symbol); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if fact.name == "" {
		return common.ErrFactInvalid.Wrap(common.ErrValueInvalid.Wrap(errors.Errorf("empty symbol")))
	}

	if !fact.decimal.OverNil() {
		return common.ErrFactInvalid.Wrap(
			common.ErrValOOR.Wrap(
				errors.Errorf("decimal must be bigger than or equal to zero, got %v", fact.decimal)))
	}

	if !fact.initialSupply.OverNil() {
		return common.ErrFactInvalid.Wrap(
			common.ErrValOOR.Wrap(
				errors.Errorf("initial supply must be bigger than or equal to zero, got %v", fact.initialSupply)))
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}
	return nil
}

func (fact RegisterModelFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact RegisterModelFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.PointFact.Bytes(),
		fact.symbol.Bytes(),
		[]byte(fact.name),
		fact.decimal.Bytes(),
		fact.initialSupply.Bytes(),
	)
}

func (fact RegisterModelFact) Name() string {
	return fact.name
}

func (fact RegisterModelFact) Symbol() types.PointSymbol {
	return fact.symbol
}

func (fact RegisterModelFact) Decimal() common.Big {
	return fact.decimal
}

func (fact RegisterModelFact) InitialSupply() common.Big {
	return fact.initialSupply
}

func (fact RegisterModelFact) InActiveContractOwnerHandlerOnly() [][2]base.Address {
	return [][2]base.Address{{fact.contract, fact.sender}}
}

type RegisterModel struct {
	extras.ExtendedOperation
}

func NewRegisterModel(fact RegisterModelFact) RegisterModel {
	return RegisterModel{
		ExtendedOperation: extras.NewExtendedOperation(RegisterModelHint, fact),
	}
}
