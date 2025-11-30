package services

/*
Функции:
	1. Создание pr
	2. Merge
	3. Переназначение пользоватля
	4. По пользователю найти Ревью

Основная сложность в написании сервиса была связана с возможным рейс кондишн.
Было исправлено за счет транзакций
*/

type XService struct {
}

func NewPullRequestService() *XService {
	return &XService{}
}
