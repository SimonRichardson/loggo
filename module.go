// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"strings"
	"sync"
)

// Do not change rootName: modules.resolve() will misbehave if it isn't "".
const rootName = ""

type module struct {
	name   string
	level  Level
	parent *module
}

// Name returns the module's name.
func (module *module) Name() string {
	if module.name == rootName {
		return "<root>"
	}
	return module.name
}

func (module *module) getEffectiveLogLevel() Level {
	// Note: the root module is guaranteed to have a
	// specified logging level, so acts as a suitable sentinel
	// for this loop.
	for {
		if level := module.level.get(); level != UNSPECIFIED {
			return level
		}
		module = module.parent
	}
	panic("unreachable")
}

type modules struct {
	mu        sync.Mutex
	rootLevel Level
	all       map[string]*module
}

// Initially the modules map only contains the root module.
func newModules(rootLevel Level) *modules {
	root := &module{
		name:  rootName,
		level: rootLevel,
	}
	return &modules{
		rootLevel: rootLevel,
		all: map[string]*module{
			rootName: root,
		},
	}
}

// get returns a Logger for the given module name,
// creating it and its parents if necessary.
func (m *modules) get(name string) *module {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Lowercase the module name, and look for it in the modules map.
	name = strings.ToLower(name)
	return m.resolve(name)
}

// resolve assumes that the modules mutex is locked.
func (m *modules) resolve(name string) *module {
	impl, found := m.all[name]
	if found {
		return impl
	}
	parentName := rootName
	if i := strings.LastIndex(name, "."); i >= 0 {
		parentName = name[0:i]
	}
	parent := m.resolve(parentName)
	impl = &module{
		name:   name,
		level:  UNSPECIFIED,
		parent: parent,
	}
	m.all[name] = impl
	return impl
}

// ResetLogging iterates through the known modules and sets the levels of all
// to UNSPECIFIED, except for <root> which is set to WARNING.
func (m *modules) resetLevels() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, module := range m.all {
		if name == rootName {
			module.level.set(m.rootLevel)
		} else {
			module.level.set(UNSPECIFIED)
		}
	}
}

// config returns information about the modules and their logging
// levels. The information is returned in the format expected by
// ConfigureLoggers. Modules with UNSPECIFIED level will not
// be included.
func (m *modules) config() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return loggerInfo(m.all)
}
