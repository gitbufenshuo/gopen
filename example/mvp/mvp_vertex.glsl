#version 330

layout (location = 0) in vec3 vert;
layout (location = 1) in vec2 vertTexCoord;
uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * view * model * vec4(vert, 1);
}
