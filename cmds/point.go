package cmds

type PointCommand struct {
	RegisterPoint RegisterPointCommand `cmd:"" name:"register-point" help:"register point to contract account"`
	Mint          MintCommand          `cmd:"" name:"mint" help:"mint point to receiver"`
	Burn          BurnCommand          `cmd:"" name:"burn" help:"burn point of target"`
	Approve       ApproveCommand       `cmd:"" name:"approve" help:"approve point to approved account"`
	Transfer      TransferCommand      `cmd:"" name:"transfer" help:"transfer point to receiver"`
}
