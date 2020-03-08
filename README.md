# Constraints

An experimental package for defining data constraints in Go.

Note that this project is in experimental stage. The interfaces are not
guaranteed to stay as they are now.

## Background

When attempting to restructure errors in my project, I realized that errors
could be classified into two classes: API errors, and data errors.

Trying to map data errors made me realized one thing: a data is considered
invalid because it violates rules, or in another term, constraints.

A good data-related error, provides information which constraint(s) the data
violates. And this information could passed through the systems down to the
end-user interface (UI).

## Design

```go
var (
    usernameFirstCharacter = Func(
        `starts with any letter from A to Z`,
		func(v string) bool {
			for _, r := range v {
				if !(r > 'a' && r < 'z') && !(r > 'A' && r < 'Z') {
					return false
				}
			}
			return true
		})
    usernameLastCharacter = Func(
        `ends with any letter from A to Z`,
		func(v string) bool {
			for _, r := range v {
				if !(r > 'a' && r < 'z') && !(r > 'A' && r < 'Z') {
					return false
				}
			}
			return true
		})
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

Then upon use:

```go
violatedConstraints := usernameConstraints.ValidateAll(`hello`)
for _, vc := range violatedConstraints {
    switch vc {
    case usernameMinLength:
        // ...
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

Ideally, we'd like to make it easy to generate validation directive for
other systems, e.g., JSON Schema.

```go
usernameJSONSchemaConstraints := usernameConstraints.JSONSchemaConstraints()
```
