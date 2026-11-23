// Package permutations is the learning set "Permutations и Combinations".
package permutations

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

//go:embed chapter5.html
var lesson5 string

func number(value int64) *int64 { return &value }

func text(value string) *string { return &value }

// Set is the content of the learning set.
var Set = sets.LearningSet{
	Slug:    "permutations-and-combinations",
	Title:   "Permutations и Combinations",
	Summary: "Вы учитесь считать варианты: Permutations, Combinations и Subsets. Затем Backtracking строит ответы, убирает повторы и применяет Pruning. В конце считаете Combinations по модулю 10⁹ + 7. Считаем операции, рабочую память и размер ответа.",
	Topics:  []string{"Combinatorics", "Recursion", "Backtracking"},
	Chapters: []sets.Chapter{
		{
			Title:   "Считаем варианты",
			Summary: "Считаем операции и память: Linear Scan, Product Rule, Sum Rule, Factorial, Permutations, Combinations и Subsets. Отличаем число вариантов от размера ответа.",
			Lesson:  lesson1,
			Tasks: []sets.Task{
				{
					URL:        "https://codeforces.com/problemset/problem/1821/A",
					Title:      "Matching",
					Difficulty: "easy",
					Topics:     []string{"Combinatorics", "Product Rule"},
					Note:       "Codeforces, рейтинг 800. n это длина String, q число знаков ?. Product Rule: знак ? на первой позиции даёт 9 вариантов, на остальных по 10; ведущий 0 даёт ответ 0. Время: O(n), один проход и не больше q умножений. Для ?? ответ 9 · 10 = 90. Память: O(1) дополнительной для индекса и счётчика при n ≤ 5, ответ одно число. Это оценка прохода без срезов; копирование суффикса требует O(n) памяти. Проверка всех замен рассматривает до 10^q вариантов, каждый длины n.",
					Statement: `<p><strong>Шаблон числа</strong> это String из цифр и/или вопросительных знаков.</p><p>Положительное (строго больше <code>0</code>) целое число подходит под шаблон, если каждый вопросительный знак в шаблоне можно заменить на цифру так, чтобы получилась десятичная запись этого числа без ведущих нулей.</p><p>Например:</p><ul><li><code>42</code> подходит под <code>4?</code>;</li><li><code>1337</code> подходит под <code>????</code>;</li><li><code>1337</code> подходит под <code>1?3?</code>;</li><li><code>1337</code> подходит под <code>1337</code>;</li><li><code>3</code> не подходит под <code>??</code>;</li><li><code>8</code> не подходит под <code>???8</code>;</li><li><code>1337</code> не подходит под <code>1?7</code>.</li></ul><p>Дан шаблон числа длиной не больше <code>5</code> символов. Посчитайте, сколько положительных (строго больше <code>0</code>) целых чисел подходят под него.</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>t</code> (<code>1 ≤ t ≤ 2 · 10<sup>4</sup></code>): количество наборов входных данных.</p><p>Каждый набор состоит из одной строки. В ней записана String <code>s</code> (<code>1 ≤ |s| ≤ 5</code>) из цифр и/или вопросительных знаков: шаблон числа для этого набора.</p><h3>Выходные данные</h3><p>Для каждого набора выведите одно целое число: количество положительных (строго больше <code>0</code>) целых чисел, которые подходят под шаблон.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>8
??
?
0
9
03
1??7
?5?
9??99</pre><p><strong>Выходные данные:</strong></p><pre>90
9
0
1
0
100
90
100</pre><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 2 · 10<sup>4</sup></code></li><li><code>1 ≤ |s| ≤ 5</code>, <code>s</code> состоит из цифр и символов <code>?</code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 512 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `8
??
?
0
9
03
1??7
?5?
9??99
`, Expected: text(`90
9
0
1
0
100
90
100
`)},
						{Input: `1
?????
`, Expected: text(`90000
`)},
						{Input: `3
?0
0?
10000
`, Expected: text(`9
0
1
`)},
						{Input: `2
?
00
`, Expected: text(`9
0
`)},
						{Input: `1
5
`, Expected: text(`1
`)},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1743/A",
					Title:      "Password",
					Difficulty: "easy",
					Topics:     []string{"Combinatorics", "Combinations"},
					Note:       "Codeforces, рейтинг 800. n это число запрещённых цифр, m = 10 − n. Выбор двух разрешённых цифр даёт C(m, 2), размещение их пар даёт 6 вариантов: ответ 3 · m · (m − 1). Время: O(n) с чтением входа; формула делает 2 вычитания и 2 умножения, O(1). При n = 8 ответ 6. Память: O(1) дополнительной для формулы, ответ одно число; сохранённый ввод учитываем отдельно. Brute Force всех четырёхзначных паролей проверяет фиксированные 10 000 вариантов. Рост 10^d возникает, если длина d тоже меняется.",
					Statement: `<p>Монокарп забыл пароль от своего мобильного телефона. Пароль состоит из <code>4</code> цифр от <code>0</code> до <code>9</code> (пароль может начинаться с цифры <code>0</code>).</p><p>Монокарп помнит, что в пароле было ровно две разные цифры, и каждая из них встречалась в пароле ровно два раза. Ещё Монокарп помнит некоторые цифры, которых точно не было в пароле.</p><p>Посчитайте, сколько разных последовательностей из <code>4</code> цифр могут быть паролем Монокарпа (то есть подходят под все условия о его пароле).</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>t</code> (<code>1 ≤ t ≤ 200</code>): количество наборов входных данных.</p><p>В первой строке каждого набора записано одно целое число <code>n</code> (<code>1 ≤ n ≤ 8</code>): количество цифр, про которые Монокарп помнит, что их не было в пароле.</p><p>Во второй строке записаны <code>n</code> разных целых чисел <code>a<sub>1</sub>, a<sub>2</sub>, ..., a<sub>n</sub></code> (<code>0 ≤ a<sub>i</sub> ≤ 9</code>): цифры, которых не было в пароле. Цифры <code>a<sub>1</sub>, a<sub>2</sub>, ..., a<sub>n</sub></code> даны по возрастанию.</p><h3>Выходные данные</h3><p>Для каждого набора выведите одно целое число: количество разных последовательностей из <code>4</code> цифр, которые подходят под условия.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>2
8
0 1 2 4 5 6 8 9
1
8</pre><p><strong>Выходные данные:</strong></p><pre>6
216</pre><p><strong>Пояснение:</strong> В первом наборе подходят такие пароли: <code>3377</code>, <code>3737</code>, <code>3773</code>, <code>7337</code>, <code>7373</code>, <code>7733</code>.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 200</code></li><li><code>1 ≤ n ≤ 8</code></li><li><code>0 ≤ a<sub>i</sub> ≤ 9</code>, все <code>a<sub>i</sub></code> разные и даны по возрастанию</li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `2
8
0 1 2 4 5 6 8 9
1
8
`, Expected: text(`6
216
`)},
						{Input: `1
7
0 1 2 3 4 5 6
`, Expected: text(`18
`)},
						{Input: `1
2
0 9
`, Expected: text(`168
`)},
						{Input: `3
4
1 3 5 7
5
0 2 4 6 8
6
0 1 2 3 4 5
`, Expected: text(`90
60
36
`)},
						{Input: `1
8
2 3 4 5 6 7 8 9
`, Expected: text(`6
`)},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/629/A",
					Title:      "Far Relative's Birthday Cake",
					Difficulty: "easy",
					Topics:     []string{"Combinatorics", "Combinations"},
					Note:       "Codeforces, рейтинг 800. n это сторона торта. В ряду с k шоколадками k · (k − 1) / 2 пар. Посчитайте шоколадки по строкам и столбцам, затем сложите числа пар. Время: O(n²), n² проверок клеток и 2n вычислений числа пар. При n = 3 это 9 клеток и 6 счётчиков; первый пример даёт 4 пары. Память: O(n) дополнительной для 2n счётчиков, O(1) для ответа. Сохранённая сетка входа занимает O(n²). Прямая проверка всех пар в каждом ряду и столбце делает n²(n − 1) проверок, O(n³).",
					Statement: `<p>Семья Дора собирается отпраздновать день рождения Фамиля Дора. Они любят Фамиля Дора, поэтому хотят сделать ему странный праздничный торт.</p><p>Торт это квадрат <code>n × n</code> из одинаковых квадратиков со стороной <code>1</code>. Каждый квадратик либо пустой, либо на нём лежит одна шоколадка. Родственники купили торт и разложили шоколадки по нему случайно. Счастье Фамиля Дора равно количеству пар клеток с шоколадками, которые лежат в одной строке или в одном столбце торта. Семья хочет узнать, каким будет счастье Фамиля.</p><p>Обратите внимание: каждую пару считают не больше одного раза, так как две разные клетки не могут лежать одновременно в одной строке и в одном столбце.</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>n</code> (<code>1 ≤ n ≤ 100</code>): длина стороны торта.</p><p>Дальше идут <code>n</code> строк, в каждой <code>n</code> символов. Пустая клетка обозначена символом <code>.</code>, клетка с шоколадкой обозначена символом <code>C</code>.</p><h3>Выходные данные</h3><p>Выведите счастье Фамиля Дора, то есть количество пар шоколадок, которые лежат в одной строке или в одном столбце.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>3
.CC
C..
C.C</pre><p><strong>Выходные данные:</strong></p><pre>4</pre><p><strong>Пояснение:</strong> Пронумеруем строки сверху вниз, а столбцы слева направо. В одной строке лежат такие пары шоколадок: <code>(1, 2)</code> и <code>(1, 3)</code>; <code>(3, 1)</code> и <code>(3, 3)</code>. В одном столбце лежат такие пары: <code>(2, 1)</code> и <code>(3, 1)</code>; <code>(1, 3)</code> и <code>(3, 3)</code>.</p><h3>Пример 2</h3><p><strong>Входные данные:</strong></p><pre>4
CC..
C..C
.CC.
.CC.</pre><p><strong>Выходные данные:</strong></p><pre>9</pre><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 100</code></li><li>каждая клетка это <code>.</code> или <code>C</code></li><li>Ограничение по времени: 1 секунда на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `3
.CC
C..
C.C
`, Expected: text(`4
`)},
						{Input: `4
CC..
C..C
.CC.
.CC.
`, Expected: text(`9
`)},
						{Input: `1
.
`, Expected: text(`0
`)},
						{Input: `1
C
`, Expected: text(`0
`)},
						{Input: `3
CCC
CCC
CCC
`, Expected: text(`18
`)},
						{Input: `2
C.
.C
`, Expected: text(`0
`)},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1674/B",
					Title:      "Dictionary",
					Difficulty: "easy",
					Topics:     []string{"Combinatorics", "Sum Rule"},
					Note:       "Codeforces, рейтинг 800. Пусть s = 26 будет размером алфавита, a и b номерами букв слова от 0. Номер равен a · (s − 1) + b + 1 − (b > a): перед каждой первой буквой блок из 25 слов. Для bc получаем 27. Время: O(1) на слово, постоянное число операций; O(t) для t слов. Память: O(1) дополнительной, ответ одно число на слово. Если хранить все t ответов, нужно O(t) памяти ответа. Построение всего словаря требует O(s²) времени и памяти, здесь 650 двухбуквенных слов.",
					Statement: `<p>В берляндском языке каждое слово состоит ровно из двух букв. При этом первая буква слова отличается от второй. Любая пара разных берляндских букв является правильным словом берляндского языка. (Берляндские буквы совпадают со строчными буквами латинского алфавита.)</p><p>В берляндском словаре есть все слова этого языка. Слова в нём идут в <strong>Lexicographic Order</strong>. Точнее, слово <code>a</code> стоит в словаре раньше слова <code>b</code>, если выполнено одно из условий:</p><ul><li>первая буква <code>a</code> меньше первой буквы <code>b</code>;</li><li>первые буквы <code>a</code> и <code>b</code> одинаковые, а вторая буква <code>a</code> меньше второй буквы <code>b</code>.</li></ul><p>Поэтому словарь выглядит так:</p><ul><li>Слово 1: <code>ab</code></li><li>Слово 2: <code>ac</code></li><li>...</li><li>Слово 25: <code>az</code></li><li>Слово 26: <code>ba</code></li><li>Слово 27: <code>bc</code></li><li>...</li><li>Слово 649: <code>zx</code></li><li>Слово 650: <code>zy</code></li></ul><p>Дано слово <code>s</code> берляндского языка. Найдите его номер в словаре.</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>t</code> (<code>1 ≤ t ≤ 650</code>): количество наборов входных данных.</p><p>Каждый набор состоит из одной строки. В ней записана String <code>s</code> ровно из двух разных строчных латинских букв (то есть правильное слово берляндского языка).</p><h3>Выходные данные</h3><p>Для каждого набора выведите одно целое число: номер слова <code>s</code> в словаре.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>7
ab
ac
az
ba
bc
zx
zy</pre><p><strong>Выходные данные:</strong></p><pre>1
2
25
26
27
649
650</pre><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 650</code></li><li><code>s</code> состоит ровно из двух разных строчных латинских букв</li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 512 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `7
ab
ac
az
ba
bc
zx
zy
`, Expected: text(`1
2
25
26
27
649
650
`)},
						{Input: `1
za
`, Expected: text(`626
`)},
						{Input: `2
yz
zy
`, Expected: text(`625
650
`)},
						{Input: `1
mn
`, Expected: text(`313
`)},
						{Input: `1
nm
`, Expected: text(`338
`)},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/number-of-good-pairs/",
					Number:     number(1512),
					Title:      "Number of Good Pairs",
					Difficulty: "easy",
					Topics:     []string{"Array", "Hash Table", "Combinations"},
					Note:       "LeetCode. n это длина Array, u число разных значений. Hash Table хранит частоты: при новом числе прибавьте прежнюю частоту к ответу, затем увеличьте её. Время: ожидаемое O(n) для n чтений и n обновлений частот. Для [1,2,3,1,1,3] это 6 шагов и ответ 4, а Brute Force сравнит 15 пар. Память: O(u) дополнительной, ответ одно число. При заданном диапазоне 1..100 можно использовать 101 счётчик: O(n) времени и O(1) дополнительной памяти.",
					Statement: `<p>Дан Array целых чисел <code>nums</code>. Верните количество <strong>хороших пар</strong>.</p><p>Пара <code>(i, j)</code> называется <em>хорошей</em>, если <code>nums[i] == nums[j]</code> и <code>i &lt; j</code>.</p><p>Метод класса <code>Solution</code>: <code>def numIdenticalPairs(self, nums: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,2,3,1,1,3]
<strong>Выход:</strong> 4
<strong>Пояснение:</strong> Есть 4 хорошие пары: (0,3), (0,4), (3,4), (2,5). Индексы считают с нуля.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [1,1,1,1]
<strong>Выход:</strong> 6
<strong>Пояснение:</strong> Каждая пара в этом Array хорошая.</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> nums = [1,2,3]
<strong>Выход:</strong> 0</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 100</code></li><li><code>1 ≤ nums[i] ≤ 100</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "numIdenticalPairs"}, Cases: []sets.TestCase{
						{Input: "[[1,2,3,1,1,3]]", Expected: text("4")},
						{Input: "[[1,1,1,1]]", Expected: text("6")},
						{Input: "[[1,2,3]]", Expected: text("0")},
						{Input: "[[7]]", Expected: text("0")},
						{Input: "[[100,100,1,100,1]]", Expected: text("4")},
						{Input: "[[2,2,2,2,2,2,2,2,2,2]]", Expected: text("45")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "pc1-product-rule", Type: "choice", Prompt: "Сколько есть четырёхзначных чисел, у которых все цифры чётные?", Options: []string{"500", "256", "625", "20"}, Answer: 0, Explain: "Чётные цифры: 0, 2, 4, 6, 8. Первая цифра не может быть 0: для неё 4 варианта, для остальных по 5. Product Rule: 4 · 5 · 5 · 5 = 500. Это 3 умножения, а не построение 500 чисел. 625 считает числа с нулём впереди, 256 запрещает 0 на всех местах, 20 складывает вместо умножения."},
				{ID: "pc1-sum-rule", Type: "choice", Prompt: "В столовой 3 супа, 4 вторых блюда и 2 десерта. Обед состоит из двух блюд: суп и второе или второе и десерт. Сколько есть разных обедов?", Options: []string{"24", "20", "96", "9"}, Answer: 1, Explain: "Product Rule: суп и второе дают 3 · 4 = 12, второе и десерт дают 4 · 2 = 8. Sum Rule: 12 + 8 = 20. Это 2 умножения и 1 сложение. 24 считает обед из трёх блюд, 9 складывает все блюда, 96 умножает два вида обедов вместо сложения."},
				{ID: "pc1-factorial", Type: "short", Prompt: "5 разных книг ставят в ряд. Одна из них словарь, он должен стоять с краю: первым или последним. Сколько есть расстановок?", Answer: "48", Explain: "Место для словаря: 2 варианта. Factorial считает Permutations остальных книг: 4! = 24. Product Rule даёт 2 · 24 = 48. Цикл factorial(4) из главы делает 4 умножения, затем нужно ещё одно умножение на 2. Все 48 порядков сохранять не требуется."},
				{ID: "pc1-medals", Type: "short", Prompt: "В забеге 8 бегунов, среди них Петя. Золото, серебро и бронзу получат трое. Сколько есть способов раздать медали так, чтобы Петя получил одну из них?", Answer: "126", Explain: "Петя получает одну из 3 медалей: 3 способа. Две другие медали по порядку получают двое из 7 оставшихся: 7 · 6 = 42. Product Rule: 3 · 42 = 126, всего 2 умножения. Проверка: всего 8 · 7 · 6 = 336, без Пети 7 · 6 · 5 = 210, 336 − 210 = 126."},
				{ID: "pc1-team", Type: "code", Prompt: "Что напечатает программа? Посчитайте без запуска.", Language: "python", Code: `count = 0
for i in range(6):
    for j in range(i + 1, 6):
        count += 1
print(count)`, Options: []string{"36", "21", "15", "30"}, Answer: 2, Explain: "j больше i, поэтому каждую пару из n = 6 чисел считаем один раз: C(6, 2) = 15. Цикл делает n(n − 1)/2 прибавлений: O(n²) времени и O(1) дополнительной памяти для индексов и счётчика в модели ограниченных чисел. Формула n · (n − 1) // 2 считает тот же ответ за постоянное число операций. 36 даёт range(6), 30 даёт условие i != j, 21 даёт range(i, 6)."},
				{ID: "pc1-symmetry", Type: "short", Prompt: "В классе 12 учеников. На экскурсию едут 10 из них. Сколько есть способов выбрать, кто едет?", Answer: "66", Explain: "Combinations: выбрать 10, кто едет, равносильно выбору 2, кто остаётся. C(12, 10) = C(12, 2) = 12 · 11 / 2 = 66. Формула из главы использует r = min(10, 2) = 2 итерации и хранит один текущий результат."},
				{ID: "pc1-subsets", Type: "short", Prompt: "В пиццу можно добавить любые из 5 начинок. Пиццу без начинок не продают. Сколько есть разных пицц?", Answer: "31", Explain: "Subsets для 5 начинок: каждую берём или нет, поэтому 2 · 2 · 2 · 2 · 2 = 32. Убираем пустой вариант: 32 − 1 = 31. Цикл от 1 делает 5 умножений на 2 и одно вычитание. Ответ один счётчик, а не 31 набор начинок."},
				{ID: "pc1-order-matters", Type: "choice", Prompt: "Из 10 учеников выбирают троих: один моет доску, второй поливает цветы, третий выносит мусор. Каждый ученик делает не больше одного дела. Сколько есть способов?", Options: []string{"120", "1000", "30", "720"}, Answer: 3, Explain: "Дела разные, поэтому порядок важен: доску моет любой из 10, цветы любой из 9 оставшихся, мусор любой из 8. Product Rule даёт 10 · 9 · 8 = 720: 2 умножения. 120 подходит, если дела одинаковые, 1000 подходит, если один ученик может делать несколько дел: тогда каждое дело может достаться любому из 10."},
			},
			Links: []sets.Link{
				{Title: "Permutations и Combinations: введение в Combinatorics", URL: "https://skillbox.ru/media/code/perestanovki-sochetaniya-i-razmeshcheniya-starterpak-po-kombinatorike-dlya-it/", Source: "Skillbox Медиа"},
				{Title: "Product Rule", URL: "https://ru.wikipedia.org/wiki/Правило_умножения_(комбинаторика)", Source: "Википедия"},
				{Title: "Combinations", URL: "https://ru.wikipedia.org/wiki/Сочетание", Source: "Википедия"},
				{Title: "Combinatorics в Python", URL: "https://habr.com/ru/articles/479816/", Source: "Хабр"},
			},
			Figures: []sets.Figure{
				{
					ID:    "pc1-linear-scan",
					Title: "Linear Scan: от сравнений к O(n)",
					Frames: []sets.FigureFrame{
						{Caption: "Начало: best = 3. Значение a[0] задаёт максимум; сравнений пока 0.", Rows: []sets.FigureRow{
							{Label: "Array", Cells: []any{3, 1, 4, 2}, Marked: []int{0}},
							{Label: "best", Cells: []any{3}},
							{Label: "Сравнения", Cells: []any{0}},
						}},
						{Caption: "i = 1: сравнили 1 > 3. Условие ложно, best остаётся 3. Выполнено 1 сравнение.", Rows: []sets.FigureRow{
							{Label: "Array", Cells: []any{3, 1, 4, 2}, Marked: []int{1}, Pointers: &sets.Pointers{{Name: "i", Cell: 1}}},
							{Label: "best", Cells: []any{3}},
							{Label: "Сравнения", Cells: []any{1}},
						}},
						{Caption: "i = 2: сравнили 4 > 3. Условие истинно, best становится 4. Выполнено 2 сравнения.", Rows: []sets.FigureRow{
							{Label: "Array", Cells: []any{3, 1, 4, 2}, Marked: []int{2}, Pointers: &sets.Pointers{{Name: "i", Cell: 2}}},
							{Label: "best", Cells: []any{4}},
							{Label: "Сравнения", Cells: []any{2}},
						}},
						{Caption: "i = 3: сравнили 2 > 4. Максимум 4. Для n = 4 получили n - 1 = 3 сравнения: O(n). Индекс i и best занимают O(1) рабочей памяти; счётчик нужен только рисунку.", Rows: []sets.FigureRow{
							{Label: "Array", Cells: []any{3, 1, 4, 2}, Marked: []int{3}, Pointers: &sets.Pointers{{Name: "i", Cell: 3}}},
							{Label: "best", Cells: []any{4}},
							{Label: "Сравнения", Cells: []any{3}},
						}},
					},
				},
			},
		},
		{
			Title:   "Backtracking: Subsets, Permutations, Combinations",
			Summary: "Backtracking: выбрать, пойти глубже, отменить. start задаёт возрастающие индексы, used отмечает занятые. Считаем вызовы, проверки, копии и память Call Stack.",
			Lesson:  lesson2,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/letter-case-permutation/",
					Number:     number(784),
					Title:      "Letter Case Permutation",
					Difficulty: "medium",
					Topics:     []string{"String", "Backtracking"},
					Note:       "LeetCode. Пусть n будет длиной String, l числом букв. Буква даёт две ветки, цифра одну. Время: O(n · 2^l) для Backtracking с общим буфером и join в листе: каждый из 2^l ответов содержит n символов. Для a1b2 это 4 ответа и 16 символов. Память: O(n) дополнительной для буфера и Call Stack; O(n · 2^l) для ответа. Цифры тоже копируются. Итеративное расширение префиксов может хранить O(n · 2^l) рабочих символов целого уровня.",
					Statement: `<p>Дана String <code>s</code>. Каждую букву в ней можно отдельно сделать строчной или заглавной. Так получаются другие String.</p><p>Верните Array всех String, которые можно так получить. Порядок в ответе <strong>любой</strong>.</p><p>Метод класса <code>Solution</code>: <code>def letterCasePermutation(self, s: str) -&gt; list[str]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> s = "a1b2"
<strong>Выход:</strong> ["a1b2","a1B2","A1b2","A1B2"]</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> s = "3z4"
<strong>Выход:</strong> ["3z4","3Z4"]</pre><h3>Ограничения</h3><ul><li><code>1 ≤ s.length ≤ 12</code></li><li><code>s</code> состоит из строчных английских букв, заглавных английских букв и цифр</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "letterCasePermutation", Order: "any"}, Cases: []sets.TestCase{
						{Input: `["a1b2"]`, Expected: text(`["A1B2","A1b2","a1B2","a1b2"]`)},
						{Input: `["3z4"]`, Expected: text(`["3Z4","3z4"]`)},
						{Input: `["12345"]`, Expected: text(`["12345"]`)},
						{Input: `["C"]`, Expected: text(`["C","c"]`)},
						{Input: `["aB"]`, Expected: text(`["AB","Ab","aB","ab"]`)},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/subsets/",
					Number:     number(78),
					Title:      "Subsets",
					Difficulty: "medium",
					Topics:     []string{"Array", "Backtracking"},
					Note:       "LeetCode. Пусть n будет числом разных элементов Array, n ≥ 1. Сохраняйте path[:] при каждом вызове. Время: O(n · 2^n): код со start делает 2^n вызовов и копирует n · 2^(n − 1) элементов. При n = 3 это 8 ответов и 12 элементов. Память: O(n) дополнительной для пути и Call Stack; O(n · 2^n) для ответа. Tree с выбором «взять или нет» делает 2^(n + 1) − 1 вызовов. Bitmask Enumeration проверяет n · 2^n битов и сохраняет тот же объём ответа.",
					Statement: `<p>Дан Array целых чисел <code>nums</code>. Все его элементы <strong>разные</strong>. Верните все <em>Subsets</em> этого Array.</p><p>Каждый результат Subsets содержит выбранные элементы исходного Array. Можно не выбрать ни одного элемента.</p><p>В ответе <strong>не должно</strong> быть одинаковых наборов. Порядок в ответе <strong>любой</strong>.</p><p>Метод класса <code>Solution</code>: <code>def subsets(self, nums: list[int]) -&gt; list[list[int]]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,2,3]
<strong>Выход:</strong> [[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [0]
<strong>Выход:</strong> [[],[0]]</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 10</code></li><li><code>-10 ≤ nums[i] ≤ 10</code></li><li>все числа в <code>nums</code> разные</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "subsets", Order: "any-nested"}, Cases: []sets.TestCase{
						{Input: "[[1,2,3]]", Expected: text("[[],[1],[1,2],[1,2,3],[1,3],[2],[2,3],[3]]")},
						{Input: "[[0]]", Expected: text("[[],[0]]")},
						{Input: "[[-1,5]]", Expected: text("[[],[-1],[-1,5],[5]]")},
						{Input: "[[4,-2,0]]", Expected: text("[[],[-2],[-2,0],[-2,0,4],[-2,4],[0],[0,4],[4]]")},
						{Input: "[[1,2,3,4]]", Expected: text("[[],[1],[1,2],[1,2,3],[1,2,3,4],[1,2,4],[1,3],[1,3,4],[1,4],[2],[2,3],[2,3,4],[2,4],[3],[3,4],[4]]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/permutations/",
					Number:     number(46),
					Title:      "Permutations",
					Difficulty: "medium",
					Topics:     []string{"Array", "Backtracking"},
					Note:       "LeetCode. Пусть n будет числом разных элементов Array. Цикл с нуля проверяет used; отмена снимает пометку и удаляет последний элемент пути. Время: O(n · n!) для проверок и копирования n! ответов длины n. При n = 3 код делает 16 вызовов, 30 проверок used и копирует 18 элементов в 6 ответов. Память: O(n) дополнительной для used, пути и Call Stack; O(n · n!) для ответа. Повторный Next Permutation также требует O(n · n!) времени с сохранением всех результатов.",
					Statement: `<p>Дан Array <code>nums</code> из разных целых чисел. Верните все <strong>Permutations</strong> этого Array. Порядок в ответе <strong>любой</strong>.</p><p>Каждый результат Permutations содержит все элементы Array в одном из возможных порядков.</p><p>Метод класса <code>Solution</code>: <code>def permute(self, nums: list[int]) -&gt; list[list[int]]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,2,3]
<strong>Выход:</strong> [[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [0,1]
<strong>Выход:</strong> [[0,1],[1,0]]</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> nums = [1]
<strong>Выход:</strong> [[1]]</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 6</code></li><li><code>-10 ≤ nums[i] ≤ 10</code></li><li>все числа в <code>nums</code> разные</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "permute", Order: "any"}, Cases: []sets.TestCase{
						{Input: "[[1,2,3]]", Expected: text("[[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]")},
						{Input: "[[0,1]]", Expected: text("[[0,1],[1,0]]")},
						{Input: "[[1]]", Expected: text("[[1]]")},
						{Input: "[[-10,10]]", Expected: text("[[-10,10],[10,-10]]")},
						{Input: "[[5,-3,0]]", Expected: text("[[-3,0,5],[-3,5,0],[0,-3,5],[0,5,-3],[5,-3,0],[5,0,-3]]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/combinations/",
					Number:     number(77),
					Title:      "Combinations",
					Difficulty: "medium",
					Topics:     []string{"Backtracking"},
					Note:       "LeetCode. n задаёт диапазон 1..n, k длину ответа. Код главы сохраняет пути длины k. Время: O(V + k · C(n, k)), где V = C(n, 0) + ... + C(n, k) вызовов; копируется k · C(n, k) элементов. Память: O(k) дополнительной для пути и Call Stack; O(k · C(n, k)) для ответа. Pruning по числу оставшихся мест даёт O(k · C(n, k)) времени. При n = k = 4 текущий код делает 16 вызовов, с Pruning 5; оба копируют 4 элемента одного ответа.",
					Statement: `<p>Даны два целых числа <code>n</code> и <code>k</code>. Верните все <strong>Combinations</strong>: наборы из <code>k</code> разных чисел диапазона <code>[1, n]</code>.</p><p>Порядок в ответе <strong>любой</strong>.</p><p>Метод класса <code>Solution</code>: <code>def combine(self, n: int, k: int) -&gt; list[list[int]]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> n = 4, k = 2
<strong>Выход:</strong> [[1,2],[1,3],[1,4],[2,3],[2,4],[3,4]]
<strong>Пояснение:</strong> Combinations для 4 чисел по 2 дают 6 наборов.
Порядок чисел внутри набора не важен: [1,2] и [2,1] задают один набор.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> n = 1, k = 1
<strong>Выход:</strong> [[1]]
<strong>Пояснение:</strong> Combinations для 1 числа по 1 дают один набор.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 20</code></li><li><code>1 ≤ k ≤ n</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "combine", Order: "any-nested"}, Cases: []sets.TestCase{
						{Input: "[4,2]", Expected: text("[[1,2],[1,3],[1,4],[2,3],[2,4],[3,4]]")},
						{Input: "[1,1]", Expected: text("[[1]]")},
						{Input: "[3,3]", Expected: text("[[1,2,3]]")},
						{Input: "[3,1]", Expected: text("[[1],[2],[3]]")},
						{Input: "[5,3]", Expected: text("[[1,2,3],[1,2,4],[1,2,5],[1,3,4],[1,3,5],[1,4,5],[2,3,4],[2,3,5],[2,4,5],[3,4,5]]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/letter-combinations-of-a-phone-number/",
					Number:     number(17),
					Title:      "Letter Combinations of a Phone Number",
					Difficulty: "medium",
					Topics:     []string{"String", "Backtracking"},
					Note:       "LeetCode. Пусть n будет числом цифр, R произведением числа букв их кнопок, R ≤ 4^n. Время: O(nR) при n ≥ 1: общий путь и join в листе записывают nR символов. Для 23 это 13 вызовов, 9 ответов и 18 символов. Память: O(n) дополнительной для пути и Call Stack, O(nR) для ответа; KEYS имеет постоянный размер. Итеративное расширение префиксов может хранить O(nR) рабочих символов. Пустой вход в примере главы даёт [] за O(1) времени и памяти.",
					Statement: `<p>Дана String из цифр от <code>2</code> до <code>9</code> включительно. Верните все String, которые получаются при выборе одной буквы для каждой цифры. Порядок в ответе <strong>любой</strong>.</p><p>Цифры соответствуют буквам так же, как на кнопках телефона (таблица ниже). Цифра 1 не соответствует ни одной букве.</p><table><thead><tr><th>Цифра</th><th>Буквы</th></tr></thead><tbody><tr><td>2</td><td><code>abc</code></td></tr><tr><td>3</td><td><code>def</code></td></tr><tr><td>4</td><td><code>ghi</code></td></tr><tr><td>5</td><td><code>jkl</code></td></tr><tr><td>6</td><td><code>mno</code></td></tr><tr><td>7</td><td><code>pqrs</code></td></tr><tr><td>8</td><td><code>tuv</code></td></tr><tr><td>9</td><td><code>wxyz</code></td></tr></tbody></table><p>Метод класса <code>Solution</code>: <code>def letterCombinations(self, digits: str) -&gt; list[str]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> digits = "23"
<strong>Выход:</strong> ["ad","ae","af","bd","be","bf","cd","ce","cf"]</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> digits = "2"
<strong>Выход:</strong> ["a","b","c"]</pre><h3>Ограничения</h3><ul><li><code>1 ≤ digits.length ≤ 4</code></li><li><code>digits[i]</code> это цифра из диапазона <code>['2', '9']</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "letterCombinations", Order: "any"}, Cases: []sets.TestCase{
						{Input: `["23"]`, Expected: text(`["ad","ae","af","bd","be","bf","cd","ce","cf"]`)},
						{Input: `["2"]`, Expected: text(`["a","b","c"]`)},
						{Input: `["7"]`, Expected: text(`["p","q","r","s"]`)},
						{Input: `["22"]`, Expected: text(`["aa","ab","ac","ba","bb","bc","ca","cb","cc"]`)},
						{Input: `["79"]`, Expected: text(`["pw","px","py","pz","qw","qx","qy","qz","rw","rx","ry","rz","sw","sx","sy","sz"]`)},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "pc2-no-pop", Type: "code", Prompt: "В коде Subsets забыли строку path.pop(). Что вернёт subsets([1, 2])?", Language: "python", Code: `def dfs(start, path):
    res.append(path[:])
    for i in range(start, len(nums)):
        path.append(nums[i])
        dfs(i + 1, path)`, Options: []string{"[[], [1], [1, 2], [2]]", "[[], [1], [1, 2], [1, 2, 2]]", "Программа не остановится"}, Answer: 1, Explain: "После ветки [1, 2] из path не убрали числа 1 и 2. В корне цикл переходит к i = 1 и добавляет 2 к path = [1, 2]. Получается [1, 2, 2]. Программа остановится: каждый вызов сдвигает индекс вправо. Здесь 4 вызова и 6 скопированных элементов; path растёт до 3 элементов."},
				{ID: "pc2-undo", Type: "code", Prompt: "В коде Permutations после path.pop() нет строки used[i] = False. Что вернёт permute([1, 2, 3])?", Language: "python", Code: `def dfs(path):
    if len(path) == len(nums):
        res.append(path[:])
        return
    for i in range(len(nums)):
        if used[i]:
            continue
        used[i] = True
        path.append(nums[i])
        dfs(path)
        path.pop()`, Options: []string{"[[1, 2, 3]]", "Шесть раз [1, 2, 3]", "Программа не остановится", "Все 6 Permutations, как обычно"}, Answer: 0, Explain: "Первая ветка доходит до [1, 2, 3] и помечает все три индекса. Пометки не снимаются, поэтому после возврата каждый индекс занят. Ответ один: [[1, 2, 3]]. Здесь 4 вызова, 9 проверок used и 3 элемента ответа. Меньшее число операций не делает неверный алгоритм подходящим."},
				{ID: "pc2-start", Type: "code", Prompt: "Subsets (LeetCode 78). В коде вызов dfs(i + 1, path) заменили на dfs(start + 1, path). Сколько объектов list окажется в res для nums = [1, 2, 3]?", Language: "python", Code: `nums = [1, 2, 3]
res = []
def dfs(start, path):
    res.append(path[:])
    for i in range(start, len(nums)):
        path.append(nums[i])
        dfs(start + 1, path)
        path.pop()
dfs(0, [])`, Options: []string{"40", "16", "8", "Recursion не закончится"}, Answer: 1, Explain: "Начало цикла зависит только от глубины. Число вызовов по уровням: 1, 3, 6, 6. Всего 16 объектов list; копируется 0 + 3 + 12 + 18 = 33 элемента. Ответ 40 был бы при трёх ветках на каждом уровне, но range(start, ...) становится короче. После [3] код снова берёт индекс 1 и даёт [3, 2]. Recursion закончится: start увеличивается. С dfs(i + 1, path) каждый следующий индекс больше предыдущего: получаем 8 правильных ответов и 12 скопированных элементов."},
				{ID: "pc2-copy", Type: "code", Prompt: "Здесь ответ сохраняют без копии: res.append(path). Что вернёт combine(4, 2)?", Language: "python", Code: `def combine(n, k):
    res = []
    def dfs(start, path):
        if len(path) == k:
            res.append(path)
            return
        for x in range(start, n + 1):
            path.append(x)
            dfs(x + 1, path)
            path.pop()
    dfs(1, [])
    return res`, Options: []string{"Все 6 наборов: [[1, 2], [1, 3], ..., [3, 4]]", "Шесть раз [3, 4]", "Шесть пустых list: [[], [], [], [], [], []]"}, Answer: 2, Explain: "В res лежат шесть ссылок на один list path. Каждый pop отменяет свой append, поэтому после dfs(1, []) path пустой. Все ссылки показывают этот пустой list. Вызовов по-прежнему 11, но 12 элементов ответов код не копирует. В пиковой точке path содержит 2 элемента, Call Stack хранит 3 вызова."},
				{ID: "pc2-used", Type: "short", Prompt: "В коде Permutations (LeetCode 46) убрали две строки: if used[i]: continue. Остальное не меняли. Сколько объектов list вернёт permute([1, 2, 3])?", Answer: "27", Explain: "На каждом из 3 уровней можно взять любое из 3 чисел: 3 · 3 · 3 = 27 ответов, среди них [1, 1, 1] и [2, 1, 2]. Проверка used запрещает повторно брать занятый индекс. Без неё код делает 1 + 3 + 9 + 27 = 40 вызовов и копирует 81 элемент. Путь всё ещё содержит не больше 3 элементов."},
				{ID: "pc2-perm-calls", Type: "short", Prompt: "Permutations (LeetCode 46), nums = [1, 2, 3]. Сколько раз всего вызывается dfs, вместе с первым вызовом dfs([])?", Answer: "16", Explain: "Один вызов соответствует одному узлу Tree. По уровням: 1, 3, 3 · 2 = 6, 6 · 1 = 6. Всего 16. Каждый из 10 нелистовых вызовов проверяет 3 индекса: 30 проверок used. Шесть листьев копируют по 3 элемента: 18. Одновременно активны не больше 4 вызовов в Call Stack."},
				{ID: "pc2-save-where", Type: "code", Prompt: "Ученик сохраняет ответ только при start == len(nums). Что вернёт subsets([1, 2, 3])?", Language: "python", Code: `def dfs(start, path):
    if start == len(nums):
        res.append(path[:])
    for i in range(start, len(nums)):
        path.append(nums[i])
        dfs(i + 1, path)
        path.pop()`, Options: []string{"[[1, 2, 3], [1, 3], [2, 3], [3]]", "[[1, 2, 3]]", "Все 8 Subsets, как обычно"}, Answer: 0, Explain: "start == 3 только там, где последним взято число 3: [1, 2, 3], [1, 3], [2, 3], [3]. Варианты без 3 не сохраняются. Для Subsets нужен ответ в каждом узле Tree. Вызовов остаётся 8, но код сохраняет только 4 ответа и копирует 8 элементов."},
				{ID: "pc2-phone", Type: "short", Prompt: `Letter Combinations of a Phone Number (LeetCode 17), digits = "22". Сколько String в ответе?`, Answer: "9", Explain: "Каждая цифра задаёт свой уровень Tree, на каждом уровне снова 3 буквы abc. Буквы могут повторяться (aa), порядок важен (ab и ba различаются). Product Rule даёт 3 · 3 = 9. Общий путь содержит до 2 символов; 9 операций join записывают 18 символов ответа."},
			},
			Links: []sets.Link{
				{Title: "Backtracking: первые примеры", URL: "https://notes.algoprog.ru/ru/backtrack/01_1_elementary.html", Source: "Заметки П. Калинина (algoprog)"},
				{Title: "Subsets: построение всех вариантов", URL: "https://habr.com/ru/companies/vk/articles/420605/", Source: "Хабр"},
				{Title: "Backtracking", URL: "https://ru.wikipedia.org/wiki/Поиск_с_возвратом", Source: "Википедия"},
				{Title: "Recursion: Subsets, Permutations и N-Queens", URL: "https://usaco.guide/bronze/complete-rec", Source: "USACO Guide (англ.)"},
			},
		},
		{
			Title:   "Повторы и Pruning",
			Summary: "Sorting и пропуск убирают повторные ответы Backtracking. Pruning останавливает бесполезные ветки. Считаем вызовы, проверки, копии и память префиксов.",
			Lesson:  lesson3,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/subsets-ii/",
					Number:     number(90),
					Title:      "Subsets II",
					Difficulty: "medium",
					Topics:     []string{"Array", "Backtracking", "Sorting"},
					Note:       "После Sorting пропускайте равное число только при i > start. n - длина Array, fᵢ - частоты, U = ∏(fᵢ + 1) - число ответов. Время: O(n log n + nU). Память: O(n) рабочая с Call Stack и Sorting; ответ O((n + 1)U). Для [1, 2, 2] получаем 6 ответов и копируем 9 элементов. Построение всех позиционных Subsets до удаления повторов требует O(n2ⁿ) времени.",
					Statement: `<p>Дан Array целых чисел <code>nums</code>. В нём могут быть повторы. Верните <em>все различные Subsets</em> этого Array, включая пустой.</p><p>В ответе <strong>не должно</strong> быть одинаковых Subsets. Порядок в ответе <strong>любой</strong>.</p><p>Метод класса <code>Solution</code>: <code>def subsetsWithDup(self, nums: list[int]) -&gt; list[list[int]]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,2,2]
<strong>Выход:</strong> [[],[1],[1,2],[1,2,2],[2],[2,2]]</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [0]
<strong>Выход:</strong> [[],[0]]</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 10</code></li><li><code>-10 ≤ nums[i] ≤ 10</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "subsetsWithDup", Order: "any-nested"}, Cases: []sets.TestCase{
						{Input: "[[1,2,2]]", Expected: text("[[],[1],[1,2],[1,2,2],[2],[2,2]]")},
						{Input: "[[0]]", Expected: text("[[],[0]]")},
						{Input: "[[2,2,2]]", Expected: text("[[],[2],[2,2],[2,2,2]]")},
						{Input: "[[4,4,4,1,4]]", Expected: text("[[],[1],[1,4],[1,4,4],[1,4,4,4],[1,4,4,4,4],[4],[4,4],[4,4,4],[4,4,4,4]]")},
						{Input: "[[-1,1,-1]]", Expected: text("[[],[-1],[-1,-1],[-1,-1,1],[-1,1],[1]]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/permutations-ii/",
					Number:     number(47),
					Title:      "Permutations II",
					Difficulty: "medium",
					Topics:     []string{"Array", "Backtracking", "Sorting"},
					Note:       "После Sorting берите равные числа по порядку индексов. n - длина Array, U - число разных ответов, V - число посещённых неполных путей. Время: O(n log n + nV + nU), цикл проверяет nV индексов. Память: O(n) рабочая; ответ O(nU). Для [1, 1, 2] это 18 проверок и 9 элементов ответа. Для n равных чисел проверок n² при U = 1; Next Permutation после Sorting обходит ответы за O(nU).",
					Statement: `<p>Дан Array чисел <code>nums</code>. В нём могут быть повторы. Верните все возможные <strong>различные</strong> Permutations. Порядок в ответе <strong>любой</strong>.</p><p>Метод класса <code>Solution</code>: <code>def permuteUnique(self, nums: list[int]) -&gt; list[list[int]]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,1,2]
<strong>Выход:</strong> [[1,1,2],[1,2,1],[2,1,1]]</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [1,2,3]
<strong>Выход:</strong> [[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 8</code></li><li><code>-10 ≤ nums[i] ≤ 10</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "permuteUnique", Order: "any"}, Cases: []sets.TestCase{
						{Input: "[[1,1,2]]", Expected: text("[[1,1,2],[1,2,1],[2,1,1]]")},
						{Input: "[[1,2,3]]", Expected: text("[[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]")},
						{Input: "[[0]]", Expected: text("[[0]]")},
						{Input: "[[2,2,2]]", Expected: text("[[2,2,2]]")},
						{Input: "[[-1,-1,3,3]]", Expected: text("[[-1,-1,3,3],[-1,3,-1,3],[-1,3,3,-1],[3,-1,-1,3],[3,-1,3,-1],[3,3,-1,-1]]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/combination-sum/",
					Number:     number(39),
					Title:      "Combination Sum",
					Difficulty: "medium",
					Topics:     []string{"Array", "Backtracking"},
					Note:       "Кандидаты положительны. Передавайте start = i; после Sorting используйте Pruning через break. m - число кандидатов, T - цель, d = ⌊T/min(candidates)⌋, V - число вызовов, R - число ответов, S - сумма их длин. Время: O(m log m + V + S). Память: O(m + d) рабочая; ответ O(R + S). Для [2, 3, 6, 7], T = 7: 10 вызовов, R = 2, S = 4. Вариант с выходом при remain < 0 делает 28 вызовов.",
					Statement: `<p>Дан Array <strong>разных</strong> положительных целых чисел <code>candidates</code> и целое число <code>target</code>. Верните Array всех <strong>различных Combinations</strong> чисел из <code>candidates</code>, сумма которых равна <code>target</code>. Порядок Combinations в ответе <strong>любой</strong>.</p><p><strong>Одно и то же</strong> число из <code>candidates</code> можно брать <strong>сколько угодно раз</strong>. Два результата различны, если хотя бы одно число встречается в них разное количество раз.</p><p>Тесты составлены так, что для каждого входа различных Combinations с суммой <code>target</code> меньше <code>150</code>.</p><p>Метод класса <code>Solution</code>: <code>def combinationSum(self, candidates: list[int], target: int) -&gt; list[list[int]]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> candidates = [2,3,6,7], target = 7
<strong>Выход:</strong> [[2,2,3],[7]]
<strong>Пояснение:</strong> 
Числа 2 и 3 есть в candidates, и 2 + 2 + 3 = 7. Число 2 можно брать несколько раз.
Число 7 есть в candidates, и 7 = 7.
Других Combinations нет.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> candidates = [2,3,5], target = 8
<strong>Выход:</strong> [[2,2,2,2],[2,3,3],[3,5]]</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> candidates = [2], target = 1
<strong>Выход:</strong> []</pre><h3>Ограничения</h3><ul><li><code>1 ≤ candidates.length ≤ 30</code></li><li><code>2 ≤ candidates[i] ≤ 40</code></li><li>все элементы <code>candidates</code> разные</li><li><code>1 ≤ target ≤ 40</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "combinationSum", Order: "any-nested"}, Cases: []sets.TestCase{
						{Input: "[[2,3,6,7],7]", Expected: text("[[2,2,3],[7]]")},
						{Input: "[[2,3,5],8]", Expected: text("[[2,2,2,2],[2,3,3],[3,5]]")},
						{Input: "[[2],1]", Expected: text("[]")},
						{Input: "[[8,7,4,3],11]", Expected: text("[[3,4,4],[3,8],[4,7]]")},
						{Input: "[[5],15]", Expected: text("[[5,5,5]]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/generate-parentheses/",
					Number:     number(22),
					Title:      "Generate Parentheses",
					Difficulty: "medium",
					Topics:     []string{"String", "Backtracking"},
					Note:       "n - число пар скобок, R - число правильных String. Время: O(nR) с копированием префиксов. Память: O(n²) рабочая у кода с s + символ, потому что Call Stack удерживает все префиксы; ответ O(nR). При n = 3 до первого сохранения жив 21 символ, ответ содержит 30 символов в 5 String. Общий list с append/pop и join в конце пути снижает рабочую память до O(n).",
					Statement: `<p>Дано число <code>n</code>: количество пар скобок. Напишите функцию, которая строит <em>все правильные String из скобок</em> с <code>n</code> парами. У каждой закрывающей скобки должна быть своя открывающая перед ней.</p><p>Метод класса <code>Solution</code>: <code>def generateParenthesis(self, n: int) -&gt; list[str]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> n = 3
<strong>Выход:</strong> ["((()))","(()())","(())()","()(())","()()()"]</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> n = 1
<strong>Выход:</strong> ["()"]</pre><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 8</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "generateParenthesis", Order: "any"}, Cases: []sets.TestCase{
						{Input: "[3]", Expected: text(`["((()))","(()())","(())()","()(())","()()()"]`)},
						{Input: "[1]", Expected: text(`["()"]`)},
						{Input: "[2]", Expected: text(`["(())","()()"]`)},
						{Input: "[4]", Expected: text(`["(((())))","((()()))","((())())","((()))()","(()(()))","(()()())","(()())()","(())(())","(())()()","()((()))","()(()())","()(())()","()()(())","()()()()"]`)},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1097/B",
					Title:      "Petr and a Combination Lock",
					Difficulty: "medium",
					Topics:     []string{"Backtracking", "Subsets"},
					Note:       "Codeforces, рейтинг 1200. n - число поворотов. Backtracking обновляет сумму при каждом выборе знака. Время: O(2ⁿ), до 2ⁿ⁺¹ − 1 вызовов. Память: O(n) для Call Stack; ответ YES/NO занимает O(1). Для [10, 20, 30] достаточно 5 вызовов; для [10, 10, 10] нужны все 15. Dynamic Programming по остаткам даёт O(360n) времени и O(360) памяти.",
					Statement: `<p>Петя только что купил новую машину. Он приехал на самую известную заправку Петербурга, чтобы заправить машину, и вдруг увидел, что бензобак закрыт кодовым замком. На замке есть шкала на <code>360</code> градусов и стрелка. Сначала стрелка указывает на ноль.</p><p>Петя позвонил продавцу машины. Продавец сказал, что нужно повернуть колесо замка ровно <code>n</code> раз. <code>i</code>-й поворот должен быть на <code>a<sub>i</sub></code> градусов, по часовой стрелке или против часовой стрелки. После всех <code>n</code> поворотов стрелка должна снова указывать на ноль.</p><p>Петя немного запутался: он не знает, какие повороты делать по часовой стрелке, а какие против часовой стрелки. Способов повернуть замок много. Помогите ему: узнайте, есть ли хотя бы один способ, при котором после всех <code>n</code> поворотов стрелка снова будет указывать на ноль.</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>n</code> (<code>1 ≤ n ≤ 15</code>): количество поворотов.</p><p>В каждой из следующих <code>n</code> строк записано одно целое число <code>a<sub>i</sub></code> (<code>1 ≤ a<sub>i</sub> ≤ 180</code>): угол <code>i</code>-го поворота в градусах.</p><h3>Выходные данные</h3><p>Если можно сделать все повороты так, чтобы после них стрелка указывала на ноль, выведите одно слово <code>YES</code>. Иначе выведите <code>NO</code>. В этом случае Петя, наверное, купит новую машину.</p><p>Каждую букву ответа можно вывести в верхнем или нижнем регистре.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>3
10
20
30</pre><p><strong>Выходные данные:</strong></p><pre>YES</pre><p><strong>Пояснение:</strong> Первый и второй повороты делаем по часовой стрелке, а третий поворот против часовой стрелки.</p><h3>Пример 2</h3><p><strong>Входные данные:</strong></p><pre>3
10
10
10</pre><p><strong>Выходные данные:</strong></p><pre>NO</pre><p><strong>Пояснение:</strong> Здесь нельзя сделать повороты так, чтобы в конце стрелка указывала на ноль.</p><h3>Пример 3</h3><p><strong>Входные данные:</strong></p><pre>3
120
120
120</pre><p><strong>Выходные данные:</strong></p><pre>YES</pre><p><strong>Пояснение:</strong> Петя может сделать все три поворота по часовой стрелке. Тогда колесо повернётся на <code>360</code> градусов по часовой стрелке, и стрелка снова будет указывать на ноль.</p><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 15</code></li><li><code>1 ≤ a<sub>i</sub> ≤ 180</code></li><li>Ограничение по времени: 1 секунда на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio", IgnoreCase: true}, Cases: []sets.TestCase{
						{Input: `3
10
20
30
`, Expected: text(`YES
`)},
						{Input: `3
10
10
10
`, Expected: text(`NO
`)},
						{Input: `3
120
120
120
`, Expected: text(`YES
`)},
						{Input: `1
180
`, Expected: text(`NO
`)},
						{Input: `15
24
24
24
24
24
24
24
24
24
24
24
24
24
24
24
`, Expected: text(`YES
`)},
						{Input: `15
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
`, Expected: text(`NO
`)},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "pc3-dup-skip", Type: "code", Prompt: "Subsets II, nums = [1, 2, 2]. В правиле пропуска стоит i > 0, а не i > start. Какие Array вернёт код?", Language: "python", Code: `nums.sort()
def dfs(start, path):
    result.append(path[:])
    for i in range(start, len(nums)):
        if i > 0 and nums[i] == nums[i - 1]:
            continue
        path.append(nums[i])
        dfs(i + 1, path)
        path.pop()`, Options: []string{"[], [1], [1, 2], [1, 2, 2], [2], [2, 2]", "[], [1], [1, 2], [1, 2, 2], [2]", "[], [1], [1, 2], [2]", "[], [1], [1, 2], [2], [2, 2]"}, Answer: 2, Explain: "Вторая двойка стоит на месте 2, и nums[2] == nums[1]. С i > 0 код пропускает её на любом уровне, даже когда i == start. Поэтому пропадают [1, 2, 2] и [2, 2]. С i > start первая копия на своём уровне остаётся: 6 Subsets и 9 скопированных элементов ответа."},
				{ID: "pc3-sort-first", Type: "short", Prompt: "Код Subsets II с пропуском (i > start and nums[i] == nums[i - 1]) запустили без nums.sort(). nums = [2, 2, 1, 2]. Сколько результатов он вернёт?", Answer: "12", Explain: "Проверка срабатывает только на первом уровне при i = 1: вторая двойка стоит рядом с первой. Код не начинает путь со второй двойки и убирает 4 результата: [2], [2, 1], [2, 2], [2, 1, 2]. Всего 16 - 4 = 12. Последняя двойка стоит после единицы, и её повторы проверка не видит. После Sorting Array [1, 2, 2, 2] даёт U = (1 + 1)(3 + 1) = 8 результатов."},
				{ID: "pc3-perm-used", Type: "code", Prompt: "Permutations II, nums = [1, 1, 2]. В правиле пропуска нет части not used[i - 1]. Что вернёт код?", Language: "python", Code: `nums.sort()
def dfs(path):
    if len(path) == len(nums):
        result.append(path[:])
        return
    for i in range(len(nums)):
        if used[i]:
            continue
        if i > 0 and nums[i] == nums[i - 1]:
            continue
        used[i] = True
        path.append(nums[i])
        dfs(path)
        path.pop()
        used[i] = False`, Options: []string{"[[1, 1, 2], [1, 2, 1], [2, 1, 1]]", "[[1, 1, 2], [1, 2, 1], [1, 1, 2], [1, 2, 1], [2, 1, 1], [2, 1, 1]]", "[[1, 2, 1], [2, 1, 1]]", "[]"}, Answer: 3, Explain: "Вторая единица стоит на месте 1, и nums[1] == nums[0]. Без not used[i - 1] код пропускает её всегда, даже когда первая единица уже в пути. Ветки [1, 2] и [2, 1] не могут взять третье число. Код вернёт пустой Array. Правило с not used[i - 1] берёт вторую единицу после первой и даёт 3 Permutations, но проверяет 18 индексов."},
				{ID: "pc3-reuse-start", Type: "code", Prompt: "Combination Sum, candidates = [2, 3], цель 7. Глубже передают start = 0, а не i. Сколько результатов окажется в result?", Language: "python", Code: `candidates = [2, 3]
result = []
def dfs(start, remain, path):
    if remain == 0:
        result.append(path[:])
        return
    for i in range(start, len(candidates)):
        if candidates[i] > remain:
            break
        path.append(candidates[i])
        dfs(0, remain - candidates[i], path)
        path.pop()
dfs(0, 7, [])`, Options: []string{"3", "1", "6", "Ни одного: Recursion не закончится"}, Answer: 0, Explain: "С start = 0 после 3 цикл снова начинает с 2. Один набор попадает в result в разных порядках: [2, 2, 3], [2, 3, 2], [3, 2, 2]. Это 3 результата, а не 3! = 6: две двойки одинаковые. Recursion закончится: каждый вызов уменьшает remain хотя бы на 2, а Pruning не допускает отрицательного остатка. С dfs(i, ...) остаётся только [2, 2, 3]."},
				{ID: "pc3-prune", Type: "code", Prompt: "Combination Sum, candidates = [5, 2, 3], цель 7. Ученик забыл candidates.sort(), но оставил break. Что вернёт код?", Language: "python", Code: `def dfs(start, remain, path):
    if remain == 0:
        result.append(path[:])
        return
    for i in range(start, len(candidates)):
        if candidates[i] > remain:
            break
        path.append(candidates[i])
        dfs(i, remain - candidates[i], path)
        path.pop()`, Options: []string{"[[5, 2], [2, 2, 3]]", "[[2, 2, 3]]", "[]", "[[5, 2]]"}, Answer: 1, Explain: "После 5 остаток 2. Цикл снова начинает с i = 0: 5 > 2, и break останавливает цикл. До двойки на месте 1 цикл не доходит, и набор [5, 2] теряется. Остаётся [2, 2, 3]. Pruning через break верен после Sorting: тогда все числа дальше не меньше текущего."},
				{ID: "pc3-paren-branch", Type: "short", Prompt: "Generate Parentheses, n = 4. Код уже построил String «(()(». Сколько готовых String даст эта ветка?", Answer: "3", Explain: "Осталось добавить одну «(» и три «)». Ветка даёт (()(())), (()()()), (()())(): 3 ответа по 8 символов. В варианте (()()))( закрытых становится 4 при 3 открытых. Правило closed < opened не даёт так сделать."},
				{ID: "pc3-paren-bug", Type: "short", Prompt: "Generate Parentheses, n = 2. В коде условие closed < opened заменили на closed < n. Сколько String попадёт в result?", Answer: "6", Explain: "Теперь «(» и «)» не зависят друг от друга: каждую можно взять не больше 2 раз. В String из 4 символов будет ровно 2 «(» и 2 «)». Таких String C(4, 2) = 6: (()), ()(), ())(, )((), )()(, ))((. Правильных из них только 2."},
			},
			Links: []sets.Link{
				{Title: "Backtracking: Pruning", URL: "https://notes.algoprog.ru/ru/backtrack/01_2_cuts.html", Source: "Заметки П. Калинина (algoprog)"},
				{Title: "Balanced Parentheses: определение", URL: "https://ru.wikipedia.org/wiki/Правильная_скобочная_последовательность", Source: "Википедия"},
				{Title: "Как работает Backtracking", URL: "https://thecode.media/backtracking/", Source: "Журнал «Код»"},
				{Title: "Balanced Parentheses: построение и Counting", URL: "https://cp-algorithms.com/combinatorics/bracket_sequences.html", Source: "cp-algorithms (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "pc3-prefix-memory",
					Title: "Generate Parentheses: память префиксов",
					Frames: []sets.FigureFrame{
						{Caption: "n = 3. Перед первым result.append Call Stack удерживает String длины 0, 1, 2, 3, 4, 5, 6. Всего 21 символ. Знак ∅ обозначает пустой String и не добавляет символ. Каждый ряд s относится к отдельному вызову.", Rows: []sets.FigureRow{
							{Label: "s0", Cells: []any{"∅"}},
							{Label: "s1", Cells: []any{"("}},
							{Label: "s2", Cells: []any{"(", "("}},
							{Label: "s3", Cells: []any{"(", "(", "("}},
							{Label: "s4", Cells: []any{"(", "(", "(", ")"}},
							{Label: "s5", Cells: []any{"(", "(", "(", ")", ")"}},
							{Label: "s6", Cells: []any{"(", "(", "(", ")", ")", ")"}},
							{Label: "Символы", Cells: []any{21}},
						}},
						{Caption: "Вариант с общим list: все 7 вызовов используют один path из 6 элементов. join создаёт String ответа из 6 символов. Рабочая память с Call Stack равна O(n), готовые String учитываются отдельно.", Rows: []sets.FigureRow{
							{Label: "path", Cells: []any{"(", "(", "(", ")", ")", ")"}},
							{Label: "Ответ", Cells: []any{"(", "(", "(", ")", ")", ")"}, Marked: []int{0, 1, 2, 3, 4, 5}},
						}},
					},
				},
			},
		},
		{
			Title:   "Порядок и номер",
			Summary: "Next Permutation, Dynamic Programming для путей, Permutation Sequence и N-Queens. Учитываем копии Python, вызовы Factorial и формат ответа.",
			Lesson:  lesson4,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/next-permutation/",
					Number:     number(31),
					Title:      "Next Permutation",
					Difficulty: "medium",
					Topics:     []string{"Permutations", "Lexicographic Order"},
					Note:       "n - длина Array. Время: O(n) на поиск двух индексов и разворот хвоста. Память: O(n) у примера со срезами Python; O(1) при развороте обменами по индексам, как требует задача. Ответ остаётся во входном Array. Для 1 2 5 4 3: 4 сравнения значений и 6 элементов в двух копиях хвоста. Sorting хвоста имеет более слабую общую оценку O(n log n).",
					Statement: `<p><strong>Permutations</strong> Array целых чисел получают, если расставить его элементы во всех возможных порядках.</p><ul><li>Например, у <code>arr = [1,2,3]</code> есть такие Permutations: <code>[1,2,3], [1,3,2], [2,1,3], [2,3,1], [3,1,2], [3,2,1]</code>.</li></ul><p><strong>Next Permutation</strong> Array: следующий вариант в <strong>Lexicographic Order</strong>. Сравнивайте элементы слева направо до первого отличия. Меньший элемент означает меньший вариант. Если большего варианта нет, расположите Array по возрастанию.</p><ul><li>Next Permutation для <code>arr = [1,2,3]</code>: <code>[1,3,2]</code>.</li><li>Next Permutation для <code>arr = [2,3,1]</code>: <code>[3,1,2]</code>.</li><li>Next Permutation для <code>arr = [3,2,1]</code>: <code>[1,2,3]</code>, потому что большего варианта нет.</li></ul><p>Дан Array целых чисел <code>nums</code>. <em>Примените Next Permutation</em> к <code>nums</code>.</p><p>Меняйте Array <strong>на месте</strong> и используйте только постоянную дополнительную память. Метод ничего не возвращает. В примерах «Выход» показывает <code>nums</code> после вызова.</p><p>Метод класса <code>Solution</code>: <code>def nextPermutation(self, nums: list[int]) -&gt; None</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,2,3]
<strong>Выход:</strong> [1,3,2]</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [3,2,1]
<strong>Выход:</strong> [1,2,3]</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> nums = [1,1,5]
<strong>Выход:</strong> [1,5,1]</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 100</code></li><li><code>0 ≤ nums[i] ≤ 100</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "nextPermutation", Check: "arg0"}, Cases: []sets.TestCase{
						{Input: "[[1,2,3]]", Expected: text("[1,3,2]")},
						{Input: "[[3,2,1]]", Expected: text("[1,2,3]")},
						{Input: "[[1,1,5]]", Expected: text("[1,5,1]")},
						{Input: "[[7]]", Expected: text("[7]")},
						{Input: "[[1,5,1]]", Expected: text("[5,1,1]")},
						{Input: "[[2,2,0,4,3,1]]", Expected: text("[2,2,1,0,3,4]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/unique-paths/",
					Number:     number(62),
					Title:      "Unique Paths",
					Difficulty: "medium",
					Topics:     []string{"Combinations", "Dynamic Programming"},
					Note:       "m - число строк сетки, n - число столбцов. Время: O(mn), включая инициализацию n ячеек и (m − 1)(n − 1) сложений. Память: O(n) для одной строки Dynamic Programming; ответ - одно число. Для 3 × 3 это 4 сложения и 3 ячейки, полная таблица хранит 9. Выбор более короткого измерения даёт O(min(m, n)) памяти. Формула Combinations требует O(min(m, n)) арифметических шагов и O(1) рабочих чисел.",
					Statement: `<p>На сетке <code>m x n</code> стоит робот. Сначала робот находится в <strong>левом верхнем углу</strong> (то есть в <code>grid[0][0]</code>). Робот хочет попасть в <strong>правый нижний угол</strong> (то есть в <code>grid[m - 1][n - 1]</code>). В любой момент робот может сделать ход только вниз или вправо.</p><p>Даны два целых числа <code>m</code> и <code>n</code>. Верните <em>количество разных путей, по которым робот может дойти до правого нижнего угла</em>.</p><p>Тесты составлены так, что ответ не больше <code>2 · 10<sup>9</sup></code>.</p><p>Метод класса <code>Solution</code>: <code>def uniquePaths(self, m: int, n: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> m = 3, n = 7
<strong>Выход:</strong> 28</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> m = 3, n = 2
<strong>Выход:</strong> 3
<strong>Пояснение:</strong> Из левого верхнего угла в правый нижний угол ведут 3 пути:
1. Вправо -&gt; Вниз -&gt; Вниз
2. Вниз -&gt; Вниз -&gt; Вправо
3. Вниз -&gt; Вправо -&gt; Вниз</pre><h3>Ограничения</h3><ul><li><code>1 ≤ m, n ≤ 100</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "uniquePaths"}, Cases: []sets.TestCase{
						{Input: "[3,7]", Expected: text("28")},
						{Input: "[3,2]", Expected: text("3")},
						{Input: "[1,1]", Expected: text("1")},
						{Input: "[1,100]", Expected: text("1")},
						{Input: "[10,10]", Expected: text("48620")},
						{Input: "[23,12]", Expected: text("193536720")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/permutation-sequence/",
					Number:     number(60),
					Title:      "Permutation Sequence",
					Difficulty: "hard",
					Topics:     []string{"Permutations", "Factorial"},
					Note:       "n - число цифр, k - номер с 1, F(s) - стоимость factorial(s). Время: O(n²) + Σ F(s), s = 0…n − 1: учитывайте сдвиги pop, рост String и повторные вызовы Factorial. Память: O(n) рабочая; ответ n символов при n ≤ 9. Для n = 4, k = 9 это 4 вызова Factorial и 4 сдвига. Предварительный расчёт Factorial убирает повторную работу, но не квадратичный предел сдвигов.",
					Statement: `<p>У Array <code>[1, 2, 3, ..., n]</code> всего <code>n!</code> разных Permutations.</p><p>Запишем все Permutations в Lexicographic Order и пронумеруем их. Для <code>n = 3</code> получим:</p><ol><li><code>"123"</code></li><li><code>"132"</code></li><li><code>"213"</code></li><li><code>"231"</code></li><li><code>"312"</code></li><li><code>"321"</code></li></ol><p>Даны <code>n</code> и <code>k</code>. Верните вариант с номером <code>k</code> в виде String.</p><p>Метод класса <code>Solution</code>: <code>def getPermutation(self, n: int, k: int) -&gt; str</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> n = 3, k = 3
<strong>Выход:</strong> "213"</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> n = 4, k = 9
<strong>Выход:</strong> "2314"</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> n = 3, k = 1
<strong>Выход:</strong> "123"</pre><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 9</code></li><li><code>1 ≤ k ≤ n!</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "getPermutation"}, Cases: []sets.TestCase{
						{Input: "[3,3]", Expected: text(`"213"`)},
						{Input: "[4,9]", Expected: text(`"2314"`)},
						{Input: "[3,1]", Expected: text(`"123"`)},
						{Input: "[1,1]", Expected: text(`"1"`)},
						{Input: "[3,6]", Expected: text(`"321"`)},
						{Input: "[9,362880]", Expected: text(`"987654321"`)},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/n-queens/",
					Number:     number(51),
					Title:      "N-Queens",
					Difficulty: "hard",
					Topics:     []string{"Backtracking", "Pruning"},
					Note:       "n - размер доски, V - число неполных посещённых расстановок, Q - число решений. Время: O(nV + nQ) в среднем у образца с Hash Table и номерами столбцов; построение досок добавляет O(n²Q). Память: O(n) рабочая для поиска; ответ O(nQ) чисел или O(n²Q) символов досок. Для n = 4: 60 проверок, 2 ответа, 8 чисел или 32 символа. Стройте доску в конечном вызове, чтобы не хранить все промежуточные номера столбцов.",
					Statement: `<p>Задача <strong>N-Queens</strong>: нужно поставить <code>n</code> ферзей на шахматную доску <code>n x n</code> так, чтобы никакие два ферзя не били друг друга (ферзь бьёт по строке, по столбцу и по диагоналям).</p><p>Дано целое число <code>n</code>. Верните <em>все различные решения N-Queens</em>. Порядок в ответе <strong>любой</strong>.</p><p>Каждое решение это отдельная расстановка ферзей на доске. В ней <code>'Q'</code> означает ферзя, а <code>'.'</code> означает пустую клетку.</p><p>Метод класса <code>Solution</code>: <code>def solveNQueens(self, n: int) -&gt; list[list[str]]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> n = 4
<strong>Выход:</strong> [[".Q..","...Q","Q...","..Q."],["..Q.","Q...","...Q",".Q.."]]
<strong>Пояснение:</strong> Для 4 ферзей есть два различных решения:
.Q..   ..Q.
...Q   Q...
Q...   ...Q
..Q.   .Q..</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> n = 1
<strong>Выход:</strong> [["Q"]]</pre><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 9</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "solveNQueens", Order: "any"}, Cases: []sets.TestCase{
						{Input: "[4]", Expected: text(`[["..Q.","Q...","...Q",".Q.."],[".Q..","...Q","Q...","..Q."]]`)},
						{Input: "[1]", Expected: text(`[["Q"]]`)},
						{Input: "[2]", Expected: text("[]")},
						{Input: "[3]", Expected: text("[]")},
						{Input: "[5]", Expected: text(`[["....Q","..Q..","Q....","...Q.",".Q..."],["....Q",".Q...","...Q.","Q....","..Q.."],["...Q.",".Q...","....Q","..Q..","Q...."],["...Q.","Q....","..Q..","....Q",".Q..."],["..Q..","....Q",".Q...","...Q.","Q...."],["..Q..","Q....","...Q.",".Q...","....Q"],[".Q...","....Q","..Q..","Q....","...Q."],[".Q...","...Q.","Q....","..Q..","....Q"],["Q....","...Q.",".Q...","....Q","..Q.."],["Q....","..Q..","....Q",".Q...","...Q."]]`)},
						{Input: "[6]", Expected: text(`[["....Q.","..Q...","Q.....",".....Q","...Q..",".Q...."],["...Q..","Q.....","....Q.",".Q....",".....Q","..Q..."],["..Q...",".....Q",".Q....","....Q.","Q.....","...Q.."],[".Q....","...Q..",".....Q","Q.....","..Q...","....Q."]]`)},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "pc4-next-perm", Type: "short", Prompt: "Какой вариант идёт сразу после 1 3 5 4 2 в Lexicographic Order? Запишите цифры без пробелов.", Answer: "14235", Explain: "Справа налево: хвост 5 4 2 убывает, а 3 < 5. Меняем 3 с самым правым числом больше 3, это 4: 1 4 5 3 2. Хвост 5 3 2 переворачиваем: 1 4 2 3 5. Next Permutation сравнивает 3 соседние пары при поиске i и 2 значения при поиске j."},
				{ID: "pc4-last-perm", Type: "short", Prompt: "Примените Next Permutation два раза подряд к 4 3 1 2. Что получится? Запишите цифры без пробелов.", Answer: "1234", Explain: "Первый раз: 1 < 2, обмен даёт 4 3 2 1. Второй раз: всё невозрастает, это последний вариант. Разворачиваем весь Array: 1 2 3 4. Во второй раз поиск i делает 3 сравнения. Разворот обменами требует 2 обмена и O(1) рабочей памяти; срезы в образце требуют O(n), где n - длина Array."},
				{ID: "pc4-grid-through", Type: "short", Prompt: "Сетка из 3 строк и 4 столбцов. Робот идёт из левого верхнего угла в правый нижний, только вправо и вниз. Сколько путей проходят через клетку во 2-й строке и 2-м столбце?", Answer: "6", Explain: "До этой клетки 1 шаг вниз и 1 вправо: C(2, 1) = 2 пути. От неё до угла 1 шаг вниз и 2 вправо: C(3, 1) = 3 пути. Любой путь до клетки можно продолжить любым путём после неё: 2 · 3 = 6. Всего путей в сетке C(5, 2) = 10."},
				{ID: "pc4-kth-perm", Type: "short", Prompt: "Permutations цифр 1, 2, 3, 4 стоят в Lexicographic Order и пронумерованы с 1. Какой вариант имеет номер 12? Запишите цифры без пробелов.", Answer: "2431", Explain: "В каждом блоке 3! = 6 вариантов. Номера с 7 по 12 начинаются с 2, и 12-й последний в этом блоке: остальные цифры 1, 3, 4 по убыванию. По формуле: (12 − 1) // 6 = 1, берём 2; остаток 5, 5 // 2 = 2, из 1, 3, 4 берём 4; остаток 1, из 1, 3 берём 3; потом 1. Без вычитания 1 получится 12 // 6 = 2 и ответ 3124, это 13-й вариант."},
				{ID: "pc4-queens-sets", Type: "choice", Prompt: "Доска 4 на 4, строки и столбцы нумеруются с 0. Ферзи стоят в клетках (0, 1) и (1, 3). В какой столбец строки 2 можно поставить ферзя?", Options: []string{"Столбцы 0 и 2", "Только столбец 2", "Только столбец 0", "Ни в один столбец"}, Answer: 2, Explain: "cols = {1, 3}, diag1 = {−1, −2}, diag2 = {1, 4}. Столбец 0: 2 − 0 = 2 и 2 + 0 = 2, обе диагонали свободны. Столбец 2: 2 + 2 = 4, эту диагональ занял ферзь (1, 3). Столбцы 1 и 3 заняты."},
				{ID: "pc4-code-equal", Type: "code", Prompt: "Дополните поиск j: самое правое число, которое больше a[i]. В Array могут быть равные числа, например 1 5 1.", Language: "python", Code: `i = len(a) - 2
while i >= 0 and a[i] >= a[i + 1]:
    i -= 1
if i >= 0:
    j = len(a) - 1
    while ___:
        j -= 1
    a[i], a[j] = a[j], a[i]
a[i + 1:] = a[i + 1:][::-1]`, Options: []string{"a[j] < a[i]", "a[j] >= a[i]", "a[j] > a[i]", "a[j] <= a[i]"}, Answer: 3, Explain: "Цикл пропускает числа, которые не больше a[i]: a[j] <= a[i]. С условием a[j] < a[i] для 1 5 1 цикл остановится на последней 1, она равна a[i] = 1. Обмен ничего не меняет, и выходит 1 1 5 вместо 5 1 1. При верном условии поиск j делает 2 сравнения. Два среза хвоста копируют 4 элемента; для Array длины n рабочая память этого кода O(n)."},
			},
			Links: []sets.Link{
				{Title: "Next Permutation: алгоритм", URL: "https://ru.wikipedia.org/wiki/Алгоритм_Нарайаны", Source: "Википедия"},
				{Title: "Permutations и Factorial Number System", URL: "https://habr.com/ru/articles/319696/", Source: "Хабр"},
				{Title: "Pascal's Triangle", URL: "https://ru.wikipedia.org/wiki/Треугольник_Паскаля", Source: "Википедия"},
				{Title: "N-Queens: задача о восьми ферзях", URL: "https://ru.wikipedia.org/wiki/Задача_о_восьми_ферзях", Source: "Википедия"},
			},
			Figures: []sets.Figure{
				{
					ID:    "pc4-rolling-row",
					Title: "Dynamic Programming: одна строка вместо всей сетки",
					Frames: []sets.FigureFrame{
						{Caption: "Сетка 3 × 3. Начальная строка row = [1, 1, 1]: по одному пути вдоль верхнего края. Храним 3 числа; сложений ещё нет.", Rows: []sets.FigureRow{
							{Label: "row", Cells: []any{1, 1, 1}},
							{Label: "Сложения", Cells: []any{0}},
						}},
						{Caption: "Вторая строка, j = 1: 1 сверху + 1 слева = 2. Обновлённый участок row[0…1] относится к текущей строке; row[2] ещё хранит прошлую. Выполнено 1 сложение, храним 3 числа.", Rows: []sets.FigureRow{
							{Label: "row", Cells: []any{1, 2, 1}, Marked: []int{1}, Pointers: &sets.Pointers{{Name: "j", Cell: 1}}, Window: &[2]int{0, 1}},
							{Label: "Сложения", Cells: []any{1}},
						}},
						{Caption: "Вторая строка, j = 2: 1 сверху + 2 слева = 3. Вся row относится ко второй строке. Выполнено 2 сложения, храним 3 числа.", Rows: []sets.FigureRow{
							{Label: "row", Cells: []any{1, 2, 3}, Marked: []int{2}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}, Window: &[2]int{0, 2}},
							{Label: "Сложения", Cells: []any{2}},
						}},
						{Caption: "Третья строка, j = 1: 2 сверху + 1 слева = 3. Последняя ячейка пока хранит 3 из второй строки. Выполнено 3 сложения, храним 3 числа.", Rows: []sets.FigureRow{
							{Label: "row", Cells: []any{1, 3, 3}, Marked: []int{1}, Pointers: &sets.Pointers{{Name: "j", Cell: 1}}, Window: &[2]int{0, 1}},
							{Label: "Сложения", Cells: []any{3}},
						}},
						{Caption: "Третья строка, j = 2: 3 сверху + 3 слева = 6. Ответ 6. Всего 4 сложения и 3 рабочие ячейки. Полная таблица потребовала бы 9 ячеек при тех же 4 сложениях.", Rows: []sets.FigureRow{
							{Label: "row", Cells: []any{1, 3, 6}, Marked: []int{2}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}, Window: &[2]int{0, 2}},
							{Label: "Сложения", Cells: []any{4}},
						}},
					},
				},
			},
		},
		{
			Title:   "Combinations по модулю",
			Summary: "Pascal's Triangle по модулю, C(n, k) через таблицы Factorial, Stars and Bars, пути на прямой и раскраски. Ответы по модулю 10⁹ + 7 и 998244353, стоимость таблиц и запросов.",
			Lesson:  lesson5,
			Tasks: []sets.Task{
				{
					URL:        "https://codeforces.com/problemset/problem/2025/B",
					Title:      "Binomial Coefficients, Kind Of",
					Difficulty: "easy",
					Topics:     []string{"Math", "Combinatorics", "Dynamic Programming"},
					Note:       "Codeforces, рейтинг 1100. Выпишите маленькую таблицу по неверной формуле и найдите закономерность: C[n][k] = 2ᵏ. Заранее посчитайте p[k] = 2ᵏ % (10⁹ + 7) для k до 10⁵. t - число пар, K - наибольшее k. Время: O(K + t). Память: O(K) для степеней. Таблица по формуле до n = 10⁵ содержит около 5 · 10⁹ чисел.",
					Statement: `<p>Недавно akshiM встретил задачу, для решения которой нужны были числа Combinations C[n][k]. Он написал код, который пишет обычно. Код выглядел так (запись на Python):</p><pre>for n in range(N):
    C[n][0] = 1
    C[n][n] = 1
    for k in range(1, n):
        C[n][k] = C[n][k - 1] + C[n - 1][k - 1]</pre><p>К сожалению, он ошибся, потому что правильная формула такая:</p><pre>C[n][k] = C[n - 1][k] + C[n - 1][k - 1]</pre><p>Но его товарищу по команде keblidA интересны значения, которые получаются по неверной формуле. Помогите ему посчитать эти числа для <code>t</code> разных пар <code>(n<sub>i</sub>, k<sub>i</sub>)</code>. Обратите внимание: считать их нужно по первой (неверной) формуле.</p><p>Значения <code>C[n<sub>i</sub>][k<sub>i</sub>]</code> могут быть очень большими, поэтому выведите их по модулю <code>10<sup>9</sup> + 7</code>.</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>5</sup></code>): количество пар. Дальше <code>t</code> пар записаны в двух строках.</p><p>Во второй строке записаны <code>t</code> целых чисел <code>n<sub>1</sub>, n<sub>2</sub>, ..., n<sub>t</sub></code> (<code>2 ≤ n<sub>i</sub> ≤ 10<sup>5</sup></code>).</p><p>В третьей строке записаны <code>t</code> целых чисел <code>k<sub>1</sub>, k<sub>2</sub>, ..., k<sub>t</sub></code> (<code>1 ≤ k<sub>i</sub> &lt; n<sub>i</sub></code>).</p><h3>Выходные данные</h3><p>Выведите <code>t</code> целых чисел <code>C[n<sub>i</sub>][k<sub>i</sub>]</code> по модулю <code>10<sup>9</sup> + 7</code>.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>7
2 5 5 100000 100000 100000 100000
1 2 3 1 33333 66666 99999</pre><p><strong>Выходные данные:</strong></p><pre>2
4
8
2
326186014
984426998
303861760</pre><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>5</sup></code></li><li><code>2 ≤ n<sub>i</sub> ≤ 10<sup>5</sup></code></li><li><code>1 ≤ k<sub>i</sub> &lt; n<sub>i</sub></code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 512 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `7
2 5 5 100000 100000 100000 100000
1 2 3 1 33333 66666 99999
`, Expected: text(`2
4
8
2
326186014
984426998
303861760
`)},
						{Input: `1
2
1
`, Expected: text(`2
`)},
						{Input: `3
3 4 4
2 1 3
`, Expected: text(`4
2
8
`)},
						{Input: `2
31 31
30 29
`, Expected: text(`73741817
536870912
`)},
						{Input: `3
6 6 6
3 4 5
`, Expected: text(`8
16
32
`)},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/find-the-n-th-value-after-k-seconds/",
					Number:     number(3179),
					Title:      "Find the N-th Value After K Seconds",
					Difficulty: "medium",
					Topics:     []string{"Array", "Math", "Simulation", "Combinatorics", "Prefix Sum"},
					Note:       "После k секунд a[i] = C(i + k, k): новое число равно сумме соседа слева и числа выше, как в Pascal's Triangle. Ответ C(n + k − 1, k) % (10⁹ + 7). Время: O(n + k + log M) для таблиц и O(1) на ответ; Simulation делает (n − 1) · k сложений. Память: O(n + k) для таблиц, Simulation хранит n чисел. Для n = 4, k = 5: C(8, 5) = 56. При n, k ≤ 1000 Simulation с % тоже успевает: 999 000 сложений.",
					Statement: `<p>Даны два целых числа <code>n</code> и <code>k</code>.</p><p>Сначала есть Array <code>a</code> из <code>n</code> целых чисел, в котором <code>a[i] = 1</code> для всех <code>0 ≤ i ≤ n - 1</code>. Каждую секунду все элементы обновляются одновременно: каждый элемент становится равен сумме всех элементов перед ним плюс сам элемент. Например, через одну секунду <code>a[0]</code> остаётся прежним, <code>a[1]</code> становится равен <code>a[0] + a[1]</code>, <code>a[2]</code> становится равен <code>a[0] + a[1] + a[2]</code>, и так далее.</p><p>Верните <strong>значение</strong> <code>a[n - 1]</code> через <code>k</code> секунд.</p><p>Ответ может быть очень большим, поэтому верните его <strong>по модулю</strong> <code>10<sup>9</sup> + 7</code>.</p><p>Метод класса <code>Solution</code>: <code>def valueAfterKSeconds(self, n: int, k: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> n = 4, k = 5
<strong>Выход:</strong> 56
<strong>Пояснение:</strong> состояние Array после каждой секунды:</pre><table><thead><tr><th>Секунда</th><th>Состояние после</th></tr></thead><tbody><tr><td>0</td><td>[1,1,1,1]</td></tr><tr><td>1</td><td>[1,2,3,4]</td></tr><tr><td>2</td><td>[1,3,6,10]</td></tr><tr><td>3</td><td>[1,4,10,20]</td></tr><tr><td>4</td><td>[1,5,15,35]</td></tr><tr><td>5</td><td>[1,6,21,56]</td></tr></tbody></table><h3>Пример 2</h3><pre><strong>Вход:</strong> n = 5, k = 3
<strong>Выход:</strong> 35
<strong>Пояснение:</strong> состояние Array после каждой секунды:</pre><table><thead><tr><th>Секунда</th><th>Состояние после</th></tr></thead><tbody><tr><td>0</td><td>[1,1,1,1,1]</td></tr><tr><td>1</td><td>[1,2,3,4,5]</td></tr><tr><td>2</td><td>[1,3,6,10,15]</td></tr><tr><td>3</td><td>[1,4,10,20,35]</td></tr></tbody></table><h3>Ограничения</h3><ul><li><code>1 ≤ n, k ≤ 1000</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "valueAfterKSeconds"}, Cases: []sets.TestCase{
						{Input: "[4,5]", Expected: text("56")},
						{Input: "[5,3]", Expected: text("35")},
						{Input: "[1,1000]", Expected: text("1")},
						{Input: "[1000,1]", Expected: text("1000")},
						{Input: "[2,2]", Expected: text("3")},
						{Input: "[1000,1000]", Expected: text("36237869")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/number-of-ways-to-reach-a-position-after-exactly-k-steps/",
					Number:     number(2400),
					Title:      "Number of Ways to Reach a Position After Exactly k Steps",
					Difficulty: "medium",
					Topics:     []string{"Math", "Dynamic Programming", "Combinatorics", "Combinations"},
					Note:       "Пусть d = |endPos − startPos|. Нужно r шагов в сторону конца и k − r обратно: r − (k − r) = d, r = (k + d)/2. Если d > k или k − d нечётно, ответ 0. Иначе C(k, r) % (10⁹ + 7). Время: O(k + log M) для таблиц. Память: O(k). Для 1 → 2 за 3 шага: d = 1, r = 2, C(3, 2) = 3. Dynamic Programming по шагам и позициям стоит O(k²) времени и O(k) памяти. Brute Force проверяет 2ᵏ путей.",
					Statement: `<p>Даны два <strong>положительных</strong> целых числа <code>startPos</code> и <code>endPos</code>. Сначала вы стоите в позиции <code>startPos</code> на <strong>бесконечной</strong> числовой прямой. За один шаг можно сдвинуться на одну позицию влево или на одну позицию вправо.</p><p>Дано положительное целое число <code>k</code>. Верните <em>число <strong>разных</strong> способов попасть из позиции </em><code>startPos</code><em> в позицию </em><code>endPos</code><em>, если сделать <strong>ровно</strong> </em><code>k</code><em> шагов</em>. Ответ может быть очень большим, поэтому верните его <strong>по модулю</strong> <code>10<sup>9</sup> + 7</code>.</p><p>Два способа разные, если порядок шагов в них не совпадает полностью.</p><p><strong>Обратите внимание:</strong> на числовой прямой есть и отрицательные числа.</p><p>Метод класса <code>Solution</code>: <code>def numberOfWays(self, startPos: int, endPos: int, k: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> startPos = 1, endPos = 2, k = 3
<strong>Выход:</strong> 3
<strong>Пояснение:</strong> из позиции 1 в позицию 2 ровно за 3 шага можно попасть тремя способами:
- 1 -&gt; 2 -&gt; 3 -&gt; 2.
- 1 -&gt; 2 -&gt; 1 -&gt; 2.
- 1 -&gt; 0 -&gt; 1 -&gt; 2.
Можно доказать, что других способов нет, поэтому ответ 3.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> startPos = 2, endPos = 5, k = 10
<strong>Выход:</strong> 0
<strong>Пояснение:</strong> из позиции 2 нельзя попасть в позицию 5 ровно за 10 шагов.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ startPos, endPos, k ≤ 1000</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "numberOfWays"}, Cases: []sets.TestCase{
						{Input: "[1,2,3]", Expected: text("3")},
						{Input: "[2,5,10]", Expected: text("0")},
						{Input: "[5,5,2]", Expected: text("2")},
						{Input: "[1,1000,999]", Expected: text("1")},
						{Input: "[7,3,6]", Expected: text("6")},
						{Input: "[500,500,1000]", Expected: text("159835829")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/number-of-sets-of-k-non-overlapping-line-segments/",
					Number:     number(1621),
					Title:      "Number of Sets of K Non-Overlapping Line Segments",
					Difficulty: "medium",
					Topics:     []string{"Math", "Dynamic Programming", "Combinatorics", "Combinations"},
					Note:       "Сдвиньте i-й отрезок вправо на i − 1: тогда 2k концов станут разными числами от 0 до n + k − 2. Ответ C(n + k − 1, 2k) % (10⁹ + 7). Время: O(n + k + log M). Память: O(n + k) для таблиц. Для n = 4, k = 2: C(5, 4) = 5. Dynamic Programming по точкам и числу отрезков с Prefix Sum стоит O(nk), без Prefix Sum O(n²k).",
					Statement: `<p>На прямой даны <code>n</code> точек. Точка с номером <code>i</code> (от <code>0</code> до <code>n - 1</code>) стоит в <code>x = i</code>. Найдите число способов нарисовать <strong>ровно</strong> <code>k</code> <strong>неперекрывающихся</strong> отрезков так, чтобы каждый отрезок покрывал две или больше точек. Концы каждого отрезка должны иметь <strong>целые координаты</strong>. Эти <code>k</code> отрезков <strong>не обязаны</strong> покрывать все <code>n</code> точек, и им <strong>можно</strong> иметь общие концы.</p><p>Верните <em>число способов нарисовать </em><code>k</code><em> неперекрывающихся отрезков</em>. Это число может быть огромным, поэтому верните его <strong>по модулю</strong> <code>10<sup>9</sup> + 7</code>.</p><p>Метод класса <code>Solution</code>: <code>def numberOfSets(self, n: int, k: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> n = 4, k = 2
<strong>Выход:</strong> 5
<strong>Пояснение:</strong> есть 5 разных способов: {(0,2),(2,3)}, {(0,1),(1,3)}, {(0,1),(2,3)}, {(1,2),(2,3)}, {(0,1),(1,2)}.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> n = 3, k = 1
<strong>Выход:</strong> 3
<strong>Пояснение:</strong> 3 способа: {(0,1)}, {(0,2)}, {(1,2)}.</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> n = 30, k = 7
<strong>Выход:</strong> 796297179
<strong>Пояснение:</strong> всего есть 3796297200 способов нарисовать 7 отрезков. Остаток от деления этого числа на 10<sup>9</sup> + 7 равен 796297179.</pre><h3>Ограничения</h3><ul><li><code>2 ≤ n ≤ 1000</code></li><li><code>1 ≤ k ≤ n - 1</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "numberOfSets"}, Cases: []sets.TestCase{
						{Input: "[4,2]", Expected: text("5")},
						{Input: "[3,1]", Expected: text("3")},
						{Input: "[30,7]", Expected: text("796297179")},
						{Input: "[2,1]", Expected: text("1")},
						{Input: "[1000,1]", Expected: text("499500")},
						{Input: "[1000,500]", Expected: text("70047606")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1081/C",
					Title:      "Colorful Bricks",
					Difficulty: "medium",
					Topics:     []string{"Math", "Combinatorics", "Dynamic Programming", "Combinations"},
					Note:       "Codeforces, рейтинг 1500. Выберите k мест смены цвета из n − 1: C(n − 1, k). Первый кирпич красят m способами, каждый кирпич смены m − 1 способами. Ответ C(n − 1, k) · m · (m − 1)ᵏ по модулю M = 998244353. Время: O(n + log M). Память: O(n) для таблиц Factorial. Для n = 3, m = 2, k = 1: 2 · 2 · 1 = 4. Dynamic Programming f[i][j] = f[i − 1][j] + f[i − 1][j − 1] · (m − 1) стоит O(nk) времени: при n = 2000 это до 4 · 10⁶ шагов, тоже успевает.",
					Statement:  `<p>В свободное время Чоути любит заниматься домашними делами. У него появилось новое задание: покрасить кирпичи во дворе.</p><p>На земле в ряд лежат <code>n</code> кирпичей. У Чоути под рукой <code>m</code> вёдер с красками разных цветов, и он покрасил каждый кирпич в один из этих <code>m</code> цветов.</p><p>Закончив покраску, Чоути остался доволен. Он отошёл назад и решил найти в кирпичах что-нибудь интересное. Он посчитал и нашёл <code>k</code> кирпичей, цвет которых отличается от цвета кирпича слева от них (первый кирпич, конечно, не считается).</p><p>Как обычно, ему нужна ваша помощь: посчитайте, сколькими способами он мог покрасить кирпичи. Два способа покраски разные, если хотя бы один кирпич в них покрашен в разные цвета. Ответ может быть довольно большим, поэтому выведите только число способов по модулю <code>998 244 353</code>.</p><h3>Входные данные</h3><p>В первой и единственной строке записаны три целых числа <code>n</code>, <code>m</code> и <code>k</code> (<code>1 ≤ n, m ≤ 2000</code>, <code>0 ≤ k ≤ n - 1</code>): число кирпичей, число цветов и число кирпичей, цвет которых отличается от цвета кирпича слева от них.</p><h3>Выходные данные</h3><p>Выведите одно целое число: число способов покрасить кирпичи по модулю <code>998 244 353</code>.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>3 3 0</pre><p><strong>Выходные данные:</strong></p><pre>3</pre><h3>Пример 2</h3><p><strong>Входные данные:</strong></p><pre>3 2 1</pre><p><strong>Выходные данные:</strong></p><pre>4</pre><h3>Примечание</h3><p>В первом примере <code>k = 0</code>, поэтому все кирпичи должны быть одного цвета. Значит, способов ровно <code>m = 3</code>.</p><p>Во втором примере пусть в вёдрах жёлтая и салатовая краска. Вот все <code>4</code> возможные раскраски (Ж - жёлтый кирпич, С - салатовый): ЖЖС, ЖСС, ССЖ, СЖЖ.</p><h3>Ограничения</h3><ul><li><code>1 ≤ n, m ≤ 2000</code></li><li><code>0 ≤ k ≤ n - 1</code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `3 3 0
`, Expected: text(`3
`)},
						{Input: `3 2 1
`, Expected: text(`4
`)},
						{Input: `1 1 0
`, Expected: text(`1
`)},
						{Input: `2000 1 1
`, Expected: text(`0
`)},
						{Input: `4 3 2
`, Expected: text(`36
`)},
						{Input: `2000 2000 1000
`, Expected: text(`161680257
`)},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "pc5-row-mod", Type: "short", Prompt: "Строка n = 5 Pascal's Triangle по модулю 4: 1, 1, 2, 2, 1, 1. Постройте из неё строку n = 6 по модулю 4, % после каждого сложения. Запишите числа строки без пробелов.", Answer: "1230321", Explain: "Края равны 1. Внутри: 1 + 1 = 2, 1 + 2 = 3, 2 + 2 = 4 даёт остаток 0, 2 + 1 = 3, 1 + 1 = 2. Проверка без модуля: строка 6 равна 1, 6, 15, 20, 15, 6, 1, и 6 % 4 = 2, 15 % 4 = 3, 20 % 4 = 0. Без % после сложения в середине останется 4, а остаток по модулю 4 всегда меньше 4. Строка требует 5 сложений, как и без модуля."},
				{ID: "pc5-pascal-cost", Type: "choice", Prompt: "Сколько примерно сложений нужно, чтобы построить Pascal's Triangle по модулю для всех строк от 0 до N = 3 · 10⁴?", Options: []string{"около 3 · 10⁴", "около 4,5 · 10⁸", "около 9 · 10⁸"}, Answer: 1, Explain: "Строка n требует n − 1 сложений, всего N(N − 1)/2 = 30 000 · 29 999 / 2 = 449 985 000, около 4,5 · 10⁸. 9 · 10⁸ - это N², но строка n короче N. 3 · 10⁴ - это число строк, а не сложений. По оценке 10⁷ шагов в секунду это около 45 с при 3 с на тест. Таблицы Factorial дают O(N + log M) подготовки и O(1) на каждое C(n, k)."},
				{ID: "pc5-stars", Type: "short", Prompt: "Сколько способов разложить 4 одинаковых шара по 3 разным коробкам, если коробки могут быть пустыми?", Answer: "15", Explain: "4 звезды и 2 палочки занимают 6 мест; выбираем места палочек: C(6, 2) = 15. Это формула C(n + k − 1, k − 1) при n = 4, k = 3. Ответ C(3, 2) = 3 считает только раскладки без пустых коробок. Brute Force по тройкам (a, b, c) с a + b + c = 4 даёт те же 15 вариантов, но число троек растёт с n и k."},
				{ID: "pc5-stars-min", Type: "choice", Prompt: "9 одинаковых конфет делят между 3 детьми, каждому хотя бы одну. Какая формула даёт число способов?", Options: []string{"C(11, 2)", "C(9, 3)", "C(8, 2)"}, Answer: 2, Explain: "Сначала дайте каждому по конфете. Остаются 6 конфет без ограничений: C(6 + 3 − 1, 2) = C(8, 2) = 28. C(11, 2) = 55 считает и пустые доли для 9 конфет. C(9, 3) = 84 выбирает 3 конфеты из 9, а не 2 границы. То же число 28 даёт выбор 2 из 8 промежутков между 9 конфетами."},
				{ID: "pc5-line", Type: "short", Prompt: "Путь на прямой из 0 в 3 ровно за 5 шагов, каждый шаг на 1 влево или вправо. Сколько путей?", Answer: "5", Explain: "Шагов вправо r, влево 5 − r, и r − (5 − r) = 3, значит r = 4. Выбираем место одного шага влево: C(5, 1) = 5. Ответ C(5, 3) = 10 неверен: 3 - это смещение, а не число шагов вправо. При чётном числе шагов в 3 не попасть, тогда ответ 0."},
				{ID: "pc5-parity", Type: "choice", Prompt: "startPos = 1, endPos = 4, k = 6. Что вернёт LeetCode 2400?", Options: []string{"0", "C(6, 5) = 6", "C(6, 3) = 20"}, Answer: 0, Explain: "Каждый шаг меняет чётность позиции. За 6 шагов чётность возвращается, а 1 и 4 разной чётности: ответ 0. Формула r = (k + d)/2 = (6 + 3)/2 = 4,5 не целая. C(6, 5) получится, если округлить 4,5 до 5, но 5 шагов вправо и 1 влево ведут из 1 в 5, а не в 4. C(6, 3) считает пути длины 6 в ту же позицию."},
				{ID: "pc5-code", Type: "code", Prompt: "Функция C(n, k) по таблицам. Что стоит на месте ___?", Language: "python", Code: `def C(n, k):
    if k < 0 or k > n:
        return 0
    return fact[n] * inv_fact[k] % M * ___ % M`, Options: []string{"fact[n - k]", "inv_fact[n - k]", "inv_fact[n]"}, Answer: 1, Explain: "C(n, k) = n! / (k! · (n − k)!), деление на (n − k)! - это умножение на inv_fact[n − k]. С fact[n − k] получится n! · (n − k)! / k!, неверно. С inv_fact[n] деление на n! отменяет n!, и ответ равен 1 / k!. После подготовки таблиц вызов делает 2 умножения и 2 операции %: O(1)."},
				{ID: "pc5-flags", Type: "short", Prompt: "Ряд из 4 флажков, 2 цвета. Ровно 2 раза соседние флажки разного цвета. Сколько раскрасок?", Answer: "6", Explain: "Места смены: 2 из 3 промежутков, C(3, 2) = 3. Первый флажок: 2 цвета, после каждой смены 2 − 1 = 1 цвет. Ответ 3 · 2 · 1 = 6: ABAA, ABBA, AABA и три раскраски с B в начале. Ответ 3 забывает цвет первого флажка. Это C(n − 1, k) · m · (m − 1) в степени k."},
			},
			Links: []sets.Link{
				{Title: "Биномиальный коэффициент", URL: "https://ru.wikipedia.org/wiki/Биномиальный_коэффициент", Source: "Википедия"},
				{Title: "Комбинаторные объекты", URL: "https://neerc.ifmo.ru/wiki/index.php?title=Комбинаторные_объекты", Source: "Викиконспекты ИТМО"},
				{Title: "Деление по модулю и таблицы Factorial", URL: "https://redlamp.me/leetcode/sets/modular-arithmetic/3", Source: "Лекции"},
				{Title: "Stars and Bars", URL: "https://cp-algorithms.com/combinatorics/stars_and_bars.html", Source: "cp-algorithms (англ.)"},
				{Title: "Binomial Coefficients по модулю", URL: "https://cp-algorithms.com/combinatorics/binomial-coefficients.html", Source: "cp-algorithms (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "pc5-pascal",
					Title: "Pascal's Triangle по модулю 7",
					Frames: []sets.FigureFrame{
						{Caption: "Строки от 0 до 4: каждое внутреннее число равно сумме двух чисел над ним. Все числа меньше 7.", Rows: []sets.FigureRow{
							{Label: "n = 0", Cells: []any{1}},
							{Label: "n = 1", Cells: []any{1, 1}},
							{Label: "n = 2", Cells: []any{1, 2, 1}},
							{Label: "n = 3", Cells: []any{1, 3, 3, 1}},
							{Label: "n = 4", Cells: []any{1, 4, 6, 4, 1}},
						}},
						{Caption: "Строки 5 и 6 по модулю 7: 10 даёт остаток 3, 15 даёт остаток 1, 20 даёт остаток 6. Остатки складывают тоже по модулю: 5 + 3 = 8 даёт остаток 1.", Rows: []sets.FigureRow{
							{Label: "n = 0", Cells: []any{1}},
							{Label: "n = 1", Cells: []any{1, 1}},
							{Label: "n = 2", Cells: []any{1, 2, 1}},
							{Label: "n = 3", Cells: []any{1, 3, 3, 1}},
							{Label: "n = 4", Cells: []any{1, 4, 6, 4, 1}},
							{Label: "n = 5", Cells: []any{1, 5, 3, 3, 5, 1}, Marked: []int{2, 3}},
							{Label: "n = 6", Cells: []any{1, 6, 1, 6, 1, 6, 1}, Marked: []int{2, 3, 4}},
						}},
						{Caption: "Строка 7: 6 + 1 = 7 и 1 + 6 = 7 дают остаток 0. Все внутренние C(7, k) делятся на 7. Строки до 7 потребовали 21 сложение.", Rows: []sets.FigureRow{
							{Label: "n = 0", Cells: []any{1}},
							{Label: "n = 1", Cells: []any{1, 1}},
							{Label: "n = 2", Cells: []any{1, 2, 1}},
							{Label: "n = 3", Cells: []any{1, 3, 3, 1}},
							{Label: "n = 4", Cells: []any{1, 4, 6, 4, 1}},
							{Label: "n = 5", Cells: []any{1, 5, 3, 3, 5, 1}},
							{Label: "n = 6", Cells: []any{1, 6, 1, 6, 1, 6, 1}},
							{Label: "n = 7", Cells: []any{1, 0, 0, 0, 0, 0, 0, 1}, Marked: []int{1, 2, 3, 4, 5, 6}},
						}},
					},
				},
				{
					ID:    "pc5-stars",
					Title: "Stars and Bars: 5 конфет, 3 ребёнка",
					Frames: []sets.FigureFrame{
						{Caption: "★★|★|★★: палочки на местах 3 и 5. Доли 2, 1, 2.", Rows: []sets.FigureRow{
							{Label: "место", Cells: []any{1, 2, 3, 4, 5, 6, 7}},
							{Label: "запись", Cells: []any{"★", "★", "|", "★", "|", "★", "★"}, Marked: []int{2, 4}},
							{Label: "доли", Cells: []any{2, 1, 2}},
						}},
						{Caption: "|★★★★★|: палочки на местах 1 и 7. Доли 0, 5, 0: пустая доля разрешена.", Rows: []sets.FigureRow{
							{Label: "место", Cells: []any{1, 2, 3, 4, 5, 6, 7}},
							{Label: "запись", Cells: []any{"|", "★", "★", "★", "★", "★", "|"}, Marked: []int{0, 6}},
							{Label: "доли", Cells: []any{0, 5, 0}},
						}},
						{Caption: "★★★★★||: палочки на местах 6 и 7. Каждый выбор 2 мест из 7 даёт одну раздачу: C(7, 2) = 21.", Rows: []sets.FigureRow{
							{Label: "место", Cells: []any{1, 2, 3, 4, 5, 6, 7}},
							{Label: "запись", Cells: []any{"★", "★", "★", "★", "★", "|", "|"}, Marked: []int{5, 6}},
							{Label: "доли", Cells: []any{5, 0, 0}},
						}},
					},
				},
				{
					ID:    "pc5-line",
					Title: "Пути на прямой: из 0 за 4 шага",
					Frames: []sets.FigureFrame{
						{Caption: "Шаг 0: один путь стоит в 0.", Rows: []sets.FigureRow{
							{Label: "позиция", Cells: []any{-4, -3, -2, -1, 0, 1, 2, 3, 4}},
							{Label: "шаг 0", Cells: []any{"", "", "", "", 1, "", "", "", ""}, Marked: []int{4}},
						}},
						{Caption: "Шаг 1: в −1 и в 1 по одному пути.", Rows: []sets.FigureRow{
							{Label: "позиция", Cells: []any{-4, -3, -2, -1, 0, 1, 2, 3, 4}},
							{Label: "шаг 0", Cells: []any{"", "", "", "", 1, "", "", "", ""}},
							{Label: "шаг 1", Cells: []any{"", "", "", 1, "", 1, "", "", ""}, Marked: []int{3, 5}},
						}},
						{Caption: "Шаг 2: в 0 ведут 2 пути. Число в клетке равно сумме двух соседей шагом раньше.", Rows: []sets.FigureRow{
							{Label: "позиция", Cells: []any{-4, -3, -2, -1, 0, 1, 2, 3, 4}},
							{Label: "шаг 0", Cells: []any{"", "", "", "", 1, "", "", "", ""}},
							{Label: "шаг 1", Cells: []any{"", "", "", 1, "", 1, "", "", ""}},
							{Label: "шаг 2", Cells: []any{"", "", 1, "", 2, "", 1, "", ""}, Marked: []int{2, 4, 6}},
						}},
						{Caption: "Шаг 3: 1, 3, 3, 1, как строка 3 Pascal's Triangle.", Rows: []sets.FigureRow{
							{Label: "позиция", Cells: []any{-4, -3, -2, -1, 0, 1, 2, 3, 4}},
							{Label: "шаг 0", Cells: []any{"", "", "", "", 1, "", "", "", ""}},
							{Label: "шаг 1", Cells: []any{"", "", "", 1, "", 1, "", "", ""}},
							{Label: "шаг 2", Cells: []any{"", "", 1, "", 2, "", 1, "", ""}},
							{Label: "шаг 3", Cells: []any{"", 1, "", 3, "", 3, "", 1, ""}, Marked: []int{1, 3, 5, 7}},
						}},
						{Caption: "Шаг 4: 1, 4, 6, 4, 1. В позицию 2 ведут 4 пути: C(4, 1) = 4. В нечётные позиции за 4 шага не попасть.", Rows: []sets.FigureRow{
							{Label: "позиция", Cells: []any{-4, -3, -2, -1, 0, 1, 2, 3, 4}},
							{Label: "шаг 0", Cells: []any{"", "", "", "", 1, "", "", "", ""}},
							{Label: "шаг 1", Cells: []any{"", "", "", 1, "", 1, "", "", ""}},
							{Label: "шаг 2", Cells: []any{"", "", 1, "", 2, "", 1, "", ""}},
							{Label: "шаг 3", Cells: []any{"", 1, "", 3, "", 3, "", 1, ""}},
							{Label: "шаг 4", Cells: []any{1, "", 4, "", 6, "", 4, "", 1}, Marked: []int{0, 2, 4, 6, 8}, Pointers: &sets.Pointers{{Name: "конец", Cell: 6}}},
						}},
					},
				},
			},
		},
	},
}
