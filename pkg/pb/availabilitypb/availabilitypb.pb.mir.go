package availabilitypb

import (
	reflect "reflect"
)

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
