package internal

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
	"google.golang.org/api/iterator"
)

type Task struct {
	Name       string
	Due        time.Time
	InProgress bool
}

func GetTask(ctx context.Context, client *firestore.Client, allItems bool) error {
	var iter *firestore.DocumentIterator
	if allItems {
		iter = client.Collection("items").Documents(ctx)
	} else {
		iter = client.Collection("items").Where("InProgress", "==", true).Documents(ctx)
	}
	defer iter.Stop()

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tIN PROGRESS\tDUE")

	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		var t Task

		if err := doc.DataTo(&t); err != nil {
			return err
		}

		inProgress := "true"

		if !t.InProgress {
			inProgress = "false"
		}

		dueStr := "XXXX-XX-XX"
		if !t.Due.IsZero() {
			dueStr = t.Due.Format("2006-01-02")
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", doc.Ref.ID, t.Name, inProgress, dueStr)

	}

	_ = w.Flush()

	return nil
}

func AddTask(ctx context.Context, client *firestore.Client, taskName string, due string) error {
	count, err := GetNumberOfTasks(ctx, client)

	if err != nil {
		return err
	}

	dueDate, err := time.Parse("2006-01-02 15:04", due)

	if err != nil {
		dueDate, err = time.Parse("2006-01-02", due)

		if err != nil {
			return err
		}
	}

	docID := fmt.Sprintf("Task %v", count+1)
	_, err = client.Collection("items").Doc(docID).Set(ctx, map[string]interface{}{
		"Name":       taskName,
		"Due":        dueDate,
		"InProgress": true,
	})

	return err
}

func GetNumberOfTasks(ctx context.Context, client *firestore.Client) (int64, error) {

	queryResult := client.Collection("items").NewAggregationQuery().WithCount("all")

	results, err := queryResult.Get(ctx)
	if err != nil {
		return 0, err
	}

	count, ok := results["all"]

	if !ok {
		return 0, nil
	}

	countValue := count.(*firestorepb.Value)

	return countValue.GetIntegerValue(), nil
}

func SetTaskDone(ctx context.Context, client *firestore.Client, docID string) error {
	_, err := client.Collection("items").Doc(docID).Set(ctx, map[string]interface{}{
		"InProgress": false,
	}, firestore.MergeAll)

	if err != nil {
		return err
	}

	return nil
}
