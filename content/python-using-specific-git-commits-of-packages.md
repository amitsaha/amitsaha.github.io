Setup.py:

Install a specific commit of a package:

```
 dependency_links=['git+https://<git-repo>@4ed6231457c244b8459037ee2224b0ef430cf766#egg=<package-name>-0'],
 ```
 
 However if a specific
 
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
