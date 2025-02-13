package main

import (
	"fmt"
	//"io"
	//"net"
	"net/http"
	//"os"
)

func main() {
	//fmt.Println("Hello, Go!")
	//fmt.Println("Hello, Go!")
	//fmt.Println("Hello, Go!")
	/*httpRequest := "GET / HTTP/1.1\n" + "Host: golang.org\n\n"
	conn, err := net.Dial("tcp", "golang.org:80")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	//посылает данные
	if _, err = conn.Write([]byte(httpRequest)); err != nil {
		fmt.Println(err)
		return
	}

	io.Copy(os.Stdout, conn)
	fmt.Println("Done")*/
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})
	http.ListenAndServe(":80", nil)
}
