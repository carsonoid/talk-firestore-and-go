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

	// pre-delete batch
	b := client.Batch()
	for i := range []int{6, 7, 8, 9, 10} {
		userID := fmt.Sprint("user-", i)
		docRef := client.Collection("customers").Doc("weave").Collection("users").Doc(userID)

		b.Delete(docRef)
	}
	b.Commit(ctx)

	// START BATCH OMIT
	batch := client.Batch()

	for i := range []int{6, 7, 8, 9, 10} {
		userID := fmt.Sprint("user-", i)
		docRef := client.Collection("customers").Doc("weave").Collection("users").Doc(userID)

		// no error because nothing has happened yet
		batch.Create(docRef, &User{ID: userID})
	}

	// commit batch
	result, err := batch.Commit(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("created", len(result), "docs")
	// END BATCH OMIT
}

func dumpDoc(docSnapshot *firestore.DocumentSnapshot) {
	user := &User{}
	err := docSnapshot.DataTo(user)
	if err != nil {
		panic(err)
	}
	litter.Dump(user)
}
