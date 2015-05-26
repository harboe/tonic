package tonic

import "path"

// Used internally to configure router, a RouterGroup is associated with a prefix
// and an array of handlers (middlewares)
type RouteGroup struct {
	Handlers     []HandleFunc
	absolutePath string
	parent       *RouteGroup
	api          *API
}

// Creates a new router group. You should add all the routes that have common middlwares or the same path prefix.
// For example, all the routes that use a common middlware for authorization could be grouped.
func (grp *RouteGroup) Group(component string, handlers ...HandleFunc) *RouteGroup {
	absPath := grp.calculateAbsolutePath(component)

	return &RouteGroup{
		Handlers:     grp.combineHandlers(handlers),
		parent:       grp,
		absolutePath: absPath,
		api:          grp.api,
	}
}

// Get is a shortcut for router.Handle("GET", path, handle)
func (grp *RouteGroup) Get(route string, handlers ...HandleFunc) {
	grp.Route(route, "GET", handlers...)
}

// Post is a shortcut for router.Handle("POST", path, handle)
func (grp *RouteGroup) Post(route string, handlers ...HandleFunc) {
	grp.Route(route, "POST", handlers...)
}

// Put is a shortcut for router.Handle("PUT", path, handle)
func (grp *RouteGroup) Put(route string, handlers ...HandleFunc) {
	grp.Route(route, "PUT", handlers...)
}

// Delete is a shortcut for router.Handle("DELETE", path, handle)
func (grp *RouteGroup) Delete(route string, handlers ...HandleFunc) {
	grp.Route(route, "DELETE", handlers...)
}

// Patch is a shortcut for router.Handle("PATCH", path, handle)
func (grp *RouteGroup) Patch(route string, handlers ...HandleFunc) {
	grp.Route(route, "PATCH", handlers...)
}

// Handle registers a new request handle and middlewares with the given path and method.
// The last handler should be the real handler, the other ones should be middlewares that can and should be shared among different routes.
// See the example code in github.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (grp *RouteGroup) Route(route, method string, handlers ...HandleFunc) {
	absRoute := grp.calculateAbsolutePath(route)
	handlers = grp.combineHandlers(handlers)

	grp.api.handleRoute(absRoute, method, handlers)
}

func (grp *RouteGroup) calculateAbsolutePath(relativePath string) string {
	if len(relativePath) == 0 {
		return grp.absolutePath
	}
	absolutePath := path.Join(grp.absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(absolutePath) != '/'
	if appendSlash {
		return absolutePath + "/"
	}
	return absolutePath
}

func (grp *RouteGroup) combineHandlers(handlers []HandleFunc) []HandleFunc {
	finalSize := len(grp.Handlers) + len(handlers)
	mergedHandlers := make([]HandleFunc, 0, finalSize)
	mergedHandlers = append(mergedHandlers, grp.Handlers...)
	return append(mergedHandlers, handlers...)
}
