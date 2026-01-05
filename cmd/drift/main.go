package main

import (
	"drift/internal/core/engine"
	"drift/internal/core/model"
	"drift/internal/utils"
	"flag"
	"fmt"
	"os"
)

var (
	jsonOut bool
	jsonPath string
	base string
)

func init() {
	flag.StringVar(&base, "base", "origin/main", "Base branch to compare against")
	flag.BoolVar(&jsonOut, "json", false, "Output results as JSON")
	flag.StringVar(&jsonPath, "path", "drift.json", "Path to write JSON output")
}

func main() {
	
	flag.Parse()

	ctx := model.NewContext()
	eng := engine.NewEngine()

	ctx.Base = base
	ctx.JSON = jsonOut
	ctx.JSONPath = jsonPath
	
	findings, err := eng.Run(ctx)

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	if len(findings) == 0 {
		fmt.Println("No changes")
		return
	}
	if jsonOut {
		if err := utils.WriteJSON(jsonPath, findings); err != nil {
			fmt.Println("Failed to write JSON:", err)
			os.Exit(1)
		}
		return
	}
	for _, f := range findings {
		fmt.Println(f.Message)
	}
}