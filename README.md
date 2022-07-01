[![Coverage Status](https://coveralls.io/repos/github/jdginn/durins-door/badge.svg?branch=main)](https://coveralls.io/github/jdginn/durins-door?branch=main)

# Durins-door is a Rosetta stone for programming languages using DWARF
Durins-door is a Go package that uses [DWARF debug information](https://dwarfstd.org/) to extract
information about compiled programs and represent that information in
a predictable way.

Some of the information available:
- Typedefs, sizes, members, etc.
- Types of variables by name
- Value of static variables in binary files or memory
- Decoding values of struct members from byte representations of structs

Durins-door is first and foremost a Go package for use in Go programs.
However, Durins-door also includes servers providing an HTTP API (coming soon)
and a gRPC API (coming soon). We have also provided some clients for
Durins-door in a few languages (coming soon):
- [dwarf-explore](https://github.com/jdginn/dwarf-explore): a TUI for exploring the contents of programs, written in golang using durins-door as an imported package
- [python client](): `coming soon`
- [go client](): `coming soon?`
- TODO: maybe others? TypeScript? C++? Maybe not?

Durins-door is supported on macOs and Ubuntu Linux.

## Why do we need Durins-door?
Programs often need to interact. In some situations, programs require some
understanding of the internals of other programs to facilitate this interaction.

Some examples we have encountered in our own careers:
- Communicating with a microcontroller using a protocol based on structs
  written in C++
- Debugging one program from another program
- Reading values out of a compiled binary
- Modifying static memory scratchpads to communicate with a program running
  on a microcontroller

This task is more difficult when the two programs are written in different
languages or the source of the program with which we wish to interact is
unavailable to us for any reason.

The typical approach is to manually encode the relevant protocols, struct
types, or addresses from one language into another. This is fragile due to:

1. Incompatible versioning:
What happens if you forget to update one of your programs to the new version?
What happens if one program needs to interact with multiple versions of the other?)
2. Human error:
What if you implement the relevant structure incorrectly in one of the programs?
3. Compiler decisions:
What if the compiler resolves one of these items of interest in an unpredictable
way (at least unpredictable before compile-time)?
4. Inefficiency:
It can be time-consuming to manually implement protocols or addresses over
and over again.

Durins-door presents a unified approach to extracting this information from the
exact binary of interest. This removes the versioning concern because we view
exactly what the binary contains. This also removes the compiler decision conern
for the same reason.

Because the API to extract the information provided by durins-door is always the same,
once a client is written, that client is now equipped to understand and work with any
of the information provided by durins-door. This avoids the concern about duplicated
effort and human error.

## What is the status of durins-door?
Durins-door is still experimental and not yet feature-complete. Currently we are only
testing it against C++ binaries and additional language support will is not yet in
the roadmap (although in general, DWARF should be standard enough for most C-like
languages to be easily supported).

## How does durins-door work?
The [DWARF standard](https://dwarfstd.org/) provides an abstract, language-independent
representation of compiled code. This is the basis for most debuggers such as gdb
or llvm-db and as such, is already heavily relied upon and proven. Provided a binary
is compiled with debug information, most aspects of the binary can be introspected
using the accompanying DWARF.

An important note is that DWARF debug information _does_ increase binary size,
sometimes dramatically. However, most compilers can locate the debug information
it is own file that can be referenced entirely independently of the actual
executable.

## Is durins-door fast?
As it is still a work in progress, it's too early to evaluate performance (we have
not even begun optimization). However in the typical use-case, durins-door has
natural leg up. We envision the typical usecase as working with low-level constructs
of a compiled language from within an interpreted language (for example, teaching
Python to understand C++). In these situations, using a compiled, relatively
efficient language like Go will tend to outperform a native implementation in the
interpreted language. The debug/dwarf package in the Go standard library is quite
good and offers efficient ways to work with DWARF.

In addition, durins-door is generally intended to interact with other programs
as a server using HTTP or gRPC APIs. This allows durins-door to run asynchronously
to a program in a language like Python, for example. Residing in its own process,
there are also opportunities for durins-door to persist long beyond the runtime
of its clients, making it easier to cache the results of queries and avoid
paying the processing penalty each time the program boots. This is especially useful
for shorter-lived programs that would typically be forced to parse the relevant
DWARF each time they run. Locating all of this responsibility in durins-door allows
us to invest in efficient caching on the server side and free users from having
to roll their own implementation on the client side.
