package backend

import (
	"testing"

	"github.com/joho/godotenv"
)

type BlogDBFetcherMock struct{}

func TestOpenConnection(t *testing.T) {
	var fetcher = BlogFetcher{}
	godotenv.Load(".testingenv")
	fetcher.OpenConnection()
}
