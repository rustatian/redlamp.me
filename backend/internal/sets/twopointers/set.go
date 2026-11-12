// Package twopointers is the learning set "Two Pointers".
package twopointers

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
	Slug:    "two-pointers",
	Title:   "Two Pointers",
	Summary: "Two Pointers сокращает число проверок в Array. Вы учите три формы приёма, доказываете каждый ход и считаете время, дополнительную память и стоимость Sorting.",
	Topics:  []string{"Two Pointers", "Array", "Sorting"},
	Chapters: []sets.Chapter{
		{
			Title:   "Навстречу",
			Summary: "Two Pointers проходит отсортированный Array с двух концов. Вы доказываете исключение пар и считаете суммы, движения и дополнительную память с учётом Sorting.",
			Lesson:  lesson1,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/two-sum-ii-input-array-is-sorted/",
					Number:     number(167),
					Title:      "Two Sum II - Input Array Is Sorted",
					Difficulty: "medium",
					Topics:     []string{"Array", "Two Pointers"},
					Note:       "n - длина уже отсортированного Array, числа ограниченного размера. Время: не более n-1 вычислений суммы, O(n). Память: O(1) дополнительной, ответ - два индекса с единицы. При малой сумме исключайте левое число, при большой - правое. Brute Force проверяет до n(n-1)/2 пар; Hash Table требует O(u) записей для u разных значений.",
					Statement: `<p>Дан Array целых чисел <code>numbers</code>. Номера элементов в нём <strong>начинаются с 1</strong>. Array уже <strong>отсортирован по неубыванию</strong> (каждое следующее число не меньше предыдущего).</p><p>Найдите <strong>два</strong> числа, сумма которых равна числу <code>target</code>. Пусть это числа <code>numbers[index<sub>1</sub>]</code> и <code>numbers[index<sub>2</sub>]</code>, где <code>1 ≤ index<sub>1</sub> &lt; index<sub>2</sub> ≤ numbers.length</code>.</p><p>Верните номера этих двух чисел <code>index<sub>1</sub></code> и <code>index<sub>2</sub></code> в виде Array <code>[index<sub>1</sub>, index<sub>2</sub>]</code> длины 2.</p><p>Тесты составлены так, что решение <strong>ровно одно</strong>. Один и тот же элемент <strong>нельзя</strong> брать два раза.</p><p>Ваше решение должно использовать только постоянную дополнительную память.</p><p>Метод класса <code>Solution</code>: <code>def twoSum(self, numbers: list[int], target: int) -&gt; list[int]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> numbers = [2,7,11,15], target = 9
<strong>Выход:</strong> [1,2]
<strong>Пояснение:</strong> Сумма 2 и 7 равна 9. Значит, index1 = 1, index2 = 2. Верните [1, 2].</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> numbers = [2,3,4], target = 6
<strong>Выход:</strong> [1,3]
<strong>Пояснение:</strong> Сумма 2 и 4 равна 6. Значит, index1 = 1, index2 = 3. Верните [1, 3].</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> numbers = [-1,0], target = -1
<strong>Выход:</strong> [1,2]
<strong>Пояснение:</strong> Сумма -1 и 0 равна -1. Значит, index1 = 1, index2 = 2. Верните [1, 2].</pre><h3>Ограничения</h3><ul><li><code>2 ≤ numbers.length ≤ 3 · 10<sup>4</sup></code></li><li><code>-1000 ≤ numbers[i] ≤ 1000</code></li><li><code>numbers</code> отсортирован по неубыванию</li><li><code>-1000 ≤ target ≤ 1000</code></li><li>тесты составлены так, что решение ровно одно</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "twoSum"}, Cases: []sets.TestCase{
						{Input: "[[2,7,11,15],9]", Expected: text("[1,2]")},
						{Input: "[[2,3,4],6]", Expected: text("[1,3]")},
						{Input: "[[-1,0],-1]", Expected: text("[1,2]")},
						{Input: "[[5,25,75],100]", Expected: text("[2,3]")},
						{Input: "[[1,2,3,4,4,9,56,90],8]", Expected: text("[4,5]")},
						{Input: "[[-1000,-3,0,1000],0]", Expected: text("[1,4]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/valid-palindrome/",
					Number:     number(125),
					Title:      "Valid Palindrome",
					Difficulty: "easy",
					Topics:     []string{"String", "Two Pointers"},
					Note:       "n - длина String из ASCII. Пропускайте символы, кроме букв и цифр, и сравнивайте края без учёта регистра. Время: индексы суммарно проходят не более n позиций, O(n). Память: O(1) дополнительной и Boolean-ответ. Фильтрация с созданием String и её разворот тоже требуют O(n) времени, но O(n) дополнительной памяти.",
					Statement: `<p>Palindrome читается одинаково слева направо и справа налево после двух действий: все заглавные буквы заменили на строчные и удалили все символы, кроме букв и цифр.</p><p>Дана String <code>s</code>. Верните <code>True</code>, если это <strong>Palindrome</strong>, иначе верните <code>False</code>.</p><p>Метод класса <code>Solution</code>: <code>def isPalindrome(self, s: str) -&gt; bool</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> s = "A man, a plan, a canal: Panama"
<strong>Выход:</strong> True
<strong>Пояснение:</strong> "amanaplanacanalpanama" читается одинаково в обоих направлениях.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> s = "race a car"
<strong>Выход:</strong> False
<strong>Пояснение:</strong> "raceacar" читается по-разному в двух направлениях.</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> s = " "
<strong>Выход:</strong> True
<strong>Пояснение:</strong> После удаления всех символов, кроме букв и цифр, s становится пустой String "".
Пустая String читается одинаково в обоих направлениях, поэтому это Palindrome.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ s.length ≤ 2 · 10<sup>5</sup></code></li><li><code>s</code> состоит только из печатаемых символов ASCII</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "isPalindrome"}, Cases: []sets.TestCase{
						{Input: `["A man, a plan, a canal: Panama"]`, Expected: text("True")},
						{Input: `["race a car"]`, Expected: text("False")},
						{Input: `[" "]`, Expected: text("True")},
						{Input: `["0P"]`, Expected: text("False")},
						{Input: `[".,"]`, Expected: text("True")},
						{Input: `["ab_a"]`, Expected: text("True")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/max-number-of-k-sum-pairs/",
					Number:     number(1679),
					Title:      "Max Number of K-Sum Pairs",
					Difficulty: "medium",
					Topics:     []string{"Array", "Two Pointers", "Sorting"},
					Note:       "n - длина Array, u - число разных значений; числа ограниченного размера. Sorting допустим: нужны количество пар, а не исходные индексы. Время: O(n log n) на Sorting плюс O(n) на проход. Память: O(n) дополнительной в худшем случае Python Sorting; ответ - счётчик. Hash Table частот даёт ожидаемое O(n) время при хороших хешах и O(u) памяти. При сумме k используйте оба числа.",
					Statement: `<p>Дан Array целых чисел <code>nums</code> и целое число <code>k</code>.</p><p>За одну операцию можно выбрать в Array два числа с суммой <code>k</code> и удалить их.</p><p>Верните <em>наибольшее количество таких операций</em>.</p><p>Метод класса <code>Solution</code>: <code>def maxOperations(self, nums: List[int], k: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,2,3,4], k = 5
<strong>Выход:</strong> 2
<strong>Пояснение:</strong> Начинаем с nums = [1,2,3,4]:
- Удаляем числа 1 и 4, получаем nums = [2,3]
- Удаляем числа 2 и 3, получаем nums = []
Пар с суммой 5 больше нет, поэтому всего 2 операции.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [3,1,3,4,3], k = 6
<strong>Выход:</strong> 1
<strong>Пояснение:</strong> Начинаем с nums = [3,1,3,4,3]:
- Удаляем первые две тройки, получаем nums = [1,4,3]
Пар с суммой 6 больше нет, поэтому всего 1 операция.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 10<sup>5</sup></code></li><li><code>1 ≤ nums[i] ≤ 10<sup>9</sup></code></li><li><code>1 ≤ k ≤ 10<sup>9</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "maxOperations"}, Cases: []sets.TestCase{
						{Input: "[[1,2,3,4],5]", Expected: text("2")},
						{Input: "[[3,1,3,4,3],6]", Expected: text("1")},
						{Input: "[[5],10]", Expected: text("0")},
						{Input: "[[2,2,2,2,2],4]", Expected: text("2")},
						{Input: "[[1,999999999],1000000000]", Expected: text("1")},
						{Input: "[[4,4,1,3,1,3,2,2,5,5,1,5,2,1,2,3,5,4],2]", Expected: text("2")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/container-with-most-water/",
					Number:     number(11),
					Title:      "Container With Most Water",
					Difficulty: "medium",
					Topics:     []string{"Array", "Two Pointers"},
					Note:       "n - число линий, высоты ограниченного размера. Площадь равна min(h[l], h[r]) · (r-l). Исключайте меньшую высоту: с ней более узкая пара не улучшит результат. Время: n-1 вычислений площади и сужений, O(n). Память: O(1) дополнительной и одно число в ответе. Brute Force проверяет n(n-1)/2 пар за O(n²) времени с O(1) памятью.",
					Statement: `<p>Дан Array целых чисел <code>height</code> длины <code>n</code>. Нарисовано <code>n</code> вертикальных линий. Концы <code>i</code>-й линии находятся в точках <code>(i, 0)</code> и <code>(i, height[i])</code>.</p><p>Найдите две линии, которые вместе с осью x образуют сосуд, в котором больше всего воды.</p><p>Верните <em>наибольшее количество воды, которое может вместить сосуд</em>.</p><p><strong>Обратите внимание:</strong> сосуд нельзя наклонять.</p><p>Метод класса <code>Solution</code>: <code>def maxArea(self, height: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> height = [1,8,6,2,5,4,8,3,7]
<strong>Выход:</strong> 49
<strong>Пояснение:</strong> Вертикальные линии заданы Array [1,8,6,2,5,4,8,3,7]. Наибольшая площадь равна 49: линии с индексами 1 и 8 дают min(8,7) · (8-1) = 49.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> height = [1,1]
<strong>Выход:</strong> 1</pre><h3>Ограничения</h3><ul><li><code>n == height.length</code></li><li><code>2 ≤ n ≤ 10<sup>5</sup></code></li><li><code>0 ≤ height[i] ≤ 10<sup>4</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "maxArea"}, Cases: []sets.TestCase{
						{Input: "[[1,8,6,2,5,4,8,3,7]]", Expected: text("49")},
						{Input: "[[1,1]]", Expected: text("1")},
						{Input: "[[0,0]]", Expected: text("0")},
						{Input: "[[4,3,2,1,4]]", Expected: text("16")},
						{Input: "[[1,2,1]]", Expected: text("2")},
						{Input: "[[2,3,4,5,18,17,6]]", Expected: text("17")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/boats-to-save-people/",
					Number:     number(881),
					Title:      "Boats to Save People",
					Difficulty: "medium",
					Topics:     []string{"Array", "Two Pointers", "Sorting"},
					Note:       "n - число людей, целые веса ограничены; в лодке не больше двух. После Sorting самый тяжёлый едет с самым лёгким, если они помещаются, иначе один. Время: O(n log n) на Sorting и O(n) на посадку. Память: O(n) дополнительной в худшем случае Python Sorting; ответ - счётчик лодок. Для весов 1..W частоты дают O(n+W) времени и O(W) памяти.",
					Statement: `<p>Дан Array <code>people</code>, где <code>people[i]</code> это вес <code>i</code>-го человека. Лодок <strong>сколько угодно</strong>, и каждая лодка выдерживает вес не больше <code>limit</code>. В одной лодке может быть не больше двух человек сразу, если их общий вес не больше <code>limit</code>.</p><p>Верните <em>наименьшее количество лодок, чтобы перевезти всех людей</em>.</p><p>Метод класса <code>Solution</code>: <code>def numRescueBoats(self, people: list[int], limit: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> people = [1,2], limit = 3
<strong>Выход:</strong> 1
<strong>Пояснение:</strong> 1 лодка: (1, 2)</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> people = [3,2,2,1], limit = 3
<strong>Выход:</strong> 3
<strong>Пояснение:</strong> 3 лодки: (1, 2), (2) и (3)</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> people = [3,5,3,4], limit = 5
<strong>Выход:</strong> 4
<strong>Пояснение:</strong> 4 лодки: (3), (3), (4), (5)</pre><h3>Ограничения</h3><ul><li><code>1 ≤ people.length ≤ 5 · 10<sup>4</sup></code></li><li><code>1 ≤ people[i] ≤ limit ≤ 3 · 10<sup>4</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "numRescueBoats"}, Cases: []sets.TestCase{
						{Input: "[[1,2],3]", Expected: text("1")},
						{Input: "[[3,2,2,1],3]", Expected: text("3")},
						{Input: "[[3,5,3,4],5]", Expected: text("4")},
						{Input: "[[5],5]", Expected: text("1")},
						{Input: "[[2,2,2,2],4]", Expected: text("2")},
						{Input: "[[1,1,1,1,1],3]", Expected: text("3")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "tp1-less-move", Type: "choice", Prompt: "a = [2, 3, 5, 8, 9, 13], target = 14. Сейчас l = 0, r = 4. Two Pointers проверяет сумму a[l] + a[r] и делает один ход. Какие пары, кроме проверенной, исключает этот ход?", Options: []string{"2 + 3, 2 + 5, 2 + 8", "3 + 9, 5 + 9, 8 + 9", "2 + 3, 2 + 5, 2 + 8, 2 + 13", "3 + 9, 5 + 9, 8 + 9, 9 + 13"}, Answer: 0, Explain: "2 + 9 = 11 меньше 14, поэтому l += 1 исключает строку числа 2 в окне: 2 + 3, 2 + 5, 2 + 8. Все суммы ≤ 11. Пара 2 + 13 уже исключена прошлым ходом r -= 1. В столбце числа 9 ещё есть ответ 5 + 9 = 14. Одно вычисление суммы исключило ещё 3 пары; таблицу хранить не нужно. Для n ограниченных по размеру чисел полный поиск занимает O(n) времени и O(1) дополнительной памяти."},
				{ID: "tp1-why-not-r", Type: "choice", Prompt: "a = [2, 5, 6, 9], target = 14. В коде ученика при сумме меньше target стоит r -= 1. Какие суммы проверит его цикл while l < r?", Options: []string{"11, 8", "11, 8, 7", "11, 14", "11, 8, 7, 4"}, Answer: 1, Explain: "2 + 9 = 11, затем 2 + 6 = 8 и 2 + 5 = 7. После трёх ходов r = l, цикл завершён. Первый неверный ход исключил ответ 5 + 9 = 14. Верный ход l += 1 находит его второй проверкой. Даже неправильный цикл имеет O(n) время и O(1) дополнительную память для n чисел ограниченного размера. Скорость не доказывает правильность."},
				{ID: "tp1-invariant", Type: "short", Prompt: "В Array 8 чисел, номера 0..7. Сейчас l = 2, r = 5, сумму a[2] + a[5] ещё не проверили. Сколько пар номеров (i, j), i < j, по инварианту ещё могут быть ответом?", Answer: "6", Explain: "В окне 4 числа: номера 2, 3, 4, 5. Из них 4 · 3 / 2 = 6 пар. Остальные 28 - 6 = 22 пары исключены. Количество кандидатов не равно числу вычисленных сумм: для n = 8 нужно не более 7 вычислений за весь поиск. Время O(n), дополнительная память O(1) при ограниченном размере чисел; пары не хранятся."},
				{ID: "tp1-code-greater", Type: "code", Prompt: "Это решение LeetCode 1679. На каком входе оно вернёт неверный ответ?", Language: "python", Code: `nums.sort()
l, r = 0, len(nums) - 1
ops = 0
while l < r:
    s = nums[l] + nums[r]
    if s == k:
        ops += 1
        l += 1
    elif s < k:
        l += 1
    else:
        r -= 1
return ops`, Options: []string{"[1, 4, 1], k = 5", "[4, 1, 4], k = 5", "[4, 1, 1, 4], k = 5"}, Answer: 0, Explain: "При сумме k код двигает только l и может повторно использовать nums[r]. [1,4,1] после Sorting даёт [1,1,4]: две суммы 1 + 4, ответ 2 вместо 1. Для [4,1,4] код вернёт 1, а для [4,1,1,4] - 2; эти ответы случайно верны. При сумме k двигайте оба указателя. Для n чисел ограниченного размера Sorting плюс проход требуют O(n log n) времени и O(n) дополнительной памяти Python в худшем случае. Счётчик ответа занимает O(1)."},
				{ID: "tp1-trace", Type: "short", Prompt: "a = [3, 4, 6, 8, 11, 12, 16], target = 17. Сколько сумм проверит Two Pointers, считая последнюю, равную 17?", Answer: "5", Explain: "3 + 16 = 19: r -= 1. 3 + 12 = 15: l += 1. 4 + 12 = 16: l += 1. 6 + 12 = 18: r -= 1. 6 + 11 = 17: успех. Получились 5 сумм, 4 хода и 3 рабочие переменные l, r, s. Для n чисел ограниченного размера это O(n) времени и O(1) дополнительной памяти, ответ - два индекса."},
				{ID: "tp1-moves", Type: "short", Prompt: "В Array 10 чисел, номера 0..9. Two Pointers нашёл пару, когда l = 3 и r = 5. Сколько ходов сделано до этого?", Answer: "7", Explain: "l прошёл от 0 до 3: 3 хода. r прошёл от 9 до 5: 4 хода. Всего 7 ходов и 8 вычислений суммы, включая успешное. Для n чисел ограниченного размера вычислений не более n-1, время O(n). Дополнительная память O(1); ответ содержит два индекса."},
				{ID: "tp1-sort-allowed", Type: "choice", Prompt: "Array не отсортирован. Нужно вернуть исходные номера (с единицы) двух чисел с суммой target. Ученик применяет Sorting и запускает код из задачи 167. На каком входе код вернёт верный ответ?", Options: []string{"[8, 1, 5], target = 6", "[3, 9, 4], target = 13", "[1, 8, 5], target = 9"}, Answer: 1, Explain: "Код возвращает номера после Sorting. [8,1,5] → [1,5,8]: ответ [1,2], нужны [2,3]. [1,8,5] → [1,5,8]: ответ [1,3], нужны [1,2]. [3,9,4] → [3,4,9]: ответ [2,3] случайно совпал. Храните (value,index): для n чисел ограниченного размера O(n log n) времени и O(n) дополнительной памяти Python. Альтернатива - Hash Table значений и индексов: ожидаемое O(n) время при хороших хешах, O(n) памяти. Ответ в обоих случаях - два индекса. В 1679 нужен только счёт, поэтому исходные индексы не нужны."},
				{ID: "tp1-container", Type: "choice", Prompt: "LeetCode 11, высоты [2, 7, 5, 6, 3]. Ученик всегда двигает указатель у более высокой линии. Какую наибольшую площадь найдёт его код?", Options: []string{"12", "9", "8", "6"}, Answer: 2, Explain: "Неверный код проверяет площади 8, 6, 4, 2: ответ 8. Правильный код двигает меньшую высоту и проверяет 8, 9, 12, 5: ответ 12. Для n = 5 оба делают 4 вычисления площади. Время O(n), дополнительная память O(1) при ограниченных высотах. Только у правильного хода есть доказательство: при той же меньшей высоте более узкий сосуд не улучшит найденную площадь."},
			},
			Links: []sets.Link{
				{Title: "Задачи LeetCode: Two Pointers", URL: "https://habr.com/ru/articles/964104/", Source: "Хабр"},
				{Title: "Two Pointers: навстречу и в одну сторону", URL: "https://habr.com/ru/articles/1020222/", Source: "Хабр"},
				{Title: "Two Pointers: Two Sum и Sliding Window", URL: "https://usaco.guide/silver/two-pointers", Source: "USACO Guide (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "tp1-table",
					Title: "Two Pointers: 5 сумм и 4 хода среди 21 пары",
					Frames: []sets.FigureFrame{
						{Caption: "l = 0, r = 6: 1 + 16 = 17. Остальные суммы столбца 6 не меньше 17. Следующий ход r -= 1 безопасен. Таблица иллюстрирует пары; алгоритм её не хранит.", Rows: []sets.FigureRow{
							{Label: `i \ j`, Cells: []any{1, 2, 3, 4, 5, 6}},
							{Label: "0", Cells: []any{5, 6, 8, 11, 13, 17}, Marked: []int{5}},
							{Label: "1", Cells: []any{"", 9, 11, 14, 16, 20}},
							{Label: "2", Cells: []any{"", "", 12, 15, 17, 21}},
							{Label: "3", Cells: []any{"", "", "", 17, 19, 23}},
							{Label: "4", Cells: []any{"", "", "", "", 22, 26}},
							{Label: "5", Cells: []any{"", "", "", "", "", 28}},
							{Label: "Сумм вычислено", Cells: []any{1}},
							{Label: "Ходов сделано", Cells: []any{0}},
							{Label: "Всего пар", Cells: []any{21}},
							{Label: "Доп. память: 3 значения", Cells: []any{"l", "r", "s", "O(1)"}},
						}},
						{Caption: "l = 0, r = 5: 1 + 12 = 13. Остальные суммы строки 0 не больше 13. Следующий ход l += 1 исключит строку без ответа.", Rows: []sets.FigureRow{
							{Label: `i \ j`, Cells: []any{1, 2, 3, 4, 5, 6}, Dimmed: []int{5}},
							{Label: "0", Cells: []any{5, 6, 8, 11, 13, 17}, Marked: []int{4}, Dimmed: []int{5}},
							{Label: "1", Cells: []any{"", 9, 11, 14, 16, 20}, Dimmed: []int{5}},
							{Label: "2", Cells: []any{"", "", 12, 15, 17, 21}, Dimmed: []int{5}},
							{Label: "3", Cells: []any{"", "", "", 17, 19, 23}, Dimmed: []int{5}},
							{Label: "4", Cells: []any{"", "", "", "", 22, 26}, Dimmed: []int{5}},
							{Label: "5", Cells: []any{"", "", "", "", "", 28}, Dimmed: []int{5}},
							{Label: "Сумм вычислено", Cells: []any{2}},
							{Label: "Ходов сделано", Cells: []any{1}},
							{Label: "Всего пар", Cells: []any{21}},
							{Label: "Доп. память: 3 значения", Cells: []any{"l", "r", "s", "O(1)"}},
						}},
						{Caption: "l = 1, r = 5: 4 + 12 = 16. Остальные суммы активной части столбца 5 не меньше 16. Следующий ход r -= 1 исключит его.", Rows: []sets.FigureRow{
							{Label: `i \ j`, Cells: []any{1, 2, 3, 4, 5, 6}, Dimmed: []int{5}},
							{Label: "0", Cells: []any{5, 6, 8, 11, 13, 17}, Dimmed: []int{0, 1, 2, 3, 4, 5}},
							{Label: "1", Cells: []any{"", 9, 11, 14, 16, 20}, Marked: []int{4}, Dimmed: []int{5}},
							{Label: "2", Cells: []any{"", "", 12, 15, 17, 21}, Dimmed: []int{5}},
							{Label: "3", Cells: []any{"", "", "", 17, 19, 23}, Dimmed: []int{5}},
							{Label: "4", Cells: []any{"", "", "", "", 22, 26}, Dimmed: []int{5}},
							{Label: "5", Cells: []any{"", "", "", "", "", 28}, Dimmed: []int{5}},
							{Label: "Сумм вычислено", Cells: []any{3}},
							{Label: "Ходов сделано", Cells: []any{2}},
							{Label: "Всего пар", Cells: []any{21}},
							{Label: "Доп. память: 3 значения", Cells: []any{"l", "r", "s", "O(1)"}},
						}},
						{Caption: "l = 1, r = 4: 4 + 10 = 14. Остальные суммы активной части строки 1 не больше 14. Следующий ход l += 1 исключит её.", Rows: []sets.FigureRow{
							{Label: `i \ j`, Cells: []any{1, 2, 3, 4, 5, 6}, Dimmed: []int{4, 5}},
							{Label: "0", Cells: []any{5, 6, 8, 11, 13, 17}, Dimmed: []int{0, 1, 2, 3, 4, 5}},
							{Label: "1", Cells: []any{"", 9, 11, 14, 16, 20}, Marked: []int{3}, Dimmed: []int{4, 5}},
							{Label: "2", Cells: []any{"", "", 12, 15, 17, 21}, Dimmed: []int{4, 5}},
							{Label: "3", Cells: []any{"", "", "", 17, 19, 23}, Dimmed: []int{4, 5}},
							{Label: "4", Cells: []any{"", "", "", "", 22, 26}, Dimmed: []int{4, 5}},
							{Label: "5", Cells: []any{"", "", "", "", "", 28}, Dimmed: []int{4, 5}},
							{Label: "Сумм вычислено", Cells: []any{4}},
							{Label: "Ходов сделано", Cells: []any{3}},
							{Label: "Всего пар", Cells: []any{21}},
							{Label: "Доп. память: 3 значения", Cells: []any{"l", "r", "s", "O(1)"}},
						}},
						{Caption: "l = 2, r = 4: 5 + 10 = 15. Получились 5 сумм и 4 хода. Ещё 14 пар исключены без вычисления, 2 не потребовались после успеха. Ответ - два индекса, O(1).", Rows: []sets.FigureRow{
							{Label: `i \ j`, Cells: []any{1, 2, 3, 4, 5, 6}, Dimmed: []int{4, 5}},
							{Label: "0", Cells: []any{5, 6, 8, 11, 13, 17}, Dimmed: []int{0, 1, 2, 3, 4, 5}},
							{Label: "1", Cells: []any{"", 9, 11, 14, 16, 20}, Dimmed: []int{0, 1, 2, 3, 4, 5}},
							{Label: "2", Cells: []any{"", "", 12, 15, 17, 21}, Marked: []int{3}, Dimmed: []int{4, 5}},
							{Label: "3", Cells: []any{"", "", "", 17, 19, 23}, Dimmed: []int{4, 5}},
							{Label: "4", Cells: []any{"", "", "", "", 22, 26}, Dimmed: []int{4, 5}},
							{Label: "5", Cells: []any{"", "", "", "", "", 28}, Dimmed: []int{4, 5}},
							{Label: "Сумм вычислено", Cells: []any{5}},
							{Label: "Ходов сделано", Cells: []any{4}},
							{Label: "Всего пар", Cells: []any{21}},
							{Label: "Доп. память: 3 значения", Cells: []any{"l", "r", "s", "O(1)"}},
						}},
					},
				},
				{
					ID:    "tp1-window",
					Title: "Two Pointers: окно и рабочая память",
					Frames: []sets.FigureFrame{
						{Caption: "Окно содержит весь Array. 1 + 16 = 17; любая другая пара с 16 внутри окна даёт не меньше 17. Следующий ход r -= 1.", Rows: []sets.FigureRow{
							{Label: "Array, n = 7", Cells: []any{1, 4, 5, 7, 10, 12, 16}, Marked: []int{0, 6}, Pointers: &sets.Pointers{{Name: "l", Cell: 0}, {Name: "r", Cell: 6}}, Window: &[2]int{0, 6}},
							{Label: "Сумм вычислено", Cells: []any{1}},
							{Label: "Ходов сделано", Cells: []any{0}},
							{Label: "Доп. память: 3 значения", Cells: []any{0, 6, 17}},
						}},
						{Caption: "1 + 12 = 13; любая другая пара с 1 внутри окна даёт не больше 13. Следующий ход l += 1. Рабочие значения в нижней строке: l, r, s.", Rows: []sets.FigureRow{
							{Label: "Array, n = 7", Cells: []any{1, 4, 5, 7, 10, 12, 16}, Marked: []int{0, 5}, Dimmed: []int{6}, Pointers: &sets.Pointers{{Name: "l", Cell: 0}, {Name: "r", Cell: 5}}, Window: &[2]int{0, 5}},
							{Label: "Сумм вычислено", Cells: []any{2}},
							{Label: "Ходов сделано", Cells: []any{1}},
							{Label: "Доп. память: 3 значения", Cells: []any{0, 5, 13}},
						}},
						{Caption: "4 + 12 = 16; любая другая активная пара с 12 даёт не меньше 16. Следующий ход r -= 1. Число рабочих значений не растёт.", Rows: []sets.FigureRow{
							{Label: "Array, n = 7", Cells: []any{1, 4, 5, 7, 10, 12, 16}, Marked: []int{1, 5}, Dimmed: []int{0, 6}, Pointers: &sets.Pointers{{Name: "l", Cell: 1}, {Name: "r", Cell: 5}}, Window: &[2]int{1, 5}},
							{Label: "Сумм вычислено", Cells: []any{3}},
							{Label: "Ходов сделано", Cells: []any{2}},
							{Label: "Доп. память: 3 значения", Cells: []any{1, 5, 16}},
						}},
						{Caption: "4 + 10 = 14; любая другая активная пара с 4 даёт не больше 14. Следующий ход l += 1. Копия окна не создаётся.", Rows: []sets.FigureRow{
							{Label: "Array, n = 7", Cells: []any{1, 4, 5, 7, 10, 12, 16}, Marked: []int{1, 4}, Dimmed: []int{0, 5, 6}, Pointers: &sets.Pointers{{Name: "l", Cell: 1}, {Name: "r", Cell: 4}}, Window: &[2]int{1, 4}},
							{Label: "Сумм вычислено", Cells: []any{4}},
							{Label: "Ходов сделано", Cells: []any{3}},
							{Label: "Доп. память: 3 значения", Cells: []any{1, 4, 14}},
						}},
						{Caption: "5 + 10 = 15: ответ [3,5] с единицы. Проверены 5 сумм из 21 возможной, сделаны 4 хода. Дополнительная память O(1); два индекса ответа тоже O(1).", Rows: []sets.FigureRow{
							{Label: "Array, n = 7", Cells: []any{1, 4, 5, 7, 10, 12, 16}, Marked: []int{2, 4}, Dimmed: []int{0, 1, 5, 6}, Pointers: &sets.Pointers{{Name: "l", Cell: 2}, {Name: "r", Cell: 4}}, Window: &[2]int{2, 4}},
							{Label: "Сумм вычислено", Cells: []any{5}},
							{Label: "Ходов сделано", Cells: []any{4}},
							{Label: "Доп. память: 3 значения", Cells: []any{2, 4, 15}},
						}},
					},
				},
			},
		},
		{
			Title:   "Читатель и писатель",
			Summary: "Two Pointers с ролями read и write сохраняет ответ в исходном Array. Вы считаете проверки и записи, включая обмены с собой, и отличаете рабочую память от памяти результата.",
			Lesson:  lesson2,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/remove-duplicates-from-sorted-array/",
					Number:     number(26),
					Title:      "Remove Duplicates from Sorted Array",
					Difficulty: "easy",
					Topics:     []string{"Array", "Two Pointers"},
					Note:       "n - длина отсортированного Array, k - число разных значений; числа ограниченного размера. Сравнивайте с nums[write-1]. Время: n-1 сравнений и k-1 записей, O(n). Память: O(1) дополнительной. Ответ занимает первые k входных ячеек; отдельная копия этого префикса для проверки занимает O(k) памяти результата и требует O(k) времени.",
					Statement: `<p>Дан Array целых чисел <code>nums</code>, который <strong>отсортирован по неубыванию</strong> (каждое следующее число не меньше предыдущего). Удалите повторы <strong>на месте</strong> (in-place), чтобы каждое число встречалось только <strong>один раз</strong>. <strong>Порядок</strong> элементов должен остаться <strong>тем же</strong>.</p><p>Пусть <code>k</code> это количество разных чисел в <code>nums</code>. После удаления повторов верните <code>k</code>.</p><p>Первые <code>k</code> элементов <code>nums</code> должны содержать разные числа в отсортированном порядке. Элементы после индекса <code>k - 1</code> не важны.</p><p><strong>Как проверяют решение.</strong> Проверяющая программа делает так:</p><ol><li>Вызывает ваш метод: <code>k = removeDuplicates(nums)</code>.</li><li>Проверяет, что <code>k</code> равно длине правильного ответа <code>expectedNums</code>.</li><li>Проверяет, что <code>nums[i] == expectedNums[i]</code> для каждого <code>i</code> от <code>0</code> до <code>k - 1</code>.</li></ol><p>Если все проверки прошли, решение <strong>принято</strong>.</p><p>Метод класса <code>Solution</code>: <code>def removeDuplicates(self, nums: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,1,2]
<strong>Выход:</strong> 2, nums = [1,2,_]
<strong>Пояснение:</strong> Метод должен вернуть k = 2, и первые два элемента nums должны быть 1 и 2.
Что лежит после первых k элементов, не важно (поэтому там стоят подчёркивания).</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [0,0,1,1,1,2,2,3,3,4]
<strong>Выход:</strong> 5, nums = [0,1,2,3,4,_,_,_,_,_]
<strong>Пояснение:</strong> Метод должен вернуть k = 5, и первые пять элементов nums должны быть 0, 1, 2, 3 и 4.
Что лежит после первых k элементов, не важно (поэтому там стоят подчёркивания).</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 3 · 10<sup>4</sup></code></li><li><code>-100 ≤ nums[i] ≤ 100</code></li><li><code>nums</code> отсортирован по неубыванию</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "removeDuplicates", Check: "prefix"}, Cases: []sets.TestCase{
						{Input: "[[1,1,2]]", Expected: text("[1,2]")},
						{Input: "[[0,0,1,1,1,2,2,3,3,4]]", Expected: text("[0,1,2,3,4]")},
						{Input: "[[7]]", Expected: text("[7]")},
						{Input: "[[-100,-100,-100]]", Expected: text("[-100]")},
						{Input: "[[-3,-1,0,2]]", Expected: text("[-3,-1,0,2]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/remove-element/",
					Number:     number(27),
					Title:      "Remove Element",
					Difficulty: "easy",
					Topics:     []string{"Array", "Two Pointers"},
					Note:       "n - длина Array, k - число элементов, не равных val; сравнение чисел стоит O(1). Время: n проверок и k записей, O(n). Память: O(1) дополнительной, ответ в первых k входных ячейках. Копия префикса требует O(k) памяти результата. При допустимой смене порядка замена значением с конца может уменьшить число записей. Повторные удаления из середины list могут занять O(n²) времени.",
					Statement: `<p>Дан Array целых чисел <code>nums</code> и целое число <code>val</code>. Удалите все вхождения <code>val</code> из <code>nums</code> <strong>на месте</strong> (in-place). Порядок элементов можно менять. Затем верните <em>количество элементов</em> <code>nums</code><em>, которые не равны</em> <code>val</code>.</p><p>Пусть <code>k</code> это количество элементов <code>nums</code>, которые не равны <code>val</code>. Чтобы решение приняли, сделайте так:</p><ul><li>Измените Array <code>nums</code> так, чтобы в первых <code>k</code> элементах лежали элементы, которые не равны <code>val</code>. Остальные элементы <code>nums</code> и размер <code>nums</code> не важны.</li><li>Верните <code>k</code>.</li></ul><p><strong>Как проверяют решение.</strong> Проверяющая программа делает так:</p><ol><li>Вызывает ваш метод: <code>k = removeElement(nums, val)</code>.</li><li>Проверяет, что <code>k</code> равно длине правильного ответа <code>expectedNums</code>. В <code>expectedNums</code> нет значений <code>val</code>, и он отсортирован.</li><li>Применяет Sorting к первым <code>k</code> элементам <code>nums</code>.</li><li>Проверяет, что <code>nums[i] == expectedNums[i]</code> для каждого <code>i</code> от <code>0</code> до <code>k - 1</code>.</li></ol><p>Если все проверки прошли, решение <strong>принято</strong>.</p><p>Метод класса <code>Solution</code>: <code>def removeElement(self, nums: list[int], val: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [3,2,2,3], val = 3
<strong>Выход:</strong> 2, nums = [2,2,_,_]
<strong>Пояснение:</strong> Метод должен вернуть k = 2, и первые два элемента nums должны быть равны 2.
Что лежит после первых k элементов, не важно (поэтому там стоят подчёркивания).</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [0,1,2,2,3,0,4,2], val = 2
<strong>Выход:</strong> 5, nums = [0,1,4,0,3,_,_,_]
<strong>Пояснение:</strong> Метод должен вернуть k = 5, и в первых пяти элементах nums должны лежать 0, 0, 1, 3 и 4.
Эти пять элементов могут идти в любом порядке.
Что лежит после первых k элементов, не важно (поэтому там стоят подчёркивания).</pre><h3>Ограничения</h3><ul><li><code>0 ≤ nums.length ≤ 100</code></li><li><code>0 ≤ nums[i] ≤ 50</code></li><li><code>0 ≤ val ≤ 100</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "removeElement", Check: "prefix", Order: "any"}, Cases: []sets.TestCase{
						{Input: "[[3,2,2,3],3]", Expected: text("[2,2]")},
						{Input: "[[0,1,2,2,3,0,4,2],2]", Expected: text("[0,0,1,3,4]")},
						{Input: "[[],0]", Expected: text("[]")},
						{Input: "[[1],1]", Expected: text("[]")},
						{Input: "[[4,5],3]", Expected: text("[4,5]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/move-zeroes/",
					Number:     number(283),
					Title:      "Move Zeroes",
					Difficulty: "easy",
					Topics:     []string{"Array", "Two Pointers"},
					Note:       "n - длина Array, k - число ненулевых элементов; числа ограниченного размера. Время: n проверок и k обменов, O(n); каждый обмен делает 2 записи, даже с собой. Память: O(1) дополнительной, результат остаётся во входном Array. Защита read != write убирает обмены с собой. Запись k ненулевых элементов и заполнение n-k нулями дают ровно n записей без такой защиты.",
					Statement: `<p>Дан Array целых чисел <code>nums</code>. Переместите все <code>0</code> в конец Array. Порядок ненулевых элементов должен остаться прежним.</p><p><strong>Обратите внимание:</strong> это нужно сделать на месте, без копии Array. Метод ничего не возвращает. В примерах «Выход» показывает <code>nums</code> после вызова.</p><p>Метод класса <code>Solution</code>: <code>def moveZeroes(self, nums: list[int]) -&gt; None</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [0,1,0,3,12]
<strong>Выход:</strong> [1,3,12,0,0]</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [0]
<strong>Выход:</strong> [0]</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 10<sup>4</sup></code></li><li><code>-2<sup>31</sup> ≤ nums[i] ≤ 2<sup>31</sup> - 1</code></li></ul><p><strong>Дополнительный вопрос:</strong> Можете ли вы сделать как можно меньше операций?</p>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "moveZeroes", Check: "arg0"}, Cases: []sets.TestCase{
						{Input: "[[0,1,0,3,12]]", Expected: text("[1,3,12,0,0]")},
						{Input: "[[0]]", Expected: text("[0]")},
						{Input: "[[1,2,3]]", Expected: text("[1,2,3]")},
						{Input: "[[0,0,1]]", Expected: text("[1,0,0]")},
						{Input: "[[-1,0,0,-2147483648,2147483647]]", Expected: text("[-1,-2147483648,2147483647,0,0]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/remove-duplicates-from-sorted-array-ii/",
					Number:     number(80),
					Title:      "Remove Duplicates from Sorted Array II",
					Difficulty: "medium",
					Topics:     []string{"Array", "Two Pointers"},
					Note:       "n - длина отсортированного Array, k - длина ответа; сравнение чисел стоит O(1). Первые два элемента оставьте, затем сравнивайте с nums[write-2]. Время: при n≥2 это n-2 сравнений и k-2 записей, O(n). Память: O(1) дополнительной; первые k входных ячеек содержат ответ. Его копия требует O(k) памяти результата. Hash Table частот добавляет O(u) памяти для u разных чисел, хотя порядок позволяет обойтись без неё.",
					Statement: `<p>Дан Array целых чисел <code>nums</code>, который <strong>отсортирован по неубыванию</strong> (каждое следующее число не меньше предыдущего). Удалите часть повторов <strong>на месте</strong> (in-place), чтобы каждое число встречалось <strong>не больше двух раз</strong>. <strong>Порядок</strong> элементов должен остаться <strong>тем же</strong>.</p><p>В некоторых языках длину Array менять нельзя. Поэтому результат должен лежать в <strong>первой части</strong> Array <code>nums</code>. Точнее: если после удаления повторов осталось <code>k</code> элементов, то первые <code>k</code> элементов <code>nums</code> должны содержать итоговый результат. Что лежит после первых <code>k</code> элементов, не важно.</p><p>Поместите итоговый результат в первые <code>k</code> ячеек <code>nums</code> и верните <code>k</code>.</p><p><strong>Не</strong> создавайте другой Array. Меняйте входной Array <strong>на месте</strong> и используйте O(1) дополнительной памяти.</p><p><strong>Как проверяют решение.</strong> Проверяющая программа делает так:</p><ol><li>Вызывает ваш метод: <code>k = removeDuplicates(nums)</code>.</li><li>Проверяет, что <code>k</code> равно длине правильного ответа <code>expectedNums</code>.</li><li>Проверяет, что <code>nums[i] == expectedNums[i]</code> для каждого <code>i</code> от <code>0</code> до <code>k - 1</code>.</li></ol><p>Если все проверки прошли, решение <strong>принято</strong>.</p><p>Метод класса <code>Solution</code>: <code>def removeDuplicates(self, nums: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,1,1,2,2,3]
<strong>Выход:</strong> 5, nums = [1,1,2,2,3,_]
<strong>Пояснение:</strong> Метод должен вернуть k = 5, и первые пять элементов nums должны быть 1, 1, 2, 2 и 3.
Что лежит после первых k элементов, не важно (поэтому там стоят подчёркивания).</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [0,0,1,1,1,1,2,3,3]
<strong>Выход:</strong> 7, nums = [0,0,1,1,2,3,3,_,_]
<strong>Пояснение:</strong> Метод должен вернуть k = 7, и первые семь элементов nums должны быть 0, 0, 1, 1, 2, 3 и 3.
Что лежит после первых k элементов, не важно (поэтому там стоят подчёркивания).</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 3 · 10<sup>4</sup></code></li><li><code>-10<sup>4</sup> ≤ nums[i] ≤ 10<sup>4</sup></code></li><li><code>nums</code> отсортирован по неубыванию</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "removeDuplicates", Check: "prefix"}, Cases: []sets.TestCase{
						{Input: "[[1,1,1,2,2,3]]", Expected: text("[1,1,2,2,3]")},
						{Input: "[[0,0,1,1,1,1,2,3,3]]", Expected: text("[0,0,1,1,2,3,3]")},
						{Input: "[[5]]", Expected: text("[5]")},
						{Input: "[[1,1]]", Expected: text("[1,1]")},
						{Input: "[[-2,-2,-2,-2,0,1,1,1]]", Expected: text("[-2,-2,0,1,1]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/sort-colors/",
					Number:     number(75),
					Title:      "Sort Colors",
					Difficulty: "medium",
					Topics:     []string{"Array", "Two Pointers", "Sorting"},
					Note:       "n - длина Array со значениями 0, 1, 2. Dutch National Flag: low, mid, high ограничивают неизвестную часть. Время: она сокращается на 1 в каждой из n итераций, O(n). Память: O(1) дополнительной, результат во входном Array. Counting Sort с тремя счётчиками тоже даёт O(n) времени и O(1) памяти за два прохода. Общий Python Sorting требует O(n log n) времени и до O(n) рабочей памяти.",
					Statement: `<p>Дан Array <code>nums</code> из <code>n</code> предметов. Каждый предмет красный, белый или синий. Расставьте их <strong>на месте</strong> (in-place) так, чтобы предметы одного цвета шли подряд, а цвета шли в порядке: красный, белый, синий.</p><p>Цвета обозначены целыми числами: 0 это красный, 1 это белый, 2 это синий.</p><p>Решите задачу без библиотечной функции Sorting. Метод ничего не возвращает. В примерах «Выход» показывает <code>nums</code> после вызова.</p><p>Метод класса <code>Solution</code>: <code>def sortColors(self, nums: list[int]) -&gt; None</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [2,0,2,1,1,0]
<strong>Выход:</strong> [0,0,1,1,2,2]
<strong>Пояснение:</strong> В Array два 0, две 1 и две 2. После Sorting на месте сначала идут все 0, потом все 1, потом все 2.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [2,0,1]
<strong>Выход:</strong> [0,1,2]
<strong>Пояснение:</strong> В Array по одному 0, 1 и 2. После Sorting на месте они идут в порядке 0, 1, 2.</pre><h3>Ограничения</h3><ul><li><code>n == nums.length</code></li><li><code>1 ≤ n ≤ 300</code></li><li><code>nums[i]</code> равно <code>0</code>, <code>1</code> или <code>2</code></li></ul><p><strong>Дополнительный вопрос:</strong> Можете ли вы придумать алгоритм, который делает один проход по Array и использует только постоянную дополнительную память?</p>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "sortColors", Check: "arg0"}, Cases: []sets.TestCase{
						{Input: "[[2,0,2,1,1,0]]", Expected: text("[0,0,1,1,2,2]")},
						{Input: "[[2,0,1]]", Expected: text("[0,1,2]")},
						{Input: "[[1]]", Expected: text("[1]")},
						{Input: "[[2,2,2]]", Expected: text("[2,2,2]")},
						{Input: "[[1,0]]", Expected: text("[0,1]")},
						{Input: "[[2,1,0,0,1,2,1]]", Expected: text("[0,0,1,1,1,2,2]")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "tp2-zones", Type: "choice", Prompt: "LeetCode 27, val = 2, nums = [2, 5, 2, 7, 9], в начале write = 0. Обработаны read = 0, 1, 2, 3; теперь read = 4. Что лежит между write и read, то есть в nums[write:read]?", Options: []string{"[2, 2]", "[5, 7]", "[2]", "[2, 7]"}, Answer: 3, Explain: "После 4 проверок и 2 записей nums = [5,7,2,7,9], write = 2, read = 4. В диапазоне nums[2:4] лежат 2 и старая копия 7; [5,7] - готовый ответ. Здесь slice обозначает диапазон, код не создаёт его копию. Для n чисел с O(1) сравнением полный проход требует O(n) времени и O(1) дополнительной памяти; результат остаётся во входе."},
				{ID: "tp2-no-overtake", Type: "code", Prompt: "LeetCode 27, порядок можно менять. Ученик пишет нужные числа с конца Array, а read идёт слева направо. Каким станет nums = [5, 6, 7] при val = 0?", Language: "python", Code: `write = len(nums) - 1
for read in range(len(nums)):
    if nums[read] != val:
        nums[write] = nums[read]
        write -= 1`, Options: []string{"[7, 6, 5]", "[5, 6, 7]", "[5, 6, 5]"}, Answer: 2, Explain: "При read = 0 значение 5 записано в ячейку 2 и уничтожило непрочитанное 7. При read = 2 код читает уже 5. Сделаны 3 проверки и 3 записи, но ответ неверен. Для n чисел с O(1) сравнением время O(n), дополнительная память O(1). Условие write ≤ read защищает непрочитанные данные в правильном шаблоне."},
				{ID: "tp2-what-changes", Type: "code", Prompt: "Для Move Zeroes (283) в шаблон Two Pointers подставили условие «число не равно 0». Что будет в nums после цикла, если nums = [0, 1, 0, 3]?", Language: "python", Code: `write = 0
for read in range(len(nums)):
    if nums[read] != 0:
        nums[write] = nums[read]
        write += 1`, Options: []string{"[1, 3, 0, 3]", "[1, 3, 0, 0]", "[1, 3]"}, Answer: 0, Explain: "Код делает 4 проверки и 2 записи: 1 в ячейку 0, затем 3 в ячейку 1. Ячейки 2 и 3 содержат прежние 0 и 3. Присваивание не меняет длину list. Допишите нули в конец или используйте обмены. Для n чисел с O(1) сравнением оба исправления требуют O(n) времени, O(1) дополнительной памяти и хранят результат во входе."},
				{ID: "tp2-code-indent", Type: "code", Prompt: "LeetCode 27. Что вернёт этот код на nums = [3, 1, 3, 3, 2], val = 3?", Language: "python", Code: `write = 0
for read in range(len(nums)):
    if nums[read] != val:
        nums[write] = nums[read]
    write += 1
return write`, Options: []string{"2", "5", "3"}, Answer: 1, Explain: "write += 1 стоит вне if, поэтому до каждой проверки write == read. Код делает 5 проверок, 2 записи в себя и возвращает 5. Правильный код увеличивает write только после записи и вернёт 2. Оба цикла для n чисел с O(1) сравнением требуют O(n) времени и O(1) дополнительной памяти; неправильный отступ нарушает смысл готового префикса."},
				{ID: "tp2-code-dedup", Type: "code", Prompt: "LeetCode 26, nums = [1, 2, 3]. Что получится после этого кода?", Language: "python", Code: `write = 1
for read in range(1, len(nums)):
    if nums[read] != nums[write]:
        nums[write] = nums[read]
        write += 1
return write`, Options: []string{"k = 3, в начале [1, 2, 3]", "k = 2, в начале [1, 2]", "k = 2, в начале [1, 3]"}, Answer: 2, Explain: "При read = 1 код сравнивает 2 с собой и пропускает его. При read = 2 число 3 отличается от nums[1] = 2 и записывается в ячейку 1. Итого 2 сравнения, 1 запись, неверный ответ [1,3]. Сравнивать нужно с nums[write-1]. Исправленный проход для n ограниченных по размеру чисел требует n-1 сравнений, O(n) времени и O(1) дополнительной памяти; ответ во входе."},
				{ID: "tp2-dedup-count", Type: "short", Prompt: "LeetCode 26, nums = [2, 2, 3, 5, 5, 5, 8, 9], в начале write = 1. Цикл уже обработал read = 1, 2, 3, 4, 5, 6. Чему равно write?", Answer: "4", Explain: "В прочитанной части четыре разных числа: 2, 3, 5, 8. Значит, write = 4 после 6 сравнений и 3 записей. Первое число не переписывали. Полный проход по n = 8 элементам сделает 7 сравнений. Для чисел с O(1) сравнением время O(n), дополнительная память O(1), ответ в исходном Array."},
				{ID: "tp2-self-write", Type: "short", Prompt: "Шаблон Two Pointers (write = 0) с условием nums[read] != 0 запускают на nums = [1, 2, 0, 3, 4]. Сколько раз строка nums[write] = nums[read] пишет число в ту же ячейку, то есть при write == read?", Answer: "2", Explain: "При read = 0 и read = 1 число пишется в свою ячейку. Ноль пропущен; затем 3 и 4 пишутся в ячейки 2 и 3. Итого 5 проверок и 4 записи, из них 2 в себя. Проверка read != write уберёт эти 2 записи. Для n чисел ограниченного размера время O(n), дополнительная память O(1); готовый префикс хранится во входе."},
			},
			Links: []sets.Link{
				{Title: "Two Pointers: краткий конспект", URL: "https://server.179.ru/wiki/files/global/two_pointers.pdf", Source: "Школа 179"},
				{Title: "Подходы к алгоритмическим задачам: Two Pointers", URL: "https://habr.com/ru/companies/avito/articles/806711/", Source: "Хабр"},
				{Title: "Two Pointers в Array: быстрый и медленный указатели", URL: "https://labuladong.online/en/algo/essential-technique/array-two-pointers-summary/", Source: "labuladong (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "tp2-dedup",
					Title: "LeetCode 26: 4 сравнения, 2 записи",
					Frames: []sets.FigureFrame{
						{Caption: "После read = 1: число 1 совпало с первым оставленным. Повтор пропущен, write = 1. Счётчики показывают выполненные сравнения и записи, а не номер кадра.", Rows: []sets.FigureRow{
							{Label: "Array, n = 5", Cells: []any{1, 1, 2, 3, 3}, Marked: []int{0, 1}, Pointers: &sets.Pointers{{Name: "write", Cell: 1}, {Name: "read", Cell: 1}}, Window: &[2]int{0, 0}},
							{Label: "Сравнений чисел", Cells: []any{1}},
							{Label: "Записей в Array", Cells: []any{0}},
							{Label: "Доп. память: 2 индекса", Cells: []any{"read", "write", "O(1)"}},
						}},
						{Caption: "После read = 2: число 2 отличается от последнего оставленного 1. Число 2 записано в ячейку 1; write = 2.", Rows: []sets.FigureRow{
							{Label: "Array, n = 5", Cells: []any{1, 2, 2, 3, 3}, Marked: []int{1, 2}, Pointers: &sets.Pointers{{Name: "write", Cell: 2}, {Name: "read", Cell: 2}}, Window: &[2]int{0, 1}},
							{Label: "Сравнений чисел", Cells: []any{2}},
							{Label: "Записей в Array", Cells: []any{1}},
							{Label: "Доп. память: 2 индекса", Cells: []any{"read", "write", "O(1)"}},
						}},
						{Caption: "После read = 3: число 3 отличается от последнего оставленного 2. Число 3 записано в ячейку 2; write = 3.", Rows: []sets.FigureRow{
							{Label: "Array, n = 5", Cells: []any{1, 2, 3, 3, 3}, Marked: []int{2, 3}, Pointers: &sets.Pointers{{Name: "write", Cell: 3}, {Name: "read", Cell: 3}}, Window: &[2]int{0, 2}},
							{Label: "Сравнений чисел", Cells: []any{3}},
							{Label: "Записей в Array", Cells: []any{2}},
							{Label: "Доп. память: 2 индекса", Cells: []any{"read", "write", "O(1)"}},
						}},
						{Caption: "После read = 4: число 3 совпало с последним оставленным. Новых записей нет. Готовый ответ содержит 3 входные ячейки.", Rows: []sets.FigureRow{
							{Label: "Array, n = 5", Cells: []any{1, 2, 3, 3, 3}, Marked: []int{2, 4}, Dimmed: []int{3}, Pointers: &sets.Pointers{{Name: "write", Cell: 3}, {Name: "read", Cell: 4}}, Window: &[2]int{0, 2}},
							{Label: "Сравнений чисел", Cells: []any{4}},
							{Label: "Записей в Array", Cells: []any{2}},
							{Label: "Доп. память: 2 индекса", Cells: []any{"read", "write", "O(1)"}},
						}},
						{Caption: "Готово: k = 3, ответ [1,2,3] в исходном Array. n-1 = 4 сравнения, k-1 = 2 записи. Дополнительная память O(1); копия ответа для проверки заняла бы 3 ячейки, O(k).", Rows: []sets.FigureRow{
							{Label: "Array, n = 5", Cells: []any{1, 2, 3, 3, 3}, Dimmed: []int{3, 4}, Pointers: &sets.Pointers{{Name: "write", Cell: 3}}, Window: &[2]int{0, 2}},
							{Label: "Сравнений чисел", Cells: []any{4}},
							{Label: "Записей в Array", Cells: []any{2}},
							{Label: "Доп. память: 2 индекса", Cells: []any{"read", "write", "O(1)"}},
						}},
					},
				},
				{
					ID:    "tp2-remove",
					Title: "LeetCode 27: 4 проверки условия, 2 записи",
					Frames: []sets.FigureFrame{
						{Caption: "После read = 0: условие nums[read] != 3 ложно. Записей нет, write = 0.", Rows: []sets.FigureRow{
							{Label: "Array, n = 4", Cells: []any{3, 2, 2, 3}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "write", Cell: 0}, {Name: "read", Cell: 0}}},
							{Label: "Проверок != val", Cells: []any{1}},
							{Label: "Записей в Array", Cells: []any{0}},
							{Label: "Доп. память: 2 индекса", Cells: []any{"read", "write", "O(1)"}},
						}},
						{Caption: "После read = 1: условие 2 != 3 истинно. Число 2 записано в ячейку 0; write = 1.", Rows: []sets.FigureRow{
							{Label: "Array, n = 4", Cells: []any{2, 2, 2, 3}, Marked: []int{0, 1}, Pointers: &sets.Pointers{{Name: "write", Cell: 1}, {Name: "read", Cell: 1}}, Window: &[2]int{0, 0}},
							{Label: "Проверок != val", Cells: []any{2}},
							{Label: "Записей в Array", Cells: []any{1}},
							{Label: "Доп. память: 2 индекса", Cells: []any{"read", "write", "O(1)"}},
						}},
						{Caption: "После read = 2: условие 2 != 3 истинно. Число 2 записано в ячейку 1; write = 2.", Rows: []sets.FigureRow{
							{Label: "Array, n = 4", Cells: []any{2, 2, 2, 3}, Marked: []int{1, 2}, Pointers: &sets.Pointers{{Name: "write", Cell: 2}, {Name: "read", Cell: 2}}, Window: &[2]int{0, 1}},
							{Label: "Проверок != val", Cells: []any{3}},
							{Label: "Записей в Array", Cells: []any{2}},
							{Label: "Доп. память: 2 индекса", Cells: []any{"read", "write", "O(1)"}},
						}},
						{Caption: "После read = 3: условие 3 != 3 ложно. Записей по-прежнему 2, write = 2.", Rows: []sets.FigureRow{
							{Label: "Array, n = 4", Cells: []any{2, 2, 2, 3}, Marked: []int{3}, Dimmed: []int{2}, Pointers: &sets.Pointers{{Name: "write", Cell: 2}, {Name: "read", Cell: 3}}, Window: &[2]int{0, 1}},
							{Label: "Проверок != val", Cells: []any{4}},
							{Label: "Записей в Array", Cells: []any{2}},
							{Label: "Доп. память: 2 индекса", Cells: []any{"read", "write", "O(1)"}},
						}},
						{Caption: "Готово: k = 2, в начале [2,2]. Выполнены n = 4 проверки условия и k = 2 записи. Дополнительная память O(1), ответ использует вход. Отдельная копия ответа заняла бы O(k) памяти.", Rows: []sets.FigureRow{
							{Label: "Array, n = 4", Cells: []any{2, 2, 2, 3}, Dimmed: []int{2, 3}, Pointers: &sets.Pointers{{Name: "write", Cell: 2}}, Window: &[2]int{0, 1}},
							{Label: "Проверок != val", Cells: []any{4}},
							{Label: "Записей в Array", Cells: []any{2}},
							{Label: "Доп. память: 2 индекса", Cells: []any{"read", "write", "O(1)"}},
						}},
					},
				},
			},
		},
		{
			Title:   "Две последовательности",
			Summary: "Two Pointers по двум Array или String: Merge, подсчёт, Subsequence и пересечение. Считаем сравнения, записи, подготовку данных и временные копии.",
			Lesson:  lesson3,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/merge-sorted-array/",
					Number:     number(88),
					Title:      "Merge Sorted Array",
					Difficulty: "easy",
					Topics:     []string{"Array", "Two Pointers", "Sorting"},
					Note:       "m и n - числа элементов двух отсортированных входов. Заполняйте nums1 с конца и остановитесь, когда закончится nums2. Время: O(m + n), не больше m + n записей. Память: O(1) дополнительная; вход O(m + n), ответ занимает готовые ячейки nums1. Отдельный результат прямого Merge требует O(m + n) ячеек.",
					Statement: `<p>Даны два Array целых чисел <code>nums1</code> и <code>nums2</code>. Каждый <strong>отсортирован по неубыванию</strong> (каждое следующее число не меньше предыдущего). Ещё даны <code>m</code> и <code>n</code>: числа значимых элементов в <code>nums1</code> и <code>nums2</code>.</p><p>Выполните <strong>Merge</strong> двух входов в один Array, отсортированный по неубыванию.</p><p>Запишите результат в <code>nums1</code>. Длина <code>nums1</code> равна <code>m + n</code>: первые <code>m</code> элементов содержат входные числа, а последние <code>n</code> элементов равны <code>0</code> и служат свободными ячейками. Длина <code>nums2</code> равна <code>n</code>.</p><p>Метод ничего не возвращает. В примерах «Выход» показывает <code>nums1</code> после вызова.</p><p>Метод класса <code>Solution</code>: <code>def merge(self, nums1: list[int], m: int, nums2: list[int], n: int) -&gt; None</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums1 = [1,2,3,0,0,0], m = 3, nums2 = [2,5,6], n = 3
<strong>Выход:</strong> [1,2,2,3,5,6]
<strong>Пояснение:</strong> Merge двух Array [1,2,3] и [2,5,6].
Элементы с индексами 0, 1 и 3 пришли из nums1.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums1 = [1], m = 1, nums2 = [], n = 0
<strong>Выход:</strong> [1]
<strong>Пояснение:</strong> Второй Array пуст. Первый уже содержит результат.</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> nums1 = [0], m = 0, nums2 = [1], n = 1
<strong>Выход:</strong> [1]
<strong>Пояснение:</strong> При m = 0 в nums1 нет значимых элементов. Ноль обозначает свободную ячейку для числа 1.</pre><h3>Ограничения</h3><ul><li><code>nums1.length == m + n</code></li><li><code>nums2.length == n</code></li><li><code>0 ≤ m, n ≤ 200</code></li><li><code>1 ≤ m + n ≤ 200</code></li><li><code>-10<sup>9</sup> ≤ nums1[i], nums2[j] ≤ 10<sup>9</sup></code></li></ul><p><strong>Дополнительный вопрос:</strong> Можете ли вы выполнить Merge за <code>O(m + n)</code> времени?</p>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "merge", Check: "arg0"}, Cases: []sets.TestCase{
						{Input: "[[1,2,3,0,0,0],3,[2,5,6],3]", Expected: text("[1,2,2,3,5,6]")},
						{Input: "[[1],1,[],0]", Expected: text("[1]")},
						{Input: "[[0],0,[1],1]", Expected: text("[1]")},
						{Input: "[[4,5,6,0,0,0],3,[1,2,3],3]", Expected: text("[1,2,3,4,5,6]")},
						{Input: "[[-5,0,0,0,0],1,[-5,-5,2,3],4]", Expected: text("[-5,-5,-5,2,3]")},
						{Input: "[[0,0,3,0,0,0],3,[-1,1,1],3]", Expected: text("[-1,0,0,1,1,3]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/is-subsequence/",
					Number:     number(392),
					Title:      "Is Subsequence",
					Difficulty: "easy",
					Topics:     []string{"String", "Two Pointers"},
					Note:       "n = len(s), m = len(t). Двигайте i по s при совпадении, j по t всегда. Время: O(m), до m сравнений, с ранним успехом при i = n. Память: O(1) дополнительная и для ответа bool; входные String занимают O(n + m). Для многих запросов сохраните позиции букв t за O(m) и используйте Binary Search.",
					Statement: `<p>Даны две String <code>s</code> и <code>t</code>. Верните <code>True</code>, если <code>s</code> является <strong>Subsequence</strong> от <code>t</code>, иначе верните <code>False</code>.</p><p><strong>Subsequence</strong> получают удалением символов исходной String без изменения порядка оставшихся. Можно не удалять ни одного символа. Например, <code>"ace"</code> является Subsequence от <code>"abcde"</code>, а <code>"aec"</code> не является.</p><p>Метод класса <code>Solution</code>: <code>def isSubsequence(self, s: str, t: str) -&gt; bool</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> s = "abc", t = "ahbgdc"
<strong>Выход:</strong> True</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> s = "axc", t = "ahbgdc"
<strong>Выход:</strong> False</pre><h3>Ограничения</h3><ul><li><code>0 ≤ s.length ≤ 100</code></li><li><code>0 ≤ t.length ≤ 10<sup>4</sup></code></li><li><code>s</code> и <code>t</code> состоят только из строчных английских букв</li></ul><p><strong>Дополнительный вопрос:</strong> Пусть String <code>s</code> много: <code>s<sub>1</sub>, s<sub>2</sub>, ..., s<sub>k</sub></code>, где <code>k ≥ 10<sup>9</sup></code>. Нужно по очереди проверить, является ли каждая из них Subsequence от <code>t</code>. Как изменить код?</p>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "isSubsequence"}, Cases: []sets.TestCase{
						{Input: `["abc","ahbgdc"]`, Expected: text("True")},
						{Input: `["axc","ahbgdc"]`, Expected: text("False")},
						{Input: `["",""]`, Expected: text("True")},
						{Input: `["a",""]`, Expected: text("False")},
						{Input: `["aab","ab"]`, Expected: text("False")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/merge-strings-alternately/",
					Number:     number(1768),
					Title:      "Merge Strings Alternately",
					Difficulty: "easy",
					Topics:     []string{"String", "Two Pointers"},
					Note:       "n и m - длины word1 и word2. Добавляйте буквы по очереди в list, затем используйте join. Время: O(n + m), n + m добавлений и столько же символов при join. Память: O(n + m) дополнительная для list; вход и итоговая String отдельно занимают O(n + m). Повторное сложение String может копировать растущие префиксы за O((n + m)²).",
					Statement: `<p>Даны две String <code>word1</code> и <code>word2</code>. Берите буквы по очереди, начиная с <code>word1</code>. Если одна String длиннее другой, допишите её оставшиеся буквы в конец результата.</p><p>Верните <em>итоговую String</em>.</p><p>Метод класса <code>Solution</code>: <code>def mergeAlternately(self, word1: str, word2: str) -&gt; str</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> word1 = "abc", word2 = "pqr"
<strong>Выход:</strong> "apbqcr"
<strong>Пояснение:</strong> Буквы чередуются так:
word1:  a   b   c
word2:    p   q   r
итог:   a p b q c r</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> word1 = "ab", word2 = "pqrs"
<strong>Выход:</strong> "apbqrs"
<strong>Пояснение:</strong> word2 длиннее, поэтому "rs" дописывается в конец.
word1:  a   b
word2:    p   q   r   s
итог:   a p b q   r   s</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> word1 = "abcd", word2 = "pq"
<strong>Выход:</strong> "apbqcd"
<strong>Пояснение:</strong> word1 длиннее, поэтому "cd" дописывается в конец.
word1:  a   b   c   d
word2:    p   q
итог:   a p b q c   d</pre><h3>Ограничения</h3><ul><li><code>1 ≤ word1.length, word2.length ≤ 100</code></li><li><code>word1</code> и <code>word2</code> состоят из строчных английских букв</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "mergeAlternately"}, Cases: []sets.TestCase{
						{Input: `["abc","pqr"]`, Expected: text(`"apbqcr"`)},
						{Input: `["ab","pqrs"]`, Expected: text(`"apbqrs"`)},
						{Input: `["abcd","pq"]`, Expected: text(`"apbqcd"`)},
						{Input: `["a","b"]`, Expected: text(`"ab"`)},
						{Input: `["x","yzzz"]`, Expected: text(`"xyzzz"`)},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/intersection-of-two-arrays/",
					Number:     number(349),
					Title:      "Intersection of Two Arrays",
					Difficulty: "easy",
					Topics:     []string{"Array", "Two Pointers", "Sorting"},
					Note:       "n и m - длины входных Array, k - число разных общих значений. После Sorting двигайте указатель у меньшего числа. При равенстве сравнивайте число только с последним в ответе. Время: O(n log n + m log m + n + m). Память: O(n + m) дополнительная для Python Sorting, вход O(n + m), ответ O(k). Hash Table через set даёт ожидаемое O(n + m) время.",
					Statement: `<p>Даны два Array целых чисел <code>nums1</code> и <code>nums2</code>. Верните <em>Array их пересечения</em>. Каждый элемент в ответе должен быть <strong>уникальным</strong>. Порядок в ответе <strong>любой</strong>.</p><p>Пересечение двух Array содержит все значения, которые есть в обоих.</p><p>Метод класса <code>Solution</code>: <code>def intersection(self, nums1: list[int], nums2: list[int]) -&gt; list[int]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums1 = [1,2,2,1], nums2 = [2,2]
<strong>Выход:</strong> [2]</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums1 = [4,9,5], nums2 = [9,4,9,8,4]
<strong>Выход:</strong> [9,4]
<strong>Пояснение:</strong> Ответ [4,9] тоже правильный.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums1.length, nums2.length ≤ 1000</code></li><li><code>0 ≤ nums1[i], nums2[i] ≤ 1000</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "intersection", Order: "any"}, Cases: []sets.TestCase{
						{Input: "[[1,2,2,1],[2,2]]", Expected: text("[2]")},
						{Input: "[[4,9,5],[9,4,9,8,4]]", Expected: text("[4,9]")},
						{Input: "[[1,3],[2,4]]", Expected: text("[]")},
						{Input: "[[0],[0]]", Expected: text("[0]")},
						{Input: "[[1000,0,7,7],[7,1000,1000]]", Expected: text("[7,1000]")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "tp3-one-compare", Type: "short", Prompt: "Merge из главы на каждом шаге берёт меньшее из a[i] и b[j]. Какое число он выдаст третьим, если a = [4, 1, 6], b = [2, 5]?", Answer: "1", Explain: "4 и 2: берём 2. 4 и 5: берём 4. 1 и 5: берём 1. Третьим выходит 1, сразу после 4. Ответ [2, 4, 1, 5, 6] не отсортирован. Merge требует готового порядка: a[i] должно быть наименьшим оставшимся числом a. Линейное число движений не исправляет нарушение этого условия."},
				{ID: "tp3-code-tail", Type: "code", Prompt: "В последней строке этого Merge ошибка. На каком входе код выдаст неверный ответ?", Language: "python", Code: `i = j = 0
out = []
while i < len(a) and j < len(b):
    if a[i] <= b[j]:
        out.append(a[i])
        i += 1
    else:
        out.append(b[j])
        j += 1
out += b[j:]`, Options: []string{"a = [5], b = [1, 2]", "a = [1, 5], b = [1, 5]", "a = [], b = [3, 4]", "a = [1, 2], b = [5]"}, Answer: 0, Explain: "При a = [5], b = [1, 2] первым закончится b: код возьмёт 1 и 2, а 5 останется в a. Строка out += b[j:] допишет пустой list. Нужно дописать оба остатка: out += a[i:] + b[j:]. На остальных входах первым заканчивается a. Срезы и их объединение требуют O(n + m) временной памяти в худшем случае, где n и m - длины входов."},
				{ID: "tp3-loop-count", Type: "short", Prompt: "Merge из главы: a = [1, 7, 8, 9], b = [2, 3]. Сколько раз выполнится тело цикла while?", Answer: "3", Explain: "Берём 1 из a, потом 2 и 3 из b. Это три итерации, три сравнения значений и три сдвига: i = 1, j = 2. Затем срез [7, 8, 9] и объединение остатков копируют по три элемента, а out получает ещё три записи. Итог: шесть записей в out; копирование остатков не двигает i и j. Одна строка Python не означает O(1) работы."},
				{ID: "tp3-count-smaller", Type: "short", Prompt: "Через Merge считаем, сколько чисел в a строго меньше b[j]. В момент, когда берём b[2] = 10, в out уже лежат [1, 2, 4, 5, 7, 9]. Какой ответ для 10?", Answer: "4", Explain: "Ответ для b[j] равен i. В out лежат i чисел из a и j чисел из b, всего i + j = 6. Из b уже взяли два числа, значит i = 6 - 2 = 4. Сам out для подсчёта не нужен: достаточно индексов и места для m ответов, где m = len(b)."},
				{ID: "tp3-equal-first", Type: "code", Prompt: "Для каждого b[j] нужно найти, сколько чисел в a строго меньше него. Что выдаст код для a = [1, 3, 3, 5] и b = [3, 5]?", Language: "python", Code: `i = j = 0
res = []
while i < len(a) and j < len(b):
    if a[i] <= b[j]:
        i += 1
    else:
        res.append(i)
        j += 1
res += [i] * (len(b) - j)`, Options: []string{"[2, 3]", "[3, 4]", "[1, 3]", "[3]"}, Answer: 1, Explain: "Код выдаст [3, 4]: условие <= учитывает равные числа. Правильный ответ [1, 3] требует a[i] < b[j]. n = len(a), m = len(b): время O(n + m), ответ O(m). Выражение [i] * (len(b) - j) создаёт временный list до m ячеек. Поэлементное добавление оставшихся ответов уменьшает дополнительную память до O(1)."},
				{ID: "tp3-merge-back", Type: "short", Prompt: "LeetCode 88: nums1 = [2, 5, 7, 0, 0], m = 3, nums2 = [3, 6]. Петя выполняет Merge слева направо прямо в nums1: i, j и write начинают с 0. Каждый шаг пишет меньшее число в nums1[write]. Какое непрочитанное число из nums1 пропадёт первым?", Answer: "5", Explain: "Сначала записываем 2 в ячейку 0. Затем 3 меньше 5: запись 3 в ячейку 1 удалит ещё не использованное число 5. При обратном Merge write = i + j + 1, поэтому при j >= 0 запись идёт правее i. Нужны только три индекса, O(1) дополнительной памяти."},
				{ID: "tp3-subsequence", Type: "choice", Prompt: "LeetCode 392. Для какой пары String s не является Subsequence от t?", Options: []string{`s = "abb", t = "babab"`, `s = "", t = "xyz"`, `s = "aba", t = "aabb"`, `s = "aa", t = "baba"`}, Answer: 2, Explain: `Для "aba" и "aabb": a найдена в t[0], b в t[2], дальше второй a нет. Наличие букв не заменяет их порядок. "abb" есть в "babab" на позициях 1, 2, 4; "aa" есть в "baba" на позициях 1, 3. Пустая s является Subsequence любой t. Проверка идёт по t один раз, до m сравнений при m = len(t), O(1) дополнительной памяти.`},
				{ID: "tp3-common-move", Type: "choice", Prompt: "Ищем общие числа отсортированных Array a = [3, 7, 9, 12] и b = [3, 5, 9, 12]. i и j стоят в начале. При равенстве записываем число и двигаем оба указателя. На втором сравнении Петя сдвинул не тот указатель. Дальше он всё делает верно. Какие общие числа он найдёт?", Options: []string{"3, 9 и 12", "Только 3", "3 и 12"}, Answer: 0, Explain: "Пары значений: (3,3), (7,5), затем после ошибочного сдвига i: (9,5), (9,9), (12,12). Это пять проверок пар и три записи. Пропущенное число 7 отсутствует в b, поэтому ответ случайно верен. На другом входе такой ход может потерять общее число. Для уникальности сравнивайте найденное значение с последним в ответе, без поиска по всему ответу."},
			},
			Links: []sets.Link{
				{Title: "Sorting: Merge двух отсортированных Array", URL: "https://algorithmica.org/ru/sorting", Source: "Алгоритмика"},
				{Title: "Merge Sort", URL: "https://ru.wikipedia.org/wiki/Сортировка_слиянием", Source: "Википедия"},
				{Title: "Two Pointers по разным Array", URL: "https://habr.com/ru/articles/883030/", Source: "Хабр"},
			},
			Figures: []sets.Figure{
				{
					ID:    "tp3-merge",
					Title: "Merge: сравнения и записи",
					Frames: []sets.FigureFrame{
						{Caption: "До первого сравнения: 1 <= 2. Следующий ход добавит 1 из a и увеличит i.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 8}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "i", Cell: 0}}},
							{Label: "b", Cells: []any{2, 3, 6, 7, 9}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "j", Cell: 0}}},
							{Label: "out", Cells: []any{"", "", "", "", "", "", "", "", ""}},
							{Label: "Сравнения", Cells: []any{0}},
							{Label: "Записи в out", Cells: []any{0}},
						}},
						{Caption: "Одно сравнение и одна запись завершены. Далее 3 > 2: добавим 2 из b, увеличим j.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 8}, Marked: []int{1}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "i", Cell: 1}}},
							{Label: "b", Cells: []any{2, 3, 6, 7, 9}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "j", Cell: 0}}},
							{Label: "out", Cells: []any{1, "", "", "", "", "", "", "", ""}},
							{Label: "Сравнения", Cells: []any{1}},
							{Label: "Записи в out", Cells: []any{1}},
						}},
						{Caption: "Два сравнения завершены. Далее 3 <= 3: при равенстве добавим число из a.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 8}, Marked: []int{1}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "i", Cell: 1}}},
							{Label: "b", Cells: []any{2, 3, 6, 7, 9}, Marked: []int{1}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "j", Cell: 1}}},
							{Label: "out", Cells: []any{1, 2, "", "", "", "", "", "", ""}},
							{Label: "Сравнения", Cells: []any{2}},
							{Label: "Записи в out", Cells: []any{2}},
						}},
						{Caption: "Три сравнения завершены. Далее 5 > 3: добавим 3 из b.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 8}, Marked: []int{2}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "i", Cell: 2}}},
							{Label: "b", Cells: []any{2, 3, 6, 7, 9}, Marked: []int{1}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "j", Cell: 1}}},
							{Label: "out", Cells: []any{1, 2, 3, "", "", "", "", "", ""}},
							{Label: "Сравнения", Cells: []any{3}},
							{Label: "Записи в out", Cells: []any{3}},
						}},
						{Caption: "Четыре сравнения завершены. Далее 5 <= 6: добавим 5 из a.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 8}, Marked: []int{2}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "i", Cell: 2}}},
							{Label: "b", Cells: []any{2, 3, 6, 7, 9}, Marked: []int{2}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}},
							{Label: "out", Cells: []any{1, 2, 3, 3, "", "", "", "", ""}},
							{Label: "Сравнения", Cells: []any{4}},
							{Label: "Записи в out", Cells: []any{4}},
						}},
						{Caption: "Пять сравнений завершены. Далее 8 > 6: добавим 6 из b.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 8}, Marked: []int{3}, Dimmed: []int{0, 1, 2}, Pointers: &sets.Pointers{{Name: "i", Cell: 3}}},
							{Label: "b", Cells: []any{2, 3, 6, 7, 9}, Marked: []int{2}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}},
							{Label: "out", Cells: []any{1, 2, 3, 3, 5, "", "", "", ""}},
							{Label: "Сравнения", Cells: []any{5}},
							{Label: "Записи в out", Cells: []any{5}},
						}},
						{Caption: "Шесть сравнений завершены. Далее 8 > 7: добавим 7 из b.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 8}, Marked: []int{3}, Dimmed: []int{0, 1, 2}, Pointers: &sets.Pointers{{Name: "i", Cell: 3}}},
							{Label: "b", Cells: []any{2, 3, 6, 7, 9}, Marked: []int{3}, Dimmed: []int{0, 1, 2}, Pointers: &sets.Pointers{{Name: "j", Cell: 3}}},
							{Label: "out", Cells: []any{1, 2, 3, 3, 5, 6, "", "", ""}},
							{Label: "Сравнения", Cells: []any{6}},
							{Label: "Записи в out", Cells: []any{6}},
						}},
						{Caption: "Семь сравнений завершены. Далее 8 <= 9: добавим 8 из a. После этого i = 4, j = 4.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 8}, Marked: []int{3}, Dimmed: []int{0, 1, 2}, Pointers: &sets.Pointers{{Name: "i", Cell: 3}}},
							{Label: "b", Cells: []any{2, 3, 6, 7, 9}, Marked: []int{4}, Dimmed: []int{0, 1, 2, 3}, Pointers: &sets.Pointers{{Name: "j", Cell: 4}}},
							{Label: "out", Cells: []any{1, 2, 3, 3, 5, 6, 7, "", ""}},
							{Label: "Сравнения", Cells: []any{7}},
							{Label: "Записи в out", Cells: []any{7}},
						}},
						{Caption: "Array a закончился: восемь сравнений и восемь сдвигов. Остаток [9] добавляется без сравнений значений. Срез и объединение временно хранят две копии одной ячейки.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 8}, Dimmed: []int{0, 1, 2, 3}},
							{Label: "b", Cells: []any{2, 3, 6, 7, 9}, Marked: []int{4}, Dimmed: []int{0, 1, 2, 3}, Pointers: &sets.Pointers{{Name: "j", Cell: 4}}},
							{Label: "out", Cells: []any{1, 2, 3, 3, 5, 6, 7, 8, ""}},
							{Label: "Сравнения", Cells: []any{8}},
							{Label: "Записи в out", Cells: []any{8}},
						}},
						{Caption: "Готово: 8 сравнений значений, 9 записей в out, 8 сдвигов указателей. i = 4, j = 4: копирование остатка не меняет индексы. Вход содержит 9 чисел; ответ содержит 9. Пик временных копий остатка - 2 ячейки, кроме постоянного состояния.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 8}, Dimmed: []int{0, 1, 2, 3}},
							{Label: "b", Cells: []any{2, 3, 6, 7, 9}, Dimmed: []int{0, 1, 2, 3, 4}, Pointers: &sets.Pointers{{Name: "j", Cell: 4}}},
							{Label: "out", Cells: []any{1, 2, 3, 3, 5, 6, 7, 8, 9}},
							{Label: "Сравнения", Cells: []any{8}},
							{Label: "Записи в out", Cells: []any{9}},
							{Label: "Пик копий остатка", Cells: []any{2}},
						}},
					},
				},
				{
					ID:    "tp3-count",
					Title: "Подсчёт: сравнения, сдвиги и ответы",
					Frames: []sets.FigureFrame{
						{Caption: "До первого сравнения: 1 < 2. Увеличим i, ответ пока не записываем.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 8}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "i", Cell: 0}}},
							{Label: "b", Cells: []any{2, 3, 6}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "j", Cell: 0}}},
							{Label: "ответ", Cells: []any{"", "", ""}},
							{Label: "Сравнения", Cells: []any{0}},
							{Label: "Сдвиги i", Cells: []any{0}},
							{Label: "Записи ответа", Cells: []any{0}},
						}},
						{Caption: "Одно сравнение завершено. Далее 3 < 2 неверно: запишем i = 1 для запроса 2.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 8}, Marked: []int{1}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "i", Cell: 1}}, Window: &[2]int{0, 0}},
							{Label: "b", Cells: []any{2, 3, 6}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "j", Cell: 0}}},
							{Label: "ответ", Cells: []any{"", "", ""}},
							{Label: "Сравнения", Cells: []any{1}},
							{Label: "Сдвиги i", Cells: []any{1}},
							{Label: "Записи ответа", Cells: []any{0}},
						}},
						{Caption: "Два сравнения завершены. Далее 3 < 3 неверно: равное число не учитываем, запишем i = 1 для запроса 3.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 8}, Marked: []int{1}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "i", Cell: 1}}, Window: &[2]int{0, 0}},
							{Label: "b", Cells: []any{2, 3, 6}, Marked: []int{1}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "j", Cell: 1}}},
							{Label: "ответ", Cells: []any{1, "", ""}},
							{Label: "Сравнения", Cells: []any{2}},
							{Label: "Сдвиги i", Cells: []any{1}},
							{Label: "Записи ответа", Cells: []any{1}},
						}},
						{Caption: "Три сравнения завершены. Далее 3 < 6: увеличим i.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 8}, Marked: []int{1}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "i", Cell: 1}}},
							{Label: "b", Cells: []any{2, 3, 6}, Marked: []int{2}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}},
							{Label: "ответ", Cells: []any{1, 1, ""}},
							{Label: "Сравнения", Cells: []any{3}},
							{Label: "Сдвиги i", Cells: []any{1}},
							{Label: "Записи ответа", Cells: []any{2}},
						}},
						{Caption: "Четыре сравнения завершены. Далее 5 < 6: увеличим i ещё раз.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 8}, Marked: []int{2}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "i", Cell: 2}}},
							{Label: "b", Cells: []any{2, 3, 6}, Marked: []int{2}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}},
							{Label: "ответ", Cells: []any{1, 1, ""}},
							{Label: "Сравнения", Cells: []any{4}},
							{Label: "Сдвиги i", Cells: []any{2}},
							{Label: "Записи ответа", Cells: []any{2}},
						}},
						{Caption: "Пять сравнений завершены. Далее 8 < 6 неверно: запишем i = 3 для запроса 6.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 8}, Marked: []int{3}, Dimmed: []int{0, 1, 2}, Pointers: &sets.Pointers{{Name: "i", Cell: 3}}, Window: &[2]int{0, 2}},
							{Label: "b", Cells: []any{2, 3, 6}, Marked: []int{2}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}},
							{Label: "ответ", Cells: []any{1, 1, ""}},
							{Label: "Сравнения", Cells: []any{5}},
							{Label: "Сдвиги i", Cells: []any{3}},
							{Label: "Записи ответа", Cells: []any{2}},
						}},
						{Caption: "Все запросы обработаны: [1, 1, 3]. Всего 6 сравнений значений, 3 сдвига i и 3 записи ответа. Вход: 4 + 3 числа. Рабочее состояние: два индекса, O(1). Ответ: 3 числа, O(m).", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 8}, Marked: []int{3}, Dimmed: []int{0, 1, 2}, Pointers: &sets.Pointers{{Name: "i", Cell: 3}}},
							{Label: "b", Cells: []any{2, 3, 6}, Dimmed: []int{0, 1, 2}},
							{Label: "ответ", Cells: []any{1, 1, 3}},
							{Label: "Сравнения", Cells: []any{6}},
							{Label: "Сдвиги i", Cells: []any{3}},
							{Label: "Записи ответа", Cells: []any{3}},
							{Label: "Рабочие индексы", Cells: []any{2}},
						}},
					},
				},
			},
		},
		{
			Title:   "Карта Two Pointers",
			Summary: "Выбираем форму Two Pointers и проверяем её условия. Учитываем Sorting, внешние циклы, Sliding Window и память временных String.",
			Lesson:  lesson4,
			Tasks: []sets.Task{
				{
					URL:        "https://codeforces.com/problemset/problem/1538/C",
					Title:      "Number of Pairs",
					Difficulty: "medium",
					Topics:     []string{"Array", "Two Pointers", "Sorting"},
					Note:       "n - длина Array одного теста. Один Sorting и два прохода count(r) - count(l - 1). Время: O(n log n + n), каждый проход проверяет max(n - 1, 0) сумм; цены тестов складываются. Память: O(n) дополнительная для Python Sorting, вход O(n), ответ O(1). Чтение всего ввода из B байт через read().split() добавляет O(B) времени и памяти для токенов.",
					Statement: `<p>Дан Array <code>a</code> из <code>n</code> целых чисел. Найдите количество пар <code>(i, j)</code> (<code>1 ≤ i &lt; j ≤ n</code>), для которых сумма <code>a<sub>i</sub> + a<sub>j</sub></code> больше или равна <code>l</code> и меньше или равна <code>r</code> (то есть <code>l ≤ a<sub>i</sub> + a<sub>j</sub> ≤ r</code>).</p><p>Например, если <code>n = 3</code>, <code>a = [5, 1, 2]</code>, <code>l = 4</code> и <code>r = 7</code>, то подходят две пары:</p><ul><li><code>i = 1</code> и <code>j = 2</code> (<code>4 ≤ 5 + 1 ≤ 7</code>);</li><li><code>i = 1</code> и <code>j = 3</code> (<code>4 ≤ 5 + 2 ≤ 7</code>).</li></ul><h3>Входные данные</h3><p>В первой строке записано целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>). Дальше идут <code>t</code> наборов входных данных.</p><p>В первой строке каждого набора записаны три целых числа <code>n, l, r</code> (<code>1 ≤ n ≤ 2 · 10<sup>5</sup></code>, <code>1 ≤ l ≤ r ≤ 10<sup>9</sup></code>): длина Array и границы для суммы пары.</p><p>Во второй строке записаны <code>n</code> целых чисел <code>a<sub>1</sub>, a<sub>2</sub>, ..., a<sub>n</sub></code> (<code>1 ≤ a<sub>i</sub> ≤ 10<sup>9</sup></code>).</p><p>Гарантируется, что сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого набора выведите одно целое число: количество пар индексов <code>(i, j)</code> (<code>i &lt; j</code>), для которых <code>l ≤ a<sub>i</sub> + a<sub>j</sub> ≤ r</code>.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>4
3 4 7
5 1 2
5 5 8
5 1 2 4 3
4 100 1000
1 1 1 1
5 9 13
2 5 5 1 1</pre><p><strong>Выходные данные:</strong></p><pre>2
7
0
1</pre><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>1 ≤ n ≤ 2 · 10<sup>5</sup></code>, <code>1 ≤ l ≤ r ≤ 10<sup>9</sup></code></li><li><code>1 ≤ a<sub>i</sub> ≤ 10<sup>9</sup></code></li><li>сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `4
3 4 7
5 1 2
5 5 8
5 1 2 4 3
4 100 1000
1 1 1 1
5 9 13
2 5 5 1 1
`, Expected: text(`2
7
0
1
`)},
						{Input: `1
1 1 2
1
`, Expected: text(`0
`)},
						{Input: `1
2 2 2
1 1
`, Expected: text(`1
`)},
						{Input: `1
4 1 1000000000
1000000000 1 1 1
`, Expected: text(`3
`)},
						{Input: `1
5 2 2
1 1 1 1 1
`, Expected: text(`10
`)},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/3sum/",
					Number:     number(15),
					Title:      "3Sum",
					Difficulty: "medium",
					Topics:     []string{"Array", "Two Pointers", "Sorting"},
					Note:       "n = len(nums), q - число разных троек ответа. После Sorting фиксируйте одно число и запускайте Two Pointers для остальных. Время: O(n²), потому что O(n) проход повторяется до n раз; Sorting стоит O(n log n). Память: O(n) дополнительная для Python Sorting; вход O(n), ответ O(q) троек. Пропускайте соседние повторы. Brute Force троек требует O(n³) времени.",
					Statement: `<p>Дан Array целых чисел <code>nums</code>. Верните все тройки <code>[nums[i], nums[j], nums[k]]</code>, для которых <code>i != j</code>, <code>i != k</code>, <code>j != k</code> и <code>nums[i] + nums[j] + nums[k] == 0</code>.</p><p>В ответе не должно быть одинаковых троек.</p><p>Метод класса <code>Solution</code>: <code>def threeSum(self, nums: list[int]) -&gt; list[list[int]]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [-1,0,1,2,-1,-4]
<strong>Выход:</strong> [[-1,-1,2],[-1,0,1]]
<strong>Пояснение:</strong>
nums[0] + nums[1] + nums[2] = (-1) + 0 + 1 = 0.
nums[1] + nums[2] + nums[4] = 0 + 1 + (-1) = 0.
nums[0] + nums[3] + nums[4] = (-1) + 2 + (-1) = 0.
Различные тройки: [-1,0,1] и [-1,-1,2].
Порядок троек в ответе и порядок чисел в тройке не важен.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [0,1,1]
<strong>Выход:</strong> []
<strong>Пояснение:</strong> Единственная возможная тройка не даёт в сумме 0.</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> nums = [0,0,0]
<strong>Выход:</strong> [[0,0,0]]
<strong>Пояснение:</strong> Единственная возможная тройка даёт в сумме 0.</pre><h3>Ограничения</h3><ul><li><code>3 ≤ nums.length ≤ 3000</code></li><li><code>-10<sup>5</sup> ≤ nums[i] ≤ 10<sup>5</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "threeSum", Order: "any-nested"}, Cases: []sets.TestCase{
						{Input: "[[-1,0,1,2,-1,-4]]", Expected: text("[[-1,-1,2],[-1,0,1]]")},
						{Input: "[[0,1,1]]", Expected: text("[]")},
						{Input: "[[0,0,0]]", Expected: text("[[0,0,0]]")},
						{Input: "[[0,0,0,0]]", Expected: text("[[0,0,0]]")},
						{Input: "[[-2,0,1,1,2]]", Expected: text("[[-2,0,2],[-2,1,1]]")},
						{Input: "[[-4,-2,-2,-2,0,1,2,2,2,3,3,4,4,6,6]]", Expected: text("[[-4,-2,6],[-4,0,4],[-4,1,3],[-4,2,2],[-2,-2,4],[-2,0,2]]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/string-compression/",
					Number:     number(443),
					Title:      "String Compression",
					Difficulty: "medium",
					Topics:     []string{"String", "Two Pointers"},
					Note:       "n = len(chars). Читайте группы и пишите символ с длиной в префикс того же Array. Время: O(n), n продвижений read и не больше n записей. Память: O(log(n + 1)) временных символов при str(count), кроме O(1) индексов; вход O(n), ответ в его префиксе. Для строгого O(1) дополнительной памяти извлекайте цифры арифметически.",
					Statement: `<p>Дан Array символов <code>chars</code>. Сожмите его по такому алгоритму:</p><p>Начните с пустой String <code>s</code>. Для каждой группы <strong>одинаковых символов подряд</strong> в <code>chars</code>:</p><ul><li>Если длина группы равна <code>1</code>, допишите символ в <code>s</code>.</li><li>Иначе допишите символ, а после него длину группы.</li></ul><p>Итоговую String <code>s</code> <strong>не нужно возвращать отдельно</strong>. Запишите её <strong>во входной Array символов <code>chars</code></strong>. Если длина группы равна <code>10</code> или больше, её цифры займут несколько ячеек.</p><p>Верните <em>длину записанного префикса</em>.</p><p>Условие требует постоянную дополнительную память. Способ записи цифр должен выполнять это требование.</p><p><strong>Примечание:</strong> символы Array после возвращённой длины не важны, на них не смотрят.</p><p>Метод класса <code>Solution</code>: <code>def compress(self, chars: list[str]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> chars = ["a","a","b","b","c","c","c"]
<strong>Выход:</strong> 6
<strong>Пояснение:</strong> Группы: "aa", "bb" и "ccc". Итоговая String: "a2b2c3".
Первые 6 символов chars после записи на месте: ["a","2","b","2","c","3"].</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> chars = ["a"]
<strong>Выход:</strong> 1
<strong>Пояснение:</strong> Единственная группа "a" содержит один символ.
Первый символ chars после записи: ["a"].</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> chars = ["a","b","b","b","b","b","b","b","b","b","b","b","b"]
<strong>Выход:</strong> 4
<strong>Пояснение:</strong> Группы: "a" и "bbbbbbbbbbbb". Итоговая String: "ab12".
Первые 4 символа chars после записи на месте: ["a","b","1","2"].</pre><h3>Ограничения</h3><ul><li><code>1 ≤ chars.length ≤ 2000</code></li><li><code>chars[i]</code> это строчная английская буква, заглавная английская буква, цифра или другой символ</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "compress", Check: "prefix"}, Cases: []sets.TestCase{
						{Input: `[["a","a","b","b","c","c","c"]]`, Expected: text(`["a","2","b","2","c","3"]`)},
						{Input: `[["a"]]`, Expected: text(`["a"]`)},
						{Input: `[["a","b","b","b","b","b","b","b","b","b","b","b","b"]]`, Expected: text(`["a","b","1","2"]`)},
						{Input: `[["a","a","a","b","b","a","a"]]`, Expected: text(`["a","3","b","2","a","2"]`)},
						{Input: `[["#","#","1"]]`, Expected: text(`["#","2","1"]`)},
						{Input: `[["z","z","z","z","z","z","z","z","z","z"]]`, Expected: text(`["z","1","0"]`)},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "tp4-law-steps", Type: "code", Prompt: "Это часть решения LeetCode 443. В chars 1000 букв. Сколько раз всего выполнится строка read += 1?", Language: "python", Code: `read = 0
while read < len(chars):
    start = read
    while read < len(chars) and chars[read] == chars[start]:
        read += 1
    # дальше пишем букву и длину группы read - start`, Options: []string{"1 000 000", "Столько, сколько в chars групп", "500 500", "1000"}, Answer: 3, Explain: "read растёт от 0 до 1000. Каждый элемент даёт одно увеличение; внутренний цикл не начинает заново с нуля. Число групп меняет число итераций внешнего цикла. Это 1000 продвижений, а не все обращения к chars: проверки границ групп тоже читают символы. При записи через str(count) нужны ещё O(log(n + 1)) временных символов, где n = len(chars)."},
				{ID: "tp4-zeroes-form", Type: "choice", Prompt: "Какая форма Two Pointers? Перенести все нули в конец Array на месте. Порядок остальных чисел не меняется: [0, 1, 0, 3, 12] превратить в [1, 3, 12, 0, 0].", Options: []string{"Навстречу", "Читатель и писатель", "Две последовательности", "Другое"}, Answer: 1, Explain: "Обмен нуля слева с числом справа может нарушить порядок. Читатель идёт по всем числам, писатель размещает ненулевые подряд, затем заполняет остаток нулями. Для этого примера пять проверок ненулевого значения, три записи ненулевых и две записи нулей. При n элементах время O(n), дополнительная память O(1), результат остаётся во входном Array. Это LeetCode 283."},
				{ID: "tp4-reverse-form", Type: "choice", Prompt: "Какая форма Two Pointers? Развернуть Array на месте: [1, 2, 3, 4, 5] превратить в [5, 4, 3, 2, 1].", Options: []string{"Навстречу", "Читатель и писатель", "Две последовательности", "Другое"}, Answer: 0, Explain: "Первое число меняется с последним, второе с предпоследним. l идёт с начала, r с конца. Для пяти чисел это два обмена и четыре записи. Для n элементов нужны floor(n/2) обменов: O(n) времени, O(1) дополнительной памяти, результат в исходном Array. Срез a[::-1] требует O(n) ячеек нового результата."},
				{ID: "tp4-diff-trap", Type: "choice", Prompt: "a = [1, 7, 9, 10], ищем два числа с разностью 2. Идём навстречу: если a[r] - a[l] больше 2, то r -= 1, если меньше 2, то l += 1. Что получится?", Options: []string{"Указатели встретятся, пара 7 и 9 потеряется", "Цикл не остановится", "Найдёт 7 и 9, если при разности больше 2 двигать l", "Найдёт 7 и 9 на третьем ходу"}, Answer: 0, Explain: "Проверки 10 - 1, 9 - 1, 7 - 1 дают три сдвига r и потерю ответа. Для положительной разности используем Two Pointers с движением вперёд: при малой разности увеличиваем правый, при большой левый; равные индексы разделяем. Здесь проверим 7 - 1, затем 9 - 7 и найдём пару. На готовом порядке n элементов дают O(n) движений, O(n) времени и O(1) дополнительной памяти."},
				{ID: "tp4-two-sum", Type: "choice", Prompt: "Two Sum, X = 6: нужно вернуть исходные индексы двух чисел. Код делает a = sorted(nums), идёт навстречу по a и возвращает [l, r], когда a[l] + a[r] = 6. Для какого nums код вернёт верные индексы?", Options: []string{"[2, 1, 5, 7]", "[1, 7, 5, 2]", "[7, 1, 5, 2]", "[1, 2, 7, 5]"}, Answer: 1, Explain: "После Sorting во всех вариантах a = [1, 2, 5, 7]. Код возвращает [0, 2]. Только в [1, 7, 5, 2] эти исходные индексы дают 6. Сохраняйте пары (значение, исходный индекс): O(n log n) времени, O(n) дополнительной памяти Python при n = len(nums). Hash Table даёт ожидаемое O(n) время и O(n) рабочую память. Ответ в обоих случаях содержит два индекса."},
				{ID: "tp4-compress-len", Type: "short", Prompt: "LeetCode 443: chars = [a, a, b, c, c, c, c, c, c, c, c, c, c, c], это две a, одна b и одиннадцать c. Какое число вернёт функция?", Answer: "6", Explain: `Итоговая String a2bc11 имеет длину 2 + 1 + 3 = 6. Группа из одной b пишется без числа. Чтение обрабатывает 14 элементов, запись занимает шесть ячеек исходного Array. При str(count) самая длинная временная String здесь "11", два символа. Общая оценка для n символов: O(n) времени и O(log(n + 1)) дополнительной памяти для цифр.`},
				{ID: "tp4-count-row", Type: "code", Prompt: "count(a, x) считает пары с суммой не больше x, Array a отсортирован. Что вписать вместо ___?", Language: "python", Code: `lo, hi = 0, len(a) - 1
res = 0
while lo < hi:
    if a[lo] + a[hi] <= x:
        res += ___
        lo += 1
    else:
        hi -= 1
return res`, Options: []string{"1", "hi", "hi - lo", "hi - lo + 1"}, Answer: 2, Explain: "Если a[lo] + a[hi] <= x, то подходят все hi - lo пар с первым индексом lo. Для a = [1, 2, 4, 5, 7], x = 8 ответ 6 получается за четыре проверки суммы. С «1» получится 2, с «hi - lo + 1» получится 8: добавятся пары числа с самим собой. Для n элементов проход делает max(n - 1, 0) проверок: O(n) времени, O(1) дополнительной памяти и одно число ответа."},
				{ID: "tp4-count-bounds", Type: "choice", Prompt: "Codeforces 1538C. count(x) - число пар с суммой не больше x. Решение выводит count(r) - count(l). На каком тесте ответ будет неверным?", Options: []string{"a = [5, 1, 2], l = 4, r = 7", "Ни на одном: ответы совпадут", "a = [2, 5, 5, 1, 1], l = 9, r = 13", "a = [5, 1, 2, 4, 3], l = 5, r = 8"}, Answer: 3, Explain: "count(l) включает подходящие пары с суммой ровно l. В [5, 1, 2, 4, 3] это 1 + 4 и 2 + 3. Ошибка даёт 9 - 4 = 5 вместо 9 - 2 = 7. Нужна формула count(r) - count(l - 1). Для n = 5 оба прохода вместе делают восемь проверок суммы; Sorting выполняется один раз перед ними."},
			},
			Links: []sets.Link{
				{Title: "Краткая теория Two Pointers", URL: "https://informatics.msk.ru/mod/page/view.php?id=16344", Source: "informatics.msk.ru"},
				{Title: "14 шаблонов задач: когда нужны Two Pointers", URL: "https://tproger.ru/translations/14-templates-to-answer-interview-questions", Source: "Tproger"},
				{Title: "Задачи с тегом Two Pointers, рейтинг от 800 до 1200", URL: "https://codeforces.com/problemset?tags=two+pointers,800-1200", Source: "Codeforces"},
				{Title: "ITMO Academy: Two Pointers", URL: "https://codeforces.com/edu/course/2/lesson/9", Source: "Codeforces EDU (нужен вход)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "tp4-compress",
					Title: "String Compression: чтения и записи",
					Frames: []sets.FigureFrame{
						{Caption: "Первая группа начинается в ячейке 0. Счётчик чтений показывает число продвижений read, без повторных проверок границ группы.", Rows: []sets.FigureRow{
							{Cells: []any{"a", "a", "a", "b", "b"}, Pointers: &sets.Pointers{{Name: "write", Cell: 0}, {Name: "read", Cell: 0}}},
							{Label: "Прочитано", Cells: []any{0}},
							{Label: "Записи", Cells: []any{0}},
						}},
						{Caption: "Группа a a a прочитана: read увеличился три раза и стоит в ячейке 3. Проверка следующей b завершила группу.", Rows: []sets.FigureRow{
							{Cells: []any{"a", "a", "a", "b", "b"}, Pointers: &sets.Pointers{{Name: "write", Cell: 0}, {Name: "read", Cell: 3}}, Window: &[2]int{0, 2}},
							{Label: "Прочитано", Cells: []any{3}},
							{Label: "Записи", Cells: []any{0}},
						}},
						{Caption: "Записали a и 3: две записи в ячейки 0 и 1. write = 2. Для str(3) нужна одна временная цифра. Ячейка 2 не входит в готовый префикс.", Rows: []sets.FigureRow{
							{Cells: []any{"a", "3", "a", "b", "b"}, Marked: []int{0, 1}, Dimmed: []int{2}, Pointers: &sets.Pointers{{Name: "write", Cell: 2}, {Name: "read", Cell: 3}}},
							{Label: "Прочитано", Cells: []any{3}},
							{Label: "Записи", Cells: []any{2}},
							{Label: "Временная цифра", Cells: []any{"3"}},
						}},
						{Caption: "Вторая группа b b прочитана: ещё два продвижения read. Всего обработано пять элементов входного Array.", Rows: []sets.FigureRow{
							{Cells: []any{"a", "3", "a", "b", "b"}, Dimmed: []int{2}, Pointers: &sets.Pointers{{Name: "write", Cell: 2}}, Window: &[2]int{3, 4}},
							{Label: "Прочитано", Cells: []any{5}},
							{Label: "Записи", Cells: []any{2}},
						}},
						{Caption: "Записали b и 2. Итог: 5 прочитанных элементов, 4 записи, две группы. Вход занимает пять ячеек, ответ использует его префикс из четырёх ячеек. Пик временных цифр при str(count) здесь равен 1; индексы и счётчики занимают O(1).", Rows: []sets.FigureRow{
							{Cells: []any{"a", "3", "b", "2", "b"}, Marked: []int{2, 3}, Dimmed: []int{4}, Pointers: &sets.Pointers{{Name: "write", Cell: 4}}},
							{Label: "Прочитано", Cells: []any{5}},
							{Label: "Записи", Cells: []any{4}},
							{Label: "Пик временных цифр", Cells: []any{1}},
						}},
					},
				},
				{
					ID:    "tp4-wrong-move",
					Title: "Two Pointers: цена и правильность хода",
					Frames: []sets.FigureFrame{
						{Caption: "До хода l = 1, r = 4: 4 + 10 = 14, меньше 15. В обоих вариантах одно вычисление суммы и один сдвиг. Проверки ==, затем < дают два сравнения значений. Верный ход l += 1 исключает 4: с оставшимися числами оно даёт не больше 14. Ошибка r -= 1 теряет пару 5 + 10. Рабочая память - два индекса и сумма, O(1); нарисованные копии Array алгоритму не нужны.", Rows: []sets.FigureRow{
							{Label: "верно", Cells: []any{1, 4, 5, 7, 10, 12, 16}, Marked: []int{2, 4}, Dimmed: []int{0, 1, 5, 6}, Pointers: &sets.Pointers{{Name: "l", Cell: 2}, {Name: "r", Cell: 4}}, Window: &[2]int{2, 4}},
							{Label: "ошибка", Cells: []any{1, 4, 5, 7, 10, 12, 16}, Marked: []int{4}, Dimmed: []int{0, 5, 6}, Pointers: &sets.Pointers{{Name: "l", Cell: 1}, {Name: "r", Cell: 3}}, Window: &[2]int{1, 3}},
							{Label: "Вычисления суммы", Cells: []any{1}},
							{Label: "Сдвиги за ход", Cells: []any{1}},
							{Label: "Рабочие значения", Cells: []any{3}},
						}},
					},
				},
			},
		},
	},
}
