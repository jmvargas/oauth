package aps

import (
	"reflect"
	"testing"
	"time"

	"github.com/markbates/goth"
)

func TestSession_GetAuthURL(t *testing.T) {
	type fields struct {
		AuthURL string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:    "empty",
			wantErr: true,
		},
		{
			name:    "with value",
			wantErr: false,
			fields: fields{
				AuthURL: "myurl",
			},
			want: "myurl",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Session{
				AuthURL: tt.fields.AuthURL,
			}
			got, err := s.GetAuthURL()
			if (err != nil) != tt.wantErr {
				t.Errorf("Session.GetAuthURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Session.GetAuthURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_Marshal(t *testing.T) {
	type fields struct {
		AuthURL      string
		AccessToken  string
		RefreshToken string
		ExpiresAt    time.Time
	}
	aTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-05-16T08:28:06.801Z")

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple",
			fields: fields{
				AuthURL:      "myurl",
				AccessToken:  "token",
				RefreshToken: "rtoken",
				ExpiresAt:    aTime,
			},
			want: `{"AuthURL":"myurl","AccessToken":"token","RefreshToken":"rtoken","ExpiresAt":"2014-05-16T08:28:06.801Z"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Session{
				AuthURL:      tt.fields.AuthURL,
				AccessToken:  tt.fields.AccessToken,
				RefreshToken: tt.fields.RefreshToken,
				ExpiresAt:    tt.fields.ExpiresAt,
			}
			if got := s.Marshal(); got != tt.want {
				t.Errorf("Session.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_String(t *testing.T) {
	type fields struct {
		AuthURL      string
		AccessToken  string
		RefreshToken string
		ExpiresAt    time.Time
	}
	aTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-05-16T08:28:06.801Z")
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple",
			fields: fields{
				AuthURL:      "myurl",
				AccessToken:  "token",
				RefreshToken: "rtoken",
				ExpiresAt:    aTime,
			},
			want: `{"AuthURL":"myurl","AccessToken":"token","RefreshToken":"rtoken","ExpiresAt":"2014-05-16T08:28:06.801Z"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Session{
				AuthURL:      tt.fields.AuthURL,
				AccessToken:  tt.fields.AccessToken,
				RefreshToken: tt.fields.RefreshToken,
				ExpiresAt:    tt.fields.ExpiresAt,
			}
			if got := s.String(); got != tt.want {
				t.Errorf("Session.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProvider_UnmarshalSession(t *testing.T) {
	type args struct {
		data string
	}
	aTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-05-16T08:28:06.801Z")
	tests := []struct {
		name    string
		args    args
		want    goth.Session
		wantErr bool
	}{
		{
			name: "simple test",
			args: args{
				data: `{
					"AuthURL": "auth",
					"AccessToken": "b",
					"RefreshToken": "a",
					"ExpiresAt": "2014-05-16T08:28:06.801Z"
				}`,
			},
			wantErr: false,
			want: &Session{
				AuthURL:      "auth",
				AccessToken:  "b",
				RefreshToken: "a",
				ExpiresAt:    aTime,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{}
			got, err := p.UnmarshalSession(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.UnmarshalSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Provider.UnmarshalSession() = %v, want %v", got, tt.want)
			}
		})
	}
}
