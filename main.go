package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	ExpectedUsername = "user"
	ExpectedPassword = "secret"
)

func main() {
	ConnectDB()

	fmt.Println("Hello world!")

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	//регистрирует нашу функцию обработчика  для обработки всех запросов GET
	//http.HandleFunc("/", handler)
	http.HandleFunc("/login", authHandler) //обработчик авторизации
	http.HandleFunc("/echo", echoHandler)  //предоставление доступа (к эхо) по авторизации
	http.HandleFunc("/addecho", addEchoHandler)
	logger.Println("Server is starting...")

	//плавное завершение работы
	server := &http.Server{ //создаём экземпляр с настраиваиваемыми тайм-аутами
		Addr:         ":8080",
		Handler:      nil,              //обработчик http
		ReadTimeout:  5 * time.Second,  //для чтения
		WriteTimeout: 10 * time.Second, //записи
		IdleTimeout:  15 * time.Second, // простоя
	}
	done := make(chan bool)         //канал для оповещения об остановке сервера
	quit := make(chan os.Signal, 1) //канал для прослушивания SIGINT и SIGTERM сигналов (Ctrl-C например)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v \n", err)
		}
		close(done)
	}()

	//вызываем server.ListenAndServe() запуск сервера
	logger.Println("Server is ready to handle requests at :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("couls not listen on :8080 %v\n", err)
	}
	<-done //сигнал, указывающий, что сервер остановился
	logger.Println("Server stopped")

}

//
//
//
//
//
//
//
//
//
//
//
////
//
//
//
//
////
//
//
//
//
////
//
//
//
//
//
/*Старый handler
func handler(w http.ResponseWriter, r *http.Request) {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Printf("Received request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	if r.Method == http.MethodGet && r.URL.Path == "/echo" {
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		log.Println(token)
		if !validateToken(token) {
			http.Error(w, "Unauthorized!", http.StatusUnauthorized)
			return
		}

		message, err := getEchomessageList()
		if err != nil {
			http.Error(w, "Error read in DB", http.StatusInternalServerError)
			return
		}

		var response string
		for _, msg := range message {
			response += fmt.Sprintf("ID: %d, Message: %s \n", msg.ID, msg.Message)
		}

		w.Write([]byte(response))
		return
	}

	if r.Method == http.MethodPost && r.URL.Path == "/login" {
		var user struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil || user.Username == "" || user.Password == "" {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		//Здесь должна быть проверка пользователя в БД
		//
		if user.Username != ExpectedUsername && user.Password != ExpectedPassword {
			http.Error(w, "Unauthorization", http.StatusUnauthorized)
			return
		}

		token, err := generateToken(user.Username)
		if err != nil {
			http.Error(w, "Could not create token", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(token))
		return

		/*body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error read query", http.StatusBadRequest)
			return
		}
		message := string(body)
		err = addEchotoDB(message)
		if err != nil {
			http.Error(w, "Error write in DB", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Message saved in DB: %s\n", message)
		return

		/
	}

	//stuff, _ := io.ReadAll(r.Body)
	//w.Write(stuff)
}*/

/*rows, err := getTestoneList()
if err != nil {
	log.Fatalf("error getting rows in tableone: %v", err)
}
*/
/*err = addRowTestone(2, "два")
if err != nil {
	log.Fatalf("Error additing rows in tableone: %v", err)
}*/

/*fmt.Println("Список всех значений в таблице testone:")
for _, row := range rows {
	fmt.Printf("ID: %d, One: %s\n", row.ID, row.One)
}*/

/*tables, err := getTable()
if err != nil {
	fmt.Println(err)
	return
}

fmt.Println("Список таблиц в БД:")
for _, table := range tables {
	fmt.Println(" - ", table)
}*/

/*package main

import (
	"fmt"
	//для управления жизненным циклом процесса завершения работы
	"fmt"
	"io"
	"log"
	"net/http"
	"os" //обработка сигналов
	//для определения сигналов, которые мы хоти прослушивать
	//"my-go-project-1/db"
)


type key int

const (
	requestIDKey key = 0
)

//добавляет уник-й идентификатор запроса к каждому запросу
func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

//регистрирует сведения о запросе (id, метод, путь, ip-адрес, UserAgent()
func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

var healthy int32

// 204 - исправно, 503 - неисправно
func healthCheck(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&healthy) == 1 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}
*/
// вызывается для каждого входящего запроса
// http.ResponseWriter - интерфейс для написания ответа
// http.Request - структура, содержащая информацию о входящем запросе

/*func handler(w http.ResponseWriter, r *http.Request) {
	//r.Method - медот HTTP(GET, POST, PUT) и др.
	//r.URL.Path - часть пути URL-адреса
	//r.RemoteAddr - IP-адрес клиента, сделавшего запрос
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Printf("Received request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	if r.Method == http.MethodGet {
		fmt.Fprintf(w, "Hello world!\n")
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	stuff, _ := io.ReadAll(r.Body)
	w.Write(stuff)
}

func main() {

	fmt.Println("Hello world!")
	//os.Stdout - направляет журналы на стандартный вывод
	//"http: " - префикс для сообщений журнала
	//log.LstdFlags - включает дату и время в каждой записи журнала
	/*logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	//регистрирует нашу функцию обработчика  для обработки всех запросов GET
	http.HandleFunc("/", handler)
	http.HandleFunc("/healthz", healthCheck)
	logger.Println("Server is starting...")


	//ген. id запроса на основе текущего времени
	nextRequestID := func() string {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	//плавное завершение работы
	server := &http.Server{ //создаём экземпляр с настраиваиваемыми тайм-аутами
		Addr:         ":8080",
		Handler:      tracing(nextRequestID)(logging(logger)(http.DefaultServeMux)), //обработчик http
		ReadTimeout:  5 * time.Second,                                               //для чтения
		WriteTimeout: 10 * time.Second,                                              //записи
		IdleTimeout:  15 * time.Second,                                              // простоя
	}
	done := make(chan bool)         //канал для оповещения об остановке сервера
	quit := make(chan os.Signal, 1) //канал для прослушивания SIGINT и SIGTERM сигналов (Ctrl-C например)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	atomic.StoreInt32(&healthy, 1)
	go func() {
		<-quit
		logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v \n", err)
		}
		close(done)
	}()

	//вызываем server.ListenAndServe() запуск сервера
	logger.Println("Server is ready to handle requests at :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("couls not listen on :8080 %v\n", err)
	}
	<-done //сигнал, указывающий, что сервер остановился
	logger.Println("Server stopped")
}
*/
