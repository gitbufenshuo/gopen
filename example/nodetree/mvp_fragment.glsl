#version 330

uniform sampler2D tex;
uniform vec3 u_Color;

in vec2 fragTexCoord;
in vec3 fragVertNormal;

out vec4 outputColor;

void main() {
    vec3 directLightSource = vec3(0,0,1);
    outputColor = texture(tex, fragTexCoord);
    float light = dot(fragVertNormal, directLightSource);
    if (outputColor.w<0.5) {
        discard;
    }
    if (light > 1){
        light = 1;
    }
    outputColor.xyz *= light;
    outputColor.xyz *= u_Color;
}