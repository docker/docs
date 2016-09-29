#
# Copyright (C) 2016 Rolf Neugebauer <rolf.neugebauer@docker.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#

"""Execute a command and re-direct stdout/stderr to a log"""

import os
import subprocess
import threading
import time

import rt.local
import rt.log

_WIN_BASH = None


def _get_win_bash():
    """On Windows, we can't use the WSL bash.exe which may be in the
    path before the MSYS bash.exe. This function tries to find
    the right bash.exe and stashes the result in _WIN_BASH"""
    global _WIN_BASH
    if _WIN_BASH:
        return _WIN_BASH

    for p in os.environ["PATH"].split(';'):
        if os.path.isfile(os.path.join(p, "bash.exe")):
            if "system32" not in p:
                _WIN_BASH = os.path.join(p, "bash.exe")
    if not _WIN_BASH:
        _WIN_BASH = "bash.exe"
    return _WIN_BASH


def _run(cmd, cwd=None, env=None):
    """Run a command and return the return code"""

    rt.log.logger.debug("Executing: %s in %s" % (cmd, cwd))

    # exec the command
    p = subprocess.Popen(cmd, cwd=cwd, env=env,
                         stdout=subprocess.PIPE, stderr=subprocess.PIPE)

    def log(stream, cb):
        """Executed by a Thread to log either stdout or stderr to log"""
        while True:
            out = stream.readline()
            if out:
                cb(out.rstrip())
            else:
                break
        return

    logger = rt.log.logger

    # kick off a thread per output stream
    out_thd = threading.Thread(
        target=log,
        args=(p.stdout, lambda s: logger.log(rt.log.LOG_LVL_STDOUT, s)))
    err_thd = threading.Thread(
        target=log,
        args=(p.stderr, lambda s: logger.log(rt.log.LOG_LVL_STDERR, s)))

    # kick the threads
    out_thd.start()
    err_thd.start()

    # wait till they terminate
    out_thd.join()
    err_thd.join()

    # wait for process to terminate (it should have) and return code
    while p.poll() is None:
        time.sleep(0.5)
    return p.returncode


def shexec(script, arg=None, cwd=None, name=None, labels=[]):
    """Execute a shell script in directory"""

    # Add RT environment variables
    env = os.environ
    env["RT_ROOT"] = rt.local.ENV_RT_ROOT
    env["RT_UTILS"] = os.path.join(rt.local.ENV_RT_ROOT, "lib/utils")
    env["RT_PROJECT_ROOT"] = rt.local.ENV_RT_PROJECT_ROOT
    env["RT_OS"] = rt.local.ENV_RT_OS
    env["RT_OS_VER"] = rt.local.ENV_RT_OS_VER
    env["RT_LABELS"] = ":".join(labels)
    env["RT_TEST_NAME"] = name if name else "UNKNOWN"
    for k in rt.local.ENV_VARS:
        env[k] = rt.local.ENV_VARS[k]

    if rt.log.ENV_RT_RESULTS:
        env["RT_RESULTS"] = rt.log.ENV_RT_RESULTS

    if os.name == 'nt':
        # On Windows execute bash.exe with the command as argument
        cmd = [_get_win_bash()]
        env["MSYS_NO_PATHCONV"] = "1"
    else:
        cmd = ['/bin/sh']

    if rt.local.ENV_EXTRA_DEBUG:
        cmd.append('-x')

    cmd.append(script)
    if arg:
        cmd.append(arg)

    return _run(cmd, cwd, env)
