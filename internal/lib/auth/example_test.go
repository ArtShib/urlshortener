package auth

import (
	"fmt"
	"log"
)

// Example сценарий использования сервиса:
// создание сервиса, генерация ID, создание токена и его проверка.
func Example() {
	// 1. Инициализируем сервис с секретным ключом
	svc := NewAuthService("my-secret-key")

	// 2. Генерируем новый ID пользователя
	userID, err := svc.GenerateUserID()
	if err != nil {
		log.Fatal(err)
	}

	// 3. Создаем токен для пользователя
	token := svc.CreateToken(userID)

	// 4. Проверяем валидность токена
	// Если токен валиден, мы можем доверять извлеченному из него ID
	if svc.ValidateToken(token) {
		extractedID := svc.GetUserID(token)

		// Проверяем, совпадает ли ID (в реальном коде это очевидно, здесь для демонстрации)
		if extractedID == userID {
			fmt.Println("Authentication workflow successful")
		}
	}

	// Output:
	// Authentication workflow successful
}

// ExampleService_CreateToken показывает, как создается токен.
// Для предсказуемости вывода здесь используются фиксированные данные.
func ExampleService_CreateToken() {
	// Используем простой секрет "secret" для примера
	svc := NewAuthService("secret")

	// Фиксированный ID пользователя (16 байт в hex)
	userID := "0102030405060708090a0b0c0d0e0f10"

	token := svc.CreateToken(userID)
	fmt.Println(token)

	// Output:
	// 0102030405060708090a0b0c0d0e0f1075b0a11c6b621648933b0aceb3a50c2552d2bbf46b473eeff0c809debade5149
}

// ExampleService_ValidateToken показывает разницу между валидным и невалидным токеном.
func ExampleService_ValidateToken() {
	svc := NewAuthService("secret")

	// Токен, созданный с секретом "secret" и ID "0102...10"
	validToken := "0102030405060708090a0b0c0d0e0f1075b0a11c6b621648933b0aceb3a50c2552d2bbf46b473eeff0c809debade5149"

	// Токен с поврежденной подписью (последний символ изменен с '9' на '0')
	invalidToken := "0102030405060708090a0b0c0d0e0f1075b0a11c6b621648933b0aceb3a50c2552d2bbf46b473eeff0c809debade5140"

	fmt.Printf("Valid token: %v\n", svc.ValidateToken(validToken))
	fmt.Printf("Invalid token: %v\n", svc.ValidateToken(invalidToken))

	// Output:
	// Valid token: true
	// Invalid token: false
}
