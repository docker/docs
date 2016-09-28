package client

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/bootstrap/config"
)

func (c *EngineClient) testMemory() error {
	meminfo, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return fmt.Errorf("Unable to determine system memory: %s", err)
	}

	value := 0.0
	for _, line := range strings.Split(string(meminfo), "\n") {
		l := strings.TrimSpace(line)
		if len(l) == 0 {
			continue
		}
		parsedLine := strings.Split(l, ":")
		if len(parsedLine) != 2 {
			continue
		}
		label := strings.ToLower(parsedLine[0])
		strValue := strings.TrimSpace(parsedLine[1])
		if label == "memtotal" {
			parsedValue := strings.Split(strValue, " ")
			if len(parsedValue) != 2 {
				continue
			}
			v, err := strconv.Atoi(parsedValue[0])
			if err != nil {
				log.Debugf("Failed to parse value in line: %s - %s", l, err)
				continue
			}
			units := parsedValue[1]

			if units != "kB" {
				// TODO If this hits, we can add parsing logic for different units
				return fmt.Errorf("Failed memory test, unexpected memory units: %s", strValue)
			}
			// Convert to GB
			value = float64(v) / 1024 / 1024
			if value > config.MinMemoryGB {
				log.Debugf("Your system meets minimum memory requirements:  %0.2f GB >= %0.2f GB",
					value, config.MinMemoryGB)
				return nil
			} else {
				log.Warnf("Your system does not have enough memory.  UCP suggests a minimum of %0.2f GB, but you only have %0.2f GB.  You may have unexpected errors.",
					config.MinMemoryGB, value)
				time.Sleep(2 * time.Second)
				return nil
			}
			break
		}
	}
	log.Debug("Unable to determine system memory")
	return fmt.Errorf("Unable to determine system memory")
}

func (c *EngineClient) testStorage() error {
	info, err := c.client.Info()
	if err != nil {
		return err
	}
	if strings.Contains(info.Driver, "devicemapper") {
		log.Warnf("Your system uses devicemapper.  We can not accurately detect available storage space.  Please make sure you have at least %0.2f GB available in %s",
			config.MinStorageGB, info.DockerRootDir)
		time.Sleep(2 * time.Second)
		return nil
	}

	sb := syscall.Statfs_t{}
	if err := syscall.Statfs("/", &sb); err != nil {
		return fmt.Errorf("Failed to determin available disk space: %s", err)
	}

	freeSpaceGB := float64(uint64(sb.Bsize)*sb.Bfree) / 1024 / 1024 / 1024

	if freeSpaceGB >= config.MinStorageGB {
		log.Debugf("Your system meets minimum storage requirements:  %0.2f GB >= %0.2f GB",
			freeSpaceGB, config.MinStorageGB)
		return nil
	}
	return fmt.Errorf("Your system does not have available disk space.  UCP requires a minimum of %0.2f GB, but you only have %0.2f GB",
		config.MinStorageGB, freeSpaceGB)
}

// Validate the base system meets minimum requirements
func (c *EngineClient) SystemValidation() error {
	log.Debug("Validating base system meets minimum requirements")
	if err := c.testMemory(); err != nil {
		return err
	}
	if err := c.testStorage(); err != nil {
		return err
	}
	return nil
}
