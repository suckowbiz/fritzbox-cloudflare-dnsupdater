package fritzbox

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateUpdateQuery(t *testing.T) {

	tests := []struct {
		failMsg        string
		rawQuery       string
		parsedValues   url.Values
		expectedZoneID string
		expectedToken  string
		expectedIP     string
		expectedError  string
	}{
		{
			failMsg:       "failed to validate illegal query",
			rawQuery:      "&&??illegal=query",
			parsedValues:  map[string][]string{},
			expectedError: "url query must not be empty",
		},
		{
			failMsg:       "failed to validate on missing token",
			rawQuery:      "zone_identifier=42",
			parsedValues:  map[string][]string{"zone_identifier": []string{"42"}},
			expectedError: "'token' must not be absent",
		},
		{
			failMsg:       "failed to validate on missing zone_identifier",
			rawQuery:      "token=4711",
			parsedValues:  map[string][]string{"token": []string{"4711"}},
			expectedError: "'zone_identifier' must not be absent",
		},
		{
			failMsg:  "failed to validate on missing zone_identifier",
			rawQuery: "ip=127.0.0.1&zone_identifier=42&token=4711",
			parsedValues: map[string][]string{"ip": []string{"127.0.0.1"}, "token": []string{"4711"},
				"zone_identifier": []string{"42"}},
			expectedIP:     "127.0.0.1",
			expectedToken:  "4711",
			expectedZoneID: "42",
			expectedError:  "",
		},
	}

	for _, test := range tests {
		parserMock := *new(ParserMock)
		parserMock.On("ParseQuery", test.rawQuery).Return(test.parsedValues, nil)

		actualZoneID, actualToken, actualIP, actualErr := validateUpdateQuery(test.rawQuery, parserMock.ParseQuery)
		assert := assert.New(t)
		if test.expectedError != "" {
			assert.EqualError(actualErr, test.expectedError, test.failMsg)
		}
		assert.Equal(test.expectedToken, actualToken)
		assert.Equal(test.expectedIP, actualIP)
		assert.Equal(test.expectedZoneID, actualZoneID)
		parserMock.AssertExpectations(t)
	}

}
