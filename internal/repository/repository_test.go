package repository

import (
	"testing"

	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestRepo(t *testing.T) {
	tests := []struct{
		name string
		url model.URL
		want string
	}{
		{
			name:    "test add and get",
			url: model.URL{
				OriginalURL: "https://www.yandex.com/",
				UUID: "a4d1as",
			},
			want: "a4d1as",
		},
	}
	for _, test := range tests{
		t.Run(test.name, func(t *testing.T){
			repo, _ := NewRepository("test.json")
			err := repo.Store(&test.url)
			if err != nil{
				assert.Errorf(t,err, "Error add")
			}
			url, err := repo.FindByShortCode(test.url.UUID)
			if err != nil{
				assert.Errorf(t,err, "longUrl is not found")
			}	
			assert.Equal(t, url.OriginalURL, test.url.OriginalURL)
		})
	}	
}
