package conf

import "testing"

func TestConf_Parse(t *testing.T) {
	type fields struct {
		Server Server
		Logger Logger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"should not error", fields{Server: Server{Debug: true}}, false},
		{"should not error", fields{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conf{
				Server: tt.fields.Server,
				Logger: tt.fields.Logger,
			}
			if err := c.Parse(); (err != nil) != tt.wantErr {
				t.Errorf("Conf.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
