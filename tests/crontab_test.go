package tests

import (
	"testing"
	"github.com/robfig/cron"
	"fmt"
)

func TestCron(t *testing.T) {
	tes1 := cron.New()
	tes1.AddFunc("0 * * * * *", func() {
		fmt.Println("helll world")
	})
	tes1.Start()
	select {}

}
