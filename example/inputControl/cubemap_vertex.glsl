#version 330

layout (location = 0) in vec3 vert;
layout (location = 1) in vec3 vertNormal;
uniform mat4 rotation;
uniform mat4 view;

out vec3 fragVertNormal;

void main() {
    vec4 wNormal = rotation * vec4(vertNormal, 1);
    fragTexCoord = vertTexCoord;
    fragVertNormal = wNormal.xyz;
    gl_Position = view * vec4(vert, 1);
}
