package main

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/sanity-io/litter"
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

	docRef := client.Collection("customers").Doc("weave")

	// START GET OMIT
	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		panic(err)
	}

	// At this point, the data has been fetched
	// and is stored inside the snapshot, but is private
	litter.Dump(docSnapshot)
	// END GET OMIT
}
