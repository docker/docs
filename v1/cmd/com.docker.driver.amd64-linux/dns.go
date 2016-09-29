package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	datakit "github.com/docker/datakit/api/go-datakit"
	"github.com/docker/pinata/v1/apple"
	"golang.org/x/net/context"
	"strings"
	"time"
)

func refreshDNSServers(ctx context.Context, client *datakit.Client) {
	dns := apple.GetDNSServers()
	search := apple.GetDNSSearchDomains()
	commitMsg := fmt.Sprintf("Settings Changed %s", time.Now().Format(time.RFC822Z))
	t, err := datakit.NewTransaction(ctx, client, "master", "refreshDNSServers")
	if err != nil {
		logrus.Fatal(err)
	}
	key := []string{driver, "slirp", "dns"}
	if len(dns) == 0 {
		logrus.Info("SC database lists no DNS servers: removing overrides from db")
		if err = t.Remove(ctx, key); err != nil {
			logrus.Fatal(err)
		}
	} else {
		logrus.Info("SC database lists DNS servers: ", strings.Join(dns, ", "))
		logrus.Info("SC database lists search domains: ", strings.Join(search, " "))
		resolvConf := ""
		for _, server := range dns {
			if resolvConf != "" {
				resolvConf += "\n"
			}
			resolvConf += fmt.Sprintf("nameserver %s", server)
		}
		if len(search) > 0 {
			resolvConf += "\nsearch " + strings.Join(search, " ") + "\n"
		}
		if err = t.Write(ctx, key, resolvConf); err != nil {
			logrus.Fatal(err)
		}
	}
	if err = t.Commit(ctx, commitMsg); err != nil {
		logrus.Fatal(err)
	}
}
