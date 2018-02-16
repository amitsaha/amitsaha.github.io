Setup.py:

Install a specific commit of a package:

```
 dependency_links=['git+https://<git-repo>@4ed6231457c244b8459037ee2224b0ef430cf766#egg=<package-name>-0'],
 ```
 
However if the package is already in `pypi`, we have a problem. So, we fool `pip` like, so:

```
import setuptools

setuptools.setup(
    name='my_awesome_cli',
    version='0.1',
    description='My Awesome CLI',
    packages=setuptools.find_packages(),
    # I fool `pip` by specifying the version number which
    # is greater than the one released in PyPi and force
    # it to look at the dependency_links where i wrongly specify
    # that i have a version which is greater than 0.1.2
    install_requires='fire>0.1.2',
    dependency_links=[
        'git+https://github.com/google/python-fire.git@9bff9d01ce16589201f57ffef27ea84744951c11#egg=fire-0.1.2.1',
    ],
    entry_points={
        'console_scripts': [
            'my-awesome-cli=my_awesome_cli.main:main'
        ],
    }
)

```

Now, if we install `pip install . --process-dependency-links`, we will see:

```
  Could not find a tag or branch '9bff9d01ce16589201f57ffef27ea84744951c11', assuming commit.
  Requested fire>0.1.2 from git+https://github.com/google/python-fire.git@9bff9d01ce16589201f57ffef27ea84744951c11#egg=fire-0.1.2.1 (from my-awesome-cli==0.1), but installing version None
```
 
 To then distribute this to  PyPI, we need to make sure that we distribute this as a source tarball, [not a wheel](https://github.com/pypa/pip/issues/3172):
 
 ```
 $ python setup.py sdist
 $ TWINE_REPOSITORY_URL=https://test.pypi.org/legacy/ TWINE_USERNAME=echorand TWINE_PASSWORD="secret" twine upload dist/*

 ```
 
 Once we have done that, we can install it, like so:
 
 ```
$ pip install my-awesome-cli==0.2 --process-dependency-links -i https://test.pypi.org/simple/

Collecting my-awesome-cli==0.2
  Downloading https://test-files.pythonhosted.org/packages/5b/ac/7c307602de47a832a7f61094515cf3f49b0d80db135c56e960133c439dff/my_awesome_cli-0.2.tar.gz
  DEPRECATION: Dependency Links processing has been deprecated and will be removed in a future release.
Collecting fire>0.1.2 (from my-awesome-cli==0.2)
  Cloning https://github.com/google/python-fire.git (to 9bff9d01ce16589201f57ffef27ea84744951c11) to /tmp/pip-build-SykxjY/fire
  Could not find a tag or branch '9bff9d01ce16589201f57ffef27ea84744951c11', assuming commit.
  Requested fire>0.1.2 from git+https://github.com/google/python-fire.git@9bff9d01ce16589201f57ffef27ea84744951c11#egg=fire-0.1.2.1 (from my-awesome-cli==0.2), but installing version None
Collecting six (from fire>0.1.2->my-awesome-cli==0.2)
  Downloading https://test-files.pythonhosted.org/packages/b3/b2/238e2590826bfdd113244a40d9d3eb26918bd798fc187e2360a8367068db/six-1.10.0.tar.gz
Building wheels for collected packages: my-awesome-cli, six
  Running setup.py bdist_wheel for my-awesome-cli ... done
  Stored in directory: /home/asaha/.cache/pip/wheels/bd/92/e5/b7aa57da22bddad31418b524cc6d64f2be1fbe480c4cc69bbf
  Running setup.py bdist_wheel for six ... done
  Stored in directory: /home/asaha/.cache/pip/wheels/5b/14/2e/85da3578e30a577e4edf5b477cbc6b8a3bdc3228d211bf9fef
Successfully built my-awesome-cli six
Installing collected packages: six, fire, my-awesome-cli
  Running setup.py install for fire ... done
Successfully installed fire-0.1.2 my-awesome-cli-0.2 six-1.10.0
(fire2) asaha@asaha-desktop:~$ my-awesome-cli
Type:        Calculator
String form: <my_awesome_cli.main.Calculator object at 0x7feecae69850>
Docstring:   A simple calculator class.

Usage:       my-awesome-cli
             my-awesome-cli double

```

 
 ```
 pip install --process-dependency-links
 ```
 
 
# requirements.txt

```
git+https://<git repo>@master
```

pip install -r requirements.txt


## Helpful links

- https://caremad.io/posts/2013/07/setup-vs-requirement/
- https://yuji.wordpress.com/2011/04/11/pip-install-specific-commit-from-git-repository/
