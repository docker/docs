#!/usr/bin/env python

# Port of Gavin Brock's Perl kcpassword generator to Python, by Tom Taylor
# <tom@tomtaylor.co.uk>.
# Perl version: http://www.brock-family.org/gavin/perl/kcpassword.html

import sys
import os

def kcpassword(passwd):
    # The magic 11 bytes - these are just repeated
    # 0x7D 0x89 0x52 0x23 0xD2 0xBC 0xDD 0xEA 0xA3 0xB9 0x1F
    key = [125,137,82,35,210,188,221,234,163,185,31]
    key_len = len(key)

    passwd = [ord(x) for x in list(passwd)]
    # pad passwd length out to an even multiple of key length
    r = len(passwd) % key_len
    if (r > 0):
        passwd = passwd + [0] * (key_len - r)

    for n in range(0, len(passwd), len(key)):
        ki = 0
        for j in range(n, min(n+len(key), len(passwd))):
            passwd[j] = passwd[j] ^ key[ki]
            ki += 1

    passwd = [chr(x) for x in passwd]
    return "".join(passwd)

if __name__ == "__main__":
    passwd = kcpassword(sys.argv[1])
    fd = os.open('/etc/kcpassword', os.O_WRONLY | os.O_CREAT, 0o600)
    file = os.fdopen(fd, 'w')
    file.write(passwd)
    file.close()
