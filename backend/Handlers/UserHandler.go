package handlers

import (
	dependencies "movieshop/backend/cmd/web/dependancies"
	"net/http"
)

type UserHandler struct {
	dep *dependencies.Dependencies
}

func CreateAccount(deps dependencies.Dependencies) *UserHandler {
	return &UserHandler{dep: deps}
}

func (u *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){

}
