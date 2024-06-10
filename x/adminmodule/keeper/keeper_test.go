package keeper_test

import (
	"fmt"

	"github.com/cosmos/admin-module/v2/app"
	"github.com/cosmos/admin-module/v2/x/adminmodule/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupKeeper() (*keeper.Keeper, sdk.Context) {
	testApp := app.GetTestApp()
	return &testApp.AdminmoduleKeeper, testApp.BaseApp.NewContext(false)
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
