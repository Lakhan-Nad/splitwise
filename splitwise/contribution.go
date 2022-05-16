package splitwise

import "fmt"

/**
 * A contribution is what a user owes to the group as whole
 * It can be simply thought as a expense user made to the group
 */
type Contribution struct {
	userId string
	amount float64
	name   string
	notes  string
	imgUrl string
}

func (e *Contribution) UserId() string {
	return e.userId
}

func (e *Contribution) Amount() float64 {
	return e.amount
}

func (e *Contribution) Name() string {
	return e.name
}

func (e *Contribution) Notes() string {
	return e.notes
}

func (e *Contribution) ImgUrl() string {
	return e.imgUrl
}

func NewContribution(userId string, amount float64, name string, notes string, imgUrl string) (*Contribution, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount cannot be zero or negative")
	}
	return &Contribution{
		userId: userId,
		amount: amount,
		name:   name,
		notes:  notes,
		imgUrl: imgUrl,
	}, nil
}
