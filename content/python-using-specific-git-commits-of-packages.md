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
