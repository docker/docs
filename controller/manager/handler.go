package manager

// XXX TODO - This appears unused - we should probably remove it
/*
import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types/events"
	"github.com/docker/orca"
	"github.com/docker/orca/utils"
)

type (
	EventHandler struct {
		Manager Manager
	}
)

func (h *EventHandler) Handle(e *events.Message) error {
	log.Infof("event: date=%d status=%s container=%s", e.Time, e.Status, e.ID[:12])
	h.logDockerEvent(e)
	return nil
}

func (h *EventHandler) logDockerEvent(e *events.Message) error {
	info, err := h.Manager.Container(e.ID)
	if err != nil {
		return err
	}

	ts, err := utils.FromUnixTimestamp(e.Time)
	if err != nil {
		return err
	}

	evt := &orca.Event{
		Type: e.Status,
		Message: fmt.Sprintf("action=%s container=%s",
			e.Status, e.ID[:12]),
		Time:          *ts,
		ContainerInfo: &info,
		Tags:          []string{"docker"},
	}
	if err := h.Manager.SaveEvent(evt); err != nil {
		return err
	}
	return nil
}
*/
