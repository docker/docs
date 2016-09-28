package manager

import (
	"reflect"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca"
)

var (
	ksEvents = datastoreVersion + "/events"
)

// Replicated from encoding/json since this isn't public
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// TODO Consider renaming this...
func (m DefaultManager) SaveEvent(event *orca.Event) error {
	// Instead of hard-coding this, try to make it dynamic so we can
	// refine the event type in the future
	fields := log.Fields{}
	fields["license_key"] = m.GetLicenseKeyID()
	ev := reflect.ValueOf(*event)
	et := ev.Type()
	msg := "event"
	for i := 0; i < ev.NumField(); i++ {
		fv := ev.Field(i)
		ft := et.Field(i)
		tags := strings.Split(ft.Tag.Get("json"), ",") // Mimic the json serialization
		if len(tags) > 0 {
			// Figure out if we omit it
			omit := false
			for _, tag := range tags {
				if tag == "omitempty" && isEmptyValue(fv) {
					omit = true
					break
				}
			}
			if omit {
				continue
			}
			// Special case the message
			if tags[0] == "message" {
				msg = fv.String()
				continue
			}
			// XXX Are we ever logging a time != now?
			if tags[0] == "time" {
				continue
			}
			// Else include
			fields[tags[0]] = fv.Interface()
		} else {
			fields[ft.Name] = fv.Interface()
		}
	}
	log.WithFields(fields).Info(msg)
	return nil
}
