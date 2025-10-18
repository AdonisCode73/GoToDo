package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"todoproj/internal"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Data struct {
	ID   string
	Name string
	Due  string
}

func GetTop3Tasks(ctx context.Context, client *firestore.Client) ([]Data, error) {
	var iter = client.Collection("items").Where("InProgress", "==", true).
		OrderBy("Due", firestore.Asc).Limit(3).Documents(ctx)

	defer iter.Stop()

	var retval []Data
	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		var t internal.Task
		if err := doc.DataTo(&t); err != nil {
			return nil, err
		}

		dueStr := "XXXX-XX-XX"
		if !t.Due.IsZero() {
			dueStr = t.Due.Format("2006-01-02")
		}

		retval = append(retval, Data{
			ID:   doc.Ref.ID,
			Name: t.Name,
			Due:  dueStr,
		})
	}

	return retval, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, client, cleanUp, err := internal.NewFirestore(ctx)

	if err != nil {
		log.Fatal(err)
	}
	defer cleanUp()

	mux := http.NewServeMux()
	mux.HandleFunc("/top3", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		tasks, err := GetTop3Tasks(ctx, client)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(tasks)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, mux))

}
