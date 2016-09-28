from doit import get_var
import itertools
import subprocess
import os
import sys
import pathlib
import hashlib
import pdb

def stream_cmd(*argv, **kwArgs):
    print('$$ ' + ' '.join(argv[0]))
    kwArgs['stdout'] = subprocess.PIPE
    kwArgs['stderr'] = subprocess.STDOUT
    kwArgs['universal_newlines'] = True
    popen = subprocess.Popen(*argv, **kwArgs)
    stdout_lines = iter(popen.stdout.readline, "")
    for stdout_line in stdout_lines:
        yield stdout_line

    popen.stdout.close()
    return_code = popen.wait()
    if return_code != 0:
        print('Error running %s, exited with error code %d; was ran with args: %s' % (argv, return_code, kwArgs))

def run(*argv, **kwArgs):
    print('$ ' + ' '.join(argv[0]))
    if 'stream' in kwArgs.keys():
        if kwArgs['stream']:
            # Such hacks, sorry. it looks like doit messes with sys.stdout, so I can't use that
            kwArgs['stdout'] = 1
            kwArgs['stderr'] = 2
        del kwArgs['stream']
    try:
        return subprocess.run(*argv, **kwArgs)
    except subprocess.CalledProcessError as e:
        out = ''
        if e.stdout is not None:

            if type(e.stdout) is str:
                out = e.stdout
            else:
                out = e.stdout.decode('utf-8')
        err = ''
        if e.stderr is not None:
            if type(e.stderr) is str:
                err = e.stderr
            else:
                err = e.stderr.decode('utf-8')
        print('Error running %s \nstdout: %s\nstderr: %s' % (e.cmd, out, err))
        raise

def hash_files(*filepaths):
    hash_md5 = hashlib.md5()
    for filepath in filepaths:
        if os.path.isfile(filepath):
            with open(filepath, 'rb') as f:
                for chunk in iter(lambda: f.read(4096), b''):
                    hash_md5.update(chunk)
        else: # assume this is a string that's not supposed to be a file
            hash_md5.update(filepath)
    return hash_md5.hexdigest()[:20]

makefileDir = subprocess.check_output('git rev-parse --show-toplevel'.split()).decode('utf-8').strip()

def read_file(p):
    try:
        with open(p, 'r') as f:
            return f.read().strip()
    except FileNotFoundError:
        return None

def mkprinter(v):
    return lambda: print(v)

apiServerName = 'dtr-api'
nginxName = 'dtr-nginx'
rethinkName = 'dtr-rethink'
notaryServerName = 'dtr-notary-server'
notarySignerName = 'dtr-notary-signer'
jobrunnerName = 'dtr-jobrunner'
registryName = 'dtr-registry'
bootstrapName = 'dtr'
componentNameRefs = [apiServerName, nginxName, rethinkName, registryName,
                     bootstrapName, notaryServerName, notarySignerName, jobrunnerName]
etcdImageName = 'docker/dtr-etcd:v2.2.4'
thirdPartyImages = [etcdImageName]
moshpitName = 'dtr-moshpit'

clidocgenName = 'clidocgen'

constantsFile = makefileDir + '/generated_constants.go'
versionFile = makefileDir + '/version'
channelFile = makefileDir + '/releaseChannel'
binariesDir = 'binaries'
cmdDir = 'cmd'
version = read_file(versionFile)
preVersion = ''
if len(version.split('-')) > 1:
    preVersion = version.split('-')[1]
shortVersion = version.split('-')[0]

gitSha = run('git rev-parse HEAD'.split(), stdout=subprocess.PIPE,
                        universal_newlines=True, check=True).stdout[:12]
channel = read_file(channelFile)

if len(channel.split('/')) == 1:
    dockerHubOrg, suffix = 'docker', channel
if len(channel.split('/')) == 3:
    url, org, suffix = channel.split('/')
    dockerHubOrg = url + '/' + org
else:
    dockerHubOrg, suffix = channel.split('/')
if suffix == 'stable':
    if preVersion == '':
        print('Warning: Trying to build stable channel with a pre-release version!')
else:
    commitCount = int(run('git rev-list HEAD --count'.split(),
                                     stdout=subprocess.PIPE,
                                     universal_newlines=True, check=True) \
                      .stdout.strip())
    maybeDirtySuffix = ''
    if run('git diff --quiet'.split()).returncode != 0:
        maybeDirtySuffix = '-dirty'
    version = '%s-%06d_g%s%s' % (version, commitCount, gitSha[:7], maybeDirtySuffix)
    # add suffix to all non-stable builds' image names
    componentNameRefs = map(lambda ref: ref+'-'+suffix, componentNameRefs)
    # XXX: there's no easy way to modify these strings in bulk :/
    # we should stop messing with names anyway... when we do that we can remove
    # this garbage
    apiServerName += '-' + suffix
    nginxName += '-' + suffix
    rethinkName += '-' + suffix
    notaryServerName += '-' + suffix
    notarySignerName += '-' + suffix
    jobrunnerName += '-' + suffix
    registryName += '-' + suffix
    bootstrapName += '-' + suffix


goListContainerName = 'golistcontainer'
swagger1_2_to_2Image = 'dtr-swagger1.2to2:latest'
apidocgenImage = 'dtr-apidocgen:latest'
internalBaseTag = hash_files(makefileDir+'/shared/base/Dockerfile')
internalBaseImage = 'dockerhubenterprise/internal-base:'+internalBaseTag
goBaseTag = hash_files(makefileDir+'/shared/GobaseDockerfile')
goBaseImage = 'dockerhubenterprise/dtr-gobase:'+goBaseTag
ginkgoTag = hash_files(makefileDir+'/shared/GinkgoDockerfile')
ginkgoImage = 'dockerhubenterprise/dtr-ginkgo:'+ginkgoTag
uiBaseUntaggedImage = 'dockerhubenterprise/dtr-ui-build'
rethinkBaseUntaggedImage = 'dockerhubenterprise/dtr-rethink-base'
notaryServerBaseUntaggedImage = 'dockerhubenterprise/dtr-notary-server-base'
notarySignerBaseUntaggedImage = 'dockerhubenterprise/dtr-notary-signer-base'

goCachePath = makefileDir+'/go-build-cache'
testFlags = ''
bindAssets = False
noExec = False
runCoverage = False
serial = False
envs = ''
disallowBuildImgPush = False
tarUCPImages = ''
pullUCPImages = False
try:
    goCachePath = get_var('goCachePath', goCachePath)
    testFlags = get_var('testFlags', '')
    userId = get_var('userId', '0')
    groupId = get_var('groupId', '0')
    bindAssets = get_var('bindAssets', False)
    noExec = get_var('noExec', False)
    runCoverage = get_var('runCoverage', False)
    serial = get_var('serial', False)
    envs = get_var('envs', '')
    disallowBuildImgPush = get_var('disallowBuildImgPush', False)
    tarUCPImages = get_var('tarUCPImages', '')
    pullUCPImages = get_var('pullUCPImages', False)
except:
    pass

def get_images():
    imageStrings = subprocess.check_output('docker images -a'.split()) \
            .decode('utf-8') \
            .strip() \
            .split('\n')[1:]
    return set(img[0]+':'+img[1] for img in [imgStr.split() for imgStr in imageStrings])

imageSet = get_images()

def container_exists(cID, expectedImgName):
    p = run(['docker', 'inspect', '-f',
                        '{{.Config.Image}} {{.State.Running}}', cID],
                       stdout=subprocess.PIPE, universal_newlines=True)
    if p.returncode != 0:
        return 'none', False
    imgName, runningStr = p.stdout.split()
    runningDesc = 'running' if runningStr == 'true' else 'stopped'
    return runningDesc, imgName == expectedImgName

def docker_pull_or_build_and_maybe_push(img, dirPath, dockerfileName):
    if img in imageSet:
        return
    if run(['docker', 'pull', img]).returncode == 0:
        imageSet.add(img)
        return
    run(['docker', 'build', '-t', img, '-f',
                    dirPath+'/'+dockerfileName, dirPath], check=True)
    imageSet.add(img)
    if disallowBuildImgPush:
        print('Built %s. Please push it to docker hub with \'docker push %s\' if \
              you are happy with it.' % (img, img))
    elif run(['docker', 'push', img]).returncode != 0:
        print('Built %s and tried to push but failed' % img)

def fix_dirs(goCachePath, cacheSubdir):
    run(['mkdir', '-p', '%s/%s' % (goCachePath, cacheSubdir)])
    run(['chown', '-R', '%s:%s' % (userId, groupId), '%s/%s' % (goCachePath, cacheSubdir)])

# TODO use docker-py?
# We use different build cache subdirectories for different binaries because it
# looks like go install is not safe to run in parallel
def go_run(argv, flags=None, runArgs=None, cacheSubdir='', asAdmin=False):
    if flags == None:
        flags = []
    if cacheSubdir != '':
        fix_dirs(goCachePath, cacheSubdir)
    cmd = ''
    if not asAdmin:
        cmd = 'docker run -u %s:%s -v %s:/go/src/github.com/docker/dhe-deploy' % \
                (userId, groupId, makefileDir)
    else:
        cmd = 'docker run -v %s:/go/src/github.com/docker/dhe-deploy' % makefileDir
    if cacheSubdir != '':
        cmd += ' -v %s:/go/pkg' % (goCachePath+cacheSubdir)

    if '-d' not in flags:
        flags = flags + ['--rm']
    cmdArgv = cmd.split() + flags + [goBaseImage] + argv
    #print('!!! ' + ' '.join(cmdArgv))
    #print(cmdArgv)
    if runArgs == None:
        runArgs = {'check': True}
    #print('runArgs', runArgs)
    return run(cmdArgv, **runArgs)

# TODO: deduplicate with go_run
def ginkgo_go_run(argv, flags=None, runArgs=None, cacheSubdir='', asAdmin=False):
    if flags == None:
        flags = []
    if cacheSubdir != '':
        fix_dirs(goCachePath, cacheSubdir)
    cmd = ''
    if not asAdmin:
        cmd = 'docker run -u %s:%s -v %s:/go/src/github.com/docker/dhe-deploy' % \
            (userId, groupId, makefileDir)
    else:
        cmd = 'docker run -v %s:/go/src/github.com/docker/dhe-deploy' % makefileDir
    if cacheSubdir != '':
        cmd += ' -v %s:/go/pkg' % (goCachePath+cacheSubdir)
    if '-d' not in flags:
        flags = flags + ['--rm']
    cmdArgv = cmd.split() + flags + [ginkgoImage] + argv
    if runArgs == None:
        runArgs = {'check': True}
    return run(cmdArgv, **runArgs)

# TODO: maybe we can switch everything to use this instead of go_run?
def go_stream_print(argv, flags=None, runArgs={}, cacheSubdir='', asAdmin=False):
    for line in go_stream(argv, flags, runArgs, cacheSubdir, asAdmin):
        sys.stdout.write(line)

def go_stream(argv, flags=None, runArgs={}, cacheSubdir='', asAdmin=False):
    if flags == None:
        flags = []
    if cacheSubdir != '':
        fix_dirs(goCachePath, cacheSubdir)
    cmd = ''
    if not asAdmin:
        cmd = 'docker run -u %s:%s -v %s:/go/src/github.com/docker/dhe-deploy' % \
                (userId, groupId, makefileDir)
    else:
        cmd = 'docker run -v %s:/go/src/github.com/docker/dhe-deploy' % makefileDir
    if cacheSubdir != '':
        cmd += ' -v %s:/go/pkg' % (goCachePath+cacheSubdir)
    if '-d' not in flags:
        flags = flags + ['--rm']
    cmdArgv = cmd.split() + flags + [goBaseImage] + argv
    return stream_cmd(cmdArgv, **runArgs)

def go_exec(argv, runArgs={}):
    cmdArgv = ['docker', 'exec', goListContainerName] + argv
    return run(cmdArgv, **runArgs)

def kill_container(cID):
    try:
        run(['docker', 'rm', '-f', cID], check=True)
    except subprocess.CalledProcessError as e:
        pass

def go_deps(path):
    goRunShCmd = '''go list -f "{{if not .Standard}}{{.ImportPath}}{{end}}" $(go list -f "{{join .Deps \\" \\"}}" %s) %s''' % (path, path)
    p = None
    if not noExec:
        p = go_exec(['sh', '-c', goRunShCmd], runArgs={'stdout': subprocess.PIPE,
                                                      'universal_newlines': True,
                                                      'check': True})
    else:
        p = go_run(['sh', '-c', goRunShCmd], runArgs={'stdout': subprocess.PIPE,
                                                      'universal_newlines': True,
                                                      'check': True})
    if p.returncode != 0:
        raise Exception('compilation failed for %s' % path)
    goContainerPkgPaths = p.stdout.split()
    # strip the beginning of the paths
    pythonContainerPkgPaths = [pkgPath[len('github.com/docker/dhe-deploy/'):]
                               for pkgPath in goContainerPkgPaths]
    # list all go files in each package
    files = [str(goFilePath) for pkgPath in pythonContainerPkgPaths for goFilePath in
             pathlib.Path(pkgPath).glob('*.go')]
    #print("deps of", path, '=>', files)
    return files


def go_build(packageRelPath, outputRelPath, goArgs=None, dockerArgs=None):
    if goArgs == None:
        goArgs = []
    if dockerArgs == None:
        dockerArgs = []

    goInstallArgv = ['go', 'install'] + \
            goArgs + ['github.com/docker/dhe-deploy/%s' % packageRelPath]
    dockerFlags = dockerArgs + ['-v', makefileDir+'/'+outputRelPath+':/go/bin']

    # hack to make a unique cache directory for each binary we try to build
    # we don't use this on circle ci because we run without parallelism there
    # so we want to be able to reuse the cache between `go install` runs
    cacheDir = ''
    if not serial:
        hash_md5 = hashlib.md5()
        hash_md5.update(packageRelPath.encode("utf-8"))
        cacheDir = '/'+hash_md5.hexdigest()[:20]

    go_run(goInstallArgv, flags=dockerFlags, cacheSubdir=cacheDir)

# when literalName is set, we use imgName directly instead of constructing the image
# name using the dtr image name conventions
def build_from_base(imgName, copyTuples, entrypoint, baseImage=internalBaseImage, literalName=False):
    cID = run(['docker', 'run', '-d', baseImage],
                         stdout=subprocess.PIPE, universal_newlines=True,
                         check=True).stdout.strip()
    for fromPath, toPath in copyTuples:
        run(['docker', 'cp', fromPath, cID+':'+toPath], check=True)

    if literalName:
        run(['docker', 'commit', '-c', 'ENTRYPOINT [%s]' % entrypoint, cID,
                        imgName], check=True)
    else:
        builtImageTag = '%s/%s:%s' % (dockerHubOrg, imgName, version)
        builtImageTagLatest = '%s/%s:latest' % (dockerHubOrg, imgName)
        run(['docker', 'commit', '-c', 'ENTRYPOINT [%s]' % entrypoint, cID,
                        builtImageTag], check=True)
        run(['docker', 'tag', builtImageTag, builtImageTagLatest],
                      check=True)

# TODO: remove CurrentDirectory in prod builds; it's used only for BIND_ASSETS
generatedConstantsText = '''package deploy

const (
    Version = "%s"
    ShortVersion = "%s"
    GitSHA = "%s"
    CurrentDirectory = "%s"
    DefaultReleaseChannel = "%s"
    DockerHubOrg = "%s"
    EtcdImageName = "%s"
    BootstrapRepo repoName = "%s"
    APIServerRepo repoName = "%s"
    NginxRepo repoName = "%s"
    RethinkRepo repoName = "%s"
    RegistryRepo repoName = "%s"
    NotaryServerRepo repoName = "%s"
    NotarySignerRepo repoName = "%s"
    JobrunnerRepo repoName = "%s"
    BindAssets bool = %s
)''' % (version, shortVersion, gitSha, makefileDir, channel, dockerHubOrg,
        etcdImageName, bootstrapName, apiServerName, nginxName,
        rethinkName, registryName, notaryServerName, notarySignerName, jobrunnerName, 'true' if bindAssets else 'false')
