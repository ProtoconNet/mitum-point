package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	BurnFactHint = hint.MustNewHint("mitum-point-burn-operation-fact-v0.0.1")
	BurnHint     = hint.MustNewHint("mitum-point-burn-operation-v0.0.1")
)

type BurnFact struct {
	PointFact
	target base.Address
	amount common.Big
}

func NewBurnFact(
	token []byte,
	sender, contract base.Address,
	currency currencytypes.CurrencyID,
	target base.Address,
	amount common.Big,
) BurnFact {
	fact := BurnFact{
		PointFact: NewPointFact(
			base.NewBaseFact(BurnFactHint, token), sender, contract, currency,
		),
		target: target,
		amount: amount,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact BurnFact) IsValid(b []byte) error {
	if err := fact.PointFact.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if err := fact.target.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if fact.contract.Equal(fact.target) {
		return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("target %v is same with contract account", fact.target)))
	}

	if !fact.amount.OverZero() {
		return common.ErrFactInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("burn amount must be over zero, got %v", fact.amount)))
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}
	return nil
}

func (fact BurnFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact BurnFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.PointFact.Bytes(),
		fact.target.Bytes(),
		fact.amount.Bytes(),
	)
}

func (fact BurnFact) Target() base.Address {
	return fact.target
}

func (fact BurnFact) Amount() common.Big {
	return fact.amount
}

func (fact BurnFact) Addresses() ([]base.Address, error) {
	var as []base.Address

	as = append(as, fact.PointFact.Sender())
	as = append(as, fact.PointFact.Contract())
	as = append(as, fact.target)

	return as, nil
}

type Burn struct {
	common.BaseOperation
}

func NewBurn(fact BurnFact) Burn {
	return Burn{BaseOperation: common.NewBaseOperation(BurnHint, fact)}
}
