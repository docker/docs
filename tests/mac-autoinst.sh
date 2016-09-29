#!/bin/bash
exec 2>&1 > /tmp/install.log

# Disable spotlight
sudo launchctl unload -w /System/Library/LaunchDaemons/com.apple.metadata.mds.plist

REG_PROXY=${1}
FILE=/Users/docker/pinata/docker.dmg
MOUNTPOINT=/tmp/docker-latest
BETA_TOKEN=ughetJWhz5aFzf5dgc8qu24Tp

echo "### Verifying DMG..."
hdiutil verify ${FILE} || exit

if [ -e "/Applications/Docker.app" ]; then
    echo "### Existing Docker.app found"
    echo "   - Uninstalling..."
    open /Applications/Docker.app --args --uninstall || exit
    echo "   - Moving to trash..."
    # see http://apple.stackexchange.com/questions/50844/how-to-move-files-to-trash-from-command-line for more options
    osascript -e 'tell app "Finder" to move the POSIX file "/Applications/Docker.app" to trash' || exit
fi

mkdir ${MOUNTPOINT}
echo "### Mounting DMG..."
hdiutil attach ${FILE} -mountpoint ${MOUNTPOINT} -nobrowse || exit

echo "### Copying to /Applications/Docker.app..."
mkdir -p /Users/docker/v1/mac/build/
cp -a ${MOUNTPOINT}/Docker.app /Users/docker/v1/mac/build/ || exit

echo "### Unmounting DMG..."
hdiutil detach ${MOUNTPOINT} || exit
rmdir ${MOUNTPOINT} || exit
rm ${FILE} || exit
