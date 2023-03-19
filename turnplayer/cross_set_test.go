package turnplayer

import (
	"testing"

	"github.com/domino14/macondo/board"
	"github.com/domino14/macondo/config"
	"github.com/domino14/macondo/game"
	pb "github.com/domino14/macondo/gen/api/proto/macondo"
	"github.com/matryer/is"
	"github.com/stretchr/testify/assert"
)

var DefaultConfig = config.DefaultConfig()

func compareCrossScores(t *testing.T, b1 *board.GameBoard, b2 *board.GameBoard) {
	dim := b1.Dim()
	dirs := []board.BoardDirection{board.HorizontalDirection, board.VerticalDirection}
	var d board.BoardDirection

	for r := 0; r < dim; r++ {
		for c := 0; c < dim; c++ {
			for _, d = range dirs {
				cs1 := b1.GetCrossScore(r, c, d)
				cs2 := b2.GetCrossScore(r, c, d)
				assert.Equal(t, cs1, cs2)
			}
		}
	}
}

type testMove struct {
	coords string
	word   string
	rack   string
	score  int
}

func TestCompareGameMove(t *testing.T) {
	is := is.New(t)
	opts := &GameOptions{
		ChallengeRule: pb.ChallengeRule_SINGLE,
	}
	players := []*pb.PlayerInfo{
		{Nickname: "JD", RealName: "Jesse"},
		{Nickname: "cesar", RealName: "César"},
	}

	// gen1 := cross_set.GaddagCrossSetGenerator{Dist: dist, Gaddag: gd}
	// gen2 := cross_set.CrossScoreOnlyGenerator{Dist: dist}

	rules1, err := game.NewBasicGameRules(
		&DefaultConfig, "America", board.CrosswordGameLayout, "english",
		game.CrossScoreAndSet, "")
	rules2, err := game.NewBasicGameRules(
		&DefaultConfig, "America", board.CrosswordGameLayout, "english",
		game.CrossScoreOnly, "")

	var testCases = []testMove{
		{"8D", "QWERTY", "QWERTYU", 62},
		{"H8", "TAEL", "TAELABC", 4},
		{"D7", "EQUALITY", "EUALITY", 90},
		{"E10", "MINE", "MINEFHI", 24},
		{"C13", "AB", "ABIIOOO", 21},
	}

	game1, err := BaseTurnPlayerFromRules(opts, players, rules1)
	is.NoErr(err)
	game2, err := BaseTurnPlayerFromRules(opts, players, rules2)
	is.NoErr(err)
	// create a move.
	for _, tc := range testCases {
		err = game1.SetCurrentRack(tc.rack)
		is.NoErr(err)
		err = game2.SetCurrentRack(tc.rack)
		is.NoErr(err)
		m1, err := game1.NewPlacementMove(game1.PlayerOnTurn(), tc.coords, tc.word)
		is.NoErr(err)
		m2, err := game2.NewPlacementMove(game2.PlayerOnTurn(), tc.coords, tc.word)
		is.NoErr(err)
		game1.PlayMove(m1, true, 0)
		game2.PlayMove(m2, true, 0)
		compareCrossScores(t, game1.Board(), game2.Board())
		assert.Equal(t, tc.score, m1.Score())
		assert.Equal(t, tc.score, m2.Score())
	}
}
