/// <reference types="vite/client" />
import { expect, it } from 'vitest';
import { glossary, splitTerms } from '../shared/glossary';
import content from '../backend/internal/sets/settest/sets.json';

// Shows each term as [word]. The texts are the text nodes of one page and share seen, so a word marked in an earlier text stays plain in a later one.
function marked(texts: string[]) {
	const seen = new Set<number>();
	return texts.map(text => splitTerms(text, seen).map(part => typeof part === 'string' ? part : `[${part.term}]`).join(''));
}

it.each([
	{ name: 'inflected form gets the tip', texts: ['Держите по инварианту.'], expected: ['Держите по [инварианту].'] },
	{ name: 'only the first match of a word', texts: ['Инвариант и снова инвариант.'], expected: ['[Инвариант] и снова инвариант.'] },
	{ name: 'each word has its own first match', texts: ['Кортеж в heap, кортеж в deque.'], expected: ['[Кортеж] в [heap], кортеж в [deque].'] },
	{ name: 'word marked in an earlier text stays plain', texts: ['Кортеж в heap.', 'Кортеж и deque.'], expected: ['[Кортеж] в [heap].', 'Кортеж и [deque].'] },
	{ name: 'word that only overlaps an earlier word stays free for a later text', texts: ['Время O(log n).', 'Память O(1).'], expected: ['Время O([log] n).', 'Память [O(1)].'] },
	{ name: 'capital at the sentence start', texts: ['Кортежи сравниваются по первому элементу.'], expected: ['[Кортежи] сравниваются по первому элементу.'] },
	{ name: 'English name at the sentence start', texts: ['Heap хранит наименьший элемент в h[0].'], expected: ['[Heap] хранит наименьший элемент в h[0].'] },
	{ name: 'heapq is not heap', texts: ['Модуль heapq и heap.'], expected: ['Модуль heapq и [heap].'] },
	{ name: 'dequeue is not deque or queue', texts: ['Слово dequeue, потом deque и queue.'], expected: ['Слово dequeue, потом [deque] и [queue].'] },
	{ name: 'code name is not a term', texts: ['collections.deque, deque([1]), deque.popleft() и queue.Queue.'], expected: ['collections.deque, deque([1]), deque.popleft() и queue.Queue.'] },
	{ name: 'deque at the sentence end', texts: ['pop from an empty deque.'], expected: ['pop from an empty [deque].'] },
	{ name: 'priority queue is one term', texts: ['Это priority queue, а не queue.'], expected: ['Это [priority queue], а не [queue].'] },
	{ name: 'monotonic deque is one term', texts: ['В monotonic deque, а потом в deque.'], expected: ['В [monotonic deque], а потом в [deque].'] },
	{ name: 'canonical hash table and array terms', texts: ['Hash Table хранит значения Array.'], expected: ['[Hash Table] хранит значения [Array].'] },
	{ name: 'difference array before array', texts: ['Difference Array и Array.'], expected: ['[Difference Array] и [Array].'] },
	{ name: 'full canonical heap topic', texts: ['Heap (Priority Queue) и Queue.'], expected: ['[Heap (Priority Queue)] и [Queue].'] },
	{ name: 'specific heap names before generic heap', texts: ['Min-Heap, Max-Heap и Heap.'], expected: ['[Min-Heap], [Max-Heap] и [Heap].'] },
	{ name: 'call stack before stack', texts: ['Call Stack и Stack.'], expected: ['[Call Stack] и [Stack].'] },
	{ name: 'monotonic queue before queue', texts: ['Monotonic Queue и Queue.'], expected: ['[Monotonic Queue] и [Queue].'] },
	{ name: 'algorithm names in Russian prose', texts: ['Сравните Linear Scan и Binary Search.'], expected: ['Сравните [Linear Scan] и [Binary Search].'] },
	{ name: 'two-word bound name before its argument', texts: ['Lower Bound(x) и Upper Bound(3).'], expected: ['[Lower Bound](x) и [Upper Bound](3).'] },
	{ name: 'pointer and window techniques', texts: ['Two Pointers и Sliding Window.'], expected: ['[Two Pointers] и [Sliding Window].'] },
	{ name: 'anagram in Russian prose, not in a method name', texts: ['findAnagrams ищет анаграммы.'], expected: ['findAnagrams ищет [анаграммы].'] },
	{ name: 'possessive algorithm name is one term', texts: ["Kadane's Algorithm находит сумму."], expected: ["[Kadane's Algorithm] находит сумму."] },
	{ name: 'enumeration and recursion terms', texts: ['Backtracking строит Subsets через Recursion.'], expected: ['[Backtracking] строит [Subsets] через [Recursion].'] },
	{ name: 'next permutation before permutations', texts: ['Next Permutation и Permutations.'], expected: ['[Next Permutation] и [Permutations].'] },
	{ name: 'next combination before combinations', texts: ['Next Combination и Combinations.'], expected: ['[Next Combination] и [Combinations].'] },
	{ name: 'merge variants before merge', texts: ['K-Way Merge, Merge Sort и Merge.'], expected: ['[K-Way Merge], [Merge Sort] и [Merge].'] },
	{ name: 'code identifiers stay plain for new terms', texts: ['hash_table, array.array, dict(), set.add, nextPermutation и Array.'], expected: ['hash_table, array.array, dict(), set.add, nextPermutation и [Array].'] },
	{ name: 'dynamic programming and pruning', texts: ['Dynamic Programming и Pruning.'], expected: ['[Dynamic Programming] и [Pruning].'] },
	{ name: 'repeated phrase keeps its full meaning', texts: ['Next Permutation, Next Permutation и Permutations.'], expected: ['[Next Permutation], Next Permutation и [Permutations].'] },
	{ name: 'repeated phrase across nodes keeps the generic term available', texts: ['Min-Heap и Merge Sort.', 'Min-Heap, Merge Sort, Heap и Merge.'], expected: ['[Min-Heap] и [Merge Sort].', 'Min-Heap, Merge Sort, [Heap] и [Merge].'] },
	{ name: 'English aliases exclude calls and dotted identifiers', texts: ['subsets([1]), itertools.permutations(a), obj.recursion, my_backtracking, permutations_count и Subsets.'], expected: ['subsets([1]), itertools.permutations(a), obj.recursion, my_backtracking, permutations_count и [Subsets].'] },
	{ name: 'distinct graph structures have separate tips', texts: ['Graph, Tree и Trie.'], expected: ['[Graph], [Tree] и [Trie].'] },
	{ name: 'binary tree before tree', texts: ['Binary Tree и Tree.'], expected: ['[Binary Tree] и [Tree].'] },
	{ name: 'distinct counting rules have separate tips', texts: ['Product Rule и Sum Rule.'], expected: ['[Product Rule] и [Sum Rule].'] },
	{ name: 'методика is not метод', texts: ['Методика урока.'], expected: ['Методика урока.'] },
	{ name: 'регистрация is not регистр', texts: ['Регистрация на сайте.'], expected: ['Регистрация на сайте.'] },
	{ name: 'регистрах is a form of регистр', texts: ['Буквы в разных регистрах.'], expected: ['Буквы в разных [регистрах].'] },
	{ name: 'регистрами is a form of регистр', texts: ['Буквы с разными регистрами.'], expected: ['Буквы с разными [регистрами].'] },
	{ name: 'линейное время', texts: ['Решите задачу за линейное время.'], expected: ['Решите задачу за [линейное время].'] },
	{ name: 'петель and кратных рёбер are forms of the graph terms', texts: ['Двудольный Graph без петель и кратных рёбер.'], expected: ['[Двудольный] [Graph] без [петель] и [кратных рёбер].'] },
	{ name: 'кратна is not кратное ребро', texts: ['Сумма кратна 3, а кратное ребро повторяет пару.'], expected: ['Сумма кратна 3, а [кратное ребро] повторяет пару.'] },
	{ name: 'связного is a form of связный', texts: ['Форма связного Graph.'], expected: ['Форма [связного] [Graph].'] },
	{ name: 'поддеревьев is a form of поддерево', texts: ['Размеры поддеревьев.'], expected: ['Размеры [поддеревьев].'] },
	{ name: 'диаметра is a form of диаметр', texts: ['Середина диаметра.'], expected: ['Середина [диаметра].'] },
	{ name: 'битовая is a form of бит, 32-битный is its own term', texts: ['Битовая запись, 32-битный тип и биты.'], expected: ['[Битовая] запись, [32-битный] тип и биты.'] },
	{ name: 'срезы is a form of срез, срезать is not', texts: ['Срезы копируют, листья срезать.'], expected: ['[Срезы] копируют, листья срезать.'] },
	{ name: 'переполнение', texts: ['Целые числа не переполняются.'], expected: ['Целые числа не [переполняются].'] },
	{ name: 'log inside O(n log n)', texts: ['Время O(n log n).'], expected: ['Время O(n [log] n).'] },
	{ name: 'log before its base', texts: ['Около n · log2 n шагов.'], expected: ['Около n · [log]2 n шагов.'] },
	{ name: 'O and log in one text', texts: ['За O(1) или за O(log n).'], expected: ['За [O(1)] или за O([log] n).'] },
	{ name: 'O with brackets inside', texts: ['Время O(n log(k)), память O(n log(k)).'], expected: ['Время O(n [log](k)), память [O(n log(k))].'] },
	{ name: 'lower-case o is not O', texts: ['Вызов foo(1) и o(n).'], expected: ['Вызов foo(1) и o(n).'] },
	{ name: 'phrase term', texts: ['Метод класса Solution: def f(self).'], expected: ['[Метод класса] Solution: def f(self).'] },
	{ name: 'word of a phrase term alone', texts: ['Четыре метода работают за O(1).'], expected: ['Четыре [метода] работают за [O(1)].'] },
	{ name: 'phrase across a line break', texts: ['Это перебор\nс возвратом.'], expected: ['Это [перебор\nс возвратом].'] },
	{ name: 'ambiguous корень stays plain', texts: ['Корень это пустой список.'], expected: ['Корень это пустой список.'] },
	{ name: 'ambiguous место stays plain', texts: ['Число стоит на месте 2.'], expected: ['Число стоит на месте 2.'] },
	{ name: 'ambiguous сочетание stays plain', texts: ['Любое сочетание двух разных букв.'], expected: ['Любое сочетание двух разных букв.'] },
	{ name: 'класс остатков in an inflected form', texts: ['Числа одного класса остатков.'], expected: ['Числа одного [класса остатков].'] },
	{ name: 'класс alone stays plain', texts: ['Метод класса Solution и класс 0.'], expected: ['[Метод класса] Solution и класс 0.'] },
	{ name: 'обратный элемент in an inflected form', texts: ['Делите через обратный элемент.'], expected: ['Делите через [обратный элемент].'] },
	{ name: 'малая теорема Ферма in an inflected form', texts: ['По малой теореме Ферма остаток равен 1.'], expected: ['По [малой теореме Ферма] остаток равен 1.'] },
	{ name: 'малая теорема without the name', texts: ['Малая теорема верна для p = 7.'], expected: ['[Малая теорема] верна для p = 7.'] },
	{ name: 'математическое ожидание in an inflected form', texts: ['Найдите математическое ожидание произведения.'], expected: ['Найдите [математическое ожидание] произведения.'] },
	{ name: 'combinatorics models are one term each', texts: ["Pascal's Triangle и Stars and Bars."], expected: ["[Pascal's Triangle] и [Stars and Bars]."] },
	{ name: 'text without hard words', texts: ['Два указателя идут навстречу.'], expected: ['Два указателя идут навстречу.'] },
	{ name: 'empty text', texts: [''], expected: [''] },
])('$name', ({ texts, expected }) => {
	expect(marked(texts)).toEqual(expected);
});

// The learning sets are Go code. The file sets.json has their text for the checks of this file.
const sets = content as unknown as { slug: string; summary: string; chapters: { summary: string; lesson: string; tasks: { title: string; statement?: string; note: string }[]; cards: { id: string; prompt: string; explain: string }[] }[] }[];

// The hard words that need a tip in shared/glossary.ts. The list also has words that no set uses yet, so a new set with such a word fails until the word gets a tip.
// An English name of a data structure is a word of its own: heapq, dequeue, collections.deque, deque(...) and deque.popleft are code.
const hardWord = /(?<!\p{L})(?:инвариант|O\(|log|метод|(?<![\d_.])(?:(?:heap|deque|queue|stack|array|string|subset|permutation|combination|anagram|bitmask)s?|kadane|monotonic|priority|hash\s+table|binary\s+tree|two\s+pointers|linear\s+scan|binary\s+search|lower\s+bound|upper\s+bound|ternary\s+search|modular\s+arithmetic|binary\s+exponentiation|pascal's\s+triangle|stars\s+and\s+bars|sliding\s+window|dynamic\s+programming|brute\s+force|prefix\s+product|difference\s+array|matri(?:x|ces)|backtracking|recursion|pruning|factorial|sorting|merge|greedy|simulation|quickselect|next\s+combination|pyodide|IndexError|KeyError|Counter|DFS|BFS)(?![\p{L}\d_(]|\.\p{L})|монотонн|невозраст|медиан|жадн|моделир|кортеж|рекурс|перебор\p{L}*\s+с\s+возврат|отсеч|трассир|подмассив|подпоследовательн|анаграмм|лексикограф|ленив|индекс|срез(?:а|у|ом|е|ы|ов|ам|ами|ах)?(?!\p{L})|факториал|перестановк|подмножеств|скользящ|ASCII|регистр|слияни|32-бит|64-бит|бит(?:ов\p{L}*|а|у|ом|е|ы|ам|ами|ах)?(?!\p{L})|двоичн|слагаем|префикс|суффикс|включени\p{L}*\s+и\s+исключени|линейн|переполн|амортиз|асимптот|ожидаем|мемоиз|итератор|компаратор|эвристик|коллизи|логарифм|треугольн\p{L}*\s+чис|повёрнут|граф|бинарн|класс\p{L}*\s+остатк|обратн\p{L}*\s+элемент|мал\p{L}*\s+теорем|математическ\p{L}*\s+ожидани|связн|двудольн|пет(?:л|ел)|кратн\p{L}*\s+р[её]б|поддерев|диаметр)/giu;
// The hard words of the text that no match of the patterns covers.
const untipped = (text: string, patterns: RegExp[]) => [...text.matchAll(hardWord)].filter(word => !patterns.some(pattern => [...text.matchAll(pattern)].some(match => match.index <= word.index && word.index < match.index + match[0].length))).map(word => word[0]);
// Code is not prose. The text of the other tags stays.
const prose = (html: string) => html.replace(/<(pre|code)\b[\s\S]*?<\/\1>/g, ' ').replace(/<[^>]+>/g, '');
it.each(sets.map(item => ({ name: item.slug, item })))('$name hard words have a tip', ({ item }) => {
	const patterns = glossary.map(entry => entry.pattern);
	const texts = [['summary', item.summary], ...item.chapters.flatMap((chapter, index) => [
		[`chapter ${index + 1} summary`, chapter.summary], [`chapter ${index + 1} lesson`, prose(chapter.lesson)],
		...chapter.tasks.flatMap(task => [[`${task.title} statement`, prose(task.statement ?? '')], [`${task.title} note`, task.note]]),
		...chapter.cards.flatMap(card => [[`${card.id} prompt`, card.prompt], [`${card.id} explain`, card.explain]]),
	])];
	expect(texts.flatMap(([place, text]) => untipped(text, patterns).map(word => `${place}: ${word}`))).toEqual([]);
});

// The Editorials are repository files in backend/internal/editorials.
const editorials = import.meta.glob<string>('../backend/internal/editorials/*.html', { query: '?raw', import: 'default', eager: true });
// The title holds the whole slug: vitest cuts a long $name.
it.each(Object.entries(editorials).map(([file, html]) => [`${file.replace(/^.*\/|\.html$/g, '')} hard words have a tip`, html]))('%s', (_, html) => {
	expect(untipped(prose(html), glossary.map(entry => entry.pattern))).toEqual([]);
});

it('each tip is plain', () => {
	expect(glossary.filter(entry => /[–—]/.test(entry.tip) || untipped(entry.tip, [entry.pattern]).length).map(entry => entry.tip)).toEqual([]);
});
