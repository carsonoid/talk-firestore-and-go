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

func main() {
	ctx := context.Background()

	client, cerr := firestore.NewClient(ctx, "firestore-and-go")
	if cerr != nil {
		panic(cerr)
	}

	// START ITER OMIT

	// We will be using a new docref for the watch
	docRef := client.Collection("customers").Doc("weavedevx")

	// create an iterator for the docRef
	iter := docRef.Snapshots(ctx)

	// iter.Next() get's the next doc update and err, if any
	// iter.Stop() cancels the iterator
	// END ITER OMIT

	// START OPS OMIT
	go func() {
		time.Sleep(time.Second)
		docRef.Create(ctx, &Customer{ID: "weavedevx"})

		time.Sleep(time.Second)
		docRef.Update(ctx, []firestore.Update{{Path: "Name", Value: "weave devx"}})

		time.Sleep(time.Second)
		docRef.Update(ctx, []firestore.Update{{Path: "Name", Value: "Weave Devx"}})

		time.Sleep(time.Second)
		docRef.Delete(ctx)

		time.Sleep(time.Second)
		iter.Stop()
	}()

	// END OPS OMIT

	// START WATCH OMIT
	for {
		docSnapshot, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			panic(err)
		}

		// Exists returns false for deleted or non-existing docs
		// The first event for a not-yet-created doc would hit this
		if !docSnapshot.Exists() {
			fmt.Println("Doc does not exist or was deleted")
			continue
		}

		cust := &Customer{}
		err = docSnapshot.DataTo(cust)
		if err != nil {
			panic(err)
		}
		litter.Dump(cust)
	}
	// END WATCH OMIT
}
