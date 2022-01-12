package resource

var ShaderMVPText ShaderText = ShaderText{
	Vertex: `#version 330

	layout (location = 0) in vec3 vert;
	layout (location = 1) in vec2 vertTexCoord;
	layout (location = 2) in vec3 vertNormal;
	uniform mat4 rotation;
	uniform mat4 model;
	uniform mat4 view;
	uniform mat4 projection;
	
	out vec2 fragTexCoord;
	out vec3 fragVertNormal;
	
	void main() {
		vec4 wNormal = rotation * vec4(vertNormal, 1);
		fragTexCoord = vertTexCoord;
		fragVertNormal = wNormal.xyz;
		gl_Position = projection * view * model * vec4(vert, 1);
	}`,
	Fragment: `#version 330

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
		outputColor.xyz *= 1.0;
	}`,
}
