package fritzbox

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestValidateUpdateQuery(t *testing.T) {

	tests := []struct {
		failMsg         string
		rawQuery        string
		parsedValues    url.Values
		expectedZoneIDs []string
		expectedToken   string
		expectedIP      string
		expectedError   string
	}{
		{
			failMsg:         "failed to validate illegal query",
			rawQuery:        "&&??illegal=query",
			parsedValues:    map[string][]string{},
			expectedError:   "url query must not be empty",
			expectedZoneIDs: []string{},
		},
		{
			failMsg:         "failed to validate on missing token",
			rawQuery:        "zone_id=42",
			parsedValues:    map[string][]string{"zone_id": {"42"}},
			expectedError:   "'token' must not be absent",
			expectedZoneIDs: []string{},
		},
		{
			failMsg:         "failed to validate on missing zone_id",
			rawQuery:        "token=4711",
			parsedValues:    map[string][]string{"token": {"4711"}},
			expectedError:   "'zone_id' must not be absent",
			expectedZoneIDs: []string{},
		},
		{
			failMsg:  "failed to validate on missing zone_id",
			rawQuery: "ip=127.0.0.1&zone_id=42&zone_id=43&token=4711",
			parsedValues: map[string][]string{"ip": {"127.0.0.1"}, "token": {"4711"},
				"zone_id": {"42", "43"}},
			expectedIP:      "127.0.0.1",
			expectedToken:   "4711",
			expectedZoneIDs: []string{"42", "43"},
			expectedError:   "",
		},
	}

	for _, test := range tests {
		parserMock := *new(ParserMock)
		parserMock.On("ParseQuery", test.rawQuery).Return(test.parsedValues, nil)

		actualZoneIDs, actualToken, actualIP, actualErr := validateUpdateQuery(test.rawQuery, parserMock.ParseQuery)
		assert := assert.New(t)
		if test.expectedError != "" {
			require.EqualError(t, actualErr, test.expectedError, test.failMsg)
		}
		assert.Equal(test.expectedToken, actualToken)
		assert.Equal(test.expectedIP, actualIP)
		assert.EqualValues(test.expectedZoneIDs, actualZoneIDs)
		parserMock.AssertExpectations(t)
	}

}
