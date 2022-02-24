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
	out vec3 fragPos;
	
	void main() {
		vec4 wNormal = rotation * vec4(vertNormal, 1);
		fragTexCoord = vertTexCoord;
		fragVertNormal = wNormal.xyz;
		gl_Position = projection * view * model * vec4(vert, 1);
		fragPos = gl_Position.xyz / gl_Position.w;
	}`,
	Fragment: `#version 330

	uniform sampler2D tex;
	uniform vec3 u_lightColor;
	uniform vec3 u_lightDirection;
	uniform vec3 u_viewPos;

	in vec2 fragTexCoord;
	in vec3 fragVertNormal;
	in vec3 fragPos;
	
	out vec4 outputColor;
	
	void main() {
		float ambientStrength = 0.1;
		float specularStrength = 0.6;

		vec3 ambient = ambientStrength * u_lightColor;

		vec4 sampleColor = texture(tex, fragTexCoord);
		if (sampleColor.w<0.5) {
			discard;
		}
		float diffuseLight = max(dot(fragVertNormal, normalize(u_lightDirection)),0);

		vec3 viewDir = normalize(u_viewPos - fragPos);
		vec3 reflectDir = reflect(normalize(u_lightDirection), fragVertNormal);  
		float spec = pow(max(dot(viewDir, reflectDir), 0.0), 32);
		vec3 specular = specularStrength * spec * u_lightColor;  
		


		sampleColor.xyz *= (ambient + diffuseLight + specular);
		outputColor = sampleColor;

	}`,
}
