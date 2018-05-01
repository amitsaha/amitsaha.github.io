#!/usr/bin/env python
# -*- coding: utf-8 -*- #
from __future__ import unicode_literals

import os

AUTHOR = u'Amit Saha'
SITENAME = u'Exploring Software and Writing about it'
SITEURL = 'http://echorand.me'
PATH = 'content'
TIMEZONE = 'Australia/Sydney'

DEFAULT_LANG = u'en'

# Feed generation is usually not desired when developing
FEED_ALL_ATOM = None
CATEGORY_FEED_ATOM = None
TRANSLATION_FEED_ATOM = None
AUTHOR_FEED_ATOM = None
AUTHOR_FEED_RSS = None

DEFAULT_PAGINATION = 10
DISPLAY_CATEGORIES_ON_MENU = False
DISPLAY_TAGS_ON_SIDEBAR = False

THEME='/tmp/pelican-svbhack/'
USER_LOGO_URL=SITEURL+'/images/logo.png'
TAGLINE=SITENAME

GOOGLE_ANALYTICS='UA-108901610-1'

STATIC_PATHS = ['images']


PLUGIN_PATHS = ['/tmp/pelican-plugins', 'plugins']
PLUGINS = [
    'code_include',
    'pelican_youtube',
    'pelican_gist',
    'share_post',
    'pelican_notes',
]

MD_EXTENSIONS =  [ 'toc', 'codehilite','extra']

NOTE_PATHS = 'notes'
NOTE_SAVE_AS = 'notes/{slug}.html'
NOTE_URL = 'notes/{slug}.html'

ARTICLE_EXCLUDES = [NOTE_PATHS]
PAGE_EXCLUDES = [NOTE_PATHS]

def basename(path):
    return os.path.basename(path)
JINJA_FILTERS = {
    'basename': basename,
}

