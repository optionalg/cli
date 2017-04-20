// This file was generated by counterfeiter
package pluginactionfakes

import (
	"sync"

	"code.cloudfoundry.org/cli/actor/pluginaction"
	"code.cloudfoundry.org/cli/util/configv3"
)

type FakeConfig struct {
	PluginRepositoriesStub        func() []configv3.PluginRepository
	pluginRepositoriesMutex       sync.RWMutex
	pluginRepositoriesArgsForCall []struct{}
	pluginRepositoriesReturns     struct {
		result1 []configv3.PluginRepository
	}
	pluginRepositoriesReturnsOnCall map[int]struct {
		result1 []configv3.PluginRepository
	}
	PluginsStub        func() []configv3.Plugin
	pluginsMutex       sync.RWMutex
	pluginsArgsForCall []struct{}
	pluginsReturns     struct {
		result1 []configv3.Plugin
	}
	pluginsReturnsOnCall map[int]struct {
		result1 []configv3.Plugin
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeConfig) PluginRepositories() []configv3.PluginRepository {
	fake.pluginRepositoriesMutex.Lock()
	ret, specificReturn := fake.pluginRepositoriesReturnsOnCall[len(fake.pluginRepositoriesArgsForCall)]
	fake.pluginRepositoriesArgsForCall = append(fake.pluginRepositoriesArgsForCall, struct{}{})
	fake.recordInvocation("PluginRepositories", []interface{}{})
	fake.pluginRepositoriesMutex.Unlock()
	if fake.PluginRepositoriesStub != nil {
		return fake.PluginRepositoriesStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.pluginRepositoriesReturns.result1
}

func (fake *FakeConfig) PluginRepositoriesCallCount() int {
	fake.pluginRepositoriesMutex.RLock()
	defer fake.pluginRepositoriesMutex.RUnlock()
	return len(fake.pluginRepositoriesArgsForCall)
}

func (fake *FakeConfig) PluginRepositoriesReturns(result1 []configv3.PluginRepository) {
	fake.PluginRepositoriesStub = nil
	fake.pluginRepositoriesReturns = struct {
		result1 []configv3.PluginRepository
	}{result1}
}

func (fake *FakeConfig) PluginRepositoriesReturnsOnCall(i int, result1 []configv3.PluginRepository) {
	fake.PluginRepositoriesStub = nil
	if fake.pluginRepositoriesReturnsOnCall == nil {
		fake.pluginRepositoriesReturnsOnCall = make(map[int]struct {
			result1 []configv3.PluginRepository
		})
	}
	fake.pluginRepositoriesReturnsOnCall[i] = struct {
		result1 []configv3.PluginRepository
	}{result1}
}

func (fake *FakeConfig) Plugins() []configv3.Plugin {
	fake.pluginsMutex.Lock()
	ret, specificReturn := fake.pluginsReturnsOnCall[len(fake.pluginsArgsForCall)]
	fake.pluginsArgsForCall = append(fake.pluginsArgsForCall, struct{}{})
	fake.recordInvocation("Plugins", []interface{}{})
	fake.pluginsMutex.Unlock()
	if fake.PluginsStub != nil {
		return fake.PluginsStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.pluginsReturns.result1
}

func (fake *FakeConfig) PluginsCallCount() int {
	fake.pluginsMutex.RLock()
	defer fake.pluginsMutex.RUnlock()
	return len(fake.pluginsArgsForCall)
}

func (fake *FakeConfig) PluginsReturns(result1 []configv3.Plugin) {
	fake.PluginsStub = nil
	fake.pluginsReturns = struct {
		result1 []configv3.Plugin
	}{result1}
}

func (fake *FakeConfig) PluginsReturnsOnCall(i int, result1 []configv3.Plugin) {
	fake.PluginsStub = nil
	if fake.pluginsReturnsOnCall == nil {
		fake.pluginsReturnsOnCall = make(map[int]struct {
			result1 []configv3.Plugin
		})
	}
	fake.pluginsReturnsOnCall[i] = struct {
		result1 []configv3.Plugin
	}{result1}
}

func (fake *FakeConfig) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.pluginRepositoriesMutex.RLock()
	defer fake.pluginRepositoriesMutex.RUnlock()
	fake.pluginsMutex.RLock()
	defer fake.pluginsMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeConfig) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ pluginaction.Config = new(FakeConfig)