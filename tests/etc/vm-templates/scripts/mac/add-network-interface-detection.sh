#!/bin/sh

# This script adds a Mac OS Launch Daemon, which runs every time the
# machine is booted. The daemon will re-detect the attached network
# interfaces. If this is not done, network devices may not work.
PLIST=/Library/LaunchDaemons/com.github.timsutton.osx-vm-templates.detectnewhardware.plist
cat <<EOF > "${PLIST}"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.github.timsutton.osx-vm-templates.detectnewhardware</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/sbin/networksetup</string>
        <string>-detectnewhardware</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>
EOF

# These should be already set as follows, but since they're required
# in order to load properly, we set them explicitly.
/bin/chmod 644 "${PLIST}"
/usr/sbin/chown root:wheel "${PLIST}"
