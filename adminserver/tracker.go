package adminserver

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy"
	notaryserver "github.com/docker/notary/server/storage"
	tuf "github.com/docker/notary/tuf/data"
	"github.com/mssola/user_agent"
	mixpanel "github.com/pdevine/go-mixpanel"
	"gopkg.in/dancannon/gorethink.v2"
)

func getRepos(as *AdminServer, usage map[string]interface{}) error {
	var publicRepos, privateRepos int
	var err error
	publicRepos, err = as.repositoryManager.CountPublicRepositories()
	if err != nil {
		return err
	}

	privateRepos, err = as.repositoryManager.CountPrivateRepositories()
	if err != nil {
		return err
	}

	usage["repos"] = publicRepos + privateRepos
	usage["public_repos"] = publicRepos
	usage["private_repos"] = privateRepos
	return nil
}

func getKernelInfo(usage map[string]interface{}) error {
	cmd := exec.Command("uname", "-srmo")
	if output, err := cmd.Output(); err != nil {
		log.Errorf("Couldn't get kernel info: %s", err)
		return err
	} else {
		usage["kernel_info"] = string(output)
		return nil
	}
}

func getStorageType(as *AdminServer, usage map[string]interface{}) error {
	s := as.settingsStore
	if registryConfig, err := s.RegistryConfig(); err != nil {
		return err
	} else if registryConfig != nil && registryConfig.Storage != nil {
		usage["storage"] = registryConfig.Storage.Type()
	}

	return nil
}

func getID(as *AdminServer) (string, error) {
	haConfig, err := as.settingsStore.HAConfig()
	if err != nil {
		log.Errorf("Failed to retrieve hub config: %s", err.Error())
		return "", err
	}
	// Safe to assume there's at least one replica config whose id we want
	replicaConfigs := haConfig.ReplicaConfig
	if len(replicaConfigs) == 0 {
		// Just to make integration tests pass since no replicas are created
		return "", nil
	}
	var replicaIDs []string
	for id := range replicaConfigs {
		replicaIDs = append(replicaIDs, id)
	}
	sort.Strings(replicaIDs)
	return replicaIDs[0], nil
}

func getReplicaCount(as *AdminServer, usage map[string]interface{}) error {
	haConfig, err := as.settingsStore.HAConfig()
	if err != nil {
		log.Errorf("Failed to retrieve hub config: %s", err.Error())
		return err
	}
	usage["replica_count"] = len(haConfig.ReplicaConfig)
	return nil
}

func getLicenseInfo(as *AdminServer, usage map[string]interface{}) {
	licenseIsValid := as.licenseChecker.IsValid()
	if licenseIsValid {
		usage["license_type"] = as.licenseChecker.LicenseType()
		usage["license_tier"] = as.licenseChecker.LicenseTier()
		if !as.AnonymizeAnalytics() {
			usage["license_id"] = as.licenseChecker.GetLicenseID()
		}
	} else {
		licenseExpired := as.licenseChecker.IsExpired()
		if licenseExpired {
			usage["license_type"] = "Expired"
			usage["license_tier"] = as.licenseChecker.LicenseTier()
		} else {
			usage["license_type"] = "Invalid license"
		}
	}
}

func getMedian(numbers []int) float64 {
	middle := len(numbers) / 2
	median := float64(numbers[middle])
	if len(numbers)%2 == 0 {
		median = (median + float64(numbers[middle-1])) / 2
	}
	return median
}

func setNotaryMetricStats(metric string, usage map[string]interface{}, numbers []int) {
	if len(numbers) == 0 {
		return
	}
	sort.Ints(numbers)
	minStat := fmt.Sprintf("notary_min_%s", metric)
	maxStat := fmt.Sprintf("notary_max_%s", metric)
	medianStat := fmt.Sprintf("notary_median_%s", metric)
	usage[minStat] = numbers[0]
	usage[maxStat] = numbers[len(numbers)-1]
	usage[medianStat] = getMedian(numbers)
}

func getRepoDelegationsTags(as *AdminServer, repo string) (minNumOfSignedTags, maxNumOfSignedTags int, err error) {
	var targetFiles []notaryserver.RDBTUFFile
	// effectively a set to keep track of which delegation roles have been already processed
	delegationRoles := make(map[string]struct{})

	query := gorethink.DB(deploy.NotaryServerDBName).Table(notaryserver.TUFFilesRethinkTable.Name).Filter(map[string]interface{}{
		"gun": repo,
	}).Filter(func(file gorethink.Term) gorethink.Term {
		// each gun could have multiple delegation roles each of which is prefixed with 'targets/'
		return file.Field("role").Match("^targets/")
	}).OrderBy(gorethink.Desc("version"))
	res, err := query.Run(as.rethinkSession)
	if err != nil {
		log.Errorf("Couldn't run the notary query: %s", err.Error())
		return 0, 0, err
	}

	err = res.All(&targetFiles)
	if err != nil {
		log.Errorf("Couldn't parse notary's rethinkdb response: %s", err.Error())
		return 0, 0, err
	}

	for _, targetFile := range targetFiles {
		var signedTargets tuf.SignedTargets
		role := targetFile.Role
		if _, alreadyProcessed := delegationRoles[role]; alreadyProcessed {
			continue
		}
		delegationRoles[role] = struct{}{}

		err = json.Unmarshal(targetFile.Data, &signedTargets)
		if err != nil {
			log.Infof("Couldn't get signed targets: %s", err.Error())
			return 0, 0, err
		}

		// targets correspond to tags
		tagCount := len(signedTargets.Signed.Targets)
		if minNumOfSignedTags == 0 {
			minNumOfSignedTags = tagCount
			maxNumOfSignedTags = tagCount
			continue
		}

		if tagCount < minNumOfSignedTags {
			minNumOfSignedTags = tagCount
		} else if tagCount > maxNumOfSignedTags {
			maxNumOfSignedTags = tagCount
		}
	}
	return minNumOfSignedTags, maxNumOfSignedTags, nil
}

func getNotaryInformation(as *AdminServer, usage map[string]interface{}) error {
	var err error
	var targetFiles []notaryserver.RDBTUFFile
	// effectively a set to keep track of which repos have been already visited
	notaryRepos := make(map[string]struct{})

	// since we don't know the version at any given time we sort in descending order of version and disregard
	// any GUNs (read: repos) after coming across and processing the first instance of any GUN
	query := gorethink.DB(deploy.NotaryServerDBName).Table(notaryserver.TUFFilesRethinkTable.Name).Filter(map[string]interface{}{
		"role": "targets",
	}).OrderBy(gorethink.Desc("version"))
	res, err := query.Run(as.rethinkSession)
	if err != nil {
		log.Errorf("Couldn't run the notary query: %s", err.Error())
		return err
	}

	err = res.All(&targetFiles)
	if err != nil {
		log.Errorf("Couldn't parse notary's rethinkdb response: %s", err.Error())
		return err
	}

	if len(targetFiles) == 0 {
		usage["notary_in_usage"] = false
		return nil
	}
	usage["notary_in_usage"] = true

	notaryDelegationRoles := []int{}
	notarySignedTags := []int{}
	notaryDelegationsTags := []int{}
	repo := ""
	for _, targetFile := range targetFiles {
		var signedTargets tuf.SignedTargets
		repo = targetFile.Gun // gun corresponds to a repo
		if _, alreadyProcessed := notaryRepos[repo]; alreadyProcessed {
			continue
		}
		notaryRepos[repo] = struct{}{}

		err := json.Unmarshal(targetFile.Data, &signedTargets)
		if err != nil {
			log.Infof("Couldn't get signed targets: %s", err.Error())
			return err
		}

		notaryDelegationRoles = append(notaryDelegationRoles, len(signedTargets.Signed.Delegations.Roles))
		notarySignedTags = append(notarySignedTags, len(signedTargets.Signed.Targets))

		if minDelegationsTagsCount, maxDelegationsTagsCount, err := getRepoDelegationsTags(as, repo); err != nil {
			log.Debugf("Couldn't get tag count for delegations of repo %s: %s", repo, err.Error())
			return err
		} else {
			notaryDelegationsTags = append(notaryDelegationsTags, minDelegationsTagsCount, maxDelegationsTagsCount)
		}
	}
	usage["notary_repos"] = len(notaryRepos)
	setNotaryMetricStats("delegation_roles", usage, notaryDelegationRoles)
	setNotaryMetricStats("signed_tags", usage, notarySignedTags)
	setNotaryMetricStats("delegations_tags", usage, notaryDelegationsTags)
	return nil
}

func (as *AdminServer) reportAnalytics() bool {
	hubConfig, err := as.settingsStore.UserHubConfig()
	if err != nil {
		return false
	} else {
		return hubConfig.ReportAnalytics
	}
}

func (as *AdminServer) AnonymizeAnalytics() bool {
	hubConfig, err := as.settingsStore.UserHubConfig()
	if err != nil {
		return false
	} else {
		return hubConfig.AnonymizeAnalytics
	}
}

func sanitizeUA(useragent string, usage map[string]interface{}) map[string]interface{} {
	ua := user_agent.New(useragent)
	usage["Operating System"] = ua.OS()

	name, version := ua.Browser()
	versionNum := strings.Split(version, ".")[0]
	usage["Browser"] = name
	usage["Browser Version"] = versionNum

	return usage
}

func (as *AdminServer) sendClientAnalytics(useragent string, sourceIP string) {
	if !as.reportAnalytics() {
		return
	}

	id, err := getID(as)
	usage := make(map[string]interface{})
	usage = sanitizeUA(useragent, usage)
	usage["ip"] = sourceIP
	// wrap around goroutine in case request takes a long time or even times out
	go func(as *AdminServer, usage map[string]interface{}) {
		if err = as.mixpanelClient.Track(id, "Client usage", usage); err != nil {
			log.WithField("error", err).Error("Can't report client analytics")
		}
	}(as, usage)
}

func (as *AdminServer) setupTracking() {
	ticker := time.NewTicker(time.Hour)
	mp := mixpanel.NewMixpanel(deploy.MixpanelToken)
	as.mixpanelClient = mp

	go func(as *AdminServer) {
		// Initial report
		reportUsage(as)
		for {
			select {
			// send report every hour if reporting is turned on
			case <-ticker.C:
				if as.reportAnalytics() {
					reportUsage(as)
				}
			}
		}
	}(as)
}

func reportUsage(as *AdminServer) {

	usage := make(map[string]interface{})
	usage["DTR_version"] = deploy.Version
	getLicenseInfo(as, usage)

	if err := getStorageType(as, usage); err != nil {
		log.Debugf("Couldn't get storage type: %s", err.Error())
		return
	}

	if err := getRepos(as, usage); err != nil {
		log.Debugf("Couldn't get repo information: %s", err.Error())
		return
	}

	if err := getKernelInfo(usage); err != nil {
		log.Debugf("Couldn't get kernel information: %s", err.Error())
		return
	}

	if err := getReplicaCount(as, usage); err != nil {
		log.Debugf("Couldn't get replica information: %s", err.Error())
		return
	}

	if err := getNotaryInformation(as, usage); err != nil {
		log.Debugf("Couldn't get notary information: %s", err.Error())
		return
	}

	id, err := getID(as)
	if err != nil {
		log.Debugf("Couldn't get replica unique instance ID: %s", err.Error())
		return
	}

	if err := as.mixpanelClient.Track(id, "Server stats", usage); err != nil {
		log.Debugf("Can't report analytics: %s", err.Error())
	}
}
