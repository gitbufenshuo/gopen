#version 330

layout (location = 0) in vec3 vert;
layout (location = 1) in vec2 vertTexCoord;
layout (location = 2) in vec3 vertNormal;
uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

out vec2 fragTexCoord;
out vec3 fragVertNormal;

void main() {
    vec4 wNormal = model * vec4(vertNormal, 1);
    fragTexCoord = vertTexCoord;
    fragVertNormal = wNormal.xyz;
    gl_Position = projection * view * model * vec4(vert, 1);
}
