#! /usr/bin/env python
""" print out the /proc/cpuinfo
    file
"""

from __future__ import print_function

with open('/proc/cpuinfo') as f:
    for line in f:
        print(line.rstrip('\n'))
