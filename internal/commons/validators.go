package commons

import "reflect"

func ContainsEmptyValues(obj interface{}) bool {
	value := reflect.ValueOf(obj)
	numOfFields := value.NumField()
	for i := 0; i < numOfFields; i++ {
		if value.Field(i).IsZero() {
			return true
		}
	}
	return false
}
