package keeper

import (
	"github.com/cosmos/admin-module/x/adminmodule/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func (k Keeper) NewMsgServerImpl() types.MsgServer {
	return &msgServer{Keeper: k}
}

var _ types.MsgServer = msgServer{}
