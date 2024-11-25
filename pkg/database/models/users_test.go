package models

import (
	"Yakudza/pkg/config"
	"Yakudza/pkg/database"
	"testing"
	"time"
)

func TestUser_Create(t *testing.T) {
	cfg := config.MustLoadByPath("../../../config/local.yaml")
	database.Init(cfg)
	type fields struct {
		ID        uint64
		Login     string
		Password  string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Создание юзера",
			fields: fields{
				Login:    "Collapse",
				Password: "syncmaster12",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:        tt.fields.ID,
				Login:     tt.fields.Login,
				Password:  tt.fields.Password,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
			}
			if err := u.Create(); err != nil {
				t.Errorf("Create() error = %v", err)
			}

			t.Logf("Создали юзера")
		})
	}
}
