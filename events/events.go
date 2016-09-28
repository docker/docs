package events

import (
	"github.com/docker/dhe-deploy/events/types"
	"github.com/docker/dhe-deploy/manager/schema"
)

func NewRepositoryEvent(eventManager schema.EventManager, userId, repositoryId string) error {
	return repositoryEvent(eventManager, types.Create, userId, repositoryId)
}

func DeleteRepositoryEvent(eventManager schema.EventManager, userId, repositoryId string) error {
	return repositoryEvent(eventManager, types.Delete, userId, repositoryId)
}

func repositoryEvent(eventManager schema.EventManager, eventType, userId, repositoryId string) error {
	event := schema.Event{
		Actor: userId,
		Type:  eventType,
		Object: schema.Object{
			Type: types.Repository,
			ID:   repositoryId,
		},
	}
	return eventManager.CreateEvent(&event)
}

func TagImageEvent(eventManager schema.EventManager, userId, repositoryId, tagId string) error {
	return tagEvent(eventManager, types.Update, userId, repositoryId, tagId)
}

func DeleteTagEvent(eventManager schema.EventManager, userId, repositoryId, tagId string) error {
	return tagEvent(eventManager, types.Delete, userId, repositoryId, tagId)
}

func tagEvent(eventManager schema.EventManager, eventType, userId, repositoryId, tagId string) error {
	fullRepoId := repositoryId + ":" + tagId
	event := schema.Event{
		Actor: userId,
		Type:  eventType,
		Object: schema.Object{
			Type: types.Tag,
			ID:   fullRepoId,
		},
	}
	return eventManager.CreateEvent(&event)
}

func GCDeleteBlobEvent(eventManager schema.EventManager, blobId string) error {
	return blobEvent(eventManager, types.Delete, types.GarbageCollector, blobId)
}

func blobEvent(eventManager schema.EventManager, eventType, userId, blobId string) error {
	event := schema.Event{
		Actor: userId,
		Type:  eventType,
		Object: schema.Object{
			Type: types.Blob,
			ID:   blobId,
		},
	}
	return eventManager.CreateEvent(&event)
}
