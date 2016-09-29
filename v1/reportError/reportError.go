package reportError

import (
	"fmt"
	"runtime"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/bugsnag/bugsnag-go"

	"github.com/docker/pinata/v1/apple/nsuserdefaults"
)

func shouldSendCrashReport() bool {
	if nsuserdefaults.KeyExists("autoSendCrashReports") == false {
		return true
	}
	result := nsuserdefaults.BoolForKey("autoSendCrashReports")
	return result
}

// Add a sysctl key value pair to the error metadata
func addSysctlUint32(MetaData *bugsnag.MetaData, key string) {
	x, err := syscall.SysctlUint32(key)
	if err != nil {
		MetaData.Add("sysctl", key, fmt.Sprintf("%#v", err))
		return
	}
	MetaData.Add("sysctl", key, fmt.Sprintf("%d", x))
}

// Initialize initializes the error reporter
func Initialize() {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:          "bd33788a5e024696711befb36d8e7b82",
		ProjectPackages: []string{"github.com/docker/pinata/[^v]*"},
		ReleaseStage:    "production",
		Synchronous:     true,
		// more configuration options
	})
	bugsnag.OnBeforeNotify(
		func(event *bugsnag.Event, config *bugsnag.Configuration) error {
			// from docker/machine/libmachine/crashreport/crash_report.go
			event.MetaData.Add("app", "compiler", fmt.Sprintf("%s (%s)", runtime.Compiler, runtime.Version()))
			event.MetaData.Add("device", "os", runtime.GOOS)
			event.MetaData.Add("device", "arch", runtime.GOARCH)
			addSysctlUint32(&event.MetaData, "machdep.cpu.family")
			addSysctlUint32(&event.MetaData, "machdep.cpu.model")
			addSysctlUint32(&event.MetaData, "machdep.cpu.extmodel")
			addSysctlUint32(&event.MetaData, "machdep.cpu.extfamily")
			addSysctlUint32(&event.MetaData, "machdep.cpu.stepping")
			addSysctlUint32(&event.MetaData, "machdep.cpu.feature_bits")
			addSysctlUint32(&event.MetaData, "machdep.cpu.leaf7_feature_bits")
			addSysctlUint32(&event.MetaData, "machdep.cpu.extfeature_bits")
			return nil
		})
}

// ReportError reports an error to bugsnag
func ReportError(message string, err error) {
	if shouldSendCrashReport() {
		err2 := bugsnag.Notify(err, bugsnag.Context{message})
		if err2 != nil {
			logrus.Println("bugsnag.Notify failed", err2)
		}
	}
}

// FatalError reports a fatal error to bugsnag
func FatalError(err error, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	ReportError(message, err)
	logrus.Fatalln(message, err)
}
