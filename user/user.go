package user

import "fmt"

type User struct {
	id    string
	name  string
	email string
	phone string
}

func (u *User) Id() string {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Phone() string {
	return u.phone
}

func New(id string, name string, email string) (*User, error) {
	exist := get(id)
	if exist != nil {
		return nil, fmt.Errorf("user %s already exists", id)
	}
	user := &User{
		id:    id,
		name:  name,
		email: email,
	}
	add(user)
	return user, nil
}

func Remove(id string) bool {
	if get(id) == nil {
		return false
	}
	remove(id)
	return true
}

func Update(id, name, email string) (bool, error) {
	return update(id, name, email)
}
