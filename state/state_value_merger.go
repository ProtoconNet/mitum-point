package state

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
	"sync"
)

type PointBalanceStateValueMerger struct {
	*common.BaseStateValueMerger
	existing PointBalanceStateValue
	add      common.Big
	remove   common.Big
	sync.Mutex
}

func NewPointBalanceStateValueMerger(height base.Height, key string, st base.State) *PointBalanceStateValueMerger {
	nst := st
	if st == nil {
		nst = common.NewBaseState(base.NilHeight, key, nil, nil, nil)
	}

	s := &PointBalanceStateValueMerger{
		BaseStateValueMerger: common.NewBaseStateValueMerger(height, nst.Key(), nst),
	}

	s.existing = NewPointBalanceStateValue(common.ZeroBig)
	if nst.Value() != nil {
		s.existing = nst.Value().(PointBalanceStateValue) //nolint:forcetypeassert //...
	}
	s.add = common.ZeroBig
	s.remove = common.ZeroBig

	return s
}

func (s *PointBalanceStateValueMerger) Merge(value base.StateValue, ops util.Hash) error {
	s.Lock()
	defer s.Unlock()

	switch t := value.(type) {
	case AddPointBalanceStateValue:
		s.add = s.add.Add(t.Amount)
	case DeductPointBalanceStateValue:
		s.remove = s.remove.Add(t.Amount)
	default:
		return errors.Errorf("unsupported point balance state value, %T", value)
	}

	s.AddOperation(ops)

	return nil
}

func (s *PointBalanceStateValueMerger) CloseValue() (base.State, error) {
	s.Lock()
	defer s.Unlock()

	newValue, err := s.closeValue()
	if err != nil {
		return nil, errors.WithMessage(err, "close PointBalanceStateValueMerger")
	}

	s.BaseStateValueMerger.SetValue(newValue)

	return s.BaseStateValueMerger.CloseValue()
}

func (s *PointBalanceStateValueMerger) closeValue() (base.StateValue, error) {
	existingAmount := s.existing.Amount

	if s.add.OverZero() {
		existingAmount = existingAmount.Add(s.add)
	}

	if s.remove.OverZero() {
		existingAmount = existingAmount.Sub(s.remove)
	}

	return NewPointBalanceStateValue(
		existingAmount,
	), nil
}
