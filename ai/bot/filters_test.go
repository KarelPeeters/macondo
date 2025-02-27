package bot

import (
	"testing"

	"github.com/domino14/macondo/config"
	"github.com/domino14/macondo/tilemapping"
	"github.com/matryer/is"
)

func TestCombinations(t *testing.T) {
	is := is.New(t)
	cfg := config.DefaultConfig()
	ld, err := tilemapping.EnglishLetterDistribution(&cfg)
	is.NoErr(err)

	scc := createSubCombos(ld)
	cmbs := combinations(ld, scc, []tilemapping.MachineLetter{1, 5, 8, 10}, true)
	is.Equal(cmbs, uint64(1121))
}
