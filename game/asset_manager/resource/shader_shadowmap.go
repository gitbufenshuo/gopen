package resource

var ShaderShadowMapText ShaderText = ShaderText{
	Vertex: `#version 330 core
	layout (location = 0) in vec3 aPos;
	
	uniform mat4 lightSpaceMatrix;
	uniform mat4 model;
	
	void main()
	{
		gl_Position = lightSpaceMatrix * model * vec4(aPos, 1.0);
	}`,
	Fragment: `#version 330 core

	void main()
	{             
		// gl_FragDepth = gl_FragCoord.z;
	}`,
}
