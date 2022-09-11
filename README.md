# go-notion

[![GitHub tag (latest
SemVer)](https://img.shields.io/github/v/tag/dstotijn/go-notion?label=go%20module)](https://github.com/dstotijn/go-notion/tags)
[![Test](https://github.com/dstotijn/go-notion/actions/workflows/test.yaml/badge.svg)](https://github.com/dstotijn/go-notion/actions/workflows/test.yaml)
[![Go
Reference](https://pkg.go.dev/badge/github.com/dstotijn/go-notion.svg)](https://pkg.go.dev/github.com/dstotijn/go-notion)
[![GitHub](https://img.shields.io/github/license/dstotijn/go-notion)](LICENSE)
[![Go Report
Card](https://goreportcard.com/badge/github.com/dstotijn/go-notion)](https://goreportcard.com/report/github.com/dstotijn/go-notion)

**go-notion** is a client for the [Notion
API](https://developers.notion.com/reference), written in
[Go](https://golang.org/).

## Features

The client supports all (non-deprecated) endpoints available in the Notion API,
as of September 4, 2022:

<details>
<summary>Databases</summary>

- [x] [Query a
      database](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.QueryDatabase)
- [x] [Create a
      database](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.CreateDatabase)
- [x] [Update
      database](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.UpdateDatabase)
- [x] [Retrieve a
    database](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.FindDatabaseByID)
</details>

<details>
<summary>Pages</summary>

- [x] [Retrieve a
      page](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.FindPageByID)
- [x] [Create a
      page](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.CreatePage)
- [x] [Update
      page](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.UpdatePage)
- [x] [Retrieve a page
    property](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.FindPagePropertyByID)
</details>

<details>
<summary>Blocks</summary>

- [x] [Retrieve a
      block](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.FindBlockByID)
- [x] [Update
      block](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.UpdateBlock)
- [x] [Retrieve block
      children](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.FindBlockChildrenByID)
- [x] [Append block
      children](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.AppendBlockChildren)
- [x] [Delete
    block](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.DeleteBlock)
</details>

<details>
<summary>Users</summary>

- [x] [Retrieve a
      user](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.FindUserByID)
- [x] [List all
      users](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.ListUsers)
- [x] [Retrieve your token's bot
    user](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.FindCurrentUser)
</details>

<details>
<summary>Search</summary>

- [x] [Search](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.Search)
</details>

<details>
<summary>Comments</summary>

- [x] [Retrieve
      comments](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.FindCommentsByBlockID)
- [x] [Create a
    comment](https://pkg.go.dev/github.com/dstotijn/go-notion#Client.CreateComment)
</details>

## Installation

```sh
$ go get github.com/dstotijn/go-notion
```

## Usage

To obtain an API key, follow Notion’s [getting started
guide](https://developers.notion.com/docs/getting-started).

```go
import "github.com/dstotijn/go-notion"

(...)

client := notion.NewClient("secret-api-key")

page, err := client.FindPageByID(context.Background(), "18d35eb5-91f1-4dcb-85b0-c340fd965015")
if err != nil {
    // Handle error...
}
```

👉 Check out the docs on
[pkg.go.dev](https://pkg.go.dev/github.com/dstotijn/go-notion) for a complete
reference and the [examples](/examples) directory for more example code.

## Status

The Notion API itself is out of beta. This library is updated periodically
following documented changes from the Notion
[changelog](https://developers.notion.com/changelog).

**Note:** This library **will** make breaking changes in its code until
`v1.0` of the module is released. There are no immediate plans for a `v1.0`
release. I want the design choices to be solidified and battle-tested more
before committing to a stable release (and the possible burden of a "v2+" Go
module should I want to introduce breaking changes).

## License

[MIT License](LICENSE)

© 2022 [David Stotijn](https://v0x.nl)
