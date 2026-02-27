package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func SetConfig() {
	// echo decides whether path is printed to stdout
	val := os.Getenv("GOZELLE_ECHO")
	if val == "" {
		os.Setenv("GOZELLE_ECHO", "false")
	} else if val == "true" {
	} else if val != "false" && val != "true" {
		fmt.Println("GOZELLE_ECHO must be true or false")
		os.Setenv("GOZELLE_ECHO", "false")
	}

	// echo decides whether path is printed to stdout
	val = os.Getenv("GOZELLE_MINIMUM")
	if val == "" {
		os.Setenv("GOZELLE_MINIMUM_WEIGHT", "0.1")
	} else {
		_, err := strconv.ParseFloat(val, 64)
		if err != nil {
			fmt.Println("GOZELLE_MINIMUM_WEIGHT must be a valid float")
			os.Setenv("GOZELLE_ECHO", "0.1")
		} else {
			os.Setenv("GOZELLE_MINIMUM_WEIGHT", val)
		}
	}

	var filePath string
	// data_dir decides where the data is stored
	val = os.Getenv("GOZELLE_DATA_DIR")
	if val == "" {
		dataDir := os.Getenv("XDG_DATA_HOME")
		if dataDir == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Printf("could not get user home directory:")
			}
			dataDir = filepath.Join(homeDir, ".local", "share")
		}
		filePath := filepath.Join(dataDir, "gozelle", "db.gob")
		os.Setenv("GOZELLE_DATA_DIR", filePath)
	} else if val != filePath {
		// check if the directory exists
		if _, err := os.Stat(val); os.IsNotExist(err) {
			// create the directory
			fmt.Println("Creating data directory:", val)
			err := os.MkdirAll(val, os.ModePerm)
			if err != nil {
				fmt.Println("Error creating directory:", err)
				os.Exit(1)
			}
		}
	}
}
