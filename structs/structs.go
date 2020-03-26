package structs

import(
	"fmt"
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
	"path"
	"reflect"
	"strings"
	"github.com/nsgtest/packages/interfaces"
)

var slice []interfaces.Interface

type Struct struct{
	File	string
	Object	interfaces.Interface
}

func (s Struct) Add(matches []int){
	s.Read()
	indices := s.Find(matches)

	if len(indices) != len(slice){
		fmt.Println("FAIL!")
		fmt.Println("Current:")
		fmt.Printf("\n")
		s.Output()
		fmt.Printf("\n")
		fmt.Printf("\nFound similiar %v!\n", reflect.TypeOf(s.Object).Name())
		panic(nil)
	}

	slice = append(slice, s.Object)
	s.Write()
}

func (s Struct) Update(matches []int){
	s.Read()
	indices := s.Find(matches)

	if len(indices) == 0{
		fmt.Println("FAIL!")
		fmt.Println("Current:")
		fmt.Printf("\n")
		s.Output()
		fmt.Printf("\n")
		fmt.Println("Found nothing similar!")
	} else {
		for _, index := range indices{
			reference := reflect.New(reflect.TypeOf(s.Object)).Elem()

			for i := 0; i < reflect.TypeOf(s.Object).NumField(); i++{
				if !reflect.ValueOf(s.Object).Field(i).IsZero(){
					reference.Field(i).Set(reflect.ValueOf(s.Object).Field(i))
				} else {
					reference.Field(i).Set(reflect.ValueOf(slice[index]).Field(i))
				}
			}

			reflect.ValueOf(&slice[index]).Elem().Set(reference)
		}
	}

	s.Write()
}

func (s Struct) Remove(matches []int){
	s.Read()
	indices := s.Find(matches)

	if len(indices) == 0{
		fmt.Println("FAIL!")
		fmt.Println("Current:")
		fmt.Printf("\n")
		s.Output()
		fmt.Printf("\n")
		fmt.Println("Found nothing similar!")
	}

	for i, index := range indices{
		slice = append(slice[:index-i], slice[index+1-i:]...)
	}

	s.Write()
}

func (s Struct) List(){
	s.Read()

	if len(slice) < 1{
		fmt.Println("FAIL!")
		fmt.Printf("%v is empty!\n", s.File)
		panic(nil)
	}

	for i, object := range slice{
		s.Object = object
		s.Output()
		if i != len(slice) - 1{
			fmt.Printf("\n")
		}
	}
}

func (s Struct) Find(matches []int) []int{
	indices := []int{}

	for i, object := range slice{
		count := 0

		for j := 0; j < reflect.TypeOf(object).NumField(); j++{
			if reflect.ValueOf(s.Object).Field(j).Interface() == reflect.ValueOf(object).Field(j).Interface(){
				count += int(math.Pow(2, float64(j)))
			}
		}

		for _, match := range matches{
			if count == match{
				indices = append(indices, i)
				break
			}
		}

	}
	return indices
}

func (s Struct) Read(){
	enc, err := ioutil.ReadFile(s.File)
	if err != nil{
		fmt.Println("FAIL!")
		fmt.Printf("Could not read from %v!\n", s.File)
		panic(err)
	}

	references := reflect.New(reflect.SliceOf(reflect.TypeOf(s.Object))).Interface()
	err = json.Unmarshal(enc, references)
	if err != nil{
		fmt.Println("FAIL!")
		fmt.Println("%v is not a JSON file!\n", s.File)
		panic(err)
	}

	slice = make([]interfaces.Interface, reflect.ValueOf(references).Elem().Len())
	for i := 0; i < reflect.ValueOf(references).Elem().Len(); i++{
		reflect.ValueOf(&slice[i]).Elem().Set(reflect.ValueOf(references).Elem().Index(i))
	}
}

func (s Struct) Write(){
	_, err := os.Stat(path.Dir(s.File))
	if os.IsNotExist(err){
		fmt.Println("FAIL!")
		fmt.Printf("Directory %v does not exist!\n", path.Dir(s.File))
		panic(err)
	}

	enc, _ := json.Marshal(slice)

	err = ioutil.WriteFile(s.File, enc, 0666)
	if err != nil{
		fmt.Println("FAIL!")
		fmt.Printf("Could not write to %v!\n", s.File)
		panic(err)
	}
}

func (s Struct) Output(){
	for i := 0; i < reflect.TypeOf(s.Object).NumField(); i++{
		fmt.Printf("%v: %v\n", reflect.TypeOf(s.Object).Field(i).Name, strings.Trim(fmt.Sprintf("%v", reflect.ValueOf(s.Object).Field(i)), "[]"))
	}
}
