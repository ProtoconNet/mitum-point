package point

import (
	"context"
	"fmt"
	"sync"

	"github.com/ProtoconNet/mitum-point/types"
	"github.com/ProtoconNet/mitum-point/utils"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-point/state"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var transferFromProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TransferFromProcessor)
	},
}

func (TransferFrom) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type TransferFromProcessor struct {
	*base.BaseOperationProcessor
}

func NewTransferFromProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		t := TransferFromProcessor{}
		e := util.StringError(utils.ErrStringCreate(fmt.Sprintf("new %T", t)))

		nopp := transferFromProcessorPool.Get()
		opp, ok := nopp.(*TransferFromProcessor)
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

func (opp *TransferFromProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(TransferFromFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", TransferFromFact{}, op.Fact())), nil
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

	if _, _, aErr, cErr := currencystate.ExistsCAccount(
		fact.Sender(), "sender", true, false, getStateFunc); aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v: sender %v is contract account", cErr, fact.Sender())), nil
	}

	_, _, aErr, cErr := currencystate.ExistsCAccount(
		fact.Contract(), "contract", true, true, getStateFunc)
	if aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", cErr)), nil
	}

	if _, _, _, cErr := currencystate.ExistsCAccount(
		fact.Receiver(), "receiver", true, false, getStateFunc); cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v: receiver %v is contract account", cErr, fact.Receiver())), nil
	}

	if _, _, aErr, cErr := currencystate.ExistsCAccount(
		fact.Target(), "target", true, false, getStateFunc); aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v: target %v is contract account", cErr, fact.Target())), nil
	}

	g := state.NewStateKeyGenerator(fact.Contract())

	st, err := currencystate.ExistsState(g.Design(), "design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceNF).Errorf("point design state for contract account %v",
				fact.Contract(),
			)), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceNF).Errorf("point design state value for contract account %v", fact.Contract())), nil
	}

	approveBoxList := design.Policy().ApproveList()

	idx := -1
	for i, apb := range approveBoxList {
		if apb.Account().Equal(fact.Target()) {
			idx = i
			break
		}
	}

	if idx < 0 {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMAccountNAth).
				Errorf("target %v has not approved any accounts in contract account %v",
					fact.Target(), fact.Contract())), nil
	}

	aprInfo := approveBoxList[idx].GetApproveInfo(fact.Sender())
	if aprInfo == nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMAccountNAth).
				Errorf("sender %v has not been approved by target %v in contract account %v",
					fact.Sender(), fact.Target(), fact.Contract())), nil
	}

	if aprInfo.Amount().Compare(fact.Amount()) < 0 {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMValueInvalid).
				Errorf("approved amount of sender %v is less than amount to transfer in contract account %v, %v < %v",
					fact.Sender(), fact.Contract(), aprInfo.Amount(), fact.Amount())), nil
	}

	st, err = currencystate.ExistsState(g.PointBalance(fact.Target()), "point balance", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateNF).
				Errorf("point balance of target %v in contract account %v", fact.Target(), fact.Contract())), nil
	}

	tb, err := state.StatePointBalanceValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateValInvalid).
				Errorf("point balance of target %v in contract account %v", fact.Target(), fact.Contract())), nil
	}

	if tb.Compare(fact.Amount()) < 0 {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMValueInvalid).
				Errorf("point balance of target %v is less than amount to transfer-from in contract account %v, %v < %v",
					fact.Target(), fact.Contract(), tb, fact.Amount())), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMSignInvalid).
				Errorf("%v", err)), nil
	}

	return ctx, nil, nil
}

func (opp *TransferFromProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError(ErrStringProcess(*opp))

	fact, _ := op.Fact().(TransferFromFact)

	g := state.NewStateKeyGenerator(fact.Contract())

	var sts []base.StateMergeValue

	v, baseErr, err := calculateCurrencyFee(fact.PointFact, getStateFunc)
	if baseErr != nil || err != nil {
		return nil, baseErr, err
	}
	if len(v) > 0 {
		sts = append(sts, v...)
	}

	st, _ := currencystate.ExistsState(g.Design(), "design", getStateFunc)
	design, _ := state.StateDesignValue(st)

	approveBoxList := design.Policy().ApproveList()

	idx := -1
	for i, apb := range approveBoxList {
		if apb.Account().Equal(fact.Target()) {
			idx = i
			break
		}
	}

	apb := approveBoxList[idx]
	am := apb.GetApproveInfo(fact.Sender()).Amount().Sub(fact.Amount())

	if am.IsZero() {
		err := apb.RemoveApproveInfo(fact.Sender())
		if err != nil {
			return nil, nil, e.Wrap(err)
		}
	} else {
		apb.SetApproveInfo(types.NewApproveInfo(fact.Sender(), am))
	}

	approveBoxList[idx] = apb

	policy := types.NewPolicy(design.Policy().TotalSupply(), approveBoxList)
	if err := policy.IsValid(nil); err != nil {
		return nil, ErrInvalid(policy, err), nil
	}

	de := types.NewDesign(design.Symbol(), design.Name(), design.Decimal(), policy)
	if err := de.IsValid(nil); err != nil {
		return nil, ErrInvalid(de, err), nil
	}

	sts = append(sts, currencystate.NewStateMergeValue(
		g.Design(),
		state.NewDesignStateValue(de),
	))

	st, err = currencystate.ExistsState(g.PointBalance(fact.Target()), "point balance", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("point balance", utils.JoinStringers(fact.Contract(), fact.Target()), err), nil
	}

	_, err = state.StatePointBalanceValue(st)
	if err != nil {
		return nil, ErrStateNotFound("point balance value", utils.JoinStringers(fact.Contract(), fact.Target()), err), nil
	}

	sts = append(sts, common.NewBaseStateMergeValue(
		g.PointBalance(fact.Target()),
		state.NewDeductPointBalanceStateValue(fact.Amount()),
		func(height base.Height, st base.State) base.StateValueMerger {
			return state.NewPointBalanceStateValueMerger(height, g.PointBalance(fact.Target()), st)
		},
	))

	smv, err := currencystate.CreateNotExistAccount(fact.Receiver(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("%w", err), nil
	} else if smv != nil {
		sts = append(sts, smv)
	}

	switch st, found, err := getStateFunc(g.PointBalance(fact.Receiver())); {
	case err != nil:
		return nil, ErrBaseOperationProcess(err, "failed to check point balance, %s, %s", fact.Contract(), fact.Receiver()), nil
	case found:
		_, err := state.StatePointBalanceValue(st)
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

func (opp *TransferFromProcessor) Close() error {
	transferFromProcessorPool.Put(opp)
	return nil
}
