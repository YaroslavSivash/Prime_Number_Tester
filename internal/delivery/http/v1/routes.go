package v1

import (
	"github.com/gorilla/mux"

	"Prime_Number_Tester/internal/delivery/http/v1/api"
	"Prime_Number_Tester/internal/service"
)

func RegisterHTTPEndpoints(router *mux.Router, numbersService service.NumbersServiceProvider) {
	apiHandler := api.NewHandler(numbersService)
	router.HandleFunc("/check-prime-number", apiHandler.CheckNumbersForPrimeNumber)
}
