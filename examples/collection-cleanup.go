package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/genproto/googleapis/type/latlng"
)

type Customer struct {
	ID   string
	Name string

	// Location is a point on the earth
	// it MUST be a pointer to saved as the right type in Firestore
	Location *latlng.LatLng
}

func main() {
	ctx := context.Background()

	client, cerr := firestore.NewClient(ctx, "firestore-and-go")
	if cerr != nil {
		panic(cerr)
	}

	// START CLEANUP OMIT
	// We will be building a firestore query for the watch
	iter := client.Collection("customers").Doc("weave").Collection("users").DocumentRefs(ctx)

	docRefs, err := iter.GetAll()
	if err != nil {
		panic(err)
	}

	for _, docRef := range docRefs {
		docRef.Delete(ctx)
		if err != nil {
			panic(err)
		}
	}
	// START CLEANUPLINE OMIT
	// Do Cleanup
	fmt.Println("Deleted", len(docRefs), "documents")
	// END CLEANUPLINE OMIT

	// END CLEANUP OMIT
}
