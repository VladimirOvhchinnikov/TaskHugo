package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/chi"
)

func main() {

	var doMain chan bool

	//os.Setenv("HOST", "127.0.0.1:1313 hugo")
	router := chi.NewRouter()

	proxy := NewReverseProxy("hugo", "1313")
	router.Use(proxy.ReverseProxy)

	router.Get("/api", handlerRoute)

	go WorkerTest3()

	http.ListenAndServe(":8080", router)

	<-doMain
}

func handlerRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	var text string = fmt.Sprintf("<!DOCTYPE html><html><head><title>Webserver</title></head><body>Hello API</body></html>")

	w.Write([]byte(text))

}

type ReverseProxy struct {
	host string
	port string
}

func NewReverseProxy(host, port string) *ReverseProxy {
	return &ReverseProxy{
		host: host,
		port: port,
	}
}

func (rp *ReverseProxy) ReverseProxy(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api" {
			fmt.Println("proxy start")
			targetURL, _ := url.Parse("http://hugo:1313")
			proxy := httputil.NewSingleHostReverseProxy(targetURL)
			proxy.ServeHTTP(w, r)

			return

		}
		handler.ServeHTTP(w, r)
	})
}

const content = `---
menu:
    before:
        name: tasks
        weight: 5
title: Обновление данных в реальном времени
---

# Задача: Обновление данных в реальном времени

Напишите воркер, который будет обновлять данные в реальном времени, на текущей странице.
Текст данной задачи менять нельзя, только время и счетчик.

Файл данной страницы: /app/static/tasks/_index.md

Должен меняться счетчик и время:

Текущее время:%s

Счетчик: %d



## Критерии приемки:
- [ ] Воркер должен обновлять данные каждые 5 секунд
- [ ] Счетчик должен увеличиваться на 1 каждые 5 секунд
- [ ] Время должно обновляться каждые 5 секунд
`

func WorkerTest() {
	var timer *time.Ticker = time.NewTicker(5 * time.Second)
	const path string = "/app/static/tasks/_index.md"
	var count int = 0
	for {
		select {
		case <-timer.C:
			{
				err := os.WriteFile(path,
					[]byte(fmt.Sprintf(content, (time.Now().Format("2006-01-02 15:04:05")), count)),
					0644)
				if err != nil {
					log.Println(err)
				}
				count++
			}

		}
	}
}

func WorkerTest2(doMain chan bool) {

	var staticText []byte = []byte("Текущее время:")
	var staticCount []byte = []byte("Счетчик:")
	var timer *time.Ticker = time.NewTicker(5 * time.Second)
	var countText int = 0
	for {
		select {
		case <-timer.C:
			{
				var newTimeText []byte = []byte(fmt.Sprintf("%s %s", staticText, time.Now().Format("2006-01-02 15:04:05")))
				var newCountText []byte = []byte(fmt.Sprintf("%s %d", staticCount, countText))

				file, err := os.ReadFile("/app/static/tasks/_index.md")
				if err != nil {
					log.Println(err)
				}

				var fileSplit [][]byte = bytes.Split(file, []byte("\n"))

				for count, fileByte := range fileSplit {
					if bytes.Contains(fileByte, staticText) {
						fileSplit[count] = newTimeText
					}
					if bytes.Contains(fileByte, staticCount) {
						fileSplit[count] = newCountText
					}
				}

				file = nil
				for count, fileByte := range fileSplit {
					if count != 0 {
						file = append(file, '\n')
					}

					file = append(file, fileByte...)
				}
				fmt.Println(string(file))

				err = os.WriteFile("/app/static/tasks/_index.md", file, 0644)
				countText++
			}
		}
	}

}

func WorkerTest3() {
	t := time.NewTicker(5 * time.Second)
	var b int
	path := "/app/static/tasks/_index.md"
	//f,_ := os.ReadFile(path)
	//  pathLockal := "C:\\Users\\kolya\\hugoKAta\\hugoproxy\\hugo\\content\\tasks\\_index.md"
	for {
		select {
		case <-t.C:
			{
				err := os.WriteFile(path, []byte(fmt.Sprintf(content, time.Now().Format("2006-01-02 15:04:05"), b)), 0644)
				if err != nil {
					log.Println(err)
				}
				b++
			}
		}
	}
}
