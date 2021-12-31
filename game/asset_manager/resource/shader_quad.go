package resource

var ShaderQuadText ShaderText = ShaderText{
	Vertex: `#version 330

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
	`,
	Fragment: `#version 330

	uniform sampler2D tex;
	uniform float light;
	
	in vec2 fragTexCoord;
	
	out vec4 outputColor;
	
	void main() {
		outputColor = texture(tex, fragTexCoord);
		outputColor.xyz *= light;
	}`,
}
