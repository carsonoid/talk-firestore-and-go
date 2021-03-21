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

	client.Collection("customers").Doc("ut-gophers").Delete(ctx)
	client.Collection("customers").Doc("ut-gophers").Collection("users").Doc("carsonoid").Delete(ctx)

	// START TRANSACTION OMIT
	// Create a new customer and new user in the customer in one transaction
	custRef := client.Collection("customers").Doc("ut-gophers")
	userRef := custRef.Collection("users").Doc("carsonoid")

	_, _ = custRef.Create(ctx, &Customer{ID: "ut-gophers", Name: "Utah Gophers"})

	err := client.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
		custSnapshot, err := t.Get(custRef)
		if err != nil {
			return err
		}

		custName, _ := custSnapshot.DataAt("Name")
		err = t.Create(userRef, &User{ID: "carsonoid", Name: fmt.Sprint("Carsonoid FROM ", custName)})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	// END TRANSACTION OMIT
	fmt.Println("Created new customer and user all at once!")
	dumpDoc(custRef)
	dumpDoc(userRef)
}

func dumpDoc(docRef *firestore.DocumentRef) {
	docSnapshot, err := docRef.Get(context.Background())
	if err != nil {
		panic(err)
	}

	litter.Dump(docSnapshot.Data())
}
