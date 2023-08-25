package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/admin-module/x/adminmodule/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (k msgServer) SubmitProposal(goCtx context.Context, msg *types.MsgSubmitProposal) (*types.MsgSubmitProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AdminKey))
	storeCreator := store.Get([]byte(msg.Proposer))
	if storeCreator == nil {
		return nil, fmt.Errorf("proposer %s must be admin to submit proposals to admin-module", msg.Proposer)
	}

	msgs, err := msg.GetMsgs()
	if err != nil {
		return nil, err
	}

	for _, msg := range msgs {
		if !k.Keeper.IsMessageWhitelisted(msg) {
			return nil, fmt.Errorf("sdk.Msg is not whitelisted: %s", msg)
		}
	}

	proposal, err := k.Keeper.SubmitProposal(ctx, msgs)
	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounter(1, types.ModuleName, "proposal")

	submitEvent := sdk.NewEvent(types.EventTypeSubmitAdminProposal, sdk.NewAttribute(govtypes.AttributeKeyProposalType, "TODO"))
	ctx.EventManager().EmitEvent(submitEvent)

	return &types.MsgSubmitProposalResponse{
		ProposalId: proposal.Id,
	}, nil
}
