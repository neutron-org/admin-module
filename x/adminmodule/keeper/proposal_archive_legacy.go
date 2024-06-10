package keeper

import (
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/admin-module/v2/x/adminmodule/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

func (k Keeper) GetArchivedProposalsLegacy(ctx sdk.Context) []*govv1beta1types.Proposal {
	proposals := make([]*govv1beta1types.Proposal, 0)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ArchiveLegacyKey))

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var proposal govv1beta1types.Proposal

		k.MustUnmarshalProposalLegacy(iterator.Value(), &proposal)
		proposals = append(proposals, &proposal)
	}

	return proposals
}

func (k Keeper) AddToArchiveLegacy(ctx sdk.Context, proposal govv1beta1types.Proposal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ArchiveLegacyKey))

	bz := k.MustMarshalProposalLegacy(proposal)

	store.Set(types.ProposalLegacyKey(proposal.ProposalId), bz)
}
