Git-LFS Notes
=============

## Authentication

Git-LFS uses HTTP Basic Auth 

### SSH

ssh git@git-server.com git-lfs-authenticate foo/bar.git download

## Building A URL

By default, Git LFS will append .git/info/lfs to the end of a Git remote url to
build the LFS server URL it will use:

Git Remote: https://git-server.com/foo/bar
LFS Server: https://git-server.com/foo/bar.git/info/lfs

Git Remote: https://git-server.com/foo/bar.git
LFS Server: https://git-server.com/foo/bar.git/info/lfs

Git Remote: git@git-server.com:foo/bar.git
LFS Server: https://git-server.com/foo/bar.git/info/lfs

Git Remote: ssh://git-server.com/foo/bar.git
LFS Server: https://git-server.com/foo/bar.git/info/lfs

## 
