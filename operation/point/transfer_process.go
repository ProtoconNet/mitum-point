package point

import (
	"context"
	"fmt"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-point/state"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var transferItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TransferItemProcessor)
	},
}

var transferProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TransferProcessor)
	},
}

func (Transfer) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	// NOTE Process is nil func
	return nil, nil, nil
}

type TransferItemProcessor struct {
	//h      util.Hash
	sender base.Address
	item   *TransferItem
}

func (opp *TransferItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	e := util.StringError("preprocess TransferItemProcessor")

	if err := opp.item.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if _, _, _, cErr := cstate.ExistsCAccount(
		opp.item.Receiver(), "receiver", true, false, getStateFunc); cErr != nil {
		return e.Wrap(
			common.ErrCAccountNA.Wrap(errors.Errorf("%v: receiver %v is contract account", cErr, opp.item.Receiver())))
	}

	return nil
}

func (opp *TransferItemProcessor) Process(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	e := util.StringError("preprocess TransferItemProcessor")

	g := state.NewStateKeyGenerator(opp.item.Contract().String())
	var sts []base.StateMergeValue
	receiver := opp.item.Receiver()
	amount := opp.item.Amount()
	smv, err := cstate.CreateNotExistAccount(receiver, getStateFunc)
	if err != nil {
		return nil, e.Wrap(err)
	} else if smv != nil {
		sts = append(sts, smv)
	}

	switch st, found, err := getStateFunc(g.PointBalance(receiver.String())); {
	case err != nil:
		return nil, e.Wrap(err)
	case found:
		_, err := state.StatePointBalanceValue(st)
		if err != nil {
			return nil, e.Wrap(err)
		}
	}

	sts = append(sts, common.NewBaseStateMergeValue(
		g.PointBalance(receiver.String()),
		state.NewAddPointBalanceStateValue(amount),
		func(height base.Height, st base.State) base.StateValueMerger {
			return state.NewPointBalanceStateValueMerger(height, g.PointBalance(receiver.String()), st)
		},
	))

	return sts, nil
}

func (opp *TransferItemProcessor) Close() {
	//opp.h = nil
	opp.item = nil

	transferItemProcessorPool.Put(opp)
}

type TransferProcessor struct {
	*base.BaseOperationProcessor
	//required map[string]common.Big
}

func NewTransferProcessor() types.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("create new TransferProcessor")

		nopp := transferProcessorPool.Get()
		opp, ok := nopp.(*TransferProcessor)
		if !ok {
			return nil, e.Wrap(errors.Errorf("expected TransferProcessor, not %T", nopp))
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
			common.ErrMPreProcess.Wrap(common.ErrMTypeMismatch).Errorf("expected %T, not %T", TransferFact{}, op.Fact()),
		), nil
	}

	required := make(map[string]common.Big)
	for i := range fact.Items() {
		v, found := required[fact.Items()[i].contract.String()]
		if !found {
			required[fact.Items()[i].contract.String()] = fact.Items()[i].Amount()
		} else {
			required[fact.Items()[i].contract.String()] = v.Add(fact.Items()[i].Amount())
		}
	}

	_, err := PrepareSenderTotalAmounts(fact.Sender().String(), required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("process Transfer; %w", err), nil
	}

	var wg sync.WaitGroup
	errChan := make(chan *base.BaseOperationProcessReasonError, len(fact.Items()))
	for i := range fact.Items() {
		wg.Add(1)
		go func(item TransferItem) {
			defer wg.Done()
			tip := transferItemProcessorPool.Get()
			t, ok := tip.(*TransferItemProcessor)
			if !ok {
				err := base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.Wrap(
						common.ErrMTypeMismatch).Errorf("expected %T, not %T", &TransferItemProcessor{}, tip))
				errChan <- &err
				return
			}

			t.sender = fact.Sender()
			t.item = &item

			if err := t.PreProcess(ctx, op, getStateFunc); err != nil {
				err := base.NewBaseOperationProcessReasonError(common.ErrMPreProcess.Errorf("%v", err))
				errChan <- &err
				return
			}
			t.Close()
		}(fact.Items()[i])
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return nil, *err, nil
		}
	}

	return ctx, nil, nil
}

func (opp *TransferProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, ok := op.Fact().(TransferFact)
	if !ok {
		return nil, base.NewBaseOperationProcessReasonError("expected %T, not %T", TransferFact{}, op.Fact()), nil
	}

	var stateMergeValues []base.StateMergeValue // nolint:prealloc
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan *base.BaseOperationProcessReasonError, len(fact.Items()))
	for i := range fact.Items() {
		wg.Add(1)
		go func(item TransferItem) {
			defer wg.Done()
			cip := transferItemProcessorPool.Get()
			c, ok := cip.(*TransferItemProcessor)
			if !ok {
				err := base.NewBaseOperationProcessReasonError("expected %T, not %T", &TransferItemProcessor{}, cip)
				errChan <- &err
				return
			}

			c.sender = fact.Sender()
			c.item = &item

			s, err := c.Process(ctx, op, getStateFunc)
			if err != nil {
				err := base.NewBaseOperationProcessReasonError("process transfer item: %w", err)
				errChan <- &err
				return
			}
			mu.Lock()
			stateMergeValues = append(stateMergeValues, s...)
			mu.Unlock()
			c.Close()
		}(fact.Items()[i])
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()
	for err := range errChan {
		if err != nil {
			return nil, *err, nil
		}
	}

	required := make(map[string]common.Big)
	for i := range fact.Items() {
		v, found := required[fact.Items()[i].contract.String()]
		if !found {
			required[fact.Items()[i].contract.String()] = fact.Items()[i].amount
		} else {
			required[fact.Items()[i].contract.String()] = v.Add(fact.Items()[i].amount)
		}
	}
	totalAmounts, _ := PrepareSenderTotalAmounts(fact.Sender().String(), required, getStateFunc)

	for key, total := range totalAmounts {
		stateMergeValues = append(
			stateMergeValues,
			common.NewBaseStateMergeValue(
				key,
				state.NewDeductPointBalanceStateValue(total),
				func(height base.Height, st base.State) base.StateValueMerger {
					return state.NewPointBalanceStateValueMerger(height, key, st)
				}),
		)
	}

	return stateMergeValues, nil, nil
}

func (opp *TransferProcessor) Close() error {
	transferProcessorPool.Put(opp)

	return nil
}

func PrepareSenderTotalAmounts(
	holder string,
	required map[string]common.Big,
	getStateFunc base.GetStateFunc,
) (map[string]common.Big, error) {
	totalAmounts := map[string]common.Big{}

	for ca, rq := range required {
		g := state.NewStateKeyGenerator(ca)
		if err := cstate.CheckExistsState(g.Design(), getStateFunc); err != nil {
			return nil, common.ErrServiceNF.Wrap(errors.Errorf("point service state for contract account %v", ca))
		}

		st, err := cstate.ExistsState(g.PointBalance(holder), fmt.Sprintf("point balance, %s", holder), getStateFunc)
		if err != nil {
			return nil, err
		}

		am, err := state.StatePointBalanceValue(st)
		if err != nil {
			return nil, err
		}
		if am.Compare(rq) < 0 {
			return nil, errors.Errorf(
				"point balance of holder %s is less than amount to transfer in contract account %s, %v < %v",
				holder, ca, am, rq)
		}

		totalAmounts[g.PointBalance(holder)] = rq
	}

	return totalAmounts, nil
}
