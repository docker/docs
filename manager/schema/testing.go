package schema

import (
	"sync"
	"time"

	"github.com/docker/dhe-deploy/shared/dtrutil"

	rethink "gopkg.in/dancannon/gorethink.v2"
)

var testOnce sync.Once
var testRethinkSession *rethink.Session
var testSchemaManager Manager
var testInitErr error
var dbTestingLock sync.Mutex

func TestSetup() (*rethink.Session, Manager, func(), error) {
	testOnce.Do(func() {
		if testInitErr = dtrutil.Poll(200*time.Millisecond, 15, func() error {
			var err error
			testRethinkSession, err = rethink.Connect(rethink.ConnectOpts{
				Addresses:     []string{"testrethink"},
				DiscoverHosts: false,
				MaxIdle:       5,
				MaxOpen:       10,
			})
			return err
		}); testInitErr != nil {
			return
		}
	})
	if testInitErr != nil {
		return nil, nil, nil, testInitErr
	}

	testSchemaManager = NewManager(func() (*rethink.Session, error) {
		return testRethinkSession, nil
	})

	if err := testSchemaManager.Initialize(); err != nil {
		return nil, nil, nil, err
	}

	if err := testSchemaManager.Migrate(); err != nil {
		return nil, nil, nil, err
	}

	dbTestingLock.Lock()
	unlockFunc := func() {
		dbTestingLock.Unlock()
	}

	return testRethinkSession, testSchemaManager, unlockFunc, nil
}
