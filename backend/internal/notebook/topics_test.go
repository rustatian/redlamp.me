package notebook

import (
	"slices"
	"strings"
	"testing"
)

// The canonical cases of test/topics.test.ts. The topics of a result give the same result again.
func TestTopics(t *testing.T) {
	tests := []struct {
		name   string
		values []any
		slug   string
		want   []string
	}{
		{name: "Hash Table aliases count once", values: []any{" HashTable ", "Хеш-таблица", "hash\t TABLE"}, want: []string{"Hash Table"}},
		{name: "Array aliases count once", values: []any{"Arrays", "Массив", "Массивы", "array"}, want: []string{"Array"}},
		{name: "String and Two Pointers", values: []any{"Строка", "Два указателя", "two pointers"}, want: []string{"String", "Two Pointers"}},
		{name: "search and window topics", values: []any{"Двоичный поиск", "Скользящее окно"}, want: []string{"Binary Search", "Sliding Window"}},
		{name: "Heap and Monotonic Queue use LeetCode names", values: []any{"Heap", "heap (priority queue)", "Monotonic deque"}, want: []string{"Heap (Priority Queue)", "Monotonic Queue"}},
		{name: "linear containers have canonical case", values: []any{"queue", "deque", "stack"}, want: []string{"Queue", "Deque", "Stack"}},
		{name: "Russian names of supported containers", values: []any{"Очередь", "Дек", "Двусторонняя очередь", "Стек"}, want: []string{"Queue", "Deque", "Stack"}},
		{name: "Russian and English priority names share one topic", values: []any{"Куча", "Priority Queue", "Очередь с приоритетом", "Приоритетная очередь"}, want: []string{"Heap (Priority Queue)"}},
		{name: "Russian Monotonic Queue", values: []any{"Монотонная очередь", "Монотонный дек"}, want: []string{"Monotonic Queue"}},
		{name: "Russian Dynamic Programming name", values: []any{"Динамическое программирование"}, want: []string{"Dynamic Programming"}},
		{name: "common Graph and String forms", values: []any{"Graph", "Граф", "Графы", "Строки"}, want: []string{"Graph Theory", "String"}},
		{name: "alternate Hash Table and Binary Search spelling", values: []any{"Хэш-таблица", "Хеш таблица", "Hash-Table", "Бинарный поиск"}, want: []string{"Hash Table", "Binary Search"}},
		{name: "selection and divide topics", values: []any{"Быстрый выбор", "Разделяй и властвуй"}, want: []string{"Quickselect", "Divide and Conquer"}},
		{name: "Greedy and Simulation", values: []any{"Жадный алгоритм", "Моделирование"}, want: []string{"Greedy", "Simulation"}},
		{name: "mathematical topics", values: []any{"Математика", "Комбинаторика", "Рекурсия", "Матрица"}, want: []string{"Math", "Combinatorics", "Recursion", "Matrix"}},
		{name: "prefix topics", values: []any{"Префиксные суммы", "Префиксное дерево"}, want: []string{"Prefix Sum", "Trie"}},
		{name: "NFC before alias lookup", values: []any{"Подсчёт", "Подсчёт"}, want: []string{"Counting"}},
		{name: "Sorting topics remain distinct", values: []any{"Сортировка", "Сортировка подсчётом", "bucket sort"}, want: []string{"Sorting", "Counting Sort", "Bucket Sort"}},
		{name: "ordered and range topics", values: []any{"Упорядоченное множество", "Минимум и максимум на отрезке"}, want: []string{"Ordered Set", "Range Minimum/Maximum Query"}},
		{name: "Graph Theory uses the current name", values: []any{"Graphs", "graph theory"}, want: []string{"Graph Theory"}},
		{name: "counting objects", values: []any{"Перестановки", "Сочетания", "Подмножества", "Факториал"}, want: []string{"Permutations", "Combinations", "Subsets", "Factorial"}},
		{name: "counting rules", values: []any{"Правило произведения", "Правило суммы"}, want: []string{"Product Rule", "Sum Rule"}},
		{name: "order and Pruning", values: []any{"Словарный порядок", "Отсечение"}, want: []string{"Lexicographic Order", "Pruning"}},
		{name: "explicit Backtracking aliases", values: []any{"перебор с возвратом", "поиск с возвратом", "backtracking"}, want: []string{"Backtracking"}},
		{name: "Enumeration and Brute Force stay distinct", values: []any{"Enumeration", "Brute Force", "Backtracking"}, want: []string{"Enumeration", "Brute Force", "Backtracking"}},
		{name: "letter case context", values: []any{"Перебор"}, slug: "letter-case-permutation", want: []string{"Backtracking"}},
		{name: "subsets context", values: []any{"Перебор"}, slug: "subsets", want: []string{"Backtracking"}},
		{name: "permutations context", values: []any{"Перебор"}, slug: "permutations", want: []string{"Backtracking"}},
		{name: "combinations context", values: []any{"Перебор"}, slug: "combinations", want: []string{"Backtracking"}},
		{name: "phone context", values: []any{"Перебор"}, slug: "letter-combinations-of-a-phone-number", want: []string{"Backtracking"}},
		{name: "repeated subsets context", values: []any{"Перебор"}, slug: "subsets-ii", want: []string{"Backtracking"}},
		{name: "repeated permutations context", values: []any{"Перебор"}, slug: "permutations-ii", want: []string{"Backtracking"}},
		{name: "combination sum context", values: []any{"Перебор"}, slug: "combination-sum", want: []string{"Backtracking"}},
		{name: "parentheses context", values: []any{"Перебор"}, slug: "generate-parentheses", want: []string{"Backtracking"}},
		{name: "combination lock context", values: []any{"Перебор"}, slug: "cf-1097-b", want: []string{"Backtracking"}},
		{name: "queens context", values: []any{"Перебор"}, slug: "n-queens", want: []string{"Backtracking"}},
		{name: "Unique Paths table context", values: []any{"Таблица", "dynamic programming"}, slug: "unique-paths", want: []string{"Dynamic Programming"}},
		{name: "Top K Frequent Words correction", values: []any{"Сортировка подсчётом", "Bucket Sort"}, slug: "top-k-frequent-words", want: []string{"Bucket Sort"}},
		{name: "known Counting Sort stays distinct", values: []any{"Counting Sort", "Bucket Sort"}, slug: "top-k-frequent-words", want: []string{"Counting Sort", "Bucket Sort"}},
		{name: "unknown English topics keep first occurrence order", values: []any{"  Custom   Technique ", "HashTable", "custom technique", "Хеш-таблица"}, want: []string{"Custom Technique", "Hash Table"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := topics(tt.values, tt.slug)
			if err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("topics(%q, %q) = %q, %v; want %q", tt.values, tt.slug, got, err, tt.want)
			}
			again := make([]any, len(tt.want))
			for i, topic := range tt.want {
				again[i] = topic
			}
			if got, err := topics(again, tt.slug); err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("topics(%q, %q) = %q, %v; want the same topics", again, tt.slug, got, err)
			}
		})
	}
}

// The cases "rejects $name with an English-name instruction" of test/topics.test.ts.
func TestTopicsRejectNonEnglishNames(t *testing.T) {
	tests := []struct {
		name   string
		values []any
		slug   string
	}{
		{name: "unknown Russian topic", values: []any{"Особая тема"}},
		{name: "unknown Chinese topic", values: []any{"数组"}},
		{name: "ambiguous Enumeration without context", values: []any{"Перебор"}},
		{name: "ambiguous Enumeration in an unrelated problem", values: []any{"Перебор"}, slug: "two-sum"},
		{name: "table without context", values: []any{"Таблица"}},
		{name: "table in an unrelated problem", values: []any{"Таблица"}, slug: "two-sum"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := topics(tt.values, tt.slug)
			checkInvalid(t, err, "Укажите название темы на английском языке.")
		})
	}
}

// The cases "keeps validation for $name" of test/topics.test.ts. Each value gives an error.
func TestTopicsValidation(t *testing.T) {
	tests := []struct {
		name   string
		values []any
	}{
		{name: "too many aliases before deduplication", values: slices.Repeat([]any{"Array"}, 13)},
		{name: "topic longer than 40 characters", values: []any{strings.Repeat("a", 41)}},
		{name: "blank topic", values: []any{"  "}},
		{name: "non-string topic", values: []any{1.0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := topics(tt.values, ""); err == nil {
				t.Errorf("topics(%v) = %q, want an error", tt.values, got)
			}
		})
	}
}
