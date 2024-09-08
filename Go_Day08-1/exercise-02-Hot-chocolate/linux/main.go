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
