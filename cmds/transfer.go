package cmds

import (
	"context"
	"github.com/ProtoconNet/mitum-point/operation/point"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-point/utils"
	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type TransferCommand struct {
	OperationCommand
	Receiver currencycmds.AddressFlag `arg:"" name:"receiver" help:"point receiver" required:"true"`
	Amount   currencycmds.BigFlag     `arg:"" name:"amount" help:"amount to transfer" required:"true"`
	receiver base.Address
}

func (cmd *TransferCommand) Run(pctx context.Context) error { // nolint:dupl
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	currencycmds.PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *TransferCommand) parseFlags() error {
	if err := cmd.OperationCommand.parseFlags(); err != nil {
		return err
	}

	receiver, err := cmd.Receiver.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid receiver format, %q", cmd.Receiver.String())
	}
	cmd.receiver = receiver

	return nil
}

func (cmd *TransferCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError(utils.ErrStringCreate("transfer operation"))

	fact := point.NewTransferFact(
		[]byte(cmd.Token),
		cmd.sender, cmd.contract,
		cmd.Currency.CID,
		cmd.receiver,
		cmd.Amount.Big,
	)

	op := point.NewTransfer(fact)
	if err := op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID()); err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
