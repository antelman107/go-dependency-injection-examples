package dic

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sarulabs/di"
	"github.com/sarulabs/dingo"

	log "log"

	service "github.com/antelman107/dependency_injection_examples/sarulabs/dingo/service"
)

// C retrieves a Container from an interface.
// The function panics if the Container can not be retrieved.
//
// The interface can be :
// - a *Container
// - an *http.Request containing a *Container in its context.Context
//   for the dingo.ContainerKey("dingo") key.
//
// The function can be changed to match the needs of your application.
var C = func(i interface{}) *Container {
	if c, ok := i.(*Container); ok {
		return c
	}
	r, ok := i.(*http.Request)
	if !ok {
		panic("could not get the container with dic.C()")
	}
	c, ok := r.Context().Value(dingo.ContainerKey("dingo")).(*Container)
	if !ok {
		panic("could not get the container from the given *http.Request in dic.C()")
	}
	return c
}

type builder struct {
	builder *di.Builder
}

// NewBuilder creates a builder that can create a Container.
// You should you NewContainer to create the container directly.
// Using NewBuilder allows you to redefine some di services though.
// This could be used for testing.
// But this behaviour is not safe, so be sure to know what you are doing.
func NewBuilder(scopes ...string) (*builder, error) {
	if len(scopes) == 0 {
		scopes = []string{di.App, di.Request, di.SubRequest}
	}
	b, err := di.NewBuilder(scopes...)
	if err != nil {
		return nil, fmt.Errorf("could not create di.Builder: %v", err)
	}
	provider := &service.Provider{}
	if err := provider.Load(); err != nil {
		return nil, fmt.Errorf("could not load definitions with the Provider (Provider from github.com/antelman107/dependency_injection_examples/dingo/service): %v", err)
	}
	for _, d := range getDiDefs(provider) {
		if err := b.Add(d); err != nil {
			return nil, fmt.Errorf("could not add di.Def in di.Builder: %v", err)
		}
	}
	return &builder{builder: b}, nil
}

// Add adds one or more definitions in the Builder.
// It returns an error if a definition can not be added.
func (b *builder) Add(defs ...di.Def) error {
	return b.builder.Add(defs...)
}

// Set is a shortcut to add a definition for an already built object.
func (b *builder) Set(name string, obj interface{}) error {
	return b.builder.Set(name, obj)
}

// Build creates a Container in the most generic scope.
func (b *builder) Build() *Container {
	return &Container{ctn: b.builder.Build()}
}

// NewContainer creates a new Container.
// If no scope is provided, di.App, di.Request and di.SubRequest are used.
// The returned Container has the most generic scope (di.App).
// The SubContainer() method should be called to get a Container in a more specific scope.
func NewContainer(scopes ...string) (*Container, error) {
	b, err := NewBuilder(scopes...)
	if err != nil {
		return nil, err
	}
	return b.Build(), nil
}

// Container represents a generated dependency injection container.
// It is a wrapper around a di.Container.
//
// A Container has a scope and may have a parent in a more generic scope
// and children in a more specific scope.
// Objects can be retrieved from the Container.
// If the requested object does not already exist in the Container,
// it is built thanks to the object definition.
// The following attempts to get this object will return the same object.
type Container struct {
	ctn di.Container
}

// Scope returns the Container scope.
func (c *Container) Scope() string {
	return c.ctn.Scope()
}

// Scopes returns the list of available scopes.
func (c *Container) Scopes() []string {
	return c.ctn.Scopes()
}

// ParentScopes returns the list of scopes wider than the Container scope.
func (c *Container) ParentScopes() []string {
	return c.ctn.ParentScopes()
}

// SubScopes returns the list of scopes that are more specific than the Container scope.
func (c *Container) SubScopes() []string {
	return c.ctn.SubScopes()
}

// Parent returns the parent Container.
func (c *Container) Parent() *Container {
	if p := c.ctn.Parent(); p != nil {
		return &Container{ctn: p}
	}
	return nil
}

// SubContainer creates a new Container in the next sub-scope
// that will have this Container as parent.
func (c *Container) SubContainer() (*Container, error) {
	sub, err := c.ctn.SubContainer()
	if err != nil {
		return nil, err
	}
	return &Container{ctn: sub}, nil
}

// SafeGet retrieves an object from the Container.
// The object has to belong to this scope or a more generic one.
// If the object does not already exist, it is created and saved in the Container.
// If the object can not be created, it returns an error.
func (c *Container) SafeGet(name string) (interface{}, error) {
	return c.ctn.SafeGet(name)
}

// Get is similar to SafeGet but it does not return the error.
// Instead it panics.
func (c *Container) Get(name string) interface{} {
	return c.ctn.Get(name)
}

// Fill is similar to SafeGet but it does not return the object.
// Instead it fills the provided object with the value returned by SafeGet.
// The provided object must be a pointer to the value returned by SafeGet.
func (c *Container) Fill(name string, dst interface{}) error {
	return c.ctn.Fill(name, dst)
}

// UnscopedSafeGet retrieves an object from the Container, like SafeGet.
// The difference is that the object can be retrieved
// even if it belongs to a more specific scope.
// To do so, UnscopedSafeGet creates a sub-container.
// When the created object is no longer needed,
// it is important to use the Clean method to delete this sub-container.
func (c *Container) UnscopedSafeGet(name string) (interface{}, error) {
	return c.ctn.UnscopedSafeGet(name)
}

// UnscopedGet is similar to UnscopedSafeGet but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGet(name string) interface{} {
	return c.ctn.UnscopedGet(name)
}

// UnscopedFill is similar to UnscopedSafeGet but copies the object in dst instead of returning it.
func (c *Container) UnscopedFill(name string, dst interface{}) error {
	return c.ctn.UnscopedFill(name, dst)
}

// Clean deletes the sub-container created by UnscopedSafeGet, UnscopedGet or UnscopedFill.
func (c *Container) Clean() error {
	return c.ctn.Clean()
}

// DeleteWithSubContainers takes all the objects saved in this Container
// and calls the Close function of their Definition on them.
// It will also call DeleteWithSubContainers on each child and remove its reference in the parent Container.
// After deletion, the Container can no longer be used.
// The sub-containers are deleted even if they are still used in other goroutines.
// It can cause errors. You may want to use the Delete method instead.
func (c *Container) DeleteWithSubContainers() error {
	return c.ctn.DeleteWithSubContainers()
}

// Delete works like DeleteWithSubContainers if the Container does not have any child.
// But if the Container has sub-containers, it will not be deleted right away.
// The deletion only occurs when all the sub-containers have been deleted manually.
// So you have to call Delete or DeleteWithSubContainers on all the sub-containers.
func (c *Container) Delete() error {
	return c.ctn.Delete()
}

// IsClosed returns true if the Container has been deleted.
func (c *Container) IsClosed() bool {
	return c.ctn.IsClosed()
}

// SafeGetConfig works like SafeGet but only for Config.
// It does not return an interface but a *service.Config.
func (c *Container) SafeGetConfig() (*service.Config, error) {
	i, err := c.ctn.SafeGet("config")
	if err != nil {
		var eo *service.Config
		return eo, err
	}
	o, ok := i.(*service.Config)
	if !ok {
		return o, errors.New("could get 'config' because the object could not be cast to *service.Config")
	}
	return o, nil
}

// GetConfig is similar to SafeGetConfig but it does not return the error.
// Instead it panics.
func (c *Container) GetConfig() *service.Config {
	o, err := c.SafeGetConfig()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetConfig works like UnscopedSafeGet but only for Config.
// It does not return an interface but a *service.Config.
func (c *Container) UnscopedSafeGetConfig() (*service.Config, error) {
	i, err := c.ctn.UnscopedSafeGet("config")
	if err != nil {
		var eo *service.Config
		return eo, err
	}
	o, ok := i.(*service.Config)
	if !ok {
		return o, errors.New("could get 'config' because the object could not be cast to *service.Config")
	}
	return o, nil
}

// UnscopedGetConfig is similar to UnscopedSafeGetConfig but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetConfig() *service.Config {
	o, err := c.UnscopedSafeGetConfig()
	if err != nil {
		panic(err)
	}
	return o
}

// Config is similar to GetConfig.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetConfig method.
// If the container can not be retrieved, it panics.
func Config(i interface{}) *service.Config {
	return C(i).GetConfig()
}

// SafeGetLogger works like SafeGet but only for Logger.
// It does not return an interface but a *log.Logger.
func (c *Container) SafeGetLogger() (*log.Logger, error) {
	i, err := c.ctn.SafeGet("logger")
	if err != nil {
		var eo *log.Logger
		return eo, err
	}
	o, ok := i.(*log.Logger)
	if !ok {
		return o, errors.New("could get 'logger' because the object could not be cast to *log.Logger")
	}
	return o, nil
}

// GetLogger is similar to SafeGetLogger but it does not return the error.
// Instead it panics.
func (c *Container) GetLogger() *log.Logger {
	o, err := c.SafeGetLogger()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetLogger works like UnscopedSafeGet but only for Logger.
// It does not return an interface but a *log.Logger.
func (c *Container) UnscopedSafeGetLogger() (*log.Logger, error) {
	i, err := c.ctn.UnscopedSafeGet("logger")
	if err != nil {
		var eo *log.Logger
		return eo, err
	}
	o, ok := i.(*log.Logger)
	if !ok {
		return o, errors.New("could get 'logger' because the object could not be cast to *log.Logger")
	}
	return o, nil
}

// UnscopedGetLogger is similar to UnscopedSafeGetLogger but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetLogger() *log.Logger {
	o, err := c.UnscopedSafeGetLogger()
	if err != nil {
		panic(err)
	}
	return o
}

// Logger is similar to GetLogger.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetLogger method.
// If the container can not be retrieved, it panics.
func Logger(i interface{}) *log.Logger {
	return C(i).GetLogger()
}
