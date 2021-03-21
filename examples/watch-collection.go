package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/sanity-io/litter"
	"google.golang.org/api/iterator"
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

	// START ITER OMIT
	// We will be building a firestore query for the watch
	query := client.Collection("customers").Doc("weave").Collection("users").Query

	// create an iterator for the docRef
	iter := query.Snapshots(ctx)

	// iter.Next() get's the next snapshot and err, if any
	// iter.Stop() cancels the iterator
	// END ITER OMIT

	// START OPS OMIT
	go func() {
		for i := range []int{1, 2, 3, 4, 5} {
			time.Sleep(time.Second)
			userID := fmt.Sprint("user-", i)
			docRef := client.Collection("customers").Doc("weave").Collection("users").Doc(userID)

			docRef.Create(ctx, &User{ID: userID}) // ignore error for demo purposes

			_, err := docRef.Update(ctx, []firestore.Update{{Path: "Name", Value: userID}})
			if err != nil {
				panic(err)
			}
		}
		time.Sleep(time.Second)

		iter.Stop()
	}()

	// END OPS OMIT

	// START WATCH OMIT
	for {
		snapshot, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			panic(err)
		}
		// ...
		// END WATCH OMIT

		// START WATCHWORK OMIT
		// for iterator.Next
		for _, change := range snapshot.Changes { // each snapshot may include multiple document changes
			switch change.Kind {
			case firestore.DocumentAdded:
				fmt.Println("ADDED")
				dumpDoc(change.Doc)
			case firestore.DocumentModified:
				fmt.Println("MODIFIED")
				dumpDoc(change.Doc)
			case firestore.DocumentRemoved:
				fmt.Println("REMOVED")
				dumpDoc(change.Doc)
			}
		}
		// END WATCHWORK OMIT
	}
}

func dumpDoc(docSnapshot *firestore.DocumentSnapshot) {
	user := &User{}
	err := docSnapshot.DataTo(user)
	if err != nil {
		panic(err)
	}
	litter.Dump(user)
}
