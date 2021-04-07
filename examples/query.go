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

	//pre-add
	client.Collection("customers").Doc("weave").Set(ctx, &Customer{ID: "weave", Name: "Weave"})
	client.Collection("customers").Doc("firestore-and-go").Set(ctx, &Customer{ID: "firestore-and-go", Name: "Firestore And Go"})

	// START QUERY OMIT
	query := client.Collection("customers").Query.
		Where("Name", "in", []string{"Weave", "Firestore And Go"}).
		Limit(2)

	docSnapshots, err := query.Documents(ctx).GetAll()
	if err != nil {
		panic(err)
	}

	fmt.Println("Got", len(docSnapshots), "docs!")

	for _, docSnapshot := range docSnapshots {
		litter.Dump(docSnapshot.Data())
	}
	// END QUERY OMIT
}
