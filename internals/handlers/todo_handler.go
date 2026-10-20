package handlers

import (
	"net/http"
	"strconv"
	"todo_with_auth/internals/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateTodoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

type UpdateTodoInput struct {
	/*  Here the payload fields should be specifically of type pointer as in the cases where only of the value is provided (title or completed),
	    the Go compiler will consider the missing value and replace them by their default values.
		So if title is missing in the payload, it will consider title as empty string or
		if completed is missing it will consider it as false - the default value provided in migration file.
		Using the pointers instead of simple types will prevent this by asserting the fact that these values are not at all provided with the help of nil pointer
	*/
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
	Id        int     `json:"id"`
}

func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value, exists := ctx.Get("user_id")
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find the user"})
			return
		}

		userID := value.(string)
		var input CreateTodoInput

		err := ctx.ShouldBindJSON(&input)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo, err := repositories.CreateTodo(pool, input.Title, input.Completed, userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, todo)
	}
}

func GetTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value, exists := ctx.Get("user_id")
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find the user"})
			return
		}

		userID := value.(string)
		todos, err := repositories.GetTodos(pool, userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusOK, todos)
	}
}

func GetTodoByIdHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value, exists := ctx.Get("user_id")
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find the user"})
			return
		}

		userID := value.(string)

		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Todo Id"})
			return
		}

		todo, err := repositories.GetTodoById(pool, id, userID)
		if err != nil {
			if err == pgx.ErrNoRows {
				// ctx.JSON(http.StatusOK, gin.H{})
				ctx.JSON(http.StatusNotFound, gin.H{"error": "No Todo found against the given Id"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, todo)
	}
}

func UpdateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value, exists := ctx.Get("user_id")
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find the user"})
			return
		}

		userID := value.(string)
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Todo Id"})
			return
		}

		var input UpdateTodoInput
		err = ctx.ShouldBindJSON(&input)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Title == nil && input.Completed == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a Title or Completed"})
			return
		}

		todo, err := repositories.GetTodoById(pool, id, userID)
		if err != nil {
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "No Todo found against the given Id"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if input.Completed == nil {
			input.Completed = &todo.Completed
		}
		if input.Title == nil {
			input.Title = &todo.Title
		}

		todo, err = repositories.UpdateTodo(pool, id, *input.Title, *input.Completed)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, todo)
	}
}

func DeleteTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value, exists := ctx.Get("user_id")
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find the user"})
			return
		}

		userID := value.(string)

		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid id"})
			return
		}

		err = repositories.DeleteTodo(pool, id, userID)
		if err != nil {
			if err.Error() == "Todo with the given id not found" {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
	}
}
