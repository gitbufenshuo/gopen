#version 330

layout (location = 0) in vec3 vert;
layout (location = 1) in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = vec4(vert, 1);
}
