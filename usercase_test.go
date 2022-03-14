package randx

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSampleUserCases(t *testing.T) {
	template := NewSampleTemplate()
	eighteenYearsAgo := time.Now().AddDate(-18, 0, 0)

	testCases := []struct {
		name    string
		user    *User
		checker func(user *User) bool
	}{
		{
			"any_user", SomeOne(template),
			func(u *User) bool { return true },
		},
		{
			"adult", SomeOne(template).WithDOB(fmt.Sprintf("birthdate < '%s'", eighteenYearsAgo.Format(time.RFC3339))),
			func(u *User) bool {
				return u.DOB.(time.Time).Before(eighteenYearsAgo)
			},
		},
		{
			"name_begin_with_nguyen", SomeOne(template).WithName("begins == 'Nguyễn'"),
			func(u *User) bool {
				return u.Name.([]string)[0] == "Nguyễn"
			},
		},
		{
			"tax_free", SomeOne(template).WithTaxPercent("0"),
			func(u *User) bool {
				return u.TaxPercent.(float64) == 0
			},
		},
		{
			"high_score", SomeOne(template).WithScore("score > 100"),
			func(u *User) bool {
				return u.Score.(float64) > 100
			},
		},
		{
			"telco_Viettel", SomeOne(template).WithPhonePrefix("^[+]8497[0-9]{7}$"),
			func(u *User) bool {
				return u.Phone.(string)[:5] == "+8497"
			},
		},
		{
			"combined", SomeOne(template).
				WithTaxPercent("tax_fee > 0.1").
				WithScore("score > 100 && score < 300").
				WithName("length == 4").
				WithDOB(fmt.Sprintf("birthdate < '%s'", eighteenYearsAgo.Format(time.RFC3339))),
			func(u *User) bool {
				return u.TaxPercent.(float64) > 0.1 &&
					u.Score.(float64) > 100 && u.Score.(float64) < 300 &&
					len(u.Name.([]string)) == 4 &&
					u.DOB.(time.Time).Before(eighteenYearsAgo)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name, tc.user)
			require.True(t, tc.checker(tc.user))
		})
	}
}
