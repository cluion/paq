package provider

import (
	"fmt"
	"sort"
	"sync"
)

// defaultRegistry is the global provider registry.
var defaultRegistry = &Registry{
	providers: make(map[string]Provider),
}

// Registry stores registered Providers in a concurrent-safe map.
type Registry struct {
	mu        sync.RWMutex
	providers map[string]Provider
}

// Register adds a Provider to the default registry.
// Panics if a provider with the same Name is already registered.
func Register(p Provider) {
	defaultRegistry.register(p)
}

// Get retrieves a Provider by name from the default registry.
func Get(name string) (Provider, bool) {
	return defaultRegistry.get(name)
}

// All returns all registered Providers sorted alphabetically by Name.
func All() []Provider {
	return defaultRegistry.all()
}

// Available returns only Providers whose Detect returns true, sorted alphabetically.
func Available() []Provider {
	return defaultRegistry.available()
}

func (r *Registry) register(p Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := p.Name()
	if _, exists := r.providers[name]; exists {
		panic(fmt.Sprintf("provider already registered: %s", name))
	}
	r.providers[name] = p
}

func (r *Registry) get(name string) (Provider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.providers[name]
	return p, ok
}

func (r *Registry) all() []Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Provider, 0, len(r.providers))
	for _, p := range r.providers {
		result = append(result, p)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name() < result[j].Name()
	})
	return result
}

func (r *Registry) available() []Provider {
	all := r.all()
	result := make([]Provider, 0, len(all))
	for _, p := range all {
		if p.Detect() {
			result = append(result, p)
		}
	}
	return result
}
