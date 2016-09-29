#!/bin/bash

OSX_VERS=$(sw_vers -productVersion | awk -F "." '{print $2}')

# Turn off hibernation and get rid of the sleepimage
pmset hibernatemode 0
rm -f /var/vm/sleepimage

# Stop the pager process and drop swap files. These will be re-created on boot.
# Starting with El Cap we can only stop the dynamic pager if SIP is disabled.
if [ "$OSX_VERS" -lt 11 ] || $(csrutil status | grep -q disabled); then
    launchctl unload /System/Library/LaunchDaemons/com.apple.dynamic_pager.plist
    sleep 5
fi
rm -rf /private/var/vm/swap*

# VMware Fusion specific items
if [ -e .vmfusion_version ] || [[ "$PACKER_BUILDER_TYPE" == vmware* ]]; then
    # Shrink the disk
    /Library/Application\ Support/VMware\ Tools/vmware-tools-cli disk shrink /
fi
