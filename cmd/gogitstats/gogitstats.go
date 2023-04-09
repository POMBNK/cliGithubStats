package main

import (
	"log"
	"os"

	"github.com/POMBNK/cliGitStats/internal/service"
	"github.com/POMBNK/cliGitStats/pkg/scanHelpers/finder"
	"github.com/POMBNK/cliGitStats/pkg/scanHelpers/writer"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load env vars %s", err.Error())
	}

	f := finder.New()
	w := writer.New()
	s := service.New(os.Getenv("BASE_PATH"), f, w)
	s.Run()

}
