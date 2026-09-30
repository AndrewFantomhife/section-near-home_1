package dialog

import "fmt"

// Служебные сообщения маршрутизации и fallback-сценариев.
const (
	messageNonTextOnly = "Я понимаю только текст и кнопки. Используйте /start, чтобы начать."

	messageDistrictButtonsOnly = "Пожалуйста, выберите район кнопками ниже."

	messageAgeButtonsOnly = "Пожалуйста, выберите возраст ребёнка кнопками ниже."

	messageSearchInProgress = "Подбор секций ещё идёт, пожалуйста, подождите несколько секунд."

	messageHelp = "Я помогу подобрать секцию для ребёнка: сначала узнаю район и возраст, " +
		"потом пройдём профориентационный тест из 8 вопросов — по его итогам порекомендую " +
		"направление (спорт или творчество) и покажу подходящие секции рядом с домом.\n\n" +
		"Доступные команды:\n" +
		"/start - начать диалог заново\n" +
		"/help - показать эту справку\n" +
		"/reset - сбросить текущий диалог"

	messageResetDone = "Диалог сброшен. Введите /start, чтобы начать заново."

	messageUnknownCommand = "Неизвестная команда. Доступны /start и /help."

	messageUnknownAction = "Неизвестное действие"

	messageInvalidAgeCallback = "Возраст должен быть от 6 до 14 лет"

	messageSearchFailed  = "Не удалось получить список секций. Попробуйте ещё раз через минуту."
	messageRestartButton = "Начать заново"
)

// Сообщения онбординга. Используются в onboarding.go.
const (
	messageWelcome = "Привет! Я помогу подобрать секцию для ребёнка.\n" +
		"Сначала уточню район и возраст, а затем пройдём короткий профориентационный тест " +
		"из 8 вопросов — по нему я порекомендую направление. Всего это займёт 3–5 минут.\n\n" +
		"Нажмите кнопку ниже, чтобы начать."

	messageWelcomeButton = "Начать подбор"

	messageAskDistrict = "В каком районе ищете секцию?\nВыберите район кнопками ниже."
)

// Сообщения результатов. Используются в results.go.
const (
	messageSearching = "Ищу подходящие секции рядом... это займёт несколько секунд."

	messageNoResults = "По вашему запросу секций не найдено.\n\n" +
		"Попробуйте изменить параметры командой /start."

	// messageDistrictFallbackNote предупреждает, что показанные секции —
	// не из района родителя: у демо-базы всего 6 секций (по одной на
	// направление), и точное совпадение район+направление есть не всегда.
	// Используется вместо тишины: родитель должен понимать, почему
	// адрес не рядом, а не решить, что бот ошибся.
	messageDistrictFallbackNote = "В районе «%s» подходящих по направлению секций пока нет — показываю ближайшие варианты в других районах:"

	messageResultsFooter = "Чтобы сравнить все секции района, откройте карту.\n" +
		"Для нового поиска введите /start."

	messageOpenMapButton = "Открыть карту"
)

// formatAskChildAge формирует вопрос о возрасте с подтверждением района.
func formatAskChildAge(district string) string {
	return fmt.Sprintf("Принял: %s.\nСколько лет ребёнку? Выберите кнопками ниже.", district)
}

// formatResultsHeader формирует заголовок списка найденных секций.
func formatResultsHeader(count int) string {
	return fmt.Sprintf("Нашёл %d подходящих секций:", count)
}

// formatDistrictFallbackNote формирует предупреждение о том, что в районе
// родителя нужного направления нет и показаны секции из других районов.
func formatDistrictFallbackNote(district string) string {
	return fmt.Sprintf(messageDistrictFallbackNote, district)
}

// formatQuizButtonsOnly формирует подсказку для текстовых сообщений,
// присланных во время прохождения теста, с указанием текущего шага.
func formatQuizButtonsOnly(step, total int) string {
	return fmt.Sprintf("Пожалуйста, отвечайте кнопками ниже (вопрос %d из %d).", step, total)
}

// formatQuizRecommendation формирует строку с рекомендацией теста и
// коротким пояснением "почему подходит" — подставляется перед списком
// секций. Для пустого или неизвестного направления возвращает "".
func formatQuizRecommendation(direction string) string {
	reason, ok := directionReasons[direction]
	if !ok {
		return ""
	}
	return fmt.Sprintf("По тесту ребёнку подходит: %s — это направление хорошо подходит %s.",
		sportDisplayName(direction), reason)
}
