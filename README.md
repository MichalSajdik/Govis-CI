We really like simplicity of travis, we want to run shell script jobs on our own servers (for example for abap style check, that needs standalone server to parse abap in),
upon PR/pushes/ whatever we'd like to github. You can think of this as plugin package to server/docker image, that enables running PR checks on it, with result report back to gh, just like travis/gh pull request builder from jenkins.


After some research, we haven't found anything open source that does this. There is something simillar writen in haskell:
https://github.com/ElvishJerricco/nix-simple-ci , but it is not flexible enough, and well it is in haskell.

GH themselfes provide several guides to do this: https://developer.github.com/v3/guides/building-a-ci-server/ , they all use ruby gh client: https://github.com/octokit/octokit.rb
detailed guide for ruby: https://www.thegreatcodeadventure.com/building-a-github-ci-server/

However we'd like to do it in Go, as Go is much better for standalone distribution, and parallelism.
There are some gh clients already existing we could leverage:

https://github.com/octokit/go-octokit
https://github.com/google/go-github
Features for PoC:

- Accepts webhook request from Github
- Reads .govis.yaml file from repo
- Runs whatever is in yaml "script:"
- Reports result back to Github
- Open Source

Future features:

- standalone server with:
- results from nodes,
- running jobs view,
- deployment to new servers
- scheduling
- connection to kuberenetes cluster with worker pods

---
# Setup:

0. Generate API TOKEN in github, go to Settings -> Developer Settings ->  Personal access token. Select at least repo:status scope. Be aware, that jobs can see this API token.

## Server:

1. Create user, and folder he owns
2. Move Govis binary from releases to that folder
3. Configuration
4. Set Environment Variables:

- `GOVIS_SECRET` - WebHook Secret
- `GOVIS_BASIC_AUTH` - Personal_access_Api_token of user with set status access to the repository
- `GOVIS_PORT` - port to start Govis on - default : `8000`
- `GOVIS_SCRIPT_RIGHTS` - if unset Govis can't run `script` from yaml directly, only from `abapMerge:` inside podman
- `GOVIS_SECURE_MODE` - if set to `negative` then we don't use token for verification with github
- `GOVIS_STATUS_BASEURL` - url with port for accessing govis-ci server
- `GOVIS_SCRIPT_TIMEOUT` - number of seconds to timeout each script run - default : 300 
- `GOVIS_HOSTNAME` - hostname on which we accept webhooks (0.0.0.0 accepts all and is default)
- `GOVIS_LOG` - if set to any string then we use verbose logging otherwise we don't use it
- `GOVIS_API` - rest api of 
    - SAP : <compnay git url>/api/v3
    - Normal GitHub (default) : https://api.github.com/
    
5. Run binary in service/nohup/screen/tmux

## Github:

1. Create new webhook in github repository that targets `IP:PORT/payload/` , and has `GOVIS_SECRET` as secret, and runs on pull request event.
2. Create `.govis.yaml` file with following structure:

```
script:
    - test command
    - another test command
abapMerge:
    - command to run sap cli in podman
```
3. Create a PR with this file.
4. After result status gets back go to Settings - Branches, and add that status to required 

- you can view results at `http:/IP:PORT/govis-ci/results/`, but only when program is running

# Contribute

 - If you don't have experience with contributing read this https://github.com/firstcontributions/first-contributions/blob/master/README.md

## Conventions 
 - `go test` must be successful on windows 10 and ubuntu/linux
 - Code must be formatted with `go fmt`
 - All new functionality must have tests, which uses `testing` library. For examples you can check *_test.go files
 - Per each `*.go` file there will be `*_test.go` with tests for the `*.go` file
 - All files are in one package
 - Pull request must contain summary of new changes/features