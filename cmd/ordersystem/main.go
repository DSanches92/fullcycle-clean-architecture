package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/DSanches92/fullcycle-clean-architecture/internal/infra/database"
	"github.com/DSanches92/fullcycle-clean-architecture/internal/infra/graph"
	"github.com/DSanches92/fullcycle-clean-architecture/internal/infra/graph/generated"
	grpcservice "github.com/DSanches92/fullcycle-clean-architecture/internal/infra/grpc/service"
	pb "github.com/DSanches92/fullcycle-clean-architecture/internal/infra/grpc/protofile"
	"github.com/DSanches92/fullcycle-clean-architecture/internal/infra/web"
	"github.com/DSanches92/fullcycle-clean-architecture/internal/infra/web/webserver"
	"github.com/DSanches92/fullcycle-clean-architecture/internal/usecase"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/go-sql-driver/mysql"
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

var (
	dbDriver     = envOr("DB_DRIVER", "mysql")
	dbHost       = envOr("DB_HOST", "localhost")
	dbPort       = envOr("DB_PORT", "3306")
	dbUser       = envOr("DB_USER", "root")
	dbPassword   = envOr("DB_PASSWORD", "root")
	dbName       = envOr("DB_NAME", "orders")
	webPort      = envOr("WEB_PORT", "3000")
	grpcPort     = envOr("GRPC_PORT", "3002")
	graphQLPort  = envOr("GRAPHQL_PORT", "3001")
)

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := database.WaitForDB(db); err != nil {
		log.Fatalf("database not ready: %v", err)
	}
	if err := database.RunMigrations(buildDSN()); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	orderRepository := database.NewOrderRepository(db)

	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepository)

	// REST server
	webOrderHandler := web.NewWebOrderHandler(createOrderUseCase, listOrdersUseCase)
	webServer := webserver.NewWebServer(webOrderHandler, db, ":"+webPort)
	webServer.Setup()
	go func() {
		log.Printf("REST server listening on :%s", webPort)
		if err := webServer.Start(); err != nil {
			log.Fatalf("REST server error: %v", err)
		}
	}()

	// gRPC server
	grpcServer := grpc.NewServer()
	orderService := grpcservice.NewOrderService(createOrderUseCase, listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen on gRPC port: %v", err)
	}
	go func() {
		log.Printf("gRPC server listening on :%s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
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
		log.Printf("GraphQL server listening on :%s", graphQLPort)
		if err := http.ListenAndServe(":"+graphQLPort, graphqlRouter); err != nil {
			log.Fatalf("GraphQL server error: %v", err)
		}
	}()

	select {}
}

func openDB() (*sql.DB, error) {
	return sql.Open(dbDriver, buildDSN())
}

func buildDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
}