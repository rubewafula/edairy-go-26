package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateMemberNo() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(999999))
}

func GenerateSupplierNo() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(999999))
}
