package keeper

import (
	"fmt"

	"github.com/cosmos/admin-module/x/adminmodule/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// SubmitProposal create new proposal given a content
func (k Keeper) SubmitProposal(ctx sdk.Context, msgs []sdk.Msg) (govv1types.Proposal, error) {
	var events sdk.Events
	proposalID, err := k.GetProposalID(ctx)
	if err != nil {
		return govv1types.Proposal{}, err
	}

	headerTime := ctx.BlockHeader().Time

	proposal, err := govv1types.NewProposal(msgs, proposalID, headerTime, headerTime, "", "", "", nil)
	if err != nil {
		return govv1types.Proposal{}, err
	}

	for idx, msg := range msgs {
		handler := k.Router().Handler(msg)

		var res *sdk.Result
		res, err := handler(ctx, msg)
		if err != nil {
			return proposal, fmt.Errorf("failed to handle %d msg in proposal %d: %w", idx, proposal.Id, err)
		}
		events = append(events, res.GetEvents()...)
	}
	proposal.Status = govv1types.StatusPassed
	k.SetProposal(ctx, proposal)
	k.SetProposalID(ctx, proposalID+1)
	k.AddToArchive(ctx, proposal)

	return proposal, nil
}

// GetProposalID gets the highest proposal ID
func (k Keeper) GetProposalID(ctx sdk.Context) (proposalID uint64, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ProposalIDKey)
	if bz == nil {
		return 0, sdkerrors.Wrap(types.ErrInvalidGenesis, "initial proposal ID hasn't been set")
	}

	proposalID = types.GetProposalIDFromBytes(bz)
	return proposalID, nil
}

// SetProposalID sets the new proposal ID to the store
func (k Keeper) SetProposalID(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ProposalIDKey, types.GetProposalIDBytes(proposalID))
}

// SetProposal set a proposal to store
func (k Keeper) SetProposal(ctx sdk.Context, proposal govv1types.Proposal) {
	store := ctx.KVStore(k.storeKey)

	bz := k.MustMarshalProposal(proposal)

	store.Set(types.ProposalKey(proposal.Id), bz)
}

// GetProposal get proposal from store by ProposalID
func (k Keeper) GetProposal(ctx sdk.Context, proposalID uint64) (govv1types.Proposal, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.ProposalKey(proposalID))
	if bz == nil {
		return govv1types.Proposal{}, false
	}

	var proposal govv1types.Proposal
	k.MustUnmarshalProposal(bz, &proposal)

	return proposal, true
}

func (k Keeper) MarshalProposal(proposal govv1types.Proposal) ([]byte, error) {
	bz, err := k.cdc.Marshal(&proposal)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func (k Keeper) UnmarshalProposal(bz []byte, proposal *govv1types.Proposal) error {
	err := k.cdc.Unmarshal(bz, proposal)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) MustMarshalProposal(proposal govv1types.Proposal) []byte {
	bz, err := k.MarshalProposal(proposal)
	if err != nil {
		panic(err)
	}
	return bz
}

func (k Keeper) MustUnmarshalProposal(bz []byte, proposal *govv1types.Proposal) {
	err := k.UnmarshalProposal(bz, proposal)
	if err != nil {
		panic(err)
	}
}
