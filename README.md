# firestore-and-go

## Utah Gophers April 2020

Carson Anderson

DevX Engineer, Weave

@carson_ops

```sh
// Get deps into your GOPATH
GO111MODULE=off go get cloud.google.com/go/firestore github.com/sanity-io/litter

// run locally
present -use_playground=false -base theme
```

## Run locally without GCP Creds using the Firestore emulator

https://firebase.google.com/docs/emulator-suite/install_and_configure

```
# Start emulator
firebase emulators:start --only firestore

# Use emulator and present
export FIRESTORE_EMULATOR_HOST="localhost:8080"
present -use_playground=false -base theme

# Open http://127.0.0.1:3999/presentation.slide
```
