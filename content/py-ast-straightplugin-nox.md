Title: A demo plugin based Python code analyser
Date: 2018-05-13 12:00
Category: python
Status: Draft

A few weeks back I wrote a analyser for [Apache Thrift IDL](https://thrift.apache.org/) in Python. We used it to enforce
some code review guidelines. When we hooked it onto [arcanist](https://secure.phabricator.com/book/phabricator/article/arcanist/) lint engine, we could give feedback to developers
at the time they were proposing a code change. The thrift parsing was done using [ptsd](https://github.com/wickman/ptsd).
The analyser was written as a single file which meant adding new rules meant changing the engine itself. I wanted to implement
a plugin based architecture for it. However, I didn't get around to do that because of other reasons.

Around the same time, I learned about [straight.plugin](http://straightplugin.readthedocs.io/en/latest/) from Nick Coghlan .
So finally, I got around to sit with it and the result is this post and a plugin based
Python code analyser with an accompanying [git repository](https://github.com/amitsaha/py_analyser).

Our analyser has two parts - the core engine and the plugins which can do various things with the code being analysed. For
the demo analyser, we will be focussed on Python classes. We will ignore everything else. And as far as the plugins
are concerned, they check if a certain condition or conditions are met by the class - in other words these are
`checkers`.

## Core analyser engine

The core analyser engine uses the [ast](https://docs.python.org/3/library/ast.html) module to create an AST node - 
basically, we call the [parse()](https://docs.python.org/3/library/ast.html#ast.parse) function
and we get back an AST Node object which we can then use to traverse through the various nodes using the
[walk](https://docs.python.org/3/library/ast.html#ast.walk) function. Here's the code for the engine at the time of
writing:

```
# analyser/main.py
...
def analyse(file_path):
    with open(file_path) as f:
        root = ast.parse(f.read())
        for node in ast.walk(root):
            if isinstance(node, ast.ClassDef):
                check_class(node)
...
```

As we walk through the tree, we check if a `node` is a Python class via `isinstance(node, ast.ClassDef)`. If it is,
we call this function `check_class` which then invokes all the checks the analyser *knows* of. So, we can write the
`check_class` function such that we have all the rules hard-coded in there or have a way to externally load the check
rules. Externally loading the rules without having to change the core engine is where `straight.plugin` comes in.

This how the `check_class` function looks like at the time of writing:

```
# analyser/main.py
def check_class(node):
    # The call() function yields a function object
    # http://straightplugin.readthedocs.io/en/latest/api.html#straight.plugin.manager.PluginManager.call
    # Not sure why..
    for p in plugins.call('check_violation', node):
        if p:
            p()
```

`plugins` is a basically a plugin registry - `straight.plugin` calls it as `PluginManager`. The `call` method
returns function objects corresponding to the function you specified. Here, I have specified `check_violation`
which expects an argument to be passed, `node`. If it finds an valid object - i.e. it found the specified
method, we call it.

How do we create the plugin registry? We use the `load` function:

```
from straight.plugin import load
plugins = load('analyser.extensions', subclasses=BaseClassCheck)
..
```

The `load` function is called with two parameters:

- The namespace for our plugins which I have set as `analyser.extensions`
- The `subclasses` kwarg specifies that we only want to load classes which are subclasses of `BaseClassCheck`.

`BaseClassCheck` is implemented as follows:


```
# 
```


## Other learnings

<class 'analyser.main.BaseClassCheck'> <class '__main__.BaseClassCheck'> False
<class 'analyser.extensions.capwords.CheckCapwords'> <class '__main__.BaseClassCheck'> False


