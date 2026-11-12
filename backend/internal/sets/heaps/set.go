// Package heaps is the learning set "Heap (Priority Queue), Queue и Deque".
package heaps

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
	Slug:    "heaps-and-queues",
	Title:   "Heap (Priority Queue), Queue и Deque",
	Summary: "Queue, Deque и Heap (Priority Queue) выбирают элементы по разным правилам. Вы используете Simulation, Monotonic Queue, K-Way Merge и Greedy. Для каждого решения считаете операции, пиковую дополнительную память и размер ответа.",
	Topics:  []string{"Queue", "Deque", "Heap (Priority Queue)", "Sliding Window"},
	Chapters: []sets.Chapter{
		{
			Title:   "Queue и Deque",
			Summary: "Queue сохраняет порядок прихода, Deque поддерживает оба конца. Вы считаете создание контейнера, сдвиги, повороты и шаги Simulation. Формула, счётчики и Josephus Recurrence помогают сократить время и память.",
			Lesson:  lesson1,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/time-needed-to-buy-tickets/",
					Number:     number(2073),
					Title:      "Time Needed to Buy Tickets",
					Difficulty: "easy",
					Topics:     []string{"Array", "Queue", "Simulation"},
					Note:       "Simulation хранит номера людей в Deque. n - число людей, b - число продаж до последнего билета человека k. Время: O(n + b), включая создание Deque и b удалений. Память: O(n) дополнительно, ответ O(1). Для [3, 1, 2], k = 0: 6 продаж, 3 возврата, пик 3 номера. Linear Scan суммы min(t[i], T) при i ≤ k и min(t[i], T - 1) после k даёт O(n) времени и O(1) дополнительной памяти; T = t[k] до изменений.",
					Statement: `<p>В очереди за билетами стоят <code>n</code> человек. Человек <code>0</code> стоит в <strong>начале</strong> очереди, человек <code>n - 1</code> в <strong>конце</strong>.</p><p>Дан Array целых чисел <code>tickets</code> длины <code>n</code>, номера в нём <strong>начинаются с 0</strong>. Человек <code>i</code> хочет купить <code>tickets[i]</code> билетов.</p><p>Покупка одного билета занимает <strong>ровно 1 секунду</strong>. Человек покупает <strong>только 1 билет за раз</strong>. Чтобы купить ещё, он встаёт в <strong>конец</strong> очереди, и это происходит <strong>мгновенно</strong>. Если человеку больше не нужны билеты, он <strong>уходит</strong> из очереди.</p><p>Верните <strong>время</strong>, за которое человек, который <strong>в начале</strong> стоял на месте <strong>k</strong> (номера с 0), купит все свои билеты.</p><p>Метод класса <code>Solution</code>: <code>def timeRequiredToBuy(self, tickets: list[int], k: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> tickets = [2,3,2], k = 2
<strong>Выход:</strong> 6
<strong>Пояснение:</strong> человек k выделен жирным.
- Сначала очередь [2,3,<strong>2</strong>].
- Человек в начале купил билет, и через 1 секунду очередь [3,<strong>2</strong>,1].
- Дальше так же: через 2 секунды очередь [<strong>2</strong>,1,2].
- Через 3 секунды очередь [1,2,<strong>1</strong>].
- Через 4 секунды очередь [2,<strong>1</strong>]. Человек в начале ушёл из очереди.
- Через 5 секунд очередь [<strong>1</strong>,1].
- Через 6 секунд очередь [1]. Человек k купил все билеты, поэтому верните 6.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> tickets = [5,1,1,1], k = 0
<strong>Выход:</strong> 8
<strong>Пояснение:</strong>
- Сначала очередь [<strong>5</strong>,1,1,1].
- Человек в начале купил билет, и через 1 секунду очередь [1,1,1,<strong>4</strong>].
- Ещё через 3 секунды, через 4 секунды от начала, очередь [<strong>4</strong>].
- Ещё через 4 секунды, через 8 секунд от начала, очередь []. Человек k купил все билеты, поэтому верните 8.</pre><h3>Ограничения</h3><ul><li><code>n == tickets.length</code></li><li><code>1 ≤ n ≤ 100</code></li><li><code>1 ≤ tickets[i] ≤ 100</code></li><li><code>0 ≤ k &lt; n</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "timeRequiredToBuy"}, Cases: []sets.TestCase{
						{Input: "[[2,3,2],2]", Expected: text("6")},
						{Input: "[[5,1,1,1],0]", Expected: text("8")},
						{Input: "[[3,1,2],0]", Expected: text("6")},
						{Input: "[[1],0]", Expected: text("1")},
						{Input: "[[100,100],1]", Expected: text("200")},
						{Input: "[[2,6,3,4,5],2]", Expected: text("12")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/number-of-students-unable-to-eat-lunch/",
					Number:     number(1700),
					Title:      "Number of Students Unable to Eat Lunch",
					Difficulty: "easy",
					Topics:     []string{"Array", "Stack", "Queue", "Simulation"},
					Note:       "После полного круга отказов от верхнего бутерброда остановите Simulation; после выдачи обнулите счётчик отказов. n - число учеников. Время: O(n²) в худшем случае, перед каждой выдачей возможен полный проход. Память: O(n) для Deque, ответ O(1). Для students = [0, 0, 1], sandwiches = [1, 1, 0]: 5 проверок, пик 3 ученика, ответ 2. Два счётчика желаний дают O(n) времени и O(1) дополнительной памяти.",
					Statement: `<p>В школьной столовой на обед дают круглые и квадратные бутерброды. Они обозначены числами <code>0</code> и <code>1</code> соответственно. Все ученики стоят в очереди. Каждый ученик хочет либо квадратный, либо круглый бутерброд.</p><p>Бутербродов в столовой столько же, сколько учеников. Бутерброды лежат <strong>стопкой</strong>. На каждом шаге:</p><ul><li>если ученик в начале очереди <strong>хочет</strong> верхний бутерброд стопки, он <strong>берёт его</strong> и уходит из очереди;</li><li>иначе он <strong>не берёт</strong> бутерброд и встаёт в конец очереди.</li></ul><p>Так продолжается, пока ни один ученик в очереди не хочет верхний бутерброд. Эти ученики остаются без обеда.</p><p>Даны два Array целых чисел <code>students</code> и <code>sandwiches</code>. <code>sandwiches[i]</code> это вид <code>i</code>-го бутерброда в стопке (<code>i = 0</code> это верх стопки). <code>students[j]</code> это бутерброд, который хочет <code>j</code>-й ученик в начальной очереди (<code>j = 0</code> это начало очереди). Верните <em>число учеников, которые останутся без обеда</em>.</p><p>Метод класса <code>Solution</code>: <code>def countStudents(self, students: list[int], sandwiches: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> students = [1,1,0,0], sandwiches = [0,1,0,1]
<strong>Выход:</strong> 0
<strong>Пояснение:</strong>
- Ученик в начале не берёт верхний бутерброд и встаёт в конец: students = [1,0,0,1].
- Ученик в начале не берёт верхний бутерброд и встаёт в конец: students = [0,0,1,1].
- Ученик в начале берёт верхний бутерброд и уходит: students = [0,1,1], sandwiches = [1,0,1].
- Ученик в начале не берёт верхний бутерброд и встаёт в конец: students = [1,1,0].
- Ученик в начале берёт верхний бутерброд и уходит: students = [1,0], sandwiches = [0,1].
- Ученик в начале не берёт верхний бутерброд и встаёт в конец: students = [0,1].
- Ученик в начале берёт верхний бутерброд и уходит: students = [1], sandwiches = [1].
- Ученик в начале берёт верхний бутерброд и уходит: students = [], sandwiches = [].
Значит, все ученики пообедали.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> students = [1,1,1,0,0,1], sandwiches = [1,0,0,0,1,1]
<strong>Выход:</strong> 3</pre><h3>Ограничения</h3><ul><li><code>1 ≤ students.length, sandwiches.length ≤ 100</code></li><li><code>students.length == sandwiches.length</code></li><li><code>sandwiches[i]</code> равно <code>0</code> или <code>1</code></li><li><code>students[i]</code> равно <code>0</code> или <code>1</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "countStudents"}, Cases: []sets.TestCase{
						{Input: "[[1,1,0,0],[0,1,0,1]]", Expected: text("0")},
						{Input: "[[1,1,1,0,0,1],[1,0,0,0,1,1]]", Expected: text("3")},
						{Input: "[[1,0,0,0],[0,0,1,1]]", Expected: text("1")},
						{Input: "[[1],[0]]", Expected: text("1")},
						{Input: "[[0],[0]]", Expected: text("0")},
						{Input: "[[1,1,0,0],[0,0,0,1]]", Expected: text("2")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1579/E1",
					Title:      "Permutation Minimization by Deque",
					Difficulty: "easy",
					Topics:     []string{"Array", "Deque", "Greedy"},
					Note:       "Codeforces, рейтинг 1000. Greedy сравнивает новое число с d[0] и выбирает конец Deque. n - число элементов одного набора. Время: O(n), n - 1 сравнений и n добавлений. Память: O(1) состояния сверх Deque ответа из O(n) чисел; копия в list добавит O(n). Для [3, 1, 2, 4]: 3 сравнения, ответ [1, 3, 2, 4]. list.insert(0, x) может дать O(n²) сдвигов. Для нескольких наборов сложите время, память считайте по наибольшему набору при обработке по одному.",
					Statement: `<p>Дана Permutation <code>p</code> размера <code>n</code>. Permutation размера <code>n</code> - это Array из <code>n</code> чисел, в котором каждое целое число от <code>1</code> до <code>n</code> встречается ровно один раз. Например, <code>[1, 4, 3, 2]</code> и <code>[4, 2, 1, 3]</code> подходят, а <code>[1, 2, 4]</code> и <code>[1, 2, 2]</code> нет.</p><p>Возьмём пустой Deque. В Deque можно добавлять элементы и в начало, и в конец. Например, в Deque лежит <code>[1, 5, 2]</code>. Если добавить 4 в начало, получится <code>[4, 1, 5, 2]</code>, а если в конец, то <code>[1, 5, 2, 4]</code>.</p><p>Элементы Permutation по очереди кладут в пустой Deque, от <code>p<sub>1</sub></code> до <code>p<sub>n</sub></code>. Перед добавлением каждого элемента можно выбрать, положить его в начало или в конец.</p><p>Например, для <code>p = [3, 1, 2, 4]</code> возможны такие действия:</p><ol><li>положить 3 в конец Deque: <code>[3]</code>;</li><li>положить 1 в начало Deque: <code>[1, 3]</code>;</li><li>положить 2 в конец Deque: <code>[1, 3, 2]</code>;</li><li>положить 4 в конец Deque: <code>[1, 3, 2, 4]</code>.</li></ol><p>Найдите наименьшую последовательность в Lexicographic Order, которая может лежать в Deque после обработки всего входа.</p><p>В Lexicographic Order последовательность <code>[x<sub>1</sub>, x<sub>2</sub>, ..., x<sub>n</sub>]</code> меньше <code>[y<sub>1</sub>, y<sub>2</sub>, ..., y<sub>n</sub>]</code>, если есть такое <code>i ≤ n</code>, что <code>x<sub>1</sub> = y<sub>1</sub></code>, <code>x<sub>2</sub> = y<sub>2</sub></code>, ..., <code>x<sub>i−1</sub> = y<sub>i−1</sub></code> и <code>x<sub>i</sub> &lt; y<sub>i</sub></code>. Другими словами, начало у последовательностей одинаковое (может быть, пустое), а следующий элемент <code>x</code> строго меньше элемента <code>y</code>. Например, <code>[1, 3, 2, 4]</code> меньше <code>[1, 3, 4, 2]</code>: после общего начала <code>[1, 3]</code> в первой идёт 2, а во второй 4.</p><h3>Входные данные</h3><p>В первой строке записано целое число <code>t</code> (<code>1 ≤ t ≤ 1000</code>): количество наборов входных данных.</p><p>Дальше идут <code>t</code> наборов. В первой строке набора записано целое число <code>n</code> (<code>1 ≤ n ≤ 2 · 10<sup>5</sup></code>): размер Permutation. Во второй строке через пробел записаны <code>n</code> целых чисел <code>p<sub>i</sub></code> (<code>1 ≤ p<sub>i</sub> ≤ n</code>, все <code>p<sub>i</sub></code> разные): элементы Permutation.</p><p>Гарантируется, что сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Выведите <code>t</code> строк, по одной для каждого набора. В строке через пробел выведите <code>n</code> чисел: наименьшую последовательность в Lexicographic Order, которая может получиться в Deque.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>5
4
3 1 2 4
3
3 2 1
3
3 1 2
2
1 2
2
2 1</pre><p><strong>Выходные данные:</strong></p><pre>1 3 2 4
1 2 3
1 3 2
1 2
1 2</pre><p><strong>Пояснение:</strong> ответ <code>[1, 3, 2, 4]</code> первого набора получают действия из примера выше.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 1000</code></li><li><code>1 ≤ n ≤ 2 · 10<sup>5</sup></code></li><li><code>p</code> это Permutation чисел от 1 до <code>n</code></li><li>сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `5
4
3 1 2 4
3
3 2 1
3
3 1 2
2
1 2
2
2 1
`, Expected: text(`1 3 2 4
1 2 3
1 3 2
1 2
1 2
`)},
						{Input: `1
1
1
`, Expected: text(`1
`)},
						{Input: `1
5
1 2 3 4 5
`, Expected: text(`1 2 3 4 5
`)},
						{Input: `1
5
5 4 3 2 1
`, Expected: text(`1 2 3 4 5
`)},
						{Input: `1
6
4 6 2 5 1 3
`, Expected: text(`1 2 4 6 5 3
`)},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/find-the-winner-of-the-circular-game/",
					Number:     number(1823),
					Title:      "Find the Winner of the Circular Game",
					Difficulty: "medium",
					Topics:     []string{"Array", "Math", "Recursion", "Queue", "Simulation"},
					Note:       "Для n друзей Simulation делает rotate(-(k - 1)) и popleft. Пусть d_s = min((k - 1) mod s, s - ((k - 1) mod s)) при длине s. Время: O(n + ∑d_s), s от 2 до n, верхняя оценка O(n²). Память: O(n) для Deque, ответ O(1). При n = 6, k = 3: 5 удалений, сумма d_s = 7, победитель 1. Josephus Recurrence J(1) = 0, J(s) = (J(s - 1) + k) mod s даёт O(n) времени и O(1) памяти при итеративном расчёте; ответ J(n) + 1.",
					Statement: `<p>В игре участвуют <code>n</code> друзей. Они сидят по кругу, их номера от <code>1</code> до <code>n</code> <strong>по часовой стрелке</strong>. Точнее: если идти по часовой стрелке от друга <code>i</code>, следующим будет друг <code>i + 1</code> для <code>1 ≤ i &lt; n</code>, а после друга <code>n</code> идёт друг <code>1</code>.</p><p>Правила игры:</p><ol><li><strong>Начните</strong> с друга <code>1</code>.</li><li>Отсчитайте <code>k</code> друзей по часовой стрелке, <strong>включая</strong> друга, с которого начали. Счёт идёт по кругу, и одного друга можно посчитать несколько раз.</li><li>Последний посчитанный друг выходит из круга и проигрывает.</li><li>Если в круге больше одного друга, вернитесь к шагу <code>2</code>. Новый счёт <strong>начинается</strong> с друга, который сидит <strong>сразу по часовой стрелке</strong> после того, кто вышел.</li><li>Иначе последний друг в круге побеждает.</li></ol><p>Даны число друзей <code>n</code> и целое число <code>k</code>. Верните <em>номер победителя</em>.</p><p>Метод класса <code>Solution</code>: <code>def findTheWinner(self, n: int, k: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> n = 5, k = 2
<strong>Выход:</strong> 3
<strong>Пояснение:</strong> шаги игры:
1) Начинаем с друга 1.
2) Считаем 2 друзей по часовой стрелке: это друзья 1 и 2.
3) Друг 2 выходит из круга. Следующий счёт начинается с друга 3.
4) Считаем 2 друзей по часовой стрелке: это друзья 3 и 4.
5) Друг 4 выходит из круга. Следующий счёт начинается с друга 5.
6) Считаем 2 друзей по часовой стрелке: это друзья 5 и 1.
7) Друг 1 выходит из круга. Следующий счёт начинается с друга 3.
8) Считаем 2 друзей по часовой стрелке: это друзья 3 и 5.
9) Друг 5 выходит из круга. Остался только друг 3, он победил.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> n = 6, k = 5
<strong>Выход:</strong> 1
<strong>Пояснение:</strong> друзья выходят в таком порядке: 5, 4, 6, 2, 3. Победил друг 1.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ k ≤ n ≤ 500</code></li></ul><p><strong>Дополнительный вопрос:</strong> можете ли вы решить задачу за O(n) времени и с O(1) дополнительной памятью?</p>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "findTheWinner"}, Cases: []sets.TestCase{
						{Input: "[5,2]", Expected: text("3")},
						{Input: "[6,5]", Expected: text("1")},
						{Input: "[1,1]", Expected: text("1")},
						{Input: "[7,3]", Expected: text("4")},
						{Input: "[7,7]", Expected: text("5")},
						{Input: "[500,1]", Expected: text("500")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/reveal-cards-in-increasing-order/",
					Number:     number(950),
					Title:      "Reveal Cards In Increasing Order",
					Difficulty: "medium",
					Topics:     []string{"Array", "Queue", "Sorting", "Simulation"},
					Note:       "Sorting карт и Simulation на Deque номеров мест. n - число карт. Время: O(n log n), включая Sorting; далее n записей ответа, 2n - 1 удалений и n - 1 возвратов номеров. Память: O(n) дополнительно для Deque и Python Sorting; ответ отдельно O(n). Для [1, 2, 3, 4]: 4 записи, 7 popleft, 3 append, пик 4 номера, ответ [1, 3, 2, 4]. list.pop(0) может сделать Simulation квадратичной; обратная Simulation на Deque также требует Sorting.",
					Statement: `<p>Дан Array целых чисел <code>deck</code>. Есть колода карт, на каждой карте записано своё число, и все числа <strong>разные</strong>. На <code>i</code>-й карте записано число <code>deck[i]</code>.</p><p>Карты можно сложить в колоду в любом порядке. Сначала все карты лежат в колоде лицом вниз (закрыты).</p><p>Повторяйте эти шаги, пока все карты не будут открыты:</p><ol><li>Возьмите верхнюю карту колоды, откройте её и уберите из колоды.</li><li>Если в колоде ещё есть карты, положите следующую верхнюю карту вниз колоды.</li><li>Если ещё есть закрытые карты, вернитесь к шагу 1. Иначе остановитесь.</li></ol><p>Верните <em>порядок карт в колоде, при котором карты откроются по возрастанию</em>.</p><p><strong>Обратите внимание:</strong> первое число ответа это верх колоды.</p><p>Метод класса <code>Solution</code>: <code>def deckRevealedIncreasing(self, deck: list[int]) -&gt; list[int]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> deck = [17,13,11,2,3,5,7]
<strong>Выход:</strong> [2,13,3,11,5,17,7]
<strong>Пояснение:</strong>
Колода пришла в порядке [17,13,11,2,3,5,7] (этот порядок не важен), и мы меняем порядок.
После этого колода [2,13,3,11,5,17,7], сверху лежит 2.
Открываем 2 и кладём 13 вниз. Колода [3,11,5,17,7,13].
Открываем 3 и кладём 11 вниз. Колода [5,17,7,13,11].
Открываем 5 и кладём 17 вниз. Колода [7,13,11,17].
Открываем 7 и кладём 13 вниз. Колода [11,17,13].
Открываем 11 и кладём 17 вниз. Колода [13,17].
Открываем 13 и кладём 17 вниз. Колода [17].
Открываем 17.
Все карты открылись по возрастанию, значит, ответ верный.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> deck = [1,1000]
<strong>Выход:</strong> [1,1000]</pre><h3>Ограничения</h3><ul><li><code>1 ≤ deck.length ≤ 1000</code></li><li><code>1 ≤ deck[i] ≤ 10<sup>6</sup></code></li><li>все числа <code>deck</code> разные</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "deckRevealedIncreasing"}, Cases: []sets.TestCase{
						{Input: "[[17,13,11,2,3,5,7]]", Expected: text("[2,13,3,11,5,17,7]")},
						{Input: "[[1,1000]]", Expected: text("[1,1000]")},
						{Input: "[[1]]", Expected: text("[1]")},
						{Input: "[[1,2,3,4]]", Expected: text("[1,3,2,4]")},
						{Input: "[[4,3,2,1,5]]", Expected: text("[1,5,2,4,3]")},
						{Input: "[[60,10,50,20,40,30]]", Expected: text("[10,40,20,60,30,50]")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "hq1-pop-end", Type: "choice", Prompt: "q = deque([4, 9, 6]). Ученик хочет взять первый элемент Queue и вызывает q.pop(). Что вернёт вызов?", Options: []string{"9", "4", "6"}, Answer: 2, Explain: "pop() берёт с конца: вернёт 6, останется [4, 9]. Для Queue нужен q.popleft(), он вернул бы 4. Оба удаления требуют O(1) времени и рабочего состояния; создание Deque из n элементов требует O(n) времени и памяти. Здесь создано 3 элемента и сделано 1 удаление."},
				{ID: "hq1-pop-zero", Type: "choice", Prompt: "В list a 200 000 чисел. Цикл вызывает a.pop(0), пока list не опустеет. Сколько сдвигов ссылок он сделает?", Options: []string{"около 10 миллиардов", "около 20 миллиардов", "около 40 миллиардов"}, Answer: 1, Explain: "Сдвигов 199 999 + 199 998 + ... + 0 = 19 999 900 000. Для n элементов это n(n - 1)/2, время O(n²), дополнительная память O(1) сверх входа, если удалённые значения не собирать. Deque делает n вызовов popleft за O(n) без сдвига оставшихся ссылок; построение нового Deque требует O(n) памяти."},
				{ID: "hq1-middle", Type: "choice", Prompt: "В Deque q находится n чисел. Для какого обращения число шагов растёт линейно с n?", Options: []string{"q[0]", "q[-1]", "q[len(q) // 2]"}, Answer: 2, Explain: "Deque проходит блоки от ближайшего конца до середины: число блоков пропорционально n, время O(n). Чтение концов стоит O(1). Дополнительная память каждого чтения O(1), ответ - одно число. Для доступа к произвольному индексу list даёт O(1); копирование туда всех n элементов стоит O(n) времени и памяти."},
				{ID: "hq1-empty", Type: "choice", Prompt: "q = deque([3]). Код вызывает q.popleft() два раза подряд. Что даст второй вызов?", Options: []string{"None", "IndexError", "0"}, Answer: 1, Explain: "Первый вызов вернул 3, Deque пуст. Второй вызывает IndexError. Перед удалением проверяйте if q: или while q:. Здесь 2 попытки удаления, пик 1 элемент. Проверка пустоты и удаление требуют O(1) времени и дополнительного состояния."},
				{ID: "hq1-tickets-code", Type: "code", Prompt: "Simulation очереди за билетами, LeetCode 2073. Какая строка стоит на месте ___?", Language: "python", Code: `q = deque(range(len(tickets)))
time = 0
while True:
    p = q.popleft()
    tickets[p] -= 1
    time += 1
    if p == k and tickets[p] == 0:
        return time
    if tickets[p] > 0:
        ___`, Options: []string{"q.append(p)", "q.appendleft(p)", "q.pop()"}, Answer: 0, Explain: "Человек встаёт в конец: q.append(p). appendleft дал бы человеку 0 купить 3 билета подряд для [3, 1, 2], k = 0, вместо правильных 6 продаж. pop удаляет другого человека. При n людях и b продажах время O(n + b), дополнительная память O(n), ответ O(1). Здесь 6 popleft, 3 append в цикле, пик 3 номера."},
				{ID: "hq1-lunch", Type: "short", Prompt: "LeetCode 1700. students = [0, 0, 1], sandwiches = [1, 1, 0], верх стопки первый в Array. Сколько учеников останется без обеда?", Answer: "2", Explain: "Два отказа, затем одна выдача оставляют [0, 0]. Следующий бутерброд снова 1. Ещё 2 отказа дают полный круг: остановка, ответ 2. Всего 5 проверок, пик 3 ученика. Simulation требует O(n²) времени в худшем случае и O(n) памяти для n учеников; два счётчика желаний дают O(n) времени и O(1) памяти."},
				{ID: "hq1-rotate-k", Type: "choice", Prompt: "d = deque([1, 2, 3, 4, 5, 6, 7]). В считалке выбывает каждый четвёртый, счёт начинается с 1. Какой вызов нужен перед d.popleft()?", Options: []string{"d.rotate(-3)", "d.rotate(3)", "d.rotate(-4)"}, Answer: 0, Explain: "rotate(-3) даёт [4, 5, 6, 7, 1, 2, 3], затем popleft удаляет 4. Два других вызова поставили бы первым 5. Для длины s CPython выбирает d = min(r mod s, s - (r mod s)): время O(1 + d), дополнительное состояние O(1). Здесь s = 7, r = -3, d = 3; создание хранит 7 элементов, после удаления живых 6."},
				{ID: "hq1-reveal", Type: "short", Prompt: "Колода сверху вниз [1, 2, 3, 4, 5]. Открываем верхнюю карту, следующую кладём вниз, и так до конца. Какая карта откроется третьей?", Answer: "5", Explain: "Открыли 1, перенесли 2; открыли 3, перенесли 4; открыли 5. До третьего открытия выполнены 5 popleft и 2 append. Весь порядок: 1, 3, 5, 4, 2. Simulation n карт на Deque стоит O(n) времени и O(n) памяти, пик здесь 5. Если собрать весь порядок, ответ займёт ещё n мест; для третьей карты нужен один результат."},
			},
			Links: []sets.Link{
				{Title: "Deque и list в Python: стоимость операций", URL: "https://habr.com/ru/companies/vk/articles/700520/", Source: "Хабр"},
				{Title: "Queue: порядок операций", URL: "https://neerc.ifmo.ru/wiki/index.php?title=%D0%9E%D1%87%D0%B5%D1%80%D0%B5%D0%B4%D1%8C", Source: "Викиконспекты ИТМО"},
				{Title: "Deque: работа с двумя концами", URL: "https://neerc.ifmo.ru/wiki/index.php?title=%D0%94%D0%B5%D0%BA", Source: "Викиконспекты ИТМО"},
				{Title: "collections.deque", URL: "https://docs.python.org/3/library/collections.html#collections.deque", Source: "Python docs (англ.)"},
				{Title: "Queue на list: ограничения", URL: "https://docs.python.org/3/tutorial/datastructures.html#using-lists-as-queues", Source: "Python tutorial (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "hq1-queue",
					Title: "Queue и Deque: вызовы и живые элементы",
					Frames: []sets.FigureFrame{
						{Caption: "q.append(5). Один вызов на конце, один живой элемент. Счётчик вызовов не включает создание пустого q.", Rows: []sets.FigureRow{
							{Label: "q", Cells: []any{5}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "голова", Cell: 0}, {Name: "хвост", Cell: 0}}},
							{Label: "вызовы / живых", Cells: []any{1, 1}},
						}},
						{Caption: "q.append(8). Всего 2 вызова, живых элементов 2.", Rows: []sets.FigureRow{
							{Label: "q", Cells: []any{5, 8}, Marked: []int{1}, Pointers: &sets.Pointers{{Name: "голова", Cell: 0}, {Name: "хвост", Cell: 1}}},
							{Label: "вызовы / живых", Cells: []any{2, 2}},
						}},
						{Caption: "q.append(2). Всего 3 вызова. Достигнут пик q: 3 элемента.", Rows: []sets.FigureRow{
							{Label: "q", Cells: []any{5, 8, 2}, Marked: []int{2}, Pointers: &sets.Pointers{{Name: "голова", Cell: 0}, {Name: "хвост", Cell: 2}}},
							{Label: "вызовы / живых", Cells: []any{3, 3}},
						}},
						{Caption: "q.popleft() вернул 5. Всего 4 вызова, живых элементов 2. Затемнённая ячейка показывает удалённое значение, а не сохранённую ссылку q.", Rows: []sets.FigureRow{
							{Label: "q", Cells: []any{5, 8, 2}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "голова", Cell: 1}, {Name: "хвост", Cell: 2}}},
							{Label: "вызовы / живых", Cells: []any{4, 2}},
						}},
						{Caption: "q.append(7). Всего 5 вызовов, живых элементов 3; пик остаётся 3.", Rows: []sets.FigureRow{
							{Label: "q", Cells: []any{8, 2, 7}, Marked: []int{2}, Pointers: &sets.Pointers{{Name: "голова", Cell: 0}, {Name: "хвост", Cell: 2}}},
							{Label: "вызовы / живых", Cells: []any{5, 3}},
						}},
						{Caption: "q.popleft() вернул 8. Всего 6 вызовов за O(1) каждый; пик 3 элемента. Вывод: два числа 5 и 8. Для m вызовов время O(m), память O(h), где h - пик длины.", Rows: []sets.FigureRow{
							{Label: "q", Cells: []any{8, 2, 7}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "голова", Cell: 1}, {Name: "хвост", Cell: 2}}},
							{Label: "вызовы / живых", Cells: []any{6, 2}},
						}},
						{Caption: "Отдельный пример: создание d копирует 3 ссылки. Затем appendleft(1) делает один вызов на конце; длина и пик d равны 4. Создание n элементов стоит O(n) времени и памяти.", Rows: []sets.FigureRow{
							{Label: "d", Cells: []any{1, 2, 4, 5}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "голова", Cell: 0}, {Name: "хвост", Cell: 3}}},
							{Label: "вызовы / живых", Cells: []any{1, 4}},
						}},
						{Caption: "d.pop() вернул 5. После создания выполнено 2 вызова, в d живых элементов 3, пик 4. Само удаление требует O(1) времени и временного состояния.", Rows: []sets.FigureRow{
							{Label: "d", Cells: []any{1, 2, 4, 5}, Dimmed: []int{3}, Pointers: &sets.Pointers{{Name: "голова", Cell: 0}, {Name: "хвост", Cell: 2}}},
							{Label: "вызовы / живых", Cells: []any{2, 3}},
						}},
					},
				},
				{
					ID:    "hq1-shift",
					Title: "list.pop(0): считаем сдвиги ссылок",
					Frames: []sets.FigureFrame{
						{Caption: "До удаления входной list содержит 5 ссылок. Сдвигов пока 0. Рабочее состояние удаления O(1), эти 5 входных мест не считаем дополнительной памятью.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{5, 8, 2, 7, 4}, Marked: []int{0}},
							{Label: "сдвиги / длина", Cells: []any{0, 5}},
						}},
						{Caption: "pop(0) вернул 5 и сдвинул 4 ссылки. Удаление всех 5 значений даст 4 + 3 + 2 + 1 = 10 сдвигов. Для n: O(n²) времени, O(1) дополнительной памяти без сбора ответа; для 200 000: 19 999 900 000 сдвигов.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{8, 2, 7, 4}, Marked: []int{0, 1, 2, 3}},
							{Label: "сдвиги / длина", Cells: []any{4, 4}},
						}},
					},
				},
				{
					ID:    "hq1-tickets",
					Title: "Simulation билетов: 6 продаж, пик 3 номера",
					Frames: []sets.FigureFrame{
						{Caption: "tickets = [3, 1, 2], k = 0. Созданы 3 номера. Продажа 1: человек 0 возвращается. Счётчики далее: продажи, append в цикле, живые номера. Пик Queue равен 3; оставшиеся билеты хранятся во входе.", Rows: []sets.FigureRow{
							{Label: "Queue", Cells: []any{1, 2, 0}, Marked: []int{2}},
							{Label: "осталось", Cells: []any{1, 2, 2}, Marked: []int{2}},
							{Label: "счётчики", Cells: []any{1, 1, 3}},
						}},
						{Caption: "Продажа 2: человек 1 уходит. Выполнено 2 popleft и 1 append; живых номеров 2.", Rows: []sets.FigureRow{
							{Label: "Queue", Cells: []any{2, 0}, Marked: []int{1}},
							{Label: "осталось", Cells: []any{2, 2}, Marked: []int{1}},
							{Label: "счётчики", Cells: []any{2, 1, 2}},
						}},
						{Caption: "Продажа 3: человек 2 возвращается. Выполнено 3 popleft и 2 append; живых номеров 2.", Rows: []sets.FigureRow{
							{Label: "Queue", Cells: []any{0, 2}, Marked: []int{0}},
							{Label: "осталось", Cells: []any{2, 1}, Marked: []int{0}},
							{Label: "счётчики", Cells: []any{3, 2, 2}},
						}},
						{Caption: "Продажа 4: человек 0 возвращается. Выполнено 4 popleft и 3 append; живых номеров 2.", Rows: []sets.FigureRow{
							{Label: "Queue", Cells: []any{2, 0}, Marked: []int{1}},
							{Label: "осталось", Cells: []any{1, 1}, Marked: []int{1}},
							{Label: "счётчики", Cells: []any{4, 3, 2}},
						}},
						{Caption: "Продажа 5: человек 2 уходит. Выполнено 5 popleft и 3 append; живой номер 1.", Rows: []sets.FigureRow{
							{Label: "Queue", Cells: []any{0}, Marked: []int{0}},
							{Label: "осталось", Cells: []any{1}, Marked: []int{0}},
							{Label: "счётчики", Cells: []any{5, 3, 1}},
						}},
						{Caption: "Продажа 6: человек 0 закончил. Queue пуст. Итог: 6 popleft, 3 append, пик 3 номера; ответ - одно число 6. Для n людей и b продаж время O(n + b), дополнительная память O(n).", Rows: []sets.FigureRow{
							{Label: "Queue", Cells: []any{0}, Dimmed: []int{0}},
							{Label: "осталось", Cells: []any{0}, Dimmed: []int{0}},
							{Label: "счётчики", Cells: []any{6, 3, 0}},
						}},
					},
				},
				{
					ID:    "hq1-rotate",
					Title: "Deque: расстояние поворота и длина",
					Frames: []sets.FigureFrame{
						{Caption: "Созданы 6 элементов за O(n), n = 6; пик d равен 6. Счёт начинается с 1, выбывает каждый третий. Поворотов и удалений пока нет.", Rows: []sets.FigureRow{
							{Label: "d", Cells: []any{1, 2, 3, 4, 5, 6}, Marked: []int{0, 1}},
							{Label: "расстояние / живых", Cells: []any{0, 6}},
						}},
						{Caption: "rotate(-2): a = (-2) mod 6 = 4, кратчайшее расстояние min(4, 2) = 2. Время O(1 + d), временное состояние O(1). Живых элементов 6; rotate(-8) дал бы то же состояние и то же расстояние.", Rows: []sets.FigureRow{
							{Label: "d", Cells: []any{3, 4, 5, 6, 1, 2}, Marked: []int{4, 5}},
							{Label: "расстояние / живых", Cells: []any{2, 6}},
						}},
						{Caption: "Один popleft за O(1) убрал 3. Живых элементов 5, пик 6. За всю игру n = 6, k = 3 будут 5 удалений и сумма расстояний 7; останется 1. Simulation хранит O(n) элементов, ответ - один номер.", Rows: []sets.FigureRow{
							{Label: "d", Cells: []any{3, 4, 5, 6, 1, 2}, Dimmed: []int{0}},
							{Label: "удалений / живых", Cells: []any{1, 5}},
						}},
					},
				},
			},
		},
		{
			Title:   "Binary Heap: приоритет и стоимость",
			Summary: "Min-Heap и Max-Heap выбирают элемент по приоритету. Вы считаете высоту Tree, построение Heapify, обновления и пиковую длину Binary Heap. Пары tuple задают правило равенства, а Min-Heap размера k хранит k лучших.",
			Lesson:  lesson2,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/last-stone-weight/",
					Number:     number(1046),
					Title:      "Last Stone Weight",
					Difficulty: "easy",
					Topics:     []string{"Array", "Heap (Priority Queue)"},
					Note:       "n - число камней. Min-Heap отрицательных весов выбирает два максимума. Время: O(n log(n + 1)), включая O(n) на копию и Heapify; кругов не больше n - 1. Память: O(n) дополнительно для h, ответ O(1). Для [6, 4, 9, 2]: 6 pop, 3 push, пик 4, ответ 1. Повторный Sorting имеет верхнюю оценку O(n² log(n + 1)) времени и O(n) рабочей памяти. Python 3.14 также поддерживает Max-Heap напрямую.",
					Statement: `<p>Дан Array целых чисел <code>stones</code>: <code>stones[i]</code> это вес <code>i</code>-го камня.</p><p>Мы играем с камнями. На каждом ходу берём <strong>два самых тяжёлых камня</strong> и разбиваем их друг о друга. Пусть веса этих камней <code>x</code> и <code>y</code>, и <code>x ≤ y</code>. Результат удара:</p><ul><li>если <code>x == y</code>, оба камня разрушаются;</li><li>если <code>x != y</code>, камень веса <code>x</code> разрушается, а камень веса <code>y</code> получает новый вес <code>y - x</code>.</li></ul><p>В конце игры остаётся <strong>не больше одного</strong> камня.</p><p>Верните <em>вес последнего камня</em>. Если камней не осталось, верните <code>0</code>.</p><p>Метод класса <code>Solution</code>: <code>def lastStoneWeight(self, stones: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> stones = [2,7,4,1,8,1]
<strong>Выход:</strong> 1
<strong>Пояснение:</strong>
Разбиваем 7 и 8, получаем 1, и Array становится [2,4,1,1,1].
Разбиваем 2 и 4, получаем 2, и Array становится [2,1,1,1].
Разбиваем 2 и 1, получаем 1, и Array становится [1,1,1].
Разбиваем 1 и 1, получаем 0, и Array становится [1]. Это вес последнего камня.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> stones = [1]
<strong>Выход:</strong> 1</pre><h3>Ограничения</h3><ul><li><code>1 ≤ stones.length ≤ 30</code></li><li><code>1 ≤ stones[i] ≤ 1000</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "lastStoneWeight"}, Cases: []sets.TestCase{
						{Input: "[[2,7,4,1,8,1]]", Expected: text("1")},
						{Input: "[[1]]", Expected: text("1")},
						{Input: "[[2,2]]", Expected: text("0")},
						{Input: "[[10,4,2,10]]", Expected: text("2")},
						{Input: "[[1000,1]]", Expected: text("999")},
						{Input: "[[3,3,3]]", Expected: text("3")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/final-array-state-after-k-multiplication-operations-i/",
					Number:     number(3264),
					Title:      "Final Array State After K Multiplication Operations I",
					Difficulty: "easy",
					Topics:     []string{"Array", "Math", "Heap (Priority Queue)", "Simulation"},
					Note:       "tuple (число, индекс) сохраняет первое место при равенстве. n - длина nums, t = k - число умножений. Время: O(n + t log(n + 1)) на Heapify и t обновлений. Память: O(n) пар дополнительно; ответ из n чисел можно записывать во входной nums. При n = 5, t = 5: 5 обновлений, пик 5 пар. Linear Scan на каждом шаге делает t(n - 1) сравнений: 20 в примере; O(n + tn) времени с выводом и O(1) рабочего состояния. При t ≤ 10 этот способ проще.",
					Statement: `<p>Даны Array целых чисел <code>nums</code>, целое число <code>k</code> и целое число <code>multiplier</code>.</p><p>Сделайте с <code>nums</code> ровно <code>k</code> операций. Каждая операция:</p><ul><li>Найдите <strong>наименьшее</strong> значение <code>x</code> в <code>nums</code>. Если наименьшее значение встречается несколько раз, выберите то, которое стоит <strong>первым</strong>.</li><li>Замените выбранное значение <code>x</code> на <code>x * multiplier</code>.</li></ul><p>Верните Array <code>nums</code> <em>после всех <code>k</code> операций</em>.</p><p>Метод класса <code>Solution</code>: <code>def getFinalState(self, nums: list[int], k: int, multiplier: int) -&gt; list[int]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [2,1,3,5,6], k = 5, multiplier = 2
<strong>Выход:</strong> [8,4,6,5,6]
<strong>Пояснение:</strong>
После операции 1: [2, 2, 3, 5, 6]
После операции 2: [4, 2, 3, 5, 6]
После операции 3: [4, 4, 3, 5, 6]
После операции 4: [4, 4, 6, 5, 6]
После операции 5: [8, 4, 6, 5, 6]</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [1,2], k = 3, multiplier = 4
<strong>Выход:</strong> [16,8]
<strong>Пояснение:</strong>
После операции 1: [4, 2]
После операции 2: [4, 8]
После операции 3: [16, 8]</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 100</code></li><li><code>1 ≤ nums[i] ≤ 100</code></li><li><code>1 ≤ k ≤ 10</code></li><li><code>1 ≤ multiplier ≤ 5</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "getFinalState"}, Cases: []sets.TestCase{
						{Input: "[[2,1,3,5,6],5,2]", Expected: text("[8,4,6,5,6]")},
						{Input: "[[1,2],3,4]", Expected: text("[16,8]")},
						{Input: "[[3,3,3],2,2]", Expected: text("[6,6,3]")},
						{Input: "[[5,5],10,1]", Expected: text("[5,5]")},
						{Input: "[[1],10,5]", Expected: text("[9765625]")},
						{Input: "[[2,3,5],3,2]", Expected: text("[8,6,5]")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1800/C2",
					Title:      "Powering the Hero (hard version)",
					Difficulty: "easy",
					Topics:     []string{"Array", "Greedy", "Heap (Priority Queue)"},
					Note:       "Codeforces, рейтинг 1100. Max-Heap временно назначает герою наибольший доступный бонус. Обмен назначений сохраняет сумму и даёт допустимый порядок Stack; герой берёт только вершину. n - число карт одного набора. Время: O(n log(n + 1)), каждый бонус входит один раз и выходит не больше раза. Память: O(n) дополнительно, ответ O(1). Для [4, 1, 0, 7, 0, 0]: 3 push, 3 pop, пик 2, сумма 12. Linear Scan для каждого героя может дать O(n²). Sorting всей колоды нарушает доступность бонусов. Для нескольких наборов сложите время; при обработке по одному память зависит от наибольшего n.",
					Statement: `<p>Это сложная версия задачи. Она отличается от простой только ограничениями на <code>n</code> и <code>t</code>.</p><p>Есть колода из <code>n</code> карт. У каждой карты есть сила. Карты бывают двух видов:</p><ul><li>карта героя: её сила всегда равна 0;</li><li>карта бонуса: её сила всегда положительна.</li></ul><p>С колодой можно делать так:</p><ul><li>взять карту сверху колоды;</li><li>если это карта бонуса, её можно положить на верх своей стопки бонусов или в сброс;</li><li>если это карта героя, к его силе прибавляется сила верхней карты из стопки бонусов (если стопка не пуста). После этого герой входит в вашу армию, а использованный бонус уходит в сброс.</li></ul><p>Соберите такими действиями армию с наибольшей суммарной силой.</p><h3>Входные данные</h3><p>В первой строке записано целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>): количество наборов входных данных.</p><p>В первой строке каждого набора записано целое число <code>n</code> (<code>1 ≤ n ≤ 2 · 10<sup>5</sup></code>): количество карт в колоде.</p><p>Во второй строке набора записаны <code>n</code> целых чисел <code>s<sub>1</sub>, s<sub>2</sub>, ..., s<sub>n</sub></code> (<code>0 ≤ s<sub>i</sub> ≤ 10<sup>9</sup></code>): силы карт сверху вниз.</p><p>Гарантируется, что сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Выведите <code>t</code> чисел, по одному для каждого набора: наибольшую суммарную силу армии.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>5
5
3 3 3 0 0
6
0 3 3 0 0 3
7
1 2 3 0 4 5 0
7
1 2 5 0 4 3 0
5
3 1 0 0 4</pre><p><strong>Выходные данные:</strong></p><pre>6
6
8
9
4</pre><p><strong>Пояснение:</strong> в первом наборе можно взять бонусы 1 и 2. Оба героя получат по 3 к силе. Если взять все бонусы, один из них останется без дела.</p><p>Во втором наборе героя сверху колоды усилить нельзя, а двух других героев можно усилить бонусами 2 и 3 и получить суммарную силу 6.</p><p>В четвёртом наборе можно взять бонусы 1, 2, 3, 5 и пропустить бонус 6. Тогда герой 4 получит от бонуса 3 силу 5, а герой 7 от бонуса 5 силу 4: 4 + 5 = 9.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>1 ≤ n ≤ 2 · 10<sup>5</sup></code></li><li><code>0 ≤ s<sub>i</sub> ≤ 10<sup>9</sup></code></li><li>сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code></li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `5
5
3 3 3 0 0
6
0 3 3 0 0 3
7
1 2 3 0 4 5 0
7
1 2 5 0 4 3 0
5
3 1 0 0 4
`, Expected: text(`6
6
8
9
4
`)},
						{Input: `3
1
0
1
5
4
0 5 0 0
`, Expected: text(`0
0
5
`)},
						{Input: `1
4
1000000000 1000000000 0 0
`, Expected: text(`2000000000
`)},
						{Input: `2
6
4 1 0 7 0 0
3
0 0 0
`, Expected: text(`12
0
`)},
						{Input: `1
7
2 9 0 3 0 0 0
`, Expected: text(`14
`)},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/kth-largest-element-in-an-array/",
					Number:     number(215),
					Title:      "Kth Largest Element in an Array",
					Difficulty: "medium",
					Topics:     []string{"Array", "Divide and Conquer", "Sorting", "Heap (Priority Queue)", "Quickselect"},
					Note:       "n - длина nums. Min-Heap размера k хранит k наибольших чисел с повторами. Время: O(n log(k + 1)), по одному посещению каждого числа и до одного обновления. Память: O(k) дополнительно, ответ O(1). Для [7, 2, 9, 4, 6, 1, 8], k = 3: 3 push, 4 сравнения с корнем, 3 замены, пик 3, ответ 7. Sorting требует O(n log n) времени и до O(n) памяти Python. Итеративный Quickselect на месте со случайным выбором опорного элемента требует ожидаемого O(n) времени, O(n²) в худшем случае и O(1) дополнительного состояния.",
					Statement: `<p>Даны Array целых чисел <code>nums</code> и целое число <code>k</code>. Верните <em><code>k</code>-й по величине элемент Array</em>.</p><p>Обратите внимание: нужен <code>k</code>-й по величине элемент в отсортированном порядке, а не <code>k</code>-й среди разных значений.</p><p>Можете ли вы решить задачу без Sorting?</p><p>Метод класса <code>Solution</code>: <code>def findKthLargest(self, nums: list[int], k: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [3,2,1,5,6,4], k = 2
<strong>Выход:</strong> 5</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [3,2,3,1,2,4,5,5,6], k = 4
<strong>Выход:</strong> 4</pre><h3>Ограничения</h3><ul><li><code>1 ≤ k ≤ nums.length ≤ 10<sup>5</sup></code></li><li><code>-10<sup>4</sup> ≤ nums[i] ≤ 10<sup>4</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "findKthLargest"}, Cases: []sets.TestCase{
						{Input: "[[3,2,1,5,6,4],2]", Expected: text("5")},
						{Input: "[[3,2,3,1,2,4,5,5,6],4]", Expected: text("4")},
						{Input: "[[1],1]", Expected: text("1")},
						{Input: "[[-1,-3,-2],3]", Expected: text("-3")},
						{Input: "[[2,9,9,4],2]", Expected: text("9")},
						{Input: "[[-10000,10000,0],2]", Expected: text("0")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/top-k-frequent-words/",
					Number:     number(692),
					Title:      "Top K Frequent Words",
					Difficulty: "medium",
					Topics:     []string{"Array", "Hash Table", "String", "Trie", "Sorting", "Heap (Priority Queue)", "Bucket Sort", "Counting"},
					Note:       "Counter, затем Heapify всех пар (-частота, слово) и k удалений. n - число слов, u - разных слов, S - всего символов, L - максимальная длина слова. Время: ожидаемое O(S + uL + kL log(u + 1)); сравнение слов при равной частоте стоит до O(L). При ограниченной L: O(n + u + k log(u + 1)). Память: O(u) рабочих ссылок и счётчиков, ответ отдельно O(k) ссылок. Для [b, a, b, c, a], k = 2: 5 обновлений Counter, 3 пары, 2 pop, ответ [a, b]. Sorting u разных слов после Counting даёт ожидаемое O(S + Lu log(u + 1)) времени и O(u) рабочей памяти; Binary Heap размера k сохраняет O(u) памяти Counter.",
					Statement: `<p>Даны Array значений String <code>words</code> и целое число <code>k</code>. Верните <em><code>k</code> самых частых слов</em>.</p><p>Расположите ответ по <strong>частоте</strong> от большей к меньшей. Слова с одинаковой частотой расположите в <strong>Lexicographic Order</strong> (по алфавиту).</p><p>Метод класса <code>Solution</code>: <code>def topKFrequent(self, words: list[str], k: int) -&gt; list[str]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> words = ["i","love","leetcode","i","love","coding"], k = 2
<strong>Выход:</strong> ["i","love"]
<strong>Пояснение:</strong> "i" и "love" это два самых частых слова.
"i" стоит раньше "love", потому что идёт раньше по алфавиту.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> words = ["the","day","is","sunny","the","the","the","sunny","is","is"], k = 4
<strong>Выход:</strong> ["the","is","sunny","day"]
<strong>Пояснение:</strong> "the", "is", "sunny" и "day" это четыре самых частых слова.
Они встречаются 4, 3, 2 и 1 раз.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ words.length ≤ 500</code></li><li><code>1 ≤ words[i].length ≤ 10</code></li><li><code>words[i]</code> состоит из маленьких английских букв</li><li><code>k</code> от 1 до числа <strong>разных</strong> слов в <code>words</code></li></ul><p><strong>Дополнительный вопрос:</strong> можете ли вы при этих ограничениях длины слова решить задачу за ожидаемое O(n log(k + 1)) времени и O(n) дополнительной памяти, где n - число слов?</p>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "topKFrequent"}, Cases: []sets.TestCase{
						{Input: `[["i","love","leetcode","i","love","coding"],2]`, Expected: text(`["i","love"]`)},
						{Input: `[["the","day","is","sunny","the","the","the","sunny","is","is"],4]`, Expected: text(`["the","is","sunny","day"]`)},
						{Input: `[["a"],1]`, Expected: text(`["a"]`)},
						{Input: `[["z","y","x"],3]`, Expected: text(`["x","y","z"]`)},
						{Input: `[["aa","a","aa","b","a","b"],3]`, Expected: text(`["a","aa","b"]`)},
						{Input: `[["b","c","a","c","b","c"],2]`, Expected: text(`["c","b"]`)},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "hq2-max-place", Type: "choice", Prompt: "В Min-Heap 6 разных чисел, места от 0 до 5. Где может стоять наибольшее число?", Options: []string{"только на месте 5, последнем", "на любом из мест 3, 4 и 5", "на месте 1 или 2"}, Answer: 1, Explain: "Максимум может быть только в листе: у мест 3, 4, 5 нет детей. В [1, 3, 2, 5, 9, 8] это 9 на месте 4. Linear Scan трёх листьев находит максимум за 2 сравнения. Для n элементов это O(n) времени, O(1) дополнительной памяти и один результат; минимум уже доступен в h[0] за O(1)."},
				{ID: "hq2-push-swaps", Type: "short", Prompt: "В Binary Heap 31 число. Новое число добавлено на место 31. Сколько обменов максимум сделает учебный подъём через обмен с родителем?", Answer: "5", Explain: "Путь индексов: 31, 15, 7, 3, 1, 0. Это 5 рёбер и не больше 5 обменов. Пик длины 32. Для h элементов подъём требует O(log(h + 1)) сравнений и O(1) временного состояния. Это счёт обменов учебной модели; heapq может перемещать ссылки без парных обменов."},
				{ID: "hq2-sift-down", Type: "short", Prompt: "h = [1, 6, 2, 7, 8, 5, 3]. Вызов heappop(h) вернул 1. На каком месте теперь стоит число 3?", Answer: "2", Explain: "Учебный спуск ставит 3 в корень и меняет с меньшим ребёнком 2. Ниже остался только ребёнок 5, и 3 < 5: h = [2, 6, 3, 7, 8, 5]. Один обмен, 3 сравнения значений: выбор ребёнка и 2 проверки с ним. Длина уменьшилась с 7 до 6. Для h элементов pop требует O(log(h + 1)) сравнений, O(1) временного состояния и возвращает одно число."},
				{ID: "hq2-heapify-none", Type: "code", Prompt: "Что случится, когда выполнится этот код?", Language: "python", Code: `import heapq

h = [5, 3, 8]
h = heapq.heapify(h)
print(h[0])`, Options: []string{"выведет 3", "выведет 5", "TypeError: у None нет h[0]"}, Answer: 2, Explain: "heapify меняет list на месте и возвращает None. Поэтому h[0] вызывает TypeError. Пишите heapq.heapify(h) без присваивания. Для n элементов Heapify стоит O(n) времени и O(1) дополнительной памяти сверх готового list. Здесь 3 элемента, один родитель; после правильного вызова корень равен 3."},
				{ID: "hq2-pop-negate", Type: "code", Prompt: "Min-Heap хранит отрицательные веса. Что стоит на месте ___, чтобы y стал весом самого тяжёлого камня, 8?", Language: "python", Code: `h = [-x for x in [5, 1, 8]]
heapq.heapify(h)
y = ___`, Options: []string{"-heapq.heappop(h)", "heapq.heappop(h)", "-h.pop()"}, Answer: 0, Explain: "После Heapify h = [-8, -1, -5]. heappop вернёт -8, и смена знака даёт 8. -h.pop() дал бы 5. Создание n отрицательных весов и Heapify требуют O(n) времени, O(n) памяти; затем один pop стоит O(log(n + 1)). Здесь пик 3 элемента, ответ - одно число. Python 3.14 позволяет хранить исходные веса через Max-Heap API."},
				{ID: "hq2-tuple-order", Type: "short", Prompt: "Число солдат в рядах 0, 1, 2, 3, 4: [3, 1, 3, 0, 1]. В Min-Heap кладут tuple (число солдат, номер ряда) в порядке рядов 4, 3, 2, 1, 0 и достают по одному. Какой номер ряда выйдет третьим?", Answer: "4", Explain: "Порядок: (0, 3), (1, 1), (1, 4), (3, 0), (3, 2). При равном числе солдат решает номер. Для третьего ответа нужны 5 push и 3 pop, пик 5 пар. Для n таких пар последовательные вставки и все удаления требуют O(n log(n + 1)) времени и O(n) памяти; сравнение двух ограниченных чисел стоит O(1)."},
				{ID: "hq2-tuple-dict", Type: "choice", Prompt: `В Min-Heap лежит (1, {"b": 2}). Код вызывает heappush(h, (1, {"a": 1})). Что случится?`, Options: []string{"TypeError: dict нельзя сравнить через <", "tuple встанут в порядке добавления", "tuple встанут по алфавиту ключей"}, Answer: 0, Explain: "Приоритеты равны, поэтому сравнение tuple доходит до двух dict и вызывает TypeError. Используйте разные номера: (приоритет, номер, объект). Тогда сравниваются только числовые поля за O(1). Для h записей вставка требует O(log(h + 1)) сравнений, а хранение пар или троек занимает O(h) ссылок."},
				{ID: "hq2-kth-replace", Type: "choice", Prompt: "Третье по величине ищут через Min-Heap размера 3. Сейчас h = [5, 8, 6], пришло 7. Что станет с Min-Heap?", Options: []string{"h = [5, 7, 6]: 8 уходит", "h = [6, 8, 7]: 5 уходит", "ничего: 7 меньше 8"}, Answer: 1, Explain: "7 больше корня 5: heapreplace(h, 7) даёт [6, 8, 7]. Это одно внешнее сравнение и одна замена; длина остаётся 3. Для размера k замена требует O(log(k + 1)) сравнений и O(1) временного состояния. Весь проход n чисел: O(n log(k + 1)) времени, O(k) дополнительной памяти, ответ O(1)."},
			},
			Links: []sets.Link{
				{Title: "Binary Heap: устройство и операции", URL: "https://ru.algorithmica.org/cs/basic-structures/heap/", Source: "Алгоритмика"},
				{Title: "Binary Heap: реализация", URL: "https://habr.com/ru/articles/112222/", Source: "Хабр"},
				{Title: "Binary Heap и Priority Queue: визуализация", URL: "https://visualgo.net/en/heap", Source: "VisuAlgo (англ.)"},
				{Title: "heapq: реализация Priority Queue", URL: "https://docs.python.org/3/library/heapq.html#priority-queue-implementation-notes", Source: "Python docs (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "hq2-tree",
					Title: "Min-Heap: 6 элементов, 5 связей",
					Frames: []sets.FigureFrame{
						{Caption: "Место 1 хранит 3; дети 5 и 9 не меньше него. Показаны 2 сравнения. Весь Min-Heap содержит 6 элементов и 5 связей. Tree и Array - два вида одного list, а не 12 рабочих мест.", Rows: []sets.FigureRow{
							{Label: "Min-Heap", Cells: []any{1, 3, 2, 5, 9, 8}, Marked: []int{1, 3, 4}, Tree: true},
							{Label: "h", Cells: []any{1, 3, 2, 5, 9, 8}, Marked: []int{1, 3, 4}},
							{Label: "связей / элементов", Cells: []any{5, 6}},
						}},
						{Caption: "Родитель места 5: (5 - 1) // 2 = 2. Здесь одно сравнение: 2 ≤ 8. Проверка всех n - 1 связей стоит O(n) времени и O(1) дополнительного состояния. Ответ - True или False.", Rows: []sets.FigureRow{
							{Label: "Min-Heap", Cells: []any{1, 3, 2, 5, 9, 8}, Marked: []int{2, 5}, Tree: true},
							{Label: "h", Cells: []any{1, 3, 2, 5, 9, 8}, Marked: []int{2, 5}},
							{Label: "связей / элементов", Cells: []any{5, 6}},
						}},
						{Caption: "h[-1] = 8, но максимум 9 стоит на месте 4. Чтение минимума h[0] = 1 требует одного доступа: O(1) времени и памяти. Поиск максимума через Linear Scan 3 листьев потребует 2 сравнений, в общем O(n) времени.", Rows: []sets.FigureRow{
							{Label: "Min-Heap", Cells: []any{1, 3, 2, 5, 9, 8}, Marked: []int{4, 5}, Tree: true},
							{Label: "h", Cells: []any{1, 3, 2, 5, 9, 8}, Marked: []int{4, 5}},
							{Label: "чтений корня / мест", Cells: []any{1, 6}},
						}},
					},
				},
				{
					ID:    "hq2-push-pop",
					Title: "Подъём и спуск: обмены учебной модели",
					Frames: []sets.FigureFrame{
						{Caption: "Добавили 0 на место 6. Длина и пик 7, обменов пока 0. Сравниваем с родителем 2. Считаем обмены учебной модели, а не внутренние перемещения heapq; два вида изображают один list.", Rows: []sets.FigureRow{
							{Label: "Min-Heap", Cells: []any{1, 3, 2, 5, 9, 8, 0}, Marked: []int{6}, Tree: true},
							{Label: "h", Cells: []any{1, 3, 2, 5, 9, 8, 0}, Marked: []int{6}},
							{Label: "обмены / длина", Cells: []any{0, 7}},
						}},
						{Caption: "Обмен 0 и 2: первый обмен. Длина 7. Теперь сравниваем 0 с родителем 1.", Rows: []sets.FigureRow{
							{Label: "Min-Heap", Cells: []any{1, 3, 0, 5, 9, 8, 2}, Marked: []int{2}, Tree: true},
							{Label: "h", Cells: []any{1, 3, 0, 5, 9, 8, 2}, Marked: []int{2}},
							{Label: "обмены / длина", Cells: []any{1, 7}},
						}},
						{Caption: "Обмен 0 и 1 завершает подъём: 2 сравнения и 2 обмена, длина 7. Временное состояние O(1). Длина пути ограничивает сравнения величиной O(log(h + 1)).", Rows: []sets.FigureRow{
							{Label: "Min-Heap", Cells: []any{0, 3, 1, 5, 9, 8, 2}, Marked: []int{0}, Tree: true},
							{Label: "h", Cells: []any{0, 3, 1, 5, 9, 8, 2}, Marked: []int{0}},
							{Label: "обмены / длина", Cells: []any{2, 7}},
						}},
						{Caption: "Начинаем pop: выбран результат 0, последний элемент 2 займёт корень. Счётчик обменов спуска начинаем с 0. До удаления длина 7, результат - одно число.", Rows: []sets.FigureRow{
							{Label: "Min-Heap", Cells: []any{0, 3, 1, 5, 9, 8, 2}, Dimmed: []int{0}, Tree: true},
							{Label: "h", Cells: []any{0, 3, 1, 5, 9, 8, 2}, Dimmed: []int{0}},
							{Label: "обмены / длина", Cells: []any{0, 7}},
						}},
						{Caption: "Перенесли 2 в корень, длина стала 6. Два сравнения: выбираем меньшего ребёнка 1 из 3 и 1, затем сравниваем 2 с 1. Обмен ещё не выполнен; пустая ячейка не хранится.", Rows: []sets.FigureRow{
							{Label: "Min-Heap", Cells: []any{2, 3, 1, 5, 9, 8, ""}, Marked: []int{0}, Tree: true},
							{Label: "h", Cells: []any{2, 3, 1, 5, 9, 8, ""}, Marked: []int{0}},
							{Label: "обмены / длина", Cells: []any{0, 6}},
						}},
						{Caption: "Обмен 2 и 1. Третье сравнение: 2 ≤ 8, спуск завершён. Итог модели: 1 обмен спуска, 3 сравнения. Пик list 7, рабочее состояние O(1); для h элементов обновление требует O(log(h + 1)) сравнений.", Rows: []sets.FigureRow{
							{Label: "Min-Heap", Cells: []any{1, 3, 2, 5, 9, 8}, Marked: []int{2}, Tree: true},
							{Label: "h", Cells: []any{1, 3, 2, 5, 9, 8}, Marked: []int{2}},
							{Label: "обмены / длина", Cells: []any{1, 6}},
						}},
					},
				},
				{
					ID:    "hq2-stones",
					Title: "Камни: 6 pop, 3 push, пик 4 элемента",
					Frames: []sets.FigureFrame{
						{Caption: "Для [6, 4, 9, 2] созданы 4 отрицательных веса, выполнен Heapify. Построение O(n), дополнительная память O(n). Счётчики pop и push относятся к циклу; два изображения показывают один h.", Rows: []sets.FigureRow{
							{Label: "Min-Heap", Cells: []any{-9, -4, -6, -2}, Marked: []int{0}, Tree: true},
							{Label: "h", Cells: []any{-9, -4, -6, -2}, Marked: []int{0}},
							{Label: "pop / push / длина", Cells: []any{0, 0, 4}},
						}},
						{Caption: "Два pop выбрали 9 и 6. Один push вернул вес 3 со знаком минус. Итого 2 pop, 1 push, длина 3; пик остаётся 4.", Rows: []sets.FigureRow{
							{Label: "Min-Heap", Cells: []any{-4, -2, -3}, Marked: []int{2}, Tree: true},
							{Label: "h", Cells: []any{-4, -2, -3}, Marked: []int{2}},
							{Label: "pop / push / длина", Cells: []any{2, 1, 3}},
						}},
						{Caption: "Выбраны 4 и 3, вернулся вес 1. Итого 4 pop, 2 push, длина 2. Каждый круг уменьшает длину хотя бы на 1.", Rows: []sets.FigureRow{
							{Label: "Min-Heap", Cells: []any{-2, -1}, Marked: []int{1}, Tree: true},
							{Label: "h", Cells: []any{-2, -1}, Marked: []int{1}},
							{Label: "pop / push / длина", Cells: []any{4, 2, 2}},
						}},
						{Caption: "Выбраны 2 и 1, вернулся вес 1. Всего 3 круга, 6 pop, 3 push, пик 4. Время O(n log(n + 1)), дополнительная память O(n), ответ - одно число 1.", Rows: []sets.FigureRow{
							{Label: "Min-Heap", Cells: []any{-1}, Marked: []int{0}, Tree: true},
							{Label: "h", Cells: []any{-1}, Marked: []int{0}},
							{Label: "pop / push / длина", Cells: []any{6, 3, 1}},
						}},
					},
				},
				{
					ID:    "hq2-top3",
					Title: "Min-Heap размера 3: сравнения и замены",
					Frames: []sets.FigureFrame{
						{Caption: "n = 7, k = 3. Добавили 7: один push, длина 1. Счётчики: push, замены, длина. Array a - вход; Min-Heap и h показывают одну рабочую структуру.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{7, 2, 9, 4, 6, 1, 8}, Pointers: &sets.Pointers{{Name: "i", Cell: 0}}},
							{Label: "Min-Heap", Cells: []any{7}, Marked: []int{0}, Tree: true},
							{Label: "h", Cells: []any{7}, Marked: []int{0}},
							{Label: "счётчики", Cells: []any{1, 0, 1}},
						}},
						{Caption: "Добавили 2: второй push, длина 2. Сравнения внутри push входят в O(log(k + 1)); отдельные сравнения с корнем начнутся после заполнения.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{7, 2, 9, 4, 6, 1, 8}, Pointers: &sets.Pointers{{Name: "i", Cell: 1}}},
							{Label: "Min-Heap", Cells: []any{2, 7}, Marked: []int{0}, Tree: true},
							{Label: "h", Cells: []any{2, 7}, Marked: []int{0}},
							{Label: "счётчики", Cells: []any{2, 0, 2}},
						}},
						{Caption: "Добавили 9: третий push. Min-Heap заполнен, пик 3 элемента; h[0] = 2 - наименьшее из трёх лучших.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{7, 2, 9, 4, 6, 1, 8}, Pointers: &sets.Pointers{{Name: "i", Cell: 2}}},
							{Label: "Min-Heap", Cells: []any{2, 7, 9}, Marked: []int{2}, Tree: true},
							{Label: "h", Cells: []any{2, 7, 9}, Marked: []int{2}},
							{Label: "счётчики", Cells: []any{3, 0, 3}},
						}},
						{Caption: "Первое внешнее сравнение: 4 > 2. Первая замена корня даёт [4, 7, 9]. Длина остаётся 3.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{7, 2, 9, 4, 6, 1, 8}, Pointers: &sets.Pointers{{Name: "i", Cell: 3}}},
							{Label: "Min-Heap", Cells: []any{4, 7, 9}, Marked: []int{0}, Tree: true},
							{Label: "h", Cells: []any{4, 7, 9}, Marked: []int{0}},
							{Label: "счётчики", Cells: []any{3, 1, 3}},
						}},
						{Caption: "Второе внешнее сравнение: 6 > 4. Вторая замена даёт [6, 7, 9]. Длина 3.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{7, 2, 9, 4, 6, 1, 8}, Pointers: &sets.Pointers{{Name: "i", Cell: 4}}},
							{Label: "Min-Heap", Cells: []any{6, 7, 9}, Marked: []int{0}, Tree: true},
							{Label: "h", Cells: []any{6, 7, 9}, Marked: []int{0}},
							{Label: "счётчики", Cells: []any{3, 2, 3}},
						}},
						{Caption: "Третье внешнее сравнение: 1 ≤ 6. Замены нет, длина 3. Посещение элемента не обязательно вызывает обновление.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{7, 2, 9, 4, 6, 1, 8}, Pointers: &sets.Pointers{{Name: "i", Cell: 5}}},
							{Label: "Min-Heap", Cells: []any{6, 7, 9}, Marked: []int{0}, Tree: true},
							{Label: "h", Cells: []any{6, 7, 9}, Marked: []int{0}},
							{Label: "счётчики", Cells: []any{3, 2, 3}},
						}},
						{Caption: "Четвёртое внешнее сравнение: 8 > 6. Третья замена даёт [7, 8, 9]. Всего 3 push, 4 внешних сравнения, 3 замены; пик 3, ответ 7. Время O(n log(k + 1)), дополнительная память O(k), ответ O(1).", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{7, 2, 9, 4, 6, 1, 8}, Pointers: &sets.Pointers{{Name: "i", Cell: 6}}},
							{Label: "Min-Heap", Cells: []any{7, 8, 9}, Marked: []int{1}, Tree: true},
							{Label: "h", Cells: []any{7, 8, 9}, Marked: []int{1}},
							{Label: "счётчики", Cells: []any{3, 3, 3}},
						}},
					},
				},
			},
		},
		{
			Title:   "Sliding Window и Monotonic Queue",
			Summary: "Monotonic Queue через Deque хранит номера кандидатов на максимум. Сравниваем работу и память с Heap, ищем окна с max − min ≤ limit и обновляем сумму вместе с границами Sliding Window.",
			Lesson:  lesson3,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/sliding-window-maximum/",
					Number:     number(239),
					Title:      "Sliding Window Maximum",
					Difficulty: "hard",
					Topics:     []string{"Array", "Queue", "Sliding Window", "Heap (Priority Queue)", "Monotonic Queue", "Range Minimum/Maximum Query"},
					Note:       "Monotonic Deque хранит номера: сзади удаляйте меньшие или равные значения, спереди истёкшие номера. n - длина nums, k - ширина окна. Время: O(n) в худшем случае, n добавлений и не больше n удалений. Память: O(k) дополнительно; ответ отдельно O(n − k + 1). Brute Force посещает (n − k + 1)k элементов; Heap с отложенным удалением требует O(n log(n + 1)) времени и O(n) рабочих записей.",
					Statement: `<p>Дан Array целых чисел <code>nums</code>. Sliding Window размера <code>k</code> движется от левого края Array к правому. Видны только <code>k</code> чисел в окне. Каждый раз окно сдвигается вправо на одну позицию.</p><p>Верните <em>максимум каждого окна</em>.</p><p>Метод класса <code>Solution</code>: <code>def maxSlidingWindow(self, nums: list[int], k: int) -&gt; list[int]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [1,3,-1,-3,5,3,6,7], k = 3
<strong>Выход:</strong> [3,3,5,5,6,7]
<strong>Пояснение:</strong> окна [1,3,-1], [3,-1,-3], [-1,-3,5], [-3,5,3], [5,3,6], [3,6,7] дают максимумы 3, 3, 5, 5, 6, 7.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [1], k = 1
<strong>Выход:</strong> [1]</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 10<sup>5</sup></code></li><li><code>-10<sup>4</sup> ≤ nums[i] ≤ 10<sup>4</sup></code></li><li><code>1 ≤ k ≤ nums.length</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "maxSlidingWindow"}, Cases: []sets.TestCase{
						{Input: "[[1,3,-1,-3,5,3,6,7],3]", Expected: text("[3,3,5,5,6,7]")},
						{Input: "[[1],1]", Expected: text("[1]")},
						{Input: "[[4,-2,9,1],4]", Expected: text("[9]")},
						{Input: "[[5,-1,3],1]", Expected: text("[5,-1,3]")},
						{Input: "[[7,7,7,2],2]", Expected: text("[7,7,7]")},
						{Input: "[[-10000,0,10000,-5],2]", Expected: text("[0,10000,10000]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/longest-continuous-subarray-with-absolute-diff-less-than-or-equal-to-limit/",
					Number:     number(1438),
					Title:      "Longest Continuous Subarray With Absolute Diff Less Than or Equal to Limit",
					Difficulty: "medium",
					Topics:     []string{"Array", "Queue", "Sliding Window", "Heap (Priority Queue)", "Ordered Set", "Monotonic Queue"},
					Note:       "Два Monotonic Deque дают максимум и минимум. Пока max − min > limit, сдвигайте l. n - длина nums, w - наибольшая проверяемая ширина до сжатия. Время: O(n) в худшем случае: 2n добавлений, не больше 2n удалений и n сдвигов l. Память: O(w) дополнительно, в худшем случае O(n); ответ отдельно O(1). Повторные max/min дают до O(n²) работы. При limit = 0 в окне остаются равные числа.",
					Statement: `<p>Даны Array целых чисел <code>nums</code> и целое число <code>limit</code>.</p><p>Верните <em>длину самого длинного непустого Subarray</em>, в котором модуль разности любых двух элементов не больше <code>limit</code>. Subarray состоит из соседних элементов Array.</p><p>Метод класса <code>Solution</code>: <code>def longestSubarray(self, nums: list[int], limit: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [8,2,4,7], limit = 4
<strong>Выход:</strong> 2
<strong>Пояснение:</strong> [2,4] и [4,7] подходят. В [8,2,4] разность 8 − 2 = 6, в [2,4,7] разность 7 − 2 = 5. Поэтому длина 3 и больше не подходит.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [10,1,2,4,7,2], limit = 5
<strong>Выход:</strong> 4
<strong>Пояснение:</strong> Subarray [2,4,7,2] подходит: max − min = 7 − 2 = 5. Добавить слева 1 нельзя: 7 − 1 = 6.</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> nums = [4,2,2,2,4,4,2,2], limit = 0
<strong>Выход:</strong> 3
<strong>Пояснение:</strong> при limit = 0 подходят только равные числа. Самый длинный Subarray [2,2,2], его длина 3.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 10<sup>5</sup></code></li><li><code>1 ≤ nums[i] ≤ 10<sup>9</sup></code></li><li><code>0 ≤ limit ≤ 10<sup>9</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "longestSubarray"}, Cases: []sets.TestCase{
						{Input: "[[8,2,4,7],4]", Expected: text("2")},
						{Input: "[[10,1,2,4,7,2],5]", Expected: text("4")},
						{Input: "[[4,2,2,2,4,4,2,2],0]", Expected: text("3")},
						{Input: "[[1],0]", Expected: text("1")},
						{Input: "[[1,1000000000],999999999]", Expected: text("2")},
						{Input: "[[1,5,6,7,8,10,6,5,6],4]", Expected: text("5")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/continuous-subarrays/",
					Number:     number(2762),
					Title:      "Continuous Subarrays",
					Difficulty: "medium",
					Topics:     []string{"Array", "Queue", "Sliding Window", "Heap (Priority Queue)", "Ordered Set", "Monotonic Queue"},
					Note:       "Для каждого r прибавляйте r − l + 1. n - длина nums, w - наибольшая ширина окна. Время: O(n) в худшем случае, границы и два Deque делают суммарно O(n) работы. Память: O(w) в общем случае; при целых значениях, limit = 2 и удалении равных сзади только O(1) живых записей. Ответ отдельно: одно целое до n(n + 1)/2. Для [5,4,2,4] нужны 4 прибавления, ответ 8; Brute Force проверяет 10 интервалов.",
					Statement: `<p>Дан Array целых чисел <code>nums</code> с нумерацией с 0. Subarray называется допустимым, если разность любых двух его элементов по модулю не больше 2.</p><p>Верните <em>общее число допустимых Subarray</em>. Subarray - непустая последовательность соседних элементов Array.</p><p>Метод класса <code>Solution</code>: <code>def continuousSubarrays(self, nums: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums = [5,4,2,4]
<strong>Выход:</strong> 8
<strong>Пояснение:</strong> длина 1: [5], [4], [2], [4]; длина 2: [5,4], [4,2], [2,4]; длина 3: [4,2,4]. Длина 4 не подходит. Всего 4 + 3 + 1 = 8.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums = [1,2,3]
<strong>Выход:</strong> 6
<strong>Пояснение:</strong> подходят [1], [2], [3], [1,2], [2,3], [1,2,3]. Всего 3 + 2 + 1 = 6.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ nums.length ≤ 10<sup>5</sup></code></li><li><code>1 ≤ nums[i] ≤ 10<sup>9</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "continuousSubarrays"}, Cases: []sets.TestCase{
						{Input: "[[5,4,2,4]]", Expected: text("8")},
						{Input: "[[1,2,3]]", Expected: text("6")},
						{Input: "[[7]]", Expected: text("1")},
						{Input: "[[1,1000000000,1]]", Expected: text("3")},
						{Input: "[[2,2,2,2]]", Expected: text("10")},
						{Input: "[[1,3,5,3,1]]", Expected: text("10")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/maximum-number-of-robots-within-budget/",
					Number:     number(2398),
					Title:      "Maximum Number of Robots Within Budget",
					Difficulty: "hard",
					Topics:     []string{"Array", "Binary Search", "Queue", "Sliding Window", "Heap (Priority Queue)", "Prefix Sum", "Monotonic Queue"},
					Note:       "Monotonic Deque даёт максимум chargeTimes; сумму runningCosts обновляйте при движении границ. Положительные стоимости позволяют сжимать окно, пока max + k · sum > budget. n - число роботов, w - наибольшая проверяемая ширина до сжатия. Время: O(n) в худшем случае; каждый номер добавляется и удаляется не больше раза. Память: O(w), в худшем случае O(n), ответ отдельно O(1). Повторный sum или max может добавить O(n²) работы.",
					Statement: `<p>Даны два Array положительных целых чисел <code>chargeTimes</code> и <code>runningCosts</code> длины <code>n</code>. Зарядка робота <code>i</code> стоит <code>chargeTimes[i]</code>, работа - <code>runningCosts[i]</code>. Дан бюджет <code>budget</code>.</p><p>Стоимость запуска <code>k</code> роботов подряд равна <code>max(chargeTimes) + k * sum(runningCosts)</code>. Максимум и сумма берутся только по выбранным роботам.</p><p>Верните <em>наибольшее число роботов подряд</em>, которых можно запустить в пределах бюджета.</p><p>Метод класса <code>Solution</code>: <code>def maximumRobots(self, chargeTimes: list[int], runningCosts: list[int], budget: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> chargeTimes = [3,6,1,3,4], runningCosts = [2,1,3,4,5], budget = 25
<strong>Выход:</strong> 3
<strong>Пояснение:</strong> первые три робота стоят max(3,6,1) + 3 * sum(2,1,3) = 6 + 3 * 6 = 24. Ни одно окно длины 4 не укладывается в бюджет.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> chargeTimes = [11,12,19], runningCosts = [10,8,7], budget = 19
<strong>Выход:</strong> 0
<strong>Пояснение:</strong> даже один робот стоит больше бюджета: 21, 20 или 26.</pre><h3>Ограничения</h3><ul><li><code>chargeTimes.length == runningCosts.length == n</code></li><li><code>1 ≤ n ≤ 5 · 10<sup>4</sup></code></li><li><code>1 ≤ chargeTimes[i], runningCosts[i] ≤ 10<sup>5</sup></code></li><li><code>1 ≤ budget ≤ 10<sup>15</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "maximumRobots"}, Cases: []sets.TestCase{
						{Input: "[[3,6,1,3,4],[2,1,3,4,5],25]", Expected: text("3")},
						{Input: "[[11,12,19],[10,8,7],19]", Expected: text("0")},
						{Input: "[[5],[5],10]", Expected: text("1")},
						{Input: "[[1,1,1,1],[1,1,1,1],17]", Expected: text("4")},
						{Input: "[[10,1,1,1],[1,1,1,1],13]", Expected: text("3")},
						{Input: "[[100000,100000],[100000,100000],1000000000000000]", Expected: text("2")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "hq3-heap-stale", Type: "choice", Prompt: "Максимум Sliding Window ищут через heapq: в Heap пары (−x, i). Номер вершины уже вне окна. Что делать, чтобы все шаги вместе стоили O(n log(n + 1))?", Options: []string{"удалить из Heap все пары, которые вышли из окна", "заново вызвать heapify", "удалять вершину, пока её номер вне окна"}, Answer: 2, Explain: "Старые пары ниже вершины пока не влияют на максимум. Каждая пара входит один раз и выходит не больше раза: O(n log(n + 1)) времени и до O(n) живых пар. heapq не имеет метода удаления по индексу исходного Array. Если каждый шаг фильтровать весь Heap и вызывать heapify, хранение O(k), но время O(nk). В отложенном варианте сам Heap может иметь размер n, поэтому один его проход нельзя безусловно оценить как O(k). Один heapify старые пары не удаляет."},
				{ID: "hq3-why-drop", Type: "choice", Prompt: "Пришло x, и Monotonic Deque удалил сзади число y ≤ x. Почему y больше не нужен для максимума?", Options: []string{"номер y уже вне окна", "x не меньше y и уйдёт из окна позже", "x стало максимумом окна"}, Answer: 1, Explain: "Пока y в окне, x тоже в окне, потому что x пришло позже. И x не меньше y. Номер y может ещё быть в окне, а x не обязательно максимум: в [9,1,5] число 5 удаляет 1, но максимум 9."},
				{ID: "hq3-store-index", Type: "choice", Prompt: "Что удобнее хранить в Monotonic Deque?", Options: []string{"номера элементов", "значения элементов", "пары (значение, длина окна)"}, Answer: 0, Explain: "По номеру видно, что элемент вышел из окна: dq[0] <= i − k. Вариант с одними значениями требует сохранять повторы и сравнивать первое значение с уходящим a[i − k]."},
				{ID: "hq3-front-check", Type: "code", Prompt: "Окно длины k = 2, a = [5,1,2,3]. Изначально dq и res пусты. В проверке спереди стоит знак <. Какое число код запишет в res на шаге i = 2?", Language: "python", Code: `for i, x in enumerate(a):
    while dq and a[dq[-1]] <= x:
        dq.pop()
    dq.append(i)
    if dq[0] < i - k:
        dq.popleft()
    if i >= k - 1:
        res.append(a[dq[0]])`, Options: []string{"2", "1", "5"}, Answer: 2, Explain: "На шаге 2 окно [1,2], его максимум 2. Сзади уходит номер 1, в Deque остаются [0,2]. Номер 0 вне окна, но 0 < 2 − 2 ложно, поэтому код пишет a[0] = 5. Нужна проверка <=. До этой записи сделаны 3 append и 1 pop; живы 2 номера. Правильный общий проход требует O(n) времени, O(k) рабочей памяти и O(n − k + 1) чисел ответа."},
				{ID: "hq3-deque-size", Type: "short", Prompt: "Monotonic Deque для максимума, a = [9,2,7,5,3,1,4], k = 4. Сколько номеров лежит в Deque после шага i = 6?", Answer: "2", Explain: "После шага 5 в Deque номера [2,3,4,5], значения [7,5,3,1]. На шаге 6 сзади уходят 1 и 3, добавляется номер 6, спереди уходит номер 2 ≤ 6 − 4. Остаются [3,6], значения [5,4]. Здесь 1 добавление и 3 удаления, после шага 2 живых номера."},
				{ID: "hq3-min-deque", Type: "choice", Prompt: "Monotonic Deque для минимума окна, a = [3,5,4,1,2], k = 3. Какие значения лежат в Deque после шага i = 2?", Options: []string{"5, 4", "3, 4", "3, 5, 4"}, Answer: 1, Explain: "Шаг 0: [3]. Шаг 1: [3,5]. Шаг 2: число 5 ≥ 4 удаляется сзади, 3 остаётся, добавляется 4. Итого 3 добавления, 1 удаление, 3 сравнения сзади; 2 живых номера. Пик тоже 2. Для всего входа время O(n), рабочая память O(k), ответ отдельно O(n − k + 1)."},
				{ID: "hq3-count", Type: "short", Prompt: "nums = [1,3,2,4]. Сколько Subarray с max − min ≤ 2?", Answer: "9", Explain: "Для каждой правой границы прибавьте r − l + 1: 1 + 2 + 3 + 3 = 9. На последнем шаге l сдвигается на 1. Четыре прибавления дают число ответов без хранения девяти интервалов; ответ занимает одно целое."},
				{ID: "hq3-robot-cost", Type: "short", Prompt: "chargeTimes = [4,2,7], runningCosts = [2,3,2]. Сколько стоит запуск всех трёх роботов по формуле max(chargeTimes) + k · sum(runningCosts)?", Answer: "28", Explain: "Максимум 7, сумма 7, k = 3. Цена 7 + 3 · 7 = 28. Один прямой расчёт для k роботов делает k − 1 сравнений максимума и k прибавлений к сумме с нуля: здесь 2 и 3. Время O(k), рабочая память O(1), ответ отдельно O(1). При движении окна сумму обновляем, а максимум берём из Deque."},
			},
			Links: []sets.Link{
				{Title: "Sliding Window", URL: "https://usaco.guide/gold/sliding-window", Source: "USACO Guide (англ.)"},
				{Title: "Minimum Stack и Minimum Queue", URL: "https://cp-algorithms.com/data_structures/stack_queue_modification.html#queue-modification-method-2", Source: "cp-algorithms (англ.)"},
				{Title: "Monotonic Queue для Sliding Window", URL: "https://labuladong.online/en/algo/data-structure/monotonic-queue/", Source: "labuladong (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "hq3-naive",
					Title: "Brute Force: a = [2,7,3,1,5,2,6,4], k = 3",
					Frames: []sets.FigureFrame{
						{Caption: "Окно [2,7,3]: максимум 7. Счётчики работы накопительные: посещения элементов и сравнения значений. Рабочая память хранит один текущий максимум; число элементов ответа считаем отдельно.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 7, 3, 1, 5, 2, 6, 4}, Marked: []int{1}, Window: &[2]int{0, 2}},
							{Label: "ответ", Cells: []any{7, "", "", "", "", ""}, Marked: []int{0}},
							{Label: "посещений", Cells: []any{3}},
							{Label: "сравнений", Cells: []any{2}},
							{Label: "макс (чисел)", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "Окно [7,3,1]: ещё 3 посещения и 2 сравнения. Максимум 7.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 7, 3, 1, 5, 2, 6, 4}, Marked: []int{1}, Window: &[2]int{1, 3}},
							{Label: "ответ", Cells: []any{7, 7, "", "", "", ""}, Marked: []int{1}},
							{Label: "посещений", Cells: []any{6}},
							{Label: "сравнений", Cells: []any{4}},
							{Label: "макс (чисел)", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{2}},
						}},
						{Caption: "Окно [3,1,5]: максимум 5.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 7, 3, 1, 5, 2, 6, 4}, Marked: []int{4}, Window: &[2]int{2, 4}},
							{Label: "ответ", Cells: []any{7, 7, 5, "", "", ""}, Marked: []int{2}},
							{Label: "посещений", Cells: []any{9}},
							{Label: "сравнений", Cells: []any{6}},
							{Label: "макс (чисел)", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{3}},
						}},
						{Caption: "Окно [1,5,2]: максимум 5.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 7, 3, 1, 5, 2, 6, 4}, Marked: []int{4}, Window: &[2]int{3, 5}},
							{Label: "ответ", Cells: []any{7, 7, 5, 5, "", ""}, Marked: []int{3}},
							{Label: "посещений", Cells: []any{12}},
							{Label: "сравнений", Cells: []any{8}},
							{Label: "макс (чисел)", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{4}},
						}},
						{Caption: "Окно [5,2,6]: максимум 6.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 7, 3, 1, 5, 2, 6, 4}, Marked: []int{6}, Window: &[2]int{4, 6}},
							{Label: "ответ", Cells: []any{7, 7, 5, 5, 6, ""}, Marked: []int{4}},
							{Label: "посещений", Cells: []any{15}},
							{Label: "сравнений", Cells: []any{10}},
							{Label: "макс (чисел)", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{5}},
						}},
						{Caption: "Окно [2,6,4]: максимум 6. Всего 6 · 3 = 18 посещений и 6 · 2 = 12 сравнений. Срезы не создаём; индексы занимают постоянную память.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 7, 3, 1, 5, 2, 6, 4}, Marked: []int{6}, Window: &[2]int{5, 7}},
							{Label: "ответ", Cells: []any{7, 7, 5, 5, 6, 6}, Marked: []int{5}},
							{Label: "посещений", Cells: []any{18}},
							{Label: "сравнений", Cells: []any{12}},
							{Label: "макс (чисел)", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{6}},
						}},
					},
				},
				{
					ID:    "hq3-deque",
					Title: "Monotonic Deque: a = [2,7,3,1,5,2,6,4], k = 3",
					Frames: []sets.FigureFrame{
						{Caption: "x = 2 входит в пустой Deque. Считаем вызовы append, удаления через pop/popleft и сравнения значений сзади. Счётчики накопительные. Живые номера Deque и числа ответа показаны после шага.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 7, 3, 1, 5, 2, 6, 4}, Pointers: &sets.Pointers{{Name: "i", Cell: 0}}, Window: &[2]int{0, 0}},
							{Label: "номера", Cells: []any{0}, Marked: []int{0}},
							{Label: "значения", Cells: []any{2}, Marked: []int{0}},
							{Label: "ответ", Cells: []any{"", "", "", "", "", ""}},
							{Label: "append", Cells: []any{1}},
							{Label: "удалений", Cells: []any{0}},
							{Label: "сравнений", Cells: []any{0}},
							{Label: "живых Deque", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{0}},
						}},
						{Caption: "x = 7: число 2 удаляется сзади. Приглушённые ячейки показывают удалённые записи и в живую память не входят.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 7, 3, 1, 5, 2, 6, 4}, Pointers: &sets.Pointers{{Name: "i", Cell: 1}}, Window: &[2]int{0, 1}},
							{Label: "номера", Cells: []any{0, 1}, Marked: []int{1}, Dimmed: []int{0}},
							{Label: "значения", Cells: []any{2, 7}, Marked: []int{1}, Dimmed: []int{0}},
							{Label: "ответ", Cells: []any{"", "", "", "", "", ""}},
							{Label: "append", Cells: []any{2}},
							{Label: "удалений", Cells: []any{1}},
							{Label: "сравнений", Cells: []any{1}},
							{Label: "живых Deque", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{0}},
						}},
						{Caption: "x = 3 добавляется сзади. Номера [1,2], значения [7,3], первый ответ 7.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 7, 3, 1, 5, 2, 6, 4}, Pointers: &sets.Pointers{{Name: "i", Cell: 2}}, Window: &[2]int{0, 2}},
							{Label: "номера", Cells: []any{1, 2}, Marked: []int{1}},
							{Label: "значения", Cells: []any{7, 3}, Marked: []int{1}},
							{Label: "ответ", Cells: []any{7, "", "", "", "", ""}, Marked: []int{0}},
							{Label: "append", Cells: []any{3}},
							{Label: "удалений", Cells: []any{1}},
							{Label: "сравнений", Cells: []any{2}},
							{Label: "живых Deque", Cells: []any{2}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "x = 1 добавляется сзади. Номера [1,2,3], значения [7,3,1]. Достигнут пик 3 живых номера.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 7, 3, 1, 5, 2, 6, 4}, Pointers: &sets.Pointers{{Name: "i", Cell: 3}}, Window: &[2]int{1, 3}},
							{Label: "номера", Cells: []any{1, 2, 3}, Marked: []int{2}},
							{Label: "значения", Cells: []any{7, 3, 1}, Marked: []int{2}},
							{Label: "ответ", Cells: []any{7, 7, "", "", "", ""}, Marked: []int{1}},
							{Label: "append", Cells: []any{4}},
							{Label: "удалений", Cells: []any{1}},
							{Label: "сравнений", Cells: []any{3}},
							{Label: "живых Deque", Cells: []any{3}},
							{Label: "чисел ответа", Cells: []any{2}},
						}},
						{Caption: "x = 5: сзади уходят 1 и 3. Сравнение 7 ≤ 5 ложно. Номер 1 ≤ 4 − 3 уходит спереди. Остаётся номер 4, ответ 5.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 7, 3, 1, 5, 2, 6, 4}, Pointers: &sets.Pointers{{Name: "i", Cell: 4}}, Window: &[2]int{2, 4}},
							{Label: "номера", Cells: []any{1, 2, 3, 4}, Marked: []int{3}, Dimmed: []int{0, 1, 2}},
							{Label: "значения", Cells: []any{7, 3, 1, 5}, Marked: []int{3}, Dimmed: []int{0, 1, 2}},
							{Label: "ответ", Cells: []any{7, 7, 5, "", "", ""}, Marked: []int{2}},
							{Label: "append", Cells: []any{5}},
							{Label: "удалений", Cells: []any{4}},
							{Label: "сравнений", Cells: []any{6}},
							{Label: "живых Deque", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{3}},
						}},
						{Caption: "x = 2 добавляется сзади. Номера [4,5], значения [5,2], ответ 5.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 7, 3, 1, 5, 2, 6, 4}, Pointers: &sets.Pointers{{Name: "i", Cell: 5}}, Window: &[2]int{3, 5}},
							{Label: "номера", Cells: []any{4, 5}, Marked: []int{1}},
							{Label: "значения", Cells: []any{5, 2}, Marked: []int{1}},
							{Label: "ответ", Cells: []any{7, 7, 5, 5, "", ""}, Marked: []int{3}},
							{Label: "append", Cells: []any{6}},
							{Label: "удалений", Cells: []any{4}},
							{Label: "сравнений", Cells: []any{7}},
							{Label: "живых Deque", Cells: []any{2}},
							{Label: "чисел ответа", Cells: []any{4}},
						}},
						{Caption: "x = 6: сзади уходят 2 и 5. Остаётся номер 6, ответ 6.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 7, 3, 1, 5, 2, 6, 4}, Pointers: &sets.Pointers{{Name: "i", Cell: 6}}, Window: &[2]int{4, 6}},
							{Label: "номера", Cells: []any{4, 5, 6}, Marked: []int{2}, Dimmed: []int{0, 1}},
							{Label: "значения", Cells: []any{5, 2, 6}, Marked: []int{2}, Dimmed: []int{0, 1}},
							{Label: "ответ", Cells: []any{7, 7, 5, 5, 6, ""}, Marked: []int{4}},
							{Label: "append", Cells: []any{7}},
							{Label: "удалений", Cells: []any{6}},
							{Label: "сравнений", Cells: []any{9}},
							{Label: "живых Deque", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{5}},
						}},
						{Caption: "x = 4: номера [6,7], значения [6,4]. Всего 8 добавлений + 6 удалений = 14 ≤ 2n. Ещё были 10 сравнений сзади, 8 проверок срока и 6 записей ответа.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{2, 7, 3, 1, 5, 2, 6, 4}, Pointers: &sets.Pointers{{Name: "i", Cell: 7}}, Window: &[2]int{5, 7}},
							{Label: "номера", Cells: []any{6, 7}, Marked: []int{1}},
							{Label: "значения", Cells: []any{6, 4}, Marked: []int{1}},
							{Label: "ответ", Cells: []any{7, 7, 5, 5, 6, 6}, Marked: []int{5}},
							{Label: "append", Cells: []any{8}},
							{Label: "удалений", Cells: []any{6}},
							{Label: "сравнений", Cells: []any{10}},
							{Label: "живых Deque", Cells: []any{2}},
							{Label: "пик Deque", Cells: []any{3}},
							{Label: "чисел ответа", Cells: []any{6}},
						}},
					},
				},
				{
					ID:    "hq3-two-deques",
					Title: "Два Monotonic Deque: nums = [4,8,5,6,7,1,2], limit = 3",
					Frames: []sets.FigureFrame{
						{Caption: "x = 4 входит в оба Deque. Окно [4], лучшая длина 1. Счётчики накапливают вызовы обоих Deque и сдвиги l. Живые записи и их пик считаем для двух Deque вместе, после сжатия; ответ хранит одно число.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{4, 8, 5, 6, 7, 1, 2}, Pointers: &sets.Pointers{{Name: "l", Cell: 0}, {Name: "r", Cell: 0}}, Window: &[2]int{0, 0}},
							{Label: "макс", Cells: []any{4}, Marked: []int{0}},
							{Label: "мин", Cells: []any{4}, Marked: []int{0}},
							{Label: "append", Cells: []any{2}},
							{Label: "удалений", Cells: []any{0}},
							{Label: "сдвигов l", Cells: []any{0}},
							{Label: "живых вместе", Cells: []any{2}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "x = 8: число 4 уходит сзади из Deque максимума. Разность 4 > 3, l = 1; число 4 уходит спереди из Deque минимума. Окно [8].", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{4, 8, 5, 6, 7, 1, 2}, Pointers: &sets.Pointers{{Name: "l", Cell: 1}, {Name: "r", Cell: 1}}, Window: &[2]int{1, 1}},
							{Label: "макс", Cells: []any{4, 8}, Marked: []int{1}, Dimmed: []int{0}},
							{Label: "мин", Cells: []any{4, 8}, Marked: []int{1}, Dimmed: []int{0}},
							{Label: "append", Cells: []any{4}},
							{Label: "удалений", Cells: []any{2}},
							{Label: "сдвигов l", Cells: []any{1}},
							{Label: "живых вместе", Cells: []any{2}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "x = 5: из Deque минимума сзади уходит 8. Окно [8,5], разность 3, лучшая длина 2.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{4, 8, 5, 6, 7, 1, 2}, Pointers: &sets.Pointers{{Name: "l", Cell: 1}, {Name: "r", Cell: 2}}, Window: &[2]int{1, 2}},
							{Label: "макс", Cells: []any{8, 5}, Marked: []int{1}},
							{Label: "мин", Cells: []any{8, 5}, Marked: []int{1}, Dimmed: []int{0}},
							{Label: "append", Cells: []any{6}},
							{Label: "удалений", Cells: []any{3}},
							{Label: "сдвигов l", Cells: []any{1}},
							{Label: "живых вместе", Cells: []any{3}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "x = 6: из Deque максимума сзади уходит 5. Окно [8,5,6], разность 3, лучшая длина 3.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{4, 8, 5, 6, 7, 1, 2}, Pointers: &sets.Pointers{{Name: "l", Cell: 1}, {Name: "r", Cell: 3}}, Window: &[2]int{1, 3}},
							{Label: "макс", Cells: []any{8, 5, 6}, Marked: []int{2}, Dimmed: []int{1}},
							{Label: "мин", Cells: []any{5, 6}, Marked: []int{1}},
							{Label: "append", Cells: []any{8}},
							{Label: "удалений", Cells: []any{4}},
							{Label: "сдвигов l", Cells: []any{1}},
							{Label: "живых вместе", Cells: []any{4}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "x = 7: из Deque максимума сзади уходит 6. Окно [8,5,6,7], лучшая длина 4. Пик двух Deque вместе - 5 номеров.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{4, 8, 5, 6, 7, 1, 2}, Pointers: &sets.Pointers{{Name: "l", Cell: 1}, {Name: "r", Cell: 4}}, Window: &[2]int{1, 4}},
							{Label: "макс", Cells: []any{8, 6, 7}, Marked: []int{2}, Dimmed: []int{1}},
							{Label: "мин", Cells: []any{5, 6, 7}, Marked: []int{2}},
							{Label: "append", Cells: []any{10}},
							{Label: "удалений", Cells: []any{5}},
							{Label: "сдвигов l", Cells: []any{1}},
							{Label: "живых вместе", Cells: []any{5}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "x = 1: сзади из Deque минимума уходят 7, 6, 5. l сдвигается с 1 до 5, спереди из Deque максимума уходят 8 и 7. Окно [1], ответ остаётся 4.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{4, 8, 5, 6, 7, 1, 2}, Pointers: &sets.Pointers{{Name: "l", Cell: 5}, {Name: "r", Cell: 5}}, Window: &[2]int{5, 5}},
							{Label: "макс", Cells: []any{8, 7, 1}, Marked: []int{2}, Dimmed: []int{0, 1}},
							{Label: "мин", Cells: []any{5, 6, 7, 1}, Marked: []int{3}, Dimmed: []int{0, 1, 2}},
							{Label: "append", Cells: []any{12}},
							{Label: "удалений", Cells: []any{10}},
							{Label: "сдвигов l", Cells: []any{5}},
							{Label: "живых вместе", Cells: []any{2}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "x = 2: из Deque максимума сзади уходит 1. Окно [1,2], ответ 4. Условие max − min проверено 7 + 5 = 12 раз за весь проход.", Rows: []sets.FigureRow{
							{Label: "nums", Cells: []any{4, 8, 5, 6, 7, 1, 2}, Pointers: &sets.Pointers{{Name: "l", Cell: 5}, {Name: "r", Cell: 6}}, Window: &[2]int{5, 6}},
							{Label: "макс", Cells: []any{1, 2}, Marked: []int{1}, Dimmed: []int{0}},
							{Label: "мин", Cells: []any{1, 2}, Marked: []int{1}},
							{Label: "append", Cells: []any{14}},
							{Label: "удалений", Cells: []any{11}},
							{Label: "сдвигов l", Cells: []any{5}},
							{Label: "живых вместе", Cells: []any{3}},
							{Label: "пик вместе", Cells: []any{5}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
					},
				},
				{
					ID:    "hq3-robots",
					Title: "Роботы, budget = 25: Monotonic Deque и сумма окна",
					Frames: []sets.FigureFrame{
						{Caption: "Окно 0..0: 3 + 1 · 2 = 5 ≤ 25. Счётчики вызовов накопительные; «обновл. суммы» считает прибавления и вычитания. Память после сжатия: живые номера Deque, одно число суммы и одно число ответа.", Rows: []sets.FigureRow{
							{Label: "charge", Cells: []any{3, 6, 1, 3, 4}, Pointers: &sets.Pointers{{Name: "l", Cell: 0}, {Name: "r", Cell: 0}}, Window: &[2]int{0, 0}},
							{Label: "cost", Cells: []any{2, 1, 3, 4, 5}, Window: &[2]int{0, 0}},
							{Label: "макс", Cells: []any{3}, Marked: []int{0}},
							{Label: "append", Cells: []any{1}},
							{Label: "удалений", Cells: []any{0}},
							{Label: "обновл. суммы", Cells: []any{1}},
							{Label: "живых Deque", Cells: []any{1}},
							{Label: "чисел суммы", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "Из Deque сзади уходит 3. Окно 0..1: 6 + 2 · 3 = 12 ≤ 25, лучшая длина 2.", Rows: []sets.FigureRow{
							{Label: "charge", Cells: []any{3, 6, 1, 3, 4}, Pointers: &sets.Pointers{{Name: "l", Cell: 0}, {Name: "r", Cell: 1}}, Window: &[2]int{0, 1}},
							{Label: "cost", Cells: []any{2, 1, 3, 4, 5}, Window: &[2]int{0, 1}},
							{Label: "макс", Cells: []any{3, 6}, Marked: []int{1}, Dimmed: []int{0}},
							{Label: "append", Cells: []any{2}},
							{Label: "удалений", Cells: []any{1}},
							{Label: "обновл. суммы", Cells: []any{2}},
							{Label: "живых Deque", Cells: []any{1}},
							{Label: "чисел суммы", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "Окно 0..2: 6 + 3 · 6 = 24 ≤ 25, лучшая длина 3. Пик Deque - 2 номера.", Rows: []sets.FigureRow{
							{Label: "charge", Cells: []any{3, 6, 1, 3, 4}, Pointers: &sets.Pointers{{Name: "l", Cell: 0}, {Name: "r", Cell: 2}}, Window: &[2]int{0, 2}},
							{Label: "cost", Cells: []any{2, 1, 3, 4, 5}, Window: &[2]int{0, 2}},
							{Label: "макс", Cells: []any{6, 1}, Marked: []int{1}},
							{Label: "append", Cells: []any{3}},
							{Label: "удалений", Cells: []any{1}},
							{Label: "обновл. суммы", Cells: []any{3}},
							{Label: "живых Deque", Cells: []any{2}},
							{Label: "чисел суммы", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "Сзади уходит 1. Окна 0..3 и 1..3 стоят 46 и 30, оба выше 25. Два вычитания дают l = 2, число 6 уходит спереди. Цена окна 2..3: 3 + 2 · 7 = 17.", Rows: []sets.FigureRow{
							{Label: "charge", Cells: []any{3, 6, 1, 3, 4}, Pointers: &sets.Pointers{{Name: "l", Cell: 2}, {Name: "r", Cell: 3}}, Window: &[2]int{2, 3}},
							{Label: "cost", Cells: []any{2, 1, 3, 4, 5}, Window: &[2]int{2, 3}},
							{Label: "макс", Cells: []any{6, 1, 3}, Marked: []int{2}, Dimmed: []int{0, 1}},
							{Label: "append", Cells: []any{4}},
							{Label: "удалений", Cells: []any{3}},
							{Label: "обновл. суммы", Cells: []any{6}},
							{Label: "живых Deque", Cells: []any{1}},
							{Label: "чисел суммы", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "Сзади уходит 3. Цена окна 2..4 равна 40 > 25. После вычитания cost[2] получаем l = 3, цена 4 + 2 · 9 = 22. Ответ 3; всего 8 проверок цены.", Rows: []sets.FigureRow{
							{Label: "charge", Cells: []any{3, 6, 1, 3, 4}, Pointers: &sets.Pointers{{Name: "l", Cell: 3}, {Name: "r", Cell: 4}}, Window: &[2]int{3, 4}},
							{Label: "cost", Cells: []any{2, 1, 3, 4, 5}, Window: &[2]int{3, 4}},
							{Label: "макс", Cells: []any{3, 4}, Marked: []int{1}, Dimmed: []int{0}},
							{Label: "append", Cells: []any{5}},
							{Label: "удалений", Cells: []any{4}},
							{Label: "обновл. суммы", Cells: []any{8}},
							{Label: "живых Deque", Cells: []any{1}},
							{Label: "пик Deque", Cells: []any{2}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
					},
				},
			},
		},
		{
			Title:   "Heap и Greedy",
			Summary: "Heap выбирает лучшее из доступного сейчас. K-Way Merge, порядок задач, удаление худшего кандидата и выбор k пар: считаем операции, построение структуры, Sorting и память результата.",
			Lesson:  lesson4,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/kth-smallest-element-in-a-sorted-matrix/",
					Number:     number(378),
					Title:      "Kth Smallest Element in a Sorted Matrix",
					Difficulty: "medium",
					Topics:     []string{"Array", "Binary Search", "Sorting", "Heap (Priority Queue)", "Matrix"},
					Note:       "K-Way Merge строк Matrix: Heap хранит (число, строка, столбец). r - число строк, k - нужный порядковый номер. Время: O(r + k log(r + 1)) в худшем случае после Heapify и k удалений; чтение всех r² чисел отдельно O(r²). Память: O(r) дополнительно; ответ отдельно O(1). Sorting всех r² чисел требует O(r² log(r + 1)) времени и O(r²) рабочих записей. Повторы считаются.",
					Statement: `<p>Дана Matrix <code>matrix</code> размера <code>n x n</code>. Каждая строка и каждый столбец отсортированы по неубыванию. Верните <em><code>k</code>-й наименьший элемент</em>. Повторяющиеся значения считаются отдельно.</p><p>Найдите решение с дополнительной памятью меньше O(n<sup>2</sup>).</p><p>Метод класса <code>Solution</code>: <code>def kthSmallest(self, matrix: list[list[int]], k: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> matrix = [[1,5,9],[10,11,13],[12,13,15]], k = 8
<strong>Выход:</strong> 13
<strong>Пояснение:</strong> порядок чисел [1,5,9,10,11,12,13,13,15], восьмое число 13.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> matrix = [[-5]], k = 1
<strong>Выход:</strong> -5</pre><h3>Ограничения</h3><ul><li><code>n == matrix.length == matrix[i].length</code></li><li><code>1 ≤ n ≤ 300</code></li><li><code>-10<sup>9</sup> ≤ matrix[i][j] ≤ 10<sup>9</sup></code></li><li>Все строки и столбцы отсортированы по неубыванию.</li><li><code>1 ≤ k ≤ n<sup>2</sup></code></li></ul><p><strong>Дополнительный вопрос:</strong> можно ли решить задачу с O(1) рабочей памяти? Можно ли получить время O(n)? Последняя оценка требует более сложного алгоритма.</p>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "kthSmallest"}, Cases: []sets.TestCase{
						{Input: "[[[1,5,9],[10,11,13],[12,13,15]],8]", Expected: text("13")},
						{Input: "[[[-5]],1]", Expected: text("-5")},
						{Input: "[[[1,2],[1,3]],2]", Expected: text("1")},
						{Input: "[[[-10,-5,0],[-8,-3,2],[-6,1,9]],5]", Expected: text("-3")},
						{Input: "[[[-1000000000,1000000000],[1000000000,1000000000]],2]", Expected: text("1000000000")},
						{Input: "[[[1,3,5],[2,4,6],[7,8,9]],6]", Expected: text("6")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/maximum-subsequence-score/",
					Number:     number(2542),
					Title:      "Maximum Subsequence Score",
					Difficulty: "medium",
					Topics:     []string{"Array", "Greedy", "Sorting", "Heap (Priority Queue)"},
					Note:       "Sorting пар по nums2 от большего к меньшему, Min-Heap из k наибольших nums1 и сохранённая сумма. n - число пар, k - размер выбора. Время: O(n log(n + 1)) в худшем случае, включая Sorting и O(n log(k + 1)) для Heap. Память: O(n + k) = O(n) дополнительно; ответ отдельно O(1). Новый sum(h) при каждом полном Heap добавляет (n − k + 1)k посещений: общая оценка O(n log(n + 1) + nk).",
					Statement: `<p>Даны два Array неотрицательных целых чисел <code>nums1</code> и <code>nums2</code> длины <code>n</code> и положительное <code>k</code>. Выберите Subsequence из <code>k</code> разных номеров. Subsequence сохраняет исходный порядок, но номера могут идти не подряд.</p><p>Счёт равен сумме выбранных <code>nums1</code>, умноженной на минимум выбранных <code>nums2</code>. Верните <em>наибольший возможный счёт</em>.</p><p>Метод класса <code>Solution</code>: <code>def maxScore(self, nums1: list[int], nums2: list[int], k: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> nums1 = [1,3,3,2], nums2 = [2,1,3,4], k = 3
<strong>Выход:</strong> 12
<strong>Пояснение:</strong> варианты номеров и счёт:
0,1,2: (1+3+3) * min(2,1,3) = 7;
0,1,3: (1+3+2) * min(2,1,4) = 6;
0,2,3: (1+3+2) * min(2,3,4) = 12;
1,2,3: (3+3+2) * min(1,3,4) = 8.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> nums1 = [4,2,3,1,1], nums2 = [7,5,10,9,6], k = 1
<strong>Выход:</strong> 30
<strong>Пояснение:</strong> номер 2 даёт 3 * 10 = 30, это наибольший счёт.</pre><h3>Ограничения</h3><ul><li><code>n == nums1.length == nums2.length</code></li><li><code>1 ≤ n ≤ 10<sup>5</sup></code></li><li><code>0 ≤ nums1[i], nums2[j] ≤ 10<sup>5</sup></code></li><li><code>1 ≤ k ≤ n</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "maxScore"}, Cases: []sets.TestCase{
						{Input: "[[1,3,3,2],[2,1,3,4],3]", Expected: text("12")},
						{Input: "[[4,2,3,1,1],[7,5,10,9,6],1]", Expected: text("30")},
						{Input: "[[1,2,3],[3,2,1],3]", Expected: text("6")},
						{Input: "[[0,0],[5,5],1]", Expected: text("0")},
						{Input: "[[5,1,1],[1,5,5],2]", Expected: text("10")},
						{Input: "[[100000],[100000],1]", Expected: text("10000000000")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1526/C2",
					Title:      "Potions (Hard Version)",
					Difficulty: "medium",
					Topics:     []string{"Array", "Greedy", "Heap (Priority Queue)"},
					Note:       "Принимайте каждое зелье и добавляйте его в Min-Heap. Если здоровье отрицательно, удалите самое вредное принятое зелье; одного удаления достаточно. n - число зелий. Время: O(n log(n + 1)) в худшем случае, n добавлений и не больше n удалений. Память: O(n) дополнительно; ответ отдельно O(1). Brute Force по маскам с проверкой порядка требует O(n · 2ⁿ) времени.",
					Statement: `<p>Это сложная версия задачи с <code>n ≤ 200000</code>. В ряд стоят <code>n</code> зелий. Зелье <code>i</code> меняет здоровье на <code>a<sub>i</sub></code>, в том числе может уменьшить его.</p><p>Здоровье начинается с 0. Вы идёте слева направо. Каждое зелье можно выпить или пропустить. Здоровье всегда должно оставаться неотрицательным. Найдите наибольшее число зелий, которые можно выпить.</p><h3>Входные данные</h3><p>В первой строке целое <code>n</code>. Во второй строке <code>n</code> целых изменений здоровья <code>a<sub>1</sub>, ..., a<sub>n</sub></code>.</p><h3>Выходные данные</h3><p>Выведите наибольшее число принятых зелий.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>6
4 -4 1 -3 1 -3</pre><p><strong>Выходные данные:</strong></p><pre>5</pre><p><strong>Пояснение:</strong> можно выпить зелья 1, 3, 4, 5 и 6. Здоровье будет 4, 5, 2, 3, 0. Все 6 выпить нельзя.</p><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 200000</code></li><li><code>-10<sup>9</sup> ≤ a<sub>i</sub> ≤ 10<sup>9</sup></code></li><li>Ограничение по времени: 1 секунда на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: `6
4 -4 1 -3 1 -3
`, Expected: text(`5
`)},
						{Input: `1
-5
`, Expected: text(`0
`)},
						{Input: `1
0
`, Expected: text(`1
`)},
						{Input: `5
5 -5 -1 -1 -1
`, Expected: text(`4
`)},
						{Input: `3
1000000000 -1000000000 -1000000000
`, Expected: text(`2
`)},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/single-threaded-cpu/",
					Number:     number(1834),
					Title:      "Single-Threaded CPU",
					Difficulty: "medium",
					Topics:     []string{"Array", "Sorting", "Heap (Priority Queue)"},
					Note:       "Sorting по приходу, Min-Heap из (длительность, номер) доступных задач. Если Heap пуст, перенесите время на следующий приход. n - число задач. Время: O(n log(n + 1)) в худшем случае: Sorting, n добавлений и n удалений. Память: O(n) дополнительно для порядка прихода, Sorting и Heap; ответ отдельно O(n) номеров. Ожидание по одному такту добавляет H шагов для H единиц простоя.",
					Statement: `<p>Даны <code>n</code> задач с номерами <code>0..n - 1</code> в виде двумерного Array <code>tasks</code>. Запись <code>tasks[i] = [enqueueTime<sub>i</sub>, processingTime<sub>i</sub>]</code> задаёт время прихода и длительность задачи.</p><p>Однопоточный процессор выполняет не больше одной задачи за раз. Когда он свободен, выбирает самую короткую из доступных задач; при равенстве длительностей - меньший номер. Если доступных задач нет, ждёт. Начатую задачу выполняет до конца без остановок.</p><p>Верните <em>порядок номеров выполненных задач</em>.</p><p>Метод класса <code>Solution</code>: <code>def getOrder(self, tasks: list[list[int]]) -&gt; list[int]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> tasks = [[1,2],[2,4],[3,2],[4,1]]
<strong>Выход:</strong> [0,2,3,1]
<strong>Пояснение:</strong> задача 0 работает с 1 до 3, задача 2 с 3 до 5, задача 3 с 5 до 6, задача 1 с 6 до 10.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> tasks = [[7,10],[7,12],[7,5],[7,4],[7,2]]
<strong>Выход:</strong> [4,3,2,0,1]
<strong>Пояснение:</strong> все задачи пришли в 7. Длительности по порядку 2, 4, 5, 10, 12. Моменты завершения 9, 13, 18, 28, 40.</pre><h3>Ограничения</h3><ul><li><code>tasks.length == n</code></li><li><code>1 ≤ n ≤ 10<sup>5</sup></code></li><li><code>1 ≤ enqueueTime<sub>i</sub>, processingTime<sub>i</sub> ≤ 10<sup>9</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "getOrder"}, Cases: []sets.TestCase{
						{Input: "[[[1,2],[2,4],[3,2],[4,1]]]", Expected: text("[0,2,3,1]")},
						{Input: "[[[7,10],[7,12],[7,5],[7,4],[7,2]]]", Expected: text("[4,3,2,0,1]")},
						{Input: "[[[5,5]]]", Expected: text("[0]")},
						{Input: "[[[1,3],[10,2],[10,2]]]", Expected: text("[0,1,2]")},
						{Input: "[[[3,1],[1,1],[2,1]]]", Expected: text("[1,2,0]")},
						{Input: "[[[1000000000,1000000000],[1,1000000000]]]", Expected: text("[1,0]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/total-cost-to-hire-k-workers/",
					Number:     number(2462),
					Title:      "Total Cost to Hire K Workers",
					Difficulty: "medium",
					Topics:     []string{"Array", "Two Pointers", "Heap (Priority Queue)", "Simulation"},
					Note:       "Два Min-Heap содержат доступных работников с концов. n - число цен, c = candidates, h = min(n, 2c), k - число наймов. Heapify начальных групп, затем пополнение с выбранной стороны при l ≤ r. Время: O(h + k log(h + 1)) в худшем случае; чтение входа отдельно O(n). Повторные вставки при построении стоят O(h log(h + 1)). Память: O(h) дополнительно; ответ отдельно O(1). Общий Sorting неверен, если не все работники доступны.",
					Statement: `<p>Дан Array целых цен <code>costs</code>. Нужно нанять ровно <code>k</code> работников за <code>k</code> раундов.</p><p>В каждом раунде доступны первые <code>candidates</code> и последние <code>candidates</code> оставшихся работников. Выберите минимальную цену; при равенстве - меньший номер. Если группы пересекаются, работник участвует в выборе один раз. Каждый работник может быть нанят только один раз. Порядок оставшихся сохраняется.</p><p>Верните <em>общую цену найма</em>.</p><p>Метод класса <code>Solution</code>: <code>def totalCost(self, costs: list[int], k: int, candidates: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> costs = [17,12,10,2,7,2,11,20,8], k = 3, candidates = 4
<strong>Выход:</strong> 11
<strong>Пояснение:</strong> по исходным номерам выбираем 3 с ценой 2, затем 5 с ценой 2, затем 4 с ценой 7. Цена 2 + 2 + 7 = 11.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> costs = [1,2,4,1], k = 3, candidates = 3
<strong>Выход:</strong> 4
<strong>Пояснение:</strong> все работники доступны. По исходным номерам выбираем 0, 3, 1. Цена 1 + 1 + 2 = 4.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ costs.length ≤ 10<sup>5</sup></code></li><li><code>1 ≤ costs[i] ≤ 10<sup>5</sup></code></li><li><code>1 ≤ k, candidates ≤ costs.length</code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "totalCost"}, Cases: []sets.TestCase{
						{Input: "[[17,12,10,2,7,2,11,20,8],3,4]", Expected: text("11")},
						{Input: "[[1,2,4,1],3,3]", Expected: text("4")},
						{Input: "[[5],1,1]", Expected: text("5")},
						{Input: "[[2,1,9,9,5,2],2,1]", Expected: text("3")},
						{Input: "[[1,2],2,2]", Expected: text("3")},
						{Input: "[[100000,100000,100000],3,1]", Expected: text("300000")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "hq4-merge-tuple", Type: "choice", Prompt: "K-Way Merge через heapq объединяет k отсортированных последовательностей. Что класть в Heap, чтобы числа выходили по возрастанию и было видно, откуда брать следующее?", Options: []string{"только число", "(число, номер последовательности, позиция)", "(позиция, число)"}, Answer: 1, Explain: "После heappop нужны источник и позиция следующего числа. Значение стоит первым, потому что tuple сравниваются с первого элемента. Ключ (позиция, число) сначала выдаст первые позиции, а не наименьшие значения."},
				{ID: "hq4-merge-next", Type: "short", Prompt: "K-Way Merge через Heap: [2,8,9], [1,4,14], [3,5,6]. Какое число попадёт в Heap сразу после удаления 4?", Answer: "14", Explain: "Число 4 находится в последовательности 1 на позиции 1. Следующее из неё - 14. Число 5 уже в Heap после удаления 3 из последовательности 2. Это один heappop и один heappush; после них живы 3 записи."},
				{ID: "hq4-merge-step", Type: "code", Prompt: "Шаг K-Way Merge: извлекли tuple и записали число в ответ. Что поставить вместо ___?", Language: "python", Code: `v, i, j = heapq.heappop(h)
out.append(v)
if j + 1 < len(lists[i]):
    heapq.heappush(h, ___)`, Options: []string{"(i, lists[i][j + 1], j + 1)", "(lists[i][j + 1], i, j)", "(lists[i][j + 1], i, j + 1)"}, Answer: 2, Explain: "Значение определяет приоритет, i указывает источник, j + 1 - новую позицию. Старое j снова добавляло бы одно и то же число. Шаг делает 1 удаление, 1 запись ответа и до 1 добавления: O(log(k + 1)) времени для k последовательностей. Весь Heap хранит O(k) записей, ответ отдельно растёт до N чисел."},
				{ID: "hq4-merge-size", Type: "choice", Prompt: "K-Way Merge объединяет 3 последовательности, всего 9 чисел. Сколько tuple лежит в Heap после полного шага, пока ни одна последовательность не закончилась?", Options: []string{"3", "9", "1"}, Answer: 0, Explain: "Одна запись от каждого источника. Между heappop и heappush записей временно 2, после шага 3. Для k источников размер O(k), шаг O(log(k + 1)), а собранный ответ отдельно содержит до N чисел."},
				{ID: "hq4-potion-drop", Type: "choice", Prompt: "Зелья [5,−4,1,−1,−2], здоровье в начале 0. После всех пяти здоровье −1. Какое одно зелье удалить из принятых, чтобы здоровье стало наибольшим?", Options: []string{"−2, последнее", "−1", "−4"}, Answer: 2, Explain: "Без −4 здоровье 3, без −2 только 1, без −1 только 0. Во всех случаях остаётся 4 зелья. Min-Heap находит вреднейшее в корне. Здесь 5 добавлений и 1 удаление, пик 5 записей; ответ отдельно одно число. Для n зелий время O(n log(n + 1)), рабочая память O(n)."},
				{ID: "hq4-cpu-idle", Type: "short", Prompt: "Задачи (приход, длительность): [[4,2],[5,3],[20,2],[21,2]], t = 0. Цикл добавляет в Heap пришедшие задачи. Если Heap пуст, выполняет t += 1. Сколько раз выполнится t += 1?", Answer: "15", Explain: "До t = 4 нужно 4 прибавления. Задачи 0 и 1 работают с 4 до 9. До t = 20 ещё 11 прибавлений: всего 15. Прямые переходы требуют 2 присваивания. Для n задач ожидание добавляет H шагов, здесь H = 15; с переходами общее время O(n log(n + 1)) вместе с Sorting, рабочая память O(n), ответ отдельно n номеров."},
				{ID: "hq4-score", Type: "short", Prompt: "nums1 = [2,3,1,5], nums2 = [3,2,9,3], k = 2. Счёт = сумма выбранных nums1 · минимум выбранных nums2. Чему равен наибольший счёт?", Answer: "21", Explain: "Пары после Sorting: (9,1), (3,5), (3,2), (2,3). Кандидаты при полном Heap: 18, 21, 16. Выбор номеров 0 и 3 даёт (2 + 5) · 3 = 21. Здесь 4 добавления, 2 удаления, 6 изменений суммы, 3 произведения; пик Heap 3 числа плюс 4 пары. Для n пар время O(n log(n + 1)), рабочая память O(n), ответ отдельно O(1)."},
				{ID: "hq4-cpu-fourth", Type: "short", Prompt: "Задачи (приход, длительность), номера с 0: [[2,5],[1,4],[3,1],[3,2],[10,1]]. Процессор берёт самую короткую из пришедших задач. Какая задача выполнится четвёртой?", Answer: "0", Explain: "Задача 1 работает с 1 до 5, задача 2 с 5 до 6, задача 3 с 6 до 8. В t = 8 задача 4 ещё недоступна, поэтому четвёртой идёт 0. Порядок [1,2,3,0,4]. Всего 5 добавлений и 5 удалений, пик Heap 3 записи, ответ 5 номеров. Sorting по длительности всех задач нарушил бы доступность."},
			},
			Links: []sets.Link{
				{Title: "Priority Queue", URL: "https://neerc.ifmo.ru/wiki/index.php?title=%D0%9F%D1%80%D0%B8%D0%BE%D1%80%D0%B8%D1%82%D0%B5%D1%82%D0%BD%D1%8B%D0%B5_%D0%BE%D1%87%D0%B5%D1%80%D0%B5%D0%B4%D0%B8", Source: "Викиконспекты ИТМО"},
				{Title: "Priority Queue", URL: "https://ru.wikipedia.org/wiki/%D0%9E%D1%87%D0%B5%D1%80%D0%B5%D0%B4%D1%8C_%D1%81_%D0%BF%D1%80%D0%B8%D0%BE%D1%80%D0%B8%D1%82%D0%B5%D1%82%D0%BE%D0%BC_(%D0%BF%D1%80%D0%BE%D0%B3%D1%80%D0%B0%D0%BC%D0%BC%D0%B8%D1%80%D0%BE%D0%B2%D0%B0%D0%BD%D0%B8%D0%B5)", Source: "Википедия"},
				{Title: "Heap (Priority Queue)", URL: "https://usaco.guide/silver/priority-queues", Source: "USACO Guide (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "hq4-merge",
					Title: "K-Way Merge: L0 = [1,5,9], L1 = [2,3,10], L2 = [4,6,7]",
					Frames: []sets.FigureFrame{
						{Caption: "Созданы 3 начальные записи, выполнен 1 heapify. Запись tuple хранит (число, источник, позиция). Счётчики heappop/heappush накапливают вызовы после Heapify. Память: число живых tuple в Heap и отдельно число элементов ответа.", Rows: []sets.FigureRow{
							{Label: "L0", Cells: []any{1, 5, 9}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "j", Cell: 0}}},
							{Label: "L1", Cells: []any{2, 3, 10}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "j", Cell: 0}}},
							{Label: "L2", Cells: []any{4, 6, 7}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "j", Cell: 0}}},
							{Label: "Heap", Cells: []any{"(1, 0, 0)", "(2, 1, 0)", "(4, 2, 0)"}, Marked: []int{0, 1, 2}, Tree: true},
							{Label: "ответ", Cells: []any{"", "", "", "", "", "", "", "", ""}},
							{Label: "heappop", Cells: []any{0}},
							{Label: "heappush", Cells: []any{0}},
							{Label: "живых Heap", Cells: []any{3}},
							{Label: "чисел ответа", Cells: []any{0}},
						}},
						{Caption: "Удаляем (1,0,0), пишем 1, добавляем (5,0,1). Память показана после полного шага.", Rows: []sets.FigureRow{
							{Label: "L0", Cells: []any{1, 5, 9}, Marked: []int{1}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "j", Cell: 1}}},
							{Label: "L1", Cells: []any{2, 3, 10}, Pointers: &sets.Pointers{{Name: "j", Cell: 0}}},
							{Label: "L2", Cells: []any{4, 6, 7}, Pointers: &sets.Pointers{{Name: "j", Cell: 0}}},
							{Label: "Heap", Cells: []any{"(2, 1, 0)", "(4, 2, 0)", "(5, 0, 1)"}, Marked: []int{2}, Tree: true},
							{Label: "ответ", Cells: []any{1, "", "", "", "", "", "", "", ""}, Marked: []int{0}},
							{Label: "heappop", Cells: []any{1}},
							{Label: "heappush", Cells: []any{1}},
							{Label: "живых Heap", Cells: []any{3}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "Удаляем (2,1,0), пишем 2, добавляем (3,1,1).", Rows: []sets.FigureRow{
							{Label: "L0", Cells: []any{1, 5, 9}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "j", Cell: 1}}},
							{Label: "L1", Cells: []any{2, 3, 10}, Marked: []int{1}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "j", Cell: 1}}},
							{Label: "L2", Cells: []any{4, 6, 7}, Pointers: &sets.Pointers{{Name: "j", Cell: 0}}},
							{Label: "Heap", Cells: []any{"(3, 1, 1)", "(5, 0, 1)", "(4, 2, 0)"}, Marked: []int{0}, Tree: true},
							{Label: "ответ", Cells: []any{1, 2, "", "", "", "", "", "", ""}, Marked: []int{1}},
							{Label: "heappop", Cells: []any{2}},
							{Label: "heappush", Cells: []any{2}},
							{Label: "живых Heap", Cells: []any{3}},
							{Label: "чисел ответа", Cells: []any{2}},
						}},
						{Caption: "Удаляем (3,1,1), пишем 3, добавляем (10,1,2).", Rows: []sets.FigureRow{
							{Label: "L0", Cells: []any{1, 5, 9}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "j", Cell: 1}}},
							{Label: "L1", Cells: []any{2, 3, 10}, Marked: []int{2}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}},
							{Label: "L2", Cells: []any{4, 6, 7}, Pointers: &sets.Pointers{{Name: "j", Cell: 0}}},
							{Label: "Heap", Cells: []any{"(4, 2, 0)", "(5, 0, 1)", "(10, 1, 2)"}, Marked: []int{2}, Tree: true},
							{Label: "ответ", Cells: []any{1, 2, 3, "", "", "", "", "", ""}, Marked: []int{2}},
							{Label: "heappop", Cells: []any{3}},
							{Label: "heappush", Cells: []any{3}},
							{Label: "живых Heap", Cells: []any{3}},
							{Label: "чисел ответа", Cells: []any{3}},
						}},
						{Caption: "Удаляем (4,2,0), пишем 4, добавляем (6,2,1).", Rows: []sets.FigureRow{
							{Label: "L0", Cells: []any{1, 5, 9}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "j", Cell: 1}}},
							{Label: "L1", Cells: []any{2, 3, 10}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}},
							{Label: "L2", Cells: []any{4, 6, 7}, Marked: []int{1}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "j", Cell: 1}}},
							{Label: "Heap", Cells: []any{"(5, 0, 1)", "(10, 1, 2)", "(6, 2, 1)"}, Marked: []int{2}, Tree: true},
							{Label: "ответ", Cells: []any{1, 2, 3, 4, "", "", "", "", ""}, Marked: []int{3}},
							{Label: "heappop", Cells: []any{4}},
							{Label: "heappush", Cells: []any{4}},
							{Label: "живых Heap", Cells: []any{3}},
							{Label: "чисел ответа", Cells: []any{4}},
						}},
						{Caption: "Удаляем (5,0,1), пишем 5, добавляем (9,0,2).", Rows: []sets.FigureRow{
							{Label: "L0", Cells: []any{1, 5, 9}, Marked: []int{2}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}},
							{Label: "L1", Cells: []any{2, 3, 10}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}},
							{Label: "L2", Cells: []any{4, 6, 7}, Dimmed: []int{0}, Pointers: &sets.Pointers{{Name: "j", Cell: 1}}},
							{Label: "Heap", Cells: []any{"(6, 2, 1)", "(10, 1, 2)", "(9, 0, 2)"}, Marked: []int{2}, Tree: true},
							{Label: "ответ", Cells: []any{1, 2, 3, 4, 5, "", "", "", ""}, Marked: []int{4}},
							{Label: "heappop", Cells: []any{5}},
							{Label: "heappush", Cells: []any{5}},
							{Label: "живых Heap", Cells: []any{3}},
							{Label: "чисел ответа", Cells: []any{5}},
						}},
						{Caption: "Удаляем (6,2,1), пишем 6, добавляем (7,2,2).", Rows: []sets.FigureRow{
							{Label: "L0", Cells: []any{1, 5, 9}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}},
							{Label: "L1", Cells: []any{2, 3, 10}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}},
							{Label: "L2", Cells: []any{4, 6, 7}, Marked: []int{2}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}},
							{Label: "Heap", Cells: []any{"(7, 2, 2)", "(10, 1, 2)", "(9, 0, 2)"}, Marked: []int{0}, Tree: true},
							{Label: "ответ", Cells: []any{1, 2, 3, 4, 5, 6, "", "", ""}, Marked: []int{5}},
							{Label: "heappop", Cells: []any{6}},
							{Label: "heappush", Cells: []any{6}},
							{Label: "живых Heap", Cells: []any{3}},
							{Label: "чисел ответа", Cells: []any{6}},
						}},
						{Caption: "Удаляем (7,2,2), пишем 7. Источник 2 закончился: нового heappush нет.", Rows: []sets.FigureRow{
							{Label: "L0", Cells: []any{1, 5, 9}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}},
							{Label: "L1", Cells: []any{2, 3, 10}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}},
							{Label: "L2", Cells: []any{4, 6, 7}, Dimmed: []int{0, 1, 2}},
							{Label: "Heap", Cells: []any{"(9, 0, 2)", "(10, 1, 2)"}, Tree: true},
							{Label: "ответ", Cells: []any{1, 2, 3, 4, 5, 6, 7, "", ""}, Marked: []int{6}},
							{Label: "heappop", Cells: []any{7}},
							{Label: "heappush", Cells: []any{6}},
							{Label: "живых Heap", Cells: []any{2}},
							{Label: "чисел ответа", Cells: []any{7}},
						}},
						{Caption: "Удаляем (9,0,2), пишем 9. Источник 0 закончился.", Rows: []sets.FigureRow{
							{Label: "L0", Cells: []any{1, 5, 9}, Dimmed: []int{0, 1, 2}},
							{Label: "L1", Cells: []any{2, 3, 10}, Dimmed: []int{0, 1}, Pointers: &sets.Pointers{{Name: "j", Cell: 2}}},
							{Label: "L2", Cells: []any{4, 6, 7}, Dimmed: []int{0, 1, 2}},
							{Label: "Heap", Cells: []any{"(10, 1, 2)"}, Tree: true},
							{Label: "ответ", Cells: []any{1, 2, 3, 4, 5, 6, 7, 9, ""}, Marked: []int{7}},
							{Label: "heappop", Cells: []any{8}},
							{Label: "heappush", Cells: []any{6}},
							{Label: "живых Heap", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{8}},
						}},
						{Caption: "Удаляем (10,1,2), пишем 10. Heap пуст. Всего 9 удалений, 6 добавлений после Heapify и 9 записей ответа; пик Heap - 3 tuple.", Rows: []sets.FigureRow{
							{Label: "L0", Cells: []any{1, 5, 9}, Dimmed: []int{0, 1, 2}},
							{Label: "L1", Cells: []any{2, 3, 10}, Dimmed: []int{0, 1, 2}},
							{Label: "L2", Cells: []any{4, 6, 7}, Dimmed: []int{0, 1, 2}},
							{Label: "Heap", Cells: []any{""}, Tree: true},
							{Label: "ответ", Cells: []any{1, 2, 3, 4, 5, 6, 7, 9, 10}, Marked: []int{8}},
							{Label: "heappop", Cells: []any{9}},
							{Label: "heappush", Cells: []any{6}},
							{Label: "живых Heap", Cells: []any{0}},
							{Label: "пик Heap", Cells: []any{3}},
							{Label: "чисел ответа", Cells: []any{9}},
						}},
					},
				},
				{
					ID:    "hq4-cpu",
					Title: "Single-Threaded CPU: выбор из доступных задач",
					Frames: []sets.FigureFrame{
						{Caption: "t = 1: пришла задача 0, она работает до 3. Рисунок Heap показывает состояние до выбора; накопленные вызовы и живые записи - после heappop. Ответ считаем в номерах задач. Порядок прихода отдельно хранит 5 номеров во всех кадрах.", Rows: []sets.FigureRow{
							{Label: "задачи", Cells: []any{0, 1, 2, 3, 4}, Marked: []int{0}},
							{Label: "приход", Cells: []any{1, 2, 3, 4, 20}},
							{Label: "Heap до выбора", Cells: []any{"(2, 0)"}, Marked: []int{0}, Tree: true},
							{Label: "порядок", Cells: []any{0, "", "", "", ""}, Marked: []int{0}},
							{Label: "heappush", Cells: []any{1}},
							{Label: "heappop", Cells: []any{1}},
							{Label: "живых после", Cells: []any{0}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "t = 3: доступны задачи 1 и 2. Задача 2 короче, работает до 5. Перед выбором достигнут пик Heap - 2 записи.", Rows: []sets.FigureRow{
							{Label: "задачи", Cells: []any{0, 1, 2, 3, 4}, Marked: []int{1, 2}, Dimmed: []int{0}},
							{Label: "приход", Cells: []any{1, 2, 3, 4, 20}},
							{Label: "Heap до выбора", Cells: []any{"(2, 2)", "(4, 1)"}, Marked: []int{0}, Tree: true},
							{Label: "порядок", Cells: []any{0, 2, "", "", ""}, Marked: []int{1}},
							{Label: "heappush", Cells: []any{3}},
							{Label: "heappop", Cells: []any{2}},
							{Label: "живых после", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{2}},
						}},
						{Caption: "t = 5: пришла задача 3. Она короче задачи 1, работает до 6.", Rows: []sets.FigureRow{
							{Label: "задачи", Cells: []any{0, 1, 2, 3, 4}, Marked: []int{1, 3}, Dimmed: []int{0, 2}},
							{Label: "приход", Cells: []any{1, 2, 3, 4, 20}},
							{Label: "Heap до выбора", Cells: []any{"(1, 3)", "(4, 1)"}, Marked: []int{0}, Tree: true},
							{Label: "порядок", Cells: []any{0, 2, 3, "", ""}, Marked: []int{2}},
							{Label: "heappush", Cells: []any{4}},
							{Label: "heappop", Cells: []any{3}},
							{Label: "живых после", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{3}},
						}},
						{Caption: "t = 6: осталась задача 1, она работает до 10.", Rows: []sets.FigureRow{
							{Label: "задачи", Cells: []any{0, 1, 2, 3, 4}, Marked: []int{1}, Dimmed: []int{0, 2, 3}},
							{Label: "приход", Cells: []any{1, 2, 3, 4, 20}},
							{Label: "Heap до выбора", Cells: []any{"(4, 1)"}, Marked: []int{0}, Tree: true},
							{Label: "порядок", Cells: []any{0, 2, 3, 1, ""}, Marked: []int{3}},
							{Label: "heappush", Cells: []any{4}},
							{Label: "heappop", Cells: []any{4}},
							{Label: "живых после", Cells: []any{0}},
							{Label: "чисел ответа", Cells: []any{4}},
						}},
						{Caption: "t = 10: Heap пуст, следующий приход в 20. Одно присваивание t = 20 пропускает 10 пустых тактов; вызовов Heap нет.", Rows: []sets.FigureRow{
							{Label: "задачи", Cells: []any{0, 1, 2, 3, 4}, Dimmed: []int{0, 2, 3, 1}},
							{Label: "приход", Cells: []any{1, 2, 3, 4, 20}},
							{Label: "Heap", Cells: []any{""}, Tree: true},
							{Label: "порядок", Cells: []any{0, 2, 3, 1, ""}},
							{Label: "heappush", Cells: []any{4}},
							{Label: "heappop", Cells: []any{4}},
							{Label: "живых после", Cells: []any{0}},
							{Label: "чисел ответа", Cells: []any{4}},
						}},
						{Caption: "t = 20: задача 4 работает до 23. Ответ [0,2,3,1,4]. Всего 5 добавлений и 5 удалений; время простоя не определяет число вызовов Heap.", Rows: []sets.FigureRow{
							{Label: "задачи", Cells: []any{0, 1, 2, 3, 4}, Marked: []int{4}, Dimmed: []int{0, 2, 3, 1}},
							{Label: "приход", Cells: []any{1, 2, 3, 4, 20}},
							{Label: "Heap до выбора", Cells: []any{"(3, 4)"}, Marked: []int{0}, Tree: true},
							{Label: "порядок", Cells: []any{0, 2, 3, 1, 4}, Marked: []int{4}},
							{Label: "heappush", Cells: []any{5}},
							{Label: "heappop", Cells: []any{5}},
							{Label: "живых после", Cells: []any{0}},
							{Label: "пик Heap", Cells: []any{2}},
							{Label: "чисел ответа", Cells: []any{5}},
						}},
					},
				},
				{
					ID:    "hq4-potions",
					Title: "Potions: [4,−3,−1,−2,2], Min-Heap хранит принятые зелья",
					Frames: []sets.FigureFrame{
						{Caption: "i = 0: принимаем 4, здоровье 4. Счётчики heappush/heappop накапливают вызовы. Память соответствует показанному состоянию: записи Heap, одно число здоровья и отдельно одно число ответа.", Rows: []sets.FigureRow{
							{Label: "зелья", Cells: []any{4, -3, -1, -2, 2}, Marked: []int{0}, Pointers: &sets.Pointers{{Name: "i", Cell: 0}}},
							{Label: "принято", Cells: []any{4}, Marked: []int{0}, Tree: true},
							{Label: "heappush", Cells: []any{1}},
							{Label: "heappop", Cells: []any{0}},
							{Label: "живых Heap", Cells: []any{1}},
							{Label: "чисел здоровья", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "i = 1: принимаем −3, здоровье 1. Минимальное принятое число −3 находится в корне.", Rows: []sets.FigureRow{
							{Label: "зелья", Cells: []any{4, -3, -1, -2, 2}, Marked: []int{1}, Pointers: &sets.Pointers{{Name: "i", Cell: 1}}},
							{Label: "принято", Cells: []any{-3, 4}, Marked: []int{0}, Tree: true},
							{Label: "heappush", Cells: []any{2}},
							{Label: "heappop", Cells: []any{0}},
							{Label: "живых Heap", Cells: []any{2}},
							{Label: "чисел здоровья", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "i = 2: принимаем −1, здоровье 0, удаление не нужно.", Rows: []sets.FigureRow{
							{Label: "зелья", Cells: []any{4, -3, -1, -2, 2}, Marked: []int{2}, Pointers: &sets.Pointers{{Name: "i", Cell: 2}}},
							{Label: "принято", Cells: []any{-3, 4, -1}, Marked: []int{2}, Tree: true},
							{Label: "heappush", Cells: []any{3}},
							{Label: "heappop", Cells: []any{0}},
							{Label: "живых Heap", Cells: []any{3}},
							{Label: "чисел здоровья", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "i = 3: принимаем −2, здоровье −2. Это промежуточное состояние до исправления; достигнут пик 4 записи.", Rows: []sets.FigureRow{
							{Label: "зелья", Cells: []any{4, -3, -1, -2, 2}, Marked: []int{3}, Pointers: &sets.Pointers{{Name: "i", Cell: 3}}},
							{Label: "принято", Cells: []any{-3, -2, -1, 4}, Marked: []int{1}, Tree: true},
							{Label: "heappush", Cells: []any{4}},
							{Label: "heappop", Cells: []any{0}},
							{Label: "живых Heap", Cells: []any{4}},
							{Label: "чисел здоровья", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "heappop удаляет −3. Здоровье −2 + 3 = 1, принято 3 зелья. Удаление текущего −2 дало бы только 0 здоровья.", Rows: []sets.FigureRow{
							{Label: "зелья", Cells: []any{4, -3, -1, -2, 2}, Dimmed: []int{1}, Pointers: &sets.Pointers{{Name: "i", Cell: 3}}},
							{Label: "принято", Cells: []any{-2, 4, -1}, Tree: true},
							{Label: "heappush", Cells: []any{4}},
							{Label: "heappop", Cells: []any{1}},
							{Label: "живых Heap", Cells: []any{3}},
							{Label: "чисел здоровья", Cells: []any{1}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "i = 4: принимаем 2, здоровье 3. Ответ 4. Всего 5 добавлений, 1 удаление и 6 изменений здоровья.", Rows: []sets.FigureRow{
							{Label: "зелья", Cells: []any{4, -3, -1, -2, 2}, Marked: []int{4}, Dimmed: []int{1}, Pointers: &sets.Pointers{{Name: "i", Cell: 4}}},
							{Label: "принято", Cells: []any{-2, 2, -1, 4}, Marked: []int{1}, Tree: true},
							{Label: "heappush", Cells: []any{5}},
							{Label: "heappop", Cells: []any{1}},
							{Label: "живых Heap", Cells: []any{4}},
							{Label: "пик Heap", Cells: []any{4}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
					},
				},
				{
					ID:    "hq4-two-heaps",
					Title: "Два Min-Heap: costs = [5,1,8,9,7,2,6], candidates = 2, k = 3",
					Frames: []sets.FigureFrame{
						{Caption: "Создаём группы L из номеров 0,1 и R из 5,6; вызываем heapify два раза. В рисунке показаны цены; ключи также хранят исходные номера. Счётчики накапливают вызовы при найме и пополнении, а также сравнения двух вершин. Живые записи считаем для L и R вместе; ответ хранит одно число.", Rows: []sets.FigureRow{
							{Label: "costs", Cells: []any{5, 1, 8, 9, 7, 2, 6}, Marked: []int{0, 1, 5, 6}, Pointers: &sets.Pointers{{Name: "l", Cell: 2}, {Name: "r", Cell: 4}}},
							{Label: "L", Cells: []any{1, 5}, Marked: []int{0, 1}, Tree: true},
							{Label: "R", Cells: []any{2, 6}, Marked: []int{0, 1}, Tree: true},
							{Label: "heappop", Cells: []any{0}},
							{Label: "heappush", Cells: []any{0}},
							{Label: "сравнений", Cells: []any{0}},
							{Label: "живых L+R", Cells: []any{4}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "Сравниваем цены 1 и 2. Нанимаем номер 1 за 1 из L. Добавляем в L номер 2 за 8, l = 3. Память дана после пополнения.", Rows: []sets.FigureRow{
							{Label: "costs", Cells: []any{5, 1, 8, 9, 7, 2, 6}, Marked: []int{2}, Dimmed: []int{1}, Pointers: &sets.Pointers{{Name: "l", Cell: 3}, {Name: "r", Cell: 4}}},
							{Label: "L", Cells: []any{5, 8}, Marked: []int{1}, Tree: true},
							{Label: "R", Cells: []any{2, 6}, Tree: true},
							{Label: "heappop", Cells: []any{1}},
							{Label: "heappush", Cells: []any{1}},
							{Label: "сравнений", Cells: []any{1}},
							{Label: "живых L+R", Cells: []any{4}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "Сравниваем 5 и 2. Нанимаем номер 5 за 2 из R. Добавляем в R номер 4 за 7, r = 3. Общая цена 3.", Rows: []sets.FigureRow{
							{Label: "costs", Cells: []any{5, 1, 8, 9, 7, 2, 6}, Marked: []int{4}, Dimmed: []int{1, 5}, Pointers: &sets.Pointers{{Name: "l", Cell: 3}, {Name: "r", Cell: 3}}},
							{Label: "L", Cells: []any{5, 8}, Tree: true},
							{Label: "R", Cells: []any{6, 7}, Marked: []int{1}, Tree: true},
							{Label: "heappop", Cells: []any{2}},
							{Label: "heappush", Cells: []any{2}},
							{Label: "сравнений", Cells: []any{2}},
							{Label: "живых L+R", Cells: []any{4}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
						{Caption: "Сравниваем 5 и 6. Нанимаем номер 0 за 5 из L, добавляем номер 3 за 9. Теперь l = 4 > r = 3. Ответ 1 + 2 + 5 = 8, пик двух Heap вместе 4 записи.", Rows: []sets.FigureRow{
							{Label: "costs", Cells: []any{5, 1, 8, 9, 7, 2, 6}, Marked: []int{3}, Dimmed: []int{1, 5, 0}, Pointers: &sets.Pointers{{Name: "l", Cell: 4}, {Name: "r", Cell: 3}}},
							{Label: "L", Cells: []any{8, 9}, Marked: []int{1}, Tree: true},
							{Label: "R", Cells: []any{6, 7}, Tree: true},
							{Label: "heappop", Cells: []any{3}},
							{Label: "heappush", Cells: []any{3}},
							{Label: "сравнений", Cells: []any{3}},
							{Label: "живых L+R", Cells: []any{4}},
							{Label: "чисел ответа", Cells: []any{1}},
						}},
					},
				},
			},
		},
	},
}
