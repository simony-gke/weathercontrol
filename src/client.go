package main

import (
    "context"
    "fmt"
    "log"

    "google.golang.org/grpc"
    pb "github.com/simony-gke/weathercontrol/proto"
)

const (
    port = ":50051"
)

func main() {
    conn, err := grpc.Dial(port, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("failed to connect to server: %v", err)
    }
    defer conn.Close()

    client := pb.NewWeatherControlServiceClient(conn)

    // Get current weather
    currentWeather, err := client.GetWeather(context.Background(), &pb.GetWeatherRequest{})
    if err != nil {
        log.Fatalf("failed to get weather: %v", err)
    }
    fmt.Printf("Current weather: %v", currentWeather)

    // Set the weather to rainy
    rainyWeather := pb.SetWeatherRequest{
        WeatherType: "Rainy",
        Intensity: 2,
    }
    setWeatherResponse, err := client.SetWeather(context.Background(), &rainyWeather)
    if err != nil {
        log.Fatalf("failed to set weather: %v", err)
    }
    fmt.Printf("Set weather response: %v", setWeatherResponse)
}
