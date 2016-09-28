package adminserver

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/robfig/cron"
)

func getTemplate() *template.Template {
	tpl := path.Join(TemplatesDir, "index.html")
	return template.Must(template.ParseFiles(tpl))
}

func makePageClass(pageName string) map[string]string {
	return map[string]string{pageName: "active"}
}

func writeJSON(rw http.ResponseWriter, payload interface{}) {
	writeJSONStatus(rw, payload, http.StatusOK)
}

func writeJSONError(rw http.ResponseWriter, err error, status int) {
	writeJSONStatus(rw, struct {
		Message string `json:"message"`
	}{
		Message: err.Error(),
	}, status)
}

func writeJSONStatus(rw http.ResponseWriter, payload interface{}, status int) {
	if payload == nil {
		payload = struct{}{}
	}
	marshalled, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.WithField("error", err).Error("Failed to encode response")
		writeJSONError(rw, fmt.Errorf("failed to encode response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(status)
	fmt.Fprint(rw, string(marshalled)+"\n")
}

func stringPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func boolPointer(b bool) *bool {
	return &b
}

// Converts the DTR cron format to the robfig/cron format with >4 fields. A DTR cron format has at most
// four fields representing hour, day, month and day of week respectively. Moreover, the shortcuts
// `@yearly`, `@annually`, `@monthly`,`@weekly`, `@daily`, `@midnight`, `@hourly` are allowed.
func transformCron(schedule string) (string, error) {
	// don't care about leading or trailing whitespace
	schedule = strings.TrimSpace(schedule)

	// disallow the @every prefix
	if strings.HasPrefix(schedule, "@every") {
		return "", fmt.Errorf("Unrecognized descriptor: %s", schedule)
	}

	// verify that an @ prefix-ed schedule is valid.
	if len(schedule) >= 1 && schedule[0] == '@' {
		_, err := cron.Parse(schedule)
		return schedule, err
	}

	// enforce four fields
	fields := strings.Fields(schedule)
	if len(fields) > 4 {
		return "", fmt.Errorf("Cron schedule has more than four fields: %v", len(fields))
	}

	transformedSchedule := fmt.Sprintf("0 0 %s", schedule)
	_, err := cron.Parse(transformedSchedule)
	if err != nil {
		return "", err
	} else {
		return transformedSchedule, nil
	}
}

// Transforms the underlying schedule string to one we want to be user visible
func untransformCron(schedule string) string {
	if schedule == "" || (len(schedule) > 1 && schedule[0] == '@') {
		return schedule
	}

	return strings.TrimPrefix(schedule, "0 0 ")
}
