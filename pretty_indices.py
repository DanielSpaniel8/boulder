f = open("./reference_top_indices.scene", "rb")
content = f.read()

num_faces = 14
num_indices = num_faces * 3


def get_ushort() -> int:
    global content
    bits = content[:2]
    content = content[2:]
    return bits[1] * 256 + bits[0]


for i in range(num_indices):
    print(get_ushort())
    if i % 3 == 2:
        print("^ tri\n")
