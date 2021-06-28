package turnero

import (
	"sort"

	"github.com/soyarielruiz/tdl-borbotones-go/server/user"
)

type Turnero struct {
	actualIndex     int
	direction       int
	usersIndex      map[string]int
	usersCollection []user.User // Ver si usamos punteros aca
}

func New(users map[string]user.User) *Turnero {

	if len(users) < 2 {
		return nil
	}

	t := new(Turnero)
	t.actualIndex = 0
	t.direction = 1
	t.usersIndex = make(map[string]int)
	t.usersCollection = make([]user.User, 0)

	keys := make([]string, 0, len(users))
	for k := range users {
		keys = append(keys, k)
	}
	// Orden alfabetico para los test
	sort.Strings(keys)
	for i, k := range keys {
		t.usersCollection = append(t.usersCollection, users[k])
		t.usersIndex[k] = i
	}

	return t
}

func (t *Turnero) CurrentUser() string {
	return t.usersCollection[t.actualIndex].PlayerId
}

func (t *Turnero) Next() {
	t.actualIndex += t.direction
	if t.direction == 1 && t.actualIndex == len(t.usersCollection) {
		t.actualIndex = 0
	} else if t.direction == -1 && t.actualIndex == -1 {
		t.actualIndex = len(t.usersCollection) - 1
	}
}

func (t *Turnero) GoTo(s string) {
	t.actualIndex = t.usersIndex[s]
}

func (t *Turnero) ChangeDirection() {
	t.direction *= -1
}

func (t *Turnero) Remove(s string) {
	cu := ""
	i := t.usersIndex[s]

	if t.actualIndex == i {
		t.Next()
		cu = t.CurrentUser()
	} else if i < t.actualIndex {
		t.actualIndex--
	}

	t.usersCollection = append(t.usersCollection[:i], t.usersCollection[i+1:]...)
	delete(t.usersIndex, s)
	for i, v := range t.usersCollection {
		t.usersIndex[v.PlayerId] = i
	}

	if cu != "" {
		t.actualIndex = t.usersIndex[cu]
	}
}

func (t *Turnero) IsUserTurn(playerId string) bool {
	return playerId == t.CurrentUser()
}
