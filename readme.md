# hlb: git + hub/lab/bucket and more
hlb is a command line tool that provides a unified interface to multiple git repository hosting services.

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

## Authentication
authenticate infomation of hlb is stored in `~/.hlb.yaml`.

### github.com & GitHub Enterprise
a. Use `hlb add-service` command
 ```Shell
 $ hlb add-service github https://github.com # or your GHE server domain
   github username: yourname
   github password:   
 ```
(Currently, add-service command only supports GitHub)

b. Add below setting to `~/.hlb.yaml`  
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
Add below setting to `~/.hlb.yaml`

```yaml
services:
   - name: gitlab.com # or your GitLab server domain
     type: gitlab
     protocol: https # or http
     oauth_token: xxxxxxxxxxxxxxxxxxxxx
```
(oauth_token can generate from [GitLab Personal access token page](https://gitlab.com/profile/personal_access_tokens))

## Commands
### hlb init
Create config file(`~/.hlb.yaml`).

### hlb add-service
Get OAuth token from git service and add to config file.

### hlb browse
#### Browse page
* `$ hlb browse`
    * Open current repository page by default browser
* `$ hlb browse issues`
    * Open issues page of current repository by browser
* `$ hlb browse issues 1` 
    * Open the page that issue ID is 1
* `$ hlb browse pull-requests`
    * Open pull requests page of current repository by browser
* `$ hlb browse pull-requests 1`
    * Open the page that pull request ID is 1

### hlb list
* list issues
* list pull-requests

## TODO
### Implement commands
- [ ] `hlb pull-request`
- [ ] `hlb fork`
- [ ] `hlb create`
- [x] `hlb browse`
- [ ] `hlb compare`
- [x] `hlb list`
- [ ] `hlb ci-status`

### Support Services
- [x] [GitHub.com](https://github.com) / GitHub Enterprise
- [x] [GitLab.com](https://gitlab.com) / Your Own GitLab Server
- [ ] [BitBucket.org](https://bitbucket.org) / BitBucket Server
- [ ] [GitBucket](https://github.com/gitbucket/gitbucket)
- [ ] [Gogs](https://gogs.io)
- [ ] [AWS CodeCommit](https://aws.amazon.com/codecommit/)
- [ ] [GCP Cloud Source Repositories](https://cloud.google.com/source-repositories/)
