#
# Copyright (C) 2016 Rolf Neugebauer <rolf.neugebauer@docker.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#

import platform

"""Miscellaneous utility functions"""


def str2labels(labelstr):
    """Convert a comma separate list of labels to two sets.

    Plain labels are placed into the labels_set set and labels
    starting with a '!' are placed in the labels_not_set set.
    """
    labels_set = set()
    labels_not_set = set()

    if not labelstr:
        return labels_set, labels_not_set

    labels = labelstr.split(',')

    for label in labels:
        l = label.strip()
        if len(l) == 0:
            continue
        if l.startswith('!'):
            labels_not_set.add(l[1:])
        else:
            labels_set.add(l)
    return labels_set, labels_not_set


#
# OS info
#
def _get_osx_info():
    """collect information for OS X based hosts"""
    info = {}
    info["os"] = "osx"

    v = platform.mac_ver()
    info["version"] = v[0]

    if v[0].startswith("10.12"):
        info["name"] = "Sierra"
    elif v[0].startswith("10.11"):
        info["name"] = "ElCapitan"
    elif v[0].startswith("10.10"):
        info["name"] = "Yosemite"
    elif v[0].startswith("10.9"):
        info["name"] = "Mavericks"
    elif v[0].startswith("10.8"):
        info["name"] = "Lion"
    elif v[0].startswith("10.7"):
        info["name"] = "Snow Leopard"
    else:
        info["name"] = "UNKNOWN"

    info["arch"] = v[2]

    return info


def _get_win_info():
    """collect information for OS X based hosts"""
    info = {}
    info["os"] = "win"

    v = platform.win32_ver()
    info["version"] = v[1]

    # Need a way to find Edition (Pro, Home etc)
    info["name"] = "Windows" + v[0] + v[2]

    a = platform.machine()
    info["arch"] = "x86_64" if a == "AMD64" else "x86_32"

    return info


def get_os_info():
    """Collect information about the local OS and stick them into a
    canonical form.  Currently supports Mac OS and Windows."""

    t = platform.system()
    if t == "Darwin":
        return _get_osx_info()
    elif t == "Windows":
        return _get_win_info()
    else:
        raise Exception("Unsupported OS")
