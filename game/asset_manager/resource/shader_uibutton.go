package resource

var ShaderUIButtonText ShaderText = ShaderText{
	Vertex: `#version 330

	layout (location = 0) in vec3 vert;
	layout (location = 1) in vec2 vertTexCoord;
	uniform mat4 model;
	uniform float sortz;
	uniform float whr;

	out vec2 fragTexCoord;
	
	void main() {
		fragTexCoord = vertTexCoord;
		gl_Position = model * vec4(vert, 1);
		gl_Position.y *= whr;
		gl_Position.z = sortz;
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

var ShaderUIButton_Bling_Text ShaderText = ShaderText{
	Vertex: `#version 330

	layout (location = 0) in vec3 vert;
	layout (location = 1) in vec2 vertTexCoord;
	uniform mat4 model;
	uniform float sortz;
	uniform float whr;

	out vec2 fragTexCoord;
	
	void main() {
		fragTexCoord = vertTexCoord;
		gl_Position = model * vec4(vert, 1);
		gl_Position.y *= whr;
		gl_Position.z = sortz;
	}
	`,
	Fragment: `#version 330

	uniform sampler2D tex;
	uniform float light;
	uniform float blingx;
	in vec2 fragTexCoord;

	out vec4 outputColor;

	void main() {
		outputColor = texture(tex, fragTexCoord);
		// if (abs(fragTexCoord.x-blingx) < 0.3){
			outputColor.xyz *= light;
		// }
	}`,
}
