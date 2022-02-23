package resource

type ShaderText struct {
	Vertex   string
	Fragment string
}

var ShaderSkyboxText ShaderText = ShaderText{
	Vertex: `#version 330

	layout (location = 0) in vec3 vert;
	uniform mat4 rotation;
	
	out vec3 textureDir;
	
	void main() {
		vec4 wNormal = rotation * vec4(vert, 1);
		textureDir = vert;
		gl_Position = wNormal.xyww;
	}`,
	Fragment: `#version 330

	uniform samplerCube tex;
	
	in vec3 textureDir;
	
	out vec4 outputColor;
	
	void main() {
		// outputColor = vec4(1,0,0,1);
		outputColor = texture(tex, textureDir);
	}`,
}
