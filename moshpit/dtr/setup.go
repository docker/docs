package dtr

import (
	"fmt"

	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"

	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/integration/apiclient"
	"github.com/docker/dhe-deploy/pkg/moshpit-framework"
	"github.com/docker/dhe-deploy/shared/dtrutil"

	enzierrors "github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/forms"
)

type SetupConfig struct {
	RepoName       string
	DTRURL         string
	DTRCA          string
	DTRInsecureTLS bool
	Username       string
	Password       string

	// These were added for mass-populating an instance
	NumUsers        int
	NumReposPerUser int
}

func Setup(ctx context.Context, setupConfig string) error {
	// ** parse config
	config := SetupConfig{}
	log := moshpit.LoggerFromCtx(ctx)
	err := yaml.Unmarshal([]byte(setupConfig), &config)
	if err != nil {
		log.Fatal(err)
	}

	// ** create dtr client
	httpClient, err := dtrutil.HTTPClient(config.DTRInsecureTLS, config.DTRCA)
	if err != nil {
		return err
	}
	api := apiclient.New(config.DTRURL, 3, httpClient)
	err = api.Login(config.Username, config.Password)
	if err != nil {
		return err
	}

	// ** create repo
	_, err = api.GetRepository(config.Username, config.RepoName)
	if err, ok := err.(*apiclient.APIError); ok {
		if len(err.Errors) > 0 && err.Errors[0].Code == errors.ErrorCodeNoSuchRepository.Code {
			// the repo hasn't been created yet, so let's create it
			_, err := api.CreateRepository(config.Username, config.RepoName, "", "", "private")
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else if err != nil {
		return err
	}

	// Create users for mass-populating
	for i := 1; i < config.NumUsers+1; i++ {
		_, err := api.EnziSession().GetAccount(fmt.Sprintf("testuser%06d", i))
		if apiErrs, ok := err.(*enzierrors.APIErrors); ok {
			if len(apiErrs.Errors) > 0 && apiErrs.Errors[0].Code == enzierrors.NoSuchAccount("").Code {
				_, err := api.EnziSession().CreateAccount(forms.CreateAccount{
					Name:     fmt.Sprintf("testuser%06d", i),
					Password: "password",
				})
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	// Create repos for mass-populating
	// ** create repo
	for i := 1; i < config.NumUsers+1; i++ {
		for j := 0; j < config.NumReposPerUser; j++ {
			_, err = api.GetRepository(fmt.Sprintf("testuser%06d", i), fmt.Sprintf("repo%d", j))
			if err, ok := err.(*apiclient.APIError); ok {
				if len(err.Errors) > 0 && err.Errors[0].Code == errors.ErrorCodeNoSuchRepository.Code {
					// the repo hasn't been created yet, so let's create it
					_, err := api.CreateRepository(fmt.Sprintf("testuser%06d", i), fmt.Sprintf("repo%d", j), "", "", "private")
					if err != nil {
						return err
					}
				} else {
					return err
				}
			} else if err != nil {
				return err
			}
		}
	}

	return nil
}
