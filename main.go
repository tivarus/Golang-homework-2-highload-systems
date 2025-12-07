package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Основная функция
func main() {
	minioClient, err := minio.New("minio:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("root", "password", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("MinIO init failed: %v", err)
	}

	bucket := "users"

	err = minioClient.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{})
	if err != nil {
		exists, errB := minioClient.BucketExists(context.Background(), bucket)
		if errB != nil || !exists {
			log.Fatalf("Cannot create bucket %s: %v", bucket, err)
		}
	}

	log.Println("MinIO connected, bucket ready:", bucket)

	// Создание роутера
	r := mux.NewRouter()

	// Определение маршрутов
	r.HandleFunc("/api/users", GetUsers).Methods("GET")
	r.HandleFunc("/api/users/{id}", GetUser).Methods("GET")
	r.HandleFunc("/api/users", CreateUser).Methods("POST")
	r.HandleFunc("/api/users/{id}", UpdateUser).Methods("PUT")
	r.HandleFunc("/api/users/{id}", DeleteUser).Methods("DELETE")

	r.Handle("/metrics", promhttp.Handler())

	withMiddlewares := rateLimitMiddleware(
		metricsMiddleware(r),
	)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", withMiddlewares); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
