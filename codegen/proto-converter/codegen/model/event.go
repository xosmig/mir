package model

type Event interface {
	Message() *Message
}

type 

type EventLeaf struct {
}

func (*EventLeaf) isEvent() {}

type EventNode interface {
	Event
	isEventNode()
}

type EventInternalNode struct {
	EventNode
}

func (*EventInternalNode) isEvent()     {}
func (*EventInternalNode) isEventNode() {}

type EventRootNode struct {
	EventNode
}

func (*EventRootNode) isEvent()     {}
func (*EventRootNode) isEventNode() {}
