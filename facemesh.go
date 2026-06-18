package main

import (
	"fmt"
	"slices"
)

func faceMesh(face GroundMesh) ([]byte, error) {
	bits := []byte{}
	if len(face.Polygon) < 3 {
		return []byte{}, fmt.Errorf("not enough points in the polygon, need at least 3")
	}
	for len(face.Polygon) > 3 {
		foundEar := false
		for i := 0; i < len(face.Polygon)-2; i++ {
			if isAnEar(i, i+1, i+2, face.Polygon) {
				a := face.Polygon[i]
				b := face.Polygon[i+1]
				c := face.Polygon[i+2]
				bits = append(bits, makeTri(a, b, c, face)...)
				face.Polygon = slices.Concat(face.Polygon[0:i+1], face.Polygon[i+2:])
				foundEar = true
				break
			}
		}
		if !foundEar {
			return []byte{}, fmt.Errorf("couldn't make a face from the polygon")
		}
	}
	a := face.Polygon[0]
	b := face.Polygon[1]
	c := face.Polygon[2]
	bits = append(bits, makeTri(a, b, c, face)...)
	return bits, nil
}

func isAnEar(a int, b int, c int, v []PolygonPoint) bool {
	if getCross(v[a], v[b], v[c]) < 0 {
		return false
	}
	// this is a bit slow, but ok for now
	for i := range v {
		if i != a && i != b && i != c {
			if isPointWithin(v[a], v[b], v[c], v[i]) {
				return false
			}
		}
	}
	return true
}

func getCross(a, b, c PolygonPoint) float64 {
	return (a.X-c.X)*(b.Y-c.Y) - (a.Y-c.Y)*(b.X-c.X)
}

func isPointWithin(a, b, c, target PolygonPoint) bool {
	if getCross(a, b, target) < 0 {
		return false
	}
	if getCross(b, c, target) < 0 {
		return false
	}
	if getCross(c, a, target) < 0 {
		return false
	}
	return true
}

func makeTri(a PolygonPoint, b PolygonPoint, c PolygonPoint, gm GroundMesh) []byte {
	bits := []byte{}
	for _, vertex := range []PolygonPoint{a, b, c} {
		// pos x y z
		bits = append(bits, encodeFloat(vertex.X)...)
		bits = append(bits, encodeFloat(vertex.Y)...)
		bits = append(bits, encodeFloat(float64(gm.MaxDepth-5))...)
		// normal x y z
		bits = append(bits, encodeFloat(-0)...)
		bits = append(bits, encodeFloat(0)...)
		bits = append(bits, encodeFloat(1)...)
		// texturemapping u v
		u := (vertex.X * textureSpaceFactor) + textureSpaceXOffset
		v := (vertex.Y * textureSpaceFactor) + textureSpaceYOffset
		bits = append(bits, encodeFloat(u)...)
		bits = append(bits, encodeFloat(v)...)
	}
	return bits
}
