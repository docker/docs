# dsinfo
Docker/System information collection script to assist with troubleshooting.

### Note

Debug logging is extremely helpful when troubleshooting problems. If your Docker daemon is not running in debug mode this script will prompt and ask if you want to enable it. If you choose so, it will attempt to enable debug mode by adjusting the Docker config file (add the **-D** flag) and restarting the docker daemon.

After troubleshooting remember to revert the changes made to the config file and restart the daemon.

### Usage

#### Running as a container:

Hint: don't use the `-t` flag to `docker run` as that will corrupt the tarfile

```
docker run --rm  \
    -v /boot:/boot \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /var/lib/docker:/var/lib/docker \
    -v /var/log:/var/log \
    -v /etc/sysconfig:/etc/sysconfig \
    -v /etc/default:/etc/default \
    dsinfo > dump.tgz
```

#### On a Linux Machine
```
curl -O https://raw.githubusercontent.com/dockersupport/dsinfo/master/dsinfo.sh
chmod +x dsinfo.sh
sudo ./dsinfo.sh
```

or

```
wget https://raw.githubusercontent.com/dockersupport/dsinfo/master/dsinfo.sh
chmod +x dsinfo.sh
sudo ./dsinfo.sh
```

#### Using Boot2Docker

This script needs to be run in the Boot2Docker virtual machine. You can ssh into your Boot2Docker VM with:

- With Bash: `ssh -i ~/.ssh/id_boot2docker docker@$(boot2docker ip)`

- With [Fish](http://fishshell.com/): `ssh -i ~/.ssh/id_boot2docker docker@(boot2docker ip)`

### Notes

This script creates the directories:
  - /tmp/dsinfo
  - /tmp/dsinfo/dhe
  - /tmp/dsinfo/dhe/generatedConfigs
  - /tmp/dsinfo/dhe/logs
  - /tmp/dsinfo/inspect
  - /tmp/dsinfo/logs

This script collects the following:
  - Docker daemon configuration and logs
  - Inspect results and logs from the last 20 containers
  - Miscellaneous system information (Output to: /tmp/dsinfo/report.md)

All files/directories are compressed to: /tmp/dsinfo.tar.gz

---------------------------------------------------------------------------------

*** Important ***

Before sharing the dsinfo.tar.gz archive, review all collected files for
private information and edit/delete if necessary.

If you do edit/remove any files, recreate the tar file with the following command:

  sudo tar -zcf /tmp/dsinfo.tar.gz /tmp/dsinfo
