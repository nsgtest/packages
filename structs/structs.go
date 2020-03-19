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

type Struct struct{
	File	string
	Object	interfaces.Interface
	Array	interfaces.Interfaces
}

func (s Struct) Add(matches []int){
	indices, counts := s.Find(matches)

	if len(indices) > 0{
		fmt.Printf("FAIL!\n\nCurrent:\n")
		s.Output()
		fmt.Printf("\n")
		for i, index := range indices{
			s.Array[index].Message(counts[i])
			fmt.Printf("\n")
			s.Object = s.Array[index]
			s.Output()
			fmt.Printf("\n")

			if i == len(indices) - 1{
				panic(nil)
			}
		}
	}

	s.Array = append(s.Array, s.Object)
	s.Write()
}

func (s Struct) Update(matches []int){
	indices, _ := s.Find(matches)
	if len(indices) == 0{
		fmt.Printf("FAIL!\n\nCurrent:\n")
		s.Output()
		fmt.Printf("\n")
		s.Object.Message(0)
		fmt.Print("\n")
	} else {
		for _, index := range indices{
			for i := 0; i < reflect.TypeOf(s.Object).NumField(); i++{
				if !reflect.ValueOf(s.Object).Field(i).IsZero() && reflect.ValueOf(s.Object).Field(i).Interface() != reflect.ValueOf(s.Array[index]).Field(i).Interface() {
					reflect.ValueOf(s.Array[index]).Field(i).Set(reflect.ValueOf(s.Object).Field(i))
				}
			}
		}
	}
}

func (s Struct) Remove(matches []int){
	indices, _ := s.Find(matches)

	if len(indices) == 0{
		fmt.Printf("FAIL!\n\nCurrent:\n")
		s.Output()
		fmt.Printf("\n")
		s.Object.Message(0)
		fmt.Print("\n")
	}

	for i, index := range indices{
		s.Array = append(s.Array[:index-i], s.Array[index+1-i:]...)
	}

	s.Write()
}

func (s Struct) List(){
	if len(s.Array) < 1{
		fmt.Printf("FAIL!\n%v is empty!\n\n", s.File)
		panic(nil)
	}

	for i, object := range s.Array{
		s.Object = object
		s.Output()
		if i != len(s.Array) - 1{
			fmt.Printf("\n")
		}
	}
}

func (s Struct) Find(matches []int) ([]int, []int){
	indices := []int{}
	counts := []int{}

	for i, object := range s.Array{
		count := 0

		for j := 0; j < reflect.TypeOf(object).NumField(); j++{
			if reflect.ValueOf(s.Object).Field(j).Interface() == reflect.ValueOf(object).Field(j).Interface(){
				count += int(math.Pow(2, float64(j)))
			}
		}

		for _, match := range matches{
			if count == match{
				indices = append(indices, i)
				counts = append(counts, count)
				break
			}
		}

	}
	return indices, counts
}

func (s Struct) Write(){
	_, err := os.Stat(path.Dir(s.File))
	if os.IsNotExist(err){
		fmt.Printf("FAIL!\nDirectory %v does not exist!\n\n", path.Dir(s.File))
		panic(err)
	}

	enc, _ := json.Marshal(s.Array)

	err = ioutil.WriteFile(s.File, enc, 0666)
	if err != nil{
		fmt.Printf("FAIL!\nCould not write to %v\n\n", s.File)
		panic(err)
	}
}

func (s Struct) Output(){
	for i := 0; i < reflect.TypeOf(s.Object).NumField(); i++{
		fmt.Printf("%v: %v\n", reflect.TypeOf(s.Object).Field(i).Name, strings.Trim(fmt.Sprintf("%v", reflect.ValueOf(s.Object).Field(i)), "[]"))
	}
}
