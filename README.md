# Constraints

An experimental package for defining data constraints in Go.

Note that this project is in experimental stage. The interfaces are not
guaranteed to stay as they are now. Everything is subject to change.

## Background

When attempting to restructure errors in my project, I realized that errors
could be classified into two classes: API errors, and data errors.

Trying to map data errors made me realized one thing: a data is considered
invalid because it violates rules, or in another term, constraints.

A good data-related error, provides information which constraint(s) the data
violates. And this information should be able to be passed through the
systems down to the end-user interface (UI). We want the constraints to
be structured so that it could be nicely passed as understandable data to
other systems, including to be humanized.

## Design

Let's start from how we want to use this package.

First, we declare the individual constraints. We should be able to declare
the constraints inline, but by assigning them to variables, we can refer
them when we are trying to figure out which constraits coming out of a
validation.

```go
var (
    // Name the variables based on the semantic rather than the
    // description of the constraint (what a constraint is for, rather than
    // what the constraint does). For the description of the constraint,
    // we put it into the constraint instance itself.
    usernameFirstCharacter = Func(
        `starts with any letter from A to Z`,
        func(v string) bool {
            for _, r := range v {
                if !(r >= 'a' && r <= 'z') && !(r >= 'A' && r <= 'Z') {
                    return false
                }
            }
            return true
        })
    usernameLastCharacter = Func(
        `ends with any letter from A to Z`,
        func(v string) bool {
            for _, r := range v {
                if !(r >= 'a' && r <= 'z') && !(r >= 'A' && r <= 'Z') {
                    return false
                }
            }
            return true
        })
    // All of these constraints in this example could be declared as a
    // regular expression pattern like this, but we are trying to design
    // a mechanism which is more readable, constructed of smaller, clear
    // rules rather than a complex pattern which relatively hard to
    // understand.
    //
    // We let this slide for now until we can find a better way to declare
    // this kind of constraint. And also, because this pattern is pretty
    // simple.
    usernameAllowedCharacters = Func(
        `allowed are A to Z, 0 to 9 and underscore`,
        regexp.MustCompile(`^[A-Za-z0-9_]+$`).MatchString))
    usernameNoConsecutiveUnderscore = NoConsecutiveRune('_')
    usernameMinLength = MinLength(6)
    usernameMaxLength = MaxLength(32)
)
```

Next we define the constraint set:

```go
var usernameConstraints = Set(
    usernameFirstCharacter,
    usernameLastCharacter,
    usernameAllowedCharacters,
    usernameMinLength,
    usernameMaxLength,
)
```

Then we can use the set:

```go
violatedConstraints := usernameConstraints.ValidateAll(`hello`)
for _, vc := range violatedConstraints {
    switch vc {
    case usernameMinLength:
        // this case must be hit.
    }
}
```

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

## References

We want to align this package with existing standards.

- https://tools.ietf.org/html/draft-handrews-json-schema-validation-02
- https://developer.mozilla.org/en-US/docs/Web/Guide/HTML/HTML5/Constraint_validation
