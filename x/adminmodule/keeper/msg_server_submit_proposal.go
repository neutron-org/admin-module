package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/admin-module/x/adminmodule/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrortypes "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (k msgServer) SubmitProposal(goCtx context.Context, msg *types.MsgSubmitProposal) (*types.MsgSubmitProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	authority := authtypes.NewModuleAddress(types.ModuleName)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AdminKey))
	storeCreator := store.Get([]byte(msg.Proposer))
	if storeCreator == nil {
		return nil, fmt.Errorf("proposer %s must be admin to submit proposals to admin-module", msg.Proposer)
	}

	msgs, err := msg.GetMsgs()
	if err != nil {
		return nil, errors.Wrap(err, "failed to submit proposal")
	}

	for _, msg := range msgs {
		signers := msg.GetSigners()
		if len(signers) != 1 {
			return nil, fmt.Errorf("should be only 1 signer in message, received: %s", msg.GetSigners())
		}
		if !signers[0].Equals(authority) {
			return nil, errors.Wrap(sdkerrortypes.ErrorInvalidSigner, signers[0].String())
		}
		if !k.Keeper.isMessageWhitelisted(msg) {
			return nil, fmt.Errorf("sdk.Msg is not whitelisted: %s", msg)
		}
	}

	proposal, err := k.Keeper.SubmitProposal(ctx, msgs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to submit proposal")
	}

	defer telemetry.IncrCounter(1, types.ModuleName, "proposal")

	submitEvent := sdk.NewEvent(types.EventTypeSubmitAdminProposal, sdk.NewAttribute(govtypes.AttributeKeyProposalType, types.EventTypeSubmitSdkMessage))
	ctx.EventManager().EmitEvent(submitEvent)

	return &types.MsgSubmitProposalResponse{
		ProposalId: proposal.Id,
	}, nil
}
