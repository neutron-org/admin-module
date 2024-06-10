package keeper

import (
	"context"

	"github.com/cosmos/admin-module/v2/x/adminmodule/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Admins(goCtx context.Context, req *types.QueryAdminsRequest) (*types.QueryAdminsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryAdminsResponse{
		Admins: k.GetAdmins(sdk.UnwrapSDKContext(goCtx)),
	}, nil
}
