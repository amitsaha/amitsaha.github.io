#!/usr/bin/env python
# -*- coding: utf-8 -*- #
from __future__ import unicode_literals

AUTHOR = u'Amit Saha'
SITENAME = u'Programming and Writing about it'
SITEURL = ''
PATH = 'content'
TIMEZONE = 'Australia/Brisbane'

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

#THEME='/pelican-themes/pelican-bootstrap3/'

TWITTER_USERNAME="echorand"
TWITTER_WIDGET_ID="608854659248181251"

ADDTHIS_PROFILE="ra-557910c01cb1eb5c"

STATIC_PATHS = ['images']


PLUGINS = [
    'pelican_youtube',
]
