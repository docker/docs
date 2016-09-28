package utils

import (
	"github.com/stretchr/testify/suite"
)

// Used to inject additional logic at various points during the suite setup
// Every OrcaTestSuite is an OTSC, with default implementations
// that generally work.  Specific test suites can override to
// change behavior.
//
// Note: You MUST call mysuite.Init(mysuite) before starting tests
//       to get everything wired up properly
type OTSC interface { // Orca Test Suite Callbacks

	// Establish the number of nodes for the suite
	GetNodeCounts() (controllerCount int, workerCount int)

	// Called before loading images on every machine
	PreImageLoad() error

	// Called before install/join on a node (after loading images)
	PreInstall(m Machine) error
	InstallArgs(m Machine) []string
	JoinReplicaArgs(m Machine, serverURL, fingerprint string) [3][]string
	JoinWorkerArgs(m Machine, serverURL, fingerprint string) [3][]string

	// Called after install/join on each node (after loading images)
	PostInstall(m Machine) error

	// Called after all nodes are joined, to perform initial configuration tasks
	InitialConfig() error
}

// Dynamic multi-node suites
type OrcaTestSuite struct {
	suite.Suite
	ControllerMachines []Machine
	WorkerMachines     []Machine

	IsLDAP             bool // When set true, this indicates this cluster is operating in LDAP mode
	IsDiscoveryMissing bool // When set true, this indicates this cluster doesn't have discovery wired up

	// This map is used to store arbitrary blobs of data during a test run
	// Examples might be backups taken during setup that are later used in specific test
	// cases that restore from that backup.  Use descriptive keys so it doesn't get confusing.
	Data map[string][]byte

	// Go doesn't support polymorphism of structs, you have to use
	// interfaces.  But you can't implement functions directly on
	// interfaces, only structs, which means you can't easily have function
	// logic in a base struct that defers to overridden functions in the
	// sub structs.  To workaround this limitation, we keep a reference to
	// "ourself" typed as the interface, so that we can call the various
	// routines within setup and get the subtypes implementation if it has
	// one.
	self OTSC
}
