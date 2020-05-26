Golfs
=====

Golfs is a Git-LFS server written in Go, that uses Google Cloud Storage to store
files. It's meant to be a replacement for Github's LFS storage so you can have
better control over where and how your files are stored.

## Setup

Here's how to get Golfs setup for yourself.

### Configure

any steps required to configure go here

### Deploy

how to deploy ( preferably makefile )

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
* **user**: 
* **github token**: 
* **url**:
* **owner**:
* **repo**:

## Future Thoughts

**Multiple backends?**
  - google cloud storage
  - amazon s3
  - others?

Right now only going to build for Google Cloud Storage, as that's all we need.

**Use database to store locks**

Need to store locks so that they can be retrieved on application
(re)start. Something like Redis should be fine, not sure SQL is required.

Don't need audit logs -- Git should be handling that already.
