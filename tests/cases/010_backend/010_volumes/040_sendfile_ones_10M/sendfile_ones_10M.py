#!/usr/bin/env python

import os

os.system("tr '\\0' '\\xff' < /dev/zero | dd bs=1k of=/tmp/ones_10M count=10240")
inf = os.open("/tmp/ones_10M", os.O_RDONLY)
outf = os.open("/tmp/sendfile_10M", os.O_CREAT | os.O_TRUNC | os.O_WRONLY)
moved = 0
while moved < 10*1024*1024:
    moved += os.sendfile(outf, inf, 0, 1 << 24)
    print("%d bytes moved" % moved)

exit_status = os.system("diff /tmp/ones_10M /tmp/sendfile_10M")
print("diff exited %d" % (exit_status >> 8))
exit(exit_status != 0)
