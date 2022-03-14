package randx

import "time"

type UserTemplate struct {
	Name       X
	Phone      X
	Score      X
	TaxPercent X
	DOB        X
}

type User struct {
	template   UserTemplate
	Name       interface{}
	Phone      interface{}
	Score      interface{}
	TaxPercent interface{}
	DOB        interface{}
}

func NewSampleTemplate() UserTemplate {
	template := UserTemplate{
		Name:       NewXWords("etc/vietnamese"),
		Phone:      NewXRegex(),
		Score:      NewXNumber(0, 1000, 1),
		TaxPercent: NewXNumber(0, 1, 0.05),
		DOB:        NewXTime(time.Now().AddDate(-200, 0, 0), time.Now(), 24*time.Hour),
	}
	template.Phone.SetFallback("^[+]84[0-9]{9}$")
	return template
}

func SomeOne(template UserTemplate) *User {
	return &User{
		template:   template,
		Name:       template.Name.Random(""),
		Phone:      template.Phone.Random(""),
		Score:      template.Score.Random(""),
		TaxPercent: template.TaxPercent.Random(""),
		DOB:        template.DOB.Random(""),
	}
}

func (user *User) WithDOB(condition string) *User {
	user.DOB = user.template.DOB.Random(condition)
	return user
}

func (user *User) WithTaxPercent(condition string) *User {
	user.TaxPercent = user.template.TaxPercent.Random(condition)
	return user
}

func (user *User) WithScore(condition string) *User {
	user.Score = user.template.Score.Random(condition)
	return user
}

func (user *User) WithPhonePrefix(regex string) *User {
	user.Phone = user.template.Phone.Random(regex)
	return user
}

func (user *User) WithName(condition string) *User {
	user.Name = user.template.Name.Random(condition)
	return user
}
