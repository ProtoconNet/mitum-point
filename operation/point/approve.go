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
	ApproveFactHint = hint.MustNewHint("mitum-point-approve-operation-fact-v0.0.1")
	ApproveHint     = hint.MustNewHint("mitum-point-approve-operation-v0.0.1")
)

type ApproveFact struct {
	PointFact
	approved base.Address
	amount   common.Big
}

func NewApproveFact(
	token []byte,
	sender, contract base.Address,
	currency currencytypes.CurrencyID,
	approved base.Address,
	amount common.Big,
) ApproveFact {
	fact := ApproveFact{
		PointFact: NewPointFact(
			base.NewBaseFact(ApproveFactHint, token), sender, contract, currency,
		),
		approved: approved,
		amount:   amount,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact ApproveFact) IsValid(b []byte) error {
	if err := fact.PointFact.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if err := fact.approved.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if fact.sender.Equal(fact.approved) {
		return common.ErrFactInvalid.Wrap(
			common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with approved", fact.sender)))
	}

	if fact.contract.Equal(fact.approved) {
		return common.ErrFactInvalid.Wrap(
			common.ErrSelfTarget.Wrap(
				errors.Errorf("approve %v is is same with contract account", fact.approved)))
	}

	if !fact.amount.OverNil() {
		return common.ErrFactInvalid.Wrap(
			common.ErrValOOR.Wrap(
				errors.Errorf("approve amount must be greater than or equal to zero, got %v", fact.amount)))
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact ApproveFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact ApproveFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.PointFact.Bytes(),
		fact.approved.Bytes(),
		fact.amount.Bytes(),
	)
}

func (fact ApproveFact) Approved() base.Address {
	return fact.approved
}

func (fact ApproveFact) Amount() common.Big {
	return fact.amount
}

func (fact ApproveFact) Addresses() ([]base.Address, error) {
	var as []base.Address

	as = append(as, fact.PointFact.Sender())
	as = append(as, fact.PointFact.Contract())
	as = append(as, fact.approved)

	return as, nil
}

type Approve struct {
	common.BaseOperation
}

func NewApprove(fact ApproveFact) Approve {
	return Approve{BaseOperation: common.NewBaseOperation(ApproveHint, fact)}
}
