package cmds

type PointCommand struct {
	RegisterPoint RegisterPointCommand `cmd:"" name:"register-model" help:"register point to contract account"`
	Mint          MintCommand          `cmd:"" name:"mint" help:"mint point to receiver"`
	Burn          BurnCommand          `cmd:"" name:"burn" help:"burn point of target"`
	Approve       ApproveCommand       `cmd:"" name:"approve" help:"approve point to approved account"`
	Approves      ApprovesCommand      `cmd:"" name:"approves" help:"approves point to approved account"`
	Transfer      TransferCommand      `cmd:"" name:"transfer" help:"transfer point to receiver"`
	Transfers     TransfersCommand     `cmd:"" name:"transfers" help:"transfers point to receiver"`
	TransferFrom  TransferFromCommand  `cmd:"" name:"transfer-from" help:"transfer point to receiver from target"`
	TransfersFrom TransfersFromCommand `cmd:"" name:"transfers-from" help:"transfers point to receiver from target"`
}
