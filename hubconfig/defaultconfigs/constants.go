package defaultconfigs

import (
	"fmt"
	"path"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/shared/containers"

	distributionconfig "github.com/docker/distribution/configuration"
	garantconfig "github.com/docker/garant/config"
)

var (
	// ReplicaIDReplacementString can be used as a placeholder for the replica
	// ID when we generate configs but want them replaced with the environment variable
	ReplicaIDReplacementString = "XXXREPLICA_IDXXX"

	DefaultHAConfig = hubconfig.HAConfig{
		LogProtocol: "internal",
		LogLevel:    "INFO",
	}
	DefaultRegistryConfig = func() (c distributionconfig.Configuration) {
		// Anonymous nested structs are bad, mmkay
		c.Version = "0.1"
		c.Storage = distributionconfig.Storage{
			"filesystem": distributionconfig.Parameters{
				"rootdirectory": "/local",
			},
			"delete": distributionconfig.Parameters{
				"enabled": true,
			},
			"maintenance": distributionconfig.Parameters{
				"readonly": map[string]interface{}{
					"enabled": false,
				},
			},
		}
		c.Notifications.Endpoints = []distributionconfig.Endpoint{}
		c.HTTP.Addr = fmt.Sprintf(":%d", deploy.StorageContainerPort)
		c.HTTP.Secret = "placeholder"
		c.Log.Level = "debug"
		// XXX: This was removed, but it should actually work now with the new networks
		// TODO: add  it back for improved logging of repository changes?
		//c.Notifications.Endpoints = []distributionconfig.Endpoint{
		//	{
		//		Name:    "dtr-audit-logs",
		//		URL:     fmt.Sprintf("http://%s/%s/", deploy.AdminServer.BridgeName(), deploy.EventsEndpointSubroute),
		//		Headers: http.Header{"X-Registry-Events": []string{"true"}},
		//	},
		//}
		return
	}()
	DefaultAuthConfig = garantconfig.Configuration{
		Version: "0.1",
		Auth: garantconfig.Auth{
			BackendName: "dtr",
		},
		Logging: garantconfig.Logging{
			Level: "debug",
		},
		HTTP: garantconfig.HTTP{
			Addr:   fmt.Sprintf(":%d", deploy.GarantPort),
			Prefix: deploy.GarantSubroute,
		},
		SigningKey: path.Join(deploy.ConfigDirPath, deploy.GarantSigningKeyFilename),
		Issuer:     "docker-trusted-registry",
	}

	DefaultMiddlewareConfig = map[string][]distributionconfig.Middleware{
		"repository": {
			{Name: "metadata"},
		},
	}

	DefaultNotaryServerConfig = hubconfig.NotaryServerConfig{
		Server: hubconfig.NotaryListeningServer{
			HTTPAddr: fmt.Sprintf(":%d", deploy.NotaryServerHTTPPort),
			NotaryTLSServerOptions: hubconfig.NotaryTLSServerOptions{
				ServerCert: containers.NotaryServerStore.CertPath(),
				ServerKey:  containers.NotaryServerStore.KeyPath(),
				ClientCA:   containers.NotaryCACertStore.CertPath(),
			},
		},

		Storage: hubconfig.NotaryStorage{
			Backend:    "rethinkdb",
			URL:        containers.Rethinkdb.BridgeName(ReplicaIDReplacementString),
			DB:         deploy.NotaryServerDBName,
			ClientCert: containers.RethinkCertStore.CertPath(),
			ClientKey:  containers.RethinkCertStore.KeyPath(),
			ServerCA:   containers.RethinkCACertStore.CertPath(),
			Username:   "server",
			Password:   "serverpass", //TODO(cyli): figure out how to randomise passwords
		},

		TrustService: hubconfig.NotaryTrustService{
			Type:         "remote",
			Hostname:     containers.NotarySigner.BridgeName(ReplicaIDReplacementString),
			Port:         fmt.Sprintf("%v", deploy.NotarySignerGRPCPort),
			KeyAlgorithm: "ecdsa",
			NotaryTLSClientOptions: hubconfig.NotaryTLSClientOptions{
				ClientCert: containers.NotaryServerStore.CertPath(),
				ClientKey:  containers.NotaryServerStore.KeyPath(),
				ServerCA:   containers.NotaryCACertStore.CertPath(),
			},
		},

		Auth: hubconfig.NotaryAuth{
			Type: "token",
			Options: hubconfig.JSONGarantOptions{
				CertBundle: path.Join(deploy.ConfigDirPath, deploy.GarantRootCertFilename),
			},
		},

		Logging: hubconfig.NotaryLogging{
			Level: "debug",
		},
	}

	DefaultNotarySignerConfig = hubconfig.NotarySignerConfig{
		Server: hubconfig.NotaryListeningServer{
			HTTPAddr: fmt.Sprintf(":%d", deploy.NotarySignerHTTPPort),
			GRPCAddr: fmt.Sprintf(":%d", deploy.NotarySignerGRPCPort),
			NotaryTLSServerOptions: hubconfig.NotaryTLSServerOptions{
				ServerCert: containers.NotarySignerStore.CertPath(),
				ServerKey:  containers.NotarySignerStore.KeyPath(),
				ClientCA:   containers.NotaryCACertStore.CertPath(),
			},
		},

		Storage: hubconfig.NotaryStorage{
			Backend:    "rethinkdb",
			URL:        containers.Rethinkdb.BridgeName(ReplicaIDReplacementString),
			DB:         deploy.NotarySignerDBName,
			ClientCert: containers.RethinkCertStore.CertPath(),
			ClientKey:  containers.RethinkCertStore.KeyPath(),
			ServerCA:   containers.RethinkCACertStore.CertPath(),
			Username:   "signer",
			Password:   "signerpass", //TODO(cyli): figure out how to randomise passwords
		},

		Logging: hubconfig.NotaryLogging{
			Level: "debug",
		},
	}
)
