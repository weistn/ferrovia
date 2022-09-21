package tracks

import "math"

// 2d Vektor.
type Vec2 [2]float64

// 3d Vektor.
type Vec3 [3]float64

func (v Vec2) Invert() Vec2 {
	return [2]float64{-v[0], -v[1]}
}

// Rotate rotates the vector v by r degrees clock-wise
func (v Vec2) Rotate(r float64) Vec2 {
	cos := math.Cos(r / -180 * math.Pi)
	sin := math.Sin(r / -180 * math.Pi)
	return [2]float64{v[0]*cos - (-v[1])*sin, -(v[0]*sin + (-v[1])*cos)}
}

func (v Vec3) Invert() Vec3 {
	return [3]float64{-v[0], -v[1], -v[2]}
}

func (v Vec3) Add2(v2 Vec2) Vec3 {
	return [3]float64{v[0] + v2[0], v[1] + v2[1], v[2]}
}

// Returns an angle in the range [0,360[.
func normalizeAngle(a float64) float64 {
	if a < 0 {
		for a < 0 {
			a += 360
		}
	} else {
		for a >= 360 {
			a -= 360
		}
	}
	return a
}
