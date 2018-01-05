#!/usr/bin/env python

"""
Utility to play around with users and passwords on a Linux system
"""

from __future__ import print_function
import pwd
import argparse
import os

def read_login_defs():

    uid_min = None
    uid_max = None

    if os.path.exists('/etc/login.defs'):
        with open('/etc/login.defs') as f:
            login_data = f.readlines()
            
        for line in login_data:
            if line.startswith('UID_MIN'):
                uid_min = int(line.split()[1].strip())
            
            if line.startswith('UID_MAX'):
                uid_max = int(line.split()[1].strip())

    return uid_min, uid_max

# Get the users from /etc/passwd
def getusers(no_system=False):

    uid_min, uid_max = read_login_defs()

    if uid_min is None:
        uid_min = 1000
    if uid_max is None:
        uid_max = 60000

    users = pwd.getpwall()
    for user in users:
        if no_system:
            if user.pw_uid >= uid_min and user.pw_uid <= uid_max:
                print('{0}:{1}'.format(user.pw_name, user.pw_shell))
        else:
            print('{0}:{1}'.format(user.pw_name, user.pw_shell))

if __name__=='__main__':

    parser = argparse.ArgumentParser(description='User/Password Utility')

    parser.add_argument('--no-system', action='store_true',dest='no_system',
                        default = False, help='Specify to omit system users')

    args = parser.parse_args()
    getusers(args.no_system)
        

