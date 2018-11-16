Title: Fedora Packaging Notes
Date: 2018-02-25 10:00
Status: Published


## Package source

https://src.fedoraproject.org/

## Setup

```
dnf install fedora-packager make
```

## Install build dependencies**

- `dnf build-dep <spec>`

## Build

```
$ fedpkg --release f27 local
```



## References

- https://fedoraproject.org/wiki/How_to_create_an_RPM_package

