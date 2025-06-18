package engine

import (
	"fmt"
	"sync"

	"github.com/RoaringBitmap/roaring/roaring64"
)

/* Max allowed value of an upper bound. */
const MaxUpperBound = 1 << 32

/* Engine represents engine, that allows to apply filters by sums and prods for subset of \Z_+. */
type Engine struct {
	/* Mutex to disallow filtering at the same time. */
	sync.Mutex

	/* Upper bound of subset of \Z_+. */
	upperBound uint64

	/* Bitmap for marking visited pairs (edges). */
	visitedPairs *roaring64.Bitmap

	/* Set of pairs for next filtering by sums. */
	pairsForSumsFilter []uint64

	/* Set of pairs for next filtering by prods. */
	pairsForProdsFilter []uint64
}

/* Creates and sets new Engine instance. */
func New(upperBound uint64) *Engine {
	checkUpperBound(upperBound)
	engine := new(Engine)

	engine.upperBound = upperBound
	engine.visitedPairs = roaring64.New()

	engine.presetPairsForSumsFilter()
	engine.presetPairsForProdsFilter()

	return engine
}

/* Getter of upperBound field. */
func (engine *Engine) GetUpperBound() uint64 {
	return engine.upperBound
}

/* Checks upperBound for allowed maximum. */
func checkUpperBound(upperBound uint64) {
	if upperBound > MaxUpperBound {
		panic(fmt.Sprint("Given upperBound (", upperBound, ") > allowed (", MaxUpperBound, ")"))
	}
}

/* Converts pair to pair id. */
func (engine *Engine) pairId(a, b uint64) uint64 {
	return engine.upperBound*(a-1) + (b - 1)
}

/* Converts pair id to pair. */
func (engine *Engine) pairById(id uint64) Pair {
	return Pair{1 + id/engine.upperBound, 1 + id%engine.upperBound}
}

/* Calculates sum of pair by its id. */
func (engine *Engine) pairSumById(id uint64) uint64 {
	return (1 + id/engine.upperBound) + (1 + id%engine.upperBound)
}

/* Calculates prod of pair by its id. */
func (engine *Engine) pairProdById(id uint64) uint64 {
	return (1 + id/engine.upperBound) * (1 + id%engine.upperBound)
}

/* Converts list of pairs ids to list of Pairs. */
func (engine *Engine) convert_ps(ps []uint64) (converted []Pair) {
	for _, pairId := range ps {
		converted = append(converted, engine.pairById(pairId))
	}
	return
}
