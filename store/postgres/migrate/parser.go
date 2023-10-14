package migrate

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"strings"
)

const (
	prefix  = "-- name: "
	comment = "--"
	newline = "\n"
)

// Statement represents a statement in the sql file.
type Statement struct {
	Name   string
	Value  string
	Driver string
}

func getStatements(filesFS embed.FS) ([]*Statement, error) {
	files, err := filesFS.ReadDir(".")
	if err != nil {
		panic(err)
	}

	var statements []*Statement

	for i := range files {
		file, err := filesFS.Open(files[i].Name())
		if err != nil {
			return nil, err
		}

		statements = append(statements, parse(file)...)
	}

	return statements, nil
}

func parse(r io.Reader) []*Statement {
	var (
		stmts []*Statement
		stmt  *Statement
	)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, prefix) {
			stmt = new(Statement)
			stmt.Name, stmt.Driver = parsePrefix(line, prefix)
			stmts = append(stmts, stmt)
		}

		if strings.HasPrefix(line, comment) {
			continue
		}

		if stmt != nil {
			stmt.Value += line + newline
		}
	}

	for _, stmt := range stmts {
		stmt.Value = strings.TrimSpace(stmt.Value)
	}

	return stmts
}

func parsePrefix(line, prefix string) (string, string) {
	line = strings.TrimPrefix(line, prefix)
	line = strings.TrimSpace(line)

	var name, driver string

	fmt.Sscanln(line, &name, &driver)

	return name, driver
}
