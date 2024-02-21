package logstruct

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
)

type SqlNullValue interface {
	Value() (driver.Value, error)
}

func Fields(str interface{}) logrus.Fields {
	lf := logrus.Fields{}

	v := reflect.ValueOf(str)
	var t reflect.Type

	// если это поинтер
	if v.Kind() == reflect.Pointer {
		// получить значение
		v = reflect.Indirect(v)
		t = v.Type()
	} else {
		t = reflect.TypeOf(str)
	}

	// не структуры не обрабатываются
	if v.Kind() != reflect.Struct {
		return lf
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.FieldByName(field.Name)
		kind := value.Kind()

		switch kind {
		case reflect.Struct:
			// считываются только тэги log
			tag, ok := field.Tag.Lookup("log")
			if ok && tag != "-" {
				// парсинг конкретных типов библиотек struct
				switch v := value.Interface().(type) {
				case time.Time:
					lf[tag] = v.String()
					continue
				case SqlNullValue:
					value, _ := v.Value()

					if value == nil {
						value = "-"
					}
					switch getV := value.(type) {
					case time.Time:
						lf[tag] = getV.String()
					default:
						lf[tag] = getV
					}
					continue
				}
			}

			templf := Fields(value.Interface())
			for k, v := range templf {
				lf[k] = v
			}
			continue
		}

		// считываются только тэги log
		tag, ok := field.Tag.Lookup("log")
		if ok && tag != "-" {
			// если филд структуры: другая структура или поинтер - пропускаются
			if kind == reflect.Pointer {
				continue
			}

			switch kind {
			case reflect.Slice:
				s := reflect.ValueOf(value.Interface())

				var arrStr string
				for i := 0; i < s.Len(); i++ {
					arrStr += fmt.Sprintf("%v", s.Index(i))
					if i+1 < s.Len() {
						arrStr += ","
					}
				}

				lf[tag] = arrStr

				continue
			case reflect.Map:
				b, err := json.Marshal(value.Interface())
				if err == nil {
					lf[tag] = string(b)

					continue
				}
			}

			lf[tag] = value.Interface()
		}

	}

	return lf
}
