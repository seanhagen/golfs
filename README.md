Golfs
=====

Golfs is a Git-LFS server written in Go, that uses Google Cloud Storage to store
files. It's meant to be a replacement for Github's LFS storage so you can have
better control over where and how your files are stored. 

At some point it will have full support for locking as well, storing the locks
in Google Datastore. 

### Configure

Create `.env.yaml`, and add the following environment variables:
 
* `GOLFS_DS_NAMESPACE` -- the Google Cloud Datastore namespace you wish to use
* `GOLFS_BUCKET` -- the Google Cloud Storage bucket you wish to use
* `GCP_PROJECT` -- the Google Cloud project identifier

### Deploy

Run `make deploy`. Look for the following section:

```
httpsTrigger:
  url: https://<region>-<project>.cloudfunctions.net/GOLFS
```

Copy that URL, you'll need it later.

### Configure Git-LFS in Repo

```
git config lfs.url https://<user>:<github token>@<url>/<owner>/<repo>
```

Alternatively edit .git/config to add:

```
[lfs]
  url = https://<user>:<github token>@<url>/<owner>/<repo>
```

What to use to replace:
* **user**: Your GitHub username
* **github token**: See `GitHub Token` below
* **url**: The URL you copied after running `make deploy` above
* **owner**: The repo owner
* **repo**: The repo name

So, for this repo, I'd construct a URL something like this:

```
https://sean:agithubtoken@github.com/seanhagen/golfs
```

#### GitHub Token

To generate a token to use with `golfs`, follow these steps:

1. Log into [Github.com](https://github.com)
2. Go to the [generate new token](https://github.com/settings/tokens/new) page.
3. Fill out the name, and select the `repo` scope: ![Instructions](/imgs/new_token.jpeg)
4. Click 'Generate token'
5. Once the token has been generated, copy it -- if you don't copy it you don't
   get to see it again, so make sure you don't forget!
6. Use the token in the git-lfs URL 


