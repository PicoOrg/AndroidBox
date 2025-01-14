package ndk

// #include <string.h>
// #include <malloc.h>
// #include <sys/system_properties.h>
import "C"
import (
	"fmt"
	"unsafe"
)

func SystemPropertySet(name, value string) (err error) {
	nameCStr, valueCStr := C.CString(name), C.CString(value)
	defer C.free(unsafe.Pointer(nameCStr))
	defer C.free(unsafe.Pointer(valueCStr))
	rc := C.__system_property_set(nameCStr, valueCStr)
	if rc == 0 {
		return nil
	} else {
		return fmt.Errorf("system_property_set error code: %d", rc)
	}
}

func SystemPropertyGet(name string) (value string, err error) {
	nameCStr := C.CString(name)
	defer C.free(unsafe.Pointer(nameCStr))
	buffer := make([]byte, C.PROP_VALUE_MAX)
	rc := C.__system_property_get(nameCStr, (*C.char)(unsafe.Pointer(&buffer[0])))
	if rc >= 0 {
		return string(buffer[:rc]), nil
	} else {
		return "", fmt.Errorf("system_property_get error code: %d", rc)
	}
}
