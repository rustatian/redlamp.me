// Package complexity is the learning set "Как понять O(n): время и память".
package complexity

import (
	_ "embed"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

//go:embed chapter1.html
var lesson1 string

// Set is the content of the learning set.
var Set = sets.LearningSet{
	Slug:    "complexity",
	Title:   "Как понять O(n): время и память",
	Summary: "Короткий урок на 15 минут. Считаем действия и дополнительную память на знакомых решениях. Восемь заданий с проверкой ответа прямо в уроке.",
	Topics:  []string{"Array", "Two Pointers"},
	Chapters: []sets.Chapter{
		{
			Title:   "Как понять O(n): время и память",
			Summary: "Максимум, копия списка, два прохода и Two Sum. От точного числа действий к O(1), O(n) и O(n²).",
			Lesson:  lesson1,
			Tasks:   []sets.Task{},
			Cards: []sets.Card{
				{ID: "complexity-scan", Type: "short", Prompt: "Поиск максимума начинает с best = a[0], затем сравнивает каждый оставшийся элемент с best. Во входе 6 чисел. Сколько сравнений выполнится?", Answer: "5", Explain: "Первое число задаёт best. Остаются 6 - 1 = 5 сравнений. Для n чисел время O(n)."},
				{ID: "complexity-memory", Type: "choice", Prompt: "Поиск максимума читает n чисел по индексам и хранит только i и best. Какова дополнительная память при росте n?", Options: []string{"O(1)", "O(n)", "O(n²)"}, Answer: 0, Explain: "Два рабочих значения перезаписываются. Их количество не растёт. Вход уже существует, результат содержит одно число."},
				{ID: "complexity-copy", Type: "choice", Prompt: "Дан список a из n чисел. Код выполняет b = a[:]. Как растут время и дополнительная память этой строки?", Options: []string{"Время O(1), память O(1)", "Время O(n), память O(n)", "Время O(n), память O(1)"}, Answer: 1, Explain: "Нужно скопировать n ссылок в новый список. Одна строка кода делает O(n) работы и занимает O(n) дополнительной памяти."},
				{ID: "complexity-passes", Type: "choice", Prompt: "Программа проходит n чисел, затем проходит их ещё раз. На каждом шаге постоянное число действий. Какова наиболее точная оценка времени из списка?", Options: []string{"O(1)", "O(n²)", "O(n)"}, Answer: 2, Explain: "Два последовательных прохода дают 2n посещений. При удвоении n их тоже вдвое больше: O(n)."},
				{ID: "complexity-pairs", Type: "short", Prompt: "В списке 5 чисел. Проверяем каждую пару индексов i < j без раннего выхода. Сколько пар проверим?", Answer: "10", Explain: "4 + 3 + 2 + 1 = 10. Для n чисел получается n(n - 1)/2 проверок: O(n²) времени."},
				{ID: "complexity-pointers", Type: "choice", Prompt: "Two Pointers ищет пару в уже упорядоченном списке. Каждая неудачная проверка сдвигает один край внутрь на одну позицию. Каковы время в худшем случае и дополнительная память?", Options: []string{"O(n²) и O(1)", "O(n) и O(n)", "O(n) и O(1)"}, Answer: 2, Explain: "Расстояние между краями уменьшается с n - 1 до нуля или до найденной пары. Проверок не больше n - 1. Храним только два индекса и сумму; ответ содержит два индекса."},
			},
			Links: []sets.Link{
				{Title: "Максимум и число вариантов", URL: "https://redlamp.me/leetcode/sets/permutations-and-combinations/1", Source: "Бомжарня"},
				{Title: "Two Pointers: навстречу", URL: "https://redlamp.me/leetcode/sets/two-pointers/1", Source: "Бомжарня"},
				{Title: "Queue, Deque и работа внутри pop(0)", URL: "https://redlamp.me/leetcode/sets/heaps-and-queues/1", Source: "Бомжарня"},
			},
		},
	},
}
