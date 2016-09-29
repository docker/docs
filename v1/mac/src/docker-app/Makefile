targetName=Docker
sourceDirectory=docker/docker
xcodeProject=docker/docker.xcodeproj
buildDirectory=build
makefilePath := $(abspath $(lastword $(MAKEFILE_LIST)))
currentDirectory := $(dir $(makefilePath))

all: build-xcode

checkenv:
	@echo "Looking for xcode ..."

build-xcode: ${sourceDirectory}/commit.swift docker/Carthage/Build/Mac/HockeySDK.framework docker/Carthage/Build/Mac/Sparkle.framework
	@echo "Building ${targetName}"
	xcodebuild -list -project ./docker/docker.xcodeproj
	xcodebuild -verbose build -project ./docker/docker.xcodeproj -configuration Release -scheme ${targetName} -derivedDataPath ${buildDirectory} -parallelizeTargets

.PHONY: CODE_SIGN_IDENTITY
CODE_SIGN_IDENTITY:
	if [ ! -e CODE_SIGN_IDENTITY ] || [ "$(shell cat CODE_SIGN_IDENTITY)" != "${CODE_SIGN_IDENTITY}" ]; then \
		echo "${CODE_SIGN_IDENTITY}" > CODE_SIGN_IDENTITY;                \
	fi

# docker-installer/Info.plist: docker-installer/Info.plist.in CODE_SIGN_IDENTITY
# 	sed 's|@SIGNING_ID@|$(shell cat CODE_SIGN_IDENTITY)|g' docker-installer/Info.plist.in | \
# 	sed 's|@DOCKER_VMNETD@|${DOCKER_VMNETD}|g' > \
# 	docker-installer/Info.plist

.PHONY: ${sourceDirectory}/commit.swift
${sourceDirectory}/commit.swift:
	sed "s|@SHA@|$(shell git show --pretty=%H --no-patch)|g" ${sourceDirectory}/commit.swift.in > ${sourceDirectory}/commit.swift


.PHONY:clean
clean:
	@echo "Cleaning ${targetName}"
	@xcodebuild -verbose clean -project ${xcodeProject} -target ${targetName} CODE_SIGN_IDENTITY="" CODE_SIGNING_REQUIRED=NO
	rm -rf build
	rm -f ${sourceDirectory}/commit.swift

re: clean all

dev:
	@echo "Starting dev for ${targetName}"
	../../scripts/openInXCode ${currentDirectory}/docker/${targetName}.xcodeproj ${targetName} ${targetName}

docker/Carthage/Build/Mac/HockeySDK.framework docker/Carthage/Build/Mac/Sparkle.framework:
	cd docker && carthage build --platform mac

depends: docker/Carthage/Build/Mac/HockeySDK.framework docker/Carthage/Build/Mac/Sparkle.framework
	@echo "Setting up dependencies"
	cd docker && carthage build --platform mac

cleandepends:
	@echo "Cleaning up dependencies"
