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
	uniform mat4 lightSpaceMatrix;
		
	out VS_OUT {
		vec3 FragPos;
		vec3 Normal;
		vec2 TexCoords;
		vec4 FragPosLightSpace;
	} vs_out;
	

	void main() {
		vec4 fragpos = model * vec4(vert, 1.0);
		vs_out.FragPos = vec3(fragpos.x,fragpos.y,fragpos.z);
		vec4 fragnormal = rotation * vec4(vertNormal, 1);
		vs_out.Normal = vec3(fragnormal.x,fragnormal.y,fragnormal.z);
		vs_out.TexCoords = vertTexCoord;
		vs_out.FragPosLightSpace = lightSpaceMatrix * vec4(vs_out.FragPos, 1.0);

		gl_Position = projection * view * model * vec4(vert, 1);
	}`,
	Fragment: `#version 330
	out vec4 outputColor;

	uniform sampler2D tex;
	uniform sampler2D u_shadowMap;

	uniform vec3 u_lightColor;
	uniform vec3 u_lightDirection;
	uniform vec3 u_viewPos;

	in VS_OUT {
		vec3 FragPos;
		vec3 Normal;
		vec2 TexCoords;
		vec4 FragPosLightSpace;
	} fs_in;
	
	float ShadowCalculation(vec4 fragPosLightSpace)
	{
		// perform perspective divide
		vec3 projCoords = fragPosLightSpace.xyz / fragPosLightSpace.w;
		// 在圆圈之内就接着检测，在圆圈之外，必然是黑的

		float chang = length(projCoords.xy);
		if (chang>0.3){
			return 1.0;
		}

		// transform to [0,1] range
		projCoords = projCoords * 0.5 + 0.5;
		if (projCoords.x<0 || projCoords.y<0){
			return 1.0;
		}
		if (projCoords.x>1 || projCoords.y>1){
			return 1.0;
		}

		// get closest depth value from light's perspective (using [0,1] range fragPosLight as coords)
		float closestDepth = texture(u_shadowMap, projCoords.xy).r; 
		// get depth of current fragment from light's perspective
		float currentDepth = projCoords.z;
		// check whether current frag pos is in shadow

		float bias = 0.005;
		float shadow = currentDepth - bias > closestDepth  ? 1.0 : 0.0;
		return shadow;
	}
		
	void main() {
		float ambientStrength = 0.1;
		float specularStrength = 6;

		vec3 ambient = ambientStrength * u_lightColor;

		vec4 sampleColor = texture(tex, fs_in.TexCoords);
		if (sampleColor.w<0.5) {
			discard;
		}
		float diffuseLight = max(dot(fs_in.Normal, normalize(u_lightDirection)),0);

		vec3 viewDir = normalize(u_viewPos - fs_in.FragPos);
		vec3 reflectDir = reflect(normalize(u_lightDirection), fs_in.Normal);  
		float spec = pow(max(dot(viewDir, reflectDir), 0.0), 32);
		vec3 specular = specularStrength * spec * u_lightColor;  
		
		// calculate shadow
		float shadow = ShadowCalculation(fs_in.FragPosLightSpace);                      
	
		vec3 lightres = (ambient + (1.0 - shadow) * (diffuseLight + specular)) * sampleColor.xyz;

		outputColor = vec4(lightres, sampleColor.w);

	}`,
}
