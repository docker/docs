#!/usr/bin/env bash
set -e

# This script looks for bundles built by make.sh, and releases them on a
# public S3 bucket.
#
# Bundles should be available for the VERSION string passed as argument.
#
# The correct way to call this script is inside a container built by the
# official Dockerfile at the root of the Docker source code. The Dockerfile,
# make.sh and release.sh should all be from the same source code revision.

set -o pipefail

# Print a usage message and exit.
usage() {
	cat >&2 <<'EOF'
To run, I need:
- to be in a container generated by the Dockerfile at the top of the Docker
  repository;
- to be provided with the name of an S3 bucket, in environment variable
  AWS_S3_BUCKET;
- to be provided with AWS credentials for this S3 bucket, in environment
  variables AWS_ACCESS_KEY and AWS_SECRET_KEY;
- the passphrase to unlock the GPG key which will sign the deb packages
  (passed as environment variable GPG_PASSPHRASE);
- a generous amount of good will and nice manners.
The canonical way to run me is to run the image produced by the Dockerfile: e.g.:"

docker run -e AWS_S3_BUCKET=get-staging.docker.io \
           -e AWS_ACCESS_KEY=AKI1234... \
           -e AWS_SECRET_KEY=sEs4mE... \
           -e GPG_PASSPHRASE=m0resEs4mE... \
           -i -t --privileged \
           docker ./hack/release.sh
EOF
	exit 1
}

[ "$AWS_S3_BUCKET" ] || usage
[ "$AWS_ACCESS_KEY" ] || usage
[ "$AWS_SECRET_KEY" ] || usage
[ "$GPG_PASSPHRASE" ] || usage
[ -d /go/src/github.com/dotcloud/docker ] || usage
cd /go/src/github.com/dotcloud/docker
[ -x hack/make.sh ] || usage

RELEASE_BUNDLES=(
	binary
	cross
	tgz
	ubuntu
)

if [ "$1" != '--release-regardless-of-test-failure' ]; then
	RELEASE_BUNDLES=( test "${RELEASE_BUNDLES[@]}" )
fi
	
VERSION=$(cat VERSION)
BUCKET=$AWS_S3_BUCKET

# These are the 2 keys we've used to sign the deb's
#   release (get.docker.io)
#	GPG_KEY="36A1D7869245C8950F966E92D8576A8BA88D21E9"
#   test    (test.docker.io)
#	GPG_KEY="740B314AE3941731B942C66ADF4FD13717AAD7D6"

setup_s3() {
	# Try creating the bucket. Ignore errors (it might already exist).
	s3cmd mb s3://$BUCKET 2>/dev/null || true
	# Check access to the bucket.
	# s3cmd has no useful exit status, so we cannot check that.
	# Instead, we check if it outputs anything on standard output.
	# (When there are problems, it uses standard error instead.)
	s3cmd info s3://$BUCKET | grep -q .
	# Make the bucket accessible through website endpoints.
	s3cmd ws-create --ws-index index --ws-error error s3://$BUCKET
}

# write_to_s3 uploads the contents of standard input to the specified S3 url.
write_to_s3() {
	DEST=$1
	F=`mktemp`
	cat > $F
	s3cmd --acl-public --mime-type='text/plain' put $F $DEST
	rm -f $F
}

s3_url() {
	case "$BUCKET" in
		get.docker.io|test.docker.io)
			echo "https://$BUCKET"
			;;
		*)
			s3cmd ws-info s3://$BUCKET | awk -v 'FS=: +' '/http:\/\/'$BUCKET'/ { gsub(/\/+$/, "", $2); print $2 }'
			;;
	esac
}

build_all() {
	if ! ./hack/make.sh "${RELEASE_BUNDLES[@]}"; then
		echo >&2
		echo >&2 'The build or tests appear to have failed.'
		echo >&2
		echo >&2 'You, as the release  maintainer, now have a couple options:'
		echo >&2 '- delay release and fix issues'
		echo >&2 '- delay release and fix issues'
		echo >&2 '- did we mention how important this is?  issues need fixing :)'
		echo >&2
		echo >&2 'As a final LAST RESORT, you (because only you, the release maintainer,'
		echo >&2 ' really knows all the hairy problems at hand with the current release'
		echo >&2 ' issues) may bypass this checking by running this script again with the'
		echo >&2 ' single argument of "--release-regardless-of-test-failure", which will skip'
		echo >&2 ' running the test suite, and will only build the binaries and packages.  Please'
		echo >&2 ' avoid using this if at all possible.'
		echo >&2
		echo >&2 'Regardless, we cannot stress enough the scarcity with which this bypass'
		echo >&2 ' should be used.  If there are release issues, we should always err on the'
		echo >&2 ' side of caution.'
		echo >&2
		exit 1
	fi
}

release_build() {
	GOOS=$1
	GOARCH=$2

	SOURCE_DIR=bundles/$VERSION/cross/$GOOS/$GOARCH
	BINARY=docker-$VERSION
	BINARY_MD5=docker-$VERSION.md5
	BINARY_SHA256=docker-$VERSION.sha256
	TGZ=docker-$VERSION.tgz

	# we need to map our GOOS and GOARCH to uname values
	# see https://en.wikipedia.org/wiki/Uname
	# ie, GOOS=linux -> "uname -s"=Linux

	S3OS=$GOOS
	case "$S3OS" in
		darwin)
			S3OS=Darwin
			;;
		freebsd)
			S3OS=FreeBSD
			;;
		linux)
			S3OS=Linux
			;;
		*)
			echo >&2 "error: can't convert $S3OS to an appropriate value for 'uname -s'"
			exit 1
			;;
	esac

	S3ARCH=$GOARCH
	case "$S3ARCH" in
		amd64)
			S3ARCH=x86_64
			;;
		386)
			S3ARCH=i386
			;;
		arm)
			S3ARCH=armel
			# someday, we might potentially support mutliple GOARM values, in which case we might get armhf here too
			;;
		*)
			echo >&2 "error: can't convert $S3ARCH to an appropriate value for 'uname -m'"
			exit 1
			;;
	esac

	S3DIR=s3://$BUCKET/builds/$S3OS/$S3ARCH

	if [ ! -x "$SOURCE_DIR/$BINARY" ]; then
		echo >&2 "error: can't find $SOURCE_DIR/$BINARY - was it compiled properly?"
		exit 1
	fi
	if [ ! -f "$TGZ" ]; then
		echo >&2 "error: can't find $TGZ - was it packaged properly?"
		exit 1
	fi

	echo "Uploading $BINARY to $S3OS/$S3ARCH/docker-$VERSION"
	s3cmd --follow-symlinks --preserve --acl-public put $SOURCE_DIR/$BINARY $S3DIR/$BINARY

	echo "Uploading $BINARY_MD5 to $S3OS/$S3ARCH/docker-$VERSION.md5"
	s3cmd --follow-symlinks --preserve --acl-public put $SOURCE_DIR/$BINARY_MD5 $S3DIR/$BINARY_MD5

	echo "Uploading $BINARY_SHA256 to $S3OS/$S3ARCH/docker-$VERSION.sha256"
	s3cmd --follow-symlinks --preserve --acl-public put $SOURCE_DIR/$BINARY_SHA256 $S3DIR/$BINARY_SHA256

	echo "Uploading $TGZ to $S3OS/$S3ARCH/docker-$VERSION.tgz"
	s3cmd --follow-symlinks --preserve --acl-public put $SOURCE_DIR/$TGZ $S3DIR/$TGZ

	if [ -z "$NOLATEST" ]; then
		echo "Copying $S3DIR/$BINARY to $S3DIR/docker-latest"
		s3cmd --acl-public cp $S3DIR/$BINARY $S3DIR/docker-latest

		echo "Copying $S3DIR/$BINARY_MD5 to $S3DIR/docker-latest.md5"
		s3cmd --acl-public cp $S3DIR/$BINARY_MD5 $S3DIR/docker-latest.md5

		echo "Copying $S3DIR/$BINARY_SHA256 to $S3DIR/docker-latest.sha256"
		s3cmd --acl-public cp $S3DIR/$BINARY_SHA256 $S3DIR/docker-latest.sha256

		echo "Copying $S3DIR/$TGZ $S3DIR/docker-latest.tgz"
		s3cmd --acl-public cp $S3DIR/$TGZ $S3DIR/docker-latest.tgz
	fi
}

# Upload the 'ubuntu' bundle to S3:
# 1. A full APT repository is published at $BUCKET/ubuntu/
# 2. Instructions for using the APT repository are uploaded at $BUCKET/ubuntu/index
release_ubuntu() {
	[ -e bundles/$VERSION/ubuntu ] || {
		echo >&2 './hack/make.sh must be run before release_ubuntu'
		exit 1
	}
	
	s3cmd sync s3://$BUCKET/ubuntu/.gnupg/ /.gnupg/ || true

	# Sign our packages
	dpkg-sig -g "--passphrase $GPG_PASSPHRASE" -k releasedocker \
		--sign builder bundles/$VERSION/ubuntu/*.deb

	# Setup the APT repo
	APTDIR=bundles/$VERSION/ubuntu/apt
	mkdir -p $APTDIR/conf $APTDIR/db
	s3cmd sync s3://$BUCKET/ubuntu/db/ $APTDIR/db/ || true
	cat > $APTDIR/conf/distributions <<EOF
Codename: docker
Components: main
Architectures: amd64 i386
EOF

	# Add the DEB package to the APT repo
	DEBFILE=bundles/$VERSION/ubuntu/lxc-docker*.deb
	reprepro -b $APTDIR includedeb docker $DEBFILE

	# Sign
	for F in $(find $APTDIR -name Release); do
		gpg -u releasedocker --passphrase $GPG_PASSPHRASE \
			--armor --sign --detach-sign \
			--output $F.gpg $F
	done

	# Upload keys
	s3cmd sync /.gnupg/ s3://$BUCKET/ubuntu/.gnupg/
	gpg --armor --export releasedocker > bundles/$VERSION/ubuntu/gpg
	s3cmd --acl-public put bundles/$VERSION/ubuntu/gpg s3://$BUCKET/gpg

	# Upload repo
	s3cmd --acl-public sync $APTDIR/ s3://$BUCKET/ubuntu/
	cat <<EOF | write_to_s3 s3://$BUCKET/ubuntu/index
# Add the repository to your APT sources
echo deb $(s3_url)/ubuntu docker main > /etc/apt/sources.list.d/docker.list
# Then import the repository key
apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9
# Install docker
apt-get update ; apt-get install -y lxc-docker

#
# Alternatively, just use the curl-able install.sh script provided at $(s3_url)
#
EOF

	# Add redirect at /ubuntu/info for URL-backwards-compatibility
	rm -rf /tmp/emptyfile && touch /tmp/emptyfile
	s3cmd --acl-public --add-header='x-amz-website-redirect-location:/ubuntu/' --mime-type='text/plain' put /tmp/emptyfile s3://$BUCKET/ubuntu/info

	echo "APT repository uploaded. Instructions available at $(s3_url)/ubuntu"
}

# Upload binaries and tgz files to S3
release_binaries() {
	[ -e bundles/$VERSION/cross/linux/amd64/docker-$VERSION ] || {
		echo >&2 './hack/make.sh must be run before release_binaries'
		exit 1
	}

	for d in bundles/$VERSION/cross/*/*; do
		GOARCH="$(basename "$d")"
		GOOS="$(basename "$(dirname "$d")")"
		release_build "$GOOS" "$GOARCH"
	done

	# TODO create redirect from builds/*/i686 to builds/*/i386

	cat <<EOF | write_to_s3 s3://$BUCKET/builds/index
# To install, run the following command as root:
curl -O $(s3_url)/builds/Linux/x86_64/docker-$VERSION && chmod +x docker-$VERSION && sudo mv docker-$VERSION /usr/local/bin/docker
# Then start docker in daemon mode:
sudo /usr/local/bin/docker -d
EOF

	# Add redirect at /builds/info for URL-backwards-compatibility
	rm -rf /tmp/emptyfile && touch /tmp/emptyfile
	s3cmd --acl-public --add-header='x-amz-website-redirect-location:/builds/' --mime-type='text/plain' put /tmp/emptyfile s3://$BUCKET/builds/info

	if [ -z "$NOLATEST" ]; then
		echo "Advertising $VERSION on $BUCKET as most recent version"
		echo $VERSION | write_to_s3 s3://$BUCKET/latest
	fi
}

# Upload the index script
release_index() {
	sed "s,url='https://get.docker.io/',url='$(s3_url)/'," hack/install.sh | write_to_s3 s3://$BUCKET/index
}

release_test() {
	if [ -e "bundles/$VERSION/test" ]; then
		s3cmd --acl-public sync bundles/$VERSION/test/ s3://$BUCKET/test/
	fi
}

setup_gpg() {
	# Make sure that we have our keys
	mkdir -p /.gnupg/
	gpg --list-keys releasedocker >/dev/null || {
		gpg --gen-key --batch <<EOF
Key-Type: RSA
Key-Length: 2048
Passphrase: $GPG_PASSPHRASE
Name-Real: Docker Release Tool
Name-Email: docker@dotcloud.com
Name-Comment: releasedocker
Expire-Date: 0
%commit
EOF
	}
}

main() {
	build_all
	setup_s3
	setup_gpg
	release_binaries
	release_ubuntu
	release_index
	release_test
}

main

echo
echo
echo "Release complete; see $(s3_url)"
echo
