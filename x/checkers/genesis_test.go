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

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.CheckersKeeper(t)
	checkers.InitGenesis(ctx, *k, genesisState)
	got := checkers.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
