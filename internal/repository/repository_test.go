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
				LongURL: "https://www.yandex.com/",
				ShortCode: "a4d1as",
			},
			want: "a4d1as",
		},
	}
	for _, test := range tests{
		t.Run(test.name, func(t *testing.T){
			repo := NewRepository()
			err := repo.Store(&test.url)
			if err != nil{
				assert.Errorf(t,err, "Error add")
			}
			url, err := repo.FindByShortCode(test.url.ShortCode)
			if err != nil{
				assert.Errorf(t,err, "longUrl is not found")
			}	
			assert.Equal(t, url.LongURL, test.url.LongURL)
		})
	}	
}
