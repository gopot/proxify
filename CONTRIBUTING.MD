# Introduction

First off, thank you for considering contributing to this project. There are many ways to contribute, from writing tutorials or blog posts, improving the documentation, submitting bug reports and feature requests or writing code. Feel welcome and read the following sections in order to know what and how to work on something.

Following these guidelines helps to communicate that you respect the time of the developers managing and developing this open source project. In return, they should reciprocate that respect in addressing your issue, assessing changes, and helping you finalize your pull requests.

# Ground Rules

1. Whatever you do for this project, it should aim to make it better. There is no need to make it worse - it could end up there by its own.
1. You should better use GitHub features available for this repository rather than third party tools, unles it is reasonable. GitHub offers very reach functionality for full project cycle.
1. You are expected to follow [Code of Conduct](./CODE_OF_CONDUCT.MD) as well as other participants. Communication plays critical role on this project and it is important for everyone (whether a user or Contributor) to keep discussions constructive, clear and polite.
1. The project language is English.

# Guidelines

## Go Code:
1. The code is fmt'ed before pushing.
1. CamelCase, explicit and meaningfull naming.
1. All hardcoded values are declared at `const( ... )` block on top of the file or at dedicated file.
1. Every entity has its own file with corresponding name.

### Go Code Comments:
1. All var/const blocks and public members are well commented.
1. Every package/folder should contain doc.go file fully and clearly describing it purpose and usage.
1. It should be explicitly mention(unless its clear from name) whether it is safe to access member concurrently.
1. Any known issue on public member should be explicitly mentioned in comment to this member as following : `// TODO(me): #5 Short issue description`


### Go Code Tesing:
1. Test(s) are better table driven.
1. It is fine to have tests for private members, but not mixed with concurrent and/or Public members tests. It would be wonderfull to have each type of tests in its own file.
1. If some member should be safe for concurrent access, it should be tested. `-race` is acceptable, but should be mentioned in documentation.

Copyright © 2017 Ivan Kostko (github.com/ivan-kostko; github.com/gopot)
