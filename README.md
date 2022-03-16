randx
====
Provides methods for picking a random value that satisfy user-readable conditions.

Which value would you choose between the boundaries?
--

Beside boundaries, middle values are the most important inputs in Software Testing. The middle value is usually a valid value, 
we can pick any and hardcode in the test, if it works all remain work (assume).

Obviously, no tester ever thinks of testing all the possible values thus random one in each execution time would make sense and be popular.
I do. But I get into inconvenience while generating random inputs for my test cases. That's why randx exists. 
Hopefully, this piece would become handy for people, especially testers.

What does randx offer?
--

Random a number

```go
	generator :=  randx.NewXNumber(0, 1000, 1)
	x := generator.Random("number < 100 || number > 200")
	// outputs of 3 execution times: 2, 16, 553
```

Random a string slice

```go
	generator :=  randx.NewXWord("etc/vietnamese")
	name1 := generator.Random("length == 4")
	name2 := generator.Random("begin == 'Nguyễn'")
	name3 := generator.Random("end == 'Vi'")
	// outputs:
	// name1 [Ngọc Phượng Hằng Thành]
	// name2 [Nguyễn Lý Phi]
	// name3 [Viên Thuận Vi]
```

Random a datetime

```go
	generator :=  randx.NewXTime(time.Now().AddDate(-100, 0, 0), time.Now(), 24*time.Hour)
	birthDay := generator.Random("birthday > '1981-01-01' && birthday < '1996-12-31'") // millennials
	// output 1982-09-04
```

Random a string as regex
```go
	generator :=  randx.NewXRegex()
	id := generator.Random("^[a-z]{6}[0-9]{2}$")
	// output yecrzi97
```

Random a user as your template

```go
	// define your own template
	template := UserTemplate{
	Name:       NewXWords("etc/vietnamese"),
	Phone:      NewXRegex(), 
	Score:      NewXNumber(0, 1000, 1), 
	TaxPercent: NewXNumber(0, 1, 0.05), 
	DOB:        NewXTime(time.Now().AddDate(-100, 0, 0), time.Now(), 24*time.Hour),
	}
	template.Phone.SetFallback("^[+]84[0-9]{9}$")
	
	// customize your chain and try
	user1 := SomeOne(template)
	user2 := SomeOne(template).WithName("begin == 'Nguyễn'")
	user3 := SomeOne(template).
		WithTaxPercent("tax_fee > 0.1").
		WithScore("score > 100 && score < 300").
		WithName("length == 4").
		WithDOB(fmt.Sprintf("birthdate < '%s'", eighteenYearsAgo.Format(time.RFC3339)))
	// outputs:
	// user1 [Lý] +84004343908 187 0.04 1882-09-07
	// user2 [Nguyễn Trâm Thanh] +84025848714 531 0.92 2004-03-29
	// user3 [My Chúc Trần Đinh] +84688014056 247 0.90 1988-03-18
```

What should be noted?
--
TBD

Reference
--
Thanks to [govaluate](https://github.com/Knetic/govaluate). randx based on a core of govaluate, the parsing expression.
And thanks to [goregen](https://github.com/zach-klippenstein/goregen). randx.XRegex just be a wrapper of goregen functions.

License
--

This project is licensed under the BSD 2-Clause License. You're free to integrate, fork, and play with this code as you feel fit without consulting the author, as long as you provide proper credit to the author in your works.
