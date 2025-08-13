package pkg

import "fmt"

func BuildTag(index int) string {
	return fmt.Sprintf("a%d", index)
}
