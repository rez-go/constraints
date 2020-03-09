# Constraints

An experimental package for defining data constraints in Go.

Note that this project is in experimental stage. The interfaces are not
guaranteed to stay as they are now. Everything is subject to change.

## Background

When attempting to restructure errors in my project, I realized that errors
could be classified into two classes: API errors, and data errors.

When we are talking about data-related errors, most are instances are coming
out from data validation. An error which coming from a validation process
should provide enough information about why the data is considered as invalid,
i.e., the rules or constraints it violated.

For example, we define `TooShortError`. We expect it to contains details
about how is the minimum before a value is considered invalid. So, the error
is actually containing information about the constraint. Rather than defining
the errors, we could just make a generic error, something like
`DataValidationError`, and make it holding the information about the
constraint(s) a data is violating. Let the presentation layer translate it
based on the consumer.

We want the constraints to be structured so that it could be nicely passed
as understandable data to other systems, including to the end-user
presentation layer, which could be translated to any human language.

## Design

Let's start from how we want to use this package.

First, we declare the individual constraints. We should be able to declare
the constraints inline, but by assigning them to variables, we can refer
them when we are trying to figure out which constraint(s) coming out of a
validation.

```go
var (
    usernameMinLength = MinLength(6)
    usernameMaxLength = MaxLength(32)
    // All these constraints in this example could be declared as a
    // single regular expression pattern, but we are trying to design
    // a mechanism which is more readable, constructed of smaller, clear
    // rules rather than putting all the rules into a complex pattern.
    //
    // We might want to find a way to declare for this kind of constraint,
    // but we will limit how far we will go before we reinvent regular
    // expression.
    usernameAllowedCharacters = Func(
        `allowed characters are A to Z, 0 to 9 and underscore`,
        regexp.MustCompile(`^[A-Za-z0-9_]+$`).MatchString))
    // Name the variables based on the semantic rather than the
    // description of the constraint (what a constraint is for, rather than
    // what the constraint does). For the description of the constraint,
    // we put it into the constraint instance itself.
    //
    // The variable name and the description are good if they sound good if
    // we merge them:
    // "username [starts with] a letter".
    usernameStartsWith = Func(
        `starts with a letter`,
        func(v string) bool {
            if v != "" {
                r, _ := utf8.DecodeRuneInString(v)
                return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
            }
            return false
        })
    usernameEndsWith = Func(
        `ends with anything but underscore`,
        func(v string) bool {
            if v != "" {
                r, _ := utf8.DecodeLastRuneInString(v)
                return r != '_'
            }
            return false
        })
    usernameNoConsecutiveUnderscore = NoConsecutiveRune('_')
)
```

Next we define the constraint set:

```go
var usernameConstraints = Set(
    usernameMinLength,
    usernameMaxLength,
    usernameAllowedCharacters,
    usernameStartsWith,
    usernameEndsWith,
    usernameNoConsecutiveUnderscore,
)
```

We could generate a decent instruction from the constraint:

```go
fmt.Printf("username: %s\n", usernameConstraints.ConstraintDescription())
```

Which would print something like `username: min length 6, max length 32,
allowed characters are ...`. Ideally, a constraint is mapped to some well
thought messages if it's going to be displayed to human.

Then we can use the set:

```go
violatedConstraints := usernameConstraints.ValidateAll(`h-llo`)
for _, vc := range violatedConstraints {
    switch vc {
    case usernameMinLength:
        // this case must be hit.
    case usernameAllowedCharacters:
        // this case must also be hit.
    }
}
```

The `violatedConstraints` contains all the constraints the data violates.
This could be a good case for end user as they could fix their input in
one go rather than back-and-forth.

Or as error:

```go
err := usernameConstraints.ValidOrError(username)
if err != nil {
    return ArgumentError("username", err)
}
```

The error contains a structured information which describes the violations so
that we can encode it and put it into response.

Ideally, we'd like to make it easy to generate validation directive for
other systems, e.g., JSON Schema.

```go
//...

jsonSchemaField := usernameConstraints.JSONSchemaField("username")

//...
```

## Hacking

- We are still in designing stage. Anything could be suggested to be
  changed.
- We limit the dependencies only to Go's stdlib.
- We won't include rules which don't have strict, static constraints.
  We might won't even include rules for standard formats.
  We will never include validations for formats like phone numbers
  (it's very dynamic), email addresses, domains and postal addresses.
  We should provide enough facilities to others to provide module for
  certain type of validation.
- We might want to limit this module to contain only the contracts and
  limited utilites, and put common constraints into their own module
  (discuss!).

## References

We want to align this package with existing standards.

- https://tools.ietf.org/html/draft-handrews-json-schema-validation-02
- https://developer.mozilla.org/en-US/docs/Web/Guide/HTML/HTML5/Constraint_validation
