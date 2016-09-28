package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/docker/dhe-deploy"
	"github.com/satori/go.uuid"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

type Event struct {
	ID          string    `gorethink:"id" json:"id"`
	PublishedAt time.Time `gorethink:"publishedAt" json:"publishedAt"`
	Actor       string    `gorethink:"actor" json:"actor"`
	Type        string    `gorethink:"type" json:"type"`
	Object      Object    `gorethink:"object" json:"object"`
}

type Object struct {
	ID   string `gorethink:"id" json:"id"`
	Type string `gorethink:"type" json:"type"`
}

type EventManager interface {
	CreateEvent(event *Event) error
	GetEvents(requestedPageEncoded string, perPage uint, publishedBefore, publishedAfter *time.Time, queryingUserId, actorId, eventType string, isAdmin bool) (events []Event, nextPageEncoded string, err error)
	Subscribe(er EventReactor) chan bool
}

type EventReactor func(e Event)

// PaginationStruct allows us to encode pagination details into a token
// for page 0, PublishedAt is right now, and we want the 0th page with any page size
// for page 1, PublishedAt must remain the same as page 0's PublishedAt. the page
// number will be 1, and the page size also must be the same as page 0's pageSize
// this allows us to calculate an offset (PageNumber*PageSize) for all events starting
// at PublishedAt
type PaginationStruct struct {
	PublishedAt time.Time `json:"publishedAt"`
	PageNumber  uint      `json:"pageNumber"`
	PageSize    uint      `json:"pageSize"`
}

type ErrDetailsPropertyMissing struct {
	Property string
}

func (e ErrDetailsPropertyMissing) Error() string {
	return fmt.Sprintf("Property '%s' is missing and required for all events", e.Property)
}

var ErrCannotQueryForOtherUser = errors.New("Users can only query for events created by themselves unless they have admin privileges")

var eventsTable = table{
	db:         deploy.DTRDBName,
	name:       "events",
	primaryKey: "id",
	secondaryIndexes: map[string][]string{
		"publishedAt":            {"publishedAt"},                  // fetch all events, ordered by publishedAt
		"type_publishedAt":       {"type", "publishedAt"},          // for getting events of a certain type, ordered by publishedAt
		"actor_publishedAt":      {"actor", "publishedAt"},         // for getting all events by a specific actor, ordered by publishedAt
		"type_actor_publishedAt": {"type", "actor", "publishedAt"}, // for getting events of specific type, by a specific actor, ordered by publishedAt
	},
}

// eventManager exports CRUDy methods for Events in the database.
type eventManager struct {
	session *rethink.Session
}

func NewEventManager(session *rethink.Session) *eventManager {
	return &eventManager{session}
}

func (m *eventManager) Subscribe(er EventReactor) chan bool {
	unsubscribe := make(chan bool)

	events, _ := eventsTable.Term().Changes().Field("new_val").Run(m.session)

	// because events.Next is blocking, we need to check for unsubscribe within the for
	// this could potentially leave a goroutine running after unsubscribing
	// but the goroutine will return when the next event comes in
	go func() {
		var event Event
		// unfortunately, it takes some time to populate events
		// calling events.Next on nil events will panic
		for events.Next(&event) {
			select {
			case <-unsubscribe:
				return
			default:
				er(event)
			}
		}
	}()

	return unsubscribe
}

func (m *eventManager) CreateEvent(event *Event) error {
	event.PublishedAt = time.Now()
	event.ID = uuid.NewV4().String()
	if _, err := eventsTable.Term().Insert(event).RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to create event in database: %s", err)
	}
	return nil
}

// GetEvents returns all events, paginated
func (m *eventManager) GetEvents(requestedPageEncoded string, perPage uint, publishedBefore, publishedAfter *time.Time, queryingUserId, actorId, eventType string, isAdmin bool) ([]Event, string, error) {
	query := eventsTable.Term()

	var publishedBeforeI, publishedAfterI interface{}

	if actorId == "" && !isAdmin {
		// if they are not admin, only query for their events
		actorId = queryingUserId
	} else if queryingUserId != actorId && !isAdmin {
		// if they give an actorId that is not theirs, and they are not admin, return error
		return nil, "", ErrCannotQueryForOtherUser
	}

	// default to getting the 0th page, starting now
	requestedPageStruct := PaginationStruct{
		PublishedAt: time.Now(),
		PageNumber:  0,
		PageSize:    perPage,
	}

	// requestedPageEncoded token supercedes published before, because:
	// suppose the previous query specifies some publishedBefore val
	// it will query, and all results will be published before publishedBefore
	// suppose we truncate the result, and the first unreturned result as id: nextid
	// event with nextid and all subsequent events will obviously have been published before publishedBefore, from the previous query
	// therefore we can use the same ceiling
	if requestedPageEncoded != "" {
		var err error
		requestedPageStruct, err = decodePaginationStruct(requestedPageEncoded)
		if err != nil {
			return nil, "", fmt.Errorf("Next pagination param is improperly formatted. It should be a next header returned by a different query.")
		}
		publishedBeforeI = requestedPageStruct.PublishedAt
	} else if publishedBefore == nil {
		publishedBeforeI = rethink.MaxVal
	} else {
		publishedBeforeI = *publishedBefore

		// also set the same ceiling on the requestedPageStruct, so we properly calculate the offset
		// for the next query
		requestedPageStruct.PublishedAt = *publishedBefore
	}
	if publishedAfter == nil {
		publishedAfterI = rethink.MinVal
	} else {
		publishedAfterI = *publishedAfter
	}

	if actorId != "" && eventType != "" {
		query = query.OrderBy(
			rethink.OrderByOpts{Index: rethink.Desc("type_actor_publishedAt")},
		).Between(
			[]interface{}{eventType, actorId, publishedAfterI},
			[]interface{}{eventType, actorId, publishedBeforeI},
		)
	} else if actorId != "" {
		query = query.OrderBy(
			rethink.OrderByOpts{Index: rethink.Desc("actor_publishedAt")},
		).Between(
			[]interface{}{actorId, publishedAfterI},
			[]interface{}{actorId, publishedBeforeI},
		)
	} else if eventType != "" {
		query = query.OrderBy(
			rethink.OrderByOpts{Index: rethink.Desc("type_publishedAt")},
		).Between(
			[]interface{}{eventType, publishedAfterI},
			[]interface{}{eventType, publishedBeforeI},
		)
	} else {
		query = query.OrderBy(
			rethink.OrderByOpts{Index: rethink.Desc("publishedAt")},
		).Between(
			[]interface{}{publishedAfterI},
			[]interface{}{publishedBeforeI},
		)
	}

	// skip results if we're not fetching the first page
	query = query.Skip(requestedPageStruct.PageNumber * requestedPageStruct.PageSize)
	return m.paginateEventQuery(query, requestedPageStruct)
}

func (m *eventManager) paginateEventQuery(query rethink.Term, requestedPageStruct PaginationStruct) (events []Event, nextPageEncoded string, err error) {
	if requestedPageStruct.PageSize > 0 {
		// fetch one more to know if there is another page
		query = query.Limit(requestedPageStruct.PageSize + 1)
	}

	cursor, err := query.Run(m.session)
	if err != nil {
		return nil, "", fmt.Errorf("EventManager: unable to query db: %s", err.Error())
	}

	if err := cursor.All(&events); err != nil {
		return nil, "", fmt.Errorf("EventManager: unable to scan query results: %s", err.Error())
	}

	if requestedPageStruct.PageSize != 0 && uint(len(events)) > requestedPageStruct.PageSize {
		// pin published at, so that the page offsets still work properly
		nextPageEncoded, err = encodePaginationStruct(requestedPageStruct.PublishedAt, requestedPageStruct.PageNumber+1, requestedPageStruct.PageSize)
		if err != nil {
			return nil, "", fmt.Errorf("EventManager: could not json marshal pagination struct: %s", err.Error())
		}

		events = events[:requestedPageStruct.PageSize]
	}

	return events, nextPageEncoded, nil
}

func encodePaginationStruct(publishedAt time.Time, pageNumber, pageSize uint) (string, error) {
	paginationStruct := PaginationStruct{
		PublishedAt: publishedAt,
		PageNumber:  pageNumber,
		PageSize:    pageSize,
	}
	bts, err := json.Marshal(paginationStruct)
	if err != nil {
		return "", err
	}
	return string(bts), nil
}

func decodePaginationStruct(encodedPaginationStruct string) (decodedPaginationStruct PaginationStruct, err error) {
	err = json.Unmarshal([]byte(encodedPaginationStruct), &decodedPaginationStruct)
	return decodedPaginationStruct, err
}
