package main

import (
    "fmt"
    "os"
    "os/exec"
    "github.com/godbus/dbus/v5"
    "flag"
)

// this program fixes the problem of not being able to toggle nightlight on permanently

// references
// https://gitlab.gnome.org/GNOME/gnome-settings-daemon/-/blob/master/data/org.gnome.settings-daemon.plugins.color.gschema.xml.in
// https://gitlab.gnome.org/GNOME/gnome-settings-daemon/-/blob/master/plugins/color/gsd-color-manager.c
// https://github.com/godbus/dbus/blob/76236955d466b078d82dcb16b7cf1dcf40ac25df/_examples/mediakeys.go#L26
// https://dbus.freedesktop.org/doc/dbus-specification.html#basic-types


const ColorInterface = "org.gnome.SettingsDaemon.Color"
const ColorPath = "/org/gnome/SettingsDaemon/Color"
const ColorPlugin = "org.gnome.settings-daemon.plugins.color"

// change these constants to match your taste
const (
    WARMEST uint32 = 1000
    WARMER uint32 = 2000
    WARM uint32 = 3000
    ON uint32 = 4000
    OFF uint32 = 6500
)

func set_gsd_property (name string, value string) error {
    cmd := exec.Command("gsettings", "set", ColorPlugin, name, value)
    _, err := cmd.Output()

    if err != nil {
        return err
    }

    return nil
}

func turn_nightlight_on_permanently() error {
    err := set_gsd_property("night-light-schedule-from", "0")
    if err != nil {
        return err
    }

    err = set_gsd_property("night-light-schedule-to", "24")
    if err != nil {
        return err
    }

    err = set_gsd_property("night-light-schedule-automatic", "false")
    if err != nil {
        return err
    }

    return nil
}

func get_current_temp (bus dbus.BusObject) (uint32, error) {
    res, err := bus.GetProperty(ColorInterface + ".Temperature")

    if err != nil {
        return 0, err
    }

    temp := res.Value().(uint32)
    return temp, nil
}

func set_current_temp (bus dbus.BusObject, temp uint32) error {
    val := dbus.MakeVariant(temp)
    err := bus.SetProperty(ColorInterface + ".Temperature", val)

    if err != nil {
        return err
    }

    return nil
}

func main() {
    // setup flags
    level_ptr := flag.String("level", "on", "the level of warmth: off, on " +
        "(default night light level), warm, warmer, or warmest")
    flag.Parse()

    var level uint32
    switch *level_ptr {
    case "off": level = OFF
    case "on": level = ON
    case "warm": level = WARM
    case "warmer": level = WARMER
    case "warmest": level = WARMEST
    default:
        // error
        fmt.Fprintln(os.Stderr, "Invalid level: must be one of off, on, " +
            "warm, warmer, or warmest.");
        os.Exit(1)
    }

    conn, err := dbus.ConnectSessionBus()
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to connect to session bus: ", err);
    }
    defer conn.Close()

    bus := conn.Object(ColorInterface, ColorPath)
    curr_temp, err := get_current_temp(bus)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Couldn't get temperature: ", err)
    } else {
        fmt.Println("Old temperature: ", curr_temp)
    }

    err = set_current_temp(bus, level)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Couldn't set temperature: ", err)
    }

    curr_temp, err = get_current_temp(bus)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Couldn't get temperature: ", err)
    } else {
        fmt.Println("New temperature: ", curr_temp)
    }

    turn_nightlight_on_permanently()
}


