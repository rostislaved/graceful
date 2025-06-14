# graceful

**graceful** — это простая библиотека для Go, позволяющая запускать несколько фоновых процессов и корректно завершать их работу при ошибках или остановке приложения.

## Возможности

- Одновременный запуск нескольких процессов.
- Корректное завершение (graceful shutdown) всех процессов при ошибке любого из них.
- Гибкое управление: можно отключать отдельные процессы.
- Поддержка собственного логгера (на базе slog).

## Пример использования

```go
package main

import (
    "context"
    "github.com/yourname/graceful"
)

func main() {
    svc1 := &MyService{}
    svc2 := &MyService{}
    
    gr := graceful.New(
        graceful.NewProcess(svc1).Enable(true),
        graceful.NewProcess(svc2).Enable(true)
    )
    
    err := gr.Start(context.Background())
    if err != nil {
    	panic(err)
    }
}

type MyService struct{}

func (s *MyService) Start(ctx context.Context) error {
    // Реализация логики
    
    return nil
}