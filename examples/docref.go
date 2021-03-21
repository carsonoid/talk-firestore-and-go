package main

import (
	"context"
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

	c, cerr := firestore.NewClient(ctx, "firestore-and-go")
	if cerr != nil {
		panic(cerr)
	}

	// cleanup ahead of possible conflicts
	c.Collection("customers").Doc("weave").Delete(ctx)

	// START CREATE OMIT
	client, err := firestore.NewClient(ctx, "firestore-and-go")
	if cerr != nil {
		panic(cerr)
	}

	docRef := client.Collection("customers").Doc("weave")

	//...
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
