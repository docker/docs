| README.md |
|:---|

# Desktop Editions CLI

## Overview
The idea of this CLI is to provide a simple tool that can integrate with our release process and generate the proper tools for desktop apps to self-update.

Most of the release info documented here is inspired by the Chrome best practices, as seen here:
https://support.google.com/chrome/a/answer/6025002?hl=en

## <a name="versions"></a>Versions
Versions show follow the following pattern: `major.minor[.build[.revision]]`
The above is inline with semver.org

## Commands available

This cli contains a few common release functions, including but not limited to:

* **list**: list current releases for the given channel and arch
* **upload**: uploads an artifact for the given channel, arch and version
* **publish**: publish an artifact for the given channel, arch and version


###### Credentials
For all commands the following AWS S3 flags can be provided:

* `--aws-access-key`: the aws s3 access key that has access to the bucket
* `--aws-secret-key`: the aws s3 secret key that has access to the bucket

Alternatively, the CLI will try to read from the regular AWS env variables, or your `.aws/credentials` file. A profile can be passed to the CLI when the credentials file is used.
* `--aws-profile`: the aws s3 profile to use within the aws credentials file (default: `release`)


### List
When getting a list of current releases, the following required options need to be provided:

* `--channel`: the channel from which the artifacts should be listed
* `--arch`: the architecture of the artifacts to list


Example:
```
docker-release --channel stable --arch mac list
```

### Upload
During an upload call, the following files should be provided:

* Artifact: The file that will be pushed down as the update (`Docker.dmg` or `InstallDocker.msi`)
* NOTES: This contains the latest changes to the app including bug fixes, new features, etc.

Required options provided during CLI call

 * `--version`: the version of the app being uploaded (follow version pattern)
 * `--channel`: the channel the artifact will live in
 * `--arch`: the architecture for the artifact

Example:
```
docker-release --channel stable --arch mac --build 1.11.1.6974 upload Docker.dmg NOTES
```

### Publish
Prior to calling publish all uploads are made without an actual update to the source file.
This allows a user to push their update without releasing it to the channel, but still allow for the file to be shared.

The publish creates the 'published' version of an artifact, which means updating the release file to let all listeners grab the latest version.

Required options provided during CLI call

 * `--version`: the version of the app being released (follow [version](#version) pattern)
 * `--channel`: the channel the artifact will live in
 * `--arch`: the architecture for the artifact
 * `--human`: the human readable version of the app (used in appcast)

 Example:
 ```
 docker-release --channel stable --arch mac --build 1.11.1.6974 --human 1.11.1-beta11 publish
 ```

## Release folder structure
All releases are organize by:

* Channel (Stable, Beta, Dev)
  * Architecture (Mac, Windows)
    * Update file (`release` contains the latest stable version details)
    * Versions folder
