#
# Copyright (C) 2016 Rolf Neugebauer <rolf.neugebauer@docker.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#

"""Main entry point to run the regression test on the local host"""

import argparse
import base64
import os
import sys

import rt.base
import rt.log


ENV_EXTRA_DEBUG = False
ENV_RT_ROOT = ""
ENV_RT_PROJECT_ROOT = ""
ENV_RT_OS = ""
ENV_RT_OS_VER = ""

ENV_VARS = {}

COMMANDS = ['list', 'info', 'run']


def _parse_env(env):
    """Parse the --env argument"""
    global ENV_VARS
    if not len(env):
        return
    for var in env.split(','):
        if not len(var):
            continue
        k, _, v = var.partition('=')
        ENV_VARS[k] = v


def main():
    """Main entry point to run regressions tests locally"""

    parser = argparse.ArgumentParser()
    parser.description = \
        "Run or provide information about local regression tests"
    parser.epilog = \
        "Arguments may also be passed in as a single base64 encoded string"

    parser.add_argument("-p", "--project", default="cases",
                        help="Top-level directory of a project's test cases")

    parser.add_argument("-r", "--resultdir", default="./_results",
                        help="Directory to place results in")

    parser.add_argument("-l", "--labels",
                        help="Labels to apply (comma separated)")

    parser.add_argument("-e", "--env", default="",
                        help="Additional environment variables" +
                        " (<VAR0>=<VAL0>,<VAR1>=<VAL1>)")

    parser.add_argument("-x", "--extra", action='store_true',
                        help="Add extra debug info to log files")

    parser.add_argument("--logger",
                        help="<host>:<port> of a socket to log stdout to")

    parser.add_argument("-v", "--verbose", action='count', default=0,
                        help="Increase verbosity level")

    parser.add_argument('cmd', nargs=1,
                        choices=COMMANDS,
                        help="Command to execute")

    parser.add_argument('tests', nargs=argparse.REMAINDER)

    # support passing arguments as a base64 encoded string
    if len(sys.argv) == 2 and \
       sys.argv[1] not in COMMANDS and sys.argv[1] not in ['-h', '--help']:
        tmp = base64.b64decode(sys.argv[1]).decode('ascii')
        argv = tmp.split()
        print("Arguments: ", argv)
        args = parser.parse_args(argv)
    else:
        args = parser.parse_args()

    _parse_env(args.env)

    global ENV_EXTRA_DEBUG
    global ENV_RT_ROOT
    global ENV_RT_PROJECT_ROOT
    global ENV_RT_OS
    global ENV_RT_OS_VER

    os_info = rt.misc.get_os_info()

    ENV_EXTRA_DEBUG = args.extra
    ENV_RT_ROOT = os.getcwd()
    ENV_RT_PROJECT_ROOT = os.path.realpath(args.project)
    ENV_RT_OS = os_info['os']
    ENV_RT_OS_VER = os_info["os"] + '_' + os_info["version"]

    # add utilities to the path
    os.environ["PATH"] += ":%s/lib/utils" % ENV_RT_ROOT

    labels, _ = rt.misc.str2labels(args.labels)

    if len(args.tests) > 1:
        sys.stderr.write("Only one test parameter supported. %d provided:\n"
                         % len(args.tests))
        sys.stderr.write("%s.\n\n" % args.tests)
        parser.print_help()
        return 1

    testpat = None if len(args.tests) == 0 else args.tests[0]

    # Check for CI
    if "CIRCLECI" in os.environ:
        labels.add("circleci")
    if "APPVEYOR" in os.environ:
        labels.add("appveyor")

    labels.add(os_info["os"])
    labels.add(os_info["os"] + '_' + os_info["version"])
    labels.add(os_info["name"])

    print ("LABELS: %s" % ", ".join(labels))
    tests = rt.base.Project(args.project, labels=labels, testpat=testpat)

    if "list" in args.cmd:
        ret = tests.list()
    elif "info" in args.cmd:
        ret = tests.info()
    elif "run" in args.cmd:
        rt.log.init(args.resultdir, args.verbose, args.logger)
        ret = tests.run()
    else:
        ret = -1

    return ret
