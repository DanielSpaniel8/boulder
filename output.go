package main

import (
	"fmt"
	"os"
)

func outputFR(topVertices, topIndices, sideVertices, sideIndices, frontVertices []byte, gm GroundMesh, outPath string) {
	// build the polygon for the GroundPolygon and CollisionShape
	var left, right, bottom, top float64
	polygonString := ""
	for _, vertex := range gm.Polygon {
		polygonString += fmt.Sprintf("                    Vertex{X %f Y %f}\n", vertex.X, vertex.Y)
		if vertex.X < left {
			left = vertex.X
		}
		if vertex.X > right {
			right = vertex.X
		}
		if vertex.Y < bottom {
			bottom = vertex.Y
		}
		if vertex.Y > top {
			top = vertex.Y
		}
	}
	// the bounding box and square are used for LocalAabbs and BoundingBoxs
	boundingBoxString := fmt.Sprintf("X %f Y %f Z -50 Width %f Height %f Depth 100", left, bottom, (-left)+right, (-bottom)+top)
	boundingSquareString := fmt.Sprintf("X %f Y %f Width %f Height %f", left, bottom, (-left)+right, (-bottom)+top)
	groundPolygonComponentString := fmt.Sprintf(groundPolygonComponent, polygonString, gm.MinDepth, gm.MaxDepth)
	topMeshString := ""
	topMeshString += fmt.Sprintf(surfaceMesh, len(topVertices)/vertexSize, len(topIndices)/triSize, gm.TopTexture, boundingBoxString, Quote(topVertices, '"'), Quote(topIndices, '"'))
	groundMeshComponentString := fmt.Sprintf(groundMeshComponent, boundingSquareString, topMeshString, len(sideVertices)/vertexSize, len(sideIndices)/triSize, gm.BottomTexture, boundingBoxString, Quote(sideVertices, '"'), Quote(sideIndices, '"'), len(frontVertices)/vertexSize, len(frontVertices)/(3*vertexSize), gm.BottomTexture, boundingBoxString, Quote(frontVertices, '"'))
	collisionShapeComponentString := fmt.Sprintf(collisionShapeComponent, polygonString, gm.MinDepth, gm.MaxDepth)
	topTextureMapping := fmt.Sprintf(textureMappingComponent, 984, gm.TopTexture, 250.0)
	bottomTextureMapping := fmt.Sprintf(textureMappingComponent, 985, gm.BottomTexture, 250.0)
	outputString := groundPolygonComponentString + groundMeshComponentString + groundMeshGeneratorComponent + collisionShapeComponentString + topTextureMapping + bottomTextureMapping + "        LocalAabb{" + boundingSquareString + "}\n"
	os.WriteFile(outPath, []byte(outputString), 0o777)
}

// 98 == 0x62 == b for boulder
// GroundPolygon = 980
// GroundMesh = 981
// GroundMeshGenerator = 982
// CollisionShape = 983
// TextureMapping 1 = 984
// TextureMapping 2 = 985

const groundPolygonComponent = `
        Component{
            ClassName : 'GroundPolygon'
            Identifier : 980
            GroundPolygonComponent{
                Polygon{
%s                    Convex : 0
                    Closed : 1
                }
                Collides : 1
                MinDepth : %f
                MaxDepth : %f
            }
        }
`

const groundMeshComponent = `
        Component{
            ClassName : 'GroundMesh'
            Identifier : 981
            GroundMeshComponent{
                LocalAabb{ %s }
%s                // side mesh
                SurfaceMesh{
                    NumVertices : %d
                    NumFaces : %d
                    Indices{ ValueType : 4 ValuesPerVertex : 1 Stride : 2 DataOffset : 0 }
                    Vertices{ ValueType : 7 ValuesPerVertex : 3 Stride : 32 DataOffset : 0 }
                    Normals{ ValueType : 7 ValuesPerVertex : 3 Stride : 32 DataOffset : 12 }
                    TexCoordSet{ ValueType : 7 ValuesPerVertex : 2 Stride : 32 DataOffset : 24 }
                    Material{
                        AmbientColor{ R : 1.0 G : 1.0 B : 1.0 A : 1.0 }
                        DiffuseColor{ R : 1.0 G : 1.0 B : 1.0 A : 1.0 }
                        SpecularColor{ R : 1.0 G : 1.0 B : 1.0 A : 1.0 }
                        Shininess : 0.0
                        Texture{ Name : '%s' PixelFormat : 1 ImageType : 2 }
                    }
                    BoundingBox{ %s }
                    VertexData : "%s"
                    IndexData : "%s"
                }
                FrontMesh{
                    NumVertices : %d
                    NumFaces : %d
                    Vertices{ ValueType : 7 ValuesPerVertex : 3 Stride : 32 DataOffset : 0 }
                    Normals{ ValueType : 7 ValuesPerVertex : 3 Stride : 32 DataOffset : 12 }
                    TexCoordSet{ ValueType : 7 ValuesPerVertex : 2 Stride : 32 DataOffset : 24 }
                    Material{
                        AmbientColor{ R : 1.0 G : 1.0 B : 1.0 A : 1.0 }
                        DiffuseColor{ R : 1.0 G : 1.0 B : 1.0 A : 1.0 }
                        SpecularColor{ R : 1.0 G : 1.0 B : 1.0 A : 1.0 }
                        Shininess : 0.0
                        Texture{ Name : '%s' PixelFormat : 1 ImageType : 2 }
                    }
                    BoundingBox{ %s }
                    VertexData : "%s"
                }
                Color{ R : 1.0 G : 1.0 B : 1.0 A : 1.0 }
            }
        }
`

const surfaceMesh = `
                // top mesh
                SurfaceMesh{
                    NumVertices : %d
                    NumFaces : %d
                    Indices{ ValueType : 4 ValuesPerVertex : 1 Stride : 2 DataOffset : 0 }
                    Vertices{ ValueType : 7 ValuesPerVertex : 3 Stride : 32 DataOffset : 0 }
                    Normals{ ValueType : 7 ValuesPerVertex : 3 Stride : 32 DataOffset : 12 }
                    TexCoordSet{ ValueType : 7 ValuesPerVertex : 2 Stride : 32 DataOffset : 24 }
                    Material{
                        AmbientColor{ R : 1.0 G : 1.0 B : 1.0 A : 1.0 }
                        DiffuseColor{ R : 1.0 G : 1.0 B : 1.0 A : 1.0 }
                        SpecularColor{ R : 1.0 G : 1.0 B : 1.0 A : 1.0 }
                        Shininess : 0.0
                        Texture{ Name : '%s' PixelFormat : 1 ImageType : 2 }
                    }
                    BoundingBox{ %s }
                    VertexData : "%s"
                    IndexData : "%s"
                }
`

const groundMeshGeneratorComponent = `
        Component{
            ClassName : 'GroundMeshGenerator'
            Identifier : 982
            GroundMeshGeneratorComponent{
                GroundPolygonId : 980
                TargetMeshId : 981
                FrontTextureMappingId : 107
                SurfaceTextureMappingId : 106
                RandomSeed : 1291618994
                HorizNoise : 0.0
                MeshType : 1
                SurfaceWidth : 80.0
                HatHeight : 25.0
                HatWidthOffset1 : 5.0
                HatWidthOffset2 : 5.0
            }
        }
`

const textureMappingComponent = `
        Component{
            ClassName : 'TextureMapping'
            Identifier : %d
            TextureMappingComponent{
                TextureName : '%s'
                Scale : %f
                Offset{ X : 0.0 Y : 0.0 }
            }
        }
`

const collisionShapeComponent = `
        Component{
            ClassName : 'CollisionShape'
            Identifier : 983
            ParentComponentIdentifier : 980
            ShapeComponent{
                Polygon{
%s
                    Convex : 0
                    Closed : 1
                }
            }
            CollisionShapeComponent{
                IsGround : 1
                MinDepth : %f
                MaxDepth : %f
                Enabled : 1
            }
        }
`
