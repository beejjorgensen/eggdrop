package aabb

// AABB is an axis-aligned bounding box
type AABB struct {
	X0, Y0, X1, Y1 int32
}

// New creates a new AABB
func New(x0, y0, x1, y1 int32) *AABB {
	return &AABB{x0, y0, x1, y1}
}

// TestCollision checks for an AABB collision
func (a *AABB) TestCollision(b *AABB) bool {
	return !(a.X0 > b.X1 || a.X1 < b.X0 || a.Y0 > b.Y1 || a.Y1 < b.Y0)
}
