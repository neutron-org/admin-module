package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/stretchr/testify/require"
)

func TestAddToArchive(t *testing.T) {
	_, ctx, keeper := setupMsgServer(t)
	keeper.SetProposalIDLegacy(sdk.UnwrapSDKContext(ctx), 1)

	tp := &govv1types.TextProposal{Title: "Test", Description: "Test Description"}
	proposal, err := keeper.SubmitProposalLegacy(sdk.UnwrapSDKContext(ctx), tp)
	require.NoError(t, err)

	keeper.AddToArchiveLegacy(sdk.UnwrapSDKContext(ctx), proposal)

	proposals := keeper.GetArchivedProposals(sdk.UnwrapSDKContext(ctx))
	require.True(t, len(proposals) == 1)

	//TODO(zavgorodnii)
	//t.Log(tp, proposals[0].GetContent())
	//require.Equal(t, tp, proposals[0].GetContent())

}
