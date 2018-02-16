Setup.py:

Install a specific commit of a package:

```
 dependency_links=['git+https://<git-repo>@4ed6231457c244b8459037ee2224b0ef430cf766#egg=<package-name>-0'],
 ```
 
 However if the package is already in `pypi`, we have a problem. So, we fool `pip` like, so:

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
