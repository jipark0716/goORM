package main

import (
    "model"
    "fmt"
)

func main() {
    apps := []model.Example{}
    query := model.NewQuery()
    query.Get(&apps)
    fmt.Printf("%#v", apps)
}
