package services

import "strconv"

func StringsToInts(values ...string) ([]int, error) {
	results := make([]int, len(values))
	for i := 0; i < len(values); i++ {
		temp, err := strconv.Atoi(values[i])
		if err != nil {
			return nil, err
		}
		results[i] = temp
	}

	return results, nil
}
