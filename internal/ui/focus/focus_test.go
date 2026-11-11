package focus

import "testing"

func TestMoveRight(t *testing.T) {
	m := NewManager()
	m.SetNodes([]Node{
		{ID: "a", Rect: Rect{X: 0, Y: 0, W: 10, H: 10}},
		{ID: "b", Rect: Rect{X: 100, Y: 0, W: 10, H: 10}},
		{ID: "c", Rect: Rect{X: 200, Y: 0, W: 10, H: 10}},
	})
	m.Focus("a")
	if !m.Move(Right) {
		t.Fatal("expected move right")
	}
	cur, _ := m.Current()
	if cur.ID != "b" {
		t.Fatalf("expected b, got %s", cur.ID)
	}
	m.Move(Right)
	cur, _ = m.Current()
	if cur.ID != "c" {
		t.Fatalf("expected c, got %s", cur.ID)
	}
	// No further right.
	if m.Move(Right) {
		t.Fatal("expected no move at edge")
	}
}

func TestMoveVertical(t *testing.T) {
	m := NewManager()
	m.SetNodes([]Node{
		{ID: "a", Rect: Rect{X: 0, Y: 0, W: 10, H: 10}},
		{ID: "b", Rect: Rect{X: 0, Y: 100, W: 10, H: 10}},
	})
	m.Focus("a")
	if !m.Move(Down) {
		t.Fatal("expected move down")
	}
	cur, _ := m.Current()
	if cur.ID != "b" {
		t.Fatalf("expected b, got %s", cur.ID)
	}
}

func TestSetNodesKeepsFocus(t *testing.T) {
	m := NewManager()
	m.SetNodes([]Node{{ID: "a", Rect: Rect{0, 0, 10, 10}}})
	m.Focus("a")
	m.SetNodes([]Node{{ID: "b", Rect: Rect{0, 0, 10, 10}}, {ID: "a", Rect: Rect{10, 0, 10, 10}}})
	cur, _ := m.Current()
	if cur.ID != "a" {
		t.Fatalf("expected focus retained on a, got %s", cur.ID)
	}
}
