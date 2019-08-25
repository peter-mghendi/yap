package models

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func Test_parseDBURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name        string
		args        args
		wantDialect string
		wantURL     string
		wantErr     bool
	}{
		{
			name:        "DB URL Parse Test 1",
			args:        args{url: "postgres://username:password@data.hostname.co.ke:5234/database"},
			wantDialect: "postgres",
			wantURL:     "dbname=database host=data.hostname.co.ke password=password port=5234 user=username sslmode=disable",
			wantErr:     false,
		},
		{
			name:        "DB URL Parse Test 2",
			args:        args{url: "postgres://username:password@data.hostname.com:5432/mydatabase"},
			wantDialect: "postgres",
			wantURL:     "dbname=mydatabase host=data.hostname.com password=password port=5432 user=username sslmode=disable",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseDBURL(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDBURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantDialect {
				t.Errorf("parseDBURL() got = %v, want %v", got, tt.wantDialect)
			}
			if got1 != tt.wantURL {
				t.Errorf("parseDBURL() got1 = %v, want %v", got1, tt.wantURL)
			}
		})
	}
}

func TestInitDB(t *testing.T) {
	if e := godotenv.Load("../.env"); e != nil {
		t.Error(e)
	}

	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "FAKE DB Connect Test",
			args:    args{url: "postgres://username:password@data.hostname.com:5432/mydatabase"},
			wantErr: true,
		},
		{
			name:    "REAL DB Connect Test",
			args:    args{url: os.Getenv("DATABASE_URL")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitDB(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("InitDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
