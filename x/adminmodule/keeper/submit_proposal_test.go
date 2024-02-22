package keeper_test

import (
	"errors"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/stretchr/testify/require"
)

var TestProposalLegacy = govv1beta1types.NewTextProposal("Test", "description")

type invalidProposalLegacyRoute struct{ govv1beta1types.TextProposal }

func (invalidProposalLegacyRoute) ProposalRoute() string { return "nonexistingroute" }

func TestGetSetProposalLegacy(t *testing.T) {
	_, ctx, keeper := setupMsgServer(t)

	// Init genesis ProposalID
	keeper.SetProposalIDLegacy(sdk.UnwrapSDKContext(ctx), 1)

	tp := TestProposalLegacy
	proposal, err := keeper.SubmitProposalLegacy(sdk.UnwrapSDKContext(ctx), tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	keeper.SetProposalLegacy(sdk.UnwrapSDKContext(ctx), proposal)

	gotProposal, ok := keeper.GetProposalLegacy(sdk.UnwrapSDKContext(ctx), proposalID)
	require.True(t, ok)
	require.True(t, proposal.Equal(gotProposal))
}

func TestSubmitProposalLegacy(t *testing.T) {
	_, ctx, keeper := setupMsgServer(t)

	// Init genesis ProposalID
	keeper.SetProposalIDLegacy(sdk.UnwrapSDKContext(ctx), 1)

	testCases := []struct {
		content     govv1beta1types.Content
		expectedErr error
	}{
		{&govv1beta1types.TextProposal{Title: "title", Description: "description"}, nil},
		// Keeper does not check the validity of title and description, no error
		{&govv1beta1types.TextProposal{Title: "", Description: "description"}, govtypes.ErrInvalidProposalContent},
		{&govv1beta1types.TextProposal{Title: strings.Repeat("1234567890", 100), Description: "description"}, govtypes.ErrInvalidProposalContent},
		{&govv1beta1types.TextProposal{Title: "title", Description: ""}, govtypes.ErrInvalidProposalContent},
		{&govv1beta1types.TextProposal{Title: "title", Description: strings.Repeat("1234567890", 1000)}, nil},
		// error only when invalid route
		{&invalidProposalLegacyRoute{}, govtypes.ErrNoProposalHandlerExists},
	}

	for i, tc := range testCases {
		_, err := keeper.SubmitProposalLegacy(sdk.UnwrapSDKContext(ctx), tc.content)
		require.True(t, errors.Is(tc.expectedErr, err), "tc #%d; got: %v, expected: %v", i, err, tc.expectedErr)
	}
}
