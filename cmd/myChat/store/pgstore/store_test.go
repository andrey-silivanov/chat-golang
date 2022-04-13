package pgstore

import (
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/store/teststore"
	_ "github.com/lib/pq"
	"testing"
)

func TestMain(m *testing.M) {
	teststore.CreateTestDB(m)
}
