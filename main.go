package main

import (
	"log"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type todo struct {
	ID        string `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

var todos = []todo{
	{ID: "1", Item: "Getting a space suit", Completed: false},
	{ID: "2", Item: "Getting a spaceship", Completed: false},
	{ID: "3", Item: "Going to Mars", Completed: false},
}

func main() {
	router := gin.Default()

	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodo)
	router.POST("/todos", addTodo)
	router.PUT("/todo/:id", updateTodo)
	router.DELETE("/todo/:id", deleteTodo)

	router.Run(":9000")
}

func addTodo(context *gin.Context) {
	var newTodo = []todo{}

	err := context.BindJSON(&newTodo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todos = append(todos, newTodo...)

	context.IndentedJSON(http.StatusCreated, newTodo)
}

func getTodos(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, todos)
}

func getTodo(context *gin.Context) {
	id := context.Param("id")
	todo, err := getTodoIndexById(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
		return
	}

	context.IndentedJSON(http.StatusOK, todos[todo])
}

func getTodoIndexById(id string) (int, error) {
	for i, t := range todos {
		if t.ID == id {
			return i, nil
		}
	}

	return -1, errors.New("Todo not found")
}

func updateTodo(context *gin.Context) {
	id := context.Param("id")

	log.Println("id", id)

	todoIndex, err := getTodoIndexById(id)
	todo := &todos[todoIndex]
	log.Println("todo", todo, todoIndex);
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
		return
	}

	if err := context.BindJSON(todo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.IndentedJSON(http.StatusOK, todo)
}

func deleteTodo(context *gin.Context) {
	id := context.Param("id")

	index, err := getTodoIndexById(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
		return
	}

	todos = append(todos[:index], todos[index+1:]...)

	context.JSON(http.StatusOK, gin.H{"message": "Todo deleted!"})
}
