package commands

import (
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
)

var releaseTemplate = `Available Releases
--------------
{{range $element := .}} * {{ $element }}
{{ end }}
`

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all releases",
	Run:   list,
}

func list(cmd *cobra.Command, args []string) {
	channelRe := ""
	if channelFlag != "" {
		channelRe = channelFlag
	}
	archRe := ""
	if archFlag != "" {
		if channelRe == "" {
			channelRe = `[a-z]+`
		}
		archRe = archFlag + `\/`
	}
	re := regexp.MustCompile(`(?i)` + archRe + channelRe + `\/` + buildFlag + `[0-9.]*\/$`)
	logrus.Debug("Regex: ", re.String())
	sess := getS3Session()
	svc := s3.New(sess)
	releases := []string{}
	err := svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(awsBucket),
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		for _, obj := range p.Contents {
			if !strings.Contains(*obj.Key, "logs/") {
				logrus.Debugf("= Checking = %s", *obj.Key)
				if re.Match([]byte(*obj.Key)) {
					releases = append(releases, *obj.Key)
				}
			}
		}
		return true
	})
	if err != nil {
		logrus.Fatal("failed to list objects: ", err)
	}
	printData(releaseTemplate, releases)
}
