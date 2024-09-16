package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
)

func ReadFile(fileName string) string {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return string(data)

}

// Decode a number from a given base to decimal
func decodeBase(value string, base int) (int, error) {
	decimalValue, err := strconv.ParseInt(value, base, 64)
	if err != nil {
		return 0, err
	}
	return int(decimalValue), nil
}

// Helper function to parse JSON number (handle float64)
func parseJSONNumber(value json.RawMessage) float64 {
	var number float64
	if err := json.Unmarshal(value, &number); err != nil {
		fmt.Println("Error parsing number:", err)
	}
	return number
}

// Lagrange interpolation to find the polynomial value at x = 0
func lagrangeInterpolation(xValues []int, yValues []int) int {
	n := len(xValues)
	var result big.Int

	for i := 0; i < n; i++ {
		var term big.Int
		term.SetInt64(int64(yValues[i]))

		for j := 0; j < n; j++ {
			if i != j {
				var denominator, numerator big.Int

				denominator.Sub(big.NewInt(int64(xValues[i])), big.NewInt(int64(xValues[j])))
				numerator.Sub(big.NewInt(0), big.NewInt(int64(xValues[j])))

				term.Mul(&term, &numerator)
				term.Div(&term, &denominator)
			}
		}
		result.Add(&result, &term)
	}

	return int(result.Int64())
}

func main() {
	filenames := []string{"testcase.json", "testcase2.json"}
	for i, filename := range filenames {

		inputJSON := ReadFile(filename)

		var data map[string]json.RawMessage
		if err := json.Unmarshal([]byte(inputJSON), &data); err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		// Unmarshal keys
		var keys map[string]json.RawMessage
		if err := json.Unmarshal(data["keys"], &keys); err != nil {
			fmt.Println("Error parsing keys:", err)
			return
		}
		n := int(parseJSONNumber(keys["n"]))
		k := int(parseJSONNumber(keys["k"]))

		if k > n {
			fmt.Println("Not enough points provided")
			return
		}

		var xValues []int
		var yValues []int

		for key, value := range data {
			if key == "keys" {
				continue
			}
			var entry map[string]string
			if err := json.Unmarshal(value, &entry); err != nil {
				fmt.Println("Error parsing entry:", err)
				return
			}
			base, err := strconv.Atoi(entry["base"])
			if err != nil {
				fmt.Println("Error converting base to integer:", err)
				return
			}
			valueStr := entry["value"]

			decodedValue, err := decodeBase(valueStr, base)
			if err != nil {
				fmt.Println("Error decoding base value:", err)
				return
			}
			x, _ := strconv.Atoi(key)

			xValues = append(xValues, x)
			yValues = append(yValues, decodedValue)
		}

		c := lagrangeInterpolation(xValues, yValues)
		fmt.Printf("The constant term 'c' for the Test Case %d is: %d\n", i+1, c)
	}
}
