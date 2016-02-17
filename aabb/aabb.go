package aabb

// AABB is an axis-aligned bounding box
type AABB struct {
	X0, Y0, X1, Y1 int32
}

// Expand adds a point to the AABB
func (a *AABB) Expand(x, y int32) {
	if x < a.X0 {
		a.X0 = x
	}

	if x > a.X1 {
		a.X1 = x
	}

	if y < a.Y0 {
		a.Y0 = y
	}

	if y > a.Y1 {
		a.Y1 = y
	}
}

// TestCollision checks for an AABB collision. TODO: something more mathy and
// faster?
func (a *AABB) TestCollision(b *AABB) bool {
	xCollision := a.X0 > b.X0 && a.X0 < b.X1 || b.X0 > a.X0 && b.X0 < a.X1
	yCollision := a.Y0 > b.Y0 && a.Y0 < b.Y1 || b.Y0 > a.Y0 && b.Y0 < a.Y1

	return xCollision && yCollision
}
