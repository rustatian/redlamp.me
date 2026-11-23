// Package prefixsums is the learning set "Prefix Sum и Difference Array".
package prefixsums

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
	Slug:    "prefix-sums",
	Title:   "Prefix Sum и Difference Array",
	Summary: "Prefix Sum заранее складывает начало Array, и сумма любого отрезка стоит одно вычитание. Вы отвечаете на запросы сумм, ищете Subarray с нужной суммой через Hash Table, прибавляете число сразу ко всему отрезку через Difference Array и переносите эти приёмы на Matrix. Для каждого решения считаете операции, дополнительную память и размер ответа.",
	Topics:  []string{"Prefix Sum", "Array", "Hash Table", "Matrix"},
	Chapters: []sets.Chapter{
		{
			Title:   "Prefix Sum на Array",
			Summary: "Prefix Sum хранит суммы начала Array, поэтому сумма отрезка стоит одно вычитание. Вы строите префиксы с нулём в начале, отвечаете на запросы с номерами с 0 и с 1, ищете опорный индекс по суммам слева и справа и заменяете суммы произведениями.",
			Lesson:  lesson1,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/find-pivot-index/",
					Number:     number(724),
					Title:      "Find Pivot Index",
					Difficulty: "easy",
					Topics:     []string{"Array", "Prefix Sum"},
					Note:       "Linear Scan с total и left: сумма справа равна total − left − nums[i]. n - длина nums. Время: O(n), один sum и один проход. Память: O(1) дополнительно, ответ O(1). Для [1, 7, 3, 6, 5, 6]: total = 28, на i = 3 слева 11 и справа 11. Срезы sum(nums[:i]) и sum(nums[i + 1:]) для каждого i дают O(n²) времени и O(n) временной памяти.",
					Statement: `<p>Дан Array целых чисел <code>nums</code>. Найдите <strong>опорный индекс</strong> этого Array.</p><p><strong>Опорный индекс</strong> это индекс, для которого сумма всех чисел <strong>строго</strong> слева от него равна сумме всех чисел <strong>строго</strong> справа от него.</p><p>Если индекс стоит на левом краю Array, сумма слева равна <code>0</code>, потому что слева нет элементов. Это верно и для правого края Array.</p><p>Верните <em><strong>самый левый опорный индекс</strong></em>. Если такого индекса нет, верните <code>-1</code>.</p><p>Метод класса <code>Solution</code>: <code>def pivotIndex(self, nums: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,7,3,6,5,6]
<strong>Выход:</strong> 3
<strong>Пояснение:</strong>
Опорный индекс 3.
Сумма слева = nums[0] + nums[1] + nums[2] = 1 + 7 + 3 = 11
Сумма справа = nums[4] + nums[5] = 5 + 6 = 11</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [1,2,3]
<strong>Выход:</strong> -1
<strong>Пояснение:</strong>
Ни один индекс не выполняет условие задачи.</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> nums = [2,1,-1]
<strong>Выход:</strong> 0
<strong>Пояснение:</strong>
Опорный индекс 0.
Сумма слева = 0 (слева от индекса 0 нет элементов)
Сумма справа = nums[1] + nums[2] = 1 + (−1) = 0</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 10<sup>4</sup></code></li><li><code>-1000 ≤ nums[i] ≤ 1000</code></li></ul><p><strong>Примечание:</strong> эта задача совпадает с LeetCode 1991.</p>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "pivotIndex"}, Cases: []sets.TestCase{
						{Input: "[[1,7,3,6,5,6]]", Expected: text("3")},
						{Input: "[[1,2,3]]", Expected: text("-1")},
						{Input: "[[2,1,-1]]", Expected: text("0")},
						{Input: "[[5]]", Expected: text("0")},
						{Input: "[[0,0,0]]", Expected: text("0")},
						{Input: "[[1,-1,3]]", Expected: text("2")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1807/D",
					Title:      "Odd Queries",
					Difficulty: "easy",
					Topics:     []string{"Array", "Prefix Sum", "Math"},
					Note:       "Codeforces, рейтинг 900. Prefix Sum с номерами с 1: сумма a[l..r] = p[r] − p[l − 1]. Новая сумма: total − (p[r] − p[l − 1]) + k · (r − l + 1). n и q - размеры одного набора. Время: O(n + q) на набор, по всем наборам O(∑n + ∑q). Память: O(n) для p, ответ q строк. Повторный подсчёт суммы отрезка для каждого запроса стоит O(n · q) в худшем случае.",
					Statement: `<p>Дан Array <code>a<sub>1</sub>, a<sub>2</sub>, ..., a<sub>n</sub></code>. Ответьте на <code>q</code> запросов такого вида:</p><ul><li>если заменить все элементы отрезка <code>a<sub>l</sub>, a<sub>l+1</sub>, ..., a<sub>r</sub></code> на <code>k</code>, будет ли сумма всего Array нечётной?</li></ul><h3>Входные данные</h3><p>Во входе несколько наборов входных данных. В первой строке записано число наборов <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>). Дальше идут описания наборов.</p><p>В первой строке каждого набора записаны <code>2</code> целых числа <code>n</code> и <code>q</code> (<code>1 ≤ n ≤ 2 · 10<sup>5</sup></code>; <code>1 ≤ q ≤ 2 · 10<sup>5</sup></code>): длина Array и число запросов.</p><p>Во второй строке набора записаны <code>n</code> целых чисел <code>a<sub>i</sub></code> (<code>1 ≤ a<sub>i</sub> ≤ 10<sup>9</sup></code>): Array <code>a</code>.</p><p>В следующих <code>q</code> строках набора записано по <code>3</code> целых числа <code>l, r, k</code> (<code>1 ≤ l ≤ r ≤ n</code>; <code>1 ≤ k ≤ 10<sup>9</sup></code>): запросы.</p><p>Гарантируется, что сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code> и сумма <code>q</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого запроса выведите <code>YES</code>, если сумма всего Array станет нечётной, и <code>NO</code> в другом случае.</p><p>Ответ можно выводить в любом регистре. Например, строки <code>yEs</code>, <code>yes</code>, <code>Yes</code> и <code>YES</code> считаются положительным ответом.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>2
5 5
2 2 1 3 2
2 3 3
2 3 4
1 5 5
1 4 9
2 4 3
10 5
1 1 1 1 1 1 1 1 1 1
3 8 13
2 5 10
3 8 10
1 10 2
1 9 100</pre><p><strong>Выходные данные:</strong></p><pre>YES
YES
YES
NO
YES
NO
NO
NO
NO
YES</pre><h3>Примечание</h3><p>Первый набор:</p><ul><li>Если элементы отрезка <code>(2, 3)</code> заменить на <code>3</code>, Array станет <code>{2, 3, 3, 3, 2}</code>. Сумма <code>2 + 3 + 3 + 3 + 2 = 13</code> нечётная, поэтому ответ <code>YES</code>.</li><li>Если элементы отрезка <code>(2, 3)</code> заменить на <code>4</code>, Array станет <code>{2, 4, 4, 3, 2}</code>. Сумма <code>2 + 4 + 4 + 3 + 2 = 15</code> нечётная, поэтому ответ <code>YES</code>.</li><li>Если элементы отрезка <code>(1, 5)</code> заменить на <code>5</code>, Array станет <code>{5, 5, 5, 5, 5}</code>. Сумма <code>5 + 5 + 5 + 5 + 5 = 25</code> нечётная, поэтому ответ <code>YES</code>.</li><li>Если элементы отрезка <code>(1, 4)</code> заменить на <code>9</code>, Array станет <code>{9, 9, 9, 9, 2}</code>. Сумма <code>9 + 9 + 9 + 9 + 2 = 38</code> чётная, поэтому ответ <code>NO</code>.</li><li>Если элементы отрезка <code>(2, 4)</code> заменить на <code>3</code>, Array станет <code>{2, 3, 3, 3, 2}</code>. Сумма <code>2 + 3 + 3 + 3 + 2 = 13</code> нечётная, поэтому ответ <code>YES</code>.</li></ul><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>1 ≤ n ≤ 2 · 10<sup>5</sup></code>, <code>1 ≤ q ≤ 2 · 10<sup>5</sup></code></li><li><code>1 ≤ a<sub>i</sub> ≤ 10<sup>9</sup></code></li><li><code>1 ≤ l ≤ r ≤ n</code>, <code>1 ≤ k ≤ 10<sup>9</sup></code></li><li>сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>, сумма <code>q</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio", IgnoreCase: true}, Cases: []sets.TestCase{
						{Input: `2
5 5
2 2 1 3 2
2 3 3
2 3 4
1 5 5
1 4 9
2 4 3
10 5
1 1 1 1 1 1 1 1 1 1
3 8 13
2 5 10
3 8 10
1 10 2
1 9 100
`, Expected: text(`YES
YES
YES
NO
YES
NO
NO
NO
NO
YES
`)},
						{Input: `1
1 2
1
1 1 2
1 1 1
`, Expected: text(`NO
YES
`)},
						{Input: `1
3 2
1000000000 1000000000 1000000000
1 3 999999999
2 2 1000000000
`, Expected: text(`YES
NO
`)},
						{Input: `2
2 1
1 2
1 1 1
2 2
1 1
1 1 1
1 2 1
`, Expected: text(`YES
NO
NO
`)},
						{Input: `1
4 2
1 1 1 1
2 3 2
1 1 2
`, Expected: text(`NO
YES
`)},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/product-of-array-except-self/",
					Number:     number(238),
					Title:      "Product of Array Except Self",
					Difficulty: "medium",
					Topics:     []string{"Array", "Prefix Sum"},
					Note:       "Prefix Product слева и справа в одном списке ans. n - длина nums. Время: O(n), два прохода. Память: O(1) дополнительно, left и right; ответ O(n). Для [1, 2, 3, 4]: после первого прохода ans = [1, 1, 2, 6], после второго [24, 12, 8, 6]. Отдельные списки префиксов и суффиксов дают O(n) дополнительной памяти. Деление на nums[i] ломается на нуле.",
					Statement: `<p>Дан Array целых чисел <code>nums</code>. Верните <em>Array</em> <code>answer</code>, <em>в котором</em> <code>answer[i]</code> <em>равно произведению всех элементов</em> <code>nums</code> <em>кроме</em> <code>nums[i]</code>.</p><p>Произведение любого префикса или суффикса <code>nums</code> <strong>гарантированно</strong> помещается в <strong>32-битное</strong> целое число.</p><p>Напишите алгоритм, который работает за <code>O(n)</code> и не использует деление.</p><p>Метод класса <code>Solution</code>: <code>def productExceptSelf(self, nums: list[int]) -&gt; list[int]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,2,3,4]
<strong>Выход:</strong> [24,12,8,6]</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [-1,1,0,-3,3]
<strong>Выход:</strong> [0,0,9,0,0]</pre><h3>Ограничения</h3><ul><li><code>2 ≤ nums.length ≤ 10<sup>5</sup></code></li><li><code>-30 ≤ nums[i] ≤ 30</code></li><li>Вход составлен так, что <code>answer[i]</code> <strong>гарантированно</strong> помещается в <strong>32-битное</strong> целое число.</li></ul><p><strong>Дополнительный вопрос:</strong> можете ли вы решить задачу с <code>O(1)</code> дополнительной памяти? (Array ответа <strong>не считается</strong> дополнительной памятью в этой оценке.)</p>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "productExceptSelf"}, Cases: []sets.TestCase{
						{Input: "[[1,2,3,4]]", Expected: text("[24,12,8,6]")},
						{Input: "[[-1,1,0,-3,3]]", Expected: text("[0,0,9,0,0]")},
						{Input: "[[0,0]]", Expected: text("[0,0]")},
						{Input: "[[2,3]]", Expected: text("[3,2]")},
						{Input: "[[0,4,5]]", Expected: text("[20,0,0]")},
						{Input: "[[-2,-3,30,-1]]", Expected: text("[90,60,-6,180]")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/433/B",
					Title:      "Kuriyama Mirai's Stones",
					Difficulty: "medium",
					Topics:     []string{"Array", "Prefix Sum", "Sorting"},
					Note:       "Codeforces, рейтинг 1200. Две Prefix Sum: по v и по sorted(v). n - число камней, m - число вопросов. Время: O(n log n + m), Sorting один раз и одно вычитание на вопрос. Память: O(n) для упорядоченной копии и двух списков префиксов; ответ m чисел. Сумма до 10^14 помещается в int Python. Пересчёт суммы для каждого вопроса стоит O(n · m).",
					Statement: `<p>Курияма Мирай победила много монстров и получила много камней, а именно <code>n</code>. Она нумерует камни от <code>1</code> до <code>n</code>. Стоимость <code>i</code>-го камня равна <code>v<sub>i</sub></code>. Курияма Мирай хочет кое-что узнать об этих камнях, поэтому она задаст вам вопросы двух видов:</p><ol><li>Она назовёт два числа <code>l</code> и <code>r</code> (<code>1 ≤ l ≤ r ≤ n</code>), и вы должны назвать ей сумму <code>v<sub>l</sub> + v<sub>l+1</sub> + ... + v<sub>r</sub></code>.</li><li>Пусть <code>u<sub>i</sub></code> это стоимость <code>i</code>-го по дешевизне камня (стоимость, которая окажется на <code>i</code>-м месте, если расставить все стоимости камней в неубывающем порядке). Теперь она назовёт два числа <code>l</code> и <code>r</code> (<code>1 ≤ l ≤ r ≤ n</code>), и вы должны назвать ей сумму <code>u<sub>l</sub> + u<sub>l+1</sub> + ... + u<sub>r</sub></code>.</li></ol><p>На каждый вопрос нужно дать верный ответ, иначе Курияма Мирай скажет «fuyukai desu» и расстроится.</p><h3>Входные данные</h3><p>В первой строке записано целое число <code>n</code> (<code>1 ≤ n ≤ 10<sup>5</sup></code>). Во второй строке записаны <code>n</code> целых чисел <code>v<sub>1</sub>, v<sub>2</sub>, ..., v<sub>n</sub></code> (<code>1 ≤ v<sub>i</sub> ≤ 10<sup>9</sup></code>): стоимости камней.</p><p>В третьей строке записано целое число <code>m</code> (<code>1 ≤ m ≤ 10<sup>5</sup></code>): число вопросов Курияма Мирай. Дальше идут <code>m</code> строк, в каждой три целых числа <code>type</code>, <code>l</code> и <code>r</code> (<code>1 ≤ l ≤ r ≤ n</code>; <code>1 ≤ type ≤ 2</code>): описание вопроса. Если <code>type</code> равен 1, выведите ответ на вопрос первого вида, иначе ответ на вопрос второго вида.</p><h3>Выходные данные</h3><p>Выведите <code>m</code> строк. В каждой строке должно быть одно целое число: ответ на вопрос Курияма Мирай. Выводите ответы в порядке вопросов во входе.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>6
6 4 2 7 2 7
3
2 3 6
1 3 4
1 1 6</pre><p><strong>Выходные данные:</strong></p><pre>24
9
28</pre><h3>Пример 2</h3><p><strong>Входные данные:</strong></p><pre>4
5 5 2 3
10
1 2 4
2 1 4
1 1 1
2 1 4
2 1 2
1 1 1
1 3 3
1 1 3
1 4 4
1 2 2</pre><p><strong>Выходные данные:</strong></p><pre>10
15
5
15
5
5
2
12
3
5</pre><h3>Примечание</h3><p>Обратите внимание: ответы на вопросы могут выйти за пределы 32-битного целого типа.</p><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 10<sup>5</sup></code></li><li><code>1 ≤ v<sub>i</sub> ≤ 10<sup>9</sup></code></li><li><code>1 ≤ m ≤ 10<sup>5</sup></code></li><li><code>1 ≤ l ≤ r ≤ n</code>, <code>1 ≤ type ≤ 2</code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `6
6 4 2 7 2 7
3
2 3 6
1 3 4
1 1 6
`, Expected: text(`24
9
28
`)},
						{Input: `4
5 5 2 3
10
1 2 4
2 1 4
1 1 1
2 1 4
2 1 2
1 1 1
1 3 3
1 1 3
1 4 4
1 2 2
`, Expected: text(`10
15
5
15
5
5
2
12
3
5
`)},
						{Input: `1
7
2
1 1 1
2 1 1
`, Expected: text(`7
7
`)},
						{Input: `3
1000000000 1 1000000000
4
1 1 3
2 1 2
2 3 3
1 2 2
`, Expected: text(`2000000001
1000000001
1000000000
1
`)},
						{Input: `3
3 1 2
4
2 1 1
1 1 1
1 1 2
2 2 3
`, Expected: text(`1
3
4
5
`)},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "ps1-formula", Type: "choice", Prompt: "a = [3, 1, 4, 1, 5, 9, 2, 6], p = [0, 3, 4, 8, 9, 14, 23, 25, 31]. Какое выражение равно a[1] + a[2] + a[3]? Номера с 0.", Options: []string{"p[3] - p[1]", "p[4] - p[1]", "p[4] - p[0]"}, Answer: 1, Explain: "Сумма a[l..r] равна p[r + 1] − p[l]: при l = 1 и r = 3 это p[4] − p[1] = 9 − 3 = 6, и 1 + 4 + 1 = 6. p[3] − p[1] = 5 теряет a[3]. p[4] − p[0] = 9 добавляет лишнее a[0] = 3."},
				{ID: "ps1-ops", Type: "short", Prompt: "Array из 10 чисел и 5 запросов суммы всего Array. Brute Force делает 10 сложений на запрос. Prefix Sum делает 10 сложений на построение и 1 вычитание на запрос. На сколько операций меньше делает Prefix Sum?", Answer: "35", Explain: "Brute Force: 5 · 10 = 50 сложений. Prefix Sum: 10 + 5 = 15 операций. Разность 50 − 15 = 35. В общем случае это O(n · q) против O(n + q). При одном запросе выигрыша нет: 10 операций против 11."},
				{ID: "ps1-zero", Type: "choice", Prompt: "Зачем список p начинается с p[0] = 0?", Options: []string{"Чтобы p занимал меньше памяти", "Чтобы Prefix Sum работала с отрицательными числами", "Чтобы запрос с l = 0 считался той же формулой p[r + 1] − p[l]"}, Answer: 2, Explain: "При l = 0 формула даёт p[r + 1] − p[0] = p[r + 1], сумму всего начала. Без нуля нужен отдельный if для l = 0. Памяти p занимает на одно число больше: n + 1 чисел. Отрицательные числа работают и без нуля: вычитание их не портит."},
				{ID: "ps1-pivot", Type: "short", Prompt: "LeetCode 724. nums = [2, 3, −1, 8, 4]. Какой индекс вернёт решение?", Answer: "3", Explain: "total = 16. i = 0: слева 0, справа 14. i = 1: слева 2, справа 11. i = 2: слева 5, справа 12. i = 3: слева 4, справа 16 − 4 − 8 = 4, суммы равны. Если забыть вычесть nums[i], на i = 3 справа получится 12, и ответ потеряется."},
				{ID: "ps1-product-code", Type: "code", Prompt: "LeetCode 238, второй проход. nums = [1, 2, 3, 4], после первого прохода ans = [1, 1, 2, 6]. Что стоит на месте ___?", Language: "python", Code: `right = 1
for i in range(n - 1, -1, -1):
    ans[i] *= right
    ___`, Options: []string{"right = nums[i]", "ans[i] *= nums[i]", "right *= nums[i]"}, Answer: 2, Explain: "right хранит произведение всех чисел правее i: 1, 4, 12, 24. Поэтому после шага нужно right *= nums[i]. С right = nums[i] хранится только одно соседнее число, и ans станет [2, 3, 8, 6]. ans[i] *= nums[i] умножает ответ на само nums[i], а его нужно исключить. Время O(n), дополнительная память O(1)."},
				{ID: "ps1-division", Type: "choice", Prompt: "nums = [−1, 1, 0, −3, 3]. Код считает произведение всех чисел total и для каждого i пишет total // nums[i]. Что произойдёт?", Options: []string{"Ответ [0, 0, 9, 0, 0]", "ZeroDivisionError при i = 2", "Ответ [0, 0, 0, 0, 0]"}, Answer: 1, Explain: "total = 0. Для i = 0 и i = 1 код пишет 0, а при i = 2 делит на nums[2] = 0 и падает с ZeroDivisionError. Верный ответ [0, 0, 9, 0, 0]: без nums[2] произведение (−1) · 1 · (−3) · 3 = 9. Prefix Product слева и справа не делит, поэтому ноль ему не мешает."},
				{ID: "ps1-odd", Type: "short", Prompt: "Codeforces 1807D. a = [3, 1, 4, 1, 5], номера с 1. Все a[2..4] заменили на k = 6. Чему равна новая сумма?", Answer: "26", Explain: "total = 14, p = [0, 3, 4, 8, 9, 14]. Сумма отрезка p[4] − p[1] = 9 − 3 = 6. Новая сумма 14 − 6 + 6 · 3 = 26, она чётная: ответ NO. Ошибка p[4] − p[2] = 5 даёт 27 и ответ YES, потому что теряет a[2]."},
				{ID: "ps1-sorted", Type: "choice", Prompt: "Codeforces 433B. v = [4, 1, 3, 2], номера с 1. Вопрос второго типа: l = 2, r = 3. Какой ответ?", Options: []string{"5", "4", "7"}, Answer: 0, Explain: "Второй тип спрашивает про упорядоченный список u = [1, 2, 3, 4]: u[2] + u[3] = 2 + 3 = 5. Ответ 4 это v[2] + v[3] = 1 + 3, вопрос первого типа. Ответ 7 = 3 + 4 берёт номера с 0. Решение хранит две Prefix Sum, по v и по u: Sorting стоит O(n log n) один раз, вопрос стоит O(1)."},
			},
			Links: []sets.Link{
				{Title: "Prefix Sum: сумма на отрезке", URL: "https://ru.algorithmica.org/cs/range-queries/prefix-sum/", Source: "Алгоритмика"},
				{Title: "Prefix Sum и запросы суммы", URL: "https://brestprog.by/topics/prefixsums/", Source: "brestprog"},
				{Title: "Prefix Sum и суффиксные суммы на Python", URL: "https://notes.algoprog.ru/shortideas/03_x_prefix_sums.html", Source: "Заметки П. Калинина (algoprog)"},
				{Title: "itertools.accumulate", URL: "https://docs.python.org/3/library/itertools.html#itertools.accumulate", Source: "Python docs (англ.)"},
				{Title: "Prefix Sum: Static Range Sum", URL: "https://usaco.guide/silver/prefix-sums", Source: "USACO Guide (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "ps1-build",
					Title: "Prefix Sum: a = [3, 1, 4, 1, 5, 9, 2, 6]",
					Frames: []sets.FigureFrame{
						{Caption: "p[0] = 0: сумма пустого начала. Сложений пока 0. В p будет n + 1 = 9 чисел.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{3, 1, 4, 1, 5, 9, 2, 6}},
							{Label: "p", Cells: []any{0, "", "", "", "", "", "", "", ""}, Marked: []int{0}},
						}},
						{Caption: "p[1] = p[0] + a[0] = 0 + 3 = 3. Сложений 1.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{3, 1, 4, 1, 5, 9, 2, 6}, Marked: []int{0}},
							{Label: "p", Cells: []any{0, 3, "", "", "", "", "", "", ""}, Marked: []int{1}},
						}},
						{Caption: "p[2] = p[1] + a[1] = 3 + 1 = 4. Сложений 2.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{3, 1, 4, 1, 5, 9, 2, 6}, Marked: []int{1}},
							{Label: "p", Cells: []any{0, 3, 4, "", "", "", "", "", ""}, Marked: []int{2}},
						}},
						{Caption: "p[3] = p[2] + a[2] = 4 + 4 = 8. Сложений 3.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{3, 1, 4, 1, 5, 9, 2, 6}, Marked: []int{2}},
							{Label: "p", Cells: []any{0, 3, 4, 8, "", "", "", "", ""}, Marked: []int{3}},
						}},
						{Caption: "p[4] = p[3] + a[3] = 8 + 1 = 9. Сложений 4.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{3, 1, 4, 1, 5, 9, 2, 6}, Marked: []int{3}},
							{Label: "p", Cells: []any{0, 3, 4, 8, 9, "", "", "", ""}, Marked: []int{4}},
						}},
						{Caption: "p[5] = p[4] + a[4] = 9 + 5 = 14. Сложений 5.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{3, 1, 4, 1, 5, 9, 2, 6}, Marked: []int{4}},
							{Label: "p", Cells: []any{0, 3, 4, 8, 9, 14, "", "", ""}, Marked: []int{5}},
						}},
						{Caption: "p[6] = p[5] + a[5] = 14 + 9 = 23. Сложений 6.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{3, 1, 4, 1, 5, 9, 2, 6}, Marked: []int{5}},
							{Label: "p", Cells: []any{0, 3, 4, 8, 9, 14, 23, "", ""}, Marked: []int{6}},
						}},
						{Caption: "p[7] = p[6] + a[6] = 23 + 2 = 25. Сложений 7.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{3, 1, 4, 1, 5, 9, 2, 6}, Marked: []int{6}},
							{Label: "p", Cells: []any{0, 3, 4, 8, 9, 14, 23, 25, ""}, Marked: []int{7}},
						}},
						{Caption: "p[8] = p[7] + a[7] = 25 + 6 = 31. Сложений 8. Готово: 8 сложений, 9 чисел в p. Время O(n), дополнительная память O(n).", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{3, 1, 4, 1, 5, 9, 2, 6}, Marked: []int{7}},
							{Label: "p", Cells: []any{0, 3, 4, 8, 9, 14, 23, 25, 31}, Marked: []int{8}},
						}},
					},
				},
				{
					ID:    "ps1-query",
					Title: "Сумма отрезка: p[r + 1] − p[l]",
					Frames: []sets.FigureFrame{
						{Caption: "a[2..5] = 4 + 1 + 5 + 9 = 19. p[6] − p[2] = 23 − 4 = 19: одно вычитание.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{3, 1, 4, 1, 5, 9, 2, 6}, Marked: []int{2, 3, 4, 5}, Pointers: &sets.Pointers{{Name: "l", Cell: 2}, {Name: "r", Cell: 5}}, Window: &[2]int{2, 5}},
							{Label: "p", Cells: []any{0, 3, 4, 8, 9, 14, 23, 25, 31}, Marked: []int{2, 6}, Pointers: &sets.Pointers{{Name: "p[l]", Cell: 2}, {Name: "p[r + 1]", Cell: 6}}},
						}},
						{Caption: "a[0..7] = 3 + 1 + 4 + 1 + 5 + 9 + 2 + 6 = 31. p[8] − p[0] = 31 − 0 = 31: одно вычитание. p[0] = 0 убирает отдельный случай l = 0.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{3, 1, 4, 1, 5, 9, 2, 6}, Marked: []int{0, 1, 2, 3, 4, 5, 6, 7}, Pointers: &sets.Pointers{{Name: "l", Cell: 0}, {Name: "r", Cell: 7}}, Window: &[2]int{0, 7}},
							{Label: "p", Cells: []any{0, 3, 4, 8, 9, 14, 23, 25, 31}, Marked: []int{0, 8}, Pointers: &sets.Pointers{{Name: "p[l]", Cell: 0}, {Name: "p[r + 1]", Cell: 8}}},
						}},
						{Caption: "a[3..3] = 1. p[4] − p[3] = 9 − 8 = 1: одно вычитание. Один элемент: l = r.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{3, 1, 4, 1, 5, 9, 2, 6}, Marked: []int{3}, Pointers: &sets.Pointers{{Name: "l", Cell: 3}, {Name: "r", Cell: 3}}, Window: &[2]int{3, 3}},
							{Label: "p", Cells: []any{0, 3, 4, 8, 9, 14, 23, 25, 31}, Marked: []int{3, 4}, Pointers: &sets.Pointers{{Name: "p[l]", Cell: 3}, {Name: "p[r + 1]", Cell: 4}}},
						}},
					},
				},
				{
					ID:    "ps1-pivot",
					Title: "Опорный индекс: nums = [1, 7, 3, 6, 5, 6], total = 28",
					Frames: []sets.FigureFrame{
						{Caption: "i = 0: слева 0, справа 28 − 0 − 1 = 27. Не равны: left += 1.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, 7, 3, 6, 5, 6}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "i", Cell: 0}}},
							{Label: "слева, справа", Cells: []any{0, 27}},
						}},
						{Caption: "i = 1: слева 1, справа 28 − 1 − 7 = 20. Не равны: left += 7.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, 7, 3, 6, 5, 6}, Marked: []int{1}, Pointers: &sets.Pointers{{Name: "i", Cell: 1}}, Window: &[2]int{0, 0}},
							{Label: "слева, справа", Cells: []any{1, 20}},
						}},
						{Caption: "i = 2: слева 8, справа 28 − 8 − 3 = 17. Не равны: left += 3.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, 7, 3, 6, 5, 6}, Marked: []int{2}, Pointers: &sets.Pointers{{Name: "i", Cell: 2}}, Window: &[2]int{0, 1}},
							{Label: "слева, справа", Cells: []any{8, 17}},
						}},
						{Caption: "i = 3: слева 11, справа 28 − 11 − 6 = 11. Равны: ответ 3. Проверено 4 индекса, хранились только total и left.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, 7, 3, 6, 5, 6}, Marked: []int{3}, Pointers: &sets.Pointers{{Name: "i", Cell: 3}}, Window: &[2]int{0, 2}},
							{Label: "слева, справа", Cells: []any{11, 11}, Marked: []int{0, 1}},
						}},
					},
				},
				{
					ID:    "ps1-product",
					Title: "Prefix Product: nums = [1, 2, 3, 4]",
					Frames: []sets.FigureFrame{
						{Caption: "После прохода слева: ans[i] = произведение чисел до i, ans = [1, 1, 2, 6]. 4 умножения.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, 2, 3, 4}},
							{Label: "ans", Cells: []any{1, 1, 2, 6}, Marked: []int{0, 1, 2, 3}},
							{Label: "right", Cells: []any{1}},
						}},
						{Caption: "i = 3: ans[3] = 6 · 1 = 6, затем right = 1 · 4 = 4.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, 2, 3, 4}, Marked: []int{3}},
							{Label: "ans", Cells: []any{1, 1, 2, 6}, Marked: []int{3}},
							{Label: "right", Cells: []any{4}},
						}},
						{Caption: "i = 2: ans[2] = 2 · 4 = 8, затем right = 4 · 3 = 12.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, 2, 3, 4}, Marked: []int{2}},
							{Label: "ans", Cells: []any{1, 1, 8, 6}, Marked: []int{2}},
							{Label: "right", Cells: []any{12}},
						}},
						{Caption: "i = 1: ans[1] = 1 · 12 = 12, затем right = 12 · 2 = 24.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, 2, 3, 4}, Marked: []int{1}},
							{Label: "ans", Cells: []any{1, 12, 8, 6}, Marked: []int{1}},
							{Label: "right", Cells: []any{24}},
						}},
						{Caption: "i = 0: ans[0] = 1 · 24 = 24, затем right = 24 · 1 = 24. Ответ [24, 12, 8, 6] без деления. Дополнительно хранились только left и right: O(1).", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, 2, 3, 4}, Marked: []int{0}},
							{Label: "ans", Cells: []any{24, 12, 8, 6}, Marked: []int{0}},
							{Label: "right", Cells: []any{24}},
						}},
					},
				},
			},
		},
		{
			Title:   "Prefix Sum и Hash Table",
			Summary: "Сумма Subarray равна разности двух префиксов. Поэтому поиск Subarray с суммой k превращается в поиск уже встреченного префикса в Hash Table. Вы считаете такие Subarray, ищете самый длинный Subarray с равным числом 0 и 1 и проверяете, есть ли Subarray с нулевой суммой.",
			Lesson:  lesson2,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/contains-duplicate/",
					Number:     number(217),
					Title:      "Contains Duplicate",
					Difficulty: "easy",
					Topics:     []string{"Array", "Hash Table", "Sorting"},
					Note:       "set seen хранит уже встреченные числа. Для каждого x сначала проверьте x in seen, потом добавьте x. n - длина nums. Время: ожидаемое O(n), до n проверок и n добавлений в set. Память: O(n) для seen, ответ True или False. Для [1, 2, 3, 1]: числа 1, 2 и 3 новые, второе 1 уже в seen, ответ True после 4 проверок. Sorting с проверкой соседей стоит O(n log n). Brute Force по всем парам стоит O(n²), при n = 10^5 почти 5 · 10^9 пар.",
					Statement: `<p>Дан Array целых чисел <code>nums</code>. Верните <code>True</code>, если какое-то значение встречается в Array <strong>хотя бы два раза</strong>, и <code>False</code>, если все элементы разные.</p><p>Метод класса <code>Solution</code>: <code>def containsDuplicate(self, nums: list[int]) -&gt; bool</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,2,3,1]
<strong>Выход:</strong> True
<strong>Пояснение:</strong>
Элемент 1 стоит на индексах 0 и 3.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [1,2,3,4]
<strong>Выход:</strong> False
<strong>Пояснение:</strong>
Все элементы разные.</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> nums = [1,1,1,3,3,4,3,2,4,2]
<strong>Выход:</strong> True</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 10<sup>5</sup></code></li><li><code>-10<sup>9</sup> ≤ nums[i] ≤ 10<sup>9</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "containsDuplicate"}, Cases: []sets.TestCase{
						{Input: "[[1,2,3,1]]", Expected: text("True")},
						{Input: "[[1,2,3,4]]", Expected: text("False")},
						{Input: "[[1,1,1,3,3,4,3,2,4,2]]", Expected: text("True")},
						{Input: "[[7]]", Expected: text("False")},
						{Input: "[[5,5]]", Expected: text("True")},
						{Input: "[[-1000000000,1000000000,0,-1000000000]]", Expected: text("True")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/subarray-sum-equals-k/",
					Number:     number(560),
					Title:      "Subarray Sum Equals K",
					Difficulty: "medium",
					Topics:     []string{"Array", "Hash Table", "Prefix Sum"},
					Note:       "cnt хранит, сколько раз встретилась каждая Prefix Sum; начните с {0: 1}. Перед записью s прибавьте cnt[s − k]. n - длина nums. Время: ожидаемое O(n), n поисков и n записей в dict. Память: O(n), до n + 1 разных префиксов; ответ одно число до n(n + 1)/2. Для [1, −1, 1, 1, −1], k = 1: ответ 7. Sliding Window здесь не работает из-за отрицательных чисел. Brute Force по всем началам с бегущей суммой стоит O(n²).",
					Statement: `<p>Дан Array целых чисел <code>nums</code> и целое число <code>k</code>. Верните <em>общее число Subarray, сумма которых равна</em> <code>k</code>.</p><p>Subarray это <strong>непустая</strong> последовательность соседних элементов Array.</p><p>Метод класса <code>Solution</code>: <code>def subarraySum(self, nums: list[int], k: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,1,1], k = 2
<strong>Выход:</strong> 2</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [1,2,3], k = 3
<strong>Выход:</strong> 2</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 2 · 10<sup>4</sup></code></li><li><code>-1000 ≤ nums[i] ≤ 1000</code></li><li><code>-10<sup>7</sup> ≤ k ≤ 10<sup>7</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "subarraySum"}, Cases: []sets.TestCase{
						{Input: "[[1,1,1],2]", Expected: text("2")},
						{Input: "[[1,2,3],3]", Expected: text("2")},
						{Input: "[[1,-1,1,1,-1],1]", Expected: text("7")},
						{Input: "[[0,0,0],0]", Expected: text("6")},
						{Input: "[[1,-2,1,-2],-1]", Expected: text("3")},
						{Input: "[[5],-5]", Expected: text("0")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/contiguous-array/",
					Number:     number(525),
					Title:      "Contiguous Array",
					Difficulty: "medium",
					Topics:     []string{"Array", "Hash Table", "Prefix Sum"},
					Note:       "Замените 0 на −1: равное число 0 и 1 значит сумму 0. Храните первую позицию каждого баланса и не обновляйте её. n - длина nums. Время: ожидаемое O(n) с dict; list из 2n + 1 мест со сдвигом n даёт O(n) и в худшем случае. Память: O(n), до n + 1 балансов; ответ одно число. Для [0, 1, 1, 1, 1, 1, 0, 0, 0]: баланс 1 впервые на позиции 3 и снова на 9, длина 6. Brute Force по всем Subarray стоит O(n²).",
					Statement: `<p>Дан Array <code>nums</code> из нулей и единиц. Верните <em>наибольшую длину Subarray, в котором поровну</em> <code>0</code> <em>и</em> <code>1</code>.</p><p>Метод класса <code>Solution</code>: <code>def findMaxLength(self, nums: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [0,1]
<strong>Выход:</strong> 2
<strong>Пояснение:</strong> [0, 1] это самый длинный Subarray, в котором поровну 0 и 1.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [0,1,0]
<strong>Выход:</strong> 2
<strong>Пояснение:</strong> [0, 1] (или [1, 0]) это самый длинный Subarray, в котором поровну 0 и 1.</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> nums = [0,1,1,1,1,1,0,0,0]
<strong>Выход:</strong> 6
<strong>Пояснение:</strong> [1,1,1,0,0,0] это самый длинный Subarray, в котором поровну 0 и 1.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 10<sup>5</sup></code></li><li><code>nums[i]</code> равно <code>0</code> или <code>1</code>.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "findMaxLength"}, Cases: []sets.TestCase{
						{Input: "[[0,1]]", Expected: text("2")},
						{Input: "[[0,1,0]]", Expected: text("2")},
						{Input: "[[0,1,1,1,1,1,0,0,0]]", Expected: text("6")},
						{Input: "[[1]]", Expected: text("0")},
						{Input: "[[1,1,1,0]]", Expected: text("2")},
						{Input: "[[0,0,1,0,0,0,1,1]]", Expected: text("6")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1915/E",
					Title:      "Romantic Glasses",
					Difficulty: "medium",
					Topics:     []string{"Array", "Hash Table", "Prefix Sum"},
					Note:       "Codeforces, рейтинг 1300. Нечётные стаканы со знаком +, чётные со знаком −. Равный сок на [l, r] значит повтор префикса со знаками. set хранит встреченные префиксы, начиная с 0. n - число стаканов набора. Время: ожидаемое O(n) на набор, по всем наборам O(∑n). Память: O(n) для set, ответ одна строка на набор. Для [1, 3, 2]: префиксы 0, 1, −2, 0, повтор 0 даёт YES. Brute Force по всем парам l, r стоит O(n²).",
					Statement: `<p>У Юлии <code>n</code> стаканов, они стоят в ряд. В <code>i</code>-м стакане <code>a<sub>i</sub></code> единиц сока. Юлия пьёт только из стаканов с нечётными номерами, а её спутник пьёт только из стаканов с чётными номерами.</p><p>Юлия хочет произвести впечатление на спутника. Для этого она хочет найти Subarray этих стаканов, в котором она и спутник выпьют одинаковое количество сока, если учитывать только стаканы этого Subarray. Помогите ей.</p><p>Точнее: узнайте, есть ли два индекса <code>l</code>, <code>r</code>, такие что <code>1 ≤ l ≤ r ≤ n</code> и <code>a<sub>l</sub> + a<sub>l+2</sub> + a<sub>l+4</sub> + ... + a<sub>r</sub> = a<sub>l+1</sub> + a<sub>l+3</sub> + ... + a<sub>r−1</sub></code>, если <code>l</code> и <code>r</code> одной чётности, и <code>a<sub>l</sub> + a<sub>l+2</sub> + a<sub>l+4</sub> + ... + a<sub>r−1</sub> = a<sub>l+1</sub> + a<sub>l+3</sub> + ... + a<sub>r</sub></code> в другом случае.</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>): число наборов входных данных.</p><p>В первой строке каждого набора записано одно целое число <code>n</code> (<code>1 ≤ n ≤ 2 · 10<sup>5</sup></code>): число стаканов.</p><p>Во второй строке набора записаны <code>n</code> целых чисел <code>a<sub>1</sub>, ..., a<sub>n</sub></code> (<code>1 ≤ a<sub>i</sub> ≤ 10<sup>9</sup></code>): количество сока в каждом стакане.</p><p>Сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого набора выведите <code>YES</code>, если такой Subarray есть, и <code>NO</code> в другом случае.</p><p>Ответ можно выводить в любом регистре. Например, строки <code>yEs</code>, <code>yes</code>, <code>Yes</code> и <code>YES</code> считаются положительным ответом.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>6
3
1 3 2
6
1 1 1 1 1 1
10
1 6 9 8 55 3 14 2 7 2
8
1 2 11 4 1 5 1 2
6
2 6 1 5 7 8
9
2 5 10 4 4 9 6 7 8</pre><p><strong>Выходные данные:</strong></p><pre>YES
YES
NO
YES
NO
YES</pre><h3>Примечание</h3><p>В первом наборе Юлия может выбрать <code>l = 1</code> и <code>r = 3</code>. Тогда она выпьет <code>a<sub>1</sub> + a<sub>3</sub> = 1 + 2 = 3</code> единицы сока, а спутник выпьет <code>a<sub>2</sub> = 3</code> единицы.</p><p>Во втором наборе Юлия может выбрать <code>l = 2</code> и <code>r = 5</code>. Тогда она выпьет <code>a<sub>3</sub> + a<sub>5</sub> = 1 + 1 = 2</code> единицы сока, а спутник выпьет <code>a<sub>2</sub> + a<sub>4</sub> = 1 + 1 = 2</code> единицы.</p><p>В третьем наборе такого Subarray нет.</p><p>В четвёртом наборе Юлия может выбрать <code>l = 2</code> и <code>r = 8</code>. Тогда она выпьет <code>a<sub>3</sub> + a<sub>5</sub> + a<sub>7</sub> = 11 + 1 + 1 = 13</code> единиц сока, а спутник выпьет <code>a<sub>2</sub> + a<sub>4</sub> + a<sub>6</sub> + a<sub>8</sub> = 2 + 4 + 5 + 2 = 13</code> единиц.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>1 ≤ n ≤ 2 · 10<sup>5</sup></code></li><li><code>1 ≤ a<sub>i</sub> ≤ 10<sup>9</sup></code></li><li>сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code></li><li>Ограничение по времени: 1 секунда на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio", IgnoreCase: true}, Cases: []sets.TestCase{
						{Input: `6
3
1 3 2
6
1 1 1 1 1 1
10
1 6 9 8 55 3 14 2 7 2
8
1 2 11 4 1 5 1 2
6
2 6 1 5 7 8
9
2 5 10 4 4 9 6 7 8
`, Expected: text(`YES
YES
NO
YES
NO
YES
`)},
						{Input: `1
1
5
`, Expected: text(`NO
`)},
						{Input: `1
2
7 7
`, Expected: text(`YES
`)},
						{Input: `1
4
1 2 3 2
`, Expected: text(`YES
`)},
						{Input: `2
3
1 2 4
2
3 5
`, Expected: text(`NO
NO
`)},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1398/C",
					Title:      "Good Subarrays",
					Difficulty: "medium",
					Topics:     []string{"String", "Hash Table", "Prefix Sum", "Counting"},
					Note:       "Codeforces, рейтинг 1600. Вычтите 1 из каждой цифры: у хорошего Subarray сумма 0. Считайте пары равных префиксов: перед записью s прибавьте cnt[s]. n - длина строки набора. Время: ожидаемое O(n) на набор, по всем наборам O(∑n). Память: O(n) для cnt, ответ одно число до n(n + 1)/2 = 5 000 050 000 при n = 10^5. Для 120: префиксы 0, 0, 1, 0, ответ 3. Brute Force по всем Subarray стоит O(n²).",
					Statement: `<p>Дан Array <code>a<sub>1</sub>, a<sub>2</sub>, ..., a<sub>n</sub></code> из целых чисел от <code>0</code> до <code>9</code>. Subarray <code>a<sub>l</sub>, a<sub>l+1</sub>, a<sub>l+2</sub>, ..., a<sub>r−1</sub>, a<sub>r</sub></code> называется хорошим, если сумма его элементов равна его длине (<code>a<sub>l</sub> + a<sub>l+1</sub> + ... + a<sub>r</sub> = r − l + 1</code>).</p><p>Например, для <code>a = [1, 2, 0]</code> есть <code>3</code> хороших Subarray: <code>a<sub>1...1</sub> = [1]</code>, <code>a<sub>2...3</sub> = [2, 0]</code> и <code>a<sub>1...3</sub> = [1, 2, 0]</code>.</p><p>Посчитайте число хороших Subarray в Array <code>a</code>.</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>t</code> (<code>1 ≤ t ≤ 1000</code>): число наборов входных данных.</p><p>В первой строке каждого набора записано одно целое число <code>n</code> (<code>1 ≤ n ≤ 10<sup>5</sup></code>): длина Array <code>a</code>.</p><p>Во второй строке набора записана строка из <code>n</code> десятичных цифр. <code>i</code>-я цифра равна <code>a<sub>i</sub></code>.</p><p>Гарантируется, что сумма <code>n</code> по всем наборам не больше <code>10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого набора выведите одно целое число: число хороших Subarray в Array <code>a</code>.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>3
3
120
5
11011
6
600005</pre><p><strong>Выходные данные:</strong></p><pre>3
6
1</pre><h3>Примечание</h3><p>Первый набор разобран в условии.</p><p>Во втором наборе <code>6</code> хороших Subarray: <code>a<sub>1...1</sub></code>, <code>a<sub>2...2</sub></code>, <code>a<sub>1...2</sub></code>, <code>a<sub>4...4</sub></code>, <code>a<sub>5...5</sub></code> и <code>a<sub>4...5</sub></code>.</p><p>В третьем наборе только один хороший Subarray: <code>a<sub>2...6</sub></code>.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 1000</code></li><li><code>1 ≤ n ≤ 10<sup>5</sup></code></li><li>строка набора состоит из <code>n</code> десятичных цифр</li><li>сумма <code>n</code> по всем наборам не больше <code>10<sup>5</sup></code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `3
3
120
5
11011
6
600005
`, Expected: text(`3
6
1
`)},
						{Input: `1
1
1
`, Expected: text(`1
`)},
						{Input: `1
1
0
`, Expected: text(`0
`)},
						{Input: `1
4
1201
`, Expected: text(`6
`)},
						{Input: `1
5
11111
`, Expected: text(`15
`)},
						{Input: `2
2
20
3
111
`, Expected: text(`1
6
`)},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "ps2-check-first", Type: "code", Prompt: "LeetCode 560. nums = [1, −1], k = 0. Что вернёт этот код?", Language: "python", Code: `cnt = {0: 1}
s = 0
ans = 0
for x in nums:
    s += x
    cnt[s] = cnt.get(s, 0) + 1
    ans += cnt.get(s - k, 0)
return ans`, Options: []string{"1", "3", "2"}, Answer: 1, Explain: "Код увеличивает cnt[s] до поиска, поэтому при k = 0 каждый префикс находит сам себя. x = 1: s = 1, cnt[1] = 1, ответ 1. x = −1: s = 0, cnt[0] = 2, ответ 1 + 2 = 3. Верный ответ 1: сумму 0 имеет только [1, −1]. Каждый лишний счёт это пара префикса с самим собой, то есть пустой отрезок. Сначала ищите s − k, потом увеличивайте cnt[s]. Время ожидаемое O(n), память O(n)."},
				{ID: "ps2-zero", Type: "short", Prompt: "LeetCode 560. nums = [2, −2, 2, −2], k = 0. Сколько Subarray имеют сумму 0?", Answer: "4", Explain: "Префиксы 0, 2, 0, 2, 0. Значение 0 встречается 3 раза и даёт 3 пары, значение 2 встречается 2 раза и даёт 1 пару: 3 + 1 = 4. Это три Subarray длины 2 и весь Array. Ответ 3 получится, если забыть весь Array и считать только Subarray длины 2."},
				{ID: "ps2-empty", Type: "short", Prompt: "LeetCode 560. nums = [3, 4, 7], k = 7. Код начинает с пустого cnt вместо {0: 1}. Какой ответ он вернёт?", Answer: "1", Explain: "Префиксы 3, 7, 14. На s = 7 нужен префикс 0, но его нет в cnt, и Subarray [3, 4] теряется. На s = 14 нужен префикс 7: он есть, это [7]. Код вернёт 1, а верный ответ 2. Ключ 0 с числом 1 означает пустое начало до первого элемента."},
				{ID: "ps2-first", Type: "choice", Prompt: "LeetCode 525 ищет самый длинный Subarray. Что хранить в Hash Table для каждого баланса?", Options: []string{"Первую позицию баланса", "Последнюю позицию баланса", "Число встреч баланса"}, Answer: 0, Explain: "Длина равна i − first[b], она больше при более ранней позиции. Для [0, 1, 0, 1] баланс 0 бывает на позициях 0, 2 и 4. С первой позицией длина 4 − 0 = 4, весь Array. С последней позицией получится 4 − 2 = 2. Число встреч нужно для подсчёта, как в LeetCode 560."},
				{ID: "ps2-romantic", Type: "short", Prompt: "Codeforces 1915E. a = [3, 1, 2, 4]. Нечётные места со знаком +, чётные со знаком −. После какого стакана сумма со знаками впервые повторяет более раннее значение?", Answer: "4", Explain: "Префиксы: 0, 3, 2, 4, 0. Значение 0 повторилось после 4-го стакана: оно было и в начале, до первого стакана. Значит стаканы с 1-го по 4-й дают сумму 0: 3 + 2 = 5 и 1 + 4 = 5, ответ задачи YES. Без начального 0 повтора нет, и решение ошибочно ответит NO. Проверка только соседних стаканов тоже ничего не находит: 3 ≠ 1, 1 ≠ 2 и 2 ≠ 4."},
				{ID: "ps2-good", Type: "short", Prompt: "Codeforces 1398C. Строка 1201. Сколько Subarray, у которых сумма цифр равна длине?", Answer: "6", Explain: "Вычтем 1: [0, 1, −1, 0]. Префиксы 0, 0, 1, 0, 0: значение 0 встречается 4 раза, это 4 · 3 / 2 = 6 пар. Подходят [1] в начале, [1] в конце, [2, 0], [1, 2, 0], [2, 0, 1] и вся строка. Без вычитания 1 пришлось бы проверять все 10 пар границ."},
				{ID: "ps2-cost", Type: "choice", Prompt: "LeetCode 560, n = 20 000. Сколько Subarray перебирает Brute Force?", Options: []string{"20 000", "около 400 миллионов", "около 200 миллионов"}, Answer: 2, Explain: "Subarray задают начало и конец, начало не правее конца: n(n + 1)/2 = 200 010 000. Даже с бегущей суммой это O(n²), больше 3 с в Python. Prefix Sum с Hash Table делает 20 000 поисков, ожидаемое O(n). 400 миллионов это n², упорядоченные пары с повтором."},
			},
			Links: []sets.Link{
				{Title: "Subarray с нулевой суммой", URL: "https://peltorator.org/posts/prefix_sums/#продвинутые-применения", Source: "peltorator"},
				{Title: "Hash Table: устройство", URL: "https://neerc.ifmo.ru/wiki/index.php?title=%D0%A5%D0%B5%D1%88-%D1%82%D0%B0%D0%B1%D0%BB%D0%B8%D1%86%D0%B0", Source: "Викиконспекты ИТМО"},
				{Title: "Число Subarray с суммой x", URL: "https://usaco.guide/problems/cses-1661-subarray-sums-ii/solution", Source: "USACO Guide (англ.)"},
				{Title: "dict.get: значение по умолчанию", URL: "https://docs.python.org/3/builtins/stdtypes.html#dict.get", Source: "Python docs (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "ps2-count",
					Title: "Subarray с суммой k = 1: nums = [1, −1, 1, 1, −1]",
					Frames: []sets.FigureFrame{
						{Caption: "Начало: cnt = {0: 1}, это пустой префикс. Ответ 0.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, -1, 1, 1, -1}},
							{Label: "s, s − k, +, ответ", Cells: []any{0, "", "", 0}},
							{Label: "ключ s", Cells: []any{0, "", ""}, Marked: []int{0}},
							{Label: "встреч", Cells: []any{1, "", ""}, Marked: []int{0}},
						}},
						{Caption: "x = 1: s = 1, ищем s − k = 0: прибавляем 1, ответ 1. Потом cnt[1] = 1.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, -1, 1, 1, -1}, Marked: []int{0}},
							{Label: "s, s − k, +, ответ", Cells: []any{1, 0, 1, 1}, Marked: []int{2}},
							{Label: "ключ s", Cells: []any{0, 1, ""}, Marked: []int{0, 1}},
							{Label: "встреч", Cells: []any{1, 1, ""}, Marked: []int{0, 1}},
						}},
						{Caption: "x = −1: s = 0, ищем s − k = −1: прибавляем 0, ответ 1. Потом cnt[0] = 2.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, -1, 1, 1, -1}, Marked: []int{1}},
							{Label: "s, s − k, +, ответ", Cells: []any{0, -1, 0, 1}, Marked: []int{2}},
							{Label: "ключ s", Cells: []any{0, 1, ""}, Marked: []int{0}},
							{Label: "встреч", Cells: []any{2, 1, ""}, Marked: []int{0}},
						}},
						{Caption: "x = 1: s = 1, ищем s − k = 0: прибавляем 2, ответ 3. Потом cnt[1] = 2.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, -1, 1, 1, -1}, Marked: []int{2}},
							{Label: "s, s − k, +, ответ", Cells: []any{1, 0, 2, 3}, Marked: []int{2}},
							{Label: "ключ s", Cells: []any{0, 1, ""}, Marked: []int{0, 1}},
							{Label: "встреч", Cells: []any{2, 2, ""}, Marked: []int{0, 1}},
						}},
						{Caption: "x = 1: s = 2, ищем s − k = 1: прибавляем 2, ответ 5. Потом cnt[2] = 1.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, -1, 1, 1, -1}, Marked: []int{3}},
							{Label: "s, s − k, +, ответ", Cells: []any{2, 1, 2, 5}, Marked: []int{2}},
							{Label: "ключ s", Cells: []any{0, 1, 2}, Marked: []int{1, 2}},
							{Label: "встреч", Cells: []any{2, 2, 1}, Marked: []int{1, 2}},
						}},
						{Caption: "x = −1: s = 1, ищем s − k = 0: прибавляем 2, ответ 7. Потом cnt[1] = 3. Итог 7: 5 поисков и 5 записей в dict, пик 3 ключа; в худшем случае n + 1 = 6.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, -1, 1, 1, -1}, Marked: []int{4}},
							{Label: "s, s − k, +, ответ", Cells: []any{1, 0, 2, 7}, Marked: []int{2}},
							{Label: "ключ s", Cells: []any{0, 1, 2}, Marked: []int{0, 1}},
							{Label: "встреч", Cells: []any{2, 3, 1}, Marked: []int{0, 1}},
						}},
					},
				},
				{
					ID:    "ps2-balance",
					Title: "Баланс 0 и 1: nums = [0, 1, 1, 1, 1, 1, 0, 0, 0]",
					Frames: []sets.FigureFrame{
						{Caption: "Позиция 1: баланс −1 впервые, запоминаем позицию 1. Лучший 0.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{0, 1, 1, 1, 1, 1, 0, 0, 0}, Marked: []int{0}},
							{Label: "баланс", Cells: []any{0, -1, "", "", "", "", "", "", "", ""}, Marked: []int{1}, Pointers: &sets.Pointers{{Name: "сейчас", Cell: 1}}},
						}},
						{Caption: "Позиция 2: баланс 0 уже был на позиции 0. Длина 2 − 0 = 2, лучший 2.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{0, 1, 1, 1, 1, 1, 0, 0, 0}, Marked: []int{1}, Window: &[2]int{0, 1}},
							{Label: "баланс", Cells: []any{0, -1, 0, "", "", "", "", "", "", ""}, Marked: []int{2}, Pointers: &sets.Pointers{{Name: "первый", Cell: 0}, {Name: "сейчас", Cell: 2}}},
						}},
						{Caption: "Позиция 3: баланс 1 впервые, запоминаем позицию 3. Лучший 2.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{0, 1, 1, 1, 1, 1, 0, 0, 0}, Marked: []int{2}},
							{Label: "баланс", Cells: []any{0, -1, 0, 1, "", "", "", "", "", ""}, Marked: []int{3}, Pointers: &sets.Pointers{{Name: "сейчас", Cell: 3}}},
						}},
						{Caption: "Позиция 4: баланс 2 впервые, запоминаем позицию 4. Лучший 2.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{0, 1, 1, 1, 1, 1, 0, 0, 0}, Marked: []int{3}},
							{Label: "баланс", Cells: []any{0, -1, 0, 1, 2, "", "", "", "", ""}, Marked: []int{4}, Pointers: &sets.Pointers{{Name: "сейчас", Cell: 4}}},
						}},
						{Caption: "Позиция 5: баланс 3 впервые, запоминаем позицию 5. Лучший 2.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{0, 1, 1, 1, 1, 1, 0, 0, 0}, Marked: []int{4}},
							{Label: "баланс", Cells: []any{0, -1, 0, 1, 2, 3, "", "", "", ""}, Marked: []int{5}, Pointers: &sets.Pointers{{Name: "сейчас", Cell: 5}}},
						}},
						{Caption: "Позиция 6: баланс 4 впервые, запоминаем позицию 6. Лучший 2.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{0, 1, 1, 1, 1, 1, 0, 0, 0}, Marked: []int{5}},
							{Label: "баланс", Cells: []any{0, -1, 0, 1, 2, 3, 4, "", "", ""}, Marked: []int{6}, Pointers: &sets.Pointers{{Name: "сейчас", Cell: 6}}},
						}},
						{Caption: "Позиция 7: баланс 3 уже был на позиции 5. Длина 7 − 5 = 2, лучший 2.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{0, 1, 1, 1, 1, 1, 0, 0, 0}, Marked: []int{6}, Window: &[2]int{5, 6}},
							{Label: "баланс", Cells: []any{0, -1, 0, 1, 2, 3, 4, 3, "", ""}, Marked: []int{7}, Pointers: &sets.Pointers{{Name: "первый", Cell: 5}, {Name: "сейчас", Cell: 7}}},
						}},
						{Caption: "Позиция 8: баланс 2 уже был на позиции 4. Длина 8 − 4 = 4, лучший 4.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{0, 1, 1, 1, 1, 1, 0, 0, 0}, Marked: []int{7}, Window: &[2]int{4, 7}},
							{Label: "баланс", Cells: []any{0, -1, 0, 1, 2, 3, 4, 3, 2, ""}, Marked: []int{8}, Pointers: &sets.Pointers{{Name: "первый", Cell: 4}, {Name: "сейчас", Cell: 8}}},
						}},
						{Caption: "Позиция 9: баланс 1 уже был на позиции 3. Длина 9 − 3 = 6, лучший 6. Итог 6: 9 шагов, в first 6 ключей.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{0, 1, 1, 1, 1, 1, 0, 0, 0}, Marked: []int{8}, Window: &[2]int{3, 8}},
							{Label: "баланс", Cells: []any{0, -1, 0, 1, 2, 3, 4, 3, 2, 1}, Marked: []int{9}, Pointers: &sets.Pointers{{Name: "первый", Cell: 3}, {Name: "сейчас", Cell: 9}}},
						}},
					},
				},
			},
		},
		{
			Title:   "Difference Array",
			Summary: "Difference Array прибавляет число ко всему отрезку двумя изменениями, а одна Prefix Sum в конце восстанавливает все значения. Вы считаете население по годам и места на рейсах, проверяете, сколько рецептов покрывают температуру, и применяете Difference Array два раза подряд.",
			Lesson:  lesson3,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/maximum-population-year/",
					Number:     number(1854),
					Title:      "Maximum Population Year",
					Difficulty: "easy",
					Topics:     []string{"Array", "Counting", "Prefix Sum"},
					Note:       "Difference Array по годам: d[birth - 1950] += 1, d[death - 1950] -= 1, год смерти не считаем. n - число людей, Y = 101 год. Время: O(n + Y). Память: O(Y) для d, ответ одно число. Для [[1950, 1961], [1960, 1971], [1970, 1981]]: население 2 впервые в 1960. Проверка каждого года по всем людям стоит O(n · Y); при этих ограничениях она тоже быстрая.",
					Statement: `<p>Дан двумерный Array целых чисел <code>logs</code>. Каждый элемент <code>logs[i] = [birth<sub>i</sub>, death<sub>i</sub>]</code> задаёт год рождения и год смерти <code>i</code>-го человека.</p><p><strong>Население</strong> года <code>x</code> это число людей, которые жили в этом году. <code>i</code>-го человека считают в населении года <code>x</code>, если <code>x</code> лежит в отрезке <code>[birth<sub>i</sub>, death<sub>i</sub> − 1]</code> <strong>включительно</strong>. Обратите внимание: человека <strong>не</strong> считают в году его смерти.</p><p>Верните <em><strong>самый ранний</strong> год с <strong>наибольшим населением</strong></em>.</p><p>Метод класса <code>Solution</code>: <code>def maximumPopulation(self, logs: list[list[int]]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> logs = [[1993,1999],[2000,2010]]
<strong>Выход:</strong> 1993
<strong>Пояснение:</strong> Наибольшее население равно 1, и 1993 это самый ранний год с таким населением.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> logs = [[1950,1961],[1960,1971],[1970,1981]]
<strong>Выход:</strong> 1960
<strong>Пояснение:</strong>
Наибольшее население равно 2, так было в 1960 и в 1970 году.
Более ранний из этих двух годов 1960.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ logs.length ≤ 100</code></li><li><code>1950 ≤ birth<sub>i</sub> &lt; death<sub>i</sub> ≤ 2050</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "maximumPopulation"}, Cases: []sets.TestCase{
						{Input: "[[[1993,1999],[2000,2010]]]", Expected: text("1993")},
						{Input: "[[[1950,1961],[1960,1971],[1970,1981]]]", Expected: text("1960")},
						{Input: "[[[1950,2050]]]", Expected: text("1950")},
						{Input: "[[[2000,2001],[2001,2002]]]", Expected: text("2000")},
						{Input: "[[[2040,2050],[2045,2050],[2049,2050]]]", Expected: text("2049")},
						{Input: "[[[1980,1990],[1985,1995],[1985,1987],[1990,2000]]]", Expected: text("1985")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/corporate-flight-bookings/",
					Number:     number(1109),
					Title:      "Corporate Flight Bookings",
					Difficulty: "medium",
					Topics:     []string{"Array", "Prefix Sum"},
					Note:       "Difference Array длины n + 1: d[first - 1] += seats, d[last] -= seats. n - число рейсов, b - число бронирований. Время: O(n + b): 2 изменения на бронирование и n сложений Prefix Sum. Память: O(n) для d, ответ O(n). Для первого примера d = [10, 45, −10, −20, 0, −25]. Цикл по всем рейсам каждого бронирования стоит O(n · b), до 4 · 10^8 прибавлений.",
					Statement: `<p>Есть <code>n</code> рейсов, они пронумерованы от <code>1</code> до <code>n</code>.</p><p>Дан Array бронирований <code>bookings</code>, где <code>bookings[i] = [first<sub>i</sub>, last<sub>i</sub>, seats<sub>i</sub>]</code> означает бронирование рейсов с <code>first<sub>i</sub></code> по <code>last<sub>i</sub></code> (<strong>включительно</strong>): на <strong>каждом рейсе</strong> этого отрезка забронировано <code>seats<sub>i</sub></code> мест.</p><p>Верните <em>Array</em> <code>answer</code> <em>длины</em> <code>n</code><em>, где</em> <code>answer[i]</code> <em>это общее число мест, забронированных на рейс</em> <code>i</code>.</p><p>Метод класса <code>Solution</code>: <code>def corpFlightBookings(self, bookings: list[list[int]], n: int) -&gt; list[int]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> bookings = [[1,2,10],[2,3,20],[2,5,25]], n = 5
<strong>Выход:</strong> [10,55,45,25,25]
<strong>Пояснение:</strong>
Номера рейсов:    1   2   3   4   5
Бронирование 1:  10  10
Бронирование 2:      20  20
Бронирование 3:      25  25  25  25
Всего мест:      10  55  45  25  25
Поэтому answer = [10,55,45,25,25]</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> bookings = [[1,2,10],[2,2,15]], n = 2
<strong>Выход:</strong> [10,25]
<strong>Пояснение:</strong>
Номера рейсов:    1   2
Бронирование 1:  10  10
Бронирование 2:      15
Всего мест:      10  25
Поэтому answer = [10,25]</pre><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 2 · 10<sup>4</sup></code></li><li><code>1 ≤ bookings.length ≤ 2 · 10<sup>4</sup></code></li><li><code>bookings[i].length == 3</code></li><li><code>1 ≤ first<sub>i</sub> ≤ last<sub>i</sub> ≤ n</code></li><li><code>1 ≤ seats<sub>i</sub> ≤ 10<sup>4</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "corpFlightBookings"}, Cases: []sets.TestCase{
						{Input: "[[[1,2,10],[2,3,20],[2,5,25]],5]", Expected: text("[10,55,45,25,25]")},
						{Input: "[[[1,2,10],[2,2,15]],2]", Expected: text("[10,25]")},
						{Input: "[[[1,1,7]],1]", Expected: text("[7]")},
						{Input: "[[[3,3,5],[1,3,2]],4]", Expected: text("[2,2,7,0]")},
						{Input: "[[[1,5,10000],[5,5,10000]],5]", Expected: text("[10000,10000,10000,10000,20000]")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/816/B",
					Title:      "Karen and Coffee",
					Difficulty: "medium",
					Topics:     []string{"Array", "Prefix Sum"},
					Note:       "Codeforces, рейтинг 1400. Difference Array даёт число рецептов для каждой температуры. Затем отметка 1 для температуры с числом не меньше k и Prefix Sum по отметкам. T = 200 000 температур, n рецептов, q вопросов. Время: O(n + T + q). Память: O(T) для d и префиксов, ответ q чисел. В первом примере подходят 92, 93, 94 и 97. Подсчёт по каждой температуре каждого вопроса стоит O(q · T).",
					Statement: `<p>Чтобы не спать и быть внимательной на уроках, Карен нужен кофе!</p><p>Карен знает толк в кофе и хочет узнать лучшую температуру, чтобы заварить идеальную чашку. Она прочитала несколько книг рецептов, в том числе всеми признанную «The Art of the Covfefe».</p><p>Карен знает <code>n</code> рецептов кофе. <code>i</code>-й рецепт советует заваривать кофе при температуре от <code>l<sub>i</sub></code> до <code>r<sub>i</sub></code> градусов включительно: так получится лучший вкус.</p><p>Карен считает температуру допустимой, если её рекомендуют не меньше <code>k</code> рецептов.</p><p>Карен часто меняет решения, поэтому она задаёт <code>q</code> вопросов. В каждом вопросе она хочет готовить кофе только при температуре от <code>a</code> до <code>b</code> градусов включительно. Скажите ей, сколько допустимых целых температур лежит в этом отрезке.</p><h3>Входные данные</h3><p>В первой строке записаны три целых числа <code>n</code>, <code>k</code> (<code>1 ≤ k ≤ n ≤ 200000</code>) и <code>q</code> (<code>1 ≤ q ≤ 200000</code>): число рецептов, наименьшее число рецептов, которые должны рекомендовать температуру, чтобы она была допустимой, и число вопросов Карен.</p><p>Следующие <code>n</code> строк описывают рецепты. <code>i</code>-я из этих строк содержит два целых числа <code>l<sub>i</sub></code> и <code>r<sub>i</sub></code> (<code>1 ≤ l<sub>i</sub> ≤ r<sub>i</sub> ≤ 200000</code>): <code>i</code>-й рецепт советует заваривать кофе при температуре от <code>l<sub>i</sub></code> до <code>r<sub>i</sub></code> градусов включительно.</p><p>Следующие <code>q</code> строк описывают вопросы. Каждая из них содержит два целых числа <code>a</code> и <code>b</code> (<code>1 ≤ a ≤ b ≤ 200000</code>): Карен хочет знать число допустимых целых температур от <code>a</code> до <code>b</code> градусов включительно.</p><h3>Выходные данные</h3><p>Для каждого вопроса выведите одно целое число в отдельной строке: число допустимых целых температур от <code>a</code> до <code>b</code> градусов включительно.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>3 2 4
91 94
92 97
97 99
92 94
93 97
95 96
90 100</pre><p><strong>Выходные данные:</strong></p><pre>3
3
0
4</pre><h3>Пример 2</h3><p><strong>Входные данные:</strong></p><pre>2 1 1
1 1
200000 200000
90 100</pre><p><strong>Выходные данные:</strong></p><pre>0</pre><h3>Примечание</h3><p>В первом примере Карен знает 3 рецепта.</p><ol><li>Первый советует заваривать кофе при температуре от 91 до 94 градусов включительно.</li><li>Второй советует заваривать кофе при температуре от 92 до 97 градусов включительно.</li><li>Третий советует заваривать кофе при температуре от 97 до 99 градусов включительно.</li></ol><p>Температура допустима, если её рекомендуют не меньше 2 рецептов.</p><p>Карен задаёт 4 вопроса.</p><p>В первом вопросе она хочет знать число допустимых целых температур от 92 до 94 градусов включительно. Их 3: температуры 92, 93 и 94 градуса допустимы.</p><p>Во втором вопросе она хочет знать число допустимых целых температур от 93 до 97 градусов включительно. Их 3: температуры 93, 94 и 97 градусов допустимы.</p><p>В третьем вопросе она хочет знать число допустимых целых температур от 95 до 96 градусов включительно. Таких температур нет.</p><p>В последнем вопросе она хочет знать число допустимых целых температур от 90 до 100 градусов включительно. Их 4: температуры 92, 93, 94 и 97 градусов допустимы.</p><p>Во втором примере Карен знает 2 рецепта.</p><ol><li>Первый, «wikiHow to make Cold Brew Coffee», советует заваривать кофе ровно при 1 градусе.</li><li>Второй, «What good is coffee that isn't brewed at at least 36.3306 times the temperature of the surface of the sun?», советует заваривать кофе ровно при 200000 градусах.</li></ol><p>Температура допустима, если её рекомендует хотя бы 1 рецепт.</p><p>В первом и единственном вопросе она хочет знать число допустимых целых температур, которые на самом деле разумны. Таких температур нет.</p><h3>Ограничения</h3><ul><li><code>1 ≤ k ≤ n ≤ 200000</code></li><li><code>1 ≤ q ≤ 200000</code></li><li><code>1 ≤ l<sub>i</sub> ≤ r<sub>i</sub> ≤ 200000</code></li><li><code>1 ≤ a ≤ b ≤ 200000</code></li><li>Ограничение по времени: 2,5 секунды на тест.</li><li>Ограничение по памяти: 512 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `3 2 4
91 94
92 97
97 99
92 94
93 97
95 96
90 100
`, Expected: text(`3
3
0
4
`)},
						{Input: `2 1 1
1 1
200000 200000
90 100
`, Expected: text(`0
`)},
						{Input: `1 1 1
1 200000
1 200000
`, Expected: text(`200000
`)},
						{Input: `3 3 2
1 10
5 6
6 20
1 20
6 6
`, Expected: text(`1
1
`)},
						{Input: `2 2 2
199999 200000
200000 200000
1 200000
199999 199999
`, Expected: text(`1
0
`)},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/295/A",
					Title:      "Greg and Array",
					Difficulty: "medium",
					Topics:     []string{"Array", "Prefix Sum"},
					Note:       "Codeforces, рейтинг 1400. Первый Difference Array по номерам операций считает cnt[i], сколько запросов применяют операцию i. Второй прибавляет d · cnt[i] на [l, r] этой операции. n, m и k - размеры входа. Время: O(n + m + k). Память: O(n + m) для двух списков разностей, ответ n чисел. Прибавка к одному элементу до 10^5 · 10^5 · 10^5 = 10^15 помещается в int Python. Для первого примера cnt = [2, 3, 2], ответ 9 18 17. Прямое применение операций стоит O(k · m · n) в худшем случае.",
					Statement: `<p>У Грега есть Array <code>a = a<sub>1</sub>, a<sub>2</sub>, ..., a<sub>n</sub></code> и <code>m</code> операций. Каждая операция имеет вид <code>l<sub>i</sub>, r<sub>i</sub>, d<sub>i</sub></code> (<code>1 ≤ l<sub>i</sub> ≤ r<sub>i</sub> ≤ n</code>). Применить операцию <code>i</code> к Array значит увеличить все элементы Array с номерами <code>l<sub>i</sub>, l<sub>i</sub> + 1, ..., r<sub>i</sub></code> на значение <code>d<sub>i</sub></code>.</p><p>Грег записал на листке бумаги <code>k</code> запросов. Каждый запрос имеет вид <code>x<sub>i</sub>, y<sub>i</sub></code> (<code>1 ≤ x<sub>i</sub> ≤ y<sub>i</sub> ≤ m</code>). Это значит, что к Array нужно применить операции с номерами <code>x<sub>i</sub>, x<sub>i</sub> + 1, ..., y<sub>i</sub></code>.</p><p>Теперь Грегу интересно, каким будет Array <code>a</code> после выполнения всех запросов. Помогите Грегу.</p><h3>Входные данные</h3><p>В первой строке записаны целые числа <code>n</code>, <code>m</code>, <code>k</code> (<code>1 ≤ n, m, k ≤ 10<sup>5</sup></code>). Во второй строке записаны <code>n</code> целых чисел <code>a<sub>1</sub>, a<sub>2</sub>, ..., a<sub>n</sub></code> (<code>0 ≤ a<sub>i</sub> ≤ 10<sup>5</sup></code>): начальный Array.</p><p>Следующие <code>m</code> строк содержат операции. Операция с номером <code>i</code> записана тремя целыми числами <code>l<sub>i</sub>, r<sub>i</sub>, d<sub>i</sub></code> (<code>1 ≤ l<sub>i</sub> ≤ r<sub>i</sub> ≤ n</code>), (<code>0 ≤ d<sub>i</sub> ≤ 10<sup>5</sup></code>).</p><p>Следующие <code>k</code> строк содержат запросы. Запрос с номером <code>i</code> записан двумя целыми числами <code>x<sub>i</sub>, y<sub>i</sub></code> (<code>1 ≤ x<sub>i</sub> ≤ y<sub>i</sub> ≤ m</code>).</p><p>Числа в строках разделены одиночными пробелами.</p><h3>Выходные данные</h3><p>В одной строке выведите <code>n</code> целых чисел <code>a<sub>1</sub>, a<sub>2</sub>, ..., a<sub>n</sub></code>: Array после выполнения всех запросов. Разделяйте выведенные числа пробелами.</p><p>Не используйте <code>%lld</code> для чтения и записи 64-битных целых чисел в C++. Лучше использовать <code>cin</code> и <code>cout</code> или <code>%I64d</code>.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>3 3 3
1 2 3
1 2 1
1 3 2
2 3 4
1 2
1 3
2 3</pre><p><strong>Выходные данные:</strong></p><pre>9 18 17</pre><h3>Пример 2</h3><p><strong>Входные данные:</strong></p><pre>1 1 1
1
1 1 1
1 1</pre><p><strong>Выходные данные:</strong></p><pre>2</pre><h3>Пример 3</h3><p><strong>Входные данные:</strong></p><pre>4 3 6
1 2 3 4
1 2 1
2 3 2
3 4 4
1 2
1 3
2 3
1 2
1 3
2 3</pre><p><strong>Выходные данные:</strong></p><pre>5 18 31 20</pre><h3>Ограничения</h3><ul><li><code>1 ≤ n, m, k ≤ 10<sup>5</sup></code></li><li><code>0 ≤ a<sub>i</sub> ≤ 10<sup>5</sup></code></li><li><code>1 ≤ l<sub>i</sub> ≤ r<sub>i</sub> ≤ n</code>, <code>0 ≤ d<sub>i</sub> ≤ 10<sup>5</sup></code></li><li><code>1 ≤ x<sub>i</sub> ≤ y<sub>i</sub> ≤ m</code></li><li>Ограничение по времени: 1,5 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `3 3 3
1 2 3
1 2 1
1 3 2
2 3 4
1 2
1 3
2 3
`, Expected: text(`9 18 17
`)},
						{Input: `1 1 1
1
1 1 1
1 1
`, Expected: text(`2
`)},
						{Input: `4 3 6
1 2 3 4
1 2 1
2 3 2
3 4 4
1 2
1 3
2 3
1 2
1 3
2 3
`, Expected: text(`5 18 31 20
`)},
						{Input: `2 1 1
5 6
1 2 0
1 1
`, Expected: text(`5 6
`)},
						{Input: `2 2 3
0 0
1 1 100000
2 2 100000
1 2
1 2
2 2
`, Expected: text(`200000 300000
`)},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "ps3-marks", Type: "choice", Prompt: "n = 5, d длины 6. Нужно прибавить 4 ко всем a[1..3], номера с 0. Какие изменения d верны?", Options: []string{"d[1] += 4 и d[3] -= 4", "d[1] += 4 и d[4] -= 4", "d[0] += 4 и d[3] -= 4"}, Answer: 1, Explain: "Плюс ставят в начало отрезка, минус сразу после конца: d[r + 1] = d[4]. Тогда Prefix Sum даёт +4 на местах 1, 2, 3 и 0 на месте 4. С d[3] -= 4 место 3 не получит прибавку."},
				{ID: "ps3-restore", Type: "short", Prompt: "d = [2, 0, 3, −2, 0, −3], n = 5. Чему равно a[3]?", Answer: "3", Explain: "a[3] = d[0] + d[1] + d[2] + d[3] = 2 + 0 + 3 − 2 = 3. Весь Array a = [2, 2, 5, 3, 3]. Ответ −2 это само d[3], разность соседей, а не значение."},
				{ID: "ps3-count", Type: "short", Prompt: "LeetCode 1109: n = 20 000 рейсов и 20 000 бронирований, каждое на все рейсы. Difference Array делает 2 изменения на бронирование и n сложений в конце. Сколько операций всего?", Answer: "60000", Explain: "2 · 20 000 + 20 000 = 60 000. Цикл по рейсам каждого бронирования сделал бы 20 000 · 20 000 = 400 000 000 прибавлений. Это O(n + b) против O(n · b)."},
				{ID: "ps3-death", Type: "choice", Prompt: "LeetCode 1854. Человек: birth = 1990, death = 1995. Год смерти не считается. Какие изменения d по годам?", Options: []string{"d[1990] += 1 и d[1996] -= 1", "d[1991] += 1 и d[1995] -= 1", "d[1990] += 1 и d[1995] -= 1"}, Answer: 2, Explain: "Человек жив в годы 1990..1994, это отрезок [birth, death − 1]. Минус стоит сразу после конца: d[death] = d[1995]. Вариант с d[1996] посчитал бы и 1995 год. В коде номер года сдвигают: d[year − 1950]."},
				{ID: "ps3-coffee", Type: "short", Prompt: "Codeforces 816B. Рецепты [1, 3], [2, 5], [3, 4], k = 2. Сколько допустимых температур на [1, 5]?", Answer: "3", Explain: "Число рецептов: 1 → 1, 2 → 2, 3 → 3, 4 → 2, 5 → 1. Не меньше 2 у температур 2, 3 и 4: ответ 3. Условие «больше k» дало бы 1. Для q вопросов нужна вторая Prefix Sum по отметкам, тогда вопрос стоит O(1)."},
				{ID: "ps3-size", Type: "code", Prompt: "n = 5, ops = [(2, 4, 7)], номера с 0. Что случится при запуске?", Language: "python", Code: `d = [0] * n
for l, r, v in ops:
    d[l] += v
    d[r + 1] -= v`, Options: []string{"IndexError: d[5] вне списка", "a[4] станет 0", "Код верный"}, Answer: 0, Explain: "При r = n − 1 = 4 код пишет в d[5], а в d есть только места 0..4: IndexError. Нужен d = [0] * (n + 1). Место d[n] не входит в ответ, оно только закрывает отрезки, которые кончаются в последней точке."},
				{ID: "ps3-online", Type: "choice", Prompt: "После каждого из q прибавлений на отрезке нужно сразу печатать a[i] для любого i. Почему одного Difference Array мало?", Options: []string{"Difference Array не умеет отрицательные v", "d портится после чтения", "Каждое чтение a[i] складывает d[0..i]: до n сложений"}, Answer: 2, Explain: "a[i] это Prefix Sum по d. Без готовых префиксов каждое чтение стоит O(n), всего O(n · q). Difference Array хорош, когда сначала идут все прибавления, а потом все чтения. Отрицательные v работают: d[l] += v и d[r + 1] -= v при любом знаке."},
			},
			Links: []sets.Link{
				{Title: "Difference Array", URL: "https://peltorator.org/posts/prefix_sums/#разностный-массив", Source: "peltorator"},
				{Title: "Спортивное программирование: глава о Difference Array", URL: "https://peltorator.ru/cp_book.pdf#page=19", Source: "peltorator, книга"},
				{Title: "Difference Arrays", URL: "https://usaco.guide/silver/more-prefix-sums#difference-arrays", Source: "USACO Guide (англ.)"},
				{Title: "Difference Array", URL: "https://labuladong.online/en/algo/data-structure/diff-array/", Source: "labuladong (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "ps3-marks",
					Title: "Difference Array: n = 6, три прибавления",
					Frames: []sets.FigureFrame{
						{Caption: "n = 6, d длины n + 1 = 7, все нули. Три прибавления: +2 на [1, 3], +5 на [2, 5], +1 на [0, 0].", Rows: []sets.FigureRow{
							{Label: "d", Cells: []any{0, 0, 0, 0, 0, 0, 0}},
							{Label: "a", Cells: []any{"", "", "", "", "", ""}},
						}},
						{Caption: "+2 на [1, 3]: d[1] += 2, d[4] -= 2. Два изменения вместо 3 прибавлений.", Rows: []sets.FigureRow{
							{Label: "d", Cells: []any{0, 2, 0, 0, -2, 0, 0}, Marked: []int{1, 4}},
							{Label: "a", Cells: []any{"", "", "", "", "", ""}, Window: &[2]int{1, 3}},
						}},
						{Caption: "+5 на [2, 5]: d[2] += 5, d[6] -= 5. Два изменения вместо 4 прибавлений.", Rows: []sets.FigureRow{
							{Label: "d", Cells: []any{0, 2, 5, 0, -2, 0, -5}, Marked: []int{2, 6}},
							{Label: "a", Cells: []any{"", "", "", "", "", ""}, Window: &[2]int{2, 5}},
						}},
						{Caption: "+1 на [0, 0]: d[0] += 1, d[1] -= 1. Два изменения вместо одного прибавления.", Rows: []sets.FigureRow{
							{Label: "d", Cells: []any{1, 1, 5, 0, -2, 0, -5}, Marked: []int{0, 1}},
							{Label: "a", Cells: []any{"", "", "", "", "", ""}, Window: &[2]int{0, 0}},
						}},
						{Caption: "Prefix Sum по d даёт a = [1, 2, 7, 7, 5, 5]. Всего 6 изменений d и 6 сложений; прямой цикл сделал бы 3 + 4 + 1 = 8 прибавлений.", Rows: []sets.FigureRow{
							{Label: "d", Cells: []any{1, 1, 5, 0, -2, 0, -5}, Dimmed: []int{6}},
							{Label: "a", Cells: []any{1, 2, 7, 7, 5, 5}, Marked: []int{0, 1, 2, 3, 4, 5}},
						}},
					},
				},
				{
					ID:    "ps3-coffee",
					Title: "Karen and Coffee: k = 2, температуры 90..100",
					Frames: []sets.FigureFrame{
						{Caption: "Рецепты [91, 94], [92, 97], [97, 99]: d[l] += 1, d[r + 1] -= 1. Метка d[100] от [97, 99] видна в последней клетке.", Rows: []sets.FigureRow{
							{Label: "°C", Cells: []any{90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100}},
							{Label: "d", Cells: []any{0, 1, 1, 0, 0, -1, 0, 1, -1, 0, -1}, Marked: []int{1, 2, 5, 7, 8, 10}},
						}},
						{Caption: "Prefix Sum по d: число рецептов для каждой температуры. Пик 2.", Rows: []sets.FigureRow{
							{Label: "°C", Cells: []any{90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100}},
							{Label: "d", Cells: []any{0, 1, 1, 0, 0, -1, 0, 1, -1, 0, -1}},
							{Label: "рецептов", Cells: []any{0, 1, 2, 2, 2, 1, 1, 2, 1, 1, 0}, Marked: []int{2, 3, 4, 7}},
						}},
						{Caption: "k = 2: подходят 92, 93, 94 и 97. Ставим 1 для подходящей температуры.", Rows: []sets.FigureRow{
							{Label: "°C", Cells: []any{90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100}},
							{Label: "рецептов", Cells: []any{0, 1, 2, 2, 2, 1, 1, 2, 1, 1, 0}},
							{Label: "подходит", Cells: []any{0, 0, 1, 1, 1, 0, 0, 1, 0, 0, 0}, Marked: []int{2, 3, 4, 7}},
						}},
						{Caption: "Вторая Prefix Sum по отметкам. Вопрос [93, 97]: c[97] − c[92] = 4 − 1 = 3.", Rows: []sets.FigureRow{
							{Label: "°C", Cells: []any{90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100}},
							{Label: "подходит", Cells: []any{0, 0, 1, 1, 1, 0, 0, 1, 0, 0, 0}, Window: &[2]int{3, 7}},
							{Label: "c", Cells: []any{0, 0, 1, 2, 3, 3, 3, 4, 4, 4, 4}, Marked: []int{2, 7}, Pointers: &sets.Pointers{{Name: "c[a − 1]", Cell: 2}, {Name: "c[b]", Cell: 7}}},
						}},
					},
				},
				{
					ID:    "ps3-greg",
					Title: "Greg and Array: два Difference Array",
					Frames: []sets.FigureFrame{
						{Caption: "Запросы [1, 2], [1, 3], [2, 3] по номерам операций: Difference Array cd длины m + 1 = 4. Prefix Sum даёт cnt = [2, 3, 2].", Rows: []sets.FigureRow{
							{Label: "cd", Cells: []any{2, 1, -1, -2}, Marked: []int{0, 1, 2, 3}},
							{Label: "cnt", Cells: []any{2, 3, 2}, Marked: []int{0, 1, 2}},
						}},
						{Caption: "Операция i прибавляет d · cnt[i]: 1 · 2 = 2 на [1, 2], 2 · 3 = 6 на [1, 3], 4 · 2 = 8 на [2, 3]. Второй Difference Array по номерам элементов.", Rows: []sets.FigureRow{
							{Label: "cnt", Cells: []any{2, 3, 2}},
							{Label: "d · cnt", Cells: []any{2, 6, 8}},
							{Label: "ad", Cells: []any{8, 8, -2, -14}, Marked: []int{0, 1, 2, 3}},
						}},
						{Caption: "Prefix Sum по ad даёт прибавки [8, 16, 14], ответ a = [9, 18, 17]. Всего 2k + 2m = 12 изменений, m + n = 6 сложений Prefix Sum и n = 3 прибавления к a.", Rows: []sets.FigureRow{
							{Label: "ad", Cells: []any{8, 8, -2, -14}, Dimmed: []int{3}},
							{Label: "прибавка", Cells: []any{8, 16, 14}},
							{Label: "a", Cells: []any{9, 18, 17}, Marked: []int{0, 1, 2}},
						}},
					},
				},
			},
		},
		{
			Title:   "Prefix Sum на Matrix",
			Summary: "Двумерная Prefix Sum хранит суммы прямоугольников от левого верхнего угла, и сумма любого прямоугольника получается из четырёх чисел. Вы считаете суммы окон с обрезкой по краям, прибавляете 1 к прямоугольникам через двумерный Difference Array и отвечаете на запросы по таблице размеров.",
			Lesson:  lesson4,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/image-smoother/",
					Number:     number(661),
					Title:      "Image Smoother",
					Difficulty: "easy",
					Topics:     []string{"Array", "Matrix", "Prefix Sum"},
					Note:       "Сумма окна 3 × 3 с обрезкой по краям: из двумерной Prefix Sum или прямым обходом до 9 клеток окна. Число клеток окна (r2 − r1 + 1)(c2 − c1 + 1), ответ sum // count. m × n - размер img. Время: O(mn). Память: O(mn) для P или O(1) при прямом обходе; ответ O(mn). Для клетки (0, 0) второго примера: 550 // 4 = 137.",
					Statement: `<p><strong>Сглаживание изображения</strong> это фильтр размера <code>3 × 3</code>. Его применяют к каждой клетке изображения: новое значение клетки равно среднему значению самой клетки и восьми клеток вокруг неё, округлённому вниз (floor). Это среднее девяти клеток окна 3 × 3. Если одной или нескольких клеток вокруг нет, их не учитывают в среднем: например, в окне у угла изображения есть только 4 клетки.</p><p>Дана Matrix целых чисел <code>img</code> размера <code>m × n</code>: оттенки серого цвета изображения. Верните <em>изображение после сглаживания каждой его клетки</em>.</p><p>Метод класса <code>Solution</code>: <code>def imageSmoother(self, img: list[list[int]]) -&gt; list[list[int]]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> img = [[1,1,1],[1,0,1],[1,1,1]]
<strong>Выход:</strong> [[0,0,0],[0,0,0],[0,0,0]]
<strong>Пояснение:</strong>
Для клеток (0,0), (0,2), (2,0), (2,2): floor(3/4) = floor(0.75) = 0
Для клеток (0,1), (1,0), (1,2), (2,1): floor(5/6) = floor(0.83333333) = 0
Для клетки (1,1): floor(8/9) = floor(0.88888889) = 0</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> img = [[100,200,100],[200,50,200],[100,200,100]]
<strong>Выход:</strong> [[137,141,137],[141,138,141],[137,141,137]]
<strong>Пояснение:</strong>
Для клеток (0,0), (0,2), (2,0), (2,2): floor((100+200+200+50)/4) = floor(137.5) = 137
Для клеток (0,1), (1,0), (1,2), (2,1): floor((200+200+50+200+100+100)/6) = floor(141.666667) = 141
Для клетки (1,1): floor((50+200+200+200+200+100+100+100+100)/9) = floor(138.888889) = 138</pre><h3>Ограничения</h3><ul><li><code>m == img.length</code></li><li><code>n == img[i].length</code></li><li><code>1 ≤ m, n ≤ 200</code></li><li><code>0 ≤ img[i][j] ≤ 255</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "imageSmoother"}, Cases: []sets.TestCase{
						{Input: "[[[1,1,1],[1,0,1],[1,1,1]]]", Expected: text("[[0,0,0],[0,0,0],[0,0,0]]")},
						{Input: "[[[100,200,100],[200,50,200],[100,200,100]]]", Expected: text("[[137,141,137],[141,138,141],[137,141,137]]")},
						{Input: "[[[7]]]", Expected: text("[[7]]")},
						{Input: "[[[1,2,3,4]]]", Expected: text("[[1,2,3,3]]")},
						{Input: "[[[255,0],[0,255]]]", Expected: text("[[127,127],[127,127]]")},
						{Input: "[[[0],[5],[10]]]", Expected: text("[[2],[5],[7]]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/matrix-block-sum/",
					Number:     number(1314),
					Title:      "Matrix Block Sum",
					Difficulty: "medium",
					Topics:     []string{"Array", "Matrix", "Prefix Sum"},
					Note:       "Двумерная Prefix Sum с нулевой строкой и столбцом. Для клетки (i, j): r1 = max(0, i − k), r2 = min(m − 1, i + k), так же для столбцов. m × n - размер mat. Время: O(mn): построение P и по 4 обращения на клетку. Память: O(mn) для P, ответ O(mn). Для первого примера answer[0][0] = 12. Brute Force складывает до min(m, 2k + 1) · min(n, 2k + 1) клеток на клетку ответа: при 100 × 100 и k = 100 это 10^4 · 10^4 = 10^8 сложений.",
					Statement: `<p>Дана Matrix <code>mat</code> размера <code>m × n</code> и целое число <code>k</code>. Верните <em>Matrix</em> <code>answer</code>, <em>где каждое</em> <code>answer[i][j]</code> <em>это сумма всех элементов</em> <code>mat[r][c]</code>, <em>для которых</em>:</p><ul><li><code>i − k ≤ r ≤ i + k</code>,</li><li><code>j − k ≤ c ≤ j + k</code>, и</li><li><code>(r, c)</code> это допустимая позиция в Matrix.</li></ul><p>Метод класса <code>Solution</code>: <code>def matrixBlockSum(self, mat: list[list[int]], k: int) -&gt; list[list[int]]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> mat = [[1,2,3],[4,5,6],[7,8,9]], k = 1
<strong>Выход:</strong> [[12,21,16],[27,45,33],[24,39,28]]</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> mat = [[1,2,3],[4,5,6],[7,8,9]], k = 2
<strong>Выход:</strong> [[45,45,45],[45,45,45],[45,45,45]]</pre><h3>Ограничения</h3><ul><li><code>m == mat.length</code></li><li><code>n == mat[i].length</code></li><li><code>1 ≤ m, n, k ≤ 100</code></li><li><code>1 ≤ mat[i][j] ≤ 100</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "matrixBlockSum"}, Cases: []sets.TestCase{
						{Input: "[[[1,2,3],[4,5,6],[7,8,9]],1]", Expected: text("[[12,21,16],[27,45,33],[24,39,28]]")},
						{Input: "[[[1,2,3],[4,5,6],[7,8,9]],2]", Expected: text("[[45,45,45],[45,45,45],[45,45,45]]")},
						{Input: "[[[5]],100]", Expected: text("[[5]]")},
						{Input: "[[[1,2,3,4]],1]", Expected: text("[[3,6,9,7]]")},
						{Input: "[[[1,1,1,1],[1,2,2,1],[1,1,1,1]],1]", Expected: text("[[5,8,8,5],[7,11,11,7],[5,8,8,5]]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/increment-submatrices-by-one/",
					Number:     number(2536),
					Title:      "Increment Submatrices by One",
					Difficulty: "medium",
					Topics:     []string{"Array", "Matrix", "Prefix Sum"},
					Note:       "Двумерный Difference Array размера (n + 1) × (n + 1): 4 изменения на запрос, затем Prefix Sum по строкам и по столбцам. n - размер Matrix, q - число запросов. Время: O(q + n²). Память: O(n²) для D, ответ n² чисел. Для первого примера ответ [[1, 1, 0], [1, 2, 1], [0, 1, 1]]. Прибавление по всем клеткам каждого запроса стоит O(q · n²) в худшем случае; Difference Array в каждой строке даёт O(q · n + n²).",
					Statement: `<p>Дано целое положительное число <code>n</code>. Сначала есть Matrix целых чисел <code>mat</code> размера <code>n × n</code> с <strong>номерами с 0</strong>, заполненная нулями.</p><p>Также дан двумерный Array целых чисел <code>queries</code>. Для каждого <code>queries[i] = [row1<sub>i</sub>, col1<sub>i</sub>, row2<sub>i</sub>, col2<sub>i</sub>]</code> выполните такую операцию:</p><ul><li>Прибавьте <code>1</code> к <strong>каждому элементу</strong> прямоугольной части Matrix с <strong>левым верхним</strong> углом <code>(row1<sub>i</sub>, col1<sub>i</sub>)</code> и <strong>правым нижним</strong> углом <code>(row2<sub>i</sub>, col2<sub>i</sub>)</code>. То есть прибавьте <code>1</code> к <code>mat[x][y]</code> для всех <code>row1<sub>i</sub> ≤ x ≤ row2<sub>i</sub></code> и <code>col1<sub>i</sub> ≤ y ≤ col2<sub>i</sub></code>.</li></ul><p>Верните <em>Matrix</em> <code>mat</code> <em>после выполнения всех запросов</em>.</p><p>Метод класса <code>Solution</code>: <code>def rangeAddQueries(self, n: int, queries: list[list[int]]) -&gt; list[list[int]]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> n = 3, queries = [[1,1,2,2],[0,0,1,1]]
<strong>Выход:</strong> [[1,1,0],[1,2,1],[0,1,1]]
<strong>Пояснение:</strong> Matrix в начале [[0,0,0],[0,0,0],[0,0,0]], после первого запроса [[0,0,0],[0,1,1],[0,1,1]], после второго запроса [[1,1,0],[1,2,1],[0,1,1]].
- Первый запрос прибавляет 1 к каждому элементу прямоугольной части с левым верхним углом (1, 1) и правым нижним углом (2, 2).
- Второй запрос прибавляет 1 к каждому элементу прямоугольной части с левым верхним углом (0, 0) и правым нижним углом (1, 1).</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> n = 2, queries = [[0,0,1,1]]
<strong>Выход:</strong> [[1,1],[1,1]]
<strong>Пояснение:</strong> Matrix в начале [[0,0],[0,0]], после первого запроса [[1,1],[1,1]].
- Первый запрос прибавляет 1 к каждому элементу Matrix.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 500</code></li><li><code>1 ≤ queries.length ≤ 10<sup>4</sup></code></li><li><code>0 ≤ row1<sub>i</sub> ≤ row2<sub>i</sub> &lt; n</code></li><li><code>0 ≤ col1<sub>i</sub> ≤ col2<sub>i</sub> &lt; n</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "rangeAddQueries"}, Cases: []sets.TestCase{
						{Input: "[3,[[1,1,2,2],[0,0,1,1]]]", Expected: text("[[1,1,0],[1,2,1],[0,1,1]]")},
						{Input: "[2,[[0,0,1,1]]]", Expected: text("[[1,1],[1,1]]")},
						{Input: "[1,[[0,0,0,0],[0,0,0,0]]]", Expected: text("[[2]]")},
						{Input: "[3,[[0,2,2,2],[2,0,2,2]]]", Expected: text("[[0,0,1],[0,0,1],[1,1,2]]")},
						{Input: "[4,[[1,1,2,2],[0,0,3,3],[2,2,2,2]]]", Expected: text("[[1,1,1,1],[1,2,2,1],[1,2,3,1],[1,1,1,1]]")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1722/E",
					Title:      "Counting Rectangles",
					Difficulty: "medium",
					Topics:     []string{"Matrix", "Prefix Sum"},
					Note:       "Codeforces, рейтинг 1600. Таблица A[h][w] хранит сумму площадей прямоугольников h × w, P - двумерная Prefix Sum по A. Строгие неравенства дают строки h_s + 1..h_b − 1 и столбцы w_s + 1..w_b − 1. C = 1000. Время: O(C² + n + q) на набор. Память: O(C²) для P, ответ q чисел. Стройте P один раз на набор, а не на запрос. В первом наборе подходит только 2 × 3, ответ 6. Проверка всех прямоугольников для каждого запроса стоит O(n · q).",
					Statement: `<p>У вас есть <code>n</code> прямоугольников, у <code>i</code>-го прямоугольника высота <code>h<sub>i</sub></code> и ширина <code>w<sub>i</sub></code>.</p><p>Вам задают <code>q</code> запросов вида <code>h<sub>s</sub> w<sub>s</sub> h<sub>b</sub> w<sub>b</sub></code>.</p><p>Для каждого запроса выведите общую площадь ваших прямоугольников, в которые помещается прямоугольник высоты <code>h<sub>s</sub></code> и ширины <code>w<sub>s</sub></code> и которые сами помещаются в прямоугольник высоты <code>h<sub>b</sub></code> и ширины <code>w<sub>b</sub></code>. Другими словами, выведите сумму <code>h<sub>i</sub> · w<sub>i</sub></code> по всем <code>i</code>, для которых <code>h<sub>s</sub> &lt; h<sub>i</sub> &lt; h<sub>b</sub></code> и <code>w<sub>s</sub> &lt; w<sub>i</sub> &lt; w<sub>b</sub></code>.</p><p>Обратите внимание: если у двух прямоугольников одинаковая высота или одинаковая ширина, то ни один из них не помещается в другой. Также прямоугольники нельзя поворачивать.</p><p>Обратите внимание: ответ в некоторых тестах не помещается в 32-битный целый тип, поэтому нужен хотя бы 64-битный целый тип вашего языка программирования (например, <code>long long</code> в C++).</p><h3>Входные данные</h3><p>В первой строке записано целое число <code>t</code> (<code>1 ≤ t ≤ 100</code>): число наборов входных данных.</p><p>В первой строке каждого набора записаны два целых числа <code>n, q</code> (<code>1 ≤ n ≤ 10<sup>5</sup></code>; <code>1 ≤ q ≤ 10<sup>5</sup></code>): число ваших прямоугольников и число запросов.</p><p>Затем идут <code>n</code> строк, в каждой два целых числа <code>h<sub>i</sub>, w<sub>i</sub></code> (<code>1 ≤ h<sub>i</sub>, w<sub>i</sub> ≤ 1000</code>): высота и ширина <code>i</code>-го прямоугольника.</p><p>Затем идут <code>q</code> строк, в каждой четыре целых числа <code>h<sub>s</sub>, w<sub>s</sub>, h<sub>b</sub>, w<sub>b</sub></code> (<code>1 ≤ h<sub>s</sub> &lt; h<sub>b</sub> ≤ 1000</code>, <code>1 ≤ w<sub>s</sub> &lt; w<sub>b</sub> ≤ 1000</code>): описание запроса.</p><p>Сумма <code>q</code> по всем наборам не больше <code>10<sup>5</sup></code>, и сумма <code>n</code> по всем наборам не больше <code>10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого набора выведите <code>q</code> строк: в <code>i</code>-й строке ответ на <code>i</code>-й запрос.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>3
2 1
2 3
3 2
1 1 3 4
5 5
1 1
2 2
3 3
4 4
5 5
3 3 6 6
2 1 4 5
1 1 2 10
1 1 100 100
1 1 3 3
3 1
999 999
999 999
999 998
1 1 1000 1000</pre><p><strong>Выходные данные:</strong></p><pre>6
41
9
0
54
4
2993004</pre><h3>Примечание</h3><p>В первом наборе только один запрос. Нужно найти сумму площадей всех прямоугольников, в которые помещается прямоугольник <code>1 × 1</code> и которые сами помещаются в прямоугольник <code>3 × 4</code>.</p><p>Подходит только прямоугольник <code>2 × 3</code>: <code>1 &lt; 2</code> (сравниваем высоты) и <code>1 &lt; 3</code> (сравниваем ширины), поэтому прямоугольник <code>1 × 1</code> помещается в него; <code>2 &lt; 3</code> (сравниваем высоты) и <code>3 &lt; 4</code> (сравниваем ширины), поэтому он помещается в прямоугольник <code>3 × 4</code>. Прямоугольник <code>3 × 2</code> слишком высокий для прямоугольника <code>3 × 4</code>. Общая площадь <code>2 · 3 = 6</code>.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 100</code></li><li><code>1 ≤ n ≤ 10<sup>5</sup></code>, <code>1 ≤ q ≤ 10<sup>5</sup></code></li><li><code>1 ≤ h<sub>i</sub>, w<sub>i</sub> ≤ 1000</code></li><li><code>1 ≤ h<sub>s</sub> &lt; h<sub>b</sub> ≤ 1000</code>, <code>1 ≤ w<sub>s</sub> &lt; w<sub>b</sub> ≤ 1000</code></li><li>сумма <code>n</code> по всем наборам не больше <code>10<sup>5</sup></code>, сумма <code>q</code> по всем наборам не больше <code>10<sup>5</sup></code></li><li>Ограничение по времени: 6 секунд на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `3
2 1
2 3
3 2
1 1 3 4
5 5
1 1
2 2
3 3
4 4
5 5
3 3 6 6
2 1 4 5
1 1 2 10
1 1 100 100
1 1 3 3
3 1
999 999
999 999
999 998
1 1 1000 1000
`, Expected: text(`6
41
9
0
54
4
2993004
`)},
						{Input: `1
1 2
5 5
4 4 6 6
5 5 6 6
`, Expected: text(`25
0
`)},
						{Input: `1
2 1
999 999
1000 1000
1 1 1000 1000
`, Expected: text(`998001
`)},
						{Input: `1
3 1
2 3
2 3
2 4
1 2 3 4
`, Expected: text(`12
`)},
						{Input: `2
1 1
10 10
9 9 11 11
1 1
20 20
9 9 11 11
`, Expected: text(`100
0
`)},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "ps4-corners", Type: "choice", Prompt: "P[i][j] хранит сумму строк 0..i − 1 и столбцов 0..j − 1. Какая формула даёт сумму строк r1..r2 и столбцов c1..c2?", Options: []string{"P[r2 + 1][c2 + 1] - P[r1][c2 + 1] - P[r2 + 1][c1]", "P[r2 + 1][c2 + 1] - P[r1][c2 + 1] - P[r2 + 1][c1] + P[r1][c1]", "P[r2][c2] - P[r1][c2] - P[r2][c1] + P[r1][c1]"}, Answer: 1, Explain: "Из большого угла вычитают полосу сверху и полосу слева. Их общий угол вычтен два раза, поэтому + P[r1][c1]. Для Matrix 3 × 4, строк 1..2 и столбцов 1..2: 54 − 6 − 15 + 1 = 34. Без + P[r1][c1] получится 33. Формула без + 1 у индексов теряет строку r2 и столбец c2."},
				{ID: "ps4-value", Type: "short", Prompt: "mat = [[1, 2, 3, 4], [5, 6, 7, 8], [9, 10, 11, 12]], у P есть нулевая строка и нулевой столбец. Чему равно P[2][3]?", Answer: "24", Explain: "P[2][3] это строки 0..1 и столбцы 0..2: 1 + 2 + 3 + 5 + 6 + 7 = 24. По формуле: mat[1][2] + P[1][3] + P[2][2] − P[1][2] = 7 + 6 + 14 − 3 = 24. Ответ 36 это P[2][4], все четыре столбца."},
				{ID: "ps4-rows", Type: "code", Prompt: "Что напечатает этот код?", Language: "python", Code: `P = [[0] * 3] * 3
P[1][1] = 5
print(P[0][1])`, Options: []string{"5", "0", "IndexError"}, Answer: 0, Explain: "[[0] * 3] * 3 создаёт три ссылки на один и тот же список. Запись в P[1][1] видна во всех строках, поэтому P[0][1] = 5. Каждую строку создают отдельно: [[0] * (n + 1) for _ in range(m + 1)]. Тогда P занимает (m + 1)(n + 1) чисел, как задумано."},
				{ID: "ps4-clip", Type: "short", Prompt: "LeetCode 1314. mat = [[2, 1, 1], [1, 3, 1], [1, 1, 1]], k = 1. Чему равно answer[0][2]?", Answer: "6", Explain: "Окно строк −1..1 и столбцов 1..3 обрезаем до строк 0..1 и столбцов 1..2: 1 + 1 + 3 + 1 = 6. Без обрезки индекс −1 в Python берёт последнюю строку, а индекс 3 даёт IndexError."},
				{ID: "ps4-smoother", Type: "short", Prompt: "LeetCode 661. img = [[1, 2], [3, 4]]. Какое число будет в каждой клетке после сглаживания?", Answer: "2", Explain: "Окно 3 × 3 каждой клетки после обрезки покрывает все 4 клетки: сумма 10, клеток 4, 10 // 4 = 2. Деление на 9 дало бы 1: клетки за краем в среднее не входят."},
				{ID: "ps4-diff", Type: "choice", Prompt: "Двумерный Difference Array D размера 4 × 4 для Matrix 3 × 3. Нужно прибавить 1 на строках 0..1 и столбцах 0..1. Какие четыре изменения верны?", Options: []string{"D[0][0] += 1, D[0][1] -= 1, D[1][0] -= 1, D[1][1] += 1", "D[0][0] += 1, D[0][2] -= 1, D[2][0] -= 1, D[2][2] -= 1", "D[0][0] += 1, D[0][2] -= 1, D[2][0] -= 1, D[2][2] += 1"}, Answer: 2, Explain: "Минусы стоят сразу после последнего столбца, D[0][2], и сразу после последней строки, D[2][0]. Клетка (2, 2) лежит ниже и правее прямоугольника: в её сумму входят + 1 из D[0][0] и оба минуса, поэтому + 1 в D[2][2] возвращает её к 0. С D[2][2] -= 1 клетка (2, 2) получит −2. Первый вариант ставит минусы внутрь прямоугольника: 1 получит только клетка (0, 0)."},
				{ID: "ps4-strict", Type: "choice", Prompt: "Codeforces 1722E. h_s = 2, h_b = 6. Какие высоты h входят в ответ?", Options: []string{"2..6", "3..5", "3..6"}, Answer: 1, Explain: "Условие строгое: 2 < h < 6, поэтому h = 3, 4, 5. В P это строки h_s + 1 = 3 .. h_b − 1 = 5. Отрезок 2..6 берёт прямоугольники, которые не помещаются, а 3..6 берёт h = 6, равное h_b."},
				{ID: "ps4-cost", Type: "short", Prompt: "LeetCode 1314: Matrix 100 × 100, k = 100. Brute Force складывает все клетки окна для каждой клетки ответа, и окно всегда вся Matrix. Сколько сложений всего?", Answer: "100000000", Explain: "10 000 клеток ответа, в каждом окне 10 000 клеток: 10^8 сложений. Двумерная Prefix Sum делает около 3 · 10 000 операций построения и по 3 на клетку ответа, O(mn)."},
			},
			Links: []sets.Link{
				{Title: "Prefix Sum на Matrix", URL: "https://brestprog.by/topics/prefixsums/#префиксные-суммы-в-матрицах", Source: "brestprog"},
				{Title: "Двумерная Prefix Sum на Python", URL: "https://habr.com/ru/articles/901190/", Source: "Хабр"},
				{Title: "Двумерный Difference Array", URL: "https://peltorator.org/posts/prefix_sums/#двумерный-разностный-массив", Source: "peltorator"},
				{Title: "2D Prefix Sums", URL: "https://usaco.guide/silver/more-prefix-sums#2d-prefix-sums", Source: "USACO Guide (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "ps4-table",
					Title: "Двумерная Prefix Sum: Matrix 3 × 4",
					Frames: []sets.FigureFrame{
						{Caption: "Строка 1 P готова: P[1][j + 1] = mat[0][j] + сверху + слева − угол. Пример: P[1][2] = 2 + 0 + 1 − 0 = 3.", Rows: []sets.FigureRow{
							{Label: "mat 0", Cells: []any{1, 2, 3, 4}, Marked: []int{0, 1, 2, 3}},
							{Label: "mat 1", Cells: []any{5, 6, 7, 8}},
							{Label: "mat 2", Cells: []any{9, 10, 11, 12}},
							{Label: "P 0", Cells: []any{0, 0, 0, 0, 0}},
							{Label: "P 1", Cells: []any{0, 1, 3, 6, 10}, Marked: []int{1, 2, 3, 4}},
							{Label: "P 2", Cells: []any{"", "", "", "", ""}},
							{Label: "P 3", Cells: []any{"", "", "", "", ""}},
						}},
						{Caption: "Строка 2 P готова: P[2][j + 1] = mat[1][j] + сверху + слева − угол. Пример: P[2][2] = 6 + 3 + 6 − 1 = 14.", Rows: []sets.FigureRow{
							{Label: "mat 0", Cells: []any{1, 2, 3, 4}},
							{Label: "mat 1", Cells: []any{5, 6, 7, 8}, Marked: []int{0, 1, 2, 3}},
							{Label: "mat 2", Cells: []any{9, 10, 11, 12}},
							{Label: "P 0", Cells: []any{0, 0, 0, 0, 0}},
							{Label: "P 1", Cells: []any{0, 1, 3, 6, 10}},
							{Label: "P 2", Cells: []any{0, 6, 14, 24, 36}, Marked: []int{1, 2, 3, 4}},
							{Label: "P 3", Cells: []any{"", "", "", "", ""}},
						}},
						{Caption: "Строка 3 P готова: P[3][j + 1] = mat[2][j] + сверху + слева − угол. Пример: P[3][2] = 10 + 14 + 15 − 6 = 33.", Rows: []sets.FigureRow{
							{Label: "mat 0", Cells: []any{1, 2, 3, 4}},
							{Label: "mat 1", Cells: []any{5, 6, 7, 8}},
							{Label: "mat 2", Cells: []any{9, 10, 11, 12}, Marked: []int{0, 1, 2, 3}},
							{Label: "P 0", Cells: []any{0, 0, 0, 0, 0}},
							{Label: "P 1", Cells: []any{0, 1, 3, 6, 10}},
							{Label: "P 2", Cells: []any{0, 6, 14, 24, 36}},
							{Label: "P 3", Cells: []any{0, 15, 33, 54, 78}, Marked: []int{1, 2, 3, 4}},
						}},
					},
				},
				{
					ID:    "ps4-query",
					Title: "Сумма прямоугольника: строки 1..2, столбцы 1..2",
					Frames: []sets.FigureFrame{
						{Caption: "Сумма строк 1..2 и столбцов 1..2: P[3][3] − P[1][3] − P[3][1] + P[1][1] = 54 − 6 − 15 + 1 = 34. Четыре обращения к P на любой прямоугольник.", Rows: []sets.FigureRow{
							{Label: "mat 0", Cells: []any{1, 2, 3, 4}},
							{Label: "mat 1", Cells: []any{5, 6, 7, 8}, Marked: []int{1, 2}, Window: &[2]int{1, 2}},
							{Label: "mat 2", Cells: []any{9, 10, 11, 12}, Marked: []int{1, 2}, Window: &[2]int{1, 2}},
							{Label: "P 0", Cells: []any{0, 0, 0, 0, 0}},
							{Label: "P 1", Cells: []any{0, 1, 3, 6, 10}, Marked: []int{1, 3}},
							{Label: "P 2", Cells: []any{0, 6, 14, 24, 36}},
							{Label: "P 3", Cells: []any{0, 15, 33, 54, 78}, Marked: []int{1, 3}},
						}},
						{Caption: "P[1][3] вычел строку 0, столбцы 0..2. P[3][1] вычел столбец 0, строки 0..2. Клетка (0, 0) вычтена два раза, поэтому + P[1][1] = 1.", Rows: []sets.FigureRow{
							{Label: "mat 0", Cells: []any{1, 2, 3, 4}, Marked: []int{0}, Dimmed: []int{1, 2}},
							{Label: "mat 1", Cells: []any{5, 6, 7, 8}, Dimmed: []int{0}},
							{Label: "mat 2", Cells: []any{9, 10, 11, 12}, Dimmed: []int{0}},
						}},
					},
				},
				{
					ID:    "ps4-clip",
					Title: "Matrix Block Sum: обрезка окна по краям",
					Frames: []sets.FigureFrame{
						{Caption: "Клетка (0, 0), k = 1: строки 0..1, столбцы 0..1 после обрезки. Сумма 12, клеток 4.", Rows: []sets.FigureRow{
							{Label: "mat 0", Cells: []any{1, 2, 3}, Marked: []int{0}, Window: &[2]int{0, 1}},
							{Label: "mat 1", Cells: []any{4, 5, 6}, Window: &[2]int{0, 1}},
							{Label: "mat 2", Cells: []any{7, 8, 9}},
						}},
						{Caption: "Клетка (1, 1), k = 1: строки 0..2, столбцы 0..2 после обрезки. Сумма 45, клеток 9.", Rows: []sets.FigureRow{
							{Label: "mat 0", Cells: []any{1, 2, 3}, Window: &[2]int{0, 2}},
							{Label: "mat 1", Cells: []any{4, 5, 6}, Marked: []int{1}, Window: &[2]int{0, 2}},
							{Label: "mat 2", Cells: []any{7, 8, 9}, Window: &[2]int{0, 2}},
						}},
						{Caption: "Клетка (2, 2), k = 1: строки 1..2, столбцы 1..2 после обрезки. Сумма 28, клеток 4. Через P каждая клетка стоит 4 обращения, а не до 9 сложений.", Rows: []sets.FigureRow{
							{Label: "mat 0", Cells: []any{1, 2, 3}},
							{Label: "mat 1", Cells: []any{4, 5, 6}, Window: &[2]int{1, 2}},
							{Label: "mat 2", Cells: []any{7, 8, 9}, Marked: []int{2}, Window: &[2]int{1, 2}},
						}},
					},
				},
				{
					ID:    "ps4-diff",
					Title: "Двумерный Difference Array: n = 3, два запроса",
					Frames: []sets.FigureFrame{
						{Caption: "+1 на строках 1..2, столбцах 1..2: четыре изменения D, + в углах (1, 1) и (3, 3), − в двух других.", Rows: []sets.FigureRow{
							{Label: "D 0", Cells: []any{0, 0, 0, 0}},
							{Label: "D 1", Cells: []any{0, 1, 0, -1}, Marked: []int{1, 3}},
							{Label: "D 2", Cells: []any{0, 0, 0, 0}},
							{Label: "D 3", Cells: []any{0, -1, 0, 1}, Marked: []int{1, 3}},
						}},
						{Caption: "+1 на строках 0..1, столбцах 0..1: четыре изменения D, + в углах (0, 0) и (2, 2), − в двух других.", Rows: []sets.FigureRow{
							{Label: "D 0", Cells: []any{1, 0, -1, 0}, Marked: []int{0, 2}},
							{Label: "D 1", Cells: []any{0, 1, 0, -1}},
							{Label: "D 2", Cells: []any{-1, 0, 1, 0}, Marked: []int{0, 2}},
							{Label: "D 3", Cells: []any{0, -1, 0, 1}},
						}},
						{Caption: "Prefix Sum в строках 0..2, столбцы 0..2. Строка 3 нужна только для отметок.", Rows: []sets.FigureRow{
							{Label: "D 0", Cells: []any{1, 1, 0, 0}, Marked: []int{0, 1, 2}, Dimmed: []int{3}},
							{Label: "D 1", Cells: []any{0, 1, 1, -1}, Marked: []int{0, 1, 2}, Dimmed: []int{3}},
							{Label: "D 2", Cells: []any{-1, -1, 0, 0}, Marked: []int{0, 1, 2}, Dimmed: []int{3}},
							{Label: "D 3", Cells: []any{0, -1, 0, 1}, Dimmed: []int{3}},
						}},
						{Caption: "Prefix Sum в столбцах 0..2, строки 0..2. Клетки 0..2 дают ответ [[1, 1, 0], [1, 2, 1], [0, 1, 1]]. Работа: 4 изменения на запрос и 2n(n − 1) = 12 сложений.", Rows: []sets.FigureRow{
							{Label: "D 0", Cells: []any{1, 1, 0, 0}, Marked: []int{0, 1, 2}, Dimmed: []int{3}},
							{Label: "D 1", Cells: []any{1, 2, 1, -1}, Marked: []int{0, 1, 2}, Dimmed: []int{3}},
							{Label: "D 2", Cells: []any{0, 1, 1, 0}, Marked: []int{0, 1, 2}, Dimmed: []int{3}},
							{Label: "D 3", Cells: []any{0, -1, 0, 1}, Dimmed: []int{0, 1, 2, 3}},
						}},
					},
				},
			},
		},
	},
}
