package point

import (
	"context"
	"fmt"
	"sync"

	"github.com/ProtoconNet/mitum-point/utils"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-point/state"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var transferProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TransferProcessor)
	},
}

func (Transfer) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type TransferProcessor struct {
	*base.BaseOperationProcessor
}

func NewTransferProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		t := TransferProcessor{}
		e := util.StringError(utils.ErrStringCreate(fmt.Sprintf("new %T", t)))

		nopp := transferProcessorPool.Get()
		opp, ok := nopp.(*TransferProcessor)
		if !ok {
			return nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(&t, nopp)))
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e.Wrap(err)
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *TransferProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(TransferFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", TransferFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	_, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCurrencyNF).Errorf("currency id %v", fact.Currency())), nil
	}

	if _, _, _, cErr := currencystate.ExistsCAccount(fact.Receiver(), "receiver", true, false, getStateFunc); cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v: receiver %v is contract account", cErr, fact.Receiver())), nil
	}

	g := state.NewStateKeyGenerator(fact.Contract().String())

	if err := currencystate.CheckExistsState(g.Design(), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateNF).
				Wrap(common.ErrMServiceNF).Errorf("point design for contract account %v",
				fact.Contract(),
			)), nil
	}

	st, err := currencystate.ExistsState(g.PointBalance(fact.Sender().String()), "point balance", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateNF).
				Errorf("point balance of sender %v in contract account %v", fact.Sender(), fact.Contract())), nil
	}

	tb, err := state.StatePointBalanceValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateValInvalid).
				Errorf("point balance of sender %v in contract account %v", fact.Sender(), fact.Contract())), nil
	}

	if tb.Compare(fact.Amount()) < 0 {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMValueInvalid).
				Errorf("point balance of sender %v is less than amount to transfer in contract account %v, %v < %v",
					fact.Sender(), fact.Contract(), tb, fact.Amount())), nil
	}

	return ctx, nil, nil
}

func (opp *TransferProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, _ := op.Fact().(TransferFact)

	g := state.NewStateKeyGenerator(fact.Contract().String())

	var sts []base.StateMergeValue

	sts = append(sts, common.NewBaseStateMergeValue(
		g.PointBalance(fact.Sender().String()),
		state.NewDeductPointBalanceStateValue(fact.Amount()),
		func(height base.Height, st base.State) base.StateValueMerger {
			return state.NewPointBalanceStateValueMerger(height, g.PointBalance(fact.Sender().String()), st)
		},
	))

	smv, err := currencystate.CreateNotExistAccount(fact.Receiver(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("%w", err), nil
	} else if smv != nil {
		sts = append(sts, smv)
	}

	switch st, found, err := getStateFunc(g.PointBalance(fact.Receiver().String())); {
	case err != nil:
		return nil, ErrBaseOperationProcess(err, "failed to check point balance, %s, %s", fact.Contract(), fact.Receiver()), nil
	case found:
		_, err = state.StatePointBalanceValue(st)
		if err != nil {
			return nil, ErrBaseOperationProcess(err, "failed to get point balance value from state, %s, %s", fact.Contract(), fact.Receiver()), nil
		}
	}

	sts = append(sts, common.NewBaseStateMergeValue(
		g.PointBalance(fact.Receiver().String()),
		state.NewAddPointBalanceStateValue(fact.Amount()),
		func(height base.Height, st base.State) base.StateValueMerger {
			return state.NewPointBalanceStateValueMerger(height, g.PointBalance(fact.Receiver().String()), st)
		},
	))

	return sts, nil, nil
}

func (opp *TransferProcessor) Close() error {
	transferProcessorPool.Put(opp)
	return nil
}
