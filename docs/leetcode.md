# LeetCode notebook

The private notebook at `https://redlamp.me/leetcode` supports one parent and one student. Use LeetCode for code runs and submissions. Use this notebook for explanations, saved attempts, homework, feedback, learning sets, and review reminders.

The notebook server is a Go program in `backend/`. On AWS, it runs as a Lambda function behind CloudFront and keeps its data in Aurora DSQL. A local server uses PostgreSQL 16. The web client uses React and Bun.

## Local setup

Install these tools:

- Bun 1.4.4 and Node.js 24.
- Go, the version in `backend/go.mod`.
- Docker, for PostgreSQL.
- golangci-lint 2.15.0, for the Go lint check.

Install the pinned dependencies in the repository root:

```bash
bun install --frozen-lockfile
```

Start PostgreSQL 16 in Docker. Aurora DSQL uses UTF-8 and the C collation only, so the local database uses them too. Text then sorts in the same order as on AWS.

```bash
docker run -d --name redlamp-backend-test -e POSTGRES_PASSWORD=test -e POSTGRES_INITDB_ARGS="--encoding=UTF8 --locale=C" -p 127.0.0.1:54329:5432 postgres:16.15
```

Set the connection string in each shell that runs the server or the checks:

```bash
export TEST_DATABASE_URL='postgres://postgres:test@127.0.0.1:54329/postgres?sslmode=disable'
```

### Start a local server

Build the site. The build copies the public site files from `public/` to `dist/` and writes the notebook assets to `dist/leetcode/assets/`.

```bash
bun run build
```

Run the next commands in `backend/`:

```bash
cd backend
```

Apply the migrations:

```bash
go run ./cmd/migrate -database-url "$TEST_DATABASE_URL"
```

The tests make their own schemas in the same database. The local server keeps its data in the schema `public`.

Generate the two accounts in an interactive terminal. Replace the example usernames with two distinct names. The first name is the parent, and the second name is the student. Use 3-40 letters, digits, dots, dashes, or underscores.

```bash
export LEETCODE_ACCOUNTS="$(go run ./cmd/credentials parent-name student-name)"
```

The command writes two generated 24-character passwords to the terminal. Save them in a password manager. Standard output contains only the usernames and the password digests, so the variable gets no password. This digest format requires the generated high-entropy passwords. If the variable does not have two valid accounts, the server starts without sign-in and answers "Аккаунты тетради не настроены."

Start the server:

```bash
PUBLIC_ORIGIN=http://localhost:8787 go run ./cmd/api -listen localhost:8787 -database-url "$TEST_DATABASE_URL" -assets ../dist
```

Open `http://localhost:8787/leetcode`. Use the origin of `PUBLIC_ORIGIN` in the browser, because the server refuses a change from another origin. `PUBLIC_ORIGIN` must be an origin in lower case without a path, for example `https://redlamp.me` or `http://localhost:8787`. If it has another form, the server does not start. For `http://localhost` and `http://127.0.0.1`, the server uses a separate session cookie without the `Secure` attribute.

The server serves the files of `dist/` and does not watch them. Run `bun run build` again after a change in `web/`, `shared/`, or `public/`. Restart the server after a change in `backend/`.

A local server has no speech service, so dictation answers "Распознавание речи сейчас недоступно. Попробуйте позже." A local server has a tutor only when `ANTHROPIC_API_KEY` contains an Anthropic API key. Set this variable before you start the server. Each tutor question then goes to Anthropic under the account of that key. See [AI tutor](#ai-tutor).

## Use the notebook

- Add a LeetCode problem. Paste its URL in **Add a problem**. The notebook fills the number, title, difficulty, and topics from LeetCode. If LeetCode does not answer, enter them yourself.
- As the student, start a draft. A draft can hold one or more approaches. The editor autosaves changed content. On a Бомжарня task, use **Отправить**: the server checks the code on hidden tests and accepts the attempt when all tests pass. On another problem, use **Решено** to save the attempt as accepted.
- The understanding of an accepted attempt is automatic: "independent" when the student did not write to the tutor about the problem, else "with help". Both roles can use **Нужно потренироваться**. It sets the understanding "needs practice" and a random review date from 15 to 45 days after today.
- Below **Solution code** of a Python approach, the notebook checks the student's draft code while he types. Ruff and the Python compiler report the problems with their line numbers. **Run tests** runs the code in the browser with Pyodide against the tests of the problem. Each test has 2 seconds. The code is not sent to LeetCode or Codeforces.
- In the code editor, a list of names opens while the student types. It shows the Python keywords, the built-in names, the names that a solution can use without an import, such as `heapq` and `List`, and the names of the code. After a dot, for example `heapq.` or `q.`, [jedi](https://jedi.readthedocs.io/) in the Pyodide worker gives the members. The browser loads jedi on the first completion after a dot. Tab or Enter inserts the selected name. Enter after a name that is typed in full adds a new line. Escape closes the list, and Ctrl+Space opens it. Ctrl+Space without a dot shows the names without jedi.
- Both roles can add up to 20 own tests to a problem in **My tests**. An own test has only the input. **Run tests** runs them after the standard tests.
- A task of a learning set has a reference solution on the server. The reference solution gives the correct answer of each own test. The result of a test shows the correct answer and the answer of the student. The result is green when the two answers are the same and red when they are different. An input outside the constraints of the task gets a note and no check. A problem outside the sets has no reference solution, so its own tests show only the answer of the student.
- A task of a learning set has its standard tests in the set. For another LeetCode problem, the standard tests are the examples of the statement from LeetCode. A problem with a `TreeNode`, a `ListNode` or a class to design has no runner. A Codeforces problem outside the sets has only own tests.
- The notebook records each test run of the student. An approach shows **Tests passed** when its latest run passed all tests and its code did not change after that run. The attempt history shows the same mark for the attempt.
- As the parent, create homework on the **Homework** page. Give it a title, instructions, and an optional due date. Use **Edit homework** to change these details later. Use **Изменить срок** next to the due date to change only the date. A new due date lets the student use **Сила Бомжа** again.
- On the homework page, use **Add a task** to add a problem from the library. Each homework task has its own state. A problem, chapter, or quiz deck can have only one active task in all homework.
- As the parent, open **Все домашние работы** to find homework, including homework with no tasks yet. Use **Ещё домашние работы** to reach older homework.
- The **Homework** page shows the tasks of all homework in three sections. **В процессе** has the active tasks of the homework with the earliest due date. **Будущее** has the active tasks of the other homework, by due date. **Принятое** has the completed tasks, with the newest first. The page does not show cancelled tasks.
- In **Принятое**, use **Распределить повторения на 30 дней**. Each problem of a completed task gets a random review date from tomorrow to 30 days after today. Two problems get the same date only when there are more than 30 problems. A problem with a review date in this period keeps its date. On the review date, **Today** shows the problem in **Пора повторить**. The link opens the recall page of a pick, so the student solves the problem again without the saved work.
- The parent can link learning sets to problems and homework in **Materials**. **Materials** also shows the set chapters that contain the problem. For homework, it shows assigned chapters and the chapters of its problems and quizzes.
- A problem task opens the problem page. The page shows the homework, the due date, the review state, and the newest attempt. As the student, open a task from **Today** or the homework page. **Решено** submits the attempt for review. The parent can request another attempt or complete the task. An attempt that passes the hidden tests completes the task. The parent can still return a completed problem task.
- A red dot before a homework or a task shows a message or a review decision from the other account that you did not read. Open the discussion of the task or the problem to remove the dot.
- Use **Recall this problem** to practise before revealing saved work. Record the result. Choose the next review date. When revealing saved work, use **Старые попытки** and **Новые попытки** to page through the saved attempts.
- **Today** also shows random recall picks from the library: three problems of the Two Pointers set and two problems of the first permutations homework. The picks change every day. A pick that you recall today leaves the list. The recall page of a pick does not show the summary or the saved work. Use **Start a new attempt** to solve the problem again. If the problem has a draft, the button saves the draft as an attempt and starts a new empty draft. This save does not change the understanding of the problem.
- **Archive problem** moves a problem to **Archived problems** in the library and blocks new homework tasks for it. **Restore problem** returns it. Restore the problem before returning a completed homework task for another attempt.
- As the student, use **Голосовой ввод** (the microphone button) to dictate Russian text into a notes field, the review note, the tutor message, or the problem description. Select the button again to stop. A recording stops after 100 seconds. The notebook puts the text at the cursor, or at the end when you did not select a position in the field. The browser sends the audio as mono 16 kHz PCM. The server sends it to Amazon Transcribe streaming with the language Russian and does not store the audio.
- If a save fails, keep the page open. Use **Retry save** or **Copy unsaved work**. Reload recovery is available in the same tab when browser session storage works. Conflicting edits require an explicit choice of the server copy.

## Learning sets

A learning set is a sequence of chapters about one subject. Each chapter has a lesson, a list of problems and a deck of quiz cards. The sets are content in the repository. A deployment publishes them. The **Бомжарня** page lists the sets. It is read-only for both roles.

A set page shows the reading links and the chapters. A chapter row shows the accepted problem count and the quiz status. A chapter page shows the chapter lesson, the problems, the quiz status and a link to the chapter quiz. A lesson can show figures. A figure with more than one step shows **Назад** and **Вперёд**. A hard word, for example "инвариант", has a dotted line at its first use in a lesson, a statement, a summary or a quiz card. Its tip shows on hover, on keyboard focus and after a tap. A problem row shows **Accepted** when the library problem has an accepted saved attempt. A problem row with an active homework task shows the task state. The state links to the problem page.

Learning explanations use Russian. Algorithm and data-structure names use English, including titles, topics, figures, quiz cards, and reading-link labels. Use the canonical topic names in `shared/topics.ts`. For example, `Hash Table`, `HashTable`, and `Хеш-таблица` refer to one `Hash Table` progress topic. A problem counts once for each topic. The progress count uses accepted saved attempts and excludes archived problems.

Each worked example explains its input sizes, counted operations, time bound, and peak auxiliary memory. State the assumptions for average or amortized bounds. Count Python copies, sorting workspace, and pending calls. Separate input storage and growing output storage from auxiliary memory. Show a concrete count table or figure. Explain a useful alternative when the chapter compares techniques. Label model growth values as model values, not measured time. Each problem note includes **Время:** and **Память:** with its variables and assumptions.

Keep each lesson below 12,000 text characters, including collapsed sections. Keep the main explanation short. Put longer comparisons in `details` elements. The figures are data in `set.go`, outside the lesson HTML. Put their main count and conclusion in the lesson text too, so the tutor also receives them. The tutor allows 50,000 lesson characters for a set and 16,000 characters for a problem chapter and its outline.

As the parent, use **Добавить главу в домашнюю работу** on a chapter page. This adds the chapter reading, its quiz deck, and its available problems together. The dialog shows active homework. Use the search box to find homework by title. Use **Назначить** to add the chapter to that homework. Use **Создать и назначить** to create homework with the chapter. The notebook adds missing problems to the library. It skips problems that are complete, archived, or already assigned. The result lists each skipped item and its reason.

To assign only a quiz, open its page as the parent and use **Добавить квиз в домашнюю работу**. A new passing attempt completes the quiz homework task. An earlier pass does not complete a new assignment. As the student, use **Прочитал главу** on an assigned chapter to complete its reading task. As the parent, use **Отменить** on the homework page to cancel an active chapter or quiz task.

For the student, a chapter page shows **Объяснить с ИИ** below the text that is selected in the lesson. Ctrl+Enter or Cmd+Enter does the same. The AI tutor then explains the selected text. When an answer quotes the lesson, **Показать в уроке** marks the passage. The parent does not see these controls.

The parent can also link a set to a problem or to homework. Choose the set in **Materials**. The link opens the set page.

The **Квиз** page shows one deck for each chapter. An attempt answers all cards of one deck. The page sends the answers to the server after the last card. If the student leaves before the last card, the notebook keeps no answers. An attempt passes with 80 % or more correct answers: 8 cards need 7, 7 cards need 6, 6 cards need 5. After a passed attempt, the page shows each card with the answer of the student, the right answer and the explanation. After a failed attempt, the page shows only the cards with a wrong answer, without the right answers. After each attempt, passed or failed, the notebook closes the deck for 1 hour. The quiz status is **Пройден** with the best passed result, or **Не пройден**, and the time of the next attempt during the pause. The parent sees all cards of a deck with the answers. A quiz round shows no AI tutor.

A learning set is a Go package in `backend/internal/sets/<name>/`. The model of a set is in `backend/internal/sets/model.go`. To add a set:

1. Create the folder `backend/internal/sets/<name>/`. A Go package name has no dashes, so the folder name can be different from the slug of the set, for example `twopointers` for `two-pointers`. Use lowercase letters, digits, and dashes in the slug.
2. Write `chapter<N>.html` for each chapter as an HTML fragment. The chapter page inserts it without changes. Start with `<p class="lead">` and one sentence. Then write sections that start with `<h2 id="sN">`. Use only these elements: `h2 h3 p ul ol li table caption thead tbody tr th td pre code strong em sup blockquote details summary a br figure`, and `<div class="rule">` for a rule. Put a picture as `<figure data-figure="<id>"></figure>` on its own line. Do not use scripts, styles, `style` attributes or event attributes. Write code blocks as `<pre><code class="language-python">`. Write a hint as `<details><summary>Подсказка 1</summary>...</details>`. Give external links `target="_blank" rel="noopener noreferrer"`.
3. Write `set.go`. Embed each lesson with `//go:embed chapter<N>.html` in a `string` variable. Export `var Set = sets.LearningSet{...}` with the `Slug`, `Title`, `Summary`, `Topics` and `Chapters` of the set, with the chapters in reading order. Give each chapter a `Title`, a `Summary`, its `Lesson`, `Tasks`, `Cards`, `Links` and `Figures`. For each task, give the problem `URL` on LeetCode or Codeforces, its `Number`, `Title`, `Difficulty`, `Topics`, `Note` and `Statement`, and the standard tests in `Tests`. Do not change the slug of a published set. The quiz attempts keep the deck as `<set slug>-<chapter number>`.
4. Describe each picture in `Figures` of its chapter, in the order of the placeholders in the lesson. A figure has an ID, a Russian title and one or more frames. A frame has a Russian caption and rows of cells. A row can have a label, marked cells, dimmed cells, pointers under the cells and a window bracket. A row with `Tree: true` draws its cells as a binary tree in `heapq` order: the children of place i are places 2i + 1 and 2i + 2. An empty string is an empty place. A tree row has no pointers and no window. A marked cell is a cell that the step works on. A dimmed cell holds a value that a move removed. Do not show a pointer that is past the end of its row. Say in the caption where that pointer is. Make each frame from a trace of the algorithm on the numbers of the picture.
5. Write `solutions.go`. Export `var References = map[string]sets.Reference{...}` with a reference solution for each task, by task slug. A `call` reference has `Arity` and `Call`, and `Call` returns the value that the tests compare. The arguments and the result of `Call` are Python values of the package `answers`: make the result from a Go number or a Go slice with `sets.Py`. A `stdio` reference has `Stdio`, and `Stdio` returns the standard output. First check the constraints of the statement with the helpers in `backend/internal/sets/constraints.go`, for example `sets.Int`, `sets.Ints` and `sets.NewReader`. A failed check panics with `sets.ErrConstraints` from `backend/internal/sets/reference.go`. The references stay on the server, so the student cannot read them.
6. Write `generators.go`. Export `var Generators = map[string]sets.Generator{...}` with a generator of hidden tests for each task, by task slug. A generator gets a seeded random function and returns 10 to 16 inputs in the test format of the task. Make the input of a call task with `sets.Call`. A set without tasks needs no `solutions.go` and no `generators.go`.
7. In `backend/internal/sets/content/content.go`, import the package. Add its `Set` to `Sets`. The order of `Sets` is the display order. Add its `References` to `References` and its `Generators` to `Generators`.
8. Find the hard words in the student text: the summaries, the lessons, the statements, the notes and the cards. For each hard word that `shared/glossary.ts` does not have, add an entry. Give the entry a pattern for the word and its forms. Give it a plain Russian tip of one or two short sentences. Add the start of the word to `hardWord` in `web/glossary.test.ts`. Do not add a word that has two meanings in the sets. Use exact English names such as `Array`, `Hash Table`, `Two Pointers`, `Heap (Priority Queue)`, `Queue`, `Deque`, and `Stack`. Use precise names such as `Min-Heap` or `Monotonic Deque` when explaining an implementation. Keep Python API names such as `heapq` and `collections.deque` unchanged. Use "очередь" only for a real line of people.
9. Write the JSON copy of the learning sets again. The web tests and the browser tests read this copy in `backend/internal/sets/settest/sets.json`. Run this command in `backend/`:

   ```bash
   go test ./internal/sets/content/ -run TestSetsEqualTheFile -update
   ```

10. Each chapter has one quiz deck. Add the decks of the new chapters to the expected list of `TestQuizDecks` in `backend/internal/notebook/quiz_test.go`. Increase the deck count in `e2e/quiz.spec.ts`.
11. Run the checks, see [Checks](#checks). In `backend/`, `go test ./...` runs the content checks of `backend/internal/notebook/content_checks_test.go` and `backend/internal/notebook/sets_references_test.go`. They check content identities, problem counts, links, lesson size and format, figure placeholders and indexes, quiz cards, English technical names, problem cost notes, and reference answers. `bun run test` checks the glossary coverage. `bun run test:browser` checks executable lesson outputs, figures, and navigation.

The file `backend/internal/sets/settest/references.json` has recorded hidden tests of earlier tasks. A task without recorded hidden tests needs no entry there. When you change the standard tests or the generator of a recorded task, change its entry in this file or remove the entry.

To add a chapter to a published set:

1. Write `chapter<N>.html` for the new chapter. Use the next number after the last chapter.
2. In `set.go` of the set, embed the new lesson and add the chapter at the end of `Chapters`. Add a new chapter only after the last chapter, because the quiz attempts keep the deck as `<set slug>-<chapter number>`.
3. For each new task, add a reference to `References` in `solutions.go` and a generator to `Generators` in `generators.go`.
4. Add the missing hard words, as in step 8 of the list above.
5. Write `sets.json` again and run the checks, as in steps 9 and 10 of the list above.

## AI tutor

The AI tutor uses the Anthropic API with the model Claude Opus 5.5 (`claude-opus-5-5`), adaptive thinking and the effort `xhigh`. When the API declines a question for this model, another Claude model of the API can answer it. The server uses the API key of the parent. On AWS, the key is in the SecureString parameter `/redlamp/anthropic-api-key`, see [Secrets](#secrets). The server reads the parameter that `ANTHROPIC_API_KEY_PARAMETER` names. When `ANTHROPIC_API_KEY` has a value, the server uses that value and does not read the parameter. The server reads the key when it starts.

Each question of the student goes to Anthropic under the API account of the parent, with the earlier messages of the chat and the context of the open page. This context can contain the problem, the LeetCode or Codeforces statement, the homework instructions, the lesson text, and the draft or the newest saved attempt of the student with its notes and code. For a problem in a set chapter, the context also has the lesson text of the first such chapter.

The notebook has no usage cap. The parent sets a spend limit in the Anthropic Console.

The tutor gives hints. It does not write the solution, also when the student asks for it. For a new question about how to solve a problem, it asks one guiding question. When the student asks for more help, it names the idea, then gives the key step on a small example, and then names the parts of the function. When the student asks it to check the code, it tells about one bug in each answer and gives a small input that shows the bug. It does not write the corrected code. The rules are in `backend/internal/notebook/tutor-prompt.txt`.

On a problem page with code, **Оценить решение с ИИ** asks the tutor for an evaluation of the code and the notes of the student. The answer ends with a score from 1 to 10. The rules of the evaluation are in `backend/internal/notebook/evaluate-prompt.txt`.

Each problem and set has one chat. The student selects **Спросить ИИ** on these pages. On a screen of 1400 px and wider, the chat panel stays open beside the page, also on the next page. On a smaller screen, the panel opens over the page. Ctrl+Enter or Cmd+Enter sends the message. On a chapter page, the student can also select text in the lesson and choose **Объяснить с ИИ**. See [Learning sets](#learning-sets).

Only the student can read and use the chats. The parent cannot read them, and the pages of the parent show no tutor.

The thinking text of the model can contain the full solution, so it never goes to the browser. The server sends only the text of the answer. While the model thinks, the server sends a sign every 10 seconds, and the panel shows "Репетитор думает." The server stops a request to Anthropic after 9 minutes.

The notebook keeps a turn only when the answer is complete. When the browser closes the connection, the server still reads the answer to its end and keeps the turn. If an answer stops, the student sends the question again. The tutor shows these errors:

| Code | Cause |
|---|---|
| `setup_required` | The server has no API key, or the Anthropic API refuses the key, the billing or the permission. |
| `ai_usage_limit` | The Anthropic API answers with a rate limit. |
| `ai_unavailable` | The Anthropic API is overloaded or has an error, or the server cannot reach it. |
| `ai_interrupted` | The answer stops before its end, reaches the token limit or takes more than 9 minutes, or the server cannot store the turn. |
| `ai_refused` | The model refuses to answer the message. |
| `ai_failed` | The Anthropic API refuses the request for another reason. |

## Levels and badges

The notebook uses the Arcade theme. Dark is the default appearance. Use the appearance button in the top bar to change to light and back. The browser keeps the choice.

The top bar shows the level and the progress to the next level. The bar is hidden on screens of 620 px or less. **Today** shows the level, the XP, and the badges.

The notebook calculates XP from saved records. Failed test runs, late homework and missed days decrease the XP. The XP is never less than 0. Archived problems keep their XP. The notebook records a visit day when the student opens a notebook page or moves to another page. A missed day is a UTC day without a visit, from the first visit day to yesterday.

| Source | XP |
|---|---|
| A problem with an accepted saved attempt, one time for each problem | easy 10, medium 20, hard 40 |
| An accepted problem with a saved attempt or a review with the understanding "independent", one time for each problem | 10 |
| A completed homework task | 15 |
| A review, one time for each problem and day | 5 |
| A correct answer in the best passed quiz attempt of each deck | 2 |
| A test run that passed all tests, one time for each problem | 10 |
| A test run with a failed test, at most three runs for each problem and day | -2 |
| A problem task without a submission, or a chapter or quiz task without completion, by its due date | -15 |
| A missed day | -5 |

The correct answers of the earlier quiz, which saved one answer for each card, keep their 2 XP and count for **50 карточек**.

Level L starts at `50 * L * (L - 1)` XP: level 2 at 100 XP, level 3 at 300 XP, level 4 at 600 XP, and level 5 at 1000 XP.

A badge stays earned, except Уровень 5: it shows the current level. The student and the parent see the same badges.

| Badge | Rule |
|---|---|
| First accept | One problem with an accepted saved attempt. |
| 3-day streak | Saved attempts on three consecutive days. The notebook uses UTC dates. |
| First hard | One hard problem with an accepted saved attempt. |
| Ten accepted | Ten problems with an accepted saved attempt. |
| Independent five | Five accepted problems with the understanding "independent". |
| Set cleared | An accepted saved attempt for each problem of one learning set. |
| Неделя без пропусков | Visit days on 7 consecutive days. |
| Месяц без пропусков | Visit days on 30 consecutive days. |
| 25 принятых | 25 problems with an accepted saved attempt. |
| 50 принятых | 50 problems with an accepted saved attempt. |
| 50 карточек | 50 correct answers in the best passed quiz attempts of the decks. |
| Десять повторений | 10 reviews, one for each problem and day. |
| Первая с Codeforces | One Codeforces problem with an accepted saved attempt. |
| С первого раза | The first test run of a problem passed all tests. |
| Два набора | An accepted saved attempt for each problem of two learning sets. |
| Вовремя | 5 homework tasks submitted or completed by the due date. |
| Уровень 5 | Level 5, which starts at 1000 XP. |

The hidden tests or **Решено** set the acceptance, and the understanding is automatic. The **Today** counts do not include archived problems, so they can be different from the XP.

## AWS setup

The site runs in one AWS account. The templates are in `infra/`. Route 53 takes the key of a DNSSEC key-signing key only from us-east-1, and CloudFront takes a certificate only from us-east-1. For these reasons, the zone stack and the edge stack are in us-east-1. The other stacks are in eu-central-1.

| Stack | Region | Template | Parameters | Outputs |
|---|---|---|---|---|
| `redlamp-dns` | us-east-1 | `infra/dns.yaml` | `AlarmEmail` | `HostedZoneId`, `NameServers` |
| `redlamp-data` | eu-central-1 | `infra/data.yaml` | none | `ClusterEndpoint`, `ClusterArn` |
| `redlamp-api` | eu-central-1 | `infra/template.yaml` | `ClusterEndpoint`, `ClusterArn`, `DistributionId`, `AlarmEmail` | `FunctionUrl`, `FunctionRoleArn` |
| `redlamp-edge` | us-east-1 | `infra/edge.yaml` | `HostedZoneId`, `FunctionUrl` | `BucketName`, `DistributionId`, `DistributionDomainName` |
| `redlamp-pipeline` | eu-central-1 | `infra/pipeline.yaml` | `AlarmEmail`, `ClusterArn` | none |

- `redlamp-dns` has the zone `redlamp.me`, the SPF and DMARC records of a domain that sends no mail, the KMS key and the key-signing key of DNSSEC, and two DNSSEC alarms on the topic `redlamp-dns-alarms`. CloudFormation keeps the zone when it deletes the stack, because the domain registration names the name servers of this zone.
- `redlamp-data` has the Aurora DSQL cluster with deletion protection and its daily backup. See [Backup and recovery](#backup-and-recovery).
- `redlamp-api` is an AWS SAM template. It has the Lambda function of `backend/cmd/api` on arm64, with 512 MB of memory and a timeout of 600 seconds. Lambda runs at most 10 instances of the function at one time. The function URL uses IAM authentication and response streaming. The stack also has a log group with 30 days of retention, three alarms on the topic `redlamp-api-alarms`, and a monthly cost alert at 10 USD for AWS Lambda and Amazon CloudFront. The alarm `redlamp-api-errors` is for a function error of Lambda. The alarm `redlamp-api-failures` is for an internal failure, for which the server answers 503 with the code `unavailable`. The alarm `redlamp-data-backup` is for a missing backup.
- `redlamp-edge` has the certificate of `redlamp.me` and `www.redlamp.me`, the site bucket `redlamp-site-<account ID>`, the CloudFront distribution with its viewer request function, and the alias records of the two names.
- `redlamp-pipeline` has the GitHub connection `redlamp-github`, the artifact bucket, the CodeBuild projects `redlamp-build` and `redlamp-deploy`, the pipeline `redlamp`, the role that CloudFormation uses for the four stacks of the site, the permissions boundary `redlamp-boundary` for the IAM roles of these stacks, and the topic `redlamp-pipeline-alarms` for a failed execution.

CloudFront sends `/leetcode`, `/leetcode.html` and `/leetcode/*` to the function URL, except `/leetcode/assets/*`. It sends the other paths to the site bucket. On the public site, a path that ends with a slash gets `index.html`, and a path without a dot gets the extension `.html`. A request to `www.redlamp.me` gets a redirect to `https://redlamp.me` with the same path and query. CloudFront signs each request to the function URL, so the browser sends the SHA-256 of the body of each change request in the header `x-amz-content-sha256`.

The function has these environment variables: `DSQL_ENDPOINT` with the cluster endpoint, `PUBLIC_ORIGIN` with `https://redlamp.me`, `LEETCODE_ACCOUNTS_PARAMETER` with `/redlamp/leetcode-accounts`, and `ANTHROPIC_API_KEY_PARAMETER` with `/redlamp/anthropic-api-key`. The function connects to the cluster as the database role `app`. `cmd/migrate` maps this role to the IAM role of the function. Sessions expire after seven days. The server permits 10 sign-in attempts in 60 seconds for each source address. Each function instance counts the attempts separately.

### Secrets

The secrets are SecureString parameters of SSM Parameter Store in eu-central-1. The owner writes both parameters. Each parameter uses the default key `aws/ssm`. The policy of this key lets each principal of the account decrypt through Parameter Store. Another key needs the permission `kms:Decrypt` in the role of the function, and the role does not have it. The function stack lets the function read them. The server reads each secret from the parameter that its variable names: `LEETCODE_ACCOUNTS_PARAMETER` and `ANTHROPIC_API_KEY_PARAMETER`. When `LEETCODE_ACCOUNTS` or `ANTHROPIC_API_KEY` has a value, the server uses that value and does not read the parameter. The server reads both secrets one time, when a function instance starts, so an instance that already runs keeps the old values.

- `/redlamp/leetcode-accounts`: the accounts secret, the standard output of `go run ./cmd/credentials`.
- `/redlamp/anthropic-api-key`: the Anthropic API key of the parent.

Write the accounts with this command. Run it in `backend/`, in an interactive terminal, with the intended usernames. The command shows the new passwords on the terminal, and only the JSON of the accounts goes into the parameter:

```bash
aws ssm put-parameter --region eu-central-1 --name /redlamp/leetcode-accounts --type SecureString --value "$(go run ./cmd/credentials parent-name student-name)"
```

Add `--overwrite` when the parameter exists. Write `/redlamp/anthropic-api-key` in the Parameter Store console. Do not type the API key in a command line, because the shell keeps the line in its history.

To use a new value at once, change the configuration of the function, for example its description. After such an update, Lambda provisions the function again, so new instances serve the requests ([UpdateFunctionConfiguration](https://docs.aws.amazon.com/lambda/latest/api/API_UpdateFunctionConfiguration.html)):

```bash
aws lambda update-function-configuration --region eu-central-1 --function-name "$(aws cloudformation describe-stack-resource --region eu-central-1 --stack-name redlamp-api --logical-resource-id ApiFunction --query StackResourceDetail.PhysicalResourceId --output text)" --description "secrets <date>"
```

The repository is public, so the alarm address is not in the templates. The stacks `redlamp-dns`, `redlamp-api` and `redlamp-pipeline` take it as the parameter `AlarmEmail`. SNS sends a message only after the address confirms the subscription. SNS deletes a subscription that is not confirmed in 48 hours. The cost alert of `redlamp-api` also goes to this address.

### First deploy

The stacks take outputs of other stacks as parameters. For the first deploy, use this order:

1. Deploy `redlamp-dns` with `AlarmEmail`. The domain registration of `redlamp.me` must name the name servers of the zone. The output `NameServers` gives them.
2. Deploy the stack `redlamp-pipeline` with `AlarmEmail` and with the capability `CAPABILITY_NAMED_IAM`. Leave `ClusterArn` empty. The stack has the policy `redlamp-boundary`, and the roles of `redlamp-data` and `redlamp-api` name this policy as their permissions boundary. CloudFormation makes the connection `redlamp-github` with the status `PENDING`.

   ```bash
   aws cloudformation deploy --region eu-central-1 --stack-name redlamp-pipeline --template-file infra/pipeline.yaml --capabilities CAPABILITY_NAMED_IAM --parameter-overrides AlarmEmail=<address>
   ```

3. Deploy `redlamp-data`.
4. Deploy the stack `redlamp-pipeline` again with the same command, and with `ClusterArn=<output ClusterArn of redlamp-data>` as its only parameter. The address stays.
5. Write the two SecureString parameters, see [Secrets](#secrets).
6. Deploy `redlamp-api` with `ClusterEndpoint` and `ClusterArn` from `redlamp-data`, and with `AlarmEmail`. Leave `DistributionId` empty. The distribution needs the function URL, so the first deploy of this stack has no distribution.
7. Deploy `redlamp-edge` with `HostedZoneId` from `redlamp-dns` and `FunctionUrl` from `redlamp-api`. ACM adds the validation records of the certificate to the zone.
8. A person completes the connection `redlamp-github` in the console. Before this step the connection cannot start the pipeline, and a merge to `main` deploys nothing.

`infra/buildspec-deploy.yml` has the commands that the pipeline runs for the four stacks of the site, with the region and the parameters of each stack. For a deploy by hand, leave out `--role-arn`, and give `AlarmEmail` on the first deploy of a stack that has this parameter. After the first deploy, each execution of the pipeline applies the migrations and deploys `redlamp-api` again with `DistributionId` from `redlamp-edge`. See [Releases](#releases).

## Releases

The sidebar shows the version from `package.json`. The versions follow [Semantic Versioning](https://semver.org/). Increase the version in `package.json` in the pull request that changes the notebook. A feature increases the minor version. A fix increases the patch version.

The pipeline `redlamp` can start after the last step of [First deploy](#first-deploy), when its connection is complete. A push to `main` then starts the pipeline in eu-central-1. An execution does not overtake an earlier one, so two pushes deploy in their order. It has three stages:

1. Source: the connection `redlamp-github` gives the files of `main`.
2. Build: the CodeBuild project `redlamp-build` runs `infra/buildspec-build.yml`. It starts PostgreSQL 16.15 in Docker and runs `go vet ./...` and `go test ./...` in `backend/`. It builds the command `migrate`, runs `bun run lint`, `bun run typecheck` and `bun run test`, and runs `sam build` in `infra/`. The artifact has `build/`, `dist/` and `infra/`.
3. Deploy: the CodeBuild project `redlamp-deploy` runs `infra/buildspec-deploy.yml`. It deploys the stacks in this order: `redlamp-dns`, `redlamp-data`, the database migrations, `redlamp-api`, `redlamp-edge`. Then it uploads `dist/` to the site bucket and clears the cache of CloudFront.

The migrations run before the new code of the function. They are the files of `backend/migrations/`, with one SQL statement in each file. `cmd/migrate` applies the files that the table `schema_migrations` does not have, in the order of their names. On AWS, it connects as the database role `admin` and maps the database role `app` to the IAM role of the function.

A failed execution sends a message to the topic `redlamp-pipeline-alarms`. CloudFormation changes the four stacks of the site with the role of the stack `redlamp-pipeline`. This role can make and change only the IAM roles of `redlamp-data` and `redlamp-api`, and only a role with the permissions boundary `redlamp-boundary`. A role with this boundary has only the rights that are also in the boundary. To give such a role a new right, add the right to the boundary in `infra/pipeline.yaml` and deploy the stack `redlamp-pipeline` by hand before the merge. The role still has the full rights of its services in the account.

After a deploy, check sign-in, sign-out, and a normal draft save on `https://redlamp.me/leetcode`. The logs of the function are in the log group of `redlamp-api`.

## Credential rotation

Make new credentials for both accounts and write them to the parameter. Run this command in `backend/`, in an interactive terminal, with the intended usernames:

```bash
aws ssm put-parameter --region eu-central-1 --name /redlamp/leetcode-accounts --type SecureString --overwrite --value "$(go run ./cmd/credentials parent-name student-name)"
```

Save the new passwords that the command shows on the terminal. Changing an account's username or password digest invalidates its existing sessions. For a local server, set `LEETCODE_ACCOUNTS` to the standard output of `go run ./cmd/credentials parent-name student-name`, as in [Local setup](#local-setup). Then restart the server.

To rotate the API key of the tutor, write the new key to the parameter `/redlamp/anthropic-api-key` in the Parameter Store console.

A function instance that already runs keeps the old accounts and the old key. Replace the instances with the command in [Secrets](#secrets).

## Backup and recovery

AWS Backup makes a backup of the cluster each day at 17:00 UTC and keeps it for 35 days, in the vault `redlamp-data`. Each backup is a full copy of the cluster. A backup needs a cluster in the Active state. A cluster without connections goes to the Idle state, and its backup fails. For this reason, the backup starts in the hours of normal use. The alarm `redlamp-data-backup` is on when no backup job completed in two days. Then open the notebook, which makes the cluster Active. The next job at 17:00 UTC makes a backup. Aurora DSQL refuses to delete the cluster while deletion protection is on.

A restore makes a new cluster from a backup. Each pipeline execution gives `redlamp-api` the endpoint and the ARN of the cluster from the outputs of `redlamp-data`. The restored cluster has the session records of the backup. Rotate both accounts after recovery to invalidate these sessions. This document does not give the steps of a restore: do one restore test first, then write the steps here.

## Checks

The Go tests and the browser tests use the local PostgreSQL and fixed test credentials. They do not use production data. Set `TEST_DATABASE_URL` first, see [Local setup](#local-setup).

Run the web checks in the repository root:

```bash
bun run lint
```

```bash
bun run typecheck
```

```bash
bun run test
```

`bun run test` builds the site, runs the test of the CloudFront function in `scripts/edge.test.ts`, and then runs the Vitest tests in `web/`.

Run the Go checks in `backend/`:

```bash
go vet ./...
```

```bash
go test ./...
```

```bash
golangci-lint run ./...
```

Without `TEST_DATABASE_URL`, `go test` skips the database tests. When the variable `CI` is set, a missing `TEST_DATABASE_URL` fails them.

Run the browser tests in the repository root. Install Chromium for Playwright one time:

```bash
bunx playwright install chromium
```

Then run the tests:

```bash
bun run test:browser
```

`bun run test:browser` builds the site, builds the Go commands `api` and `migrate` to `.playwright/bin/`, and runs the Playwright tests in Chromium. Each test starts its own Go server with a new schema in the database of `TEST_DATABASE_URL`. These servers have no tutor key and no speech service, and they cannot reach leetcode.com. The report is in `.playwright/report/`. The results, with the screenshots of failed tests, are in `.playwright/results/`.

The workflow `Checks` runs the web checks in the job `checks` and the browser tests in the job `browser`, for each push and pull request. The workflow `Backend` runs the Go checks in the job `go` when a change touches `backend/`, `shared/answers.cases.json` or the workflow file.

The notebook assets in `/leetcode/assets/` are public. The server checks the session before it serves notebook data, lessons, or the application shell.
