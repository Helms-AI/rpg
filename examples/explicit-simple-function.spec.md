# slugify

A utility function to convert text into URL-friendly slugs.

## Target Languages

- go
- python
- typescript

## Functions

### slugify

Converts a string into a URL-friendly slug.

**accepts:**
- text: Text
- separator: Text (defaults to "-")

**returns:** Text

**logic:**
```
convert text to lowercase
replace all whitespace with the separator
remove all characters that are not letters, numbers, or the separator
collapse multiple consecutive separators into one
trim separators from start and end
return the result
```

## Tests

### slugify

#### test: converts simple text
given: "Hello World"
expect: "hello-world"

#### test: handles special characters
given: "Hello, World! How are you?"
expect: "hello-world-how-are-you"

#### test: uses custom separator
given:
- text: "Hello World"
- separator: "_"
expect: "hello_world"

#### test: handles unicode
given: "Café Au Lait"
expect: "cafe-au-lait"
