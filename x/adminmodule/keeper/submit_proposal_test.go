package keeper_test

import (
	"testing"

	"github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/admin-module/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetProposal(t *testing.T) {
	testApp := app.GetTestApp()
	keeper := testApp.AdminmoduleKeeper
	bankKeeper := testApp.BankKeeper

	acc1 := sdk.AccAddress("acc1")
	acc2 := sdk.AccAddress("acc2")
	coins := sdk.NewCoins(sdk.NewInt64Coin("denom", 10))

	ctx := testApp.NewContext(false, types.Header{})

	keeper.SetProposalID(sdk.UnwrapSDKContext(ctx), 1)

	if err := bankKeeper.MintCoins(ctx, banktypes.ModuleName, coins); err != nil {
		t.Fatal(err.Error())
	}

	if err := bankKeeper.SendCoinsFromModuleToAccount(ctx, banktypes.ModuleName, acc1, coins); err != nil {
		t.Fatal(err.Error())
	}

	msgs := []sdk.Msg{banktypes.NewMsgSend(acc1, acc2, coins)}
	proposal, err := keeper.SubmitProposal(sdk.UnwrapSDKContext(ctx), msgs)
	require.NoError(t, err)

	proposalID := proposal.Id

	gotProposal, ok := keeper.GetProposal(sdk.UnwrapSDKContext(ctx), proposalID)
	require.True(t, ok)
	require.Equal(t, proposal, gotProposal)
}
