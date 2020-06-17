package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "github.com/felipeweb/grpc-tad/proto/service"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	addr := ":6000"
	clientAddr := fmt.Sprintf("localhost%s", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to initialize TCP listen: %v", err)
	}
	defer lis.Close() // nolint

	go runGRPC(lis)
	runHTTP(clientAddr)
}

func runGRPC(lis net.Listener) {
	creds, err := credentials.NewServerTLSFromFile("server/server-cert.pem", "server/server-key.pem")
	if err != nil {
		log.Fatalf("Failed to setup tls: %v", err)
	}

	server := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(authInterceptor),
	)
	pb.RegisterUserServiceServer(server, newServer())

	log.Printf("gRPC Listening on %s\n", lis.Addr().String())
	server.Serve(lis) // nolint
}

func runHTTP(clientAddr string) {
	runtime.HTTPError = customHTTPError

	addr := ":6001"
	creds, err := credentials.NewClientTLSFromFile("server/server-cert.pem", "")
	if err != nil {
		log.Fatalf("gateway cert load error: %s", err)
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	mux := runtime.NewServeMux()
	if err := pb.RegisterUserServiceHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
	log.Printf("HTTP Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

type server struct {
	users map[string]pb.User
}

func newServer() server {
	return server{
		users: make(map[string]pb.User),
	}
}

func (s server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*empty.Empty, error) {
	log.Println("Creating user...")
	user := req.GetUser()

	if user.Username == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username cannot be empty")
	}

	if user.Role == "" {
		return nil, status.Errorf(codes.InvalidArgument, "role cannot be empty")
	}

	s.users[user.Username] = *user

	log.Println("User created!")
	return &empty.Empty{}, nil
}

func (s server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	log.Println("Getting user!")

	if req.Username == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username cannot be empty")
	}

	u, exists := s.users[req.Username]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	log.Println("User found!")
	return &u, nil
}

func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing context metadata")
	}
	if len(meta["authorization"]) != 1 {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	if meta["authorization"][0] != "nuveosummit" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return handler(ctx, req)
}

type errorBody struct {
	Err string `json:"error,omitempty"`
}

func customHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`
	var jErr error
	e, ok := status.FromError(err)
	if ok {
		w.Header().Set("Content-type", marshaler.ContentType())
		w.WriteHeader(runtime.HTTPStatusFromCode(e.Code()))
		jErr = json.NewEncoder(w).Encode(errorBody{
			Err: e.Message(),
		})
	}

	if jErr != nil {
		w.Write([]byte(fallback)) //nolint
	}
}
