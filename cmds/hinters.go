package cmds

import (
	ccmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-point/operation/point"
	"github.com/ProtoconNet/mitum-point/state"
	"github.com/ProtoconNet/mitum-point/types"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var AddedHinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: types.ApproveBoxHint, Instance: types.ApproveBox{}},
	{Hint: types.ApproveInfoHint, Instance: types.ApproveInfo{}},
	{Hint: types.PolicyHint, Instance: types.Policy{}},
	{Hint: types.DesignHint, Instance: types.Design{}},

	{Hint: state.DesignStateValueHint, Instance: state.DesignStateValue{}},
	{Hint: state.PointBalanceStateValueHint, Instance: state.PointBalanceStateValue{}},

	{Hint: point.RegisterModelHint, Instance: point.RegisterModel{}},
	{Hint: point.MintHint, Instance: point.Mint{}},
	{Hint: point.BurnHint, Instance: point.Burn{}},
	{Hint: point.ApproveHint, Instance: point.Approve{}},
	{Hint: point.ApprovesHint, Instance: point.Approves{}},
	{Hint: point.ApprovesItemHint, Instance: point.ApprovesItem{}},
	{Hint: point.TransferHint, Instance: point.Transfer{}},
	{Hint: point.TransfersHint, Instance: point.Transfers{}},
	{Hint: point.TransfersItemHint, Instance: point.TransfersItem{}},
	{Hint: point.TransferFromHint, Instance: point.TransferFrom{}},
	{Hint: point.TransfersFromHint, Instance: point.TransfersFrom{}},
	{Hint: point.TransfersFromItemHint, Instance: point.TransfersFromItem{}},
}

var AddedSupportedHinters = []encoder.DecodeDetail{
	{Hint: point.RegisterModelFactHint, Instance: point.RegisterModelFact{}},
	{Hint: point.MintFactHint, Instance: point.MintFact{}},
	{Hint: point.BurnFactHint, Instance: point.BurnFact{}},
	{Hint: point.ApproveFactHint, Instance: point.ApproveFact{}},
	{Hint: point.ApprovesFactHint, Instance: point.ApprovesFact{}},
	{Hint: point.TransferFactHint, Instance: point.TransferFact{}},
	{Hint: point.TransfersFactHint, Instance: point.TransfersFact{}},
	{Hint: point.TransferFromFactHint, Instance: point.TransferFromFact{}},
	{Hint: point.TransfersFromFactHint, Instance: point.TransfersFromFact{}},
}

func init() {
	Hinters = append(Hinters, ccmds.Hinters...)
	Hinters = append(Hinters, AddedHinters...)

	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, ccmds.SupportedProposalOperationFactHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, AddedSupportedHinters...)
}

func LoadHinters(encs *encoder.Encoders) error {
	for i := range Hinters {
		if err := encs.AddDetail(Hinters[i]); err != nil {
			return errors.Wrap(err, "add hinter to encoder")
		}
	}

	for i := range SupportedProposalOperationFactHinters {
		if err := encs.AddDetail(SupportedProposalOperationFactHinters[i]); err != nil {
			return errors.Wrap(err, "add supported proposal operation fact hinter to encoder")
		}
	}

	return nil
}
