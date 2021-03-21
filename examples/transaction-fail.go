package main

import (
	"context"
	"errors"
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

type User struct {
	ID   string
	Name string
}

func main() {
	ctx := context.Background()

	client, cerr := firestore.NewClient(ctx, "firestore-and-go")
	if cerr != nil {
		panic(cerr)
	}

	// START TRANSACTION OMIT
	// Create a new customer and new user in the customer in one transaction
	custRef := client.Collection("customers").Doc("ut-gophers")

	err := client.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
		err := t.Set(custRef, &Customer{ID: "DO_NOT_SET_ME"})
		if err != nil {
			return err
		}

		return errors.New("test error, cancel transaction")
	})
	if err != nil {
		fmt.Println("ERROR", err)
	}

	docSnapshot, _ := custRef.Get(ctx)
	cust := &Customer{}
	err = docSnapshot.DataTo(cust)
	if err != nil {
		panic(err)
	}
	litter.Dump(cust)
	// END TRANSACTION OMIT
}

func must(_ interface{}, err error) {
	if err != nil {
		panic(err)
	}
}
