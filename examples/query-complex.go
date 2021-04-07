package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/genproto/googleapis/type/latlng"
)

type Customer struct {
	ID   string
	Name string

	CreateDate time.Time

	// Location is a point on the earth
	// it MUST be a pointer to saved as the right type in Firestore
	Location *latlng.LatLng
}

func main() {
	ctx := context.Background()

	os.Unsetenv("FIRESTORE_EMULATOR_HOST")

	client, cerr := firestore.NewClient(ctx, "firestore-and-go")
	if cerr != nil {
		panic(cerr)
	}

	//pre-add
	client.Collection("customers").Doc("weave").Set(ctx, &Customer{ID: "weave", Name: "Weave", CreateDate: time.Now()})
	client.Collection("customers").Doc("firestore-and-go").Set(ctx, &Customer{ID: "firestore-and-go", CreateDate: time.Now(), Name: "Firestore And Go"})

	// START QUERY OMIT
	query := client.Collection("customers").Query.
		Where("ExpireDate", ">", time.Now()).
		Where("OnboardDate", "<", time.Now().AddDate(-1, 0, 0)).
		Limit(2)

	_, err := query.Documents(ctx).GetAll()
	fmt.Println(err.Error())

	query = client.Collection("customers").Query.
		Where("ExpireDate", ">", time.Now()).
		OrderBy("OnboardDate", firestore.Asc).
		Limit(2)
	_, err = query.Documents(ctx).GetAll()
	fmt.Println(err.Error())

	query = client.Collection("customers").Query.
		Where("ExpireDate", "==", time.Now()).
		OrderBy("OnboardDate", firestore.Asc).
		Limit(2)

	_, err = query.Documents(ctx).GetAll()
	fmt.Println(err.Error())
	// END QUERY OMIT

	// START QUERY VALID OMIT
	client.Collection("customers").Query.
		Where("First", "==", "Carson").
		Where("Last", "==", "Carson")
	// END QUERY VALID OMIT
}
