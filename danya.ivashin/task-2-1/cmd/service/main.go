package main

import (
	"fmt"

	"github.com/Danya-byte/task-2-1/internal/department"
)

func scanWorkerRequest() (*department.ChangeRequest, error) {
	req := new(department.ChangeRequest)
	if _, err := fmt.Scan(&req.Operator); err != nil {
		return nil, fmt.Errorf("scan operator: %w", err)
	}

	if _, err := fmt.Scan(&req.Temperature); err != nil {
		return nil, fmt.Errorf("scan temperature: %w", err)
	}

	return req, nil
}

func processDepartment() error {
	var countWorkers int
	if _, err := fmt.Scan(&countWorkers); err != nil {
		return fmt.Errorf("scan count workers: %w", err)
	}

	dep := department.New()

	for range countWorkers {
		req, err := scanWorkerRequest()
		if err != nil {
			return fmt.Errorf("process worker: %w", err)
		}

		optimum, err := dep.Recalculate(req)
		if err != nil {
			return fmt.Errorf("recalculate temperature: %w", err)
		}

		fmt.Println(optimum)
	}

	return nil
}

func main() {
	var countDepartments int
	if _, err := fmt.Scan(&countDepartments); err != nil {
		fmt.Printf("scan count departments: %s\n", err)

		return
	}

	for range countDepartments {
		if err := processDepartment(); err != nil {
			fmt.Printf("process department: %s\n", err)

			return
		}
	}
}
