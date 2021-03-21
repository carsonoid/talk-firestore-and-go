package main

import (
	"context"
	"time"

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

	// START UPDATE OMIT
	updates := []firestore.Update{
		// Updates can be by path
		{Path: "Name", Value: "Weave"},

		// ...even when adding new fields
		{Path: "PatchedAt", Value: time.Now()},

		// ...even with new nested fields
		{Path: "Hiring.URL", Value: "https://boards.greenhouse.io/weavehq"},
		{Path: "Hiring.ReferredBy", Value: "Carson Anderson"},

		// Paths that contain any of ".˜*/[]" must be wrapped in a custom type
		{FieldPath: firestore.FieldPath{"[my]", "sp.ecial", "/key**"}, Value: "Company"},
	}

	// apply all updates to the targeted doc
	_, _ = docRef.Update(ctx, updates)

	// In Go, you MUST get as a map to not drop unknown fields
	docSnapshot, _ := docRef.Get(context.Background())
	litter.Dump(docSnapshot.Data())
	// END UPDATE OMIT
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
