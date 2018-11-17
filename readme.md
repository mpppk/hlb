# hlb: git + hub/lab/bucket and more
hlb is a command line tool that provides unified & interactive interface to multiple git repository hosting services.

[![CircleCI](https://circleci.com/gh/mpppk/hlb/tree/master.svg?style=svg)](https://circleci.com/gh/mpppk/hlb/tree/master)
[![Build status](https://ci.appveyor.com/api/projects/status/9jw7n8ruxseys95n/branch/master?svg=true)](https://ci.appveyor.com/project/mpppk/hlb/branch/master)
[![codebeat badge](https://codebeat.co/badges/544129f2-79a9-4641-8399-f06581cd2c53)](https://codebeat.co/projects/github-com-mpppk-hlb-master)

![hlb_ibrowse.gif](https://raw.githubusercontent.com/wiki/mpppk/hlb/images/hlb_ibrowse.gif)

## Features
* Cross Platform
* Support multi git repository hosting services
* [hub](https://hub.github.com) command compatible 
* Interactive command

## Commands
### hlb browse
* `$ hlb browse`
    * Open current repository page by default browser
* `$ hlb browse issues`
    * Open issues page of current repository by browser
* `$ hlb browse issues 1` 
    * Open the page that issue ID is 1
* `$ hlb browse pull-requests` or `$ hlb browse merge-requests`
    * Open pull-requests/merge-requests page of current repository by browser
* `$ hlb browse pull-requests 1`
    * Open the page that pull-requests/merge-requests ID is 1

### hlb ibrowse (interactive browse)
![hlb_ibrowse](https://i.gyazo.com/510fe10751129f1716b3a99b1a5014ec.png)
![hlb_ibrowse.gif](https://raw.githubusercontent.com/wiki/mpppk/hlb/images/hlb_ibrowse.gif)

### hlb create
![hlb_create](https://i.gyazo.com/56d7fe0535e79819c22ec4248fcfabc4.png)
![hlb_create_and_browse.gif](https://raw.githubusercontent.com/wiki/mpppk/hlb/images/hlb_create_and_browse.gif)

### hlb init
Create config file to `~/.config/hlb/.hlb.yaml`.

### hlb add-service
Get OAuth token from git service and add to config file.

## Installation
### Homebrew
```Shell
$ brew tap mpppk/mpppk
$ brew install hlb
```

### Standalone
Download from [release page](https://github.com/mpppk/hlb/releases) and put it anywhere in your executable path.

### Source
```Shell
$ go get github.com/mpppk/hlb
```

## Update
v0.0.3 or greater has `selfupdate` command for easy updating.

![hlb_selfupdate](https://raw.githubusercontent.com/wiki/mpppk/hlb/images/hlb_selfupdate.gif)

## Authentication
authenticate infomation of hlb is stored in `~/.config/hlb/.hlb.yaml`.

### github.com & GitHub Enterprise
a. Use `hlb add-service` command
 ```Shell
 $ hlb add-service github https://github.com # or your GHE server domain
   github username: yourname
   github password:   
 ```
(Currently, add-service command only supports GitHub)

b. Add below setting to `~/.config/hlb/.hlb.yaml`  
(If file does not exist yet, execute `hlb init` first)
```yaml
services:
   - name: github.com # or your GHE server domain
     type: github
     protocol: https # or http
     oauth_token: xxxxxxxxxxxxxxxxxx
```
(oauth_token can generate from [GitHub Personal access token page](https://github.com/settings/tokens))

### gitlab.com & your GitLab Server 
Add below setting to `~/.config/hlb/.hlb.yaml`

```yaml
services:
   - name: gitlab.com # or your GitLab server domain
     type: gitlab
     protocol: https # or http
     oauth_token: xxxxxxxxxxxxxxxxxxxxx
```
(oauth_token can generate from [GitLab Personal access token page](https://gitlab.com/profile/personal_access_tokens))

## TODO
### hub compatibility
- [x] `hlb pull-request`(experimental)
- [ ] `hlb fork`
- [x] `hlb create`
- [x] `hlb browse`
- [ ] `hlb compare`
- [ ] `hlb ci-status`

### Support Services
- [x] [GitHub.com](https://github.com) / GitHub Enterprise
- [x] [GitLab.com](https://gitlab.com) / Your Own GitLab Server
- [ ] [BitBucket.org](https://bitbucket.org) / BitBucket Server
- [ ] [GitBucket](https://github.com/gitbucket/gitbucket)
- [ ] [Gogs](https://gogs.io)
- [ ] [AWS CodeCommit](https://aws.amazon.com/codecommit/)
- [ ] [GCP Cloud Source Repositories](https://cloud.google.com/source-repositories/)
