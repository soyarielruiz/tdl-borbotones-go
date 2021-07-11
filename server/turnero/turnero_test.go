package turnero_test

import (
	"testing"

	"github.com/soyarielruiz/tdl-borbotones-go/server/turnero"
	"github.com/soyarielruiz/tdl-borbotones-go/server/user"
)

func TestTurneroCrearConMapaVacio(t *testing.T) {
	m := make(map[string]*user.User)
	if turnero.New(m) != nil {
		t.Fail()
	}
}

func TestTurneroNext(t *testing.T) {
	m := make(map[string]*user.User)
	m["lucho"] = &user.User{PlayerId: "lucho"}
	m["santi"] = &user.User{PlayerId: "santi"}
	tu := turnero.New(m)
	if tu.CurrentUser() != "lucho" {
		t.Errorf("Expected lucho - Actual %s", tu.CurrentUser())
	}
}

func TestTurneroNextAll(t *testing.T) {
	m := make(map[string]*user.User)
	m["lucho"] = &user.User{PlayerId: "lucho"}
	m["santi"] = &user.User{PlayerId: "santi"}
	tu := turnero.New(m)
	if tu.CurrentUser() != "lucho" {
		t.Errorf("Expected lucho - Actual %s", tu.CurrentUser())
	}
	tu.Next()
	if tu.CurrentUser() != "santi" {
		t.Errorf("Expected santi - Actual %s", tu.CurrentUser())
	}
	tu.Next()
	if tu.CurrentUser() != "lucho" {
		t.Errorf("Expected lucho - Actual %s", tu.CurrentUser())
	}
}

func TestTurneroGoToAndNext(t *testing.T) {
	m := make(map[string]*user.User)
	m["ari"] = &user.User{PlayerId: "ari"}
	m["ele"] = &user.User{PlayerId: "ele"}
	m["lucho"] = &user.User{PlayerId: "lucho"}
	m["santi"] = &user.User{PlayerId: "santi"}
	tu := turnero.New(m)
	tu.GoTo("lucho")
	if tu.CurrentUser() != "lucho" {
		t.Errorf("Expected lucho - Actual %s", tu.CurrentUser())
	}
}

func TestTurneroChangeDirection(t *testing.T) {
	m := make(map[string]*user.User)
	m["ari"] = &user.User{PlayerId: "ari"}
	m["ele"] = &user.User{PlayerId: "ele"}
	m["lucho"] = &user.User{PlayerId: "lucho"}
	m["santi"] = &user.User{PlayerId: "santi"}
	tu := turnero.New(m)
	tu.ChangeDirection()
	if tu.CurrentUser() != "ari" {
		t.Errorf("Expected ari - Actual %s", tu.CurrentUser())
	}
	tu.Next()
	if tu.CurrentUser() != "santi" {
		t.Errorf("Expected santi - Actual %s", tu.CurrentUser())
	}
}

func TestTurneroTurnFirstRemoveLastRightDirection(t *testing.T) {
	m := make(map[string]*user.User)
	m["ari"] = &user.User{PlayerId: "ari"}
	m["ele"] = &user.User{PlayerId: "ele"}
	m["lucho"] = &user.User{PlayerId: "lucho"}
	tu := turnero.New(m)
	tu.Remove("lucho")
	if tu.CurrentUser() != "ari" {
		t.Errorf("Expected ari - Actual %s", tu.CurrentUser())
	}
	tu.Next()
	if tu.CurrentUser() != "ele" {
		t.Errorf("Expected ele - Actual %s", tu.CurrentUser())
	}
	tu.Next()
	if tu.CurrentUser() != "ari" {
		t.Errorf("Expected ari - Actual %s", tu.CurrentUser())
	}
}

func TestTurneroTurnFirstRemoveFirstRightDirection(t *testing.T) {
	m := make(map[string]*user.User)
	m["ari"] = &user.User{PlayerId: "ari"}
	m["ele"] = &user.User{PlayerId: "ele"}
	m["lucho"] = &user.User{PlayerId: "lucho"}
	tu := turnero.New(m)
	tu.Remove("ari")
	if tu.CurrentUser() != "ele" {
		t.Errorf("Expected ele - Actual %s", tu.CurrentUser())
	}
}

func TestTurneroTurnMiddleRemoveMiddleRightDirection(t *testing.T) {
	m := make(map[string]*user.User)
	m["ele"] = &user.User{PlayerId: "ele"}
	m["lucho"] = &user.User{PlayerId: "lucho"}
	m["santi"] = &user.User{PlayerId: "santi"}
	tu := turnero.New(m)
	tu.Next()
	tu.Remove("lucho")
	if tu.CurrentUser() != "santi" {
		t.Errorf("Expected santi - Actual %s", tu.CurrentUser())
	}
}

func TestTurneroTurnLastRemoveFirstRightDirection(t *testing.T) {
	m := make(map[string]*user.User)
	m["ele"] = &user.User{PlayerId: "ele"}
	m["lucho"] = &user.User{PlayerId: "lucho"}
	m["santi"] = &user.User{PlayerId: "santi"}
	tu := turnero.New(m)
	tu.GoTo("santi")
	tu.Remove("ele")
	if tu.CurrentUser() != "santi" {
		t.Errorf("Expected santi - Actual %s", tu.CurrentUser())
	}
}

func TestTurneroTurnFirstRemoveLastLeftDirection(t *testing.T) {
	m := make(map[string]*user.User)
	m["ari"] = &user.User{PlayerId: "ari"}
	m["ele"] = &user.User{PlayerId: "ele"}
	m["lucho"] = &user.User{PlayerId: "lucho"}
	tu := turnero.New(m)
	tu.ChangeDirection()
	tu.Remove("lucho")
	if tu.CurrentUser() != "ari" {
		t.Errorf("Expected ari - Actual %s", tu.CurrentUser())
	}
	tu.Next()
	if tu.CurrentUser() != "ele" {
		t.Errorf("Expected ele - Actual %s", tu.CurrentUser())
	}
}

func TestTurneroTurnFirstRemoveFirstLeftDirection(t *testing.T) {
	m := make(map[string]*user.User)
	m["ari"] = &user.User{PlayerId: "ari"}
	m["ele"] = &user.User{PlayerId: "ele"}
	m["lucho"] = &user.User{PlayerId: "lucho"}
	tu := turnero.New(m)
	tu.ChangeDirection()
	tu.Remove("ari")
	if tu.CurrentUser() != "lucho" {
		t.Errorf("Expected lucho - Actual %s", tu.CurrentUser())
	}
}

func TestTurneroTurnMiddleRemoveMiddleLeftDirection(t *testing.T) {
	m := make(map[string]*user.User)
	m["ele"] = &user.User{PlayerId: "ele"}
	m["lucho"] = &user.User{PlayerId: "lucho"}
	m["santi"] = &user.User{PlayerId: "santi"}
	tu := turnero.New(m)
	tu.Next()
	tu.ChangeDirection()
	tu.Remove("lucho")
	if tu.CurrentUser() != "ele" {
		t.Errorf("Expected ele - Actual %s", tu.CurrentUser())
	}
}

func TestTurneroTurnLastRemoveFirstLeftDirection(t *testing.T) {
	m := make(map[string]*user.User)
	m["ele"] = &user.User{PlayerId: "ele"}
	m["lucho"] = &user.User{PlayerId: "lucho"}
	m["santi"] = &user.User{PlayerId: "santi"}
	tu := turnero.New(m)
	tu.ChangeDirection()
	tu.GoTo("santi")
	tu.Remove("ele")
	if tu.CurrentUser() != "santi" {
		t.Errorf("Expected santi - Actual %s", tu.CurrentUser())
	}
}
