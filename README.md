Необходимо написать in-memory кэш, который будет по ключу (uuid пользователя) возвращать профиль и список его заказов.
1. У кэша должен быть TTL (2 сек)
2. Кэшем может пользоваться функция(-и), которая работает с заказами (добавляет/обновляет/удаляет). Если TTL истек, то возвращается nil. При апдейте TTL снова устанавливается 2 сек. Методы должны быть потокобезопасными
3. Автоматическая очистка истекших записей

````
type Profile struct {
UUID string
Name string
Orders []*Order
}

type Order struct {
UUID string
Value any
CreatedAt time.Time
UpdatedAt time.Time
}
````