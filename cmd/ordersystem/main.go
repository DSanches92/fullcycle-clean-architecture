package main

import (
	"log"
	"net"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/DSanches92/fullcycle-clean-architecture/configs"
	"github.com/DSanches92/fullcycle-clean-architecture/internal/infra/database"
	"github.com/DSanches92/fullcycle-clean-architecture/internal/infra/graph"
	"github.com/DSanches92/fullcycle-clean-architecture/internal/infra/graph/generated"
	pb "github.com/DSanches92/fullcycle-clean-architecture/internal/infra/grpc/protofile"
	grpcservice "github.com/DSanches92/fullcycle-clean-architecture/internal/infra/grpc/service"
	"github.com/DSanches92/fullcycle-clean-architecture/internal/infra/web"
	"github.com/DSanches92/fullcycle-clean-architecture/internal/infra/web/webserver"
	"github.com/DSanches92/fullcycle-clean-architecture/internal/usecase"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	envs, err := configs.LoadEnvironment()
	if err != nil {
		log.Fatalf("Falha ao carregar as variáveis de ambiente: %v", err)
	}

	databaseConn, err := database.DatabaseConn(envs.Database)
	if err != nil {
		log.Fatalf("Falha conexão banco de dados: %v", err)
	}
	defer databaseConn.Close()

	orderRepository := database.NewOrderRepository(databaseConn)

	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepository)

	// REST server
	webOrderHandler := web.NewWebOrderHandler(createOrderUseCase, listOrdersUseCase)
	webServer := webserver.NewWebServer(webOrderHandler, databaseConn, ":"+envs.RestPort)
	webServer.Setup()
	go func() {
		log.Printf("Servidor REST ouvindo :%s", envs.RestPort)
		if err := webServer.Start(); err != nil {
			log.Fatalf("Erro servidor REST: %v", err)
		}
	}()

	// gRPC server
	grpcServer := grpc.NewServer()
	orderService := grpcservice.NewOrderService(createOrderUseCase, listOrdersUseCase)

	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+envs.GrpcPort)
	if err != nil {
		log.Fatalf("Falha ao ouvir a porta no servidor gRPC: %v", err)
	}
	go func() {
		log.Printf("Servidor gRPC ouvindo :%s", envs.GrpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Erro servidor gRPC: %v", err)
		}
	}()

	// GraphQL server
	graphResolver := &graph.Resolver{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graphResolver}))

	graphqlRouter := chi.NewRouter()
	graphqlRouter.Handle("/", playground.Handler("GraphQL playground", "/query"))
	graphqlRouter.Handle("/query", srv)
	go func() {
		log.Printf("Servidor GraphQL ouvindo :%s", envs.GraphQLPort)
		if err := http.ListenAndServe(":"+envs.GraphQLPort, graphqlRouter); err != nil {
			log.Fatalf("Erro servidor GraphQL: %v", err)
		}
	}()

	select {}
}
