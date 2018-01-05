#!/usr/bin/env python

"""
Python interface to the /proc file system.
Although this can be used as a replacement for cat /proc/... on the command line,
its really aimed to be an interface to /proc for other Python programs.

As long as the object you are looking for exists in /proc
and is readable (you have permission and if you are reading a file,
its contents are alphanumeric, this program will find it). If its a
directory, it will return a list of all the files in that directory
(and its sub-dirs) which you can then read using the same function.


Example usage:

Read /proc/cpuinfo:

$ ./readproc.py proc.cpuinfo

Read /proc/meminfo:

$ ./readproc.py proc.meminfo

Read /proc/cmdline:

$ ./readproc.py proc.cmdline

Read /proc/1/cmdline:

$ ./readproc.py proc.1.cmdline

Read /proc/net/dev:

$ ./readproc.py proc.net.dev

Comments/Suggestions:

Amit Saha <@echorand>
<http://echorand.me>

"""

from __future__ import print_function
import os
import sys
import re

def toitem(path):
    """ Convert /foo/bar to foo.bar """
    path = path.lstrip('/').replace('/','.')
    return path

def todir(item):
    """ Convert foo.bar to /foo/bar"""
    # TODO: breaks if there is a directory whose name is foo.bar (for
    # eg. conf.d/), but we don't have to worry as long as we are using
    # this for reading /proc
    return '/' + item.replace('.','/')

def readproc(item):
    """ 
    Resolves proc.foo.bar items to /proc/foo/bar and returns the
    appropriate data.
    1. If its a file, simply return the lines in this file as a list
    2. If its a directory, return the files in this directory in the
    proc.foo.bar style as a list, so that this function can then be
    called to retrieve the contents
    """
    item = todir(item) 

    if not os.path.exists(item):
        return 'Non-existent object'
    
    if os.path.isfile(item):
        # its a little tricky here. We don't want to read huge binary
        # files and return the contents. We will probably not need it
        # in the usual case.
        # utilities like 'file' on Linux and the Python interface to
        # libmagic are useless when it comes to files in /proc for
        # detecting the mime type, since the these are not on-disk
        # files. 
        # Searching, i find this solution which seems to be a
        # reasonable assumption. If we find a '\0' in the first 1024
        # bytes of a file, we declare it as binary and return an empty string
        # however, some of the files in /proc which contain text may
        # also contain the null byte as a constituent character.
        # Hence, I use a RE expression that matches against any
        # combination of alphanumeric characters
        # If any of these conditions suffice, we read the file's contents
        
        pattern = re.compile('\w*')
        try:
            with open(item) as f:
                chunk = f.read(1024)
                if '\0' not in chunk or pattern.match(chunk) is not None:
                    f.seek(0)
                    data = f.readlines()
                    return data
                else:
                    return '{0} is binary'.format(item)
        except IOError:
            return 'Error reading object'

    if os.path.isdir(item):
        data = []
        for dir_path, dir_name, files in os.walk(item):
            for file in files:
                data.append(toitem(os.path.join(dir_path, file)))
        return data

if __name__=='__main__':
    
    if len(sys.argv)>1:
        data = readproc(sys.argv[1])
    else:
        data = readproc('proc')
    
    if type(data) == list:
        for line in data:
            print(line)
    else:
        print(data)
