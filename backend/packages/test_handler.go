package backend

import (
	backend "backend/packages"
	"testing"

	"github.com/joho/godotenv"
)

type BlogDBFetcherMock struct{}

func TestOpenConnection(t *testing.T) {
	godotenv.Load(".testingenv")
	backend.OpenConnection()
}
