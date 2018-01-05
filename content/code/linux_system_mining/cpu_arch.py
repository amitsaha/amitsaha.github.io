#! /usr/bin/env python

""" Find the real bit architecture
"""

from __future__ import print_function

with open('/proc/cpuinfo') as f:
    for line in f:
        # Ignore the blank line separating the information between
        # details about two processing units
        if line.strip():
            if line.rstrip('\n').startswith('flags') \
                    or line.rstrip('\n').startswith('Features'):
                if 'lm' in line.rstrip('\n').split():
                    print('64-bit')
                else:
                    print('32-bit')
