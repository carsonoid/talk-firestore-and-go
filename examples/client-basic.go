package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

func main() {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, "firestore-and-go")
	if err != nil {
		panic(err)
	}

	// NOTE: the .Close method doesn't *need* to be called on program exit
	// You only really need it if you are creating ephemeral clients
	defer client.Close()

	fmt.Println(*client)
}
