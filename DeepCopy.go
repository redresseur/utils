package utils

import "reflect"

func DeepCopy(dstValue reflect.Value, srcValue reflect.Value, paramType reflect.Type){
	switch paramType.Kind() {
	case reflect.Ptr:
		if srcValue.Elem().Kind() == reflect.Struct{
			if dstValue.Pointer() == 0{
				dstValue.Set(reflect.New(srcValue.Elem().Type()))
			}
			DeepCopy(dstValue.Elem(), srcValue.Elem(), srcValue.Elem().Type())
		}

	case reflect.Slice:
		slice_line := srcValue.Len()
		if slice_line == 0{
			break
		}

		dstValue.Set(reflect.MakeSlice(srcValue.Type() , slice_line, slice_line))
		for i := 0; i < srcValue.Len(); i++{
			DeepCopy(dstValue.Index(i), srcValue.Index(i), srcValue.Index(i).Type())
		}
	case reflect.Struct:

		for i := 0; i < srcValue.NumField(); i++ {
			DeepCopy(dstValue.Field(i), srcValue.Field(i), srcValue.Field(i).Type())
		}
	default:
		dstValue.Set(srcValue)
	}

	return
}
