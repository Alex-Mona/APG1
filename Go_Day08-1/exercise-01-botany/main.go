package main

import (
    "fmt"
    "reflect"
)

type UnknownPlant struct {
    FlowerType  string
    LeafType    string
    Color       int `color_scheme:"rgb"`
}

type AnotherUnknownPlant struct {
    FlowerColor int
    LeafType    string
    Height      int `unit:"inches"`
}

func describePlant(plant interface{}) {
    val := reflect.ValueOf(plant)
    typ := reflect.TypeOf(plant)

    for i := 0; i < val.NumField(); i++ {
        field := typ.Field(i)
        fieldValue := val.Field(i)

        // Считываем теги
        tag := field.Tag

        // Собираем имя поля и значение
        name := field.Name
        value := fieldValue.Interface()

        // Формируем строку для вывода с учетом тегов
        if colorTag, ok := tag.Lookup("color_scheme"); ok {
            fmt.Printf("%s(%s):%v\n", name, colorTag, value)
        } else if unitTag, ok := tag.Lookup("unit"); ok {
            fmt.Printf("%s(%s):%v\n", name, unitTag, value)
        } else {
            fmt.Printf("%s:%v\n", name, value)
        }
    }
}

func main() {
    plant1 := UnknownPlant{
        FlowerType: "rose",
        LeafType:   "oval",
        Color:      255,
    }

    plant2 := AnotherUnknownPlant{
        FlowerColor: 10,
        LeafType:    "lanceolate",
        Height:      15,
    }

    fmt.Println("Описание первого растения:")
    describePlant(plant1)

    fmt.Println("\nОписание второго растения:")
    describePlant(plant2)
}
