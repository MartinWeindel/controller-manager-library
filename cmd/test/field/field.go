package field

import (
	"bytes"
	"fmt"
	"github.com/gardener/controller-manager-library/pkg/fieldpath"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"reflect"
)

type S1 struct {
	Field1 S2
}

type S2 struct {
	Field2 string
	Field3 *string
	Field4 []string

	Field5 *S3

	Field6 []S3
}
type S3 struct {
	FieldA string
	FieldB string
}

type Config struct {
	Users map[string]*User `json:"users,omitempty"`
}
type User struct {
	Token *string `json:"token,omitempty"`
}

func FieldMain() {
	data, err := ioutil.ReadFile("local/test.yaml")
	if err == nil {

		d := yaml.NewDecoder(bytes.NewBuffer(data))

		for true {
			var doc interface{}
			err := d.Decode(&doc)
			if err != nil {
				if err != io.EOF {
					fmt.Printf("cannot parse yaml: %e", err)
					return
				}
				break
			}
			fmt.Printf("parsed doc\n")
		}
	} else {
		fmt.Printf("cannot read: %s\n", err)
	}
	array := [2]string{}

	slice := array[0:1]
	slice[0] = "slice0"
	fmt.Printf("initial %d s[0]: %s, a[0]: %s\n", cap(slice), slice[0], array[0])
	slice = append(slice, "slice01")
	fmt.Printf("append  %d s[1]: %s, a[1]: %s\n", cap(slice), slice[1], array[1])
	slice = append(slice, "slice02")
	slice[0] = "slice02"
	fmt.Printf("append  %d s[0]: %s, a[0]: %s\n", cap(slice), slice[0], array[0])

	slice1 := append(slice, "slice1")
	slice2 := append(slice, "slice2")
	fmt.Printf("twice   %d s1[3]: %s, s2[3]: %s\n", cap(slice1), slice1[3], slice2[3])

	slice11 := append(slice1, "slice11")
	slice12 := append(slice1, "slice12")
	fmt.Printf("twice   %d s11[4]: %s, s12[4]: %s\n", cap(slice11), slice11[4], slice12[4])

	//////////////////////////////////////////////////////////////////////////

	n, err := fieldpath.Compile(".Field1.Field5.FieldA")

	if err != nil {
		fmt.Printf("ERR: %s\n", err)
	} else {
		fmt.Printf(": %s\n", n.String())
	}

	s1 := S1{}

	str := "test"
	s1.Field1.Field2 = str
	s1.Field1.Field3 = &str
	s1.Field1.Field4 = []string{str}

	val, err := n.Get(&s1)

	if err != nil {
		fmt.Printf("ERR: %s\n", err)
	} else {
		fmt.Printf("= %v\n", fieldpath.Value(val))
	}

	err = n.Set(&s1, "NEW")

	if err != nil {
		fmt.Printf("ERR: %s\n", err)
	} else {
		val, _ := n.Get(&s1)
		fmt.Printf(":= %v\n", fieldpath.Value(val))
	}

	err = n.Set(&s1, 5)
	fmt.Printf("ERR: %s\n", err)
	err = n.ValidateType(&s1, reflect.TypeOf(5))
	fmt.Printf("ERR: %s\n", err)

	fmt.Printf("===============GET\n")

	f := fieldpath.RequiredField(&S1{}, ".Field1.Field4[0]")
	v, err := f.Get(s1)
	if err != nil {
		fmt.Printf("ERR: %s\n", err)
	} else {
		fmt.Printf(" %v == %v\n", s1.Field1.Field4[0], v)
	}
	fmt.Printf("===============SET\n")

	f = fieldpath.RequiredField(&S1{}, ".Field1.Field4[1]")
	err = f.Set(&s1, "NEWER")
	if err != nil {
		fmt.Printf("ERR: %s\n", err)
	} else {
		fmt.Printf(" %v\n", s1.Field1.Field4[1])
	}

	fmt.Printf("===============SELECT\n")
	f = fieldpath.RequiredField(&S1{}, ".Field1.Field6[1].FieldA")
	err = f.Set(&s1, "name")
	if err != nil {
		fmt.Printf("ERR: %s\n", err)
	} else {
		fmt.Printf(" %v\n", s1.Field1.Field6[1].FieldA)
	}

	n, err = fieldpath.Compile(".Field1.Field6[.FieldA=\"name\"].FieldB")
	if err != nil {
		fmt.Printf("ERR: %s\n", err)
	} else {
		fmt.Printf(": %s\n", n.String())
	}

	n, _ = fieldpath.Compile(".Field1.Field6[.FieldA=\"name2\"].FieldB")
	err = n.Set(&s1, "TEST")
	if err != nil {
		fmt.Printf("ERR: %s\n", err)
	} else {
		fmt.Printf(": %s %s\n", s1.Field1.Field6[2].FieldA, s1.Field1.Field6[2].FieldB)
	}

	f, err = fieldpath.NewField(&S1{}, ".Field1.Field6[.FieldA=\"name2\"].FieldB")
	if err != nil {
		fmt.Printf("ERR: %s\n", err)
	} else {
		fmt.Printf("field %s\n", f)
	}
}
