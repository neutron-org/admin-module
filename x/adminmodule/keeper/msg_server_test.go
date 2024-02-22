package keeper_test

import (
	"context"

	"github.com/cosmos/admin-module/x/adminmodule/keeper"
	"github.com/cosmos/admin-module/x/adminmodule/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer() (types.MsgServer, context.Context, *keeper.Keeper) {
	k, ctx := setupKeeper()
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx), k
}
