package resource

var ShaderCubemapText ShaderText = ShaderText{
	Vertex: `#version 330

	layout (location = 0) in vec3 vert;
	uniform mat4 rotation;
	uniform mat4 model;
	uniform mat4 view;
	uniform mat4 projection;
	
	out vec3 textureDir;
	
	void main() {
		textureDir = vert;
		vec4 pos = projection * view * model * vec4(vert, 1);
		gl_Position = pos;
	}`,
	Fragment: `#version 330

	uniform samplerCube tex;
	
	in vec3 textureDir;
	
	out vec4 outputColor;
	
	void main() {
		outputColor = texture(tex, textureDir);
	}`,
}
