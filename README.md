# go-pocketsmith

Go client library for the [PocketSmith API](https://developers.pocketsmith.com/).

> ⚠️ Work in progress — API may change.

## Installation

```sh
go get github.com/rhyselsmore/go-pocketsmith
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/rhyselsmore/go-pocketsmith"
)

func main() {
	client, err := pocketsmith.New(os.Getenv("POCKETSMITH_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	me, err := client.GetMe(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// List uncategorized transactions
	p := pocketsmith.ListTransactionsParams{Uncategorized: true}
	page, err := client.ListTransactionsInUser(ctx, me.ID, p)
	if err != nil {
		log.Fatal(err)
	}

	for _, tx := range page.Items {
		fmt.Printf("%s %s %.2f\n", tx.Date, tx.Payee, tx.Amount)
	}

	// Paginate through results
	for page.PageInfo.HasNext() {
		p.Page++
		page, err = client.ListTransactionsInUser(ctx, me.ID, p)
		if err != nil {
			log.Fatal(err)
		}
		for _, tx := range page.Items {
			fmt.Printf("%s %s %.2f\n", tx.Date, tx.Payee, tx.Amount)
		}
	}
}
```

## License

[MIT](LICENSE)