package main

import (
	"fmt"
)

const textureSpaceFactor float64 = 1.0 / 250.0
const textureSpaceXOffset float64 = 0.5
const textureSpaceYOffset float64 = 0.5

const triSize = 3 * 2    // 3 UShorts per triangle
const vertexSize = 8 * 4 // 8 floats per vertex

var gm = GroundMesh{
	Polygon: []PolygonPoint{
		// >>> original obj9#14
		// {-109.481705, 38.326042},
		// {-34.278633, -38.326042},
		// {4.485405, -3.3574295},
		// {53.21099, -9.89341},
		// {78.1648, 9.575951},
		// {109.4817, 35.590324},

		// >>> triangle
		// {-50, 30},
		// {0, -30},
		// {50, 30},

		// >>> curve top
		// {-80, 30},
		// {-50, -20},
		// {0, -30},
		// {50, -20},
		// {80, 30},
		// {40, 33},
		// {0, 35},
		// {-40, 33},

		// >>> two top meshes
		// {-50, 100},
		// {-50, -50},
		// {100, -50},
		// {100, 40},
		// {50, 50},
		// {0, 50},
		// {0, 100},

		// >>> B top segment
		// {-2.749579, 380.182159},
		// {4.946239, 266.730316},
		// {80.493423, 311.578369},
		// {128.331375, 305.506958},
		// {151.524948, 271.435486},
		// {191.530624, 212.355042},
		// {217.867828, 281.741730},
		// {198.082489, 351.050354},
		// {144.633545, 375.730835},
		// >>> B middle segment
		// {4.946239, 266.730316},
		// {5.626010, 123.984055},
		// {81.521408, 164.346451},
		// {131.110580, 163.591156},
		// {156.122574, 140.186340},
		// {240.519104, 138.642395},
		// {226.777542, 184.076096},
		// {191.530609, 212.355042},
		// {151.524948, 271.435486},
		// {132.555817, 213.289001},
		// {78.933609, 210.272003},
		// {80.493423, 311.578369},
		// >>> B bottom segment
		// {5.626010, 123.984055},
		// {-1.557067, 9.726846},
		// {124.868614, 12.579068},
		// {199.843750, 36.732552},
		// {231.795013, 73.789246},
		// {240.519104, 138.642395},
		// {156.122574, 140.186340},
		// {153.665985, 91.722458},
		// {126.651886, 72.750023},
		// {81.568466, 69.202629},
		// {81.521408, 164.346451},

		// >>> o top
		// {119.504181, 148.362946},
		// {57.756130, 186.673950},
		// {-21.714581, 170.808746},
		// {-58.755249, 102.656525},
		// {-6.165573, 97.430336},
		// {13.189793, 129.832687},
		// {59.940540, 133.794006},
		// {78.298187, 89.992195},
		// {134.010864, 75.454102},
		// >>> o bottom
		{134.010864, 75.454102},
		{78.298187, 89.992195},
		{54.981716, 55.064114},
		{8.336624, 55.802517},
		{-6.165573, 97.430336},
		{-58.755249, 102.656525},
		{-40.860992, 29.603432},
		{22.276726, -5.513626},
		{89.160629, 10.512428},

		// >>> u
		// {0, 60},
		// {1.5, 18.5},
		// {9.1, 6.1},
		// {25.8, 1},
		// {36.5, 5.1},
		// {41.9, 11.9},
		// {44.2, 0.3},
		// {76.2, 2.5},
		// {73.5, 14.3},
		// {68.4, 58.6},
		// {41.6, 58},
		// {44.4, 30.5},
		// {42.7, 22},
		// {36.3, 17},
		// {29.2, 24.9},
		// {26.1, 58.1},
	},
	MinDepth:      -100,
	MaxDepth:      100,
	TopAngle:      60,
	GenerateTop:   false,
	TopTexture:    "grass_subtle",
	BottomTexture: "maybegood",
}

func main() {
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
	outputFR(topVertices, topIndices, sideVertices, sideIndices, face, gm)
}
