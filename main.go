package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go handleSignals(cancel) // запускаем обрабочик сигналов, передаем функию отмены контекста

	if err := startServer(ctx); err != nil { // запускаем сервер
		log.Fatal(err)
	}
}

func handleSignals(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal)      // создаем канал для помещения в него сигнала
	signal.Notify(sigChan, os.Interrupt) // подписываемся на получение сигнала завершения, указав созданный канал
	for {
		sig := <-sigChan // если в канал поступает какой-нибудь сигнал, проверяем его
		switch sig {
		case os.Interrupt: // если это сигнал завершения,
			cancel() // выполняем cancel контекста, которая прекратит работу всех горутин контекста
			return
		}
	}
}

func startServer(ctx context.Context) error {
	lissenerAddr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", lissenerAddr)
	if err != nil {
		return err
	}

	defer l.Close()

	for {
		select {
		case <-ctx.Done(): // если пришел сигнал о завершенни, останавливаем работу
			log.Println("server stopped")
			return nil
		default: // если сигнала о завершении нет, устанавливаем таймаут 1с (
			if err := l.SetDeadline(time.Now().Add(time.Second)); err != nil {
				return err
			}

			_, err := l.Accept()
			if err != nil {
				if os.IsTimeout(err) { // если ошибка вызвана нашим таймаутом, перезапускаем цикл
					continue
				}
				return err
			}

			log.Println("new client connected")
		}
	}
}
