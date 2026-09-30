package dialog

import (
	"context"
	"fmt"
	"strings"

	"section-near-home/backend/internal/clients/maxapi"
)

// quizOption — один вариант ответа на вопрос теста. Direction — ключ
// направления из каталога sportOptions (onboarding.go), к которому
// относится этот вариант; за него начисляется балл этому направлению.
type quizOption struct {
	Text      string
	Direction string
}

type quizQuestion struct {
	Text    string
	Options []quizOption
}

// quizDirectionOrder — фиксированный порядок направлений: используется
// везде, где нужен детерминированный вывод (равенство баллов, клавиатуры).
var quizDirectionOrder = []string{"football", "judo", "swimming", "art", "theatre"}

// quizQuestions — профориентационный тест из 8 вопросов. В каждом вопросе
// ровно 5 вариантов, по одному на каждое направление каталога, без
// "правильных" ответов. Содержание и порядок соответствуют утверждённой
// методике теста ACT.
var quizQuestions = []quizQuestion{
	{
		Text: "Что ребёнку больше всего нравится делать в свободное время?",
		Options: []quizOption{
			{Text: "Бегать, играть в мяч, быть в движении", Direction: "football"},
			{Text: "Соревноваться, мериться силами", Direction: "judo"},
			{Text: "Плескаться в воде, купаться", Direction: "swimming"},
			{Text: "Рисовать, лепить, мастерить", Direction: "art"},
			{Text: "Придумывать истории, наряжаться, показывать кому-то что-то", Direction: "theatre"},
		},
	},
	{
		Text: "Как ребёнок реагирует на соревнования?",
		Options: []quizOption{
			{Text: "Обожает командные матчи и голы", Direction: "football"},
			{Text: "Любит личные поединки и разряды", Direction: "judo"},
			{Text: "Не любит соревноваться, важнее сам процесс", Direction: "swimming"},
			{Text: "Не участвует в соревнованиях, любит творческий процесс", Direction: "art"},
			{Text: "Любит выступления и аплодисменты", Direction: "theatre"},
		},
	},
	{
		Text: "Какая обстановка ребёнку комфортнее всего?",
		Options: []quizOption{
			{Text: "Открытое поле, стадион", Direction: "football"},
			{Text: "Татами, зал единоборств", Direction: "judo"},
			{Text: "Бассейн, вода", Direction: "swimming"},
			{Text: "Спокойная студия с красками/материалами", Direction: "art"},
			{Text: "Сцена, зрительный зал", Direction: "theatre"},
		},
	},
	{
		Text: "Что говорят про ребёнка воспитатели или учителя?",
		Options: []quizOption{
			{Text: "Активный, лидер в подвижных играх", Direction: "football"},
			{Text: "Дисциплинированный, любит чёткие правила", Direction: "judo"},
			{Text: "Спокойный, не боится воды", Direction: "swimming"},
			{Text: "Усидчивый, любит рукоделие/рисование", Direction: "art"},
			{Text: "Артистичный, любит быть в центре внимания", Direction: "theatre"},
		},
	},
	{
		Text: "Любимое занятие ребёнка дома",
		Options: []quizOption{
			{Text: "Играть в мяч во дворе", Direction: "football"},
			{Text: "Бороться, показывать приёмы", Direction: "judo"},
			{Text: "Купаться в ванной/бассейне", Direction: "swimming"},
			{Text: "Рисовать, лепить, собирать конструктор", Direction: "art"},
			{Text: "Устраивать домашние спектакли, петь, наряжаться", Direction: "theatre"},
		},
	},
	{
		Text: "Какой герой из мультфильмов/книг ближе ребёнку?",
		Options: []quizOption{
			{Text: "Командный герой-спортсмен", Direction: "football"},
			{Text: "Герой-боец, побеждающий благодаря дисциплине", Direction: "judo"},
			{Text: "Герой, связанный с водой/морем", Direction: "swimming"},
			{Text: "Герой-художник/изобретатель", Direction: "art"},
			{Text: "Герой-артист/певец", Direction: "theatre"},
		},
	},
	{
		Text: "Как ребёнок реагирует на физическую нагрузку?",
		Options: []quizOption{
			{Text: "С радостью бегает и не устаёт", Direction: "football"},
			{Text: "Любит силовые упражнения, растяжку", Direction: "judo"},
			{Text: "Любит плавные, ритмичные движения", Direction: "swimming"},
			{Text: "Предпочитает спокойные, усидчивые занятия", Direction: "art"},
			{Text: "Двигается вместе с эмоциями: танцы, мимика", Direction: "theatre"},
		},
	},
	{
		Text: "Что важнее всего в занятиях для ребёнка?",
		Options: []quizOption{
			{Text: "Быть частью команды", Direction: "football"},
			{Text: "Развивать характер и самодисциплину", Direction: "judo"},
			{Text: "Укреплять здоровье и уверенность в воде", Direction: "swimming"},
			{Text: "Творческое самовыражение", Direction: "art"},
			{Text: "Внимание и эмоции публики", Direction: "theatre"},
		},
	},
}

// directionReasons — короткие формулировки "почему подходит" по направлению.
// Подставляются рядом со ссылкой на секцию по итогам теста.
var directionReasons = map[string]string{
	"football": "активным и командным детям, которые любят двигаться и играть в команде",
	"judo":     "дисциплинированным и целеустремлённым детям, которым важен порядок и постепенный рост уровня",
	"swimming": "спокойным детям, которым важно физическое развитие без сильной конкуренции",
	"art":      "усидчивым и творческим детям, которые любят создавать что-то своими руками",
	"theatre":  "эмоциональным и артистичным детям, которые не боятся быть в центре внимания",
}

// startQuiz обнуляет счёт предыдущего прохождения (если было, например,
// после /start в середине теста) и задаёт первый вопрос.
func (machine *Machine) startQuiz(ctx context.Context, userID int64) error {
	state := machine.store.GetOrCreate(userID)
	state.QuizIndex = 0
	state.ScoreFootball = 0
	state.ScoreJudo = 0
	state.ScoreSwimming = 0
	state.ScoreArt = 0
	state.ScoreTheatre = 0
	state.CurrentState = StateQuiz
	machine.store.Set(state)

	return machine.askQuizQuestion(ctx, userID, 0)
}

// askQuizQuestion отправляет вопрос теста с номером index (0-based)
// и клавиатурой из его вариантов ответа.
func (machine *Machine) askQuizQuestion(ctx context.Context, userID int64, index int) error {
	question := quizQuestions[index]
	buttons := buildQuizKeyboard(question)
	text := fmt.Sprintf("Вопрос %d из %d.\n%s", index+1, len(quizQuestions), question.Text)

	if err := machine.maxClient.SendMessageWithKeyboard(ctx, userID, text, buttons); err != nil {
		return err
	}

	state := machine.store.GetOrCreate(userID)
	state.QuizIndex = index
	state.CurrentState = StateQuiz
	machine.store.Set(state)
	return nil
}

// buildQuizKeyboard формирует клавиатуру вопроса: один вариант ответа —
// одна строка, потому что формулировки вариантов длиннее одного слова.
func buildQuizKeyboard(question quizQuestion) [][]maxapi.Button {
	buttons := make([][]maxapi.Button, 0, len(question.Options))
	for _, option := range question.Options {
		buttons = append(buttons, []maxapi.Button{{
			Text:    option.Text,
			Payload: payloadPrefixQuiz + option.Direction,
		}})
	}
	return buttons
}

// processQuizAnswer засчитывает балл направлению, к которому относится
// выбранный вариант, и либо задаёт следующий вопрос, либо, если это был
// последний, завершает тест подбором результата.
func (machine *Machine) processQuizAnswer(ctx context.Context, userID int64, direction string) error {
	state := machine.store.GetOrCreate(userID)

	if !isKnownSportKey(direction) {
		machine.logger.WarnContext(ctx, "unknown quiz answer payload", "direction", direction)
		return machine.askQuizQuestion(ctx, userID, state.QuizIndex)
	}

	state = addQuizPoint(state, direction)
	machine.store.Set(state)

	nextIndex := state.QuizIndex + 1
	if nextIndex >= len(quizQuestions) {
		return machine.finishQuiz(ctx, userID)
	}
	return machine.askQuizQuestion(ctx, userID, nextIndex)
}

// finishQuiz определяет направление(я)-победители по числу набранных
// баллов. Если оно одно — сразу передаёт управление подбору секций
// (state.Sport = направление, дальше как в обычном сценарии).
// Если несколько направлений набрали максимум баллов поровну — по методике
// теста показываются все они сразу, без дополнительных вопросов.
func (machine *Machine) finishQuiz(ctx context.Context, userID int64) error {
	state := machine.store.GetOrCreate(userID)
	winners := quizWinners(state)

	if len(winners) == 1 {
		state.Sport = winners[0]
		machine.store.Set(state)
		return machine.showResults(ctx, userID)
	}

	return machine.finishQuizTie(ctx, userID, state, winners)
}

// finishQuizTie отправляет сообщение со всеми направлениями, набравшими
// максимум баллов поровну: короткое пояснение "почему подходит" и кнопка
// со ссылкой на сайт для каждого — родитель выбирает, что ближе.
func (machine *Machine) finishQuizTie(ctx context.Context, userID int64, state UserState, winners []string) error {
	var text strings.Builder
	text.WriteString("Тест одинаково хорошо показал сразу несколько направлений:\n\n")
	buttons := make([][]maxapi.Button, 0, len(winners))
	for _, direction := range winners {
		text.WriteString("• ")
		text.WriteString(sportDisplayName(direction))
		if reason, ok := directionReasons[direction]; ok {
			text.WriteString(" — подходит ")
			text.WriteString(reason)
		}
		text.WriteString("\n")
		buttons = append(buttons, []maxapi.Button{{
			Text:    "Секции: " + sportDisplayName(direction),
			LinkURL: machine.buildSiteURL(direction, state),
		}})
	}
	text.WriteString("\nВыберите, что больше откликается ребёнку, и переходите к подбору конкретной секции.")

	if err := machine.maxClient.SendMessageWithKeyboard(ctx, userID, text.String(), buttons); err != nil {
		return err
	}

	completed := machine.store.GetOrCreate(userID)
	completed.CurrentState = StateCompleted
	machine.store.Set(completed)
	return nil
}

// addQuizPoint начисляет один балл направлению ответа. Неизвестный ключ
// молча игнорируется — вызывающий код (processQuizAnswer) уже проверяет
// direction через isKnownSportKey до вызова.
func addQuizPoint(state UserState, direction string) UserState {
	switch direction {
	case "football":
		state.ScoreFootball++
	case "judo":
		state.ScoreJudo++
	case "swimming":
		state.ScoreSwimming++
	case "art":
		state.ScoreArt++
	case "theatre":
		state.ScoreTheatre++
	}
	return state
}

// quizScore возвращает текущий счёт направления.
func quizScore(state UserState, direction string) int {
	switch direction {
	case "football":
		return state.ScoreFootball
	case "judo":
		return state.ScoreJudo
	case "swimming":
		return state.ScoreSwimming
	case "art":
		return state.ScoreArt
	case "theatre":
		return state.ScoreTheatre
	default:
		return 0
	}
}

// quizWinners возвращает направления с максимальным числом баллов.
// Обычно одно, но при равенстве — несколько, в порядке quizDirectionOrder.
func quizWinners(state UserState) []string {
	best := -1
	for _, direction := range quizDirectionOrder {
		if score := quizScore(state, direction); score > best {
			best = score
		}
	}

	winners := make([]string, 0, len(quizDirectionOrder))
	for _, direction := range quizDirectionOrder {
		if quizScore(state, direction) == best {
			winners = append(winners, direction)
		}
	}
	return winners
}
