package keeper

import (
	"fmt"

	"github.com/cosmos/admin-module/x/adminmodule/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// SubmitProposal create new proposal given a content
func (k Keeper) SubmitProposalLegacy(ctx sdk.Context, content govv1beta1types.Content) (govv1beta1types.Proposal, error) {
	if !k.rtr.HasRoute(content.ProposalRoute()) {
		return govv1beta1types.Proposal{}, sdkerrors.Wrap(govtypes.ErrNoProposalHandlerExists, content.ProposalRoute())
	}

	cacheCtx, _ := ctx.CacheContext()
	handler := k.rtr.GetRoute(content.ProposalRoute())
	if err := handler(cacheCtx, content); err != nil {
		return govv1beta1types.Proposal{}, sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, err.Error())
	}

	proposalID, err := k.GetProposalIDLegacy(ctx)
	if err != nil {
		return govv1beta1types.Proposal{}, err
	}

	headerTime := ctx.BlockHeader().Time

	// submitTime and depositEndTime would not be used
	proposal, err := govv1beta1types.NewProposal(content, proposalID, headerTime, headerTime)
	if err != nil {
		return govv1beta1types.Proposal{}, err
	}

	k.SetProposalLegacy(ctx, proposal)
	k.InsertActiveProposalQueueLegacy(ctx, proposalID)
	k.SetProposalIDLegacy(ctx, proposalID+1)

	logger := k.Logger(ctx)
	logger.Info(
		"LEGACY proposal processed via keeper",
		"proposal", proposal.ProposalId,

	)

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

	logger := k.Logger(ctx)
	logger.Info(
		"LEGACY proposal set in storage",
		"proposal", proposal.ProposalId,

	)
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

	logger := k.Logger(ctx)
	logger.Info(
		"LEGACY proposal get from storage",
		"proposal", proposalID,

	)

	return proposal, true
}

// InsertActiveProposalQueueLegacy inserts a ProposalID into the active proposal queue
func (k Keeper) InsertActiveProposalQueueLegacy(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ActiveProposalLegacyQueueKey(proposalID), types.GetProposalIDBytes(proposalID))
}

// RemoveFromActiveProposalQueueLegacy removes a proposalID from the Active Proposal Queue
func (k Keeper) RemoveFromActiveProposalQueueLegacy(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ActiveProposalLegacyQueueKey(proposalID))
}

// IterateActiveProposalsQueueLegacy iterates over the proposals in the active proposal queue
// and performs a callback function
func (k Keeper) IterateActiveProposalsQueueLegacy(ctx sdk.Context, cb func(proposal govv1beta1types.Proposal) (stop bool)) {
	iterator := k.ActiveProposalQueueIteratorLegacy(ctx)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID := types.GetProposalIDFromBytes(iterator.Value())
		proposal, found := k.GetProposalLegacy(ctx, proposalID)
		if !found {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

// ActiveProposalQueueIteratorLegacy returns an sdk.Iterator for all the proposals in the Active Queue
func (k Keeper) ActiveProposalQueueIteratorLegacy(ctx sdk.Context) sdk.Iterator {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ActiveProposalLegacyQueuePrefix)
	return prefixStore.Iterator(nil, nil)
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
