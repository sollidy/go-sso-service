package main

import "fmt"

type Person struct {
	Id   int
	Name string
}

func main() {
	list := []Person{
		{1, "John"},
		{2, "Jane"},
		{3, "Joe"},
		{4, "Jill"},
	}
	fmt.Printf("Type: %T, Length: %d, Capacity: %d\n", list, len(list), cap(list))

	list = append(list, Person{5, "Jim"})
	fmt.Printf("Type: %T, Length: %d, Capacity: %d\n", list, len(list), cap(list))

	pMaps := make(map[int]Person, len(list))

	for _, val := range list {
		if _, ok := pMaps[val.Id]; !ok {
			pMaps[val.Id] = val
		}
	}

	fmt.Printf("Type: %T, Length: %d \n", pMaps, len(pMaps))

	for key, val := range pMaps {
		fmt.Println(key)
		fmt.Println(val)
	}
	fmt.Println("---")
	fmt.Println(findInMaps(pMaps, 1))
	fmt.Println(findInMaps(pMaps, 6))

	res := findInMaps(pMaps, 1)
	fmt.Println(res.Name)
}

func findInMaps(maps map[int]Person, id int) *Person {
	if val, ok := maps[id]; ok {
		return &val
	}
	return nil
}
