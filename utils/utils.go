package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

var HOME_DIR string = os.Getenv("HOME")
var MANGO_CONFIG_DIR string = path.Join(HOME_DIR, ".config", "mango")
var PATH_TO_TODOS string = path.Join(MANGO_CONFIG_DIR, "todos.json")

func CreateConfigFile() {
	if _, err := os.Stat(MANGO_CONFIG_DIR); err != nil {
		err := os.Mkdir(MANGO_CONFIG_DIR, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(PATH_TO_TODOS); err != nil {
		file, err := os.Create(PATH_TO_TODOS)
		if err != nil {
			log.Fatal(err)
		}
		file.Write(bytes.NewBufferString("[]").Bytes())
		fmt.Println("Initialized mango data")
	}
}

func GetTodo(id string, todos *TodoList) (*Todo, error) {
	for i := 0; i < len(*todos); i++ {
		if (*todos)[i].Id == id {
			return (*todos)[i], nil
		}
		if len(*(*todos)[i].Todos) != 0 {
			todo, err := GetTodo(id, (*todos)[i].Todos)
			if err == nil {
				return todo, nil
			}
		}
	}
	return nil, errors.New("Element not found of id: '" + id + "'")
}

func GetTodos() *TodoList {
	content, err := os.ReadFile(PATH_TO_TODOS)
	if err != nil {
		log.Fatal(err)
	}
	var todos TodoList
	err = json.Unmarshal(content, &todos)
	if err != nil {
		log.Fatal(err)
	}

	return &todos
}

func SaveTodos(todos *TodoList) {
	todos = TidyTodos(todos)
	// json byte representation
	v, err := json.Marshal(todos)
	if err != nil {
		log.Fatal(err)
	}
	var indentedData bytes.Buffer
	json.Indent(&indentedData, v, "", "    ")
	// Create / Open file
	f, err := os.Create(PATH_TO_TODOS)
	if err != nil {
		log.Fatal(err)
	}
	// Save to file
	f.Write(indentedData.Bytes())
}

/*
Add todo to the slice provided, and returns a copy of the slice with that todo
*/
func AddTodo(todo *Todo, todos *TodoList) *TodoList {
	newTodos := make(TodoList, len(*todos))
	copy(newTodos, *todos)
	newTodos = append(newTodos, todo)
	return &newTodos
}

func AssignId(parentId string) string { // TODO clean and optimize
	todos := GetTodos()
	if len(parentId) == 0 {
		// get ID list
		idList := make([]string, 0)
		for i := 0; i < len(*todos); i++ {
			idList = append(idList, (*todos)[i].Id)
		}

		// Sort it
		sort.Slice(idList, func(i, j int) bool {
			lastIdInt_i, err := strconv.Atoi(idList[i])
			if err != nil {
				log.Fatal(err)
			}
			lastIdInt_j, err := strconv.Atoi(idList[j])
			if err != nil {
				log.Fatal(err)
			}
			return lastIdInt_i < lastIdInt_j
		})

		// Check missing id
		for i := 0; i < len(idList)-1; i++ {
			id1, err := strconv.Atoi(idList[i])
			if err != nil {
				log.Fatal(err)
			}
			id2, err := strconv.Atoi(idList[i+1])
			if err != nil {
				log.Fatal(err)
			}
			if math.Abs(float64(id1-id2)) > 1 {
				return strconv.Itoa(int(math.Min(float64(id1), float64(id2)) + 1))
			}
		}
		return strconv.Itoa(len(idList) + 1)

	} else {
		// Search for the parent
		parentTodo, err := GetTodo(parentId, todos)
		if err != nil {
			log.Fatal(err)
		}
		todos := (*parentTodo).Todos

		// get ID list
		idList := make([]string, 0)
		for i := 0; i < len(*todos); i++ {
			idList = append(idList, (*todos)[i].Id)
		}

		// Sort it
		sort.Slice(idList, func(i, j int) bool {
			splits_i := strings.Split(idList[i], "-")
			lastId_i := splits_i[len(splits_i)-1]
			lastIdInt_i, err := strconv.Atoi(lastId_i)
			if err != nil {
				log.Fatal(err)
			}

			splits_j := strings.Split(idList[j], "-")
			lastId_j := splits_j[len(splits_j)-1]
			lastIdInt_j, err := strconv.Atoi(lastId_j)
			if err != nil {
				log.Fatal(err)
			}

			return lastIdInt_i < lastIdInt_j
		})
		for i := 0; i < len(idList)-1; i++ {
			splits_1 := strings.Split(idList[i], "-")
			lastId_1 := splits_1[len(splits_1)-1]
			id1, err := strconv.Atoi(lastId_1)
			if err != nil {
				log.Fatal(err)
			}
			splits_2 := strings.Split(idList[i], "-")
			lastId_2 := splits_2[len(splits_2)-1]
			id2, err := strconv.Atoi(lastId_2)
			if err != nil {
				log.Fatal(err)
			}
			if math.Abs(float64(id1-id2)) > 1 {
				return parentTodo.Id + "-" + strconv.Itoa(int(math.Min(float64(id1), float64(id2))+1))
			}
		}
		return parentTodo.Id + "-" + strconv.Itoa(len(idList)+1)
	}
}

func ChangeTodoState(state bool, todo *Todo) {
	todo.IsDone = state
}

func DeleteTodo(todoId string, todos *TodoList) *TodoList {
	newTodos := make(TodoList, 0)
	for i := 0; i < len(*todos); i++ {
		todo := (*todos)[i]
		if todo.Id == todoId {
			continue
		}
		if len(*todo.Todos) != 0 {
			todo.Todos = DeleteTodo(todoId, todo.Todos)
		}
		newTodos = append(newTodos, todo)
	}
	return &newTodos
}

func TidyTodos(todos *TodoList) *TodoList {
	todos = todos.SortByDeadline()
	newTodos := make(TodoList, 0)
	for i := 0; i < len(*todos); i++ {
		todo := (*todos)[i]
		splits := strings.Split(todo.Id, "-")
		prefixPart := strings.Join(splits[:len(splits)-1], "-")
		if len(prefixPart) != 0 {
			prefixPart = prefixPart + "-"
		} else {
			prefixPart = ""
		}
		todo.Id = prefixPart + strconv.Itoa(i+1)
		if len(*todo.Todos) != 0 {
			todo.Todos = TidyTodos(todo.Todos)
		}
		newTodos = append(newTodos, todo)
	}
	return &newTodos
}
