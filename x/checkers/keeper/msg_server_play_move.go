package keeper

import (
	"context"
	"strconv"
	"strings"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	rules "github.com/zjing20/checkers/x/checkers/rules"
	"github.com/zjing20/checkers/x/checkers/types"
)

func (k msgServer) PlayMove(goCtx context.Context, msg *types.MsgPlayMove) (*types.MsgPlayMoveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
    ctx.Logger().Info("zjz:0")
	// TODO: Handling the message

	storedGame, found := k.Keeper.GetStoredGame(ctx, msg.IdValue)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "game not found %s", msg.IdValue)
	}

	// Is the game already won? Can happen when it is forfeited.
	if storedGame.Winner != rules.NO_PLAYER.Color {
		return nil, types.ErrGameFinished
	}

	var player rules.Player
	if strings.Compare(storedGame.Red, msg.Creator) == 0 {
		player = rules.RED_PLAYER
	} else if strings.Compare(storedGame.Black, msg.Creator) == 0 {
		player = rules.BLACK_PLAYER
	} else {
		return nil, types.ErrCreatorNotPlayer
	}
	ctx.Logger().Info(fmt.Sprintf("zjz: player %s", player))
	game, err := storedGame.ParseGame()
	if err != nil {
		panic("test_zj0")//err.Error()) // test_zj
	}
	ctx.Logger().Info("zjz:1")
	if !game.TurnIs(player) {
		return nil, types.ErrNotPlayerTurn
	}
	ctx.Logger().Info("zjz:2")
	// Make the player pay the wager at the beginning
	err = k.Keeper.CollectWager(ctx, &storedGame)
	if err != nil {
		return nil, err
	}
	ctx.Logger().Info("zjz:3")
	ctx.Logger().Info(fmt.Sprintf("zjz:move from :%d,%d to %d,%d ",msg.FromX,msg.FromY,msg.ToX,msg.ToY))
	captured, moveErr := game.Move(
		rules.Pos{
			X: int(msg.FromX),
			Y: int(msg.FromY),
		},
		rules.Pos{
			X: int(msg.ToX),
			Y: int(msg.ToY),
		},
		ctx.Logger(),
	)
	if moveErr != nil {
		return nil, sdkerrors.Wrapf(moveErr, types.ErrWrongMove.Error())
	}
	ctx.Logger().Info("zjz:4")
	storedGame.MoveCount++
	storedGame.Deadline = types.FormatDeadline(types.GetNextDeadline(ctx))
	storedGame.Winner = game.Winner().Color

	// Send to the back of the FIFO
	nextGame, found := k.Keeper.GetNextGame(ctx)
	if !found {
		panic("NextGame not found")
	}
	ctx.Logger().Info("zjz:5")
	if storedGame.Winner == rules.NO_PLAYER.Color {
		k.Keeper.SendToFifoTail(ctx, &storedGame, &nextGame)
	} else {
		k.Keeper.RemoveFromFifo(ctx, &storedGame, &nextGame)

		// Pay the winnings to the winner
		k.Keeper.MustPayWinnings(ctx, &storedGame)
		winnerInfo, _ := k.Keeper.MustRegisterPlayerWin(ctx, &storedGame)
		k.Keeper.MustAddToLeaderboard(ctx, winnerInfo)
	}
	ctx.Logger().Debug("zjz:6")
	// Save for the next play move
	storedGame.Game = game.String()
	storedGame.Turn = game.Turn.Color
	k.Keeper.SetStoredGame(ctx, storedGame)
	k.Keeper.SetNextGame(ctx, nextGame)

	ctx.GasMeter().ConsumeGas(types.PlayMoveGas, "Play a move")

	// What to emit
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.PlayMoveEventKey),
			sdk.NewAttribute(types.PlayMoveEventCreator, msg.Creator),
			sdk.NewAttribute(types.PlayMoveEventIdValue, msg.IdValue),
			sdk.NewAttribute(types.PlayMoveEventCapturedX, strconv.FormatInt(int64(captured.X), 10)),
			sdk.NewAttribute(types.PlayMoveEventCapturedY, strconv.FormatInt(int64(captured.Y), 10)),
			sdk.NewAttribute(types.PlayMoveEventWinner, game.Winner().Color),
		),
	)
	ctx.Logger().Debug("zjz:7")
	return &types.MsgPlayMoveResponse{
		IdValue:   msg.IdValue,
		CapturedX: int64(captured.X),
		CapturedY: int64(captured.Y),
		Winner:    game.Winner().Color,
	}, nil
}
