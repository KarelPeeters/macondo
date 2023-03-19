package montecarlo

import (
	"testing"

	"github.com/matryer/is"
)

// func TestPassZTest(t *testing.T) {
// 	is := is.New(t)
// 	is.Equal(passTest(450, 10000, 460, 10000, Stop95), false)
// 	is.True(passTest(450, 10, 400, 5, Stop95))
// 	is.True(passTest(450, 10, 400, 5, Stop99))
// 	is.Equal(passTest(450, 10, 450, 5, Stop95), false)
// 	// 53% win chances with a stdev of 0.01 beats 50% win chances with a stdev of 0.01
// 	// at the 95% confidence level, but not at the 99% confidence level.
// 	is.True(passTest(0.53, 0.0001, 0.50, 0.0001, Stop95))
// 	is.Equal(passTest(0.53, 0.0001, 0.50, 0.0001, Stop99), false)
// }

// func TestZVal(t *testing.T) {
// 	is := is.New(t)
// 	is.Equal(zValStdev(10, 5, 10, 2), float64(0))
// 	is.True(math.Abs(zValStdev(450, 100, 460, 100)-(0.07071)) < 0.0001)
// }

func TestPassTest(t *testing.T) {
	is := is.New(t)
	is.True(passTest(30, 1, 27.9, 1))
	is.True(!passTest(30, 1, 28.0, 1))
}
