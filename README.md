# `dnsp`: A DNS Proxy [![wercker status](https://app.wercker.com/status/f42156d5f4e863ebe8cf0c311bd7800a/s/master "wercker status")](https://app.wercker.com/project/bykey/f42156d5f4e863ebe8cf0c311bd7800a)


### Installation

```sh
go get -u github.com/gophergala/dnsp
go install github.com/gophergala/dnsp
```


### Usage

```sh
dnsp -h
```


### Running with a non-root user

Because `dnsp` binds to port 53 by default, it requires to be run with a
privileged user on most systems. To avoid having to run `dnsp` with sudo, you
can set the `setuid` and `setgid` access right flags on the compiled
executable:

```
sudo mkdir -p /usr/local/bin
sudo cp $GOPATH/bin/dnsp
sudo chmod ug+s /usr/local/bin/dnsp
```

While `dnsp` will still run with root privileges, at least now we can run it
with a non-admin user (someone who is not in the `sudoers` group).
