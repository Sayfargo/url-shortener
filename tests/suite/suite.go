package suite

import (
	"context"
	"fmt"
	"testing"
	"time"

	urlv1 "github.com/Sayfargo/protos/gen/go/url"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	UrlShortenerClient urlv1.UrlShortenerClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", "localhost", "44044"), opts...)

	if err != nil {
		t.Fatalf("failed to create a new client server : %v", err)
	}

	t.Cleanup(func() {
		t.Helper()
		cancel()
		conn.Close()
	})

	return ctx, &Suite{
		T:                  t,
		UrlShortenerClient: urlv1.NewUrlShortenerClient(conn),
	}

}
