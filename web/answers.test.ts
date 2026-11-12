import { expect, it } from 'vitest';
import { parsePy, reprPy, same } from '../shared/answers';
import type { TestSpec } from '../shared/leetcode';
import cases from '../shared/answers.cases.json';

// The Go copy of shared/answers.ts reads the same cases (backend/internal/answers), so the two copies stay equal.
// scripts/golden-answers.ts wrote the file. The values of the first list are the results of Python ast.literal_eval and repr.
it.each(cases.parse)('parses $name', ({ text, value }) => {
	expect(parsePy(text)).toEqual(value);
});

it.each(cases.roundTrip)('parses and prints $name', ({ text, repr }) => {
	expect(reprPy(parsePy(text))).toBe(repr);
});

it.each(cases.parseErrors)('does not parse $name', ({ text }) => {
	expect(() => parsePy(text)).toThrow(SyntaxError);
});

it.each(cases.print)('prints $name', ({ value, text }) => {
	expect(reprPy(value)).toBe(text);
});

it('does not print a value without a Python literal', () => {
	expect(() => reprPy(undefined)).toThrow(TypeError);
});

it.each(cases.compare)('compares $name', ({ spec, expected, output, equal }) => {
	expect(same(spec as TestSpec, expected, output)).toBe(equal);
});
