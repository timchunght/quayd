# Quayd

Quayd extends the Quay.io <-> GitHub integration. It:

1. Excepts WebHooks from Quay.io and creates appropriate [GitHub Commit Statuses](https://developer.github.com/v3/repos/statuses/).
   ![](https://s3.amazonaws.com/ejholmes.github.com/A72Nj.png)
2. Tags the Quay.io build with the git sha.

**TODO**

* GitHub backend implementation of StatusesRepository.
* GitHub backend implementation of CommitResolver.
* Command.
