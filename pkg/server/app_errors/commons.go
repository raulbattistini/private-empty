package app_errors

import (
	"errors"
)

var ErrInternalServeError = errors.New("Internal Server Error")
var ErrInternalServeErrorStr = errors.New("Internal Server Error").Error()
var ErrInvalidLimit = errors.New("O número selecionado para busca é inválido")
var ErrInvalidLimitStr = errors.New("O número selecionado para busca é inválido").Error()
var NoLimitValInReq = errors.New("Não foi fornecido um número válido para consulta.")
var NoLimitValInReqStr = errors.New("Não foi fornecido um número válido para consulta.").Error()
var InvalidMetricRequest = errors.New("A métrica solicitada é inválida ou não está inclusa.")
var InvalidMetricRequestStr = errors.New("A métrica solicitada é inválida ou não está inclusa.").Error()
var InvalidLimitNumRequest = errors.New("Valor inválido, o limite deve ser um número.")
var InvalidLimitNumRequestStr = errors.New("Valor inválido, o limite deve ser um número.").Error()

type ErrorResponse struct {
	Status  *int   `json:"status,omitempty"`
	Message string `json:"message"`
}
