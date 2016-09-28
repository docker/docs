package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/blang/semver"
	"github.com/docker/dhe-deploy/ha-integration/ha_utils"
	integrationutil "github.com/docker/dhe-deploy/integration/util"
	"github.com/docker/dhe-deploy/manager/versions"
	"github.com/docker/dhe-deploy/shared/containers"

	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	v2_0_0 = semver.Version{Major: 2, Minor: 0, Patch: 0}
	v2_1_0 = semver.Version{Major: 2, Minor: 1, Patch: 0, Pre: []semver.PRVersion{{VersionStr: "alpha"}}}
)

// GetExpectedDBTables lists the expected DBs in the rethink cluster, and which tables it should have -
// would be better if there were a different source of truth.
// Ideally, the Table objects in `rethinkutil` and the private one in `manager` can be merged and/or
// use an interface instead, that the Notary ones can implement.
// There are DTRTables in `manager`, `NotaryTables` in manager, and some job runner tables in
// jobrunner, and the `eventsTable` in `events.go`.
// `migrations.go` could also export the lists from the first rethink schema to the
// second, except it also seems to be missing some tables (properties table which is in DTRTables)?
func GetExpectedDBTables(version semver.Version) map[string][]string {
	switch {
	// note this order is important because we want to hit one before the other
	case version.GTE(v2_1_0):
		return map[string][]string{
			"dtr2": {
				"repositories",
				"namespace_team_access",
				"repository_team_access",
				"client_tokens",
				"events",
				"properties",
			},
			"jobrunner": {
				"joblogs",
				"crons",
				"action_configs",
				"jobs",
			},
			"notaryserver": {
				"tuf_files",
			},
			"notarysigner": {
				"private_keys",
			},
		}
	case version.GTE(v2_0_0):
		return map[string][]string{
			"dtr2": {
				"repositories",
				"namespace_team_access",
				"repository_team_access",
				"client_tokens",
				"properties",
			},
		}
	}
	return nil
}

// GetExpectedVolumes lists the expected volumes that should be in the cluster
func GetExpectedVolumes(version semver.Version) []containers.Volume {
	var volumes []containers.Volume
	if version.GTE(v2_0_0) {
		volumes = append(volumes,
			containers.EtcdVolume,
			containers.RethinkVolume,
			containers.RegistryVolume,
			containers.CAVolume,
		)
	}
	if version.GTE(v2_1_0) {
		volumes = append(volumes, containers.NotaryVolume)
	}
	return volumes
}

// GetAnyDTR returns true if there are either any dtr containers or volumes on the machine
func GetAnyDTR(machine ha_utils.Machine) (bool, error) {
	dtrContainers, err := machine.MachineSSH("sudo docker ps | grep 'dtr-' | cat")
	if err != nil {
		return false, err
	}

	dtrVolumes, err := machine.MachineSSH("sudo docker volume ls | grep 'dtr-' | cat")
	if err != nil {
		return false, err
	}

	return dtrVolumes != "" || dtrContainers != "", nil
}

// RequireAllDTRContainers returns true if all the expected DTR containers are there and running with the correct images
func RequireAllDTRContainers(machine ha_utils.Machine, replicaID string, imageArgs DTRImageArgs) (bool, error) {
	expectedContainers := make(map[string]string) // match those with the expected container names
	for _, expectedImageName := range imageArgs.DTRImages {
		for _, currentContainer := range containers.AllContainers {
			if strings.HasPrefix(expectedImageName, fmt.Sprintf("%s/dtr-%s", imageArgs.DTRRepo, currentContainer.Name)) {
				expectedContainers[fmt.Sprintf("dtr-%s-%s", currentContainer.Name, replicaID)] = expectedImageName
				break
			}
		}
	}

	// we need to map image names/tags to image IDs in case when we do a docker ps the images are listed as IDs rather
	// than the image name + tag
	dtrImages, err := machine.MachineSSH(`sudo docker images | grep dtr- | awk '{print $1":"$2" "$3}' | cat`)
	if err != nil {
		return false, err
	}
	imageNameToID := make(map[string]string)
	for _, outputLine := range strings.Split(dtrImages, "\n") {
		tokens := strings.Fields(strings.TrimSpace(outputLine))
		if len(tokens) != 2 {
			return false, fmt.Errorf("docker images output line unparsable: %s (%d tokens)", outputLine, len(tokens))
		}
		imageNameToID[strings.TrimSpace(tokens[0])] = strings.TrimSpace(tokens[1])
	}

	dtrContainers, err := machine.MachineSSH("sudo docker ps --format '{{.Names}} {{.Image}}' | grep 'dtr-' | cat")
	if err != nil {
		return false, err
	}

	foundContainersAndImages := make(map[string]string)
	for _, outputLine := range strings.Split(dtrContainers, "\n") {
		tokens := strings.Fields(strings.TrimSpace(outputLine))
		if len(tokens) != 2 {
			return false, fmt.Errorf("docker ps output line unparsable: %s (%d tokens)", outputLine, len(tokens))
		}
		foundContainersAndImages[strings.TrimSpace(tokens[0])] = strings.TrimSpace(tokens[1])
	}

	for expectedContainerName, expectedImageName := range expectedContainers {
		foundImageID, ok := imageNameToID[expectedImageName]
		if !ok {
			return false, fmt.Errorf("Image %s not even present on %s", expectedImageName, machine.GetName())
		}

		containerImage, ok := foundContainersAndImages[expectedContainerName]
		if !ok {
			return false, fmt.Errorf("Container %s not present on %s", expectedContainerName, machine.GetName())
		}

		if containerImage != expectedImageName && containerImage != foundImageID {
			return false, fmt.Errorf("Container %s using image %s instead of the expected %s (%s), on %s",
				expectedContainerName, containerImage, expectedImageName, foundImageID, machine.GetName())
		}
	}

	return true, nil
}

// RequireAllDTRVolumes returns true if all the expected DTR volumes are there
func RequireAllDTRVolumes(machine ha_utils.Machine, imageArgs DTRImageArgs) (bool, error) {
	currVersion, err := versions.TagToSemver(imageArgs.DTRTag)
	if err != nil {
		return false, err
	}

	dtrVolumes, err := machine.MachineSSH("sudo docker volume ls -q | grep 'dtr-' | cat")
	if err != nil {
		return false, err
	}

	volumeNames := strings.Split(dtrVolumes, "\n")
	for i, volumeString := range volumeNames {
		volumeNames[i] = strings.TrimSpace(volumeString)
	}

	expectedVolumes := make(map[string]bool)
	for _, expected := range GetExpectedVolumes(currVersion) {
		expectedVolumes[expected.Name] = false
		for _, foundName := range volumeNames {
			if strings.HasPrefix(foundName, expected.Name) {
				expectedVolumes[expected.Name] = true
			}
		}
	}

	for name, found := range expectedVolumes {
		if !found {
			return false, fmt.Errorf("%s volume not present", name)
		}
	}
	return true, nil
}

// FullSmokeTest confirms that the expected cluster setup is functional
// (and additionally that nodes that are not expected ot have DTR have neither containers nor volumes for it)
func FullSmokeTest(machines []ha_utils.Machine, imageArgs DTRImageArgs, t require.TestingT) {
	// TODO: Load-balancer health check (currently misconfigured?)
	// TODO: Push/Pull test - still load balancer is misconfigured
	// TODO: /meta/settings for ReplicaConfig

	PresenceTest(machines, imageArgs, t)
	APITest(machines, t)
	ClusterStatusTest(imageArgs, t)
}

// PresenceTest requires that every machine that is a DTR replica has all the required
// DTR containers and volumes.  Otherwise, it should not have any DTR containers or volumes.
func PresenceTest(machines []ha_utils.Machine, imageArgs DTRImageArgs, t require.TestingT) {
	for _, machine := range machines {
		if replicaID, ok := DefaultInstallation.Replicas[machine.GetName()]; !ok {
			// Confirm there is no DTR node on this machine
			replicaID, err := GetReplicaID(machine)
			assert.Nil(t, err)
			assert.Equal(t, "", replicaID)

			anyDTR, err := GetAnyDTR(machine)
			assert.Nil(t, err)
			assert.Equal(t, false, anyDTR)
		} else {
			alternateReplicaID, err := GetReplicaID(machine)
			require.Nil(t, err)

			// Confirm that the replicaID is what we would expect
			require.Equal(t, alternateReplicaID, replicaID)

			allDTRContainers, err := RequireAllDTRContainers(machine, replicaID, imageArgs)
			require.Nil(t, err)
			require.True(t, allDTRContainers)

			allDTRVolumes, err := RequireAllDTRVolumes(machine, imageArgs)
			require.Nil(t, err)
			require.True(t, allDTRVolumes)
		}
	}
}

// APITest requires that we can log in and get/set licenses on the cluster.
func APITest(machines []ha_utils.Machine, t require.TestingT) {
	// If we have a DTR
	if len(DefaultInstallation.Replicas) > 0 {
		// Try inserting a license
		err := DefaultInstallation.API.Login(ha_utils.GetAdminUser(), ha_utils.GetAdminPassword())
		require.Nil(t, err)
		licenseConfig, err := integrationutil.GetOnlineLicense()
		require.Nil(t, err)
		_, err = DefaultInstallation.API.SetLicenseSettings(licenseConfig)
		require.Nil(t, err)
	}
}

// ClusterStatusTest requires that all replicas in the cluster are healthy and if
// the expected tables are present.
func ClusterStatusTest(imageArgs DTRImageArgs, t require.TestingT) {
	// Skip if we don't have a DTR
	if len(DefaultInstallation.Replicas) == 0 {
		return
	}

	currVersion, err := versions.TagToSemver(imageArgs.DTRTag)
	require.NoError(t, err, "Unable to get the current semantic version from the current DTR tag")

	expectedTables := GetExpectedDBTables(currVersion)
	require.NotNil(t, expectedTables, "Unable to find list of expected tables for %s, interpreted as version %s", imageArgs.DTRTag, currVersion.String())

	checks := []func(*responses.ClusterStatus) error{
		checkReplicaHealth,
		checkEtcdMembership,
		checkRethinkClusterMembership,
		func(r *responses.ClusterStatus) error { return checkRethinkTables(r, expectedTables) },
	}

	// since we are checking cluster health, it may take a little while for everything to stabilize,
	// especially right after cluster install, so give it up to 2 minutes
	for i := 0; i < 8; i++ {
		if err != nil {
			logrus.Infof("Error checking cluster status.  Will retry...  %s", err.Error())
			time.Sleep(15 * time.Second)
		}

		err := DefaultInstallation.API.Login(ha_utils.GetAdminUser(), ha_utils.GetAdminPassword())
		require.NoError(t, err, "can't log into the cluster")
		status, err := DefaultInstallation.API.GetClusterStatus()
		require.NoError(t, err, "could not get cluster status")

		for _, check := range checks {
			if err = check(status); err != nil {
				break // break the check loop, not the retry loop
			}
		}
	}
	require.NoError(t, err, "cluster status hasn't stablized within 2 minutes")
}

func checkReplicaHealth(status *responses.ClusterStatus) error {
	// -- all replias should be healthy --
	// replica health is a map of replica IDs to health status ("OK", hopefully)
	if len(status.ReplicaHealth) != len(DefaultInstallation.Replicas) {
		return fmt.Errorf("Health check only lists %d replicas - expecting %d",
			len(status.ReplicaHealth), len(DefaultInstallation.Replicas))
	}
	for _, replicaID := range DefaultInstallation.Replicas {
		replicaStatus, ok := status.ReplicaHealth[replicaID]
		if !ok {
			return fmt.Errorf("Health of %s not listed", replicaID)
		}
		if replicaStatus != "OK" {
			return fmt.Errorf("%s not healthy: %s", replicaID, replicaStatus)
		}
	}
	return nil
}

func checkEtcdMembership(status *responses.ClusterStatus) error {
	// -- all the etcd containers in each replica should be in the etcd memberlist --
	untypedList, ok := status.EtcdStatus["members"].([]interface{})
	if !ok {
		return fmt.Errorf("cluster status missing etcd member list")
	}
	found := make(map[string]struct{})
	for _, memberOpaque := range untypedList {
		memberMap, ok := memberOpaque.(map[string]interface{})
		if !ok {
			return fmt.Errorf("cluster status etcd member list has an invalid object")
		}

		// if the item doesn't exist in the map, the casting the nil object to a string results in
		// the empty string
		memberName, ok := memberMap["name"].(string)
		if !ok {
			return fmt.Errorf("cluster status etcd member list has an invalid object")
		}
		found[memberName] = struct{}{}
	}

	if len(found) != len(DefaultInstallation.Replicas) {
		return fmt.Errorf("%d members in the etcd cluster; expected %d", len(found), len(DefaultInstallation.Replicas))
	}
	for machineName, replicaID := range DefaultInstallation.Replicas {
		if _, ok := found[fmt.Sprintf("dtr-etcd-%s", replicaID)]; !ok {
			return fmt.Errorf("etcd from %s (replica %s) not a member of the etcd cluster", machineName, replicaID)
		}
	}
	return nil
}

func checkRethinkClusterMembership(status *responses.ClusterStatus) error {
	// --  all rethink containers in each replica should be in the rethink cluster --
	untypedList, ok := status.RethinkSystemTables["server_status"].([]interface{})
	if !ok {
		return fmt.Errorf("cluster status missing rethink server status list")
	}
	found := make(map[string]struct{})
	for _, memberOpaque := range untypedList {
		memberMap, ok := memberOpaque.(map[string]interface{})
		if !ok {
			return fmt.Errorf("cluster status rethink server status list has an invalid object")
		}

		// if the item doesn't exist in the map, the casting the nil object to a string results in
		// the empty string
		memberName, ok := memberMap["name"].(string)
		if !ok {
			return fmt.Errorf("cluster status rethink server status list has an invalid object")
		}
		found[memberName] = struct{}{}
	}

	if len(found) != len(DefaultInstallation.Replicas) {
		return fmt.Errorf("%d members in the rethink cluster; expected %d", len(found), len(DefaultInstallation.Replicas))
	}
	for machineName, replicaID := range DefaultInstallation.Replicas {
		if _, ok := found[fmt.Sprintf("dtr_rethinkdb_%s", replicaID)]; !ok {
			return fmt.Errorf("rethinkdb from %s (replica %s) not a member of the rethink cluster", machineName, replicaID)
		}
	}
	return nil
}

func checkRethinkTables(status *responses.ClusterStatus, expectedTables map[string][]string) error {
	// --  all the DBs and tables should be in the rethink cluster --
	untypedList, ok := status.RethinkSystemTables["table_config"].([]interface{})
	if !ok {
		return fmt.Errorf("cluster status missing rethink system tables list")
	}
	// map of DB names to a map of table names to the shards object, which gives information
	// about replicas per table
	foundTables := make(map[string]map[string][]interface{})

	for _, tableOpaque := range untypedList {
		tableMap, ok := tableOpaque.(map[string]interface{})
		if !ok {
			return fmt.Errorf("cluster status rethink system table list has an invalid object")
		}

		dbName, ok := tableMap["db"].(string)
		if !ok {
			return fmt.Errorf("cluster status rethink system table list contains invalid object that doesn't have db")
		}

		tableName, ok := tableMap["name"].(string)
		if !ok {
			return fmt.Errorf("cluster status rethink system table list contains invalid object that doesn't have name")
		}

		shardsOpaque, ok := tableMap["shards"].([]interface{})
		if !ok || len(shardsOpaque) == 0 {
			return fmt.Errorf("cluster status rethink system table list contains invalid object that doesn't have shards")
		}

		if foundTables[dbName] == nil {
			foundTables[dbName] = make(map[string][]interface{})
		}
		// we want to save the shard info for the tables we found because we only care about testing the replication factor
		// for tables we care about
		foundTables[dbName][tableName] = shardsOpaque
	}

	for dbName, tables := range expectedTables {
		tableMap, ok := foundTables[dbName]
		if !ok {
			return fmt.Errorf("could not find expected DB %s", dbName)
		}
		for _, tableName := range tables {
			shardsOpaque, ok := tableMap[tableName]
			if !ok {
				return fmt.Errorf("%s table %s not found in rethink cluster", dbName, tableName)
			}

			// Ok, we found a table we expected.  Its replication factor should be <#DTR-replicas>
			for i, shardOpaque := range shardsOpaque {
				shardMap, ok := shardOpaque.(map[string]interface{})
				if !ok {
					return fmt.Errorf("cluster status rethink system table list contains invalid object with an invalid shards object")
				}
				replicas, ok := shardMap["replicas"].([]interface{})
				if !ok {
					return fmt.Errorf("cluster status rethink system table list contains invalid object with an invalid shards object")
				}

				if len(replicas) != len(DefaultInstallation.Replicas) {
					return fmt.Errorf("%d replicas for table %s/%s (shard %d) in the rethink cluster; expected %d",
						len(replicas), dbName, tableName, i+1, len(DefaultInstallation.Replicas))

				}
			}

		}
	}
	return nil
}
