# APG1

## Go_Day00-1:
1. Реализация алгоритма Anscombe's quartet.
## Go_Day01-1: 
1. Чтение XML/JSON,
2. сравнение баз данных XML/JSON,
3. создание дампов файловых систем.
## Go_Day02-1: 
1. Реализация утилит find,
2. wc,
3. xargs,
4. инструмент ротации журналов(архивация).
## Go_Day03-1: 
1. Чтение базы данных с помощью Elasticsearch,
2. HTML UI для базы данных,
3. реализация обработчика API вместо простого HTML,
4. реализация функционала поиска трех ближайших ресторанов,
5. JWT.
## Go_Day04-1: 
1. Спецификация протокола между торговым автоматом и сервером,
2. аутентификации TLS Minica,
3. вернуть текст из C кода в ASCII-формате в качестве текста в поле "спасибо" в ответе при запросе через cURL.
## Go_Day05-1:
1. Функция areToysBalanced, которая будет получать указатель на корень дерева в качестве аргумента
2. Функция unrollGarland который также получает указатель на корневой узел
3. Функция getNCoolestPresents, которая, имея несортированный срез Presents и целое число n, вернет отсортированный срез (desc) "самых крутых" из списка.
4. Функция grabPresents, которая получает срез экземпляров Present и емкость вашего жесткого диска
## Go_Day06-1:
1. Генерация файла PNG размером 300x300 пикселей
2. Создание блога и спользованием баз данных PostreSQL, панель администратора, базовая поддержка разметки ### в сгенерированном HTML, пагинация
3. Реализация ограничения скорости, так что если более сотни клиентов в секунду пытаются получить к нему доступ, они должны получить ответ "429 Too Many Requests"
## Go_Day07-1:
1. Написать несколько тестов (в *_test.goфайлах) для кода данного в проекте, которые покажут, что он выдает неправильные результаты. Также нужно написать отдельную функцию (вы должны назвать ее minCoins2), которая будет иметь те же параметры, но будет успешно обрабатывать эти случаи
2. Тестирования кода на производительность используя встроенные инструменты Go - пакет pprof который поддерживает профилирование по CPU, памяти, блокировкам (горутины) и другим параметрам.
3. Генерация документации godoc
## Go_Day08-1:
1. Использование указателей и пакет unsafe, что позволяет работать с памятью на более низком уровне, как в C.
2. Функция describePlant, которая будет принимать любой вид растений в рамках задания(она должна работать со структурами разных типов), а затем вывести все поля как пары ключ-значение, разделенные запятой
3. Создать пустое окно Mac OS GUI/Linux по умолчанию (размером 300x200) с заголовком "Школа 21". 
## Go_Day09-1:
1. Функция sleepSort, которая принимает несортированный срез целых чисел и возвращает целочисленный канал, где эти числа будут записаны одно за другим в отсортированном порядке
2. Функция crawlWeb которая будет принимать входной канал (для отправки URL) и возвращать другой канал для результатов сканирования (указателей на тело веб-страницы в виде строки). Кроме того, в любой момент времени не должно быть более 8 горутин, опрашивающих страницы параллельно.
3. Функцию multiplex, которая должна быть вариативной (принимать переменное количество аргументов). Она должна принимать каналы ( chan interface{}) в качестве аргументов и возвращать один канал того же типа.
## Go_Team00-2:
1. Реализация gRPC сервер, который генерирует и отправляет случайные данные частот через потоковый сервис. Он включает генерацию случайных данных, настройку gRPC сервера и логирование ошибок.
2. Реализация клиент на Go подключается к gRPC серверу и обрабатывает поток данных частот для обнаружения аномалий
3. Проект демонстрирует использование gRPC для стриминга данных частот и их обработки в тестовой среде. Проект включает реализацию сервера gRPC, клиента и тестов, которые проверяют основную логику работы с данными частот, их запись в базу данных SQLite и проверку на аномалии.
## Go_Team01-1/warehouse-cli:
1. Итак, здесь нам нужно реализовать две программы - одну, которая будет клиентом, и одну, которая будет экземпляром базы данных. Всякий раз, когда вы запускаете новый экземпляр, вы должны иметь возможность указать ему на существующий экземпляр, поэтому после получения heartbeat он отправит свой хост и порт всем остальным работающим узлам, и все будут знать нового парня.
Если узел экземпляра запущен с фактором репликации, отличным от существующих узлов, он должен это обнаружить и автоматически выйти из строя без присоединения к кластеру. Это означает, что фактор репликации, вероятно, также должен быть включен в heartbeat.
2. Три операции — GET, SET и DELETE, использование строк UUID4 в качестве ключей артефактов
3. Обновление логики из задач 00/01. Теперь мы вводим концепции узлов Лидер и Последователь
