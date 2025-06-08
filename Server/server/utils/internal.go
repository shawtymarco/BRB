package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/sandertv/gophertunnel/minecraft/text"
)

func ReadConfig[T any](file string) (T, error) {
	var zero T
	if _, err := os.Stat(file); os.IsNotExist(err) {
		data, err := json.Marshal(zero)
		if err != nil {
			return zero, fmt.Errorf("encode default config: %v", err)
		}
		if err := os.WriteFile(file, data, 0644); err != nil {
			return zero, fmt.Errorf("create default config: %v", err)
		}
		return zero, err
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return zero, fmt.Errorf("read config: %v", err)
	}
	if err := json.Unmarshal(data, &zero); err != nil {
		return zero, fmt.Errorf("decode config: %v", err)
	}
	return zero, err
}

func Filter[T comparable](slice []T, callable func(v T) bool) []T {
	var newSlice []T

	for _, v := range slice {
		if callable(v) {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}

func RandChance(n int) bool {
	return n > rand.Intn(100)
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func RandShuffle[T any](slice []T) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

func IntToRoman(num int) string {
	roman := ""
	numbers := []int{1, 4, 5, 9, 10, 40, 50, 90, 100, 400, 500, 900, 1000}
	romans := []string{"I", "IV", "V", "IX", "X", "XL", "L", "XC", "C", "CD", "D", "CM", "M"}
	index := len(romans) - 1

	for num > 0 {
		for numbers[index] <= num {
			roman += romans[index]
			num -= numbers[index]
		}
		index -= 1
	}

	return roman
}

func Rad2deg(rad float64) float64 {
	return rad * (180 / math.Pi)
}

func ShortenNumber(num float64) string {
	units := []string{"", text.Colourf("<yellow>K</yellow>"), text.Colourf("<dark-red>M</dark-red>"), text.Colourf("<purple>B</purple>"), text.Colourf("<black>T</black>"), text.Colourf("<red>QD</red>"), text.Colourf("<dark-rerd>QT</dark-red>")}
	for i := 0; i < len(units); i++ {
		if num >= 1000 {
			num /= 1000
		} else {
			return fmt.Sprintf("%.1f%s", num, units[i])
		}
	}
	return fmt.Sprintf("%.1f%s", num, units[len(units)-1]) // Handle numbers beyond trillions
}

func IsPrime(n int64) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}

	if n%2 == 0 || n%3 == 0 {
		return false
	}

	for i := int64(5); i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}

	return true
}

func NextPrime(n int64) int64 {
	if n <= 1 {
		return 2
	}

	prime := n + 1
	for !IsPrime(prime) {
		prime++
	}
	return prime
}

func Patch(url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("PATCH", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return http.DefaultClient.Do(req)
}

func Delete(url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return http.DefaultClient.Do(req)
}

func Question[T any](cond bool, o1, o2 T) T {
	if cond {
		return o1
	}
	return o2
}
