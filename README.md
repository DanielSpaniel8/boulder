# Building
install the go toolchain and run `go build`

# Usage
make a directory called `boulders` and make one or more .gmesh files in it with this format:

```
// x and y positions of the vertices (these have to go in anti-clockwise order)
Vertex[
    -2.60 -2.15
    -2.28 -1.31
    -1.45 -2.25
    -1.43 -2.89
    -0.52 -2.84
    -0.53 -2.20
    0.21 -1.34
    0.55 -2.16
    -1.03 -4.33
]
// the depth of the back of the groundmesh
MinDepth -45
// the depth of the front of the groundmesh
MaxDepth 45
// the angle either side of 0 degrees within which an edge is considered a top edge and will get a top segment
// so basically if the line between two vertices is, for example, 13.2 degrees it will have grass on it
TopAngle 20
// whether or not to make top (grass) segments at all
GenerateTop true
// the texture to use for top segments
TopTexture "graveyard_grass_2x"
// the texture to use for the front and sides
BottomTexture "graveyard_ground_2x"
```

then run the executable (`boulder` for linux, `boulder.exe` for windows) and for every input file there should be an output file in the `boulder_out` directory with the .boulder extension. this will contain all the components necessary for a groundmesh object (GroundPolygon, GroundMesh, GroundMeshGenerator, TexureMapping and CollisionShape) as well as a LocalAabb message. you can copy this in to a filerift file and put in an object definition like this:

```
Object{
    Identifier "sp_gm"
    // paste it here
    Position{
        X 213
        Y 42
    }
    Depth 0
    
```

or use the $source template

if you know how to use blender you can use the `blender_boulder.py` script. make an object with at least one face, make it face upward since the x and y coordinates will be used. select the face and stay in edit mode. go to the "scripting" workspace and paste the script into the text editor, then run it. you should see another text data-block appear with the vertex coordinates
