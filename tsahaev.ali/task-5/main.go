package main

import (
	"context"
	"fmt"
	"time"

	"task-5/pkg/conveyer"
	"task-5/pkg/handlers"
)

func main() {
	fmt.Println("=== ПРОСТОЙ ТЕСТ КОНВЕЙЕРА ===")

	c := conveyer.New(10)

	c.RegisterDecorator(handlers.PrefixDecoratorFunc, "input", "output")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		fmt.Println("Запускаем конвейер...")
		errCh <- c.Run(ctx)
	}()

	time.Sleep(500 * time.Millisecond)
	fmt.Println("Конвейер запущен")

	fmt.Println("\nТест 1: Отправляем 'hello'")
	if err := c.Send("input", "hello"); err != nil {
		fmt.Printf("ОШИБКА Send: %v\n", err)
	} else {
		fmt.Println("Send успешен")
	}

	time.Sleep(300 * time.Millisecond)

	fmt.Println("\nТест 2: Получаем данные")
	if data, err := c.Recv("output"); err != nil {
		fmt.Printf("ОШИБКА Recv: %v\n", err)
	} else if data != "" {
		fmt.Printf("УСПЕХ! Получено: %s\n", data)
	}

	fmt.Println("\nТест 3: Отправляем 'no decorator' (должна быть ошибка)")
	if err := c.Send("input", "no decorator"); err != nil {
		fmt.Printf("ОШИБКА Send: %v\n", err)
	} else {
		fmt.Println("Send успешен (но обработчик должен вернуть ошибку)")
	}

	time.Sleep(500 * time.Millisecond)

	cancel()

	select {
	case err := <-errCh:
		if err != nil && err != context.Canceled {
			fmt.Printf("\nКонвейер завершился с ошибкой: %v\n", err)
		} else {
			fmt.Println("\nКонвейер завершился успешно")
		}
	case <-time.After(1 * time.Second):
		fmt.Println("\nТаймаут ожидания конвейера")
	}

	fmt.Println("=== ТЕСТ ЗАВЕРШЕН ===")
}
