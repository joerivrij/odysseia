package app

import (
	"context"
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"

	pb "github.com/odysseia-greek/plato/proto"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	config := configs.PtolemaiosConfig{}
	handler := PtolemaiosHandler{Config: &config}
	pb.RegisterPtolemaiosServer(s, &handler)
	go func() {
		if err := s.Serve(lis); err != nil {
			glg.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestSayHello(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewPtolemaiosClient(conn)
	resp, err := client.GetSecret(ctx, &pb.VaultRequest{})
	if err != nil {
		t.Fatalf("SayHello failed: %v", err)
	}
	glg.Printf("Response: %+v", resp)
	// Test for output here.
}
