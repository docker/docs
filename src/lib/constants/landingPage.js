/* eslint-disable max-len */
//------------------------------------------------------------------------------
// DOCKER DATA CENTER
//------------------------------------------------------------------------------
export const DDC = {
  title: 'On-premises Containers as a Service with Docker Datacenter',
  DTR: {
    name: 'Trusted Registry',
    description: 'Docker Trusted Registry allows you to store and manage your Docker images on-premises or in your virtual private cloud to support security or regulatory compliance requirements.',
  },
  UCP: {
    name: 'Universal Control Plane',
    description: 'Docker Universal Control Plane (UCP) is an on-premises tool that allows IT ops teams to manage, deploy and scale their applications across their dockerized environment.',
  },
  CS: {
    name: 'CS Engine',
    description: 'Docker commercially supported engine (CS engine) provides patches, hot fixes, and support for previous engine versions (3 versions) for all Docker engines installed in an enterprise’s environment.',
  },
  whatsIncluded: [
    'Universal Control Plane (with embedded Swarm)',
    'Docker Trusted Registry',
    'Commercial Support for Docker Engines',
  ],
};

//------------------------------------------------------------------------------
// Categories
//------------------------------------------------------------------------------
export const ANALYTICS_CATEGORY = 'Analytics';
export const APPLICATION_FRAMEWORK_CATEGORY = 'Application Frameworks';
export const APPLICATION_INFRASTRUCTURE_CATEGORY = 'Application Infrastructure';
export const APPLICATION_SERVICES_CATEGORY = 'Application Services';
export const BASE_CATEGORY = 'Base Images';
export const DATABASE_CATEGORY = 'Databases';
export const FEATURED_CATEGORY = 'Featured Images';
export const LANGUAGES_CATEGORY = 'Programming Languages';
export const MESSAGING_CATEGORY = 'Messaging Services';
export const OS_CATEGORY = 'Operating Systems';
export const STORAGE_CATEGORY = 'Storage';
export const TOOLS_CATEGORY = 'DevOps Tools';

export const categoryDescriptions = {
  analytics: {
    name: ANALYTICS_CATEGORY,
    description: 'Understand your application and users.',
  },
  application_framework: {
    name: APPLICATION_FRAMEWORK_CATEGORY,
    description: 'Libraries and structure for your application development.',
  },
  application_infrastructure: {
    name: APPLICATION_INFRASTRUCTURE_CATEGORY,
    description: 'Runtime services to scale your application.',
  },
  application_services: {
    name: APPLICATION_SERVICES_CATEGORY,
    description: 'Fully featured solutions to assist your application, or stand alone.',
  },
  base: {
    name: BASE_CATEGORY,
    description: 'Build on these secure and solid foundations.',
  },
  database: {
    name: DATABASE_CATEGORY,
    description: 'Persist your application and user data.',
  },
  featured: {
    name: FEATURED_CATEGORY,
    description: 'Exemplary solutions that stand out.',
  },
  languages: {
    name: LANGUAGES_CATEGORY,
    description: 'Assemble your development environment.',
  },
  messaging: {
    name: MESSAGING_CATEGORY,
    description: 'Secure and efficient communication for your components.',
  },
  os: {
    name: OS_CATEGORY,
    description: 'Secure and customizable low level functionality.',
  },
  storage: {
    name: STORAGE_CATEGORY,
    description: 'Save your data, whatever its form.',
  },
  tools: {
    name: TOOLS_CATEGORY,
    description: 'Build, deploy, and monitor your application.',
  },
};

//------------------------------------------------------------------------------
// Articles
//------------------------------------------------------------------------------
export const helpArticles = {
  small: {
    name: 'How to Create Great Images for the Docker Store',
    description: 'Learn how to create great Docker Images for your software and delight your customers.',
  },
  large: {
    name: 'Why You Should List Your Content on Docker Store',
    description: 'Distribute software to Docker\'s fast-growing customer base. Customers discover, install, and purchase your software directly from the Docker store.',
  },
  blogPost: {
    name: 'Introducing Docker Store Private Beta',
    description: 'The Store offers better discovery, security, trust and reputation indicators - all the things you like about the Official Repos, and more. The Store is about enabling new publishers, large and small, to publish their containers on the Docker Store, while leveraging a publish pipeline to ensure the quality.',
  },
};
