package helper

import "golang.org/x/crypto/bcrypt"

// HashPassword Hash password using the bcrypt hashing algorithm
func HashPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Hash password with bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

// CompareHashAndPassword Check if two passwords match using Bcrypt's CompareHashAndPassword
// which return nil on success and an error on failure.
func CompareHashAndPassword(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}
