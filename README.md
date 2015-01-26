# Quayd

Quayd extends the Quay.io <-> GitHub integration. It:

1. Excepts WebHooks from Quay.io and creates appropriate [GitHub Commit Statuses](https://developer.github.com/v3/repos/statuses/).
   ![](https://s3.amazonaws.com/ejholmes.github.com/A72Nj.png)
2. Tags the Quay.io build with the git sha.

## Usage

Run quayd to start the server, providing it with a github api token that has the **repo** scope.

```console
$ quayd -port=8080 -github-token=1234
```

Now, create some webhooks on Quay.io that POST to "/quayd/\<status\>"

![](https://s3.amazonaws.com/ejholmes.github.com/0mIUw.png)
