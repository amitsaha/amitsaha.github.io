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

THEME='/tmp/pelican-themes/pelican-svbhack/'


STATIC_PATHS = ['images']


PLUGIN_PATHS = ['/tmp/pelican-plugins']
PLUGINS = ['code_include', 'pelican_youtube']
