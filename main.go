package main

import (
    "fmt"
    "os"
    "github.com/godbus/dbus/v5"
)

const ColorInterface = "org.gnome.SettingsDaemon.Color"
const ColorPath = "/org/gnome/SettingsDaemon/Color"

func get_current_temp (bus dbus.BusObject) (uint32, error) {
    res, err := bus.GetProperty(ColorInterface + ".Temperature")

    if err != nil {
        return 0, err
    }

    temp := res.Value().(uint32)
    return temp, nil
}

func main() {
    conn, err := dbus.ConnectSessionBus()
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to connect to session bus: ", err);
    }
    defer conn.Close()

    // var curr_temp uint32
    bus := conn.Object(ColorInterface, ColorPath)
    curr_temp, err := get_current_temp(bus)
    if err != nil {
        fmt.Fprintln(os.Stderr, "error: ", err)
    }
    fmt.Println(curr_temp)
}


