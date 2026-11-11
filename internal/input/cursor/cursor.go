package cursor

// Cursor is a software pointer for cursor mode.
type Cursor struct {
	X, Y   int
	Width  int
	Height int
	Step   int // pixels per move
	Active bool
}

// New creates a cursor centered on a screen of w x h.
func New(w, h int) *Cursor {
	return &Cursor{X: w / 2, Y: h / 2, Width: w, Height: h, Step: 24, Active: false}
}

// Enable turns cursor mode on.
func (c *Cursor) Enable() { c.Active = true }

// Disable turns cursor mode off.
func (c *Cursor) Disable() { c.Active = false }

// Toggle toggles cursor mode.
func (c *Cursor) Toggle() { c.Active = !c.Active }

// Move moves the cursor by dx,dy pixels, clamped to screen.
func (c *Cursor) Move(dx, dy int) {
	if !c.Active {
		return
	}
	c.X += dx
	c.Y += dy
	if c.X < 0 {
		c.X = 0
	}
	if c.Y < 0 {
		c.Y = 0
	}
	if c.X > c.Width {
		c.X = c.Width
	}
	if c.Y > c.Height {
		c.Y = c.Height
	}
}

// MoveStep moves by one step in a direction.
func (c *Cursor) MoveStep(dir string) {
	switch dir {
	case "left":
		c.Move(-c.Step, 0)
	case "right":
		c.Move(c.Step, 0)
	case "up":
		c.Move(0, -c.Step)
	case "down":
		c.Move(0, c.Step)
	}
}
