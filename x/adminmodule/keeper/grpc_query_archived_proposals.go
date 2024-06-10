package keeper

import (
	"context"

	"github.com/cosmos/admin-module/v2/x/adminmodule/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ArchivedProposals(goCtx context.Context, req *types.QueryArchivedProposalsRequest) (*types.QueryArchivedProposalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	proposals := k.GetArchivedProposals(sdk.UnwrapSDKContext(goCtx))
	return &types.QueryArchivedProposalsResponse{
		Proposals: proposals,
	}, nil
}

func (k Keeper) ArchivedProposalsLegacy(goCtx context.Context, req *types.QueryArchivedProposalsLegacyRequest) (*types.QueryArchivedProposalsLegacyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	proposals := k.GetArchivedProposalsLegacy(sdk.UnwrapSDKContext(goCtx))
	return &types.QueryArchivedProposalsLegacyResponse{
		ProposalsLegacy: proposals,
	}, nil
}
