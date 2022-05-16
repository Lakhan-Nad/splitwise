package splits

import (
	"fmt"
	"math"
)

// EXACT SPLIT
type ExactSplitStrategy struct {
}

func (s *ExactSplitStrategy) CalculateShare(total float64, splits []*Split) []float32 {
	m := make([]float32, len(splits))
	for idx, split := range splits {
		m[idx] = split.Value()
	}
	return m
}

func (s *ExactSplitStrategy) ValidateSplits(total float64, splits []*Split) (bool, error) {
	totalValue := 0.0
	for _, split := range splits {
		totalValue += float64(split.Value())
	}
	if totalValue == total {
		return true, nil
	}
	return false, fmt.Errorf("total value of splits %f doesn't match given total amount %f", totalValue, total)
}

func (s *ExactSplitStrategy) Name() string {
	return "EXACT"
}

// PERCENT SPLIT
type PercentSplitStrategy struct {
}

func (s *PercentSplitStrategy) CalculateShare(total float64, splits []*Split) []float32 {
	m := make([]float32, len(splits))
	for idx, split := range splits {
		m[idx] = (float32(total * float64(split.Value()/100)))
	}
	return m
}

func (s *PercentSplitStrategy) ValidateSplits(total float64, splits []*Split) (bool, error) {
	totalPercent := 0.0
	for _, split := range splits {
		totalPercent += float64(split.Value())
	}
	if math.Abs(100.0-totalPercent) == 0.0 {
		return true, nil
	}
	return false, fmt.Errorf("percents doesn't add upto 100, total is %f", totalPercent)
}

func (s *PercentSplitStrategy) Name() string {
	return "PERCENT"
}

// Equal Split
type EqualSplitStrategy struct {
}

func (s *EqualSplitStrategy) CalculateShare(total float64, splits []*Split) []float32 {
	m := make([]float32, len(splits))
	totalCount := float64(len(splits))
	for idx := range splits {
		m[idx] = float32(total / totalCount)
	}
	return m
}

func (s *EqualSplitStrategy) ValidateSplits(total float64, splits []*Split) (bool, error) {
	return true, nil
}

func (s *EqualSplitStrategy) Name() string {
	return "EQUAL"
}

// Share Split
type ShareSplitStrategy struct {
}

func (s *ShareSplitStrategy) CalculateShare(total float64, splits []*Split) []float32 {
	m := make([]float32, len(splits))
	totalCount := 0.0
	for _, split := range splits {
		totalCount += float64(split.Value())
	}
	for idx, s := range splits {
		m[idx] = float32(total * (float64(s.Value()) / totalCount))
	}
	return m
}

func (s *ShareSplitStrategy) ValidateSplits(total float64, splits []*Split) (bool, error) {
	return true, nil
}

func (s *ShareSplitStrategy) Name() string {
	return "SHARE"
}
