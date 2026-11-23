// Package trees is the learning set "Tree и Dynamic Programming".
package trees

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
	Slug:    "trees",
	Title:   "Tree и Dynamic Programming",
	Summary: "Tree связывает n вершин n − 1 рёбрами без циклов. Вы строите список соседей, обходите Tree без Recursion, считаете поддеревья в порядке «дети раньше родителя» и решаете задачи Dynamic Programming на поддеревьях. В конце находите пути, диаметр и центр. Для каждого решения считаете время и дополнительную память.",
	Topics:  []string{"Tree", "Depth-First Search", "Breadth-First Search", "Dynamic Programming", "Graph Theory"},
	Chapters: []sets.Chapter{
		{
			Title:   "Tree: вершины, рёбра и обход",
			Summary: "Tree из n вершин имеет n − 1 рёбер и один путь между любыми двумя вершинами. Вы строите список соседей по рёбрам или по Array родителей, выбираете корень и находите родителя и глубину каждой вершины. BFS и DFS со Stack обходят Tree без Recursion, поэтому путь из 100 000 вершин не вызывает ошибку.",
			Lesson:  lesson1,
			Tasks: []sets.Task{
				{
					URL:        "https://leetcode.com/problems/find-center-of-star-graph/",
					Number:     number(1791),
					Title:      "Find Center of Star Graph",
					Difficulty: "easy",
					Topics:     []string{"Tree", "Graph Theory"},
					Note:       "Центр входит в каждое ребро, поэтому это общая вершина первых двух рёбер. n - число вершин. Время: O(1), не больше двух сравнений. Память: O(1) дополнительной, ответ одно число. Подсчёт степеней тоже работает: время O(n), память O(n) для n счётчиков. Для [[1, 2], [2, 3], [4, 2]] общая вершина первых двух рёбер равна 2. Проверка каждой вершины по всем рёбрам даёт O(n²).",
					Statement: `<p>Дан Graph в форме <strong>звезды</strong> из <code>n</code> вершин с номерами от <code>1</code> до <code>n</code>. Рёбра этого Graph не имеют направления. Graph в форме звезды - это Graph, в котором есть одна вершина-<strong>центр</strong> и <strong>ровно</strong> <code>n - 1</code> рёбер: они соединяют центр с каждой другой вершиной.</p><p>Дан двумерный Array целых чисел <code>edges</code>. Запись <code>edges[i] = [u<sub>i</sub>, v<sub>i</sub>]</code> означает, что между вершинами <code>u<sub>i</sub></code> и <code>v<sub>i</sub></code> есть ребро. Верните центр этой звезды.</p><p>Метод класса <code>Solution</code>: <code>def findCenter(self, edges: list[list[int]]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> edges = [[1,2],[2,3],[4,2]]
<strong>Выход:</strong> 2
<strong>Пояснение:</strong> вершина 2 соединена со всеми другими вершинами, поэтому 2 - центр.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> edges = [[1,2],[5,1],[1,3],[1,4]]
<strong>Выход:</strong> 1</pre><h3>Ограничения</h3><ul><li><code>3 ≤ n ≤ 10<sup>5</sup></code></li><li><code>edges.length == n - 1</code></li><li><code>edges[i].length == 2</code></li><li><code>1 ≤ u<sub>i</sub>, v<sub>i</sub> ≤ n</code></li><li><code>u<sub>i</sub> != v<sub>i</sub></code></li><li>Рёбра <code>edges</code> задают правильный Graph в форме звезды.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "findCenter"}, Cases: []sets.TestCase{
						{Input: "[[[1,2],[2,3],[4,2]]]", Expected: text("2")},
						{Input: "[[[1,2],[5,1],[1,3],[1,4]]]", Expected: text("1")},
						{Input: "[[[1,3],[3,2]]]", Expected: text("3")},
						{Input: "[[[4,7],[7,1],[7,2],[3,7],[7,5],[6,7]]]", Expected: text("7")},
						{Input: "[[[1,5],[5,2],[3,5],[5,4]]]", Expected: text("5")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1843/C",
					Title:      "Sum in Binary Tree",
					Difficulty: "easy",
					Topics:     []string{"Tree", "Math", "Bit Manipulation"},
					Note:       "Codeforces, рейтинг 800. Родитель вершины v равен v // 2, поэтому путь к корню: v, v // 2, ..., 1. Время: O(log n) на набор, не больше 54 шагов при n ≤ 10^16; всего O(t log n). Память: O(1) дополнительной, ответ одно число на набор. Для 10: 10 + 5 + 2 + 1 = 18. Делите целым делением //: n / 2 в float теряет точность при n > 2^53.",
					Statement: `<p>Ваня очень любит математику. Однажды он решал очередную задачу и придумал интересное Tree. Это Tree строится так.</p><p>Сначала в Tree есть только одна вершина с номером <code>1</code>: корень Tree. Затем Ваня добавляет к ней двух детей и даёт им номера подряд: <code>2</code> и <code>3</code>. После этого он добавляет детей к вершинам в порядке возрастания их номеров, начиная с <code>2</code>, и даёт детям наименьшие свободные номера. В результате у Вани получается бесконечное Binary Tree с корнем в вершине <code>1</code>: у каждой вершины ровно два ребёнка, а номера вершин идут подряд по уровням.</p><p>Ваня захотел узнать сумму номеров вершин на пути от вершины <code>1</code> до вершины с номером <code>n</code> в этом Tree. Ваня не любит считать, поэтому он попросил вас найти эту сумму.</p><h3>Входные данные</h3><p>В первой строке записано целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>): количество наборов входных данных.</p><p>Дальше идут <code>t</code> строк с описанием наборов. В каждой строке записано одно целое число <code>n</code> (<code>1 ≤ n ≤ 10<sup>16</sup></code>): номер вершины, для которой Ваня хочет посчитать сумму номеров вершин на пути от корня до этой вершины.</p><h3>Выходные данные</h3><p>Для каждого набора выведите одно целое число: искомую сумму.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>6
3
10
37
1
10000000000000000
15</pre><p><strong>Выходные данные:</strong></p><pre>4
18
71
1
19999999999999980
26</pre><h3>Примечание</h3><p>В первом наборе примера на пути от корня до вершины <code>3</code> лежат две вершины, <code>1</code> и <code>3</code>. Их сумма равна <code>4</code>.</p><p>Во втором наборе примера на пути от корня до вершины с номером <code>10</code> лежат вершины <code>1</code>, <code>2</code>, <code>5</code>, <code>10</code>. Сумма их номеров равна <code>1 + 2 + 5 + 10 = 18</code>.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>1 ≤ n ≤ 10<sup>16</sup></code></li><li>Ограничение по времени: 1 секунда на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "6\n3\n10\n37\n1\n10000000000000000\n15\n", Expected: text("4\n18\n71\n1\n19999999999999980\n26\n")},
						{Input: "1\n1\n", Expected: text("1\n")},
						{Input: "4\n2\n4\n8\n1024\n", Expected: text("3\n7\n15\n2047\n")},
						{Input: "3\n7\n6\n5\n", Expected: text("11\n10\n8\n")},
						{Input: "2\n9999999999999999\n9007199254740995\n", Expected: text("19999999999999963\n18014398509481987\n")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/time-needed-to-inform-all-employees/",
					Number:     number(1376),
					Title:      "Time Needed to Inform All Employees",
					Difficulty: "medium",
					Topics:     []string{"Tree", "Depth-First Search", "Breadth-First Search"},
					Note:       "Постройте списки подчинённых по manager и пройдите Tree от headID со Stack: подчинённый узнаёт новость в момент time начальника плюс informTime начальника. n - число сотрудников, h - наибольшая глубина. Время: O(n), каждый сотрудник один раз входит в Stack. Память: O(n) для списков, времён и Stack. Для цепочки из 5 сотрудников, где у 4 начальников informTime 1000, ответ 4000. Подъём от каждого сотрудника к headID даёт O(n · h), на цепочке O(n²).",
					Statement: `<p>В компании <code>n</code> сотрудников, у каждого свой номер от <code>0</code> до <code>n - 1</code>. Глава компании - сотрудник с номером <code>headID</code>.</p><p>У каждого сотрудника есть один прямой начальник. Начальники даны в Array <code>manager</code>: <code>manager[i]</code> - прямой начальник сотрудника <code>i</code>, <code>manager[headID] = -1</code>. Гарантируется, что отношения подчинения образуют Tree.</p><p>Глава компании хочет сообщить всем сотрудникам срочную новость. Он сообщит её своим прямым подчинённым, они сообщат своим подчинённым, и так дальше, пока все сотрудники не узнают новость.</p><p>Сотруднику <code>i</code> нужно <code>informTime[i]</code> минут, чтобы сообщить новость всем своим прямым подчинённым. То есть через <code>informTime[i]</code> минут все его прямые подчинённые могут начать передавать новость дальше.</p><p>Верните <em>число минут</em>, которое нужно, чтобы все сотрудники узнали срочную новость.</p><p>Метод класса <code>Solution</code>: <code>def numOfMinutes(self, n: int, headID: int, manager: list[int], informTime: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> n = 1, headID = 0, manager = [-1], informTime = [0]
<strong>Выход:</strong> 0
<strong>Пояснение:</strong> глава компании - единственный сотрудник компании.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> n = 6, headID = 2, manager = [2,2,-1,2,2,2], informTime = [0,0,1,0,0,0]
<strong>Выход:</strong> 1
<strong>Пояснение:</strong> глава компании с номером 2 - прямой начальник всех сотрудников компании. Ему нужна 1 минута, чтобы сообщить новость всем.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 10<sup>5</sup></code></li><li><code>0 ≤ headID &lt; n</code></li><li><code>manager.length == n</code></li><li><code>0 ≤ manager[i] &lt; n</code></li><li><code>manager[headID] == -1</code></li><li><code>informTime.length == n</code></li><li><code>0 ≤ informTime[i] ≤ 1000</code></li><li><code>informTime[i] == 0</code>, если у сотрудника <code>i</code> нет подчинённых.</li><li><strong>Гарантируется</strong>, что новость можно сообщить всем сотрудникам.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "numOfMinutes"}, Cases: []sets.TestCase{
						{Input: "[1,0,[-1],[0]]", Expected: text("0")},
						{Input: "[6,2,[2,2,-1,2,2,2],[0,0,1,0,0,0]]", Expected: text("1")},
						{Input: "[7,0,[-1,0,0,1,1,2,4],[1,2,7,0,3,0,0]]", Expected: text("8")},
						{Input: "[5,4,[1,2,3,4,-1],[0,1000,1000,1000,1000]]", Expected: text("4000")},
						{Input: "[4,3,[3,3,3,-1],[0,0,0,9]]", Expected: text("9")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/reorder-routes-to-make-all-paths-lead-to-the-city-zero/",
					Number:     number(1466),
					Title:      "Reorder Routes to Make All Paths Lead to the City Zero",
					Difficulty: "medium",
					Topics:     []string{"Tree", "Depth-First Search", "Breadth-First Search", "Graph Theory"},
					Note:       "Для дороги a → b запишите (b, 1) в adj[a] и (a, 0) в adj[b]. BFS от 0 складывает флаги рёбер к новым вершинам: флаг 1 значит, что дорога ведёт от 0. n - число городов. Время: O(n), 2(n − 1) записей соседей. Память: O(n) для списков, флагов и Deque. В первом примере флаг 1 у дорог 0 → 1, 1 → 3 и 4 → 5: ответ 3. Отдельный поиск пути до 0 для каждой дороги даёт O(n²).",
					Statement: `<p>Есть <code>n</code> городов с номерами от <code>0</code> до <code>n - 1</code> и <code>n - 1</code> дорог. Между двумя разными городами можно проехать только одним путём: эта сеть дорог образует Tree. В прошлом году министерство транспорта решило сделать на дорогах одностороннее движение, потому что дороги слишком узкие.</p><p>Дороги заданы списком <code>connections</code>: запись <code>connections[i] = [a<sub>i</sub>, b<sub>i</sub>]</code> означает дорогу из города <code>a<sub>i</sub></code> в город <code>b<sub>i</sub></code>.</p><p>В этом году в столице (город <code>0</code>) пройдёт большое событие, и многие люди хотят приехать в этот город.</p><p>Ваша задача: изменить направление некоторых дорог так, чтобы из каждого города можно было доехать до города <code>0</code>. Верните <strong>наименьшее</strong> число изменённых дорог.</p><p><strong>Гарантируется</strong>, что после изменения из каждого города можно доехать до города <code>0</code>.</p><p>Метод класса <code>Solution</code>: <code>def minReorder(self, n: int, connections: list[list[int]]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> n = 6, connections = [[0,1],[1,3],[2,3],[4,0],[4,5]]
<strong>Выход:</strong> 3</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> n = 5, connections = [[1,0],[1,2],[3,2],[3,4]]
<strong>Выход:</strong> 2</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> n = 3, connections = [[1,0],[2,0]]
<strong>Выход:</strong> 0</pre><h3>Ограничения</h3><ul><li><code>2 ≤ n ≤ 5 · 10<sup>4</sup></code></li><li><code>connections.length == n - 1</code></li><li><code>connections[i].length == 2</code></li><li><code>0 ≤ a<sub>i</sub>, b<sub>i</sub> ≤ n - 1</code></li><li><code>a<sub>i</sub> != b<sub>i</sub></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "minReorder"}, Cases: []sets.TestCase{
						{Input: "[6,[[0,1],[1,3],[2,3],[4,0],[4,5]]]", Expected: text("3")},
						{Input: "[5,[[1,0],[1,2],[3,2],[3,4]]]", Expected: text("2")},
						{Input: "[3,[[1,0],[2,0]]]", Expected: text("0")},
						{Input: "[2,[[0,1]]]", Expected: text("1")},
						{Input: "[5,[[0,1],[1,2],[2,3],[3,4]]]", Expected: text("4")},
						{Input: "[7,[[1,0],[1,2],[4,2],[3,4],[0,5],[6,5]]]", Expected: text("2")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/862/B",
					Title:      "Mahmoud and Ehab and the bipartiteness",
					Difficulty: "medium",
					Topics:     []string{"Tree", "Depth-First Search", "Graph Theory"},
					Note:       "Codeforces, рейтинг 1300. Покрасьте вершины по чётности глубины BFS: a вершин одного цвета, b другого. Новые рёбра можно проводить только между цветами: всего a · b пар, n − 1 из них уже рёбра. Время: O(n). Память: O(n) для списка соседей, цветов и Deque. Ответ до 2,5 · 10^9. Для пути из 5 вершин a = 3, b = 2, ответ 6 − 4 = 2. Проверка всех пар вершин даёт O(n²).",
					Statement: `<p>Махмуд и Эхаб продолжают свои приключения! Все в злой стране знают, что Доктор Зло любит двудольные Graph, особенно Tree.</p><p>Tree - это связный Graph без циклов. Двудольный Graph - это Graph, вершины которого можно разбить на 2 множества так, что для каждого ребра (u, v) этого Graph вершины u и v лежат в разных множествах. Более строгие определения Tree и двудольного Graph даны в примечании ниже.</p><p>Доктор Зло дал Махмуду и Эхабу Tree из n вершин и попросил добавить к нему рёбра так, чтобы Graph остался двудольным. Кроме того, после добавления рёбер Graph должен быть простым: без петель и кратных рёбер. Какое наибольшее число рёбер они могут добавить?</p><p>Петля - это ребро, которое соединяет вершину с ней самой. В Graph нет кратных рёбер, если между каждой парой вершин не больше одного ребра. Цикл и петля - не одно и то же.</p><h3>Входные данные</h3><p>В первой строке записано целое число <code>n</code>: число вершин Tree (<code>1 ≤ n ≤ 10<sup>5</sup></code>).</p><p>В следующих <code>n − 1</code> строках записаны целые числа <code>u</code> и <code>v</code> (<code>1 ≤ u, v ≤ n</code>, <code>u ≠ v</code>): рёбра Tree.</p><p>Гарантируется, что данный Graph - это Tree.</p><h3>Выходные данные</h3><p>Выведите одно целое число: наибольшее число рёбер, которое Махмуд и Эхаб могут добавить к Tree и выполнить условия.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>3
1 2
1 3</pre><p><strong>Выходные данные:</strong></p><pre>0</pre><h3>Пример 2</h3><p><strong>Входные данные:</strong></p><pre>5
1 2
2 3
3 4
4 5</pre><p><strong>Выходные данные:</strong></p><pre>2</pre><h3>Примечание</h3><p>Определение Tree: статья Википедии «Tree (graph theory)». Определение двудольного Graph: статья Википедии «Bipartite graph».</p><p>В первом примере без петель и кратных рёбер можно добавить только ребро (2, 3). Но с этим ребром Graph перестанет быть двудольным, поэтому ответ 0.</p><p>Во втором примере Махмуд и Эхаб могут добавить рёбра (1, 4) и (2, 5).</p><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 10<sup>5</sup></code></li><li><code>1 ≤ u, v ≤ n</code>, <code>u ≠ v</code></li><li>рёбра образуют Tree</li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "3\n1 2\n1 3\n", Expected: text("0\n")},
						{Input: "5\n1 2\n2 3\n3 4\n4 5\n", Expected: text("2\n")},
						{Input: "1\n", Expected: text("0\n")},
						{Input: "7\n2 1\n1 3\n4 2\n2 5\n6 3\n5 7\n", Expected: text("6\n")},
						{Input: "10\n1 2\n2 3\n3 4\n4 5\n5 6\n6 7\n7 8\n8 9\n9 10\n", Expected: text("16\n")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "tr1-parity", Type: "short", Prompt: "Codeforces 862B. Tree - путь из 6 вершин. Сколько рёбер можно добавить, чтобы Graph остался двудольным и без кратных рёбер?", Answer: "4", Explain: "Глубины от конца пути 0, 1, 2, 3, 4, 5: 3 вершины чётной глубины и 3 нечётной. Новое ребро может соединять только вершины разной чётности: 3 · 3 = 9 пар, 5 из них уже рёбра. Ответ 9 − 5 = 4. Ответ 9 забывает рёбра, которые уже есть. Время O(n), память O(n)."},
				{ID: "tr1-degree-sum", Type: "short", Prompt: "В Tree 50 вершин. Чему равна сумма степеней всех вершин?", Answer: "98", Explain: "Каждое из 49 рёбер добавляет 1 к степени обоих концов: 2 · 49 = 98. Ответ 49 считает каждое ребро только один раз."},
				{ID: "tr1-parent-check", Type: "code", Prompt: "BFS по Tree без Array посещённых вершин. Что стоит на месте ___?", Language: "python", Code: `parent = [-1] * n
q = deque([0])
while q:
    v = q.popleft()
    for u in adj[v]:
        if u != ___:
            parent[u] = v
            q.append(u)`, Options: []string{"parent[v]", "parent[u]", "v"}, Answer: 0, Explain: "У вершины v уже найден только её родитель, поэтому его и пропускаем. С parent[u] новая вершина 1 добавит назад вершину 0: parent[0] = -1, и BFS не закончится. Проверка u != v верна всегда, ошибка та же. Время O(n), память O(n)."},
				{ID: "tr1-depth", Type: "short", Prompt: "Корень 0, parent = [-1, 0, 1, 1, 2, 0]. Какая глубина у вершины 4?", Answer: "3", Explain: "Подъём 4 → 2 → 1 → 0 делает 3 шага. Ответ 2 - номер родителя, а не глубина. Ответ 1 считает только шаг к родителю. Подъём стоит O(depth) времени и O(1) памяти."},
				{ID: "tr1-recursion", Type: "choice", Prompt: "Tree - путь 0 − 1 − 2 − ... − 99 999. Рекурсивный dfs(0, -1) из урока запускают в Python с пределом по умолчанию. Что произойдёт?", Options: []string{"Обход закончится за O(n) без ошибок", "RecursionError: вложенных вызовов больше предела около 1000", "Обход сделает O(n²) шагов и закончится"}, Answer: 1, Explain: "Каждый вызов ждёт конца вызова для следующей вершины, глубина 100 000. Python останавливает такой код на глубине около 1000. Число шагов было бы O(n), но до конца код не дойдёт. BFS с Deque или DFS со Stack обходят путь за O(n) без предела глубины."},
				{ID: "tr1-stack-order", Type: "choice", Prompt: "adj = [[1, 2], [0, 3], [0], [1]]. DFS со Stack из урока начинает с вершины 0. В каком порядке он снимает вершины из Stack?", Options: []string{"0, 1, 3, 2", "0, 1, 2, 3", "0, 2, 1, 3"}, Answer: 2, Explain: "После 0 в Stack лежат 1 и 2, pop() берёт последнюю: 2. У 2 новых соседей нет, затем 1, потом её ребёнок 3. Порядок 0, 1, 2, 3 даёт BFS, а не Stack. 0, 1, 3, 2 даёт рекурсивный dfs: он идёт по соседям в порядке списка. Stack снимает последнюю добавленную вершину."},
				{ID: "tr1-stack-peak", Type: "short", Prompt: "Звезда: центр 0 и вершины 1, 2, ..., 9. DFS со Stack начинает с 0. Сколько вершин лежит в Stack в самый большой момент?", Answer: "9", Explain: "После pop() вершины 0 в Stack ложатся все 9 соседей. Дальше каждый pop() только уменьшает Stack. Поэтому пик Stack у звезды из n вершин равен n − 1, дополнительная память O(n). Ответ 2 верен для T, но не для звезды."},
				{ID: "tr1-reorder", Type: "short", Prompt: "LeetCode 1466. connections = [[1, 0], [1, 2], [3, 2]]. Сколько дорог нужно повернуть, чтобы из каждого города можно было доехать до 0?", Answer: "1", Explain: "BFS от 0: дорога 1 → 0 уже ведёт к 0. Дорога 1 → 2 ведёт от 0, её поворачиваем. Дорога 3 → 2 ведёт к 0 через 2 и 1. Ответ 2 поворачивает и 3 → 2, но она уже ведёт в нужную сторону. Время O(n), память O(n)."},
			},
			Links: []sets.Link{
				{Title: "Tree: эквивалентные определения", URL: "https://neerc.ifmo.ru/wiki/index.php?title=%D0%94%D0%B5%D1%80%D0%B5%D0%B2%D0%BE,_%D1%8D%D0%BA%D0%B2%D0%B8%D0%B2%D0%B0%D0%BB%D0%B5%D0%BD%D1%82%D0%BD%D1%8B%D0%B5_%D0%BE%D0%BF%D1%80%D0%B5%D0%B4%D0%B5%D0%BB%D0%B5%D0%BD%D0%B8%D1%8F", Source: "Викиконспекты ИТМО"},
				{Title: "Хранение Graph: список соседей", URL: "https://ru.algorithmica.org/cs/graph-traversals/storing-graphs/", Source: "Алгоритмика"},
				{Title: "BFS: обход в ширину", URL: "https://brestprog.by/topics/bfs/", Source: "brestprog"},
				{Title: "Recursion и DFS", URL: "https://brestprog.by/topics/dfs/", Source: "brestprog"},
				{Title: "sys.setrecursionlimit", URL: "https://docs.python.org/3/library/sys.html#sys.setrecursionlimit", Source: "Python docs (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "tr1-tree",
					Title: "Tree T: рёбра, степени и родители",
					Frames: []sets.FigureFrame{
						{Caption: "Tree T из 7 вершин и 6 рёбер, корень 0 сверху. Между вершинами 3 и 5 один путь: 3, 1, 0, 2, 5.", Rows: []sets.FigureRow{
							{Label: "Tree", Cells: []any{0, 1, 2, 3, 4, 5, "", "", "", 6}, Tree: true},
						}},
						{Caption: "Список рёбер: 6 пар (a, b), это n − 1 при n = 7. Порядок пар и порядок концов не важен.", Rows: []sets.FigureRow{
							{Label: "a", Cells: []any{0, 0, 1, 1, 2, 4}},
							{Label: "b", Cells: []any{1, 2, 3, 4, 5, 6}},
						}},
						{Caption: "Степень - число рёбер у вершины. Сумма 12 = 2 · 6. У вершин 3, 5 и 6 степень 1.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "степень", Cells: []any{2, 3, 2, 1, 2, 1, 1}, Marked: []int{3, 5, 6}},
						}},
						{Caption: "Корень 0: Array родителей из 6 чисел и −1 у корня. Подъём 6 → 4 → 1 → 0 делает 3 шага, это глубина вершины 6.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "родитель", Cells: []any{-1, 0, 0, 1, 1, 2, 4}, Marked: []int{6}},
							{Label: "глубина", Cells: []any{0, 1, 1, 2, 2, 2, 3}, Marked: []int{6}},
						}},
					},
				},
				{
					ID:    "tr1-bfs",
					Title: "BFS от корня 0",
					Frames: []sets.FigureFrame{
						{Caption: "popleft() вернул 0. Соседи 1 и 2 получили родителя 0 и глубину 1. Знак ? - вершина, которую BFS ещё не нашёл.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "родитель", Cells: []any{-1, 0, 0, "?", "?", "?", "?"}, Marked: []int{1, 2}},
							{Label: "глубина", Cells: []any{0, 1, 1, "?", "?", "?", "?"}, Marked: []int{1, 2}},
							{Label: "Deque", Cells: []any{1, 2}},
						}},
						{Caption: "Вершина 1: сосед 0 - её родитель, его пропускаем. Вершины 3 и 4 получили глубину 2. В Deque 3 вершины: это пик.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "родитель", Cells: []any{-1, 0, 0, 1, 1, "?", "?"}, Marked: []int{3, 4}},
							{Label: "глубина", Cells: []any{0, 1, 1, 2, 2, "?", "?"}, Marked: []int{3, 4}},
							{Label: "Deque", Cells: []any{2, 3, 4}},
						}},
						{Caption: "Вершина 2: новый сосед только 5, глубина 2.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "родитель", Cells: []any{-1, 0, 0, 1, 1, 2, "?"}, Marked: []int{5}},
							{Label: "глубина", Cells: []any{0, 1, 1, 2, 2, 2, "?"}, Marked: []int{5}},
							{Label: "Deque", Cells: []any{3, 4, 5}},
						}},
						{Caption: "У вершины 3 нет новых соседей.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "родитель", Cells: []any{-1, 0, 0, 1, 1, 2, "?"}},
							{Label: "глубина", Cells: []any{0, 1, 1, 2, 2, 2, "?"}},
							{Label: "Deque", Cells: []any{4, 5}},
						}},
						{Caption: "Вершина 4: найдена вершина 6 с глубиной 3.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "родитель", Cells: []any{-1, 0, 0, 1, 1, 2, 4}, Marked: []int{6}},
							{Label: "глубина", Cells: []any{0, 1, 1, 2, 2, 2, 3}, Marked: []int{6}},
							{Label: "Deque", Cells: []any{5, 6}},
						}},
						{Caption: "У вершины 5 нет новых соседей.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "родитель", Cells: []any{-1, 0, 0, 1, 1, 2, 4}},
							{Label: "глубина", Cells: []any{0, 1, 1, 2, 2, 2, 3}},
							{Label: "Deque", Cells: []any{6}},
						}},
						{Caption: "Deque пуст. 7 popleft и 12 чтений соседей: время O(n). Порядок 0, 1, 2, 3, 4, 5, 6.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "родитель", Cells: []any{-1, 0, 0, 1, 1, 2, 4}},
							{Label: "глубина", Cells: []any{0, 1, 1, 2, 2, 2, 3}},
							{Label: "Deque", Cells: []any{"пусто"}},
						}},
					},
				},
				{
					ID:    "tr1-stack",
					Title: "DFS со Stack от вершины 0",
					Frames: []sets.FigureFrame{
						{Caption: "pop() вернул 0. В Stack легли 1 и 2, сверху 2.", Rows: []sets.FigureRow{
							{Label: "Stack", Cells: []any{1, 2}, Marked: []int{0, 1}},
							{Label: "порядок", Cells: []any{0}, Marked: []int{0}},
						}},
						{Caption: "pop() вернул 2: последнюю добавленную вершину. Легла 5.", Rows: []sets.FigureRow{
							{Label: "Stack", Cells: []any{1, 5}, Marked: []int{1}},
							{Label: "порядок", Cells: []any{0, 2}, Marked: []int{1}},
						}},
						{Caption: "pop() вернул 5, новых соседей нет.", Rows: []sets.FigureRow{
							{Label: "Stack", Cells: []any{1}},
							{Label: "порядок", Cells: []any{0, 2, 5}, Marked: []int{2}},
						}},
						{Caption: "pop() вернул 1. Легли 3 и 4.", Rows: []sets.FigureRow{
							{Label: "Stack", Cells: []any{3, 4}, Marked: []int{0, 1}},
							{Label: "порядок", Cells: []any{0, 2, 5, 1}, Marked: []int{3}},
						}},
						{Caption: "pop() вернул 4. Легла 6.", Rows: []sets.FigureRow{
							{Label: "Stack", Cells: []any{3, 6}, Marked: []int{1}},
							{Label: "порядок", Cells: []any{0, 2, 5, 1, 4}, Marked: []int{4}},
						}},
						{Caption: "pop() вернул 6.", Rows: []sets.FigureRow{
							{Label: "Stack", Cells: []any{3}},
							{Label: "порядок", Cells: []any{0, 2, 5, 1, 4, 6}, Marked: []int{5}},
						}},
						{Caption: "pop() вернул 3, Stack пуст. Родитель всегда раньше детей. Пик Stack: 2 вершины.", Rows: []sets.FigureRow{
							{Label: "Stack", Cells: []any{"пусто"}},
							{Label: "порядок", Cells: []any{0, 2, 5, 1, 4, 6, 3}, Marked: []int{6}},
						}},
					},
				},
			},
		},
		{
			Title:   "Поддеревья: дети раньше родителя",
			Summary: "Размер, сумма и другие счётчики поддерева складываются из ответов детей. Вы проходите порядок BFS с конца, и каждая вершина передаёт готовый ответ родителю. Так за O(n) считаются размеры поддеревьев, части Tree без вершины или без ребра и перевозки по рёбрам.",
			Lesson:  lesson2,
			Tasks: []sets.Task{
				{
					URL:        "https://codeforces.com/problemset/problem/1676/G",
					Title:      "White-Black Balanced Subtrees",
					Difficulty: "medium",
					Topics:     []string{"Tree", "Depth-First Search", "Dynamic Programming"},
					Note:       "Codeforces, рейтинг 1300. a_i < i, поэтому номера уже дают порядок «родитель раньше детей». Пройдите i от n до 2 и прибавьте баланс вершины i (W = +1, B = −1) к балансу a_i; поддерево сбалансировано, если его баланс 0. Время: O(n) на набор, O(сумма n) всего. Память: O(n) для родителей и балансов. Для 1 1 2 3 3 5 и WBBWWBW баланс 0 у вершин 2 и 3, ответ 2. Подъём от каждой вершины к корню даёт O(n²) на пути.",
					Statement: `<p>Дано Tree с корнем из <code>n</code> вершин с номерами от <code>1</code> до <code>n</code>. Корень Tree: вершина <code>1</code>. Также дана строка <code>s</code> с цветом каждой вершины: если <code>s<sub>i</sub> = B</code>, то вершина <code>i</code> чёрная, а если <code>s<sub>i</sub> = W</code>, то вершина <code>i</code> белая.</p><p>Поддерево называется <strong>сбалансированным</strong>, если в нём число белых вершин равно числу чёрных вершин. Посчитайте число сбалансированных поддеревьев.</p><p>Tree - это связный Graph без циклов, рёбра которого не имеют направления. Tree с корнем - это Tree, в котором выбрана одна вершина: корень. В этой задаче у всех Tree корень <code>1</code>.</p><p>Tree задано Array родителей <code>a<sub>2</sub>, ..., a<sub>n</sub></code> из <code>n − 1</code> чисел: <code>a<sub>i</sub></code> - родитель вершины с номером <code>i</code> для всех <code>i = 2, ..., n</code>. Родитель вершины <code>u</code> - это следующая вершина на простом пути от <code>u</code> к корню.</p><p>Поддерево вершины <code>u</code> - это множество всех вершин, у которых простой путь к корню проходит через <code>u</code>. Например, возьмём Tree с <code>n = 7</code>, <code>a = [1, 1, 2, 3, 3, 5]</code> и <code>s = WBBWWBW</code>. Вершина <code>7</code> лежит в поддереве вершины <code>3</code>, потому что простой путь <code>7 → 5 → 3 → 1</code> проходит через <code>3</code>. Вершина входит в своё поддерево, а поддерево корня - это всё Tree. В этом Tree поддерево вершины <code>3</code> сбалансировано.</p><h3>Входные данные</h3><p>В первой строке записано целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>): количество наборов входных данных.</p><p>В первой строке каждого набора записано целое число <code>n</code> (<code>2 ≤ n ≤ 4000</code>): число вершин Tree.</p><p>Во второй строке каждого набора записаны <code>n − 1</code> целых чисел <code>a<sub>2</sub>, ..., a<sub>n</sub></code> (<code>1 ≤ a<sub>i</sub> &lt; i</code>): родители вершин <code>2, ..., n</code>.</p><p>В третьей строке каждого набора записана строка <code>s</code> длины <code>n</code> из символов <code>B</code> и <code>W</code>: цвета вершин Tree.</p><p>Гарантируется, что сумма значений <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого набора выведите одно целое число: число сбалансированных поддеревьев.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>3
7
1 1 2 3 3 5
WBBWWBW
2
1
BW
8
1 2 3 4 5 6 7
BWBWBWBW</pre><p><strong>Выходные данные:</strong></p><pre>2
1
4</pre><h3>Примечание</h3><p>Первый набор описан в условии. Сбалансированы только поддеревья вершин <code>2</code> и <code>3</code>.</p><p>Во втором наборе сбалансировано только поддерево вершины <code>1</code>.</p><p>В третьем наборе сбалансированы только поддеревья вершин <code>1</code>, <code>3</code>, <code>5</code> и <code>7</code>.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>2 ≤ n ≤ 4000</code></li><li><code>1 ≤ a<sub>i</sub> &lt; i</code></li><li>строка <code>s</code> длины <code>n</code> из символов <code>B</code> и <code>W</code></li><li>сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code></li><li>Ограничение по времени: 1 секунда на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "3\n7\n1 1 2 3 3 5\nWBBWWBW\n2\n1\nBW\n8\n1 2 3 4 5 6 7\nBWBWBWBW\n", Expected: text("2\n1\n4\n")},
						{Input: "1\n2\n1\nBB\n", Expected: text("0\n")},
						{Input: "1\n4\n1 1 1\nBWWB\n", Expected: text("1\n")},
						{Input: "1\n7\n1 1 2 2 3 5\nWBBWBWW\n", Expected: text("3\n")},
						{Input: "1\n3\n1 2\nWBW\n", Expected: text("1\n")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/count-nodes-with-the-highest-score/",
					Number:     number(2049),
					Title:      "Count Nodes With the Highest Score",
					Difficulty: "medium",
					Topics:     []string{"Array", "Tree", "Depth-First Search"},
					Note:       "Сначала посчитайте размеры поддеревьев с конца порядка «родитель раньше детей». Без вершины v остаются части size[c] для детей c и n − size[v] сверху, если эта часть не пуста. n - число вершин. Время: O(n). Память: O(n) для детей, порядка и размеров. Оценка не больше (n / 3)^3, меньше 4 · 10^13, в Python это обычный int. Для [-1, 2, 0, 2, 0] оценки 3, 4, 2, 4, 4, ответ 3. Удаление каждой вершины с новым обходом даёт O(n²).",
					Statement: `<p>Дано <strong>Binary Tree</strong> из <code>n</code> вершин с корнем <code>0</code>. Вершины имеют номера от <code>0</code> до <code>n - 1</code>. Tree задано Array целых чисел <code>parents</code>, номера в нём <strong>начинаются с 0</strong>: <code>parents[i]</code> - родитель вершины <code>i</code>. Вершина <code>0</code> - корень, поэтому <code>parents[0] == -1</code>.</p><p>У каждой вершины есть <strong>оценка</strong>. Чтобы найти оценку вершины, представьте, что эту вершину и все её рёбра <strong>удалили</strong>. Tree распадётся на одно или несколько <strong>непустых</strong> поддеревьев. <strong>Размер</strong> поддерева - это число вершин в нём. <strong>Оценка</strong> вершины - это <strong>произведение размеров</strong> всех этих поддеревьев.</p><p>Верните <em><strong>число</strong> вершин с <strong>наибольшей оценкой</strong></em>.</p><p>Метод класса <code>Solution</code>: <code>def countHighestScoreNodes(self, parents: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> parents = [-1,2,0,2,0]
<strong>Выход:</strong> 3
<strong>Пояснение:</strong>
- Оценка вершины 0: 3 · 1 = 3
- Оценка вершины 1: 4 = 4
- Оценка вершины 2: 1 · 1 · 2 = 2
- Оценка вершины 3: 4 = 4
- Оценка вершины 4: 4 = 4
Наибольшая оценка 4, её имеют три вершины: 1, 3 и 4.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> parents = [-1,2,0]
<strong>Выход:</strong> 2
<strong>Пояснение:</strong>
- Оценка вершины 0: 2 = 2
- Оценка вершины 1: 2 = 2
- Оценка вершины 2: 1 · 1 = 1
Наибольшая оценка 2, её имеют две вершины: 0 и 1.</pre><h3>Ограничения</h3><ul><li><code>n == parents.length</code></li><li><code>2 ≤ n ≤ 10<sup>5</sup></code></li><li><code>parents[0] == -1</code></li><li><code>0 ≤ parents[i] ≤ n - 1</code> для <code>i != 0</code></li><li><code>parents</code> задаёт правильное Binary Tree.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "countHighestScoreNodes"}, Cases: []sets.TestCase{
						{Input: "[[-1,2,0,2,0]]", Expected: text("3")},
						{Input: "[[-1,2,0]]", Expected: text("2")},
						{Input: "[[-1,0]]", Expected: text("2")},
						{Input: "[[-1,0,0,1,1,2,4]]", Expected: text("1")},
						{Input: "[[-1,0,1,2,3]]", Expected: text("3")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1843/D",
					Title:      "Apple Tree",
					Difficulty: "medium",
					Topics:     []string{"Tree", "Depth-First Search", "Dynamic Programming", "Combinatorics"},
					Note:       "Codeforces, рейтинг 1200. Корень 1 не лист при n ≥ 2: у него есть ребёнок. leaves[v] = 1 у вершины без детей, иначе сумма leaves по детям; ответ запроса leaves[x] · leaves[y]. Время: O(n + q) на набор. Память: O(n) для соседей, порядка и счётчиков. Ответ меньше (2 · 10^5)^2 = 4 · 10^10. В первом наборе первого примера leaves[3] = 2 (листья 4 и 5) и leaves[4] = 1, поэтому ответ на запрос 3 4 равен 2. Новый обход для каждого запроса даёт O(nq).",
					Statement: `<p>У Тимофея в саду растёт яблоня. Это Tree с корнем из <code>n</code> вершин, корень в вершине <code>1</code> (вершины имеют номера от <code>1</code> до <code>n</code>). Tree - это связный Graph без петель и кратных рёбер.</p><p>Это Tree очень необычное: оно растёт корнем вверх. Но для Tree программистов это вполне нормально.</p><p>Яблоня ещё молодая, поэтому на ней вырастут только два яблока. Яблоки вырастут в некоторых вершинах (эти вершины могут совпадать). Когда яблоки вырастут, Тимофей начнёт трясти яблоню, пока яблоки не упадут. Каждый раз, когда Тимофей трясёт яблоню, с каждым яблоком происходит следующее.</p><p>Пусть яблоко сейчас в вершине <code>u</code>.</p><ul><li>Если у вершины <code>u</code> есть ребёнок, яблоко переходит в него (если таких вершин несколько, яблоко может перейти в любую из них).</li><li>Иначе яблоко падает с Tree.</li></ul><p>Можно показать, что за конечное время оба яблока упадут с Tree.</p><p>У Тимофея есть <code>q</code> предположений о том, в каких вершинах могут вырасти яблоки. В каждом предположении он считает, что яблоки вырастут в вершинах <code>x</code> и <code>y</code>, и хочет узнать число пар вершин (<code>a</code>, <code>b</code>), из которых яблоки могут упасть с Tree. Здесь <code>a</code> - вершина, из которой упадёт яблоко из вершины <code>x</code>, а <code>b</code> - вершина, из которой упадёт яблоко из вершины <code>y</code>. Помогите ему.</p><h3>Входные данные</h3><p>В первой строке записано целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>): количество наборов входных данных.</p><p>В первой строке каждого набора записано целое число <code>n</code> (<code>2 ≤ n ≤ 2 · 10<sup>5</sup></code>): число вершин Tree.</p><p>Дальше идут <code>n − 1</code> строк с описанием Tree. В строке <code>i</code> записаны два целых числа <code>u<sub>i</sub></code> и <code>v<sub>i</sub></code> (<code>1 ≤ u<sub>i</sub>, v<sub>i</sub> ≤ n</code>, <code>u<sub>i</sub> ≠ v<sub>i</sub></code>): ребро Tree.</p><p>В следующей строке записано одно целое число <code>q</code> (<code>1 ≤ q ≤ 2 · 10<sup>5</sup></code>): число предположений Тимофея.</p><p>В каждой из следующих <code>q</code> строк записаны два целых числа <code>x<sub>i</sub></code> и <code>y<sub>i</sub></code> (<code>1 ≤ x<sub>i</sub>, y<sub>i</sub> ≤ n</code>): вершины, в которых вырастут яблоки по предположению <code>i</code>.</p><p>Гарантируется, что сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>. Также гарантируется, что сумма <code>q</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого предположения Тимофея выведите в отдельной строке число упорядоченных пар вершин, из которых яблоки могут упасть с Tree, если предположение верно.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>2
5
1 2
3 4
5 3
3 2
4
3 4
5 1
4 4
1 3
3
1 2
1 3
3
1 1
2 3
3 1</pre><p><strong>Выходные данные:</strong></p><pre>2
2
1
4
4
1
2</pre><h3>Пример 2</h3><p><strong>Входные данные:</strong></p><pre>2
5
5 1
1 2
2 3
4 3
2
5 5
5 1
5
3 2
5 3
2 1
4 2
3
4 3
2 1
4 2</pre><p><strong>Выходные данные:</strong></p><pre>1
2
1
4
2</pre><h3>Примечание</h3><p>В первом примере:</p><ul><li>Для первого предположения яблоки могут упасть из двух пар вершин: (4, 4), (5, 4).</li><li>Для второго предположения тоже две пары: (5, 4), (5, 5).</li><li>Для третьего предположения только одна пара: (4, 4).</li><li>Для четвёртого предположения 4 пары: (4, 4), (4, 5), (5, 4), (5, 5).</li></ul><p>Во втором наборе первого примера для первого предположения яблоки могут упасть из 4 пар вершин: (2, 3), (2, 2), (3, 2), (3, 3). Для второго предположения возможна только одна пара: (2, 3). Для третьего предположения две пары: (3, 2), (3, 3).</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>2 ≤ n ≤ 2 · 10<sup>5</sup></code></li><li><code>1 ≤ u<sub>i</sub>, v<sub>i</sub> ≤ n</code>, <code>u<sub>i</sub> ≠ v<sub>i</sub></code></li><li><code>1 ≤ q ≤ 2 · 10<sup>5</sup></code></li><li><code>1 ≤ x<sub>i</sub>, y<sub>i</sub> ≤ n</code></li><li>сумма <code>n</code> и сумма <code>q</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code> каждая</li><li>Ограничение по времени: 4 секунды на тест.</li><li>Ограничение по памяти: 512 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "2\n5\n1 2\n3 4\n5 3\n3 2\n4\n3 4\n5 1\n4 4\n1 3\n3\n1 2\n1 3\n3\n1 1\n2 3\n3 1\n", Expected: text("2\n2\n1\n4\n4\n1\n2\n")},
						{Input: "2\n5\n5 1\n1 2\n2 3\n4 3\n2\n5 5\n5 1\n5\n3 2\n5 3\n2 1\n4 2\n3\n4 3\n2 1\n4 2\n", Expected: text("1\n2\n1\n4\n2\n")},
						{Input: "1\n2\n1 2\n3\n1 1\n1 2\n2 2\n", Expected: text("1\n1\n1\n")},
						{Input: "2\n3\n1 2\n2 3\n2\n1 1\n2 1\n5\n1 2\n1 3\n4 1\n1 5\n3\n1 1\n1 3\n3 3\n", Expected: text("1\n1\n16\n4\n1\n")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/minimum-fuel-cost-to-report-to-the-capital/",
					Number:     number(2477),
					Title:      "Minimum Fuel Cost to Report to the Capital",
					Difficulty: "medium",
					Topics:     []string{"Tree", "Depth-First Search", "Breadth-First Search", "Graph Theory"},
					Note:       "Через дорогу над городом v едут size[v] представителей, им нужно size[v] / seats машин с округлением вверх: (size[v] + seats − 1) // seats. Сложите по всем v ≠ 0. n - число городов. Время: O(n). Память: O(n) для соседей, порядка и размеров. Во втором примере над городами 1, 3, 2, 4, 5, 6 едут 3, 2, 1, 2, 1, 1 человек: 2 + 1 + 1 + 1 + 1 + 1 = 7. При seats = 1 ответ равен сумме глубин, на пути n(n − 1)/2. Подъём каждого представителя до столицы даёт O(n²) на пути.",
					Statement: `<p>Сеть дорог страны имеет форму Tree, то есть связного Graph без циклов, рёбра которого не имеют направления. В стране <code>n</code> городов с номерами от <code>0</code> до <code>n - 1</code> и ровно <code>n - 1</code> дорог. Столица - город <code>0</code>. Дан двумерный Array целых чисел <code>roads</code>: запись <code>roads[i] = [a<sub>i</sub>, b<sub>i</sub>]</code> означает, что есть <strong>двусторонняя дорога</strong> между городами <code>a<sub>i</sub></code> и <code>b<sub>i</sub></code>.</p><p>Представители всех городов едут на встречу. Встреча проходит в столице.</p><p>В каждом городе есть одна машина. Дано целое число <code>seats</code>: число мест в каждой машине.</p><p>Представитель может ехать на машине своего города или пересесть в другую машину и ехать вместе с другим представителем. Поездка между двумя городами по дороге стоит один литр топлива.</p><p>Верните <em>наименьшее число литров топлива, чтобы все представители доехали до столицы</em>.</p><p>Метод класса <code>Solution</code>: <code>def minimumFuelCost(self, roads: list[list[int]], seats: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> roads = [[0,1],[0,2],[0,3]], seats = 5
<strong>Выход:</strong> 3
<strong>Пояснение:</strong>
- Представитель<sub>1</sub> едет прямо в столицу и тратит 1 литр топлива.
- Представитель<sub>2</sub> едет прямо в столицу и тратит 1 литр топлива.
- Представитель<sub>3</sub> едет прямо в столицу и тратит 1 литр топлива.
Всего нужно 3 литра топлива.
Можно доказать, что 3 - наименьшее нужное число литров.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> roads = [[3,1],[3,2],[1,0],[0,4],[0,5],[4,6]], seats = 2
<strong>Выход:</strong> 7
<strong>Пояснение:</strong>
- Представитель<sub>2</sub> едет прямо в город 3 и тратит 1 литр топлива.
- Представитель<sub>2</sub> и представитель<sub>3</sub> едут вместе в город 1 и тратят 1 литр топлива.
- Представитель<sub>2</sub> и представитель<sub>3</sub> едут вместе в столицу и тратят 1 литр топлива.
- Представитель<sub>1</sub> едет прямо в столицу и тратит 1 литр топлива.
- Представитель<sub>5</sub> едет прямо в столицу и тратит 1 литр топлива.
- Представитель<sub>6</sub> едет прямо в город 4 и тратит 1 литр топлива.
- Представитель<sub>4</sub> и представитель<sub>6</sub> едут вместе в столицу и тратят 1 литр топлива.
Всего нужно 7 литров топлива.
Можно доказать, что 7 - наименьшее нужное число литров.</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> roads = [], seats = 1
<strong>Выход:</strong> 0
<strong>Пояснение:</strong> никому из представителей не нужно ехать в столицу.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 10<sup>5</sup></code></li><li><code>roads.length == n - 1</code></li><li><code>roads[i].length == 2</code></li><li><code>0 ≤ a<sub>i</sub>, b<sub>i</sub> &lt; n</code></li><li><code>a<sub>i</sub> != b<sub>i</sub></code></li><li><code>roads</code> задаёт правильное Tree.</li><li><code>1 ≤ seats ≤ 10<sup>5</sup></code></li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "minimumFuelCost"}, Cases: []sets.TestCase{
						{Input: "[[[0,1],[0,2],[0,3]],5]", Expected: text("3")},
						{Input: "[[[3,1],[3,2],[1,0],[0,4],[0,5],[4,6]],2]", Expected: text("7")},
						{Input: "[[],1]", Expected: text("0")},
						{Input: "[[[0,1],[1,2],[2,3]],1]", Expected: text("6")},
						{Input: "[[[0,1],[0,2],[1,3],[1,4],[2,5],[4,6]],2]", Expected: text("7")},
						{Input: "[[[0,1],[0,2],[1,3],[1,4],[2,5],[4,6]],100000]", Expected: text("6")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/maximum-number-of-k-divisible-components/",
					Number:     number(2872),
					Title:      "Maximum Number of K-Divisible Components",
					Difficulty: "hard",
					Topics:     []string{"Tree", "Depth-First Search"},
					Note:       "Суммы поддеревьев считайте с конца порядка BFS. Если сумма поддерева v делится на k, ребро над v можно разрезать: остальная часть тоже делится на k, потому что делится общая сумма. Ответ равен числу таких вершин, корень тоже считаем. Время: O(n). Память: O(n) для соседей, порядка и сумм. Суммы до 3 · 10^13. Во втором примере суммы поддеревьев вершин 1, 2 и 0 равны 6, 9 и 18: ответ 3. Новый подсчёт суммы части для каждого ребра даёт O(n²).",
					Statement: `<p>Дано Tree из <code>n</code> вершин с номерами от <code>0</code> до <code>n - 1</code>, рёбра которого не имеют направления. Дано целое число <code>n</code> и двумерный Array целых чисел <code>edges</code> длины <code>n - 1</code>: запись <code>edges[i] = [a<sub>i</sub>, b<sub>i</sub>]</code> означает, что в Tree есть ребро между вершинами <code>a<sub>i</sub></code> и <code>b<sub>i</sub></code>.</p><p>Также дан Array целых чисел <code>values</code> длины <code>n</code>, номера в нём <strong>начинаются с 0</strong>: <code>values[i]</code> - <strong>значение</strong> вершины <code>i</code>. И дано целое число <code>k</code>.</p><p>Удалим из Tree какие-то рёбра. Tree распадётся на части: в одной части лежат вершины, между которыми остался путь по рёбрам. <strong>Значение части</strong> - это сумма значений её вершин. <strong>Правильное разбиение</strong> Tree получается, если удалить любое множество рёбер (может быть, пустое) так, что значение каждой части делится на <code>k</code>.</p><p>Верните <em><strong>наибольшее число частей</strong> в правильном разбиении</em>.</p><p>Метод класса <code>Solution</code>: <code>def maxKDivisibleComponents(self, n: int, edges: list[list[int]], values: list[int], k: int) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> n = 5, edges = [[0,2],[1,2],[1,3],[2,4]], values = [1,8,1,4,4], k = 6
<strong>Выход:</strong> 2
<strong>Пояснение:</strong> удалим ребро между вершинами 1 и 2. Разбиение правильное, потому что:
- значение части с вершинами 1 и 3 равно values[1] + values[3] = 12;
- значение части с вершинами 0, 2 и 4 равно values[0] + values[2] + values[4] = 6.
Можно показать, что ни в одном другом правильном разбиении нет больше 2 частей.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> n = 7, edges = [[0,1],[0,2],[1,3],[1,4],[2,5],[2,6]], values = [3,0,6,1,5,2,1], k = 3
<strong>Выход:</strong> 3
<strong>Пояснение:</strong> удалим ребро между вершинами 0 и 2 и ребро между вершинами 0 и 1. Разбиение правильное, потому что:
- значение части с вершиной 0 равно values[0] = 3;
- значение части с вершинами 2, 5 и 6 равно values[2] + values[5] + values[6] = 9;
- значение части с вершинами 1, 3 и 4 равно values[1] + values[3] + values[4] = 6.
Можно показать, что ни в одном другом правильном разбиении нет больше 3 частей.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 3 · 10<sup>4</sup></code></li><li><code>edges.length == n - 1</code></li><li><code>edges[i].length == 2</code></li><li><code>0 ≤ a<sub>i</sub>, b<sub>i</sub> &lt; n</code></li><li><code>values.length == n</code></li><li><code>0 ≤ values[i] ≤ 10<sup>9</sup></code></li><li><code>1 ≤ k ≤ 10<sup>9</sup></code></li><li>Сумма <code>values</code> делится на <code>k</code>.</li><li>Входные данные таковы, что <code>edges</code> задаёт правильное Tree.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "maxKDivisibleComponents"}, Cases: []sets.TestCase{
						{Input: "[5,[[0,2],[1,2],[1,3],[2,4]],[1,8,1,4,4],6]", Expected: text("2")},
						{Input: "[7,[[0,1],[0,2],[1,3],[1,4],[2,5],[2,6]],[3,0,6,1,5,2,1],3]", Expected: text("3")},
						{Input: "[1,[],[0],1]", Expected: text("1")},
						{Input: "[4,[[0,1],[1,2],[1,3]],[0,0,0,0],7]", Expected: text("4")},
						{Input: "[3,[[0,1],[1,2]],[5,7,9],1]", Expected: text("3")},
						{Input: "[7,[[0,1],[0,2],[1,3],[1,4],[2,5],[4,6]],[2,1,4,3,5,2,1],3]", Expected: text("4")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "tr2-rest", Type: "short", Prompt: "В Tree 10 вершин, корень 0. У вершины v ≠ 0 size[v] = 3. Удалим ребро между v и её родителем. Сколько вершин в части с корнем?", Answer: "7", Explain: "Поддерево v уходит целиком: 3 вершины. В части с корнем остаётся n − size[v] = 10 − 3 = 7. Ответ 3 - это часть с вершиной v. Подсчёт стоит O(1) после размеров."},
				{ID: "tr2-no-reversed", Type: "choice", Prompt: "Ученик забыл reversed: for v in order: if parent[v] != -1: size[parent[v]] += size[v]. Путь 0 − 1 − 2, parent = [-1, 0, 1], order = [0, 1, 2], сначала size = [1, 1, 1]. Что получится в size?", Options: []string{"[3, 2, 1]", "[1, 1, 1]", "[2, 2, 1]"}, Answer: 2, Explain: "При v = 1 к size[0] прибавляется size[1] = 1, хотя вершина 2 ещё не учтена: size[0] = 2. Потом v = 2 даёт size[1] = 2. Верные размеры [3, 2, 1] даёт только проход с конца, где дети раньше родителя. [1, 1, 1] был бы без прибавлений."},
				{ID: "tr2-visits", Type: "short", Prompt: "Brute Force считает размер каждого поддерева отдельным обходом. Tree - путь из 100 вершин, корень на конце. Сколько посещений вершин он сделает?", Answer: "5050", Explain: "Поддеревья имеют размеры 100, 99, ..., 1, сумма 100 · 101 / 2 = 5050: время O(n²). Проход с конца делает 99 прибавлений: O(n). Ответ 100 считает только один обход."},
				{ID: "tr2-score", Type: "short", Prompt: "LeetCode 2049. parents = [-1, 0, 0, 1, 1, 2, 4]. Чему равна оценка вершины 2?", Answer: "5", Explain: "size[2] = 2. Без вершины 2 остаются её поддерево {5} и верхняя часть из 7 − 2 = 5 вершин: 1 · 5 = 5. Ответ 1 забывает верхнюю часть n − size[v]."},
				{ID: "tr2-fuel", Type: "short", Prompt: "LeetCode 2477. Путь 0 − 1 − 2 − 3 − 4, столица 0, seats = 3. Сколько литров нужно?", Answer: "5", Explain: "Над вершинами 1, 2, 3, 4 едут 4, 3, 2 и 1 человек. Для 4 человек нужно (4 + 2) // 3 = 2 машины, а 3, 2 и 1 человек едут на одной машине: 2 + 1 + 1 + 1 = 5. Ответ 10 = 4 + 3 + 2 + 1 верен только при seats = 1."},
				{ID: "tr2-leaves", Type: "choice", Prompt: "Codeforces 1843D. n = 3, рёбра 1 − 2 и 2 − 3, корень 1. Сколько листьев в поддереве вершины 1?", Options: []string{"2", "1", "3"}, Answer: 1, Explain: "У корня 1 есть ребёнок 2, поэтому 1 не лист, хотя у него одно ребро. Лист один: вершина 3. Ответ 2 считает листом каждую вершину с одним ребром."},
				{ID: "tr2-loop", Type: "code", Prompt: "Codeforces 1676G: родители a_i < i, bal = +1 для W и −1 для B. Какая строка стоит на месте ___?", Language: "python", Code: `bal = [0] + [1 if c == "W" else -1 for c in s]
___
    bal[par[i]] += bal[i]`, Options: []string{"for i in range(n, 1, -1):", "for i in range(2, n + 1):", "for i in range(1, n):"}, Answer: 0, Explain: "a_i < i, поэтому ребёнок имеет больший номер. Проход от n до 2 передаёт ответ ребёнка раньше, чем родитель передаст свой. Для пути 1 − 2 − 3 и WBW проход от 2 до n передаст bal[2] в bal[1] до прибавления bal[3] и найдёт 2 вместо 1. range(1, n) тоже идёт вперёд, берёт корень 1 и пропускает вершину n. Время O(n), память O(n)."},
				{ID: "tr2-divisible", Type: "short", Prompt: "LeetCode 2872. Путь 0 − 1 − 2, values = [6, 3, 3], k = 3. Какое наибольшее число частей?", Answer: "3", Explain: "Суммы поддеревьев: вершина 2: 3, вершина 1: 6, вершина 0: 12. Все делятся на 3, поэтому режем оба ребра: части 6, 3 и 3. Ответ 1 не режет ни одного ребра."},
			},
			Links: []sets.Link{
				{Title: "Обход Tree: прямой и обратный порядок", URL: "https://ru.wikipedia.org/wiki/%D0%9E%D0%B1%D1%85%D0%BE%D0%B4_%D0%B4%D0%B5%D1%80%D0%B5%D0%B2%D0%B0", Source: "Википедия"},
				{Title: "DFS: время входа и выхода", URL: "https://ru.algorithmica.org/cs/graph-traversals/dfs/", Source: "Алгоритмика"},
				{Title: "Размеры поддеревьев", URL: "https://usaco.guide/silver/intro-tree", Source: "USACO Guide (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "tr2-sizes",
					Title: "Размеры поддеревьев с конца порядка BFS",
					Frames: []sets.FigureFrame{
						{Caption: "Сначала size = 1 у каждой вершины. Идём по порядку BFS с конца: 6, 5, 4, 3, 2, 1.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "родитель", Cells: []any{-1, 0, 0, 1, 1, 2, 4}},
							{Label: "size", Cells: []any{1, 1, 1, 1, 1, 1, 1}},
						}},
						{Caption: "v = 6: size[4] = 1 + 1 = 2.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 6}}},
							{Label: "родитель", Cells: []any{-1, 0, 0, 1, 1, 2, 4}},
							{Label: "size", Cells: []any{1, 1, 1, 1, 2, 1, 1}, Marked: []int{4, 6}},
						}},
						{Caption: "v = 5: size[2] = 1 + 1 = 2.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 5}}},
							{Label: "родитель", Cells: []any{-1, 0, 0, 1, 1, 2, 4}},
							{Label: "size", Cells: []any{1, 1, 2, 1, 2, 1, 1}, Marked: []int{2, 5}},
						}},
						{Caption: "v = 4: size[4] уже готов и равен 2, поэтому size[1] = 1 + 2 = 3.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 4}}},
							{Label: "родитель", Cells: []any{-1, 0, 0, 1, 1, 2, 4}},
							{Label: "size", Cells: []any{1, 3, 2, 1, 2, 1, 1}, Marked: []int{1, 4}},
						}},
						{Caption: "v = 3: size[1] = 3 + 1 = 4. Поддерево 1: вершины 1, 3, 4, 6.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 3}}},
							{Label: "родитель", Cells: []any{-1, 0, 0, 1, 1, 2, 4}},
							{Label: "size", Cells: []any{1, 4, 2, 1, 2, 1, 1}, Marked: []int{1, 3}},
						}},
						{Caption: "v = 2: size[0] = 1 + 2 = 3.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 2}}},
							{Label: "родитель", Cells: []any{-1, 0, 0, 1, 1, 2, 4}},
							{Label: "size", Cells: []any{3, 4, 2, 1, 2, 1, 1}, Marked: []int{0, 2}},
						}},
						{Caption: "v = 1: size[0] = 3 + 4 = 7 = n. 6 прибавлений, по одному на ребро: время O(n).", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 1}}},
							{Label: "родитель", Cells: []any{-1, 0, 0, 1, 1, 2, 4}},
							{Label: "size", Cells: []any{7, 4, 2, 1, 2, 1, 1}, Marked: []int{0, 1}},
						}},
					},
				},
				{
					ID:    "tr2-cut",
					Title: "Без вершины 1: три части",
					Frames: []sets.FigureFrame{
						{Caption: "Удалим вершину 1 и её 3 ребра. size[1] = 4.", Rows: []sets.FigureRow{
							{Label: "Tree", Cells: []any{0, 1, 2, 3, 4, 5, "", "", "", 6}, Marked: []int{1}, Tree: true},
						}},
						{Caption: "Дети дают части size[3] = 1 и size[4] = 2, сверху остаётся 7 − 4 = 3 вершины. Произведение 1 · 2 · 3 = 6.", Rows: []sets.FigureRow{
							{Label: "Tree", Cells: []any{0, 1, 2, 3, 4, 5, "", "", "", 6}, Dimmed: []int{1}, Tree: true},
							{Label: "часть", Cells: []any{"3", "4 и 6", "0, 2, 5"}},
							{Label: "размер", Cells: []any{1, 2, 3}},
						}},
					},
				},
				{
					ID:    "tr2-fuel",
					Title: "Машины на рёбрах",
					Frames: []sets.FigureFrame{
						{Caption: "seats = 1: каждый едет на своей машине. Литров 4 + 2 + 1 + 2 + 1 + 1 = 11, это сумма глубин.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{1, 2, 3, 4, 5, 6}},
							{Label: "size", Cells: []any{4, 2, 1, 2, 1, 1}},
							{Label: "машин", Cells: []any{4, 2, 1, 2, 1, 1}},
						}},
						{Caption: "seats = 2: над вершиной 1 едут 4 человека, им нужно (4 + 1) // 2 = 2 машины. Всего 7 литров.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{1, 2, 3, 4, 5, 6}},
							{Label: "size", Cells: []any{4, 2, 1, 2, 1, 1}},
							{Label: "машин", Cells: []any{2, 1, 1, 1, 1, 1}, Marked: []int{0}},
						}},
					},
				},
			},
		},
		{
			Title:   "Dynamic Programming на поддеревьях",
			Summary: "Состояние вершины хранит лучший ответ её поддерева. Вы выбираете: взять вершину или пропустить её, какой конец отрезка записать в вершину, какую цепочку вниз продолжить. Две самые длинные цепочки детей дают самый длинный путь через вершину.",
			Lesson:  lesson3,
			Tasks: []sets.Task{
				{
					URL:        "https://codeforces.com/problemset/problem/1900/C",
					Title:      "Anji's Binary Tree",
					Difficulty: "medium",
					Topics:     []string{"Tree", "Depth-First Search", "Dynamic Programming"},
					Note:       "Codeforces, рейтинг 1300. dp[v] - наименьшее число замен букв на пути от v до листа: 0 у листа, иначе минимум dp[l] + (s[v] ≠ 'L') и dp[r] + (s[v] ≠ 'R') по существующим детям. Считайте dp с конца порядка BFS от вершины 1, ответ dp[1]. Время: O(n) на набор. Память: O(n) для детей, порядка и dp. В пятом наборе примера путь к листу 3 стоит одну замену, а к листьям 6 и 7 две. Подъём от каждого листа к корню даёт O(n²) на длинной цепочке с листьями.",
					Statement: `<p>Кексич пишет Анджи, а она читает сообщения и не отвечает. Через общего друга он узнал, что Анджи очень любит Binary Tree, и решил решить её задачу, чтобы она обратила на него внимание.</p><p>Анджи дала Кексичу Binary Tree из <code>n</code> вершин. Вершина <code>1</code> - корень, у неё нет родителя. У всех остальных вершин ровно один родитель. У каждой вершины может быть до <code>2</code> детей: левый ребёнок и правый ребёнок. Для каждой вершины Анджи называет Кексичу номера её левого и правого ребёнка или говорит, что такого ребёнка нет.</p><p>Кроме того, на каждой вершине написана буква <code>s<sub>i</sub></code>: <code>U</code>, <code>L</code> или <code>R</code>.</p><p>Кексич начинает путь в корне и на каждом ходу делает следующее:</p><ul><li>Если на текущей вершине буква <code>U</code>, он переходит в её родителя. Если родителя нет, он ничего не делает.</li><li>Если на текущей вершине буква <code>L</code>, он переходит в её левого ребёнка. Если такого ребёнка нет, он ничего не делает.</li><li>Если на текущей вершине буква <code>R</code>, он переходит в её правого ребёнка. Если такого ребёнка нет, он ничего не делает.</li></ul><p>Перед началом пути Кексич может делать операции: одна операция заменяет букву на одной вершине любой из букв <code>U</code>, <code>L</code>, <code>R</code>. Найдите наименьшее число операций перед началом пути, после которых Кексич во время пути в какой-то момент попадёт в лист. Лист - вершина без детей. Неважно, в какой лист он попадёт. Неважно, останется ли он в листе: ему нужно только попасть в него. И неважно, сколько ходов ему понадобится, чтобы попасть в лист.</p><p>Помогите Кексичу решить задачу Анджи, чтобы покорить её сердце и пригласить её в Чачак.</p><h3>Входные данные</h3><p>Каждый тест состоит из нескольких наборов входных данных. В первой строке записано одно целое число <code>t</code> (<code>1 ≤ t ≤ 5 · 10<sup>4</sup></code>): количество наборов. Дальше идут описания наборов.</p><p>В первой строке каждого набора записано одно целое число <code>n</code> (<code>2 ≤ n ≤ 3 · 10<sup>5</sup></code>): число вершин Tree.</p><p>Во второй строке каждого набора записана строка <code>s</code> из <code>n</code> символов: буквы на вершинах. Гарантируется, что <code>s</code> состоит только из символов <code>U</code>, <code>L</code> и <code>R</code>.</p><p>В <code>i</code>-й из следующих <code>n</code> строк записаны два целых числа <code>l<sub>i</sub></code> и <code>r<sub>i</sub></code> (<code>0 ≤ l<sub>i</sub>, r<sub>i</sub> ≤ n</code>): номера левого и правого ребёнка вершины <code>i</code>. Если <code>l<sub>i</sub> = 0</code>, у вершины <code>i</code> нет левого ребёнка. Если <code>r<sub>i</sub> = 0</code>, у вершины <code>i</code> нет правого ребёнка. Гарантируется, что эти данные задают правильное Binary Tree с корнем <code>1</code>.</p><p>Гарантируется, что сумма <code>n</code> по всем наборам не больше <code>3 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого набора выведите одно целое число: наименьшее число операций, после которых Кексич попадёт в лист.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>5
3
LRU
2 3
0 0
0 0
3
ULR
3 2
0 0
0 0
2
LU
0 2
0 0
4
RULR
3 0
0 0
0 4
2 0
7
LLRRRLU
5 2
3 6
0 0
7 0
4 0
0 0
0 0</pre><p><strong>Выходные данные:</strong></p><pre>0
1
1
3
1</pre><h3>Примечание</h3><p>В первом наборе у вершины <code>1</code> левый ребёнок <code>2</code> и правый ребёнок <code>3</code>. У вершин <code>2</code> и <code>3</code> нет детей, поэтому они листья. На вершине <code>1</code> написана буква <code>L</code>, поэтому Кексич перейдёт в вершину <code>2</code>, и операции не нужны.</p><p>Во втором наборе у вершины <code>1</code> левый ребёнок <code>3</code> и правый ребёнок <code>2</code>. Вершины <code>2</code> и <code>3</code> - листья. На вершине <code>1</code> написана буква <code>U</code>, поэтому Кексич должен заменить её на <code>L</code> или <code>R</code>, чтобы попасть в лист.</p><p>В третьем наборе у вершины <code>1</code> есть только правый ребёнок: вершина <code>2</code>. На вершине <code>1</code> написана буква <code>L</code>, поэтому Кексич должен заменить её на <code>R</code>, иначе он останется в вершине <code>1</code>.</p><p>В четвёртом наборе он может заменить <code>3</code> буквы, чтобы на вершинах были буквы <code>LURL</code>. Тогда он попадёт в вершину <code>2</code>.</p><p>В пятом наборе <code>3</code> листа: <code>3</code>, <code>6</code> и <code>7</code>. Чтобы попасть в лист <code>6</code> или в лист <code>7</code>, нужно заменить <code>2</code> буквы. Но если заменить букву на вершине <code>1</code> на <code>R</code>, он попадёт в лист <code>3</code>, поэтому ответ <code>1</code>.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 5 · 10<sup>4</sup></code></li><li><code>2 ≤ n ≤ 3 · 10<sup>5</sup></code></li><li>строка <code>s</code> длины <code>n</code> из символов <code>U</code>, <code>L</code> и <code>R</code></li><li><code>0 ≤ l<sub>i</sub>, r<sub>i</sub> ≤ n</code>, данные задают Binary Tree с корнем <code>1</code></li><li>сумма <code>n</code> по всем наборам не больше <code>3 · 10<sup>5</sup></code></li><li>Ограничение по времени: 2,5 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "5\n3\nLRU\n2 3\n0 0\n0 0\n3\nULR\n3 2\n0 0\n0 0\n2\nLU\n0 2\n0 0\n4\nRULR\n3 0\n0 0\n0 4\n2 0\n7\nLLRRRLU\n5 2\n3 6\n0 0\n7 0\n4 0\n0 0\n0 0\n", Expected: text("0\n1\n1\n3\n1\n")},
						{Input: "2\n2\nUU\n2 0\n0 0\n2\nRL\n0 2\n0 0\n", Expected: text("1\n0\n")},
						{Input: "1\n5\nLUUUU\n2 3\n4 0\n0 0\n5 0\n0 0\n", Expected: text("1\n")},
						{Input: "1\n4\nUUUU\n2 0\n0 3\n4 0\n0 0\n", Expected: text("3\n")},
						{Input: "1\n4\nRURU\n2 3\n0 0\n0 4\n0 0\n", Expected: text("0\n")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/maximum-score-after-applying-operations-on-a-tree/",
					Number:     number(2925),
					Title:      "Maximum Score After Applying Operations on a Tree",
					Difficulty: "medium",
					Topics:     []string{"Tree", "Depth-First Search", "Dynamic Programming"},
					Note:       "На каждом пути от корня к листу должно остаться хотя бы одно значение. keep[v] - наименьшая сумма оставленных значений в поддереве v: values[v] у листа, иначе min(values[v], сумма keep детей). Ответ: сумма всех values минус keep[0]. Время: O(n). Память: O(n) для соседей, порядка и keep. Ответ до 2 · 10^4 · 10^9 = 2 · 10^13, в Python это обычный int. В первом примере keep[0] = 5, ответ 16 − 5 = 11. Новый обход каждого поддерева для его суммы даёт O(n²).",
					Statement: `<p>Дано Tree из <code>n</code> вершин с номерами от <code>0</code> до <code>n - 1</code>, рёбра которого не имеют направления, с корнем в вершине <code>0</code>. Дан двумерный Array целых чисел <code>edges</code> длины <code>n - 1</code>: запись <code>edges[i] = [a<sub>i</sub>, b<sub>i</sub>]</code> означает, что в Tree есть ребро между вершинами <code>a<sub>i</sub></code> и <code>b<sub>i</sub></code>.</p><p>Также дан Array целых чисел <code>values</code> длины <code>n</code>, номера в нём <strong>начинаются с 0</strong>: <code>values[i]</code> - <strong>значение</strong> вершины <code>i</code>.</p><p>Сначала ваш счёт равен 0. За одну операцию вы выбираете вершину <code>i</code>, <strong>прибавляете</strong> <code>values[i]</code> к счёту и <strong>записываете</strong> в <code>values[i]</code> число 0. Все три шага выполняются вместе как одна операция.</p><p>Tree <strong>здоровое</strong>, если для <strong>каждого</strong> листа сумма <strong>текущих</strong> значений на пути от корня до этого листа <strong>не равна</strong> 0.</p><p>Верните <strong>наибольший счёт</strong>, который можно получить, если сделать эту операцию любое число раз так, чтобы в конце Tree осталось <strong>здоровым</strong>.</p><p>Метод класса <code>Solution</code>: <code>def maximumScoreAfterOperations(self, edges: list[list[int]], values: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> edges = [[0,1],[0,2],[0,3],[2,4],[4,5]], values = [5,2,5,2,1,1]
<strong>Выход:</strong> 11
<strong>Пояснение:</strong> применим операцию к вершинам 1, 2, 3, 4 и 5, значения станут [5,0,0,0,0,0]. Листья: вершины 1, 3 и 5.
- Сумма значений на пути от 0 до 1 равна 5.
- Сумма значений на пути от 0 до 3 равна 5.
- Сумма значений на пути от 0 до 5 равна 5.
У каждого листа сумма на пути не равна 0, поэтому Tree здоровое. Счёт - сумма исходных значений выбранных вершин: 2 + 5 + 2 + 1 + 1 = 11.
Можно показать, что 11 - наибольший счёт, который можно получить этими операциями.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> edges = [[0,1],[0,2],[1,3],[1,4],[2,5],[2,6]], values = [20,10,9,7,4,3,5]
<strong>Выход:</strong> 40
<strong>Пояснение:</strong> применим операцию к вершинам 0, 2, 3 и 4, значения станут [0,10,0,0,0,3,5]. Листья: вершины 3, 4, 5 и 6.
- Сумма значений на пути от 0 до 3 равна 10.
- Сумма значений на пути от 0 до 4 равна 10.
- Сумма значений на пути от 0 до 5 равна 3.
- Сумма значений на пути от 0 до 6 равна 5.
У каждого листа сумма на пути не равна 0, поэтому Tree здоровое. Счёт - сумма исходных значений выбранных вершин: 20 + 9 + 7 + 4 = 40.
Можно показать, что 40 - наибольший счёт, который можно получить этими операциями.</pre><h3>Ограничения</h3><ul><li><code>2 ≤ n ≤ 2 · 10<sup>4</sup></code></li><li><code>edges.length == n - 1</code></li><li><code>edges[i].length == 2</code></li><li><code>0 ≤ a<sub>i</sub>, b<sub>i</sub> &lt; n</code></li><li><code>values.length == n</code></li><li><code>1 ≤ values[i] ≤ 10<sup>9</sup></code></li><li>Входные данные таковы, что <code>edges</code> задаёт правильное Tree.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "maximumScoreAfterOperations"}, Cases: []sets.TestCase{
						{Input: "[[[0,1],[0,2],[0,3],[2,4],[4,5]],[5,2,5,2,1,1]]", Expected: text("11")},
						{Input: "[[[0,1],[0,2],[1,3],[1,4],[2,5],[2,6]],[20,10,9,7,4,3,5]]", Expected: text("40")},
						{Input: "[[[0,1]],[5,3]]", Expected: text("5")},
						{Input: "[[[0,1],[1,2]],[4,1,7]]", Expected: text("11")},
						{Input: "[[[0,1],[0,2],[1,3],[1,4],[2,5],[4,6]],[6,2,4,3,1,5,2]]", Expected: text("17")},
						{Input: "[[[0,1],[0,2],[0,3]],[1,10,10,10]]", Expected: text("30")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/minimum-increments-to-equalize-leaf-paths/",
					Number:     number(3593),
					Title:      "Minimum Increments to Equalize Leaf Paths",
					Difficulty: "medium",
					Topics:     []string{"Array", "Tree", "Depth-First Search", "Dynamic Programming"},
					Note:       "best[v] - наибольшая сумма cost на пути от v вниз до листа: cost[v] плюс наибольший best среди детей. Каждый ребёнок c, у которого best[c] меньше наибольшего best среди детей того же родителя, требует одного увеличения: поднимите cost[c] на разницу. Ответ - число таких детей. Время: O(n). Память: O(n) для соседей, порядка и best. Суммы до 10^5 · 10^9 = 10^14. В третьем примере у корня дети 4 (best 7) и 1 (best 5): одно увеличение. Новый поиск самого длинного пути вниз от каждой вершины даёт O(n²).",
					Statement: `<p>Дано целое число <code>n</code> и Tree из <code>n</code> вершин с номерами от <code>0</code> до <code>n - 1</code>, рёбра которого не имеют направления, с корнем в вершине <code>0</code>. Tree задано двумерным Array <code>edges</code> длины <code>n - 1</code>: запись <code>edges[i] = [u<sub>i</sub>, v<sub>i</sub>]</code> означает ребро между вершинами <code>u<sub>i</sub></code> и <code>v<sub>i</sub></code>.</p><p>У каждой вершины <code>i</code> есть стоимость <code>cost[i]</code>: столько стоит пройти через эту вершину.</p><p><strong>Счёт</strong> пути - сумма стоимостей всех вершин на этом пути.</p><p>Нужно сделать счета всех путей <strong>от корня до листа</strong> <strong>равными</strong>. Для этого можно <strong>увеличить</strong> стоимость любого числа вершин на <strong>любое неотрицательное</strong> число.</p><p>Верните <strong>наименьшее</strong> число вершин, стоимость которых нужно увеличить, чтобы счета всех путей от корня до листа стали равными.</p><p>Метод класса <code>Solution</code>: <code>def minIncrease(self, n: int, edges: list[list[int]], cost: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> n = 3, edges = [[0,1],[0,2]], cost = [2,1,3]
<strong>Выход:</strong> 1
<strong>Пояснение:</strong> есть два пути от корня до листа:
- Путь 0 → 1 имеет счёт 2 + 1 = 3.
- Путь 0 → 2 имеет счёт 2 + 3 = 5.
Чтобы счета всех путей стали равны 5, увеличим стоимость вершины 1 на 2.
Увеличена только одна вершина, поэтому ответ 1.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> n = 3, edges = [[0,1],[1,2]], cost = [5,1,4]
<strong>Выход:</strong> 0
<strong>Пояснение:</strong> есть только один путь от корня до листа:
- Путь 0 → 1 → 2 имеет счёт 5 + 1 + 4 = 10.
Путь только один, поэтому все счета уже равны, и ответ 0.</pre><h3>Пример 3</h3><pre><strong>Вход:</strong> n = 5, edges = [[0,4],[0,1],[1,2],[1,3]], cost = [3,4,1,1,7]
<strong>Выход:</strong> 1
<strong>Пояснение:</strong> есть три пути от корня до листа:
- Путь 0 → 4 имеет счёт 3 + 7 = 10.
- Путь 0 → 1 → 2 имеет счёт 3 + 4 + 1 = 8.
- Путь 0 → 1 → 3 имеет счёт 3 + 4 + 1 = 8.
Чтобы счета всех путей стали равны 10, увеличим стоимость вершины 1 на 2. Поэтому ответ 1.</pre><h3>Ограничения</h3><ul><li><code>2 ≤ n ≤ 10<sup>5</sup></code></li><li><code>edges.length == n - 1</code></li><li><code>edges[i] == [u<sub>i</sub>, v<sub>i</sub>]</code></li><li><code>0 ≤ u<sub>i</sub>, v<sub>i</sub> &lt; n</code></li><li><code>cost.length == n</code></li><li><code>1 ≤ cost[i] ≤ 10<sup>9</sup></code></li><li>Входные данные таковы, что <code>edges</code> задаёт правильное Tree.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "minIncrease"}, Cases: []sets.TestCase{
						{Input: "[3,[[0,1],[0,2]],[2,1,3]]", Expected: text("1")},
						{Input: "[3,[[0,1],[1,2]],[5,1,4]]", Expected: text("0")},
						{Input: "[5,[[0,4],[0,1],[1,2],[1,3]],[3,4,1,1,7]]", Expected: text("1")},
						{Input: "[4,[[0,1],[0,2],[0,3]],[1,5,5,5]]", Expected: text("0")},
						{Input: "[4,[[0,1],[0,2],[0,3]],[1,2,3,4]]", Expected: text("2")},
						{Input: "[7,[[0,1],[0,2],[1,3],[1,4],[2,5],[4,6]],[1,2,3,4,1,2,5]]", Expected: text("2")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/1528/A",
					Title:      "Parsa's Humongous Tree",
					Difficulty: "medium",
					Topics:     []string{"Tree", "Depth-First Search", "Dynamic Programming", "Greedy"},
					Note:       "Codeforces, рейтинг 1600. Есть лучший ответ, в котором каждое a_v равно l_v или r_v. dp[v][0] и dp[v][1] - лучшая сумма в поддереве v при a_v = l_v и при a_v = r_v. Если у родителя значение x, ребёнок c добавляет max(dp[c][0] + |x − l_c|, dp[c][1] + |x − r_c|). Время: O(n) на набор. Память: O(n) для соседей, порядка и dp. Ответ одного набора меньше 10^5 · 10^9 = 10^14. Во втором наборе примера a_2 = 4 даёт 3 + 5 = 8, ответ 8. Brute Force по всем 2^n выборам концов слишком долгий уже при n = 40.",
					Statement: `<p>У Парсы есть огромное Tree из <code>n</code> вершин.</p><p>На каждой вершине <code>v</code> он написал два целых числа <code>l<sub>v</sub></code> и <code>r<sub>v</sub></code>.</p><p>Чтобы Tree Парсы выглядело ещё величественнее, Нима хочет дать каждой вершине <code>v</code> число <code>a<sub>v</sub></code> (<code>l<sub>v</sub> ≤ a<sub>v</sub> ≤ r<sub>v</sub></code>) так, чтобы красота Tree была наибольшей.</p><p>У Нимы странное чувство красоты. Он называет <strong>красотой</strong> Tree сумму <code>|a<sub>u</sub> − a<sub>v</sub>|</code> по всем рёбрам <code>(u, v)</code> Tree.</p><p>Tree Парсы слишком большое, и Нима не может найти наибольшую красоту сам. Найдите наибольшую возможную красоту Tree Парсы.</p><h3>Входные данные</h3><p>В первой строке записано целое число <code>t</code> (<code>1 ≤ t ≤ 250</code>): количество наборов входных данных. Дальше идут описания наборов.</p><p>В первой строке каждого набора записано одно целое число <code>n</code> (<code>2 ≤ n ≤ 10<sup>5</sup></code>): число вершин Tree Парсы.</p><p>В <code>i</code>-й из следующих <code>n</code> строк записаны два целых числа <code>l<sub>i</sub></code> и <code>r<sub>i</sub></code> (<code>1 ≤ l<sub>i</sub> ≤ r<sub>i</sub> ≤ 10<sup>9</sup></code>).</p><p>В каждой из следующих <code>n − 1</code> строк записаны два целых числа <code>u</code> и <code>v</code> (<code>1 ≤ u, v ≤ n</code>, <code>u ≠ v</code>): в Tree Парсы есть ребро между вершинами <code>u</code> и <code>v</code>.</p><p>Гарантируется, что данный Graph - Tree.</p><p>Гарантируется, что сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого набора выведите наибольшую возможную красоту Tree Парсы.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>3
2
1 6
3 8
1 2
3
1 3
4 6
7 9
1 2
2 3
6
3 14
12 20
12 19
2 12
10 17
3 17
3 2
6 5
1 5
2 6
4 6</pre><p><strong>Выходные данные:</strong></p><pre>7
8
62</pre><h3>Примечание</h3><p>В первом наборе можно взять, например, <code>a = {1, 8}</code>: красота <code>|1 − 8| = 7</code>.</p><p>Во втором наборе можно взять, например, <code>a = {1, 5, 9}</code>: красота <code>|1 − 5| + |5 − 9| = 8</code>.</p><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 250</code></li><li><code>2 ≤ n ≤ 10<sup>5</sup></code></li><li><code>1 ≤ l<sub>i</sub> ≤ r<sub>i</sub> ≤ 10<sup>9</sup></code></li><li><code>1 ≤ u, v ≤ n</code>, <code>u ≠ v</code>, рёбра задают Tree</li><li>сумма <code>n</code> по всем наборам не больше <code>2 · 10<sup>5</sup></code></li><li>Ограничение по времени: 1 секунда на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "3\n2\n1 6\n3 8\n1 2\n3\n1 3\n4 6\n7 9\n1 2\n2 3\n6\n3 14\n12 20\n12 19\n2 12\n10 17\n3 17\n3 2\n6 5\n1 5\n2 6\n4 6\n", Expected: text("7\n8\n62\n")},
						{Input: "3\n2\n5 5\n5 5\n1 2\n3\n1 10\n4 6\n1 10\n1 2\n2 3\n3\n1 1000000000\n1 1000000000\n1 1000000000\n2 1\n3 2\n", Expected: text("0\n12\n1999999998\n")},
						{Input: "1\n4\n1 1\n10 10\n3 3\n1 100\n1 2\n1 3\n1 4\n", Expected: text("110\n")},
						{Input: "1\n4\n1 100\n50 50\n50 50\n1 1\n1 2\n1 3\n1 4\n", Expected: text("199\n")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/longest-path-with-different-adjacent-characters/",
					Number:     number(2246),
					Title:      "Longest Path With Different Adjacent Characters",
					Difficulty: "hard",
					Topics:     []string{"Array", "String", "Tree", "Depth-First Search"},
					Note:       "down[v] - число вершин самой длинной цепочки вниз от v, где соседние буквы разные. Учитывайте только детей c с s[c] ≠ s[v]: down[v] = 1 + top1, а путь через v содержит 1 + top1 + top2 вершин, где top1 ≥ top2 - две наибольшие down таких детей (0, если ребёнка нет). Ответ - максимум по всем v. Время: O(n). Память: O(n) для детей, порядка, top1 и top2. В первом примере путь 0 → 1 → 3 даёт 3. BFS от каждой вершины даёт O(n²).",
					Statement: `<p>Дано <strong>Tree</strong> (то есть связный Graph без циклов, рёбра которого не имеют направления) с <strong>корнем</strong> в вершине <code>0</code> из <code>n</code> вершин с номерами от <code>0</code> до <code>n - 1</code>. Tree задано Array <code>parent</code> длины <code>n</code>, номера в нём <strong>начинаются с 0</strong>: <code>parent[i]</code> - родитель вершины <code>i</code>. Вершина <code>0</code> - корень, поэтому <code>parent[0] == -1</code>.</p><p>Также дана строка <code>s</code> длины <code>n</code>: <code>s[i]</code> - буква вершины <code>i</code>.</p><p>Верните <em>длину <strong>самого длинного пути</strong> в Tree, на котором никакие две <strong>соседние</strong> вершины не имеют одинаковых букв.</em></p><p>Метод класса <code>Solution</code>: <code>def longestPath(self, parent: list[int], s: str) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> parent = [-1,0,0,1,1,2], s = "abacbe"
<strong>Выход:</strong> 3
<strong>Пояснение:</strong> самый длинный путь, на котором у каждых двух соседних вершин разные буквы: 0 -&gt; 1 -&gt; 3. Длина этого пути 3, поэтому ответ 3.
Можно доказать, что более длинного подходящего пути нет.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> parent = [-1,0,0,0], s = "aabc"
<strong>Выход:</strong> 3
<strong>Пояснение:</strong> самый длинный путь, на котором у каждых двух соседних вершин разные буквы: 2 -&gt; 0 -&gt; 3. Длина этого пути 3, поэтому ответ 3.</pre><h3>Ограничения</h3><ul><li><code>n == parent.length == s.length</code></li><li><code>1 ≤ n ≤ 10<sup>5</sup></code></li><li><code>0 ≤ parent[i] ≤ n - 1</code> для всех <code>i ≥ 1</code></li><li><code>parent[0] == -1</code></li><li><code>parent</code> задаёт правильное Tree.</li><li><code>s</code> состоит только из строчных английских букв.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "longestPath"}, Cases: []sets.TestCase{
						{Input: "[[-1,0,0,1,1,2],\"abacbe\"]", Expected: text("3")},
						{Input: "[[-1,0,0,0],\"aabc\"]", Expected: text("3")},
						{Input: "[[-1],\"z\"]", Expected: text("1")},
						{Input: "[[-1,0,1,2],\"aaaa\"]", Expected: text("1")},
						{Input: "[[-1,0,1,1,2,3],\"aabcde\"]", Expected: text("5")},
						{Input: "[[-1,0,0,1,1,2,4],\"abcabab\"]", Expected: text("5")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "tr3-take", Type: "short", Prompt: "Путь 0 − 1 − 2 − 3, w = [3, 5, 4, 1]. Нельзя брать две соседние вершины. Какая наибольшая сумма?", Answer: "7", Explain: "Вершины 0 и 2 не соседи: 3 + 4 = 7. Если сначала взять самое большое число 5, то вершины 0 и 2 брать нельзя, остаётся вершина 3: 5 + 1 = 6. Два состояния у корня 0: взять = 3 + 4 = 7, пропустить = 6."},
				{ID: "tr3-take-code", Type: "code", Prompt: "Что стоит на месте ___ в переходе «взять или пропустить»?", Language: "python", Code: `for v in reversed(order):
    p = parent[v]
    if p != -1:
        take[p] += ___
        skip[p] += max(take[v], skip[v])`, Options: []string{"take[v]", "skip[v]", "max(take[v], skip[v])"}, Answer: 1, Explain: "Если взять родителя, ребёнка брать нельзя: прибавляем skip[v]. take[v] разрешил бы две соседние вершины. max(...) нужен только у пропущенного родителя. Время O(n), память O(n)."},
				{ID: "tr3-skip-bug", Type: "choice", Prompt: "Ученик пишет skip[p] += skip[v] вместо max(take[v], skip[v]). Путь 0 − 1, w = [1, 9], корень 0. Какой ответ даст его код?", Options: []string{"9", "10", "1"}, Answer: 2, Explain: "take = [1, 9], а skip[0] = skip[1] = 0, поэтому ответ max(1, 0) = 1. Верно skip[0] = max(9, 0) = 9, ответ 9. Ответ 10 берёт обе соседние вершины."},
				{ID: "tr3-ends", Type: "choice", Prompt: "Codeforces 1528A. Вершина с отрезком [2, 7] соединена с вершинами, у которых значения 4 и 6. Какое значение вершины даёт наибольшую сумму |x − 4| + |x − 6|?", Options: []string{"2", "4", "7"}, Answer: 0, Explain: "x = 2 даёт 2 + 4 = 6, x = 7 даёт 3 + 1 = 4, x = 4 даёт 0 + 2 = 2. Лучший x всегда на конце отрезка, но не обязательно на большем: здесь на левом. Поэтому нужны оба состояния, l и r. Ответ 7 берёт больший конец без проверки."},
				{ID: "tr3-chain", Type: "short", Prompt: "parent = [-1, 0, 1, 1, 3, 0]. Сколько вершин в самой длинной цепочке вниз от вершины 1?", Answer: "3", Explain: "Цепочка 1 → 3 → 4: down[1] = 1 + down[3] = 1 + 2 = 3. Ответ 4 - это путь 2, 1, 3, 4 через вершину 1 из двух цепочек: 1 + 2 + 1."},
				{ID: "tr3-through", Type: "short", Prompt: "Звезда: центр 0 и 4 листа. Сколько вершин в самом длинном пути через вершину 0?", Answer: "3", Explain: "Путь спускается только к двум детям: 1 + top1 + top2 = 1 + 1 + 1 = 3. Ответ 5 складывает всех детей, но путь не может разветвляться."},
				{ID: "tr3-letters", Type: "short", Prompt: "LeetCode 2246. parent = [-1, 0, 1, 0, 3], s = \"abbba\". Сколько вершин в самом длинном пути с разными соседними буквами?", Answer: "4", Explain: "Ребро 1 − 2 соединяет b и b, его использовать нельзя. Остаются рёбра 0 − 1, 0 − 3 и 3 − 4: путь 1, 0, 3, 4 из 4 вершин. Ответ 5 берёт и вершину 2."},
				{ID: "tr3-keep", Type: "short", Prompt: "LeetCode 2925. Путь 0 − 1 − 2, values = [4, 1, 7]. Какой наибольший счёт?", Answer: "11", Explain: "На единственном пути нужно оставить хотя бы одно значение, выгоднее оставить 1: keep[0] = min(4, min(1, 7)) = 1. Счёт 12 − 1 = 11. Ответ 12 оставляет путь с суммой 0."},
			},
			Links: []sets.Link{
				{Title: "Динамика по поддеревьям", URL: "https://neerc.ifmo.ru/wiki/index.php?title=%D0%94%D0%B8%D0%BD%D0%B0%D0%BC%D0%B8%D0%BA%D0%B0_%D0%BF%D0%BE_%D0%BF%D0%BE%D0%B4%D0%B4%D0%B5%D1%80%D0%B5%D0%B2%D1%8C%D1%8F%D0%BC", Source: "Викиконспекты ИТМО"},
				{Title: "Взять или пропустить вершину", URL: "https://apmath-spbu.github.io/basics/dynamic_programming/independent_set_of_maximum_weight_in_tree/index.html", Source: "Курс СПбГУ"},
				{Title: "Независимое множество", URL: "https://ru.wikipedia.org/wiki/%D0%97%D0%B0%D0%B4%D0%B0%D1%87%D0%B0_%D0%BE_%D0%BD%D0%B5%D0%B7%D0%B0%D0%B2%D0%B8%D1%81%D0%B8%D0%BC%D0%BE%D0%BC_%D0%BC%D0%BD%D0%BE%D0%B6%D0%B5%D1%81%D1%82%D0%B2%D0%B5", Source: "Википедия"},
				{Title: "Dynamic Programming on Trees", URL: "https://usaco.guide/gold/dp-trees", Source: "USACO Guide (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "tr3-take-skip",
					Title: "Взять или пропустить: w = [5, 2, 4, 3, 1, 6, 2]",
					Frames: []sets.FigureFrame{
						{Caption: "Сначала взять[v] = w[v] и пропустить[v] = 0. Листья 3, 5 и 6 уже готовы.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "w", Cells: []any{5, 2, 4, 3, 1, 6, 2}},
							{Label: "взять", Cells: []any{5, 2, 4, 3, 1, 6, 2}},
							{Label: "пропустить", Cells: []any{0, 0, 0, 0, 0, 0, 0}},
						}},
						{Caption: "v = 6: взять[4] += пропустить[6] = 0, пропустить[4] += max(2, 0) = 2.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 6}}},
							{Label: "w", Cells: []any{5, 2, 4, 3, 1, 6, 2}},
							{Label: "взять", Cells: []any{5, 2, 4, 3, 1, 6, 2}, Marked: []int{4}},
							{Label: "пропустить", Cells: []any{0, 0, 0, 0, 2, 0, 0}, Marked: []int{4}},
						}},
						{Caption: "v = 5: пропустить[2] = max(6, 0) = 6.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 5}}},
							{Label: "w", Cells: []any{5, 2, 4, 3, 1, 6, 2}},
							{Label: "взять", Cells: []any{5, 2, 4, 3, 1, 6, 2}, Marked: []int{2}},
							{Label: "пропустить", Cells: []any{0, 0, 6, 0, 2, 0, 0}, Marked: []int{2}},
						}},
						{Caption: "v = 4: взять[1] = 2 + 2 = 4, пропустить[1] = max(1, 2) = 2.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 4}}},
							{Label: "w", Cells: []any{5, 2, 4, 3, 1, 6, 2}},
							{Label: "взять", Cells: []any{5, 4, 4, 3, 1, 6, 2}, Marked: []int{1}},
							{Label: "пропустить", Cells: []any{0, 2, 6, 0, 2, 0, 0}, Marked: []int{1}},
						}},
						{Caption: "v = 3: пропустить[1] = 2 + max(3, 0) = 5.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 3}}},
							{Label: "w", Cells: []any{5, 2, 4, 3, 1, 6, 2}},
							{Label: "взять", Cells: []any{5, 4, 4, 3, 1, 6, 2}, Marked: []int{1}},
							{Label: "пропустить", Cells: []any{0, 5, 6, 0, 2, 0, 0}, Marked: []int{1}},
						}},
						{Caption: "v = 2: взять[0] = 5 + 6 = 11, пропустить[0] = max(4, 6) = 6.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 2}}},
							{Label: "w", Cells: []any{5, 2, 4, 3, 1, 6, 2}},
							{Label: "взять", Cells: []any{11, 4, 4, 3, 1, 6, 2}, Marked: []int{0}},
							{Label: "пропустить", Cells: []any{6, 5, 6, 0, 2, 0, 0}, Marked: []int{0}},
						}},
						{Caption: "v = 1: взять[0] = 11 + 5 = 16, пропустить[0] = 6 + max(4, 5) = 11. Ответ max(16, 11) = 16: вершины 0, 3, 5, 6.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 1}}},
							{Label: "w", Cells: []any{5, 2, 4, 3, 1, 6, 2}},
							{Label: "взять", Cells: []any{16, 4, 4, 3, 1, 6, 2}, Marked: []int{0}},
							{Label: "пропустить", Cells: []any{11, 5, 6, 0, 2, 0, 0}, Marked: []int{0}},
						}},
					},
				},
				{
					ID:    "tr3-ends",
					Title: "Только концы отрезка: путь 0 − 1 − 2",
					Frames: []sets.FigureFrame{
						{Caption: "a_1 = 4: соседям выгоднее 10, |4 − 10| = 6 с каждой стороны. Сумма 12.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2}},
							{Label: "l", Cells: []any{1, 4, 1}},
							{Label: "r", Cells: []any{10, 6, 10}},
							{Label: "a", Cells: []any{10, 4, 10}, Marked: []int{1}},
						}},
						{Caption: "a_1 = 6: соседям выгоднее 1, по 5 с каждой стороны. Сумма 10. Ответ 12.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2}},
							{Label: "l", Cells: []any{1, 4, 1}},
							{Label: "r", Cells: []any{10, 6, 10}},
							{Label: "a", Cells: []any{1, 6, 1}, Marked: []int{1}},
						}},
					},
				},
				{
					ID:    "tr3-chains",
					Title: "Цепочки вниз и путь через вершину",
					Frames: []sets.FigureFrame{
						{Caption: "v = 6, лист: down[6] = 1 + 0 = 1. Это top1 у родителя 4. Знак ? - вершина, до которой цикл ещё не дошёл.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 6}}},
							{Label: "down", Cells: []any{"?", "?", "?", "?", "?", "?", 1}, Marked: []int{6}},
							{Label: "top1", Cells: []any{0, 0, 0, 0, 1, 0, 0}, Marked: []int{4}},
							{Label: "top2", Cells: []any{0, 0, 0, 0, 0, 0, 0}},
						}},
						{Caption: "v = 5, лист: down[5] = 1 + 0 = 1. Это top1 у родителя 2.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 5}}},
							{Label: "down", Cells: []any{"?", "?", "?", "?", "?", 1, 1}, Marked: []int{5}},
							{Label: "top1", Cells: []any{0, 0, 1, 0, 1, 0, 0}, Marked: []int{2}},
							{Label: "top2", Cells: []any{0, 0, 0, 0, 0, 0, 0}},
						}},
						{Caption: "v = 4: down[4] = 1 + 1 = 2, это top1 у вершины 1.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 4}}},
							{Label: "down", Cells: []any{"?", "?", "?", "?", 2, 1, 1}, Marked: []int{4}},
							{Label: "top1", Cells: []any{0, 2, 1, 0, 1, 0, 0}, Marked: []int{1}},
							{Label: "top2", Cells: []any{0, 0, 0, 0, 0, 0, 0}},
						}},
						{Caption: "v = 3: лист даёт 1, это top2 у вершины 1. Путь через 1: 1 + 2 + 1 = 4 вершины.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 3}}},
							{Label: "down", Cells: []any{"?", "?", "?", 1, 2, 1, 1}, Marked: []int{3}},
							{Label: "top1", Cells: []any{0, 2, 1, 0, 1, 0, 0}, Marked: []int{1}},
							{Label: "top2", Cells: []any{0, 1, 0, 0, 0, 0, 0}, Marked: []int{1}},
						}},
						{Caption: "v = 2: down[2] = 1 + 1 = 2, это top1 у корня 0.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 2}}},
							{Label: "down", Cells: []any{"?", "?", 2, 1, 2, 1, 1}, Marked: []int{2}},
							{Label: "top1", Cells: []any{2, 2, 1, 0, 1, 0, 0}, Marked: []int{0}},
							{Label: "top2", Cells: []any{0, 1, 0, 0, 0, 0, 0}},
						}},
						{Caption: "v = 1: down[1] = 1 + 2 = 3 больше top1[0] = 2. Старый top1 переходит в top2: top1[0] = 3, top2[0] = 2.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 1}}},
							{Label: "down", Cells: []any{"?", 3, 2, 1, 2, 1, 1}, Marked: []int{1}},
							{Label: "top1", Cells: []any{3, 2, 1, 0, 1, 0, 0}, Marked: []int{0}},
							{Label: "top2", Cells: []any{2, 1, 0, 0, 0, 0, 0}, Marked: []int{0}},
						}},
						{Caption: "v = 0: через корень 1 + 3 + 2 = 6 вершин, путь 6, 4, 1, 0, 2, 5. Это самый длинный путь: 5 рёбер.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}, Pointers: &sets.Pointers{{Name: "v", Cell: 0}}},
							{Label: "down", Cells: []any{4, 3, 2, 1, 2, 1, 1}, Marked: []int{0}},
							{Label: "top1", Cells: []any{3, 2, 1, 0, 1, 0, 0}, Marked: []int{0}},
							{Label: "top2", Cells: []any{2, 1, 0, 0, 0, 0, 0}, Marked: []int{0}},
						}},
					},
				},
			},
		},
		{
			Title:   "Пути, диаметр и центр",
			Summary: "Путь между двумя вершинами Tree единственный, его восстанавливают по Array родителей. Два BFS находят диаметр, самый длинный путь, а его середина даёт центр с наименьшей высотой. Вы считаете листья и диаметр после соединения двух Tree.",
			Lesson:  lesson4,
			Tasks: []sets.Task{
				{
					URL:        "https://codeforces.com/problemset/problem/1905/B",
					Title:      "Begginer's Zelda",
					Difficulty: "easy",
					Topics:     []string{"Tree", "Greedy"},
					Note:       "Codeforces, рейтинг 1100. Лист - вершина степени 1. Одна операция сжимает путь и убирает не больше двух листьев, поэтому для L листьев нужно не меньше ⌈L / 2⌉ операций, и столько хватает. Время: O(n) на набор для степеней. Память: O(n) для n счётчиков. Для звезды с 5 листьями нужно 3 операции. Проверка путей между всеми парами листьев даёт O(n²) на звезде.",
					Statement: `<p>Дано Tree<sup>†</sup>. За одну zelda-операцию можно сделать следующее:</p><ul><li>выбрать две вершины Tree <code>u</code> и <code>v</code>;</li><li>сжать все вершины на пути от <code>u</code> до <code>v</code> в одну вершину. Другими словами, все вершины на пути от <code>u</code> до <code>v</code> удаляются из Tree, и создаётся новая вершина <code>w</code>. Затем каждая вершина <code>s</code>, у которой было ребро к какой-то вершине на пути от <code>u</code> до <code>v</code>, получает ребро к вершине <code>w</code>.</li></ul><p>Найдите наименьшее число zelda-операций, после которых в Tree останется только одна вершина.</p><p><sup>†</sup>Tree - это связный Graph без циклов, рёбра которого не имеют направления.</p><h3>Входные данные</h3><p>Каждый тест состоит из нескольких наборов входных данных. В первой строке записано одно целое число <code>t</code> (<code>1 ≤ t ≤ 10<sup>4</sup></code>): количество наборов. Дальше идут описания наборов.</p><p>В первой строке каждого набора записано одно целое число <code>n</code> (<code>2 ≤ n ≤ 10<sup>5</sup></code>): число вершин.</p><p>В <code>i</code>-й из следующих <code>n − 1</code> строк записаны два целых числа <code>u<sub>i</sub></code> и <code>v<sub>i</sub></code> (<code>1 ≤ u<sub>i</sub>, v<sub>i</sub> ≤ n</code>, <code>u<sub>i</sub> ≠ v<sub>i</sub></code>): номера вершин, которые соединяет <code>i</code>-е ребро.</p><p>Гарантируется, что эти рёбра образуют Tree.</p><p>Гарантируется, что сумма <code>n</code> по всем наборам не больше <code>10<sup>5</sup></code>.</p><h3>Выходные данные</h3><p>Для каждого набора выведите одно целое число: наименьшее число zelda-операций, после которых в Tree останется только одна вершина.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>4
4
1 2
1 3
3 4
9
3 1
3 5
3 2
5 6
6 7
7 8
7 9
6 4
7
1 2
1 3
2 4
4 5
3 6
2 7
6
1 2
1 3
1 4
4 5
2 6</pre><p><strong>Выходные данные:</strong></p><pre>1
3
2
2</pre><h3>Примечание</h3><p>В первом наборе достаточно одной zelda-операции для вершин <code>2</code> и <code>4</code>.</p><p>Во втором наборе можно сделать такие zelda-операции:</p><ol><li><code>u = 2, v = 1</code>. Пусть новая вершина получит номер <code>w = 10</code>;</li><li><code>u = 4, v = 9</code>. Пусть новая вершина получит номер <code>w = 11</code>;</li><li><code>u = 8, v = 10</code>. После этой операции в Tree только одна вершина.</li></ol><h3>Ограничения</h3><ul><li><code>1 ≤ t ≤ 10<sup>4</sup></code></li><li><code>2 ≤ n ≤ 10<sup>5</sup></code></li><li><code>1 ≤ u<sub>i</sub>, v<sub>i</sub> ≤ n</code>, <code>u<sub>i</sub> ≠ v<sub>i</sub></code></li><li>рёбра образуют Tree</li><li>сумма <code>n</code> по всем наборам не больше <code>10<sup>5</sup></code></li><li>Ограничение по времени: 1 секунда на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "4\n4\n1 2\n1 3\n3 4\n9\n3 1\n3 5\n3 2\n5 6\n6 7\n7 8\n7 9\n6 4\n7\n1 2\n1 3\n2 4\n4 5\n3 6\n2 7\n6\n1 2\n1 3\n1 4\n4 5\n2 6\n", Expected: text("1\n3\n2\n2\n")},
						{Input: "3\n2\n1 2\n6\n1 2\n1 3\n1 4\n1 5\n1 6\n7\n2 1\n1 3\n4 2\n2 5\n6 3\n5 7\n", Expected: text("1\n3\n2\n")},
						{Input: "1\n4\n2 1\n2 3\n2 4\n", Expected: text("2\n")},
						{Input: "1\n5\n1 2\n2 3\n3 4\n4 5\n", Expected: text("1\n")},
					}},
				},
				{
					URL:        "https://codeforces.com/problemset/problem/690/C2",
					Title:      "Brain Network (medium)",
					Difficulty: "medium",
					Topics:     []string{"Tree", "Depth-First Search", "Graph Theory"},
					Note:       "Codeforces, рейтинг 1500. Задержка мозга - диаметр Tree в рёбрах. BFS от вершины 1 находит самую дальнюю вершину a, второй BFS от a находит самую дальнюю вершину b, и расстояние до b - ответ. Время: O(n), два BFS. Память: O(n) для списка соседей, расстояний и Deque. Во втором примере BFS от 1 находит вершину 5 на расстоянии 3, BFS от 5 находит вершину 1 на расстоянии 3: ответ 3. BFS от каждой вершины даёт O(n²).",
					Statement: `<p>Дальнейшие исследования мыслительных процессов зомби дали интересные результаты. Как мы знаем из предыдущей задачи, нервная система зомби состоит из <code>n</code> мозгов и <code>m</code> мозговых связей, которые соединяют некоторые пары мозгов. Учёные заметили, что умственные способности зомби зависят в основном от устройства его нервной системы. Точнее, назовём расстоянием между двумя мозгами <code>u</code> и <code>v</code> (<code>1 ≤ u, v ≤ n</code>) наименьшее число мозговых связей, через которые проходит мысль при передаче между этими двумя мозгами. <strong>Задержка мозга</strong> зомби - это наибольшее расстояние между двумя его мозгами. Учёные предполагают, что задержка мозга - главный параметр, от которого зависит, насколько умён зомби. Помогите им проверить это предположение: напишите программу, которая считает задержку мозга нервной системы.</p><p>Можно считать, что каждая нервная система во входных данных правильная, то есть выполняет условия (1) и (2) из простой версии задачи:</p><ol><li>Любые два мозга могут обмениваться мыслями, может быть, через другие мозги.</li><li>Нет лишних мозговых связей: если удалить любую связь, условие (1) перестанет выполняться.</li></ol><p>Поэтому мозги и связи правильной нервной системы образуют Tree.</p><h3>Входные данные</h3><p>В первой строке записаны два целых числа <code>n</code> и <code>m</code> через пробел (<code>1 ≤ n, m ≤ 100000</code>): число мозгов (у них удобные номера от <code>1</code> до <code>n</code>) и число мозговых связей в нервной системе. В следующих <code>m</code> строках описаны мозговые связи. Каждая связь задана парой мозгов <code>a b</code>, которые она соединяет (<code>1 ≤ a, b ≤ n</code> и <code>a ≠ b</code>).</p><h3>Выходные данные</h3><p>Выведите одно число: задержку мозга.</p><h3>Пример 1</h3><p><strong>Входные данные:</strong></p><pre>4 3
1 2
1 3
1 4</pre><p><strong>Выходные данные:</strong></p><pre>2</pre><h3>Пример 2</h3><p><strong>Входные данные:</strong></p><pre>5 4
1 2
2 3
3 4
3 5</pre><p><strong>Выходные данные:</strong></p><pre>3</pre><h3>Ограничения</h3><ul><li><code>1 ≤ n, m ≤ 100000</code></li><li><code>1 ≤ a, b ≤ n</code>, <code>a ≠ b</code></li><li>нервная система правильная: мозги и связи образуют Tree</li><li>Ограничение по времени: 2 секунды на тест.</li><li>Ограничение по памяти: 256 мегабайт на тест.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{
						{Input: "4 3\n1 2\n1 3\n1 4\n", Expected: text("2\n")},
						{Input: "5 4\n1 2\n2 3\n3 4\n3 5\n", Expected: text("3\n")},
						{Input: "2 1\n1 2\n", Expected: text("1\n")},
						{Input: "7 6\n2 1\n1 3\n4 2\n2 5\n6 3\n5 7\n", Expected: text("5\n")},
						{Input: "6 5\n1 2\n2 3\n2 4\n3 5\n4 6\n", Expected: text("4\n")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/minimum-height-trees/",
					Number:     number(310),
					Title:      "Minimum Height Trees",
					Difficulty: "medium",
					Topics:     []string{"Tree", "Breadth-First Search", "Graph Theory"},
					Note:       "Корни с наименьшей высотой лежат в середине диаметра. Два BFS находят концы a и b диаметра, Array родителей второго BFS восстанавливает путь из d рёбер. При чётном d центр один, при нечётном два. Время: O(n). Память: O(n) для соседей, расстояний и родителей. Во втором примере путь 0 − 3 − 4 − 5, d = 3: центры 3 и 4. BFS от каждой вершины даёт O(n²).",
					Statement: `<p>Tree - это Graph, рёбра которого не имеют направления, и в нём любые две вершины соединены <em>ровно</em> одним путём. Другими словами, любой связный Graph без простых циклов - это Tree.</p><p>Дано Tree из <code>n</code> вершин с номерами от <code>0</code> до <code>n - 1</code> и Array <code>edges</code> из <code>n - 1</code> рёбер: запись <code>edges[i] = [a<sub>i</sub>, b<sub>i</sub>]</code> означает, что в Tree есть ребро без направления между вершинами <code>a<sub>i</sub></code> и <code>b<sub>i</sub></code>. Корнем можно выбрать любую вершину Tree. Если выбрать корнем вершину <code>x</code>, у Tree с этим корнем будет высота <code>h</code>. Среди всех таких Tree с корнем те, у которых высота наименьшая (то есть <code>min(h)</code>), называют <strong>Tree наименьшей высоты</strong> (MHT).</p><p>Верните <em>список номеров корней всех <strong>MHT</strong></em>. Ответ можно вернуть в <strong>любом порядке</strong>.</p><p><strong>Высота</strong> Tree с корнем - это число рёбер на самом длинном пути вниз от корня до листа.</p><p>Метод класса <code>Solution</code>: <code>def findMinHeightTrees(self, n: int, edges: list[list[int]]) -&gt; list[int]</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> n = 4, edges = [[1,0],[1,2],[1,3]]
<strong>Выход:</strong> [1]
<strong>Пояснение:</strong> если корень - вершина 1, высота Tree равна 1. Это единственное MHT.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> n = 6, edges = [[3,0],[3,1],[3,2],[3,4],[5,4]]
<strong>Выход:</strong> [3,4]</pre><h3>Ограничения</h3><ul><li><code>1 ≤ n ≤ 2 · 10<sup>4</sup></code></li><li><code>edges.length == n - 1</code></li><li><code>0 ≤ a<sub>i</sub>, b<sub>i</sub> &lt; n</code></li><li><code>a<sub>i</sub> != b<sub>i</sub></code></li><li>Все пары <code>(a<sub>i</sub>, b<sub>i</sub>)</code> различны.</li><li><strong>Гарантируется</strong>, что входные данные задают Tree и в нём <strong>нет повторённых</strong> рёбер.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "findMinHeightTrees", Order: "any"}, Cases: []sets.TestCase{
						{Input: "[4,[[1,0],[1,2],[1,3]]]", Expected: text("[1]")},
						{Input: "[6,[[3,0],[3,1],[3,2],[3,4],[5,4]]]", Expected: text("[3,4]")},
						{Input: "[1,[]]", Expected: text("[0]")},
						{Input: "[2,[[0,1]]]", Expected: text("[0,1]")},
						{Input: "[7,[[0,1],[0,2],[1,3],[1,4],[2,5],[4,6]]]", Expected: text("[0,1]")},
						{Input: "[5,[[0,1],[1,2],[2,3],[3,4]]]", Expected: text("[2]")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/most-profitable-path-in-a-tree/",
					Number:     number(2467),
					Title:      "Most Profitable Path in a Tree",
					Difficulty: "medium",
					Topics:     []string{"Array", "Tree", "Depth-First Search", "Breadth-First Search", "Graph Theory"},
					Note:       "Сначала BFS от 0: родитель и глубина каждой вершины, глубина - время Алисы. Время Боба: шаги по родителям от bob до 0. Доход Алисы в вершине: amount, если она пришла раньше, amount / 2 при равенстве, 0 позже. Сложите доход сверху вниз и возьмите максимум по листьям. Лист - вершина v ≠ 0 степени 1: корень 0 со степенью 1 не лист. n - число вершин, h - наибольшая глубина. Время: O(n). Память: O(n) для соседей, родителей, глубин и доходов. В первом примере: −2 + 2 + 0 + 6 = 6. Отдельный проход по пути к каждому листу даёт O(n · h).",
					Statement: `<p>Дано Tree из <code>n</code> вершин с номерами от <code>0</code> до <code>n - 1</code>, рёбра которого не имеют направления, с корнем в вершине <code>0</code>. Дан двумерный Array целых чисел <code>edges</code> длины <code>n - 1</code>: запись <code>edges[i] = [a<sub>i</sub>, b<sub>i</sub>]</code> означает, что в Tree есть ребро между вершинами <code>a<sub>i</sub></code> и <code>b<sub>i</sub></code>.</p><p>В каждой вершине <code>i</code> стоят ворота. Также дан Array чётных целых чисел <code>amount</code>, где <code>amount[i]</code> - это:</p><ul><li>цена, которую нужно заплатить, чтобы открыть ворота в вершине <code>i</code>, если <code>amount[i]</code> отрицательно;</li><li>иначе денежная награда за открытие ворот в вершине <code>i</code>.</li></ul><p>Игра идёт так:</p><ul><li>Сначала Алиса стоит в вершине <code>0</code>, а Боб стоит в вершине <code>bob</code>.</li><li>Каждую секунду Алиса и Боб <strong>оба</strong> переходят в соседнюю вершину. Алиса идёт к какому-то <strong>листу</strong>, а Боб идёт к вершине <code>0</code>.</li><li>В <strong>каждой</strong> вершине на своём пути Алиса и Боб либо платят за открытие ворот в этой вершине, либо получают награду. При этом:<ul><li>Если ворота <strong>уже открыты</strong>, платить не нужно, и награды тоже нет.</li><li>Если Алиса и Боб приходят в вершину <strong>одновременно</strong>, они делят цену или награду за открытие ворот. Другими словами, если цена открытия ворот равна <code>c</code>, то Алиса и Боб платят по <code>c / 2</code>. Так же, если награда равна <code>c</code>, каждый из них получает <code>c / 2</code>.</li></ul></li><li>Когда Алиса приходит в лист, она останавливается. Так же, когда Боб приходит в вершину <code>0</code>, он останавливается. Эти события <strong>не зависят</strong> друг от друга.</li></ul><p>Верните <em><strong>наибольший</strong> чистый доход Алисы, если она идёт к лучшему листу</em>.</p><p>Метод класса <code>Solution</code>: <code>def mostProfitablePath(self, edges: list[list[int]], bob: int, amount: list[int]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> edges = [[0,1],[1,2],[1,3],[3,4]], bob = 3, amount = [-2,4,2,-4,6]
<strong>Выход:</strong> 6
<strong>Пояснение:</strong> игра идёт так:
- Сначала Алиса в вершине 0, Боб в вершине 3. Каждый открывает ворота своей вершины.
  Чистый доход Алисы теперь равен −2.
- Алиса и Боб оба переходят в вершину 1.
  Они пришли одновременно, поэтому открывают ворота вместе и делят награду.
  Доход Алисы становится −2 + (4 / 2) = 0.
- Алиса переходит в вершину 3. Боб уже открыл её ворота, поэтому доход Алисы не меняется.
  Боб переходит в вершину 0 и останавливается.
- Алиса переходит в вершину 4 и открывает там ворота. Её доход становится 0 + 6 = 6.
Теперь ни Алиса, ни Боб не могут идти дальше, и игра заканчивается.
Алиса не может получить больший доход.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> edges = [[0,1]], bob = 1, amount = [-7280,2350]
<strong>Выход:</strong> -7280
<strong>Пояснение:</strong> Алиса идёт по пути 0 -&gt; 1, а Боб по пути 1 -&gt; 0.
Поэтому Алиса открывает ворота только в вершине 0. Её чистый доход равен −7280.</pre><h3>Ограничения</h3><ul><li><code>2 ≤ n ≤ 10<sup>5</sup></code></li><li><code>edges.length == n - 1</code></li><li><code>edges[i].length == 2</code></li><li><code>0 ≤ a<sub>i</sub>, b<sub>i</sub> &lt; n</code></li><li><code>a<sub>i</sub> != b<sub>i</sub></code></li><li><code>edges</code> задаёт правильное Tree.</li><li><code>1 ≤ bob &lt; n</code></li><li><code>amount.length == n</code></li><li><code>amount[i]</code> - <strong>чётное</strong> целое число из отрезка <code>[-10<sup>4</sup>, 10<sup>4</sup>]</code>.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "mostProfitablePath"}, Cases: []sets.TestCase{
						{Input: "[[[0,1],[1,2],[1,3],[3,4]],3,[-2,4,2,-4,6]]", Expected: text("6")},
						{Input: "[[[0,1]],1,[-7280,2350]]", Expected: text("-7280")},
						{Input: "[[[0,1],[0,2],[1,3],[1,4],[2,5],[4,6]],6,[-2,4,2,-4,6,8,10]]", Expected: text("8")},
						{Input: "[[[0,1],[1,2]],2,[2,6,4]]", Expected: text("5")},
						{Input: "[[[0,1]],1,[10,20]]", Expected: text("10")},
						{Input: "[[[0,1],[1,2]],1,[2,4,-10]]", Expected: text("-8")},
					}},
				},
				{
					URL:        "https://leetcode.com/problems/find-minimum-diameter-after-merging-two-trees/",
					Number:     number(3203),
					Title:      "Find Minimum Diameter After Merging Two Trees",
					Difficulty: "hard",
					Topics:     []string{"Tree", "Depth-First Search", "Breadth-First Search", "Graph Theory"},
					Note:       "Соедините центры двух Tree: новый диаметр max(d1, d2, ⌈d1 / 2⌉ + ⌈d2 / 2⌉ + 1), меньше получить нельзя. Каждый диаметр считают двумя BFS. Время: O(n + m). Память: O(n + m) для соседей и расстояний. В первом примере d1 = 2, d2 = 1: max(2, 1, 1 + 1 + 1) = 3. Brute Force по всем n · m рёбрам соединения с BFS от каждой вершины слишком долгий.",
					Statement: `<p>Даны два Tree, рёбра которых <strong>не имеют направления</strong>. В первом <code>n</code> вершин с номерами от <code>0</code> до <code>n - 1</code>, во втором <code>m</code> вершин с номерами от <code>0</code> до <code>m - 1</code>. Даны два двумерных Array целых чисел <code>edges1</code> и <code>edges2</code> длины <code>n - 1</code> и <code>m - 1</code>: запись <code>edges1[i] = [a<sub>i</sub>, b<sub>i</sub>]</code> означает, что в первом Tree есть ребро между вершинами <code>a<sub>i</sub></code> и <code>b<sub>i</sub></code>, а запись <code>edges2[i] = [u<sub>i</sub>, v<sub>i</sub>]</code> означает, что во втором Tree есть ребро между вершинами <code>u<sub>i</sub></code> и <code>v<sub>i</sub></code>.</p><p>Нужно соединить ребром одну вершину первого Tree с одной вершиной второго Tree.</p><p>Верните <strong>наименьший</strong> возможный <strong>диаметр</strong> получившегося Tree.</p><p><strong>Диаметр</strong> Tree - это длина <em>самого длинного</em> пути между двумя вершинами Tree.</p><p>Метод класса <code>Solution</code>: <code>def minimumDiameterAfterMerge(self, edges1: list[list[int]], edges2: list[list[int]]) -&gt; int</code>.</p><h3>Пример 1</h3><pre><strong>Вход:</strong> edges1 = [[0,1],[0,2],[0,3]], edges2 = [[0,1]]
<strong>Выход:</strong> 3
<strong>Пояснение:</strong> Tree с диаметром 3 получится, если соединить вершину 0 первого Tree с любой вершиной второго Tree.</pre><h3>Пример 2</h3><pre><strong>Вход:</strong> edges1 = [[0,1],[0,2],[0,3],[2,4],[2,5],[3,6],[2,7]], edges2 = [[0,1],[0,2],[0,3],[2,4],[2,5],[3,6],[2,7]]
<strong>Выход:</strong> 5
<strong>Пояснение:</strong> Tree с диаметром 5 получится, если соединить вершину 0 первого Tree с вершиной 0 второго Tree.</pre><h3>Ограничения</h3><ul><li><code>1 ≤ n, m ≤ 10<sup>5</sup></code></li><li><code>edges1.length == n - 1</code></li><li><code>edges2.length == m - 1</code></li><li><code>edges1[i].length == edges2[i].length == 2</code></li><li><code>edges1[i] = [a<sub>i</sub>, b<sub>i</sub>]</code></li><li><code>0 ≤ a<sub>i</sub>, b<sub>i</sub> &lt; n</code></li><li><code>edges2[i] = [u<sub>i</sub>, v<sub>i</sub>]</code></li><li><code>0 ≤ u<sub>i</sub>, v<sub>i</sub> &lt; m</code></li><li>Входные данные таковы, что <code>edges1</code> и <code>edges2</code> задают правильные Tree.</li></ul>`,
					Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "call", Method: "minimumDiameterAfterMerge"}, Cases: []sets.TestCase{
						{Input: "[[[0,1],[0,2],[0,3]],[[0,1]]]", Expected: text("3")},
						{Input: "[[[0,1],[0,2],[0,3],[2,4],[2,5],[3,6],[2,7]],[[0,1],[0,2],[0,3],[2,4],[2,5],[3,6],[2,7]]]", Expected: text("5")},
						{Input: "[[],[]]", Expected: text("1")},
						{Input: "[[[0,1],[1,2],[2,3],[3,4]],[]]", Expected: text("4")},
						{Input: "[[[0,1],[0,2],[1,3],[1,4],[2,5],[4,6]],[[0,1],[1,2],[2,3]]]", Expected: text("6")},
					}},
				},
			},
			Cards: []sets.Card{
				{ID: "tr4-path-len", Type: "short", Prompt: "Tree T урока. Сколько рёбер на пути между вершинами 3 и 6?", Answer: "3", Explain: "Путь 3 − 1 − 4 − 6: от 3 вверх к 1, затем вниз к 4 и 6. Ответ 5 идёт через корень 0, но такой путь проходит вершину 1 дважды."},
				{ID: "tr4-first-bfs", Type: "choice", Prompt: "Первый BFS в Tree T начинается от вершины 2. Какую самую дальнюю вершину он найдёт?", Options: []string{"6", "4", "5"}, Answer: 0, Explain: "Расстояния от 2: до 6 четыре ребра (2 − 0 − 1 − 4 − 6), до 3 и 4 три, до 5 одно. Первый BFS найдёт 6, второй BFS от 6 найдёт 5 на расстоянии 5. Вершина 5 тоже конец диаметра, но она соседка вершины 2. Вершина 4 лежит на расстоянии 3."},
				{ID: "tr4-diam", Type: "short", Prompt: "Tree: путь 0 − 1 − 2 − ... − 9 и вершина 10, соединённая с вершиной 5. Чему равен диаметр в рёбрах?", Answer: "9", Explain: "От 10 до 0 шесть рёбер, до 9 пять. Длиннее путь 0 − ... − 9 из 9 рёбер. Ответ 10 добавляет ребро к 10, но путь не может пройти 10 и оба конца сразу."},
				{ID: "tr4-center", Type: "choice", Prompt: "Диаметр Tree равен 6 рёбрам. Сколько центров и какая высота при корне в центре?", Options: []string{"2 центра, высота 3", "1 центр, высота 4", "1 центр, высота 3"}, Answer: 2, Explain: "На пути из 6 рёбер 7 вершин, середина одна: четвёртая вершина. От неё до каждого конца 3 ребра, высота 6 / 2 = 3. Два центра бывают только при нечётной длине диаметра."},
				{ID: "tr4-merge", Type: "short", Prompt: "LeetCode 3203. Диаметры двух Tree равны 4 и 4. Каков наименьший диаметр после соединения одним ребром?", Answer: "5", Explain: "Соединяем центры: ⌈4/2⌉ + ⌈4/2⌉ + 1 = 5, и это больше 4. Ответ 9 = 4 + 4 + 1 соединяет концы диаметров."},
				{ID: "tr4-zelda", Type: "short", Prompt: "Codeforces 1905B. Звезда: центр и 7 листьев. Сколько операций нужно?", Answer: "4", Explain: "Одна операция убирает не больше двух листьев: нужно не меньше ⌈7/2⌉ = 4. Три операции уберут только 6 листьев. Ответ 7 сжимает каждый лист отдельно."},
				{ID: "tr4-latency", Type: "choice", Prompt: "Codeforces 690C2. Рёбра 1 − 2, 1 − 3, 3 − 4, 4 − 5. Один BFS от вершины 1 даёт наибольшее расстояние 3. Чему равна задержка мозга, то есть диаметр?", Options: []string{"3", "4", "5"}, Answer: 1, Explain: "Самая дальняя от 1 вершина 5. Второй BFS от 5 находит вершину 2 на расстоянии 4: путь 5 − 4 − 3 − 1 − 2. Ответ 3 - только высота при корне 1, а вершина 1 не конец диаметра. Ответ 5 считает вершины пути, а не рёбра."},
				{ID: "tr4-brute", Type: "choice", Prompt: "Диаметр ищут BFS от каждой вершины. n = 20 000. Сколько посещений вершин всего?", Options: []string{"около 40 тысяч", "около 400 миллионов", "около 20 тысяч"}, Answer: 1, Explain: "n обходов по n вершин: 20 000 · 20 000 = 4 · 10^8, время O(n²). Два BFS делают 40 тысяч посещений, это и есть первый вариант: O(n)."},
			},
			Links: []sets.Link{
				{Title: "Диаметр и центр дерева", URL: "https://neerc.ifmo.ru/wiki/index.php?title=%D0%90%D0%BB%D0%B3%D0%BE%D1%80%D0%B8%D1%82%D0%BC%D1%8B_%D0%BD%D0%B0_%D0%B4%D0%B5%D1%80%D0%B5%D0%B2%D1%8C%D1%8F%D1%85", Source: "Викиконспекты ИТМО"},
				{Title: "Почему два обхода находят диаметр", URL: "https://habr.com/ru/articles/305912/", Source: "Хабр"},
				{Title: "BFS и расстояния", URL: "https://ru.algorithmica.org/cs/shortest-paths/bfs/", Source: "Алгоритмика"},
				{Title: "Ответ для каждого корня", URL: "https://usaco.guide/gold/all-roots", Source: "USACO Guide (англ.)"},
			},
			Figures: []sets.Figure{
				{
					ID:    "tr4-path",
					Title: "Путь от 3 до 5",
					Frames: []sets.FigureFrame{
						{Caption: "BFS от 3: у каждой вершины родитель в сторону 3 и расстояние до 3. dist[5] = 4.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "parent", Cells: []any{1, 3, 0, -1, 1, 2, 4}},
							{Label: "dist", Cells: []any{2, 1, 3, 0, 2, 4, 3}, Marked: []int{3, 5}},
						}},
						{Caption: "От 5 по parent: 5 → 2 → 0 → 1 → 3. После разворота путь 3, 1, 0, 2, 5 из 4 рёбер.", Rows: []sets.FigureRow{
							{Label: "подъём", Cells: []any{5, 2, 0, 1, 3}},
							{Label: "путь", Cells: []any{3, 1, 0, 2, 5}},
						}},
					},
				},
				{
					ID:    "tr4-diameter",
					Title: "Диаметр двумя BFS",
					Frames: []sets.FigureFrame{
						{Caption: "BFS от 0: самая дальняя вершина 6, расстояние 3.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "от 0", Cells: []any{0, 1, 1, 2, 2, 2, 3}, Marked: []int{6}},
						}},
						{Caption: "BFS от 6: самая дальняя вершина 5, расстояние 5. Это длина диаметра.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "от 6", Cells: []any{3, 2, 4, 3, 1, 5, 0}, Marked: []int{5}},
						}},
						{Caption: "Путь 6, 4, 1, 0, 2, 5 из 5 рёбер. В середине вершины 1 и 0: при корне в любой из них высота 3.", Rows: []sets.FigureRow{
							{Label: "диаметр", Cells: []any{6, 4, 1, 0, 2, 5}, Marked: []int{2, 3}},
						}},
					},
				},
				{
					ID:    "tr4-peel",
					Title: "Срезаем листья слоями",
					Frames: []sets.FigureFrame{
						{Caption: "Листья 3, 5 и 6 имеют степень 1. Удаляем их.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "степень", Cells: []any{2, 3, 2, 1, 2, 1, 1}, Marked: []int{3, 5, 6}},
						}},
						{Caption: "У вершин 1, 2 и 4 степень уменьшилась. Новые листья 2 и 4, удаляем их.", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "степень", Cells: []any{2, 2, 1, 1, 1, 1, 1}, Marked: []int{2, 4}, Dimmed: []int{3, 5, 6}},
						}},
						{Caption: "Остались вершины 0 и 1: это центры. Каждая вершина удалена один раз, время O(n).", Rows: []sets.FigureRow{
							{Label: "вершина", Cells: []any{0, 1, 2, 3, 4, 5, 6}},
							{Label: "степень", Cells: []any{1, 1, 1, 1, 1, 1, 1}, Marked: []int{0, 1}, Dimmed: []int{2, 3, 4, 5, 6}},
						}},
					},
				},
			},
		},
	},
}
