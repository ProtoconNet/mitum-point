package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	"github.com/ProtoconNet/mitum-currency/v3/types"
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

var MaxTransferFromItems = 100

type TransferFromFact struct {
	base.BaseFact
	sender base.Address
	items  []TransferFromItem
}

func NewTransferFromFact(
	token []byte,
	sender base.Address,
	items []TransferFromItem,
) TransferFromFact {
	fact := TransferFromFact{
		BaseFact: base.NewBaseFact(TransferFromFactHint, token),
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact TransferFromFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if l := len(fact.items); l < 1 {
		return common.ErrFactInvalid.Wrap(common.ErrArrayLen.Wrap(errors.Errorf("empty items for TransferFromFact")))
	} else if l > int(MaxTransferFromItems) {
		return common.ErrFactInvalid.Wrap(
			common.ErrArrayLen.Wrap(errors.Errorf("items over allowed, %d > %d", l, MaxTransferFromItems)))
	}

	if err := util.CheckIsValiders(nil, false,
		fact.BaseFact,
		fact.sender,
	); err != nil {
		return err
	}

	founds := map[string]struct{}{}
	for _, item := range fact.items {
		if err := item.IsValid(nil); err != nil {
			return common.ErrFactInvalid.Wrap(err)
		}

		if fact.sender.Equal(item.contract) {
			return common.ErrFactInvalid.Wrap(
				common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
		}

		if _, found := founds[item.contract.String()+"-"+item.target.String()+"-"+item.receiver.String()]; found {
			return common.ErrFactInvalid.Wrap(
				common.ErrDupVal.Wrap(
					errors.Errorf(
						"target account %v and receiver account %v in contract account %v",
						item.target, item.receiver, item.contract)))
		}

		founds[item.contract.String()+"-"+item.target.String()+"-"+item.receiver.String()] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}
	return nil
}

func (fact TransferFromFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact TransferFromFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))
	for i := range fact.items {
		is[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact TransferFromFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact TransferFromFact) Sender() base.Address {
	return fact.sender
}

func (fact TransferFromFact) Items() []TransferFromItem {
	return fact.items
}

func (fact TransferFromFact) Addresses() ([]base.Address, error) {
	var as []base.Address

	for i := range fact.items {
		if ads, err := fact.items[i].Addresses(); err != nil {
			return nil, err
		} else {
			as = append(as, ads...)
		}
	}

	as = append(as, fact.Sender())

	return as, nil
}

func (fact TransferFromFact) FeeBase() map[types.CurrencyID][]common.Big {
	required := make(map[types.CurrencyID][]common.Big)

	for i := range fact.items {
		zeroBig := common.ZeroBig
		cid := fact.items[i].Currency()
		var amsTemp []common.Big
		if ams, found := required[cid]; found {
			ams = append(ams, zeroBig)
			required[cid] = ams
		} else {
			amsTemp = append(amsTemp, zeroBig)
			required[cid] = amsTemp
		}
	}

	return required
}

func (fact TransferFromFact) FeePayer() base.Address {
	return fact.sender
}

func (fact TransferFromFact) FeeItemCount() (uint, bool) {
	return uint(len(fact.items)), extras.HasItem
}

func (fact TransferFromFact) FactUser() base.Address {
	return fact.sender
}

func (fact TransferFromFact) Signer() base.Address {
	return fact.sender
}

func (fact TransferFromFact) ActiveContract() []base.Address {
	var arr []base.Address
	for i := range fact.items {
		arr = append(arr, fact.items[i].contract)
	}
	return arr
}

type TransferFrom struct {
	extras.ExtendedOperation
}

func NewTransferFrom(fact TransferFromFact) TransferFrom {
	return TransferFrom{
		ExtendedOperation: extras.NewExtendedOperation(TransferFromHint, fact),
	}
}
