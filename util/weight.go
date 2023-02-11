package util

import (
	"errors"
	"github.com/linguohua/titan/api"
	"math"
	"math/rand"
	"sort"
	"time"
)

// Possible errors returned by NewChooser, preventing the creation of a Chooser
// with unsafe runtime states.
var (
	// If the sum of provided api.DownloadInfoResult weights exceed the maximum integer value
	// for the current platform (e.g. math.MaxInt32 or math.MaxInt64), then
	// the internal running total will overflow, resulting in an imbalanced
	// distribution generating improper results.
	errWeightOverflow = errors.New("sum of Choice Weights exceeds max int")
	// If there are no Choices available to the Chooser with a weight >= 1,
	// there are no valid choices and Pick would produce a runtime panic.
	errNoValidChoices = errors.New("zero Choices with Weight >= 1")
)

// A Chooser caches many possible Choices in a structure designed to improve
// performance on repeated calls for weighted random selection.
type Chooser struct {
	data   []*api.DownloadInfoResult
	totals []int
	max    int
}

// NewChooser initializes a new Chooser for picking from the provided choices.
func NewChooser(choices ...*api.DownloadInfoResult) (*Chooser, error) {
	sort.Slice(choices, func(i, j int) bool {
		return choices[i].Weight < choices[j].Weight
	})

	totals := make([]int, len(choices))
	runningTotal := 0
	for i, c := range choices {
		weight := c.Weight
		if weight < 0 {
			continue // ignore negative weights, can never be picked
		}

		if (math.MaxInt64 - runningTotal) <= weight {
			return nil, errWeightOverflow
		}
		runningTotal += weight
		totals[i] = runningTotal
	}

	if runningTotal < 1 {
		return nil, errNoValidChoices
	}

	return &Chooser{data: choices, totals: totals, max: runningTotal}, nil
}

// Pick returns a single weighted random api.DownloadInfoResult from the Chooser.
// Utilizes global rand as the source of randomness.
func (c *Chooser) Pick() *api.DownloadInfoResult {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(c.max) + 1
	i := sort.SearchInts(c.totals, r)
	return c.data[i]
}
