package engine

import (
	"github.com/RoaringBitmap/roaring/roaring64"
)

/* Sets up first edge for sums filter. */
func (engine *Engine) presetPairsForSumsFilter() {
	preset := make([]uint64, 0, 4)
	preset = append(preset, engine.pairId(1, 1))
	if engine.upperBound == 2 {
		preset = append(preset, engine.pairId(1, 2))
		preset = append(preset, engine.pairId(2, 2))
	} else if engine.upperBound > 2 {
		n := engine.upperBound
		preset = append(preset, engine.pairId(1, 2))
		preset = append(preset, engine.pairId(n-1, n))
		preset = append(preset, engine.pairId(n, n))
	}
	engine.pairsForSumsFilter = preset
}

/* Sets up first edge for prods filter. */
func (engine *Engine) presetPairsForProdsFilter() {
	preset := make([]uint64, 0)
	n := engine.upperBound

	seen := roaring64.New()
	duplicated := roaring64.New()

	for a := uint64(1); a <= n; a++ {
		for b := a; b <= n; b++ {
			p := a * b
			if seen.Contains(p) {
				duplicated.Add(p)
			} else {
				seen.Add(p)
			}
		}
	}

	for a := uint64(1); a <= n; a++ {
		for b := a; b <= n; b++ {
			if !duplicated.Contains(a * b) {
				preset = append(preset, engine.pairId(a, b))
			}
		}
	}

	engine.pairsForProdsFilter = preset
}
