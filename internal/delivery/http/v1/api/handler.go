package api

import "Prime_Number_Tester/internal/service"

type Handler struct {
	numbersService service.NumbersServiceProvider
}

func NewHandler(numbersService service.NumbersServiceProvider) *Handler {
	return &Handler{numbersService: numbersService}
}
