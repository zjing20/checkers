package checkers_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zjing20/checkers/testutil/keeper"
	"github.com/zjing20/checkers/testutil/nullify"
	"github.com/zjing20/checkers/x/checkers"
	"github.com/zjing20/checkers/x/checkers/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		NextGame: &types.NextGame{
			IdValue: 32,
		},
		StoredGameList: []types.StoredGame{
			{
				Index: "0",
			},
			{
				Index: "1",
			},
		},
		PlayerInfoList: []types.PlayerInfo{
			{
				Index: "0",
			},
			{
				Index: "1",
			},
		},
		Leaderboard: &types.Leaderboard{
			Winners: "4",
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.CheckersKeeper(t)
	checkers.InitGenesis(ctx, *k, genesisState)
	got := checkers.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.NextGame, got.NextGame)
	require.ElementsMatch(t, genesisState.StoredGameList, got.StoredGameList)
	require.ElementsMatch(t, genesisState.PlayerInfoList, got.PlayerInfoList)
	require.Equal(t, genesisState.Leaderboard, got.Leaderboard)
	// this line is used by starport scaffolding # genesis/test/assert
}
