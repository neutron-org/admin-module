package keeper_test

import (
	"cosmossdk.io/store/metrics"
	"fmt"
	dbm "github.com/cosmos/cosmos-db"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/admin-module/x/adminmodule/keeper"
	"github.com/cosmos/admin-module/x/adminmodule/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/stretchr/testify/require"
)

func setupKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	// TODO Add more routes
	rtr := govv1types.NewRouter()
	rtr.AddRoute(govtypes.RouterKey, govv1types.ProposalHandler)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	//cdc := codec.NewProtoCodec(registry)

	types.RegisterInterfaces(registry)
	k := keeper.NewKeeper(
		codec.NewProtoCodec(registry),
		storeKey,
		memStoreKey,
		rtr,
		func(govv1types.Content) bool { return true },
	)
	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	return k, ctx
}

// Using for setting admins before tests
func InitTestAdmins(k *keeper.Keeper, ctx sdk.Context, genesisAdmins []string) error {
	// Removing old admins
	oldAdmins := k.GetAdmins(ctx)
	for _, admin := range oldAdmins {
		err := k.RemoveAdmin(ctx, admin)
		if err != nil {
			return fmt.Errorf("Error removing admin %s\n, error: %s", admin, err)
		}
	}

	// Setting new admins
	for _, admin := range genesisAdmins {
		k.SetAdmin(ctx, admin)
	}
	return nil
}
