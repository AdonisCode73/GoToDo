package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"todoproj/internal"
)

const (
	ADD  = "add"
	LIST = "list"
	DONE = "done"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Refer to manual for progam usage")
	}

	args := os.Args

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, client, cleanUp, err := internal.NewFirestore(ctx)
	defer cleanUp()

	if err != nil {
		log.Fatal(err)
	}

	switch args[1] {
	case LIST:
		listCmd := flag.NewFlagSet("list", flag.ExitOnError)
		list := listCmd.Bool("all", false, "List all tasks")

		listCmd.Parse(os.Args[2:])

		err := internal.GetTask(ctx, client, *list)

		if err != nil {
			log.Fatal(err)
		}

	case ADD:
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		taskName := addCmd.String("name", "", "Name of the task")
		due := addCmd.String("due", "", "Due date/time for the task (YYYY-MM-DD) OR (YYYY-MM-DD HH:MM)")

		addCmd.Parse(os.Args[2:])

		if *taskName == "" {
			log.Fatal("Error: task name required")
		}

		if *due == "" {
			log.Fatal("Error: due date required")
		}

		err := internal.AddTask(context.Background(), client, *taskName, *due)

		if err != nil {
			log.Fatal(err)
		}

	case DONE:
		doneCmd := flag.NewFlagSet("done", flag.ExitOnError)
		docID := doneCmd.String("docID", "", "ID of the document")

		doneCmd.Parse(os.Args[2:])

		if *docID == "" {
			log.Fatal("Error: document ID required")
		}

		internal.SetTaskDone(ctx, client, *docID)

	default:
		log.Fatal("Refer to manual for progam usage")
	}

}
