package randx

import (
	"fmt"
	"testing"
	"time"
)

func TestUserCases(t *testing.T) {
	template := NewSampleTemplate()
	eighteenYearsAgo := time.Now().AddDate(-18, 0, 0).Format(time.RFC3339)

	testCases := []struct {
		name string
		user *User
	}{
		{
			"any_user",
			SomeOne(template),
		},
		{
			"adult",
			SomeOne(template).
				WithDOB(fmt.Sprintf("birthdate < '%s'", eighteenYearsAgo)),
		},
		{
			"name_begin",
			SomeOne(template).
				WithName("begins == 'Nguyễn'"),
		},
		{
			"tax_free",
			SomeOne(template).
				WithTaxPercent("0"),
		},
		{
			"high_score",
			SomeOne(template).
				WithScore("score > 100"),
		},
		{
			"telco_Viettel",
			SomeOne(template).
				WithPhonePrefix("^[+]8497[0-9]{7}$"),
		},
		{
			"combined",
			SomeOne(template).
				WithTaxPercent("tax_fee > 0.1").
				WithScore("score > 100 && score < 300").
				WithName("length == 4").
				WithDOB(fmt.Sprintf("birthdate < '%s'", eighteenYearsAgo)),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name, tc.user)
		})
	}
}
