# Go - Validator

This package provides a validations for Go applications.

### Validation Rules

To use them, the tag format is `validate:"<Name>"`.

| Name | Note
|---|---
| OmitEmpty | Omit rest of validations if value is empty
| Required | Must be non-zero value
| AlphaDash | Must be alpha characters or numerics or -_
| AlphaDashDot | Must be alpha characters or numerics, -_ or .
| Size(int) | Fixed length
| MinSize(int) | Minimum length
| MaxSize(int) | Maximum length
| Email | Must be E-mail address
| Url | Must be HTTP/HTTPS URL address
| Include(string) | Must contain
| Exclude(string) | Must not contain

To combine multiple rules: `govalidator:"Required;MinSize(10)".`

### Example

This is an exemple of how to use **govalidator**:

    package main
    
    import (
        "fmt"
    
        "github.com/rodkranz/govalidator"
    )
    
    type User struct {
        FirstName  *string `validate:"Required" alias:"Name"`
        SecondName string  `validate:"Required" alias:"Second"`
    }
    
    func main() {
        Validate(&User{FirstName: String("Rodrigo"), SecondName: "Lopes"})
        Validate(&User{})
    }
    
    func Validate(u *User) {
        isOk, errs := govalidator.Validate(u)
    
        if isOk {
            fmt.Printf("The struct %T is valid!\n", u)
        }
    
        if errs != nil {
            for i := range errs {
                fmt.Printf("The field %s is wrong: %v \n", errs[i].Field(), errs[i].Message)
            }
        }
    }
    
    func String(s string) *string {
        return &s
    }


Result:

    $ go run main.go
    The struct *main.User is valid!
    The field Name is wrong: Required
    The field Second is wrong: Required


### Custom Validation Rules

If you need to use the specific rule for you case, you can use the `govalidator.Rule` or `govalidator.ParamRule`.

Suppose you want to check if advert id is blocked:

    govalidator.AddParamRule(&govalidator.ParamRule{
        IsMatch: func(rule string) bool {
            return strings.HasPrefix(rule, "AdBlockIDs(")
        },
        IsValid: func(errs govalidator.Errors, rule, name, alias string, v interface{}) (bool, govalidator.Errors) {
            idAdvert, ok := v.(int)
            if !ok {
                return false, errs
            }

            ids := rule[11: len(rule)-1]
            for _, idString := range strings.Split(ids, ",") {
                id, _ := strconv.Atoi(idString)
                if id == idAdvert {
                    errs.Add(name, alias, "AdIdAllowed", "This Advert id is blocked")
                    return false, errs
                }
            }

            return true, errs
        },
    })

If your rule is simple, you can also use `govalidator.AddRule`, it accepts type `govalidator.Rule`:


    govalidator.AddRule(&govalidator.Rule{
        IsMatch: func(rule string) bool {
            return rule == "AdExists"
        },
        IsValid: func(errs govalidator.Errors, name, alias string, v interface{}) (bool, govalidator.Errors) {
            idAdvert, ok := v.(int)
            if !ok {
                return false, errs
            }

            if idAdvert < 10 {
                errs.Add(name, alias, "AdExists", "This advert not exist in our database.")
            }

            return true, errs
        },
    })


