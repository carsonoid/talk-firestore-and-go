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

	// START GET OMIT
	// Do the document Get once, this is the only part of the opertion which actually talks to firestore
	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		panic(err)
	}

	cust := &Customer{}
	err = docSnapshot.DataTo(cust)
	if err != nil {
		panic(err)
	}
	fmt.Print("extracted into a struct:\n", litter.Sdump(cust), "\n\n")

	name, err := docSnapshot.DataAt("Name")
	if err != nil {
		panic(err)
	}
	fmt.Print("single field by name: ", name, "\n")

	data := docSnapshot.Data()
	fmt.Print("\nextracted as map[string]interface data:\n", litter.Sdump(data), "\n")
	// END GET OMIT
}
