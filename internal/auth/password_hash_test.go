package auth

import (
	"testing"
)

func TestPasswordHash(t *testing.T) {
	cases := []struct {
		password string
	}{
		{
	  		password: "",
		},
		{
			password: "test",
		},
		{
			password: "12345_",
		},
		{
			password: "abcd1234!",
		},
	}

	for _, c := range cases {
		hashedPassword, err := HashPassword(c.password)
		if err != nil {
			t.Errorf("Password hash failed: %v", err)
		}

		if err = CheckPasswordHash(hashedPassword, c.password); err != nil {
			t.Errorf("Passwords do not match: %v", err)
		}
	}

}