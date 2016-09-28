package adminserver

import (
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/distribution/configuration"
)

var spewConfig = &spew.ConfigState{DisableMethods: true}

type RegistrySettings struct {
	Config  *string               `json:"config,omitempty"`
	Storage configuration.Storage `json:"storage,omitempty"`
}

func (s RegistrySettings) String() string {
	return spewConfig.Sprintf("%#v", s)
}

func (s RegistrySettings) GoString() string {
	return spewConfig.Sprintf("%#v", s)
}

type LicenseSettings struct {
	IsValid     bool      `json:"is_valid"`
	AutoRefresh bool      `json:"auto_refresh"`
	Expiration  time.Time `json:"expiration"`
	KeyID       string    `json:"key_id"`
	LicenseTier string    `json:"tier"`
	LicenseType string    `json:"type"`
}

func (s LicenseSettings) String() string {
	return spewConfig.Sprintf("%#v", s)
}

func (s LicenseSettings) GoString() string {
	return spewConfig.Sprintf("%#v", s)
}

type storageConfig struct {
	Field       string `json:"field"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsRequired  bool   `json:"isRequired"`
}

var storageConfigurations = map[string][]storageConfig{
	"filesystem": {
		{
			"rootdirectory",
			"Storage directory",
			"The root directory tree in which all registry files will be stored",
			true,
		},
	},
	"s3": {
		{
			"region",
			"AWS region",
			"The name of the AWS region in which you would like to store objects (for example us-east-1)",
			true,
		},
		{
			"regionendpoint",
			"AWS region endpoint URL",
			"An alternate S3 endpoint URL for AWS compatible storage.  Note: you can leave this empty if using S3 directly",
			false,
		},
		{
			"bucket",
			"Bucket name",
			"The name of your S3 bucket where you wish to store objects (needs to already be created prior to driver initialization)",
			true,
		},
		{
			"accesskey",
			"Access Key",
			"Your AWS access key",
			false,
		},
		{
			"secretkey",
			"Secret Key",
			"Your AWS secret key. Note: you can provide empty strings for your access and secret keys if you plan on running the driver on an EC2 instance using IAM role grant credentials",
			false,
		},
		{
			"rootdirectory",
			"Root Directory",
			"The root directory tree in which all registry files will be stored. Defaults to the empty string (bucket root)",
			false,
		},
		{
			"v4auth",
			"Signature Version 4 Auth",
			"Indicates whether the registry uses Version 4 of AWS's authentication. Generally, you should set this to true",
			false,
		},
		{
			"secure",
			"Use HTTPS",
			"Indicates whether to use HTTPS instead of HTTP",
			false,
		},
	},
	"azure": {
		{
			"accountname",
			"Account name",
			"Name of the Azure storage account",
			true,
		},
		{
			"accountkey",
			"Account key",
			"Primary or secondary key for the storage account",
			true,
		},
		{
			"container",
			"Container",
			"Name of the root storage container in which all registry data will be stored. Must comply the storage container name requirements",
			true,
		},
		{
			"realm",
			"Realm",
			`Domain name suffix for the storage service API endpoint. Defaults to core.windows.net. For example, the realm for "Azure in China" would be core.chinacloudapi.cn and the realm for "Azure Government" would be core.usgovcloudapi.net`,
			false,
		},
	},
	// transcribed from https://github.com/docker/distribution/blob/master/docs/storage-drivers/swift.md
	"swift": {
		{
			"authurl",
			"Authorization URL",
			"URL for obtaining an auth token.",
			true,
		},
		{
			"username",
			"Username",
			"Your OpenStack username",
			true,
		},
		{
			"password",
			"Password",
			"Your OpenStack password",
			true,
		},
		{
			"container",
			"Container",
			"The name of your Swift container where you wish to store objects",
			true,
		},
		{
			"tenant",
			"Tenant Name",
			"(Optional) Your OpenStack tenant name. You can either use Tenant or TenantID",
			false,
		},
		{
			"tenantid",
			"Tenant ID",
			"(Optional) Your OpenStack tenant id. You can either use Tenant or TenantID",
			false,
		},
		{
			"domain",
			"Domain Name",
			"(Optional) Your OpenStack domain name for identity v3 API. You can either use Domain or DomainID",
			false,
		},
		{
			"domainid",
			"Domain ID",
			"(Optional) Your OpenStack domain id for identity v3 API. You can either use Domain or DomainID",
			false,
		},
		{
			"insecureskipverify",
			"Skip TLS Verification",
			"(Optional) Set this to true to skip TLS verification for your OpenStack provider. The driver uses false by default.",
			false,
			// TODO support checkboxes or something
		},
		{
			"region",
			"Region",
			"(Optional) Specify the OpenStack region name in which you would like to store objects (for example \"fr\").",
			false,
		},
		{
			"chunksize",
			"Chunk size",
			"(Optional) Specify the segment size for Dynamic Large Objects uploads (performed by WriteStream) to Swift. The default is 5 MB. You might experience better performance for larger chunk sizes depending on the speed of your connection to Swift.",
			false,
		},
		{
			"prefix",
			"Prefix",
			"(Optional) Supply the root directory tree in which to store all registry files. Defaults to the empty string which is the container's root.",
			false,
		},
	},
}
