
#version 330

in vec3 vert;
uniform float frame;

void main() {
    gl_Position = vec4(vert.x + frame, vert.y * frame, vert.z, 1.0);
}
