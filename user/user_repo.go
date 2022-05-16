package user

import "fmt"

var (
	userMap = make(map[string]*User)
)

func add(user *User) (bool, error) {
	_, ok := userMap[user.Id()]
	if ok {
		return false, fmt.Errorf("user %s already exists", user.Id())
	}
	userMap[user.Id()] = user
	return true, nil
}

func get(id string) *User {
	return userMap[id]
}

func remove(id string) {
	delete(userMap, id)
}

func update(id, name, email string) (bool, error) {
	user, ok := userMap[id]
	if !ok {
		return false, fmt.Errorf("user %s not found", id)
	}
	user.name = name
	user.email = email
	return true, nil
}
