
#version 330

out vec4 outputColor;
uniform float frame;

void main() {
    outputColor = vec4(frame, 0, 0, 1.0);
}
