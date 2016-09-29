#!/usr/bin/env python

import os, tempfile

dir = tempfile.mkdtemp()
os.mkdir(dir + "/subdir")
os.system("touch %s/foo" % dir)
os.rename(dir + "/foo", dir + "/subdir/foo")
