package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	rules "github.com/zjing20/checkers/x/checkers/rules"
	"github.com/zjing20/checkers/x/checkers/types"
)

func (k msgServer) CreateGame(goCtx context.Context, msg *types.MsgCreateGame) (*types.MsgCreateGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx.Logger().Info("zjz: -2")
	// TODO: Handling the message
	nextGame, found := k.Keeper.GetNextGame(ctx)
	if !found {
		panic("NextGame not found")
	}

	newIndex := strconv.FormatUint(nextGame.IdValue, 10)

	storedGame := types.StoredGame{
		Creator:   msg.Creator,
		Index:     newIndex,
		Game:      rules.New().String(),
		Turn:      rules.BLACK_PLAYER.Color, // test_zj
		Red:       msg.Red,
		Black:     msg.Black,
		MoveCount: 0,
		Deadline:  types.FormatDeadline(types.GetNextDeadline(ctx)),
		Winner:    rules.NO_PLAYER.Color,
		Wager:     msg.Wager,
		Token:     msg.Token,
	}

	err := storedGame.Validate()
	if err != nil {
		return nil, err
	}

	k.Keeper.SendToFifoTail(ctx, &storedGame, &nextGame)
	k.Keeper.SetStoredGame(ctx, storedGame)

	nextGame.IdValue++
	k.Keeper.SetNextGame(ctx, nextGame)

	ctx.GasMeter().ConsumeGas(types.CreateGameGas, "Create game")

	// What to emit
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.StoredGameEventKey),
			sdk.NewAttribute(types.StoredGameEventCreator, msg.Creator),
			sdk.NewAttribute(types.StoredGameEventIndex, newIndex),
			sdk.NewAttribute(types.StoredGameEventRed, msg.Red),
			sdk.NewAttribute(types.StoredGameEventBlack, msg.Black),
			sdk.NewAttribute(types.StoredGameEventTurn, rules.BLACK_PLAYER.Color), // test_zj
			sdk.NewAttribute(types.StoredGameEventWager, strconv.FormatUint(msg.Wager, 10)),
			sdk.NewAttribute(types.StoredGameEventToken, msg.Token),
		),
	)

	return &types.MsgCreateGameResponse{
		IdValue: newIndex,
	}, nil
}
