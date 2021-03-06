// Code generated by counterfeiter. DO NOT EDIT.
package libfakes

import (
	"sync"

	"github.com/rdrdog/bldr/pkg/lib"
)

type FakeLibProvider struct {
	GetDockerLibStub        func() lib.Docker
	getDockerLibMutex       sync.RWMutex
	getDockerLibArgsForCall []struct {
	}
	getDockerLibReturns struct {
		result1 lib.Docker
	}
	getDockerLibReturnsOnCall map[int]struct {
		result1 lib.Docker
	}
	GetGitLibStub        func() lib.Git
	getGitLibMutex       sync.RWMutex
	getGitLibArgsForCall []struct {
	}
	getGitLibReturns struct {
		result1 lib.Git
	}
	getGitLibReturnsOnCall map[int]struct {
		result1 lib.Git
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeLibProvider) GetDockerLib() lib.Docker {
	fake.getDockerLibMutex.Lock()
	ret, specificReturn := fake.getDockerLibReturnsOnCall[len(fake.getDockerLibArgsForCall)]
	fake.getDockerLibArgsForCall = append(fake.getDockerLibArgsForCall, struct {
	}{})
	stub := fake.GetDockerLibStub
	fakeReturns := fake.getDockerLibReturns
	fake.recordInvocation("GetDockerLib", []interface{}{})
	fake.getDockerLibMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeLibProvider) GetDockerLibCallCount() int {
	fake.getDockerLibMutex.RLock()
	defer fake.getDockerLibMutex.RUnlock()
	return len(fake.getDockerLibArgsForCall)
}

func (fake *FakeLibProvider) GetDockerLibCalls(stub func() lib.Docker) {
	fake.getDockerLibMutex.Lock()
	defer fake.getDockerLibMutex.Unlock()
	fake.GetDockerLibStub = stub
}

func (fake *FakeLibProvider) GetDockerLibReturns(result1 lib.Docker) {
	fake.getDockerLibMutex.Lock()
	defer fake.getDockerLibMutex.Unlock()
	fake.GetDockerLibStub = nil
	fake.getDockerLibReturns = struct {
		result1 lib.Docker
	}{result1}
}

func (fake *FakeLibProvider) GetDockerLibReturnsOnCall(i int, result1 lib.Docker) {
	fake.getDockerLibMutex.Lock()
	defer fake.getDockerLibMutex.Unlock()
	fake.GetDockerLibStub = nil
	if fake.getDockerLibReturnsOnCall == nil {
		fake.getDockerLibReturnsOnCall = make(map[int]struct {
			result1 lib.Docker
		})
	}
	fake.getDockerLibReturnsOnCall[i] = struct {
		result1 lib.Docker
	}{result1}
}

func (fake *FakeLibProvider) GetGitLib() lib.Git {
	fake.getGitLibMutex.Lock()
	ret, specificReturn := fake.getGitLibReturnsOnCall[len(fake.getGitLibArgsForCall)]
	fake.getGitLibArgsForCall = append(fake.getGitLibArgsForCall, struct {
	}{})
	stub := fake.GetGitLibStub
	fakeReturns := fake.getGitLibReturns
	fake.recordInvocation("GetGitLib", []interface{}{})
	fake.getGitLibMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeLibProvider) GetGitLibCallCount() int {
	fake.getGitLibMutex.RLock()
	defer fake.getGitLibMutex.RUnlock()
	return len(fake.getGitLibArgsForCall)
}

func (fake *FakeLibProvider) GetGitLibCalls(stub func() lib.Git) {
	fake.getGitLibMutex.Lock()
	defer fake.getGitLibMutex.Unlock()
	fake.GetGitLibStub = stub
}

func (fake *FakeLibProvider) GetGitLibReturns(result1 lib.Git) {
	fake.getGitLibMutex.Lock()
	defer fake.getGitLibMutex.Unlock()
	fake.GetGitLibStub = nil
	fake.getGitLibReturns = struct {
		result1 lib.Git
	}{result1}
}

func (fake *FakeLibProvider) GetGitLibReturnsOnCall(i int, result1 lib.Git) {
	fake.getGitLibMutex.Lock()
	defer fake.getGitLibMutex.Unlock()
	fake.GetGitLibStub = nil
	if fake.getGitLibReturnsOnCall == nil {
		fake.getGitLibReturnsOnCall = make(map[int]struct {
			result1 lib.Git
		})
	}
	fake.getGitLibReturnsOnCall[i] = struct {
		result1 lib.Git
	}{result1}
}

func (fake *FakeLibProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getDockerLibMutex.RLock()
	defer fake.getDockerLibMutex.RUnlock()
	fake.getGitLibMutex.RLock()
	defer fake.getGitLibMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeLibProvider) recordInvocation(key string, args []interface{}) {
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

var _ lib.LibProvider = new(FakeLibProvider)
