#!/usr/bin/env python

import os

f = os.open("/tmp/open_flags", os.O_RDONLY | os.O_CREAT | os.O_NOFOLLOW)
os.close(f)
f = os.open("/tmp/open_flags", os.O_RDONLY | os.O_SYNC)
