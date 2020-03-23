package app

import (
	"github.com/shohrukh56/jwt/pkg/jwt"
	"github.com/shohrukh56/mux/pkg/mux"
	"github.com/shohrukh56/rest/pkg/rest"
	"github.com/shohrukh56/sales/pkg/core/sales"
	"github.com/shohrukh56/sales/pkg/core/token"
	jwt2 "github.com/shohrukh56/sales/pkg/mux/middleware/jwt"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	router      *mux.ExactMux
	purchaseSvc *sales.Service
	secret      jwt.Secret
}

func NewServer(router *mux.ExactMux, purchaseSvc *sales.Service, secret jwt.Secret) *Server {
	return &Server{router: router, purchaseSvc: purchaseSvc, secret: secret}
}

func (s Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s Server) Start() {
	s.InitRoutes()
}

func (s Server) handlePurchaseList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		list, err := s.purchaseSvc.BuyingList(request.Context())
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		err = rest.WriteJSONBody(writer, &list)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
		}
	}
}
func (s Server) handleSales() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		context, ok := mux.FromContext(request.Context(), "id")
		if !ok {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(context)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		get := request.Header.Get("Content-Type")
		if get != "application/json" {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		pur := sales.Purchase{}
		err = rest.ReadJSONBody(request, &pur)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		payload := jwt2.FromContext(request.Context()).(*token.Payload)
		pur.User_id = payload.Id
		if id == 0 {
			err = s.purchaseSvc.AddNewPurchase(request.Context(), pur)
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Print(err)
				return
			}
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		if id > 0 {
			err = s.purchaseSvc.UpdateBought(request.Context(), int64(id), pur)
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}
func (s Server) handleSalesByUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		context, ok := mux.FromContext(request.Context(), "id")
		if !ok {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(context)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		payload := jwt2.FromContext(request.Context()).(*token.Payload)
		if id == 0 {
			prod, err := s.purchaseSvc.BuyByUserID(request.Context(), payload.Id)
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Print(err)
				return
			}
			err = rest.WriteJSONBody(writer, &prod)
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Print(err)
			}
			return
		}
		if id > 0 {
			for _, role := range payload.Roles {
				if role == "Admin" {
					prod, err := s.purchaseSvc.BuyByUserID(request.Context(), int64(id))
					if err != nil {
						http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
						log.Print(err)
						return
					}
					err = rest.WriteJSONBody(writer, &prod)
					if err != nil {
						http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
						log.Print(err)
					}
					return
				}
			}
		}
		http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	}
}
func (s Server) handleSalesByID() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		context, ok := mux.FromContext(request.Context(), "id")
		if !ok {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		purchase_ID, err := strconv.Atoi(context)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		purchaseByID, err := s.purchaseSvc.BuyByID(request.Context(), int64(purchase_ID))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		err = rest.WriteJSONBody(writer, &purchaseByID)
		if err != nil {
			log.Print(err)
		}
	}
}

func (s Server) handleSalesDelete() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		context, ok := mux.FromContext(request.Context(), "id")
		if !ok {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		purchase_ID, err := strconv.Atoi(context)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		err = s.purchaseSvc.RemoveByID(request.Context(), int64(purchase_ID))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		writer.WriteHeader(http.StatusNoContent)

	}
}
