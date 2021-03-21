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
	userRef := custRef.Collection("users").Doc("carsonoid")

	err := client.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
		custSnap, err := t.Get(custRef)
		if err != nil {
			return err
		}

		userSnap, err := t.Get(userRef)
		if err != nil {
			return err
		}

		fmt.Println("The Docs look like this right now:")
		fmt.Println(custSnap.Data())
		fmt.Println(userSnap.Data())
		return nil
	}, firestore.ReadOnly) // Run the transaction in readonly mode
	if err != nil {
		panic(err)
	}
	// END TRANSACTION OMIT
}

func dumpDoc(docSnapshot *firestore.DocumentSnapshot) {
	user := &User{}
	err := docSnapshot.DataTo(user)
	if err != nil {
		panic(err)
	}
	litter.Dump(user)
}
