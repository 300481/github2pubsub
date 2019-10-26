# github2pubsub

A bridge from Github webhooks to Google PubSub

## Environment Settings

Environment Variable     |Description                                 |Type / valid values
-------------------------|--------------------------------------------|-----------------------
`GCP_TOPIC_NAME`         |The PubSub topic name                       |*String*
`GCP_CREATE_TOPIC`       |Allow creation of topic if not exists       |*String* `TRUE`/`FALSE`
`GCP_PROJECT_ID`         |The Google Project ID                       |*String*
`GITHUB_SECRET`          |The Github Secret                           |*String*
`GITHUB_EVENTS`          |The Github Event Types, separated by `/`    |*String*
