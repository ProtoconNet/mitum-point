package point

import (
	"context"
	"fmt"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"sync"

	"github.com/ProtoconNet/mitum-point/types"
	"github.com/ProtoconNet/mitum-point/utils"

	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	extstate "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-point/state"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var registerModelProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RegisterModelProcessor)
	},
}

func (RegisterModel) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type RegisterModelProcessor struct {
	*base.BaseOperationProcessor
}

func NewRegisterModelProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		t := RegisterModelProcessor{}
		e := util.StringError(utils.ErrStringCreate(fmt.Sprintf("new %T", t)))

		nopp := registerModelProcessorPool.Get()
		opp, ok := nopp.(*RegisterModelProcessor)
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

func (opp *RegisterModelProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(RegisterModelFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", RegisterModelFact{}, op.Fact())), nil
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

	g := state.NewStateKeyGenerator(fact.Contract())

	if found, _ := currencystate.CheckNotExistsState(g.Design(), getStateFunc); found {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceE).Errorf("point design for contract account %v", fact.Contract())), nil
	}

	return ctx, nil, nil
}

func (opp *RegisterModelProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError(ErrStringProcess(*opp))

	fact, ok := op.Fact().(RegisterModelFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(RegisterModelFact{}, op.Fact())))
	}

	g := state.NewStateKeyGenerator(fact.Contract())

	var sts []base.StateMergeValue

	policy := types.NewPolicy(fact.InitialSupply(), []types.ApproveBox{})
	if err := policy.IsValid(nil); err != nil {
		return nil, ErrInvalid(policy, err), nil
	}

	design := types.NewDesign(fact.Symbol(), fact.Name(), fact.Decimal(), policy)
	if err := design.IsValid(nil); err != nil {
		return nil, ErrInvalid(design, err), nil
	}

	sts = append(sts, currencystate.NewStateMergeValue(
		g.Design(),
		state.NewDesignStateValue(design),
	))

	st, err := currencystate.ExistsState(extstate.StateKeyContractAccount(fact.Contract()), "contract account", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("contract", fact.Contract().String(), err), nil
	}

	ca, err := extstate.StateContractAccountValue(st)
	if err != nil {
		return nil, ErrStateNotFound("contract value", fact.Contract().String(), err), nil
	}
	nca := ca.SetIsActive(true)

	sts = append(sts, currencystate.NewStateMergeValue(
		extstate.StateKeyContractAccount(fact.Contract()),
		extstate.NewContractAccountStateValue(nca),
	))

	if fact.InitialSupply().OverZero() {
		sts = append(sts, common.NewBaseStateMergeValue(
			g.PointBalance(fact.Sender()),
			state.NewAddPointBalanceStateValue(fact.InitialSupply()),
			func(height base.Height, st base.State) base.StateValueMerger {
				return state.NewPointBalanceStateValueMerger(
					height,
					g.PointBalance(fact.Sender()),
					st,
				)
			},
		))
	}

	return sts, nil, nil
}

func (opp *RegisterModelProcessor) Close() error {
	registerModelProcessorPool.Put(opp)
	return nil
}
