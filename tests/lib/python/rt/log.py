#
# Copyright (C) 2016 Rolf Neugebauer <rolf.neugebauer@docker.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#

"""Logging and results reporting for the RT framework"""

# This configures the logging for the RT framework.  The default is
# that everything is logged to a logfile and results are written to a
# CSV file, if a log directory is specified, and only warnings and
# higher are displayed on the console. There are also per test log
# files for easier debugging.  Finally, the console log can optionally
# be redirected to a socket, provided something is listening on the
# other end.
#
# In addition to the standard log levels we register a couple of
# additional log levels:
# - LOG_LVL_STDERR: This is where stderr is logged to from commands invoked
# - LOG_LVL_STDOUT: This is where stdout is logged to from commands invoked
#
# - LOG_LVL_SKIP: Used to log the result of tests which have been skipped
# - LOG_LVL_PASS: Used to log the result of tests which pass
# - LOG_LVL_FAIL: Used to log the result of tests which fail
# - LOG_LVL_SUM: Used to log a summary of results
#
# This module deliberately does not implement a class and stores state
# in module global variables to make it easier to call it from
# anywhere in the framework without having to pass object references
# around.


import csv
import datetime
import logging
import logging.handlers
import os
import pickle
import platform
import struct
import subprocess
import time
import uuid
import sys

import rt.misc
from rt.colour import COLOUR_RESET, COLOUR_BOLD, COLOUR_FOREGROUND
from rt.colour import COLOUR_RED, COLOUR_GREEN, COLOUR_YELLOW, COLOUR_GREY
from rt.colour import COLOUR_MAGENTA

if sys.version_info < (3, 0):
    import SocketServer as socketserver
else:
    import socketserver


# Test return codes
TEST_RC_PASS = 0  # Test passed
TEST_RC_FAIL = 1  # Test failed
TEST_RC_SKIP = 2  # Test was skipped
TEST_RC_CANCEL = 3  # Test was cancelled

LOG_LVL_STDERR = logging.DEBUG + 5
LOG_LVL_STDOUT = logging.DEBUG + 6

LOG_LVL_SKIP = logging.WARNING + 4
LOG_LVL_PASS = logging.WARNING + 5
LOG_LVL_CANCEL = logging.WARNING + 6
LOG_LVL_FAIL = logging.WARNING + 7
LOG_LVL_SUM = logging.WARNING + 8


logger = logging.getLogger('rt')


_LVL_COLOURS = {
    # Standard log levels
    logging.CRITICAL: COLOUR_RESET + COLOUR_BOLD + COLOUR_RED,
    logging.ERROR: COLOUR_RESET + COLOUR_RED,
    logging.WARNING: COLOUR_RESET + COLOUR_YELLOW,
    logging.INFO: COLOUR_RESET + COLOUR_FOREGROUND,
    logging.DEBUG: COLOUR_RESET + COLOUR_GREY,

    # Custom log-levels
    LOG_LVL_STDERR: COLOUR_RESET + COLOUR_RED,
    LOG_LVL_STDOUT: COLOUR_RESET + COLOUR_FOREGROUND,

    LOG_LVL_SKIP: COLOUR_RESET + COLOUR_YELLOW,
    LOG_LVL_PASS: COLOUR_RESET + COLOUR_GREEN,
    LOG_LVL_CANCEL: COLOUR_RESET + COLOUR_MAGENTA,
    LOG_LVL_FAIL: COLOUR_RESET + COLOUR_BOLD + COLOUR_RED,
}


# Other state
_UUID = uuid.uuid4()
ENV_RT_RESULTS = None
_START_TIME = datetime.datetime.utcnow()
_CSV_FILE = None
_CSV_WRITER = None
_SOCKET_HANDLER = None


class _ColouredFormatter(logging.Formatter):
    """Format a log messages with colour"""

    # This is a bit of a hack as we use the normal formatter but
    # override the level name with a colourised level name. We can
    # make it a bit more elaborate, but this seems the quickest way to
    # do what we want.
    def format(self, record):
        """Replace the logname with a colourised version"""
        if record.levelno in _LVL_COLOURS:
            record.levelname = _LVL_COLOURS[record.levelno] + \
                               "[%-6s]" % record.levelname + COLOUR_RESET
        return logging.Formatter.format(self, record)


class _LogRecvHandler(socketserver.StreamRequestHandler):
    """A simple handler for log records"""
    formatter = _ColouredFormatter()

    def handle(self):
        """The format is 4 Bytes followed by a pickled log record. Unpickle,
        format and print."""
        while True:
            chunk = self.connection.recv(4)
            if len(chunk) < 4:
                break
            slen = struct.unpack(">L", chunk)[0]
            chunk = self.connection.recv(slen)
            while len(chunk) < slen:
                chunk = chunk + self.connection.recv(slen - len(chunk))
            obj = pickle.loads(chunk)
            record = logging.makeLogRecord(obj)
            # run through formatter to get the colourised level name
            msg = self.formatter.format(record)
            print("%s %s" % (record.levelname, msg))


def log_svr(port):
    """A simple log server for remote logging. Handles one connection and
    then exits."""
    addr = ("", port)
    socketserver.TCPServer.allow_reuse_address = True
    server = socketserver.TCPServer(addr, _LogRecvHandler)
    server.handle_request()
    server.server_close()


def init(log_dir=None, verbose=0, remote=None):
    """Initialise logging for the RT framework."""

    # Default LOG level is debug
    logger.setLevel(logging.DEBUG)

    # Add new log levels for stdout/stderr, and results
    logging.addLevelName(LOG_LVL_STDERR, 'STDERR')
    logging.addLevelName(LOG_LVL_STDOUT, 'STDOUT')

    logging.addLevelName(LOG_LVL_SKIP, 'SKIP')
    logging.addLevelName(LOG_LVL_PASS, 'PASS')
    logging.addLevelName(LOG_LVL_CANCEL, 'CANCEL')
    logging.addLevelName(LOG_LVL_FAIL, 'FAIL')
    logging.addLevelName(LOG_LVL_SUM, 'SUMMARY')

    # Create a handler for the logfile if requested
    if log_dir:
        global ENV_RT_RESULTS
        log_fn = None
        ENV_RT_RESULTS = os.path.join(log_dir, _UUID.__str__())
        ENV_RT_RESULTS = os.path.abspath(ENV_RT_RESULTS)
        if not os.path.isdir(ENV_RT_RESULTS):
            os.makedirs(ENV_RT_RESULTS)
        log_fn = os.path.join(ENV_RT_RESULTS, "TESTS.log")

        logf = logging.FileHandler(log_fn)
        logf.setLevel(logging.DEBUG)
        # Add milliseconds with a . (instead of a ,)
        logf_form = logging.Formatter(
            fmt='[%(levelname)-8s] %(asctime)s.%(msecs)03d: %(message)s',
            datefmt='%Y-%m-%d %H:%M:%S')
        # always use UTC
        logf.converter = time.gmtime
        logf.setFormatter(logf_form)
        logger.addHandler(logf)

        # open CSV file
        global _CSV_FILE
        global _CSV_WRITER
        cols = ["ID", "Timestamp", "Duration", "Name", "Result", "Message"]
        fn = os.path.join(ENV_RT_RESULTS, "TESTS.csv")
        _CSV_FILE = open(fn, 'w')
        _CSV_WRITER = csv.writer(_CSV_FILE)
        _CSV_WRITER.writerow(cols)

    # Create console handler, which by default only logs results
    if remote:
        host, port = remote.split(':')
        logc = logging.handlers.SocketHandler(host, int(port))
        global _SOCKET_HANDLER
        _SOCKET_HANDLER = logc
    else:
        logc = logging.StreamHandler(sys.stdout)
    if verbose == 0:
        logc.setLevel(LOG_LVL_PASS)
    elif verbose == 1:
        logc.setLevel(LOG_LVL_SKIP)
    elif verbose == 2:
        logc.setLevel(logging.INFO)
    else:
        logc.setLevel(logging.DEBUG)
    logc_form = _ColouredFormatter('%(levelname)s %(message)s')
    logc.setFormatter(logc_form)
    logger.addHandler(logc)
    return


# Map return codes to log levels
_RC_TO_LVL = {
    TEST_RC_PASS: LOG_LVL_PASS,
    TEST_RC_FAIL: LOG_LVL_FAIL,
    TEST_RC_SKIP: LOG_LVL_SKIP,
    TEST_RC_CANCEL: LOG_LVL_CANCEL,
}


_TEST_HANDLER = None


def _darwin_asl(msg):
    cmd = ["syslog", "-s",
           "-k", "Facility", "com.docker.rt-local",
           "-k", "Level", "Notice",
           "-k", "Sender", "Docker",
           "-k", "Message", msg]
    subprocess.call(cmd)


def _test_start_darwin(name):
    _darwin_asl("Start test: %s" % name)


def test_start(name):
    """Create an additional logger for the test with @name"""
    if not ENV_RT_RESULTS:
        return
    global _TEST_HANDLER
    log_fn = os.path.join(ENV_RT_RESULTS, "%s.log" % name)
    _TEST_HANDLER = logging.FileHandler(log_fn)
    _TEST_HANDLER.setLevel(logging.DEBUG)
    # Add milliseconds with a . (instead of a ,)
    logf_form = logging.Formatter(
        fmt='[%(levelname)-8s] %(asctime)s.%(msecs)03d: %(message)s',
        datefmt='%Y-%m-%d %H:%M:%S')
    # always use UTC
    _TEST_HANDLER.converter = time.gmtime
    _TEST_HANDLER.setFormatter(logf_form)
    logger.addHandler(_TEST_HANDLER)
    if platform.system() == "Darwin":
        _test_start_darwin(name)
    return


def _test_end_darwin(name):
    _darwin_asl("End test: %s" % name)


def test_end(name):
    """Stop logger for a individual test"""
    global _TEST_HANDLER
    if platform.system() == "Darwin":
        _test_end_darwin(name)
    if not _TEST_HANDLER:
        return
    logger.removeHandler(_TEST_HANDLER)
    _TEST_HANDLER.close()
    _TEST_HANDLER = None
    return


def _os_logs_darwin(name):
    asllog_fn = os.path.join(ENV_RT_RESULTS, "%s.asl.log" % name)
    cmd = ["syslog", "-F",
           "$Time $Host $(Sender)[$(Facility)][$(PID)]<$((Level)(str))>: " +
           "$Message",
           "-k", "Sender",  "Seq", "Docker", "-o",
           "-k", "Sender",  "Seq", "docker", "-o",
           "-k", "Message", "Seq", "Docker", "-o",
           "-k", "Message", "Seq", "docker"]
    with open(asllog_fn, "w") as asllog:
        subprocess.call(cmd, stdout=asllog)


def os_logs(name):
    if platform.system() == "Darwin":
        _os_logs_darwin(name)


def result(name, rc, duration=0, message=None):
    """Log the result of a test. @name is the test name, @rc must be one
    of TEST_RC_*. @duration and @message is optional. @duration is
    expected to be a float of the number of seconds.
    """

    ts = datetime.datetime.utcnow()
    _CSV_WRITER.writerow([_UUID, ts.isoformat(' '), duration, name, rc,
                          "" if not message else message])

    msg = "%-50s" % name
    if duration:
        msg += " %7.3fs" % duration
    if message:
        msg += " %s" % message

    logger.log(_RC_TO_LVL[rc], msg)
    return


def _read_version():
    """A project can create a VERSION.txt in ENV_RT_RESULTS. This function
    reads it and returns it as a string (or UNKNOWN/UNDEFINED if not
    present)"""

    if not ENV_RT_RESULTS:
        return "UNKNOWN"
    fn = os.path.join(ENV_RT_RESULTS, "VERSION.txt")
    if not os.path.isfile(fn):
        return "UNDEFINED"
    f = open(fn, 'r')
    version = f.readline()
    f.close()
    return version.strip()


def finish(results):
    """Log a summary of tests"""

    end_time = datetime.datetime.utcnow()
    diff = end_time - _START_TIME
    duration = diff.total_seconds()
    version = _read_version()

    logger.log(LOG_LVL_SUM, "Logdir:    %s" % _UUID)
    logger.log(LOG_LVL_SUM, "Version:   %s" % version)
    logger.log(LOG_LVL_SUM, "Passed:    %d" % results[TEST_RC_PASS])
    logger.log(LOG_LVL_SUM, "Failed:    %d" % results[TEST_RC_FAIL])
    logger.log(LOG_LVL_SUM, "Cancelled: %d" % results[TEST_RC_CANCEL])
    logger.log(LOG_LVL_SUM, "Skipped:   %d" % results[TEST_RC_SKIP])
    logger.log(LOG_LVL_SUM, "Duration:  %.2fs" % duration)

    if _SOCKET_HANDLER is not None:
        _SOCKET_HANDLER.close()

    if ENV_RT_RESULTS:
        cols = ["ID", "Version", "Start Time", "End Time", "Duration",
                "Passed", "Failed", "Skipped",
                "Labels", "OS", "OS Name", "OS Version"]
        fn = os.path.join(ENV_RT_RESULTS, "SUMMARY.csv")
        f = open(fn, 'w')
        writer = csv.writer(f)

        os_info = rt.misc.get_os_info()
        res = []
        res.append(_UUID)
        res.append(version)
        res.append(_START_TIME.isoformat(' '))
        res.append(end_time.isoformat(' '))
        res.append(duration)
        res.append(results[TEST_RC_PASS])
        res.append(results[TEST_RC_FAIL])
        res.append(results[TEST_RC_SKIP])
        res.append(os_info["os"])
        res.append(os_info["name"])
        res.append(os_info["version"])

        writer.writerow(cols)
        writer.writerow(res)

        f.close()

    _CSV_FILE.close()
    return
