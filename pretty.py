import struct

f = open("reference_top.scene", "rb")
content = f.read()


def get_float() -> float:
    global content
    bits = content[:4]
    content = content[4:]
    return struct.unpack("<f", bits)[0]


num_vertices = 20

for i in range(num_vertices):
    print("X " + str(get_float()))
    print("Y " + str(get_float()))
    print("Z " + str(get_float()))
    print("X " + str(get_float()))
    print("Y " + str(get_float()))
    print("Z " + str(get_float()))
    print("U " + str(get_float()))
    print("V " + str(get_float()) + "\n")
