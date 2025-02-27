package bot

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/domino14/macondo/equity"
	"github.com/domino14/macondo/game"
	"github.com/domino14/macondo/kwg"
	"github.com/domino14/macondo/montecarlo"
	"github.com/domino14/macondo/move"
	"github.com/domino14/macondo/movegen"
	"github.com/domino14/macondo/tilemapping"
)

const InferencesSimLimit = 400

// Elite bot uses Monte Carlo simulations to rank plays, plays an endgame,
// a pre-endgame (when ready).

// BestPlay picks the highest play by win percentage. It uses montecarlo
// and some other smart things to figure it out.
func eliteBestPlay(ctx context.Context, p *BotTurnPlayer) (*move.Move, error) {

	var moves []*move.Move
	// First determine what stage of the game we are in.
	tr := p.Game.Bag().TilesRemaining()
	// We don't necessarily know the number of tiles on our opponent's rack.
	opp := p.Game.RackFor(p.Game.NextPlayer()).NumTiles()
	// If this is an annotated game, we may not have full rack info.
	unseen := int(opp) + tr
	// Assume our own rack is fully known, however. So if unseen == 7, the bag
	// is empty and we should assign the oppRack accordingly.
	useEndgame := false
	endgamePlies := 0
	simPlies := 0

	if unseen <= 7 {
		useEndgame = true
		if tr > 0 {
			log.Debug().Msg("assigning all unseen to opp")
			// bag is actually empty. Assign all of unseen to the opponent.
			mls := make([]tilemapping.MachineLetter, tr)
			err := p.Game.Bag().Draw(tr, mls)
			if err != nil {
				return nil, err
			}
			for _, t := range mls {
				p.Game.RackFor(p.Game.NextPlayer()).Add(t)
			}
		}
		// Just some sort of estimate
		endgamePlies = unseen + int(p.Game.RackFor(p.Game.PlayerOnTurn()).NumTiles())
	} else if unseen > 7 && unseen <= 14 {
		// at some point check for the specific case of 1 or 2 PEG when
		// the code is ready.
		moves = p.GenerateMoves(80)
		simPlies = unseen
	} else {
		moves = p.GenerateMoves(40)
		if p.minSimPlies > 2 {
			simPlies = p.minSimPlies
		} else {
			simPlies = 2
		}
	}
	log.Debug().Int("simPlies", simPlies).
		Int("simThreads", p.simThreads).
		Int("endgamePlies", endgamePlies).
		Bool("useEndgame", useEndgame).
		Int("unseen", unseen).
		Int("consideredMoves", len(moves)).Msg("elite-player")

	if useEndgame {
		return endGameBest(ctx, p, endgamePlies)
	} else {
		return nonEndgameBest(ctx, p, simPlies, moves)
	}

}

func endGameBest(ctx context.Context, p *BotTurnPlayer, endgamePlies int) (*move.Move, error) {
	if !hasEndgame(p.botType) {
		// Just return the static best play if we don't have an endgame engine.
		return p.GenerateMoves(1)[0], nil
	}
	gd, err := kwg.Get(p.Game.Config(), p.Game.LexiconName())
	if err != nil {
		return nil, err
	}
	// make a copy of the game.
	gameCopy := p.Game.Copy()
	gameCopy.SetBackupMode(game.SimulationMode)
	gameCopy.SetStateStackLength(endgamePlies)
	gen1 := movegen.NewGordonGenerator(gd, gameCopy.Board(), p.Game.Rules().LetterDistribution())
	gen2 := movegen.NewGordonGenerator(gd, gameCopy.Board(), p.Game.Rules().LetterDistribution())
	err = p.endgamer.Init(gen1, gen2, gameCopy, p.Game.Config())
	if err != nil {
		return nil, err
	}
	v, seq, err := p.endgamer.Solve(ctx, endgamePlies)
	if err != nil {
		return nil, err
	}
	log.Debug().Float32("best-endgame-val", v).Interface("seq", seq).Msg("endgame-solve-done")
	return seq[0], nil
}

func nonEndgameBest(ctx context.Context, p *BotTurnPlayer, simPlies int, moves []*move.Move) (*move.Move, error) {
	// use montecarlo if we have it.
	if !hasSimming(p.botType) {
		return moves[0], nil
	}
	var inferTimeout context.Context
	var cancel context.CancelFunc
	if HasInfer(p.botType) {
		log.Debug().Msg("running inference..")
		p.inferencer.Init(p.Game, p.simmerCalcs, p.Config())
		if p.simThreads != 0 {
			p.inferencer.SetThreads(p.simThreads)
		}
		err := p.inferencer.PrepareFinder(p.Game.RackFor(p.Game.PlayerOnTurn()).TilesOn())
		if err != nil {
			// ignore all errors and move on.
			log.Debug().AnErr("inference-prepare-error", err).Msg("probably-ok")
		} else {
			inferTimeout, cancel = context.WithTimeout(context.Background(),
				time.Duration(5*int(time.Second)))
			defer cancel()
			err = p.inferencer.Infer(inferTimeout)
			if err != nil {
				// ignore all errors and move on.
				log.Debug().AnErr("inference-error", err).Msg("probably-ok")
			}
		}
	}

	p.simmer.Init(p.Game, p.simmerCalcs, p.simmerCalcs[0].(*equity.CombinedStaticCalculator), p.Config())
	if p.simThreads != 0 {
		p.simmer.SetThreads(p.simThreads)
	}
	p.simmer.PrepareSim(simPlies, moves)
	p.simmer.SetStoppingCondition(montecarlo.Stop99)

	if HasInfer(p.botType) && len(p.inferencer.Inferences()) > InferencesSimLimit {
		log.Debug().Int("inferences", len(p.inferencer.Inferences())).Msg("using inferences in sim")
		p.simmer.SetInferences(p.inferencer.Inferences(), montecarlo.InferenceCycle)
	}

	// Simulate is a blocking play:
	err := p.simmer.Simulate(ctx)
	if err != nil {
		return nil, err
	}
	plays := p.simmer.WinningPlays()
	log.Debug().Interface("winning-move", plays[0].Move().String()).Msg("sim-done")
	return plays[0].Move(), nil

}
