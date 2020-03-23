package app

import (
	"github.com/shohrukh56/sales/pkg/core/token"
	"github.com/shohrukh56/sales/pkg/mux/middleware/authenticated"
	"github.com/shohrukh56/sales/pkg/mux/middleware/authorized"
	"github.com/shohrukh56/sales/pkg/mux/middleware/jwt"
	"github.com/shohrukh56/sales/pkg/mux/middleware/logger"
	"reflect"
)

func (s Server) InitRoutes() {

	s.router.GET(
		"/api/purchases",
		s.handlePurchaseList(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		authorized.Authorized([]string{"Admin"}, jwt.FromContext),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("get list"),
	)
	s.router.DELETE(
		"/api/purchases/{id}",
		s.handlePurchaseDelete(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		authorized.Authorized([]string{"Admin"}, jwt.FromContext),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("removing"),
	)

	s.router.GET(
		"/api/purchases/{id}",
		s.handlePurchaseByID(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("get sales by id"),
	)
	s.router.GET(
		"/api/purchases/users/{id}",
		s.handlePurchaseByUser(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("get sales by user_id"),
	)

	s.router.POST(
		"/api/purchases/{id}",
		s.handlePurchase(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("post new sales"),
	)

}