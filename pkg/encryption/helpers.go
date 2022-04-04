package encryption

import "fmt"

func ValidateAESSecret(plainSecret string) error {
	k := len([]byte(plainSecret))
	switch k {
	default:
		return fmt.Errorf("only key size of either 16, 24, or 32 bytes allowed")
	case 16, 24, 32:
		return nil
	}
}
