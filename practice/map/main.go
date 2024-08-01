package main

import "fmt"

type Person struct {
	Id   int
	Name string
}

func main() {

	list := createList()
	pMaps := make(map[int]Person, len(list))

	// Iterate over each `Person` in the `list`
	for _, val := range list {
		// Check if the `Person` with the same `Id` already exists in `pMaps`
		if _, ok := pMaps[val.Id]; !ok {
			// If it doesn't exist, add it to `pMaps` with the `Id` as the key and the `Person` as the value
			pMaps[val.Id] = val
		}
	}

	fmt.Printf("Type: %T, Length: %d \n", pMaps, len(pMaps))

	for key, val := range pMaps {
		fmt.Println(key)
		fmt.Println(val)
	}

	fmt.Println("---")
	fmt.Println(findInMaps(pMaps, 1).Name)
	fmt.Println(findInMaps(pMaps, 6))
}

func findInMaps(maps map[int]Person, id int) *Person {
	if val, ok := maps[id]; ok {
		return &val
	}
	return nil
}

func createList() []Person {
	list := []Person{
		{1, "John"},
		{2, "Jane"},
		{3, "Joe"},
		{4, "Jill"},
	}
	fmt.Printf("Type: %T, Length: %d, Capacity: %d\n", list, len(list), cap(list))

	list = append(list, Person{5, "Jim"})
	fmt.Printf("Type: %T, Length: %d, Capacity: %d\n", list, len(list), cap(list))

	return list
}
