// Package modular is the learning set "Modular Arithmetic".
package modular

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

func number(value int64) *int64 { return &value }

func text(value string) *string { return &value }

// Set is the content of the learning set.
var Set = sets.LearningSet{
	Slug:    "modular-arithmetic",
	Title:   "Modular Arithmetic",
	Summary: "Остаток от деления помогает работать с огромными числами. Вы считаете суммы и произведения по модулю, возводите в степень через Binary Exponentiation и делите через обратный элемент. Для каждого решения считаете операции и дополнительную память.",
	Topics:  []string{"Math", "Number Theory", "Hash Table", "Prefix Sum", "Combinatorics"},
	Chapters: []sets.Chapter{
		{
			Title:   "Остатки и классы остатков",
			Summary: "Python % с отрицательными числами, остаток суммы и произведения, последние цифры. Числа с равным остатком образуют класс: по классам считаем пары и Subarray с суммой, кратной k.",
			Lesson:  lesson1,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/binary-prefix-divisible-by-5/",
					Number:     number(1018),
					Title:      "Binary Prefix Divisible By 5",
					Difficulty: "easy",
					Topics:     []string{"Array", "Bit Manipulation", "Math"},
					Note:       "Храните только r = x % 5: новый бит даёт r = (2r + bit) % 5. n - длина nums. Время: O(n), один проход, все числа меньше 10. Память: O(1) дополнительной, ответ - n значений bool, O(n). Для [1, 0, 1, 0]: r = 1, 2, 0, 0. Полное число x растёт до n двоичных цифр: каждый шаг тогда стоит O(n) работы с цифрами, всего O(n²).",
					Statement: `<p>Дан Array <code>nums</code> из нулей и единиц. Номера элементов в нём <strong>начинаются с 0</strong>.</p><p>Обозначим через <code>x<sub>i</sub></code> число, двоичная запись которого - Subarray <code>nums[0..i]</code> (старший бит слева, младший справа).</p><ul><li>Например, если <code>nums = [1,0,1]</code>, то <code>x<sub>0</sub> = 1</code>, <code>x<sub>1</sub> = 2</code> и <code>x<sub>2</sub> = 5</code>.</li></ul><p>Верните Array значений bool <code>answer</code>, где <code>answer[i]</code> равно <code>True</code>, если <code>x<sub>i</sub></code> делится на <code>5</code>.</p><p>Метод класса <code>Solution</code>: <code>def prefixesDivBy5(self, nums: list[int]) -&gt; list[bool]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [0,1,1]
<strong>Выход:</strong> [True,False,False]
<strong>Пояснение:</strong> двоичные записи входных чисел 0, 01 и 011, в десятичной записи это 0, 1 и 3.
Только первое число делится на 5, поэтому только answer[0] равно True.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [1,1,1]
<strong>Выход:</strong> [False,False,False]</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 10<sup>5</sup></code></li><li><code>nums[i]</code> равно <code>0</code> или <code>1</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "prefixesDivBy5"}, Cases: []sets.TestCase{
						{Input: "[[0,1,1]]", Expected: text("[True,False,False]")},
						{Input: "[[1,1,1]]", Expected: text("[False,False,False]")},
						{Input: "[[1,0,1,0]]", Expected: text("[False,False,True,True]")},
						{Input: "[[0]]", Expected: text("[True]")},
						{Input: "[[1,1,1,1,0,1,1,1,1,0]]", Expected: text("[False,False,False,True,True,False,False,False,True,True]")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1490/B",
					Title:      "Balanced Remainders",
					Difficulty: "easy",
					Topics:     []string{"Math", "Counting", "Greedy"},
					Note:       "Codeforces, рейтинг 1000. Ход +1 переводит число в следующий класс остатков: 0 → 1 → 2 → 0. Посчитайте c0, c1, c2 и переносите излишек класса в следующий класс, пока все счётчики не станут n/3. n - длина набора. Время: O(n) на подсчёт и O(1) на переносы. Память: O(1) дополнительной, три счётчика. Для [0, 0, 0]: c = (3, 0, 0), два числа переходят в класс 1, одно из них дальше в класс 2, ответ 2 + 1 = 3. Simulation по одному ходу с пересчётом всех остатков стоит O(n) на ход и до O(n²) всего.",
					Statement: `<p>Дано число <code>n</code>, которое делится на <code>3</code>, и Array <code>a[1..n]</code>. За один ход можно увеличить любой элемент Array на единицу. Точнее: вы выбираете номер <code>i</code> (<code>1 ≤ i ≤ n</code>) и заменяете <code>a<sub>i</sub></code> на <code>a<sub>i</sub> + 1</code>. Один и тот же номер <code>i</code> можно выбирать в разных ходах несколько раз.</p><p>Обозначим через <code>c<sub>0</sub></code>, <code>c<sub>1</sub></code> и <code>c<sub>2</sub></code> количество чисел Array <code>a</code>, которые дают остаток <code>0</code>, <code>1</code> и <code>2</code> при делении на <code>3</code>. Скажем, что у Array <code>a</code> <strong>сбалансированные остатки</strong>, если числа <code>c<sub>0</sub></code>, <code>c<sub>1</sub></code> и <code>c<sub>2</sub></code> равны.</p><p>Например, если <code>n = 6</code> и <code>a = [0, 2, 5, 5, 4, 8]</code>, возможны такие ходы:</p><ul><li>сначала <code>c<sub>0</sub> = 1</code>, <code>c<sub>1</sub> = 1</code> и <code>c<sub>2</sub> = 4</code>, эти числа не равны. Увеличим <code>a<sub>3</sub></code>, теперь <code>a = [0, 2, 6, 5, 4, 8]</code>;</li><li><code>c<sub>0</sub> = 2</code>, <code>c<sub>1</sub> = 1</code> и <code>c<sub>2</sub> = 3</code>, эти числа не равны. Увеличим <code>a<sub>6</sub></code>, теперь <code>a = [0, 2, 6, 5, 4, 9]</code>;</li><li><code>c<sub>0</sub> = 3</code>, <code>c<sub>1</sub> = 1</code> и <code>c<sub>2</sub> = 2</code>, эти числа не равны. Увеличим <code>a<sub>1</sub></code>, теперь <code>a = [1, 2, 6, 5, 4, 9]</code>;</li><li><code>c<sub>0</sub> = 2</code>, <code>c<sub>1</sub> = 2</code> и <code>c<sub>2</sub> = 2</code>, эти числа равны. Значит, у Array <code>a</code> сбалансированные остатки.</li></ul><p>Найдите наименьшее число ходов, после которых у Array <code>a</code> будут сбалансированные остатки.</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>). Дальше идут <code>t</code> наборов входных данных.</p><p>В первой строке набора записано одно целое число <code>n</code> (<code>3 ≤ n ≤ 3 · 10<sup>4</sup></code>): длина Array <code>a</code>. Гарантируется, что <code>n</code> делится на <code>3</code>.</p><p>В следующей строке записаны <code>n</code> целых чисел <code>a<sub>1</sub>, a<sub>2</sub>, ..., a<sub>n</sub></code> (<code>0 ≤ a<sub>i</sub> ≤ 100</code>).</p><p>Гарантируется, что сумма <code>n</code> по всем наборам не больше <code>150 000</code>.</p><h3>Выходные данные</h3><p>Для каждого набора выведите одно целое число: наименьшее число ходов, после которых у Array <code>a</code> будут сбалансированные остатки.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>4
6
0 2 5 5 4 8
6
2 0 2 1 0 0
9
7 1 3 4 2 10 3 9 6
6
0 1 2 3 4 5</pre><p><strong>Выходные данные:</strong></p><pre>3
1
3
0</pre><h3>Примечание</h3><p>Первый набор разобран в условии.</p><p>Во втором наборе нужен один ход для <code>i = 2</code>.</p><p>В третьем наборе нужны три хода:</p><ul><li>первый ход: <code>i = 9</code>;</li><li>второй ход: <code>i = 9</code>;</li><li>третий ход: <code>i = 2</code>.</li></ul><p>В четвёртом наборе числа <code>c<sub>0</sub></code>, <code>c<sub>1</sub></code> и <code>c<sub>2</sub></code> равны с самого начала, поэтому у Array <code>a</code> уже сбалансированные остатки.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>3 ≤ n ≤ 3 · 10<sup>4</sup></code>, <code>n</code> делится на 3</li><li><code>0 ≤ a<sub>i</sub> ≤ 100</code></li><li>сумма <code>n</code> по всем наборам не больше <code>150 000</code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `4
6
0 2 5 5 4 8
6
2 0 2 1 0 0
9
7 1 3 4 2 10 3 9 6
6
0 1 2 3 4 5
`, Expected: text(`3
1
3
0
`)},
						{Input: `2
3
0 0 0
3
2 2 2
`, Expected: text(`3
3
`)},
						{Input: `1
6
1 1 1 1 1 1
`, Expected: text(`6
`)},
						{Input: `1
6
0 3 6 2 5 8
`, Expected: text(`3
`)},
						{Input: `1
3
100 100 99
`, Expected: text(`1
`)},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/check-if-array-pairs-are-divisible-by-k/",
					Number:     number(1497),
					Title:      "Check If Array Pairs Are Divisible by k",
					Difficulty: "medium",
					Topics:     []string{"Array", "Hash Table", "Counting", "Math"},
					Note:       "Число с остатком r нужно соединить с числом с остатком (k − r) % k. Нужны равные cnt[r] и cnt[k − r], чётный cnt[0], а при чётном k ещё чётный cnt[k/2]. n - длина arr. Время: O(n + k): создание cnt и один проход. Память: O(k) для счётчиков, ответ O(1). Python даёт −1 % 3 = 2, поэтому отрицательные числа сразу попадают в нужный класс. Поиск пары для каждого числа по остальным стоит O(n²).",
					Statement: `<p>Дан Array целых чисел <code>arr</code> чётной длины <code>n</code> и целое число <code>k</code>.</p><p>Нужно разбить Array ровно на <code>n / 2</code> пар так, чтобы сумма каждой пары делилась на <code>k</code>.</p><p>Верните <code>True</code>, если это можно сделать, иначе верните <code>False</code>.</p><p>Метод класса <code>Solution</code>: <code>def canArrange(self, arr: list[int], k: int) -&gt; bool</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> arr = [1,2,3,4,5,10,6,7,8,9], k = 5
<strong>Выход:</strong> True
<strong>Пояснение:</strong> пары (1,9), (2,8), (3,7), (4,6) и (5,10).</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> arr = [1,2,3,4,5,6], k = 7
<strong>Выход:</strong> True
<strong>Пояснение:</strong> пары (1,6), (2,5) и (3,4).</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> arr = [1,2,3,4,5,6], k = 10
<strong>Выход:</strong> False
<strong>Пояснение:</strong> проверка всех возможных пар показывает, что arr нельзя разбить на 3 пары с суммой, кратной 10.</pre><h3>Ограничения</h3><ul><li><code>arr.length == n</code></li><li><code>1 ≤ n ≤ 10<sup>5</sup></code></li><li><code>n</code> чётное</li><li><code>-10<sup>9</sup> ≤ arr[i] ≤ 10<sup>9</sup></code></li><li><code>1 ≤ k ≤ 10<sup>5</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "canArrange"}, Cases: []sets.TestCase{
						{Input: "[[1,2,3,4,5,10,6,7,8,9],5]", Expected: text("True")},
						{Input: "[[1,2,3,4,5,6],7]", Expected: text("True")},
						{Input: "[[1,2,3,4,5,6],10]", Expected: text("False")},
						{Input: "[[-1,4,-2,-4],3]", Expected: text("True")},
						{Input: "[[4,2,1,3],4]", Expected: text("False")},
						{Input: "[[-1000000000,999999999],1]", Expected: text("True")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/subarray-sums-divisible-by-k/",
					Number:     number(974),
					Title:      "Subarray Sums Divisible by K",
					Difficulty: "medium",
					Topics:     []string{"Array", "Hash Table", "Prefix Sum", "Math"},
					Note:       "Сумма nums[l..r] делится на k, если префиксы p[l] и p[r + 1] имеют равные остатки. Считайте остатки префиксов в Array длины k и начните с cnt[0] = 1 для пустого префикса p[0] = 0. n - длина nums. Время: O(n + k). Память: O(k) для счётчиков, ответ O(1). Для [−2, 5, 1, −4, 3], k = 3 остатки префиксов 0, 1, 0, 1, 0, 0, ответ 6 + 1 = 7. Brute Force проверяет n(n + 1)/2 Subarray: при n = 30 000 это 450 015 000.",
					Statement: `<p>Дан Array целых чисел <code>nums</code> и целое число <code>k</code>. Верните число непустых <strong>Subarray</strong>, сумма которых делится на <code>k</code>.</p><p><strong>Subarray</strong> - это часть Array, элементы которой идут <strong>подряд</strong>.</p><p>Метод класса <code>Solution</code>: <code>def subarraysDivByK(self, nums: list[int], k: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [4,5,0,-2,-3,1], k = 5
<strong>Выход:</strong> 7
<strong>Пояснение:</strong> сумму, кратную k = 5, имеют 7 Subarray:
[4, 5, 0, -2, -3, 1], [5], [5, 0], [5, 0, -2, -3], [0], [0, -2, -3], [-2, -3]</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [5], k = 9
<strong>Выход:</strong> 0</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 3 · 10<sup>4</sup></code></li><li><code>-10<sup>4</sup> ≤ nums[i] ≤ 10<sup>4</sup></code></li><li><code>2 ≤ k ≤ 10<sup>4</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "subarraysDivByK"}, Cases: []sets.TestCase{
						{Input: "[[4,5,0,-2,-3,1],5]", Expected: text("7")},
						{Input: "[[5],9]", Expected: text("0")},
						{Input: "[[0,0,0],2]", Expected: text("6")},
						{Input: "[[-1,2,9],2]", Expected: text("2")},
						{Input: "[[2,-2,2,-2,2],4]", Expected: text("6")},
						{Input: "[[-10000,-10000],10000]", Expected: text("3")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1931/D",
					Title:      "Divisible Pairs",
					Difficulty: "medium",
					Topics:     []string{"Hash Table", "Math", "Number Theory", "Combinatorics"},
					Note:       "Codeforces, рейтинг 1300. Пара подходит, если a_i % x = (−a_j) % x и a_i % y = a_j % y. Идите слева направо и храните в Hash Table число встреч ключа (a % x, a % y). n - длина набора. Время: O(n) в среднем при O(1) операциях Hash Table. Память: O(n) ключей. Проверка всех пар: n(n − 1)/2, при n = 200 000 это 19 999 900 000.",
					Statement: `<p>У Поликарпа есть два любимых целых числа <code>x</code> и <code>y</code> (они могут быть равны). Ещё он нашёл Array <code>a</code> длины <code>n</code>.</p><p>Поликарп называет пару номеров <code>(i, j)</code> (<code>1 ≤ i &lt; j ≤ n</code>) красивой, если:</p><ul><li><code>a<sub>i</sub> + a<sub>j</sub></code> делится на <code>x</code>;</li><li><code>a<sub>i</sub> - a<sub>j</sub></code> делится на <code>y</code>.</li></ul><p>Например, если <code>x = 5</code>, <code>y = 2</code>, <code>n = 6</code> и <code>a = [1, 2, 7, 4, 9, 6]</code>, то красивых пар только две:</p><ul><li><code>(1, 5)</code>: <code>a<sub>1</sub> + a<sub>5</sub> = 1 + 9 = 10</code> (10 делится на 5) и <code>a<sub>1</sub> - a<sub>5</sub> = 1 - 9 = -8</code> (−8 делится на 2);</li><li><code>(4, 6)</code>: <code>a<sub>4</sub> + a<sub>6</sub> = 4 + 6 = 10</code> (10 делится на 5) и <code>a<sub>4</sub> - a<sub>6</sub> = 4 - 6 = -2</code> (−2 делится на 2).</li></ul><p>Найдите число красивых пар в Array <code>a</code>.</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>): количество наборов входных данных. Дальше идут описания наборов.</p><p>В первой строке набора записаны три целых числа <code>n</code>, <code>x</code> и <code>y</code> (<code>2 ≤ n ≤ 2 · 10<sup>5</sup></code>, <code>1 ≤ x, y ≤ 10<sup>9</sup></code>): размер Array и любимые числа Поликарпа.</p><p>Во второй строке набора записаны <code>n</code> целых чисел <code>a<sub>1</sub>, a<sub>2</sub>, ..., a<sub>n</sub></code> (<code>1 ≤ a<sub>i</sub> ≤ 10<sup>9</sup></code>): элементы Array.</p><p>Гарантируется, что сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого набора выведите одно целое число: количество красивых пар в Array <code>a</code>.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>7
6 5 2
1 2 7 4 9 6
7 9 5
1 10 15 3 8 12 15
9 4 10
14 10 2 2 11 11 13 5 6
9 5 6
10 7 6 7 9 7 7 10 10
9 6 2
4 9 7 1 2 2 13 3 15
9 2 3
14 6 1 15 12 15 8 2 15
10 5 7
13 3 3 2 12 11 3 7 13 14</pre><p><strong>Выходные данные:</strong></p><pre>2
0
1
3
5
7
0</pre><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>2 ≤ n ≤ 2 · 10<sup>5</sup></code></li><li><code>1 ≤ x, y ≤ 10<sup>9</sup></code></li><li><code>1 ≤ a<sub>i</sub> ≤ 10<sup>9</sup></code></li><li>сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `7
6 5 2
1 2 7 4 9 6
7 9 5
1 10 15 3 8 12 15
9 4 10
14 10 2 2 11 11 13 5 6
9 5 6
10 7 6 7 9 7 7 10 10
9 6 2
4 9 7 1 2 2 13 3 15
9 2 3
14 6 1 15 12 15 8 2 15
10 5 7
13 3 3 2 12 11 3 7 13 14
`, Expected: text(`2
0
1
3
5
7
0
`)},
						{Input: `1
2 1 1
1 1
`, Expected: text(`1
`)},
						{Input: `1
4 2 3
5 5 5 5
`, Expected: text(`6
`)},
						{Input: `1
3 1000000000 1000000000
1000000000 1000000000 1
`, Expected: text(`1
`)},
						{Input: `1
3 3 1
3 6 9
`, Expected: text(`3
`)},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "ma1-neg-mod", Type: "choice", Prompt: "Что выведет print(-8 % 3, -8 // 3) в Python?", Options: []string{"-2 -2", "1 -3", "2 -2"}, Answer: 1, Explain: "Python округляет частное вниз: −8 // 3 = −3. Тогда −8 = 3 · (−3) + 1, и остаток 1 лежит от 0 до 2. Ответ −2 −2 дают C++ и Java: там частное округляют к нулю. В ответе 2 −2 остаток 8 % 3 = 2 взят без знака. Одна операция % с числами из ограничений стоит O(1) времени и O(1) памяти."},
				{ID: "ma1-last-digit", Type: "short", Prompt: "Какая последняя цифра у произведения 37 · 48 · 59? Ответьте одной цифрой.", Answer: "4", Explain: "Последняя цифра - остаток от деления на 10. Умножаем остатки: 7 · 8 = 56, остаток 6; 6 · 9 = 54, остаток 4. Полное произведение 104 784 тоже кончается на 4. Ответ 9 получится, если взять последнюю цифру только у последнего множителя. С % после каждого шага r остаётся меньше 10: n умножений и O(1) памяти."},
				{ID: "ma1-binary-code", Type: "code", Prompt: "LeetCode 1018. Код хранит r = x % 5, где x - двоичное число из прочитанных битов. Что стоит на месте ___?", Language: "python", Code: `r = 0
out = []
for bit in nums:
    r = ___
    out.append(r == 0)`, Options: []string{"(r + bit) % 5", "2 * r + bit % 5", "(2 * r + bit) % 5"}, Answer: 2, Explain: "Новый бит сдвигает число: x становится 2x + bit, и новый остаток (2r + bit) % 5. Вариант (r + bit) % 5 забывает сдвиг: для [1, 0] он даст r = 1 вместо 2. Во втором варианте % применяется только к bit, и r растёт без границы. Время O(n), дополнительная память O(1), ответ - n значений bool."},
				{ID: "ma1-pair-need", Type: "choice", Prompt: "k = 6, новое число x = −8. Какой остаток должен быть у второго числа пары, чтобы сумма делилась на 6?", Options: []string{"4", "2", "−2"}, Answer: 1, Explain: "−8 % 6 = 4, потому что −8 = 6 · (−2) + 4. Второму числу нужен остаток (−x) % k = 8 % 6 = 2. Проверка: −8 + 2 = −6 делится на 6. Ответ 4 - остаток самого x, а остатка −2 в Python при k > 0 нет. Поиск в cnt стоит O(1)."},
				{ID: "ma1-half-class", Type: "choice", Prompt: "a = [3, 9, 5, 15], k = 6. Сколько пар i < j имеют сумму, кратную 6?", Options: []string{"9", "0", "3"}, Answer: 2, Explain: "Остатки 3, 3, 5, 3. Числу с остатком 3 нужен второй остаток (−3) % 6 = 3, тот же класс. Три числа класса 3 дают 3 · 2 / 2 = 3 пары: 3 + 9 = 12, 3 + 15 = 18, 9 + 15 = 24. Числу 5 нужен остаток 1, таких чисел нет. Ответ 0 получится, если соединять остаток r только с другим остатком k − r и забыть класс k/2. 9 = 3 · 3 считает каждую пару дважды и пары числа с самим собой. Один проход: O(n + k) времени и O(k) памяти."},
				{ID: "ma1-prefix-classes", Type: "short", Prompt: "nums = [1, 4, 3, 1, 1, 3], k = 4. Сколько Subarray имеют сумму, кратную 4?", Answer: "7", Explain: "Префиксы p = 0, 1, 5, 8, 9, 10, 13 дают остатки 0, 1, 1, 0, 1, 2, 1. Префиксы 1, 5, 9 и 13 разные, но остаток у них один: класс из 4 префиксов даёт 4 · 3 / 2 = 6 Subarray. Класс 0 из p[0] и p[3] даёт ещё 1, это [1, 4, 3]. Ответ 7. Без p[0] = 0 код потеряет [1, 4, 3] и вернёт 6. Время O(n), дополнительная память O(min(n + 1, k))."},
				{ID: "ma1-cpp-rem", Type: "short", Prompt: "Код на C++ берёт остаток со знаком: −2 % 3 = −2. Сколько Subarray с суммой, кратной 3, он найдёт для nums = [−2, 3, −1]?", Answer: "1", Explain: "Префиксы p = 0, −2, 1, 0. В Python остатки 0, 1, 1, 0: два класса по 2 префикса, верный ответ 2, это [3] и весь Array. В C++ остатки 0, −2, 1, 0: префиксы −2 и 1 попали в разные классы, и код потерял [3]. Поэтому в C++ пишут ((p % k) + k) % k. Время O(n), дополнительная память O(min(n + 1, k))."},
				{ID: "ma1-two-keys", Type: "choice", Prompt: "Codeforces 1931D, x = 3, y = 2. Для нового числа a = 7 какой ключ ищем среди прежних чисел?", Options: []string{"(2, 1)", "(1, 1)", "(2, 0)"}, Answer: 0, Explain: "Сумма делится на 3, если у второго числа остаток (−7) % 3 = 2. Разность делится на 2, если у второго числа остаток 7 % 2 = 1. Ключ (2, 1): например, 5 подходит, 5 + 7 = 12 и 5 − 7 = −2. Ключ (1, 1) - остатки самого a. С ключом (2, 0) сумма делится на 3, но разность нечётная: 2 + 7 = 9 и 2 − 7 = −5. Один поиск в Hash Table стоит O(1) в среднем."},
			},
			Links: []sets.Link{
				{Title: "Деление с остатком", URL: "https://ru.wikipedia.org/wiki/Деление_с_остатком", Source: "Википедия"},
				{Title: "Сравнение по модулю и классы остатков", URL: "https://ru.wikipedia.org/wiki/Сравнение_по_модулю", Source: "Википедия"},
				{Title: "Операции по модулю", URL: "https://brestprog.by/topics/modulo/", Source: "brestprog"},
				{Title: "Prefix Sum и Hash Table", URL: "https://redlamp.me/leetcode/sets/prefix-sums/2", Source: "Лекции"},
				{Title: "Операторы % и // в Python", URL: "https://docs.python.org/3/reference/expressions.html#binary-arithmetic-operations", Source: "Python docs (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "ma1-binary",
					Title: "Остаток растущего двоичного числа",
					Frames: []sets.FigureFrame{
						{Caption: "Бит 1: x = 1, r = (2 · 0 + 1) % 5 = 1. Не делится на 5.", Rows: []sets.FigureRow{
							{Label: "бит", Cells: []any{1, 0, 1, 0}, Marked: []int{0}},
							{Label: "x", Cells: []any{1, "", "", ""}},
							{Label: "r = x % 5", Cells: []any{1, "", "", ""}, Marked: []int{0}},
							{Label: "ответ", Cells: []any{"False", "", "", ""}},
						}},
						{Caption: "Бит 0: x = 2, r = (2 · 1 + 0) % 5 = 2. Не делится.", Rows: []sets.FigureRow{
							{Label: "бит", Cells: []any{1, 0, 1, 0}, Marked: []int{1}},
							{Label: "x", Cells: []any{1, 2, "", ""}},
							{Label: "r = x % 5", Cells: []any{1, 2, "", ""}, Marked: []int{1}},
							{Label: "ответ", Cells: []any{"False", "False", "", ""}},
						}},
						{Caption: "Бит 1: x = 5, r = (2 · 2 + 1) % 5 = 0. Делится: True.", Rows: []sets.FigureRow{
							{Label: "бит", Cells: []any{1, 0, 1, 0}, Marked: []int{2}},
							{Label: "x", Cells: []any{1, 2, 5, ""}},
							{Label: "r = x % 5", Cells: []any{1, 2, 0, ""}, Marked: []int{2}},
							{Label: "ответ", Cells: []any{"False", "False", "True", ""}},
						}},
						{Caption: "Бит 0: x = 10, r = (2 · 0 + 0) % 5 = 0. Код хранит только r меньше 5, а x растёт на один бит за шаг.", Rows: []sets.FigureRow{
							{Label: "бит", Cells: []any{1, 0, 1, 0}, Marked: []int{3}},
							{Label: "x", Cells: []any{1, 2, 5, 10}},
							{Label: "r = x % 5", Cells: []any{1, 2, 0, 0}, Marked: []int{3}},
							{Label: "ответ", Cells: []any{"False", "False", "True", "True"}},
						}},
					},
				},
				{
					ID:    "ma1-pairs",
					Title: "Пары с суммой, кратной 4: один проход",
					Frames: []sets.FigureFrame{
						{Caption: "x = 4, остаток 0. Второму числу пары нужен остаток (−4) % 4 = 0, cnt[0] = 0. Затем cnt[0] = 1.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{4, 7, 6, 1, 8, -3, 2}, Marked: []int{0}},
							{Label: "a % 4", Cells: []any{0, "", "", "", "", "", ""}},
							{Label: "cnt до шага", Cells: []any{0, 0, 0, 0}, Marked: []int{0}},
							{Label: "пар", Cells: []any{0}},
						}},
						{Caption: "x = 7, остаток 3. Нужен остаток 1, таких чисел нет. Затем cnt[3] = 1.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{4, 7, 6, 1, 8, -3, 2}, Marked: []int{1}},
							{Label: "a % 4", Cells: []any{0, 3, "", "", "", "", ""}},
							{Label: "cnt до шага", Cells: []any{1, 0, 0, 0}, Marked: []int{1}},
							{Label: "пар", Cells: []any{0}},
						}},
						{Caption: "x = 6, остаток 2. Нужен остаток 2, таких чисел нет. Затем cnt[2] = 1.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{4, 7, 6, 1, 8, -3, 2}, Marked: []int{2}},
							{Label: "a % 4", Cells: []any{0, 3, 2, "", "", "", ""}},
							{Label: "cnt до шага", Cells: []any{1, 0, 0, 1}, Marked: []int{2}},
							{Label: "пар", Cells: []any{0}},
						}},
						{Caption: "x = 1, остаток 1. Нужен остаток 3: пара 7 + 1 = 8. Пар 1.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{4, 7, 6, 1, 8, -3, 2}, Marked: []int{3}},
							{Label: "a % 4", Cells: []any{0, 3, 2, 1, "", "", ""}},
							{Label: "cnt до шага", Cells: []any{1, 0, 1, 1}, Marked: []int{3}},
							{Label: "пар", Cells: []any{1}},
						}},
						{Caption: "x = 8, остаток 0. Нужен остаток 0: пара 4 + 8 = 12. Пар 2.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{4, 7, 6, 1, 8, -3, 2}, Marked: []int{4}},
							{Label: "a % 4", Cells: []any{0, 3, 2, 1, 0, "", ""}},
							{Label: "cnt до шага", Cells: []any{1, 1, 1, 1}, Marked: []int{0}},
							{Label: "пар", Cells: []any{2}},
						}},
						{Caption: "x = −3, остаток 1: −3 = 4 · (−1) + 1. Нужен остаток 3: пара 7 + (−3) = 4. Пар 3.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{4, 7, 6, 1, 8, -3, 2}, Marked: []int{5}},
							{Label: "a % 4", Cells: []any{0, 3, 2, 1, 0, 1, ""}},
							{Label: "cnt до шага", Cells: []any{2, 1, 1, 1}, Marked: []int{3}},
							{Label: "пар", Cells: []any{3}},
						}},
						{Caption: "x = 2, остаток 2. Нужен остаток 2: пара 6 + 2 = 8. Пар 4: 7 шагов вместо 21 проверки, в конце cnt = [2, 2, 2, 1].", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{4, 7, 6, 1, 8, -3, 2}, Marked: []int{6}},
							{Label: "a % 4", Cells: []any{0, 3, 2, 1, 0, 1, 2}},
							{Label: "cnt до шага", Cells: []any{2, 2, 1, 1}, Marked: []int{2}},
							{Label: "пар", Cells: []any{4}},
						}},
					},
				},
				{
					ID:    "ma1-prefix",
					Title: "Остатки префиксов: Subarray с суммой, кратной 3",
					Frames: []sets.FigureFrame{
						{Caption: "Пустой префикс p[0] = 0 с остатком 0 записан до прохода.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{-2, 5, 1, -4, 3}},
							{Label: "p", Cells: []any{0, "", "", "", "", ""}},
							{Label: "p % 3", Cells: []any{0, "", "", "", "", ""}, Marked: []int{0}},
							{Label: "cnt после шага", Cells: []any{1, 0, 0}, Marked: []int{0}},
							{Label: "ответ", Cells: []any{0}},
						}},
						{Caption: "x = −2: p[1] = −2, остаток 1. Раньше остатка 1 не было: ответ 0.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{-2, 5, 1, -4, 3}, Marked: []int{0}},
							{Label: "p", Cells: []any{0, -2, "", "", "", ""}},
							{Label: "p % 3", Cells: []any{0, 1, "", "", "", ""}, Marked: []int{1}},
							{Label: "cnt после шага", Cells: []any{1, 1, 0}, Marked: []int{1}},
							{Label: "ответ", Cells: []any{0}},
						}},
						{Caption: "x = 5: p[2] = 3, остаток 0, как у p[0]. Subarray [−2, 5] с суммой 3. Ответ 1.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{-2, 5, 1, -4, 3}, Marked: []int{1}},
							{Label: "p", Cells: []any{0, -2, 3, "", "", ""}},
							{Label: "p % 3", Cells: []any{0, 1, 0, "", "", ""}, Marked: []int{2}},
							{Label: "cnt после шага", Cells: []any{2, 1, 0}, Marked: []int{0}},
							{Label: "ответ", Cells: []any{1}},
						}},
						{Caption: "x = 1: p[3] = 4, остаток 1, как у p[1] = −2. Subarray [5, 1] с суммой 6. Ответ 2.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{-2, 5, 1, -4, 3}, Marked: []int{2}},
							{Label: "p", Cells: []any{0, -2, 3, 4, "", ""}},
							{Label: "p % 3", Cells: []any{0, 1, 0, 1, "", ""}, Marked: []int{3}},
							{Label: "cnt после шага", Cells: []any{2, 2, 0}, Marked: []int{1}},
							{Label: "ответ", Cells: []any{2}},
						}},
						{Caption: "x = −4: p[4] = 0, остаток 0. Два прежних префикса с остатком 0: ответ 4.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{-2, 5, 1, -4, 3}, Marked: []int{3}},
							{Label: "p", Cells: []any{0, -2, 3, 4, 0, ""}},
							{Label: "p % 3", Cells: []any{0, 1, 0, 1, 0, ""}, Marked: []int{4}},
							{Label: "cnt после шага", Cells: []any{3, 2, 0}, Marked: []int{0}},
							{Label: "ответ", Cells: []any{4}},
						}},
						{Caption: "x = 3: p[5] = 3, остаток 0. Три прежних префикса: ответ 7. Класс 0 из 4 префиксов дал 6 Subarray, класс 1 из 2 префиксов дал 1.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{-2, 5, 1, -4, 3}, Marked: []int{4}},
							{Label: "p", Cells: []any{0, -2, 3, 4, 0, 3}},
							{Label: "p % 3", Cells: []any{0, 1, 0, 1, 0, 0}, Marked: []int{5}},
							{Label: "cnt после шага", Cells: []any{4, 2, 0}, Marked: []int{0}},
							{Label: "ответ", Cells: []any{7}},
						}},
					},
				},
			},
		},
		{
			Title:   "Binary Exponentiation",
			Summary: "Степень по модулю за O(log b) умножений вместо b. Разбираем шаги цикла, встроенную pow(a, b, m), циклы последних цифр и показатель, записанный цифрами.",
			Lesson:  lesson2,
			Tasks: []sets.Task{
				{
					URL:        "https://codeforces.com/problemset/problem/1594/B",
					Title:      "Special Numbers",
					Difficulty: "easy",
					Topics:     []string{"Math", "Bit Manipulation"},
					Note:       "Codeforces, рейтинг 1100. Особые числа соответствуют двоичной записи k: единица в разряде i добавляет n в степени i. Идите по битам k, как в Binary Exponentiation, и храните p = n в степени i по модулю 10⁹ + 7. Время: O(log k) на набор, до 30 шагов при k ≤ 10⁹. Память: O(1). Для n = 3, k = 4 (двоичная запись 100) ответ 3² = 9. Brute Force по особым числам по порядку до k-го делает до 10⁹ шагов на набор.",
					Statement: `<p>Теофанис очень любит последовательности целых положительных чисел, поэтому его учитель (Yeltsa Kcir) дал ему задачу о последовательности, которая состоит только из особых чисел.</p><p>Назовём положительное число <strong>особым</strong>, если его можно записать как сумму различных неотрицательных степеней числа <code>n</code>. Например, для <code>n = 4</code> число <code>17</code> особое, потому что его можно записать как <code>4<sup>0</sup> + 4<sup>2</sup> = 1 + 16 = 17</code>, а число <code>9</code> не особое.</p><p>Теофанис просит вас помочь ему найти <code>k</code>-е особое число, если все особые числа выписаны по возрастанию. Это число может быть очень большим, поэтому выведите его остаток от деления на <code>10<sup>9</sup> + 7</code>.</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>): количество наборов входных данных.</p><p>Единственная строка каждого набора содержит два целых числа <code>n</code> и <code>k</code> (<code>2 ≤ n ≤ 10<sup>9</sup></code>; <code>1 ≤ k ≤ 10<sup>9</sup></code>).</p><h3>Выходные данные</h3><p>Для каждого набора выведите одно целое число: <code>k</code>-е по возрастанию особое число по модулю <code>10<sup>9</sup> + 7</code>.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>3
3 4
2 12
105 564</pre><p><strong>Выходные данные:</strong></p><pre>9
12
3595374</pre><h3>Примечание</h3><p>Для <code>n = 3</code> последовательность особых чисел начинается так: <code>[1, 3, 4, 9, ...]</code>.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>2 ≤ n ≤ 10<sup>9</sup></code></li><li><code>1 ≤ k ≤ 10<sup>9</sup></code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `3
3 4
2 12
105 564
`, Expected: text(`9
12
3595374
`)},
						{Input: `1
2 1
`, Expected: text(`1
`)},
						{Input: `1
2 1000000000
`, Expected: text(`1000000000
`)},
						{Input: `2
10 7
10 8
`, Expected: text(`111
1000
`)},
						{Input: `1
1000000000 1000000000
`, Expected: text(`165493999
`)},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/double-modular-exponentiation/",
					Number:     number(2961),
					Title:      "Double Modular Exponentiation",
					Difficulty: "medium",
					Topics:     []string{"Array", "Math", "Simulation"},
					Note:       "Для каждой четвёрки сначала x = pow(a, b, 10), затем pow(x, c, m). q - длина variables, B - наибольший показатель. Время: O(q log B) умножений чисел не больше 1000. Память: O(1) дополнительной, ответ до q номеров. Для [2, 3, 3, 10]: 2³ % 10 = 8, 8³ % 10 = 2. Запись a ** b строит число до 3001 цифры; при этих ограничениях это ещё быстро, но время и память растут вместе с b.",
					Statement: `<p>Дан двумерный Array <code>variables</code>, номера в нём <strong>начинаются с 0</strong>: <code>variables[i] = [a<sub>i</sub>, b<sub>i</sub>, c<sub>i</sub>, m<sub>i</sub>]</code>. Дано также целое число <code>target</code>.</p><p>Номер <code>i</code> называется <strong>хорошим</strong>, если выполнены условия:</p><ul><li><code>0 ≤ i &lt; variables.length</code></li><li><code>((a<sub>i</sub><sup>b<sub>i</sub></sup> % 10)<sup>c<sub>i</sub></sup>) % m<sub>i</sub> == target</code></li></ul><p>Верните Array из <strong>хороших</strong> номеров <strong>в любом порядке</strong>.</p><p>Метод класса <code>Solution</code>: <code>def getGoodIndices(self, variables: list[list[int]], target: int) -&gt; list[int]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> variables = [[2,3,3,10],[3,3,3,1],[6,1,1,4]], target = 2
<strong>Выход:</strong> [0,2]
<strong>Пояснение:</strong> для каждого номера i в Array variables:
1) для номера 0, variables[0] = [2,3,3,10], (2<sup>3</sup> % 10)<sup>3</sup> % 10 = 2.
2) для номера 1, variables[1] = [3,3,3,1], (3<sup>3</sup> % 10)<sup>3</sup> % 1 = 0.
3) для номера 2, variables[2] = [6,1,1,4], (6<sup>1</sup> % 10)<sup>1</sup> % 4 = 2.
Поэтому ответ [0,2].</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> variables = [[39,3,1000,1000]], target = 17
<strong>Выход:</strong> []
<strong>Пояснение:</strong> для каждого номера i в Array variables:
1) для номера 0, variables[0] = [39,3,1000,1000], (39<sup>3</sup> % 10)<sup>1000</sup> % 1000 = 1.
Поэтому ответ [].</pre><h3>Ограничения</h3><ul><li><code>1 ≤ variables.length ≤ 100</code></li><li><code>variables[i] == [a<sub>i</sub>, b<sub>i</sub>, c<sub>i</sub>, m<sub>i</sub>]</code></li><li><code>1 ≤ a<sub>i</sub>, b<sub>i</sub>, c<sub>i</sub>, m<sub>i</sub> ≤ 10<sup>3</sup></code></li><li><code>0 ≤ target ≤ 10<sup>3</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "getGoodIndices", Order: "any"}, Cases: []sets.TestCase{
						{Input: "[[[2,3,3,10],[3,3,3,1],[6,1,1,4]],2]", Expected: text("[0,2]")},
						{Input: "[[[39,3,1000,1000]],17]", Expected: text("[]")},
						{Input: "[[[5,5,5,1],[7,1,1,1]],0]", Expected: text("[0,1]")},
						{Input: "[[[1000,1000,1000,1000]],0]", Expected: text("[0]")},
						{Input: "[[[7,1000,1000,13],[3,999,1000,7]],1]", Expected: text("[0]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/count-good-numbers/",
					Number:     number(1922),
					Title:      "Count Good Numbers",
					Difficulty: "medium",
					Topics:     []string{"Math", "Product Rule"},
					Note:       "На чётную позицию подходят 5 чётных цифр, на нечётную 4 простые цифры. Product Rule: pow(5, (n + 1) // 2, M) · pow(4, n // 2, M) % M, M = 10⁹ + 7. Время: O(log n) умножений, около 50 шагов на степень при n = 10¹⁵. Память: O(1). Для n = 3: 5 · 4 · 5 = 100. Цикл из n умножений при n = 10¹⁵ не закончится за 3 с.",
					Statement: `<p>String из цифр называется <strong>хорошей</strong>, если цифры на <strong>чётных</strong> позициях (номера позиций <strong>начинаются с 0</strong>) <strong>чётные</strong>, а цифры на <strong>нечётных</strong> позициях <strong>простые</strong> (<code>2</code>, <code>3</code>, <code>5</code> или <code>7</code>).</p><ul><li>Например, <code>"2582"</code> хорошая: цифры <code>2</code> и <code>8</code> на чётных позициях чётные, а цифры <code>5</code> и <code>2</code> на нечётных позициях простые. А <code>"3245"</code> <strong>не</strong> хорошая: цифра <code>3</code> стоит на чётной позиции, но не чётная.</li></ul><p>Дано целое число <code>n</code>. Верните <strong>общее</strong> число хороших String из цифр длины <code>n</code>. Ответ может быть большим, поэтому <strong>верните его остаток от деления на</strong> <code>10<sup>9</sup> + 7</code>.</p><p><strong>String из цифр</strong> состоит из цифр от <code>0</code> до <code>9</code> и может начинаться с нулей.</p><p>Метод класса <code>Solution</code>: <code>def countGoodNumbers(self, n: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> n = 1
<strong>Выход:</strong> 5
<strong>Пояснение:</strong> хорошие String длины 1: "0", "2", "4", "6", "8".</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> n = 4
<strong>Выход:</strong> 400</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> n = 50
<strong>Выход:</strong> 564908303</pre><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 10<sup>15</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "countGoodNumbers"}, Cases: []sets.TestCase{
						{Input: "[1]", Expected: text("5")},
						{Input: "[4]", Expected: text("400")},
						{Input: "[50]", Expected: text("564908303")},
						{Input: "[2]", Expected: text("20")},
						{Input: "[3]", Expected: text("100")},
						{Input: "[1000000000000000]", Expected: text("711414395")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/super-pow/",
					Number:     number(372),
					Title:      "Super Pow",
					Difficulty: "medium",
					Topics:     []string{"Math", "Divide and Conquer"},
					Note:       "Показатель записан цифрами: если b = 10x + d, то a в степени b равно (a в степени x) в степени 10, умноженному на a в степени d. Идите по цифрам слева направо: r = pow(r, 10, 1337) · pow(a, d, 1337) % 1337. L - длина b. Время: O(L), на цифру до 8 шагов Binary Exponentiation. Память: O(1). Для a = 2, b = [1, 0]: r = 2, затем 2¹⁰ = 1024. Перевод b в одно число Python и pow тоже работает: это около 6644 двоичных цифр показателя.",
					Statement: `<p>Вычислите <code>a<sup>b</sup> % 1337</code>, где <code>a</code> - положительное целое число, а <code>b</code> - очень большое положительное целое число, записанное как Array цифр.</p><p>Метод класса <code>Solution</code>: <code>def superPow(self, a: int, b: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> a = 2, b = [3]
<strong>Выход:</strong> 8
<strong>Пояснение:</strong> Array b = [3] задаёт показатель 3. Поэтому a<sup>3</sup> = 2<sup>3</sup> = 8, и ответ 8.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> a = 2, b = [1,0]
<strong>Выход:</strong> 1024
<strong>Пояснение:</strong> Array b = [1, 0] задаёт показатель 10. Поэтому a<sup>10</sup> = 2<sup>10</sup> = 1024, и ответ 1024.</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> a = 1, b = [4,3,3,8,5,2]
<strong>Выход:</strong> 1
<strong>Пояснение:</strong> Array b = [4, 3, 3, 8, 5, 2] задаёт положительный показатель. Так как a = 1, любая положительная степень 1 равна 1. Поэтому ответ 1.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ a ≤ 2<sup>31</sup> - 1</code></li><li><code>1 ≤ b.length ≤ 2000</code></li><li><code>0 ≤ b[i] ≤ 9</code></li><li>в <code>b</code> нет ведущих нулей</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "superPow"}, Cases: []sets.TestCase{
						{Input: "[2,[3]]", Expected: text("8")},
						{Input: "[2,[1,0]]", Expected: text("1024")},
						{Input: "[1,[4,3,3,8,5,2]]", Expected: text("1")},
						{Input: "[1337,[1]]", Expected: text("0")},
						{Input: "[2147483647,[2,0,0]]", Expected: text("1198")},
						{Input: "[2,[1,0,0,0]]", Expected: text("345")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1236/B",
					Title:      "Alice and the List of Presents",
					Difficulty: "medium",
					Topics:     []string{"Math", "Combinatorics", "Product Rule"},
					Note:       "Codeforces, рейтинг 1500. Каждый вид подарка выбирает непустой набор коробок: 2ᵐ − 1 вариантов. По Product Rule ответ (2ᵐ − 1)ⁿ. Время: две pow, O(log n + log m) умножений, около 60 шагов. Память: O(1). Для n = 2, m = 2: 3² = 9. Вычитание по модулю: (pow(2, m, M) − 1) % M. Цикл из n или m умножений делает до 10⁹ шагов и не уложится в 3 с.",
					Statement:  `<p>В последние дни Алиса получила много подарков. Она решила упаковать их в коробки и отправить друзьям.</p><p>Есть <code>n</code> видов подарков. Подарки одного вида одинаковые (два подарка одного вида нельзя различить). Подарки разных видов различны (два подарка разных видов можно различить). Подарков каждого вида у Алисы очень много, поэтому можно считать, что у неё бесконечно много подарков каждого вида.</p><p>Также есть <code>m</code> коробок. Все они для разных людей, поэтому все коробки различны (считайте, что на коробках написаны имена <code>m</code> друзей). Например, такие два способа различны: подарок первого вида лежит в первой коробке, но не во второй; подарок первого вида лежит во второй коробке, но не в первой.</p><p>Алиса хочет упаковать подарки по таким правилам:</p><ul><li>в одну коробку она не положит больше одного подарка каждого вида, поэтому в каждой коробке подарки разных видов (каждая коробка содержит Subset из <code>n</code> видов, пустые коробки разрешены);</li><li>подарок каждого вида должен лежать хотя бы в одной коробке.</li></ul><p>Теперь Алиса хочет знать, сколько существует разных способов упаковать подарки. Помогите ей посчитать это число. Ответ может быть очень большим, поэтому выведите его по модулю <code>10<sup>9</sup> + 7</code>.</p><p>Смотрите примеры и пояснения к ним.</p><h3>Входные данные</h3><p>В первой строке записаны два целых числа <code>n</code> и <code>m</code> через пробел (<code>1 ≤ n, m ≤ 10<sup>9</sup></code>): число видов подарков и число коробок у Алисы.</p><h3>Выходные данные</h3><p>Выведите одно целое число: число способов упаковать подарки по правилам Алисы, по модулю <code>10<sup>9</sup> + 7</code>.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>1 3</pre><p><strong>Выходные данные:</strong></p><pre>7</pre><h3>Пример 2</h3><p><strong>Входные данные:</strong></p><pre>2 2</pre><p><strong>Выходные данные:</strong></p><pre>9</pre><h3>Примечание</h3><p>В первом примере семь способов упаковать подарки:</p><ul><li><code>{1}{}{}</code></li><li><code>{}{1}{}</code></li><li><code>{}{}{1}</code></li><li><code>{1}{1}{}</code></li><li><code>{}{1}{1}</code></li><li><code>{1}{}{1}</code></li><li><code>{1}{1}{1}</code></li></ul><p>Во втором примере девять способов упаковать подарки:</p><ul><li><code>{}{1,2}</code></li><li><code>{1}{2}</code></li><li><code>{1}{1,2}</code></li><li><code>{2}{1}</code></li><li><code>{2}{1,2}</code></li><li><code>{1,2}{}</code></li><li><code>{1,2}{1}</code></li><li><code>{1,2}{2}</code></li><li><code>{1,2}{1,2}</code></li></ul><p>Например, способ <code>{2}{2}</code> неверный: подарок первого вида должен лежать хотя бы в одной коробке.</p><h3>Ограничения</h3><ul><li><code>1 ≤ n, m ≤ 10<sup>9</sup></code></li><li>Ограничение по времени: 1 секунда на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `1 3
`, Expected: text(`7
`)},
						{Input: `2 2
`, Expected: text(`9
`)},
						{Input: `1 1
`, Expected: text(`1
`)},
						{Input: `5 2
`, Expected: text(`243
`)},
						{Input: `1000000000 1000000000
`, Expected: text(`751201557
`)},
						{Input: `1 500000003
`, Expected: text(`0
`)},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "ma2-loop-cost", Type: "choice", Prompt: "Простой цикл r = r * a % m выполняется b раз. Сколько он работает при b = 10¹², если Python делает около 10⁷ таких шагов в секунду?", Options: []string{"около 0,1 с", "около 10⁵ с, больше суток", "около 40 с"}, Answer: 1, Explain: "10¹² / 10⁷ = 10⁵ секунд, а в сутках 86 400 с: это больше суток. Число 40 - это шаги Binary Exponentiation: у 10¹² 40 двоичных цифр. При 3 с на скрытый тест нужен способ O(log b): до 80 умножений и O(1) памяти."},
				{ID: "ma2-steps", Type: "short", Prompt: "Сколько раз выполнится тело цикла while в power(a, b, m) из урока при b = 100?", Answer: "7", Explain: "Каждый шаг делит b на 2: 100, 50, 25, 12, 6, 3, 1, затем 0. Это 7 шагов, по числу цифр в двоичной записи 100 = 1100100. Ответ 6 получится, если забыть шаг с b = 1: на нём r ещё умножается на a. Время O(log b), дополнительная память O(1)."},
				{ID: "ma2-mults", Type: "short", Prompt: "Сколько раз при b = 100 выполнится строка r = r * a % m?", Answer: "3", Explain: "Строка работает только при нечётном b: 25, 3 и 1. Это единицы двоичной записи 1100100. Ответ 7 - это число шагов цикла, но при чётном b (100, 50, 12, 6) r не меняется. Ещё 7 умножений a = a * a % m; всего 10 умножений вместо 100 в простом цикле."},
				{ID: "ma2-trace", Type: "choice", Prompt: "power(2, 11, 10). Какое значение у r после второго шага цикла?", Options: []string{"4", "6", "8"}, Answer: 2, Explain: "Шаг 1: b = 11 нечётно, r = 2, a = 4, b = 5. Шаг 2: b = 5 нечётно, r = 2 · 4 = 8, a = 16 % 10 = 6, b = 2. Число 6 - это a после шага, а не r. Весь вызов делает 4 шага и возвращает 8: 2¹¹ = 2048."},
				{ID: "ma2-code", Type: "code", Prompt: "Какое условие стоит на месте ___ в Binary Exponentiation?", Language: "python", Code: `def power(a, b, m):
    r = 1
    a %= m
    while b > 0:
        if ___:
            r = r * a % m
        a = a * a % m
        b //= 2
    return r`, Options: []string{"b % 2 == 1", "a % 2 == 1", "r % 2 == 1"}, Answer: 0, Explain: "Множитель a нужен, когда младшая двоичная цифра b равна 1. Для b = 13 = 1101 это шаги 1, 3 и 4. Условие с a или r не связано с показателем: с a % 2 == 1 вызов power(3, 2, 7) вернёт 3 вместо 2. Время O(log b), дополнительная память O(1)."},
				{ID: "ma2-last-digit", Type: "short", Prompt: "Какая последняя цифра у 3²⁰²⁷?", Answer: "7", Explain: "Последние цифры степеней 3 повторяются через 4: 3, 9, 7, 1. 2027 = 4 · 506 + 3, поэтому цифра та же, что у 3³ = 27, то есть 7. Ответ 1 даёт ошибка номера: остаток 3 означает третью степень, а не число с номером 3 от нуля в списке 3, 9, 7, 1. pow(3, 2027, 10) даёт тот же ответ за 11 шагов: у 2027 одиннадцать двоичных цифр."},
				{ID: "ma2-subtract", Type: "choice", Prompt: "x = pow(3, n, M), y = pow(2, n, M). Какая запись даёт остаток 3ⁿ − 2ⁿ по модулю M от 0 до M − 1?", Options: []string{"x - y", "abs(x - y)", "(x - y) % M"}, Answer: 2, Explain: "Остатки могут стоять в другом порядке: при M = 7, n = 2 получаем x = 2, y = 4, x − y = −2. (−2) % 7 = 5, и правда 9 − 4 = 5. abs даст 2, это неверно. Лишняя операция % стоит O(1)."},
			},
			Links: []sets.Link{
				{Title: "Бинарное возведение в степень", URL: "https://ru.algorithmica.org/cs/algebra/binpow/", Source: "Алгоритмика"},
				{Title: "Быстрое возведение в степень", URL: "https://neerc.ifmo.ru/wiki/index.php?title=Быстрое_возведение_в_степень", Source: "Викиконспекты ИТМО"},
				{Title: "Алгоритмы быстрого возведения в степень", URL: "https://ru.wikipedia.org/wiki/Алгоритмы_быстрого_возведения_в_степень", Source: "Википедия"},
				{Title: "Функция pow", URL: "https://docs.python.org/3/builtins/functions.html#pow", Source: "Python docs (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "ma2-binpow",
					Title: "3¹³ по модулю 7: Binary Exponentiation",
					Frames: []sets.FigureFrame{
						{Caption: "Шаг 1: b = 13 нечётно, r = 1 · 3 = 3. Затем a = 3 · 3 % 7 = 2, b = 6.", Rows: []sets.FigureRow{
							{Label: "b", Cells: []any{13, "", "", ""}, Marked: []int{0}},
							{Label: "b % 2", Cells: []any{1, "", "", ""}, Marked: []int{0}},
							{Label: "a в начале шага", Cells: []any{3, "", "", ""}, Marked: []int{0}},
							{Label: "r после шага", Cells: []any{3, "", "", ""}, Marked: []int{0}},
						}},
						{Caption: "Шаг 2: b = 6 чётно, r не меняется. a = 2 · 2 % 7 = 4, b = 3.", Rows: []sets.FigureRow{
							{Label: "b", Cells: []any{13, 6, "", ""}, Marked: []int{1}},
							{Label: "b % 2", Cells: []any{1, 0, "", ""}, Marked: []int{1}},
							{Label: "a в начале шага", Cells: []any{3, 2, "", ""}, Marked: []int{1}},
							{Label: "r после шага", Cells: []any{3, 3, "", ""}, Marked: []int{1}},
						}},
						{Caption: "Шаг 3: b = 3 нечётно, r = 3 · 4 % 7 = 5. Затем a = 4 · 4 % 7 = 2, b = 1.", Rows: []sets.FigureRow{
							{Label: "b", Cells: []any{13, 6, 3, ""}, Marked: []int{2}},
							{Label: "b % 2", Cells: []any{1, 0, 1, ""}, Marked: []int{2}},
							{Label: "a в начале шага", Cells: []any{3, 2, 4, ""}, Marked: []int{2}},
							{Label: "r после шага", Cells: []any{3, 3, 5, ""}, Marked: []int{2}},
						}},
						{Caption: "Шаг 4: b = 1 нечётно, r = 5 · 2 % 7 = 3, затем b = 0. Ответ 3. 4 возведения в квадрат и 3 умножения r вместо 13 умножений.", Rows: []sets.FigureRow{
							{Label: "b", Cells: []any{13, 6, 3, 1}, Marked: []int{3}},
							{Label: "b % 2", Cells: []any{1, 0, 1, 1}, Marked: []int{3}},
							{Label: "a в начале шага", Cells: []any{3, 2, 4, 2}, Marked: []int{3}},
							{Label: "r после шага", Cells: []any{3, 3, 5, 3}, Marked: []int{3}},
						}},
					},
				},
				{
					ID:    "ma2-cycle",
					Title: "Последние цифры степеней 7",
					Frames: []sets.FigureFrame{
						{Caption: "n от 1 до 4: цифры 7, 9, 3, 1. Следующая цифра - остаток (цифра · 7) % 10.", Rows: []sets.FigureRow{
							{Label: "n", Cells: []any{1, 2, 3, 4, 5, 6, 7, 8}},
							{Label: "7ⁿ", Cells: []any{7, 49, 343, 2401, 16807, 117649, 823543, 5764801}},
							{Label: "последняя цифра", Cells: []any{7, 9, 3, 1, 7, 9, 3, 1}, Marked: []int{0, 1, 2, 3}, Window: &[2]int{0, 3}},
						}},
						{Caption: "n от 5 до 8: те же цифры, период 4. Для 7²⁰²⁶: 2026 = 4 · 506 + 2, цифра 9.", Rows: []sets.FigureRow{
							{Label: "n", Cells: []any{1, 2, 3, 4, 5, 6, 7, 8}},
							{Label: "7ⁿ", Cells: []any{7, 49, 343, 2401, 16807, 117649, 823543, 5764801}},
							{Label: "последняя цифра", Cells: []any{7, 9, 3, 1, 7, 9, 3, 1}, Marked: []int{4, 5, 6, 7}, Window: &[2]int{4, 7}},
						}},
					},
				},
			},
		},
		{
			Title:   "Деление по модулю",
			Summary: "Деление заменяем умножением на обратный элемент. Малая теорема Ферма даёт его через pow(b, p − 2, p), и одна pow считает ответ-дробь P/Q. Таблицы Factorial и обратных Factorial делят на k! и считают число разных порядков букв.",
			Lesson:  lesson3,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/prime-arrangements/",
					Number:     number(1175),
					Title:      "Prime Arrangements",
					Difficulty: "easy",
					Topics:     []string{"Math", "Factorial", "Permutations"},
					Note:       "Пусть k - число простых чисел от 1 до n. Они стоят на k простых номерах в любом порядке, остальные n − k чисел на остальных номерах. Ответ k! · (n − k)! по модулю 10⁹ + 7. Время: O(n√n) для проверки простоты делителями и O(n) умножений. Память: O(1). Для n = 5: простые 2, 3, 5, ответ 3! · 2! = 12. Factorial считайте с % после каждого умножения: 100! без модуля имеет 158 цифр.",
					Statement: `<p>Верните число Permutations чисел от <code>1</code> до <code>n</code>, в которых простые числа стоят на простых номерах (номера начинаются с 1).</p><p><em>(Напомним: целое число простое, если оно больше 1 и его нельзя записать как произведение двух целых положительных чисел, каждое из которых меньше него.)</em></p><p>Ответ может быть большим, поэтому <strong>верните его остаток от деления на <code>10<sup>9</sup> + 7</code></strong>.</p><p>Метод класса <code>Solution</code>: <code>def numPrimeArrangements(self, n: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> n = 5
<strong>Выход:</strong> 12
<strong>Пояснение:</strong> например, [1,2,5,4,3] подходит, а [5,2,3,4,1] не подходит: простое число 5 стоит на номере 1.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> n = 100
<strong>Выход:</strong> 682289015</pre><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 100</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "numPrimeArrangements"}, Cases: []sets.TestCase{
						{Input: "[5]", Expected: text("12")},
						{Input: "[100]", Expected: text("682289015")},
						{Input: "[1]", Expected: text("1")},
						{Input: "[2]", Expected: text("1")},
						{Input: "[3]", Expected: text("2")},
						{Input: "[10]", Expected: text("17280")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/2008/F",
					Title:      "Sakurako's Box",
					Difficulty: "medium",
					Topics:     []string{"Math", "Number Theory", "Probability and Statistics"},
					Note:       "Codeforces, рейтинг 1400. Пусть S - сумма чисел набора. Сумма произведений пар равна (S² − сумма квадратов) / 2, а пар n(n − 1) / 2. Значит, ответ (S² − сумма квадратов) / (n(n − 1)): умножьте на pow(n(n − 1) % M, M − 2, M), M = 10⁹ + 7. Время: O(n + log M) на набор. Память: O(1) дополнительной при чтении по одному числу. Для [1, 1, 2]: 5/3, ответ 666666673. Проверка всех пар делает n(n − 1)/2 умножений: при n = 200 000 почти 2 · 10¹⁰.",
					Statement: `<p>У Сакурако есть коробка с <code>n</code> шарами. У каждого шара есть значение. Она хочет поспорить с подругой: если подруга случайно вытащит из коробки два шара (это два разных шара, но их значения могут быть равны), то произведение их значений будет равно числу, которое загадала Сакурако.</p><p>У Сакурако учёная степень по теории вероятностей, поэтому она знает, что лучше всего загадать математическое ожидание. Но она забыла, как его считать. Помогите Сакурако: найдите математическое ожидание произведения двух элементов Array.</p><p>Можно показать, что математическое ожидание имеет вид <code>P/Q</code>, где <code>P</code> и <code>Q</code> - целые неотрицательные числа и <code>Q ≠ 0</code>. Выведите <code>P · Q<sup>−1</sup></code> по модулю <code>10<sup>9</sup> + 7</code>, где <code>Q<sup>−1</sup></code> - обратный элемент <code>Q</code>.</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>): количество наборов входных данных.</p><p>В первой строке каждого набора записано одно целое число <code>n</code> (<code>2 ≤ n ≤ 2 · 10<sup>5</sup></code>): число элементов Array.</p><p>Во второй строке набора записаны <code>n</code> целых чисел <code>a<sub>1</sub>, a<sub>2</sub>, ..., a<sub>n</sub></code> (<code>0 ≤ a<sub>i</sub> ≤ 10<sup>9</sup></code>): элементы Array.</p><p>Гарантируется, что сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого набора выведите значение <code>P · Q<sup>−1</sup></code> по модулю <code>10<sup>9</sup> + 7</code>.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>3
3
3 2 3
4
2 2 2 4
5
1 2 3 4 5</pre><p><strong>Выходные данные:</strong></p><pre>7
6
500000012</pre><h3>Примечание</h3><p>В первом наборе подруга Сакурако может вытащить такие пары шаров: <code>(a<sub>1</sub>, a<sub>2</sub>)</code>, <code>(a<sub>1</sub>, a<sub>3</sub>)</code>, <code>(a<sub>2</sub>, a<sub>3</sub>)</code>. Их произведения равны <code>3 · 2 = 6</code>, <code>3 · 3 = 9</code>, <code>3 · 2 = 6</code>, поэтому математическое ожидание равно <code>(6 + 9 + 6) / 3 = 7</code>.</p><p>Во втором наборе подруга Сакурако может вытащить такие пары шаров: <code>(a<sub>1</sub>, a<sub>2</sub>)</code>, <code>(a<sub>1</sub>, a<sub>3</sub>)</code>, <code>(a<sub>1</sub>, a<sub>4</sub>)</code>, <code>(a<sub>2</sub>, a<sub>3</sub>)</code>, <code>(a<sub>2</sub>, a<sub>4</sub>)</code>, <code>(a<sub>3</sub>, a<sub>4</sub>)</code>. Их произведения равны <code>2 · 2 = 4</code>, <code>2 · 2 = 4</code>, <code>2 · 4 = 8</code>, <code>2 · 2 = 4</code>, <code>2 · 4 = 8</code>, <code>2 · 4 = 8</code>, поэтому математическое ожидание равно <code>(4 + 4 + 8 + 4 + 8 + 8) / 6 = 36 / 6 = 6</code>.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>2 ≤ n ≤ 2 · 10<sup>5</sup></code></li><li><code>0 ≤ a<sub>i</sub> ≤ 10<sup>9</sup></code></li><li>сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `3
3
3 2 3
4
2 2 2 4
5
1 2 3 4 5
`, Expected: text(`7
6
500000012
`)},
						{Input: `1
2
0 0
`, Expected: text(`0
`)},
						{Input: `1
2
1000000000 1000000000
`, Expected: text(`49
`)},
						{Input: `1
3
1 1 2
`, Expected: text(`666666673
`)},
						{Input: `1
4
0 1000000000 0 1000000000
`, Expected: text(`166666676
`)},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1569/C",
					Title:      "Jury Meeting",
					Difficulty: "medium",
					Topics:     []string{"Math", "Combinatorics", "Factorial", "Permutations"},
					Note:       "Codeforces, рейтинг 1500. Если наибольшее число встречается хотя бы два раза, подходят все n! порядков. Иначе пусть k - число членов жюри с числом на 1 меньше максимума. Порядок плохой, если член жюри с максимумом идёт после всех этих k членов: таких порядков n!/(k + 1). Ответ n! − n!/(k + 1) по модулю M = 998244353: деление заменяют умножением на обратный элемент k + 1. Время: O(n + log M) на набор. Память: O(n) для чисел набора. Для [2, 2, 3]: 3! − 3!/3 = 4. Остаток n! % M нельзя делить на k + 1 оператором //: для n = 13 и k = 3 получится неверный ответ.",
					Statement: `<p><code>n</code> человек собрались на заседание жюри будущего соревнования. <code>i</code>-й член жюри придумал <code>a<sub>i</sub></code> задач, и члены жюри хотят рассказать их друг другу.</p><p>Сначала жюри выбирает порядок, в котором члены жюри будут рассказывать задачи. Пусть это Permutation <code>p</code> чисел от <code>1</code> до <code>n</code> (Array длины <code>n</code>, в котором каждое целое число от <code>1</code> до <code>n</code> встречается ровно один раз).</p><p>Затем обсуждение идёт так:</p><ul><li>если у члена жюри <code>p<sub>1</sub></code> остались задачи, он рассказывает остальным одну задачу. Иначе его пропускают;</li><li>если у члена жюри <code>p<sub>2</sub></code> остались задачи, он рассказывает остальным одну задачу. Иначе его пропускают;</li><li>...;</li><li>если у члена жюри <code>p<sub>n</sub></code> остались задачи, он рассказывает остальным одну задачу. Иначе его пропускают;</li><li>если у кого-то ещё остались задачи, всё повторяется с начала. Иначе обсуждение заканчивается.</li></ul><p>Permutation <code>p</code> называется <strong>хорошей</strong>, если ни один член жюри не рассказывает две или больше своих задач подряд.</p><p>Посчитайте число хороших Permutations. Ответ может быть очень большим, поэтому выведите его по модулю <code>998 244 353</code>.</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>): количество наборов входных данных.</p><p>В первой строке набора записано одно целое число <code>n</code> (<code>2 ≤ n ≤ 2 · 10<sup>5</sup></code>): число членов жюри.</p><p>Во второй строке набора записаны <code>n</code> целых чисел <code>a<sub>1</sub>, a<sub>2</sub>, ..., a<sub>n</sub></code> (<code>1 ≤ a<sub>i</sub> ≤ 10<sup>9</sup></code>): <code>a<sub>i</sub></code> - число задач, которые придумал <code>i</code>-й член жюри.</p><p>Сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого набора выведите одно целое число: число хороших Permutations по модулю <code>998 244 353</code>.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>4
2
1 2
3
5 5 5
4
1 3 3 7
6
3 4 2 1 3 3</pre><p><strong>Выходные данные:</strong></p><pre>1
6
0
540</pre><h3>Примечание</h3><p>Разберём первый набор примера. Есть две Permutations: <code>p = [1, 2]</code> и <code>p = [2, 1]</code>. Для <code>p = [1, 2]</code> обсуждение идёт так:</p><ol><li>первый член жюри рассказывает задачу;</li><li>второй член жюри рассказывает задачу;</li><li>у первого члена жюри задач не осталось, его пропускают;</li><li>второй член жюри рассказывает задачу.</li></ol><p>Второй член жюри рассказал две задачи подряд, поэтому эта Permutation не хорошая.</p><p>Для <code>p = [2, 1]</code> обсуждение идёт так:</p><ol><li>второй член жюри рассказывает задачу;</li><li>первый член жюри рассказывает задачу;</li><li>второй член жюри рассказывает задачу.</li></ol><p>Эта Permutation хорошая.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>2 ≤ n ≤ 2 · 10<sup>5</sup></code></li><li><code>1 ≤ a<sub>i</sub> ≤ 10<sup>9</sup></code></li><li>сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `4
2
1 2
3
5 5 5
4
1 3 3 7
6
3 4 2 1 3 3
`, Expected: text(`1
6
0
540
`)},
						{Input: `1
2
7 7
`, Expected: text(`2
`)},
						{Input: `1
3
1 1 3
`, Expected: text(`0
`)},
						{Input: `1
3
2 2 3
`, Expected: text(`4
`)},
						{Input: `1
5
1000000000 999999999 1 1 1
`, Expected: text(`60
`)},
						{Input: `1
13
5 4 4 4 1 1 1 1 1 1 1 1 1
`, Expected: text(`677288188
`)},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/count-anagrams/",
					Number:     number(2514),
					Title:      "Count Anagrams",
					Difficulty: "hard",
					Topics:     []string{"Hash Table", "Math", "String", "Combinatorics", "Counting", "Factorial"},
					Note:       "Для слова длины L с частотами букв c1, c2, … число разных порядков букв равно L! / (c1! · c2! · …). Перемножьте ответы слов. Деление замените умножением на inv_fact из таблицы. n - длина s, M = 10⁹ + 7. Время: O(n + log M): таблицы и один проход. Память: O(n) для двух таблиц. Для \"too hot\": 3!/2! = 3 и 3! = 6, ответ 18. Точные Factorial при n = 10⁵ имеют сотни тысяч цифр: это медленнее и требует больше памяти.",
					Statement: `<p>Дана String <code>s</code> из одного или нескольких слов. Соседние слова разделены одним пробелом <code>' '</code>.</p><p>String <code>t</code> называется <strong>анаграммой</strong> String <code>s</code>, если <code>i</code>-е слово <code>t</code> - это Permutation букв <code>i</code>-го слова <code>s</code>.</p><ul><li>Например, <code>"acb dfe"</code> - анаграмма <code>"abc def"</code>, а <code>"def cab"</code> и <code>"adc bef"</code> - нет.</li></ul><p>Верните число <strong>разных анаграмм</strong> <code>s</code>. Ответ может быть очень большим, поэтому верните его <strong>по модулю</strong> <code>10<sup>9</sup> + 7</code>.</p><p>Метод класса <code>Solution</code>: <code>def countAnagrams(self, s: str) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> s = "too hot"
<strong>Выход:</strong> 18
<strong>Пояснение:</strong> вот некоторые анаграммы данной String: "too hot", "oot hot", "oto toh", "too toh" и "too oht".</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> s = "aa"
<strong>Выход:</strong> 1
<strong>Пояснение:</strong> у данной String только одна анаграмма.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ s.length ≤ 10<sup>5</sup></code></li><li><code>s</code> состоит из строчных английских букв и пробелов <code>' '</code></li><li>между соседними словами ровно один пробел</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "countAnagrams"}, Cases: []sets.TestCase{
						{Input: "[\"too hot\"]", Expected: text("18")},
						{Input: "[\"aa\"]", Expected: text("1")},
						{Input: "[\"a\"]", Expected: text("1")},
						{Input: "[\"abc def\"]", Expected: text("36")},
						{Input: "[\"aab bba ccc\"]", Expected: text("9")},
						{Input: "[\"aaaaaaaaaaaaaaabbbbbbbbbbbbbbb\"]", Expected: text("155117520")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "ma3-divide-rem", Type: "choice", Prompt: "18 // 6 = 3, и остаток от деления 3 на 7 тоже 3. Ученик считает (18 % 7) // (6 % 7). Что получится?", Options: []string{"3", "4", "0"}, Answer: 2, Explain: "18 % 7 = 4, 6 % 7 = 6, 4 // 6 = 0. Деление не переносится на остатки, как сумма и произведение. Правильно умножить 4 на обратный элемент 6 по модулю 7: это 6, потому что 6 · 6 = 36 даёт остаток 1. Тогда 4 · 6 = 24 даёт остаток 3. Ответ 4 - только остаток числа 18."},
				{ID: "ma3-inverse", Type: "short", Prompt: "Найдите обратный элемент числа 3 по модулю 11: число x от 1 до 10, для которого 3 · x даёт остаток 1.", Answer: "4", Explain: "3 · 4 = 12 = 11 + 1. Малая теорема Ферма даёт то же: 3⁹ = 19 683 = 11 · 1789 + 4. Brute Force проверяет x от 1 до m − 1 за O(m): для m = 10⁹ + 7 это миллиард проверок, а pow(3, m − 2, m) делает до 60 умножений. Ответ 1 получается из pow(3, 10, 11): показатель 11 − 1 всегда даёт 1, а нужен 11 − 2. Проверка: 3 · 1 = 3, остаток не 1."},
				{ID: "ma3-fermat-exp", Type: "choice", Prompt: "M = 998244353 - простое число. Какой вызов даёт обратный элемент числа b, если b не делится на M?", Options: []string{"pow(b, M - 1, M)", "pow(b, M - 2, M)", "pow(b, M, M)"}, Answer: 1, Explain: "По малой теореме Ферма b в степени M − 1 даёт остаток 1. Значит, b · b в степени M − 2 даёт 1, и это обратный элемент. pow(b, M − 1, M) всегда возвращает 1, а pow(b, M, M) возвращает b % M. Время O(log M), около 60 умножений."},
				{ID: "ma3-no-inverse", Type: "choice", Prompt: "По какому модулю у числа 4 нет обратного элемента?", Options: []string{"10", "7", "9"}, Answer: 0, Explain: "4 · x всегда чётное, а чётное число при делении на 10 даёт чётный остаток, не 1. По модулю 7 обратный элемент 2: 4 · 2 = 8; по модулю 9 он равен 7: 4 · 7 = 28 = 27 + 1. Обратный элемент есть, когда у числа и модуля нет общего делителя больше 1."},
				{ID: "ma3-half", Type: "short", Prompt: "Ответ задачи 1/2. Какое число нужно вывести по модулю 11?", Answer: "6", Explain: "Нужен x, для которого 2 · x даёт остаток 1 при делении на 11: x = 6, потому что 12 = 11 + 1. Через малую теорему Ферма: 2⁹ = 512 = 11 · 46 + 6. По модулю 10⁹ + 7 то же 1/2 записывают как 500000004. Ответы 0 (1 // 2) и 5 (11 // 2) неверны: 2 · 0 и 2 · 5 дают остатки 0 и 10, а не 1."},
				{ID: "ma3-inv-code", Type: "code", Prompt: "Таблица обратных Factorial строится от конца. Что стоит на месте ___?", Language: "python", Code: `inv_fact[N] = pow(fact[N], M - 2, M)
for i in range(N, 0, -1):
    inv_fact[i - 1] = ___`, Options: []string{"inv_fact[i] * (i - 1) % M", "inv_fact[i] // i", "inv_fact[i] * i % M"}, Answer: 2, Explain: "1/(i − 1)! = i / i!, поэтому inv_fact[i − 1] = inv_fact[i] · i. По модулю 7: inv_fact[6] = 6, inv_fact[5] = 6 · 6 % 7 = 1. С множителем i − 1 получится inv_fact[5] = 6 · 5 % 7 = 2 вместо 1, а inv_fact[0] = 0. Целочисленное деление остатков тоже даёт неверные значения. Весь проход: N умножений и одна pow, O(N + log M)."},
				{ID: "ma3-table-cost", Type: "choice", Prompt: "N = 2 · 10⁵, M = 10⁹ + 7. Сколько примерно умножений по модулю нужно, если считать каждый inv_fact[i] отдельным pow(fact[i], M − 2, M)?", Options: []string{"около 10⁷", "около 2 · 10⁵", "около 2 · 10¹⁴"}, Answer: 0, Explain: "Число M − 2 = 1 000 000 005 имеет 30 двоичных цифр, из них 15 единиц. Поэтому каждый pow делает 30 возведений в квадрат и 15 умножений r, всего 45. 2 · 10⁵ вызовов дают около 9 · 10⁶ умножений. Проход назад от одной pow делает около 2 · 10⁵ умножений, в 45 раз меньше. 2 · 10¹⁴ - это Brute Force: кандидаты от 1 до M − 1 для каждого из 2 · 10⁵ чисел."},
				{ID: "ma3-level", Type: "short", Prompt: "Сколько разных String длины 5 можно составить из букв String \"level\", если каждую букву использовать столько раз, сколько она встречается?", Answer: "30", Explain: "Букв 5: l два раза, e два раза, v один раз. Ответ 5! / (2! · 2! · 1!) = 120 / 4 = 30. По модулю деление заменяют умножением на inv_fact[2] два раза. Без деления каждое слово посчитано 4 раза, и ответ 120 неверен."},
			},
			Links: []sets.Link{
				{Title: "Деление по модулю", URL: "https://ru.algorithmica.org/cs/modular/reciprocal/", Source: "Алгоритмика"},
				{Title: "Малая теорема Ферма", URL: "https://ru.wikipedia.org/wiki/Малая_теорема_Ферма", Source: "Википедия"},
				{Title: "Расширенный алгоритм Евклида", URL: "https://ru.algorithmica.org/cs/modular/extended-euclid/", Source: "Алгоритмика"},
				{Title: "Обратный элемент по модулю", URL: "https://cp-algorithms.com/algebra/module-inverse.html", Source: "cp-algorithms (англ.)"},
				{Title: "Factorial и обратные Factorial", URL: "https://usaco.guide/gold/combo", Source: "USACO Guide (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "ma3-fermat",
					Title: "Степени по модулю 7",
					Frames: []sets.FigureFrame{
						{Caption: "3⁶ = 729 = 7 · 104 + 1. Остаток 1, как говорит малая теорема Ферма при p = 7.", Rows: []sets.FigureRow{
							{Label: "k", Cells: []any{1, 2, 3, 4, 5, 6}},
							{Label: "3ᵏ % 7", Cells: []any{3, 2, 6, 4, 5, 1}, Marked: []int{5}},
						}},
						{Caption: "3⁵ % 7 = 5, и 3 · 5 = 15 даёт остаток 1. Значит, 5 - обратный элемент 3: pow(3, 5, 7) = 5.", Rows: []sets.FigureRow{
							{Label: "k", Cells: []any{1, 2, 3, 4, 5, 6}},
							{Label: "3ᵏ % 7", Cells: []any{3, 2, 6, 4, 5, 1}, Marked: []int{4}},
						}},
						{Caption: "У 2 остаток 1 появляется уже при k = 3, и 2⁶ % 7 тоже 1. Обратный элемент 2 равен 2⁵ % 7 = 4: 2 · 4 = 8.", Rows: []sets.FigureRow{
							{Label: "k", Cells: []any{1, 2, 3, 4, 5, 6}},
							{Label: "3ᵏ % 7", Cells: []any{3, 2, 6, 4, 5, 1}},
							{Label: "2ᵏ % 7", Cells: []any{2, 4, 1, 2, 4, 1}, Marked: []int{4}},
						}},
					},
				},
				{
					ID:    "ma3-tables",
					Title: "Таблицы fact и inv_fact по модулю 7",
					Frames: []sets.FigureFrame{
						{Caption: "fact[i] = fact[i − 1] · i % 7, всего 6 умножений. 4! = 24 даёт остаток 3, 5! = 120 даёт остаток 1, 6! = 720 даёт остаток 6.", Rows: []sets.FigureRow{
							{Label: "i", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "fact", Cells: []any{1, 1, 2, 6, 3, 1, 6}, Marked: []int{4, 5, 6}},
							{Label: "inv_fact", Cells: []any{"", "", "", "", "", "", ""}},
						}},
						{Caption: "Одна pow: inv_fact[6] = 6⁵ % 7 = 6. Проверка: 6 · 6 = 36 = 5 · 7 + 1.", Rows: []sets.FigureRow{
							{Label: "i", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "fact", Cells: []any{1, 1, 2, 6, 3, 1, 6}},
							{Label: "inv_fact", Cells: []any{"", "", "", "", "", "", 6}, Marked: []int{6}},
						}},
						{Caption: "Назад: inv_fact[5] = 6 · 6 % 7 = 1, inv_fact[4] = 1 · 5 = 5.", Rows: []sets.FigureRow{
							{Label: "i", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "fact", Cells: []any{1, 1, 2, 6, 3, 1, 6}},
							{Label: "inv_fact", Cells: []any{"", "", "", "", 5, 1, 6}, Marked: []int{4, 5}},
						}},
						{Caption: "inv_fact[3] = 5 · 4 % 7 = 6, inv_fact[2] = 6 · 3 % 7 = 4, inv_fact[1] = 4 · 2 % 7 = 1, inv_fact[0] = 1.", Rows: []sets.FigureRow{
							{Label: "i", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "fact", Cells: []any{1, 1, 2, 6, 3, 1, 6}},
							{Label: "inv_fact", Cells: []any{1, 1, 4, 6, 5, 1, 6}, Marked: []int{0, 1, 2, 3}},
						}},
						{Caption: "fact[i] · inv_fact[i] даёт остаток 1 для каждого i. Всего 12 умножений и одна pow вместо 7 вызовов pow.", Rows: []sets.FigureRow{
							{Label: "i", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "fact", Cells: []any{1, 1, 2, 6, 3, 1, 6}},
							{Label: "inv_fact", Cells: []any{1, 1, 4, 6, 5, 1, 6}},
							{Label: "fact · inv_fact % 7", Cells: []any{1, 1, 1, 1, 1, 1, 1}, Marked: []int{0, 1, 2, 3, 4, 5, 6}},
						}},
					},
				},
			},
		},
	},
}
