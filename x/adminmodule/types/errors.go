package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/adminmodule module sentinel errors
var (
	ErrInvalidGenesis = sdkerrors.Register(ModuleName, 1, "invalid genesis state")
	// this line is used by starport scaffolding # ibc/errors
)
