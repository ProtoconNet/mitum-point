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
	TransferFromFactHint = hint.MustNewHint("mitum-point-transfer-from-operation-fact-v0.0.1")
	TransferFromHint     = hint.MustNewHint("mitum-point-transfer-from-operation-v0.0.1")
)

type TransferFromFact struct {
	PointFact
	receiver base.Address
	target   base.Address
	amount   common.Big
}

func NewTransferFromFact(
	token []byte,
	sender, contract base.Address,
	currency currencytypes.CurrencyID,
	receiver, target base.Address,
	amount common.Big,
) TransferFromFact {
	fact := TransferFromFact{
		PointFact: NewPointFact(
			base.NewBaseFact(TransferFromFactHint, token), sender, contract, currency,
		),
		receiver: receiver,
		target:   target,
		amount:   amount,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact TransferFromFact) IsValid(b []byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(fact))

	if err := fact.PointFact.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if err := fact.receiver.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if err := fact.target.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if fact.contract.Equal(fact.receiver) {
		return e.Wrap(errors.Errorf("contract address is same with receiver, %s", fact.receiver))
	}

	if fact.contract.Equal(fact.target) {
		return e.Wrap(errors.Errorf("contract address is same with target, %s", fact.target))
	}

	if fact.receiver.Equal(fact.target) {
		return e.Wrap(errors.Errorf("target is same with receiver, %s", fact.receiver))
	}

	if fact.sender.Equal(fact.target) {
		return e.Wrap(errors.Errorf("sender is same with target, %s", fact.target))
	}

	if !fact.amount.OverZero() {
		return e.Wrap(errors.Errorf("zero amount"))
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}
	return nil
}

func (fact TransferFromFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact TransferFromFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.PointFact.Bytes(),
		fact.receiver.Bytes(),
		fact.target.Bytes(),
		fact.amount.Bytes(),
	)
}

func (fact TransferFromFact) Receiver() base.Address {
	return fact.receiver
}

func (fact TransferFromFact) Target() base.Address {
	return fact.target
}

func (fact TransferFromFact) Amount() common.Big {
	return fact.amount
}

func (fact TransferFromFact) Addresses() ([]base.Address, error) {
	var as []base.Address

	as = append(as, fact.PointFact.Sender())
	as = append(as, fact.PointFact.Contract())
	as = append(as, fact.receiver)
	as = append(as, fact.target)

	return as, nil
}

type TransferFrom struct {
	common.BaseOperation
}

func NewTransferFrom(fact TransferFromFact) TransferFrom {
	return TransferFrom{BaseOperation: common.NewBaseOperation(TransferFromHint, fact)}
}
