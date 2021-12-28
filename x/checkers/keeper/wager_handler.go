package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	rules "github.com/zjing20/checkers/x/checkers/rules"
	"github.com/zjing20/checkers/x/checkers/types"
//	accountTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Returns an error if the player has not enough funds.
func (k *Keeper) CollectWager(ctx sdk.Context, storedGame *types.StoredGame) error {
	// Make the player pay the wager at the beginning
	if storedGame.MoveCount == 0 {
		// Black plays first
		ctx.Logger().Info("zjz:2.0")
		black, err := storedGame.GetBlackAddress()
		ctx.Logger().Info(black.String())
		ctx.Logger().Info("zjz:2.1")
		if err != nil {
			panic(err.Error())
		}
		ctx.Logger().Info("zjz:2.2")
		// recipientAcc := k.account.GetModuleAccount(ctx, types.ModuleName) // zj_test
		// addr, perms := k.account.GetModuleAddressAndPermissions(types.ModuleName)
		// ctx.Logger().Info(fmt.Sprintf("zjz: perm len %d",len(perms)))
		// if addr == nil || strings.TrimSpace(addr.String()) == "" {
		// 	ctx.Logger().Info("zjz: address empty")
		// }
		// ctx.Logger().Info(fmt.Sprintf("zjz: address not empty %s", addr.String()))
		// acc := k.account.GetAccount(ctx, addr)
		// if acc != nil {
		// 	macc, ok := acc.(accountTypes.ModuleAccountI)
		// 	if !ok {
		// 		ctx.Logger().Info("zjz: no module account")
		// 	}
		// 	ctx.Logger().Info("zjz: if ",macc)
		// } else {
		// 	macc := accountTypes.NewEmptyModuleAccount(types.ModuleName, perms...)
		// 	ctx.Logger().Info(fmt.Sprintf("zjz: else %s",macc.String()))
		// }
		// ctx.Logger().Info("zjz:2.2.5")
		// ctx.Logger().Info(recipientAcc.GetAddress().String()) //zj_test
		// ctx.Logger().Info(black.String()) //zj_test
		// red, err := storedGame.GetRedAddress() //zj_test
		// k.bank.SendCoins(ctx, black, red, sdk.NewCoins(storedGame.GetWagerCoin()))  //zj_test
		err = k.bank.SendCoinsFromAccountToModule(ctx, black, types.ModuleName, sdk.NewCoins(storedGame.GetWagerCoin())) 
		ctx.Logger().Info("zjz:2.3")
		if err != nil {
			return sdkerrors.Wrapf(err, types.ErrBlackCannotPay.Error())
		}
	} else if storedGame.MoveCount == 1 {
		// Red plays second
		// red, err := storedGame.GetBlackAddress() // zj_test
		red, err := storedGame.GetRedAddress()
		if err != nil {
			panic(err.Error())
		}
		err = k.bank.SendCoinsFromAccountToModule(ctx, red, types.ModuleName, sdk.NewCoins(storedGame.GetWagerCoin()))
		if err != nil {
			return sdkerrors.Wrapf(err, types.ErrRedCannotPay.Error())
		}
	}
	return nil
}

// Game must have a valid winner.
func (k *Keeper) MustPayWinnings(ctx sdk.Context, storedGame *types.StoredGame) {
	// Pay the winnings to the winner
	winnerAddress, found, err := storedGame.GetWinnerAddress()
	if err != nil {
		panic(err.Error())
	}
	if !found {
		panic(fmt.Sprintf(types.ErrCannotFindWinnerByColor.Error(), storedGame.Winner))
	}
	winnings := storedGame.GetWagerCoin()
	if storedGame.MoveCount == 0 {
		panic(types.ErrNothingToPay.Error())
	} else if 1 < storedGame.MoveCount {
		winnings = winnings.Add(winnings)
	}
	err = k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, winnerAddress, sdk.NewCoins(winnings))
	if err != nil {
		panic(types.ErrCannotPayWinnings.Error())
	}
}

// Game must be in a state where it can be refunded.
func (k *Keeper) MustRefundWager(ctx sdk.Context, storedGame *types.StoredGame) {
	// Refund wager to black player if red rejects after black played
	if storedGame.MoveCount == 1 {
		black, err := storedGame.GetBlackAddress()
		if err != nil {
			panic(err.Error())
		}
		err = k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, black, sdk.NewCoins(storedGame.GetWagerCoin()))
		if err != nil {
			panic(fmt.Sprintf(types.ErrCannotRefundWager.Error(), rules.BLACK_PLAYER.Color))
		}
	} else if storedGame.MoveCount == 0 {
		// Do nothing
	} else {
		// TODO Implement a draw mechanism.
		panic(fmt.Sprintf(types.ErrNotInRefundState.Error(), storedGame.MoveCount))
	}
}
