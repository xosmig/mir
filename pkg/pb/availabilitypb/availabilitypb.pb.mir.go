package availabilitypb

import (
	reflect "reflect"
)

func (*RequestCertOrigin) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*RequestCertOrigin_ContextStore)(nil)),
		reflect.TypeOf((*RequestCertOrigin_Dsl)(nil)),
	}
}

func (*VerifyCertOrigin) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*VerifyCertOrigin_ContextStore)(nil)),
		reflect.TypeOf((*VerifyCertOrigin_Dsl)(nil)),
	}
}
