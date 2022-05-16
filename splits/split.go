package splits

type Split struct {
	value        float32
	userId       string
	share        float32
	total        float64
	strategyName string
}

func (s *Split) Share() float32 {
	return s.share
}

func (s *Split) Value() float32 {
	return s.value
}

func (s *Split) UserId() string {
	return s.userId
}

func New(value float32, userId string) *Split {
	return &Split{
		value:  value,
		userId: userId,
	}
}

func NewZeroSplit(userId string) *Split {
	return &Split{
		userId: userId,
		value:  0,
	}
}
