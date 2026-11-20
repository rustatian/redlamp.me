import type { Client } from 'pg';
import { test, expect, signIn } from './fixtures';

// cf-678-a of backend/internal/bank/problems-001.json: the smallest x above n that k divides. Its statement has 3 examples.
const solution = 'n, k = map(int, input().split())\nprint((n // k + 1) * k)\n';
const minute = 60_000, iso = (time: number) => new Date(time).toISOString();

// A contest of 90 minutes with cf-678-a at position 1. With solvedAfter, an accepted submission comes that many milliseconds after the start.
async function addContest(db: Client, startedAt: number, solvedAfter?: number) {
	const id = crypto.randomUUID();
	await db.query('INSERT INTO contests (id,started_at,ends_at) VALUES ($1,$2,$3)', [id, iso(startedAt), iso(startedAt + 90 * minute)]);
	await db.query("INSERT INTO contest_problems (contest_id,position,slug,title,rating,topic,difficulty,points) VALUES ($1,1,'cf-678-a','Johny Likes Numbers',800,'Implementation','easy',3)", [id]);
	if (solvedAfter !== undefined) await db.query("INSERT INTO contest_submissions (id,contest_id,position,verdict,code,created_at) VALUES ($1,$2,1,'accepted',$3,$4)", [crypto.randomUUID(), id, solution, iso(startedAt + solvedAfter)]);
	return id;
}

// The rating starts at 800. One solved 800 problem: E = 0.5, so the change is round(32 × 0.5) = +16. Easy points: 3. Finish: the accepted submission 10 minutes after the start.
const finishedRow = [/Лучший$/, '1 из 1', '3', '10:00', '+16', '816'];

test('the skull window closes with Нет and the setup starts a contest', async ({ page }) => {
	await signIn(page);
	await page.goto('/leetcode/contests');
	const start = page.getByRole('button', { name: 'В мясорубку', exact: true });
	await start.click();
	await expect(page.getByRole('dialog', { name: 'Готов умереть?' })).toBeVisible();
	await page.getByRole('button', { name: 'Нет', exact: true }).click();
	await expect(page.getByRole('dialog')).toHaveCount(0);
	await start.click();
	await page.getByRole('button', { name: 'Да', exact: true }).click();
	const dialog = page.getByRole('dialog', { name: 'Выбор задач' }), topics = dialog.getByRole('group', { name: 'Темы' }), topic = topics.getByRole('button', { name: /^Implementation x\d$/ });
	await topic.click();
	await expect(topic).toHaveText('Implementation x1');
	await expect(dialog.getByRole('listitem')).toHaveCount(1);
	await dialog.getByRole('button', { name: 'Убрать задачу 1: Implementation', exact: true }).click();
	await expect(topic).toHaveText('Implementation x0');
	await expect(dialog.getByRole('listitem')).toHaveCount(0);
	// Implementation has problems of each difficulty, so 4 presses fill the 4 places: 1 easy, 2 medium and 1 hard.
	for (let i = 0; i < 4; i++) await topic.click();
	await expect(topic).toHaveText('Implementation x4');
	await expect(topics.getByRole('button', { disabled: false })).toHaveCount(0);
	await dialog.getByRole('button', { name: 'Начать', exact: true }).click();
	await expect(page).toHaveURL(/\/leetcode\/contests\/[a-f0-9-]{36}$/);
	await expect(page.getByRole('timer')).toHaveText(/^(89|90):\d\d$/);
});

test('a wrong submission shows only the test number, a right one is accepted', async ({ page, db }) => {
	test.slow();
	await signIn(page);
	const id = await addContest(db, Date.now());
	await page.goto(`/leetcode/contests/${id}/1`);
	const code = page.getByLabel('Код решения', { exact: true }), send = page.getByRole('button', { name: 'Отправить', exact: true });
	// Each answer is at least k, and k is at least 1, so 0 is wrong on the first hidden test.
	await code.fill('print(0)\n');
	await send.click();
	await expect(page.getByRole('alert', { name: 'Результат проверки' })).toHaveText('Неверный ответ на тесте 1', { timeout: 60_000 });
	await expect(page.getByText(/Верный ответ/)).toHaveCount(0);
	await code.fill(solution);
	await send.click();
	await expect(page.getByText('Принято', { exact: true })).toBeVisible({ timeout: 60_000 });
	await expect(send).toHaveCount(0);
});

test('a reload keeps the code', async ({ page, db }) => {
	await signIn(page);
	const id = await addContest(db, Date.now());
	await page.goto(`/leetcode/contests/${id}/1`);
	const code = page.getByLabel('Код решения', { exact: true }), status = page.getByRole('status', { name: 'Состояние сохранения' });
	await code.fill(solution);
	// The status is Сохранено also before the change, so the check waits for the change first.
	await expect(status).toHaveText('Не сохранено');
	await expect(status).toHaveText('Сохранено');
	await page.reload();
	await expect(code.locator('.cm-line')).toHaveText(['n, k = map(int, input().split())', 'print((n // k + 1) * k)', '']);
});

test('the result and the rating show after the end', async ({ page, db }) => {
	await signIn(page);
	await addContest(db, Date.now() - 120 * minute, 10 * minute);
	await page.goto('/leetcode/contests');
	await expect(page.locator('.stat')).toHaveText(['816 Рейтинг', '816 Лучший рейтинг', '3 Очки лучшего контеста, финиш 10:00', '1 Контестов'], { useInnerText: true });
	await expect(page.getByRole('row').nth(1).getByRole('cell')).toHaveText(finishedRow);
});

test('the parent sees the results without the start button', async ({ page, db }) => {
	await signIn(page, 'parent');
	await addContest(db, Date.now() - 120 * minute, 10 * minute);
	await page.goto('/leetcode/contests');
	await expect(page.getByRole('row').nth(1).getByRole('cell')).toHaveText(finishedRow);
	await expect(page.getByRole('button', { name: 'В мясорубку', exact: true })).toHaveCount(0);
});

test('the problem page turns read-only at the end without a reload', async ({ page, db }) => {
	test.slow();
	await signIn(page);
	const id = await addContest(db, Date.now() + 10_000 - 90 * minute);
	await page.goto(`/leetcode/contests/${id}/1`);
	const code = page.getByLabel('Код решения', { exact: true }), send = page.getByRole('button', { name: 'Отправить', exact: true });
	await page.getByLabel('Ввод', { exact: true }).fill('10 4');
	await page.getByRole('button', { name: 'Добавить тест', exact: true }).click();
	await code.fill(solution);
	await expect(send).toBeVisible();
	await expect(code).not.toBeEditable({ timeout: 15_000 });
	await expect(send).toHaveCount(0);
	// The clock hides at its end. A remount of the page at that time drops the own tests, because they stay only in the page.
	await expect(page.getByRole('timer')).toHaveCount(0);
	await expect(page.getByRole('list', { name: 'Мои тесты' }).getByRole('listitem')).toHaveText(['Тест 4 Ввод: 10 4 Удалить'], { useInnerText: true });
});

test('the contest page loads the result at the end and the clock hides', async ({ page, db }) => {
	test.slow();
	await signIn(page);
	const id = await addContest(db, Date.now() + 10_000 - 90 * minute);
	await page.goto(`/leetcode/contests/${id}`);
	await expect(page.getByRole('timer')).toBeVisible();
	await expect(page.locator('.stat')).toHaveText(['0 Очки'], { useInnerText: true });
	// No solved problem: the finish is 0 and the change is round(32 × (0 - 0.5)) = -16.
	await expect(page.locator('.stat')).toHaveText(['0 Очки', '00:00 Финиш', '-16 Изменение рейтинга'], { useInnerText: true, timeout: 15_000 });
	await expect(page.getByText('Контест закончился.', { exact: true })).toBeVisible();
	await expect(page.getByRole('timer')).toHaveCount(0);
});

test('own tests of a contest problem show only the output', async ({ page, db }) => {
	test.slow();
	await signIn(page);
	const id = await addContest(db, Date.now());
	await page.goto(`/leetcode/contests/${id}/1`);
	await page.getByLabel('Ввод', { exact: true }).fill('10 4');
	await page.getByRole('button', { name: 'Добавить тест', exact: true }).click();
	await page.getByLabel('Код решения', { exact: true }).fill(solution);
	await page.getByRole('button', { name: 'Запустить тесты', exact: true }).click();
	// The 3 examples have an answer. The own test has none, so the count leaves it out.
	await expect(page.getByText(/^Пройдено/)).toHaveText('Пройдено 3 из 3', { timeout: 30_000 });
	await expect(page.getByRole('list', { name: 'Результаты тестов' }).getByRole('listitem', { name: 'Тест 4' })).toHaveText('Тест 4: без проверки Ввод: 10 4 Ваш ответ: 12', { useInnerText: true });
});
