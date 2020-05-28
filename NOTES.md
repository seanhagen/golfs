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

## How It Works

- assumptions:
  - running as google cloud function
  - storing files in google cloud storage
  - storing lock information in datastore
  - github is being used for git

- required configuration:
  - lock storage, google datastore namespace
      
- optional config
  - allow locking ( if turned off will fake locking )
  - lock timeout ( default 5 min )

- routes
  - locks
    - create
    - list locks
    - list locks for verification
    - delete 
  - object/batch
    - used for upload & download requests
