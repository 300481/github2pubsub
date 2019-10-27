# github2pubsub

A bridge from Github webhooks to Google PubSub

## Environment Variables

Environment Variable     |Description                                 |Valid values
-------------------------|--------------------------------------------|-----------------------
`GCP_TOPIC_NAME`         |The PubSub topic name                       |*String*
`GCP_CREATE_TOPIC`       |Allow creation of topic if not exists       |*String* `TRUE`/`FALSE`
`GCP_PROJECT_ID`         |The Google Project ID                       |*String*
`GITHUB_SECRET`          |The Github Secret                           |*String*
`GITHUB_EVENTS`          |The Github Event Types, separated by `/`    |*String* f.e. `ping/push`

For more information about possible events see the [sourcecode](https://github.com/go-playground/webhooks/blob/v5/github/github.go#L30).

## Installation

* `git clone https://github.com/300481/github2pubsub.git`
* `cd deployments`
* ```bash
export GCP_TOPIC_NAME=yourtopicname
export GCP_CREATE_TOPIC=TRUE
export GCP_PROJECT_ID=yourprojectid
export GITHUB_SECRET=yourgithubsecret
export GITHUB_EVENTS=ping/yourwantedevents
```
