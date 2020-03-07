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
