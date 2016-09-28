package sanitizers

import (
	"testing"

	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/defaultconfigs"
	"github.com/docker/dhe-deploy/hubconfig/memory"
)

func TestSanitizedLogConfig(t *testing.T) {
	testCases := []struct {
		unsanitized, sanitized hubconfig.HAConfig
		expectSetError         bool
	}{
		{
			// invalid log protocol
			unsanitized: hubconfig.HAConfig{
				LogProtocol: "blah",
			},
			expectSetError: true,
		},
		{
			// default log level
			unsanitized: hubconfig.HAConfig{},
			sanitized: hubconfig.HAConfig{
				LogProtocol: "internal",
				LogLevel:    "INFO",
			},
		},
		{
			// no host provided
			unsanitized: hubconfig.HAConfig{
				LogProtocol: "tcp",
			},
			expectSetError: true,
		},
		{
			// host provided
			unsanitized: hubconfig.HAConfig{
				LogProtocol: "tcp",
				LogHost:     "host",
			},
			sanitized: hubconfig.HAConfig{
				LogProtocol: "tcp",
				LogHost:     "host",
				LogLevel:    "INFO",
			},
		},
		{
			// valid log level
			unsanitized: hubconfig.HAConfig{
				LogLevel: "DEBUG",
			},
			sanitized: hubconfig.HAConfig{
				LogProtocol: "internal",
				LogLevel:    "DEBUG",
			},
		},
		{
			// invalid log level
			unsanitized: hubconfig.HAConfig{
				LogLevel: "derp",
			},
			expectSetError: true,
		},
		{
			// valid tls protocol
			unsanitized: hubconfig.HAConfig{
				LogProtocol:      "tcp+tls",
				LogHost:          "host",
				LogTLSCACert:     "ca-cert",
				LogTLSCert:       "cert",
				LogTLSKey:        "key",
				LogTLSSkipVerify: true,
			},
			sanitized: hubconfig.HAConfig{
				LogProtocol:      "tcp+tls",
				LogHost:          "host",
				LogTLSCACert:     "ca-cert",
				LogTLSCert:       "cert",
				LogTLSKey:        "key",
				LogTLSSkipVerify: true,
				LogLevel:         "INFO",
			},
		},
		{
			// valid tls protocol without tls certs - for unchanged certs
			unsanitized: hubconfig.HAConfig{
				LogProtocol: "tcp+tls",
				LogHost:     "host",
			},
			sanitized: hubconfig.HAConfig{
				LogProtocol: "tcp+tls",
				LogHost:     "host",
				LogLevel:    "INFO",
			},
		},
		{
			// valid tcp protocol with tls certs - see docker run
			unsanitized: hubconfig.HAConfig{
				LogProtocol:      "tcp",
				LogHost:          "host",
				LogTLSCACert:     "ca-cert",
				LogTLSCert:       "cert",
				LogTLSKey:        "key",
				LogTLSSkipVerify: false,
			},
			sanitized: hubconfig.HAConfig{
				LogProtocol:      "tcp",
				LogHost:          "host",
				LogTLSCACert:     "ca-cert",
				LogTLSCert:       "cert",
				LogTLSKey:        "key",
				LogTLSSkipVerify: false,
				LogLevel:         "INFO",
			},
		},
	}
	for _, testCase := range testCases {
		ss := LogSanitizingSettingsStore{memory.NewSettingsStore()}

		testCaseCombined := defaultconfigs.DefaultHAConfig
		if testCase.unsanitized.LogLevel != "" {
			testCaseCombined.LogLevel = testCase.unsanitized.LogLevel
		}
		if testCase.unsanitized.LogProtocol != "" {
			testCaseCombined.LogProtocol = testCase.unsanitized.LogProtocol
		}
		if testCase.unsanitized.LogHost != "" {
			testCaseCombined.LogHost = testCase.unsanitized.LogHost
		}
		err := ss.SetHAConfig(&testCaseCombined)
		if testCase.expectSetError {
			if err == nil {
				t.Fatalf("Expected setting of config to fail: %#v", testCase.unsanitized)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for test: %v error: %s", testCase.unsanitized, err)
			}
			retrieved, err := ss.HAConfig()
			if err != nil {
				t.Error(err)
			}
			if retrieved == nil {
				t.Error("Retrieved nil config")
			}

			assertHAConfigEquals(t, *retrieved, testCase.sanitized)

		}
	}
}

func assertHAConfigEquals(t *testing.T, received, expected hubconfig.HAConfig) {
	// we don't want to test all fields
	if received.LogLevel != expected.LogLevel || received.LogProtocol != expected.LogProtocol {
		t.Fatalf("Configs not equal!\nReceived: %#v\nExpected: %#v", received, expected)
	}
}
