package aps

import (
	"reflect"
	"testing"

	"github.com/markbates/goth"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestNew(t *testing.T) {
	type args struct {
		clientKey   string
		secret      string
		callbackURL string
		scopes      []string
	}
	tests := []struct {
		name string
		args args
		want *Provider
	}{
		{
			name: "Without scope",
			args: args{
				clientKey:   "Foo",
				secret:      "Bar",
				callbackURL: "/baz",
			},
			want: &Provider{
				ClientKey:   "Foo",
				Secret:      "Bar",
				CallbackURL: "/baz",
			},
		},
		{
			name: "With scope",
			args: args{
				clientKey:   "Bar",
				secret:      "Foo",
				callbackURL: "/bee",
				scopes:      []string{"a"},
			},
			want: &Provider{
				ClientKey:   "Bar",
				Secret:      "Foo",
				CallbackURL: "/bee",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			got := New(tt.args.clientKey, tt.args.secret, tt.args.callbackURL, tt.args.scopes...)
			a.Equal(tt.want.ClientKey, got.ClientKey)
			a.Equal(tt.want.Secret, got.Secret)
			a.Equal(tt.want.CallbackURL, got.CallbackURL)
		})
	}
}
func TestProvider_RefreshTokenAvailable(t *testing.T) {
	a := assert.New(t)
	p := &Provider{}
	result := p.RefreshTokenAvailable()
	a.True(result)
}

func TestProvider_Name(t *testing.T) {
	p := &Provider{}
	if got := p.Name(); got != "aps" {
		t.Errorf("Provider.Name() = %v, want aps", got)
	}
}

func TestProvider_Debug(t *testing.T) {
	type args struct {
		debug bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "True",
			args: args{
				debug: true,
			},
		},
		{
			name: "False",
			args: args{
				debug: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{}
			p.Debug(tt.args.debug)
		})
	}
}

func TestProvider_BeginAuth(t *testing.T) {
	type fields struct {
		config *oauth2.Config
		prompt oauth2.AuthCodeOption
	}
	type args struct {
		state string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    goth.Session
		wantErr bool
	}{
		{
			name: "First",
			fields: fields{
				config: &oauth2.Config{
					ClientID: "id",

					ClientSecret: "secret",
					Endpoint: oauth2.Endpoint{
						AuthURL:  "auth",
						TokenURL: "token",
					},
					RedirectURL: "redir",
					Scopes:      []string{},
				},
			},
			args: args{
				state: "",
			},
			want: &Session{
				AuthURL: "auth?client_id=id\u0026redirect_uri=redir\u0026response_type=code",
			},
			wantErr: false,
		},
		{
			name: "With scope",
			fields: fields{
				config: &oauth2.Config{
					ClientID: "id",

					ClientSecret: "secret",
					Endpoint: oauth2.Endpoint{
						AuthURL:  "auth",
						TokenURL: "token",
					},
					RedirectURL: "redir",
					Scopes:      []string{},
				},
			},
			args: args{
				state: "mystate",
			},
			want: &Session{
				AuthURL: "auth?client_id=id\u0026redirect_uri=redir\u0026response_type=code\u0026state=mystate",
			},
			wantErr: false,
		},
		{
			name: "With scope and prompt",
			fields: fields{
				config: &oauth2.Config{
					ClientID: "id",

					ClientSecret: "secret",
					Endpoint: oauth2.Endpoint{
						AuthURL:  "auth",
						TokenURL: "token",
					},
					RedirectURL: "redir",
					Scopes:      []string{},
				},
				prompt: oauth2.SetAuthURLParam("prompt", "one two"),
			},
			args: args{
				state: "mystate",
			},
			want: &Session{
				AuthURL: "auth?client_id=id\u0026prompt=one+two\u0026redirect_uri=redir\u0026response_type=code\u0026state=mystate",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{
				config: tt.fields.config,
				prompt: tt.fields.prompt,
			}
			got, err := p.BeginAuth(tt.args.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.BeginAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Provider.BeginAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProvider_SetPrompt(t *testing.T) {
	type args struct {
		prompt []string
	}

	tests := []struct {
		name string
		args args
		want oauth2.AuthCodeOption
	}{
		{
			name: "Empty",
			args: args{
				prompt: []string{},
			},
			want: nil,
		},
		{
			name: "Single value",
			args: args{
				prompt: []string{"one"},
			},
			want: oauth2.SetAuthURLParam("prompt", "one"),
		},
		{
			name: "Two values",
			args: args{
				prompt: []string{"one", "two"},
			},
			want: oauth2.SetAuthURLParam("prompt", "one two"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{}
			p.SetPrompt(tt.args.prompt...)
			if !reflect.DeepEqual(p.prompt, tt.want) {
				t.Errorf("Provider.SetPrompt() = %v, want %v", p.prompt, tt.want)
			}
		})
	}
}

func Test_newConfig(t *testing.T) {
	type args struct {
		provider *Provider
		scopes   []string
	}
	tests := []struct {
		name string
		args args
		want *oauth2.Config
	}{
		{
			name: "First test",
			args: args{
				provider: &Provider{
					ClientKey:   "key",
					Secret:      "secret",
					CallbackURL: "/foo",
				},
			},
			want: &oauth2.Config{
				ClientID:     "key",
				ClientSecret: "secret",
				RedirectURL:  "/foo",
				Endpoint: oauth2.Endpoint{
					AuthURL:  authURL,
					TokenURL: tokenURL,
				},
				Scopes: []string{},
			},
		},
		{
			name: "Another test",
			args: args{
				provider: &Provider{
					ClientKey:   "keyFoo",
					Secret:      "secretBar",
					CallbackURL: "/foobar",
				},
			},
			want: &oauth2.Config{
				ClientID:     "keyFoo",
				ClientSecret: "secretBar",
				RedirectURL:  "/foobar",
				Endpoint: oauth2.Endpoint{
					AuthURL:  authURL,
					TokenURL: tokenURL,
				},
				Scopes: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newConfig(tt.args.provider, tt.args.scopes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
