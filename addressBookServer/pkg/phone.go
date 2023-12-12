package pkg

import (
	"regexp"
	"strconv"
	"strings"
)

// Функция для приведения номера телефона к стандартному формату (e.g. 8(800)555-35-35 -> 78005553535)
func PhoneNormalize(phone string) (normalizedPhone string, err error) {
	myErr := NewMyError("pkg: PhoneNormalize(phone string)")
	var builder strings.Builder
	for i := range phone {
		l := phone[i]
		if l >= '0' && l <= '9' {
			builder.WriteByte(l)
		}
	}
	normalizedPhone = builder.String()
	// Проверка на пустой номер телефона
	if builder.Len() == 0 {
		return "", myErr.Wrap(nil, "Empty phone number")
	}
	if normalizedPhone[0] == '8' {
		normalizedPhone = "7" + normalizedPhone[1:]
	}

	if normalizedPhone[0] != '7' {
		return "", myErr.Wrap(nil, "Incorrect phone number: "+normalizedPhone)
	}
	if len(normalizedPhone) != 11 {
		return "", myErr.Wrap(nil, "Incorrect len of phone number: "+strconv.Itoa(len(normalizedPhone)))
	}
	return normalizedPhone, nil
}

// Ещё один вариант реализации данной функции
// С помощью регулярных выражений
func PhoneNormalize3(phone string) (normalizedPhone string, err error) {
	myErr := NewMyError("pkg: PhoneNormalize3(phone string)")
	normalizedPhone = regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	// Проверка на пустой номер телефона
	if normalizedPhone == "" {
		return "", myErr.Wrap(nil, "Empty phone number")
	}
	// Проверка на то, что телефон состоит из 10 цифр
	if matched, _ := regexp.MatchString("^\\d{10}$", normalizedPhone); matched {
		normalizedPhone = "7" + normalizedPhone
	}
	// Проверка на то, что первыя цифра 8
	if string(normalizedPhone[0]) == "8" {
		normalizedPhone = "7" + normalizedPhone[1:]
	}
	// Проверяем, корректность номера телефона (состоит только из цифр, имеет правильную длину и начинается с 7)
	if matched, _ := regexp.MatchString("^7\\d{10}$", normalizedPhone); !matched {
		return "", myErr.Wrap(nil, "Incorrect phone number: "+phone)
	}
	return normalizedPhone, nil
}
