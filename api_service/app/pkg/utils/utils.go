package utils

import "fmt"

func ConvertMapIntoStrings(From map[string]string) string {
	res := " "
	for key, value := range From {
		res += fmt.Sprintf("\"%s\": %s,", key, value)
	}
	return res
}
