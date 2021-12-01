#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;
in vec3 fragVertNormal;

out vec4 outputColor;

void main() {
    vec3 directLightSource = vec3(0,0,1);
    outputColor = texture(tex, fragTexCoord);
    float light = dot(fragVertNormal, directLightSource) + 0.3;
    outputColor.xyz *= light;
}