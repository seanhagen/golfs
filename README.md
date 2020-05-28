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

Create a Google Cloud Service Account key, and save it to `account.json` in the
project root.

Run `make deploy`. Look for the following section:

```
httpsTrigger:
  url: https://<region>-<project>.cloudfunctions.net/GOLFS
```

Copy that URL, you'll need it later.

### Configure Git-LFS in Repo

#### First Time

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

So to generate a URL for the repo `github.com/examplecom/go-lfs-repo`

```
https://seanhagen:agithubtoken@github.com/examplecom/go-lfs-repo
```

#### After Clone

After cloning a repo that uses `golfs` to store objects, you'll probably see a
message like this:

```
Error downloading object: lock.png (a2ca9a3): Smudge error: Error downloading lock.png (a2ca9a3cabe684801034f55e1c865631f50328349c6e23fc8bb4f585bb3ebd92): [a2ca9a3cabe684801034f55e1c865631f50328349c6e23fc8bb4f585bb3ebd92] Object does not exist on the server: [404] Object does not exist on the server

Errors logged to /home/sean/tmp/lfs-testing/.git/lfs/logs/20200528T115811.722416529.log
Use `git lfs logs last` to view the log.
error: external filter 'git-lfs filter-process' failed
fatal: lock.png: smudge filter lfs failed
warning: Clone succeeded, but checkout failed.
```

That's fine! Because `golfs` hasn't been setup on this clone yet, it wasn't able
to pull the objects. It tried to pull them from the GitHub LFS storage but they
don't exist there so that failed.

Simply follow the same steps in [First Time](#first-time) to set up the
URL. Then run these commands inside the repo:

1. `git lfs pull`
2. `git add .`

And your repo should be good to go!

#### Transfering

Follow the same steps in [First Time](#first-time) to set up the git-lfs URL,
then run `git lfs push --all origin <branch>` to re-push all the objects.

### GitHub Token

To generate a token to use with `golfs`, follow these steps:

1. Log into [Github.com](https://github.com)
2. Go to the [generate new token](https://github.com/settings/tokens/new) page.
3. Fill out the name, and select the `repo` scope: ![Instructions](/imgs/new_token.jpeg)
4. Click 'Generate token'
5. Once the token has been generated, copy it -- if you don't copy it you don't
   get to see it again, so make sure you don't forget!
6. Use the token in the git-lfs URL 

This token is only used to check that the user has the correct permissions when
attempting to upload or download objects ( `write` or `admin` for upload, `read`
or `admin` for download ).
