Функция `grabPresents` решает задачу о рюкзаке (или задачу о 0/1 рюкзаке), где цель — выбрать набор подарков с максимальной суммой ценностей, не превышая заданную вместимость рюкзака. Вот как работает эта функция:

### Общее Описание

Функция принимает два аргумента:
- `presents` — срез структур `Present`, представляющий подарки, каждый из которых имеет значение `Value` и размер `Size`.
- `capacity` — вместимость рюкзака.

Функция возвращает срез подарков, которые следует выбрать, чтобы максимизировать суммарное значение, не превышая при этом вместимость рюкзака.

### Работа Функции

1. **Инициализация Динамической Таблицы**:
   ```go
   n := len(presents)
   dp := make([][]int, n+1)
   for i := range dp {
       dp[i] = make([]int, capacity+1)
   }
   ```
   - Создается двумерный срез `dp` размером `(n+1) x (capacity+1)`, где `dp[i][w]` будет содержать максимальное значение, которое можно получить, используя первые `i` подарков при вместимости `w`.

2. **Заполнение Таблицы**:
   ```go
   for i := 1; i <= n; i++ {
       for w := 0; w <= capacity; w++ {
           if presents[i-1].Size <= w {
               dp[i][w] = max(dp[i-1][w], dp[i-1][w-presents[i-1].Size]+presents[i-1].Value)
           } else {
               dp[i][w] = dp[i-1][w]
           }
       }
   }
   ```
   - Перебираем все подарки (от 1 до `n`) и все возможные размеры рюкзака (от 0 до `capacity`).
   - Если текущий подарок можно поместить в рюкзак (`presents[i-1].Size <= w`), то выбираем максимальное значение между:
     - Неиспользованием текущего подарка (значение из предыдущего состояния `dp[i-1][w]`).
     - Использованием текущего подарка (значение из предыдущего состояния `dp[i-1][w-presents[i-1].Size]` плюс ценность текущего подарка).
   - Если текущий подарок больше, чем оставшееся место в рюкзаке (`presents[i-1].Size > w`), просто копируем значение из предыдущего состояния `dp[i-1][w]`.

3. **Обратный Ход для Определения Подарков**:
   ```go
   result := []Present{}
   w := capacity
   for i := n; i > 0 && w > 0; i-- {
       if dp[i][w] != dp[i-1][w] {
           result = append(result, presents[i-1])
           w -= presents[i-1].Size
       }
   }
   ```
   - Начинаем с последнего подарка и последнего возможного размера рюкзака.
   - Если значение в `dp[i][w]` отличается от значения `dp[i-1][w]`, это означает, что текущий подарок был выбран (так как его добавление привело к увеличению значения).
   - Добавляем текущий подарок в результат и уменьшаем оставшееся место в рюкзаке на размер текущего подарка.

4. **Возврат Результата**:
   ```go
   return result
   ```
   - Возвращается список подарков, которые обеспечивают максимальное значение при данной вместимости рюкзака.

### Пример

Допустим, у нас есть 3 подарка и вместимость рюкзака 5:

```go
  presents := []Present{
        {Value: 10, Size: 5},
        {Value: 40, Size: 4},
        {Value: 30, Size: 6},
        {Value: 50, Size: 3},
    }
    capacity := 10
```

Функция `grabPresents` выполнит следующие шаги:
1. Построит таблицу `dp`, чтобы определить максимальные значения для всех возможных размеров рюкзака.
2. Использует таблицу `dp`, чтобы определить, какие подарки следует выбрать, чтобы достичь максимальной ценности, не превышая вместимость рюкзака.

Если вы выполните функцию `grabPresents` с указанными данными, результатом будет список подарков, который обеспечивает максимальное значение в рюкзаке вместимостью 5.