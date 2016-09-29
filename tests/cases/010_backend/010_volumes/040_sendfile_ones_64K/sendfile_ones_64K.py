#!/usr/bin/env python

import os

os.system("tr '\\0' '\\xff' < /dev/zero | dd bs=1k of=/tmp/ones_64K count=64")
inf = os.open("/tmp/ones_64K", os.O_RDONLY)
outf = os.open("/tmp/sendfile_64K", os.O_CREAT | os.O_TRUNC | os.O_WRONLY)
moved = 0
while moved < 1 << 16:
    moved += os.sendfile(outf, inf, 0, 1 << 16)
    print("%d bytes moved" % moved)

exit_status = os.system("diff /tmp/ones_64K /tmp/sendfile_64K")
print("diff exited %d" % (exit_status >> 8))
exit(exit_status != 0)
