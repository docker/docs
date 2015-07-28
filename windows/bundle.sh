#!/bin/sh
set -e

# Script to grab binaries that are going to be bundled with windows installer.
# Note to maintainers: Update versions used below with newer releases

# boot2dockerIso=1.8.0
docker=1.8.0-rc1
dockerMachine=0.4.0-rc1
kitematic=0.8.0-rc2
vbox=5.0.0
vboxRev=101573
msysGit=1.9.5-preview20150319

boot2dockerIsoSrc=boot2docker
dockerBucket=test.docker.com

set -x
rm -rf bundle
mkdir bundle
cd bundle

(
	mkdir -p docker
	cd docker

	curl -fsSL -o docker.exe "https://${dockerBucket}/builds/Windows/x86_64/docker-${docker}.exe"
	curl -fsSL -o docker-machine.exe "https://github.com/docker/machine/releases/download/v${dockerMachine}/docker-machine_windows-amd64.exe"
)

(
	mkdir -p kitematic
	cd kitematic

	curl -fsSL -o kitematic-setup.exe "https://github.com/kitematic/kitematic/releases/download/v${kitematic}/KitematicSetup-${kitematic}-Windows-Alpha.exe"
)

(
	mkdir -p Boot2Docker
	cd Boot2Docker

	# curl -fsSL -o boot2docker.iso "https://github.com/${boot2dockerIsoSrc}/boot2docker/releases/download/v${boot2dockerIso}/boot2docker.iso"
	curl -fsSL -o boot2docker-virtualbox.iso "https://s3.amazonaws.com/toolbox-rc/boot2docker-virtualbox-1.8.0-dev.iso"
)

(
	mkdir -p msysGit
	cd msysGit

	curl -fsSL -o Git.exe "https://github.com/msysgit/msysgit/releases/download/Git-${msysGit}/Git-${msysGit}.exe"
)

(
	mkdir -p VirtualBox
	cd VirtualBox

	# http://www.virtualbox.org/manual/ch02.html
	curl -fsSL -o virtualbox.exe "http://download.virtualbox.org/virtualbox/${vbox}/VirtualBox-${vbox}-${vboxRev}-Win.exe"

	virtualbox.exe -extract -silent -path .
	rm virtualbox.exe # not neeeded after extraction

	rm *x86.msi # not needed as we only install 64-bit
	mv *_amd64.msi VirtualBox_amd64.msi # remove version number
)
