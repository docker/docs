#!/bin/sh

TOOLS_PATH="/Users/$USERNAME/darwin.iso"
# VMware Fusion specific items
if [ -e .vmfusion_version ] || [[ "$PACKER_BUILDER_TYPE" == vmware* ]]; then
    if [ ! -e "$TOOLS_PATH" ]; then
        echo "Couldn't locate uploaded tools iso at $TOOLS_PATH!"
        exit 1
    fi

    TMPMOUNT=`/usr/bin/mktemp -d /tmp/vmware-tools.XXXX`
    hdiutil attach "$TOOLS_PATH" -mountpoint "$TMPMOUNT"

    INSTALLER_PKG="$TMPMOUNT/Install VMware Tools.app/Contents/Resources/VMware Tools.pkg"
    if [ ! -e "$INSTALLER_PKG" ]; then
        echo "Couldn't locate VMware installer pkg at $INSTALLER_PKG!"
        exit 1
    fi

    echo "Installing VMware tools.."
    installer -pkg "$TMPMOUNT/Install VMware Tools.app/Contents/Resources/VMware Tools.pkg" -target /

    # This usually fails
    hdiutil detach "$TMPMOUNT"
    rm -rf "$TMPMOUNT"
    rm -f "$TOOLS_PATH"

    # Point Linux shared folder root to that used by OS X guests,
    # useful for the Hashicorp vmware_fusion Vagrant provider plugin
    mkdir /mnt
    ln -sf /Volumes/VMware\ Shared\ Folders /mnt/hgfs
fi
