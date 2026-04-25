package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func verificarServicio(url, servicio string, wg *sync.WaitGroup, canal chan float64) {
	defer wg.Done()
	inicio := time.Now()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error en %s: %s\n", url, err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error en %s: %s\n", url, err)
		return
	}

	defer resp.Body.Close()

	duracion := time.Since(inicio).Seconds()
	fmt.Printf("Servicio: %s, URL: %s, Status: %s, Duración: %f\n", servicio, url, resp.Status, duracion)

	canal <- duracion

	tiempo := make(chan float64, 3)
	tiempo <- 0.45
	tiempo <- 0.54
	tiempo <- 0.32
}

func main() {
	configData, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Printf("Error al leer el archivo de configuración: %s\n", err)
		return
	}

	var config map[string]interface{}
	err = json.Unmarshal(configData, &config)
	if err != nil {
		fmt.Printf("Error al analizar el archivo de configuración: %s\n", err)
		return
	}

	limite := config["limite"].(float64)
	servicios := config["servicios"].(map[string]interface{})

	var wg sync.WaitGroup
	fmt.Println("Iniciando verificación de servicios...")

	canalTiempos := make(chan float64, 4)

	for servicio, url := range servicios {
		wg.Add(1)
		go verificarServicio(url.(string), servicio, &wg, canalTiempos)
	}

	fmt.Println("Verificación de servicios completada.")
	wg.Wait()
	close(canalTiempos)

	var Sumatotal float64

	for i := 0; i < len(servicios); i++ {
		tiempoReal := <-canalTiempos
		Sumatotal += tiempoReal
		if tiempoReal > limite {
			fmt.Printf("El servicio ha superado el límite de tiempo: %f segundos\n", tiempoReal)
		} else {
			fmt.Printf("Tiempo de respuesta recibido: %f segundos\n", tiempoReal)

		}
	}

	promedio := Sumatotal / 4
	fmt.Printf("Tiempo promedio de respuesta: %f segundos\n", promedio)

	//f, _ := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	//f.WriteString("El promedio fue: " + fmt.Sprintf("%f", promedio) + " segundos\n")
	//f.Close()
}
