package clinic

import (
	"slices"
	"strings"
)

const maxDist = 3

func levenshtein(a, b string, maxDist int) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)

	// Убеждаемся что a — более короткая строка (меньше итераций внешнего цикла)
	if la > lb {
		ra, rb = rb, ra
		la, lb = lb, la
	}

	// Разница длин уже больше maxDist — результат точно хуже
	if lb-la > maxDist {
		return maxDist + 1
	}

	row := make([]int, lb+1)
	for j := range row {
		row[j] = j
	}

	for i := 1; i <= la; i++ {
		prev := row[0]
		row[0] = i

		jFrom := i - maxDist
		if jFrom < 1 {
			jFrom = 1
		}

		jTo := i + maxDist
		if jTo > lb {
			jTo = lb
		}

		if jFrom > 1 {
			row[jFrom-1] = maxDist + 1
		}

		minInRow := maxDist + 1

		for j := jFrom; j <= jTo; j++ {
			tmp := row[j]

			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}

			val := prev + cost // замена / совпадение
			if del := row[j] + 1; del < val {
				val = del // удаление
			}

			if ins := row[j-1] + 1; ins < val {
				val = ins // вставка
			}

			row[j] = val

			if val < minInRow {
				minInRow = val
			}

			prev = tmp
		}

		// Вся строка вышла за порог — дальше лучше не станет
		if minInRow > maxDist {
			return maxDist + 1
		}
	}

	if row[lb] > maxDist {
		return maxDist + 1
	}

	return row[lb]
}

func FindSimilar(query string, dict []string) []string {
	var result []string
	dict = Map[string, string](dict, func(s string) string {
		return strings.ToLower(s)
	})

	for _, word := range dict {
		if levenshtein(query, word, maxDist) <= maxDist {
			result = append(result, word)
		}
	}
	result = slices.Compact(result)
	return result
}

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}
