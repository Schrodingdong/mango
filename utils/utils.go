package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

var PATH_TO_TODOS string = "./todos.json"

func GetTodo(id string) (*Todo, error) {
	todos := GetTodos()
	todo, err := getTodoRecursive(id, todos)
	if err != nil {
		return nil, err
	} else {
		return todo, nil
	}
}

func getTodoRecursive(id string, todos *TodoList) (*Todo, error) {
	for i := 0; i < len(*todos); i++ {
		if (*todos)[i].Id == id {
			return (*todos)[i], nil
		}
		if len(*(*todos)[i].Todos) != 0 {
			todo, err := getTodoRecursive(id, (*todos)[i].Todos)
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

func saveTodos(todos *TodoList) {
	// json byte representation
	v, err := json.Marshal(todos)
	if err != nil {
		log.Fatal(err)
	}
	var indentedData bytes.Buffer
	json.Indent(&indentedData, v, "", "    ")
	// Create / Open file
	f, err := os.Create("./todos.json")
	if err != nil {
		log.Fatal(err)
	}
	// Save to file
	f.Write(indentedData.Bytes())
}

func AddTodo(todo *Todo, parentId string) {
	todos := GetTodos()
	if len(parentId) == 0 {
		newTodos := append(*todos, todo)
		todos = &newTodos
	} else {
		parentTodo, err := getTodoRecursive(parentId, todos)
		if err != nil {
			log.Fatal(err)
		}
		parentTodo.AddTodo(todo)
	}
	saveTodos(todos)
}

func AssignId(parentId string) string { // TODO clean and optimize
	if len(parentId) == 0 {
		todos := GetTodos()

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
		parentTodo, err := GetTodo(parentId)
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

func ChangeTodoState(state bool, todoId string) {
	todos := GetTodos()
	todo, err := getTodoRecursive(todoId, todos)
	if err != nil {
		log.Fatal(err)
	}
	todo.IsDone = state
	saveTodos(todos)
}
