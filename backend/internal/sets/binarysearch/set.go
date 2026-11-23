// Package binarysearch is the learning set "Binary Search".
package binarysearch

import (
	_ "embed"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

//go:embed chapter1.html
var lesson1 string

//go:embed chapter2.html
var lesson2 string

//go:embed chapter3.html
var lesson3 string

//go:embed chapter4.html
var lesson4 string

func number(value int64) *int64 { return &value }

func text(value string) *string { return &value }

// Set is the content of the learning set.
var Set = sets.LearningSet{
	Slug:    "binary-search",
	Title:   "Binary Search",
	Summary: "Binary Search отбрасывает половину кандидатов за одну проверку и находит число, границу или наименьший подходящий ответ за O(log n) шагов. Вы держите инвариант lo и hi, ищете Lower Bound и Upper Bound через bisect, ищете ответ по монотонной проверке и работаете с горой, повёрнутым Array и Matrix. Для каждого решения считаете проверки, дополнительную память и размер ответа.",
	Topics:  []string{"Binary Search", "Array", "Sorting", "Prefix Sum", "Matrix"},
	Chapters: []sets.Chapter{
		{
			Title:   "Binary Search в упорядоченном Array",
			Summary: "Binary Search держит границы lo и hi так, что искомое число, если оно есть, лежит между ними. Каждая проверка середины отбрасывает половину кандидатов. Вы считаете проверки, находите ошибки границ и ищете число в Array, который задан формулой.",
			Lesson:  lesson1,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/binary-search/",
					Number:     number(704),
					Title:      "Binary Search",
					Difficulty: "easy",
					Topics:     []string{"Array", "Binary Search"},
					Note:       "Binary Search держит индекс target, если он есть, в отрезке [lo, hi]. n = len(nums). Время: O(log n), не больше ⌊log₂ n⌋ + 1 проверок середины, для n = 10⁴ это 14. Память: O(1) дополнительно: lo, hi, mid; ответ O(1). Linear Scan даёт O(n), до 10⁴ сравнений. Для [−1, 0, 3, 5, 9, 12] и 9: середины 3 и 9, ответ 4.",
					Statement: `<p>Дан Array целых чисел <code>nums</code>, упорядоченный по возрастанию, и целое число <code>target</code>. Напишите функцию, которая ищет <code>target</code> в <code>nums</code>. Если <code>target</code> есть в <code>nums</code>, верните его индекс. Иначе верните <code>-1</code>.</p><p>Нужен алгоритм со временем <code>O(log n)</code>.</p><p>Метод класса <code>Solution</code>: <code>def search(self, nums: list[int], target: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [-1,0,3,5,9,12], target = 9
<strong>Выход:</strong> 4
<strong>Пояснение:</strong> 9 есть в nums, его индекс 4.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [-1,0,3,5,9,12], target = 2
<strong>Выход:</strong> -1
<strong>Пояснение:</strong> 2 нет в nums, поэтому верните -1.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 10<sup>4</sup></code></li><li><code>-10<sup>4</sup> &lt; nums[i], target &lt; 10<sup>4</sup></code></li><li>все числа в <code>nums</code> <strong>разные</strong></li><li><code>nums</code> упорядочен по возрастанию</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "search"}, Cases: []sets.TestCase{
						{Input: "[[-1,0,3,5,9,12],9]", Expected: text("4")},
						{Input: "[[-1,0,3,5,9,12],2]", Expected: text("-1")},
						{Input: "[[5],5]", Expected: text("0")},
						{Input: "[[1,3],0]", Expected: text("-1")},
						{Input: "[[-9999,-5,0,7,9999],9999]", Expected: text("4")},
						{Input: "[[1,3,5,7,9,11,13,15,17],1]", Expected: text("0")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/valid-perfect-square/",
					Number:     number(367),
					Title:      "Valid Perfect Square",
					Difficulty: "easy",
					Topics:     []string{"Math", "Binary Search"},
					Note:       "Binary Search по x на отрезке [1, num]: x · x растёт вместе с x. Время: O(log num), не больше 31 проверки при num ≤ 2³¹ − 1. Память: O(1) дополнительно, ответ bool. Linear Scan x = 1, 2, ... делает до 46 341 проверки: O(√num). Для num = 16: середины 8 и 4, ответ True.",
					Statement: `<p>Дано положительное целое число <code>num</code>. Верните <code>True</code>, <em>если</em> <code>num</code> <em>является точным квадратом, иначе верните</em> <code>False</code>.</p><p><strong>Точный квадрат</strong> - это целое число, которое равно квадрату целого числа. Другими словами, это произведение некоторого целого числа на само себя.</p><p>Нельзя использовать встроенные функции библиотек, например <code>sqrt</code>.</p><p>Метод класса <code>Solution</code>: <code>def isPerfectSquare(self, num: int) -&gt; bool</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> num = 16
<strong>Выход:</strong> True
<strong>Пояснение:</strong> Верните True, потому что 4 * 4 = 16 и 4 - целое число.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> num = 14
<strong>Выход:</strong> False
<strong>Пояснение:</strong> Верните False, потому что 3.742 * 3.742 = 14 и 3.742 - не целое число.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ num ≤ 2<sup>31</sup> - 1</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "isPerfectSquare"}, Cases: []sets.TestCase{
						{Input: "[16]", Expected: text("True")},
						{Input: "[14]", Expected: text("False")},
						{Input: "[1]", Expected: text("True")},
						{Input: "[2147395600]", Expected: text("True")},
						{Input: "[2147483647]", Expected: text("False")},
						{Input: "[2147395599]", Expected: text("False")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/minimum-common-value/",
					Number:     number(2540),
					Title:      "Minimum Common Value",
					Difficulty: "easy",
					Topics:     []string{"Array", "Hash Table", "Two Pointers", "Binary Search"},
					Note:       "Для каждого x из nums1 по порядку ищите x в nums2 через Binary Search; первое найденное x и есть ответ. n = len(nums1), m = len(nums2). Время: O(n log m), до n · (⌊log₂ m⌋ + 1) проверок, для n = m = 10⁵ это 1,7 · 10⁶. Память: O(1) дополнительно, ответ O(1). Two Pointers из набора «Two Pointers» дают O(n + m) времени и O(1) памяти. set(nums2) даёт O(n + m) в среднем, но O(m) памяти.",
					Statement: `<p>Даны два Array целых чисел <code>nums1</code> и <code>nums2</code>, упорядоченные по неубыванию (каждое следующее число не меньше предыдущего). Верните <em><strong>наименьшее общее целое число</strong> этих двух Array</em>. Если у <code>nums1</code> и <code>nums2</code> нет общего целого числа, верните <code>-1</code>.</p><p>Целое число называют <strong>общим</strong> для <code>nums1</code> и <code>nums2</code>, если оно встречается в каждом из двух Array <strong>хотя бы один раз</strong>.</p><p>Метод класса <code>Solution</code>: <code>def getCommon(self, nums1: list[int], nums2: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums1 = [1,2,3], nums2 = [2,4]
<strong>Выход:</strong> 2
<strong>Пояснение:</strong> Наименьшее общее число двух Array равно 2, поэтому верните 2.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums1 = [1,2,3,6], nums2 = [2,3,4,5]
<strong>Выход:</strong> 2
<strong>Пояснение:</strong> Общих чисел два: 2 и 3. Наименьшее из них 2, поэтому верните 2.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums1.length, nums2.length ≤ 10<sup>5</sup></code></li><li><code>1 ≤ nums1[i], nums2[j] ≤ 10<sup>9</sup></code></li><li><code>nums1</code> и <code>nums2</code> упорядочены по <strong>неубыванию</strong></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "getCommon"}, Cases: []sets.TestCase{
						{Input: "[[1,2,3],[2,4]]", Expected: text("2")},
						{Input: "[[1,2,3,6],[2,3,4,5]]", Expected: text("2")},
						{Input: "[[1],[1]]", Expected: text("1")},
						{Input: "[[1,2,3],[4,5,6]]", Expected: text("-1")},
						{Input: "[[5,5,5,7],[1,7,7]]", Expected: text("7")},
						{Input: "[[1000000000],[1,2,1000000000]]", Expected: text("1000000000")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/192/A",
					Title:      "Funky Numbers",
					Difficulty: "medium",
					Topics:     []string{"Binary Search", "Brute Force", "Implementation"},
					Note:       "Codeforces, рейтинг 1300. Выпишите треугольные числа T(k) < n в list и для каждого ищите n − T(k) в этом list через Binary Search. K - число таких T(k), K ≤ 44 720 при n ≤ 10⁹. Время: O(K log K), около 44 720 · 16 ≈ 7 · 10⁵ проверок. Память: O(K) для list, ответ O(1). Brute Force по всем парам даёт O(K²) ≈ 2 · 10⁹ сравнений и не успеет. Для n = 256: 3 + 253 = T(2) + T(22), ответ YES.",
					Statement:  `<p>Как вы хорошо знаете, в этом году самые модные числа - это так называемые треугольные числа, то есть целые числа, которые можно записать как <code>k(k + 1) / 2</code>, где <code>k</code> - некоторое целое положительное число. А самые крутые числа - те, которые можно записать как сумму двух треугольных чисел.</p><p>Известный хипстер Андрей обожает всё модное и крутое, но, к сожалению, плохо знает математику. Дано число <code>n</code>. Помогите ему узнать, можно ли записать это число как сумму двух треугольных чисел (не обязательно разных)!</p><h3>Входные данные</h3><p>В первой строке записано целое число <code>n</code> (<code>1 ≤ n ≤ 10<sup>9</sup></code>).</p><h3>Выходные данные</h3><p>Выведите «YES» (без кавычек), если <code>n</code> можно записать как сумму двух треугольных чисел, иначе выведите «NO» (без кавычек).</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>256</pre><p><strong>Выходные данные:</strong></p><pre>YES</pre><h3>Пример 2</h3><p><strong>Входные данные:</strong></p><pre>512</pre><p><strong>Выходные данные:</strong></p><pre>NO</pre><h3>Примечание</h3><p>В первом примере <code>256 = 2 · 3 / 2 + 22 · 23 / 2</code>.</p><p>Во втором примере число 512 нельзя записать как сумму двух треугольных чисел.</p><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 10<sup>9</sup></code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "256\n", Expected: text("YES\n")},
						{Input: "512\n", Expected: text("NO\n")},
						{Input: "1\n", Expected: text("NO\n")},
						{Input: "2\n", Expected: text("YES\n")},
						{Input: "999999999\n", Expected: text("YES\n")},
						{Input: "1000000000\n", Expected: text("NO\n")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "bs1-checks", Type: "short", Prompt: "Упорядоченный Array содержит 1000 чисел. Сколько проверок середины сделает Binary Search в худшем случае?", Answer: "10", Explain: "Перед каждой проверкой остаётся кандидатов: 1000, 500, 250, 125, 62, 31, 15, 7, 3, 1. Это 10 проверок, то есть ⌊log₂ 1000⌋ + 1 = 9 + 1. Ответ 9 забывает последнюю проверку, когда остался один кандидат. Linear Scan в худшем случае проверит все 1000 чисел. Дополнительная память Binary Search O(1)."},
				{ID: "bs1-trace", Type: "short", Prompt: "a = [1, 4, 6, 9, 13, 20, 27], target = 6. Код из урока начинает с lo = 0, hi = 6. Сколько проверок середины он сделает?", Answer: "3", Explain: "mid = 3: a[3] = 9 > 6, hi = 2. mid = 1: a[1] = 4 < 6, lo = 2. mid = 2: a[2] = 6, ответ 2. Всего 3 проверки. Ответ 2 получится, если не считать последнюю проверку, которая нашла число."},
				{ID: "bs1-while", Type: "choice", Prompt: "В коде из урока цикл записан как while lo < hi. На каком входе функция вернёт −1, хотя target есть в Array?", Options: []string{"a = [1, 5], target = 1", "a = [1, 5, 9], target = 5", "a = [5], target = 5"}, Answer: 2, Explain: "Для [5] сразу lo = hi = 0, условие lo < hi ложно, и цикл не начнётся: ответ −1 вместо 0. Для [1, 5] и [1, 5, 9] первая же середина равна target. Ошибка видна, когда в отрезке остался один кандидат; условие lo ≤ hi проверяет и его."},
				{ID: "bs1-lo", Type: "code", Prompt: "Ветка a[mid] < target. Какая строка стоит на месте ___?", Language: "python", Code: `lo, hi = 0, len(a) - 1
while lo <= hi:
    mid = (lo + hi) // 2
    if a[mid] == target:
        return mid
    if a[mid] < target:
        ___
    else:
        hi = mid - 1
return -1`, Options: []string{"lo = mid", "lo = mid + 1", "hi = mid + 1"}, Answer: 1, Explain: "a[mid] < target, поэтому mid и всё левее не ответ: lo = mid + 1. С lo = mid вход a = [1, 3], target = 3 даёт lo = 0, hi = 1, mid = 0 снова и снова: цикл не закончится. hi = mid + 1 не меняет lo и не сужает отрезок: на том же входе тоже снова mid = 0."},
				{ID: "bs1-queries", Type: "short", Prompt: "В упорядоченном Array 300 000 чисел. Нужно ответить на 200 000 запросов «есть ли x». Сколько проверок середины сделают все Binary Search вместе в худшем случае?", Answer: "3800000", Explain: "Один поиск делает ⌊log₂ 300 000⌋ + 1 = 18 + 1 = 19 проверок, ведь 2¹⁸ = 262 144 ≤ 300 000 < 524 288 = 2¹⁹. Всего 19 · 200 000 = 3 800 000. Ответ 3 600 000 забывает последнюю проверку каждого поиска. Linear Scan на каждый запрос дал бы до 300 000 · 200 000 = 6 · 10¹⁰ сравнений. Дополнительная память O(1) на запрос, ответы хранятся отдельно."},
				{ID: "bs1-slice", Type: "choice", Prompt: "Recursion со срезом: если a[mid] < target, функция вызывает search(a[mid + 1:], target), где mid = (len(a) - 1) // 2. В a 1000 упорядоченных чисел, target больше всех. Сколько ссылок скопируют все срезы вместе?", Options: []string{"994", "10", "10 000"}, Answer: 0, Explain: "Срезы имеют длины 500, 250, 125, 62, 31, 15, 7, 3, 1: вместе 994 ссылки, почти n. Все они живы, пока идут вложенные вызовы: время и дополнительная память O(n), а не O(log n). Версия с числами lo и hi ничего не копирует. 10 - это число проверок середины для длин 1000, 500, ..., 1, а не число копий. 10 000 получилось бы, если бы каждый из 10 вызовов копировал все 1000 ссылок."},
				{ID: "bs1-square", Type: "short", Prompt: "LeetCode 367, num = 30. Поиск идёт по x от lo = 1 до hi = 30, mid = (lo + hi) // 2. Сколько проверок середины он сделает до ответа False?", Answer: "5", Explain: "Середины 15, 7, 3, 5, 6 дают квадраты 225, 49, 9, 25, 36. После 36 > 30 отрезок [6, 5] пуст: 30 не квадрат. Ответ 4 забывает последнюю проверку: отрезок [6, 6] ещё не пуст, и mid = 6 нужно проверить. Для num около 2³¹ Binary Search делает до 31 проверки, а Linear Scan по x до 46 341."},
				{ID: "bs1-funky", Type: "choice", Prompt: "Codeforces 192A, но пусть n ≤ 10⁶. Треугольных чисел меньше n не больше 1413. Сколько примерно действий сделают Brute Force по всем парам и Binary Search второго числа для каждого первого?", Options: []string{"около 1413 и около 11", "около 10⁶ и около 20", "около 2 · 10⁶ и около 1,6 · 10⁴"}, Answer: 2, Explain: "Пар 1413 · 1413 ≈ 2 · 10⁶. Binary Search ищет n − T(i) за ⌊log₂ 1413⌋ + 1 = 10 + 1 = 11 проверок, ведь 2¹⁰ = 1024 ≤ 1413 < 2048; всего 1413 · 11 = 15 543. 11 - проверки одного поиска, а не всего решения. При n ≤ 10⁹ разница ещё больше: около 2 · 10⁹ пар против 7 · 10⁵ проверок. list треугольных чисел занимает O(√n) памяти."},
			},
			Links: []sets.Link{
				{Title: "Binary Search в упорядоченном Array", URL: "https://ru.wikipedia.org/wiki/%D0%94%D0%B2%D0%BE%D0%B8%D1%87%D0%BD%D1%8B%D0%B9_%D0%BF%D0%BE%D0%B8%D1%81%D0%BA", Source: "Википедия"},
				{Title: "Binary Search по целым числам", URL: "https://neerc.ifmo.ru/wiki/index.php?title=%D0%A6%D0%B5%D0%BB%D0%BE%D1%87%D0%B8%D1%81%D0%BB%D0%B5%D0%BD%D0%BD%D1%8B%D0%B9_%D0%B4%D0%B2%D0%BE%D0%B8%D1%87%D0%BD%D1%8B%D0%B9_%D0%BF%D0%BE%D0%B8%D1%81%D0%BA", Source: "Викиконспекты ИТМО"},
				{Title: "Binary Search: идея и реализация", URL: "https://ikcprog.github.io/topics/binsearch/", Source: "Школа олимпиадного программирования ИКЦ"},
				{Title: "Почему в Binary Search легко ошибиться", URL: "https://codechick.io/journal/binary-search-pochemu-vse-pishut-s-bagom", Source: "Журнал CodeChick"},
			},
			Figures: []sets.Figure{
				{
					ID:    "bs1-search",
					Title: "Binary Search: поиск 23",
					Frames: []sets.FigureFrame{
						{Caption: "lo = 0, hi = 9, mid = 4. a[4] = 16 меньше 23, поэтому 23 может быть только правее: lo = 5.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 5, 8, 12, 16, 23, 38, 56, 72, 91}, Marked: []int{4}, Pointers: &sets.Pointers{{Name: "lo", Cell: 0}, {Name: "mid", Cell: 4}, {Name: "hi", Cell: 9}}, Window: &[2]int{0, 9}},
						}},
						{Caption: "lo = 5, hi = 9, mid = 7. a[7] = 56 больше 23: hi = 6. Затемнены индексы, которые больше не кандидаты.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 5, 8, 12, 16, 23, 38, 56, 72, 91}, Marked: []int{7}, Dimmed: []int{0, 1, 2, 3, 4}, Pointers: &sets.Pointers{{Name: "lo", Cell: 5}, {Name: "mid", Cell: 7}, {Name: "hi", Cell: 9}}, Window: &[2]int{5, 9}},
						}},
						{Caption: "lo = 5, hi = 6, mid = 5. a[5] = 23: ответ 5 после 3 проверок. Linear Scan проверил бы 6 чисел.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 5, 8, 12, 16, 23, 38, 56, 72, 91}, Marked: []int{5}, Dimmed: []int{0, 1, 2, 3, 4, 7, 8, 9}, Pointers: &sets.Pointers{{Name: "lo", Cell: 5}, {Name: "mid", Cell: 5}, {Name: "hi", Cell: 6}}, Window: &[2]int{5, 6}},
						}},
					},
				},
				{
					ID:    "bs1-miss",
					Title: "Binary Search: числа 40 нет",
					Frames: []sets.FigureFrame{
						{Caption: "mid = 4: 16 < 40, lo = 5.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 5, 8, 12, 16, 23, 38, 56, 72, 91}, Marked: []int{4}, Pointers: &sets.Pointers{{Name: "lo", Cell: 0}, {Name: "mid", Cell: 4}, {Name: "hi", Cell: 9}}, Window: &[2]int{0, 9}},
						}},
						{Caption: "mid = 7: 56 > 40, hi = 6.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 5, 8, 12, 16, 23, 38, 56, 72, 91}, Marked: []int{7}, Dimmed: []int{0, 1, 2, 3, 4}, Pointers: &sets.Pointers{{Name: "lo", Cell: 5}, {Name: "mid", Cell: 7}, {Name: "hi", Cell: 9}}, Window: &[2]int{5, 9}},
						}},
						{Caption: "mid = 5: 23 < 40, lo = 6.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 5, 8, 12, 16, 23, 38, 56, 72, 91}, Marked: []int{5}, Dimmed: []int{0, 1, 2, 3, 4, 7, 8, 9}, Pointers: &sets.Pointers{{Name: "lo", Cell: 5}, {Name: "mid", Cell: 5}, {Name: "hi", Cell: 6}}, Window: &[2]int{5, 6}},
						}},
						{Caption: "mid = 6: 38 < 40, lo = 7. Это 4-я проверка: для n = 10 больше не бывает.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 5, 8, 12, 16, 23, 38, 56, 72, 91}, Marked: []int{6}, Dimmed: []int{0, 1, 2, 3, 4, 5, 7, 8, 9}, Pointers: &sets.Pointers{{Name: "lo", Cell: 6}, {Name: "mid", Cell: 6}, {Name: "hi", Cell: 6}}, Window: &[2]int{6, 6}},
						}},
						{Caption: "lo = 7, hi = 6: отрезок пуст, и цикл while lo ≤ hi закончился. 40 стояло бы между 38 и 56. Ответ −1.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 5, 8, 12, 16, 23, 38, 56, 72, 91}, Dimmed: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, Pointers: &sets.Pointers{{Name: "lo", Cell: 7}, {Name: "hi", Cell: 6}}},
						}},
					},
				},
				{
					ID:    "bs1-squares",
					Title: "Квадраты без хранения: num = 50",
					Frames: []sets.FigureFrame{
						{Caption: "lo = 1, hi = 50, mid = 25: 625 > 50, hi = 24.", Rows: []sets.FigureRow{
							{Label: "lo / mid / hi", Cells: []any{1, 25, 50}, Marked: []int{1}},
							{Label: "mid · mid", Cells: []any{625}},
						}},
						{Caption: "mid = 12: 144 > 50, hi = 11.", Rows: []sets.FigureRow{
							{Label: "lo / mid / hi", Cells: []any{1, 12, 24}, Marked: []int{1}},
							{Label: "mid · mid", Cells: []any{144}},
						}},
						{Caption: "mid = 6: 36 < 50, lo = 7.", Rows: []sets.FigureRow{
							{Label: "lo / mid / hi", Cells: []any{1, 6, 11}, Marked: []int{1}},
							{Label: "mid · mid", Cells: []any{36}},
						}},
						{Caption: "mid = 9: 81 > 50, hi = 8.", Rows: []sets.FigureRow{
							{Label: "lo / mid / hi", Cells: []any{7, 9, 11}, Marked: []int{1}},
							{Label: "mid · mid", Cells: []any{81}},
						}},
						{Caption: "mid = 7: 49 < 50, lo = 8.", Rows: []sets.FigureRow{
							{Label: "lo / mid / hi", Cells: []any{7, 7, 8}, Marked: []int{1}},
							{Label: "mid · mid", Cells: []any{49}},
						}},
						{Caption: "mid = 8: 64 > 50, hi = 7.", Rows: []sets.FigureRow{
							{Label: "lo / mid / hi", Cells: []any{8, 8, 8}, Marked: []int{1}},
							{Label: "mid · mid", Cells: []any{64}},
						}},
						{Caption: "lo = 8 больше hi = 7: целого x с x · x = 50 нет. 6 проверок, ответ False. Квадраты 49 и 64 стоят по обе стороны от 50.", Rows: []sets.FigureRow{
							{Label: "lo / hi", Cells: []any{8, 7}},
						}},
					},
				},
			},
		},
		{
			Title:   "Lower Bound, Upper Bound и bisect",
			Summary: "Часто нужно не само число, а граница: первая позиция, где условие начинает выполняться. Lower Bound и Upper Bound находят её за O(log n), а модуль bisect делает это в Python. Вы считаете числа в отрезке, отвечаете на много запросов и ищете номер группы по Prefix Sum.",
			Lesson:  lesson2,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/search-insert-position/",
					Number:     number(35),
					Title:      "Search Insert Position",
					Difficulty: "easy",
					Topics:     []string{"Array", "Binary Search"},
					Note:       "Ответ - Lower Bound: первый индекс i с nums[i] ≥ target, или n, если такого нет. n = len(nums). Время: O(log n), не больше 14 проверок для n = 10⁴. Память: O(1) дополнительно, ответ O(1). bisect_left(nums, target) делает то же. Linear Scan делает до 10⁴ сравнений. Для [1, 3, 5, 6] и 2: середины 5, 3, 1, ответ 1.",
					Statement: `<p>Дан упорядоченный Array из разных целых чисел и число <code>target</code>. Если <code>target</code> есть в Array, верните его индекс. Если нет, верните индекс, на котором он стоял бы после вставки с сохранением порядка.</p><p>Нужен алгоритм со временем <code>O(log n)</code>.</p><p>Метод класса <code>Solution</code>: <code>def searchInsert(self, nums: list[int], target: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,3,5,6], target = 5
<strong>Выход:</strong> 2</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [1,3,5,6], target = 2
<strong>Выход:</strong> 1</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> nums = [1,3,5,6], target = 7
<strong>Выход:</strong> 4</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 10<sup>4</sup></code></li><li><code>-10<sup>4</sup> ≤ nums[i] ≤ 10<sup>4</sup></code></li><li>все числа в <code>nums</code> <strong>разные</strong> и упорядочены по <strong>возрастанию</strong></li><li><code>-10<sup>4</sup> ≤ target ≤ 10<sup>4</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "searchInsert"}, Cases: []sets.TestCase{
						{Input: "[[1,3,5,6],5]", Expected: text("2")},
						{Input: "[[1,3,5,6],2]", Expected: text("1")},
						{Input: "[[1,3,5,6],7]", Expected: text("4")},
						{Input: "[[1,3,5,6],0]", Expected: text("0")},
						{Input: "[[7],7]", Expected: text("0")},
						{Input: "[[-10000,10000],0]", Expected: text("1")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/706/B",
					Title:      "Interesting drink",
					Difficulty: "easy",
					Topics:     []string{"Array", "Binary Search", "Sorting"},
					Note:       "Codeforces, рейтинг 1100. Упорядочите цены один раз; ответ дня - bisect_right(x, m), число цен ≤ m. n - число магазинов, q - число дней. Время: O(n log n + q log n). Память: O(n) для упорядоченной копии и O(q) для ответов. Linear Scan каждого дня даёт O(nq) = 10¹⁰ сравнений при n = q = 10⁵. Счётчик цен с Prefix Sum по значениям даёт O(n + q + V) при V = 10⁵. Для цен [3, 10, 8, 6, 11] и m = 10: ответ 4, bisect_left дал бы 3.",
					Statement: `<p>Василий любит отдыхать после тяжёлой работы, поэтому его часто можно встретить в баре рядом. Как и все программисты, он любит знаменитый напиток «Beecola». Его продают в <code>n</code> разных магазинах города. Известно, что бутылка напитка в магазине <code>i</code> стоит <code>x<sub>i</sub></code> монет.</p><p>Василий хочет покупать любимый напиток <code>q</code> дней подряд. Он знает, что в <code>i</code>-й день сможет потратить <code>m<sub>i</sub></code> монет. Для каждого дня он хочет узнать, в скольких разных магазинах он сможет купить бутылку «Beecola».</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>n</code> (<code>1 ≤ n ≤ 100 000</code>): количество магазинов в городе, где продают любимый напиток Василия.</p><p>Во второй строке записаны <code>n</code> целых чисел <code>x<sub>i</sub></code> (<code>1 ≤ x<sub>i</sub> ≤ 100 000</code>): цены бутылки напитка в <code>i</code>-м магазине.</p><p>В третьей строке записано одно целое число <code>q</code> (<code>1 ≤ q ≤ 100 000</code>): количество дней, в которые Василий хочет купить напиток.</p><p>Дальше идут <code>q</code> строк. В <code>i</code>-й из них записано одно целое число <code>m<sub>i</sub></code> (<code>1 ≤ m<sub>i</sub> ≤ 10<sup>9</sup></code>): количество монет, которое Василий может потратить в <code>i</code>-й день.</p><h3>Выходные данные</h3><p>Выведите <code>q</code> целых чисел. <code>i</code>-е из них должно быть равно количеству магазинов, где Василий сможет купить бутылку напитка в <code>i</code>-й день.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>5
3 10 8 6 11
4
1
10
3
11</pre><p><strong>Выходные данные:</strong></p><pre>0
4
1
5</pre><h3>Примечание</h3><p>В первый день Василий не сможет купить напиток ни в одном магазине.</p><p>Во второй день Василий может купить напиток в магазинах 1, 2, 3 и 4.</p><p>В третий день Василий может купить напиток только в магазине номер 1.</p><p>Наконец, в последний день Василий может купить напиток в любом магазине.</p><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 100 000</code></li><li><code>1 ≤ x<sub>i</sub> ≤ 100 000</code></li><li><code>1 ≤ q ≤ 100 000</code></li><li><code>1 ≤ m<sub>i</sub> ≤ 10<sup>9</sup></code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "5\n3 10 8 6 11\n4\n1\n10\n3\n11\n", Expected: text("0\n4\n1\n5\n")},
						{Input: "1\n5\n3\n4\n5\n6\n", Expected: text("0\n1\n1\n")},
						{Input: "4\n2 2 2 2\n3\n1\n2\n1000000000\n", Expected: text("0\n4\n4\n")},
						{Input: "6\n7 3 7 1 10 3\n4\n3\n6\n7\n9\n", Expected: text("3\n3\n5\n5\n")},
						{Input: "3\n100000 1 50000\n2\n99999\n100000\n", Expected: text("2\n3\n")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/find-first-and-last-position-of-element-in-sorted-array/",
					Number:     number(34),
					Title:      "Find First and Last Position of Element in Sorted Array",
					Difficulty: "medium",
					Topics:     []string{"Array", "Binary Search"},
					Note:       "Ответ [bisect_left, bisect_right − 1], если bisect_left < bisect_right, иначе [−1, −1]. n = len(nums). Время: O(log n), два поиска по ⌊log₂ n⌋ + 1 проверок, для n = 10⁵ до 34. Память: O(1) дополнительно, ответ из двух чисел. Поиск одного вхождения и шаги в стороны дают O(n) на Array из одинаковых чисел.",
					Statement: `<p>Дан Array целых чисел <code>nums</code>, упорядоченный по неубыванию (каждое следующее число не меньше предыдущего). Найдите начальную и конечную позиции заданного числа <code>target</code>.</p><p>Если <code>target</code> нет в Array, верните <code>[-1, -1]</code>.</p><p>Нужен алгоритм со временем <code>O(log n)</code>.</p><p>Метод класса <code>Solution</code>: <code>def searchRange(self, nums: list[int], target: int) -&gt; list[int]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [5,7,7,8,8,10], target = 8
<strong>Выход:</strong> [3,4]</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [5,7,7,8,8,10], target = 6
<strong>Выход:</strong> [-1,-1]</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> nums = [], target = 0
<strong>Выход:</strong> [-1,-1]</pre><h3>Ограничения</h3><ul><li><code>0 ≤ nums.length ≤ 10<sup>5</sup></code></li><li><code>-10<sup>9</sup> ≤ nums[i] ≤ 10<sup>9</sup></code></li><li><code>nums</code> упорядочен по неубыванию</li><li><code>-10<sup>9</sup> ≤ target ≤ 10<sup>9</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "searchRange"}, Cases: []sets.TestCase{
						{Input: "[[5,7,7,8,8,10],8]", Expected: text("[3,4]")},
						{Input: "[[5,7,7,8,8,10],6]", Expected: text("[-1,-1]")},
						{Input: "[[],0]", Expected: text("[-1,-1]")},
						{Input: "[[1],1]", Expected: text("[0,0]")},
						{Input: "[[2,2,2,2,2],2]", Expected: text("[0,4]")},
						{Input: "[[-1000000000,0,1000000000],1000000000]", Expected: text("[2,2]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/successful-pairs-of-spells-and-potions/",
					Number:     number(2300),
					Title:      "Successful Pairs of Spells and Potions",
					Difficulty: "medium",
					Topics:     []string{"Array", "Two Pointers", "Binary Search", "Sorting"},
					Note:       "Упорядочите potions. Пара успешна, если potion ≥ ⌈success / spell⌉ = (success + spell − 1) // spell; число таких зелий m − bisect_left(potions, need). n = len(spells), m = len(potions). Время: O(m log m + n log m). Память: O(m) для упорядоченной копии, ответ O(n). Brute Force проверяет n · m пар: 10¹⁰ при n = m = 10⁵. Для spell = 3, success = 7: need = 3, подходят 3 зелья из [1, 2, 3, 4, 5].",
					Statement: `<p>Даны два Array положительных целых чисел <code>spells</code> и <code>potions</code> длины <code>n</code> и <code>m</code>. <code>spells[i]</code> - сила <code>i</code>-го заклинания, <code>potions[j]</code> - сила <code>j</code>-го зелья.</p><p>Дано также целое число <code>success</code>. Пару из заклинания и зелья называют <strong>успешной</strong>, если <strong>произведение</strong> их сил <strong>не меньше</strong> <code>success</code>.</p><p>Верните Array <code>pairs</code> длины <code>n</code>, где <code>pairs[i]</code> - число <strong>зелий</strong>, которые дают успешную пару с <code>i</code>-м заклинанием.</p><p>Метод класса <code>Solution</code>: <code>def successfulPairs(self, spells: list[int], potions: list[int], success: int) -&gt; list[int]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> spells = [5,1,3], potions = [1,2,3,4,5], success = 7
<strong>Выход:</strong> [4,0,3]
<strong>Пояснение:</strong>
- заклинание 0: 5 * [1,2,3,4,5] = [5,<strong>10</strong>,<strong>15</strong>,<strong>20</strong>,<strong>25</strong>]. Успешных пар 4.
- заклинание 1: 1 * [1,2,3,4,5] = [1,2,3,4,5]. Успешных пар 0.
- заклинание 2: 3 * [1,2,3,4,5] = [3,6,<strong>9</strong>,<strong>12</strong>,<strong>15</strong>]. Успешных пар 3.
Поэтому ответ [4,0,3].</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> spells = [3,1,2], potions = [8,5,8], success = 16
<strong>Выход:</strong> [2,0,2]
<strong>Пояснение:</strong>
- заклинание 0: 3 * [8,5,8] = [<strong>24</strong>,15,<strong>24</strong>]. Успешных пар 2.
- заклинание 1: 1 * [8,5,8] = [8,5,8]. Успешных пар 0.
- заклинание 2: 2 * [8,5,8] = [<strong>16</strong>,10,<strong>16</strong>]. Успешных пар 2.
Поэтому ответ [2,0,2].</pre><h3>Ограничения</h3><ul><li><code>n == spells.length</code></li><li><code>m == potions.length</code></li><li><code>1 ≤ n, m ≤ 10<sup>5</sup></code></li><li><code>1 ≤ spells[i], potions[i] ≤ 10<sup>5</sup></code></li><li><code>1 ≤ success ≤ 10<sup>10</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "successfulPairs"}, Cases: []sets.TestCase{
						{Input: "[[5,1,3],[1,2,3,4,5],7]", Expected: text("[4,0,3]")},
						{Input: "[[3,1,2],[8,5,8],16]", Expected: text("[2,0,2]")},
						{Input: "[[1],[1],1]", Expected: text("[1]")},
						{Input: "[[100000],[100000],10000000000]", Expected: text("[1]")},
						{Input: "[[2,7],[3,1,4,2],15]", Expected: text("[0,2]")},
						{Input: "[[4,4],[2,2,3,1],8]", Expected: text("[3,3]")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/474/B",
					Title:      "Worms",
					Difficulty: "medium",
					Topics:     []string{"Array", "Binary Search", "Prefix Sum"},
					Note:       "Codeforces, рейтинг 1200. Prefix Sum p[i] - номер последнего червяка в горке i + 1; ответ для q - bisect_left(p, q) + 1. n - число горок, m - число вопросов. Время: O(n + m log n). Память: O(n) для p, ответ O(m). Проход по горкам для каждого вопроса даёт O(nm) = 10¹⁰. Array владельцев каждого номера даёт O(S + m) времени и O(S) памяти, S ≤ 10⁶. Для a = [2, 7, 3, 4, 9] и q = 11: p = [2, 9, 12, 16, 25], ответ 3.",
					Statement: `<p>У Крота время обеда. Его друг Сурок приготовил ему к обеду интересную игру.</p><p>Сурок принёс Кроту <code>n</code> горок червяков, которые стоят по порядку. В <code>i</code>-й горке <code>a<sub>i</sub></code> червяков. Сурок пометил всех червяков подряд идущими целыми числами: червяки первой горки получили номера от <code>1</code> до <code>a<sub>1</sub></code>, червяки второй горки получили номера от <code>a<sub>1</sub> + 1</code> до <code>a<sub>1</sub> + a<sub>2</sub></code>, и так далее. Пример поможет понять это лучше.</p><p>Крот не может съесть всех червяков (Сурок принёс очень много), и, как все знают, Крот слепой. Поэтому Сурок называет ему номера самых сочных червяков. Сурок отдаст Кроту червяка, только если Крот правильно скажет, в какой горке лежит этот червяк.</p><p>Бедный Крот просит вашей помощи. Для всех сочных червяков, которых назвал Сурок, скажите Кроту правильные ответы.</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>n</code> (<code>1 ≤ n ≤ 10<sup>5</sup></code>): количество горок.</p><p>Во второй строке записаны <code>n</code> целых чисел <code>a<sub>1</sub>, a<sub>2</sub>, ..., a<sub>n</sub></code> (<code>1 ≤ a<sub>i</sub> ≤ 10<sup>3</sup></code>, <code>a<sub>1</sub> + a<sub>2</sub> + ... + a<sub>n</sub> ≤ 10<sup>6</sup></code>), где <code>a<sub>i</sub></code> - количество червяков в <code>i</code>-й горке.</p><p>В третьей строке записано одно целое число <code>m</code> (<code>1 ≤ m ≤ 10<sup>5</sup></code>): количество сочных червяков, которых назвал Сурок.</p><p>В четвёртой строке записаны <code>m</code> целых чисел <code>q<sub>1</sub>, q<sub>2</sub>, ..., q<sub>m</sub></code> (<code>1 ≤ q<sub>i</sub> ≤ a<sub>1</sub> + a<sub>2</sub> + ... + a<sub>n</sub></code>): номера сочных червяков.</p><h3>Выходные данные</h3><p>Выведите <code>m</code> строк. В <code>i</code>-й строке должно быть одно целое число: номер горки, в которой лежит червяк с номером <code>q<sub>i</sub></code>.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>5
2 7 3 4 9
3
1 25 11</pre><p><strong>Выходные данные:</strong></p><pre>1
5
3</pre><h3>Примечание</h3><p>Для входных данных примера:</p><ul><li>червяки с номерами из <code>[1, 2]</code> лежат в первой горке;</li><li>червяки с номерами из <code>[3, 9]</code> лежат во второй горке;</li><li>червяки с номерами из <code>[10, 12]</code> лежат в третьей горке;</li><li>червяки с номерами из <code>[13, 16]</code> лежат в четвёртой горке;</li><li>червяки с номерами из <code>[17, 25]</code> лежат в пятой горке.</li></ul><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 10<sup>5</sup></code></li><li><code>1 ≤ a<sub>i</sub> ≤ 10<sup>3</sup></code></li><li><code>a<sub>1</sub> + a<sub>2</sub> + ... + a<sub>n</sub> ≤ 10<sup>6</sup></code></li><li><code>1 ≤ m ≤ 10<sup>5</sup></code></li><li><code>1 ≤ q<sub>i</sub> ≤ a<sub>1</sub> + a<sub>2</sub> + ... + a<sub>n</sub></code></li><li>Ограничение по времени: 1 секунда на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "5\n2 7 3 4 9\n3\n1 25 11\n", Expected: text("1\n5\n3\n")},
						{Input: "1\n1\n1\n1\n", Expected: text("1\n")},
						{Input: "3\n1 1 1\n3\n3 2 1\n", Expected: text("3\n2\n1\n")},
						{Input: "4\n3 1 4 2\n5\n3 4 5 8 9\n", Expected: text("1\n2\n3\n3\n4\n")},
						{Input: "2\n1000 1000\n3\n1000 1001 2000\n", Expected: text("1\n2\n2\n")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "bs2-lower", Type: "short", Prompt: "a = [2, 4, 4, 4, 7, 9]. Чему равен Lower Bound(4), то есть bisect_left(a, 4)?", Answer: "1", Explain: "Первый индекс с a[i] ≥ 4 - это 1. bisect_right(a, 4) = 4: первый индекс с a[i] > 4. Четвёрок 4 − 1 = 3. Ответ 3 - это индекс последней четвёрки, а не Lower Bound."},
				{ID: "bs2-count", Type: "short", Prompt: "a = [1, 2, 2, 5, 5, 5, 8, 9]. Сколько чисел лежит в отрезке от 2 до 5 включительно? Посчитайте bisect_right(a, 5) − bisect_left(a, 2).", Answer: "5", Explain: "bisect_right(a, 5) = 6, bisect_left(a, 2) = 1, разность 5: числа 2, 2, 5, 5, 5. С bisect_left(a, 5) = 3 получилось бы 2: пятёрки потерялись бы. Два поиска стоят O(log n) времени и O(1) памяти."},
				{ID: "bs2-side", Type: "choice", Prompt: "Codeforces 706B. Цены x = [3, 6, 8, 10, 11] упорядочены, у Василия m = 10 монет. Какой вызов даёт число магазинов, где ему хватит денег?", Options: []string{"bisect_left(x, m)", "bisect_right(x, m)", "bisect_right(x, m) − 1"}, Answer: 1, Explain: "Нужны цены ≤ 10. Первая позиция с ценой > 10 - это bisect_right = 4, и она равна числу подходящих цен. bisect_left = 3 не считает магазин с ценой ровно 10. bisect_right − 1 = 3 - индекс последней подходящей цены, а не их число."},
				{ID: "bs2-loop", Type: "code", Prompt: "Какая строка стоит на месте ___ в Lower Bound?", Language: "python", Code: `def lower_bound(a, x):
    lo, hi = 0, len(a)
    while lo < hi:
        mid = (lo + hi) // 2
        if a[mid] >= x:
            ___
        else:
            lo = mid + 1
    return lo`, Options: []string{"hi = mid - 1", "lo = mid", "hi = mid"}, Answer: 2, Explain: "a[mid] ≥ x, значит mid может быть ответом, и его нельзя выбросить: hi = mid. hi = mid − 1 теряет ответ: для a = [1, 3], x = 3 функция вернёт 0 вместо 1. lo = mid двигает не ту границу: на том же входе lo = 1 повторяется, и цикл не закончится."},
				{ID: "bs2-steps", Type: "short", Prompt: "Lower Bound в Array из 1 000 000 чисел, шаблон while lo < hi с hi = len(a). Сколько проверок середины он сделает в худшем случае?", Answer: "20", Explain: "Длина hi − lo уменьшается не меньше чем вдвое: 1 000 000, 500 000, ..., 3, 1, 0. Это 20 проверок: 2¹⁹ = 524 288 ≤ 1 000 000 < 2²⁰. Ответ 19 = ⌊log₂ 10⁶⌋ забывает последнюю проверку при длине 1. Linear Scan проверил бы до миллиона чисел."},
				{ID: "bs2-potions", Type: "short", Prompt: "LeetCode 2300. potions = [1, 2, 3, 4, 5] упорядочены, success = 7, spell = 2. Сколько зелий дают успешную пару?", Answer: "2", Explain: "2 · p ≥ 7 при p ≥ ⌈7 / 2⌉ = 4. bisect_left(potions, 4) = 3, подходят 5 − 3 = 2 зелья: 4 и 5. Деление 7 // 2 = 3 без округления вверх дало бы 3 зелья, но 2 · 3 = 6 < 7."},
				{ID: "bs2-worms", Type: "short", Prompt: "Codeforces 474B. Горки a = [3, 1, 4, 2]. В какой горке лежит червяк с номером 8?", Answer: "3", Explain: "Prefix Sum p = [3, 4, 8, 10]: горка 3 содержит номера 5..8. Первый индекс с p[i] ≥ 8 равен 2, ответ 2 + 1 = 3. bisect_right(p, 8) = 3 дал бы горку 4: номер, равный концу горки, ушёл бы в следующую горку."},
				{ID: "bs2-cost", Type: "choice", Prompt: "n = 1 000 цен и q = 1 000 000 дней. Что даёт модель числа сравнений?", Options: []string{"Sorting и bisect_right: около 10⁷ сравнений вместо 10⁹", "Linear Scan каждый день: около 10⁶ сравнений всего", "Sorting и bisect_right: около 10⁴ сравнений"}, Answer: 0, Explain: "Linear Scan сравнивает все n цен в каждый из q дней: 1000 · 10⁶ = 10⁹. Sorting делает около n · 10 = 10⁴ сравнений один раз, затем ⌊log₂ 1000⌋ + 1 = 10 проверок на день: 10 · 10⁶ = 10⁷. Вариант 10⁴ забывает запросы: при q намного больше n почти вся работа уходит на них. Это модель, а не секунды. Дополнительная память: копия sorted(x) из n ссылок и q ответов."},
			},
			Links: []sets.Link{
				{Title: "Левый и правый Binary Search", URL: "https://acm.khpnets.info/wiki/%D0%91%D0%B8%D0%BD%D0%B0%D1%80%D0%BD%D1%8B%D0%B9_%D0%BF%D0%BE%D0%B8%D1%81%D0%BA", Source: "Олимпиадное программирование в УлГТУ"},
				{Title: "Первое и последнее вхождение числа", URL: "https://habr.com/ru/articles/1017108/", Source: "Хабр"},
				{Title: "Левая и правая граница: разбор задач", URL: "https://habr.com/ru/articles/443244/", Source: "Хабр"},
				{Title: "bisect: позиция в упорядоченном list", URL: "https://docs.python.org/3/library/bisect.html", Source: "Python docs (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "bs2-lower",
					Title: "Lower Bound(3) в [1, 3, 3, 3, 5, 8]",
					Frames: []sets.FigureFrame{
						{Caption: "lo = 0, hi = 6: hi стоит после последней клетки. mid = 3: a[3] = 3 ≥ 3, поэтому индексы 3, 4, 5 дают «да». Ответ не правее 3: hi = 3.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 3, 3, 5, 8}, Marked: []int{3}, Pointers: &sets.Pointers{{Name: "lo", Cell: 0}, {Name: "mid", Cell: 3}}, Window: &[2]int{0, 5}},
							{Label: "a[i] ≥ 3", Cells: []any{"?", "?", "?", "да", "да", "да"}},
						}},
						{Caption: "mid = 1: a[1] = 3 ≥ 3, hi = 1.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 3, 3, 5, 8}, Marked: []int{1}, Pointers: &sets.Pointers{{Name: "lo", Cell: 0}, {Name: "mid", Cell: 1}, {Name: "hi", Cell: 3}}, Window: &[2]int{0, 3}},
							{Label: "a[i] ≥ 3", Cells: []any{"?", "да", "да", "да", "да", "да"}},
						}},
						{Caption: "mid = 0: a[0] = 1 < 3, lo = 1.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 3, 3, 5, 8}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "lo", Cell: 0}, {Name: "mid", Cell: 0}, {Name: "hi", Cell: 1}}, Window: &[2]int{0, 1}},
							{Label: "a[i] ≥ 3", Cells: []any{"нет", "да", "да", "да", "да", "да"}},
						}},
						{Caption: "lo = hi = 1: Lower Bound(3) = 1 после 3 проверок. Upper Bound(3) = 4, поэтому троек 4 − 1 = 3.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 3, 3, 5, 8}, Pointers: &sets.Pointers{{Name: "ответ", Cell: 1}}},
							{Label: "a[i] ≥ 3", Cells: []any{"нет", "да", "да", "да", "да", "да"}},
						}},
					},
				},
				{
					ID:    "bs2-queries",
					Title: "Interesting drink: bisect_right по ценам",
					Frames: []sets.FigureFrame{
						{Caption: "m = 1: ни одна цена не ≤ 1. bisect_right = 0, ответ 0.", Rows: []sets.FigureRow{
							{Label: "x", Cells: []any{3, 6, 8, 10, 11}, Pointers: &sets.Pointers{{Name: "bisect_right", Cell: 0}}},
							{Label: "x ≤ 1", Cells: []any{"нет", "нет", "нет", "нет", "нет"}},
						}},
						{Caption: "m = 10: подходят 3, 6, 8 и 10. bisect_right = 4. bisect_left дал бы 3 и потерял магазин с ценой 10.", Rows: []sets.FigureRow{
							{Label: "x", Cells: []any{3, 6, 8, 10, 11}, Marked: []int{0, 1, 2, 3}, Pointers: &sets.Pointers{{Name: "bisect_right", Cell: 4}}},
							{Label: "x ≤ 10", Cells: []any{"да", "да", "да", "да", "нет"}},
						}},
						{Caption: "m = 3: подходит одна цена, bisect_right = 1.", Rows: []sets.FigureRow{
							{Label: "x", Cells: []any{3, 6, 8, 10, 11}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "bisect_right", Cell: 1}}},
							{Label: "x ≤ 3", Cells: []any{"да", "нет", "нет", "нет", "нет"}},
						}},
						{Caption: "m = 11: подходят все 5 цен. bisect_right = 5, это позиция после конца ряда. Ответы 0, 4, 1, 5, как в примере.", Rows: []sets.FigureRow{
							{Label: "x", Cells: []any{3, 6, 8, 10, 11}, Marked: []int{0, 1, 2, 3, 4}},
							{Label: "x ≤ 11", Cells: []any{"да", "да", "да", "да", "да"}},
						}},
					},
				},
				{
					ID:    "bs2-worms",
					Title: "Worms: Prefix Sum и bisect_left",
					Frames: []sets.FigureFrame{
						{Caption: "p[i] - номер последнего червяка в горках 1..i + 1. Prefix Sum строго растёт, потому что каждое a[i] ≥ 1.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 7, 3, 4, 9}},
							{Label: "p", Cells: []any{2, 9, 12, 16, 25}},
						}},
						{Caption: "q = 11: первое p[i] ≥ 11 - это p[2] = 12. Червяк 11 в горке 3.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 7, 3, 4, 9}, Marked: []int{2}},
							{Label: "p", Cells: []any{2, 9, 12, 16, 25}, Marked: []int{2}, Pointers: &sets.Pointers{{Name: "горка 3", Cell: 2}}},
						}},
						{Caption: "q = 9: p[1] = 9 ≥ 9, червяк 9 последний в горке 2. bisect_right дал бы горку 3.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 7, 3, 4, 9}, Marked: []int{1}},
							{Label: "p", Cells: []any{2, 9, 12, 16, 25}, Marked: []int{1}, Pointers: &sets.Pointers{{Name: "горка 2", Cell: 1}}},
						}},
						{Caption: "q = 25: p[4] = 25, горка 5. Ответы 1, 5, 3 из примера получаются так же.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 7, 3, 4, 9}, Marked: []int{4}},
							{Label: "p", Cells: []any{2, 9, 12, 16, 25}, Marked: []int{4}, Pointers: &sets.Pointers{{Name: "горка 5", Cell: 4}}},
						}},
					},
				},
			},
		},
		{
			Title:   "Binary Search по ответу",
			Summary: "Иногда ответ - число, и проверить «подходит ли x» легче, чем найти x сразу. Если проверка монотонная, Binary Search находит наименьший или наибольший подходящий x. Вы выбираете границы ответа, пишете проверку за O(n) и считаете общее время O(n log C), где C - число кандидатов.",
			Lesson:  lesson3,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/arranging-coins/",
					Number:     number(441),
					Title:      "Arranging Coins",
					Difficulty: "easy",
					Topics:     []string{"Math", "Binary Search"},
					Note:       "Binary Search по ответу: наибольшее k с k(k + 1) / 2 ≤ n на отрезке [1, n]. mid = (lo + hi + 1) // 2, потому что ищем наибольшее подходящее k. Время: O(log n), не больше 31 проверки при n ≤ 2³¹ − 1. Память: O(1) дополнительно, ответ O(1). Linear Scan k = 1, 2, ... делает до 65 536 шагов: O(√n). Для n = 8: 3 · 4 / 2 = 6 ≤ 8 < 10 = 4 · 5 / 2, ответ 3.",
					Statement: `<p>У вас есть <code>n</code> монет, и вы хотите построить из них лестницу. Лестница состоит из <code>k</code> рядов, в <code>i</code>-м ряду ровно <code>i</code> монет. Последний ряд лестницы <strong>может быть</strong> неполным.</p><p>Дано целое число <code>n</code>. Верните <em>число <strong>полных рядов</strong> лестницы, которую вы построите</em>.</p><p>Метод класса <code>Solution</code>: <code>def arrangeCoins(self, n: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> n = 5
<strong>Выход:</strong> 2
<strong>Пояснение:</strong> Ряды 1 и 2 полные, в них 1 + 2 = 3 монеты. Третий ряд неполный: в нём 2 монеты из 3. Поэтому ответ 2.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> n = 8
<strong>Выход:</strong> 3
<strong>Пояснение:</strong> Ряды 1, 2 и 3 полные, в них 1 + 2 + 3 = 6 монет. Четвёртый ряд неполный: в нём 2 монеты из 4. Поэтому ответ 3.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 2<sup>31</sup> - 1</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "arrangeCoins"}, Cases: []sets.TestCase{
						{Input: "[5]", Expected: text("2")},
						{Input: "[8]", Expected: text("3")},
						{Input: "[1]", Expected: text("1")},
						{Input: "[3]", Expected: text("2")},
						{Input: "[2147483647]", Expected: text("65535")},
						{Input: "[2147450879]", Expected: text("65534")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1873/E",
					Title:      "Building an Aquarium",
					Difficulty: "easy",
					Topics:     []string{"Array", "Binary Search"},
					Note:       "Codeforces, рейтинг 1100. Проверка вода(h) = сумма max(0, h − a[i]) ≤ x монотонная: при большем h воды не меньше. Ищите наибольшее h на [1, M + x], M = max(a). Время: O(n log(M + x)), не больше 31 проверки по n столбцов на набор, по всем наборам около 6,2 · 10⁶ шагов. Память: O(1) сверх входа. Linear Scan по h = 1, 2, ... даёт до 2 · 10⁹ проверок. Для a = [3, 1, 2, 4, 6, 2, 5], x = 9: вода(4) = 8, вода(5) = 13, ответ 4.",
					Statement: `<p>Вы любите рыбок, поэтому решили построить аквариум. У вас есть коралл из <code>n</code> столбцов, высота <code>i</code>-го столбца <code>a<sub>i</sub></code> единиц. Вокруг коралла вы строите бак так:</p><ul><li>Выберите целое число <code>h ≥ 1</code>: высоту бака. Постройте стенки высоты <code>h</code> по обеим сторонам бака.</li><li>Затем налейте в бак воду так, чтобы высота каждого столбца стала <code>h</code>. Если коралл выше <code>h</code>, в этот столбец воду не наливают.</li></ul><p>Например, при <code>a = [3, 1, 2, 4, 6, 2, 5]</code> и высоте <code>h = 4</code> нужно всего <code>8</code> единиц воды: в столбцы высоты 3, 1, 2 и 2 нальют 1, 3, 2 и 2 единицы, а в столбцы высоты 4, 6 и 5 воду не нальют.</p><p>Воды можно использовать не больше <code>x</code> единиц, а бак нужен как можно выше. Какую наибольшую высоту <code>h</code> можно выбрать?</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>): количество наборов входных данных.</p><p>В первой строке каждого набора записаны два целых положительных числа <code>n</code> и <code>x</code> (<code>1 ≤ n ≤ 2 · 10<sup>5</sup></code>; <code>1 ≤ x ≤ 10<sup>9</sup></code>): количество столбцов коралла и наибольшее количество воды, которое можно использовать.</p><p>Во второй строке каждого набора записаны через пробел <code>n</code> целых чисел <code>a<sub>i</sub></code> (<code>1 ≤ a<sub>i</sub> ≤ 10<sup>9</sup></code>): высоты столбцов коралла.</p><p>Сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого набора выведите одно целое положительное число <code>h</code> (<code>h ≥ 1</code>): наибольшую высоту бака, при которой для заполнения бака нужно не больше <code>x</code> единиц воды.</p><p>Можно доказать, что при этих ограничениях такое <code>h</code> всегда существует.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>5
7 9
3 1 2 4 6 2 5
3 10
1 1 1
4 1
1 4 3 4
6 1984
2 6 5 9 1 8
1 1000000000
1</pre><p><strong>Выходные данные:</strong></p><pre>4
4
2
335
1000000001</pre><h3>Примечание</h3><p>Первый набор описан в условии. При <code>h = 4</code> нужно <code>8</code> единиц воды, а если поднять <code>h</code> до <code>5</code>, нужно <code>13</code> единиц, это больше <code>x = 9</code>. Поэтому <code>h = 4</code> оптимально.</p><p>Во втором наборе можно выбрать <code>h = 4</code> и добавить по <code>3</code> единицы в каждый столбец, всего <code>9</code> единиц воды. Можно показать, что это оптимально.</p><p>В третьем наборе можно выбрать <code>h = 2</code> и использовать всю воду, поэтому это оптимально.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>1 ≤ n ≤ 2 · 10<sup>5</sup></code></li><li><code>1 ≤ x ≤ 10<sup>9</sup></code></li><li><code>1 ≤ a<sub>i</sub> ≤ 10<sup>9</sup></code></li><li>сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "5\n7 9\n3 1 2 4 6 2 5\n3 10\n1 1 1\n4 1\n1 4 3 4\n6 1984\n2 6 5 9 1 8\n1 1000000000\n1\n", Expected: text("4\n4\n2\n335\n1000000001\n")},
						{Input: "1\n1 1\n5\n", Expected: text("6\n")},
						{Input: "1\n3 1\n2 2 2\n", Expected: text("2\n")},
						{Input: "1\n2 1000000000\n1000000000 1000000000\n", Expected: text("1500000000\n")},
						{Input: "2\n4 6\n5 1 1 5\n3 3\n1 2 3\n", Expected: text("4\n3\n")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/koko-eating-bananas/",
					Number:     number(875),
					Title:      "Koko Eating Bananas",
					Difficulty: "medium",
					Topics:     []string{"Array", "Binary Search"},
					Note:       "Наименьшая скорость k, при которой часы(k) = сумма ⌈p / k⌉ ≤ h. ⌈p / k⌉ в Python - (p + k - 1) // k. Проверка монотонная: при большей скорости часов не больше. Ищите k на [1, M], M = max(piles), n = len(piles). Время: O(n log M), не больше 30 проверок по n делений. Память: O(1) дополнительно. Linear Scan по k = 1, 2, ... даёт O(nM), до 10¹³ делений. Для [3, 6, 7, 11], h = 8: часы(3) = 10, часы(4) = 8, ответ 4.",
					Statement: `<p>Koko любит бананы. Перед ней <code>n</code> горок бананов, в <code>i</code>-й горке <code>piles[i]</code> бананов. Охранники ушли и вернутся через <code>h</code> часов.</p><p>Koko может выбрать скорость <code>k</code>: сколько бананов в час она ест. Каждый час она выбирает одну горку и съедает из неё <code>k</code> бананов. Если в горке меньше <code>k</code> бананов, она съедает их все и в этот час больше не ест.</p><p>Koko любит есть медленно, но хочет съесть все бананы до возвращения охранников.</p><p>Верните <strong>наименьшее</strong> целое <code>k</code>, при котором она съест все бананы <strong>не больше чем за</strong> <code>h</code> часов.</p><p>Метод класса <code>Solution</code>: <code>def minEatingSpeed(self, piles: list[int], h: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> piles = [3,6,7,11], h = 8
<strong>Выход:</strong> 4
<strong>Пояснение:</strong> При скорости 4 банана в час Koko ест горки за 1, 2, 2 и 3 часа, всего 8 часов.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> piles = [30,11,23,4,20], h = 5
<strong>Выход:</strong> 30
<strong>Пояснение:</strong> Каждую горку нужно съесть ровно за 1 час, поэтому наименьшая скорость 30 бананов в час.</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> piles = [30,11,23,4,20], h = 6
<strong>Выход:</strong> 23
<strong>Пояснение:</strong> При скорости 23 банана в час Koko ест первую горку за 2 часа, а каждую другую за 1 час, всего 6 часов.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ piles.length ≤ 10<sup>4</sup></code></li><li><code>piles.length ≤ h ≤ 10<sup>9</sup></code></li><li><code>1 ≤ piles[i] ≤ 10<sup>9</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "minEatingSpeed"}, Cases: []sets.TestCase{
						{Input: "[[3,6,7,11],8]", Expected: text("4")},
						{Input: "[[30,11,23,4,20],5]", Expected: text("30")},
						{Input: "[[30,11,23,4,20],6]", Expected: text("23")},
						{Input: "[[1],1]", Expected: text("1")},
						{Input: "[[1000000000],2]", Expected: text("500000000")},
						{Input: "[[2,2,2],1000000000]", Expected: text("1")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/capacity-to-ship-packages-within-d-days/",
					Number:     number(1011),
					Title:      "Capacity To Ship Packages Within D Days",
					Difficulty: "medium",
					Topics:     []string{"Array", "Binary Search", "Greedy"},
					Note:       "Наименьшая ёмкость c, при которой Greedy-погрузка по порядку укладывается в days дней. Ищите c на [max(weights), sum(weights)]. n = len(weights), S = sum(weights). Время: O(n log S), не больше 25 проверок по n весов при S ≤ 2,5 · 10⁷. Память: O(1) дополнительно. Linear Scan по c от max(weights) даёт O(nS). Для [3, 2, 2, 4, 1, 4], days = 3: ёмкость 5 даёт 4 дня, 6 даёт 3 дня, ответ 6.",
					Statement: `<p>На ленте конвейера лежат грузы. Их нужно перевезти из одного порта в другой не больше чем за <code>days</code> дней.</p><p>Вес <code>i</code>-го груза на ленте равен <code>weights[i]</code>. Каждый день на корабль грузят грузы с ленты в порядке <code>weights</code>. Нельзя грузить больше, чем позволяет грузоподъёмность корабля.</p><p>Верните наименьшую грузоподъёмность корабля, при которой все грузы с ленты будут перевезены не больше чем за <code>days</code> дней.</p><p>Метод класса <code>Solution</code>: <code>def shipWithinDays(self, weights: list[int], days: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> weights = [1,2,3,4,5,6,7,8,9,10], days = 5
<strong>Выход:</strong> 15
<strong>Пояснение:</strong> Грузоподъёмность 15 - наименьшая, при которой все грузы можно перевезти за 5 дней так:
1-й день: 1, 2, 3, 4, 5
2-й день: 6, 7
3-й день: 8
4-й день: 9
5-й день: 10

Грузы нужно везти в заданном порядке. Поэтому корабль грузоподъёмности 14 и деление на части (2, 3, 4, 5), (1, 6, 7), (8), (9), (10) не подходят.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> weights = [3,2,2,4,1,4], days = 3
<strong>Выход:</strong> 6
<strong>Пояснение:</strong> Грузоподъёмность 6 - наименьшая, при которой все грузы можно перевезти за 3 дня так:
1-й день: 3, 2
2-й день: 2, 4
3-й день: 1, 4</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> weights = [1,2,3,1,1], days = 4
<strong>Выход:</strong> 3
<strong>Пояснение:</strong>
1-й день: 1
2-й день: 2
3-й день: 3
4-й день: 1, 1</pre><h3>Ограничения</h3><ul><li><code>1 ≤ days ≤ weights.length ≤ 5 · 10<sup>4</sup></code></li><li><code>1 ≤ weights[i] ≤ 500</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "shipWithinDays"}, Cases: []sets.TestCase{
						{Input: "[[1,2,3,4,5,6,7,8,9,10],5]", Expected: text("15")},
						{Input: "[[3,2,2,4,1,4],3]", Expected: text("6")},
						{Input: "[[1,2,3,1,1],4]", Expected: text("3")},
						{Input: "[[7],1]", Expected: text("7")},
						{Input: "[[1,2,3,4],4]", Expected: text("4")},
						{Input: "[[5,5,5,5],1]", Expected: text("20")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1201/C",
					Title:      "Maximum Median",
					Difficulty: "medium",
					Topics:     []string{"Array", "Binary Search", "Greedy", "Math", "Sorting"},
					Note:       "Codeforces, рейтинг 1400. Упорядочите a. Медиана m достижима, если сумма max(0, m − a[i]) по i от n // 2 до n − 1 не больше k. Ищите наибольшее m на [a[n // 2], a[n // 2] + k]. Время: O(n log n + n log k), не больше 30 проверок по (n + 1) / 2 чисел. Память: O(n) для упорядоченной копии, ответ O(1). Подъём медианы на 1 за шаг даёт до k = 10⁹ шагов. Для [1, 2, 1, 1, 1], k = 5: для 3 нужно 5 операций, для 4 нужно 8, ответ 3.",
					Statement: `<p>Дан Array <code>a</code> из <code>n</code> целых чисел, <code>n</code> нечётное. С ним можно делать такую операцию:</p><ul><li>Выбрать один элемент Array (например, <code>a<sub>i</sub></code>) и увеличить его на <code>1</code>, то есть заменить его на <code>a<sub>i</sub> + 1</code>.</li></ul><p>Нужно сделать медиану Array как можно больше, если можно сделать не больше <code>k</code> операций.</p><p>Медиана Array нечётной длины - это средний элемент после упорядочивания Array по неубыванию. Например, медиана Array <code>[1, 5, 2, 3, 5]</code> равна <code>3</code>.</p><h3>Входные данные</h3><p>В первой строке записаны два целых числа <code>n</code> и <code>k</code> (<code>1 ≤ n ≤ 2 · 10<sup>5</sup></code>, <code>n</code> нечётное, <code>1 ≤ k ≤ 10<sup>9</sup></code>): количество элементов Array и наибольшее количество операций.</p><p>Во второй строке записаны <code>n</code> целых чисел <code>a<sub>1</sub>, a<sub>2</sub>, ..., a<sub>n</sub></code> (<code>1 ≤ a<sub>i</sub> ≤ 10<sup>9</sup></code>).</p><h3>Выходные данные</h3><p>Выведите одно целое число: наибольшую медиану, которую можно получить после операций.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>3 2
1 3 5</pre><p><strong>Выходные данные:</strong></p><pre>5</pre><h3>Пример 2</h3><p><strong>Входные данные:</strong></p><pre>5 5
1 2 1 1 1</pre><p><strong>Выходные данные:</strong></p><pre>3</pre><h3>Пример 3</h3><p><strong>Входные данные:</strong></p><pre>7 7
4 1 2 4 3 4 4</pre><p><strong>Выходные данные:</strong></p><pre>5</pre><h3>Примечание</h3><p>В первом примере можно дважды увеличить второй элемент. Тогда Array станет <code>[1, 5, 5]</code>, и его медиана равна <code>5</code>.</p><p>Во втором примере лучше всего увеличить второе число, а затем третье и пятое. Так ответ равен <code>3</code>.</p><p>В третьем примере можно сделать четыре операции: увеличить первый, четвёртый, шестой и седьмой элементы. Тогда Array станет <code>[5, 1, 2, 5, 3, 5, 5]</code>, и медиана будет <code>5</code>.</p><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 2 · 10<sup>5</sup></code>, <code>n</code> нечётное</li><li><code>1 ≤ k ≤ 10<sup>9</sup></code></li><li><code>1 ≤ a<sub>i</sub> ≤ 10<sup>9</sup></code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "3 2\n1 3 5\n", Expected: text("5\n")},
						{Input: "5 5\n1 2 1 1 1\n", Expected: text("3\n")},
						{Input: "7 7\n4 1 2 4 3 4 4\n", Expected: text("5\n")},
						{Input: "1 1000000000\n1000000000\n", Expected: text("2000000000\n")},
						{Input: "3 1\n5 5 5\n", Expected: text("5\n")},
						{Input: "5 3\n1 1 9 9 9\n", Expected: text("10\n")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "bs3-mono", Type: "choice", Prompt: "LeetCode 875: piles = [5, 9, 12, 14], h = 13. Binary Search по скорости k ищет первое «да» проверки. Какая проверка для этого подходит?", Options: []string{"hours(k) == h", "hours(k) <= h", "hours(k) % 2 == 1"}, Answer: 1, Explain: "hours(3) = 2 + 3 + 4 + 5 = 14, hours(4) = 2 + 3 + 3 + 4 = 12: ни одно k не даёт ровно 13 часов, и проверка == везде даёт «нет». При большей скорости часов не больше, поэтому из hours(k) ≤ h следует hours(k + 1) ≤ h: первое «да» при k = 4. Нечётность часов скачет: hours(2) = 21, hours(3) = 14, hours(5) = 9."},
				{ID: "bs3-hours", Type: "short", Prompt: "piles = [4, 9, 10, 13]. Сколько часов нужно Koko при скорости k = 4?", Answer: "11", Explain: "⌈4 / 4⌉ + ⌈9 / 4⌉ + ⌈10 / 4⌉ + ⌈13 / 4⌉ = 1 + 3 + 3 + 4 = 11. Деление p // 4 без округления вверх даёт 1 + 2 + 2 + 3 = 8 и теряет неполные часы. Одна проверка проходит n горок: O(n) времени, O(1) памяти."},
				{ID: "bs3-checks", Type: "short", Prompt: "LeetCode 875, max(piles) = 10⁹. Сколько проверок hours(k) в худшем случае сделает поиск наименьшего k на отрезке [1, 10⁹]?", Answer: "30", Explain: "Каждая проверка оставляет не больше половины кандидатов, округлённой вверх: 10⁹, 5 · 10⁸, ..., 2, 1. 2²⁹ < 10⁹ ≤ 2³⁰, поэтому нужно 30 проверок. Ответ 29 = ⌊log₂ 10⁹⌋ округляет вниз: после 29 проверок могут остаться 2 кандидата. Каждая проходит n горок, для n = 10⁴ это 3 · 10⁵ делений. Linear Scan по k = 1, 2, ... мог бы сделать 10⁹ проверок."},
				{ID: "bs3-mid", Type: "code", Prompt: "Поиск наибольшей высоты h. Какое выражение стоит на месте ___?", Language: "python", Code: `lo, hi = 1, max(a) + w
while lo < hi:
    mid = ___
    if water(a, mid) <= w:
        lo = mid
    else:
        hi = mid - 1
return lo`, Options: []string{"(lo + hi) // 2", "(lo + hi) // 2 - 1", "(lo + hi + 1) // 2"}, Answer: 2, Explain: "При lo = 4, hi = 5 середина (4 + 5) // 2 = 4. Если h = 4 подходит, lo = 4 снова, и цикл не закончится. С + 1 середина 5, и отрезок уменьшается при любом ответе проверки. Вариант − 1 даёт 3 < lo и выходит за отрезок."},
				{ID: "bs3-hi", Type: "choice", Prompt: "Codeforces 1873E: a = [5, 7], x = 10. Поиск наибольшего h взял hi = max(a) = 7. Что пойдёт не так?", Options: []string{"Поиск вернёт 7, а ответ 11", "Проверка water станет медленной", "Цикл while lo < hi не закончится"}, Answer: 0, Explain: "вода(11) = 6 + 4 = 10 ≤ 10, вода(12) = 7 + 5 = 12 > 10: ответ 11. На [1, 7] подходит каждое h, и поиск вернёт 7. Ответ лежит в [1, max(a) + x] = [1, 17]. Узкий hi не замедляет поиск и не ломает цикл, но теряет ответ."},
				{ID: "bs3-ship", Type: "short", Prompt: "LeetCode 1011: weights = [3, 2, 2, 4, 1, 4], ёмкость 5. Сколько дней займёт Greedy-погрузка по порядку?", Answer: "4", Explain: "Дни [3, 2], [2], [4, 1], [4]. В примере 2 задачи days = 3: 4 > 3, поэтому ёмкость 5 не подходит. Ёмкость 6 даёт [3, 2], [2, 4], [1, 4]: ответ 6. Сравнение load + w >= c вместо > c закрывает день при ровно 5: [3], [2, 2], [4], [1], [4], это 5 дней. Проверка проходит Array один раз: O(n) времени, O(1) памяти."},
				{ID: "bs3-median", Type: "short", Prompt: "Codeforces 1201C: n = 7, a = [2, 8, 1, 5, 3, 9, 4]. Сколько операций нужно, чтобы медиана стала не меньше 8?", Answer: "7", Explain: "После Sorting a = [1, 2, 3, 4, 5, 8, 9], медиана a[3] = 4. Поднять до 8 нужно a[3] и всё правее: 4 + 3 + 0 + 0 = 7. Ответ 4 поднимает только медиану: тогда a = [1, 2, 3, 5, 8, 8, 9], и медиана станет 5. Числа левее медианы менять не нужно."},
				{ID: "bs3-total", Type: "choice", Prompt: "Codeforces 1201C: n = 99 999, k = 10⁹. Поиск медианы m идёт на [a[n // 2], a[n // 2] + k], и каждая проверка проходит (n + 1) / 2 = 50 000 чисел. Сколько шагов в модели сделают все проверки?", Options: []string{"около 30", "около 30 · 5 · 10⁴ = 1,5 · 10⁶", "около 10⁹ · 5 · 10⁴"}, Answer: 1, Explain: "Кандидатов k + 1 = 10⁹ + 1, и 2²⁹ < 10⁹ + 1 < 2³⁰: не больше 30 проверок. Каждая проходит 50 000 чисел: 30 · 50 000 = 1,5 · 10⁶ шагов. 30 - число проверок, а не всех шагов. Подъём медианы на 1 за шаг дал бы до 10⁹ проверок. Sorting до поиска стоит ещё около n · 17 ≈ 1,7 · 10⁶ сравнений."},
			},
			Links: []sets.Link{
				{Title: "Binary Search по ответу", URL: "https://algorithmica.org/ru/binary-search", Source: "Алгоритмика"},
				{Title: "Binary Search по ответу на задаче Hamburgers", URL: "https://brestprog.by/topics/binsearch/", Source: "brestprog"},
				{Title: "Binary Search по монотонной функции", URL: "https://usaco.guide/silver/binary-search", Source: "USACO Guide (англ.)"},
				{Title: "Binary Search по произвольному условию", URL: "https://cp-algorithms.com/num_methods/binary_search.html", Source: "cp-algorithms (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "bs3-hours",
					Title: "Koko: часы для каждой скорости",
					Frames: []sets.FigureFrame{
						{Caption: "piles = [3, 6, 7, 11]. При большей скорости часов не больше. Для h = 8 первое «да» при k = 4.", Rows: []sets.FigureRow{
							{Label: "k", Cells: []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, Marked: []int{3}},
							{Label: "часы", Cells: []any{27, 15, 10, 8, 8, 6, 5, 5, 5, 5, 4}, Marked: []int{3}},
							{Label: "часы ≤ 8", Cells: []any{"нет", "нет", "нет", "да", "да", "да", "да", "да", "да", "да", "да"}, Marked: []int{3}},
						}},
						{Caption: "Для h = 6 первое «да» при k = 6. Условие снова идёт «нет», потом только «да»: проверка монотонная.", Rows: []sets.FigureRow{
							{Label: "k", Cells: []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, Marked: []int{5}},
							{Label: "часы", Cells: []any{27, 15, 10, 8, 8, 6, 5, 5, 5, 5, 4}, Marked: []int{5}},
							{Label: "часы ≤ 6", Cells: []any{"нет", "нет", "нет", "нет", "нет", "да", "да", "да", "да", "да", "да"}, Marked: []int{5}},
						}},
					},
				},
				{
					ID:    "bs3-koko",
					Title: "Koko: Binary Search по скорости",
					Frames: []sets.FigureFrame{
						{Caption: "lo = 1, hi = 11, mid = 6: 6 часов ≤ 8, подходит. Ответ не больше 6: hi = 6.", Rows: []sets.FigureRow{
							{Label: "k", Cells: []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, Marked: []int{5}, Pointers: &sets.Pointers{{Name: "lo", Cell: 0}, {Name: "mid", Cell: 5}, {Name: "hi", Cell: 10}}, Window: &[2]int{0, 10}},
							{Label: "часы(mid)", Cells: []any{6}},
						}},
						{Caption: "mid = 3: 10 часов > 8, не подходит: lo = 4.", Rows: []sets.FigureRow{
							{Label: "k", Cells: []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, Marked: []int{2}, Dimmed: []int{6, 7, 8, 9, 10}, Pointers: &sets.Pointers{{Name: "lo", Cell: 0}, {Name: "mid", Cell: 2}, {Name: "hi", Cell: 5}}, Window: &[2]int{0, 5}},
							{Label: "часы(mid)", Cells: []any{10}},
						}},
						{Caption: "mid = 5: 8 часов ≤ 8, hi = 5.", Rows: []sets.FigureRow{
							{Label: "k", Cells: []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, Marked: []int{4}, Dimmed: []int{0, 1, 2, 6, 7, 8, 9, 10}, Pointers: &sets.Pointers{{Name: "lo", Cell: 3}, {Name: "mid", Cell: 4}, {Name: "hi", Cell: 5}}, Window: &[2]int{3, 5}},
							{Label: "часы(mid)", Cells: []any{8}},
						}},
						{Caption: "mid = 4: 8 часов ≤ 8, hi = 4.", Rows: []sets.FigureRow{
							{Label: "k", Cells: []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, Marked: []int{3}, Dimmed: []int{0, 1, 2, 5, 6, 7, 8, 9, 10}, Pointers: &sets.Pointers{{Name: "lo", Cell: 3}, {Name: "mid", Cell: 3}, {Name: "hi", Cell: 4}}, Window: &[2]int{3, 4}},
							{Label: "часы(mid)", Cells: []any{8}},
						}},
						{Caption: "lo = hi = 4: ответ 4 после 4 проверок. Для max(piles) = 10⁹ проверок было бы не больше 30.", Rows: []sets.FigureRow{
							{Label: "k", Cells: []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, Dimmed: []int{0, 1, 2, 4, 5, 6, 7, 8, 9, 10}, Pointers: &sets.Pointers{{Name: "ответ", Cell: 3}}},
							{Label: "часы(ответ)", Cells: []any{8}},
						}},
					},
				},
				{
					ID:    "bs3-aquarium",
					Title: "Аквариум: наибольшая высота h",
					Frames: []sets.FigureFrame{
						{Caption: "h = 4: вода 1 + 3 + 2 + 2 = 8 ≤ 9.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{3, 1, 2, 4, 6, 2, 5}},
							{Label: "вода при h = 4", Cells: []any{1, 3, 2, 0, 0, 2, 0}, Marked: []int{0, 1, 2, 5}},
						}},
						{Caption: "h = 5: вода 2 + 4 + 3 + 1 + 3 = 13 > 9. Чем выше h, тем больше воды.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{3, 1, 2, 4, 6, 2, 5}},
							{Label: "вода при h = 5", Cells: []any{2, 4, 3, 1, 0, 3, 0}, Marked: []int{0, 1, 2, 3, 5}},
						}},
						{Caption: "Поиск на [1, 15], 15 = max(a) + x. mid = (1 + 15 + 1) // 2 = 8: 33 > 9, hi = 7.", Rows: []sets.FigureRow{
							{Label: "h", Cells: []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Marked: []int{7}, Pointers: &sets.Pointers{{Name: "lo", Cell: 0}, {Name: "mid", Cell: 7}, {Name: "hi", Cell: 14}}, Window: &[2]int{0, 14}},
							{Label: "вода(mid)", Cells: []any{33}},
						}},
						{Caption: "mid = (1 + 7 + 1) // 2 = 4: 8 ≤ 9, lo = 4.", Rows: []sets.FigureRow{
							{Label: "h", Cells: []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Marked: []int{3}, Dimmed: []int{7, 8, 9, 10, 11, 12, 13, 14}, Pointers: &sets.Pointers{{Name: "lo", Cell: 0}, {Name: "mid", Cell: 3}, {Name: "hi", Cell: 6}}, Window: &[2]int{0, 6}},
							{Label: "вода(mid)", Cells: []any{8}},
						}},
						{Caption: "mid = 6: 19 > 9, hi = 5.", Rows: []sets.FigureRow{
							{Label: "h", Cells: []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Marked: []int{5}, Dimmed: []int{0, 1, 2, 7, 8, 9, 10, 11, 12, 13, 14}, Pointers: &sets.Pointers{{Name: "lo", Cell: 3}, {Name: "mid", Cell: 5}, {Name: "hi", Cell: 6}}, Window: &[2]int{3, 6}},
							{Label: "вода(mid)", Cells: []any{19}},
						}},
						{Caption: "mid = (4 + 5 + 1) // 2 = 5: 13 > 9, hi = 4. Без + 1 середина была бы 4, и lo = 4 повторялось бы без конца.", Rows: []sets.FigureRow{
							{Label: "h", Cells: []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Marked: []int{4}, Dimmed: []int{0, 1, 2, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}, Pointers: &sets.Pointers{{Name: "lo", Cell: 3}, {Name: "mid", Cell: 4}, {Name: "hi", Cell: 4}}, Window: &[2]int{3, 4}},
							{Label: "вода(mid)", Cells: []any{13}},
						}},
						{Caption: "lo = hi = 4: ответ 4 после 4 проверок.", Rows: []sets.FigureRow{
							{Label: "h", Cells: []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Dimmed: []int{0, 1, 2, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}, Pointers: &sets.Pointers{{Name: "ответ", Cell: 3}}},
							{Label: "вода(ответ)", Cells: []any{8}},
						}},
					},
				},
			},
		},
		{
			Title:   "Необычный порядок: гора, повёрнутый Array, Matrix",
			Summary: "Binary Search работает и без полностью упорядоченного Array, если по середине можно понять, в какой половине ответ. Вы ищете вершину горы по соседней точке, дно долины по знаку разности, минимум и число в повёрнутом Array и отрицательные числа в упорядоченной Matrix.",
			Lesson:  lesson4,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/count-negative-numbers-in-a-sorted-matrix/",
					Number:     number(1351),
					Title:      "Count Negative Numbers in a Sorted Matrix",
					Difficulty: "easy",
					Topics:     []string{"Array", "Binary Search", "Matrix"},
					Note:       "Ряды не возрастают, поэтому в каждом ряду ищите индекс j первого отрицательного числа через Binary Search; ряд даёт n − j отрицательных. m рядов, n столбцов. Время: O(m log n), не больше 7 проверок на ряд при n = 100. Память: O(1) дополнительно. Проход по всем клеткам даёт O(mn); обход лесенкой из левого нижнего угла даёт O(m + n), он использует порядок столбцов. Для примера 1: границы 3, 3, 2, 0, ответ 8.",
					Statement: `<p>Дана Matrix <code>grid</code> размера <code>m × n</code>. Её числа не возрастают и в каждом ряду, и в каждом столбце. Верните <em>количество <strong>отрицательных</strong> чисел в</em> <code>grid</code>.</p><p>Метод класса <code>Solution</code>: <code>def countNegatives(self, grid: list[list[int]]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> grid = [[4,3,2,-1],[3,2,1,-1],[1,1,-1,-2],[-1,-1,-2,-3]]
<strong>Выход:</strong> 8
<strong>Пояснение:</strong> В Matrix 8 отрицательных чисел.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> grid = [[3,2],[1,0]]
<strong>Выход:</strong> 0</pre><h3>Ограничения</h3><ul><li><code>m == grid.length</code></li><li><code>n == grid[i].length</code></li><li><code>1 ≤ m, n ≤ 100</code></li><li><code>-100 ≤ grid[i][j] ≤ 100</code></li></ul><p><strong>Дополнительный вопрос:</strong> можете ли вы решить задачу за <code>O(n + m)</code>?</p>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "countNegatives"}, Cases: []sets.TestCase{
						{Input: "[[[4,3,2,-1],[3,2,1,-1],[1,1,-1,-2],[-1,-1,-2,-3]]]", Expected: text("8")},
						{Input: "[[[3,2],[1,0]]]", Expected: text("0")},
						{Input: "[[[-1]]]", Expected: text("1")},
						{Input: "[[[5,1,0]]]", Expected: text("0")},
						{Input: "[[[-1,-2],[-3,-4],[-5,-6]]]", Expected: text("6")},
						{Input: "[[[3],[0],[-2],[-5]]]", Expected: text("2")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/peak-index-in-a-mountain-array/",
					Number:     number(852),
					Title:      "Peak Index in a Mountain Array",
					Difficulty: "medium",
					Topics:     []string{"Array", "Binary Search"},
					Note:       "Условие arr[i] > arr[i + 1] идёт «нет, ..., нет, да, ..., да»; вершина - первый индекс «да». Ищите его на [0, n − 1]. Время: O(log n), не больше 17 проверок для n = 10⁵. Память: O(1) дополнительно. Linear Scan даёт O(n). Ternary Search делает около 28 шагов для n = 10⁵. Для [1, 3, 6, 9, 7, 4, 2]: середины 9, 3, 6, ответ 3.",
					Statement: `<p>Дан Array целых чисел <code>arr</code> длины <code>n</code>, который является <strong>горой</strong>: значения растут до <strong>вершины</strong>, а потом убывают. Точнее, <code>arr[0] &lt; arr[1] &lt; ... &lt; arr[p] &gt; arr[p + 1] &gt; ... &gt; arr[n − 1]</code> для некоторого индекса <code>p</code>, <code>0 &lt; p &lt; n − 1</code>.</p><p>Верните индекс вершины.</p><p>Нужен алгоритм со временем <code>O(log(n))</code>.</p><p>Метод класса <code>Solution</code>: <code>def peakIndexInMountainArray(self, arr: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> arr = [0,1,0]
<strong>Выход:</strong> 1</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> arr = [0,2,1,0]
<strong>Выход:</strong> 1</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> arr = [0,10,5,2]
<strong>Выход:</strong> 1</pre><h3>Ограничения</h3><ul><li><code>3 ≤ arr.length ≤ 10<sup>5</sup></code></li><li><code>0 ≤ arr[i] ≤ 10<sup>6</sup></code></li><li><code>arr</code> <strong>гарантированно</strong> является горой</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "peakIndexInMountainArray"}, Cases: []sets.TestCase{
						{Input: "[[0,1,0]]", Expected: text("1")},
						{Input: "[[0,2,1,0]]", Expected: text("1")},
						{Input: "[[0,10,5,2]]", Expected: text("1")},
						{Input: "[[1,3,6,9,7,4,2]]", Expected: text("3")},
						{Input: "[[0,1,2,3,1]]", Expected: text("3")},
						{Input: "[[2,4,6,8,10,12,1]]", Expected: text("5")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/2009/E",
					Title:      "Klee's SUPER DUPER LARGE Array!!!",
					Difficulty: "medium",
					Topics:     []string{"Binary Search", "Math"},
					Note:       "Codeforces, рейтинг 1400. S(i) = i · k + i(i − 1) / 2, d(i) = 2S(i) − S(n) растёт. Binary Search находит первый i с d(i) ≥ 0, ответ min(|d(i − 1)|, |d(i)|). Время: O(log n) на набор, не больше 30 проверок; для t = 10⁴ наборов около 3 · 10⁵. Память: O(1) на набор, ответы O(t). Linear Scan по i даёт O(n) на набор, до 10⁹. Для n = 5, k = 3: d = [−19, −11, −1, 11, 25], ответ 1.",
					Statement: `<p>У Klee есть Array <code>a</code> длины <code>n</code>, в котором по порядку записаны целые числа <code>[k, k + 1, ..., k + n − 1]</code>. Klee хочет выбрать индекс <code>i</code> (<code>1 ≤ i ≤ n</code>) так, чтобы значение <code>x = |a<sub>1</sub> + a<sub>2</sub> + ... + a<sub>i</sub> − a<sub>i+1</sub> − ... − a<sub>n</sub>|</code> было наименьшим. Для любого целого числа <code>z</code> запись <code>|z|</code> обозначает модуль числа <code>z</code>.</p><p>Выведите наименьшее возможное значение <code>x</code>.</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>): количество наборов входных данных.</p><p>Каждый набор содержит два целых числа <code>n</code> и <code>k</code> (<code>2 ≤ n, k ≤ 10<sup>9</sup></code>): длину Array и его первый элемент.</p><h3>Выходные данные</h3><p>Для каждого набора выведите наименьшее значение <code>x</code> на отдельной строке.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>4
2 2
7 2
5 3
1000000000 1000000000</pre><p><strong>Выходные данные:</strong></p><pre>1
5
1
347369930</pre><h3>Примечание</h3><p>В первом наборе <code>a = [2, 3]</code>. При <code>i = 1</code> получается <code>x = |2 − 3| = 1</code>. Можно показать, что это наименьшее возможное значение <code>x</code>.</p><p>В третьем наборе <code>a = [3, 4, 5, 6, 7]</code>. При <code>i = 3</code> получается <code>x = |3 + 4 + 5 − 6 − 7| = 1</code>. Можно показать, что это наименьшее возможное значение <code>x</code>.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>2 ≤ n, k ≤ 10<sup>9</sup></code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "4\n2 2\n7 2\n5 3\n1000000000 1000000000\n", Expected: text("1\n5\n1\n347369930\n")},
						{Input: "1\n2 1000000000\n", Expected: text("1\n")},
						{Input: "1\n3 2\n", Expected: text("1\n")},
						{Input: "1\n1000000000 2\n", Expected: text("357502304\n")},
						{Input: "1\n4 10\n", Expected: text("4\n")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/find-minimum-in-rotated-sorted-array/",
					Number:     number(153),
					Title:      "Find Minimum in Rotated Sorted Array",
					Difficulty: "medium",
					Topics:     []string{"Array", "Binary Search"},
					Note:       "Числа левой части больше nums[−1], правой части не больше. Минимум - первый индекс i с nums[i] ≤ nums[−1]. Время: O(log n), не больше 13 проверок для n = 5000. Память: O(1) дополнительно. min(nums) даёт O(n). Для [4, 5, 6, 7, 0, 1, 2]: середины 7, 1, 0, ответ 0.",
					Statement: `<p>Array длины <code>n</code>, упорядоченный по возрастанию, <strong>повернули</strong> от <code>1</code> до <code>n</code> раз. Например, Array <code>nums = [0,1,2,4,5,6,7]</code> мог стать таким:</p><ul><li><code>[4,5,6,7,0,1,2]</code>, если его повернули <code>4</code> раза;</li><li><code>[0,1,2,4,5,6,7]</code>, если его повернули <code>7</code> раз.</li></ul><p>Один <strong>поворот</strong> Array <code>[a[0], a[1], a[2], ..., a[n-1]]</code> даёт Array <code>[a[n-1], a[0], a[1], a[2], ..., a[n-2]]</code>.</p><p>Дан упорядоченный повёрнутый Array <code>nums</code> из <strong>разных</strong> чисел. Верните <em>наименьший элемент этого Array</em>.</p><p>Нужен алгоритм со временем <code>O(log n)</code>.</p><p>Метод класса <code>Solution</code>: <code>def findMin(self, nums: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [3,4,5,1,2]
<strong>Выход:</strong> 1
<strong>Пояснение:</strong> Исходный Array [1,2,3,4,5] повернули 3 раза.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [4,5,6,7,0,1,2]
<strong>Выход:</strong> 0
<strong>Пояснение:</strong> Исходный Array [0,1,2,4,5,6,7] повернули 4 раза.</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> nums = [11,13,15,17]
<strong>Выход:</strong> 11
<strong>Пояснение:</strong> Исходный Array [11,13,15,17] повернули 4 раза.</pre><h3>Ограничения</h3><ul><li><code>n == nums.length</code></li><li><code>1 ≤ n ≤ 5000</code></li><li><code>-5000 ≤ nums[i] ≤ 5000</code></li><li>все числа <code>nums</code> <strong>разные</strong></li><li><code>nums</code> упорядочен и повёрнут от <code>1</code> до <code>n</code> раз</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "findMin"}, Cases: []sets.TestCase{
						{Input: "[[3,4,5,1,2]]", Expected: text("1")},
						{Input: "[[4,5,6,7,0,1,2]]", Expected: text("0")},
						{Input: "[[11,13,15,17]]", Expected: text("11")},
						{Input: "[[7]]", Expected: text("7")},
						{Input: "[[2,1]]", Expected: text("1")},
						{Input: "[[2,3,4,5,6,7,8,1]]", Expected: text("1")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/search-in-rotated-sorted-array/",
					Number:     number(33),
					Title:      "Search in Rotated Sorted Array",
					Difficulty: "medium",
					Topics:     []string{"Array", "Binary Search"},
					Note:       "Найдите индекс минимума p, как в LeetCode 153. Если target ≤ nums[−1], ищите в [p, n − 1], иначе в [0, p − 1]: обе части упорядочены. Время: O(log n), не больше 13 + 13 = 26 проверок для n = 5000. Память: O(1) дополнительно. Linear Scan даёт O(n). Для [4, 5, 6, 7, 0, 1, 2] и 0: p = 4, 0 ≤ 2, поиск в [4, 6], ответ 4.",
					Statement: `<p>Дан Array целых чисел <code>nums</code>, упорядоченный по возрастанию. Все числа в нём <strong>разные</strong>.</p><p>Перед вызовом функции <code>nums</code>, <strong>возможно, повернули влево</strong> на неизвестный индекс <code>k</code> (<code>1 ≤ k &lt; nums.length</code>), и получился Array <code>[nums[k], nums[k+1], ..., nums[n-1], nums[0], nums[1], ..., nums[k-1]]</code> (<strong>индексы с 0</strong>). Например, <code>[0,1,2,4,5,6,7]</code> можно повернуть влево на <code>3</code> индекса и получить <code>[4,5,6,7,0,1,2]</code>.</p><p>Дан Array <code>nums</code> <strong>после</strong> возможного поворота и целое число <code>target</code>. Верните <em>индекс</em> <code>target</code><em>, если он есть в</em> <code>nums</code><em>, или</em> <code>-1</code><em>, если его нет в</em> <code>nums</code>.</p><p>Нужен алгоритм со временем <code>O(log n)</code>.</p><p>Метод класса <code>Solution</code>: <code>def search(self, nums: list[int], target: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [4,5,6,7,0,1,2], target = 0
<strong>Выход:</strong> 4</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [4,5,6,7,0,1,2], target = 3
<strong>Выход:</strong> -1</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> nums = [1], target = 0
<strong>Выход:</strong> -1</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 5000</code></li><li><code>-10<sup>4</sup> ≤ nums[i] ≤ 10<sup>4</sup></code></li><li>все числа <code>nums</code> <strong>разные</strong></li><li><code>nums</code> упорядочен по возрастанию и, возможно, повёрнут</li><li><code>-10<sup>4</sup> ≤ target ≤ 10<sup>4</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "search"}, Cases: []sets.TestCase{
						{Input: "[[4,5,6,7,0,1,2],0]", Expected: text("4")},
						{Input: "[[4,5,6,7,0,1,2],3]", Expected: text("-1")},
						{Input: "[[1],0]", Expected: text("-1")},
						{Input: "[[1],1]", Expected: text("0")},
						{Input: "[[15,18,22,3,7,10],18]", Expected: text("1")},
						{Input: "[[6,7,1,2,3,4,5],5]", Expected: text("6")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "bs4-peak", Type: "short", Prompt: "arr = [2, 5, 8, 11, 9, 3]. Функция peak из урока начинает с lo = 0, hi = 5. Какой индекс она вернёт?", Answer: "3", Explain: "Проверки: mid = 2 (8 < 11, lo = 3), mid = 4 (9 > 3, hi = 4), mid = 3 (11 > 9, hi = 3). Условие arr[i] > arr[i + 1] впервые верно при i = 3. Ответ 4 - индекс первого числа после вершины."},
				{ID: "bs4-hi", Type: "choice", Prompt: "Почему в поиске вершины горы hi = len(arr) − 1, а не len(arr)?", Options: []string{"Иначе поиск станет O(n)", "Иначе arr[mid + 1] может выйти за конец Array", "Иначе вершина в индексе 0 потеряется"}, Answer: 1, Explain: "Проверка читает arr[mid + 1]. Для [0, 1, 2, 3, 1] и hi = 5 после lo = 3 получится mid = (3 + 5) // 2 = 4 и чтение arr[5]: IndexError. Последний индекс не может быть вершиной горы, поэтому hi = n − 1 ничего не теряет. Время O(log n) в обоих вариантах."},
				{ID: "bs4-rot", Type: "short", Prompt: "nums = [6, 7, 9, 15, 19, 2, 3]. Чему равен индекс минимума, первый i с nums[i] ≤ nums[−1]?", Answer: "5", Explain: "nums[−1] = 3. Числа 6, 7, 9, 15, 19 больше 3, число 2 первое не больше 3: индекс 5. Ответ 4 - индекс наибольшего числа 19, конец левой части."},
				{ID: "bs4-side", Type: "choice", Prompt: "nums = [15, 18, 22, 3, 7, 10], target = 18. Минимум уже найден: p = 3. Где искать target?", Options: []string{"в [3, 5], ведь 18 > nums[p] = 3", "во всём Array: он не упорядочен", "в [0, 2], ведь 18 > nums[−1] = 10"}, Answer: 2, Explain: "В правой части [3, 7, 10] все числа ≤ nums[−1] = 10. 18 > 10, поэтому 18 может стоять только в левой части [15, 18, 22]. Там обычный Binary Search находит индекс 1. Linear Scan по всему Array стоит O(n)."},
				{ID: "bs4-dups", Type: "choice", Prompt: "Почему правило nums[mid] ≤ nums[−1] не работает с повторами, например для [1, 1, 1, 0, 1]?", Options: []string{"При nums[mid] = nums[−1] минимум может быть и левее, и правее mid", "С повторами цикл while lo < hi не закончится", "Нужно сравнивать nums[mid] с nums[0], а не с nums[−1]"}, Answer: 0, Explain: "Для [1, 1, 1, 0, 1] mid = 2, nums[2] = 1 ≤ 1, правило ставит hi = 2 и теряет 0 в индексе 3. В [1, 0, 1, 1, 1] такое же равенство стоит правее минимума. Цикл закончится: hi = mid или lo = mid + 1 сужают отрезок на каждом шаге. Сравнение с nums[0] тоже не помогает: на [1, 1, 1, 0, 1] условие nums[i] < nums[0] идёт «нет, нет, нет, да, нет», а не «нет, ..., да». Задачи главы гарантируют разные числа; с повторами в худшем случае нужно O(n) проверок."},
				{ID: "bs4-klee", Type: "short", Prompt: "Codeforces 2009E: n = 6, k = 2, a = [2, 3, 4, 5, 6, 7], сумма 27. Чему равен ответ x?", Answer: "1", Explain: "S(3) = 9, S(4) = 14. d(3) = 18 − 27 = −9, d(4) = 28 − 27 = 1. Первый i с d(i) ≥ 0 равен 4, ответ min(9, 1) = 1. Брать только |d(i)| нельзя: для n = 5, k = 3 из урока |d(4)| = 11, а ответ |d(3)| = 1."},
				{ID: "bs4-rows", Type: "short", Prompt: "Matrix 3 × 5, каждый ряд не возрастает. Binary Search нашёл в рядах индексы первых отрицательных чисел 4, 2 и 0. Сколько отрицательных чисел в Matrix?", Answer: "9", Explain: "В ряду с границей j отрицательных n − j: 1 + 3 + 5 = 9. Сумма границ 4 + 2 + 0 = 6 - это число неотрицательных, 15 − 6 = 9 сходится. Каждый поиск стоит O(log n), всего O(m log n) времени и O(1) памяти."},
				{ID: "bs4-index", Type: "short", Prompt: "LeetCode 74: Matrix из 3 рядов по 4 числа читают ряд за рядом как один упорядоченный Array. В каком ряду (с 0) стоит элемент с номером k = 9?", Answer: "2", Explain: "Ряд k // 4 = 2, столбец k % 4 = 1. 9 // 3 = 3 - ошибка: делить нужно на длину ряда, а не на число рядов. Один Binary Search по 12 номерам делает не больше 4 проверок."},
			},
			Links: []sets.Link{
				{Title: "Поиск в повёрнутом Array", URL: "https://sprintcode.pro/ru/blog/search-rotated-sorted-array", Source: "SprintCode.pro"},
				{Title: "Повёрнутый Array и пик: задачи на Python", URL: "https://thecode.media/binarnyj-poisk-python/", Source: "Журнал «Код»"},
				{Title: "Ternary Search", URL: "https://neerc.ifmo.ru/wiki/index.php?title=%D0%A2%D1%80%D0%BE%D0%B8%D1%87%D0%BD%D1%8B%D0%B9_%D0%BF%D0%BE%D0%B8%D1%81%D0%BA", Source: "Викиконспекты ИТМО"},
				{Title: "Поиск в упорядоченной Matrix", URL: "https://neerc.ifmo.ru/wiki/index.php?title=%D0%9F%D0%BE%D0%B8%D1%81%D0%BA_%D0%B2_%D0%BC%D0%B0%D1%82%D1%80%D0%B8%D1%86%D0%B5", Source: "Викиконспекты ИТМО"},
			},
			Figures: []sets.Figure{
				{
					ID:    "bs4-rows",
					Title: "Отрицательные числа в рядах Matrix",
					Frames: []sets.FigureFrame{
						{Caption: "Ряд 0 не возрастает. Binary Search за 2 проверки находит индекс первого отрицательного числа: 3. Отрицательных 4 − 3 = 1.", Rows: []sets.FigureRow{
							{Label: "ряд 0", Cells: []any{4, 3, 2, -1}, Marked: []int{3}, Pointers: &sets.Pointers{{Name: "граница", Cell: 3}}},
							{Label: "ряд 1", Cells: []any{3, 2, 1, -1}},
							{Label: "ряд 2", Cells: []any{1, 1, -1, -2}},
							{Label: "ряд 3", Cells: []any{-1, -1, -2, -3}},
						}},
						{Caption: "Ряд 1: граница 3 за 2 проверки, всего 2 отрицательных.", Rows: []sets.FigureRow{
							{Label: "ряд 0", Cells: []any{4, 3, 2, -1}, Marked: []int{3}},
							{Label: "ряд 1", Cells: []any{3, 2, 1, -1}, Marked: []int{3}, Pointers: &sets.Pointers{{Name: "граница", Cell: 3}}},
							{Label: "ряд 2", Cells: []any{1, 1, -1, -2}},
							{Label: "ряд 3", Cells: []any{-1, -1, -2, -3}},
						}},
						{Caption: "Ряд 2: граница 2 за 2 проверки, всего 4.", Rows: []sets.FigureRow{
							{Label: "ряд 0", Cells: []any{4, 3, 2, -1}, Marked: []int{3}},
							{Label: "ряд 1", Cells: []any{3, 2, 1, -1}, Marked: []int{3}},
							{Label: "ряд 2", Cells: []any{1, 1, -1, -2}, Marked: []int{2, 3}, Pointers: &sets.Pointers{{Name: "граница", Cell: 2}}},
							{Label: "ряд 3", Cells: []any{-1, -1, -2, -3}},
						}},
						{Caption: "Ряд 3: граница 0 за 3 проверки. Ответ 1 + 1 + 2 + 4 = 8 после 9 проверок; Linear Scan прочитал бы 16 клеток.", Rows: []sets.FigureRow{
							{Label: "ряд 0", Cells: []any{4, 3, 2, -1}, Marked: []int{3}},
							{Label: "ряд 1", Cells: []any{3, 2, 1, -1}, Marked: []int{3}},
							{Label: "ряд 2", Cells: []any{1, 1, -1, -2}, Marked: []int{2, 3}},
							{Label: "ряд 3", Cells: []any{-1, -1, -2, -3}, Marked: []int{0, 1, 2, 3}, Pointers: &sets.Pointers{{Name: "граница", Cell: 0}}},
						}},
					},
				},
				{
					ID:    "bs4-peak",
					Title: "Вершина горы: сравнение с соседом",
					Frames: []sets.FigureFrame{
						{Caption: "lo = 0, hi = 6, mid = 3: 9 > 7, спуск уже начался. Вершина не правее 3: hi = 3.", Rows: []sets.FigureRow{
							{Label: "arr", Cells: []any{1, 3, 6, 9, 7, 4, 2}, Marked: []int{3, 4}, Pointers: &sets.Pointers{{Name: "lo", Cell: 0}, {Name: "mid", Cell: 3}, {Name: "hi", Cell: 6}}, Window: &[2]int{0, 6}},
							{Label: "arr[i] > arr[i + 1]", Cells: []any{"нет", "нет", "нет", "да", "да", "да", ""}, Marked: []int{3}},
						}},
						{Caption: "mid = 1: 3 < 6, это подъём. Вершина правее: lo = 2.", Rows: []sets.FigureRow{
							{Label: "arr", Cells: []any{1, 3, 6, 9, 7, 4, 2}, Marked: []int{1, 2}, Dimmed: []int{4, 5, 6}, Pointers: &sets.Pointers{{Name: "lo", Cell: 0}, {Name: "mid", Cell: 1}, {Name: "hi", Cell: 3}}, Window: &[2]int{0, 3}},
							{Label: "arr[i] > arr[i + 1]", Cells: []any{"нет", "нет", "нет", "да", "да", "да", ""}, Marked: []int{1}},
						}},
						{Caption: "mid = 2: 6 < 9, lo = 3.", Rows: []sets.FigureRow{
							{Label: "arr", Cells: []any{1, 3, 6, 9, 7, 4, 2}, Marked: []int{2, 3}, Dimmed: []int{0, 1, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "lo", Cell: 2}, {Name: "mid", Cell: 2}, {Name: "hi", Cell: 3}}, Window: &[2]int{2, 3}},
							{Label: "arr[i] > arr[i + 1]", Cells: []any{"нет", "нет", "нет", "да", "да", "да", ""}, Marked: []int{2}},
						}},
						{Caption: "lo = hi = 3: вершина arr[3] = 9 после 3 проверок.", Rows: []sets.FigureRow{
							{Label: "arr", Cells: []any{1, 3, 6, 9, 7, 4, 2}, Marked: []int{3}, Pointers: &sets.Pointers{{Name: "вершина", Cell: 3}}},
							{Label: "arr[i] > arr[i + 1]", Cells: []any{"нет", "нет", "нет", "да", "да", "да", ""}, Marked: []int{3}},
						}},
					},
				},
				{
					ID:    "bs4-valley",
					Title: "Klee: дно долины |d(i)|",
					Frames: []sets.FigureFrame{
						{Caption: "n = 5, k = 3, a = [3, 4, 5, 6, 7], сумма 25. d(i) растёт, а |d(i)| сначала убывает, потом растёт. mid = 3: d(3) = −1 < 0, lo = 4.", Rows: []sets.FigureRow{
							{Label: "i", Cells: []any{1, 2, 3, 4, 5}, Marked: []int{2}, Pointers: &sets.Pointers{{Name: "lo", Cell: 0}, {Name: "mid", Cell: 2}, {Name: "hi", Cell: 4}}, Window: &[2]int{0, 4}},
							{Label: "d(i)", Cells: []any{-19, -11, -1, 11, 25}, Marked: []int{2}},
							{Label: "|d(i)|", Cells: []any{19, 11, 1, 11, 25}},
						}},
						{Caption: "mid = 4: d(4) = 11 ≥ 0, hi = 4.", Rows: []sets.FigureRow{
							{Label: "i", Cells: []any{1, 2, 3, 4, 5}, Marked: []int{3}, Dimmed: []int{0, 1, 2}, Pointers: &sets.Pointers{{Name: "lo", Cell: 3}, {Name: "mid", Cell: 3}, {Name: "hi", Cell: 4}}, Window: &[2]int{3, 4}},
							{Label: "d(i)", Cells: []any{-19, -11, -1, 11, 25}, Marked: []int{3}},
							{Label: "|d(i)|", Cells: []any{19, 11, 1, 11, 25}},
						}},
						{Caption: "Первый i с d(i) ≥ 0 равен 4. Сравните |d(3)| = 1 и |d(4)| = 11: ответ 1, как в примере.", Rows: []sets.FigureRow{
							{Label: "i", Cells: []any{1, 2, 3, 4, 5}, Marked: []int{2, 3}, Pointers: &sets.Pointers{{Name: "i − 1", Cell: 2}, {Name: "i", Cell: 3}}},
							{Label: "d(i)", Cells: []any{-19, -11, -1, 11, 25}, Marked: []int{2, 3}},
							{Label: "|d(i)|", Cells: []any{19, 11, 1, 11, 25}, Marked: []int{2}},
						}},
					},
				},
				{
					ID:    "bs4-rotated",
					Title: "Минимум повёрнутого Array",
					Frames: []sets.FigureFrame{
						{Caption: "nums[−1] = 2. mid = 3: 7 > 2, минимум правее: lo = 4.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{4, 5, 6, 7, 0, 1, 2}, Marked: []int{3}, Pointers: &sets.Pointers{{Name: "lo", Cell: 0}, {Name: "mid", Cell: 3}, {Name: "hi", Cell: 6}}, Window: &[2]int{0, 6}},
							{Label: "nums[i] ≤ 2", Cells: []any{"нет", "нет", "нет", "нет", "да", "да", "да"}, Marked: []int{3}},
						}},
						{Caption: "mid = 5: 1 ≤ 2, минимум не правее 5: hi = 5.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{4, 5, 6, 7, 0, 1, 2}, Marked: []int{5}, Dimmed: []int{0, 1, 2, 3}, Pointers: &sets.Pointers{{Name: "lo", Cell: 4}, {Name: "mid", Cell: 5}, {Name: "hi", Cell: 6}}, Window: &[2]int{4, 6}},
							{Label: "nums[i] ≤ 2", Cells: []any{"нет", "нет", "нет", "нет", "да", "да", "да"}, Marked: []int{5}},
						}},
						{Caption: "mid = 4: 0 ≤ 2, hi = 4.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{4, 5, 6, 7, 0, 1, 2}, Marked: []int{4}, Dimmed: []int{0, 1, 2, 3, 6}, Pointers: &sets.Pointers{{Name: "lo", Cell: 4}, {Name: "mid", Cell: 4}, {Name: "hi", Cell: 5}}, Window: &[2]int{4, 5}},
							{Label: "nums[i] ≤ 2", Cells: []any{"нет", "нет", "нет", "нет", "да", "да", "да"}, Marked: []int{4}},
						}},
						{Caption: "lo = hi = 4: минимум nums[4] = 0 после 3 проверок. Части [4, 5, 6, 7] и [0, 1, 2] упорядочены.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{4, 5, 6, 7, 0, 1, 2}, Marked: []int{4}, Pointers: &sets.Pointers{{Name: "минимум", Cell: 4}}},
							{Label: "nums[i] ≤ 2", Cells: []any{"нет", "нет", "нет", "нет", "да", "да", "да"}, Marked: []int{4}},
						}},
					},
				},
			},
		},
	},
}
