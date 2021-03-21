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
	docSnapshot, _ := docRef.Get(context.Background())

	// change doc since the last GET
	_, err := docRef.Update(ctx, []firestore.Update{{Path: "Name", Value: "Weave Inc. Again"}})
	if err != nil {
		panic(err)
	}

	// Firestore is eventually consistent, wait for a while
	time.Sleep(time.Second)

	// try to patch the doc based on the original last update time
	_, err = docRef.Update(ctx, []firestore.Update{{Path: "Name", Value: "My Weave"}},
		// firestore.LastUpdateTime will not write unless the "LastUpdateTime" matches the given
		// * Useful for concurrent write detection
		firestore.LastUpdateTime(docSnapshot.UpdateTime),
	)
	if err != nil {
		panic(err)
	}
	// END UPDATE OMIT

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
