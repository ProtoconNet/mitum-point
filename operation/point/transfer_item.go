package point

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var TransferItemHint = hint.MustNewHint("mitum-point-transfer-item-v0.0.1")

type TransferItem struct {
	hint.BaseHinter
	contract base.Address
	receiver base.Address
	amount   common.Big
	currency types.CurrencyID
}

func NewTransferItem(
	contract base.Address, receiver base.Address, amount common.Big, currency types.CurrencyID,
) TransferItem {
	return TransferItem{
		BaseHinter: hint.NewBaseHinter(TransferItemHint),
		contract:   contract,
		receiver:   receiver,
		amount:     amount,
		currency:   currency,
	}
}

func (it TransferItem) IsValid([]byte) error {
	if err := it.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false, it.contract, it.receiver); err != nil {
		return err
	}

	if it.receiver.Equal(it.contract) {
		return common.ErrSelfTarget.Wrap(errors.Errorf("receiver %v is same with contract account", it.receiver))
	}

	if !it.amount.OverZero() {
		return common.ErrFactInvalid.Wrap(
			common.ErrValOOR.Wrap(errors.Errorf("transfer amount must be over zero, got %v", it.amount)))
	}

	return util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.receiver,
		it.currency,
	)
}

func (it TransferItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.receiver.Bytes(),
		it.amount.Bytes(),
		it.currency.Bytes(),
	)
}

func (it TransferItem) Contract() base.Address {
	return it.contract
}

func (it TransferItem) Receiver() base.Address {
	return it.receiver
}

func (it TransferItem) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = it.receiver
	return as, nil
}

func (it TransferItem) Amount() common.Big {
	return it.amount
}

func (it TransferItem) Currency() types.CurrencyID {
	return it.currency
}
