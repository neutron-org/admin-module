package keeper

import (
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/admin-module/v2/x/adminmodule/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

func (k Keeper) GetArchivedProposals(ctx sdk.Context) []*govv1types.Proposal {
	proposals := make([]*govv1types.Proposal, 0)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ArchiveKey))

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var proposal govv1types.Proposal

		k.MustUnmarshalProposal(iterator.Value(), &proposal)
		proposals = append(proposals, &proposal)
	}

	return proposals
}

func (k Keeper) AddToArchive(ctx sdk.Context, proposal govv1types.Proposal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ArchiveKey))

	bz := k.MustMarshalProposal(proposal)

	store.Set(types.ProposalKey(proposal.Id), bz)
}
