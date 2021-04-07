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

	// START SET OMIT
	// Set will completely overwrite the document with the given data by default
	// ... but you can use merge options to allow smarter merging

	// create untyped map data
	mapData := map[string]interface{}{
		"Location": &latlng.LatLng{
			Latitude:  40.4162205,
			Longitude: -111.8718743,
		},
	}

	_, err := docRef.Set(ctx, mapData,
		// MergeAll causes all paths to be overwritten
		// but only works with map-based data in Go
		firestore.MergeAll,
	)
	if err != nil {
		panic(err)
	}
	// END SET OMIT

	getAndDump(docRef)
}

func getAndDump(docRef *firestore.DocumentRef) {
	cust := &Customer{}
	docSnapshot, err := docRef.Get(context.Background())
	if err != nil {
		panic(err)
	}
	err = docSnapshot.DataTo(cust)
	if err != nil {
		panic(err)
	}
	litter.Dump(cust)
}

func getAndDumpData(docRef *firestore.DocumentRef) {
	docSnapshot, err := docRef.Get(context.Background())
	if err != nil {
		panic(err)
	}
	litter.Dump(docSnapshot.Data())
}
