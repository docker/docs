#!/bin/sh
date > /etc/box_build_time
OSX_VERS=$(sw_vers -productVersion | awk -F "." '{print $2}')

# Set computer/hostname
COMPNAME=osx-10_${OSX_VERS}
scutil --set ComputerName ${COMPNAME}
scutil --set HostName ${COMPNAME}.docker.com

# Create a group and assign the user to it
dseditgroup -o create "$USERNAME"
dseditgroup -o edit -a "$USERNAME" "$USERNAME"

echo "==> Writing SSH keys"
# Add SSH Keys
mkdir "/Users/$USERNAME/.ssh"
chmod 700 "/Users/$USERNAME/.ssh"
echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDH5NTlIlJPBkZrP/dfSoi5wWni+WbtXi3najQ8L26OoLIn0P6DooTZA9htPj/AMLnfXCFVR+I2kEz3Xq539ZA/R1Vv8n3N2mIjGMDowCP6wqb7mWMXoa+Fchz7HS6tiHfjhJGhIj/jrdod0AiKlINoyqzB+BsVBwy50bbjBwopk7236Qg5Wt2v7ZBND+wauREAmsCU4H7p0IwlqW/X2GxigbKGJrPeDrHP+4ADMoGGNnlU7rTXXC7Wm4VYaeDOEAKKik2UJr0W6xVy5ec32pAoD4zmwXa42UQMdDhPX1i61cK54theieLk3GOF2pi98ALdFhtUMnnygoEYILj1J5t3 docker@rt-host.local" > "/Users/$USERNAME/.ssh/authorized_keys"
chmod 600 "/Users/$USERNAME/.ssh/authorized_keys"
chown -R "$USERNAME" "/Users/$USERNAME/.ssh"

echo "==> Setting Power Saving Settings"
# Set up power saving
pmset sleep 0
pmset displaysleep 0
pmset disksleep 0
