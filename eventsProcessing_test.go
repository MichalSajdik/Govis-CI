package main

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestSuiteForEventsProcessing struct {
	suite.Suite
}

func (t *TestSuiteForEventsProcessing) SetupSuite() {
}

func TestEventsProcessingTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuiteForEventsProcessing))
}

// check that payload endpoint return status 200 for ping event
func (t *TestSuiteForEventsProcessing) TestPingEventResponse() {

	GithubSecret = "secret"

	req, err := http.NewRequest("POST", "/payload", bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		t.FailNow(err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("X-GitHub-Event", "ping")
	//x-hub-signature is HMAC sha1 : https://gchq.github.io/CyberChef/#recipe=HMAC(%7B'option':'UTF8','string':'secret'%7D,'SHA1')
	req.Header.Set("X-Hub-Signature", "sha1=5d61605c3feea9799210ddcb71307d4ba264225f")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(EventsProcessing)

	// Our handlers satisfy http.IndexHandler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.FailNow("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func (t *TestSuiteForEventsProcessing) TestParsePushEvent() {
	payloadJSON := PayloadFromPushEvent
	testCommit := Commit{
		Owner:       "I516366",
		RepoName:    "testRepositoryForGovis",
		IdOfCommit:  "0f1f5ad2fe1c602c36cf1f217a6cd6e8adbecba7",
		StatusesUrl: "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/statuses/0f1f5ad2fe1c602c36cf1f217a6cd6e8adbecba7",
	}
	htmlURL, branchName, commit, err := parsePushEvent([]byte(payloadJSON))
	t.Equal("<compnay git url>/I516366/testRepositoryForGovis", htmlURL)
	t.Equal("master", branchName)
	t.Equal(testCommit, commit)
	t.Equal(err, nil, "parsePushEvent returned error")
}

func (t *TestSuiteForEventsProcessing) TestParsePullRequestEvent() {
	payloadJSON := PayloadFromPullRequestEvent
	testCommit := Commit{
		Owner:       "",
		RepoName:    "Govis-CI",
		IdOfCommit:  "",
		StatusesUrl: "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/statuses/",
	}
	htmlURL, branchName, commit, err := parsePushEvent([]byte(payloadJSON))
	t.Equal("<compnay git url>/FXUBRQ-QE/Govis-CI", htmlURL)
	t.Equal("", branchName)
	t.Equal(testCommit, commit)
	t.Equal(err, nil, "parsePullRequestEvent returned error")
}

// PayloadFromPushEvent && PayloadFromPullRequestEvent examples to test

const PayloadFromPushEvent = `{
"ref": "refs/heads/master",
"before": "6cc9bc24a653c0ae77e4e4fcd616bf17cf9aa955",
"after": "0f1f5ad2fe1c602c36cf1f217a6cd6e8adbecba7",
"created": false,
"deleted": false,
"forced": false,
"base_ref": null,
"compare": "<compnay git url>/I516366/testRepositoryForGovis/compare/6cc9bc24a653...0f1f5ad2fe1c",
"commits": [
{
"id": "0f1f5ad2fe1c602c36cf1f217a6cd6e8adbecba7",
"tree_id": "cd58cf787f653301dd26cf315b1c7fd996b8d7f9",
"distinct": true,
"message": "Update tmp.txt",
"timestamp": "2019-10-10T13:17:35+02:00",
"url": "<compnay git url>/I516366/testRepositoryForGovis/commit/0f1f5ad2fe1c602c36cf1f217a6cd6e8adbecba7",
"author": {
"name": "Michal Sajdik",
"email": "michal.sajdik@sap.com",
"username": "I516366"
},
"committer": {
"name": "GitHub Enterprise",
"email": "noreply-github@sap.com"
},
"added": [

],
"removed": [

],
"modified": [
"tmp.txt"
]
}
],
"head_commit": {
"id": "0f1f5ad2fe1c602c36cf1f217a6cd6e8adbecba7",
"tree_id": "cd58cf787f653301dd26cf315b1c7fd996b8d7f9",
"distinct": true,
"message": "Update tmp.txt",
"timestamp": "2019-10-10T13:17:35+02:00",
"url": "<compnay git url>/I516366/testRepositoryForGovis/commit/0f1f5ad2fe1c602c36cf1f217a6cd6e8adbecba7",
"author": {
"name": "Michal Sajdik",
"email": "michal.sajdik@sap.com",
"username": "I516366"
},
"committer": {
"name": "GitHub Enterprise",
"email": "noreply-github@sap.com"
},
"added": [

],
"removed": [

],
"modified": [
"tmp.txt"
]
},
"repository": {
"id": 284779,
"node_id": "MDEwOlJlcG9zaXRvcnkyODQ3Nzk=",
"name": "testRepositoryForGovis",
"full_name": "I516366/testRepositoryForGovis",
"private": false,
"owner": {
"name": "I516366",
"email": "michal.sajdik@sap.com",
"login": "I516366",
"id": 51335,
"node_id": "MDQ6VXNlcjUxMzM1",
"avatar_url": "<compnay git url>/avatars/u/51335?",
"gravatar_id": "",
"url": "<compnay git url>/api/v3/users/I516366",
"html_url": "<compnay git url>/I516366",
"followers_url": "<compnay git url>/api/v3/users/I516366/followers",
"following_url": "<compnay git url>/api/v3/users/I516366/following{/other_user}",
"gists_url": "<compnay git url>/api/v3/users/I516366/gists{/gist_id}",
"starred_url": "<compnay git url>/api/v3/users/I516366/starred{/owner}{/repo}",
"subscriptions_url": "<compnay git url>/api/v3/users/I516366/subscriptions",
"organizations_url": "<compnay git url>/api/v3/users/I516366/orgs",
"repos_url": "<compnay git url>/api/v3/users/I516366/repos",
"events_url": "<compnay git url>/api/v3/users/I516366/events{/privacy}",
"received_events_url": "<compnay git url>/api/v3/users/I516366/received_events",
"type": "User",
"site_admin": false
},
"html_url": "<compnay git url>/I516366/testRepositoryForGovis",
"description": null,
"fork": false,
"url": "<compnay git url>/I516366/testRepositoryForGovis",
"forks_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/forks",
"keys_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/keys{/key_id}",
"collaborators_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/collaborators{/collaborator}",
"teams_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/teams",
"hooks_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/hooks",
"issue_events_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/issues/events{/number}",
"events_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/events",
"assignees_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/assignees{/user}",
"branches_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/branches{/branch}",
"tags_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/tags",
"blobs_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/git/blobs{/sha}",
"git_tags_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/git/tags{/sha}",
"git_refs_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/git/refs{/sha}",
"trees_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/git/trees{/sha}",
"statuses_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/statuses/{sha}",
"languages_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/languages",
"stargazers_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/stargazers",
"contributors_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/contributors",
"subscribers_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/subscribers",
"subscription_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/subscription",
"commits_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/commits{/sha}",
"git_commits_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/git/commits{/sha}",
"comments_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/comments{/number}",
"issue_comment_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/issues/comments{/number}",
"contents_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/contents/{+path}",
"compare_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/compare/{base}...{head}",
"merges_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/merges",
"archive_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/{archive_format}{/ref}",
"downloads_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/downloads",
"issues_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/issues{/number}",
"pulls_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/pulls{/number}",
"milestones_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/milestones{/number}",
"notifications_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/notifications{?since,all,participating}",
"labels_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/labels{/name}",
"releases_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/releases{/id}",
"deployments_url": "<compnay git url>/api/v3/repos/I516366/testRepositoryForGovis/deployments",
"created_at": 1561561905,
"updated_at": "2019-09-27T05:39:57Z",
"pushed_at": 1570706255,
"git_url": "git://<compnay ssh git url>/I516366/testRepositoryForGovis.git",
"ssh_url": "git@<compnay ssh git url>:I516366/testRepositoryForGovis.git",
"clone_url": "<compnay git url>/I516366/testRepositoryForGovis.git",
"svn_url": "<compnay git url>/I516366/testRepositoryForGovis",
"homepage": null,
"size": 53,
"stargazers_count": 0,
"watchers_count": 0,
"language": "Go",
"has_issues": true,
"has_projects": true,
"has_downloads": true,
"has_wiki": true,
"has_pages": false,
"forks_count": 1,
"mirror_url": null,
"archived": false,
"disabled": false,
"open_issues_count": 1,
"license": null,
"forks": 1,
"open_issues": 1,
"watchers": 0,
"default_branch": "master",
"stargazers": 0,
"master_branch": "master"
},
"pusher": {
"name": "I516366",
"email": "michal.sajdik@sap.com"
},
"enterprise": {
"id": 1,
"slug": "sap-se",
"name": "SAP SE",
"node_id": "MDg6QnVzaW5lc3Mx",
"avatar_url": "<compnay git url>/avatars/b/1?",
"description": null,
"website_url": null,
"html_url": "<compnay git url>/businesses/sap-se",
"created_at": "2019-03-16T05:31:15Z",
"updated_at": "2019-03-16T05:31:15Z"
},
"sender": {
"login": "I516366",
"id": 51335,
"node_id": "MDQ6VXNlcjUxMzM1",
"avatar_url": "<compnay git url>/avatars/u/51335?",
"gravatar_id": "",
"url": "<compnay git url>/api/v3/users/I516366",
"html_url": "<compnay git url>/I516366",
"followers_url": "<compnay git url>/api/v3/users/I516366/followers",
"following_url": "<compnay git url>/api/v3/users/I516366/following{/other_user}",
"gists_url": "<compnay git url>/api/v3/users/I516366/gists{/gist_id}",
"starred_url": "<compnay git url>/api/v3/users/I516366/starred{/owner}{/repo}",
"subscriptions_url": "<compnay git url>/api/v3/users/I516366/subscriptions",
"organizations_url": "<compnay git url>/api/v3/users/I516366/orgs",
"repos_url": "<compnay git url>/api/v3/users/I516366/repos",
"events_url": "<compnay git url>/api/v3/users/I516366/events{/privacy}",
"received_events_url": "<compnay git url>/api/v3/users/I516366/received_events",
"type": "User",
"site_admin": false
}
}`

const PayloadFromPullRequestEvent = `{
  "action": "review_requested",
  "number": 100,
  "pull_request": {
    "url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/pulls/100",
    "id": 1786716,
    "node_id": "MDExOlB1bGxSZXF1ZXN0MTc4NjcxNg==",
    "html_url": "<compnay git url>/FXUBRQ-QE/Govis-CI/pull/100",
    "diff_url": "<compnay git url>/FXUBRQ-QE/Govis-CI/pull/100.diff",
    "patch_url": "<compnay git url>/FXUBRQ-QE/Govis-CI/pull/100.patch",
    "issue_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/issues/100",
    "number": 100,
    "state": "open",
    "locked": false,
    "title": "Testing uses testify framework",
    "user": {
      "login": "I516366",
      "id": 51335,
      "node_id": "MDQ6VXNlcjUxMzM1",
      "avatar_url": "<compnay git url>/avatars/u/51335?",
      "gravatar_id": "",
      "url": "<compnay git url>/api/v3/users/I516366",
      "html_url": "<compnay git url>/I516366",
      "followers_url": "<compnay git url>/api/v3/users/I516366/followers",
      "following_url": "<compnay git url>/api/v3/users/I516366/following{/other_user}",
      "gists_url": "<compnay git url>/api/v3/users/I516366/gists{/gist_id}",
      "starred_url": "<compnay git url>/api/v3/users/I516366/starred{/owner}{/repo}",
      "subscriptions_url": "<compnay git url>/api/v3/users/I516366/subscriptions",
      "organizations_url": "<compnay git url>/api/v3/users/I516366/orgs",
      "repos_url": "<compnay git url>/api/v3/users/I516366/repos",
      "events_url": "<compnay git url>/api/v3/users/I516366/events{/privacy}",
      "received_events_url": "<compnay git url>/api/v3/users/I516366/received_events",
      "type": "User",
      "site_admin": false
    },
    "body": "",
    "created_at": "2019-10-14T14:19:50Z",
    "updated_at": "2019-10-15T06:40:43Z",
    "closed_at": null,
    "merged_at": null,
    "merge_commit_sha": "bd9ce97e4223570443201a58a469a35681512c61",
    "assignee": null,
    "assignees": [

    ],
    "requested_reviewers": [
      {
        "login": "i337562",
        "id": 21508,
        "node_id": "MDQ6VXNlcjIxNTA4",
        "avatar_url": "<compnay git url>/avatars/u/21508?",
        "gravatar_id": "",
        "url": "<compnay git url>/api/v3/users/i337562",
        "html_url": "<compnay git url>/i337562",
        "followers_url": "<compnay git url>/api/v3/users/i337562/followers",
        "following_url": "<compnay git url>/api/v3/users/i337562/following{/other_user}",
        "gists_url": "<compnay git url>/api/v3/users/i337562/gists{/gist_id}",
        "starred_url": "<compnay git url>/api/v3/users/i337562/starred{/owner}{/repo}",
        "subscriptions_url": "<compnay git url>/api/v3/users/i337562/subscriptions",
        "organizations_url": "<compnay git url>/api/v3/users/i337562/orgs",
        "repos_url": "<compnay git url>/api/v3/users/i337562/repos",
        "events_url": "<compnay git url>/api/v3/users/i337562/events{/privacy}",
        "received_events_url": "<compnay git url>/api/v3/users/i337562/received_events",
        "type": "User",
        "site_admin": false
      }
    ],
    "requested_teams": [

    ],
    "labels": [

    ],
    "milestone": null,
    "commits_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/pulls/100/commits",
    "review_comments_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/pulls/100/comments",
    "review_comment_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/pulls/comments{/number}",
    "comments_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/issues/100/comments",
    "statuses_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/statuses/bb8eec32ab0f9a00564fd7e1142736ee87a8a5b2",
    "head": {
      "label": "FXUBRQ-QE:upgrade/testing/testify",
      "ref": "upgrade/testing/testify",
      "sha": "bb8eec32ab0f9a00564fd7e1142736ee87a8a5b2",
      "user": {
        "login": "FXUBRQ-QE",
        "id": 19402,
        "node_id": "MDEyOk9yZ2FuaXphdGlvbjE5NDAy",
        "avatar_url": "<compnay git url>/avatars/u/19402?",
        "gravatar_id": "",
        "url": "<compnay git url>/api/v3/users/FXUBRQ-QE",
        "html_url": "<compnay git url>/FXUBRQ-QE",
        "followers_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/followers",
        "following_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/following{/other_user}",
        "gists_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/gists{/gist_id}",
        "starred_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/starred{/owner}{/repo}",
        "subscriptions_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/subscriptions",
        "organizations_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/orgs",
        "repos_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/repos",
        "events_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/events{/privacy}",
        "received_events_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/received_events",
        "type": "Organization",
        "site_admin": false
      },
      "repo": {
        "id": 284943,
        "node_id": "MDEwOlJlcG9zaXRvcnkyODQ5NDM=",
        "name": "Govis-CI",
        "full_name": "FXUBRQ-QE/Govis-CI",
        "private": false,
        "owner": {
          "login": "FXUBRQ-QE",
          "id": 19402,
          "node_id": "MDEyOk9yZ2FuaXphdGlvbjE5NDAy",
          "avatar_url": "<compnay git url>/avatars/u/19402?",
          "gravatar_id": "",
          "url": "<compnay git url>/api/v3/users/FXUBRQ-QE",
          "html_url": "<compnay git url>/FXUBRQ-QE",
          "followers_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/followers",
          "following_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/following{/other_user}",
          "gists_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/gists{/gist_id}",
          "starred_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/starred{/owner}{/repo}",
          "subscriptions_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/subscriptions",
          "organizations_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/orgs",
          "repos_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/repos",
          "events_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/events{/privacy}",
          "received_events_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/received_events",
          "type": "Organization",
          "site_admin": false
        },
        "html_url": "<compnay git url>/FXUBRQ-QE/Govis-CI",
        "description": null,
        "fork": false,
        "url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI",
        "forks_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/forks",
        "keys_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/keys{/key_id}",
        "collaborators_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/collaborators{/collaborator}",
        "teams_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/teams",
        "hooks_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/hooks",
        "issue_events_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/issues/events{/number}",
        "events_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/events",
        "assignees_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/assignees{/user}",
        "branches_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/branches{/branch}",
        "tags_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/tags",
        "blobs_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/git/blobs{/sha}",
        "git_tags_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/git/tags{/sha}",
        "git_refs_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/git/refs{/sha}",
        "trees_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/git/trees{/sha}",
        "statuses_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/statuses/{sha}",
        "languages_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/languages",
        "stargazers_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/stargazers",
        "contributors_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/contributors",
        "subscribers_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/subscribers",
        "subscription_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/subscription",
        "commits_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/commits{/sha}",
        "git_commits_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/git/commits{/sha}",
        "comments_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/comments{/number}",
        "issue_comment_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/issues/comments{/number}",
        "contents_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/contents/{+path}",
        "compare_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/compare/{base}...{head}",
        "merges_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/merges",
        "archive_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/{archive_format}{/ref}",
        "downloads_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/downloads",
        "issues_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/issues{/number}",
        "pulls_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/pulls{/number}",
        "milestones_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/milestones{/number}",
        "notifications_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/notifications{?since,all,participating}",
        "labels_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/labels{/name}",
        "releases_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/releases{/id}",
        "deployments_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/deployments",
        "created_at": "2019-06-27T07:26:55Z",
        "updated_at": "2019-10-03T08:05:35Z",
        "pushed_at": "2019-10-15T06:19:34Z",
        "git_url": "git://<compnay ssh git url>/FXUBRQ-QE/Govis-CI.git",
        "ssh_url": "git@<compnay ssh git url>:FXUBRQ-QE/Govis-CI.git",
        "clone_url": "<compnay git url>/FXUBRQ-QE/Govis-CI.git",
        "svn_url": "<compnay git url>/FXUBRQ-QE/Govis-CI",
        "homepage": null,
        "size": 234,
        "stargazers_count": 2,
        "watchers_count": 2,
        "language": "Go",
        "has_issues": true,
        "has_projects": true,
        "has_downloads": true,
        "has_wiki": true,
        "has_pages": false,
        "forks_count": 2,
        "mirror_url": null,
        "archived": false,
        "disabled": false,
        "open_issues_count": 11,
        "license": null,
        "forks": 2,
        "open_issues": 11,
        "watchers": 2,
        "default_branch": "master"
      }
    },
    "base": {
      "label": "FXUBRQ-QE:master",
      "ref": "master",
      "sha": "6229f9739d140d6a3579a2c6a03f13810b71cf22",
      "user": {
        "login": "FXUBRQ-QE",
        "id": 19402,
        "node_id": "MDEyOk9yZ2FuaXphdGlvbjE5NDAy",
        "avatar_url": "<compnay git url>/avatars/u/19402?",
        "gravatar_id": "",
        "url": "<compnay git url>/api/v3/users/FXUBRQ-QE",
        "html_url": "<compnay git url>/FXUBRQ-QE",
        "followers_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/followers",
        "following_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/following{/other_user}",
        "gists_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/gists{/gist_id}",
        "starred_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/starred{/owner}{/repo}",
        "subscriptions_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/subscriptions",
        "organizations_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/orgs",
        "repos_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/repos",
        "events_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/events{/privacy}",
        "received_events_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/received_events",
        "type": "Organization",
        "site_admin": false
      },
      "repo": {
        "id": 284943,
        "node_id": "MDEwOlJlcG9zaXRvcnkyODQ5NDM=",
        "name": "Govis-CI",
        "full_name": "FXUBRQ-QE/Govis-CI",
        "private": false,
        "owner": {
          "login": "FXUBRQ-QE",
          "id": 19402,
          "node_id": "MDEyOk9yZ2FuaXphdGlvbjE5NDAy",
          "avatar_url": "<compnay git url>/avatars/u/19402?",
          "gravatar_id": "",
          "url": "<compnay git url>/api/v3/users/FXUBRQ-QE",
          "html_url": "<compnay git url>/FXUBRQ-QE",
          "followers_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/followers",
          "following_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/following{/other_user}",
          "gists_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/gists{/gist_id}",
          "starred_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/starred{/owner}{/repo}",
          "subscriptions_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/subscriptions",
          "organizations_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/orgs",
          "repos_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/repos",
          "events_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/events{/privacy}",
          "received_events_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/received_events",
          "type": "Organization",
          "site_admin": false
        },
        "html_url": "<compnay git url>/FXUBRQ-QE/Govis-CI",
        "description": null,
        "fork": false,
        "url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI",
        "forks_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/forks",
        "keys_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/keys{/key_id}",
        "collaborators_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/collaborators{/collaborator}",
        "teams_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/teams",
        "hooks_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/hooks",
        "issue_events_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/issues/events{/number}",
        "events_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/events",
        "assignees_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/assignees{/user}",
        "branches_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/branches{/branch}",
        "tags_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/tags",
        "blobs_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/git/blobs{/sha}",
        "git_tags_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/git/tags{/sha}",
        "git_refs_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/git/refs{/sha}",
        "trees_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/git/trees{/sha}",
        "statuses_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/statuses/{sha}",
        "languages_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/languages",
        "stargazers_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/stargazers",
        "contributors_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/contributors",
        "subscribers_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/subscribers",
        "subscription_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/subscription",
        "commits_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/commits{/sha}",
        "git_commits_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/git/commits{/sha}",
        "comments_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/comments{/number}",
        "issue_comment_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/issues/comments{/number}",
        "contents_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/contents/{+path}",
        "compare_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/compare/{base}...{head}",
        "merges_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/merges",
        "archive_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/{archive_format}{/ref}",
        "downloads_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/downloads",
        "issues_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/issues{/number}",
        "pulls_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/pulls{/number}",
        "milestones_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/milestones{/number}",
        "notifications_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/notifications{?since,all,participating}",
        "labels_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/labels{/name}",
        "releases_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/releases{/id}",
        "deployments_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/deployments",
        "created_at": "2019-06-27T07:26:55Z",
        "updated_at": "2019-10-03T08:05:35Z",
        "pushed_at": "2019-10-15T06:19:34Z",
        "git_url": "git://<compnay ssh git url>/FXUBRQ-QE/Govis-CI.git",
        "ssh_url": "git@<compnay ssh git url>:FXUBRQ-QE/Govis-CI.git",
        "clone_url": "<compnay git url>/FXUBRQ-QE/Govis-CI.git",
        "svn_url": "<compnay git url>/FXUBRQ-QE/Govis-CI",
        "homepage": null,
        "size": 234,
        "stargazers_count": 2,
        "watchers_count": 2,
        "language": "Go",
        "has_issues": true,
        "has_projects": true,
        "has_downloads": true,
        "has_wiki": true,
        "has_pages": false,
        "forks_count": 2,
        "mirror_url": null,
        "archived": false,
        "disabled": false,
        "open_issues_count": 11,
        "license": null,
        "forks": 2,
        "open_issues": 11,
        "watchers": 2,
        "default_branch": "master"
      }
    },
    "_links": {
      "self": {
        "href": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/pulls/100"
      },
      "html": {
        "href": "<compnay git url>/FXUBRQ-QE/Govis-CI/pull/100"
      },
      "issue": {
        "href": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/issues/100"
      },
      "comments": {
        "href": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/issues/100/comments"
      },
      "review_comments": {
        "href": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/pulls/100/comments"
      },
      "review_comment": {
        "href": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/pulls/comments{/number}"
      },
      "commits": {
        "href": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/pulls/100/commits"
      },
      "statuses": {
        "href": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/statuses/bb8eec32ab0f9a00564fd7e1142736ee87a8a5b2"
      }
    },
    "author_association": "COLLABORATOR",
    "draft": false,
    "merged": false,
    "mergeable": true,
    "rebaseable": true,
    "mergeable_state": "unstable",
    "merged_by": null,
    "comments": 1,
    "review_comments": 0,
    "maintainer_can_modify": false,
    "commits": 2,
    "additions": 390,
    "deletions": 335,
    "changed_files": 16
  },
  "requested_reviewer": {
    "login": "i337562",
    "id": 21508,
    "node_id": "MDQ6VXNlcjIxNTA4",
    "avatar_url": "<compnay git url>/avatars/u/21508?",
    "gravatar_id": "",
    "url": "<compnay git url>/api/v3/users/i337562",
    "html_url": "<compnay git url>/i337562",
    "followers_url": "<compnay git url>/api/v3/users/i337562/followers",
    "following_url": "<compnay git url>/api/v3/users/i337562/following{/other_user}",
    "gists_url": "<compnay git url>/api/v3/users/i337562/gists{/gist_id}",
    "starred_url": "<compnay git url>/api/v3/users/i337562/starred{/owner}{/repo}",
    "subscriptions_url": "<compnay git url>/api/v3/users/i337562/subscriptions",
    "organizations_url": "<compnay git url>/api/v3/users/i337562/orgs",
    "repos_url": "<compnay git url>/api/v3/users/i337562/repos",
    "events_url": "<compnay git url>/api/v3/users/i337562/events{/privacy}",
    "received_events_url": "<compnay git url>/api/v3/users/i337562/received_events",
    "type": "User",
    "site_admin": false
  },
  "repository": {
    "id": 284943,
    "node_id": "MDEwOlJlcG9zaXRvcnkyODQ5NDM=",
    "name": "Govis-CI",
    "full_name": "FXUBRQ-QE/Govis-CI",
    "private": false,
    "owner": {
      "login": "FXUBRQ-QE",
      "id": 19402,
      "node_id": "MDEyOk9yZ2FuaXphdGlvbjE5NDAy",
      "avatar_url": "<compnay git url>/avatars/u/19402?",
      "gravatar_id": "",
      "url": "<compnay git url>/api/v3/users/FXUBRQ-QE",
      "html_url": "<compnay git url>/FXUBRQ-QE",
      "followers_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/followers",
      "following_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/following{/other_user}",
      "gists_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/gists{/gist_id}",
      "starred_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/starred{/owner}{/repo}",
      "subscriptions_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/subscriptions",
      "organizations_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/orgs",
      "repos_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/repos",
      "events_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/events{/privacy}",
      "received_events_url": "<compnay git url>/api/v3/users/FXUBRQ-QE/received_events",
      "type": "Organization",
      "site_admin": false
    },
    "html_url": "<compnay git url>/FXUBRQ-QE/Govis-CI",
    "description": null,
    "fork": false,
    "url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI",
    "forks_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/forks",
    "keys_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/keys{/key_id}",
    "collaborators_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/collaborators{/collaborator}",
    "teams_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/teams",
    "hooks_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/hooks",
    "issue_events_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/issues/events{/number}",
    "events_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/events",
    "assignees_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/assignees{/user}",
    "branches_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/branches{/branch}",
    "tags_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/tags",
    "blobs_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/git/blobs{/sha}",
    "git_tags_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/git/tags{/sha}",
    "git_refs_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/git/refs{/sha}",
    "trees_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/git/trees{/sha}",
    "statuses_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/statuses/{sha}",
    "languages_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/languages",
    "stargazers_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/stargazers",
    "contributors_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/contributors",
    "subscribers_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/subscribers",
    "subscription_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/subscription",
    "commits_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/commits{/sha}",
    "git_commits_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/git/commits{/sha}",
    "comments_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/comments{/number}",
    "issue_comment_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/issues/comments{/number}",
    "contents_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/contents/{+path}",
    "compare_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/compare/{base}...{head}",
    "merges_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/merges",
    "archive_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/{archive_format}{/ref}",
    "downloads_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/downloads",
    "issues_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/issues{/number}",
    "pulls_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/pulls{/number}",
    "milestones_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/milestones{/number}",
    "notifications_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/notifications{?since,all,participating}",
    "labels_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/labels{/name}",
    "releases_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/releases{/id}",
    "deployments_url": "<compnay git url>/api/v3/repos/FXUBRQ-QE/Govis-CI/deployments",
    "created_at": "2019-06-27T07:26:55Z",
    "updated_at": "2019-10-03T08:05:35Z",
    "pushed_at": "2019-10-15T06:19:34Z",
    "git_url": "git://<compnay ssh git url>/FXUBRQ-QE/Govis-CI.git",
    "ssh_url": "git@<compnay ssh git url>:FXUBRQ-QE/Govis-CI.git",
    "clone_url": "<compnay git url>/FXUBRQ-QE/Govis-CI.git",
    "svn_url": "<compnay git url>/FXUBRQ-QE/Govis-CI",
    "homepage": null,
    "size": 234,
    "stargazers_count": 2,
    "watchers_count": 2,
    "language": "Go",
    "has_issues": true,
    "has_projects": true,
    "has_downloads": true,
    "has_wiki": true,
    "has_pages": false,
    "forks_count": 2,
    "mirror_url": null,
    "archived": false,
    "disabled": false,
    "open_issues_count": 11,
    "license": null,
    "forks": 2,
    "open_issues": 11,
    "watchers": 2,
    "default_branch": "master"
  },
  "organization": {
    "login": "FXUBRQ-QE",
    "id": 19402,
    "node_id": "MDEyOk9yZ2FuaXphdGlvbjE5NDAy",
    "url": "<compnay git url>/api/v3/orgs/FXUBRQ-QE",
    "repos_url": "<compnay git url>/api/v3/orgs/FXUBRQ-QE/repos",
    "events_url": "<compnay git url>/api/v3/orgs/FXUBRQ-QE/events",
    "hooks_url": "<compnay git url>/api/v3/orgs/FXUBRQ-QE/hooks",
    "issues_url": "<compnay git url>/api/v3/orgs/FXUBRQ-QE/issues",
    "members_url": "<compnay git url>/api/v3/orgs/FXUBRQ-QE/members{/member}",
    "public_members_url": "<compnay git url>/api/v3/orgs/FXUBRQ-QE/public_members{/member}",
    "avatar_url": "<compnay git url>/avatars/u/19402?",
    "description": "Quality Engineers from FXperience teams in Brno"
  },
  "enterprise": {
    "id": 1,
    "slug": "sap-se",
    "name": "SAP SE",
    "node_id": "MDg6QnVzaW5lc3Mx",
    "avatar_url": "<compnay git url>/avatars/b/1?",
    "description": null,
    "website_url": null,
    "html_url": "<compnay git url>/businesses/sap-se",
    "created_at": "2019-03-16T05:31:15Z",
    "updated_at": "2019-03-16T05:31:15Z"
  },
  "sender": {
    "login": "I516366",
    "id": 51335,
    "node_id": "MDQ6VXNlcjUxMzM1",
    "avatar_url": "<compnay git url>/avatars/u/51335?",
    "gravatar_id": "",
    "url": "<compnay git url>/api/v3/users/I516366",
    "html_url": "<compnay git url>/I516366",
    "followers_url": "<compnay git url>/api/v3/users/I516366/followers",
    "following_url": "<compnay git url>/api/v3/users/I516366/following{/other_user}",
    "gists_url": "<compnay git url>/api/v3/users/I516366/gists{/gist_id}",
    "starred_url": "<compnay git url>/api/v3/users/I516366/starred{/owner}{/repo}",
    "subscriptions_url": "<compnay git url>/api/v3/users/I516366/subscriptions",
    "organizations_url": "<compnay git url>/api/v3/users/I516366/orgs",
    "repos_url": "<compnay git url>/api/v3/users/I516366/repos",
    "events_url": "<compnay git url>/api/v3/users/I516366/events{/privacy}",
    "received_events_url": "<compnay git url>/api/v3/users/I516366/received_events",
    "type": "User",
    "site_admin": false
  }
}`
