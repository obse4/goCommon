package test

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	comm "github.com/obse4/goCommon"
	"github.com/obse4/goCommon/kafka"
)

func TestGoCommon(t *testing.T) {
	t.Log("goCommon test")
	comm.Init("/Users/wanghan/github/goCommon/config.yml")

	comm.KafkaProducer["producer1"].SendMessage(kafka.ProducerMessage{
		Topic: "test",
		Value: "obse4^-^",
	})

	// 在另一个goroutine中，模拟发送Ctrl+C操作，即发送SIGINT信号
	go func() {
		time.AfterFunc(5*time.Second, func() {
			fmt.Println("Sending SIGINT signal...")

			process, err := os.FindProcess(os.Getpid())

			if err != nil {
				fmt.Printf("Failed to find process: %v", err)
			}

			err = process.Signal(syscall.SIGINT) // 向当前进程发送SIGINT信号
			if err != nil {
				fmt.Printf("Failed to send SIGINT signal: %v", err)
			}
		})
	}()

	comm.Run()
}
