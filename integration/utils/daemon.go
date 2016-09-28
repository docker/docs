package utils

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"strings"
	"time"
)

func StopDockerDaemon(machine Machine) error {
	if _, err := machine.MachineSSH("which systemctl"); err == nil {
		log.Info("Appears to be a systemctl based system")
		cmd := "sudo systemctl stop docker.service"
		out, err := machine.MachineSSH(cmd)
		if err != nil {
			return err
		}
		log.Info(out)
		return nil
	} else { // TODO Other variants?
		log.Info("Appears to be a boot2docker based")
		log.Infof("Stopping daemon on %s", machine.GetName())
		out, err := machine.MachineSSH("sudo /usr/local/etc/init.d/docker stop")
		if err != nil {
			return err
		}
		log.Info(out)
		log.Infof("Waiting for it to stop on %s", machine.GetName())
		for i := 0; i < 60; i++ {
			// NOTE: We could use the /var/run/docker.pid instead of "docker daemon"...
			out, _ := machine.MachineSSH(`ps -e | grep "docker daemon" | grep -v grep ; true`)
			if strings.TrimSpace(out) == "" {
				log.Infof("It has Stopped on %s", machine.GetName())
				return nil
			}
			log.Debug(out) // XXX comment out
			time.Sleep(1 * time.Second)
		}
		return fmt.Errorf("Daemon never stopped")
	}
}
func StartDockerDaemon(machine Machine) error {
	engineClient, err := machine.GetClient()
	if err != nil {
		return err
	}
	if _, err := machine.MachineSSH("which systemctl"); err == nil {
		log.Info("Appears to be a systemctl based system")
		cmd := "sudo systemctl start docker.service"
		out, err := machine.MachineSSH(cmd)
		if err != nil {
			return err
		}
		log.Info(out)
		return nil
	} else { // TODO Other variants?
		log.Info("Appears to be a boot2docker based")
		// HACK / Workaround
		// Our two socat's for cfssl don't exit cleanly and the daemon
		// forgets to remove it's pid file, leading to a failure to restart
		machine.MachineSSH("sudo rm -f /var/run/docker.pid")
		time.Sleep(500 * time.Millisecond)

		log.Infof("Starting the daemon on %s", machine.GetName())
		out, err := machine.MachineSSH("sudo /usr/local/etc/init.d/docker start")
		if err != nil {
			return err
		}
		log.Info(out)
	}
	log.Infof("Waiting for the daemon to come back on %s", machine.GetName())
	// Now wait for the daemon to restart and be happy
	// This can take a long time - the daemon sometimes takes a while to get etcd
	// started, and all 3 etcds have to recover, talk to eachother, then it'll come back
	for i := 0; i < 120; i++ {
		_, err := engineClient.Info()
		if err == nil {
			log.Infof("Daemon on %s came back", machine.GetName())
			return nil
		}
		time.Sleep(time.Second)
	}
	// Grab the tail of the daemon log since we're hosed...
	out, _ := machine.MachineSSH("tail /var/log/docker.log")
	log.Info(out)
	return fmt.Errorf("Daemon on %s never came back", machine.GetName())
}

func RestartDockerDaemon(machine Machine) error {
	engineClient, err := machine.GetClient()
	if err != nil {
		return err
	}
	if _, err := machine.MachineSSH("which systemctl"); err == nil {
		log.Info("Appears to be a systemctl based system")
		cmd := "sudo systemctl restart docker.service"
		out, err := machine.MachineSSH(cmd)
		if err != nil {
			return err
		}
		log.Info(out)
	} else { // TODO Other variants?
		log.Info("Appears to be a boot2docker based")
		log.Infof("Stopping daemon on %s", machine.GetName())
		out, err := machine.MachineSSH("sudo /usr/local/etc/init.d/docker stop")
		if err != nil {
			return err
		}
		log.Infof("Waiting for it to stop on %s", machine.GetName())
		for i := 0; i < 60; i++ {
			// NOTE: We could use the /var/run/docker.pid instead of "docker daemon"...
			out, _ := machine.MachineSSH(`ps -e | grep "docker daemon" | grep -v grep; true`)
			if strings.TrimSpace(out) == "" {
				log.Infof("It has Stopped on %s", machine.GetName())
				break
			}
			time.Sleep(time.Second)
		} // Proceed anyway even if we didn't detect stop, and hope for the best

		// HACK / Workaround
		// Our two socat's for cfssl don't exit cleanly and the daemon
		// forgets to remove it's pid file, leading to a failure to restart
		time.Sleep(500 * time.Millisecond)
		machine.MachineSSH("sudo rm -f /var/run/docker.pid")
		time.Sleep(500 * time.Millisecond)

		log.Infof("Starting the daemon on %s", machine.GetName())
		out, err = machine.MachineSSH("sudo /usr/local/etc/init.d/docker start")
		if err != nil {
			return err
		}
		log.Info(out)
	}
	log.Infof("Waiting for the daemon to come back on %s", machine.GetName())
	// Now wait for the daemon to restart and be happy
	// This can take a long time - the daemon sometimes takes a while to get etcd
	// started, and all 3 etcds have to recover, talk to eachother, then it'll come back
	for i := 0; i < 120; i++ {
		_, err := engineClient.Info()
		if err == nil {
			log.Infof("Daemon on %s came back", machine.GetName())
			return nil
		}
		time.Sleep(time.Second)
	}
	// Grab the tail of the daemon log since we're hosed...
	out, _ := machine.MachineSSH("tail /var/log/docker.log")
	log.Info(out)
	return fmt.Errorf("Daemon on %s never came back", machine.GetName())
}
