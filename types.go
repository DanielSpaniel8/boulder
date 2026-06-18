package main

import (
	"fmt"
	"math"
)

type Vector2 struct{ X, Y float64 }
type Vector3 struct{ X, Y, Z float64 }

type PolygonPoint struct{ X, Y float64 }

type Vertex struct {
	X      float64
	Y      float64
	Z      float64
	Normal Vector3
	U      float64
	V      float64
}

type Tri struct{ A, B, C Vector3 }

type TriData struct {
	Tri            Tri
	Normal         Vector3
	TextureMapping []Vector2
}

type Quad struct{ A, B, C, D Vector3 }

type Aabb struct{ X, Y, Width, Height float64 }

type GroundMesh struct {
	Polygon       []PolygonPoint
	MinDepth      float64
	MaxDepth      float64
	Aabb          Aabb
	TopAngle      float64
	GenerateTop   bool
	TopTexture    string
	BottomTexture string
}

func (v1 Vector3) delta(v2 Vector3) Vector3 {
	return Vector3{v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z}
}

func distance(v1, v2 Vector2) float64 {
	dx := v2.X - v1.X
	dy := v2.Y - v1.Y
	s := math.Pow(dx, 2) + math.Pow(dy, 2)
	return math.Sqrt(s)
}

func normalize(v Vector2) Vector2 {
	len := math.Sqrt(v.X*v.X + v.Y*v.Y)
	if len == 0 {
		return Vector2{}
	}
	return Vector2{v.X / len, v.Y / len}
}

func (tri Tri) surfaceNormal() Vector3 {
	u := tri.B.delta(tri.A)
	v := tri.C.delta(tri.A)
	return Vector3{
		X: (u.Y * v.Z) - (u.Z * v.Y),
		Y: (u.Z * v.X) - (u.X * v.Z),
		Z: (u.X * v.Y) - (u.Y * v.X),
	}
}

func (v Vector3) split() (float64, float64, float64) {
	return v.X, v.Y, v.Z
}

func (v Vector3) normalized() Vector3 {
	magnitude := math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2) + math.Pow(v.Z, 2))
	v.X = v.X / magnitude
	v.Y = v.Y / magnitude
	v.Z = v.Z / magnitude
	return v
}

func (td TriData) encode() []byte {
	bits := []byte{}
	vertices := []Vector3{td.Tri.A, td.Tri.B, td.Tri.C}
	for i := range 3 {
		vertex := vertices[i]
		mapping := td.TextureMapping[i]
		bits = append(bits, encodeFloat(vertex.X)...)
		// fmt.Printf("vertex.X: %f\n", vertex.X)
		bits = append(bits, encodeFloat(vertex.Y)...)
		// fmt.Printf("vertex.Y: %f\n", vertex.Y)
		bits = append(bits, encodeFloat(vertex.Z)...)
		// fmt.Printf("vertex.Z: %f\n", vertex.Z)
		bits = append(bits, encodeFloat(td.Normal.X)...)
		// fmt.Printf("td.Normal.X: %f\n", td.Normal.X)
		bits = append(bits, encodeFloat(td.Normal.Y)...)
		// fmt.Printf("td.Normal.Y: %f\n", td.Normal.Y)
		bits = append(bits, encodeFloat(td.Normal.Z)...)
		// fmt.Printf("td.Normal.Z: %f\n", td.Normal.Z)
		bits = append(bits, encodeFloat(mapping.X)...)
		fmt.Printf("mapping.X: %f\n", mapping.X)
		bits = append(bits, encodeFloat(mapping.Y)...)
		fmt.Printf("mapping.Y: %f\n", mapping.Y)
	}
	return bits
}
