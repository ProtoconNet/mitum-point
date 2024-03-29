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

type TransferFromCommand struct {
	OperationCommand
	Receiver currencycmds.AddressFlag `arg:"" name:"receiver" help:"point receiver" required:"true"`
	Target   currencycmds.AddressFlag `arg:"" name:"target" help:"target approving" required:"true"`
	Amount   currencycmds.BigFlag     `arg:"" name:"amount" help:"amount to transfer" required:"true"`
	receiver base.Address
	target   base.Address
}

func (cmd *TransferFromCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *TransferFromCommand) parseFlags() error {
	if err := cmd.OperationCommand.parseFlags(); err != nil {
		return err
	}

	receiver, err := cmd.Receiver.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid receiver format, %q", cmd.Receiver.String())
	}
	cmd.receiver = receiver

	target, err := cmd.Target.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid target format, %q", cmd.Target.String())
	}
	cmd.target = target

	return nil
}

func (cmd *TransferFromCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError(utils.ErrStringCreate("transfer-from operation"))

	fact := point.NewTransferFromFact(
		[]byte(cmd.Token),
		cmd.sender, cmd.contract,
		cmd.Currency.CID,
		cmd.receiver,
		cmd.target,
		cmd.Amount.Big,
	)

	op := point.NewTransferFrom(fact)
	if err := op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID()); err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
