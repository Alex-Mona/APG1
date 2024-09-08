Для создания оконного приложения на Linux с использованием **GTK** в Go можно использовать **Cgo**, чтобы напрямую вызывать функции GTK из Go-кода. Делаю такую реализацию, т.к MacOs под рукой нет, и нужно это исключительно ради проверки.
### Шаги реализации с использованием GTK:

1. Убедитесь, что у вас установлены все зависимости для работы с **GTK**.

Для Ubuntu:
```bash
sudo apt-get install libgtk-3-dev
```

2. Создайте файл `main.go`, который будет использовать **GTK** для создания окна с заголовком "School 21" и размером 300x200.

### Шаг 1: Написание кода

Создайте файл `main.go` с следующим содержимым:

```go
package main

/*
#cgo pkg-config: gtk+-3.0
#include <gtk/gtk.h>

// Функция для создания и отображения окна
static void create_window() {
    GtkWidget *window;

    // Инициализация GTK
    gtk_init(NULL, NULL);

    // Создаем новое окно
    window = gtk_window_new(GTK_WINDOW_TOPLEVEL);

    // Устанавливаем заголовок окна
    gtk_window_set_title(GTK_WINDOW(window), "School 21");

    // Устанавливаем размер окна 300x200
    gtk_window_set_default_size(GTK_WINDOW(window), 300, 200);

    // Подключаем сигнал закрытия окна
    g_signal_connect(window, "destroy", G_CALLBACK(gtk_main_quit), NULL);

    // Отображаем окно
    gtk_widget_show_all(window);

    // Запуск основного цикла GTK
    gtk_main();
}
*/
import "C"

func main() {
    // Вызываем функцию для создания и отображения окна
    C.create_window()
}
```

### Шаг 2: Компиляция и запуск

1. Компиляция программы:

   ```bash
   go build -o GtkApp
   ```

   Эта команда компилирует программу, используя GTK-библиотеки.

2. Запустите скомпилированное приложение:

   ```bash
   ./GtkApp
   ```

   Программа создаст окно размером 300x200 с заголовком "School 21".

### Объяснение кода:

1. **`gtk_init(NULL, NULL)`** — инициализирует библиотеку GTK.
2. **`gtk_window_new(GTK_WINDOW_TOPLEVEL)`** — создаёт новое главное окно.
3. **`gtk_window_set_title`** — устанавливает заголовок окна "School 21".
4. **`gtk_window_set_default_size`** — устанавливает размеры окна 300x200.
5. **`g_signal_connect(window, "destroy", G_CALLBACK(gtk_main_quit), NULL)`** — подключает обработчик для закрытия окна.
6. **`gtk_widget_show_all`** — отображает все виджеты в окне.
7. **`gtk_main()`** — запускает основной цикл обработки событий GTK.

### Заключение

Это примерно то, что вышло бы на MacOs. Вот такой вот прикол... 08 задания. Что делать тем у кого MacOs нету? Не знаю... Делать вот таким образом наверное...