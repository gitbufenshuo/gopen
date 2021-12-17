#version 330

uniform samplerCube tex;

in vec3 fragVertNormal;

out vec4 outputColor;

void main() {
    outputColor = texture(tex, fragVertNormal);
}