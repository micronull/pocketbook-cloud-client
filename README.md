# pocketbook-cloud-client
PocketBook Cloud API Client

## Example

```go
package main

import (
	"context"
	"encoding/json"
	"log"

	pbc "github.com/micronull/pocketbook-cloud-client"
)

func main() {
	cli := pbc.New(
		pbc.WithClientID("qNAx1RDb"),
		pbc.WithClientSecret("K3YYSjCgDJNoWKdGVOyO1mrROp3MMZqqRNXNXTmh"),
	)
	ctx := context.Background()

	prvs, err := cli.Providers(ctx, "you.mail.box@some.com")
	if err != nil {
		log.Fatal(err)
	}

	for _, prv := range prvs {
		req := pbc.LoginRequest{
			ShopID:   prv.ShopID,
			UserName: "you.mail.box@some.com",
			Password: "you.password",
			Provider: prv.Alias,
		}

		tkn, err := cli.Login(ctx, req)
		if err != nil {
			log.Fatal(err)
		}

		books, err := cli.Books(ctx, tkn.AccessToken, 0, 0)
		if err != nil {
			log.Fatal(err)
		}

		if books.Total == 0 {
			continue
		}

		books, err = cli.Books(ctx, tkn.AccessToken, books.Total, 0)
		if err != nil {
			log.Fatal(err)
		}

		js, _ := json.MarshalIndent(books, "", "  ")

		log.Println(string(js))
	}
}
```