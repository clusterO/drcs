## DCRS

Track software revisions and allow different networks to share work on a project.
dcrs use peer to peer approach, each peer has a working copy of the codebase, copies
are synchronized by patches exchange between peers.

I'll try here to implement the basic features of a dcrs, commit, push, pull, view history,
merging and reverting changes using golang.

The project will have to main subsystems, one for networking, that implement the communication
between the peers, and the crs.

#### refs

- http://en.wikipedia.org/wiki/Distributed_revision_control
- http://git-scm.com
- http://wiki.bazaar.canonical.com
- https://tom.preston-werner.com/2009/05/19/the-git-parable.html
