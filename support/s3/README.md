Log formatter from S3
=====================

- `sync.sh` will populate the `data/` directory with the latest support logs
- `tohtml.sh` will process them into HTML in the `html/` directory.

TODO:

- classify log files so they can be processed with different containers (E.g. fuse)
- better index page
- add file sizes
- add infra for grep scripts to do autotriage
- plumb through a docker push to our private repo with the results


## Requesting Access To AWS

To get the necessary account, please ask @dave-tucker, @balrajs, @ijc25 or anyone else who has admin access on AWS...

- Your account needs to be created
- You need to be added to the `pinata` group.
- The `docker-pinata-support` bucket policy needs to have your username added to the whitelist

## Using AWS

You will be given a CSV with your username, password and login URL
Please login and change your password. You may also take this opportunity to set up 2FA.

Once logged in, follow [these instructions](http://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html#Using_CreateAccessKey) to create an Access Key. This will let you use the AWS API and CLI tools.

To create the `~/.aws/credentials` file, your best bet is to first install the AWS CLI.

```
# brew install pip if you don't have it already.
pip install awscli
```

Then you can `aws configure` and follow the instructions.

If you have more than one set of AWS credentials, you may want to consider creating a different "profile" for each one. In this case:

```
aws configure --profile pinata
```

In order for your file to be bind-mountable as a `~/.aws/credentials` file, you may want to maintain a copy at `~/.aws/pinata-credentials` with your Pinata credentials as the `default` profile.

# Fetch

Fetch is a tool that fetches the logs for a provided platform/uuid and optionally a report
It connects to S3 and outputs the logs to the mounted `/logs` directory

Example Usage:

```bash
$ docker build -t pinata:fetch -f Dockerfile.fetch .

$ docker run -it --rm \
	-v /my/logs/dir:/logs  \
	-v /home/me/.aws/credentials:/root/.aws/credentials \
	pinata:fetch $platform $uuid $report
```

