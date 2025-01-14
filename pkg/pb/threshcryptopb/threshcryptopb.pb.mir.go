package threshcryptopb

import (
	reflect "reflect"
)

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_SignShare)(nil)),
		reflect.TypeOf((*Event_SignShareResult)(nil)),
		reflect.TypeOf((*Event_VerifyShare)(nil)),
		reflect.TypeOf((*Event_VerifyShareResult)(nil)),
		reflect.TypeOf((*Event_VerifyFull)(nil)),
		reflect.TypeOf((*Event_VerifyFullResult)(nil)),
		reflect.TypeOf((*Event_Recover)(nil)),
		reflect.TypeOf((*Event_RecoverResult)(nil)),
	}
}

func (*SignShareOrigin) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*SignShareOrigin_ContextStore)(nil)),
		reflect.TypeOf((*SignShareOrigin_Dsl)(nil)),
	}
}

func (*VerifyShareOrigin) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*VerifyShareOrigin_ContextStore)(nil)),
		reflect.TypeOf((*VerifyShareOrigin_Dsl)(nil)),
	}
}

func (*VerifyFullOrigin) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*VerifyFullOrigin_ContextStore)(nil)),
		reflect.TypeOf((*VerifyFullOrigin_Dsl)(nil)),
	}
}

func (*RecoverOrigin) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*RecoverOrigin_ContextStore)(nil)),
		reflect.TypeOf((*RecoverOrigin_Dsl)(nil)),
	}
}
