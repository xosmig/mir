package mockmodules

import (
	"github.com/budougumi0617/cmpmock"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/modules/mockmodules/internal/mock_internal"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/util/cmputil"
	"github.com/golang/mock/gomock"
	"github.com/xosmig/placeholders"
)

// MockPassiveModule is a slightly more user-friendly wrapper around gomock_modules.MockPassiveModule.
type MockPassiveModule struct {
	impl *mock_internal.MockModuleImpl
}

func NewMockPassiveModule(ctrl *gomock.Controller) *MockPassiveModule {
	return &MockPassiveModule{mock_internal.NewMockModuleImpl(ctrl)}
}

// ApplyEvents applies a list of input events to the module, making it advance its state
// and returns a (potentially empty) list of output events that the application of the input events results in.
func (m *MockPassiveModule) ApplyEvents(events *events.EventList) (*events.EventList, error) {
	return modules.ApplyEventsSequentially(events, m.impl.Event)
}

// ImplementsModule only serves the purpose of indicating that this is a Module and must not be called.
func (m *MockPassiveModule) ImplementsModule() {}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPassiveModule) EXPECT() *MockPassiveModuleMockRecorder {
	return &MockPassiveModuleMockRecorder{m.impl.EXPECT()}
}

// MockPassiveModuleMockRecorder is the mock recorder for MockPassiveModule.
type MockPassiveModuleMockRecorder struct {
	implRecorder *mock_internal.MockModuleImplMockRecorder
}

// Event indicates that the module is expected to receive an event that matches the argument.
// The argument is either an event (potentially, with placeholders created with placeholders.Make) or a gomock.Matcher.
// If the argument is an event, it is matched against the events received by the mock module using cmp.Diff with
// options placeholders.Ignore() and cmputil.IgnoreAllUnexported().
func (mr *MockPassiveModuleMockRecorder) Event(arg0 interface{}) *gomock.Call {
	if expectedEvent, ok := arg0.(*eventpb.Event); ok {
		// Instead of letting gomock to compare events with reflect.DeepEqual, compare them using the go-cmp package,
		// ignoring the unexported fields and with enabled placeholders.
		arg0 = cmpmock.DiffEq(expectedEvent, placeholders.Ignore(), cmputil.IgnoreAllUnexported())
	}
	return mr.implRecorder.Event(arg0)
}
