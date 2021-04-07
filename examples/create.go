package main

import (
	"context"
	"fmt"
	"time"

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

	// cleanup ahead of possible conflicts
	client.Collection("customers").Doc("weave").Delete(ctx)

	// START CREATE OMIT
	docRef := client.Collection("customers").Doc("weave")

	result, err := docRef.Create(ctx, &Customer{
		ID:   "weave",
		Name: "weave",
		Location: &latlng.LatLng{
			Latitude:  40.4162205,
			Longitude: -111.8718743,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("result:\n%s\n%#v\n", *result, *result)
	// END CREATE OMIT
}

// START WriteResult OMIT

// Create(ctx context.Context, data interface{}) (_ *firestore.WriteResult, err error) {

//// A WriteResult is returned by methods that write documents.
type WriteResult struct {
	// The time at which the document was updated, or created if it did not
	// previously exist. Writes that do not actually change the document do
	// not change the update time.
	UpdateTime time.Time
}

// END WriteResult OMIT
