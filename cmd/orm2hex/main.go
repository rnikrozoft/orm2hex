package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rnikrozoft/orm2hex/internal/generator"
	"github.com/rnikrozoft/orm2hex/internal/parser"
)

func main() {
	outDir := flag.String("out", "./repository", "output directory for generated files")
	withCtx := flag.Bool("ctx", false, "Use context in repository methods")
	rawSQL := flag.Bool("raw", false, "Use raw SQL instead of ORM methods")
	orm := flag.String("orm", "gorm", "ORM to use (gorm or bun)")

	flag.Parse()

	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	structs, err := parser.ScanStructs(projectRoot)
	if err != nil {
		log.Fatal(err)
	}

	if len(structs) == 0 {
		log.Fatal("No structs found in project")
	}

	config := generator.GeneratorConfig{
		ORM:     *orm,
		RawSQL:  *rawSQL,
		WithCtx: *withCtx,
	}

	for _, s := range structs {
		err := generator.GenerateHexCRUD(s, config, *outDir)
		if err != nil {
			log.Printf("Failed to generate repository for %s: %v", s.Name, err)
		} else {
			fmt.Printf("Generated repository for %s using %s\n", s.Name, *orm)
		}
	}
}
