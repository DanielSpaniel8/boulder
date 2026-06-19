# Building
install the go toolchain and run `go build`

# Usage
make a directory called `boulders` and make one or more files in it with this format:

```
// x and y positions of the vertices (these have to go in clockwise order)
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
// 
```
