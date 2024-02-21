package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/admin-module/x/adminmodule/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// SubmitProposal create new proposal given a content
func (k Keeper) SubmitProposalLegacy(ctx sdk.Context, content govv1beta1types.Content) (govv1beta1types.Proposal, error) {
	if !k.legacyRouter.HasRoute(content.ProposalRoute()) {
		return govv1beta1types.Proposal{}, sdkerrors.Wrap(govtypes.ErrNoProposalHandlerExists, content.ProposalRoute())
	}

	err := content.ValidateBasic()
	if err != nil {
		return govv1beta1types.Proposal{}, sdkerrors.Wrap(err, "failed to validate proposal content")
	}

	proposalID, err := k.GetProposalIDLegacy(ctx)
	if err != nil {
		return govv1beta1types.Proposal{}, err
	}

	headerTime := ctx.BlockHeader().Time
	// submitTime and depositEndTime would not be used
	proposal, err := govv1beta1types.NewProposal(content, proposalID, headerTime, headerTime)
	if err != nil {
		return govv1beta1types.Proposal{}, sdkerrors.Wrap(err, "failed to create proposal struct")
	}

	handler := k.RouterLegacy().GetRoute(proposal.ProposalRoute())
	// The proposal handler may execute state mutating logic depending
	// on the proposal content. If the handler fails, no state mutation
	// is written and the error message is returned.
	err = handler(ctx, content)
	if err == nil {
		proposal.Status = govv1beta1types.StatusPassed

	} else {
		return govv1beta1types.Proposal{}, sdkerrors.Wrap(err, "failed to execute proposal proposal struct")
	}
	k.SetProposalLegacy(ctx, proposal)
	k.SetProposalIDLegacy(ctx, proposalID+1)
	k.AddToArchiveLegacy(ctx, proposal)

	return proposal, nil
}

// GetProposalIDLegacy gets the highest proposal ID
func (k Keeper) GetProposalIDLegacy(ctx sdk.Context) (proposalID uint64, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ProposalIDKeyLegacy)
	if bz == nil {
		return 0, sdkerrors.Wrap(types.ErrInvalidGenesis, "initial proposal ID hasn't been set")
	}

	proposalID = types.GetProposalIDFromBytes(bz)
	return proposalID, nil
}

// SetProposalIDLegacy sets the new proposal ID to the store
func (k Keeper) SetProposalIDLegacy(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ProposalIDKeyLegacy, types.GetProposalIDBytes(proposalID))
}

// SetProposalLegacy set a proposal to store
func (k Keeper) SetProposalLegacy(ctx sdk.Context, proposal govv1beta1types.Proposal) {
	store := ctx.KVStore(k.storeKey)

	bz := k.MustMarshalProposalLegacy(proposal)

	store.Set(types.ProposalLegacyKey(proposal.ProposalId), bz)
}

// GetProposalLegacy get proposal from store by ProposalID
func (k Keeper) GetProposalLegacy(ctx sdk.Context, proposalID uint64) (govv1beta1types.Proposal, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.ProposalLegacyKey(proposalID))
	if bz == nil {
		return govv1beta1types.Proposal{}, false
	}

	var proposal govv1beta1types.Proposal
	k.MustUnmarshalProposalLegacy(bz, &proposal)

	return proposal, true
}

func (k Keeper) MarshalProposalLegacy(proposal govv1beta1types.Proposal) ([]byte, error) {
	bz, err := k.cdc.Marshal(&proposal)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func (k Keeper) UnmarshalProposalLegacy(bz []byte, proposal *govv1beta1types.Proposal) error {
	err := k.cdc.Unmarshal(bz, proposal)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) MustMarshalProposalLegacy(proposal govv1beta1types.Proposal) []byte {
	bz, err := k.MarshalProposalLegacy(proposal)
	if err != nil {
		panic(err)
	}
	return bz
}

func (k Keeper) MustUnmarshalProposalLegacy(bz []byte, proposal *govv1beta1types.Proposal) {
	err := k.UnmarshalProposalLegacy(bz, proposal)
	if err != nil {
		panic(err)
	}
}
