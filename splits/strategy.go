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
	shares := strategy.CalculateShare(totalAmount, splits)
	for idx, split := range splits {
		split.share = shares[idx]
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
