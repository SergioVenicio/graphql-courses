package sql

import (
	"os"
)

func ReadQuery(queryName string) (string, error) {
	queryString, err := os.ReadFile("./sql/queries/" + queryName + ".sql")
	if err != nil {
		return "", err
	}

	return string(queryString), nil
}
