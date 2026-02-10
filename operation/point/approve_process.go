package point

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-point/state"
	"github.com/ProtoconNet/mitum-point/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var approveItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ApproveItemProcessor)
	},
}

var approveProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ApproveProcessor)
	},
}

func (Approve) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	// NOTE Process is nil func
	return nil, nil, nil
}

type ApproveItemProcessor struct {
	//h      util.Hash
	sender  base.Address
	item    *ApproveItem
	designs map[string]types.Design
}

func (opp *ApproveItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	e := util.StringError("preprocess ApproveItemProcessor")

	if err := opp.item.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if _, _, _, cErr := cstate.ExistsCAccount(
		opp.item.Approved(), "approved", true, false, getStateFunc); cErr != nil {
		return e.Wrap(
			common.ErrCAccountNA.Wrap(
				errors.Errorf("%v: approved %v is contract account", cErr, opp.item.Approved())))
	}

	keyGenerator := state.NewStateKeyGenerator(opp.item.Contract().String())

	if st, err := cstate.ExistsState(
		keyGenerator.Design(), "design", getStateFunc); err != nil {
		return e.Wrap(common.ErrServiceNF.Wrap(errors.Errorf("point service state for contract account %v",
			opp.item.Contract(),
		)))
	} else if design, err := state.StateDesignValue(st); err != nil {
		return e.Wrap(common.ErrServiceNF.Wrap(errors.Errorf("point service state value for contract account %v",
			opp.item.Contract(),
		)))
	} else if apb := design.Policy().GetApproveBox(opp.sender); apb == nil {
		if opp.item.Amount().IsZero() {
			return e.Wrap(common.ErrValueInvalid.Wrap(errors.Errorf("sender %v has not approved any accounts", opp.sender)))
		}
	} else if aprInfo := apb.GetApproveInfo(opp.item.Approved()); aprInfo == nil {
		if opp.item.Amount().IsZero() {
			return e.Wrap(common.ErrValueInvalid.Wrap(errors.Errorf("approved account %v has not been approved",
				opp.item.Approved())))
		}
	}
	if err := cstate.CheckExistsState(keyGenerator.PointBalance(opp.sender.String()), getStateFunc); err != nil {
		return e.Wrap(
			common.ErrStateNF.Wrap(
				errors.Errorf("point balance for sender %v in contract account %v", opp.sender, opp.item.Contract())))
	}

	return nil
}

func (opp *ApproveItemProcessor) Process(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	e := util.StringError("preprocess ApproveItemProcessor")

	//g := tstate.NewStateKeyGenerator(opp.item.Contract().String())
	var sts []base.StateMergeValue

	smv, err := cstate.CreateNotExistAccount(opp.item.Approved(), getStateFunc)
	if err != nil {
		return nil, e.Wrap(err)
	} else if smv != nil {
		sts = append(sts, smv)
	}

	design, _ := opp.designs[opp.item.Contract().String()]

	//st, _ := cstate.ExistsState(g.Design(), "design", getStateFunc)
	//design, _ := tstate.StateDesignValue(st)
	apb := design.Policy().GetApproveBox(opp.sender)
	if apb == nil {
		a := types.NewApproveBox(opp.sender, []types.ApproveInfo{types.NewApproveInfo(opp.item.Approved(), opp.item.Amount())})
		apb = &a
	} else {
		if opp.item.Amount().IsZero() {
			err := apb.RemoveApproveInfo(opp.item.Approved())
			if err != nil {
				return nil, e.Wrap(errors.Errorf("remove approved, %s: %w", opp.item.Approved().String(), err))
			}
		} else {
			apbInfo := apb.GetApproveInfo(opp.item.Approved())
			if apbInfo == nil {
				apb.SetApproveInfo(types.NewApproveInfo(opp.item.Approved(), opp.item.Amount()))
			} else {
				apb.SetApproveInfo(types.NewApproveInfo(opp.item.Approved(), apbInfo.Amount().Add(opp.item.Amount())))
			}

		}
	}

	policy := design.Policy()
	policy.MergeApproveBox(*apb)
	if err := policy.IsValid(nil); err != nil {
		return nil, ErrInvalid(policy, err)
	}
	de := types.NewDesign(design.Symbol(), design.Name(), design.Decimal(), policy)
	if err := de.IsValid(nil); err != nil {
		return nil, ErrInvalid(de, err)
	}
	opp.designs[opp.item.Contract().String()] = de

	//sts = append(sts, cstate.NewStateMergeValue(
	//	g.Design(),
	//	tstate.NewDesignStateValue(de),
	//))

	return sts, nil
}

func (opp *ApproveItemProcessor) Close() {
	//opp.h = nil
	opp.item = nil

	approveItemProcessorPool.Put(opp)
}

type ApproveProcessor struct {
	*base.BaseOperationProcessor
}

func NewApproveProcessor() ctypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("create new ApproveProcessor")

		nopp := approveProcessorPool.Get()
		opp, ok := nopp.(*ApproveProcessor)
		if !ok {
			return nil, e.Wrap(errors.Errorf("expected ApproveProcessor, not %T", nopp))
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

func (opp *ApproveProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(ApproveFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMTypeMismatch).Errorf("expected %T, not %T", ApproveFact{}, op.Fact()),
		), nil
	}

	for i := range fact.items {
		tip := approveItemProcessorPool.Get()
		t, ok := tip.(*ApproveItemProcessor)
		if !ok {
			return ctx, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(
					common.ErrMTypeMismatch).Errorf("expected %T, not %T", &ApproveItemProcessor{}, tip)), nil

		}

		item := fact.items[i]
		t.sender = fact.Sender()
		t.item = &item

		if err := t.PreProcess(ctx, op, getStateFunc); err != nil {
			return ctx, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.
					Errorf("%v", err)), nil
		}
		t.Close()

	}

	return ctx, nil, nil
}

func (opp *ApproveProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, ok := op.Fact().(ApproveFact)
	if !ok {
		return nil, base.NewBaseOperationProcessReasonError("expected %T, not %T", ApproveFact{}, op.Fact()), nil
	}

	designs := make(map[string]types.Design)
	for i := range fact.items {
		keyGenerator := state.NewStateKeyGenerator(fact.items[i].Contract().String())
		st, _ := cstate.ExistsState(keyGenerator.Design(), "design", getStateFunc)
		design, _ := state.StateDesignValue(st)
		if _, found := designs[fact.items[i].contract.String()]; !found {
			designs[fact.items[i].contract.String()] = *design
		}
	}

	var stateMergeValues []base.StateMergeValue // nolint:prealloc
	for i := range fact.items {
		cip := approveItemProcessorPool.Get()
		c, ok := cip.(*ApproveItemProcessor)
		if !ok {
			return nil, base.NewBaseOperationProcessReasonError("expected %T, not %T", &ApproveItemProcessor{}, cip), nil
		}

		item := fact.items[i]
		c.sender = fact.Sender()
		c.item = &item
		c.designs = designs

		s, err := c.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("process approve item: %w", err), nil
		}
		stateMergeValues = append(stateMergeValues, s...)

		c.Close()
	}

	for ca, de := range designs {
		g := state.NewStateKeyGenerator(ca)
		stateMergeValues = append(stateMergeValues, cstate.NewStateMergeValue(
			g.Design(),
			state.NewDesignStateValue(de),
		))
	}

	return stateMergeValues, nil, nil
}

func (opp *ApproveProcessor) Close() error {
	approveProcessorPool.Put(opp)

	return nil
}
