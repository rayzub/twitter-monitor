package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rayzub/twitter-monitor/src/core"
)

// @todo: add logging!
func main() {
	isDev := flag.Bool("dev", false, "run monitor with dev environment values")
	flag.Parse()


	if err := LoadAndValidateEnv(*isDev); err != nil {
		return
	}


	ctx := context.Background()
	core := core.New(ctx)
	if err := core.BotClient.OpenGateway(ctx); err != nil {
		panic(err)
	}
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