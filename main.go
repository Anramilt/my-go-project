package main

import (
	"context" //для управления жизненным циклом процесса завершения работы
	"fmt"
	"log"
	"net/http"
	"os" //обработка сигналов
	"os/signal"
	"strconv"
	"syscall" //для определения сигналов, которые мы хоти прослушивать
	"time"
)

type key int

const (
	requestIDKey key = 0
)

/*добавляет уник-й идентификатор запроса к каждому запросу*/
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

/*регистрирует сведения о запросе (id, метод, путь, ip-адрес, UserAgent()*/
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

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// вызывается для каждого входящего запроса
// http.ResponseWriter - интерфейс для написания ответа
// http.Request - структура, содержащая информацию о входящем запросе
func handler(w http.ResponseWriter, r *http.Request) {
	//r.Method - медот HTTP(GET, POST, PUT) и др.
	//r.URL.Path - часть пути URL-адреса
	//r.RemoteAddr - IP-адрес клиента, сделавшего запрос
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Printf("Received request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	fmt.Fprintf(w, "Hello world!\n")
}

func main() {
	//os.Stdout - направляет журналы на стандартный вывод
	//"http: " - префикс для сообщений журнала
	//log.LstdFlags - включает дату и время в каждой записи журнала
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	//регистрирует нашу функцию обработчика  для обработки всех запросов GET
	http.HandleFunc("GET /", handler)
	http.HandleFunc("GET /healthz", healthCheck)
	logger.Println("Server is starting...")
	/*err := http.ListenAndServe(":8080", nil) //запускает сервер на порту
	if err != nil {
		logger.Fatal("ListenAndServe", err)
	}*/

	//ген. id запроса на основе текущего времени
	nextRequestID := func() string {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	//плавное завершение работы
	server := &http.Server{ //создаём экземпляр с настраиваиваемыми тайм-аутами
		Addr:         ":8080",
		Handler:      tracing(nextRequestID)(logging(logger)(http.DefaultServeMux)), //обработчик http
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
