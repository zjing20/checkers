package rules

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	BOARD_DIM = 3
	RED       = "red"
	BLACK     = "black"
	ROW_SEP   = "|"
)

type Player struct {
	Color string
}

type Piece struct {
	Player Player
}

var PieceStrings = map[Player]string{
	RED_PLAYER:   "r",
	BLACK_PLAYER: "b",
	NO_PLAYER:    "*",
}

var NO_PIECE = Piece{NO_PLAYER}

var WINNER = NO_PLAYER
var StringPieces = map[string]Piece{
	"r": Piece{RED_PLAYER},
	"b": Piece{BLACK_PLAYER},
	// "R": Piece{RED_PLAYER},
	// "B": Piece{BLACK_PLAYER},
	"*": NO_PIECE,
}

type Pos struct {
	X int
	Y int
}

var NO_POS = Pos{-1, -1}

var BLACK_PLAYER = Player{BLACK}
var RED_PLAYER = Player{RED}
var NO_PLAYER = Player{
	Color: "NO_PLAYER",
}


var Players = map[string]Player{
	RED:   RED_PLAYER,
	BLACK: BLACK_PLAYER,
}

var Opponents = map[Player]Player{
	BLACK_PLAYER: RED_PLAYER,
	RED_PLAYER:   BLACK_PLAYER,
}


func (game *Game) init() {

	// Initialize usable spaces
	for y := 0; y < BOARD_DIM; y++ {
		for x := 0; x < BOARD_DIM; x++ {
			game.Pieces[Pos{X: x, Y: y}] = Piece{NO_PLAYER}
		}
	}


}

type Game struct {
	Pieces map[Pos]Piece
	Turn   Player
	MoveCount int
}

func New() *Game {
	pieces := make(map[Pos]Piece)
	game := &Game{pieces, BLACK_PLAYER, 0}
	game.init()
	return game
}


func (game *Game) GetPieceAt(pos Pos) Player {
	if pos.X<0 || pos.Y<0 || pos.X >=BOARD_DIM || pos.Y >=BOARD_DIM {
		return NO_PLAYER
	}
	return game.Pieces[pos].Player
}

func (game *Game) TurnIs(player Player) bool {
	return game.Turn == player
}

func (game *Game) Winner() Player {
	return WINNER

}

func (game *Game) GetWinner(pos Pos) Player {
	player := game.GetPieceAt(pos)
	if pos.X==0 && pos.Y==0{
		p1:= Pos{pos.X+1, pos.Y+1}
		p2:= Pos{pos.X+2, pos.Y+2}
		if game.GetPieceAt(p1) == player && game.GetPieceAt(p2) == player {
			return player
		}
	}
	if pos.X==1 && pos.Y==1{
		p1:= Pos{pos.X+1, pos.Y+1}
		p2:= Pos{pos.X-1, pos.Y-1}
		if game.GetPieceAt(p1) == player && game.GetPieceAt(p2) == player {
			return player
		}
	}
	if pos.X==2 && pos.Y==2{
		p1:= Pos{pos.X-1, pos.Y-1}
		p2:= Pos{pos.X-2, pos.Y-2}
		if game.GetPieceAt(p1) == player && game.GetPieceAt(p2) == player {
			return player
		}
	}
	if pos.X==0 {
		p1:= Pos{pos.X+1, pos.Y}
		p2:= Pos{pos.X+2, pos.Y}
		if game.GetPieceAt(p1) == player && game.GetPieceAt(p2) == player {
			return player
		}
	}
	if pos.X==1 {
		p1:= Pos{pos.X+1, pos.Y}
		p2:= Pos{pos.X-2, pos.Y}
		if game.GetPieceAt(p1) == player && game.GetPieceAt(p2) == player {
			return player
		}
	}
	if pos.X==2 {
		p1:= Pos{pos.X-1, pos.Y}
		p2:= Pos{pos.X-2, pos.Y}
		if game.GetPieceAt(p1) == player && game.GetPieceAt(p2) == player {
			return player
		}
	}
	if pos.Y==0 {
		p1:= Pos{pos.X, pos.Y+1}
		p2:= Pos{pos.X, pos.Y+2}
		if game.GetPieceAt(p1) == player && game.GetPieceAt(p2) == player {
			return player
		}
	}
	if pos.Y==1 {
		p1:= Pos{pos.X, pos.Y-1}
		p2:= Pos{pos.X, pos.Y+1}
		if game.GetPieceAt(p1) == player && game.GetPieceAt(p2) == player {
			return player
		}
	}
	if pos.Y==2 {
		p1:= Pos{pos.X, pos.Y-1}
		p2:= Pos{pos.X, pos.Y-2}
		if game.GetPieceAt(p1) == player && game.GetPieceAt(p2) == player {
			return player
		}
	}	
	return NO_PLAYER
}

func (game *Game) ValidMove(dst Pos) bool {
	if dst.X<0 || dst.Y<0 || dst.X >=BOARD_DIM || dst.Y >=BOARD_DIM || game.GetPieceAt(dst)!=NO_PLAYER {
		return false
	}
	return true
}


func (game *Game) updateTurn(dst Pos) {
	game.Turn = Opponents[game.Turn]

}


func (game *Game) playerHasMove(player Player) bool {
	if game.MoveCount < BOARD_DIM * BOARD_DIM {
		return true
	}
	return false
}


func (game *Game) Move(player Player, dst Pos, check bool, logger log.Logger) (captured Pos, err error) {
	logger.Info("zjz: 3.-1")
	logger.Info(game.GetPieceAt(dst).Color)
	if !game.ValidMove(dst) {
		return NO_POS, errors.New(fmt.Sprintf("%v is not a valid position", dst))
	}
	logger.Info("zjz: 3.0")
	logger.Info(game.Turn.Color)
	logger.Info(player.Color)
	if !game.TurnIs(player) {
		return NO_POS, errors.New(fmt.Sprintf("Not %v's turn", player))
	}
	game.Pieces[dst] = Piece{player}
	logger.Info("zjz: 3.1")
	if !check {
		game.MoveCount++
		logger.Info("zjz: 3.2")
		winner:=game.GetWinner(dst)
		if winner != NO_PLAYER {
			WINNER = winner
		}
		logger.Info("zjz: 3.3")
	}
	logger.Info("zjz: 3.4")
	game.updateTurn(dst)
	logger.Info("zjz:3.9")
	logger.Info("zjz:3.10")
	return
}

func (game *Game) String() string {
	var buf bytes.Buffer
	for y := 0; y < BOARD_DIM; y++ {
		for x := 0; x < BOARD_DIM; x++ {
			pos := Pos{x, y}
			if game.GetPieceAt(pos)!= NO_PLAYER {
				piece := game.Pieces[pos]
				val := PieceStrings[piece.Player]
				buf.WriteString(val)
			} else {
				buf.WriteString(PieceStrings[NO_PLAYER])
			}
		}
		if y < (BOARD_DIM - 1) {
			buf.WriteString(ROW_SEP)
		}
	}
	return buf.String()
}

func ParsePiece(s string) (Piece, bool) {
	piece, ok := StringPieces[s]
	return piece, ok
}

func Parse(s string, count int) (*Game, error) {
	if len(s) != BOARD_DIM*BOARD_DIM+(BOARD_DIM-1) {
		return nil, errors.New(fmt.Sprintf("invalid board string: %v", s))
	}
	pieces := make(map[Pos]Piece)
	result := &Game{pieces, BLACK_PLAYER, count}
	for y, row := range strings.Split(s, ROW_SEP) {
		for x, c := range strings.Split(row, "") {
			if x >= BOARD_DIM || y >= BOARD_DIM {
				return nil, errors.New(fmt.Sprintf("invalid board, piece out of bounds: %v, %v", x, y))
			}
			if piece, ok := ParsePiece(c); !ok {
				return nil, errors.New(fmt.Sprintf("invalid board, invalid piece at %v, %v", x, y))
			} else {
				result.Pieces[Pos{x, y}] = piece
			} 
			
		}
	}
	return result, nil
}
