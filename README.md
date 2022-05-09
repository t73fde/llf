# sxpf - S-Expression Framework

This is a small framework to work with simple
[s-expressions](https://en.wikipedia.org/wiki/S-expression). In contrast to
some other implementations, there are two types of atoms: symbols and string.
Symbols are case-insensitive, a string may contains any sequence of unicode
characters, using UTF-8 encoding. An expression of the form `()`, `(A)`, `(A
B)`, ..., is a list, where A and B are s-expressions. There is no pair type.

The framework contains types, functions, and methods to create s-expressions,
to encode them as a string, to write them somewhere, and to evaluate them.
Evaluation creates a third atom type, which currently cannot be encoded fully
in a s-expression: forms (aka functions). Forms can be special, their
arguments are not evaluated before calling the form.

## Note

* Cyclic structures are currently not supported. Creating them will likely
  lead to a stack obverflow.

## Usage

* [Zettelstore](https://zettelstore.de) and its [client
  library](https://zettelstore.de/client/) use this framework for transferring
  data and to help encoding some other data formats.
* My students will use it for various lectures / seminars.
