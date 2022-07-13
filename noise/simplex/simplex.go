package simplex

func Noise1D(x float32) float32 {
	var (
		i0 = floor(x)
		i1 = i0 + 1
		x0 = x - float32(i0)
		x1 = x0 - 1
	)
	return 0.395 * (noise1(i0, x0) + noise1(i1, x1))
}

func noise1(i int32, x float32) float32 {
	t := 1 - x*x
	if t < 0 {
		return 0
	}
	t *= t
	return t * t * grad1(hash(i), x)
}

func grad1(hash uint8, x float32) float32 {
	var (
		h    = hash & 0x0F
		grad = 1 + float32(h&7)
	)
	if (h & 8) != 0 {
		grad = -grad
	}
	return grad * x
}

func Noise2D(x, y float32) float32 {
	const (
		F2 = 0.36602540
		G2 = 0.211324865
	)
	var (
		s  = (x + y) * F2
		xs = x + s
		ys = y + s
		i  = floor(xs)
		j  = floor(ys)

		t            = float32(i+j) * G2
		X0           = float32(i) - t
		Y0           = float32(j) - t
		x0           = x - X0
		y0           = y - Y0
		i1, j1 int32 = 0, 1
	)
	if x0 > y0 {
		i1, j1 = j1, i1
	}
	var (
		x1 = x0 - float32(i1)*G2
		y1 = y0 - float32(j1)*G2
		x2 = x0 - 1 + (2 * G2)
		y2 = y0 - 1 + (2 * G2)

		gi0 = hash(i + int32(hash(j)))
		gi1 = hash(i + i1 + int32(hash(j+j1)))
		gi2 = hash(i + 1 + int32(hash(j+1)))
	)
	return 45.23065 *
		noise2(gi0, x0, y0) *
		noise2(gi1, x1, y1) *
		noise2(gi2, x2, y2)
}

func noise2(gi uint8, x, y float32) float32 {
	t := 0.5 - x*x - y*y
	if t < 0 {
		return 0
	}
	t *= t
	return t * t * grad2(gi, x, y)
}

func grad2(hash uint8, x, y float32) float32 {
	h := hash & 0x3F
	if h < 4 {
		x, y = y, x
	}
	if h&1 > 0 {
		x = -x
	}
	if h&2 > 0 {
		y *= -2
	} else {
		y *= 2
	}
	return x + y

}

func Noise3D(x, y, z float32) float32 {

	// ???

	panic("TODO")
}

func grad3(hash uint8, x, y, z float32) float32 {
	var (
		h = hash & 15
		u = x
		v = z
	)
	if h < 8 {
		u = y
	}
	if h < 4 {
		v = y
	} else if h == 12 || h == 14 {
		v = x
	}
	if h&1 > 0 {
		u = -u
	}
	if h&2 > 0 {
		v = -v
	}
	return u + v
}

func floor(f float32) int32 {
	n := int32(f)
	if float32(n) < f {
		return n - 1
	}
	return n
}

var perm = [256]uint8{
	151, 160, 137, 91, 90, 15, 131, 13, 201, 95, 96, 53, 194, 233, 7, 225,
	140, 36, 103, 30, 69, 142, 8, 99, 37, 240, 21, 10, 23, 190, 6, 148,
	247, 120, 234, 75, 0, 26, 197, 62, 94, 252, 219, 203, 117, 35, 11, 32,
	57, 177, 33, 88, 237, 149, 56, 87, 174, 20, 125, 136, 171, 168, 68, 175,
	74, 165, 71, 134, 139, 48, 27, 166, 77, 146, 158, 231, 83, 111, 229, 122,
	60, 211, 133, 230, 220, 105, 92, 41, 55, 46, 245, 40, 244, 102, 143, 54,
	65, 25, 63, 161, 1, 216, 80, 73, 209, 76, 132, 187, 208, 89, 18, 169,
	200, 196, 135, 130, 116, 188, 159, 86, 164, 100, 109, 198, 173, 186, 3, 64,
	52, 217, 226, 250, 124, 123, 5, 202, 38, 147, 118, 126, 255, 82, 85, 212,
	207, 206, 59, 227, 47, 16, 58, 17, 182, 189, 28, 42, 223, 183, 170, 213,
	119, 248, 152, 2, 44, 154, 163, 70, 221, 153, 101, 155, 167, 43, 172, 9,
	129, 22, 39, 253, 19, 98, 108, 110, 79, 113, 224, 232, 178, 185, 112, 104,
	218, 246, 97, 228, 251, 34, 242, 193, 238, 210, 144, 12, 191, 179, 162, 241,
	81, 51, 145, 235, 249, 14, 239, 107, 49, 192, 214, 31, 181, 199, 106, 157,
	184, 84, 204, 176, 115, 121, 50, 45, 127, 4, 150, 254, 138, 236, 205, 93,
	222, 114, 67, 29, 24, 72, 243, 141, 128, 195, 78, 66, 215, 61, 156, 180,
}

func hash(n int32) uint8 {
	return perm[uint8(n)]
}
