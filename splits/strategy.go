package splits

import "fmt"

type SplitStrategy interface {
	Name() string
	CalculateShare(float64, []*Split) []float32
	ValidateSplits(float64, []*Split) (bool, error)
}

func UseSplitStrategy(strategy SplitStrategy, totalAmount float64, splits []*Split) error {
	valid, err := strategy.ValidateSplits(totalAmount, splits)
	if !valid {
		return err
	}
	strategy.CalculateShare(totalAmount, splits)
	for _, split := range splits {
		split.strategyName = strategy.Name()
		split.total = totalAmount
	}
	return nil
}

func GetStrategyFromName(name string) (SplitStrategy, error) {
	switch name {
	case "EXACT":
		return &ExactSplitStrategy{}, nil
	case "PERCENT":
		return &PercentSplitStrategy{}, nil
	case "EQUAL":
		return &EqualSplitStrategy{}, nil
	case "SHARE":
		return &ShareSplitStrategy{}, nil
	default:
		return nil, fmt.Errorf("unknown strategy %s", name)
	}
}
