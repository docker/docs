'use strict';

export default {
	server: {
		type: 'Docker Trusted Registry',
		name: 'Server Starter Edition includes',
		includes: [
			'1 instance of Docker Trusted Registry',
			'10 Docker Engines with commercial support for the servers hosting your application',
			'Email support service levels for your Docker software'
		],
		redirect_value: 'eval',
		notes: 'License keys and commercially supported Docker Engine software are distributed and managed within your Docker Hub account.'
	},
	trial: {
		type: 'a trial of Docker Datacenter',
		name: 'Docker Datacenter Trial includes access to',
		includes: [
			'Trusted Registry',
			'Universal Control Plane',
			'Commercial Support for Docker Engine'
		],
		redirect_value: 'eval',
		notes: 'License keys and commercially supported Docker Engine software are distributed and managed within your Docker Hub account.'
	},
	cloud: {
		type: 'Docker Cloud Subscription',
		name: 'Cloud Starter Edition includes',
		includes: [
			'20 Private Repos',
			'10 Docker Engines with commercial support for the servers hosting your application',
			'Email support service levels for your Docker software'
		],
		redirect_value: 'cloud_starter',
		notes: 'License keys and commercially supported Docker Engine software are distributed and managed within your Docker Hub account.'
	},
	micro: {
		type: 'our Private Repo Plans',
		name: 'Micro plan includes',
		includes: [
			'5 Private Repos',
			'5 Parallel Builds',
			'Community Hub Support'
		],
		redirect_value: 'micro',
		notes: null
	},
	small: {
		type: 'our Private Repo Plans',
		name: 'Small plan includes',
		includes: [
			'10 Private Repos',
			'10 Parallel Builds',
			'Community Hub Support'
		],
		redirect_value: 'small',
		notes: null
	},
	medium: {
		type: 'our Private Repo Plans',
		name: 'Medium plan includes',
		includes: [
			'20 Private Repos',
			'20 Parallel Builds',
			'Community Hub Support'
		],
		redirect_value: 'medium',
		notes: null
	},
	large: {
		type: 'our Private Repo Plans',
		name: 'Large plan includes',
		includes: [
			'50 Private Repos',
			'50 Parallel Builds',
			'Community Hub Support'
		],
		redirect_value: 'large',
		notes: null
	},
	xlarge: {
		type: 'our Private Repo Plans',
		name: 'XLarge plan includes',
		includes: [
			'100 Private Repos',
			'100 Parallel Builds',
			'Community Hub Support'
		],
		redirect_value: 'xlarge',
		notes: null
	},
	xxlarge: {
		type: 'our Private Repo Plans',
		name: 'XX-Large plan includes',
		includes: [
			'250 Private Repos',
			'250 Parallel Builds',
			'Community Hub Support'
		],
		redirect_value: 'xxlarge',
		notes: null
	}
};
