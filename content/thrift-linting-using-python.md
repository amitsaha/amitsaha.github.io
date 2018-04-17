Title: Linting Apache Thrift IDL with Python (and making it work with `arcanist`)
Date: 2018-04-10 23:00
Category: python
Status: Draft

In this post, we learn about how we can setup linting for [Apache Thrift]() IDLs using Python. The 
motivation for setting up linting can be ensuring that we do catch errors before the code
chage is submitted for review, ensuring code review expectations are met, having a 
consistent code base, giving immediate feedback to developers and others.


## Parsing Thrift

The first step to creating linting rules is to of course parse the Thrift IDL. 
The [ptsd](https://github.com/wickman/ptsd) Python package works great to parse thrift IDLs.
Given a thrift definition file, you will get a `tree` object which then
includes the various objects that your thrift IDL defines:

- 


https://secure.phabricator.com/book/phabricator/article/arcanist_lint_script_and_regex/
