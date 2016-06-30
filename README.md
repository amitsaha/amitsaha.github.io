```
docker build --build-arg git_username="Amit Saha" --build-arg git_useremail="amitsaha.in@gmail.com" --build-arg gid=`    id -g` --build-arg group=`id --group` --build-arg uid=`id --user` --build-arg user=`id -u -n` -t amitsaha/pelican  .
```

```
docker run -v `pwd`:/site amitsaha/pelican
```

