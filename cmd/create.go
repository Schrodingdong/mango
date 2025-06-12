package cmd

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/schrodingdong/mango/utils"
	"github.com/spf13/cobra"
)

var description string
var deadlineString string
var parentId string

func init() {
	// Flags
	createCmd.PersistentFlags().StringVar(&description, "description", "", "Todo description")
	createCmd.PersistentFlags().StringVar(&deadlineString, "deadline", "", "Absolute deadline (DateTime format, e.g., '2025-05-24 21:41:23')")
	createCmd.PersistentFlags().StringVar(&parentId, "parent", "", "Parent todo id")
	createCmd.PersistentFlags().BoolVar(&isDone, "done", false, "Todo status")
	rootCmd.AddCommand(createCmd)
}

var zero time.Time

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new todo",
	Run: func(cmd *cobra.Command, args []string) {
		todos := utils.GetTodos()
		// Get title from args
		if len(args) == 0 {
			log.Fatal("Not enough args")
		}
		title := args[0]

		// Add deadline
		todoDeadline, err := parseDeadline(deadlineString)
		if err != nil {
			log.Fatal(err)
		}
		if time.Now().After(todoDeadline) {
			log.Fatal("Can't have a deadline before now")
		}

		// Init todo data
		todo := utils.Todo{
			Id:          utils.AssignId(parentId, todos),
			Title:       title,
			Description: description,
			Deadline:    todoDeadline,
			IsDone:      isDone,
			Todos:       &utils.TodoList{},
		}

		// Add todo
		if len(parentId) == 0 {
			todos = utils.AddTodo(&todo, todos)
		} else {
			parentTodo, err := utils.GetTodo(parentId, todos)
			if err != nil {
				log.Fatal(err)
			}
			parentTodo.Todos = utils.AddTodo(&todo, parentTodo.Todos)
			parentTodo.AdaptDeadlines()
		}
		utils.SaveTodos(todos)

		// Print added todo details
		todo.PrintTodoDetail()
	},
}

/*
Parse a dedadline string into a `time.Time` object. if an empty string is provided, we get t=0ns.
*/
func parseDeadline(deadlineString string) (time.Time, error) {
	re := regexp.MustCompile(`(\d+d)?(\d+h)?(\d+m)?(\d+s)?`)

	// check if its dmhs format or normale timedate
	m := re.FindString(deadlineString)
	if len(m) == 0 {
		deadline, err := time.ParseInLocation(time.DateTime, deadlineString, time.Local)
		if err != nil {
			return zero, errors.New("wrong timedate syntax. Use either 'xdxhxmxs' or 'yyyy:mm:dd hh:mm:ss'")
		}
		return deadline, nil
	}
	matches := re.FindAllStringSubmatch(deadlineString, -1)
	reformatedTime := ""
	for i := 0; i < len(matches); i++ {
		reformatedTime += matches[i][0]
	}
	matches = re.FindAllStringSubmatch(reformatedTime, -1)
	if len(matches) > 1 {
		return zero, errors.New("wrong timedate syntax. make sure to follow this order 'xdxhxmxs'")
	}
	var (
		days    int
		hours   int
		minutes int
		seconds int
	)
	for i := 1; i < len(matches[0]); i++ {
		v := matches[0][i]
		if len(v) == 0 {
			continue
		}
		code := v[len(v)-1]
		value, err := strconv.Atoi(v[:len(v)-1])
		if err != nil {
			fmt.Println(err)
		}
		switch code {
		case 'd':
			days = value
		case 'h':
			hours = value
		case 'm':
			minutes = value
		case 's':
			seconds = value
		default:
			return zero, errors.New("wrong timedate syntax. Use either 'xdxhxmxs' or 'yyyy:mm:dd hh:mm:ss'")
		}
	}
	deadline := getFutureTimestamp(days, hours, minutes, seconds)
	return deadline, nil
}

func getFutureTimestamp(days int, hours int, minutes int, seconds int) time.Time {
	SEC_TO_NANO := 1_000_000_000
	MIN_TO_NANO := 60 * SEC_TO_NANO
	HOUR_TO_NANO := 60 * MIN_TO_NANO
	DAY_TO_NANO := 24 * HOUR_TO_NANO
	duration :=
		seconds*SEC_TO_NANO +
			minutes*MIN_TO_NANO +
			hours*HOUR_TO_NANO +
			days*DAY_TO_NANO
	return time.Now().Add(time.Duration(duration))
}
