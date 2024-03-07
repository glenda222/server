package router

import (
	"net/http"
	"./healthcheck"
	"./middleware"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Router() http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		//AllowedOrigins: []string{"https://dtb35dqgksiev.cloudfront.net"},
		AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		// AllowCredentials: true,
	})
	router := mux.NewRouter()

	router.HandleFunc("/api/heartbeat", healthcheck.Heartbeat("/api/heartbeat")).Methods("GET", "HEAD")

	router.HandleFunc("/api/parameters", middleware.GetAllParameters).Methods("GET")

	router.HandleFunc("/api/businessunit", middleware.GetAllBusinessUnit).Methods("GET")
	router.HandleFunc("/api/businessunit/{id}", middleware.GetBusinessUnit).Methods("GET")
	router.HandleFunc("/api/businessunit", middleware.CreateBusinessUnit).Methods("POST")
	router.HandleFunc("/api/businessunit/{id}", middleware.UpdateBusinessUnit).Methods("PUT")
	router.HandleFunc("/api/businessunit/{id}", middleware.DeleteBusinessUnit).Methods("DELETE")

	router.HandleFunc("/api/applicationgroup", middleware.GetAllApplicationGroup).Methods("GET")
	router.HandleFunc("/api/applicationgroup/{id}", middleware.GetApplicationGroup).Methods("GET")
	router.HandleFunc("/api/applicationgroup", middleware.CreateApplicationGroup).Methods("POST")
	router.HandleFunc("/api/applicationgroup/{id}", middleware.UpdateApplicationGroup).Methods("PUT")
	router.HandleFunc("/api/applicationgroup/{id}", middleware.DeleteApplicationGroup).Methods("DELETE")

	router.HandleFunc("/api/application", middleware.GetAllApplication).Methods("GET")
	router.HandleFunc("/api/application/{id}", middleware.GetApplication).Methods("GET")
	router.HandleFunc("/api/application", middleware.CreateApplication).Methods("POST")
	router.HandleFunc("/api/application/{id}", middleware.UpdateApplication).Methods("PUT")
	router.HandleFunc("/api/application/import", middleware.UpdateManyApplication).Methods("POST")
	router.HandleFunc("/api/application/{id}", middleware.DeleteApplication).Methods("DELETE")
	router.HandleFunc("/api/applicationByOwner/{id}", middleware.GetAllApplicationByOwner).Methods("GET")
	router.HandleFunc("/api/log", middleware.GetAllLog).Methods("GET")
	router.HandleFunc("/api/log/{id}", middleware.GetLog).Methods("GET")
	router.HandleFunc("/api/log/{id}", middleware.DeleteLog).Methods("DELETE")

	//user routes
	router.HandleFunc("/api/register", middleware.RegisterUser).Methods("POST")

	router.HandleFunc("/api/login", middleware.LoginUser).Methods("POST")
	router.HandleFunc("/api/listAllUsers", middleware.ListAllUsers).Methods("GET")
	router.HandleFunc("/api/validateToken", middleware.ValidateToken).Methods("POST")
	router.HandleFunc("/api/requestPasswordReset", middleware.RequestPasswordReset).Methods("POST")
	router.HandleFunc("/api/resetPassword", middleware.ResetPassword).Methods("POST")
	router.HandleFunc("/api/resetPasswordForUser", middleware.ResetPasswordForUser).Methods("POST")

	// router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// })

	router.HandleFunc("/api/updateUserRole/{id}", middleware.UpdateUserRole).Methods("PUT")
	router.HandleFunc("/api/deleteUser/{id}", middleware.DeleteUser).Methods("DELETE")

	return c.Handler(handlers.CompressHandler(router))
}
