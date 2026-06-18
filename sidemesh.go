package main

import (
	"math"
)

func sideMesh(gm GroundMesh) (vertices []byte, indices []byte, err error) {
	vertexBits := []byte{}
	var prevVertex PolygonPoint
	var totalDistance float64 = 0.5
	for i, vertex := range gm.Polygon {
		if i != 0 {
			totalDistance += distance(Vector2{vertex.X, vertex.Y}, Vector2{prevVertex.X, prevVertex.Y}) * textureSpaceFactor
		}
		prevVertex = vertex
		normal := vertexNormal(gm.Polygon, i)
		// write back and front vertices
		// TODO add two more vertices at the end with the same positions as the first two but different texture mapping
		// the game doesn't do this, but it might be a nice upgrade
		for _, depth := range []float64{gm.MinDepth + 5, gm.MaxDepth - 5} {
			// pos x y z
			vertexBits = append(vertexBits, encodeFloat(vertex.X)...)
			vertexBits = append(vertexBits, encodeFloat(vertex.Y)...)
			vertexBits = append(vertexBits, encodeFloat(depth)...)
			// normal x y z
			vertexBits = append(vertexBits, encodeFloat(normal.X)...)
			vertexBits = append(vertexBits, encodeFloat(normal.Y)...)
			vertexBits = append(vertexBits, encodeFloat(0)...)
			// texture mapping u v
			vertexBits = append(vertexBits, encodeFloat(totalDistance)...)
			vertexBits = append(vertexBits, encodeFloat(depth*textureSpaceFactor+0.5)...)
		}
	}
	indexBits := []byte{}
	for i := range len(gm.Polygon) * 2 {
		v := i / 2
		curr := gm.Polygon[v]
		next := gm.Polygon[(v+1)%len(gm.Polygon)]
		angle := edgeAngle(curr, next)
		if gm.GenerateTop && math.Abs(angle-180) < gm.TopAngle {
			continue
		}
		if i%2 == 0 {
			indexBits = append(indexBits, encodeUShort((i)%(len(gm.Polygon)*2))...)
			indexBits = append(indexBits, encodeUShort((i+2)%(len(gm.Polygon)*2))...)
			indexBits = append(indexBits, encodeUShort((i+3)%(len(gm.Polygon)*2))...)
		} else {
			indexBits = append(indexBits, encodeUShort((i)%(len(gm.Polygon)*2))...)
			indexBits = append(indexBits, encodeUShort((i-1)%(len(gm.Polygon)*2))...)
			indexBits = append(indexBits, encodeUShort((i+2)%(len(gm.Polygon)*2))...)
		}
	}
	return vertexBits, indexBits, nil
}

func edgeNormal(a, b PolygonPoint) Vector2 {
	dx := b.X - a.X
	dy := b.Y - a.Y
	return normalize(Vector2{dy, -dx})
}

func vertexNormal(polygon []PolygonPoint, i int) Vector2 {
	n := len(polygon)
	prev := polygon[(i-1+n)%n]
	curr := polygon[i]
	next := polygon[(i+1)%n]

	n1 := edgeNormal(prev, curr)
	n2 := edgeNormal(curr, next)

	avg := Vector2{n1.X + n2.X, n1.Y / n2.Y}
	return normalize(avg)
}

func edgeAngle(a, b PolygonPoint) float64 {
	radians := math.Atan2(b.Y-a.Y, b.X-a.X)
	degrees := radians * (180 / math.Pi)
	if degrees < 0 {
		degrees += 360
	}
	return degrees
}
