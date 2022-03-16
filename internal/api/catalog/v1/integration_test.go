package catalogv1

import (
	"testing"
)

func TestPostInstall_Validate(t *testing.T) {
	type fields struct {
		Type  string
		Title string
		Body  string
	}
	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "invalid type",
			fields: fields{
				Type: "foo",
			},
			wantErr:    true,
			wantErrMsg: "invalid type",
		},
		{
			name: "markdown type without body",
			fields: fields{
				Type: "markdown",
			},
			wantErr:    true,
			wantErrMsg: "body cannot be empty for type markdown",
		},
		{
			name: "markdown type with body",
			fields: fields{
				Type: "markdown",
				Body: "foo",
			},
		},
		{
			name: "markdown type with body and title",
			fields: fields{
				Type:  "markdown",
				Body:  "foo",
				Title: "bar",
			},
			wantErr:    true,
			wantErrMsg: "title must be empty for type markdown",
		},
		{
			name: "section type without title",
			fields: fields{
				Type: "section",
			},
		},
		{
			name: "section type with title",
			fields: fields{
				Type:  "section",
				Title: "foo",
			},
		},
		{
			name: "section type with body and title",
			fields: fields{
				Type:  "section",
				Body:  "foo",
				Title: "bar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PostInstall{
				Type:  tt.fields.Type,
				Title: tt.fields.Title,
				Body:  tt.fields.Body,
			}
			err := p.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PostInstall.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErrMsg != err.Error() {
				t.Errorf("PostInstall.Validate() error msg = %v, wantErrMsg %v", err.Error(), tt.wantErrMsg)
			}
		})
	}
}
