# ghe-api

## usage

1. edit to `pkg/env/.env.yaml`
1. `go run cmd/main.go`

### option

- `-f [type]` is filtering event type.
  - IssueCommentEvent
  - IssuesEvent
  - PullRequestEvent
  - PullRequestReviewCommentEvent
  - ...And many more. See the [officialDoc](https://docs.github.com/ja/developers/webhooks-and-events/github-event-types) for details.

## attention

[official](https://docs.github.com/en/github-ae@latest/rest/reference/activity), the maximum number is 300, but in reality, it seems to be less than that.
> Events support pagination, however the per_page option is unsupported. The fixed page size is 30 items. Fetching up to ten pages is supported, for a total of 300 events. For information, see "Traversing with pagination."

-> In fact, up to 267.
