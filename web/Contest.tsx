import { useEffect, useRef, useState } from 'react';
import type { ContestDetail, ContestProblem, ContestVerdict, OwnTest, ProblemTests, SessionUser, TestSpec } from '../shared/leetcode';
import { api, ApiError, message } from './api';
import { hiddenOutputs } from './acceptance';
import { CodeField } from './Code';
import { points } from './Contests';
import { clock } from './contestTime';
import { difficultyLabels } from './Problems';
import { failures, Runner } from './Runner';
import { Lesson } from './Sets';
import { useResource } from './useResource';
import { Panes } from './Workspace';

// The verdict names only the failed test: the server sends no hidden input and no hidden answer.
const verdictText = (item: ContestVerdict) => item.verdict === 'accepted' ? 'Принято' : `${failures[item.verdict]} ${item.test}${item.error ? `: ${item.error}` : ''}`;

// The topbar clock shows the time left. At the end it loads this page again, so the page shows the result.
export function ContestPage({ id }: { id: string }) {
	const { data, error } = useResource<ContestDetail>(`/contests/${id}`);
	if (!data) return <div className="stack"><h1>Контест</h1><p role="status">{error || 'Открываем контест...'}</p></div>;
	const stats = [{ label: 'Очки', value: data.score }, ...(data.active ? [] : [{ label: 'Финиш', value: clock(data.finishSeconds) }, { label: 'Изменение рейтинга', value: data.ratingChange!.toLocaleString('ru-RU', { signDisplay: 'exceptZero' }) }])];
	return <div className="stack"><div className="heading"><div><p className="small"><a href="/leetcode/contests">Мясорубка</a></p><h1>Контест {new Date(data.startedAt).toLocaleString('ru-RU', { dateStyle: 'short', timeStyle: 'short' })}</h1><p>{data.active ? 'Контест идёт.' : 'Контест закончился.'}</p></div></div>
		<div className="stats">{stats.map(item => <div className="stat" key={item.label}><strong>{item.value}</strong><span>{item.label}</span></div>)}</div>
		<section className="card"><header className="card-head"><h2>Задачи</h2></header><div className="results"><table><thead><tr><th scope="col">№</th><th scope="col">Задача</th><th scope="col">Сложность</th><th scope="col">Очки</th><th scope="col">Состояние</th><th scope="col">Неверные отправки</th></tr></thead>
			<tbody>{data.problems.map(p => <tr key={p.position}><td>{p.position}</td><td><a href={`/leetcode/contests/${id}/${p.position}`}>{p.title}</a></td><td><span className={`tag ${p.difficulty}`}>{difficultyLabels[p.difficulty]}</span></td><td>{p.points}</td><td>{p.solved ? <span className="tag good">Решена</span> : 'Не решена'}</td><td>{p.wrong}</td></tr>)}</tbody></table></div></section>
	</div>;
}

export function ContestProblemPage({ id, position, user }: { id: string; position: number; user: SessionUser }) {
	const view = useResource<ContestProblem>(`/contests/${id}/problems/${position}`);
	if (!view.data) return <p role="status">{view.error || 'Открываем задачу...'}</p>;
	return <ProblemEditor path={`/contests/${id}/problems/${position}`} contestId={id} problem={view.data} student={user.role === 'student'} reload={view.reload} />;
}

// The code stays in this component, so a reload of the problem after a submission keeps the editor.
function ProblemEditor({ path, contestId, problem, student, reload }: { path: string; contestId: string; problem: ContestProblem; student: boolean; reload: () => void }) {
	// The end in the client clock, from the server time of the response, so a wrong client clock has no effect.
	const [deadline] = useState(() => Date.now() + Date.parse(problem.endsAt) - Date.parse(problem.now));
	const [over, setOver] = useState(() => Date.now() >= deadline), [code, setCode] = useState(problem.code), [status, setStatus] = useState<'saved' | 'saving' | 'unsaved'>('saved');
	const [saveError, setSaveError] = useState(''), [error, setError] = useState(''), [busy, setBusy] = useState(false), [progress, setProgress] = useState(''), [verdict, setVerdict] = useState<ContestVerdict | null>(null);
	const [own, setOwn] = useState<OwnTest[]>([]), [input, setInput] = useState('');
	const latest = useRef(problem.code), stored = useRef(problem.code), pending = useRef(Promise.resolve()), timer = useRef<ReturnType<typeof setTimeout>>(undefined), statement = useRef<HTMLElement>(null);
	const editable = student && !over, solved = problem.solved || verdict?.verdict === 'accepted';
	const spec: TestSpec = { mode: 'stdio', ignoreCase: problem.ignoreCase }, tests: ProblemTests = { spec, source: 'set', standard: problem.examples.map(e => ({ input: e.input, expected: e.output })), own, skeleton: '' };
	// The editor turns read-only at the end without a remount, so a check that started before the end goes on.
	useEffect(() => { const end = setTimeout(() => setOver(true), deadline - Date.now()); return () => clearTimeout(end); }, [deadline]);
	// A page change sends the code that waits for the autosave.
	useEffect(() => () => { if (latest.current !== stored.current) void save(); }, []);
	// One save at a time, so an older code never replaces a newer one. A save sends the newest code.
	function save(): Promise<void> {
		clearTimeout(timer.current);
		pending.current = pending.current.then(async () => {
			const sent = latest.current;
			if (sent === stored.current) return;
			setStatus('saving');
			try { await api(`${path}/code`, { method: 'PUT', body: JSON.stringify({ code: sent }) }); stored.current = sent; setSaveError(''); }
			catch (error) { if (error instanceof ApiError && error.code === 'contest_over') setOver(true); else setSaveError(message(error)); }
			setStatus(latest.current === stored.current ? 'saved' : 'unsaved');
		});
		return pending.current;
	}
	function change(value: string) {
		latest.current = value; setCode(value); setStatus(value === stored.current ? 'saved' : 'unsaved');
		clearTimeout(timer.current);
		timer.current = setTimeout(() => void save(), 1500);
	}
	// A failed save does not stop the submission: the submission carries the code.
	async function submit() {
		setBusy(true); setError(''); setVerdict(null);
		try {
			if (Date.now() < deadline) await save();
			const sent = latest.current, { inputs } = await api<{ inputs: string[] }>(`${path}/check`, { method: 'POST' });
			const outputs = await hiddenOutputs(sent, spec, inputs, (done, total) => setProgress(`Проверяем: ${done} из ${total}`));
			setVerdict(await api<ContestVerdict>(`${path}/submissions`, { method: 'POST', body: JSON.stringify({ code: sent, outputs }) }));
			stored.current = sent; setStatus(latest.current === sent ? 'saved' : 'unsaved');
		} catch (error) { if (error instanceof ApiError && error.code === 'contest_over') setOver(true); setError(message(error)); }
		finally { setBusy(false); setProgress(''); reload(); }
	}
	const example = problem.examples[0]?.input.trimEnd(), placeholder = `Входные данные, как в условии${example ? `, например:\n${example}` : ''}`;
	const left = <section className="stack"><h2>Условие</h2>{problem.statement === null ? <p className="notice">Задача убрана из банка</p> : <Lesson html={problem.statement} ref={statement} />}<p className="small">Источник: Codeforces, набор open-r1/codeforces (CC BY 4.0)</p>{problem.url && <p><a href={problem.url} target="_blank" rel="noopener noreferrer">Задача на Codeforces ↗</a></p>}</section>;
	// The tip of a button shows below it, so the verdict is beside the button.
	const after = <>{editable && !solved && <button type="button" className="primary" disabled={busy || !code.trim()} onClick={() => void submit()} data-tip="Проверить решение на скрытых тестах. Тесты сайт не показывает.">Отправить</button>}{progress && <p role="status">{progress}</p>}{verdict && (verdict.verdict === 'accepted' ? <p role="status"><span className="tag good">Принято</span></p> : <p className="error" role="alert" aria-label="Результат проверки">{verdictText(verdict)}</p>)}</>;
	const right = <div className="stack">
		{editable && <p role="status" aria-label="Состояние сохранения" aria-live="polite">{{ saved: 'Сохранено', saving: 'Сохраняем', unsaved: 'Не сохранено' }[status]}</p>}
		{(error || saveError) && <p className="error" role="alert">{error || saveError}</p>}
		{/* A disabled fieldset does not lock the code editor, which is not a form control. */}
		<fieldset disabled={busy} className="stack"><CodeField value={code} readOnly={!editable} locked={busy} onChange={change} />
			{student && <><Runner code={code} tests={tests} live={editable} before={null} after={after} />
				<section className="card card-pad stack"><h2>Мои тесты</h2><p>Тесты остаются только на этой странице.</p>
					{own.length > 0 ? <ol className="messages" aria-label="Мои тесты">{own.map((test, index) => <li key={test.id}><strong>Тест {tests.standard.length + index + 1}</strong><span className="prose">Ввод: <code>{test.input}</code></span><button type="button" aria-label={`Удалить тест ${tests.standard.length + index + 1}`} onClick={() => setOwn(own.filter(item => item !== test))} data-tip="Удалить этот тест.">Удалить</button></li>)}</ol> : <p>Своих тестов пока нет.</p>}
					<form className="stack" onSubmit={event => { event.preventDefault(); setOwn([...own, { id: crypto.randomUUID(), input, createdAt: new Date().toISOString() }]); setInput(''); }}>
						<label>Ввод<textarea value={input} required maxLength={10_000} rows={placeholder.split('\n').length} placeholder={placeholder} onChange={event => setInput(event.target.value)} /></label>
						<button data-tip="Добавить ваш тест на этой странице.">Добавить тест</button>
					</form>
				</section></>}
		</fieldset>
		<section className="card card-pad stack"><h2>Отправки</h2>{problem.submissions.length ? <ol className="messages" aria-label="Отправки">{problem.submissions.map((item, index) => <li key={index} className={item.verdict === 'accepted' ? 'pass' : 'fail'}>{new Date(item.createdAt).toLocaleTimeString('ru-RU')} {verdictText(item)}</li>)}</ol> : <p>Отправок пока нет.</p>}</section>
	</div>;
	return <div className="stack">
		<div className="heading"><div><p className="small"><a href={`/leetcode/contests/${contestId}`}>Контест</a></p><h1>{problem.position}. {problem.title}</h1><div className="tags"><span className={`tag ${problem.difficulty}`}>{difficultyLabels[problem.difficulty]}</span><span className="small">{points(problem.points)}</span>{solved && <span className="tag good">Решена</span>}</div></div></div>
		<Panes left={left} right={right} />
	</div>;
}
