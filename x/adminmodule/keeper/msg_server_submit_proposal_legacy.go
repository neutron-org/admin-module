package keeper

import (
	"context"
	"errors"

	"fmt"

	"github.com/cosmos/admin-module/x/adminmodule/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (k msgServer) SubmitProposalLegacy(goCtx context.Context, msg *types.MsgSubmitProposalLegacy) (*types.MsgSubmitProposalLegacyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AdminKey))
	storeCreator := store.Get([]byte(msg.Proposer))
	if storeCreator == nil {
		return nil, fmt.Errorf("proposer %s must be admin to submit proposals to admin-module", msg.Proposer)
	}

	content := msg.GetContent()
	if !k.Keeper.IsProposalTypeWhitelisted(content) {
		return nil, errors.New("proposal content is not whitelisted")
	}

	proposal, err := k.Keeper.SubmitProposalLegacy(ctx, content)
	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounter(1, types.ModuleName, "proposal legacy")

	submitEvent := sdk.NewEvent(types.EventTypeSubmitAdminProposal, sdk.NewAttribute(govtypes.AttributeKeyProposalType, msg.GetContent().ProposalType()))
	ctx.EventManager().EmitEvent(submitEvent)

	return &types.MsgSubmitProposalLegacyResponse{
		ProposalId: proposal.ProposalId,
	}, nil
}
