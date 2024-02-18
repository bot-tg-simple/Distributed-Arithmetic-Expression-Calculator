Это схема приложения

#### Structs:

- **Operation**:
  - Name: string
  - Duration: int
  - StartTime: time.Time
  - Status: string

- **Resource**:
  - Name: string
  - Operation: string
  - Duration: int

- **Task**:
  - ID: string
  - Expression: string

- **Result**:
  - ID: string
  - Result: float64

- **Expression**:
  - ID: string
  - Expression: string
  - Status: string
  - CreatedAt: time.Time
  - UpdatedAt: time.Time
  - Result: *float64

#### Methods:

- **NewOrchestrator(db *sqlx.DB, operations []Operation, resources []Resource, workers int, timeout time.Duration) *Orchestrator**:
  - Создает новый экземпляр Orchestrator с указанными параметрами.

- **Start()**:
  - Запускает Orchestrator и инициализирует таблицы в базе данных.

- **Stop()**:
  - Останавливает работу Orchestrator и завершает все рабочие горутины.

- **AddExpression(expression string) (string, error)**:
  - Добавляет новое выражение в базу данных и возвращает его ID.

- **GetExpressions() ([]Expression, error)**:
  - Получает все выражения из базы данных.

- **GetExpressionByID(id string) (*Expression, error)**:
  - Получает выражение по ID из базы данных.

- **GetOperations() ([]Operation, error)**:
  - Получает список операций.

- **GetResources() ([]Resource, error)**:
  - Получает список ресурсов.

- **GetTask(ctx context.Context) (*Task, error)**:
  - Получает задачу для выполнения из базы данных.

- **UpdateExpressionResult(id string, result float64) error**:
  - Обновляет результат выполнения выражения в базе данных.

- **worker()**:
  - Рабочая горутина, выполняющая задачи.

- **processTasks()**:
  - Обрабатывает задачи и отправляет их на выполнение.

- **processResults()**:
  - Обрабатывает результаты выполнения задач.

- **evaluateExpression(expression string) (float64, error)**:
  - Вычисляет результат выражения.

- **ping(c *gin.Context)**:
  - Проверяет соединение с сервером.

- **GetRunningOperations() []Operation**:
  - Получает список операций, которые находятся в статусе "running".

### API Endpoints:

1. **POST /expression**:
   - Добавляет новое выражение для выполнения.

2. **GET /expressions**:
   - Получает все выражения.

3. **GET /expression/:id**:
   - Получает конкретное выражение по ID.

4. **GET /operations**:
   - Получает список операций.

5. **GET /resources**:
   - Получает список ресурсов.

6. **GET /task**:
   - Получает текущие выполняющиеся операции.

7. **POST /result**:
   - Обновляет результат выполнения выражения.

8. **POST /operation-duration**:
   - Обновляет длительность операции.

### Ping Mechanism:

- **Ping Mechanism**:
  - Проверяет соединение с сервером и обновляет время последнего пинга.

### Server:

- **Server**:
  - Запускает сервер на порту 8080 и обрабатывает все API-запросы.

### Схема работы:

1. **Инициализация**:
   - Создание экземпляра Orchestrator с базой данных, операциями, ресурсами, работниками и таймаутом.
   - Запуск Orchestrator.

2. **API-запросы**:
   - Обработка запросов на добавление выражения, получение выражений, операций, ресурсов, текущих операций и обновление результатов выполнения.

3. **Рабочие горутины**:
   - Рабочая горутина выполняет задачи, вычисляет результаты и обновляет статус операций.

4. **Обработка результатов**:
   - Результаты выполнения задач обновляются в базе данных.

5. **Ping Mechanism**:
   - Проверка соединения с сервером и обновление времени последнего пинга.

6. **Остановка**:
   - Остановка Orchestrator и завершение работы всех горутин.
