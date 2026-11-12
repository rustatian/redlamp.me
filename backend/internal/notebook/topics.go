package notebook

import (
	"slices"
	"strings"

	"golang.org/x/text/unicode/norm"
)

// topicAliases is the table of shared/topics.ts: the LeetCode name of a topic and its other spellings.
var topicAliases = []struct {
	name    string
	aliases []string
}{
	{"Hash Table", []string{"HashTable", "Hash-Table", "Хеш-таблица", "Хэш-таблица", "Хеш таблица", "Хэш таблица"}},
	{"Array", []string{"Arrays", "Массив", "Массивы"}},
	{"String", []string{"Строка", "Строки"}},
	{"Two Pointers", []string{"Два указателя"}},
	{"Binary Search", []string{"Двоичный поиск", "Бинарный поиск"}},
	{"Sliding Window", []string{"Скользящее окно"}},
	{"Sorting", []string{"Сортировка"}},
	{"Heap (Priority Queue)", []string{"Heap", "Priority Queue", "Куча", "Очередь с приоритетом", "Приоритетная очередь"}},
	{"Monotonic Queue", []string{"Monotonic deque", "Монотонная очередь", "Монотонный дек"}},
	{"Queue", []string{"Очередь"}},
	{"Deque", []string{"Дек", "Двусторонняя очередь"}},
	{"Stack", []string{"Стек"}},
	{"Quickselect", []string{"Быстрый выбор"}},
	{"Divide and Conquer", []string{"Разделяй и властвуй"}},
	{"Greedy", []string{"Жадный алгоритм"}},
	{"Simulation", []string{"Моделирование"}},
	{"Math", []string{"Математика"}},
	{"Combinatorics", []string{"Комбинаторика"}},
	{"Recursion", []string{"Рекурсия"}},
	{"Matrix", []string{"Матрица"}},
	{"Prefix Sum", []string{"Префиксные суммы"}},
	{"Trie", []string{"Префиксное дерево"}},
	{"Counting", []string{"Подсчёт"}},
	{"Counting Sort", []string{"Сортировка подсчётом"}},
	{"Bucket Sort", nil},
	{"Ordered Set", []string{"Упорядоченное множество"}},
	{"Range Minimum/Maximum Query", []string{"Минимум и максимум на отрезке"}},
	{"Graph Theory", []string{"Graph", "Graphs", "Граф", "Графы"}},
	{"Permutations", []string{"Перестановки"}},
	{"Combinations", []string{"Сочетания"}},
	{"Subsets", []string{"Подмножества"}},
	{"Factorial", []string{"Факториал"}},
	{"Product Rule", []string{"Правило произведения"}},
	{"Sum Rule", []string{"Правило суммы"}},
	{"Lexicographic Order", []string{"Словарный порядок"}},
	{"Pruning", []string{"Отсечение"}},
	{"Backtracking", []string{"перебор с возвратом", "поиск с возвратом"}},
	{"Dynamic Programming", []string{"Динамическое программирование"}},
	{"Enumeration", nil},
	{"Brute Force", nil},
}

// topicNames maps each spelling in lower case to the LeetCode name.
var topicNames = func() map[string]string {
	names := map[string]string{}
	for _, topic := range topicAliases {
		names[strings.ToLower(topic.name)] = topic.name
		for _, alias := range topic.aliases {
			names[strings.ToLower(norm.NFC.String(alias))] = topic.name
		}
	}
	return names
}()

// backtrackingProblems are the problems where the Russian word "перебор" means backtracking.
var backtrackingProblems = []string{"letter-case-permutation", "subsets", "permutations", "combinations", "letter-combinations-of-a-phone-number", "subsets-ii", "permutations-ii", "combination-sum", "generate-parentheses", "cf-1097-b", "n-queens"}

// normalizeTopic returns the LeetCode name of a topic. slug is the problem of the topic, or "".
func normalizeTopic(value, slug string) string {
	name := strings.Join(strings.Fields(norm.NFC.String(value)), " ")
	key := strings.ToLower(name)
	switch {
	case key == "перебор" && slices.Contains(backtrackingProblems, slug):
		return "Backtracking"
	case key == "таблица" && slug == "unique-paths":
		return "Dynamic Programming"
	case key == "сортировка подсчётом" && slug == "top-k-frequent-words":
		return "Bucket Sort"
	}
	if canonical, ok := topicNames[key]; ok {
		return canonical
	}
	return name
}

// normalizeTopics normalizes each topic and keeps the first of the topics that differ only in letter case.
func normalizeTopics(values []string, slug string) []string {
	result := make([]string, 0, len(values))
	seen := map[string]bool{}
	for _, value := range values {
		topic := normalizeTopic(value, slug)
		if key := strings.ToLower(topic); !seen[key] {
			seen[key] = true
			result = append(result, topic)
		}
	}
	return result
}

// topics checks a decoded JSON value: a list of at most 12 topic names. It returns the normalized names, which are printable ASCII.
func topics(value any, slug string) ([]string, error) {
	items, ok := value.([]any)
	if !ok || len(items) > 12 {
		return nil, invalid("Укажите не больше 12 тем.")
	}
	values := make([]string, len(items))
	for i, item := range items {
		topic, err := text(item, 40, true)
		if err != nil {
			return nil, err
		}
		values[i] = topic
	}
	result := normalizeTopics(values, slug)
	for _, topic := range result {
		for i := range len(topic) {
			if topic[i] < 0x20 || topic[i] > 0x7e {
				return nil, invalid("Укажите название темы на английском языке.")
			}
		}
	}
	return result, nil
}

// searchText is the form of a text for the search: NFC and lower case.
func searchText(value string) string {
	return strings.ToLower(norm.NFC.String(value))
}
