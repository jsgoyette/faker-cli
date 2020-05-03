package handler

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"syreclabs.com/go/faker"
)

// HandleName handles the boolean map if `faker name` is called
func HandleName(opts map[string]interface{}) string {
	if opts["first"].(bool) {
		return faker.Name().FirstName()
	}
	if opts["last"].(bool) {
		return faker.Name().LastName()
	}
	return faker.Name().String()
}

// HandlePhone handles the boolean map if `faker phone` is called
func HandlePhone(opts map[string]interface{}) string {
	return faker.PhoneNumber().String()
}

func shortMode(opts map[string]interface{}) bool {
	return opts["--short"].(bool)
}

func getInt(from interface{}) (int, bool) {
	str, ok := from.(string)
	if !ok {
		return 0, false
	}
	val, err := strconv.ParseInt(str, 0, 64)
	if err != nil {
		panic(err)
	}
	return int(val), true
}

// HandleAddress handles `faker (city|state|zip-code|country)`
func HandleAddress(opts map[string]interface{}) string {
	if opts["country"].(bool) {
		if shortMode(opts) {
			return faker.Address().CountryCode()
		}
		return faker.Address().Country()
	}
	if opts["city"].(bool) {
		return faker.Address().City()
	}
	if opts["state"].(bool) {
		if shortMode(opts) {
			return faker.Address().StateAbbr()
		}
		return faker.Address().State()
	}
	if opts["street"].(bool) {
		return faker.Address().StreetAddress()
	}
	if opts["street2"].(bool) {
		return faker.Address().SecondaryAddress()
	}
	if opts["postal-code"].(bool) || opts["zip"].(bool) {
		if opts["--state"] != nil {
			faker.Address().PostcodeByState(opts["--state"].(string))
		}
		return faker.Address().Postcode()
	}
	return faker.Address().String()
}

// HandleSex handles `faker sex [--short]`
func HandleSex(opts map[string]interface{}) string {
	sexMap := map[bool]string{
		true:  "FEMALE",
		false: "MALE",
	}

	// random seed on a per call basis
	rand.Seed(time.Now().UnixNano())
	// binary gender easier to implement for now
	sex := sexMap[rand.Int()%2 == 0]
	if shortMode(opts) {
		return string(sex[0])
	}
	return sex
}

// HandleAdult handles `faker adult []` arguments
// TODO Pass in country codes for legal adult age
// current assumption is an adult age is 18 years or older
func HandleAdult(opts map[string]interface{}) string {

	// see this for formatting details:
	// https://golang.org/pkg/time/#Time.Format
	dateFormat := "2006-01-02"
	maxAge := 69
	if opts["--fmt"] != nil {
		dateFormat = opts["--fmt"].(string)
	}
	if opts["--max-age"] != nil {
		argAge, ok := getInt(opts["--max-age"])
		if !ok {
			panic("<max> must be an integer")
		}
		maxAge = argAge
	}
	dob := faker.Date().Birthday(18, maxAge)

	if opts["age"].(bool) {
		return fmt.Sprintf("%d", time.Now().Year()-dob.Year())
	}
	if opts["dob"].(bool) {
		if opts["-Y"].(bool) {
			return fmt.Sprintf("%d", dob.Year())
		}
		if opts["-M"].(bool) {
			return fmt.Sprintf("%d", dob.Month())
		}
		if opts["-D"].(bool) {
			return fmt.Sprintf("%d", dob.Day())
		}

		return dob.Format(dateFormat)
	}
	return "<nil>"
}

// HandleEmail handles `faker email`
func HandleEmail(opts map[string]interface{}) string {
	return faker.Internet().Email()
}

// HandlePassword handles `faker password` generation, allows a max and min length
// default is 8-24
func HandlePassword(opts map[string]interface{}) string {
	var ok bool
	min, max := 8, 24
	if opts["<min>"] != nil {
		min, ok = getInt(opts["<min>"])
		if !ok {
			panic("<min> must be an integer")
		}
	}
	if opts["<max>"] != nil {
		max, ok = getInt(opts["<max>"])
		if !ok {
			panic("<max> must be an integer")
		}
	}
	return faker.Internet().Password(min, max)
}
