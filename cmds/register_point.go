package cmds

import (
	"context"
	"github.com/ProtoconNet/mitum-point/operation/point"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-point/utils"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type RegisterPointCommand struct {
	OperationCommand
	Symbol        PointIDFlag          `arg:"" name:"symbol" help:"point symbol" required:"true"`
	Name          string               `arg:"" name:"name" help:"point name" required:"true"`
	InitialSupply currencycmds.BigFlag `arg:"" name:"initial-supply" help:"initial supply of point" required:"true"`
}

func (cmd *RegisterPointCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *RegisterPointCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError(utils.ErrStringCreate("register-point operation"))

	fact := point.NewRegisterPointFact(
		[]byte(cmd.Token),
		cmd.sender, cmd.contract,
		cmd.Currency.CID, cmd.Symbol.CID,
		cmd.Name,
		cmd.InitialSupply.Big,
	)

	op := point.NewRegisterPoint(fact)
	if err := op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID()); err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
