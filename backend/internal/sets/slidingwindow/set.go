// Package slidingwindow is the learning set "Sliding Window и Deque".
package slidingwindow

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
	Slug:    "sliding-window",
	Title:   "Sliding Window и Deque",
	Summary: "Sliding Window проходит по Array один раз: правая граница добавляет элементы, левая убирает их. Вы обновляете сумму, счётчики и Hash Table окна, считаете окна без перебора и применяете Monotonic Deque вместе с Prefix Sum и Dynamic Programming. Для каждого решения считаете время и дополнительную память.",
	Topics:  []string{"Sliding Window", "Hash Table", "Counting", "Prefix Sum", "Monotonic Queue", "Dynamic Programming"},
	Chapters: []sets.Chapter{
		{
			Title:   "Sliding Window фиксированной длины",
			Summary: "Окно длины k сдвигается на один шаг: один элемент входит, один выходит. Сумму, счётчики, set и Hash Table окна вы обновляете за O(1), а не считаете заново за O(k). Отдельно считаете копии срезов и размер ответа.",
			Lesson:  lesson1,
			Tasks: []sets.Task{
				{
					URL:        "https://codeforces.com/problemset/problem/1690/D",
					Title:      "Black and White Stripe",
					Difficulty: "easy",
					Topics:     []string{"String", "Sliding Window"},
					Note:       "Codeforces, рейтинг 1000. Sliding Window длины k хранит число белых клеток окна. n - длина полосы одного набора. Время: O(n) на набор, O(Σn) на весь вход: k чтений для первого окна и по 2 на каждый сдвиг. Память: O(1) дополнительно без среза; срез s[:k] добавит O(k). Для BBWBW и k = 3 окна содержат 1, 1 и 2 белые клетки, ответ 1. Brute Force считает каждое окно заново за O(nk).",
					Statement: `<p>Дана полоска клетчатой бумаги длины <code>n</code>. Каждая клетка белая или чёрная.</p><p>Какое наименьшее число клеток нужно перекрасить из белого цвета в чёрный, чтобы на полоске был отрезок из <code>k</code> чёрных клеток подряд?</p><p>Если во входных данных отрезок из <code>k</code> чёрных клеток подряд уже есть, выведите 0.</p><h3>Входные данные</h3><p>В первой строке записано целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>): количество наборов входных данных.</p><p>Дальше идут <code>t</code> наборов. В первой строке набора записаны два целых числа <code>n</code> и <code>k</code> (<code>1 ≤ k ≤ n ≤ 2 · 10<sup>5</sup></code>). Во второй строке записаны буквы 'W' (белая клетка) и 'B' (чёрная клетка). Длина этой строки равна <code>n</code>.</p><p>Гарантируется, что сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого из <code>t</code> наборов выведите одно целое число: наименьшее число клеток, которые нужно перекрасить из белого цвета в чёрный, чтобы на полоске был отрезок из <code>k</code> чёрных клеток подряд.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>4
5 3
BBWBW
5 5
BBWBW
5 1
BBWBW
1 1
W</pre><p><strong>Выходные данные:</strong></p><pre>1
2
0
1</pre><h3>Примечание</h3><p>В первом наборе <code>s = "BBWBW"</code> и <code>k = 3</code>. Достаточно перекрасить <code>s<sub>3</sub></code> и получить <code>s = "BBBBW"</code>. В ней есть отрезок длины <code>k = 3</code> из букв 'B'.</p><p>Во втором наборе <code>s = "BBWBW"</code> и <code>k = 5</code>. Достаточно перекрасить <code>s<sub>3</sub></code> и <code>s<sub>5</sub></code> и получить <code>s = "BBBBB"</code>. В ней есть отрезок длины <code>k = 5</code> из букв 'B'.</p><p>В третьем наборе <code>s = "BBWBW"</code> и <code>k = 1</code>. В <code>s</code> уже есть отрезок длины <code>k = 1</code> из букв 'B'.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>1 ≤ k ≤ n ≤ 2 · 10<sup>5</sup></code></li><li>строка набора состоит из <code>n</code> букв 'W' и 'B'</li><li>сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "4\n5 3\nBBWBW\n5 5\nBBWBW\n5 1\nBBWBW\n1 1\nW\n", Expected: text("1\n2\n0\n1\n")},
						{Input: "1\n1 1\nB\n", Expected: text("0\n")},
						{Input: "1\n5 5\nWWWWW\n", Expected: text("5\n")},
						{Input: "1\n6 2\nWWWWBB\n", Expected: text("0\n")},
						{Input: "2\n7 3\nBWBWBWB\n7 4\nWBBWBBW\n", Expected: text("1\n1\n")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/contains-duplicate-ii/",
					Number:     number(219),
					Title:      "Contains Duplicate II",
					Difficulty: "easy",
					Topics:     []string{"Array", "Hash Table", "Sliding Window"},
					Note:       "set хранит последние k чисел. n = len(nums). Время: в среднем O(n): n проверок, n добавлений, не больше n удалений. Память: O(min(n, k + 1)) дополнительно, ответ O(1). Для [1, 2, 3, 1, 2, 3] и k = 2 равные числа стоят на расстоянии 3, ответ False. Brute Force сравнивает пары на расстоянии до k: O(nk).",
					Statement: `<p>Дан Array целых чисел <code>nums</code> и целое число <code>k</code>. Верните <code>True</code>, если в Array есть два <strong>разных индекса</strong> <code>i</code> и <code>j</code>, для которых <code>nums[i] == nums[j]</code> и <code>abs(i - j) ≤ k</code>. Иначе верните <code>False</code>.</p><p>Метод класса <code>Solution</code>: <code>def containsNearbyDuplicate(self, nums: list[int], k: int) -&gt; bool</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,2,3,1], k = 3
<strong>Выход:</strong> True</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [1,0,1,1], k = 1
<strong>Выход:</strong> True</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> nums = [1,2,3,1,2,3], k = 2
<strong>Выход:</strong> False</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 10<sup>5</sup></code></li><li><code>−10<sup>9</sup> ≤ nums[i] ≤ 10<sup>9</sup></code></li><li><code>0 ≤ k ≤ 10<sup>5</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "containsNearbyDuplicate"}, Cases: []sets.TestCase{
						{Input: "[[1,2,3,1],3]", Expected: text("True")},
						{Input: "[[1,0,1,1],1]", Expected: text("True")},
						{Input: "[[1,2,3,1,2,3],2]", Expected: text("False")},
						{Input: "[[1,1],0]", Expected: text("False")},
						{Input: "[[1,2,1],1]", Expected: text("False")},
						{Input: "[[-1000000000,7,-1000000000],2]", Expected: text("True")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/find-all-anagrams-in-a-string/",
					Number:     number(438),
					Title:      "Find All Anagrams in a String",
					Difficulty: "medium",
					Topics:     []string{"Hash Table", "String", "Sliding Window"},
					Note:       "Два list по 26 счётчиков: для p и для окна длины m = len(p). n = len(s). Время: O(m + 26n): при сдвиге два обновления счётчиков, затем сравнение 26 чисел. Память: O(1) дополнительно для алфавита из 26 букв, ответ до n − m + 1 номеров. Для cbaebabacd и abc подходят окна 0 и 6. Sorting каждого окна даёт O(n · m log m).",
					Statement: `<p>Даны две String <code>s</code> и <code>p</code>. Верните Array всех начальных индексов тех подстрок <code>s</code>, которые являются Anagram строки <code>p</code>. Ответ можно вернуть в <strong>любом порядке</strong>.</p><p>Anagram строки - строка из тех же букв, возможно в другом порядке: каждая буква взята ровно столько раз, сколько в исходной строке.</p><p>Метод класса <code>Solution</code>: <code>def findAnagrams(self, s: str, p: str) -&gt; list[int]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> s = "cbaebabacd", p = "abc"
<strong>Выход:</strong> [0,6]
<strong>Пояснение:</strong>
Подстрока с начальным индексом 0 - "cba", это Anagram строки "abc".
Подстрока с начальным индексом 6 - "bac", это Anagram строки "abc".</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> s = "abab", p = "ab"
<strong>Выход:</strong> [0,1,2]
<strong>Пояснение:</strong>
Подстрока с начальным индексом 0 - "ab", это Anagram строки "ab".
Подстрока с начальным индексом 1 - "ba", это Anagram строки "ab".
Подстрока с начальным индексом 2 - "ab", это Anagram строки "ab".</pre><h3>Ограничения</h3><ul><li><code>1 ≤ s.length, p.length ≤ 3 · 10<sup>4</sup></code></li><li><code>s</code> и <code>p</code> состоят только из строчных английских букв</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "findAnagrams", Order: "any"}, Cases: []sets.TestCase{
						{Input: `["cbaebabacd","abc"]`, Expected: text("[0,6]")},
						{Input: `["abab","ab"]`, Expected: text("[0,1,2]")},
						{Input: `["a","ab"]`, Expected: text("[]")},
						{Input: `["abc","cba"]`, Expected: text("[0]")},
						{Input: `["aaaa","aa"]`, Expected: text("[0,1,2]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/maximum-sum-of-distinct-subarrays-with-length-k/",
					Number:     number(2461),
					Title:      "Maximum Sum of Distinct Subarrays With Length K",
					Difficulty: "medium",
					Topics:     []string{"Array", "Hash Table", "Sliding Window"},
					Note:       "dict счётчиков и сумма окна длины k. n = len(nums). Окно подходит, когда len(cnt) == k. Время: в среднем O(n). Память: O(k) ключей дополнительно, ответ O(1). Ключ со счётчиком 0 нужно удалить. Для [1, 5, 4, 2, 9, 9, 9] и k = 3 подходят суммы 10, 11 и 15, ответ 15. Новый set для каждого окна даёт O(nk).",
					Statement: `<p>Дан Array целых чисел <code>nums</code> и целое число <code>k</code>. Найдите наибольшую сумму Subarray среди всех Subarray из <code>nums</code>, которые выполняют условия:</p><ul><li>длина Subarray равна <code>k</code>;</li><li>все элементы Subarray <strong>разные</strong>.</li></ul><p>Верните <em>наибольшую сумму среди всех Subarray, которые выполняют условия</em>. Если ни один Subarray не выполняет условия, верните <code>0</code>.</p><p><em>Subarray - непустая последовательность соседних элементов Array.</em></p><p>Метод класса <code>Solution</code>: <code>def maximumSubarraySum(self, nums: list[int], k: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,5,4,2,9,9,9], k = 3
<strong>Выход:</strong> 15
<strong>Пояснение:</strong> Subarray длины 3 в nums:
- [1,5,4] выполняет условия, его сумма 10.
- [5,4,2] выполняет условия, его сумма 11.
- [4,2,9] выполняет условия, его сумма 15.
- [2,9,9] не выполняет условия: элемент 9 повторяется.
- [9,9,9] не выполняет условия: элемент 9 повторяется.
Верните 15: это наибольшая сумма среди Subarray, которые выполняют условия.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [4,4,4], k = 3
<strong>Выход:</strong> 0
<strong>Пояснение:</strong> Subarray длины 3 в nums:
- [4,4,4] не выполняет условия: элемент 4 повторяется.
Верните 0: ни один Subarray не выполняет условия.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ k ≤ nums.length ≤ 10<sup>5</sup></code></li><li><code>1 ≤ nums[i] ≤ 10<sup>5</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "maximumSubarraySum"}, Cases: []sets.TestCase{
						{Input: "[[1,5,4,2,9,9,9],3]", Expected: text("15")},
						{Input: "[[4,4,4],3]", Expected: text("0")},
						{Input: "[[3,1,2],1]", Expected: text("3")},
						{Input: "[[9,9,1,2],2]", Expected: text("10")},
						{Input: "[[1,2,2,3,4],2]", Expected: text("7")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/961/B",
					Title:      "Lecture Sleep",
					Difficulty: "medium",
					Topics:     []string{"Array", "Sliding Window"},
					Note:       "Codeforces, рейтинг 1200. Основа - теоремы минут с t = 1, Sliding Window длины k ищет наибольший выигрыш минут сна. Время: O(n). Память: O(1) дополнительно, вход O(n). Для примера основа 6, выигрыши окон 5, 5, 10, 9, ответ 16. Brute Force пересчитывает каждое окно: O(nk).",
					Statement: `<p>Вы и ваш друг Миша слушаете лекцию по математическому анализу. Лекция длится <code>n</code> минут. В <code>i</code>-ю минуту лектор рассказывает <code>a<sub>i</sub></code> теорем.</p><p>Мише очень интересен математический анализ, но ему трудно не спать всю лекцию. Дан Array <code>t</code> поведения Миши. Если в <code>i</code>-ю минуту лекции Миша спит, то <code>t<sub>i</sub></code> равно 0, иначе <code>t<sub>i</sub></code> равно 1. Когда Миша не спит, он записывает все теоремы, которые ему рассказывают: <code>a<sub>i</sub></code> теорем в <code>i</code>-ю минуту. Иначе он ничего не записывает.</p><p>Вы знаете секретный способ, который не даёт Мише спать <code>k</code> минут подряд. Но использовать его можно только один раз. Начать можно в начале любой минуты от 1 до <code>n − k + 1</code>. Если начать в минуту <code>i</code>, то Миша не спит в минуты <code>j</code> от <code>i</code> до <code>i + k − 1</code> и записывает все теоремы, которые рассказывает лектор.</p><p>Найдите наибольшее число теорем, которые Миша сможет записать, если вы один раз используете секретный способ, чтобы его разбудить.</p><h3>Входные данные</h3><p>В первой строке записаны два целых числа <code>n</code> и <code>k</code> (<code>1 ≤ k ≤ n ≤ 10<sup>5</sup></code>): длина лекции в минутах и число минут, которые Миша может не спать благодаря способу.</p><p>Во второй строке записаны <code>n</code> целых чисел <code>a<sub>1</sub>, a<sub>2</sub>, ..., a<sub>n</sub></code> (<code>1 ≤ a<sub>i</sub> ≤ 10<sup>4</sup></code>): число теорем, которые лектор рассказывает в <code>i</code>-ю минуту.</p><p>В третьей строке записаны <code>n</code> целых чисел <code>t<sub>1</sub>, t<sub>2</sub>, ..., t<sub>n</sub></code> (<code>0 ≤ t<sub>i</sub> ≤ 1</code>): поведение Миши в <code>i</code>-ю минуту лекции.</p><h3>Выходные данные</h3><p>Выведите одно целое число: наибольшее число теорем, которые Миша сможет записать, если вы один раз используете секретный способ, чтобы его разбудить.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>6 3
1 3 5 2 5 4
1 1 0 1 0 0</pre><p><strong>Выходные данные:</strong></p><pre>16</pre><h3>Примечание</h3><p>В примере лучше всего использовать секретный способ в начале третьей минуты. Тогда Миша запишет 16 теорем.</p><h3>Ограничения</h3><ul><li><code>1 ≤ k ≤ n ≤ 10<sup>5</sup></code></li><li><code>1 ≤ a<sub>i</sub> ≤ 10<sup>4</sup></code></li><li><code>0 ≤ t<sub>i</sub> ≤ 1</code></li><li>Ограничение по времени: 1 секунда на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "6 3\n1 3 5 2 5 4\n1 1 0 1 0 0\n", Expected: text("16\n")},
						{Input: "1 1\n5\n0\n", Expected: text("5\n")},
						{Input: "3 2\n1 2 3\n1 1 1\n", Expected: text("6\n")},
						{Input: "5 2\n5 1 1 4 4\n1 0 1 0 0\n", Expected: text("14\n")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "sw1-update", Type: "short", Prompt: "a = [4, 2, 7, 1, 6, 3], k = 3. Сумма первого окна [4, 2, 7] равна 13. Чему равна сумма окна после двух сдвигов вправо?", Answer: "14", Explain: "Первый сдвиг: вышло 4, вошло 1, 13 − 4 + 1 = 10. Второй сдвиг: вышло 2, вошло 6, 10 − 2 + 6 = 14. Это окно [7, 1, 6]. Ответ 10 был бы после одного сдвига. Два обновления стоят 4 действия вместо 6 чтений двух новых окон."},
				{ID: "sw1-cost", Type: "choice", Prompt: "n = 100 000, k = 30 000. Сколько чтений элементов делает Sliding Window, чтобы найти суммы всех окон?", Options: []string{"около 2 100 000 000", "около 30 000", "около 170 000"}, Answer: 2, Explain: "Первое окно читает k = 30 000 чисел. Каждый из n − k = 70 000 сдвигов читает 2 числа, это 140 000. Всего 170 000. Около 2,1 миллиарда чтений делает Brute Force: 70 001 окно по 30 000 чисел. 30 000 - это только первое окно."},
				{ID: "sw1-out-index", Type: "code", Prompt: "Окно длины k кончается на r. Какой номер нужно вычесть при сдвиге?", Language: "python", Code: `s = sum(a[:k])
for r in range(k, len(a)):
    s += a[r] - a[___]`, Options: []string{"r - k", "r - k + 1", "r - 1"}, Answer: 0, Explain: "Окно с концом r занимает номера r − k + 1..r. Прошлое окно занимало r − k..r − 1, поэтому вышел a[r − k]. Номер r − k + 1 ещё в окне: для a = [2, 1, 5, 1], k = 3 с ним сумма стала бы 8 − 1 + 1 = 8, а верно 8 − 2 + 1 = 7."},
				{ID: "sw1-set-trace", Type: "choice", Prompt: "LeetCode 219: nums = [4, 1, 2, 3, 4, 1], k = 2. Какие числа лежат в set после шага i = 4?", Options: []string{"{2, 3}", "{3, 4}", "{2, 3, 4}"}, Answer: 1, Explain: "После шага 3 в set лежат 2 и 3. На шаге 4 числа 4 в set нет: его добавляем, затем удаляем nums[2] = 2. Остаются 3 и 4, это nums[3] и nums[4]. Вариант {2, 3, 4} забывает удаление: в set станет больше k чисел, и код сможет найти повтор на расстоянии больше k. {2, 3} лежит в set после шага 3, а не после шага 4."},
				{ID: "sw1-anagram-count", Type: "short", Prompt: "s = \"baabab\", p = \"ab\". Сколько номеров вернёт findAnagrams?", Answer: "4", Explain: "Окна длины 2: ba, aa, ab, ba, ab. Подходят все, кроме aa: номера 0, 2, 3, 4. Каждое из 5 окон меняет не больше двух счётчиков и сравнивает два list по 26 чисел. Ответ 5 получится, если считать aa подходящим окном, но в нём нет b."},
				{ID: "sw1-dict-zero", Type: "choice", Prompt: "LeetCode 2461: окно подходит, когда len(cnt) == k. Зачем удалять ключ, когда его счётчик стал 0?", Options: []string{"иначе сумма total станет неверной", "иначе dict займёт O(n²) памяти", "иначе len(cnt) считает число, которого уже нет в окне"}, Answer: 2, Explain: "nums = [1, 1, 2, 3], k = 2. В окне [2, 3] без удаления останется ключ 1 со счётчиком 0, и len(cnt) = 3 ≠ 2: окно с суммой 5 пропадёт, ответ станет 3. Сумма total меняется отдельно и от ключей не зависит. Без удаления dict растёт до числа разных значений входа, это O(n), а не O(n²)."},
				{ID: "sw1-lecture", Type: "short", Prompt: "Codeforces 961B: a = [2, 4, 1, 3], t = [1, 0, 0, 1], k = 2. Сколько теорем запишет Миша?", Answer: "10", Explain: "Основа: минуты с t = 1 дают 2 + 3 = 5. Выигрыши окон длины 2: 0 + 4 = 4, 4 + 1 = 5, 1 + 0 = 1. Ответ 5 + 5 = 10. Ответ 11 получится, если прибавить полную сумму окна [2, 4]: теоремы первой минуты будут посчитаны дважды."},
				{ID: "sw1-slice-memory", Type: "choice", Prompt: "Вызов sum(a[:k]) считает сумму первого окна. Сколько ссылок хранит срез a[:k], пока sum его читает?", Options: []string{"1", "k", "n"}, Answer: 1, Explain: "Срез a[:k] - новый list из k ссылок: это O(k) памяти и O(k) времени на копию. Одна переменная s нужна циклу по первым k номерам: он не копирует окно и требует O(1) памяти. n ссылок хранила бы копия всего Array a[:]."},
			},
			Links: []sets.Link{
				{Title: "Sliding Window: окна фиксированной и переменной длины", URL: "https://repetitor.1c.ru/informatics/metod-skolzyashchego-okna/", Source: "1С:Репетитор"},
				{Title: "Sliding Window: среднее окон и шаблоны задач", URL: "https://habr.com/ru/articles/347378/", Source: "Хабр"},
				{Title: "Sliding Window: среднее каждого окна длины k", URL: "https://krilovskiy.com/posts/algo-patterns-sliding-window/", Source: "Блог В. Крыловского"},
				{Title: "collections.Counter", URL: "https://docs.python.org/3/library/collections.html#collections.Counter", Source: "Python docs (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "sw1-sum",
					Title: "Сумма окна: a = [2, 1, 5, 1, 3, 2], k = 3",
					Frames: []sets.FigureFrame{
						{Caption: "Первое окно [2, 1, 5]: 3 чтения, сумма 8.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 1, 5, 1, 3, 2}, Marked: []int{0, 1, 2}, Window: &[2]int{0, 2}},
							{Label: "сумма", Cells: []any{8}},
							{Label: "лучшая", Cells: []any{8}},
						}},
						{Caption: "Вошло 1, вышло 2: 8 − 2 + 1 = 7. Лучшая сумма пока 8.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 1, 5, 1, 3, 2}, Marked: []int{3}, Dimmed: []int{0}, Window: &[2]int{1, 3}},
							{Label: "сумма", Cells: []any{7}},
							{Label: "лучшая", Cells: []any{8}},
						}},
						{Caption: "Вошло 3, вышло 1: 7 − 1 + 3 = 9. Это новая лучшая сумма.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 1, 5, 1, 3, 2}, Marked: []int{4}, Dimmed: []int{1}, Window: &[2]int{2, 4}},
							{Label: "сумма", Cells: []any{9}},
							{Label: "лучшая", Cells: []any{9}},
						}},
						{Caption: "Вошло 2, вышло 5: 9 − 5 + 2 = 6. Всего 3 + 3 · 2 = 9 чтений против 12 у Brute Force.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 1, 5, 1, 3, 2}, Marked: []int{5}, Dimmed: []int{2}, Window: &[2]int{3, 5}},
							{Label: "сумма", Cells: []any{6}},
							{Label: "лучшая", Cells: []any{9}},
						}},
					},
				},
				{
					ID:    "sw1-anagram",
					Title: "Счётчики букв: s = \"cbaebabacd\", p = \"abc\"",
					Frames: []sets.FigureFrame{
						{Caption: "Окно cba: счётчики a, b, c равны 1, как у p. Номер 0 входит в ответ. Остальные 21 счётчик равны 0, их не показываем.", Rows: []sets.FigureRow{
							{Label: "s", Cells: []any{"c", "b", "a", "e", "b", "a", "b", "a", "c", "d"}, Marked: []int{2}, Window: &[2]int{0, 2}},
							{Label: "буква", Cells: []any{"a", "b", "c", "d", "e"}},
							{Label: "have", Cells: []any{1, 1, 1, 0, 0}},
							{Label: "need", Cells: []any{1, 1, 1, 0, 0}},
							{Label: "ответ", Cells: []any{0}},
						}},
						{Caption: "Вошла e, вышла c: изменились два счётчика. Счётчик c стал 0, окно bae не подходит.", Rows: []sets.FigureRow{
							{Label: "s", Cells: []any{"c", "b", "a", "e", "b", "a", "b", "a", "c", "d"}, Marked: []int{3}, Dimmed: []int{0}, Window: &[2]int{1, 3}},
							{Label: "буква", Cells: []any{"a", "b", "c", "d", "e"}},
							{Label: "have", Cells: []any{1, 1, 0, 0, 1}},
							{Label: "need", Cells: []any{1, 1, 1, 0, 0}},
							{Label: "ответ", Cells: []any{0}},
						}},
						{Caption: "Окно bac на номерах 6..8: счётчики снова равны need. Номер 6 входит в ответ.", Rows: []sets.FigureRow{
							{Label: "s", Cells: []any{"c", "b", "a", "e", "b", "a", "b", "a", "c", "d"}, Marked: []int{8}, Dimmed: []int{5}, Window: &[2]int{6, 8}},
							{Label: "буква", Cells: []any{"a", "b", "c", "d", "e"}},
							{Label: "have", Cells: []any{1, 1, 1, 0, 0}},
							{Label: "need", Cells: []any{1, 1, 1, 0, 0}},
							{Label: "ответ", Cells: []any{0, 6}},
						}},
						{Caption: "Последнее окно acd не подходит. Ответ [0, 6]: 8 окон и 8 сравнений по 26 счётчиков.", Rows: []sets.FigureRow{
							{Label: "s", Cells: []any{"c", "b", "a", "e", "b", "a", "b", "a", "c", "d"}, Marked: []int{9}, Dimmed: []int{6}, Window: &[2]int{7, 9}},
							{Label: "буква", Cells: []any{"a", "b", "c", "d", "e"}},
							{Label: "have", Cells: []any{1, 0, 1, 1, 0}},
							{Label: "need", Cells: []any{1, 1, 1, 0, 0}},
							{Label: "ответ", Cells: []any{0, 6}},
						}},
					},
				},
				{
					ID:    "sw1-lecture",
					Title: "Выигрыш окна: Codeforces 961B, k = 3",
					Frames: []sets.FigureFrame{
						{Caption: "Номера минут с 0. Основа: Миша сам записывает 1 + 3 + 2 = 6 теорем. Выигрыш первого окна, номера 0..2: 0 + 0 + 5 = 5.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 2, 5, 4}},
							{Label: "t", Cells: []any{1, 1, 0, 1, 0, 0}},
							{Label: "выигрыш", Cells: []any{0, 0, 5, 0, 5, 4}, Marked: []int{0, 1, 2}, Window: &[2]int{0, 2}},
							{Label: "основа", Cells: []any{6}},
							{Label: "лучший выигрыш", Cells: []any{5}},
						}},
						{Caption: "Вошёл номер 3, вышел номер 0: 5 − 0 + 0 = 5.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 2, 5, 4}},
							{Label: "t", Cells: []any{1, 1, 0, 1, 0, 0}},
							{Label: "выигрыш", Cells: []any{0, 0, 5, 0, 5, 4}, Marked: []int{3}, Dimmed: []int{0}, Window: &[2]int{1, 3}},
							{Label: "основа", Cells: []any{6}},
							{Label: "лучший выигрыш", Cells: []any{5}},
						}},
						{Caption: "Вошёл номер 4, вышел номер 1: 5 − 0 + 5 = 10. Лучший выигрыш 10.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 2, 5, 4}},
							{Label: "t", Cells: []any{1, 1, 0, 1, 0, 0}},
							{Label: "выигрыш", Cells: []any{0, 0, 5, 0, 5, 4}, Marked: []int{4}, Dimmed: []int{1}, Window: &[2]int{2, 4}},
							{Label: "основа", Cells: []any{6}},
							{Label: "лучший выигрыш", Cells: []any{10}},
						}},
						{Caption: "Вошёл номер 5, вышел номер 2: 10 − 5 + 4 = 9. Ответ 6 + 10 = 16.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{1, 3, 5, 2, 5, 4}},
							{Label: "t", Cells: []any{1, 1, 0, 1, 0, 0}},
							{Label: "выигрыш", Cells: []any{0, 0, 5, 0, 5, 4}, Marked: []int{5}, Dimmed: []int{2}, Window: &[2]int{3, 5}},
							{Label: "основа", Cells: []any{6}},
							{Label: "лучший выигрыш", Cells: []any{10}},
						}},
					},
				},
			},
		},
		{
			Title:   "Sliding Window переменной длины",
			Summary: "Правая граница расширяет окно, левая сжимает его, пока окно снова не станет допустимым. Обе границы идут только вправо, поэтому весь проход стоит O(n). Вы ищете самое длинное и самое короткое окно и проверяете, когда этот шаблон работает.",
			Lesson:  lesson2,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/max-consecutive-ones/",
					Number:     number(485),
					Title:      "Max Consecutive Ones",
					Difficulty: "easy",
					Topics:     []string{"Array", "Sliding Window"},
					Note:       "Шаблон окна при k = 0: ноль переносит l за себя. n = len(nums). Время: O(n). Память: O(1) дополнительно. Для [1, 1, 0, 1, 1, 1] окна из единиц имеют длины 2 и 3, ответ 3. Brute Force с подсчётом от каждого начала: O(n²) на Array из единиц.",
					Statement: `<p>Дан двоичный Array <code>nums</code>: в нём только числа 0 и 1. Верните <em>наибольшее число единиц <code>1</code>, которые стоят в Array подряд</em>.</p><p>Метод класса <code>Solution</code>: <code>def findMaxConsecutiveOnes(self, nums: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,1,0,1,1,1]
<strong>Выход:</strong> 3
<strong>Пояснение:</strong> Единицы стоят подряд в двух первых цифрах и в трёх последних цифрах. Наибольшее число единиц подряд равно 3.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [1,0,1,1,0,1]
<strong>Выход:</strong> 2</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 10<sup>5</sup></code></li><li><code>nums[i]</code> равно <code>0</code> или <code>1</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "findMaxConsecutiveOnes"}, Cases: []sets.TestCase{
						{Input: "[[1,1,0,1,1,1]]", Expected: text("3")},
						{Input: "[[1,0,1,1,0,1]]", Expected: text("2")},
						{Input: "[[0]]", Expected: text("0")},
						{Input: "[[1]]", Expected: text("1")},
						{Input: "[[1,1,0,1,1,1,1]]", Expected: text("4")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/max-consecutive-ones-iii/",
					Number:     number(1004),
					Title:      "Max Consecutive Ones III",
					Difficulty: "medium",
					Topics:     []string{"Array", "Sliding Window"},
					Note:       "Окно хранит число нулей; пока их больше k, l идёт вправо. n = len(nums). Время: O(n), r и l делают не больше n шагов каждый. Память: O(1) дополнительно. Для [1, 0, 1, 1, 0, 0, 1] и k = 1 ответ 4: окно [1, 0, 1, 1]. Brute Force по всем началам: O(n²).",
					Statement: `<p>Дан двоичный Array <code>nums</code> и целое число <code>k</code>. Верните <em>наибольшее число единиц <code>1</code>, которые стоят в Array подряд, если можно заменить на 1 не больше</em> <code>k</code> <em>нулей</em> <code>0</code>.</p><p>Метод класса <code>Solution</code>: <code>def longestOnes(self, nums: list[int], k: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,1,1,0,0,0,1,1,1,1,0], k = 2
<strong>Выход:</strong> 6
<strong>Пояснение:</strong> [1,1,1,0,0,1,1,1,1,1,1]
Нули на местах 5 и 10 (номера с 0) заменены на 1. Самый длинный Subarray из единиц занимает места 5..10.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [0,0,1,1,0,0,1,1,1,0,1,1,0,0,0,1,1,1,1], k = 3
<strong>Выход:</strong> 10
<strong>Пояснение:</strong> [0,0,1,1,1,1,1,1,1,1,1,1,0,0,0,1,1,1,1]
Нули на местах 4, 5 и 9 (номера с 0) заменены на 1. Самый длинный Subarray из единиц занимает места 2..11.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 10<sup>5</sup></code></li><li><code>nums[i]</code> равно <code>0</code> или <code>1</code></li><li><code>0 ≤ k ≤ nums.length</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "longestOnes"}, Cases: []sets.TestCase{
						{Input: "[[1,1,1,0,0,0,1,1,1,1,0],2]", Expected: text("6")},
						{Input: "[[0,0,1,1,0,0,1,1,1,0,1,1,0,0,0,1,1,1,1],3]", Expected: text("10")},
						{Input: "[[0,0,0],0]", Expected: text("0")},
						{Input: "[[0,0,0],3]", Expected: text("3")},
						{Input: "[[1,0,1,0,1],1]", Expected: text("3")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/fruit-into-baskets/",
					Number:     number(904),
					Title:      "Fruit Into Baskets",
					Difficulty: "medium",
					Topics:     []string{"Array", "Hash Table", "Sliding Window"},
					Note:       "dict счётчиков видов окна; пока видов больше двух, l идёт вправо, ключ со счётчиком 0 удаляем. n = len(fruits). Время: O(n): r и l делают не больше n шагов каждый. Память: O(1) дополнительно, в dict не больше трёх ключей. Для [1, 2, 3, 2, 2] окно [2, 3, 2, 2] даёт ответ 4. Перебор начал с проходом вправо: O(n²), на 10⁵ деревьях двух видов это 5 000 050 000 шагов.",
					Statement: `<p>Вы пришли на ферму, где фруктовые деревья стоят в один ряд слева направо. Деревья заданы Array целых чисел <code>fruits</code>: <code>fruits[i]</code> - <strong>вид</strong> фрукта, который растёт на <code>i</code>-м дереве.</p><p>Вы хотите собрать как можно больше фруктов. Но у хозяина есть строгие правила, и их нужно выполнять:</p><ul><li>У вас только <strong>две</strong> корзины, и в каждой корзине могут лежать фрукты только <strong>одного вида</strong>. Число фруктов в одной корзине не ограничено.</li><li>Вы начинаете с любого дерева и идёте вправо. С <strong>каждого</strong> дерева, включая первое, нужно сорвать <strong>ровно один фрукт</strong>. Сорванный фрукт должен подходить к одной из ваших корзин.</li><li>Когда вы дошли до дерева, фрукт которого не подходит к вашим корзинам, нужно остановиться.</li></ul><p>Дан Array целых чисел <code>fruits</code>. Верните <em><strong>наибольшее</strong> число фруктов, которые можно собрать</em>.</p><p>Метод класса <code>Solution</code>: <code>def totalFruit(self, fruits: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> fruits = [1,2,1]
<strong>Выход:</strong> 3
<strong>Пояснение:</strong> Можно собрать фрукты со всех 3 деревьев.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> fruits = [0,1,2,2]
<strong>Выход:</strong> 3
<strong>Пояснение:</strong> Можно собрать фрукты с деревьев [1,2,2].
Если начать с первого дерева, получится собрать фрукты только с деревьев [0,1].</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> fruits = [1,2,3,2,2]
<strong>Выход:</strong> 4
<strong>Пояснение:</strong> Можно собрать фрукты с деревьев [2,3,2,2].
Если начать с первого дерева, получится собрать фрукты только с деревьев [1,2].</pre><h3>Ограничения</h3><ul><li><code>1 ≤ fruits.length ≤ 10<sup>5</sup></code></li><li><code>0 ≤ fruits[i] &lt; fruits.length</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "totalFruit"}, Cases: []sets.TestCase{
						{Input: "[[1,2,1]]", Expected: text("3")},
						{Input: "[[0,1,2,2]]", Expected: text("3")},
						{Input: "[[1,2,3,2,2]]", Expected: text("4")},
						{Input: "[[0]]", Expected: text("1")},
						{Input: "[[2,2,2,2]]", Expected: text("4")},
						{Input: "[[0,1,1,2,1]]", Expected: text("4")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1692/E",
					Title:      "Binary Deque",
					Difficulty: "medium",
					Topics:     []string{"Array", "Sliding Window"},
					Note:       "Codeforces, рейтинг 1200. Удаления с концов оставляют окно; ищем самое длинное окно с суммой не больше s. n - длина набора. Время: O(n) на набор, O(Σn) на весь вход. Память: O(1) дополнительно, вход O(n). Если сумма всех чисел меньше s, ответ −1. Для 0 1 0 1 1 1 0 0 1 и s = 3 самое длинное окно имеет длину 6, ответ 9 − 6 = 3. Перебор числа удалений слева и справа: O(n²).",
					Statement: `<p>У Славика есть Array длины <code>n</code>, в котором только нули и единицы. За одну операцию он удаляет первый или последний элемент Array.</p><p>Какое наименьшее число операций нужно сделать Славику, чтобы после всех операций сумма Array была равна <code>s</code>? Если сумму <code>s</code> нельзя получить ни при каком числе операций, выведите -1.</p><h3>Входные данные</h3><p>В первой строке записано целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>): количество наборов входных данных.</p><p>В первой строке каждого набора записаны два целых числа <code>n</code> и <code>s</code> (<code>1 ≤ n, s ≤ 2 · 10<sup>5</sup></code>): длина Array и нужная сумма элементов.</p><p>Во второй строке набора записаны <code>n</code> целых чисел <code>a<sub>i</sub></code> (<code>0 ≤ a<sub>i</sub> ≤ 1</code>): элементы Array.</p><p>Гарантируется, что сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого набора выведите одно целое число: наименьшее число операций, после которых сумма Array равна <code>s</code>, или -1, если Array с суммой <code>s</code> получить нельзя.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>7
3 1
1 0 0
3 1
1 1 0
9 3
0 1 0 1 1 1 0 0 1
6 4
1 1 1 1 1 1
5 1
0 0 1 1 0
16 2
1 1 0 0 1 0 0 1 1 0 0 0 0 0 1 1
6 3
1 0 1 0 0 0</pre><p><strong>Выходные данные:</strong></p><pre>0
1
3
2
2
7
-1</pre><h3>Примечание</h3><p>В первом наборе сумма всего Array с самого начала равна 1, поэтому операции не нужны.</p><p>Во втором наборе сумма Array равна 2, а нужна сумма 1, поэтому нужно удалить первый элемент. Array станет <code>[1, 0]</code>, его сумма равна 1.</p><p>В третьем наборе сумма Array равна 5, а нужна сумма 3. Такую сумму можно получить, если удалить два первых элемента и последний элемент: всего три операции. Array станет <code>[0, 1, 1, 1, 0, 0]</code>, его сумма равна 3.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>1 ≤ n, s ≤ 2 · 10<sup>5</sup></code></li><li><code>0 ≤ a<sub>i</sub> ≤ 1</code></li><li>сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "7\n3 1\n1 0 0\n3 1\n1 1 0\n9 3\n0 1 0 1 1 1 0 0 1\n6 4\n1 1 1 1 1 1\n5 1\n0 0 1 1 0\n16 2\n1 1 0 0 1 0 0 1 1 0 0 0 0 0 1 1\n6 3\n1 0 1 0 0 0\n", Expected: text("0\n1\n3\n2\n2\n7\n-1\n")},
						{Input: "1\n1 1\n1\n", Expected: text("0\n")},
						{Input: "1\n5 1\n0 0 0 0 0\n", Expected: text("-1\n")},
						{Input: "1\n6 1\n1 0 0 0 0 1\n", Expected: text("1\n")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/minimum-window-substring/",
					Number:     number(76),
					Title:      "Minimum Window Substring",
					Difficulty: "hard",
					Topics:     []string{"Hash Table", "String", "Sliding Window"},
					Note:       "need считает недостающие символы t, missing - их общее число. Ответ обновляем внутри while. m = len(s), n = len(t). Время: O(m + n). Память: O(1) дополнительно для 52 букв, ответ O(длина окна). Для ADOBECODEBANC и ABC ответ BANC. Проверка всех окон: O(m²) окон.",
					Statement: `<p>Даны две String <code>s</code> и <code>t</code> длины <code>m</code> и <code>n</code>. Верните <em><strong>самое короткое окно</strong>: такую</em> <strong><em>подстроку</em></strong> <em>строки</em> <code>s</code><em>, в которой есть каждый символ строки</em> <code>t</code> <em>(<strong>с повторами</strong>)</em>. Если такой подстроки нет, верните <em>пустую строку</em> <code>""</code>.</p><p>Тесты составлены так, что ответ <strong>единственный</strong>.</p><p><em>Подстрока - непрерывная непустая часть строки.</em></p><p>Метод класса <code>Solution</code>: <code>def minWindow(self, s: str, t: str) -&gt; str</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> s = "ADOBECODEBANC", t = "ABC"
<strong>Выход:</strong> "BANC"
<strong>Пояснение:</strong> В самом коротком окне "BANC" есть 'A', 'B' и 'C' из строки t.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> s = "a", t = "a"
<strong>Выход:</strong> "a"
<strong>Пояснение:</strong> Вся строка s - самое короткое окно.</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> s = "a", t = "aa"
<strong>Выход:</strong> ""
<strong>Пояснение:</strong> В окне должны быть обе буквы 'a' из t.
В самом длинном окне s только одна 'a', поэтому верните пустую строку.</pre><h3>Ограничения</h3><ul><li><code>m == s.length</code></li><li><code>n == t.length</code></li><li><code>1 ≤ m, n ≤ 10<sup>5</sup></code></li><li><code>s</code> и <code>t</code> состоят из заглавных и строчных английских букв</li></ul><p><strong>Дополнительный вопрос:</strong> можете ли вы найти алгоритм, который работает за <code>O(m + n)</code> времени?</p>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "minWindow"}, Cases: []sets.TestCase{
						{Input: `["ADOBECODEBANC","ABC"]`, Expected: text(`"BANC"`)},
						{Input: `["a","a"]`, Expected: text(`"a"`)},
						{Input: `["a","aa"]`, Expected: text(`""`)},
						{Input: `["bba","ab"]`, Expected: text(`"ba"`)},
						{Input: `["cabwefgewcwaefgcf","cae"]`, Expected: text(`"cwae"`)},
						{Input: `["aA","A"]`, Expected: text(`"A"`)},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "sw2-moves", Type: "short", Prompt: "Шаблон LeetCode 1004: nums = [0, 0, 1, 0, 1], k = 1. Сколько раз выполнится l += 1 за весь проход?", Answer: "2", Explain: "На r = 1 в окне два нуля: l убирает nums[0] и встаёт на 1. На r = 3 снова два нуля: l убирает nums[1] и встаёт на 2. Всего 2 шага l. Шагов l за весь проход не больше n = 5, поэтому вложенный while не даёт O(n²). Ответ 3 - это длина лучшего окна, а не число шагов l."},
				{ID: "sw2-while", Type: "choice", Prompt: "В шаблоне внутри for стоит while. Почему всё время O(n), а не O(n²)?", Options: []string{"while выполняется не больше одного раза на шаг", "l только растёт и за весь проход сдвигается не больше n раз", "окно никогда не длиннее k"}, Answer: 1, Explain: "Один шаг r может сдвинуть l много раз: в разделе 1 на r = 5 было 3 сдвига. Но l не уменьшается и не выходит за n, поэтому все while вместе делают не больше n сдвигов. В задаче 1004 окно может быть длиннее k: k ограничивает число нулей, а не длину."},
				{ID: "sw2-negative", Type: "choice", Prompt: "Шаблон «сжимать, пока сумма больше S = 2» ищет самое длинное окно. Для какого Array он даст неверную длину?", Options: []string{"[3, 1, 2]", "[0, 0, 5]", "[4, −3, 1]"}, Answer: 2, Explain: "Для [4, −3, 1] шаблон убирает 4 и отвечает 2, а окно [4, −3, 1] имеет сумму 2 и длину 3. В двух других Array числа не отрицательные: расширение не уменьшает сумму, и шаблон верен. Там ответы 1 и 2."},
				{ID: "sw2-abba", Type: "code", Prompt: "Что вернёт этот код для s = \"tmmzuxt\"? Верный ответ LeetCode 3 равен 5.", Language: "python", Code: `last = {}
l = 0
best = 0
for r, c in enumerate(s):
    if c in last:
        l = last[c] + 1
    last[c] = r
    best = max(best, r - l + 1)`, Options: []string{"5", "6", "7"}, Answer: 1, Explain: "На r = 2 повтор m даёт l = 2. На r = 6 символ t был на месте 0, и код ставит l = 1: граница идёт назад, и окно mmzuxt содержит повтор. Код вернёт 6. Нужно l = max(l, last[c] + 1). Ответ 7 был бы совсем без проверки повторов."},
				{ID: "sw2-shortest", Type: "choice", Prompt: "LeetCode 76: где обновлять ответ, чтобы найти самое короткое хорошее окно?", Options: []string{"после while, как для самого длинного окна", "только при r = len(s) − 1", "внутри while, пока окно хорошее"}, Answer: 2, Explain: "После while окно уже плохое: в нём не хватает символа. Хорошие окна видны только внутри while, перед каждым сжатием. Для s = \"bdab\", t = \"ab\" окно ab появляется внутри while на шаге r = 3. Проверка только в конце пропустит окна, которые кончаются раньше."},
				{ID: "sw2-binary-deque", Type: "short", Prompt: "Codeforces 1692E: a = [1, 0, 1, 0, 0, 1], s = 2. Сколько удалений нужно?", Answer: "1", Explain: "Сумма всего Array 3. Самое длинное окно с суммой 2 имеет длину 5: номера 0..4 или 1..5. Ответ 6 − 5 = 1: удалить последний или первый элемент. Ответ 5 - это длина окна, а не число удалений."},
				{ID: "sw2-missing", Type: "short", Prompt: "Код раздела 5 для s = \"abcab\", t = \"ba\". Сколько раз выполнится l += 1?", Answer: "4", Explain: "На r = 1 окно ab хорошее: l убирает a, окно стало плохим, l = 1. На r = 3 окно bca хорошее: l убирает b, окно стало плохим, l = 2. На r = 4 окно cab хорошее: l убирает c, окно ab ещё хорошее, потом l убирает a, l = 4. Всего 4 шага. Шагов l не больше len(s) = 5."},
				{ID: "sw2-memory", Type: "choice", Prompt: "Самое длинное окно без повторов в list из n целых чисел, dict счётчиков окна. В list σ разных чисел, и σ может быть до n. Какая дополнительная память?", Options: []string{"O(min(n, σ))", "O(1)", "O(n²)"}, Answer: 0, Explain: "dict хранит по ключу на каждое разное число окна: их не больше n и не больше σ. Для list из n разных чисел это n ключей. O(1) верно только для постоянного небольшого алфавита, например для 95 символов строки LeetCode 3, а здесь σ растёт вместе с n. O(n²) было бы при хранении всех окон."},
			},
			Links: []sets.Link{
				{Title: "Sliding Window: подстрока без повторов", URL: "https://habr.com/ru/articles/1007886/", Source: "Хабр"},
				{Title: "Sliding Window: самое короткое окно подстроки", URL: "https://sprintcode.pro/ru/blog/minimum-window-substring", Source: "SprintCode.pro"},
				{Title: "Sliding Window: общий шаблон для задач 76, 438 и 3", URL: "https://labuladong.online/en/algo/essential-technique/sliding-window-framework/", Source: "labuladong (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "sw2-template",
					Title: "Шаблон окна: nums = [1, 0, 1, 1, 0, 0, 1], k = 1",
					Frames: []sets.FigureFrame{
						{Caption: "r = 3: окно [1, 0, 1, 1] содержит один ноль. Длина 4, это лучший ответ.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, 0, 1, 1, 0, 0, 1}, Marked: []int{3}, Pointers: &sets.Pointers{{Name: "l", Cell: 0}, {Name: "r", Cell: 3}}, Window: &[2]int{0, 3}},
							{Label: "нулей", Cells: []any{1}},
							{Label: "лучшая", Cells: []any{4}},
						}},
						{Caption: "r = 4: второй ноль. l убирает 1 и 0 и встаёт на номер 2. Окно [1, 1, 0], длина 3.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, 0, 1, 1, 0, 0, 1}, Marked: []int{4}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "l", Cell: 2}, {Name: "r", Cell: 4}}, Window: &[2]int{2, 4}},
							{Label: "нулей", Cells: []any{1}},
							{Label: "лучшая", Cells: []any{4}},
						}},
						{Caption: "r = 5: снова два нуля. l делает 3 шага до номера 5. Окно [0], длина 1.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, 0, 1, 1, 0, 0, 1}, Marked: []int{5}, Dimmed: []int{2, 3, 4}, Pointers: &sets.Pointers{{Name: "l", Cell: 5}, {Name: "r", Cell: 5}}, Window: &[2]int{5, 5}},
							{Label: "нулей", Cells: []any{1}},
							{Label: "лучшая", Cells: []any{4}},
						}},
						{Caption: "r = 6: окно [0, 1], длина 2. Ответ 4. За весь проход r сделал 7 шагов, l - 5.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, 0, 1, 1, 0, 0, 1}, Marked: []int{6}, Pointers: &sets.Pointers{{Name: "l", Cell: 5}, {Name: "r", Cell: 6}}, Window: &[2]int{5, 6}},
							{Label: "нулей", Cells: []any{1}},
							{Label: "лучшая", Cells: []any{4}},
						}},
					},
				},
				{
					ID:    "sw2-negative",
					Title: "Отрицательные числа: a = [5, −4, 1], S = 2",
					Frames: []sets.FigureFrame{
						{Caption: "r = 0: сумма 5 > 2. Шаблон убирает 5, окно пустое, l = 1.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{5, -4, 1}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "r", Cell: 0}}},
							{Label: "сумма", Cells: []any{0}},
							{Label: "лучшая", Cells: []any{0}},
						}},
						{Caption: "r = 1: окно [−4], сумма −4, длина 1.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{5, -4, 1}, Marked: []int{1}, Pointers: &sets.Pointers{{Name: "l", Cell: 1}, {Name: "r", Cell: 1}}, Window: &[2]int{1, 1}},
							{Label: "сумма", Cells: []any{-4}},
							{Label: "лучшая", Cells: []any{1}},
						}},
						{Caption: "r = 2: окно [−4, 1], сумма −3, длина 2. Шаблон отвечает 2.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{5, -4, 1}, Marked: []int{2}, Pointers: &sets.Pointers{{Name: "l", Cell: 1}, {Name: "r", Cell: 2}}, Window: &[2]int{1, 2}},
							{Label: "сумма", Cells: []any{-3}},
							{Label: "лучшая", Cells: []any{2}},
						}},
						{Caption: "Brute Force проверяет все 6 окон: окно [5, −4, 1] имеет сумму 2 ≤ S и длину 3. Шаблон ответил 2: начало 0 отброшено зря, после −4 сумма снова стала меньше.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{5, -4, 1}, Marked: []int{0, 1, 2}, Window: &[2]int{0, 2}},
							{Label: "сумма", Cells: []any{2}},
							{Label: "лучшая", Cells: []any{2}},
							{Label: "Brute Force", Cells: []any{3}},
						}},
					},
				},
				{
					ID:    "sw2-min",
					Title: "Самое короткое окно: s = \"bdab\", t = \"ab\"",
					Frames: []sets.FigureFrame{
						{Caption: "r = 2: в окне bda есть a и b. Окно хорошее, ответ bda, длина 3.", Rows: []sets.FigureRow{
							{Label: "s", Cells: []any{"b", "d", "a", "b"}, Marked: []int{2}, Pointers: &sets.Pointers{{Name: "l", Cell: 0}, {Name: "r", Cell: 2}}, Window: &[2]int{0, 2}},
							{Label: "не хватает", Cells: []any{0}},
							{Label: "лучшее", Cells: []any{"bda"}},
						}},
						{Caption: "Убрали b: теперь не хватает одного символа, l = 1.", Rows: []sets.FigureRow{
							{Label: "s", Cells: []any{"b", "d", "a", "b"}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "l", Cell: 1}, {Name: "r", Cell: 2}}, Window: &[2]int{1, 2}},
							{Label: "не хватает", Cells: []any{1}},
							{Label: "лучшее", Cells: []any{"bda"}},
						}},
						{Caption: "r = 3: окно dab снова хорошее. Длина 3 не лучше.", Rows: []sets.FigureRow{
							{Label: "s", Cells: []any{"b", "d", "a", "b"}, Marked: []int{3}, Pointers: &sets.Pointers{{Name: "l", Cell: 1}, {Name: "r", Cell: 3}}, Window: &[2]int{1, 3}},
							{Label: "не хватает", Cells: []any{0}},
							{Label: "лучшее", Cells: []any{"bda"}},
						}},
						{Caption: "Убрали d: окно ab всё ещё хорошее. Новый ответ ab, длина 2.", Rows: []sets.FigureRow{
							{Label: "s", Cells: []any{"b", "d", "a", "b"}, Dimmed: []int{1}, Pointers: &sets.Pointers{{Name: "l", Cell: 2}, {Name: "r", Cell: 3}}, Window: &[2]int{2, 3}},
							{Label: "не хватает", Cells: []any{0}},
							{Label: "лучшее", Cells: []any{"ab"}},
						}},
						{Caption: "Убрали a: не хватает одного символа, l = 3. Проход закончен, ответ ab.", Rows: []sets.FigureRow{
							{Label: "s", Cells: []any{"b", "d", "a", "b"}, Dimmed: []int{2}, Pointers: &sets.Pointers{{Name: "l", Cell: 3}, {Name: "r", Cell: 3}}, Window: &[2]int{3, 3}},
							{Label: "не хватает", Cells: []any{1}},
							{Label: "лучшее", Cells: []any{"ab"}},
						}},
					},
				},
			},
		},
		{
			Title:   "Подсчёт окон",
			Summary: "Чтобы посчитать хорошие Subarray, не выписывайте их. Для каждой правой границы r хорошие начала образуют отрезок: прибавьте его длину. «Ровно K» считаем как «не больше K» минус «не больше K − 1».",
			Lesson:  lesson3,
			Tasks: []sets.Task{
				{
					URL:        "https://codeforces.com/problemset/problem/1840/C",
					Title:      "Ski Resort",
					Difficulty: "easy",
					Topics:     []string{"Array", "Math", "Sliding Window"},
					Note:       "Codeforces, рейтинг 1000. Для каждого дня r прибавляем max(0, L − k + 1), где L - число подходящих дней подряд, которые кончаются на r. Время: O(n) на набор, O(Σn) на весь вход. Память: O(1) дополнительно. Ответ до n(n + 1)/2 = 20 000 100 000, больше 32-битного целого. Для 0 3 −2 5 −4 −4 при k = 1 и q = 3 блоки из 3 и 2 дней дают 6 + 3 = 9. Перебор всех отрезков: O(n²).",
					Statement: `<p>Дима Ватрушин - учитель математики в школе. За хорошую работу его отправили в отпуск на <code>n</code> дней. Дима давно мечтает поехать на горнолыжный курорт, поэтому он хочет выделить несколько дней подряд и покататься на лыжах. Отпуск требует хорошей подготовки, поэтому Дима поедет только на <code>k</code> дней или больше.</p><p>Дан Array <code>a</code> с прогнозом погоды на курорте: в <code>i</code>-й день температура будет <code>a<sub>i</sub></code> градусов.</p><p>Дима родился в Сибири, поэтому он может поехать в отпуск, только если температура во все дни отпуска не поднимается выше <code>q</code> градусов.</p><p>К сожалению, Дима так увлёкся абстрактной алгеброй, что забыл, как считать. Он просит вас помочь ему и посчитать число способов выбрать даты отпуска на курорте.</p><h3>Входные данные</h3><p>В первой строке записано целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>): количество наборов входных данных.</p><p>Дальше идут описания наборов.</p><p>В первой строке набора записаны три целых числа <code>n</code>, <code>k</code>, <code>q</code> (<code>1 ≤ n ≤ 2 · 10<sup>5</sup></code>, <code>1 ≤ k ≤ n</code>, <code>−10<sup>9</sup> ≤ q ≤ 10<sup>9</sup></code>): длина Array <code>a</code>, наименьшее число дней на курорте и наибольшая удобная для Димы температура.</p><p>Во второй строке набора записаны <code>n</code> целых чисел <code>a<sub>1</sub>, a<sub>2</sub>, a<sub>3</sub>, ..., a<sub>n</sub></code> (<code>−10<sup>9</sup> ≤ a<sub>i</sub> ≤ 10<sup>9</sup></code>): температура на курорте.</p><p>Сумма всех значений <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Выведите <code>t</code> целых чисел. Каждое число - ответ для своего набора: число способов, которыми Дима может выбрать даты отпуска на курорте.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>7
3 1 15
-5 0 -10
5 3 -33
8 12 9 0 5
4 3 12
12 12 10 15
4 1 -5
0 -1 2 5
5 5 0
3 -1 4 -5 -3
1 1 5
5
6 1 3
0 3 -2 5 -4 -4</pre><p><strong>Выходные данные:</strong></p><pre>6
0
1
0
0
1
9</pre><h3>Примечание</h3><p>В первом наборе примера Дима может поехать в любой день, поэтому ему подходят даты [1], [2], [3], [1, 2], [2, 3], [1, 2, 3].</p><p>Во втором и четвёртом наборах примера Дима не может поехать ни в какой день из-за высокой температуры, поэтому подходящих дат нет.</p><p>В третьем наборе примера Диме подходят только даты [1, 2, 3].</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>1 ≤ k ≤ n ≤ 2 · 10<sup>5</sup></code></li><li><code>−10<sup>9</sup> ≤ q, a<sub>i</sub> ≤ 10<sup>9</sup></code></li><li>сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code></li><li>Ограничение по времени: 1 секунда на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "7\n3 1 15\n-5 0 -10\n5 3 -33\n8 12 9 0 5\n4 3 12\n12 12 10 15\n4 1 -5\n0 -1 2 5\n5 5 0\n3 -1 4 -5 -3\n1 1 5\n5\n6 1 3\n0 3 -2 5 -4 -4\n", Expected: text("6\n0\n1\n0\n0\n1\n9\n")},
						{Input: "1\n1 1 0\n0\n", Expected: text("1\n")},
						{Input: "1\n4 2 5\n5 5 5 5\n", Expected: text("6\n")},
						{Input: "1\n5 2 0\n0 0 1 0 0\n", Expected: text("2\n")},
						{Input: "1\n6 1 -1\n-1 -1 -1 -1 -1 -1\n", Expected: text("21\n")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/subarray-product-less-than-k/",
					Number:     number(713),
					Title:      "Subarray Product Less Than K",
					Difficulty: "medium",
					Topics:     []string{"Array", "Sliding Window"},
					Note:       "Окно хранит произведение; при конце r прибавляем r − l + 1. n = len(nums). Время: O(n): n умножений, не больше n делений. Память: O(1) дополнительно. При k ≤ 1 ответ 0. Для [10, 5, 2, 6] и k = 100 вклады 1, 2, 2, 3, ответ 8. Перебор начал: O(n²).",
					Statement: `<p>Дан Array целых чисел <code>nums</code> и целое число <code>k</code>.</p><p>Верните число непрерывных Subarray, в которых <strong>произведение</strong> всех элементов <strong>Subarray</strong> <strong>строго меньше</strong> <code>k</code>.</p><p>Метод класса <code>Solution</code>: <code>def numSubarrayProductLessThanK(self, nums: list[int], k: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [10,5,2,6], k = 100
<strong>Выход:</strong> 8
<strong>Пояснение:</strong> Произведение меньше 100 у 8 Subarray:
[10], [5], [2], [6], [10, 5], [5, 2], [2, 6], [5, 2, 6]
Subarray [10, 5, 2] не подходит: его произведение 100 не строго меньше k.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [1,2,3], k = 0
<strong>Выход:</strong> 0</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 3 · 10<sup>4</sup></code></li><li><code>1 ≤ nums[i] ≤ 1000</code></li><li><code>0 ≤ k ≤ 10<sup>6</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "numSubarrayProductLessThanK"}, Cases: []sets.TestCase{
						{Input: "[[10,5,2,6],100]", Expected: text("8")},
						{Input: "[[1,2,3],0]", Expected: text("0")},
						{Input: "[[1,1,1],1]", Expected: text("0")},
						{Input: "[[1,1,1],2]", Expected: text("6")},
						{Input: "[[1000],1000]", Expected: text("0")},
						{Input: "[[2,5,3,10],30]", Expected: text("6")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/number-of-substrings-containing-all-three-characters/",
					Number:     number(1358),
					Title:      "Number of Substrings Containing All Three Characters",
					Difficulty: "medium",
					Topics:     []string{"Hash Table", "String", "Sliding Window"},
					Note:       "Сжимаем окно, пока в нём есть a, b и c; после while прибавляем l. n = len(s). Время: O(n). Память: O(1): три счётчика. Для abcabc вклады 0, 0, 1, 2, 3, 4, ответ 10. Перебор начал: O(n²).",
					Statement: `<p>Дана строка <code>s</code>, в которой есть только символы <em>a</em>, <em>b</em> и <em>c</em>.</p><p>Верните число подстрок, в которых <strong>хотя бы</strong> по одному разу есть каждый из символов <em>a</em>, <em>b</em> и <em>c</em>.</p><p>Метод класса <code>Solution</code>: <code>def numberOfSubstrings(self, s: str) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> s = "abcabc"
<strong>Выход:</strong> 10
<strong>Пояснение:</strong> Хотя бы по одному разу символы a, b и c есть в подстроках "abc", "abca", "abcab", "abcabc", "bca", "bcab", "bcabc", "cab", "cabc" и "abc" (<strong>ещё раз</strong>).</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> s = "aaacb"
<strong>Выход:</strong> 3
<strong>Пояснение:</strong> Хотя бы по одному разу символы a, b и c есть в подстроках "aaacb", "aacb" и "acb".</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> s = "abc"
<strong>Выход:</strong> 1</pre><h3>Ограничения</h3><ul><li><code>3 ≤ s.length ≤ 5 · 10<sup>4</sup></code></li><li><code>s</code> состоит только из символов <code>'a'</code>, <code>'b'</code> и <code>'c'</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "numberOfSubstrings"}, Cases: []sets.TestCase{
						{Input: `["abcabc"]`, Expected: text("10")},
						{Input: `["aaacb"]`, Expected: text("3")},
						{Input: `["abc"]`, Expected: text("1")},
						{Input: `["aaa"]`, Expected: text("0")},
						{Input: `["cccbbbaaa"]`, Expected: text("9")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1736/C1",
					Title:      "Good Subarrays (Easy Version)",
					Difficulty: "medium",
					Topics:     []string{"Array", "Sliding Window"},
					Note:       "Codeforces, рейтинг 1300. Элемент a[r] запрещает начала левее r − a[r] + 1: l = max(l, r − a[r] + 1), прибавляем r − l + 1. Время: O(n) на набор, O(Σn) на весь вход. Память: O(1) дополнительно. Ответ до n(n + 1)/2. Для 2 1 4 3 вклады 1, 1, 2, 3, ответ 7. Перебор отрезков: O(n²).",
					Statement: `<p>Это простая версия задачи. В этой версии нет запросов, и во входе несколько наборов.</p><p>Array <code>b</code> длины <code>m</code> называется хорошим, если для всех <code>i</code> его <code>i</code>-й элемент не меньше <code>i</code>. Другими словами, <code>b</code> хороший тогда и только тогда, когда <code>b<sub>i</sub> ≥ i</code> для всех <code>i</code> (<code>1 ≤ i ≤ m</code>).</p><p>Дан Array <code>a</code> из <code>n</code> целых положительных чисел. Найдите число пар номеров <code>(l, r)</code>, где <code>1 ≤ l ≤ r ≤ n</code>, для которых Array <code>[a<sub>l</sub>, a<sub>l+1</sub>, ..., a<sub>r</sub>]</code> хороший.</p><h3>Входные данные</h3><p>Каждый тест содержит несколько наборов входных данных. В первой строке записано число наборов <code>t</code> (<code>1 ≤ t ≤ 2 · 10<sup>5</sup></code>). Дальше идут описания наборов.</p><p>В первой строке набора записано целое число <code>n</code> (<code>1 ≤ n ≤ 2 · 10<sup>5</sup></code>): длина Array <code>a</code>.</p><p>Во второй строке набора записаны <code>n</code> целых чисел <code>a<sub>1</sub>, a<sub>2</sub>, ..., a<sub>n</sub></code> через пробел (<code>1 ≤ a<sub>i</sub> ≤ n</code>): Array <code>a</code>.</p><p>Гарантируется, что сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого набора выведите число подходящих пар номеров.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>3
3
1 2 3
3
1 1 1
4
2 1 4 3</pre><p><strong>Выходные данные:</strong></p><pre>6
3
7</pre><h3>Примечание</h3><p>В первом наборе все Subarray Array <code>a</code> хорошие, поэтому подходят все пары.</p><p>Во втором наборе подходят пары <code>(1, 1)</code>, <code>(2, 2)</code> и <code>(3, 3)</code>. Например, при <code>(l, r) = (1, 2)</code> Array <code>b = [1, 1]</code> не хороший, потому что <code>b<sub>2</sub> &lt; 2</code>.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 2 · 10<sup>5</sup></code></li><li><code>1 ≤ n ≤ 2 · 10<sup>5</sup></code></li><li><code>1 ≤ a<sub>i</sub> ≤ n</code></li><li>сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code></li><li>Ограничение по времени: 1 секунда на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "3\n3\n1 2 3\n3\n1 1 1\n4\n2 1 4 3\n", Expected: text("6\n3\n7\n")},
						{Input: "1\n1\n1\n", Expected: text("1\n")},
						{Input: "1\n5\n5 5 5 5 5\n", Expected: text("15\n")},
						{Input: "1\n5\n1 2 1 2 3\n", Expected: text("9\n")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/subarrays-with-k-different-integers/",
					Number:     number(992),
					Title:      "Subarrays with K Different Integers",
					Difficulty: "hard",
					Topics:     []string{"Array", "Hash Table", "Sliding Window", "Counting"},
					Note:       "Ровно k = не больше k минус не больше k − 1; каждый проход прибавляет r − l + 1. n = len(nums). Время: в среднем O(n), два прохода. Память: O(min(n, k + 1)) ключей dict. Для [1, 2, 1, 2, 3] и k = 2 получаем 12 − 5 = 7. Перебор начал с set: O(n²).",
					Statement: `<p>Дан Array целых чисел <code>nums</code> и целое число <code>k</code>. Верните <em>число <strong>хороших Subarray</strong> Array</em> <code>nums</code>.</p><p><strong>Хороший Array</strong> - это Array, в котором ровно <code>k</code> разных целых чисел.</p><ul><li>Например, в <code>[1,2,3,1,2]</code> есть <code>3</code> разных целых числа: <code>1</code>, <code>2</code> и <code>3</code>.</li></ul><p><strong>Subarray</strong> - это <strong>непрерывная</strong> часть Array.</p><p>Метод класса <code>Solution</code>: <code>def subarraysWithKDistinct(self, nums: list[int], k: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,2,1,2,3], k = 2
<strong>Выход:</strong> 7
<strong>Пояснение:</strong> Ровно 2 разных целых числа есть в Subarray [1,2], [2,1], [1,2], [2,3], [1,2,1], [2,1,2], [1,2,1,2]</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [1,2,1,3,4], k = 3
<strong>Выход:</strong> 3
<strong>Пояснение:</strong> Ровно 3 разных целых числа есть в Subarray [1,2,1,3], [2,1,3], [1,3,4].</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 2 · 10<sup>4</sup></code></li><li><code>1 ≤ nums[i], k ≤ nums.length</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "subarraysWithKDistinct"}, Cases: []sets.TestCase{
						{Input: "[[1,2,1,2,3],2]", Expected: text("7")},
						{Input: "[[1,2,1,3,4],3]", Expected: text("3")},
						{Input: "[[1],1]", Expected: text("1")},
						{Input: "[[1,1,1],1]", Expected: text("6")},
						{Input: "[[1,2,3],3]", Expected: text("1")},
						{Input: "[[1,2,1,2],3]", Expected: text("0")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "sw3-product", Type: "short", Prompt: "LeetCode 713: nums = [3, 1, 2], k = 6. Сколько Subarray с произведением меньше 6?", Answer: "5", Explain: "r = 0: окно [3], +1. r = 1: окно [3, 1], произведение 3, +2. r = 2: произведение 6 не меньше 6, убираем 3, окно [1, 2], +2. Всего 5. Ответ 6 получится, если сравнивать ≤ вместо <: у окна [3, 1, 2] произведение ровно 6."},
				{ID: "sw3-why", Type: "choice", Prompt: "Окно [l, r] - самое длинное хорошее окно с концом r, и любая его часть хорошая. Почему к ответу прибавляют r − l + 1?", Options: []string{"так считают только самое длинное окно с концом r", "так считают окна, которые кончатся позже", "хороши все окна [i, r] для i от l до r, их r − l + 1"}, Answer: 2, Explain: "Часть хорошего окна хорошая, поэтому каждое начало i от l до r даёт хорошее окно с концом r. Их ровно r − l + 1. Только самое длинное окно дало бы +1 на шаг, а не r − l + 1: для nums = [10, 5, 2, 6], k = 100 ответ стал бы 4 вместо 8. Окна с другим концом считают на своих шагах r."},
				{ID: "sw3-at-least", Type: "short", Prompt: "LeetCode 1358: s = \"aabcb\". Сколько подстрок содержат все три буквы a, b и c?", Answer: "4", Explain: "Буква c стоит только на месте 3, поэтому конец не левее 3. Начало не правее 1, чтобы взять a. Начала 0 и 1, концы 3 и 4: 2 · 2 = 4. Код раздела 4 прибавляет 0, 0, 0, 2, 2. Ответ 2 получится, если считать только окна с концом 3."},
				{ID: "sw3-exact", Type: "short", Prompt: "LeetCode 992: nums = [1, 2, 2, 1], k = 2. Сколько окон ровно с двумя разными числами?", Answer: "5", Explain: "Не больше двух разных - все 10 окон. Не больше одного: [1], [2], [2], [2, 2], [1], всего 5. Ответ 10 − 5 = 5: окна [1, 2], [1, 2, 2], [2, 2, 1], [2, 1], [1, 2, 2, 1]. Ответ 10 - это только первый проход."},
				{ID: "sw3-ski", Type: "short", Prompt: "Codeforces 1840C: a = [1, 2, 3, 9, 1, 1], q = 3, k = 2. Сколько вариантов отдыха?", Answer: "4", Explain: "Дни с a ≤ 3 идут блоками [1, 2, 3] и [1, 1]. В первом блоке L = 1, 2, 3 даёт max(0, L − 1): 0, 1, 2. Во втором блоке 0 и 1. Всего 4. Ответ 9 получится без условия «не короче k»: тогда считаются и окна из одного дня."},
				{ID: "sw3-bound", Type: "code", Prompt: "Codeforces 1736C1: места в Subarray считают с 1. Что поставить вместо ___?", Language: "python", Code: `l = 0
res = 0
for r in range(n):
    l = max(l, ___)
    res += r - l + 1`, Options: []string{"r - a[r] + 1", "a[r] - r", "r - a[r]"}, Answer: 0, Explain: "С началом l элемент a[r] стоит на месте r − l + 1. Условие a[r] ≥ r − l + 1 даёт l ≥ r − a[r] + 1. Для a = [3, 1, 3, 2] границы −2, 1, 0, 2, и код считает 1 + 1 + 2 + 2 = 6. С r − a[r] граница на 1 левее: при a = [1, 1] код возьмёт окно [1, 1], где второй элемент меньше 2."},
				{ID: "sw3-part", Type: "choice", Prompt: "nums = [1, 2, 1], k = 2. Окно с номерами 0..2 содержит ровно 2 разных числа. Его часть с номером 2 - только число 1. Что это значит?", Options: []string{"окно с номерами 1..2 тоже плохое", "часть хорошего окна может быть плохой, поэтому прибавлять r − l + 1 нельзя", "хватит одного прохода: сжимать окно, пока в нём больше K разных чисел, и прибавлять r − l + 1"}, Answer: 1, Explain: "Правило раздела 1 требует, чтобы каждая часть хорошего окна была хорошей. Для «ровно K» это неверно: часть [1] плохая. Окно с номерами 1..2 - это [2, 1], в нём 2 разных числа, оно хорошее. Один проход со сжатием при «больше K» считает окна, где не больше K разных: для [1, 2, 1] он даст 6, а верный ответ 3. Поэтому считаем разность двух проходов «не больше»."},
				{ID: "sw3-big", Type: "choice", Prompt: "Codeforces 1840C: n = 150 000, все дни подходят, k = 1. Каков ответ?", Options: []string{"150 000", "11 250 075 000", "22 500 000 000"}, Answer: 1, Explain: "Каждый шаг r прибавляет r + 1, сумма 1 + 2 + ... + n = n(n + 1)/2 = 11 250 075 000. Это больше 2 147 483 647, предела 32-битного целого. В Python число растёт само. 150 000 - это только окна из одного дня, а 22 500 000 000 - это n² без деления на 2."},
			},
			Links: []sets.Link{
				{Title: "Подсчёт окон: Subarray с произведением меньше k", URL: "https://www.geeksforgeeks.org/dsa/number-subarrays-product-less-k/", Source: "GeeksforGeeks (англ.)"},
				{Title: "Ровно K разных: разность двух подсчётов", URL: "https://www.geeksforgeeks.org/dsa/count-of-subarrays-having-exactly-k-distinct-elements/", Source: "GeeksforGeeks (англ.)"},
				{Title: "Subarrays with K Different Integers", URL: "https://takeuforward.org/blogs/data-structure-and-algorithm/subarrays-with-k-different-integers", Source: "takeUforward (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "sw3-count",
					Title: "Окна с концом r: nums = [10, 5, 2, 6], k = 100",
					Frames: []sets.FigureFrame{
						{Caption: "r = 0: окно [10], произведение 10 < 100. Прибавили 1.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{10, 5, 2, 6}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "l", Cell: 0}, {Name: "r", Cell: 0}}, Window: &[2]int{0, 0}},
							{Label: "произведение", Cells: []any{10}},
							{Label: "прибавили", Cells: []any{1}},
							{Label: "ответ", Cells: []any{1}},
						}},
						{Caption: "r = 1: произведение 50. Хорошие окна [5] и [10, 5]: прибавили 2.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{10, 5, 2, 6}, Marked: []int{1}, Pointers: &sets.Pointers{{Name: "l", Cell: 0}, {Name: "r", Cell: 1}}, Window: &[2]int{0, 1}},
							{Label: "произведение", Cells: []any{50}},
							{Label: "прибавили", Cells: []any{2}},
							{Label: "ответ", Cells: []any{3}},
						}},
						{Caption: "r = 2: произведение 100 не меньше k. Убрали 10, l = 1. Окна [2] и [5, 2]: прибавили 2.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{10, 5, 2, 6}, Marked: []int{2}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "l", Cell: 1}, {Name: "r", Cell: 2}}, Window: &[2]int{1, 2}},
							{Label: "произведение", Cells: []any{10}},
							{Label: "прибавили", Cells: []any{2}},
							{Label: "ответ", Cells: []any{5}},
						}},
						{Caption: "r = 3: произведение 60. Окна [6], [2, 6], [5, 2, 6]: прибавили 3. Ответ 8.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{10, 5, 2, 6}, Marked: []int{3}, Pointers: &sets.Pointers{{Name: "l", Cell: 1}, {Name: "r", Cell: 3}}, Window: &[2]int{1, 3}},
							{Label: "произведение", Cells: []any{60}},
							{Label: "прибавили", Cells: []any{3}},
							{Label: "ответ", Cells: []any{8}},
						}},
					},
				},
				{
					ID:    "sw3-at-least",
					Title: "Хотя бы по одной букве: s = \"abcabc\"",
					Frames: []sets.FigureFrame{
						{Caption: "r = 2: окно abc хорошее, while убирает a, l = 1. Хорошее начало одно: 0. Прибавили 1.", Rows: []sets.FigureRow{
							{Label: "s", Cells: []any{"a", "b", "c", "a", "b", "c"}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "l", Cell: 1}, {Name: "r", Cell: 2}}, Window: &[2]int{1, 2}},
							{Label: "a b c", Cells: []any{0, 1, 1}},
							{Label: "прибавили", Cells: []any{1}},
							{Label: "ответ", Cells: []any{1}},
						}},
						{Caption: "r = 3: while убирает b, l = 2. Хорошие начала 0 и 1. Прибавили 2.", Rows: []sets.FigureRow{
							{Label: "s", Cells: []any{"a", "b", "c", "a", "b", "c"}, Marked: []int{0, 1}, Pointers: &sets.Pointers{{Name: "l", Cell: 2}, {Name: "r", Cell: 3}}, Window: &[2]int{2, 3}},
							{Label: "a b c", Cells: []any{1, 0, 1}},
							{Label: "прибавили", Cells: []any{2}},
							{Label: "ответ", Cells: []any{3}},
						}},
						{Caption: "r = 4: while убирает c, l = 3. Хорошие начала 0, 1, 2. Прибавили 3.", Rows: []sets.FigureRow{
							{Label: "s", Cells: []any{"a", "b", "c", "a", "b", "c"}, Marked: []int{0, 1, 2}, Pointers: &sets.Pointers{{Name: "l", Cell: 3}, {Name: "r", Cell: 4}}, Window: &[2]int{3, 4}},
							{Label: "a b c", Cells: []any{1, 1, 0}},
							{Label: "прибавили", Cells: []any{3}},
							{Label: "ответ", Cells: []any{6}},
						}},
						{Caption: "r = 5: while убирает a, l = 4. Прибавили 4. Ответ 1 + 2 + 3 + 4 = 10.", Rows: []sets.FigureRow{
							{Label: "s", Cells: []any{"a", "b", "c", "a", "b", "c"}, Marked: []int{0, 1, 2, 3}, Pointers: &sets.Pointers{{Name: "l", Cell: 4}, {Name: "r", Cell: 5}}, Window: &[2]int{4, 5}},
							{Label: "a b c", Cells: []any{0, 1, 1}},
							{Label: "прибавили", Cells: []any{4}},
							{Label: "ответ", Cells: []any{10}},
						}},
					},
				},
				{
					ID:    "sw3-exact",
					Title: "Ровно 2 разных: nums = [1, 2, 1, 2, 3]",
					Frames: []sets.FigureFrame{
						{Caption: "r = 0: оба окна [1]. Прибавили 1 и 1, разность 0.", Rows: []sets.FigureRow{
							{Label: "не больше 2", Cells: []any{1, 2, 1, 2, 3}, Marked: []int{0}, Window: &[2]int{0, 0}},
							{Label: "не больше 1", Cells: []any{1, 2, 1, 2, 3}, Marked: []int{0}, Window: &[2]int{0, 0}},
							{Label: "суммы", Cells: []any{1, 1}},
							{Label: "ровно 2", Cells: []any{0}},
						}},
						{Caption: "r = 1: окно [1, 2] даёт 2, окно [2] даёт 1. Разность сумм 3 − 2 = 1.", Rows: []sets.FigureRow{
							{Label: "не больше 2", Cells: []any{1, 2, 1, 2, 3}, Marked: []int{1}, Window: &[2]int{0, 1}},
							{Label: "не больше 1", Cells: []any{1, 2, 1, 2, 3}, Marked: []int{1}, Dimmed: []int{0}, Window: &[2]int{1, 1}},
							{Label: "суммы", Cells: []any{3, 2}},
							{Label: "ровно 2", Cells: []any{1}},
						}},
						{Caption: "r = 2: прибавили 3 и 1. Суммы 6 и 3.", Rows: []sets.FigureRow{
							{Label: "не больше 2", Cells: []any{1, 2, 1, 2, 3}, Marked: []int{2}, Window: &[2]int{0, 2}},
							{Label: "не больше 1", Cells: []any{1, 2, 1, 2, 3}, Marked: []int{2}, Dimmed: []int{1}, Window: &[2]int{2, 2}},
							{Label: "суммы", Cells: []any{6, 3}},
							{Label: "ровно 2", Cells: []any{3}},
						}},
						{Caption: "r = 3: прибавили 4 и 1. Суммы 10 и 4.", Rows: []sets.FigureRow{
							{Label: "не больше 2", Cells: []any{1, 2, 1, 2, 3}, Marked: []int{3}, Window: &[2]int{0, 3}},
							{Label: "не больше 1", Cells: []any{1, 2, 1, 2, 3}, Marked: []int{3}, Dimmed: []int{2}, Window: &[2]int{3, 3}},
							{Label: "суммы", Cells: []any{10, 4}},
							{Label: "ровно 2", Cells: []any{6}},
						}},
						{Caption: "r = 4: число 3 даёт третье разное значение, первое окно сжимается до номера 3. Прибавили 2 и 1. Ответ 12 − 5 = 7.", Rows: []sets.FigureRow{
							{Label: "не больше 2", Cells: []any{1, 2, 1, 2, 3}, Marked: []int{4}, Dimmed: []int{0, 1, 2}, Window: &[2]int{3, 4}},
							{Label: "не больше 1", Cells: []any{1, 2, 1, 2, 3}, Marked: []int{4}, Dimmed: []int{3}, Window: &[2]int{4, 4}},
							{Label: "суммы", Cells: []any{12, 5}},
							{Label: "ровно 2", Cells: []any{7}},
						}},
					},
				},
			},
		},
		{
			Title:   "Monotonic Deque: Prefix Sum, Dynamic Programming и круг",
			Summary: "Monotonic Deque из главы о максимуме окна работает с любыми прошлыми значениями: с Prefix Sum и со значениями Dynamic Programming. Окно задаёт, какие прошлые номера ещё можно взять, и может идти по кругу. Deque хранит только лучшие номера.",
			Lesson:  lesson4,
			Tasks: []sets.Task{
				{
					URL:        "https://codeforces.com/problemset/problem/1744/C",
					Title:      "Traffic Light",
					Difficulty: "easy",
					Topics:     []string{"String", "Sliding Window"},
					Note:       "Codeforces, рейтинг 1000. Удвоенная строка и проход справа налево с номером ближайшей g. n - длина строки набора. Время: O(n) на набор, O(Σn) на весь вход. Память: O(n) дополнительно для удвоенной строки. Для rggry и r ожидания 1 и 3, ответ 3. Ходьба от каждого места c до g: O(n²).",
					Statement: `<p>Вы стоите на необычном перекрёстке со странным светофором. У светофора три цвета: красный (r), жёлтый (y) и зелёный (g). Известно, что светофор повторяет свои цвета каждые <code>n</code> секунд, и в <code>i</code>-ю секунду горит цвет <code>s<sub>i</sub></code>.</p><p>Так порядок цветов задаёт строка. Например, если <code>s = "rggry"</code>, светофор работает так: красный, зелёный, зелёный, красный, жёлтый, красный, зелёный, зелёный, красный, жёлтый, ... и так далее.</p><p>Строже: дана строка <code>s<sub>1</sub>, s<sub>2</sub>, ..., s<sub>n</sub></code> длины <code>n</code>. В первую секунду горит цвет <code>s<sub>1</sub></code>, во вторую - <code>s<sub>2</sub></code>, ..., в <code>n</code>-ю секунду горит цвет <code>s<sub>n</sub></code>, в <code>(n + 1)</code>-ю секунду снова горит цвет <code>s<sub>1</sub></code>, и так далее.</p><p>Вам нужно перейти дорогу, а это можно сделать, только когда горит зелёный.</p><p>Вы знаете, какой цвет горит сейчас, но не знаете, какая сейчас секунда. Найдите наименьшее время, через которое вы гарантированно перейдёте дорогу.</p><p>Считайте, что дорогу вы переходите мгновенно.</p><p>Например, при <code>s = "rggry"</code> и текущем цвете r есть два варианта: зелёный загорится через <code>1</code> секунду или через <code>3</code>. Поэтому ответ равен <code>3</code>: через столько секунд вы гарантированно перейдёте дорогу, если сейчас горит r.</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>): количество наборов входных данных.</p><p>Дальше идут описания наборов.</p><p>В первой строке набора записаны целое число <code>n</code> и символ <code>c</code> (<code>1 ≤ n ≤ 2 · 10<sup>5</sup></code>, <code>c</code> - один из цветов светофора r, y или g): длина строки <code>s</code> и текущий цвет светофора.</p><p>Во второй строке набора записана строка <code>s</code> длины <code>n</code> из букв r, y и g.</p><p>Гарантируется, что символ g есть в строке <code>s</code> и символ <code>c</code> есть в строке <code>s</code>.</p><p>Гарантируется, что сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого набора выведите наименьшее число секунд, через которое вы гарантированно перейдёте дорогу.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>6
5 r
rggry
1 g
g
3 r
rrg
5 y
yrrgy
7 r
rgrgyrg
9 y
rrrgyyygy</pre><p><strong>Выходные данные:</strong></p><pre>3
0
2
4
1
4</pre><h3>Примечание</h3><p>Первый набор разобран в условии.</p><p>Во втором наборе горит зелёный, поэтому дорогу можно перейти сразу.</p><p>В третьем наборе, если красный горел во вторую секунду, зелёный придётся ждать одну секунду. Если красный горел в первую секунду, зелёный придётся ждать две секунды.</p><p>В четвёртом наборе дольше всего придётся ждать зелёный, если начать ждать в пятую секунду.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>1 ≤ n ≤ 2 · 10<sup>5</sup></code></li><li><code>c</code> - одна из букв r, y и g</li><li><code>s</code> состоит из <code>n</code> букв r, y и g; в <code>s</code> есть g и <code>c</code></li><li>сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code></li><li>Ограничение по времени: 1 секунда на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "6\n5 r\nrggry\n1 g\ng\n3 r\nrrg\n5 y\nyrrgy\n7 r\nrgrgyrg\n9 y\nrrrgyyygy\n", Expected: text("3\n0\n2\n4\n1\n4\n")},
						{Input: "1\n2 r\nrg\n", Expected: text("1\n")},
						{Input: "1\n4 y\ngyyy\n", Expected: text("3\n")},
						{Input: "1\n3 r\nrgr\n", Expected: text("2\n")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/maximum-sum-circular-subarray/",
					Number:     number(918),
					Title:      "Maximum Sum Circular Subarray",
					Difficulty: "medium",
					Topics:     []string{"Array", "Prefix Sum", "Monotonic Queue", "Dynamic Programming"},
					Note:       "Prefix Sum удвоенного Array и Monotonic Deque для минимума p[i] в окне j − n..j − 1. n = len(nums). Время: O(n). Память: O(n) дополнительно: p из 2n + 1 чисел и Deque. Для [5, −3, 5] ответ 10. Kadane's Algorithm даёт O(n) и O(1) памяти, если отдельно проверить случай, когда все числа отрицательные. Перебор начала и длины: O(n²).",
					Statement: `<p>Дан <strong>круговой Array</strong> целых чисел <code>nums</code> длины <code>n</code>. Верните <em>наибольшую возможную сумму непустого <strong>Subarray</strong></em> <code>nums</code>.</p><p><strong>Круговой Array</strong> - это Array, у которого конец соединён с началом. Строже: следующий элемент после <code>nums[i]</code> - это <code>nums[(i + 1) % n]</code>, а предыдущий элемент перед <code>nums[i]</code> - это <code>nums[(i − 1 + n) % n]</code>.</p><p><strong>Subarray</strong> может взять каждый элемент Array <code>nums</code> не больше одного раза. Строже: для Subarray <code>nums[i], nums[i + 1], ..., nums[j]</code> нет таких <code>i ≤ k1</code>, <code>k2 ≤ j</code>, что <code>k1 % n == k2 % n</code>.</p><p>Метод класса <code>Solution</code>: <code>def maxSubarraySumCircular(self, nums: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,-2,3,-2]
<strong>Выход:</strong> 3
<strong>Пояснение:</strong> Subarray [3] имеет наибольшую сумму 3.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [5,-3,5]
<strong>Выход:</strong> 10
<strong>Пояснение:</strong> Subarray [5,5] имеет наибольшую сумму 5 + 5 = 10.</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> nums = [-3,-2,-3]
<strong>Выход:</strong> -2
<strong>Пояснение:</strong> Subarray [-2] имеет наибольшую сумму -2.</pre><h3>Ограничения</h3><ul><li><code>n == nums.length</code></li><li><code>1 ≤ n ≤ 3 · 10<sup>4</sup></code></li><li><code>−3 · 10<sup>4</sup> ≤ nums[i] ≤ 3 · 10<sup>4</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "maxSubarraySumCircular"}, Cases: []sets.TestCase{
						{Input: "[[1,-2,3,-2]]", Expected: text("3")},
						{Input: "[[5,-3,5]]", Expected: text("10")},
						{Input: "[[-3,-2,-3]]", Expected: text("-2")},
						{Input: "[[-30000]]", Expected: text("-30000")},
						{Input: "[[2,-10,3,-10,2]]", Expected: text("4")},
						{Input: "[[1,1,1]]", Expected: text("3")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/jump-game-vi/",
					Number:     number(1696),
					Title:      "Jump Game VI",
					Difficulty: "medium",
					Topics:     []string{"Array", "Dynamic Programming", "Monotonic Queue"},
					Note:       "dp[i] = nums[i] + max(dp[i − k..i − 1]); максимум окна даёт Monotonic Deque. n = len(nums). Время: O(n). Память: O(n) для dp и O(min(n, k + 1)) для Deque. Для [1, −1, −2, 4, −7, 3] и k = 2 dp = [1, 0, −1, 4, −3, 7], ответ 7. Перебор k прошлых dp: O(nk).",
					Statement: `<p>Дан Array целых чисел <code>nums</code> с <strong>индексами от 0</strong> и целое число <code>k</code>.</p><p>Сначала вы стоите на индексе <code>0</code>. За один ход можно прыгнуть вперёд не больше чем на <code>k</code> шагов и не выйти за границы Array. То есть с индекса <code>i</code> можно прыгнуть на любой индекс из отрезка <code>[i + 1, min(n − 1, i + k)]</code> <strong>включительно</strong>.</p><p>Вы хотите попасть на последний индекс Array (индекс <code>n − 1</code>). Ваш <strong>счёт</strong> - это <strong>сумма</strong> всех <code>nums[j]</code> для каждого индекса <code>j</code>, на котором вы побывали.</p><p>Верните <em><strong>наибольший счёт</strong>, который можно получить</em>.</p><p>Метод класса <code>Solution</code>: <code>def maxResult(self, nums: list[int], k: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,-1,-2,4,-7,3], k = 2
<strong>Выход:</strong> 7
<strong>Пояснение:</strong> Выберите индексы 0, 1, 3, 5: прыжки дают Subsequence [1,-1,4,3]. Сумма равна 7.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [10,-5,-2,4,0,3], k = 3
<strong>Выход:</strong> 17
<strong>Пояснение:</strong> Выберите индексы 0, 3, 5: прыжки дают Subsequence [10,4,3]. Сумма равна 17.</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> nums = [1,-5,-20,4,-1,3,-6,-3], k = 2
<strong>Выход:</strong> 0</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length, k ≤ 10<sup>5</sup></code></li><li><code>−10<sup>4</sup> ≤ nums[i] ≤ 10<sup>4</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "maxResult"}, Cases: []sets.TestCase{
						{Input: "[[1,-1,-2,4,-7,3],2]", Expected: text("7")},
						{Input: "[[10,-5,-2,4,0,3],3]", Expected: text("17")},
						{Input: "[[1,-5,-20,4,-1,3,-6,-3],2]", Expected: text("0")},
						{Input: "[[5],1]", Expected: text("5")},
						{Input: "[[1,-1,-1,-1,1],1]", Expected: text("-1")},
						{Input: "[[1,-5,-5,-5,2],100]", Expected: text("3")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1941/E",
					Title:      "Rudolf and k Bridges",
					Difficulty: "medium",
					Topics:     []string{"Array", "Dynamic Programming", "Monotonic Queue", "Sliding Window"},
					Note:       "Codeforces, рейтинг 1600. Для строки dp[j] = a[j] + 1 + min(dp[j − d − 1..j − 1]) с Monotonic Deque, потом окно из k соседних строк. Время: O(n · m) на набор, O(Σ n · m) на весь вход. Память: O(m) для одной строки и O(n) для цен строк. В первом примере цены строк 8, 4 и 8, при k = 1 ответ 4. Перебор прошлых опор: O(m · d) на строку.",
					Statement: `<p>Бернард любит ходить в гости к Рудольфу, но всегда опаздывает. Дело в том, что Бернарду нужно переправляться через реку на пароме. Рудольф решил помочь другу решить эту задачу.</p><p>Река - это клетчатое поле из <code>n</code> строк и <code>m</code> столбцов. На пересечении <code>i</code>-й строки и <code>j</code>-го столбца записано число <code>a<sub>i,j</sub></code>: глубина в этой клетке. Все клетки первого и последнего столбцов - берега реки, поэтому их глубина равна <code>0</code>.</p><p>Рудольф может выбрать строку <code>(i, 1), (i, 2), ..., (i, m)</code> и построить над ней мост. В каждой клетке строки он может поставить опору моста. Опора в клетке <code>(i, j)</code> стоит <code>a<sub>i,j</sub> + 1</code>. Опоры нужно поставить так, чтобы выполнялись условия:</p><ol><li>опора стоит в клетке <code>(i, 1)</code>;</li><li>опора стоит в клетке <code>(i, m)</code>;</li><li>расстояние между любыми двумя соседними опорами не больше <code>d</code>. Расстояние между опорами <code>(i, j<sub>1</sub>)</code> и <code>(i, j<sub>2</sub>)</code> равно <code>|j<sub>1</sub> − j<sub>2</sub>| − 1</code>.</li></ol><p>Строить только один мост скучно. Поэтому Рудольф решил построить <code>k</code> мостов на соседних строках реки: выбрать некоторое <code>i</code> (<code>1 ≤ i ≤ n − k + 1</code>) и независимо построить мост на каждой из строк <code>i, i + 1, ..., i + k − 1</code>. Помогите Рудольфу сделать общую стоимость установки опор наименьшей.</p><h3>Входные данные</h3><p>В первой строке записано одно целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>3</sup></code>): количество наборов входных данных. Дальше идут описания наборов.</p><p>В первой строке набора записаны четыре целых числа <code>n</code>, <code>m</code>, <code>k</code> и <code>d</code> (<code>1 ≤ k ≤ n ≤ 100</code>, <code>3 ≤ m ≤ 2 · 10<sup>5</sup></code>, <code>1 ≤ d ≤ m</code>): число строк и число столбцов поля, число мостов и наибольшее расстояние между опорами.</p><p>Дальше идут <code>n</code> строк. В <code>i</code>-й строке записаны <code>m</code> целых чисел <code>a<sub>i,j</sub></code> (<code>0 ≤ a<sub>i,j</sub> ≤ 10<sup>6</sup></code>, <code>a<sub>i,1</sub> = a<sub>i,m</sub> = 0</code>): глубины клеток реки.</p><p>Гарантируется, что сумма <code>n · m</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого набора выведите одно число: наименьшую общую стоимость установки опор.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>5
3 11 1 4
0 1 2 3 4 5 4 3 2 1 0
0 1 2 3 2 1 2 3 3 2 0
0 1 2 3 5 5 5 5 5 2 0
4 4 2 1
0 3 3 0
0 2 1 0
0 1 2 0
0 3 3 0
4 5 2 5
0 1 1 1 0
0 2 2 2 0
0 2 1 1 0
0 3 2 1 0
1 8 1 1
0 10 4 8 4 4 2 0
4 5 3 2
0 8 4 4 0
0 3 4 8 0
0 8 1 10 0
0 10 1 5 0</pre><p><strong>Выходные данные:</strong></p><pre>4
8
4
15
14</pre><h3>Примечание</h3><p>В первом наборе выгоднее всего построить мост на второй строке.</p><p>Во втором наборе выгоднее всего построить мосты на второй и третьей строках. Опоры встанут в клетки <code>(2, 3)</code>, <code>(3, 2)</code> и на берегах реки.</p><p>В третьем наборе опоры можно поставить только на берегах реки.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>3</sup></code></li><li><code>1 ≤ k ≤ n ≤ 100</code></li><li><code>3 ≤ m ≤ 2 · 10<sup>5</sup></code></li><li><code>1 ≤ d ≤ m</code></li><li><code>0 ≤ a<sub>i,j</sub> ≤ 10<sup>6</sup></code>, <code>a<sub>i,1</sub> = a<sub>i,m</sub> = 0</code></li><li>сумма <code>n · m</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "5\n3 11 1 4\n0 1 2 3 4 5 4 3 2 1 0\n0 1 2 3 2 1 2 3 3 2 0\n0 1 2 3 5 5 5 5 5 2 0\n4 4 2 1\n0 3 3 0\n0 2 1 0\n0 1 2 0\n0 3 3 0\n4 5 2 5\n0 1 1 1 0\n0 2 2 2 0\n0 2 1 1 0\n0 3 2 1 0\n1 8 1 1\n0 10 4 8 4 4 2 0\n4 5 3 2\n0 8 4 4 0\n0 3 4 8 0\n0 8 1 10 0\n0 10 1 5 0\n", Expected: text("4\n8\n4\n15\n14\n")},
						{Input: "1\n1 3 1 1\n0 5 0\n", Expected: text("2\n")},
						{Input: "1\n1 4 1 1\n0 5 7 0\n", Expected: text("8\n")},
						{Input: "1\n3 5 2 1\n0 4 1 4 0\n0 1 9 1 0\n0 5 5 5 0\n", Expected: text("10\n")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/shortest-subarray-with-sum-at-least-k/",
					Number:     number(862),
					Title:      "Shortest Subarray with Sum at Least K",
					Difficulty: "hard",
					Topics:     []string{"Array", "Prefix Sum", "Monotonic Queue", "Sliding Window"},
					Note:       "Prefix Sum и Monotonic Deque с возрастающими p: спереди забираем начала с p[j] − p[i] ≥ k, сзади удаляем p ≥ p[j]. n = len(nums). Время: O(n). Память: O(n) дополнительно. Для [2, −1, 2] и k = 3 ответ 3. Шаблон главы 2 ошибается из-за отрицательных чисел; перебор: O(n²).",
					Statement: `<p>Дан Array целых чисел <code>nums</code> и целое число <code>k</code>. Верните <em>длину самого короткого непустого <strong>Subarray</strong></em> <code>nums</code><em> с суммой не меньше </em><code>k</code>. Если такого <strong>Subarray</strong> нет, верните <code>-1</code>.</p><p><strong>Subarray</strong> - это <strong>непрерывная</strong> часть Array.</p><p>Метод класса <code>Solution</code>: <code>def shortestSubarray(self, nums: list[int], k: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1], k = 1
<strong>Выход:</strong> 1</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [1,2], k = 4
<strong>Выход:</strong> -1</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> nums = [2,-1,2], k = 3
<strong>Выход:</strong> 3</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 10<sup>5</sup></code></li><li><code>−10<sup>5</sup> ≤ nums[i] ≤ 10<sup>5</sup></code></li><li><code>1 ≤ k ≤ 10<sup>9</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "shortestSubarray"}, Cases: []sets.TestCase{
						{Input: "[[1],1]", Expected: text("1")},
						{Input: "[[1,2],4]", Expected: text("-1")},
						{Input: "[[2,-1,2],3]", Expected: text("3")},
						{Input: "[[-1,-1],1]", Expected: text("-1")},
						{Input: "[[84,-37,32,40,95],167]", Expected: text("3")},
						{Input: "[[17,85,93,-45,-21],150]", Expected: text("2")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "sw4-window", Type: "choice", Prompt: "LeetCode 918: в удвоенном Array ищем Subarray с концом j и суммой p[j] − p[i]. Какие начала i разрешены, чтобы каждый элемент круга был взят не больше одного раза?", Options: []string{"0 ≤ i ≤ j − 1", "j − n ≤ i ≤ j − 1", "j − n + 1 ≤ i ≤ j"}, Answer: 1, Explain: "Subarray с номерами i..j − 1 имеет длину j − i, на круге она не больше n: i ≥ j − n. Длина хотя бы 1: i ≤ j − 1. Вариант 0 ≤ i разрешит длину до 2n, и для [1, 1, 1] ответ станет 6 вместо 3. Вариант с i = j даёт пустой Subarray."},
				{ID: "sw4-prefix", Type: "short", Prompt: "nums = [3, −1, 2], p[0] = 0. Чему равна сумма nums[1..2], посчитанная через p?", Answer: "1", Explain: "p = [0, 3, 2, 4]. Сумма nums[1..2] = p[3] − p[1] = 4 − 3 = 1, то есть −1 + 2. p[2] − p[1] = −1 - это только nums[1]. Одна разность стоит O(1), а p строим один раз за O(n)."},
				{ID: "sw4-front", Type: "choice", Prompt: "LeetCode 862: почему номер i можно убрать спереди, когда p[j] − p[i] ≥ k?", Options: []string{"p[i] больше всех остальных в Deque", "после этого сумма окна станет отрицательной", "для более поздних j окно с началом i будет длиннее, а ответ для i уже записан"}, Answer: 2, Explain: "Длина j − i растёт вместе с j. Для начала i лучший конец - первый j с суммой ≥ k, его длину мы уже записали. В Deque p возрастают, поэтому p[i] спереди наименьший, а не наибольший. Каждый номер уходит спереди не больше одного раза: все удаления вместе стоят O(n)."},
				{ID: "sw4-dp", Type: "short", Prompt: "LeetCode 1696: nums = [2, −3, −1, 5], k = 2. Чему равен ответ?", Answer: "6", Explain: "dp[0] = 2. dp[1] = −3 + 2 = −1. dp[2] = −1 + max(2, −1) = 1. dp[3] = 5 + max(−1, 1) = 6. Путь 0, 2, 3: 2 − 1 + 5 = 6. Путь по всем местам даёт 3. Ответ 7 = 2 + 5 был бы при прыжке на 3 места, а k = 2."},
				{ID: "sw4-kadane", Type: "choice", Prompt: "LeetCode 918, формула max(лучший Subarray, total − худший Subarray) без проверки, total - сумма всех чисел. Что она даст для [−5, −1, −2]?", Options: []string{"0", "−1", "−8"}, Answer: 0, Explain: "total = −8, худший Subarray - весь Array с суммой −8. total − (−8) = 0, это пустой Subarray. Лучший обычный Subarray −1, и max(−1, 0) = 0. Верный ответ −1: если все числа отрицательные, берите лучший обычный Subarray. −8 - это сумма всего Array."},
				{ID: "sw4-bridge", Type: "short", Prompt: "Codeforces 1941E: строка [0, 3, 1, 2, 0], d = 1. Какова цена моста на этой строке?", Answer: "4", Explain: "dp[0] = 1. dp[1] = 1 + 4 = 5. dp[2] = min(1, 5) + 2 = 3. dp[3] = min(5, 3) + 3 = 6. dp[4] = min(3, 6) + 1 = 4. Опоры в столбцах 0, 2, 4. Если взять окно j − d..j − 1, опоры встанут во все столбцы: 1 + 4 + 2 + 3 + 1 = 11. Расстояние |j1 − j2| − 1 ≤ d разрешает пропустить d столбцов."},
				{ID: "sw4-sign", Type: "code", Prompt: "LeetCode 862: Deque хранит номера с возрастающими p. Какой знак поставить вместо ___?", Language: "python", Code: `for j in range(len(p)):
    while dq and p[j] - p[dq[0]] >= k:
        best = min(best, j - dq.popleft())
    while dq and p[dq[-1]] ___ p[j]:
        dq.pop()
    dq.append(j)`, Options: []string{">=", "<=", "<"}, Answer: 0, Explain: "Номер сзади с p ≥ p[j] хуже начала j: j ближе к будущим концам и даёт не меньшую сумму. После удаления p в Deque возрастают. Со знаком <= уйдут лучшие начала: для [2, −1, 2], k = 3 Deque потеряет номер 0, и код вернёт −1 вместо 3."},
				{ID: "sw4-circle", Type: "short", Prompt: "Codeforces 1744C: s = \"grr\", горит r. Через сколько секунд можно гарантированно перейти дорогу?", Answer: "2", Explain: "Удвоенная строка grrgrr. От места 1 до g на месте 3 две секунды, от места 2 одна секунда. Гарантия - наибольшее ожидание: 2. Без удвоения после места 2 нет g, и код не найдёт ожидание для r в конце строки."},
			},
			Links: []sets.Link{
				{Title: "Prefix Sum, минимум в окне и Deque", URL: "https://informatics.msk.ru/pluginfile.php/376779/mod_resource/content/0/Linear_algorithms.pdf", Source: "informatics.msk.ru"},
				{Title: "Kadane's Algorithm и Dynamic Programming", URL: "https://repetitor.1c.ru/informatics/algoritm-poiska-maksimalnoy-podposledovatelnosti/", Source: "1С:Репетитор"},
				{Title: "Subarray с наибольшей суммой через Prefix Sum", URL: "https://cp-algorithms.com/others/maximum_average_segment.html", Source: "cp-algorithms (англ.)"},
				{Title: "Maximum Subarray Problem", URL: "https://en.wikipedia.org/wiki/Maximum_subarray_problem", Source: "Википедия (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "sw4-circle",
					Title: "Круг через удвоение: nums = [5, −3, 5]",
					Frames: []sets.FigureFrame{
						{Caption: "j = 1: начало i = 0, сумма p[1] − p[0] = 5.", Rows: []sets.FigureRow{
							{Label: "удвоенный", Cells: []any{5, -3, 5, 5, -3, 5}, Marked: []int{0}},
							{Label: "p", Cells: []any{0, 5, 2, 7, 12, 9, 14}, Marked: []int{0, 1}, Pointers: &sets.Pointers{{Name: "j", Cell: 1}}, Window: &[2]int{0, 0}},
							{Label: "Deque номера", Cells: []any{0, 1}},
							{Label: "лучшая", Cells: []any{5}},
						}},
						{Caption: "j = 2: сумма 2 − 0 = 2. p[1] = 5 ≥ p[2] = 2, номер 1 уходит сзади.", Rows: []sets.FigureRow{
							{Label: "удвоенный", Cells: []any{5, -3, 5, 5, -3, 5}, Marked: []int{0, 1}},
							{Label: "p", Cells: []any{0, 5, 2, 7, 12, 9, 14}, Marked: []int{0, 2}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}, Window: &[2]int{0, 1}},
							{Label: "Deque номера", Cells: []any{0, 2}},
							{Label: "лучшая", Cells: []any{5}},
						}},
						{Caption: "j = 3: окно начал 0..2, минимум p[0] = 0, сумма 7.", Rows: []sets.FigureRow{
							{Label: "удвоенный", Cells: []any{5, -3, 5, 5, -3, 5}, Marked: []int{0, 1, 2}},
							{Label: "p", Cells: []any{0, 5, 2, 7, 12, 9, 14}, Marked: []int{0, 3}, Pointers: &sets.Pointers{{Name: "j", Cell: 3}}, Window: &[2]int{0, 2}},
							{Label: "Deque номера", Cells: []any{0, 2, 3}},
							{Label: "лучшая", Cells: []any{7}},
						}},
						{Caption: "j = 4: номер 0 вне окна 1..3 и уходит спереди. Минимум p[2] = 2, сумма 12 − 2 = 10: это 5 + 5 через конец круга.", Rows: []sets.FigureRow{
							{Label: "удвоенный", Cells: []any{5, -3, 5, 5, -3, 5}, Marked: []int{2, 3}},
							{Label: "p", Cells: []any{0, 5, 2, 7, 12, 9, 14}, Marked: []int{2, 4}, Pointers: &sets.Pointers{{Name: "j", Cell: 4}}, Window: &[2]int{1, 3}},
							{Label: "Deque номера", Cells: []any{2, 3, 4}},
							{Label: "лучшая", Cells: []any{10}},
						}},
						{Caption: "j = 5: сумма 9 − 2 = 7. p[4] = 12 ≥ 9, номер 4 уходит сзади.", Rows: []sets.FigureRow{
							{Label: "удвоенный", Cells: []any{5, -3, 5, 5, -3, 5}, Marked: []int{2, 3, 4}},
							{Label: "p", Cells: []any{0, 5, 2, 7, 12, 9, 14}, Marked: []int{2, 5}, Pointers: &sets.Pointers{{Name: "j", Cell: 5}}, Window: &[2]int{2, 4}},
							{Label: "Deque номера", Cells: []any{2, 3, 5}},
							{Label: "лучшая", Cells: []any{10}},
						}},
						{Caption: "j = 6: номер 2 вне окна 3..5. Минимум p[3] = 7, сумма 7. Ответ 10.", Rows: []sets.FigureRow{
							{Label: "удвоенный", Cells: []any{5, -3, 5, 5, -3, 5}, Marked: []int{3, 4, 5}},
							{Label: "p", Cells: []any{0, 5, 2, 7, 12, 9, 14}, Marked: []int{3, 6}, Pointers: &sets.Pointers{{Name: "j", Cell: 6}}, Window: &[2]int{3, 5}},
							{Label: "Deque номера", Cells: []any{3, 5, 6}},
							{Label: "лучшая", Cells: []any{10}},
						}},
					},
				},
				{
					ID:    "sw4-shortest",
					Title: "Самый короткий Subarray: nums = [2, −1, 2], k = 3",
					Frames: []sets.FigureFrame{
						{Caption: "j = 0: p[0] = 0, Deque [0].", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{2, -1, 2}},
							{Label: "p", Cells: []any{0, 2, 1, 3}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "j", Cell: 0}}},
							{Label: "Deque номера", Cells: []any{0}},
							{Label: "Deque p", Cells: []any{0}},
							{Label: "лучшая длина", Cells: []any{"нет"}},
						}},
						{Caption: "j = 1: p[1] − p[0] = 2 < 3. p растёт, номер 1 встаёт сзади.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{2, -1, 2}},
							{Label: "p", Cells: []any{0, 2, 1, 3}, Marked: []int{1}, Pointers: &sets.Pointers{{Name: "j", Cell: 1}}},
							{Label: "Deque номера", Cells: []any{0, 1}},
							{Label: "Deque p", Cells: []any{0, 2}},
							{Label: "лучшая длина", Cells: []any{"нет"}},
						}},
						{Caption: "j = 2: p[1] = 2 ≥ p[2] = 1. Номер 1 уходит сзади: начало 2 не хуже и ближе.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{2, -1, 2}},
							{Label: "p", Cells: []any{0, 2, 1, 3}, Marked: []int{2}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}},
							{Label: "Deque номера", Cells: []any{0, 2}},
							{Label: "Deque p", Cells: []any{0, 1}},
							{Label: "лучшая длина", Cells: []any{"нет"}},
						}},
						{Caption: "j = 3: p[3] − p[0] = 3 ≥ 3, длина 3, номер 0 уходит спереди. p[3] − p[2] = 2 < 3. Ответ 3.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{2, -1, 2}, Marked: []int{0, 1, 2}},
							{Label: "p", Cells: []any{0, 2, 1, 3}, Marked: []int{3}, Pointers: &sets.Pointers{{Name: "j", Cell: 3}}},
							{Label: "Deque номера", Cells: []any{2, 3}},
							{Label: "Deque p", Cells: []any{1, 3}},
							{Label: "лучшая длина", Cells: []any{3}},
						}},
					},
				},
				{
					ID:    "sw4-dp",
					Title: "Dynamic Programming с Deque: nums = [1, −1, −2, 4, −7, 3], k = 2",
					Frames: []sets.FigureFrame{
						{Caption: "i = 1: лучший прошлый dp[0] = 1, dp[1] = −1 + 1 = 0.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, -1, -2, 4, -7, 3}, Marked: []int{1}},
							{Label: "dp", Cells: []any{1, 0, "", "", "", ""}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "i", Cell: 1}}, Window: &[2]int{0, 0}},
							{Label: "Deque номера", Cells: []any{0, 1}},
						}},
						{Caption: "i = 2: окно dp[0..1], максимум 1, dp[2] = −2 + 1 = −1.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, -1, -2, 4, -7, 3}, Marked: []int{2}},
							{Label: "dp", Cells: []any{1, 0, -1, "", "", ""}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "i", Cell: 2}}, Window: &[2]int{0, 1}},
							{Label: "Deque номера", Cells: []any{0, 1, 2}},
						}},
						{Caption: "i = 3: номер 0 вне окна 1..2 и уходит спереди. dp[3] = 4 + 0 = 4. Номера 2 и 1 уходят сзади: 4 больше.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, -1, -2, 4, -7, 3}, Marked: []int{3}},
							{Label: "dp", Cells: []any{1, 0, -1, 4, "", ""}, Marked: []int{1}, Pointers: &sets.Pointers{{Name: "i", Cell: 3}}, Window: &[2]int{1, 2}},
							{Label: "Deque номера", Cells: []any{3}},
						}},
						{Caption: "i = 4: dp[4] = −7 + 4 = −3. Deque [3, 4].", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, -1, -2, 4, -7, 3}, Marked: []int{4}},
							{Label: "dp", Cells: []any{1, 0, -1, 4, -3, ""}, Marked: []int{3}, Pointers: &sets.Pointers{{Name: "i", Cell: 4}}, Window: &[2]int{2, 3}},
							{Label: "Deque номера", Cells: []any{3, 4}},
						}},
						{Caption: "i = 5: dp[5] = 3 + 4 = 7. Ответ 7: путь по местам 0, 1, 3, 5.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{1, -1, -2, 4, -7, 3}, Marked: []int{5}},
							{Label: "dp", Cells: []any{1, 0, -1, 4, -3, 7}, Marked: []int{3}, Pointers: &sets.Pointers{{Name: "i", Cell: 5}}, Window: &[2]int{3, 4}},
							{Label: "Deque номера", Cells: []any{5}},
						}},
					},
				},
			},
		},
	},
}
