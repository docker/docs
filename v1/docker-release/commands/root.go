package commands

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var debugFlag, prodFlag bool
var awsBucket, awsKey, awsSecret, awsRegion, awsProfile string
var channelFlag, archFlag, buildFlag, humanFlag, sha1Flag string

// RootCmd is the main cobra command for all sub-commands
var RootCmd = &cobra.Command{
	Use:           "docker-release",
	Short:         "Create and manage desktop editions releases.",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debugFlag {
			logrus.SetLevel(logrus.DebugLevel)
		}
		logrus.Debugf("Debug flag: %t - Prod: %t", debugFlag, prodFlag)
		// Define the different endpoints
		s3URL = fmt.Sprintf(stageDomainTemplate, "")
		// Set all vars for the commands
		awsKey = viper.GetString("aws-access-key")
		awsSecret = viper.GetString("aws-secret-key")
		awsBucket = viper.GetString("aws-bucket")
		awsRegion = viper.GetString("aws-region")
		awsProfile = viper.GetString("aws-profile")
		if awsBucket == "" {
			awsBucket = s3StageBucket
		}
		if awsRegion == "" {
			awsRegion = s3Region
		}
		if awsProfile == "" {
			awsProfile = s3Profile
		}
		if prodFlag {
			awsBucket = s3LiveBucket
			s3URL = fmt.Sprintf(liveDomainTemplate, "")
		}
		if sha1Flag == "" {
			sha1Flag = os.Getenv("APPVEYOR_REPO_COMMIT")
		}
		s3DomainURL = fmt.Sprintf(s3DomainTemplate, awsBucket, getReleasePath(true))

		logrus.Debugf("URL: %s - Bucket: %s", s3URL, awsBucket)
	},
}

func init() {
	// Common flags
	RootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "D", false, "Enable debug mode")
	RootCmd.PersistentFlags().BoolVarP(&prodFlag, "prod", "p", false, "Use production environment")
	// AWS specific flags with Viper Binding
	RootCmd.PersistentFlags().String("aws-access-key", "", "AWS Access Key")
	viper.BindEnv("aws-access-key", "AWS_KEY")
	viper.BindPFlag("aws-access-key", RootCmd.PersistentFlags().Lookup("aws-access-key"))

	RootCmd.PersistentFlags().String("aws-secret-key", "", "AWS Secret Key")
	viper.BindEnv("aws-secret-key", "AWS_SECRET")
	viper.BindPFlag("aws-secret-key", RootCmd.PersistentFlags().Lookup("aws-secret-key"))

	RootCmd.PersistentFlags().String("aws-bucket", "", "AWS Bucket")
	viper.BindEnv("aws-bucket", "AWS_BUCKET")
	viper.BindPFlag("aws-bucket", RootCmd.PersistentFlags().Lookup("aws-bucket"))

	RootCmd.PersistentFlags().String("aws-region", "", "AWS Region")
	viper.BindEnv("aws-region", "AWS_DEFAULT_REGION")
	viper.BindPFlag("aws-region", RootCmd.PersistentFlags().Lookup("aws-region"))

	RootCmd.PersistentFlags().String("aws-profile", "", "AWS Profile")
	viper.BindEnv("aws-profile", "AWS_PROFILE")
	viper.BindPFlag("aws-profile", RootCmd.PersistentFlags().Lookup("aws-profile"))

	// RELEASE specific flags
	RootCmd.PersistentFlags().StringVar(&archFlag, "arch", os.Getenv("RELEASE_ARCH"), "Use `ARCH` architecture desired (mac or win)")
	RootCmd.PersistentFlags().StringVar(&channelFlag, "channel", os.Getenv("RELEASE_CHANNEL"), "Use `CHANNEL` channel for the release")
	RootCmd.PersistentFlags().StringVar(&buildFlag, "build", os.Getenv("RELEASE_BUILD"), "Use `BUILD` semver build version (sparkle version)")
	RootCmd.PersistentFlags().StringVar(&humanFlag, "human", os.Getenv("RELEASE_HUMAN"), "Use human version (sparkle shortVersion)")
	RootCmd.PersistentFlags().StringVar(&sha1Flag, "sha1", os.Getenv("CIRCLE_SHA1"), "Use `SHA1` for metadata (e80f2323)")

	RootCmd.AddCommand(listCmd, uploadCmd, publishCmd)
}
