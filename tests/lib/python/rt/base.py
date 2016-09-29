#
# Copyright (C) 2016 Rolf Neugebauer <rolf.neugebauer@docker.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#

"""Base classes for tests and groups of tests"""

import datetime
import io
import os

import rt.lexec
import rt.log
import rt.misc
from rt.colour import COLOUR_RESET, COLOUR_GREEN, COLOUR_YELLOW

RT_TEST_CANCEL = 253
_VALID_TEST_NAMES = set(('test.sh', ))


def _ret_to_rc(ret):
    """Convert a shell return code to a canonical return code"""
    if ret == 0:
        return rt.log.TEST_RC_PASS
    elif ret == RT_TEST_CANCEL:
        return rt.log.TEST_RC_CANCEL
    else:
        return rt.log.TEST_RC_FAIL


def _name_strip(dirname):
    """Directories may be prefixed with numbers (and - or _) to force
    ordering. This functions removes leading numerals and seperators"""
    return dirname.lstrip('0123456789-_')


def _parse_file(f, tag, multi=False):
    """Parse a file for tags of the form '# <TAG>: and return the string
    after colon. If @multi is specified, scan for multiple occurrences
    and concatenate the strings. Return None if @tag was not found."""
    # start at the start
    f.seek(0, 0)
    res = ""
    for line in f:
        if not line.startswith('#'):
            continue
        line = line.strip('#').strip()
        if line.startswith("%s:" % tag):
            _, _, tmp = line.partition(':')
            if not multi:
                return tmp.strip()
            else:
                res += " %s" % tmp
    return None if len(res) == 0 else res.strip()


class _Test(object):
    """A test object"""

    def __init__(self, group, path, name, cmd):
        """Initialise the object"""
        self.group = group
        self.path = path
        self.name = name
        self.cmd = cmd
        f = io.open(os.path.join(self.path, self.cmd), encoding="utf8")
        self.summary = _parse_file(f, "SUMMARY")
        self.authors = _parse_file(f, "AUTHOR", True)
        repeat = _parse_file(f, "REPEAT")
        lstr = _parse_file(f, "LABELS")
        self.labels, self.not_labels = rt.misc.str2labels(lstr)
        f.close()

        self.repeat = 1
        if repeat:
            for e in repeat.split():
                if ":" in e:
                    l, _, r = e.partition(':')
                    if l in self.group.labels:
                        self.repeat = int(r)
                else:
                    self.repeat = int(e)

        # dump info into the log
        rt.log.logger.debug("Test: name=%s" % self.name)
        rt.log.logger.debug("Test: path=%s" % self.path)
        rt.log.logger.debug("Test: summary=%s" % self.summary)
        rt.log.logger.debug("Test: authors=%s" % self.authors)
        rt.log.logger.debug("Test: labels=%s" % self.labels)
        rt.log.logger.debug("Test: !labels=%s" % self.not_labels)
        rt.log.logger.debug("Test: repeat=%d" % self.repeat)

        return

    def check_labels(self):
        """Check if test should be run based on labels"""
        nl = self.not_labels & self.group.labels
        if not len(nl) == 0:
            return 0, ", ".join(['!' + x for x in nl])

        l = self.labels & self.group.labels
        if (len(self.labels) == 0 or not len(l) == 0):
            return 1, ", ".join(l)
        else:
            return 0, ", ".join(self.labels)

    def run(self):
        """Run the test"""
        check, info = self.check_labels()
        if check == 0:
            rt.log.result(self.name, rt.log.TEST_RC_SKIP, message=info)
            return rt.log.TEST_RC_SKIP

        for i in range(0, self.repeat):
            if self.repeat == 1:
                testName = self.name
            else:
                testName = "%s-%03d" % (self.name, i)

            rt.log.test_start(testName)
            rt.log.logger.info("Running Test %s in %s" % (testName, self.path))

            starttime = datetime.datetime.now()

            ret = rt.lexec.shexec(self.cmd, arg=None, cwd=self.path,
                                  name=self.name, labels=self.group.labels)

            endtime = datetime.datetime.now()
            diff = endtime - starttime

            rc = _ret_to_rc(ret)
            rt.log.result(testName, rc, duration=diff.total_seconds())
            rt.log.os_logs(testName)
            rt.log.test_end(testName)
        return rc


class _Group(object):
    """A group of tests"""

    def __init__(self, parent, path, name=None, labels=set(), testpat=None):
        """Initialise the object"""

        self.parent = parent
        self.name = name
        self.path = path
        self.labels = labels
        self.testpat = testpat
        self.repeat = 1
        self.summary = ""

        self.items = []

        self.cmd = None
        lstr = None
        fn = os.path.join(self.path, "group.sh")
        if os.path.isfile(fn):
            self.cmd = "group.sh"
            f = io.open(fn, encoding="utf8")
            self.summary = _parse_file(f, "SUMMARY")
            lstr = _parse_file(f, "LABELS")
            if not self.name:
                self.name = _parse_file(f, "NAME")
            f.close()
        self.cfg_labels, self.cfg_not_labels = rt.misc.str2labels(lstr)

        for f in os.listdir(path):
            cpath = os.path.join(path, f)
            if not os.path.isdir(cpath) or f.startswith('_'):
                continue

            if self.name:
                cname = self.name + '.' + _name_strip(f)
            else:
                cname = _name_strip(f)

            # look into sub-directories and determine if they are
            # tests or sub-groups. Add them
            files = os.listdir(cpath)
            tmp = set(files) & _VALID_TEST_NAMES
            if not len(tmp) == 0:
                # Test: Check if the name matches the test pattern. If not skip
                if self.testpat is not None and \
                   not cname.startswith(self.testpat):
                    continue
                t = _Test(self, cpath, cname, tmp.pop())
                self.items.append(t)
            else:
                # must be a group, recurse
                g = _Group(self, cpath, cname, labels, testpat)
                if len(g):
                    self.items.append(g)

        return

    def _check_labels(self):
        """Check if the group should be run based on labels"""
        nl = self.cfg_not_labels & self.labels
        if not len(nl) == 0:
            return 0, ", ".join(['!' + x for x in nl])

        l = self.cfg_labels & self.labels
        if (len(self.cfg_labels) == 0 or not len(l) == 0):
            return 1, ", ".join(l)
        else:
            return 0, ", ".join(self.cfg_labels)

    def __len__(self):
        """Return the number of elements in the group"""
        return len(self.items)

    def list(self):
        """Print a list of tests in this group (recursively)"""
        res, info = self._check_labels()
        if res == 0:
            pre = COLOUR_RESET + COLOUR_YELLOW + "SKIP" + COLOUR_RESET
            print ("%s %-50s      %s" % (pre, self.name, info))
            return

        for item in self.items:
            if isinstance(item, _Test):
                res, info = item.check_labels()
                if res == 0:
                    pre = COLOUR_RESET + COLOUR_YELLOW + "SKIP" + COLOUR_RESET
                else:
                    pre = COLOUR_RESET + COLOUR_GREEN + "RUN " + COLOUR_RESET
                post = info
                repeat = "    " if item.repeat == 1 else "[%2d]" % item.repeat
                print ("%s %-50s %s %s" % (pre, item.name, repeat, post))
            else:
                item.list()

    def info(self):
        """Print a list of tests in this group (recursively)"""
        if self.summary:
            print ("%-45s %s" % (self.name, self.summary))

        for item in self.items:
            if isinstance(item, _Test):
                print ("%-45s %s" % (item.name, item.summary))
            else:
                item.info()

    def _init(self):
        """Run group initialiser"""
        if not self.cmd:
            return 0
        rt.log.logger.info("%s::ginit()" % self.name)
        ret = rt.lexec.shexec(self.cmd, arg="init", cwd=self.path,
                              name=self.name + ".init", labels=self.labels)
        return _ret_to_rc(ret)

    def _deinit(self):
        """Run group de-initialiser"""
        if not self.cmd:
            return 0
        rt.log.logger.info("%s::gdeinit()" % self.name)
        ret = rt.lexec.shexec(self.cmd, arg="deinit", cwd=self.path,
                              name=self.name + ".deinit", labels=self.labels)
        return _ret_to_rc(ret)

    def run(self):
        """Run tests in this group and its sub-groups"""

        rt.log.logger.info("Entering Group %s in %s" % (self.name, self.path))

        # Count Passed, failed and skipped tests
        res = {rt.log.TEST_RC_PASS: 0,
               rt.log.TEST_RC_FAIL: 0,
               rt.log.TEST_RC_SKIP: 0,
               rt.log.TEST_RC_CANCEL: 0}

        # Skip group if the labels say so
        rc, info = self._check_labels()
        if rc == 0:
            res[rt.log.TEST_RC_SKIP] = 1
            rt.log.result(self.name, rt.log.TEST_RC_SKIP, message=info)
            return 0, res

        rc = self._init()
        if not rc == rt.log.TEST_RC_PASS:
            res[rc] = 1
            rt.log.result(self.name, rc, message="Init failed")
            return 1, res

        for item in self.items:
            if isinstance(item, _Test):
                rc = item.run()
                res[rc] += 1
            else:
                _, ret = item.run()
                for rc in res.keys():
                    res[rc] += ret[rc]

        rc = self._deinit()
        if not rc == rt.log.TEST_RC_PASS:
            res[rc] = 1
            rt.log.result(self.name, rc, message="De-init failed")

        err = res[rt.log.TEST_RC_FAIL]
        return err, res


class Project(_Group):
    """A project is the top-level group of a project"""

    def __init__(self, path, labels=set(), testpat=None):
        """Init the project object"""
        _Group.__init__(self, self, path, None, labels, testpat)
        return

    def run(self):
        """Run the test in the project"""
        rt.log.logger.info("Start")

        rc, res = _Group.run(self)

        rt.log.logger.info("End")
        rt.log.finish(res)
        return rc
