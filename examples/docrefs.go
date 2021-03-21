package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

func main() {
	ctx := context.Background()

	// START DOCREF OMIT
	// Get a client
	client, err := firestore.NewClient(ctx, "firestore-and-go")
	if err != nil {
		panic(err)
	}

	// the ideal way to create references is by alternating
	// .Collection and .Doc calls on the client
	docRef := client.Collection("customers").Doc("weave")
	fmt.Println(docRef.Path)

	// It is possible to include / to separate collections and docs
	fmt.Println(client.Doc("customers/weave").Path)
	fmt.Println(client.Collection("customers/weave/users").Doc("bob").Path)
	fmt.Println(client.Doc("customers/weave/users/bob").Path)
	// ... BUT malformed references can cause nil references to be returned
	fmt.Println(client.Doc("customers//users"))       // empty component
	fmt.Println(client.Doc("customers/weave/users"))  // odd number of components to doc
	fmt.Println(client.Collection("customers/weave")) // even number of components to coll
	// END DOCREF OMIT
}
