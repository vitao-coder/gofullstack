package controllers

import "github.com/vitao-coder/gofullstack/api/middlewares"

func (s *Servidor) inicializarRotas() {

	//Home Routes
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	//Login Routes
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Usuarios routes
	s.Router.HandleFunc("/usuarios/criar", middlewares.SetMiddlewareJSON(s.CriarUsuario)).Methods("POST")
	s.Router.HandleFunc("/usuarios", middlewares.SetMiddlewareJSON(s.GetUsuarios)).Methods("GET")
	s.Router.HandleFunc("/usuarios/{id}", middlewares.SetMiddlewareJSON(s.GetUsuario)).Methods("GET")
	s.Router.HandleFunc("/usuarios/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthenticacao(s.AlterarUsuario))).Methods("PUT")
	s.Router.HandleFunc("/usuarios/{id}", middlewares.SetMiddlewareAuthenticacao(s.DeletarUsuario)).Methods("DELETE")

	//Posts routes
	s.Router.HandleFunc("/posts/criar", middlewares.SetMiddlewareJSON(s.CriarPost)).Methods("POST")
	s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(s.GetPosts)).Methods("GET")
	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(s.GetPost)).Methods("GET")
	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthenticacao(s.AlterarUsuario))).Methods("PUT")
	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareAuthenticacao(s.DeletarPost)).Methods("DELETE")
}
