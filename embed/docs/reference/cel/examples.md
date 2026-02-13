---
title: CEL Function Examples
metaTitle: "Reference: CEL Examples"
description: "Example usage of nearly all CEL functions in Twisp runtime."
showNext: false
showPrev: false
showToc: true
---

Twisp pervasively uses [Common Expression Language](https://cel.dev/) for computation across many apis:

1. Index Value and Filtering Definitions
2. Security Policy Conditions
3. Webhook Filters
4. Calculations
5. Velocity Controls
6. Tran Codes


In addition to the standard CEL library, Twisp exposes a number of custom functions that you might find useful as you build your core accounting system. On this page you can find extensive executable examples you can use to discover what the CEL functions can do for you.

## Constructors & Conversions

These functions allow you to create instances of specific data types or convert values between different types. This is essential for manipulating data within CEL expressions effectively.

### Bool

Converts a string ("true" or "false") to a boolean value.

```graphql
# func bool(bool) bool
# func bool(string) bool
mutation BoolConversion {
  evaluate(
    expressions: {
      fromStringTrue: "bool('true')" # true
      fromStringFalse: "bool('false')" # false
    }
  )
}
```

### Bytes

Converts strings or UUIDs into byte sequences. Useful for hashing or encoding functions.

```graphql
# func bytes(bytes) bytes
# func bytes(string) bytes
# func bytes(twisp.type.v1.UUID) bytes
mutation BytesConversion {
  evaluate(
    expressions: {
      fromString: "string(bytes('hello'))" # Should show bytes representation or convert back for check: "hello"
      fromUuid: "hex.EncodeToString(bytes(uuid('123e4567-e89b-12d3-a456-426614174000')))" # Hex of UUID bytes: "123e4567e89b12d3a456426614174000"
    }
  )
}
```

### Date

Constructs a Date object from a string (YYYY-MM-DD format) or a Timestamp.

```graphql
# func date(twisp.type.v1.Date) twisp.type.v1.Date
# func date(google.protobuf.Timestamp) twisp.type.v1.Date
# func date(string) twisp.type.v1.Date
mutation DateConversion {
  evaluate(
    expressions: {
      fromString: "string(date('2023-10-27'))" # "2023-10-27"
      fromTimestamp: "string(date(timestamp('2023-10-27T15:00:00Z')))" # "2023-10-27"
    }
  )
}
```

### Decimal

Constructs a high-precision Decimal object from a string representation or converts a Money object to its Decimal value.

```graphql
# func decimal(twisp.type.v1.Decimal) twisp.type.v1.Decimal
# func decimal(twisp.type.v1.Money) twisp.type.v1.Decimal
# func decimal(string) twisp.type.v1.Decimal
# func (twisp.type.v1.Money) decimal() twisp.type.v1.Decimal
mutation DecimalConversion {
  evaluate(
    expressions: {
      fromString: "string(decimal('123.456'))" # "123.456"
      fromMoney: "string(decimal(money('99.99', 'USD')))" # "99.99"
      fromMoneyMethod: "string(money('50.00', 'EUR').decimal())" # "50.00"
    }
  )
}
```

### Double

Converts integer, unsigned integer, or string representations to a double-precision floating-point number.

```graphql
# func double(double) double
# func double(int) double
# func double(string) double
# func double(uint) double
mutation DoubleConversion {
  evaluate(
    expressions: {
      fromInt: "double(5)" # 5.0
      fromString: "double('3.14')" # 3.14
      fromUint: "double(10u)" # 10.0
    }
  )
}
```

### Duration

Constructs a Duration object from a string representation (e.g., "1h30m15s").

```graphql
# func duration(google.protobuf.Duration) google.protobuf.Duration
# func duration(int) google.protobuf.Duration
# func duration(string) google.protobuf.Duration
mutation DurationConversion {
  evaluate(
    expressions: {
      fromString: "string(duration('1h30m15s'))" # "5415s" (canonical string representation)
    }
  )
}
```

### Dyn

Converts a value to a dynamic type, useful for functions accepting heterogeneous lists or for working with JSON-like structures.

```graphql
# func dyn(A) dyn
mutation DynExample {
  evaluate(expressions: { result: "dyn('hello') == dyn('hello')" }) # Example: true (demonstrates dyn conversion and comparison)
}
```

### Int

Converts various types (double, duration, string, timestamp, uint) to an integer. Note that conversion from double truncates the decimal part. Timestamps convert to epoch seconds, durations to total seconds.

```graphql
# func int(int) int
# func int(double) int
# func int(google.protobuf.Duration) int
# func int(string) int
# func int(google.protobuf.Timestamp) int
# func int(uint) int
mutation IntConversion {
  evaluate(
    expressions: {
      fromDouble: "int(3.9)" # 3 (truncates)
      fromString: "int('123')" # 123
      fromUint: "int(100u)" # 100
      fromTimestamp: "int(timestamp(0))" # 0 (epoch seconds)
      fromDuration: "int(duration('60s'))" # 60 (total seconds)
    }
  )
}
```

### Money

Constructs a Money object, typically requiring an amount (as string, decimal, or int) and a currency code string.

```graphql
# func money(twisp.type.v1.Money) twisp.type.v1.Money
# func money(string, string) twisp.type.v1.Money
# func money(twisp.type.v1.Decimal, string) twisp.type.v1.Money
# func money(int, string) twisp.type.v1.Money
mutation MoneyConversion {
  evaluate(
    expressions: {
      fromString: "money('123.45', 'USD')" # "123.45 USD"
      fromDecimal: "money(decimal('99.99'), 'EUR')" # "99.99 EUR"
      fromInt: "money(100, 'JPY')" # "100 JPY" (no decimals for JPY often)
    }
  )
}
```

### String

Converts various data types (boolean, bytes, double, duration, int, timestamp, uint, date, decimal, UUID) into their string representation.

```graphql
# func string(string) string
# func string(bool) string
# func string(bytes) string
# func string(double) string
# func string(google.protobuf.Duration) string
# func string(int) string
# func string(google.protobuf.Timestamp) string
# func string(uint) string
# func string(twisp.type.v1.Date) string
# func string(twisp.type.v1.Decimal) string
# func string(twisp.type.v1.UUID) string
mutation StringConversion {
  evaluate(
    expressions: {
      fromBool: "string(true)" # "true"
      fromBytes: "string(b'hello')" # "hello"
      fromDouble: "string(3.14)" # "3.14"
      fromDuration: "string(duration('1m'))" # "60s" or similar
      fromInt: "string(123)" # "123"
      fromTimestamp: "string(timestamp(0))" # "1970-01-01T00:00:00Z"
      fromUint: "string(10u)" # "10"
      fromDecimal: "string(decimal('1.2'))" # "1.2"
      fromUuid: "string(uuid.New())" # Example: "a1b2c3d4-..."
    }
  )
}
```

### Timestamp

Constructs a Timestamp object from a string (RFC3339 format) or an integer (interpreted as epoch seconds or milliseconds, depending on implementation).

```graphql
# func timestamp(google.protobuf.Timestamp) google.protobuf.Timestamp
# func timestamp(int) google.protobuf.Timestamp
# func timestamp(string) google.protobuf.Timestamp
mutation TimestampConversion {
  evaluate(
    expressions: {
      fromString: "string(timestamp('2023-10-27T12:00:00Z'))" # "2023-10-27T12:00:00Z"
      fromEpoch: "string(timestamp(1678886400))" # Interpreted as epoch seconds: "2023-03-15T13:20:00Z"
    }
  )
}
```

### Type

Returns the type of a given value. Useful for introspection or conditional logic based on type.

```graphql
# func type(A) type(A)
mutation TypeExample {
  evaluate(
    expressions: {
      intType: "type(123) == int" # true
      stringType: "type('abc') == string" # true
    }
  )
}
```

### Uint

Converts double, integer, or string representations to an unsigned integer. Conversion from double truncates.

```graphql
# func uint(uint) uint
# func uint(double) uint
# func uint(int) uint
# func uint(string) uint
mutation UintConversion {
  evaluate(
    expressions: {
      fromDouble: "uint(3.9)" # 3u (truncates)
      fromString: "uint('123')" # 123u
      fromInt: "uint(100)" # 100u
    }
  )
}
```

### UUID

Constructs a UUID object from its string or byte representation.

```graphql
# func uuid(twisp.type.v1.UUID) twisp.type.v1.UUID
# func uuid(string) twisp.type.v1.UUID
# func uuid(bytes) twisp.type.v1.UUID
mutation UuidConversion {
  evaluate(
    expressions: {
      fromString: "string(uuid('123e4567-e89b-12d3-a456-426614174000'))" # "123e4567-e89b-12d3-a456-426614174000"
      fromBytes: "string(uuid(b'\\x12>Eg\\xe8\\x9b\\x12\\xd3\\xa4VBC\\x14\\x17@\\x00'))" # "123e4567-e89b-12d3-a456-426614174000"
    }
  )
}
```

## Operators

CEL supports common operators for arithmetic, comparison, logic, and collection access.

### Arithmetic Operators

#### Addition / Concatenation (+)

Performs addition for numeric types (int, double, decimal, money) and duration/timestamp combinations. Concatenates strings, bytes, and lists.

```graphql
# func _+_(...) (Addition/Concatenation Operator)
mutation AddOperator {
  evaluate(
    expressions: {
      int: "5 + 3" # 8
      double: "1.5 + 2.5" # 4.0
      string: "'hello' + ' ' + 'world'" # "hello world"
      list: "[1, 2] + [3, 4]" # [1, 2, 3, 4]
      timestampDuration: "string(timestamp('2023-01-01T00:00:00Z') + duration('24h'))" # "2023-01-02T00:00:00Z"
      durationDuration: "string(duration('1h') + duration('30m'))" # "5400s" or "1h30m"
      decimal: "string(decimal('1.1') + decimal('2.2'))" # "3.3"
    }
  )
}
```

#### Subtraction (-)

Performs subtraction for numeric types (int, double, decimal, money), duration/timestamp combinations, and between two timestamps (resulting in a duration).

```graphql
# func _-_(...) (Subtraction Operator)
mutation SubtractOperator {
  evaluate(
    expressions: {
      int: "10 - 3" # 7
      double: "5.5 - 1.5" # 4.0
      timestampDuration: "string(timestamp('2023-01-02T00:00:00Z') - duration('24h'))" # "2023-01-01T00:00:00Z"
      timestampTimestamp: "string(timestamp('2023-01-02T00:00:00Z') - timestamp('2023-01-01T00:00:00Z'))" # "86400s" or "24h"
      durationDuration: "string(duration('1h') - duration('15m'))" # "2700s" or "45m"
      decimal: "string(decimal('3.3') - decimal('1.1'))" # "2.2"
    }
  )
}
```

#### Multiplication (*)

Performs multiplication for numeric types (int, uint, double, decimal).

```graphql
# func _*_(double, double) double
# func _*_(int, int) int
# func _*_(uint, uint) uint (Multiplication Operator)
mutation MultiplicationOperator {
  evaluate(
    expressions: {
      int: "5 * 3" # 15
      uint: "5u * 3u" # 15u
      double: "1.5 * 2.0" # 3.0
      # decimal: decimal('1.5') * decimal('2') # decimal('3.0')
    }
  )
}
```

#### Division (/)

Performs division for numeric types (int, uint, double, decimal). Integer and uint division truncates towards zero.

```graphql
# func _/_(double, double) double
# func _/_(int, int) int
# func _/_(uint, uint) uint (Division Operator)
mutation DivisionOperator {
  evaluate(
    expressions: {
      int: "10 / 3" # 3 (integer division)
      uint: "10u / 3u" # 3u
      double: "10.0 / 4.0" # 2.5
      # decimal: decimal('10') / decimal('4') # decimal('2.5')
    }
  )
}
```

#### Modulo (%)

Calculates the remainder of integer or unsigned integer division.

```graphql
# func _%_(int, int) int
# func _%_(uint, uint) uint (Modulo Operator)
mutation ModuloOperator {
  evaluate(
    expressions: {
      int: "10 % 3" # 1
      uint: "10u % 3u" # 1u
    }
  )
}
```

#### Unary Negation (-)

Negates an integer or double value.

```graphql
# func -_(double) double
# func -_(int) int (Unary Negation)
mutation UnaryNegation {
  evaluate(
    expressions: {
      int: "-5" # -5
      double: "-3.14" # -3.14
    }
  )
}
```

### Comparison Operators

These operators compare two values and return a boolean result. They work across various compatible types (numbers, strings, timestamps, durations, decimals, money, bytes).

#### Equality (==)

Checks if two values are equal. For lists and maps, this typically performs a deep equality check.

```graphql
# func _==_(A, A) bool (Equality Operator)
mutation EqualityOperator {
  evaluate(
    expressions: {
      int: "5 == 5" # true
      string: "'hello' == 'hello'" # true
      bool: "true == !false" # true
      list: "[1, 2] == [1, 2]" # true (usually deep equality)
    }
  )
}
```

#### Inequality (!=)

Checks if two values are not equal.

```graphql
# func _!=_(A, A) bool (Inequality Operator)
mutation InequalityOperator {
  evaluate(
    expressions: {
      int: "5 != 10" # true
      string: "'hello' != 'world'" # true
      list: "[1] != [2]" # true
    }
  )
}
```

#### Less Than (<)

Checks if the left operand is strictly less than the right operand.

```graphql
# func _<_(...) bool (Comparison Operator)
mutation LessThanOperator {
  evaluate(
    expressions: {
      int: "5 < 10" # true
      double: "3.14 < 3.15" # true
      string: "'apple' < 'banana'" # true
      timestamp: "timestamp('2023-01-01T00:00:00Z') < timestamp('2023-01-02T00:00:00Z')" # true
      decimal: "decimal('1.2') < decimal('1.3')" # true
    }
  )
}
```

#### Less Than or Equal (<=)

Checks if the left operand is less than or equal to the right operand.

```graphql
# func _<=_(...) bool (Less Than or Equal Operator)
mutation LessThanOrEqualOperator {
  evaluate(
    expressions: {
      int: "5 <= 5" # true
      double: "3.14 <= 3.14" # true
      string: "'apple' <= 'apple'" # true
    }
  )
}
```

#### Greater Than (>)

Checks if the left operand is strictly greater than the right operand.

```graphql
# func _>_(...) bool (Greater Than Operator)
mutation GreaterThanOperator {
  evaluate(
    expressions: {
      int: "10 > 5" # true
      double: "3.15 > 3.14" # true
      string: "'banana' > 'apple'" # true
      decimal: "decimal('1.3') > decimal('1.2')" # true
    }
  )
}
```

#### Greater Than or Equal (>=)

Checks if the left operand is greater than or equal to the right operand.

```graphql
# func _>=_(...) bool (Greater Than or Equal Operator)
mutation GreaterThanOrEqualOperator {
  evaluate(
    expressions: {
      int: "5 >= 5" # true
      double: "3.15 >= 3.14" # true
      string: "'b' >= 'a'" # true
    }
  )
}
```

### Logical Operators

#### Logical AND (&&)

Returns true if both boolean operands are true, otherwise false. Short-circuits (does not evaluate the right operand if the left is false).

```graphql
# func _&&_(bool, bool) bool (Logical AND)
mutation LogicalAnd {
  evaluate(expressions: { result: "true && (1 < 2)" }) # Example: true
}
```

#### Logical OR (||)

Returns true if at least one boolean operand is true, otherwise false. Short-circuits (does not evaluate the right operand if the left is true).

```graphql
# func _||_(bool, bool) bool (Logical OR)
mutation LogicalOr {
  evaluate(expressions: { result: "false || (1 < 2)" }) # Example: true
}
```

#### Logical NOT (!)

Inverts a boolean value (true becomes false, false becomes true).

```graphql
# func !_(bool) bool (Logical NOT)
mutation LogicalNot {
  evaluate(expressions: { result: "!false" }) # Example: true
}
```

### Conditional (Ternary) Operator (?:)

Evaluates a boolean condition. If true, returns the second operand; if false, returns the third operand.

```graphql
# func _?_:_(bool, A, A) A
mutation TernaryOperator {
  evaluate(expressions: { result: "true ? 'yes' : 'no'" }) # Example: "yes"
}
```

### Collection Operators

#### In Operator (in)

Checks for membership. For lists, it checks if an element exists. For maps, it checks if a key exists. `@in` and `_in_` are alternative syntaxes.

```graphql
# func @in(A, list(A)) bool
# func in(A, list(A)) bool
# func _in_(A, list(A)) bool
mutation InOperatorList {
  evaluate(expressions: { result: "2 in [1, 2, 3]" }) # Example: true
}

# func @in(A, map(A, B)) bool
# func in(A, map(A, B)) bool
# func _in_(A, map(A, B)) bool
mutation InOperatorMap {
  evaluate(expressions: { result: "'b' in {'a': 1, 'b': 2}" }) # Example: true (checks for key existence)
}
```

#### Index Operator ([])

Accesses elements in lists by integer index or values in maps by key. Also works on optional lists/maps, returning an optional value. Accessing out-of-bounds index or non-existent key results in an error unless used on an optional type.

```graphql
# func _[_](list(A), int) A
mutation ListIndex {
  evaluate(expressions: { result: "[10, 20, 30][1]" }) # Example: 20
}

# func _[_](map(A, B), A) B
mutation MapIndex {
  evaluate(expressions: { result: "{'a': 1, 'b': 2}['a']" }) # Example: 1
}

# func _[_](optional_type(list(V)), int) optional_type(V)
mutation OptionalListIndex {
  evaluate(
    expressions: {
      present: "optional.of([10, 20])[0].value()" # 10
      absent: "optional.of([10, 20])[2].hasValue()" # false (index out of bounds returns empty optional)
    }
  )
}

# func _[_](optional_type(map(K, V)), K) optional_type(V)
mutation OptionalMapIndex {
  evaluate(
    expressions: {
      present: "optional.of({'a': 1})['a'].value()" # 1
      absent: "optional.of({'a': 1})['b'].hasValue()" # false (key not found returns empty optional)
    }
  )
}
```

#### Safe Index Operator ([?])

Accesses elements in lists or maps like the standard index operator, but *always* returns an optional type. Returns an empty optional instead of an error for out-of-bounds indices or non-existent keys. Works on both regular and optional collections.

```graphql
# func _[?_](list(V), int) optional_type(V)
mutation SafeListIndexPresent {
  evaluate(expressions: { result: "[10, 20][?0].value()" }) # Example: 10
}
mutation SafeListIndexAbsent {
  evaluate(expressions: { result: "[10, 20][?2].hasValue()" }) # Example: false
}

# func _[?_](optional_type(list(V)), int) optional_type(V)
mutation SafeOptionalListIndex {
  evaluate(
    expressions: {
      present: "optional.of([10, 20])[?0].value()" # 10
      absentIndex: "optional.of([10, 20])[?2].hasValue()" # false
      absentList: "optional.none()[?0].hasValue()" # false
    }
  )
}

# func _[?_](map(K, V), K) optional_type(V)
mutation SafeMapIndexPresent {
  evaluate(expressions: { result: "{'a': 1}[?'a'].value()" }) # Example: 1
}
mutation SafeMapIndexAbsent {
  evaluate(expressions: { result: "{'a': 1}[?'b'].hasValue()" }) # Example: false
}

# func _[?_](optional_type(map(K, V)), K) optional_type(V)
mutation SafeOptionalMapIndex {
  evaluate(
    expressions: {
      present: "optional.of({'a': 1})[?'a'].value()" # 1
      absentKey: "optional.of({'a': 1})[?'b'].hasValue()" # false
      absentMap: "optional.none()[?'a'].hasValue()" # false
    }
  )
}
```

#### Optional Field Access (.?)

Safely accesses fields on dynamic types (`dyn`). If the field exists, it returns an optional containing the field's value. If the field does not exist, it returns an empty optional instead of an error.

```graphql
# func _?._(dyn, string) optional_type(V)
mutation OptionalFieldAccess {
  evaluate(
    expressions: {
      present: "dyn({'a': 1}).?a.value()" # 1
      absent: "dyn({'a': 1}).?b.hasValue()" # false
    }
  )
}
```

## Optional Type Helpers

Functions for creating and working with optional types, which represent values that may or may not be present.

### optional.of

Creates an optional containing the given value.

```graphql
# func optional.of(V) optional_type(V)
mutation OptionalOf {
  evaluate(expressions: { result: "optional.of('value').value()" }) # Example: "value"
}
```

### optional.none

Creates an empty optional (representing no value).

```graphql
# func optional.none() optional_type(V)
mutation OptionalNone {
  evaluate(expressions: { result: "optional.none().hasValue()" }) # Example: false
}
```

### optional.ofNonZeroValue

Creates an optional containing the value if it's not the zero-value for its type (e.g., not 0 for int, not "" for string, not false for bool). Otherwise, creates an empty optional.

```graphql
# func optional.ofNonZeroValue(V) optional_type(V)
mutation OptionalOfNonZeroValueInt {
  evaluate(
    expressions: {
      zero: "optional.ofNonZeroValue(0).hasValue()" # false
      nonZero: "optional.ofNonZeroValue(5).hasValue()" # true
    }
  )
}
mutation OptionalOfNonZeroValueString {
  evaluate(
    expressions: {
      empty: "optional.ofNonZeroValue('').hasValue()" # false
      nonEmpty: "optional.ofNonZeroValue('hello').hasValue()" # true
    }
  )
}
```

### (optional) hasValue

Checks if the optional contains a value (returns true) or is empty (returns false).

```graphql
# func (optional_type(V)) hasValue() bool
mutation OptionalHasValue {
  evaluate(
    expressions: {
      present: "optional.of('hello').hasValue()" # true
      absent: "optional.none().hasValue()" # false
    }
  )
}
```

### (optional) value

Extracts the value from the optional. **Important:** This will cause an error if the optional is empty. Use `hasValue` to check first, or use `orValue`.

```graphql
# func (optional_type(V)) value() V
mutation OptionalValue {
  evaluate(expressions: { result: "optional.of('hello').value()" }) # Example: "hello" (Errors if optional is empty)
}
```

### (optional) or

Takes two optionals. Returns the first optional if it has a value, otherwise returns the second optional.

```graphql
# func (optional_type(V)) or(optional_type(V)) optional_type(V)
mutation OptionalOr {
  evaluate(
    expressions: {
      first: "optional.of(1).or(optional.of(2)).value()" # 1
      second: "optional.none().or(optional.of(2)).value()" # 2
      bothNone: "optional.none().or(optional.none()).hasValue()" # false
    }
  )
}
```

### (optional) orValue

Extracts the value from the optional if it's present. If the optional is empty, returns the provided default value instead.

```graphql
# func (optional_type(V)) orValue(V) V
mutation OptionalOrValue {
  evaluate(
    expressions: {
      has: "optional.of(10).orValue(0)" # 10
      none: "optional.none().orValue(0)" # 0
    }
  )
}
```

## String Manipulation

A comprehensive set of functions for working with strings, mirroring many functions from Go's `strings` package and common string methods.

### Case Conversion

#### strings.ToLower

Converts the entire string to lowercase, respecting Unicode rules.

```graphql
# func strings.ToLower(string) string
mutation StringsToLower {
  evaluate(expressions: { result: "strings.ToLower('HELLO WORLD')" }) # Example: "hello world"
}
```

#### strings.ToUpper

Converts the entire string to uppercase, respecting Unicode rules.

```graphql
# func strings.ToUpper(string) string
mutation StringsToUpper {
  evaluate(expressions: { result: "strings.ToUpper('lowercase')" }) # Example: "LOWERCASE"
}
```

#### strings.ToTitle

Converts the string to title case (often equivalent to uppercase for simple ASCII). Unicode rules apply.

```graphql
# func strings.ToTitle(string) string
mutation StringsToTitle {
  evaluate(expressions: { result: "strings.ToTitle('loud noises')" }) # Example: "LOUD NOISES" (Often same as ToUpper for ASCII)
}
```

#### strings.Title

Converts the string to title case, where the first letter of each word is capitalized. Word boundaries are Unicode-aware.

```graphql
# func strings.Title(string) string
# Note: Title casing rules are language specific and can be complex.
mutation StringsTitle {
  evaluate(expressions: { result: "strings.Title('war and peace')" }) # Example: "War And Peace" (simple case)
}
```

#### (string) lowerAscii

Converts only ASCII characters in the string to lowercase. Faster than `strings.ToLower` but doesn't handle non-ASCII characters.

```graphql
# func (string) lowerAscii() string
mutation StringLowerAscii {
  evaluate(expressions: { result: "'HELLO WORLD 123'.lowerAscii()" }) # Example: "hello world 123"
}
```

#### (string) upperAscii

Converts only ASCII characters in the string to uppercase. Faster than `strings.ToUpper` but doesn't handle non-ASCII characters.

```graphql
# func (string) upperAscii() string
mutation StringUpperAscii {
  evaluate(expressions: { result: "'hello world 123'.upperAscii()" }) # Example: "HELLO WORLD 123"
}
```

### Trimming Whitespace and Characters

#### strings.TrimSpace

Removes leading and trailing whitespace (as defined by Unicode) from the string.

```graphql
# func strings.TrimSpace(string) string
mutation StringsTrimSpace {
  evaluate(
    expressions: { result: "strings.TrimSpace(' \\t\\n Hello \\n\\t ') " }
  ) # Example: "Hello"
}
```

#### strings.Trim

Removes leading and trailing characters specified in the `cutset` string.

```graphql
# func strings.Trim(string, string) string
mutation StringsTrim {
  evaluate(expressions: { result: "strings.Trim('.,!Hello!,.', '.,!')" }) # Example: "Hello"
}
```

#### strings.TrimLeft

Removes leading characters specified in the `cutset` string.

```graphql
# func strings.TrimLeft(string, string) string
mutation StringsTrimLeft {
  evaluate(expressions: { result: "strings.TrimLeft('...Hello...', '.')" }) # Example: "Hello..."
}
```

#### strings.TrimRight

Removes trailing characters specified in the `cutset` string.

```graphql
# func strings.TrimRight(string, string) string
mutation StringsTrimRight {
  evaluate(expressions: { result: "strings.TrimRight('...Hello...', '.')" }) # Example: "...Hello"
}
```

#### strings.TrimPrefix

Removes the specified prefix from the beginning of the string, if present.

```graphql
# func strings.TrimPrefix(string, string) string
mutation StringsTrimPrefix {
  evaluate(expressions: { result: "strings.TrimPrefix('__main__', '__')" }) # Example: "main__"
}
```

#### strings.TrimSuffix

Removes the specified suffix from the end of the string, if present.

```graphql
# func strings.TrimSuffix(string, string) string
mutation StringsTrimSuffix {
  evaluate(expressions: { result: "strings.TrimSuffix('filename.txt', '.txt')" }) # Example: "filename"
}
```

#### (string) trim

Removes leading and trailing whitespace from the string (equivalent to `strings.TrimSpace`).

```graphql
# func (string) trim() string
mutation StringTrim {
  evaluate(expressions: { result: "'  whitespace  '.trim()" }) # Example: "whitespace"
}
```

### Searching and Indexing

#### strings.Contains

Checks if the string contains the specified substring.

```graphql
# func strings.Contains(string, string) bool
mutation StringsContains {
  evaluate(expressions: { result: "strings.Contains('banana', 'nan')" }) # Example: true
}
```

#### strings.ContainsAny

Checks if the string contains any character from the specified `chars` string.

```graphql
# func strings.ContainsAny(string, string) bool
mutation StringsContainsAny {
  evaluate(expressions: { result: "strings.ContainsAny('team', 'eiou')" }) # Example: true ('e' and 'a' are vowels)
}
```

#### strings.HasPrefix

Checks if the string starts with the specified prefix.

```graphql
# func strings.HasPrefix(string, string) bool
mutation StringsHasPrefix {
  evaluate(expressions: { result: "strings.HasPrefix('__main__', '__')" }) # Example: true
}
```

#### strings.HasSuffix

Checks if the string ends with the specified suffix.

```graphql
# func strings.HasSuffix(string, string) bool
mutation StringsHasSuffix {
  evaluate(expressions: { result: "strings.HasSuffix('image.jpg', '.jpg')" }) # Example: true
}
```

#### strings.Index

Finds the index of the first occurrence of the substring within the string. Returns -1 if not found.

```graphql
# func strings.Index(string, string) int
mutation StringsIndex {
  evaluate(expressions: { result: "strings.Index('banana', 'na')" }) # Example: 2
}
```

#### strings.IndexAny

Finds the index of the first occurrence of any character from the `chars` string. Returns -1 if no character is found.

```graphql
# func strings.IndexAny(string, string) int
mutation StringsIndexAny {
  evaluate(expressions: { result: "strings.IndexAny('chicken', 'aeiou')" }) # Example: 2 (index of 'i')
}
```

#### strings.LastIndex

Finds the index of the last occurrence of the substring within the string. Returns -1 if not found.

```graphql
# func strings.LastIndex(string, string) int
mutation StringsLastIndex {
  evaluate(expressions: { result: "strings.LastIndex('banana', 'na')" }) # Example: 4
}
```

#### strings.LastIndexAny

Finds the index of the last occurrence of any character from the `chars` string. Returns -1 if no character is found.

```graphql
# func strings.LastIndexAny(string, string) int
mutation StringsLastIndexAny {
  evaluate(expressions: { result: "strings.LastIndexAny('banana', 'ab')" }) # Example: 5 ('a' at index 5)
}
```

#### (string) contains

Checks if the string contains the specified substring (method form).

```graphql
# func (string) contains(string) bool
mutation StringContains {
  evaluate(expressions: { result: "'banana'.contains('nan')" }) # Example: true
}
```

#### (string) startsWith

Checks if the string starts with the specified prefix (method form).

```graphql
# func (string) startsWith(string) bool
mutation StringStartsWith {
  evaluate(expressions: { result: "'filename.txt'.startsWith('file')" }) # Example: true
}
```

#### (string) endsWith

Checks if the string ends with the specified suffix (method form).

```graphql
# func (string) endsWith(string) bool
mutation StringEndsWith {
  evaluate(expressions: { result: "'image.png'.endsWith('.png')" }) # Example: true
}
```

#### (string) indexOf

Finds the index of the first occurrence of the substring, optionally starting the search from a given index (method form).

```graphql
# func (string) indexOf(string) int
mutation StringIndexOfSimple {
  evaluate(expressions: { result: "'banana'.indexOf('na')" }) # Example: 2
}

# func (string) indexOf(string, int) int
mutation StringIndexOfWithStart {
  evaluate(expressions: { result: "'banana'.indexOf('na', 3)" }) # Example: 4 (start search from index 3)
}
```

#### (string) lastIndexOf

Finds the index of the last occurrence of the substring, optionally searching backwards from a given index (method form).

```graphql
# func (string) lastIndexOf(string) int
mutation StringLastIndexOfSimple {
  evaluate(expressions: { result: "'banana'.lastIndexOf('a')" }) # Example: 5
}

# func (string) lastIndexOf(string, int) int
mutation StringLastIndexOfWithStart {
  evaluate(expressions: { result: "'banana'.lastIndexOf('a', 4)" }) # Example: 3 (search backwards from index 4)
}
```

### Replacing Substrings

#### strings.Replace

Replaces the first `n` occurrences of `old` with `new`. If `n` is negative, all occurrences are replaced.

```graphql
# func strings.Replace(string, string, string, int) string
mutation StringsReplace {
  evaluate(expressions: { result: "strings.Replace('banana', 'a', 'o', 2)" }) # Example: "bonona" (replace first 2 'a's)
}
```

#### strings.ReplaceAll

Replaces all occurrences of `old` with `new`.

```graphql
# func strings.ReplaceAll(string, string, string) string
mutation StringsReplaceAll {
  evaluate(expressions: { result: "strings.ReplaceAll('banana', 'a', 'o')" }) # Example: "bonono"
}
```

#### (string) replace

Replaces occurrences of `old` with `new`. If `n` is provided and non-negative, replaces the first `n` occurrences. If `n` is omitted or negative, replaces all occurrences (method form).

```graphql
# func (string) replace(string, string) string
mutation StringReplaceSimple {
  evaluate(expressions: { result: "'banana'.replace('a', 'o')" }) # Example: "bonono" (replaces all)
}

# func (string) replace(string, string, int) string
mutation StringReplaceN {
  evaluate(expressions: { result: "'banana'.replace('a', 'o', 2)" }) # Example: "bonona" (replaces first 2)
}
```

### Splitting and Joining

#### strings.Split

Splits the string into a list of strings using the separator `sep`.

```graphql
# func (string) split(string) list(string)
mutation StringSplit {
  evaluate(expressions: { result: "'a-b-c'.split('-')" }) # Example: ["a", "b", "c"]
}
```

#### strings.SplitN

Splits the string by `sep` into at most `n` substrings. If `n` > 0, the last substring will contain the unsplit remainder.

```graphql
# func (string) split(string, int) list(string)
mutation StringSplitN {
  evaluate(expressions: { result: "'a-b-c'.split('-', 2)" }) # Example: ["a", "b-c"] (split into n substrings)
}
```

#### (list<string>) join

Concatenates a list of strings into a single string, optionally using a separator. Default separator is empty string.

```graphql
# func (list(string)) join() string
mutation ListStringJoinDefault {
  evaluate(expressions: { result: "['a', 'b', 'c'].join()" }) # Example: "abc"
}

# func (list(string)) join(string) string
mutation ListStringJoinSeparator {
  evaluate(expressions: { result: "['a', 'b', 'c'].join('-')" }) # Example: "a-b-c"
}
```

### Substrings and Characters

#### (string) substring

Extracts a portion of the string. With one argument `start`, takes characters from `start` to the end. With `start` and `end`, takes characters from `start` up to (but not including) `end`.

```graphql
# func (string) substring(int) string
mutation StringSubstringFrom {
  evaluate(expressions: { result: "'abcdef'.substring(2)" }) # Example: "cdef"
}

# func (string) substring(int, int) string
mutation StringSubstringRange {
  evaluate(expressions: { result: "'abcdef'.substring(1, 4)" }) # Example: "bcd" (exclusive end index)
}
```

#### (string) charAt

Returns the character (as a string) at the specified zero-based index.

```graphql
# func (string) charAt(int) string
mutation StringCharAt {
  evaluate(expressions: { result: "'hello'.charAt(1)" }) # Example: "e"
}
```

#### (string) take

Returns the first `n` characters of the string.

```graphql
# func (string) take(int) string
mutation StringTake {
  evaluate(expressions: { result: "'abcdef'.take(3)" }) # Example: "abc" (first n chars)
}
```

#### (string) drop

Returns the string with the first `n` characters removed (equivalent to `substring(n)`).

```graphql
# func (string) drop(int) string
mutation StringDrop {
  evaluate(expressions: { result: "'abcdef'.drop(2)" }) # Example: "cdef" (same as substring(n))
}
```

### Other String Functions

#### strings.Count

Counts the non-overlapping occurrences of a substring within the string.

```graphql
# func strings.Count(string, string) int
mutation StringsCount {
  evaluate(expressions: { result: "strings.Count('banana', 'a')" }) # Example: 3
}
```

#### strings.EqualFold

Compares two strings using case-insensitive comparison (Unicode-aware).

```graphql
# func strings.EqualFold(string, string) bool
mutation StringsEqualFold {
  evaluate(expressions: { result: "strings.EqualFold('GoLang', 'golang')" }) # Example: true
}
```

#### strings.Compare

Compares two strings lexicographically. Returns -1 if s1 < s2, 0 if s1 == s2, 1 if s1 > s2.

```graphql
# func strings.Compare(string, string) int
mutation StringsCompare {
  evaluate(
    expressions: {
      less: "strings.Compare('a', 'b')" # -1
      equal: "strings.Compare('a', 'a')" # 0
      greater: "strings.Compare('b', 'a')" # 1
    }
  )
}
```

#### strings.Repeat

Repeats the string `count` times.

```graphql
# func strings.Repeat(string, int) string
mutation StringsRepeat {
  evaluate(expressions: { result: "strings.Repeat('=', 5)" }) # Example: "====="
}
```

#### strings.ToValidUTF8

Replaces invalid UTF-8 byte sequences in the string with the specified replacement string.

```graphql
# func strings.ToValidUTF8(string, string) string
mutation StringsToValidUTF8 {
  evaluate(expressions: { result: "strings.ToValidUTF8('abc\\xffdef', '?')" }) # Example: "abc?def"
}
```

#### (string) format

Formats the string using placeholders (like `%s`, `%d`) and a list of dynamic arguments. Placeholder syntax depends on implementation (often similar to `sprintf`).

```graphql
# func (string) format(list(dyn)) string
# Note: format specifiers like %s, %d are implementation specific. Using a generic example.
mutation StringFormat {
  evaluate(
    expressions: {
      result: "'Hello, %s! You are %d.'.format(['World', dyn(30)])"
    }
  ) # Example: "Hello, World! You are 30." (Assuming %s, %d)
}
```

#### (string) reverse

Reverses the order of characters in the string.

```graphql
# func (string) reverse() string
mutation StringReverse {
  evaluate(expressions: { result: "'hello'.reverse()" }) # Example: "olleh"
}
```

#### (string) quote

Returns a double-quoted Go string literal representing the input string, with backslash escapes for control characters and quotes.

```graphql
# func strings.quote(string) string
mutation StringsQuote {
  evaluate(expressions: { result: "strings.quote('Hello \"World\"')" }) # Example: "\"Hello \\\"World\\\"\""
}
```

#### (string) size / size(string)

Returns the number of characters (runes) in the string.

```graphql
# func size(string) int / func (string) size() int
mutation SizeString {
  evaluate(expressions: { result: "size('hello')" }) # Example: 5
}
```

## Date & Time

Functions for working with timestamps, durations, and dates.

### Timestamp Manipulation

#### time.Now

Returns the current timestamp.

```graphql
# func time.Now() google.protobuf.Timestamp
mutation TimeNow {
  evaluate(expressions: { result: "string(time.Now())" }) # Example: Current timestamp string like "2023-10-27T10:30:00Z" (will vary)
}
```

#### (timestamp) getHours

Extracts the hour (0-23) from the timestamp, optionally in a specified timezone. Default is UTC.

```graphql
# func (google.protobuf.Timestamp) getHours() int
mutation TimestampGetHoursUTC {
  evaluate(
    expressions: { result: "timestamp('2023-10-27T15:04:05Z').getHours()" }
  ) # Example: 15
}

# func (google.protobuf.Timestamp) getHours(string) int
mutation TimestampGetHoursTimezone {
  evaluate(
    expressions: {
      result: "timestamp('2023-10-27T15:04:05Z').getHours('America/New_York')"
    }
  ) # Example: 11
}
```

#### (timestamp) getMinutes

Extracts the minute (0-59) from the timestamp, optionally in a specified timezone.

```graphql
# func (google.protobuf.Timestamp) getMinutes() int
mutation TimestampGetMinutesUTC {
  evaluate(
    expressions: { result: "timestamp('2023-10-27T15:04:05Z').getMinutes()" }
  ) # Example: 4
}

# func (google.protobuf.Timestamp) getMinutes(string) int
mutation TimestampGetMinutesTimezone {
  evaluate(
    expressions: {
      result: "timestamp('2023-10-27T15:04:05Z').getMinutes('America/New_York')"
    }
  ) # Example: 4
}
```

#### (timestamp) getSeconds

Extracts the second (0-59) from the timestamp, optionally in a specified timezone (though seconds are usually timezone-independent).

```graphql
# func (google.protobuf.Timestamp) getSeconds() int
mutation TimestampGetSecondsUTC {
  evaluate(
    expressions: {
      result: "timestamp('2023-10-27T15:04:05.123Z').getSeconds()"
    }
  ) # Example: 5
}

# func (google.protobuf.Timestamp) getSeconds(string) int
mutation TimestampGetSecondsTimezone {
  evaluate(
    expressions: {
      result: "timestamp('2023-10-27T15:04:05.123Z').getSeconds('America/New_York')"
    }
  ) # Example: 5 (seconds are usually timezone independent)
}
```

#### (timestamp) getMilliseconds

Extracts the millisecond (0-999) part of the timestamp.

```graphql
# func (google.protobuf.Timestamp) getMilliseconds() int
mutation TimestampGetMillisecondsUTC {
  evaluate(
    expressions: {
      result: "timestamp('2023-10-27T15:04:05.123Z').getMilliseconds()"
    }
  ) # Example: 123
}

# func (google.protobuf.Timestamp) getMilliseconds(string) int
mutation TimestampGetMillisecondsTimezone {
  evaluate(
    expressions: {
      result: "timestamp('2023-10-27T15:04:05.123Z').getMilliseconds('America/New_York')"
    }
  ) # Example: 123 (milliseconds are timezone independent)
}
```

#### (timestamp) getDayOfMonth / getDate

Extracts the day of the month (1-31) from the timestamp, optionally in a specified timezone. `getDate` is an alias.

```graphql
# func (google.protobuf.Timestamp) getDayOfMonth() int
# func (google.protobuf.Timestamp) getDate() int (alias)
mutation TimestampGetDayOfMonthUTC {
  evaluate(
    expressions: { result: "timestamp('2023-10-27T15:04:05Z').getDayOfMonth()" }
  ) # Example: 27
}

# func (google.protobuf.Timestamp) getDayOfMonth(string) int
# func (google.protobuf.Timestamp) getDate(string) int (alias)
mutation TimestampGetDayOfMonthTimezone {
  evaluate(
    expressions: {
      result: "timestamp('2023-11-01T02:00:00Z').getDayOfMonth('America/Los_Angeles')"
    }
  ) # Example: 31 (Oct 31st in PST)
}
```

#### (timestamp) getDayOfWeek

Extracts the day of the week (Sunday=0, Monday=1, ..., Saturday=6) from the timestamp, optionally in a specified timezone.

```graphql
# func (google.protobuf.Timestamp) getDayOfWeek() int
mutation TimestampGetDayOfWeekUTC {
  evaluate(
    expressions: { result: "timestamp('2023-10-27T00:00:00Z').getDayOfWeek()" }
  ) # Example: 5 (Friday)
}

# func (google.protobuf.Timestamp) getDayOfWeek(string) int
mutation TimestampGetDayOfWeekTimezone {
  evaluate(
    expressions: {
      result: "timestamp('2023-10-27T00:00:00Z').getDayOfWeek('America/Los_Angeles')"
    }
  ) # Example: 4 (Thursday in PST)
}
```

#### (timestamp) getDayOfYear

Extracts the day of the year (1-366) from the timestamp, optionally in a specified timezone.

```graphql
# func (google.protobuf.Timestamp) getDayOfYear() int
mutation TimestampGetDayOfYearUTC {
  evaluate(
    expressions: { result: "timestamp('2023-10-27T00:00:00Z').getDayOfYear()" }
  ) # Example: 300
}

# func (google.protobuf.Timestamp) getDayOfYear(string) int
mutation TimestampGetDayOfYearTimezone {
  evaluate(
    expressions: {
      result: "timestamp('2023-01-01T02:00:00Z').getDayOfYear('America/Los_Angeles')"
    }
  ) # Example: 365 (Dec 31st of previous year in PST)
}
```

#### (timestamp) getMonth

Extracts the month (January=0, February=1, ..., December=11) from the timestamp, optionally in a specified timezone. **Note:** This is 0-indexed, unlike `date.getMonth()`.

```graphql
# func (google.protobuf.Timestamp) getMonth() int
mutation TimestampGetMonthUTC {
  evaluate(
    expressions: { result: "timestamp('2023-10-27T15:04:05Z').getMonth()" }
  ) # Example: 9 (October)
}

# func (google.protobuf.Timestamp) getMonth(string) int
mutation TimestampGetMonthTimezone {
  evaluate(
    expressions: {
      result: "timestamp('2023-01-01T02:00:00Z').getMonth('America/Los_Angeles')"
    }
  ) # Example: 11 (December in PST)
}
```

#### (timestamp) getFullYear

Extracts the full year (e.g., 2023) from the timestamp, optionally in a specified timezone.

```graphql
# func (google.protobuf.Timestamp) getFullYear() int
mutation TimestampGetFullYearUTC {
  evaluate(
    expressions: { result: "timestamp('2023-10-27T15:04:05Z').getFullYear()" }
  ) # Example: 2023
}

# func (google.protobuf.Timestamp) getFullYear(string) int
mutation TimestampGetFullYearTimezone {
  evaluate(
    expressions: {
      result: "timestamp('2023-01-01T02:00:00Z').getFullYear('America/Los_Angeles')"
    }
  ) # Example: 2022 (Dec 31st of previous year in PST)
}
```

### Duration Manipulation

#### (duration) getHours

Calculates the total number of full hours represented by the duration.

```graphql
# func (google.protobuf.Duration) getHours() int
mutation DurationGetHours {
  evaluate(expressions: { result: "duration('3h15m').getHours()" }) # Example: 3
}
```

#### (duration) getMinutes

Calculates the total number of full minutes represented by the duration.

```graphql
# func (google.protobuf.Duration) getMinutes() int
mutation DurationGetMinutes {
  evaluate(expressions: { result: "duration('1h30m').getMinutes()" }) # Example: 90
}
```

#### (duration) getSeconds

Calculates the total number of seconds represented by the duration (including fractional parts if the underlying duration has them, but returns an int).

```graphql
# func (google.protobuf.Duration) getSeconds() int
mutation DurationGetSeconds {
  evaluate(expressions: { result: "duration('3m45s').getSeconds()" }) # Example: 225
}
```

#### (duration) getMilliseconds

Calculates the total number of milliseconds represented by the duration.

```graphql
# func (google.protobuf.Duration) getMilliseconds() int
mutation DurationGetMilliseconds {
  evaluate(expressions: { result: "duration('1s500ms').getMilliseconds()" }) # Example: 1500
}
```

### Date Manipulation

#### date.Today

Returns the current date, optionally in a specified timezone.

```graphql
# func date.Today() twisp.type.v1.Date
mutation DateTodayUTC {
  evaluate(expressions: { result: "string(date.Today())" }) # Example: Current date string in UTC e.g., "2023-10-27" (will vary)
}

# func date.Today(string) twisp.type.v1.Date
mutation DateTodayTimezone {
  evaluate(expressions: { result: "string(date.Today('America/New_York'))" }) # Example: Current date string in NY e.g., "2023-10-27" (will vary)
}
```

#### (date) getDay

Extracts the day of the month (1-31) from the date.

```graphql
# func (twisp.type.v1.Date) getDay() int
mutation DateGetDay {
  evaluate(expressions: { result: "date('2023-10-27').getDay()" }) # Example: 27
}
```

#### (date) getMonth

Extracts the month (January=1, February=2, ..., December=12) from the date. **Note:** This is 1-indexed, unlike `timestamp.getMonth()`.

```graphql
# func (twisp.type.v1.Date) getMonth() int
mutation DateGetMonth {
  evaluate(expressions: { result: "date('2023-10-27').getMonth()" }) # Example: 10 (October)
}
```

#### (date) getYear

Extracts the full year (e.g., 2023) from the date.

```graphql
# func (twisp.type.v1.Date) getYear() int
mutation DateGetYear {
  evaluate(expressions: { result: "date('2023-10-27').getYear()" }) # Example: 2023
}
```

#### (date) toString

Converts the Date object to its string representation (YYYY-MM-DD).

```graphql
# func (twisp.type.v1.Date) toString() string
mutation DateToString {
  evaluate(expressions: { result: "date('2023-10-27').toString()" }) # Example: "2023-10-27"
}
```

### Calendar Functions (cal)

#### cal.WeekOfYear

Calculates calendar-specific week information (likely week number, year). The exact map structure and behavior depends on the `int` parameters (e.g., first day of week, minimum days in first week).

```graphql
# func cal.WeekOfYear(google.protobuf.Timestamp, int, int) map(, )
# Note: The markdown signature for map is incomplete. Assuming map[string]int or similar based on common patterns.
# This function's exact map structure isn't defined, providing a basic call.
mutation CalWeekOfYear {
  evaluate(
    expressions: {
      result: "cal.WeekOfYear(timestamp('2023-01-01T00:00:00Z'), 1, 1)"
    }
  ) # Result depends on implementation details (e.g., {'year': 2023, 'week': 1})
}
```

#### cal.ISOWeekOfYear

Calculates the ISO 8601 week date (year and week number). The exact map structure returned depends on implementation.

```graphql
# func cal.ISOWeekOfYear(google.protobuf.Timestamp) map(, )
# Map structure unclear from signature. Basic call.
mutation CalISOWeekOfYear {
  evaluate(
    expressions: {
      result: "cal.ISOWeekOfYear(timestamp('2023-01-01T00:00:00Z'))"
    }
  ) # Result depends on implementation (e.g., {'year': 2022, 'week': 52})
}
```

#### cal.Quarter

Calculates the quarter (1-4) of the year for the given timestamp. The mode `int` might define quarter boundaries (e.g., fiscal vs calendar).

```graphql
# func cal.Quarter(google.protobuf.Timestamp, int) int
# Mode int might define quarter boundaries. Using 1 as default (calendar).
mutation CalQuarter {
  evaluate(
    expressions: { result: "cal.Quarter(timestamp('2023-07-01T00:00:00Z'), 1)" }
  ) # Example: 3
}
```

## Decimal Arithmetic

Functions for performing calculations with high-precision decimal numbers. These are crucial for financial calculations where floating-point inaccuracies are unacceptable.

### Basic Arithmetic

#### decimal.Add

Adds two decimal numbers.

```graphql
# func decimal.Add(twisp.type.v1.Decimal, twisp.type.v1.Decimal) twisp.type.v1.Decimal
mutation DecimalAdd {
  evaluate(
    expressions: {
      result: "string(decimal.Add(decimal('1.1'), decimal('2.2')))"
    }
  ) # Example: "3.3"
}
```

#### decimal.Sub

Subtracts the second decimal from the first.

```graphql
# func decimal.Sub(twisp.type.v1.Decimal, twisp.type.v1.Decimal) twisp.type.v1.Decimal
mutation DecimalSub {
  evaluate(
    expressions: {
      result: "string(decimal.Sub(decimal('1.23'), decimal('0.23')))"
    }
  ) # Example: "1"
}
```

#### decimal.Mul

Multiplies two decimal numbers.

```graphql
# func decimal.Mul(twisp.type.v1.Decimal, twisp.type.v1.Decimal) twisp.type.v1.Decimal
mutation DecimalMul {
  evaluate(
    expressions: { result: "string(decimal.Mul(decimal('1.5'), decimal('2')))" }
  ) # Example: "3" or "3.0"
}
```

#### decimal.Quo

Divides the first decimal by the second, up to a specified precision (`uint`). The exact division rules (rounding) depend on implementation context.

```graphql
# func decimal.Quo(twisp.type.v1.Decimal, twisp.type.v1.Decimal, uint) twisp.type.v1.Decimal
mutation DecimalQuo {
  evaluate(
    expressions: {
      result: "string(decimal.Quo(decimal('10'), decimal('3'), 4u))"
    }
  ) # Example: "3.3333" (precision 4)
}
```

#### decimal.QuoInteger

Performs integer division (truncates the result) between two decimals. The `uint` likely relates to precision context but might not directly affect the truncated result.

```graphql
# func decimal.QuoInteger(twisp.type.v1.Decimal, twisp.type.v1.Decimal, uint) twisp.type.v1.Decimal
mutation DecimalQuoInteger {
  evaluate(
    expressions: {
      result: "string(decimal.QuoInteger(decimal('10'), decimal('3'), 8u))"
    }
  ) # Example: "3"
}
```

#### decimal.Rem

Calculates the remainder of the division between two decimals. The `uint` might specify rounding mode for the remainder calculation.

```graphql
# func decimal.Rem(twisp.type.v1.Decimal, twisp.type.v1.Decimal, uint) twisp.type.v1.Decimal
# Mode uint might specify rounding for remainder. Using 0 as default.
mutation DecimalRem {
  evaluate(
    expressions: {
      result: "string(decimal.Rem(decimal('10.5'), decimal('3'), 0u))"
    }
  ) # Example: "1.5"
}
```

### Exponents, Roots, and Logs

#### decimal.Pow

Raises the first decimal to the power of the second decimal.

```graphql
# func decimal.Pow(twisp.type.v1.Decimal, twisp.type.v1.Decimal) twisp.type.v1.Decimal
mutation DecimalPow {
  evaluate(
    expressions: { result: "string(decimal.Pow(decimal('2'), decimal('3')))" }
  ) # Example: "8"
}
```

#### decimal.Exp

Raises the decimal to the power of an unsigned integer exponent.

```graphql
# func decimal.Exp(twisp.type.v1.Decimal, uint) twisp.type.v1.Decimal
mutation DecimalExp {
  evaluate(expressions: { result: "string(decimal.Exp(decimal('1.23'), 2u))" }) # Example: "1.5129" (1.23^2)
}
```

#### decimal.Sqrt

Calculates the square root of the decimal.

```graphql
# func decimal.Sqrt(twisp.type.v1.Decimal) twisp.type.v1.Decimal
mutation DecimalSqrt {
  evaluate(expressions: { result: "string(decimal.Sqrt(decimal('9')))" }) # Example: "3"
}
```

#### decimal.Cbrt

Calculates the cube root of the decimal.

```graphql
# func decimal.Cbrt(twisp.type.v1.Decimal) twisp.type.v1.Decimal
mutation DecimalCbrt {
  evaluate(expressions: { result: "string(decimal.Cbrt(decimal('27')))" }) # Example: "3"
}
```

#### decimal.Ln

Calculates the natural logarithm (base e) of the decimal.

```graphql
# func decimal.Ln(twisp.type.v1.Decimal) twisp.type.v1.Decimal
mutation DecimalLn {
  evaluate(expressions: { result: "string(decimal.Ln(decimal('10')))" }) # Example: Natural log of 10 (approx "2.302585...")
}
```

#### decimal.Log10

Calculates the base-10 logarithm of the decimal.

```graphql
# func decimal.Log10(twisp.type.v1.Decimal) twisp.type.v1.Decimal
mutation DecimalLog10 {
  evaluate(expressions: { result: "string(decimal.Log10(decimal('100')))" }) # Example: "2"
}
```

### Rounding and Precision

#### decimal.Round

Rounds the decimal to a specified number of places (`int`) using a named rounding mode (`string`, e.g., "half_up", "half_even").

```graphql
# func decimal.Round(twisp.type.v1.Decimal, string, int) twisp.type.v1.Decimal
# Rounding modes might include "half_up", "half_even", etc.
mutation DecimalRound {
  evaluate(
    expressions: {
      result: "string(decimal.Round(decimal('1.2345'), 'half_up', 2))"
    }
  ) # Example: "1.23"
}
```

#### decimal.Quantize

Rounds the decimal to have the same exponent (scale) as a specified number of places (`int`). Requires a rounding mode (`uint`, check implementation for specific mode values, e.g., 4=HALF_UP). This effectively sets the number of decimal places.

```graphql
# func decimal.Quantize(twisp.type.v1.Decimal, int, uint) twisp.type.v1.Decimal
# e.g. exponent -2 precision of 3.
mutation DecimalQuantize {
  evaluate(
    expressions: { result: "string(decimal.Quantize(decimal('1.236'), -2, 3u))" }
  ) # Example: "1.24"
}
```

#### decimal.Ceil

Rounds the decimal up to the nearest integer.

```graphql
# func decimal.Ceil(twisp.type.v1.Decimal) twisp.type.v1.Decimal
mutation DecimalCeil {
  evaluate(expressions: { result: "string(decimal.Ceil(decimal('1.23')))" }) # Example: "2"
}
```

#### decimal.Reduce

Removes trailing zeros from the fractional part of the decimal without changing its value.

```graphql
# func decimal.Reduce(twisp.type.v1.Decimal) twisp.type.v1.Decimal
mutation DecimalReduce {
  evaluate(expressions: { result: "string(decimal.Reduce(decimal('1.2300')))" }) # Example: "1.23"
}
```

### Comparison and Sign

#### decimal.Cmp

Compares two decimals. Returns a decimal value: -1 if d1 < d2, 0 if d1 == d2, 1 if d1 > d2.

```graphql
# func decimal.Cmp(twisp.type.v1.Decimal, twisp.type.v1.Decimal) twisp.type.v1.Decimal
# Returns -1, 0, or 1 as Decimal
mutation DecimalCmp {
  evaluate(
    expressions: {
      less: "string(decimal.Cmp(decimal('1'), decimal('2')))" # "-1"
      equal: "string(decimal.Cmp(decimal('2'), decimal('2')))" # "0"
      greater: "string(decimal.Cmp(decimal('3'), decimal('2')))" # "1"
    }
  )
}
```

#### decimal.Abs

Returns the absolute value of the decimal.

```graphql
# func decimal.Abs(twisp.type.v1.Decimal) twisp.type.v1.Decimal
mutation DecimalAbs {
  evaluate(expressions: { result: "string(decimal.Abs(decimal('-1.23')))" }) # Example: "1.23"
}
```

#### decimal.Neg

Returns the negation of the decimal.

```graphql
# func decimal.Neg(twisp.type.v1.Decimal) twisp.type.v1.Decimal
mutation DecimalNeg {
  evaluate(expressions: { result: "string(decimal.Neg(decimal('1.23')))" }) # Example: "-1.23"
}
```

### Conversion

#### (decimal) toString

Converts the Decimal object to its string representation.

```graphql
# func (twisp.type.v1.Decimal) toString() string
mutation DecimalToString {
  evaluate(expressions: { result: "decimal('123.456').toString()" }) # Example: "123.456"
}
```

## Money Arithmetic

Functions for performing calculations with Money objects, ensuring currency consistency.

### money.Add

Adds two Money objects. Requires both operands to have the same currency.

```graphql
# func money.Add(twisp.type.v1.Money, twisp.type.v1.Money) twisp.type.v1.Money
mutation MoneyAdd {
  evaluate(
    expressions: {
      result: "money.Add(money('10', 'USD'), money('5.50', 'USD'))"
    }
  ) # Example: "15.50 USD"
}
```

### money.Sub

Subtracts the second Money object from the first. Requires both operands to have the same currency.

```graphql
# func money.Sub(twisp.type.v1.Money, twisp.type.v1.Money) twisp.type.v1.Money
mutation MoneySub {
  evaluate(
    expressions: {
      result: "money.Sub(money('10.00', 'USD'), money('5.50', 'USD'))"
    }
  ) # Example: "4.50 USD"
}
```

### money.Mul

Multiplies a Money object by a scalar factor (provided as a string representing a decimal).

```graphql
# func money.Mul(twisp.type.v1.Money, string) twisp.type.v1.Money
mutation MoneyMul {
  evaluate(
    expressions: { result: "money.Mul(money('10.50', 'USD'), '2')" }
  ) # Example: "21.00 USD"
}
```

### money.Div

Divides a Money object by a scalar factor (provided as a string representing a decimal).

```graphql
# func money.Div(twisp.type.v1.Money, string) twisp.type.v1.Money
mutation MoneyDiv {
  evaluate(
    expressions: { result: "money.Div(money('21.00', 'USD'), '2')" }
  ) # Example: "10.50 USD"
}
```

### (money) currency

Returns the currency code of the Money object as a string.

```graphql
# func (twisp.type.v1.Money) currency() string
mutation MoneyCurrency {
  evaluate(expressions: { result: "money('10.00', 'EUR').currency()" }) # Example: "EUR"
}
```

### (money) decimal

Converts the Money object to its Decimal value (amount only, currency is discarded).

```graphql
# func (twisp.type.v1.Money) decimal() twisp.type.v1.Decimal
mutation MoneyDecimalConversion {
  evaluate(expressions: {
      fromMoneyMethod: "string(money('50.00', 'EUR').decimal())" # "50.00"
  })
}
```

## Math Functions

Standard mathematical functions operating primarily on `double` values. These mirror functions from Go's `math` package.

### Basic Operations

#### math.Abs

Returns the absolute value of a double.

```graphql
# func math.Abs(double) double
mutation MathAbs {
  evaluate(expressions: { result: "math.Abs(-5.5)" }) # Example: 5.5
}
```

#### math.Max

Returns the larger of two double values.

```graphql
# func math.Max(double, double) double
mutation MathMax {
  evaluate(expressions: { result: "math.Max(5.0, 3.0)" }) # Example: 5.0
}
```

#### math.Min

Returns the smaller of two double values.

```graphql
# func math.Min(double, double) double
mutation MathMin {
  evaluate(expressions: { result: "math.Min(5.0, 3.0)" }) # Example: 3.0
}
```

#### math.Dim

Returns the maximum of `x-y` or `0`. (Positive difference).

```graphql
# func math.Dim(double, double) double
mutation MathDim {
  evaluate(
    expressions: {
      pos: "math.Dim(5.0, 2.0)" # 3.0 (x-y)
      neg: "math.Dim(2.0, 5.0)" # 0.0 (max(x-y, 0))
    }
  )
}
```

#### math.Copysign

Returns a value with the magnitude of the first argument and the sign of the second argument.

```graphql
# func math.Copysign(double, double) double
mutation MathCopysign {
  evaluate(expressions: { result: "math.Copysign(5.0, -1.0)" }) # Example: -5.0
}
```

#### math.Signbit

Reports whether the number is negative or negative zero.

```graphql
# func math.Signbit(double) bool
mutation MathSignbit {
  evaluate(
    expressions: {
      neg: "math.Signbit(-5.0)" # true
      pos: "math.Signbit(5.0)" # false
    }
  )
}
```

### Rounding and Truncation

#### math.Ceil

Returns the least integer value greater than or equal to the input.

```graphql
# func math.Ceil(double) double
mutation MathCeil {
  evaluate(expressions: { result: "math.Ceil(3.1)" }) # Example: 4.0
}
```

#### math.Floor

Returns the greatest integer value less than or equal to the input.

```graphql
# func math.Floor(double) double
mutation MathFloor {
  evaluate(expressions: { result: "math.Floor(3.9)" }) # Example: 3.0
}
```

#### math.Round

Returns the nearest integer, rounding half away from zero.

```graphql
# func math.Round(double) double
mutation MathRound {
  evaluate(expressions: { result: "math.Round(3.5)" }) # Example: 4.0
}
```

#### math.RoundToEven

Returns the nearest integer, rounding ties to the nearest even integer (Banker's rounding).

```graphql
# func math.RoundToEven(double) double
mutation MathRoundToEven {
  evaluate(
    expressions: {
      halfUp: "math.RoundToEven(2.5)" # 2.0 (rounds to nearest even integer)
      halfDown: "math.RoundToEven(3.5)" # 4.0
    }
  )
}
```

#### math.Trunc

Returns the integer part of the double (truncates towards zero).

```graphql
# func math.Trunc(double) double
mutation MathTrunc {
  evaluate(expressions: { result: "math.Trunc(3.14159)" }) # Example: 3.0
}
```

### Powers, Roots, and Logs

#### math.Pow

Returns x**y, the base-x exponential of y.

```graphql
# func math.Pow(double, double) double
mutation MathPow {
  evaluate(expressions: { result: "math.Pow(2.0, 3.0)" }) # Example: 8.0
}
```

#### math.Pow10

Returns 10**n, the base-10 exponential of n (n is int).

```graphql
# func math.Pow10(int) double
mutation MathPow10 {
  evaluate(expressions: { result: "math.Pow10(3)" }) # Example: 1000.0
}
```

#### math.Sqrt

Returns the square root of x.

```graphql
# func math.Sqrt(double) double
mutation MathSqrt {
  evaluate(expressions: { result: "math.Sqrt(9.0)" }) # Example: 3.0
}
```

#### math.Cbrt

Returns the cube root of x.

```graphql
# func math.Cbrt(double) double
mutation MathCbrt {
  evaluate(expressions: { result: "math.Cbrt(27.0)" }) # Example: 3.0
}
```

#### math.Log

Returns the natural logarithm (base e) of x.

```graphql
# func math.Log(double) double
mutation MathLog {
  evaluate(expressions: { result: "math.Log(10.0)" }) # Example: ~2.302585 (natural log)
}
```

#### math.Log10

Returns the base-10 logarithm of x.

```graphql
# func math.Log10(double) double
mutation MathLog10 {
  evaluate(expressions: { result: "math.Log10(100.0)" }) # Example: 2.0
}
```

#### math.Log2

Returns the base-2 logarithm of x.

```graphql
# func math.Log2(double) double
mutation MathLog2 {
  evaluate(expressions: { result: "math.Log2(8.0)" }) # Example: 3.0
}
```

#### math.Log1p

Returns the natural logarithm of 1+x. More accurate than `Log(1+x)` for small x.

```graphql
# func math.Log1p(double) double
mutation MathLog1p {
  evaluate(expressions: { result: "math.Log1p(0.0)" }) # Example: 0.0 (ln(1+0))
}
```

#### math.Exp

Returns e**x, the base-e exponential of x.

```graphql
# func math.Exp(double) double
mutation MathExp {
  evaluate(expressions: { result: "math.Exp(1.0)" }) # Example: ~2.71828 (e^1)
}
```

#### math.Exp2

Returns 2**x, the base-2 exponential of x.

```graphql
# func math.Exp2(double) double
mutation MathExp2 {
  evaluate(expressions: { result: "math.Exp2(3.0)" }) # Example: 8.0 (2^3)
}
```

#### math.Expm1

Returns e**x - 1. More accurate than `Exp(x) - 1` for small x.

```graphql
# func math.Expm1(double) double
mutation MathExpm1 {
  evaluate(expressions: { result: "math.Expm1(0.0)" }) # Example: 0.0 (exp(0)-1)
}
```

### Trigonometry and Geometry

#### math.Hypot

Returns Sqrt(p*p + q*q), the length of the hypotenuse of a right triangle with legs p and q.

```graphql
# func math.Hypot(double, double) double
mutation MathHypot {
  evaluate(expressions: { result: "math.Hypot(3.0, 4.0)" }) # Example: 5.0
}
```

### Other Math Functions

#### math.Mod

Returns the floating-point remainder of x/y. The result has the same sign as x.

```graphql
# func math.Mod(double, double) double
mutation MathMod {
  evaluate(expressions: { result: "math.Mod(10.0, 3.0)" }) # Example: 1.0
}
```

#### math.Remainder

Returns the IEEE 754 floating-point remainder of x/y.

```graphql
# func math.Remainder(double, double) double
mutation MathRemainder {
  evaluate(expressions: { result: "math.Remainder(10.5, 3.0)" }) # Example: 1.5
}
```

#### math.FMA

Returns x * y + z, computed with only one rounding. (Fused Multiply-Add).

```graphql
# func math.FMA(double, double, double) double
mutation MathFMA {
  evaluate(expressions: { result: "math.FMA(2.0, 3.0, 4.0)" }) # Example: 10.0 (2*3 + 4)
}
```

## Random Number Generation (rand)

Functions for generating pseudo-random numbers, based on Go's `math/rand`.

### Integers

#### rand.Int

Returns a non-negative pseudo-random 64-bit integer.

```graphql
# func rand.Int() int
mutation RandInt {
  evaluate(expressions: { result: "rand.Int()" }) # Example: A random int64 value (e.g., 1234567890123456789)
}
```

#### rand.Int31

Returns a non-negative pseudo-random 31-bit integer as an int.

```graphql
# func rand.Int31() int
mutation RandInt31 {
  evaluate(expressions: { result: "rand.Int31()" }) # Example: Random non-negative int32 (e.g., 1073741823)
}
```

#### rand.Int63

Returns a non-negative pseudo-random 63-bit integer as an int.

```graphql
# func rand.Int63() int
mutation RandInt63 {
  evaluate(expressions: { result: "rand.Int63()" }) # Example: Random non-negative int64 (e.g., 4611686018427387903)
}
```

#### rand.Intn

Returns a non-negative pseudo-random integer in the range [0, n). Panics if n <= 0.

```graphql
# func rand.Intn(int) int
mutation RandIntn {
  evaluate(expressions: { result: "rand.Intn(100)" }) # Example: Random int between 0 and 99
}
```

#### rand.Int31n

Returns a non-negative pseudo-random 31-bit integer in the range [0, n). Panics if n <= 0.

```graphql
# func rand.Int31n(int) int
mutation RandInt31n {
  evaluate(expressions: { result: "rand.Int31n(1000)" }) # Random int32 between 0 and 999
}
```

#### rand.Int63n

Returns a non-negative pseudo-random 63-bit integer in the range [0, n). Panics if n <= 0.

```graphql
# func rand.Int63n(int) int
mutation RandInt63n {
  evaluate(expressions: { result: "rand.Int63n(100)" }) # Example: A random int64 between 0 and 99
}
```

### Unsigned Integers

#### rand.Uint32

Returns a pseudo-random 32-bit value as a uint.

```graphql
# func rand.Uint32() uint
mutation RandUint32 {
  evaluate(expressions: { result: "rand.Uint32()" }) # Example: Random uint32 (e.g., 2147483647u)
}
```

#### rand.Uint64

Returns a pseudo-random 64-bit value as a uint.

```graphql
# func rand.Uint64() uint
mutation RandUint64 {
  evaluate(expressions: { result: "rand.Uint64()" }) # Example: A random uint64 value (e.g., 9223372036854775807u)
}
```

### Floating-Point Numbers

#### rand.Float32

Returns a pseudo-random float32 in [0.0, 1.0). Returned as a double.

```graphql
# func rand.Float32() double
mutation RandFloat32 {
  evaluate(expressions: { result: "rand.Float32()" }) # Example: Random float32 (returned as double) between 0.0 and 1.0 (e.g., 0.1234567)
}
```

#### rand.Float64

Returns a pseudo-random float64 in [0.0, 1.0).

```graphql
# func rand.Float64() double
mutation RandFloat64 {
  evaluate(expressions: { result: "rand.Float64()" }) # Example: Random float64 between 0.0 and 1.0 (e.g., 0.987654321)
}
```

#### rand.NormFloat64

Returns a normally distributed float64 value with mean 0 and standard deviation 1.

```graphql
# func rand.NormFloat64() double
mutation RandNormFloat64 {
  evaluate(expressions: { result: "rand.NormFloat64()" }) # Example: Random normally distributed float64 (e.g., -0.54321)
}
```

#### rand.ExpFloat64

Returns an exponentially distributed float64 value with rate parameter 1.

```graphql
# func rand.ExpFloat64() double
mutation RandExpFloat64 {
  evaluate(expressions: { result: "rand.ExpFloat64()" }) # Example: Random exponentially distributed float64 (e.g., 1.2345)
}
```

## Hashing & Checksums

Functions for calculating various hash digests. These typically operate on `bytes`.

### MD5

Calculates the MD5 hash (128-bit). **Note:** MD5 is cryptographically broken and should not be used for security purposes.

```graphql
# func md5.Sum(bytes) bytes
mutation Md5Sum {
  evaluate(expressions: { result: "hex.EncodeToString(md5.Sum(b'hello'))" }) # Example: "5d41402abc4b2a76b9719d911017c592"
}
```

### SHA-1

Calculates the SHA-1 hash (160-bit). **Note:** SHA-1 is also considered cryptographically weak and should be avoided for security contexts.

```graphql
# func sha1.Sum(bytes) bytes
mutation Sha1Sum {
  evaluate(expressions: { result: "hex.EncodeToString(sha1.Sum(b'hello'))" }) # Example: "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"
}
```

### SHA-256

Calculates the SHA-256 hash (256-bit). A standard secure hashing algorithm.

```graphql
# func sha256.Sum256(bytes) bytes
mutation Sha256Sum256 {
  evaluate(
    expressions: { result: "hex.EncodeToString(sha256.Sum256(b'hello'))" }
  ) # Example: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
}
```

### SHA-512

Calculates the SHA-512 hash (512-bit). A standard secure hashing algorithm.

```graphql
# func sha512.Sum512(bytes) bytes
mutation Sha512Sum512 {
  evaluate(
    expressions: { result: "hex.EncodeToString(sha512.Sum512(b'hello'))" }
  ) # Example: "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043"
}
```

### System Hashes

#### system.hash.xx

Calculates a non-cryptographic xxHash64 hash, returning a `uint`. Useful for fast hashing where collision resistance is less critical than speed.

```graphql
# func system.hash.xx(bytes) uint
mutation SystemHashXx {
  evaluate(expressions: { result: "system.hash.xx(b'data')" }) # Example: xxHash64 uint result (e.g., 17241709254074915385u)
}
```

#### system.hash.jump

Implements the Jump Consistent Hash algorithm. Given a `uint` key and an `int` number of buckets, returns a bucket index (`int` from 0 to num_buckets-1) consistently assigned to that key. Useful for distributing items across shards/buckets.

```graphql
# func system.hash.jump(uint, int) int
mutation SystemHashJump {
  evaluate(expressions: { result: "system.hash.jump(123456789u, 10)" }) # Example: Jump consistent hash bucket index (0-9) (e.g., 7)
}
```

## Encoding & Decoding

Functions for converting data between different representations (e.g., bytes to hex, base64).

### Hexadecimal

#### hex.EncodeToString

Encodes a byte sequence into its hexadecimal string representation.

```graphql
# func hex.EncodeToString(bytes) string
mutation HexEncodeToString {
  evaluate(
    expressions: { result: "hex.EncodeToString(b'\\xDE\\xAD\\xBE\\xEF')" }
  ) # Example: "deadbeef"
}
```

#### hex.DecodeString

Decodes a hexadecimal string representation back into a byte sequence.

```graphql
# func hex.DecodeString(string) bytes
mutation HexDecodeString {
  evaluate(expressions: { result: "string(hex.DecodeString('68656c6c6f'))" }) # Example: "hello" (bytes converted back to string for display)
}
```

### Base64

#### base64.encode

Encodes a byte sequence into its Base64 string representation (standard encoding).

```graphql
# func base64.encode(bytes) string
mutation Base64Encode {
  evaluate(expressions: { result: "base64.encode(b'hello world')" }) # Example: "aGVsbG8gd29ybGQ="
}
```

#### base64.decode

Decodes a Base64 string representation back into a byte sequence.

```graphql
# func base64.decode(string) bytes
mutation Base64Decode {
  evaluate(expressions: { result: "string(base64.decode('aGVsbG8gd29ybGQ='))" }) # Example: "hello world" (bytes converted back to string)
}
```

### HTML Escaping

#### html.EscapeString

Escapes special HTML characters (`<`, `>`, `&`, `"`, `'`) into their corresponding HTML entities. Useful for preventing XSS when embedding strings in HTML.

```graphql
# func html.EscapeString(string) string
mutation HtmlEscapeString {
  evaluate(
    expressions: {
      result: "html.EscapeString('<script>alert(\"XSS\")</script>')"
    }
  ) # Example: "&lt;script&gt;alert(&#34;XSS&#34;)&lt;/script&gt;"
}
```

#### html.UnescapeString

Unescapes HTML entities (like `&lt;`, `&amp;`, `&#34;`) back into their original characters.

```graphql
# func html.UnescapeString(string) string
mutation HtmlUnescapeString {
  evaluate(
    expressions: {
      result: "html.UnescapeString('&lt;tag&gt; &amp; &quot;quote&quot;')"
    }
  ) # Example: "<tag> & \"quote\""
}
```

### URL Escaping

#### url.PathEscape

Escapes a string so it can be safely placed in a URL path segment. Spaces become `%20`, `/` remains unescaped.

```graphql
# func url.PathEscape(string) string
mutation UrlPathEscape {
  evaluate(expressions: { result: "url.PathEscape('a b/c?d=e')" }) # Example: "a%20b/c%3Fd=e"
}
```

#### url.QueryEscape

Escapes a string so it can be safely placed in a URL query parameter value. Spaces become `+` (or `%20` depending on context/implementation).

```graphql
# func url.QueryEscape(string) string
mutation UrlQueryEscape {
  evaluate(
    expressions: { result: "url.QueryEscape('key=value&another key=a/b')" }
  ) # Example: "key%3Dvalue%26another+key%3Da%2Fb"
}
```

### JSON Marshaling

#### json.Marshal

Converts a CEL value (often a map or list, represented as `dyn` or `google.protobuf.Any`) into its JSON byte representation.

```graphql
# func json.Marshal(google.protobuf.Any) bytes
# Requires a proto Any type, harder to demo simply. Example assumes a map can be marshaled.
mutation JsonMarshal {
  evaluate(expressions: { result: "string(json.Marshal({'key': 'value'}))" }) # Example: "{'key': 'value'}" (or similar JSON bytes as string)
}
```

## Path Manipulation (path)

Functions for working with slash-separated file paths, mirroring Go's `path` package.

### path.Base

Returns the last element of the path (the filename or directory name).

```graphql
# func path.Base(string) string
mutation PathBase {
  evaluate(expressions: { result: "path.Base('/usr/local/file.txt')" }) # Example: "file.txt"
}
```

### path.Clean

Returns the shortest path name equivalent to the path by purely lexical processing. It applies the following rules iteratively until no further processing can be done:
1. Replace multiple slashes with a single slash.
2. Eliminate each `.` path name element (the current directory).
3. Eliminate each inner `..` path name element (the parent directory) along with the preceding non-`..` element.
4. Eliminate `..` elements that begin a rooted path: that is, replace `/..` by `/` at the beginning of a path.

```graphql
# func path.Clean(string) string
mutation PathClean {
  evaluate(expressions: { result: "path.Clean('/a/b/../c//d')" }) # Example: "/a/c/d"
}
```

### path.Dir

Returns all but the last element of the path, typically the path's directory.

```graphql
# func path.Dir(string) string
mutation PathDir {
  evaluate(expressions: { result: "path.Dir('/usr/local/bin')" }) # Example: "/usr/local"
}
```

### path.Ext

Returns the file name extension used by path. The extension is the suffix beginning at the final dot in the final element of path; it is empty if there is no dot.

```graphql
# func path.Ext(string) string
mutation PathExt {
  evaluate(expressions: { result: "path.Ext('archive.tar.gz')" }) # Example: ".gz"
}
```

### path.IsAbs

Reports whether the path is absolute (begins with `/`).

```graphql
# func path.IsAbs(string) bool
mutation PathIsAbs {
  evaluate(
    expressions: {
      absolute: "path.IsAbs('/home/user')" # true
      relative: "path.IsAbs('data/file')" # false
    }
  )
}
```

## UUID Generation & Manipulation (uuid)

Functions specifically for creating UUIDs.

### uuid.New

Generates a new Version 4 (random) UUID.

```graphql
# func uuid.New() twisp.type.v1.UUID
mutation UuidNew {
  evaluate(expressions: { result: "string(uuid.New())" }) # Example: A new Version 4 UUID string (e.g., "f47ac10b-58cc-4372-a567-0e02b2c3d479")
}
```

### uuid.NewMD5

Generates a Version 3 (MD5 hash-based) UUID from a namespace UUID and a name (bytes).

```graphql
# func uuid.NewMD5(twisp.type.v1.UUID, bytes) twisp.type.v1.UUID
mutation UuidNewMD5 {
  evaluate(expressions: { result: "string(uuid.NewMD5(uuid.New(), b'data'))" }) # Example: MD5-based UUID string (e.g., "a1b...")
}
```

### uuid.NewSHA1

Generates a Version 5 (SHA-1 hash-based) UUID from a namespace UUID and a name (bytes).

```graphql
# func uuid.NewSHA1(twisp.type.v1.UUID, bytes) twisp.type.v1.UUID
mutation UuidNewSHA1 {
  evaluate(expressions: { result: "string(uuid.NewSHA1(uuid.New(), b'data'))" }) # Example: SHA1-based UUID string (e.g., "c2d...")
}
```

### (uuid) toString

Converts the UUID object to its standard string representation.

```graphql
# func (twisp.type.v1.UUID) toString() string
mutation UuidToString {
  evaluate(
    expressions: {
      result: "uuid('123e4567-e89b-12d3-a456-426614174000').toString()"
    }
  ) # Example: "123e4567-e89b-12d3-a456-426614174000"
}
```

## Finance Functions (finance)

Specialized functions for financial calculations. Examples provided are basic calls; understanding the financial formulas is necessary for correct usage. *Note: The exact signatures and behavior might vary.*

### finance.PrincipalPayment

Calculates the principal portion of a loan payment for a given period. Arguments typically include rate per period, payment period number, total number of periods, present value (loan amount), future value, and payment type (0 for end of period, 1 for beginning).

```graphql
# func finance.PrincipalPayment(double, int, int, double, double, int) double
mutation FinancePrincipalPayment {
  evaluate(
    expressions: {
      # Example: Principal for period 1 of a $200k, 30yr (360mo) loan at 5% APR (0.05/12 monthly)
      result: "finance.PrincipalPayment(0.05/12.0, 1, 360, 200000.0, 0.0, 0)"
    }
  ) # Example financial calculation result (e.g., ~239.98)
}
```

### finance.InterestPayment

Calculates the interest portion of a loan payment for a given period. Arguments are similar to `PrincipalPayment`.

```graphql
# func finance.InterestPayment(double, int, int, double, double, int) double
mutation FinanceInterestPayment {
  evaluate(
    expressions: {
       # Example: Interest for period 1 of a $200k, 30yr (360mo) loan at 5% APR (0.05/12 monthly)
      result: "finance.InterestPayment(0.05/12.0, 1, 360, 200000.0, 0.0, 0)"
    }
  ) # Example financial calculation result (e.g., ~833.33)
}
```

### finance.Payment

Calculates the total payment (principal + interest) per period for a loan or annuity. Arguments include rate per period, number of periods, present value, future value, and payment type.

```graphql
# func finance.Payment(double, int, double, double, int) double
mutation FinancePayment {
  evaluate(
    expressions: {
      # Example: Total payment for a $200k, 30yr (360mo) loan at 5% APR (0.05/12 monthly)
      result: "finance.Payment(0.05/12.0, 360, 200000.0, 0.0, 0)"
    }
  ) # Example financial calculation result (e.g., ~1073.31)
}
```

*(Other finance functions like Rate, PresentValue, FutureValue, Depreciation, etc., would follow a similar pattern with example calls based on their expected arguments)*

## System & Utility Functions

Miscellaneous functions for tasks like printing, type checking, and collection manipulation.

### print

Outputs the provided message and value to the system logs (side-effect) and returns the value itself. Useful for debugging CEL expressions.

```graphql
# func print(string, google.protobuf.Any) google.protobuf.Any
# Print is for side-effects (logging), returns its second argument.
mutation PrintExample {
  evaluate(expressions: { result: "print('Input value:', 123)" }) # Example: 123 (and logs 'Input value: 123' server-side)
}
```

### system.mergeMap

Merges two maps (represented as `dyn`). Keys from the second map overwrite keys in the first map if they overlap.

```graphql
# func system.mergeMap(dyn, dyn) dyn
mutation SystemMergeMap {
  # Merges second map into first, overwriting common keys
  evaluate(
    expressions: {
      result: "system.mergeMap({'a': 1, 'b': 2}, {'b': 3, 'c': 4})"
    }
  ) # Example: {'a': 1, 'b': 3, 'c': 4}
}
```

### system.unique

Removes duplicate UUIDs from a list of UUIDs. The order of the remaining unique elements is not guaranteed.

```graphql
# func system.unique(list(twisp.type.v1.UUID)) list(twisp.type.v1.UUID)
mutation SystemUnique {
  evaluate(
    expressions: {
      # Constructing UUIDs inline for example
      ids: "[uuid('11111111-1111-1111-1111-111111111111'), uuid('22222222-2222-2222-2222-222222222222'), uuid('11111111-1111-1111-1111-111111111111')]"
      uniqueIds: "system.unique([uuid('11111111-1111-1111-1111-111111111111'), uuid('22222222-2222-2222-2222-222222222222'), uuid('11111111-1111-1111-1111-111111111111')])"
      # uniqueIds result should be the first two UUIDs (order might vary)
      # e.g., [uuid('11111111-1111-1111-1111-111111111111'), uuid('22222222-2222-2222-2222-222222222222')]
    }
  )
}
```

### Size Functions

Calculates the size (number of elements or characters) of various types.

```graphql
# func size(bytes) int / func (bytes) size() int
mutation SizeBytes {
  evaluate(expressions: { result: "size(b'hello')" }) # Example: 5
}

# func size(list(A)) int / func (list(A)) size() int
mutation SizeList {
  evaluate(expressions: { result: "size([1, 2, 3])" }) # Example: 3
}

# func size(map(A, B)) int / func (map(A, B)) size() int
mutation SizeMap {
  evaluate(expressions: { result: "size({'a': 1, 'b': 2})" }) # Example: 2
}

# func size(string) int / func (string) size() int
mutation SizeString {
  evaluate(expressions: { result: "size('hello')" }) # Example: 5
}
```

### Regex Matching

#### matches

Performs a regular expression match using RE2 syntax. Can be used as a global function or a method on strings.

```graphql
# func matches(string, string) bool
# This is a global function, often used as receiver method.
mutation GlobalMatches {
  evaluate(expressions: { result: "matches('abc', '^a.*c$')" }) # Example: true
}

# func (string) matches(string) bool
mutation StringMatches {
  evaluate(expressions: { result: "'abc'.matches('^a.*c$')" }) # Example: true
}
```

