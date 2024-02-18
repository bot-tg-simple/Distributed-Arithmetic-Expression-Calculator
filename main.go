package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Knetic/govaluate"
	_ "github.com/Knetic/govaluate"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Operation struct {
	Name      string `db:"name"`
	Duration  int    `db:"duration"`
	StartTime time.Time
	Status    string `db:"status"`
}

type Resource struct {
	Name      string `db:"name"`
	Operation string `db:"operation"`
	Duration  int    `db:"duration"`
}

type Task struct {
	ID         string
	Expression string
}

type Result struct {
	ID     string
	Result float64
}

type Orchestrator struct {
	db            *sqlx.DB
	taskQueue     chan Task
	resultQueue   chan Result
	operations    []Operation
	resources     []Resource
	workers       int
	workerWg      sync.WaitGroup
	shutdown      chan struct{}
	shutdownWg    sync.WaitGroup
	expressionMtx sync.Mutex
	timeout       time.Duration
	lastPing      time.Time
}

type Expression struct {
	ID         string    `db:"id"`
	Expression string    `db:"expression"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Result     *float64  `db:"result"`
}

func NewOrchestrator(db *sqlx.DB, operations []Operation, resources []Resource, workers int, timeout time.Duration) *Orchestrator {
	return &Orchestrator{
		db:          db,
		taskQueue:   make(chan Task),
		resultQueue: make(chan Result),
		operations:  operations,
		resources:   resources,
		workers:     workers,
		timeout:     timeout,
		lastPing:    time.Now(),
		shutdown:    make(chan struct{}),
	}
}

func (o *Orchestrator) Start() {
	o.createTables()

	for i := 0; i < o.workers; i++ {
		o.workerWg.Add(1)
		go o.worker()
	}

	go o.processTasks()
	go o.processResults()
}

func (o *Orchestrator) Stop() {
	close(o.shutdown)
	o.workerWg.Wait()
	close(o.resultQueue)
	o.shutdownWg.Wait()
}

func (o *Orchestrator) createTables() {
	_, err := o.db.Exec(`
CREATE TABLE IF NOT EXISTS operations (
name TEXT PRIMARY KEY,
duration INTEGER NOT NULL
);
`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = o.db.Exec(`
CREATE TABLE IF NOT EXISTS expressions (
id TEXT PRIMARY KEY,
expression TEXT NOT NULL,
status TEXT NOT NULL,
created_at DATETIME NOT NULL,
updated_at DATETIME NOT NULL,
result REAL NULL
);
`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = o.db.Exec(`
CREATE TABLE IF NOT EXISTS resources (
name TEXT PRIMARY KEY,
operation TEXT NOT NULL
);
`)
	if err != nil {
		log.Fatal(err)
	}
}

func (o *Orchestrator) AddExpression(expression string) (string, error) {
	if expression == "" {
		return "", fmt.Errorf("expression is required")
	}

	id := uuid.New().String()

	stmt, err := o.db.Prepare("INSERT INTO expressions (id, expression, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return "", fmt.Errorf("failed to create expression: %v", err)
	}
	_, err = stmt.Exec(id, expression, "pending", time.Now(), time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to create expression: %v", err)
	}

	return id, nil
}

func (o *Orchestrator) GetExpressions() ([]Expression, error) {
	var expressions []Expression
	err := o.db.Select(&expressions, "SELECT id, expression, status, created_at, updated_at, result FROM expressions")
	if err != nil {
		return nil, fmt.Errorf("failed to get expressions: %v", err)
	}

	return expressions, nil
}

func (o *Orchestrator) GetExpressionByID(id string) (*Expression, error) {
	var expression Expression
	err := o.db.Get(&expression, "SELECT * FROM expressions WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get expression: %v", err)
	}

	return &expression, nil
}

func (o *Orchestrator) GetOperations() ([]Operation, error) {
	return o.operations, nil
}

func (o *Orchestrator) GetResources() ([]Resource, error) {
	return o.resources, nil
}

func (o *Orchestrator) GetTask(ctx context.Context) (*Task, error) {
	var expression Expression
	err := o.db.GetContext(ctx, &expression, "SELECT * FROM expressions WHERE status = 'running' LIMIT 1")
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	return &Task{
		ID:         expression.ID,
		Expression: expression.Expression,
	}, nil
}

func (o *Orchestrator) UpdateExpressionResult(id string, result float64) error {
	_, err := o.db.Exec("UPDATE expressions SET status = 'completed', result = ? WHERE id = ?", result, id)
	if err != nil {
		return fmt.Errorf("failed to update expression: %v", err)
	}

	return nil
}

func (o *Orchestrator) worker() {
	defer o.workerWg.Done()

	for {
		select {
		case task := <-o.taskQueue:
			for i, op := range o.operations {
				if op.Name == task.Expression {
					o.operations[i].StartTime = time.Now()
					o.operations[i].Status = "running"
					break
				}
			}

			result, err := o.evaluateExpression(task.Expression)
			if err != nil {
				log.Println("Failed to evaluate expression:", err)
				continue
			}

			o.resultQueue <- Result{
				ID:     task.ID,
				Result: result,
			}

			for i, op := range o.operations {
				if op.Name == task.Expression {
					o.operations[i].StartTime = time.Time{}
					o.operations[i].Status = "completed"
					break
				}
			}
		case <-o.shutdown:
			return
		}
	}
}

func (o *Orchestrator) processTasks() {
	for {
		select {
		case <-o.shutdown:
			close(o.taskQueue)
			return
		default:
			task, err := o.GetTask(context.Background())
			if err != nil {
				log.Println("Failed to get task:", err)
				time.Sleep(1 * time.Second)
				continue
			}

			if task == nil {
				time.Sleep(1 * time.Second)
				continue
			}

			o.taskQueue <- *task
		}
	}
}

func (o *Orchestrator) processResults() {
	o.shutdownWg.Add(1)
	defer o.shutdownWg.Done()

	for result := range o.resultQueue {
		err := o.UpdateExpressionResult(result.ID, result.Result)
		if err != nil {
			log.Println("Failed to update expression result:", err)
		}
	}
}

func (o *Orchestrator) evaluateExpression(expression string) (float64, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return 0, err
	}

	result, err := expr.Evaluate(nil)
	if err != nil {
		return 0, err
	}

	resultFloat, ok := result.(float64)
	if !ok {
		return 0, fmt.Errorf("invalid result type")
	}

	// Проверка на деление на ноль
	if resultFloat == math.Inf(1) || resultFloat == math.Inf(-1) {
		return 0, fmt.Errorf("division by zero")
	}

	return resultFloat, nil
}

func (orchestrator *Orchestrator) ping(c *gin.Context) {
	orchestrator.lastPing = time.Now()

	go func() {
		for {
			time.Sleep(1 * time.Minute)

			if time.Since(orchestrator.lastPing) > 1*time.Minute {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Connection to the server is lost"})
				break
			}
		}
	}()
}

func (o *Orchestrator) GetRunningOperations() []Operation {
	var runningOperations []Operation

	for _, op := range o.operations {
		if op.Status == "running" {
			runningOperations = append(runningOperations, op)
		}
	}

	return runningOperations
}

func main() {
	db, err := sqlx.Open("sqlite3", "./expressions.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	operations := []Operation{
		{Name: "+", Duration: 5},
		{Name: "-", Duration: 10},
		{Name: "*", Duration: 15},
		{Name: "/", Duration: 20},
	}

	resources := []Resource{
		{Name: "Resource1", Operation: "+"},
		{Name: "Resource2", Operation: "-"},
		{Name: "Resource3", Operation: "*"},
		{Name: "Resource4", Operation: "/"},
	}

	orchestrator := NewOrchestrator(db, operations, resources, 2, 10*time.Second)
	orchestrator.Start()

	router.POST("/expression", func(c *gin.Context) {
		orchestrator.ping(c)

		expression := c.PostForm("expression")

		id, err := orchestrator.AddExpression(expression)
		if err != nil {
			if err.Error() == "division by zero" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Division by zero is not allowed"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		go func() {
			task := Task{
				ID:         id,
				Expression: expression,
			}
			orchestrator.taskQueue <- task
		}()

		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	router.GET("/expressions", func(c *gin.Context) {
		orchestrator.ping(c)

		expressions, err := orchestrator.GetExpressions()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, expressions)
	})

	router.GET("/expression/:id", func(c *gin.Context) {
		orchestrator.ping(c)

		id := c.Param("id")

		expression, err := orchestrator.GetExpressionByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, expression)
	})

	router.GET("/operations", func(c *gin.Context) {
		orchestrator.ping(c)

		operations, err := orchestrator.GetOperations()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, operations)
	})

	router.GET("/resources", func(c *gin.Context) {
		orchestrator.ping(c)

		resources, err := orchestrator.GetResources()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resources)
	})

	router.GET("/task", func(c *gin.Context) {
		orchestrator.ping(c)

		runningOperations := orchestrator.GetRunningOperations()
		if len(runningOperations) == 0 {
			c.JSON(http.StatusNoContent, gin.H{"message": "No running operations"})
			return
		}

		c.JSON(http.StatusOK, runningOperations)
	})

	router.POST("/result", func(c *gin.Context) {
		orchestrator.ping(c)

		id := c.PostForm("id")
		resultStr := c.PostForm("resultStr")

		if id == "" || resultStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid result data"})
			return
		}
		result, err := strconv.ParseFloat(resultStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid result data"})
			return
		}

		err = orchestrator.UpdateExpressionResult(id, result)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Result updated successfully"})
	})

	router.POST("/operation-duration", func(c *gin.Context) {
		orchestrator.ping(c)

		operation := c.PostForm("operation")
		durationStr := c.PostForm("duration")
		if operation == "" || durationStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Operation and duration are required"})
			return
		}

		duration, err := strconv.Atoi(durationStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid duration"})
			return
		}

		for i, op := range orchestrator.operations {
			if op.Name == operation {
				orchestrator.operations[i].Duration = duration
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Operation duration updated successfully"})
	})

	go func() {
		for {
			time.Sleep(1 * time.Minute)

			if time.Since(orchestrator.lastPing) > 1*time.Minute {
				router.GET("/ping", func(c *gin.Context) {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Connection to the server is lost: the item will stop displaying after 1 minute"})
					c.Writer.Flush()
					os.Exit(1)
				})
				log.Println("Connection to the server is lost: the item will stop displaying after 1 minute")
				break
			}
		}
	}()

	router.Run(":8080")
}
