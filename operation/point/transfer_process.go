package point

import (
	"context"
	"fmt"
	"sync"

	"github.com/ProtoconNet/mitum-point/utils"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-currency/v3/state/currency"
	extstate "github.com/ProtoconNet/mitum-currency/v3/state/extension"
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
	e := util.StringError(ErrStringPreProcess(*opp))

	fact, ok := op.Fact().(TransferFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(TransferFact{}, op.Fact())))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("sender", fact.Sender().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "contract account cannot transfer point, %s", fact.Sender().String()), nil
	}

	if err := currencystate.CheckExistsState(extstate.StateKeyContractAccount(fact.Contract()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "contract not found, %s", fact.Contract().String()), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("currency", fact.Currency().String(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Receiver()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("receiver", fact.Receiver().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Receiver()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "contract account cannot receive points, %s", fact.Receiver().String()), nil
	}

	g := state.NewStateKeyGenerator(fact.Contract())

	if err := currencystate.CheckExistsState(g.Design(), getStateFunc); err != nil {
		return nil, ErrStateNotFound("point design", fact.Contract().String(), err), nil
	}

	st, err := currencystate.ExistsState(g.PointBalance(fact.Sender()), "key of point balance", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("point balance", utils.JoinStringers(fact.Contract(), fact.Sender()), err), nil
	}

	tb, err := state.StatePointBalanceValue(st)
	if err != nil {
		return nil, ErrStateNotFound("point balance value", utils.JoinStringers(fact.Contract(), fact.Sender()), err), nil
	}

	if tb.Compare(fact.Amount()) < 0 {
		return nil, ErrBaseOperationProcess(
			nil,
			"point balance is less than amount to transfer, %s < %s, %s, %s",
			tb, fact.Amount(), fact.Contract(), fact.Sender(),
		), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, ErrBaseOperationProcess(err, "invalid signing"), nil
	}

	return ctx, nil, nil
}

func (opp *TransferProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError(ErrStringProcess(*opp))

	fact, ok := op.Fact().(TransferFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(TransferFact{}, op.Fact())))
	}

	g := state.NewStateKeyGenerator(fact.Contract())

	var sts []base.StateMergeValue

	v, baseErr, err := calculateCurrencyFee(fact.PointFact, getStateFunc)
	if baseErr != nil || err != nil {
		return nil, baseErr, err
	}
	if len(v) > 0 {
		sts = append(sts, v...)
	}

	st, err := currencystate.ExistsState(g.PointBalance(fact.Sender()), "key of point balance", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("point balance", utils.JoinStringers(fact.Contract(), fact.Sender()), err), nil
	}

	_, err = state.StatePointBalanceValue(st)
	if err != nil {
		return nil, ErrStateNotFound("point balance value", utils.JoinStringers(fact.Contract(), fact.Sender()), err), nil
	}

	sts = append(sts, common.NewBaseStateMergeValue(
		g.PointBalance(fact.Sender()),
		state.NewDeductPointBalanceStateValue(fact.Amount()),
		func(height base.Height, st base.State) base.StateValueMerger {
			return state.NewPointBalanceStateValueMerger(height, g.PointBalance(fact.Sender()), st)
		},
	))

	switch st, found, err := getStateFunc(g.PointBalance(fact.Receiver())); {
	case err != nil:
		return nil, ErrBaseOperationProcess(err, "failed to check point balance, %s, %s", fact.Contract(), fact.Receiver()), nil
	case found:
		_, err = state.StatePointBalanceValue(st)
		if err != nil {
			return nil, ErrBaseOperationProcess(err, "failed to get point balance value from state, %s, %s", fact.Contract(), fact.Receiver()), nil
		}
	}

	sts = append(sts, common.NewBaseStateMergeValue(
		g.PointBalance(fact.Receiver()),
		state.NewAddPointBalanceStateValue(fact.Amount()),
		func(height base.Height, st base.State) base.StateValueMerger {
			return state.NewPointBalanceStateValueMerger(height, g.PointBalance(fact.Receiver()), st)
		},
	))

	return sts, nil, nil
}

func (opp *TransferProcessor) Close() error {
	transferProcessorPool.Put(opp)
	return nil
}
