package dialog

import (
	"context"
	"strconv"

	"section-near-home/backend/internal/clients/maxapi"
)

// sportOption связывает отображаемое название кнопки с ключом секции в db.json.
type sportOption struct {
	DisplayName string
	CatalogKey  string
}

// sportOptions является каталогом видов спорта для клавиатуры выбора.
// CatalogKey обязан совпадать со значением поля sport в mock-api/db.json:
// в текущей базе представлено 5 видов (6 секций), список ниже отражает
// именно их, иначе часть кнопок всегда возвращала бы «секций не найдено».
// При пополнении db.json новыми видами спорта — дополнить и этот список.
var sportOptions = []sportOption{
	{DisplayName: "Плавание", CatalogKey: "swimming"},
	{DisplayName: "Футбол", CatalogKey: "football"},
	{DisplayName: "Дзюдо", CatalogKey: "judo"},
	{DisplayName: "Рисование", CatalogKey: "art"},
	{DisplayName: "Театр", CatalogKey: "theatre"},
}

// availableDistricts содержит районы пилотного запуска.
// Список обязан совпадать со значениями поля district в mock-api/db.json.
var availableDistricts = []string{"Ленинский", "Октябрьский", "Фрунзенский"}

func (machine *Machine) startOnboarding(ctx context.Context, userID int64) error {
	buttons := buildDistrictKeyboard()

	if err := machine.maxClient.SendMessageWithKeyboard(ctx, userID, messageWelcome, buttons); err != nil {
		return err
	}

	state := machine.store.GetOrCreate(userID)
	state.District = ""
	state.ChildAge = 0
	state.Sport = ""
	state.QuizIndex = 0
	state.ScoreFootball = 0
	state.ScoreJudo = 0
	state.ScoreSwimming = 0
	state.ScoreArt = 0
	state.ScoreTheatre = 0
	state.CurrentState = StateAskDistrict
	machine.store.Set(state)
	return nil
}

func buildDistrictKeyboard() [][]maxapi.Button {
	buttons := make([][]maxapi.Button, 0, len(availableDistricts))
	for _, district := range availableDistricts {
		buttons = append(buttons, []maxapi.Button{{
			Text:    district,
			Payload: payloadPrefixDistrict + district,
		}})
	}
	return buttons
}

// processDistrict валидирует выбранный район и переходит к шагу возраста.
// Район вне белого списка считается подделанным payload: клавиатура повторяется.
func (machine *Machine) processDistrict(ctx context.Context, userID int64, district string) error {
	if !isKnownDistrict(district) {
		machine.logger.WarnContext(ctx, "unknown district payload", "district", district)
		return machine.startOnboarding(ctx, userID)
	}

	state := machine.store.GetOrCreate(userID)
	state.District = district
	state.CurrentState = StateAskChildAge
	machine.store.Set(state)

	return machine.askChildAge(ctx, userID)
}

// askChildAge запрашивает возраст ребёнка клавиатурой в две строки: 6-10 и 11-14.
func (machine *Machine) askChildAge(ctx context.Context, userID int64) error {
	state := machine.store.GetOrCreate(userID)
	buttons := buildAgeKeyboard()

	if err := machine.maxClient.SendMessageWithKeyboard(
		ctx, userID, formatAskChildAge(state.District), buttons); err != nil {
		return err
	}

	state.CurrentState = StateAskChildAge
	machine.store.Set(state)
	return nil
}

// processChildAge сохраняет возраст и запускает профориентационный тест,
// который и определит рекомендованное направление (см. quiz.go).
// Диапазон уже проверен в machine.go при разборе payload.
func (machine *Machine) processChildAge(ctx context.Context, userID int64, age int) error {
	state := machine.store.GetOrCreate(userID)
	state.ChildAge = age
	machine.store.Set(state)

	return machine.startQuiz(ctx, userID)
}

// isKnownDistrict возвращает true, если район присутствует в белом списке.
func isKnownDistrict(district string) bool {
	for _, knownDistrict := range availableDistricts {
		if knownDistrict == district {
			return true
		}
	}
	return false
}

// isKnownSportKey возвращает true, если ключ спорта присутствует в каталоге.
func isKnownSportKey(sportKey string) bool {
	_, found := findSportOption(sportKey)
	return found
}

// findSportOption ищет опцию каталога по ключу спорта.
func findSportOption(sportKey string) (sportOption, bool) {
	for _, option := range sportOptions {
		if option.CatalogKey == sportKey {
			return option, true
		}
	}
	return sportOption{}, false
}

// buildAgeKeyboard формирует клавиатуру возраста в две строки: 6-10 и 11-14.
func buildAgeKeyboard() [][]maxapi.Button {
	firstRow := make([]maxapi.Button, 0, 5)
	for age := minChildAge; age <= 10; age++ {
		firstRow = append(firstRow, maxapi.Button{
			Text:    strconv.Itoa(age),
			Payload: payloadPrefixAge + strconv.Itoa(age),
		})
	}

	secondRow := make([]maxapi.Button, 0, 4)
	for age := 11; age <= maxChildAge; age++ {
		secondRow = append(secondRow, maxapi.Button{
			Text:    strconv.Itoa(age),
			Payload: payloadPrefixAge + strconv.Itoa(age),
		})
	}

	return [][]maxapi.Button{firstRow, secondRow}
}
