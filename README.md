# Distributed Revision Control System (DRCS) & Dynamic Deployment System

This project aims to implement the DRCS and DDS systems in Go programming language. The primary goal is to develop a stable implementation of these systems and subsequently benchmark them against similar systems implemented in other programming languages.

## DRCS: Distributed Revision Control System

The DRCS is a system designed to track software revisions and enable different networks to collaborate on a project. It follows a peer-to-peer approach where each peer maintains a working copy of the codebase. The copies are synchronized through the exchange of patches between peers.

In this project, we will focus on implementing the fundamental features of a DRCS, including commit, push, pull, viewing history, merging, and reverting changes using the Go programming language.

### DDS : Dynamic Deployment System

A tool for automated deployment and build.

#### Refs

- [Distributed version control](http://en.wikipedia.org/wiki/Distributed_revision_control)
- [Git](http://git-scm.com)
- [The Git Parable](https://tom.preston-werner.com/2009/05/19/the-git-parable.html)
- [Software deployment](https://www.wikiwand.com/en/Software_deployment)

#### External libraries

- [equalfile](https://github.com/udhos/equalfile)
- [three way merge](https://github.com/charlesvdv/go-three-way-merge)

#### TODO

- [x] Functions test & refactor
- [_] Add only modified files
- [_] Store SHA1 as tree
- [_] Add manual conflicts resolution
- [x] Implement network connect
- [_] Store files as blobs
- [_] Implement branches & branches mergin
- [x] Unit tests

#### Usage

```shell
go run . -help
go run . -init y
go run . -add <file | directory>
go run . -commit "commit message"
go run . -clone localhost:8181/<package name>
```

## Installation

Instructions for installing and setting up the project will be provided in the project's documentation.

## License

This project is licensed under the [MIT License](LICENSE).

Please refer to the [LICENSE](LICENSE) file for more details.

## Contribution

Contributions to this project are welcome. If you encounter any issues or have suggestions for improvements, please feel free to open an issue or submit a pull request.
