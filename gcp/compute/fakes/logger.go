package fakes

import "sync"

type Logger struct {
	DebugfCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Message string
			A       []interface {
			}
		}
		Stub func(string, ...interface {
		})
	}
	DebuglnCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Message string
		}
		Stub func(string)
	}
	PrintfCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			M string
			A []interface {
			}
		}
		Stub func(string, ...interface {
		})
	}
	PromptWithDetailsCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			ResourceType string
			ResourceName string
		}
		Returns struct {
			Proceed bool
		}
		Stub func(string, string) bool
	}
}

func (f *Logger) Debugf(param1 string, param2 ...interface {
}) {
	f.DebugfCall.Lock()
	defer f.DebugfCall.Unlock()
	f.DebugfCall.CallCount++
	f.DebugfCall.Receives.Message = param1
	f.DebugfCall.Receives.A = param2
	if f.DebugfCall.Stub != nil {
		f.DebugfCall.Stub(param1, param2...)
	}
}
func (f *Logger) Debugln(param1 string) {
	f.DebuglnCall.Lock()
	defer f.DebuglnCall.Unlock()
	f.DebuglnCall.CallCount++
	f.DebuglnCall.Receives.Message = param1
	if f.DebuglnCall.Stub != nil {
		f.DebuglnCall.Stub(param1)
	}
}
func (f *Logger) Printf(param1 string, param2 ...interface {
}) {
	f.PrintfCall.Lock()
	defer f.PrintfCall.Unlock()
	f.PrintfCall.CallCount++
	f.PrintfCall.Receives.M = param1
	f.PrintfCall.Receives.A = param2
	if f.PrintfCall.Stub != nil {
		f.PrintfCall.Stub(param1, param2...)
	}
}
func (f *Logger) PromptWithDetails(param1 string, param2 string) bool {
	f.PromptWithDetailsCall.Lock()
	defer f.PromptWithDetailsCall.Unlock()
	f.PromptWithDetailsCall.CallCount++
	f.PromptWithDetailsCall.Receives.ResourceType = param1
	f.PromptWithDetailsCall.Receives.ResourceName = param2
	if f.PromptWithDetailsCall.Stub != nil {
		return f.PromptWithDetailsCall.Stub(param1, param2)
	}
	return f.PromptWithDetailsCall.Returns.Proceed
}
