package threshcryptopb

import (
	contextstorepb "github.com/filecoin-project/mir/pkg/pb/contextstorepb"
	dslpb "github.com/filecoin-project/mir/pkg/pb/dslpb"
)

type Event_Type = isEvent_Type

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func (w *Event_SignShare) Unwrap() *SignShare {
	return w.SignShare
}

func (w *Event_SignShareResult) Unwrap() *SignShareResult {
	return w.SignShareResult
}

func (w *Event_VerifyShare) Unwrap() *VerifyShare {
	return w.VerifyShare
}

func (w *Event_VerifyShareResult) Unwrap() *VerifyShareResult {
	return w.VerifyShareResult
}

func (w *Event_VerifyFull) Unwrap() *VerifyFull {
	return w.VerifyFull
}

func (w *Event_VerifyFullResult) Unwrap() *VerifyFullResult {
	return w.VerifyFullResult
}

func (w *Event_Recover) Unwrap() *Recover {
	return w.Recover
}

func (w *Event_RecoverResult) Unwrap() *RecoverResult {
	return w.RecoverResult
}

type SignShareOrigin_Type = isSignShareOrigin_Type

type SignShareOrigin_TypeWrapper[T any] interface {
	SignShareOrigin_Type
	Unwrap() *T
}

func (w *SignShareOrigin_ContextStore) Unwrap() *contextstorepb.Origin {
	return w.ContextStore
}

func (w *SignShareOrigin_Dsl) Unwrap() *dslpb.Origin {
	return w.Dsl
}

type VerifyShareOrigin_Type = isVerifyShareOrigin_Type

type VerifyShareOrigin_TypeWrapper[T any] interface {
	VerifyShareOrigin_Type
	Unwrap() *T
}

func (w *VerifyShareOrigin_ContextStore) Unwrap() *contextstorepb.Origin {
	return w.ContextStore
}

func (w *VerifyShareOrigin_Dsl) Unwrap() *dslpb.Origin {
	return w.Dsl
}

type VerifyFullOrigin_Type = isVerifyFullOrigin_Type

type VerifyFullOrigin_TypeWrapper[T any] interface {
	VerifyFullOrigin_Type
	Unwrap() *T
}

func (w *VerifyFullOrigin_ContextStore) Unwrap() *contextstorepb.Origin {
	return w.ContextStore
}

func (w *VerifyFullOrigin_Dsl) Unwrap() *dslpb.Origin {
	return w.Dsl
}

type RecoverOrigin_Type = isRecoverOrigin_Type

type RecoverOrigin_TypeWrapper[T any] interface {
	RecoverOrigin_Type
	Unwrap() *T
}

func (w *RecoverOrigin_ContextStore) Unwrap() *contextstorepb.Origin {
	return w.ContextStore
}

func (w *RecoverOrigin_Dsl) Unwrap() *dslpb.Origin {
	return w.Dsl
}
