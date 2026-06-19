import bpy
import bmesh


def get_polygon_xy_anticlockwise():
    obj = bpy.context.active_object
    if not obj or obj.type != "MESH":
        print("Error: No mesh object selected!")
        return

    # Get active face
    bm = None
    face = None

    if bpy.context.mode == "EDIT_MESH":
        bm = bmesh.from_edit_mesh(obj.data)
        face = bm.faces.active
        if not face:
            selected = [f for f in bm.faces if f.select]
            face = selected[0] if selected else None
    else:
        # Switch to edit mode temporarily
        bpy.ops.object.mode_set(mode="EDIT")
        bm = bmesh.from_edit_mesh(obj.data)
        face = bm.faces.active
        bpy.ops.object.mode_set(mode="OBJECT")

    if not face:
        print("Error: No face selected!")
        return

    # Create or clear output text block
    text_name = "Polygon_Vertices_XY"
    if text_name in bpy.data.texts:
        text_block = bpy.data.texts[text_name]
        text_block.clear()
    else:
        text_block = bpy.data.texts.new(text_name)

    output = []
    # output.append(f"Face Index: {face.index} | Sides: {len(face.verts)}")
    # output.append("Vertex Positions (X, Y) - Anti-clockwise order:\n")

    for i, vert in enumerate(face.verts):
        co = vert.co
        line = f"    {co.x:.6f} {co.y:.6f}"
        output.append(line)

    # Write to text block
    text_block.write("\n".join(output))

    # Show the text block in the editor
    for area in bpy.context.screen.areas:
        if area.type == "TEXT_EDITOR":
            for space in area.spaces:
                if space.type == "TEXT_EDITOR":
                    space.text = text_block
                    break

    print(f"✅ Done! Results written to Text block: '{text_name}'")


# Run it
get_polygon_xy_anticlockwise()
