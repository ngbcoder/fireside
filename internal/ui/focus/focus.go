package focus

import "math"

// Direction is a navigation direction.
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

// Rect is an axis-aligned rectangle in screen pixels.
type Rect struct {
	X, Y, W, H int
}

func (r Rect) Center() (float64, float64) {
	return float64(r.X) + float64(r.W)/2, float64(r.Y) + float64(r.H)/2
}

// Node is a focusable element.
type Node struct {
	ID   string
	Rect Rect
}

// Manager holds the set of focusable nodes and the current focus.
type Manager struct {
	nodes   []Node
	current int // index into nodes, -1 if none
}

// NewManager creates an empty focus manager.
func NewManager() *Manager {
	return &Manager{current: -1}
}

// SetNodes replaces the focusable set. Keeps current focus if still present.
func (m *Manager) SetNodes(nodes []Node) {
	prevID := ""
	if m.current >= 0 && m.current < len(m.nodes) {
		prevID = m.nodes[m.current].ID
	}
	m.nodes = append([]Node(nil), nodes...)
	m.current = -1
	if prevID != "" {
		for i := range m.nodes {
			if m.nodes[i].ID == prevID {
				m.current = i
				return
			}
		}
	}
	if len(m.nodes) > 0 {
		m.current = 0
	}
}

// Nodes returns a copy of the current node set.
func (m *Manager) Nodes() []Node {
	out := make([]Node, len(m.nodes))
	copy(out, m.nodes)
	return out
}

// Current returns the currently focused node, if any.
func (m *Manager) Current() (Node, bool) {
	if m.current < 0 || m.current >= len(m.nodes) {
		return Node{}, false
	}
	return m.nodes[m.current], true
}

// CurrentIndex returns the index of the focused node.
func (m *Manager) CurrentIndex() int { return m.current }

// Focus sets focus to the node with the given id.
func (m *Manager) Focus(id string) bool {
	for i := range m.nodes {
		if m.nodes[i].ID == id {
			m.current = i
			return true
		}
	}
	return false
}

// Move moves focus in the given direction using nearest-neighbor selection.
// Returns true if focus changed.
func (m *Manager) Move(dir Direction) bool {
	if m.current < 0 || len(m.nodes) == 0 {
		if len(m.nodes) > 0 {
			m.current = 0
			return true
		}
		return false
	}
	cur := m.nodes[m.current]
	cx, cy := cur.Rect.Center()
	var bestIdx int = -1
	bestScore := math.MaxFloat64
	for i := range m.nodes {
		if i == m.current {
			continue
		}
		n := m.nodes[i]
		nx, ny := n.Rect.Center()
		var primary, secondary float64
		var ok bool
		switch dir {
		case Right:
			if nx > cx {
				primary = nx - cx
				secondary = math.Abs(ny - cy)
				ok = true
			}
		case Left:
			if nx < cx {
				primary = cx - nx
				secondary = math.Abs(ny - cy)
				ok = true
			}
		case Down:
			if ny > cy {
				primary = ny - cy
				secondary = math.Abs(nx - cx)
				ok = true
			}
		case Up:
			if ny < cy {
				primary = cy - ny
				secondary = math.Abs(nx - cx)
				ok = true
			}
		}
		if !ok {
			continue
		}
		// Weighted score: primary distance dominates, secondary breaks ties.
		score := primary + 0.5*secondary
		if score < bestScore {
			bestScore = score
			bestIdx = i
		}
	}
	if bestIdx < 0 {
		return false
	}
	m.current = bestIdx
	return true
}
