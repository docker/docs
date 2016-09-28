package mock_test

import (
	"fmt"
	"time"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/orca"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/dockerhub"
	registry "github.com/docker/orca/registry/mock"
)

var (
	TestContainerAccessLabel = "testLabel"
	TestContainerAccessOwner = "testOwner"
	TestOrcaUser             = TestContainerAccessOwner
	TestContainerId          = "1234567890abcdefg"
	TestContainerName        = "test-container"
	TestContainerImage       = "test-image"
	TestRegistryConfig       = orca.RegistryConfig{
		Type:     "mock",
		ID:       "0",
		Name:     "Test Registry",
		URL:      "https://127.0.0.1",
		Insecure: true,
	}
	TestRepositoryUser  = "testuser"
	TestRegistryAddress = "127.0.0.1"
	TestRepositoryName  = fmt.Sprintf("%s/%s/testing", TestRegistryAddress, TestRepositoryUser)
	TestRegistryAuth    = fmt.Sprintf("{\"serveraddress\": \"%s\"}", TestRegistryAddress)
	TestRegistry        = &registry.MockRegistry{
		RegistryConfig: TestRegistryConfig,
	}
	TestContainers = []types.Container{
		{
			ID:    TestContainerId,
			Image: TestContainerImage,
			Labels: map[string]string{
				orca.UCPAccessLabel: TestContainerAccessOwner,
			},
		},
	}
	TestContainerInfo = &types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:      TestContainerId,
			Created: string(time.Now().UnixNano()),
			Name:    TestContainerName,
			Image:   TestContainerImage,
		},
		Config: &container.Config{
			Labels: map[string]string{
				orca.UCPOwnerLabel: TestContainerAccessOwner,
			},
		},
	}
	TestNode = &orca.Node{
		ID:   "0",
		Name: "testnode",
		Addr: "tcp://127.0.0.1:3375",
	}
	TestAccount = &auth.Account{
		Username: TestContainerAccessOwner,
		Password: "test",
		Admin:    false,
		Role:     auth.RestrictedControl,
	}
	TestTeam = &auth.Team{
		Id:          "0",
		Name:        "testTeam",
		Description: "Test Team",
		ManagedMembers: []string{
			TestAccount.Username,
		},
	}
	TestEvent = &orca.Event{
		Type:          "test-event",
		ContainerInfo: TestContainerInfo,
		Message:       "test message",
		Tags:          []string{"test-tag"},
	}
	TestWebhookKey = &dockerhub.WebhookKey{
		Image: "ehazlett/test",
		Key:   "abcdefg",
	}
	TestConsoleSession = &orca.ConsoleSession{
		ContainerID: "abcdefg",
		Token:       "1234567890",
	}
	TestContainerLogsToken = &orca.ContainerLogsToken{
		ContainerID: "12345",
		Token:       "1234567890",
	}
)

func getTestContainerInfo(id, name, image string, labels map[string]string) types.ContainerJSON {
	return types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:      id,
			Created: string(time.Now().UnixNano()),
			Name:    name,
			Image:   image,
		},
		Config: &container.Config{
			Labels: labels,
		},
	}
}

func getTestContainers() []types.ContainerJSON {
	return []types.ContainerJSON{
		getTestContainerInfo(TestContainerId, TestContainerName, TestContainerImage, nil),
	}
}

func getTestEvents() []*orca.Event {
	return []*orca.Event{
		TestEvent,
	}
}
