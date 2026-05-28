// Package main is the temporal lesson `l3_sagas_compensation` homework scaffold for Vibe Learn.
//
// Задача: BookingWorkflow (сага): BookFlight/BookHotel/BookCar + компенсации в ОБРАТНОМ порядке при сбое.
// Реализуй workflow и активности ниже — сигнатуры и тестовая поверхность
// фиксированы; CI (.github/workflows/ci.yml) гоняет `go vet` и `go test ./...`.
// Подробности и критерии приёмки — в README.md.
//
// SDK: go.temporal.io/sdk (worker + workflow + activity).
// Воркер подключается к Temporal по TEMPORAL_ADDRESS (дефолт localhost:7233 —
// совпадает с docker-compose.yml) и слушает task queue из TaskQueue().
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// ----- config -----

// envOr returns the env var for `key` if set, else `fallback`.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// TemporalAddress — адрес Temporal frontend. Дефолт совпадает с docker-compose.yml.
func TemporalAddress() string {
	return envOr("TEMPORAL_ADDRESS", "localhost:7233")
}

// TaskQueue — очередь задач, которую слушает воркер этого урока.
func TaskQueue() string {
	return envOr("TEMPORAL_TASK_QUEUE", "lesson-l3_sagas_compensation-tq")
}

// ----- Workflow: BookingWorkflow -----
//
// Оркеструет активности ниже. Тело — TODO: добавь ExecuteActivity-шаги,
// ActivityOptions (StartToCloseTimeout, RetryPolicy) и обработку ошибок
// согласно README.md. Должно оставаться ДЕТЕРМИНИРОВАННЫМ (никаких
// time.Now/rand/итераций по map — используй workflow.Now/SideEffect).
func BookingWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	logger.Info("BookingWorkflow started", "taskQueue", TaskQueue())

	// TODO #1: вызови активность BookFlight через workflow.ExecuteActivity.
	// var bookflightRes string
	// if err := workflow.ExecuteActivity(ctx, BookFlight).Get(ctx, &bookflightRes); err != nil {
	// 	return err
	// }
	// TODO #2: вызови активность BookHotel через workflow.ExecuteActivity.
	// var bookhotelRes string
	// if err := workflow.ExecuteActivity(ctx, BookHotel).Get(ctx, &bookhotelRes); err != nil {
	// 	return err
	// }
	// TODO #3: вызови активность BookCar через workflow.ExecuteActivity.
	// var bookcarRes string
	// if err := workflow.ExecuteActivity(ctx, BookCar).Get(ctx, &bookcarRes); err != nil {
	// 	return err
	// }
	// TODO #4: вызови активность CancelFlight через workflow.ExecuteActivity.
	// var cancelflightRes string
	// if err := workflow.ExecuteActivity(ctx, CancelFlight).Get(ctx, &cancelflightRes); err != nil {
	// 	return err
	// }
	// TODO #5: вызови активность CancelHotel через workflow.ExecuteActivity.
	// var cancelhotelRes string
	// if err := workflow.ExecuteActivity(ctx, CancelHotel).Get(ctx, &cancelhotelRes); err != nil {
	// 	return err
	// }

	return nil
}

// ----- Activity #1: BookFlight -----
//
// забронировать рейс; компенсация — CancelFlight
func BookFlight(ctx context.Context) (string, error) {
	// TODO: implement
	return "", fmt.Errorf("BookFlight: not implemented")
}

// ----- Activity #2: BookHotel -----
//
// забронировать отель; компенсация — CancelHotel
func BookHotel(ctx context.Context) (string, error) {
	// TODO: implement
	return "", fmt.Errorf("BookHotel: not implemented")
}

// ----- Activity #3: BookCar -----
//
// забронировать авто; флакает по флагу — провоцирует откат саги
func BookCar(ctx context.Context) (string, error) {
	// TODO: implement
	return "", fmt.Errorf("BookCar: not implemented")
}

// ----- Activity #4: CancelFlight -----
//
// идемпотентная компенсация рейса (с ретраями)
func CancelFlight(ctx context.Context) (string, error) {
	// TODO: implement
	return "", fmt.Errorf("CancelFlight: not implemented")
}

// ----- Activity #5: CancelHotel -----
//
// идемпотентная компенсация отеля (с ретраями)
func CancelHotel(ctx context.Context) (string, error) {
	// TODO: implement
	return "", fmt.Errorf("CancelHotel: not implemented")
}

// ----- main entry: register worker + run with graceful shutdown -----

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("Vibe Learn — temporal lesson %s scaffold up", "l3_sagas_compensation")
	log.Printf("temporal address: %s  task queue: %s", TemporalAddress(), TaskQueue())
	log.Printf("Реализуй workflow и активности, затем `go test ./...`. README.md содержит задачу.")

	c, err := client.Dial(client.Options{HostPort: TemporalAddress()})
	if err != nil {
		log.Fatalf("unable to create Temporal client (is `docker compose up -d` running?): %v", err)
	}
	defer c.Close()

	w := worker.New(c, TaskQueue(), worker.Options{})
	w.RegisterWorkflow(BookingWorkflow)
	w.RegisterActivity(BookFlight)
	w.RegisterActivity(BookHotel)
	w.RegisterActivity(BookCar)
	w.RegisterActivity(CancelFlight)
	w.RegisterActivity(CancelHotel)

	// Graceful shutdown so `go run .` is interactive — worker.InterruptCh()
	// stops the worker on Ctrl-C / SIGTERM.
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("worker stopped with error: %v", err)
	}
}
