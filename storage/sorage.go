package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"phishing-trainer/models"
	"sync"
)

var (
	mu sync.RWMutex
)

// SaveSimulationResults сохраняет все записи simulation_results для указанного сайта
func SaveSimulationResults(site string, results []models.SimulationResult) error {
	mu.Lock()
	defer mu.Unlock()
	file := filepath.Join("storage", "data", site+"_sim.json")
	return writeJSON(file, results)
}

// LoadSimulationResults загружает записи
func LoadSimulationResults(site string) ([]models.SimulationResult, error) {
	mu.RLock()
	defer mu.RUnlock()
	file := filepath.Join("storage", "data", site+"_sim.json")
	var results []models.SimulationResult
	err := readJSON(file, &results)
	if os.IsNotExist(err) {
		return []models.SimulationResult{}, nil
	}
	return results, err
}

// SaveAVWarningStats сохраняет статистику предупреждений
func SaveAVWarningStats(site string, stats []models.AVWarningStat) error {
	mu.Lock()
	defer mu.Unlock()
	file := filepath.Join("storage", "data", site+"_av.json")
	return writeJSON(file, stats)
}

// LoadAVWarningStats загружает статистику предупреждений
func LoadAVWarningStats(site string) ([]models.AVWarningStat, error) {
	mu.RLock()
	defer mu.RUnlock()
	file := filepath.Join("storage", "data", site+"_av.json")
	var stats []models.AVWarningStat
	err := readJSON(file, &stats)
	if os.IsNotExist(err) {
		return []models.AVWarningStat{}, nil
	}
	return stats, err
}

// AppendSimulationResult добавляет одну запись к существующим
func AppendSimulationResult(site string, result models.SimulationResult) error {
	results, err := LoadSimulationResults(site)
	if err != nil {
		return err
	}
	results = append(results, result)
	return SaveSimulationResults(site, results)
}

// AppendAVWarningStat добавляет одну запись предупреждения
func AppendAVWarningStat(site string, stat models.AVWarningStat) error {
	stats, err := LoadAVWarningStats(site)
	if err != nil {
		return err
	}
	stats = append(stats, stat)
	return SaveAVWarningStats(site, stats)
}

// UpdateSimulationResult обновляет существующую запись по visit_id
func UpdateSimulationResult(site string, visitID string, updateFn func(*models.SimulationResult)) error {
	results, err := LoadSimulationResults(site)
	if err != nil {
		return err
	}
	for i, r := range results {
		if r.VisitID == visitID {
			updateFn(&results[i])
			return SaveSimulationResults(site, results)
		}
	}
	return fmt.Errorf("visit_id not found")
}

// UpdateAVWarningStat обновляет запись предупреждения по visit_id
func UpdateAVWarningStat(site string, visitID string, updateFn func(*models.AVWarningStat)) error {
	stats, err := LoadAVWarningStats(site)
	if err != nil {
		return err
	}
	for i, s := range stats {
		if s.VisitID == visitID {
			updateFn(&stats[i])
			return SaveAVWarningStats(site, stats)
		}
	}
	return fmt.Errorf("visit_id not found")
}

// writeJSON записывает данные в файл (создавая папку при необходимости)
func writeJSON(file string, v interface{}) error {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0644)
}

// readJSON читает данные из файла
func readJSON(file string, v interface{}) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
