package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	pb "github.com/simony-gke/weathercontrol/proto"
)

// Define the port
const (
	port = ":50051"
)

// weatherData stores the current weather information
type weatherData struct {
	weatherType string
	intensity   int32
}

// server is used to implement WeatherControlService
type server struct {
	mu         sync.RWMutex // Protects weather
	weather    weatherData
	pb.UnimplementedWeatherControlServiceServer
}

// GetWeather implements WeatherControlService.GetWeather
func (s *server) GetWeather(ctx context.Context, in *pb.GetWeatherRequest) (*pb.GetWeatherResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &pb.GetWeatherResponse{
		WeatherType: s.weather.weatherType,
		Intensity:   s.weather.intensity,
	}, nil
}

// SetWeather implements WeatherControlService.SetWeather
func (s *server) SetWeather(ctx context.Context, in *pb.SetWeatherRequest) (*pb.SetWeatherResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Basic validation (you can add more sophisticated checks here)
	if in.WeatherType == "" {
		return &pb.SetWeatherResponse{Success: false}, fmt.Errorf("weather type cannot be empty")
	}
	if in.Intensity < 0 {
		return &pb.SetWeatherResponse{Success: false}, fmt.Errorf("intensity cannot be negative")
	}

	s.weather.weatherType = in.WeatherType
	s.weather.intensity = in.Intensity
	return &pb.SetWeatherResponse{Success: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Initialize the server with default weather
	weatherServer := &server{
		weather: weatherData{
			weatherType: "Sunny",
			intensity:   1,
		},
	}

	pb.RegisterWeatherControlServiceServer(s, weatherServer)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
