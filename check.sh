#!/bin/bash
set -e

# check if codesign and openssl are installed
which codesign > /dev/null 2>&1 || { echo "codesign is not installed"; exit 1; }
which openssl > /dev/null 2>&1 || { echo "openssl is not installed"; exit 1; }

# if no parameter is provided, use the default binary file
binary_file="$1"
if [ $# -ne 1 ]; then
    binary_file="/Applications/Docker.app/Contents/Library/LaunchServices/com.docker.vmnetd"
fi

# create a folder with temp folder to extract the certificates into
timestamp=$(date "+%Y%m%d_%H%M%S")
folder_name="${TMPDIR:-/tmp}/docker-desktop_cert_check_$timestamp"
mkdir -p "$folder_name"
cd "$folder_name"

# extract and check if the extraction was successful
if ! codesign -d --extract-certificates "$binary_file" >/dev/null 2>&1
then
    echo "Failed to extract certificates from $binary_file"
    exit 1
fi

# run the OpenSSL command and capture the output
certificate_details=$(openssl x509 -noout -serial -subject -issuer -dates -in codesign0)
binary_name=$(basename "$binary_file")

cd ..
rm -rf "$folder_name"

echo "-----------------------------------------------------------------"
echo "Certificate details for $binary_name:"
echo "$certificate_details" | tr '\n' ',' | tr -s ',' | tr ',' '\n' | while read -r line; do
    echo " $line"
done

echo "-----------------------------------------------------------------"
echo ""

# check for specific serial numbers
if [[ "$certificate_details" == *"serial=1316FD127D9A5715176591F85FFC3C66"* ]]; then
    echo "$binary_name is signed with a revoked certificate"
    echo "please download and install a new version of Docker Desktop"
    exit 1
elif [[ "$certificate_details" == *"serial=3EC22E699630083A"* ]]; then
    echo "$binary_name is signed with a correct certificate"
    exit 0
else
    echo "$binary_name is signed with an unknown certificate"
    echo "please download and install a new version of Docker Desktop"
    exit 1
fi