package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rayzub/twitter-monitor/src/core"
)

// @todo: add logging!
func main() {
	sc := make(chan os.Signal, 1)
	isDev := flag.Bool("dev", false, "run monitor with dev environment values")
	flag.Parse()


	if err := LoadAndValidateEnv(*isDev); err != nil {
		panic(err)
	}
	ctx := context.Background()
	if err := core.New(ctx); err != nil {
		panic(err)
	}
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}	

func LoadAndValidateEnv(isDev bool) error {
	requiredEnvs := []string{"CSRF_TOKEN", "BEARER_TOKEN", "AUTH_TOKEN", "WEBHOOK", "BOT_TOKEN", "REQUEST_CHANNEL_ID"}

	fileExt := "prod"
	if isDev {
		fileExt = "dev"
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := godotenv.Load(fmt.Sprintf("%s/.env.%s", cwd, fileExt)); err != nil {
		return err
	}

	for _, env := range requiredEnvs {
		if _, ok := os.LookupEnv(env); !ok {
			return fmt.Errorf("missing required env: %s", env)
		}
	}

	return nil
}