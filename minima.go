package minima

import (
	"context"
	"log"
	"net/http"
	"time"
)

/**
@info The framework structure
@property {*http.Server} [server] The net/http stock server
@property {bool} [started] Whether the server has started or not
@property {*time.Duration} [Timeout] The router's breathing time
@property {*Router} [router] The core router instance running with the server
@property {map[string]interface{}} [properties] The properties for the server instance
@property {*Config} [Config] The core config file for middlewares and router instances
@property {*time.Duration} [drain] The router's drain time
*/
type minima struct {
	server     *http.Server
	started    bool
	Timeout    time.Duration
	router     *Router
	properties map[string]interface{}
	Config     *Config
	Middleware *Plugins
	drain      time.Duration
}

/**
@info Make a new default minima instance
@example `
func main() {
	app := minima.New()

	app.Get("/", func(res *minima.Response, req *minima.Request) {
		res.Status(200).Send("Hello World")
	})

	app.Listen(":3000")
}
`
@returns {minima}
*/
func New() *minima {
	return &minima{
		router:     NewRouter(),
		Config:     NewConfig(),
		Middleware: use(),
		drain:      0,
	}
}

/**
@info Starts the actual http server
@param {string} [addr] The port for the server instance to run on
@returns {error}
*/
func (m *minima) Listen(addr string) error {
	if m.started {
		log.Panicf("Minimia's instance is already running at %s.", m.server.Addr)
	}
	m.server = &http.Server{Addr: addr, Handler: m}
	m.started = true

	return m.server.ListenAndServe()

}

/**
@info Injects the actual minima server logic to http
@param {http.ResponseWriter} [w] The net/http response instance
@param {http.Request} [r] The net/http request instance
@returns {}
*/
func (m *minima) ServeHTTP(w http.ResponseWriter, q *http.Request) {
	f, params, match := m.router.routes[q.Method].Get(q.URL.Path)

	if match {
		if err := q.ParseForm(); err != nil {
			log.Printf("Error parsing form: %s", err)
			return
		}

		res := response(w, q, &m.properties)
		req := request(q)
		req.Params = params

		m.Middleware.ServePlugin(res, req)
		f(res, req)
	} else {
		res := response(w, q, &m.properties)
		req := request(q)
		if m.router.notfound != nil {
			m.router.notfound(res, req)
		} else {
			w.Write([]byte("No matching route found"))
		}

	}
}

/**
@info Adds route with Get method
@param {string} [path] The route path
@param {...Handler} [handler] The handler for the given route
@returns {*minima}
*/
func (m *minima) Get(path string, handler Handler) *minima {
	m.router.Get(path, handler)
	return m
}

/**
@info Adds route with Put method
@param {string} [path] The route path
@param {...Handler} [handler] The handler for the given route
@returns {*minima}
*/
func (m *minima) Put(path string, handler Handler) *minima {
	m.router.Put(path, handler)
	return m
}

/**
@info Adds route with Options method
@param {string} [path] The route path
@param {...Handler} [handler] The handler for the given route
@returns {*minima}
*/
func (m *minima) Options(path string, handler Handler) *minima {
	m.router.Options(path, handler)
	return m
}

/**
@info Adds route with Head method
@param {string} [path] The route path
@param {...Handler} [handler] The handler for the given route
@returns {*minima}
*/
func (m *minima) Head(path string, handler Handler) *minima {
	m.router.Head(path, handler)
	return m
}

/**
@info Adds route with Delete method
@param {string} [path] The route path
@param {...Handler} [handler] The handler for the given route
@returns {*minima}
*/
func (m *minima) Delete(path string, handler Handler) *minima {
	m.router.Delete(path, handler)
	return m
}

/**
@info Adds route with Patch method
@param {string} [path] The route path
@param {...Handler} [handler] The handler for the given route
@returns {*minima}
*/
func (m *minima) Patch(path string, handler Handler) *minima {
	m.router.Patch(path, handler)
	return m
}

/**
@info Adds route with Post method
@param {string} [path] The route path
@param {...Handler} [handler] The handler for the given route
@returns {*minima}
*/
func (m *minima) Post(path string, handler Handler) *minima {
	m.router.Post(path, handler)
	return m
}

/**
@info Injects the given handler to middleware stack
@param {Handler} [handler] Minima handler instance
@returns {*minima}
*/
func (m *minima) Use(handler Handler) *minima {
	m.Middleware.AddPlugin(handler)
	return m
}

/**
@info Injects the NotFound handler to the minima instance
@param {Handler} [handler] Minima handler instance
@returns {*minima}
*/
func (m *minima) NotFound(handler Handler) *minima {
	m.router.NotFound(handler)
	return m
}

/**
@info Injects the routes from the router to core stack
@param {*Router} [router] Minima router instance
@returns {*minima}
*/
func (m *minima) UseRouter(router *Router) *minima {
	m.router.UseRouter(router)
	return m

}

/**
@info Mounts router to a specific path
@param {string} [path] The route path
@param {*Router} [router] Minima router instance
@returns {*minima}
*/
func (m *minima) Mount(path string, router *Router) *minima {
	m.router.Mount(path, router)
	return m

}

/**
@info Injects middlewares and routers directly to core instance
@param {*Config} [config] The config instance
@returns {*minima}
*/
func (m *minima) UseConfig(config *Config) *minima {
	for _, v := range config.Middleware {
		m.Middleware.plugin = append(m.Middleware.plugin, &Middleware{handler: v})
	}
	for _, rt := range config.Router {
		m.router.UseRouter(rt)
	}
	return m
}

/**
@info The drain timeout for the core instance
@param {time.Duration} [time] The time period for drain
@returns {*minima}
*/
func (m *minima) ShutdownTimeout(t time.Duration) *minima {
	m.drain = t
	return m
}

/**
@info Shutdowns the core instance
@param {context.Context} [ctx] The context for shutdown
@returns {error}
*/
func (m *minima) Shutdown(ctx context.Context) error {
	log.Println("Stopping the server")
	return m.server.Shutdown(ctx)
}

/**
@info Declares prop for core properties
@param {string} [key] Key for the prop
@param {interface{}} [value] Value of the prop
@returns {*minima}
*/
func (m *minima) SetProp(key string, value interface{}) *minima {
	m.properties[key] = value
	return m
}

/**
@info Gets prop from core properties
@param {string} [key] Key for the prop
@returns {interface{}}
*/
func (m *minima) GetProp(key string) interface{} {
	return m.properties[key]
}
