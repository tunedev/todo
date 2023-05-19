package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tunedev/todo"
)

var todoFileName = ".todo.json"

func main() {
	// use custom file name if present in the env Variable
	if os.Getenv("TODO_FILENAME") != ""{
		todoFileName = os.Getenv("TODO_FILENAME")
	}
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. \nDeveloped for learning go by building powerful cammand line tools\n", os.Args[0],
		)
		fmt.Fprintf(flag.CommandLine.Output(),"Copyright 2023\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "additionaly the add flag can simple be used by calling the %s -add, then return key you would easily enter a new line", os.Args[0])
	}

	add := flag.Bool("add", false, "Add Task to the ToDo list")
	list := flag.Bool("list", false, "List all uncompleted tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	del := flag.Int("del", 0, "Item to be deleted")
	v := flag.Bool("v", false, "List all task verbose")
	V := flag.Bool("V", false, "List all task verbose")
	u := flag.Bool("u", false, "List only uncompleted tasks")

	flag.Parse()

	l := &todo.List{}

	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *list:
		// List current todo items
		fmt.Print(l)
	case *v:
	case *V:
		// List current todo items
		result := ""
		for k, t := range *l {
			prefix := "[ ] "
			if t.Done {
				prefix = "[X] "
			}
			result += fmt.Sprintf("%s%d: %s, %v\n", prefix, k+1, t.Task, t.CreatedAt)
		}
		fmt.Print(result)
	case *u:
		// List current todo items
		result := ""
		for k, t := range *l {
			prefix := "[ ] "
			if t.Done {
				continue
			}
			result += fmt.Sprintf("%s%d: %s\n", prefix, k+1, t.Task)
		}
		fmt.Print(result)
	case *complete > 0:
		// Complete the given item
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *del > 0:
		if err := l.Delete(*del); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *add:
		// Add new task
		task, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		l.Add(task)

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case len(os.Args) == 1:
		fmt.Print(l)
		fmt.Println()
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	default:
		fmt.Fprintln(os.Stderr,"invalid option")
		os.Exit(1)
	}
}

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, ""), nil
	}
	s := bufio.NewScanner(r)
	s.Scan()
	if err := s.Err(); err != nil {
		return "", err
	}

	if len(s.Text()) == 0 {
		return "", fmt.Errorf("Task cannot be blank")
	}

	return s.Text(), nil
}