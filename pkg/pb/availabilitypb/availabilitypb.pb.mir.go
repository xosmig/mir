package availabilitypb

import (
	reflect "reflect"
)

type Event_Type = isEvent_Type

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_RequestCert)(nil)),
		reflect.TypeOf((*Event_NewCert)(nil)),
		reflect.TypeOf((*Event_VerifyCert)(nil)),
		reflect.TypeOf((*Event_CertVerified)(nil)),
		reflect.TypeOf((*Event_RequestTransactions)(nil)),
		reflect.TypeOf((*Event_ProvideTransactions)(nil)),
	}
}

type RequestCertOrigin_Type = isRequestCertOrigin_Type

func (*RequestCertOrigin) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*RequestCertOrigin_ContextStore)(nil)),
		reflect.TypeOf((*RequestCertOrigin_Dsl)(nil)),
	}
}

type VerifyCertOrigin_Type = isVerifyCertOrigin_Type

func (*VerifyCertOrigin) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*VerifyCertOrigin_ContextStore)(nil)),
		reflect.TypeOf((*VerifyCertOrigin_Dsl)(nil)),
	}
}
