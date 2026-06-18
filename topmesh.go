package main

import (
	"math"
)

func topMesh(gm GroundMesh) (vertexBits, indexBits []byte, err error) {
	var indexOffset int
	var uOffset float64
	for i := range gm.Polygon {
		if !isTopSegment(gm, i) {
			continue
		}
		vertices, indices := makeTopMeshSegment(gm, i, indexOffset, uOffset)
		indexOffset += 20
		vertexBits = append(vertexBits, vertices...)
		indexBits = append(indexBits, indices...)
		if isTopSegment(gm, i+1) {
			uOffset += distance(Vector2(gm.Polygon[i]), Vector2(gm.Polygon[(i+1)%len(gm.Polygon)]))
		} else {
			uOffset = 0
		}
	}
	return vertexBits, indexBits, nil
}

func makeTopMeshSegment(gm GroundMesh, i int, indexOffset int, uOffset float64) (vertexData, indexData []byte) {
	curr := gm.Polygon[i]
	next := gm.Polygon[(i+1)%len(gm.Polygon)]
	left := next.X
	right := curr.X
	if !isTopSegment(gm, i+1) {
		left -= 3
	}
	if !isTopSegment(gm, i-1) {
		right += 3
	}
	vertices := getVertices(left, right, next.Y, curr.Y, float64(gm.MinDepth), float64(gm.MaxDepth), uOffset)
	for _, vertex := range vertices {
		vertexData = append(vertexData, encodeFloat(vertex.X)...)
		vertexData = append(vertexData, encodeFloat(vertex.Y)...)
		vertexData = append(vertexData, encodeFloat(vertex.Z)...)
		vertexData = append(vertexData, encodeFloat(vertex.Normal.X)...)
		vertexData = append(vertexData, encodeFloat(vertex.Normal.Y)...)
		vertexData = append(vertexData, encodeFloat(vertex.Normal.Z)...)
		vertexData = append(vertexData, encodeFloat(vertex.U)...)
		vertexData = append(vertexData, encodeFloat(vertex.V)...)
	}
	for _, tri := range topIndices {
		indexData = append(indexData, encodeUShort(tri[0]+indexOffset)...)
		indexData = append(indexData, encodeUShort(tri[1]+indexOffset)...)
		indexData = append(indexData, encodeUShort(tri[2]+indexOffset)...)
	}

	return vertexData, indexData
}

func getVertices(left, right, leftHeight, rightHeight, minDepth, maxDepth float64, uOffset float64) []Vertex {
	leftHeight += 0.05 // make sure the front face doesn't clip through the top
	rightHeight += 0.05
	upNormal := Tri{
		A: Vector3{left, leftHeight, minDepth},
		B: Vector3{left, leftHeight, maxDepth},
		C: Vector3{right, rightHeight, maxDepth},
	}.surfaceNormal().normalized()
	downNormal := Tri{
		A: Vector3{right, rightHeight, maxDepth},
		B: Vector3{left, leftHeight, maxDepth},
		C: Vector3{left, leftHeight, minDepth},
	}.surfaceNormal().normalized()
	forwardNormal := Vector3{0, 0, 1}.normalized()
	leftPoint := Vector2{left, leftHeight}
	rightPoint := Vector2{right, rightHeight}
	width := distance(leftPoint, rightPoint)
	frontHeight1 := distance(Vector2{maxDepth, 0}, Vector2{maxDepth + 5, -10})
	frontHeight2 := distance(Vector2{maxDepth + 5, -10}, Vector2{maxDepth, -25})
	// there doesn't seem to be any way to procedurally make the vertices here, so here they all are ripped from a groundmesh and typed by hand
	return []Vertex{
		{right, rightHeight, minDepth, upNormal, tex(uOffset), tex(minDepth)},                                            // 0 right top back, up
		{right, rightHeight, maxDepth, upNormal, tex(uOffset), tex(maxDepth)},                                            // 1 right top front, up
		{right, rightHeight - 10, maxDepth + 5, forwardNormal, tex(uOffset - 10), tex(maxDepth + 5)},                     // 2 right middle front, forward
		{right, rightHeight - 25, maxDepth, downNormal, tex(uOffset - 25), tex(maxDepth)},                                // 3 right bottom front, down
		{right, rightHeight - 25, minDepth, downNormal, tex(uOffset - 25), tex(minDepth)},                                // 4 right bottom back, down
		{right, rightHeight - 10, minDepth - 5, forwardNormal, tex(uOffset - 10), tex(minDepth - 5)},                     // 5 right middle back, forward
		{right, rightHeight, minDepth, upNormal, tex(uOffset), tex(minDepth)},                                            // 6 = 0
		{right, rightHeight, maxDepth, upNormal, tex(uOffset), tex(maxDepth)},                                            // 7 = 1
		{right, rightHeight - 10, maxDepth + 5, forwardNormal, tex(uOffset), tex(maxDepth + frontHeight1)},               // 8 = 2, but different texture mapping
		{right, rightHeight - 25, maxDepth, downNormal, tex(uOffset), tex(maxDepth + frontHeight1 + frontHeight2)},       // 9 = 3, but different texture mapping
		{left, leftHeight, minDepth, upNormal, tex(uOffset + width), tex(minDepth)},                                      // 10 left top back, up
		{left, leftHeight, maxDepth, upNormal, tex(uOffset + width), tex(maxDepth)},                                      // 11 left top front, up
		{left, leftHeight - 10, maxDepth + 5, forwardNormal, tex(uOffset + width), tex(maxDepth + frontHeight1)},         // 12 left middle front, forward
		{left, leftHeight - 25, maxDepth, downNormal, tex(uOffset + width), tex(maxDepth + frontHeight1 + frontHeight2)}, // 13 left bottom front, down
		{left, leftHeight, minDepth, upNormal, tex(uOffset + width), tex(minDepth)},                                      // 14 = 10, but different texture mapping
		{left, leftHeight, maxDepth, upNormal, tex(uOffset + width), tex(maxDepth)},                                      // 15 = 11
		{left, leftHeight - 10, maxDepth + 5, forwardNormal, tex(uOffset + width + 10), tex(maxDepth + 5)},               // 16 = 12, but different texture mapping
		{left, leftHeight - 25, maxDepth, downNormal, tex(uOffset + width + 25), tex(maxDepth)},                          // 17 = 13, but different texture mapping
		{left, leftHeight - 25, minDepth, downNormal, tex(uOffset + width + 25), tex(minDepth)},                          // 18 left bottom back, down
		{left, leftHeight - 10, minDepth - 5, forwardNormal, tex(uOffset + width + 10), tex(minDepth - 5)},               // 19 left middle back, forward
	}
}

func isTopSegment(gm GroundMesh, i int) bool {
	l := len(gm.Polygon)
	if i < 0 {
		i += l
	}
	curr := gm.Polygon[i%l]
	next := gm.Polygon[(i+1)%l]
	angle := edgeAngle(curr, next)
	return math.Abs(angle-180) < gm.TopAngle
}

var topIndices = [][]int{
	// there also doesn't seem to be a way to make these procedurally
	{0, 4, 5},
	{0, 3, 4},
	{0, 1, 3},
	{1, 2, 3},
	{6, 10, 11},
	{7, 6, 11},
	{7, 11, 12},
	{8, 7, 12},
	{8, 12, 13},
	{9, 8, 13},
	{14, 19, 18},
	{14, 18, 17},
	{14, 17, 15},
	{17, 16, 15},
}
