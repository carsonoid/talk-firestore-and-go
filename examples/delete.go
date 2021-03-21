package main

import (
	"context"
	"fmt"

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

	// START DELETE OMIT
	// Create a slice of paths and values
	_, err := docRef.Delete(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("DELETED!")

	// try to get delete do
	_, err = docRef.Get(ctx)
	if err != nil {
		panic(err)
	}
	// END DELETE OMIT

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
