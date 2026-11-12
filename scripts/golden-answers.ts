// One-time generator of shared/answers.cases.json: the cases that the web test and the Go test of the Python values both read.
// Run from the repository root: bun scripts/golden-answers.ts
// A case with an expected value in this file is a documented result of Python ast.literal_eval and repr. The run checks it against shared/answers.ts.
// A case without one gets its result from shared/answers.ts.
import { writeFileSync } from 'node:fs';
import path from 'node:path';
import { parsePy, reprPy, same } from '../shared/answers';
import type { TestSpec } from '../shared/leetcode';

const call: TestSpec = { mode: 'call', method: 'solve' }, any: TestSpec = { ...call, order: 'any' }, nested: TestSpec = { ...call, order: 'any-nested' };
const stdio: TestSpec = { mode: 'stdio' }, noCase: TestSpec = { mode: 'stdio', ignoreCase: true };

const parse = [
	{ name: 'arguments with a nested list', text: '[[2, 7, 11, 15], 9]', value: [[2, 7, 11, 15], 9] },
	{ name: 'negative numbers', text: '[-1, -100, 0]', value: [-1, -100, 0] },
	{ name: 'floats', text: '[1.5, -0.25, .5, 2e3]', value: [1.5, -0.25, 0.5, 2000] },
	{ name: 'True, False and None', text: '[True, False, None]', value: [true, false, null] },
	{ name: 'double quote in single quotes', text: `['a"b']`, value: ['a"b'] },
	{ name: 'escaped quotes', text: `["a\\"b", 'c\\'d']`, value: ['a"b', "c'd"] },
	{ name: 'escaped backslash, line break and tab', text: `'a\\\\b\\nc\\td'`, value: 'a\\b\nc\td' },
	{ name: 'Cyrillic string', text: '"привет"', value: 'привет' },
	{ name: 'tuples as lists', text: '((1, 2), (3,), ())', value: [[1, 2], [3], []] },
	{ name: 'parentheses without a comma', text: '(5)', value: 5 },
	{ name: 'empty containers', text: '[[], {}, ""]', value: [[], {}, ''] },
	{ name: 'dict with string keys', text: `{'a': [1], "b": None}`, value: { a: [1], b: null } },
	{ name: 'trailing commas and line breaks', text: '[\n  [1, 2,],\n  3,\n]\n', value: [[1, 2], 3] },
];

// The text is parsed and printed again. The printed text shows the number format, the order of dict keys and the quotes.
const roundTrip = [
	{ name: 'zero', text: '0' },
	{ name: 'negative zero', text: '-0' },
	{ name: 'number with leading zeros', text: '007' },
	{ name: 'number with a point and no fraction', text: '5.' },
	{ name: 'negative fraction without a zero', text: '-.5' },
	{ name: 'float with a zero fraction', text: '2.0' },
	{ name: 'float with trailing zeros', text: '2.50' },
	{ name: 'one tenth', text: '0.1' },
	{ name: 'exponent with an upper-case mark', text: '1E5' },
	{ name: 'exponent with a plus sign', text: '1.5e+3' },
	{ name: 'the smallest number in plain form', text: '0.000001' },
	{ name: 'the largest number in exponent form below one', text: '1e-7' },
	{ name: 'small number with a fraction in exponent form', text: '1.5e-10' },
	{ name: 'the largest number in plain form', text: '123456789012345678901' },
	{ name: 'the smallest number in exponent form above one', text: '1e21' },
	{ name: 'integer above the exact range of a float', text: '9007199254740993' },
	{ name: 'number above the float range', text: '1e400' },
	{ name: 'negative number above the float range', text: '-1e400' },
	{ name: 'number below the float precision', text: '1e-400' },
	{ name: 'long fraction', text: '3.141592653589793238' },
	{ name: 'dict keeps the order of its keys', text: `{'b': 1, 'a': 2}` },
	{ name: 'dict with keys that are whole numbers', text: `{'b': 0, '2': 'x', '1': 'y', '10': 1, '01': 2, '-1': 3, '1.5': 4}` },
	{ name: 'dict with a repeated key', text: `{'a': 1, 'b': 2, 'a': 3}` },
	{ name: 'dict with a repeated number key', text: `{'b': 1, '7': 2, 'b': 3, '7': 4}` },
	{ name: 'dict with an empty key', text: `{'': 1}` },
	{ name: 'dict in a dict', text: `{'a': {'c': [], 'b': ()}}` },
	{ name: 'dict key with the largest list index', text: `{'x': 1, '4294967294': 2, '4294967295': 3}` },
	{ name: 'no-break space between tokens', text: '[1, 2]' },
	{ name: 'vertical tab and form feed between tokens', text: '[1,\v\f2]' },
	{ name: 'byte order mark before the value', text: '﻿[1]' },
	{ name: 'line separator between tokens', text: '[1, 2]' },
	{ name: 'string with a tab, a return and a line break', text: `'a\\tb\\rc\\nd'` },
	{ name: 'string with only a double quote', text: `'"'` },
	{ name: 'string with only a single quote', text: `"'"` },
	{ name: 'string with both quotes', text: `'it\\'s "x"'` },
	{ name: 'string with a backslash', text: `'a\\\\b'` },
	{ name: 'string with a raw tab', text: `'a\tb'` },
	{ name: 'string with a character outside the basic plane', text: `'😀'` },
	{ name: 'string with a no-break space', text: `'a b'` },
	{ name: 'string with a delete character', text: `'a\x7fb'` },
	{ name: 'string with a raw control character', text: `'a\x01b'` },
	{ name: 'empty string in double quotes', text: '""' },
	{ name: 'nested parentheses without a comma', text: '((1))' },
	{ name: 'tuple of one tuple', text: '((1,),)' },
	{ name: 'empty tuple', text: '()' },
	{ name: 'list with a tuple and a number in parentheses', text: '[(1, 2), (3)]' },
	{ name: 'tuple with a trailing comma', text: '(1, 2,)' },
	{ name: 'names next to marks', text: '[True,False,None]' },
	{ name: 'white space around the value', text: ' \n\t[ 1 , 2 ] \n' },
];

const parseErrors = [
	{ name: 'empty text', text: ' ' },
	{ name: 'unclosed list', text: '[1, 2' },
	{ name: 'two commas', text: '[1,, 2]' },
	{ name: 'text after the value', text: '[1] 2' },
	{ name: 'JSON literal', text: '[true]' },
	{ name: 'dict with a number key', text: '{1: 2}' },
	{ name: 'set', text: '{1, 2}' },
	{ name: 'unclosed string', text: "'abc" },
	{ name: 'unsupported escape', text: "'\\x41'" },
	{ name: 'expression', text: '1 + 2' },
	{ name: 'no text', text: '' },
	{ name: 'items without a comma', text: '[1 2]' },
	{ name: 'dict without a colon', text: `{'a' 1}` },
	{ name: 'dict without a value', text: `{'a': }` },
	{ name: 'dict with a list key', text: `{[1]: 2}` },
	{ name: 'dict with two commas', text: `{'a': 1,, }` },
	{ name: 'minus sign apart from the number', text: '- 5' },
	{ name: 'plus sign before a number', text: '+5' },
	{ name: 'two numbers', text: '1 2' },
	{ name: 'name with a digit after it', text: 'True1' },
	{ name: 'name with a letter after it', text: 'Nonee' },
	{ name: 'name in lower case', text: 'none' },
	{ name: 'comma after the value', text: 'None,' },
	{ name: 'close mark without an open mark', text: '[1,]]' },
	{ name: 'list closed as a tuple', text: '[1)' },
	{ name: 'comma in an empty list', text: '[,]' },
	{ name: 'comma in an empty tuple', text: '(,)' },
	{ name: 'line break in a string', text: `'a\nb'` },
	{ name: 'string that ends with a backslash', text: `"\\"` },
	{ name: 'escape that Python has and the parser does not', text: `'\\a'` },
	{ name: 'escaped line break', text: `'a\\\nb'` },
	{ name: 'next-line character as white space', text: '[1,\u00852]' },
	{ name: 'hexadecimal number', text: '0x10' },
	{ name: 'number with an underscore', text: '1_000' },
	{ name: 'exponent without digits', text: '1e' },
	{ name: 'a point alone', text: '.' },
	{ name: 'bytes literal', text: `b'a'` },
];

const print = [
	{ name: 'nested list', value: [[1, 2], [3]], text: '[[1, 2], [3]]' },
	{ name: 'negative number', value: -5, text: '-5' },
	{ name: 'float', value: 0.5, text: '0.5' },
	{ name: 'True, False and None', value: [true, false, null], text: '[True, False, None]' },
	{ name: 'string', value: 'abc', text: "'abc'" },
	{ name: 'string with a single quote', value: "it's", text: `"it's"` },
	{ name: 'string with both quotes', value: `it's "x"`, text: `'it\\'s "x"'` },
	{ name: 'string with a backslash and a line break', value: 'a\\b\nc', text: "'a\\\\b\\nc'" },
	{ name: 'string with a tab and a control character', value: '\t\x01', text: "'\\t\\x01'" },
	{ name: 'Cyrillic string', value: 'да', text: "'да'" },
	{ name: 'empty containers', value: [[], {}, ''], text: "[[], {}, '']" },
	{ name: 'dict', value: { a: 1, b: [true] }, text: "{'a': 1, 'b': [True]}" },
];

interface Comparison { name: string; spec: TestSpec; expected: string; output: string; equal?: boolean }
const compare: Comparison[] = [
	{ name: 'same list in another format', spec: call, expected: '[1, 2]', output: '[1,2]', equal: true },
	{ name: 'exact order with another order', spec: call, expected: '[1, 2]', output: '[2, 1]', equal: false },
	{ name: 'tuple for a list', spec: call, expected: '[1, 2]', output: '(1, 2)', equal: true },
	{ name: 'float with a zero fraction for an int', spec: call, expected: '2', output: '2.0', equal: true },
	{ name: 'string for a boolean', spec: call, expected: 'True', output: "'True'", equal: false },
	{ name: 'number for a boolean', spec: call, expected: 'True', output: '1', equal: false },
	{ name: 'number for a boolean in a dictionary', spec: call, expected: "{'a': True}", output: "{'a': 1}", equal: false },
	{ name: 'any order with another outer order', spec: any, expected: '[[1, 2], [3]]', output: '[[3], [1, 2]]', equal: true },
	{ name: 'any order with another inner order', spec: any, expected: '[[1, 2], [3]]', output: '[[3], [2, 1]]', equal: false },
	{ name: 'any order with other duplicates', spec: any, expected: '[1, 1, 2]', output: '[1, 2, 2]', equal: false },
	{ name: 'any order with None', spec: any, expected: '[1]', output: 'None', equal: false },
	{ name: 'any order with mixed output types', spec: any, expected: '[0, 1]', output: '[0, None]', equal: false },
	{ name: 'any order with a set', spec: any, expected: '[1, 2]', output: '{1, 2}', equal: false },
	{ name: 'any-nested order with other inner and outer orders', spec: nested, expected: '[[1, 2], [3]]', output: '[[3], [2, 1]]', equal: true },
	{ name: 'any-nested order with a missing list', spec: nested, expected: '[[1, 2], [3]]', output: '[[1, 2]]', equal: false },
	{ name: 'output that is not a Python literal', spec: call, expected: '[1]', output: '<object at 0x1>', equal: false },
	{ name: 'stdio with other white space', spec: stdio, expected: '1 2\n3\n', output: '1  2 3', equal: true },
	{ name: 'stdio with joined tokens', spec: stdio, expected: '1 2\n', output: '12\n', equal: false },
	{ name: 'stdio keeps case by default', spec: stdio, expected: 'YES\n', output: 'yes\n', equal: false },
	{ name: 'stdio ignores letter case when allowed', spec: noCase, expected: 'YES\n', output: 'yEs\n', equal: true },
	{ name: 'stdio ignores case and white space', spec: noCase, expected: 'YES NO\n', output: 'yes\nNo', equal: true },
	{ name: 'stdio with another answer when case is ignored', spec: noCase, expected: 'YES\n', output: 'no\n', equal: false },
	{ name: 'stdio with split tokens when case is ignored', spec: noCase, expected: 'YES\n', output: 'y es\n', equal: false },
	{ name: 'stdio with extra tokens when case is ignored', spec: noCase, expected: 'YES\n', output: 'yes no\n', equal: false },
	{ name: 'exponent form for a plain number', spec: call, expected: '1000', output: '1e3' },
	{ name: 'float that is not exact', spec: call, expected: '0.3', output: '0.30000000000000004' },
	{ name: 'dict with another key order', spec: call, expected: "{'a': 1, 'b': 2}", output: "{'b': 2, 'a': 1}" },
	{ name: 'dict with number keys in another order', spec: call, expected: "{'1': 1, '2': 2}", output: "{'2': 2, '1': 1}" },
	{ name: 'string in other quotes', spec: call, expected: "'ab'", output: '"ab"' },
	{ name: 'string with another letter case', spec: call, expected: "'ab'", output: "'AB'" },
	{ name: 'expected answer that is not a Python literal', spec: call, expected: '<object>', output: '<object>' },
	{ name: 'any order with strings', spec: any, expected: "['b', 'a', 'Б', 'А']", output: "['А', 'Б', 'a', 'b']" },
	{ name: 'any order with a number for a string', spec: any, expected: "['1']", output: '[1]' },
	{ name: 'any order with an empty list', spec: any, expected: '[]', output: '()' },
	{ name: 'any order with a dict', spec: any, expected: '[1]', output: "{'a': 1}" },
	{ name: 'any-nested order with an item that is not a list', spec: nested, expected: '[[1], [2]]', output: '[[1], 2]' },
	{ name: 'any-nested order with tuples', spec: nested, expected: '[[1, 2], [3, 4]]', output: '((4, 3), (2, 1))' },
	{ name: 'any-nested order with other inner duplicates', spec: nested, expected: '[[1, 1, 2]]', output: '[[1, 2, 2]]' },
	{ name: 'stdio with empty outputs', spec: stdio, expected: '', output: ' \n' },
	{ name: 'stdio with a no-break space between tokens', spec: stdio, expected: '1 2', output: '1 2' },
	{ name: 'stdio with a byte order mark before the tokens', spec: stdio, expected: '1 2', output: '﻿1 2' },
	{ name: 'stdio with a next-line character between tokens', spec: stdio, expected: '1 2', output: '1\u00852' },
	{ name: 'stdio with a vertical tab between tokens', spec: stdio, expected: '1 2', output: '1\v2' },
	{ name: 'stdio ignores the case of Cyrillic letters', spec: noCase, expected: 'ДА', output: 'да' },
	{ name: 'stdio ignores the case of a capital I with a dot', spec: noCase, expected: 'İ', output: 'i̇' },
	{ name: 'stdio compares a capital I with a dot and a small i', spec: noCase, expected: 'İ', output: 'i' },
	{ name: 'stdio ignores the case of a final sigma', spec: noCase, expected: 'ΟΔΟΣ', output: 'οδος' },
	{ name: 'stdio compares a final sigma and a middle sigma', spec: noCase, expected: 'ΟΔΟΣ', output: 'οδοσ' },
	{ name: 'stdio ignores the case of a sharp s', spec: noCase, expected: 'STRASSE', output: 'straße' },
	{ name: 'stdio compares numbers as text', spec: stdio, expected: '1.0', output: '1' },
];

function failure(text: string): boolean {
	try { parsePy(text); return false; } catch (error) { return error instanceof SyntaxError; }
}
function check(condition: boolean, message: string): void {
	if (!condition) throw new Error(message);
}

for (const { name, text, value } of parse) check(JSON.stringify(parsePy(text)) === JSON.stringify(value), `parse: ${name}`);
for (const { name, text } of parseErrors) check(failure(text), `parseErrors: ${name}`);
for (const { name, value, text } of print) check(reprPy(value) === text, `print: ${name}`);
for (const { name, spec, expected, output, equal } of compare) check(equal === undefined || same(spec, expected, output) === equal, `compare: ${name}`);
const names = [parse, roundTrip, parseErrors, print, compare].map(list => list.map(item => item.name));
for (const list of names) check(new Set(list).size === list.length, `repeated name in ${list[0]}`);

const cases = {
	parse,
	roundTrip: roundTrip.map(({ name, text }) => ({ name, text, repr: reprPy(parsePy(text)) })),
	parseErrors,
	print,
	compare: compare.map(({ name, spec, expected, output }) => ({ name, spec, expected, output, equal: same(spec, expected, output) })),
};
writeFileSync(path.join(import.meta.dirname, '../shared/answers.cases.json'), `${JSON.stringify(cases, null, '\t')}\n`);
console.log(`Wrote ${Object.values(cases).reduce((count, list) => count + list.length, 0)} cases.`);
