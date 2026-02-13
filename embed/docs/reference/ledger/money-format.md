---
title: Money Formatting
metaTitle: "Reference: Money Formatting"
description: Reference for money formatting options for amounts in a Twisp Ledger.
showNext: false
showPrev: true
---

## The Basics

The [Money]($gql:object:Money) type in Twisp represents a monetary amount. Fields of this type can be represented (via its `formatted` field) in various ways depending on the options provided in the [MoneyFormatInput]($gql:input:MoneyFormatInput).

This reference provides a review of the different ways to format monetary amounts in Twisp, including formatting whole numbers, minimum and maximum digits, rounding modes, locales, grouping, currency minimum digits, and currency displays.

{% callout %}
Unless otherwise indicated, all examples below use the `en-US` locale.
{% /callout %}

## Minor Units

Setting the `minDigits` and `maxDigits` determines how many digits of the minor units to display.
When `minDigits` not specified, it will use the default fractional digits for the currency.

| Units       | Currency | Min Digits | Max Digits | Formatted   |
|-------------|----------|------------|------------|-------------|
| 9           | USD      | _DEFAULT_  | _DEFAULT_  | `$9.00`     |
| 9           | USD      | _DEFAULT_  | 0          | `$9`        |
| 9           | USD      | 0          | _DEFAULT_  | `$9`        |
| 9           | USD      | 1          | _DEFAULT_  | `$9.0`      |
| 9           | USD      | 2          | _DEFAULT_  | `$9.00`     |
| 9           | USD      | 3          | _DEFAULT_  | `$9.000`    |
| 9.123456789 | USD      | _DEFAULT_  | _DEFAULT_  | `$9.123457` |
| 9.123456789 | USD      | _DEFAULT_  | 0          | `$9`        |
| 9.123456789 | USD      | _DEFAULT_  | 0          | `$9.1`      |
| 9.123456789 | USD      | _DEFAULT_  | 0          | `$9.12`     |
| 9.123456789 | USD      | _DEFAULT_  | 0          | `$9.123`    |

## Rounding

The `roundingMode` option defines the rounding behavior when the fractional units exceed the `maxDigits`.

| Units | Currency | Rounding Mode | Max Digits | Formatted |
|-------|----------|---------------|------------|-----------|
| 9.005 | USD      | `DOWN`        | 2          | `$9.00`   |
| 9.005 | USD      | `HALF_DOWN`   | 2          | `$9.00`   |
| 9.005 | USD      | `UP`          | 2          | `$9.01`   |
| 9.005 | USD      | `HALF_UP`     | 2          | `$9.01`   |
| 9.006 | USD      | `DOWN`        | 2          | `$9.00`   |
| 9.006 | USD      | `HALF_DOWN`   | 2          | `$9.01`   |
| 9.006 | USD      | `UP`          | 2          | `$9.01`   |
| 9.006 | USD      | `HALF_UP`     | 2          | `$9.01`   |

## Currency Display

Use the `currencyDisplay` to change the currency indicator.

| Units    | Currency | Currency Display | Formatted      |
|----------|----------|------------------|----------------|
| 12345.67 | USD      | _DEFAULT_        | `$12345.67`    |
| 12345.67 | USD      | `CODE`           | `USD 12345.67` |
| 12345.67 | USD      | `NONE`           | `12345.67`     |
| 12345.67 | USD      | `SYMBOL`         | `$12345.67`    |

## Other Locales

Changing the `locale` option formats the amount according to the standards of that locale.

| Units               | Currency | Locale (Country)   | Formatted              |
|---------------------|----------|--------------------|------------------------|
| 123456789.123456789 | USD      | `zh-CN` (China)    | `US$123456789.123457`  |
| 123456789.123456789 | USD      | `es-CO` (Colombia) | `US$ 123456789,123457` |
| 123456789.123456789 | USD      | `fr-FR` (France)   | `123456789,123457 $US` |
| 123456789.123456789 | USD      | `de-DE` (Germany)  | `123456789,123457 $`   |
| 123456789.123456789 | USD      | `hi-IN` (India)    | `$123456789.123457`    |
| 123456789.123456789 | USD      | `ja-JP` (Japan)    | `$123456789.123457`    |
| 123456789.123456789 | USD      | `ar-AE` (UAE)      | `US$ 123456789.123457` |
| 123456789.123456789 | USD      | `en-GB` (UK)       | `US$123456789.123457`  |
| 123456789.123456789 | USD      | `en-US` (USA)      | `$123456789.123457`    |
| 123456789.123456789 | EUR      | `zh-CN` (China)    | `€123456789.123457`    |
| 123456789.123456789 | EUR      | `es-CO` (Colombia) | `€ 123456789,123457`   |
| 123456789.123456789 | EUR      | `fr-FR` (France)   | `123456789,123457 €`   |
| 123456789.123456789 | EUR      | `de-DE` (Germany)  | `123456789,123457 €`   |
| 123456789.123456789 | EUR      | `hi-IN` (India)    | `€123456789.123457`    |
| 123456789.123456789 | EUR      | `ja-JP` (Japan)    | `€123456789.123457`    |
| 123456789.123456789 | EUR      | `ar-AE` (UAE)      | `€ 123456789.123457`   |
| 123456789.123456789 | EUR      | `en-GB` (UK)       | `€123456789.123457`    |
| 123456789.123456789 | EUR      | `en-US` (USA)      | `€123456789.123457`    |

## Grouping

When `groupDigits` is set to `true`, digits will be grouped according to the standards for each locale.

| Units               | Currency | Locale (Country)   | Formatted                |
|---------------------|----------|--------------------|--------------------------|
| 123456789.123456789 | USD      | `zh-CN` (China)    | `US$123,456,789.123457`  |
| 123456789.123456789 | USD      | `es-CO` (Colombia) | `US$ 123.456.789,123457` |
| 123456789.123456789 | USD      | `fr-FR` (France)   | `123 456 789,123457 $US` |
| 123456789.123456789 | USD      | `de-DE` (Germany)  | `123.456.789,123457 $`   |
| 123456789.123456789 | USD      | `hi-IN` (India)    | `$12,34,56,789.123457`   |
| 123456789.123456789 | USD      | `ja-JP` (Japan)    | `$123,456,789.123457`    |
| 123456789.123456789 | USD      | `ar-AE` (UAE)      | `US$ 123,456,789.123457` |
| 123456789.123456789 | USD      | `en-GB` (UK)       | `US$123,456,789.123457`  |
| 123456789.123456789 | USD      | `en-US` (USA)      | `$123,456,789.123457`    |

## Minimum Digits by Currency

The default setting for `minDigits` changes depending on which currency is used because different currencies have a different number of minor units used.

- US dollars (`USD`) use 2 minor units (i.e. cents).
- Jordanian Dinars (`JOD`) use 3 minor units.
- Uganda Shillings (`UGX`) use no minor units.

`UGX` is the code for the Uganda shilling, which uses 0 minor units.

| Units | Currency | Formatted   |
|-------|----------|-------------|
| 9     | USD      | `$9.00`     |
| 9     | JOD      | `JOD 9.000` |
| 9     | UGX      | `UGX 9`     |
