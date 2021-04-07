package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/sanity-io/litter"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/type/latlng"
)

func getClient() *firestore.Client {
	// START CLIENT OMIT
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, "firestore-and-go")
	if err != nil {
		panic(err)
	}
	// END CLIENT OMIT

	return client
}

func getClientAdvanced() *firestore.Client {
	// START ADVCLIENT OMIT
	ctx := context.Background()

	userHome, _ := os.UserHomeDir()
	gcpCredPath := filepath.Join(userHome, ".gcp", "creds")

	client, err := firestore.NewClient(ctx, "firestore-and-go",
		option.WithRequestReason("demo"),        // Supply reason to audit logs
		option.WithCredentialsFile(gcpCredPath)) // override default args
	if err != nil {
		panic(err)
	}
	// END ADVCLIENT OMIT

	return client
}

// this is possible, but I'n not sure why you would use it
func getClientViaFirebase() *firestore.Client {
	ctx := context.Background()

	conf := &firebase.Config{ProjectID: "firestore-and-go"}
	fbApp, err := firebase.NewApp(ctx, conf)
	if err != nil {
		panic(err)
	}

	client, err := fbApp.Firestore(ctx)
	if err != nil {
		panic(err)
	}

	return client
}

type Address struct{}

// START CUSTOMER OMIT
// import "google.golang.org/genproto/googleapis/type/latlng"

type Customer struct {
	ID   string
	Name string

	Nicknames []string          `firestore:"-,omitempty"`
	Features  map[string]string `firestore:"feats"`
	Addresses []Address         `firestore:"-"`

	Sibling *Customer
	Child   struct {
		ID string
	}

	// Location is a point on the earth
	// it MUST be a pointer to saved as the right type in Firestore
	Location *latlng.LatLng
}

// END CUSTOMER OMIT

func main() {
	// cleanup in case of previous failures
	getClient().Collection("customers").Doc("weave").Delete(context.Background())

	ctx := context.Background()

	client := getClient()

	// START DOCREF OMIT
	// the ideal way to create references is by alternating
	// .Collection and .Doc calls
	docRef := client.Collection("customers").Doc("weave")

	// It is possible to include / to separate collections and docs
	docRef = client.Doc("customers/weave")

	// ... but malformed references can cause nil references to be returned
	//     Ex: below is a reference to a collection, not a document
	//         many of the firestore cli functions handle this and print
	//         a warning, but your code might not
	fmt.Println("bad ref is:", client.Doc("customers/weave/users"))
	// result: bad ref is: <nil>
	// END DOCREF OMIT

	// START CREATE OMIT
	fmt.Println("CREATE")
	_, err := docRef.Create(ctx, &Customer{
		ID:   "weave",
		Name: "weave",
		Location: &latlng.LatLng{
			Latitude:  40.4162205,
			Longitude: -111.8718743,
		},
	})
	// END CREATE OMIT
	if err != nil {
		panic(err)
	}

	// START GET OMIT
	fmt.Println()
	fmt.Println("GET")
	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		panic(err)
	}

	// END GET OMIT

	// START Data As Map OMIT
	data := docSnapshot.Data()
	litter.Dump(data)
	// END  Data As Map OMIT

	fmt.Println()
	// START GET - To Struct OMIT
	fmt.Println("GET - To Struct")
	cust := &Customer{}
	err = docSnapshot.DataTo(cust)
	if err != nil {
		panic(err)
	}
	litter.Dump(cust)
	// END GET - To Struct OMIT

	fmt.Println()
	// START UPDATE OMIT
	fmt.Println("UPDATE")
	updates := []firestore.Update{
		// Updates can be by path
		{Path: "Name", Value: "Weave"},
		// ...even when adding new fields
		{Path: "CreatedOn", Value: time.Now()},
		// ...even with new nested fields
		{Path: "Hiring.URL", Value: "https://boards.greenhouse.io/weavehq"},
		{Path: "Hiring.ReferredBy", Value: "Carson Anderson"},

		// Paths that contain any of ".˜*/[]" must be wrapped in a custom type
		{FieldPath: firestore.FieldPath{"[my]", "sp.ecial", "/key**"}, Value: "Company"},
	}
	_, err = docRef.Update(ctx, updates)
	if err != nil {
		panic(err)
	}
	getAndDump(docRef)
	getAndDumpData(docRef)
	// END UPDATE OMIT

	// // You can also set update "preconditions" for more update safety
	// _, err = docRef.Update(ctx, []firestore.Update{{Path: "PassedPreconditions", Value: true}},
	// 	// firestore.Exists will not write unless the doc exists
	// 	// * Useful for not updating missing docs
	// 	firestore.Exists,

	// 	// firestore.LastUpdateTime will not write unless the "LastUpdateTime" matches the given
	// 	// * Useful for concurrent write detection
	// 	firestore.LastUpdateTime(result.UpdateTime))
	// if err != nil {
	// 	panic(err)
	// }
	// getAndDump(docRef)

	fmt.Println()
	// START SET OMIT
	fmt.Println("SET")
	// Set will completely overwrite the document with the given data by default
	_, err = docRef.Set(ctx, &Customer{
		Name: "Weave",
	})
	if err != nil {
		panic(err)
	}
	getAndDump(docRef)

	fmt.Println()
	// START SET - Merge Selected OMIT
	fmt.Println("SET - Merge Selected")
	// ... but you can use merge options to allow smarter merging
	_, err = docRef.Set(ctx, &Customer{Name: "", ID: "weave"},
		// Merge allows you to specify N fieldPaths which should be
		// be merged. Any fields that exist in the data but are specified will be ignored
		firestore.Merge(firestore.FieldPath{"ID"}),
	)
	if err != nil {
		panic(err)
	}
	getAndDump(docRef)

	fmt.Println()
	// START SET - MergeAll OMIT
	fmt.Println("SET - MergeAll")
	// ... but you can use merge options to allow smarter merging
	_, err = docRef.Set(ctx, map[string]interface{}{"Name": "Weave"},
		// MergeAll causes all paths to be overwritten
		// but only works with map-based data in Go
		firestore.MergeAll,
	)
	if err != nil {
		panic(err)
	}
	getAndDump(docRef)

	fmt.Println()
	// START DELETE OMIT
	fmt.Println("DELETE")
	_, err = docRef.Delete(ctx)
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
