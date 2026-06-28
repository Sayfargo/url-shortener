package tests

import (
	"testing"

	urlv1 "github.com/Sayfargo/protos/gen/go/url"
	"github.com/Sayfargo/url-shortener/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateShortCode_FailCases(t *testing.T) {
	ctx, conn := suite.New(t)

	testcases := []struct {
		name string
		url  string
	}{
		{
			name: "invalid url",
			url:  "shokoladka snickers123",
		},
		{
			name: "url with a typo",
			url:  "https:/fbig.pen",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			request := urlv1.CreateShortCodeRequest{
				Url: tt.url,
			}

			_, err := conn.UrlShortenerClient.Create(ctx, &request)
			assert.Error(t, err)

		})
	}

}

func TestCreateShortCode_Success(t *testing.T) {
	ctx, conn := suite.New(t)

	request := &urlv1.CreateShortCodeRequest{
		Url: "https://www.youtube.com/watch?v=erHz7d0v8Vw&t=783s",
	}

	resp, err := conn.UrlShortenerClient.Create(ctx, request)

	require.NoError(t, err)

	result := resp.GetShortCode()
	assert.NotEmpty(t, result)
}

func TestGetShortCode_Success(t *testing.T) {
	ctx, conn := suite.New(t)

	expected := "https://github.com/jackc/pgx"

	reqCreate := &urlv1.CreateShortCodeRequest{
		Url: expected,
	}

	respCreate, err := conn.UrlShortenerClient.Create(ctx, reqCreate)

	require.NoError(t, err)

	shortCode := respCreate.GetShortCode()
	require.NotEmpty(t, shortCode)

	reqGet := &urlv1.GetUrlRequest{
		ShortCode: shortCode,
	}

	respGet, err := conn.UrlShortenerClient.Get(ctx, reqGet)

	require.NoError(t, err)

	assert.Equal(t, expected, respGet.Url)

}
