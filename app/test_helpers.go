package app

import (
	"encoding/json"
	"time"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	tmtypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func GenesisStateWithValSet(app *App) GenesisState {
	privVal := mock.NewPV()
	pubKey, _ := privVal.GetPubKey()
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	senderPrivKey.PubKey().Address()
	acc := authtypes.NewBaseAccountWithAddress(senderPrivKey.PubKey().Address().Bytes())

	//////////////////////
	balances := []banktypes.Balance{}
	genesisState := NewDefaultGenesisState(app.AppCodec)
	genAccs := []authtypes.GenesisAccount{acc}
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec.MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.DefaultPowerReduction
	initValPowers := []abci.ValidatorUpdate{}

	for _, val := range valSet.Validators {
		pk, _ := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		pkAny, _ := codectypes.NewAnyWithValue(pk)
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   sdkmath.LegacyNewDec(1),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec()),
			MinSelfDelegation: sdkmath.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress().String(), sdk.ValAddress(val.Address).String(), sdkmath.LegacyNewDec(1)))

		// add initial validator powers so consumer InitGenesis runs correctly
		pub, _ := val.ToProto()
		initValPowers = append(initValPowers, abci.ValidatorUpdate{
			Power:  val.VotingPower,
			PubKey: pub.PubKey,
		})
	}
	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(stakingtypes.DefaultParams(), validators, delegations)
	genesisState[stakingtypes.ModuleName] = app.AppCodec.MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	for range delegations {
		// add delegated tokens to total supply
		totalSupply = totalSupply.Add(sdk.NewCoin(sdk.DefaultBondDenom, bondAmt))
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, bondAmt)},
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(
		banktypes.DefaultGenesisState().Params,
		balances,
		totalSupply,
		[]banktypes.Metadata{},
		[]banktypes.SendEnabled{},
	)
	genesisState[banktypes.ModuleName] = app.AppCodec.MustMarshalJSON(bankGenesis)

	_, err := tmtypes.PB2TM.ValidatorUpdates(initValPowers)
	if err != nil {
		panic("failed to get vals")
	}

	return genesisState
}

func GetTestApp() *App {
	db := dbm.NewMemDB()
	testApp := New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		sims.EmptyAppOptions{},
	)

	genesisState := GenesisStateWithValSet(&testApp)
	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	testApp.InitChain(
		&abci.RequestInitChain{
			ChainId:         "",
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: sims.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	return &testApp
}
