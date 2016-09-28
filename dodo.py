import project.util as util
import pdb
import os
import subprocess
import itertools
import pathlib
import sys
from doit.tools import run_once, result_dep
from doit.action import CmdAction

# DO NOT SET THIS otherwise `doit forget` won't forget everything
#DOIT_CONFIG = {'default_tasks': ['']}

def task_update_submodules():
    return {'actions': ['git submodule update --init --recursive'],
            'uptodate': [run_once]}

# TODO: get rid of generated constants and wire up direct dependencies on git hash instead of
# everything depending on 'file_dep': util.constantsFile,
# Also, we should make sure it reads the git hash only once, so if you commit during a build it's okay
def task_update_constants():
    def update_constants():
        if util.read_file(util.constantsFile) != util.generatedConstantsText:
            with open(util.constantsFile, 'w') as f:
                f.write(util.generatedConstantsText)
            os.chown(util.constantsFile, int(util.userId), int(util.groupId))
    return {'actions': [update_constants],
            'targets': [util.constantsFile]}

def task_print():
    for k, v in [('DOCKER_HUB_ORG', util.dockerHubOrg),
                 ('BOOTSTRAP_NAME', util.bootstrapName),
                 ('VERSION', util.version)]:
        yield {'basename': 'print-'+k,
                'actions': [util.mkprinter('OUTPUT: ' + k + '=' + v)]}

def task_print_components():
    def print_components():
        print('OUTPUT:', ' '.join(map(lambda x: util.dockerHubOrg + '/' + x, util.componentNameRefs)))
    return {'actions': [print_components]}

def task_base_img():
    def build_image():
        util.docker_pull_or_build_and_maybe_push(
            util.internalBaseImage, 'shared/base', 'Dockerfile')
        return util.internalBaseImage
    return {'actions': [build_image],
            'file_dep': ['shared/base/Dockerfile']}

def task_gofmt():
    def gofmt():
        util.go_run(['sh', '-c', 'test -z "$(gofmt -s -l .| grep -v \'vendor/\|generated_constants.go\' | tee /dev/stderr)"'])
    return {'actions': [gofmt],
            'file_dep': [util.constantsFile]}

def task_apidocgen_swagger_1_12_to_2_img():
    deps = ['apidocgen/swagger1.2to2/Dockerfile',
            'apidocgen/swagger1.2to2/main.js',
            'apidocgen/swagger1.2to2/package.json']
    return {'actions': ['docker build -t ' + util.swagger1_2_to_2Image + ' apidocgen/swagger1.2to2'],
            'file_dep': deps}

def task_apidocgen_bin():
    def build_bin():
        util.go_build('./apidocgen/apidocgen', util.binariesDir)
    return {'actions': [build_bin],
            'file_dep': [util.constantsFile],
            'targets': [util.binariesDir + '/apidocgen'],
            'uptodate': [result_dep('go_base_img')],
            'calc_dep': ['getapidocgendep']}

def task_apidocgen_img():
    def build_image():
        util.build_from_base(util.apidocgenImage,
                        [(util.binariesDir + '/apidocgen', '/bin/apidocgen'),
                         ('apidocgen/api_intro.md', '/swagger/')],
                        '"/bin/apidocgen"', literalName=True)
    return {'actions': [build_image],
            'file_dep': ['apidocgen/api_intro.md', util.binariesDir + '/apidocgen'],
            'uptodate': [result_dep('base_img')]}

def task_apidocgen_docs_tar():
    return {'actions': ['mkdir -p gen && docker run --rm -v ' + util.makefileDir + '/apidocgen/gen:/out ' + util.apidocgenImage],
            'task_dep': ['apidocgen_img'],
            'targets': ['apidocgen/gen/docs.tar']}

def task_apidocgen_swagger2_json():
    return {'actions': ['cd apidocgen/gen && tar -xvf docs.tar && docker run --rm -v ' + util.makefileDir + '/apidocgen/gen:/app/gen ' + util.swagger1_2_to_2Image + ' > swagger2.json'],
            'task_dep': ['apidocgen_swagger_1_12_to_2_img'],
            'file_dep': ['apidocgen/gen/docs.tar'],
            'targets': ['apidocgen/gen/swagger2.json']}

def task_apidocgen_output_template():
    fileDeps = list(filter(os.path.isfile, map(str, itertools.chain(
        pathlib.Path('apidocgen').glob('template/**/*'),
    ))))
    return {'actions': ['mkdir -p apidocgen/output/v' + util.shortVersion + '&& cp -r apidocgen/template/* apidocgen/output/v' + util.shortVersion],
            'file_dep': fileDeps,
            'targets': ['apidocgen/output/v' + util.shortVersion + '/index.html']} # we use an example file because we can't use a directory :(

def task_apidocgen_main_js():
    return {'actions': ['cd apidocgen && sed -f sed.txt < template/main.js > output/v' + util.shortVersion + '/main.js'],
            'file_dep': [
                'apidocgen/gen/swagger2.json',
                'apidocgen/template/main.js',
                'apidocgen/sed.txt',
                'apidocgen/output/v' + util.shortVersion + '/index.html', # we use an example file because we can't use a directory :(
            ],
            'targets': ['apidocgen/output/v' + util.shortVersion + '/main.js']}

# TODO: this needs way more constants
def task_apidocgen_index_md():
    return {'actions': ['sed \'s/DTRVERSION/' + util.shortVersion + '/\' < apidocgen/template/index.md > apidocgen/output/v' + util.shortVersion + '/index.md'],
            'file_dep': [
                'apidocgen/template/index.md',
                'apidocgen/output/v' + util.shortVersion + '/index.html', # we use an example file because we can't use a directory :(
            ],
            'targets': ['apidocgen/output/v' + util.shortVersion + '/index.md']}

def task_apidocgen_swagger_ui_js():
    target = 'apidocgen/output/v' + util.shortVersion + '/swagger-ui.js'
    return {'actions': ['cp apidocgen/output/v' + util.shortVersion + '/swagger-ui.js.original ' + target + ' && patch ' + target + ' < apidocgen/template/swagger-fix.diff'],
            'file_dep': [
                'apidocgen/template/swagger-fix.diff',
                'apidocgen/output/v' + util.shortVersion + '/index.html', # we use an example file because we can't use a directory :(
            ],
            'targets': [target]}

def task_apidocgen_swagger_ui_min_js():
    target = 'apidocgen/output/v' + util.shortVersion + '/swagger-ui.min.js'
    return {'actions': ['docker run --rm -v ' + util.makefileDir + '/apidocgen:/work minty/uglifyjs /work/output/v' + util.shortVersion + '/swagger-ui.js > ' + target],
            'file_dep': ['apidocgen/output/v' + util.shortVersion + '/swagger-ui.js'],
            'targets': [target]}


def task_api_docs_gen():
    return {'actions': [],
            'file_dep': ['apidocgen/output/v' + util.shortVersion + '/main.js',
                         'apidocgen/output/v' + util.shortVersion + '/index.md',
                         'apidocgen/output/v' + util.shortVersion + '/swagger-ui.min.js'],
            'task_dep': ['apidocgen_output_template']}

def task_build_integration():
    def build_integration():
        util.go_run(['sh', '-c', '''go install ./integration/...;
for p in $(go list -f "{{if lt 0 (len .TestGoFiles) }}{{.ImportPath}}{{end}}" $(go list ./integration/...));
    do go test -c $p;
done
'''])
    return {'actions': [build_integration],
            'file_dep': [util.constantsFile],
            'targets': ['integration.test']}

def task_test_integration():
    def test_integration():
        # TODO: handle conifg in python in some more reasonable way
        var_list = {}
        with open('integration/.env', 'r') as f:
            for line in f:
                if line[0] != '#':
                    if len(line.split('=')) == 2:
                        var, val = line.split('=')
                        var_list[var] = val
        # overwrite from the env only envs that have default values
        with open('integration/local.env', 'r') as f:
            for line in f:
                if line[0] != '#':
                    if len(line.split('=')) == 2:
                        var, val = line.split('=')
                        if var in var_list.keys():
                            var_list[var] = val
        util.run(['mkdir', '-p', 'integration/results'])
        with open('integration/results/.env', 'w') as f:
            for k, v in var_list.items():
                f.write(k + '=' + v)

        # TODO: maybe we can just let the shell resolve those env vars in this case
        p = util.run(['docker', 'run',
            '-d',
            '--privileged',
            '--env-file', 'integration/results/.env',
            '--entrypoint', 'sh',
            'docker:1.10-dind',
            '-c', 'dockerd-entrypoint.sh --debug=true --insecure-registry ' +
            var_list['DTR_HOST'] + ' --insecure-registry ' + var_list['DTR_HOST'] + ':1337'], stdout=subprocess.PIPE)
        DIND_HOST = p.stdout.strip().decode('utf-8').strip()

        out_file = open('integration/results/test-results.txt', 'w')
        try:
            # The docker socket is mounted in and the DIND host is passed as an environment variable because
            # of the tests that need to connect to the socket in order to docker-exec into the DIND host
            # using the engine-api, since this container does not have a docker client binary.
            # Future alternatives:
            #   1. the test can just run a `docker:1.xx` container within the DIND container (the image names should match)
            #   1. the integration test container can just copy the docker client binary from the DIND container
            for line in util.go_stream(['go', 'test',
                        '-timeout', '1h',
                        '-p', '1',
                        '-v'
                    ] + util.testFlags.split(' '), flags=[
                        '-v', '/var/run/docker.sock:/var/run/docker.sock',
                        '--link', DIND_HOST + ':docker_daemon',
                        '-e', 'DIND_HOST=' + DIND_HOST,
                        '-w', '/go/src/github.com/docker/dhe-deploy/integration',
                        '--env-file', util.makefileDir + '/integration/results/.env'],
                    asAdmin=True):
                # bypass buffering by doit
                sys.stdout.write(line)
                out_file.write(line)
        finally:
            out_file.close()
            util.run(['docker', 'rm', '-f', DIND_HOST])
            util.go_stream_print(['go2xunit',
                    '-fail',
                    '-input', 'integration/results/test-results.txt',
                    '-output', 'integration/results/test-results.xml'
                ], flags=[
                    '-v', util.makefileDir +
                    '/integration/results:/results'
                ], asAdmin=True)
            util.run(['chmod', '777', '-R', 'integration/results'])
    return {'actions': [test_integration],
            'file_dep': [util.constantsFile],
            'uptodate': [False]} # always run

def task_gen_get_depss():
    for name, package in [('getbootstrapdep', './bootstrap/bootstrap'),
                          ('getapidocgendep', './apidocgen/apidocgen'),
                          ('getregistrydep', './registry/registry'),
                          ('getclidocgendep', './clidocgen'),
                          ('getarggendep', './rethinkdb/arggen'),
                          ('getmoshpitdep', './moshpit/cmd/moshpit'),
                          ('getdnsfixdep', './nginx/dnsfix'),
                          ('getjobrunnerdep', './jobrunner/jobrunner'),
                          ('getjobrunnergcdep', './cmd/gc'),
                          ('getapidep', './adminserver/server'),
                          ('gettagmigrationdep', './' + util.cmdDir + '/tagmigration'),
                          ('getgcdep', './' + util.cmdDir + '/gc')]:
        def getgetdeps(pkg):
            return lambda: {'file_dep': util.go_deps(pkg)}
        yield {'basename': name,
               'actions': [getgetdeps(package)],
               'uptodate': [result_dep('golistcontainer'), False]}

def task_go_base_img():
    def build_go_base_img():
        util.docker_pull_or_build_and_maybe_push(util.goBaseImage,
                                                 util.makefileDir+'/shared',
                                                 'GobaseDockerfile')
        return util.goBaseImage
    return {'actions': [build_go_base_img],
            'file_dep': ['shared/GobaseDockerfile']}

def task_images():
    return {'actions': [],
            'task_dep':['api_img','nginx_img','rethink_img','registry_img','bootstrap_img','notary_server_img','notary_signer_img','jobrunner_img']}

def task_tar():
    def run_tar():
        tarName = './dtr-%s.tar' % (util.version)
        ucpImages = []
        if util.tarUCPImages != '':
            ucpVersion = util.tarUCPImages.split(':')[1]
            tarName = './ucp-%s_dtr-%s.tar' % (ucpVersion, util.version)
            p = util.run(['docker', 'run', '--rm', util.tarUCPImages, 'images', '--list'], stdout=subprocess.PIPE)
            ucpImages = p.stdout.strip().decode('utf-8').split('\n')
            print('Taring with UCP images: %s' % ucpImages)
            if util.pullUCPImages:
                for img in ucpImages:
                    util.run(['docker', 'pull', img])
        else:
            print('Taring without images.')
        print('tar name: %s' % tarName)
        util.run(['docker', 'save', '-o', tarName] + list(map(lambda x: '%s/%s:%s' % (util.dockerHubOrg, x, util.version), util.componentNameRefs)) + util.thirdPartyImages + ucpImages)
    return {'actions': [run_tar],
            'task_dep':['images']}

def task_test():
    def run_test():
        # TODO: ignore the result of this
        util.run(['time', 'docker', 'rm', '-f', 'testrethink'])
        util.run(['time', 'docker', 'run', '--name', 'testrethink', '-d', 'jlhawn/rethinkdb:2.3.0', '--bind', 'all'])

        # TODO this can probably be make cleaner in python
        p = util.run(['sh', '-c', 'find . -name Godeps -prune -o -name vendor -prune -o \
            -name integration -prune -o -name moshpit -prune -o -name ha-integration -prune -o -name node_modules -prune -o \
            -name \'*_test.go\' -exec dirname {} \; | uniq'],
                stdout=subprocess.PIPE)
        files = p.stdout.strip()

        # Only run the Coverage when we are on Circle, don't do it in local dev environment
        go_cmd = 'go'
        if util.runCoverage:
            go_cmd = 'gocov'

        try:
            with open('test-coverage.txt', 'w') as f:
                for line in util.go_stream(['sh', '-c', go_cmd + ' test -v ' +
                        util.testFlags + ' ' + ' '.join(files.decode('utf-8').split('\n'))],
                        flags = ['--link', 'testrethink:testrethink'], asAdmin=True):
                    # this is really hacky and unreliable but we mix stdout and
                    # stderr so this is the only good way to do it
                    if util.runCoverage and line[0] == '{':
                        f.write(line)
                    else:
                        print(line, end='')
            if util.runCoverage:
                # convert JSON to Jenkins-friendly XML for codecov
                util.go_run(['sh', '-c', 'gocov-xml < test-coverage.txt > coverage.xml'],
                    runArgs={'stream': True, 'check': True}, asAdmin=True)
        finally:
            util.run(['time', 'docker', 'rm', '-f', 'testrethink'])
    return {'actions': [run_test],
            'task_dep':['go_base_img']}


def task_ginkgo_img():
    def build_ginkgo_img():
        util.docker_pull_or_build_and_maybe_push(util.ginkgoImage,
                                                 util.makefileDir+'/shared',
                                                 'GinkgoDockerfile')
        return util.ginkgoImage
    return {'actions': [build_ginkgo_img],
            'file_dep': ['shared/GinkgoDockerfile']}

def task_test_ginkgo():
    def run_test():
        util.ginkgo_go_run(['ginkgo', '-trace', '-v'] +
                util.testFlags.split(' '),
                asAdmin=True,
                runArgs={'stream': True, 'check': True})
    return {'actions': [run_test],
            'task_dep':['ginkgo_img']}

def task_golistcontainer():
    def run_golistcontainer():
        containerState, correctImage = \
                util.container_exists(util.goListContainerName, util.goBaseImage)
        if containerState == 'running' and correctImage:
            return util.goBaseImage
        if containerState != 'none' and (not correctImage or containerState ==
                                         'stopped'):
            util.kill_container(util.goListContainerName)
        util.go_run(['/sbin/init'], flags=['-d',
                                           '--name', util.goListContainerName,
                                           '--restart=always'])
        return util.goBaseImage
    # this runs every time because we always need to check if it's running before doing anything else
    # we should not use result_dep('go_base_img') because it's not relevant if the base image changed
    # TODO: make checking the existence of a container a separate task and depend on its result instead
    return {'uptodate': [False],
            'task_dep': ['go_base_img'],
            'actions': [run_golistcontainer]}

def task_bootstrap_img():
    def build_image():
        util.build_from_base(util.bootstrapName,
                        [(util.binariesDir + '/bootstrap', '/bin/dtr')],
                        '"bin/dtr"')
    return {'actions': [build_image],
            'file_dep': [util.binariesDir + '/bootstrap',
                         util.constantsFile],
            'uptodate': [result_dep('base_img')]}

def task_bootstrap_bin():
    def build_bin():
        util.go_build('./bootstrap/bootstrap', util.binariesDir)
    return {'actions': [build_bin],
            'targets': [util.binariesDir + '/bootstrap'],
            'uptodate': [result_dep('go_base_img')],
            'calc_dep': ['getbootstrapdep']}

def task_registry_img():
    def build_image():
        util.build_from_base(util.registryName,
                        [('shared/init', '/init'),
                         (util.binariesDir + '/registry', '/bin/registry'),
                         ('registry/container/start.sh', '/start.sh'),
                         ('registry/container/confs', '/etc/confd/conf.d/'),
                         ('registry/container/templates', '/etc/confd/'),
                         ('registry/container/confd.toml', '/etc/confd/confd.toml')],
                        '"/init", "--skip-runit", "/start.sh"')
    containerPath = pathlib.Path('registry/container')
    containerDepsGen = itertools.chain(
        containerPath.glob('templates/**/*'),
        containerPath.glob('confs/**/*'),
        ['registry/container/start.sh',
         util.constantsFile,
         util.binariesDir + '/registry',
         'registry/container/confd.toml'])
    containerDepsPaths = list(map(str, containerDepsGen))
    return {'actions': [build_image],
            'file_dep': containerDepsPaths,
            'uptodate': [result_dep('base_img')]}

def task_registry_bin():
    def build_bin():
        util.go_build('./registry/registry', util.binariesDir,
                      goArgs=['-tags', 'include_oss include_gcs'])
    return {'actions': [build_bin],
            'targets': [util.binariesDir + '/registry'],
            'uptodate': [result_dep('go_base_img')],
            'calc_dep': ['getregistrydep']}

def task_clidocgen_bin():
    def build_bin():
        util.go_build('./clidocgen', util.binariesDir)
    return {'actions': [build_bin],
            'targets': [util.binariesDir + '/clidocgen'],
            'uptodate': [result_dep('go_base_img')],
            'calc_dep': ['getclidocgendep']}

def task_clidocgen_img():
    def build_image():
        util.build_from_base(util.clidocgenName,
                        [(util.binariesDir + '/clidocgen', '/bin/clidocgen'),
                         ('clidocgen/templates', '/templates')],
                        '"/bin/clidocgen"', literalName=True)
    containerPath = pathlib.Path('clidocgen/container')
    containerDepsGen = itertools.chain(
        containerPath.glob('templates/**/*'),
        [util.constantsFile,
         util.binariesDir + '/clidocgen'])
    containerDepsPaths = list(map(str, containerDepsGen))
    return {'actions': [build_image],
            'file_dep': containerDepsPaths,
            'uptodate': [result_dep('base_img')]}

# TODO: run the build system with --user ${id -u}
def task_clidocs():
    def build_image():
        util.run(['docker', 'run', '--rm', '-v', util.makefileDir + '/docs/reference:/output', util.clidocgenName])
        util.run(['cp', util.makefileDir + '/clidocgen/templates/menu.md', util.makefileDir + '/docs/reference'])
    return {'actions': [build_image],
            'task_dep':['clidocgen_img']}

def task_nginx_img():
    def build_image():
        util.build_from_base(util.nginxName,
                        [('shared/init', '/init'),
                         (util.binariesDir + '/nginx', '/usr/local/nginx/'),
                         (util.binariesDir + '/dnsfix', '/dnsfix'),
                         ('nginx/container/start.sh', '/start.sh'),
                         ('nginx/container/no_license', '/no_license/no_license'),
                         ('nginx/container/confs', '/etc/confd/conf.d/'),
                         ('nginx/container/templates', '/etc/confd/'),
                         ('nginx/container/confd.toml', '/etc/confd/confd.toml')],
                        '"/init", "--skip-runit", "/start.sh"')
    containerPath = pathlib.Path('nginx/container')
    containerDepsGen = itertools.chain(
        containerPath.glob('templates/**/*'),
        containerPath.glob('confs/**/*'),
        ['nginx/container/start.sh',
         util.constantsFile,
         util.binariesDir + '/dnsfix',
         'nginx/container/no_license',
         'nginx/container/confd.toml'])
    containerDepsPaths = list(map(str, containerDepsGen))
    return {'actions': [build_image],
            'file_dep': containerDepsPaths,
            # nginx_bin is necessary because we can't depend on a directory
            'uptodate': [result_dep('base_img'), result_dep('nginx_bin')]}

def task_nginx_bin():
    def build_bin():
        util.go_run(['sh', '-c', 'cd /tmp/nginx-${NGINX_VERSION} && '+
                    'make install && '+
                    'ln -sf /dev/stdout /usr/local/nginx/logs/access.log && '+
                    'ln -sf /dev/stderr /usr/local/nginx/logs/error.log'],
                    flags=['-v',
                           util.makefileDir + '/' + util.binariesDir + '/nginx:/usr/local/nginx'])
        # go base image is the only dependency/versioning here
        return util.goBaseImage
    return {'actions': [build_bin],
            'targets': [util.binariesDir + '/nginx/sbin/nginx'],
            'uptodate': [result_dep('go_base_img')]}

def task_dnsfix_bin():
    def build_bin():
        util.go_build('./nginx/dnsfix', util.binariesDir)
    return {'actions': [build_bin],
            'targets': [util.binariesDir + '/dnsfix'],
            'uptodate': [result_dep('go_base_img')],
            'calc_dep': ['getdnsfixdep']}

def task_jobrunner_img():
    def build_image():
        util.build_from_base(util.jobrunnerName,
                        [('shared/init', '/init'),
                         (util.binariesDir + '/jobrunner_gc/gc', '/bin/gc'),
                         (util.binariesDir + '/tagmigration', '/bin/tagmigration'),
                         (util.binariesDir + '/jobrunner', '/bin/jobrunner')],
                        '"/init", "--skip-runit", "/bin/jobrunner", "worker"')
    containerDepsGen = itertools.chain(
        [util.constantsFile,
         util.binariesDir + '/jobrunner_gc/gc',
         util.binariesDir + '/tagmigration',
         util.binariesDir + '/jobrunner'])
    containerDepsPaths = list(map(str, containerDepsGen))
    return {'actions': [build_image],
            'file_dep': containerDepsPaths,
            'uptodate': [result_dep('base_img')]}

def task_jobrunner_bin():
    def build_bin():
        util.go_build('./jobrunner/jobrunner', util.binariesDir)
    return {'actions': [build_bin],
            'targets': [util.binariesDir + '/jobrunner'],
            'uptodate': [result_dep('go_base_img')],
            'calc_dep': ['getjobrunnerdep']}

def task_jobrunner_gc_bin():
    def build_bin():
        util.go_build('./cmd/gc', util.binariesDir + '/jobrunner_gc')
    return {'actions': [build_bin],
            'targets': [util.binariesDir + '/jobrunner_gc/gc'],
            'uptodate': [result_dep('go_base_img')],
            'calc_dep': ['getjobrunnergcdep']}

def task_tagmigration_bin():
    def build_bin():
        util.go_build(util.cmdDir + '/tagmigration', util.binariesDir)
    return {'actions': [build_bin],
            'targets': [util.binariesDir + '/tagmigration'],
            'uptodate': [result_dep('go_base_img')],
            'calc_dep': ['gettagmigrationdep']}

def task_api_img():
    uiPath = pathlib.Path('adminserver/ui')
    pathsToMoveToUi = list(map(str, itertools.chain(
        [uiPath.joinpath('src/index.html'),
         uiPath.joinpath('src/favicon.ico'),
         uiPath.joinpath('src/styles'),
         uiPath.joinpath('src/img'),
         uiPath.joinpath('src/fonts'),
         uiPath.joinpath('src/bundle.js')],
    )))
    pathsToMoveToSwagger = list(map(str, itertools.chain(
        uiPath.glob('swagger-ui/dist/*')
    )))
    def build_image():
        util.build_from_base(util.apiServerName,
                             [(hostPath, '/swagger') for hostPath in
                              pathsToMoveToSwagger] + \
                             [(util.binariesDir + '/server', '/bin/'),
                              (util.binariesDir + '/gc', '/bin/'),
                              ('adminserver/ui/swagger.html', '/swagger/'),
                              ('apidocgen/api_intro.md', '/swagger/'),
                             ] + \
                             [(hostPath, '/ui') for hostPath in
                              pathsToMoveToUi],
                             '"/bin/server"')
        return
    # XXX: running isfile on this list means only existing files will be considered
    # dependencies. That's why we need to list bundle.js and others explicitly in the list below
    fileDeps = list(filter(os.path.isfile, map(str, itertools.chain(
        pathsToMoveToUi,
        pathsToMoveToSwagger,
        uiPath.glob('swagger-ui/dist/**/*'),
        uiPath.glob('src/styles/**/*'),
        uiPath.glob('src/img/**/*'),
        uiPath.glob('src/fonts/**/*'),
        [uiPath.joinpath('swagger.html'),
         'apidocgen/api_intro.md']
    ))))
    return {'actions': [build_image],
            'file_dep': [util.binariesDir + '/server',
                         util.binariesDir + '/gc',
                         'adminserver/ui/src/bundle.js',
                         util.constantsFile]+fileDeps,
            'uptodate': [result_dep('base_img')]}

def task_ui_base_img():
    hashArgs = ['adminserver/ui/package.json',
                'adminserver/ui/npm-shrinkwrap.json',
                'adminserver/ui/compile.sh',
                'adminserver/ui/Dockerfile']
    fileDeps = hashArgs[:]
    for f in pathlib.Path('adminserver/ui/private-deps').glob('*'):
        if os.path.isdir(str(f)+'/.git'):
            p = util.run('git rev-parse HEAD'.split(),
                               stdout=subprocess.PIPE, cwd=str(f),
                               universal_newlines=True, check=True)
            submodHash = p.stdout.strip()
            hashArgs += [submodHash]
            fileDeps += list(path for path in f.glob('**/*') if
                             os.path.isfile(str(path)))
    uiBaseTag = util.hash_files(*hashArgs)
    uiBaseImage = util.uiBaseUntaggedImage+':'+uiBaseTag

    def build_image():
        util.docker_pull_or_build_and_maybe_push(
            uiBaseImage, 'adminserver/ui',
            'Dockerfile')
        return {'tag': uiBaseTag, 'image': uiBaseImage}
    return {'actions': [build_image],
            'file_dep': ['adminserver/ui/Dockerfile',
                         'adminserver/ui/package.json']+fileDeps}

def task_api_ui():
    def build_ui(uiBaseImage):
        with open('adminserver/ui/src/scripts/version.js', 'w') as versionF:
            versionF.write("'use strict';\nexport const version = '%s';\n" % \
                           util.version)
        os.chown('adminserver/ui/src/scripts/version.js', int(util.userId), int(util.groupId))
        if util.suffix == 'dev':
            return util.run(('docker run -u %s:%s -e DEV=1 -tv %s/adminserver/ui:/build %s' % \
                    (util.userId, util.groupId, util.makefileDir, uiBaseImage)).split()).returncode == 0
        else:
            return util.run(('docker run -u %s:%s -tv %s/adminserver/ui:/build %s' % \
                             (util.userId, util.groupId, util.makefileDir, uiBaseImage)).split()).returncode == 0
    uiPath = pathlib.Path('adminserver/ui')
    fileDeps = list(filter(os.path.isfile, map(str, itertools.chain(
        uiPath.glob('src/img/**/*'),
        uiPath.glob('src/fonts/**/*'),
        uiPath.glob('src/styles/**/*'),
        uiPath.glob('src/scripts/**/*.css'),
        [p for p in uiPath.glob('src/scripts/**/*.js') if not
         str(p).endswith('src/scripts/version.js')],
        [uiPath.joinpath('src/webpack.config.js'),
         uiPath.joinpath('src/webpack.dev.config.js')],
    ))))
    return {'actions': [build_ui],
            'getargs': {'uiBaseImage': ('ui_base_img', 'image')},
            'targets': ['adminserver/ui/src/bundle.js'],
            'file_dep': fileDeps,
            'uptodate': [result_dep('ui_base_img')]}

def task_gc_bin():
    def build_bin():
        util.go_build(util.cmdDir + '/gc', util.binariesDir,
                      goArgs=['-tags', 'include_oss include_gcs'])
    return {'actions': [build_bin],
            'targets': [util.binariesDir + '/gc'],
            'uptodate': [result_dep('go_base_img')],
            'calc_dep': ['getgcdep']}

def task_api_bin():
    def build_bin():
        util.go_build('./adminserver/server', util.binariesDir,
                      goArgs=['-tags', 'include_oss include_gcs'])
    return {'actions': [build_bin],
            'targets': [util.binariesDir + '/server'],
            'uptodate': [result_dep('go_base_img')],
            'calc_dep': ['getapidep']}

def task_rethink_base_img():
    rethinkBaseTag = util.hash_files('rethinkdb/container/Dockerfile')
    rethinkBaseImage = util.rethinkBaseUntaggedImage+':'+rethinkBaseTag
    def build_image():
        util.docker_pull_or_build_and_maybe_push(
            rethinkBaseImage, 'rethinkdb/container', 'Dockerfile')
        return {'tag': rethinkBaseTag, 'image': rethinkBaseImage}
    return {'actions': [build_image],
            'uptodate': [result_dep('go_base_img')],
            'file_dep': ['rethinkdb/container/Dockerfile']}

def task_rethink_img():
    def build_image(rethinkBaseImage):
        util.build_from_base(util.rethinkName,
                             [('rethinkdb/container/start.sh', '/'),
                              (util.binariesDir + '/arggen', '/bin/')],
                             '"/start.sh"',
                             baseImage=rethinkBaseImage)
        return
    containerDepsGen = itertools.chain([
         util.constantsFile,
         util.binariesDir + '/arggen',
         'rethinkdb/container/Dockerfile'])
    containerDepsPaths = list(map(str, containerDepsGen))
    return {'actions': [build_image],
            'file_dep': containerDepsPaths,
            'getargs': {'rethinkBaseImage': ('rethink_base_img', 'image')},
            'uptodate': [result_dep('rethink_base_img')]}

# used for the rethink container
def task_arggen_bin():
    def build_bin():
        util.go_build('./rethinkdb/arggen', util.binariesDir)
    return {'actions': [build_bin],
            'targets': [util.binariesDir + '/arggen'],
            'uptodate': [result_dep('go_base_img')],
            'calc_dep': ['getarggendep']}

def task_moshpit_bin():
    def build_bin():
        util.go_build('./moshpit/cmd/moshpit', util.binariesDir)
    return {'actions': [build_bin],
            'targets': [util.binariesDir + '/moshpit'],
            'uptodate': [result_dep('go_base_img')],
            'calc_dep': ['getmoshpitdep']}

def task_moshpit_img():
    def build_image():
        util.build_from_base(util.moshpitName,
                        [(util.binariesDir + '/moshpit', '/bin/moshpit')],
                        '"/bin/moshpit"')
    return {'actions': [build_image],
            'file_dep': [
                'moshpit/container/Dockerfile',
                util.constantsFile,
                util.binariesDir + '/moshpit'],
            'uptodate': [result_dep('base_img')]}

def task_run_moshpit_dropper():
    def run_moshpit_dropper():
        # TODO: consider not using latest tag
        util.run(["docker", "push", "dockerhubenterprise/dtr-moshpit:latest"], check=True)
        # TODO: make the config file location configurable
        # TODO: use the right tag from the image build
        util.run(["docker", "run", "--rm", "-v", util.makefileDir + "/moshpit/config.yml:/config/config.yml", util.dockerHubOrg + "/" + util.moshpitName + ":latest", "--debug", "dropper", "--config-file", "/config/config.yml"], check=True)
    return {'actions': [run_moshpit_dropper],
            'task_dep':['moshpit_img']}

# TODO build notary ourselves
# TODO wire up dependencies for this properly
def task_notary_server_bin():
    def build_bin():
        util.run(['docker', 'run', '-u', '%s:%s' % (util.userId, util.groupId), '-t', '--rm', '-v', '%s:/dhe-deploy' % util.makefileDir,
                '--entrypoint', 'cp', 'dockersecurity/notary_autobuilds:server',
                '/go/bin/notary-server', '/dhe-deploy/' + util.binariesDir + '/notary_server'])
    return {'actions': [build_bin],
            'targets': [util.binariesDir + '/notary_server'],
            'uptodate': [run_once]}

def task_notary_server_img():
    def build_image():
        util.build_from_base(util.notaryServerName,
                        [('shared/init', '/init'),
                         ('notary_server/container/start.sh', '/start.sh'),
                         (util.binariesDir + '/notary_server', '/bin/notary_server'),
                         ('notary_server/container/confs', '/etc/confd/conf.d/'),
                         ('notary_server/container/templates', '/etc/confd/'),
                         ('notary_server/container/confd.toml', '/etc/confd/confd.toml')],
                        '"/init", "--skip-runit", "/start.sh"')
    containerPath = pathlib.Path('notary_server/container')
    containerDepsGen = itertools.chain(
        containerPath.glob('templates/**/*'),
        containerPath.glob('confs/**/*'),
        ['notary_server/container/start.sh',
         util.constantsFile,
         util.binariesDir + '/notary_server',
         'notary_server/container/confd.toml'])
    containerDepsPaths = list(map(str, containerDepsGen))
    return {'actions': [build_image],
            'file_dep': containerDepsPaths,
            'uptodate': [result_dep('base_img')]}

# TODO wire up dependencies for this properly
# TODO build notary ourselves
def task_notary_signer_bin():
    def build_bin():
        util.run(['docker', 'run', '-u', '%s:%s' % (util.userId, util.groupId), '-t', '--rm', '-v', '%s:/dhe-deploy' % util.makefileDir,
                '--entrypoint', 'cp', 'dockersecurity/notary_autobuilds:signer',
                '/go/bin/notary-signer', '/dhe-deploy/' + util.binariesDir + '/notary_signer'])
    return {'actions': [build_bin],
            'targets': [util.binariesDir + '/notary_signer'],
            'uptodate': [run_once]}

def task_notary_signer_img():
    def build_image():
        util.build_from_base(util.notarySignerName,
                        [('shared/init', '/init'),
                         ('notary_signer/container/start.sh', '/start.sh'),
                         (util.binariesDir + '/notary_signer', '/bin/notary_signer'),
                         ('notary_signer/container/confs', '/etc/confd/conf.d/'),
                         ('notary_signer/container/templates', '/etc/confd/'),
                         ('notary_signer/container/confd.toml', '/etc/confd/confd.toml')],
                        '"/init", "--skip-runit", "/start.sh"')
    containerPath = pathlib.Path('notary_signer/container')
    containerDepsGen = itertools.chain(
        containerPath.glob('templates/**/*'),
        containerPath.glob('confs/**/*'),
        ['notary_signer/container/start.sh',
         util.constantsFile,
         util.binariesDir + '/notary_signer',
         'notary_signer/container/confd.toml'])
    containerDepsPaths = list(map(str, containerDepsGen))
    return {'actions': [build_image],
            'file_dep': containerDepsPaths,
            'uptodate': [result_dep('base_img')]}
