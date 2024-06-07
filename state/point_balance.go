package state

import (
	"fmt"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-point/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	PointBalanceStateValueHint = hint.MustNewHint("mitum-point-balance-state-value-v0.0.1")
	PointBalanceSuffix         = "pointbalance"
)

type PointBalanceStateValue struct {
	hint.BaseHinter
	Amount common.Big
}

func NewPointBalanceStateValue(amount common.Big) PointBalanceStateValue {
	return PointBalanceStateValue{
		BaseHinter: hint.NewBaseHinter(PointBalanceStateValueHint),
		Amount:     amount,
	}
}

func (s PointBalanceStateValue) Hint() hint.Hint {
	return s.BaseHinter.Hint()
}

func (s PointBalanceStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(s))

	if err := s.BaseHinter.IsValid(PointBalanceStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if !s.Amount.OverNil() {
		return e.Wrap(errors.Errorf("nil big"))
	}

	return nil
}

func (s PointBalanceStateValue) HashBytes() []byte {
	return s.Amount.Bytes()
}

func StatePointBalanceValue(st base.State) (common.Big, error) {
	e := util.ErrNotFound.Errorf(ErrStringStateNotFound(st.Key()))

	v := st.Value()
	if v == nil {
		return common.NilBig, e.Wrap(errors.Errorf("nil value"))
	}

	s, ok := v.(PointBalanceStateValue)
	if !ok {
		return common.NilBig, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(PointBalanceStateValue{}, v)))
	}

	return s.Amount, nil
}

type AddPointBalanceStateValue struct {
	Amount common.Big
}

func NewAddPointBalanceStateValue(amount common.Big) AddPointBalanceStateValue {
	return AddPointBalanceStateValue{
		Amount: amount,
	}
}

func (b AddPointBalanceStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid AddPointBalanceStateValue")

	if err := util.CheckIsValiders(nil, false, b.Amount); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (b AddPointBalanceStateValue) HashBytes() []byte {
	return b.Amount.Bytes()
}

type DeductPointBalanceStateValue struct {
	Amount common.Big
}

func NewDeductPointBalanceStateValue(amount common.Big) DeductPointBalanceStateValue {
	return DeductPointBalanceStateValue{
		Amount: amount,
	}
}

func (b DeductPointBalanceStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid DeductPointBalanceStateValue")

	if err := util.CheckIsValiders(nil, false, b.Amount); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (b DeductPointBalanceStateValue) HashBytes() []byte {
	return b.Amount.Bytes()
}

func StateKeyPointBalance(contract base.Address, address base.Address) string {
	return fmt.Sprintf("%s:%s:%s", StateKeyPointPrefix(contract), address, PointBalanceSuffix)
}
