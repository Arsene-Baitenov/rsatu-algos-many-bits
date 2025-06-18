package engine

import (
	"sync"
)

/* Filters non trivial sums. */
func (engine *Engine) FilterNonTrivialSums() {
	engine.Lock()
	defer engine.Unlock()

	res := engine.filterSingleSumPairs()
	res = engine.markPairs(res)
	engine.pairsForProdsFilter = res
}

/* Filters non trivial prods. */
func (engine *Engine) FilterNonTrivialProds() {
	engine.Lock()
	defer engine.Unlock()

	res := engine.filterSingleProdPairs()
	res = engine.markPairs(res)
	engine.pairsForSumsFilter = res
}

/* Computes set of pairs with a single sum after filtrations. */
func (engine *Engine) ComputePairsBySums() []Pair {
	engine.Lock()
	defer engine.Unlock()

	res := removeDuplicates(engine.filterSingleSumPairs())
	engine.pairsForProdsFilter = res

	return engine.convert_ps(res)
}

/* Computes set of pairs with a single prod after filtrations. */
func (engine *Engine) ComputePairsByProds() []Pair {
	engine.Lock()
	defer engine.Unlock()

	res := removeDuplicates(engine.filterSingleProdPairs())
	engine.pairsForSumsFilter = res

	return engine.convert_ps(res)
}

/* Marks given pairs. */
func (engine *Engine) markPairs(ps []uint64) []uint64 {
	res := make([]uint64, 0, len(ps))
	for _, r := range ps {
		if engine.visitedPairs.Contains(r) {
			continue
		}

		engine.visitedPairs.Add(r)
		res = append(res, r)
	}

	return res
}

/* Removes duplicates from given slice of uint64. */
func removeDuplicates(ps []uint64) (nps []uint64) {
	m := make(map[uint64]bool)
	for _, p := range ps {
		if !m[p] {
			nps = append(nps, p)
			m[p] = true
		}
	}
	return
}

/* Filters pairs with single sums. */
func (engine *Engine) filterSingleSumPairs() []uint64 {
	resCh := make(chan uint64, len(engine.pairsForSumsFilter))
	var wg sync.WaitGroup
	for _, pairId := range engine.pairsForSumsFilter {
		wg.Add(1)
		go func(pairId uint64) {
			defer wg.Done()
			if r, success := engine.processSum(pairId); success {
				resCh <- r
			}
		}(pairId)
	}
	wg.Wait()
	close(resCh)

	res := make([]uint64, 0, len(resCh))
	for r := range resCh {
		res = append(res, r)
	}
	return res
}

/* Checks that given pair has single sum. */
func (engine *Engine) processSum(pairId uint64) (nextPairId uint64, success bool) {
	s := engine.pairSumById(pairId)
	var singlePairId uint64 = 0
	isSingle := false

	for i := uint64(1); i+i <= s; i++ {
		if curr := engine.pairId(i, s-i); s-i <= engine.upperBound && !engine.visitedPairs.Contains(curr) {
			if !isSingle {
				singlePairId, isSingle = curr, true
			} else {
				isSingle = false
				break
			}
		}
	}
	return singlePairId, isSingle
}

/* Filters pairs with single prods. */
func (engine *Engine) filterSingleProdPairs() []uint64 {
	resCh := make(chan uint64, len(engine.pairsForProdsFilter))
	var wg sync.WaitGroup
	for _, pairId := range engine.pairsForProdsFilter {
		wg.Add(1)
		go func(pairId uint64) {
			defer wg.Done()
			if r, success := engine.processProd(pairId); success {
				resCh <- r
			}
		}(pairId)
	}
	wg.Wait()
	close(resCh)

	res := make([]uint64, 0, len(resCh))
	for r := range resCh {
		res = append(res, r)
	}
	return res
}

/* Checks that given pair has single prod. */
func (engine *Engine) processProd(pairId uint64) (nextPairId uint64, success bool) {
	p := engine.pairProdById(pairId)
	var singlePairId uint64 = 0
	isSingle := false

	for i := uint64(1); i*i <= p; i++ {
		if p%i != 0 {
			continue
		}
		if o := p / i; o <= engine.upperBound {
			if curr := engine.pairId(i, o); !engine.visitedPairs.Contains(curr) {
				if !isSingle {
					singlePairId, isSingle = curr, true
				} else {
					isSingle = false
					break
				}
			}
		}
	}
	return singlePairId, isSingle
}
