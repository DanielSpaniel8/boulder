package main

import (
	"boulder/token"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const textureSpaceFactor float64 = 1.0 / 250.0
const textureSpaceXOffset float64 = 0.5
const textureSpaceYOffset float64 = 0.5

const triSize = 3 * 2    // 3 UShorts per triangle
const vertexSize = 8 * 4 // 8 floats per vertex

func main() {
	var paths []string
	filepath.WalkDir("./boulders", func(path string, d os.DirEntry, err error) error {
		if !d.IsDir() && filepath.Ext(path) == ".gmesh" {
			paths = append(paths, path)
		}
		return nil
	})
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("reading gmesh file: %s\n", err)
		}
		fmt.Printf("reading from %s\n", path)
		tokens := token.Tokenize(content)
		gm := parseGmesh(tokens)
		base := filepath.Base(path)
		name := base[:len(base)-len(filepath.Ext(path))]
		makeGM(gm, "boulders/"+name+".boulder")
		fmt.Printf("outputting to %s\n", name+".boulder")
	}
}

func makeGM(gm GroundMesh, outPath string) {
	aabb := Aabb{0, 0, 0, 0}
	for _, vertex := range gm.Polygon {
		if vertex.X < aabb.X {
			aabb.X = vertex.X
		}
		if vertex.Y < aabb.Y {
			aabb.Y = vertex.Y
		}
	}
	for _, vertex := range gm.Polygon {
		if vertex.X+(-aabb.X) > aabb.Width {
			aabb.Width = vertex.X + (-aabb.X)
		}
		if vertex.Y+(-aabb.Y) > aabb.Height {
			aabb.Height = vertex.Y + (-aabb.Y)
		}
	}
	aabb.X = aabb.X - 10
	aabb.Y = aabb.Y - 10
	aabb.Width = aabb.Width + 20
	aabb.Height = aabb.Height + 20
	gm.Aabb = aabb
	// top
	var topVertices, topIndices []byte
	var err error
	if gm.GenerateTop {
		topVertices, topIndices, err = topMesh(gm)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("topmesh: %d verts, %d tris\n", len(topVertices)/(4*8), len(topIndices)/6)
	}
	// sides
	sideVertices, sideIndices, err := sideMesh(gm)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("sidemesh: %d verts, %d faces\n", len(sideVertices)/(4*8), len(sideIndices)/(2*3))
	// front
	face, err := faceMesh(gm)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("facemesh: %d verts, %d faces\n", len(face)/(4*8), len(face)/(4*8*3))
	outputFR(topVertices, topIndices, sideVertices, sideIndices, face, gm, outPath)
}

func parseGmesh(tokens []token.Token) GroundMesh {
	gm := GroundMesh{
		Polygon:       []PolygonPoint{},
		MinDepth:      -45,
		MaxDepth:      45,
		Aabb:          Aabb{},
		TopAngle:      20,
		GenerateTop:   true,
		TopTexture:    "grass_subtle",
		BottomTexture: "maybegood",
	}
	var inRecord bool
	var inCollection bool
	var inVertexPair bool
	var firstVertexValue float64
	var name string
	for _, tok := range tokens {
		switch tok.Type {
		case token.TokenName:
			if !inRecord {
				inRecord = true
			}
			name = tok.Value
		case token.TokenNumber:
			n, err := strconv.ParseFloat(tok.Value, 64)
			if err == nil {
				switch name {
				case "Vertex":
					if inVertexPair {
						gm.Polygon = append(gm.Polygon, PolygonPoint{firstVertexValue, n})
						inVertexPair = false
					} else {
						firstVertexValue = n
						inVertexPair = true
					}
				case "MinDepth":
					gm.MinDepth = n
				case "MaxDepth":
					gm.MaxDepth = n
				case "TopAngle":
					gm.TopAngle = n
				}
			}
			inRecord = inCollection
		case token.TokenString:
			if name == "TopTexture" {
				gm.TopTexture = tok.Value
			}
			if name == "BottomTexture" {
				gm.BottomTexture = tok.Value
			}
			inRecord = inCollection
		case token.TokenBoolean:
			if name == "GenerateTop" {
				gm.GenerateTop = tok.Value == "true"
			}
			inRecord = inCollection
		case token.TokenCollectionStart:
			inCollection = true
		case token.TokenMessageEnd:
			inCollection = false
			inRecord = false
		}
	}
	return gm
}
