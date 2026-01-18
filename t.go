package main

import (
    "fmt"
    "time"
)

func main() {

    go func() {
        for {
            fmt.Println("hi im running")
            time.Sleep(3 * time.Second)
        }
    }()

    select {}
}
