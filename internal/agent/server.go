package agent

import (
	"github.com/philipslstwoyears/calculator-go/internal/dto"
	"github.com/philipslstwoyears/calculator-go/internal/storage"
	"github.com/philipslstwoyears/calculator-go/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT_AGENT")
	if config.Addr == "" {
		config.Addr = "8081"
	}
	return config
}

type Application struct {
	proto.UnimplementedCalcServiceServer
	config  *Config
	Storage storage.Storage
	input   chan dto.Expression
}

func New(s storage.Storage, input chan dto.Expression) *Application {
	return &Application{
		config:  ConfigFromEnv(),
		Storage: s,
		input:   input,
	}
}
func (a *Application) RunServer() error {
	lis, err := net.Listen("tcp", "0.0.0.0:"+a.config.Addr)
	if err != nil {
		return err
	}
	log.Printf("[AGENT SERVER] listening for gRPC, addr: %s", a.config.Addr)
	privateGRPC := grpc.NewServer()
	proto.RegisterCalcServiceServer(privateGRPC, a)

	return privateGRPC.Serve(lis)
}
