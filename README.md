# sxpf - S-Expression Framework

This is a framework to work with [s-expressions](https://en.wikipedia.org/wiki/S-expression).
In contrast to some other implementations, there are two types of atoms:
symbols and strings. Symbols are case-insensitive. A string may contain any
sequence of unicode characters, using UTF-8 encoding. An expression of the
form `()`, `(A)`, `(A B)`, ..., is a pair list, where A and B are
s-expressions. If it has the form `(A . B)` it is a pair. An expression of the
form `[]`, `[A]`, `[A B]`, ..., is a vector, where A and B are s-expressions.

The framework contains types, functions, and methods to create s-expressions,
to encode them as a string, and to evaluate them. Evaluation creates a third
atom type, which currently cannot be encoded fully in a s-expression: forms
(aka functions). Forms can be *special*. In this case their arguments are not
evaluated before calling the form.

## Syntax
* `;` starts a comment that lasts until end of line
* String = `"CHAR*"` is a sequence of characters, delimited by quotes
    * CHAR = any unicode character except from category C ("control") or an
       escaped character (ECHAR).
    * ECHAR a character sequence that starts with a backspace `\``
        * `\t` = tabulator (code 9)
        * `\n` = new line (code 10)
        * `\r` = carriage return (code 13)
        * `\xMN` = character with code MN (hex digits)
        * `\uMNOP` = character with code MNOP (hex digits)
        * `\UMNOPQR` = character with code MNOPQR (hex digits)
        * any other character, excapt category C = the character itself, e.g.
          `\\` = backslash, `\"` = quote.
* Symbol = a sequence of characters, except category C and Z ("separator"),
  and except `"`, `(`, `)`, `[`, `]`, `;`, `.`.
* Pair = `(` Z\* (s-expression (Z\* s-sexpression)\* (Z\* `.` Z\* s-expression)?)? Z\* `)`
* Vector = `[` Z\* (s-expression (Z\* s-expression)\*)? Z\* `]`
* Z = any unicode of category Z

## Note

* Cyclic structures are currently not supported. Creating them will likely
  lead to a stack overflow when working with them.

## Usage

* [Zettelstore](https://zettelstore.de) and its [client
  library](https://zettelstore.de/client/) use this framework for transferring
  data and to help encoding some other data formats.
* My students will use it for various lectures / seminars.
