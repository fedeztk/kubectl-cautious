# kubectl-cautious 

<!-- # TODO: add exhaustive description -->

Cautiously run kubectl commands, no more accidental deletions!

Supports regexes and is configured via a yaml file under ~/.kube.

This plugin is thought to be used in production environments where you want to be extra cautious when running kubectl commands. It allows pattern matching for contexts and actions, and it can be configured to run in dry-run mode.

Adding an alias to your shell may be a good idea in order to reduce cognitive load when running kubectl commands.

```sh
alias kubectl='kubectl cautious'
# or
alias k='kubectl cautious'
```


https://github.com/user-attachments/assets/0934e030-bdac-49b9-8d8c-69292d4f416e




## Quick Start

Install using krew:
```sh
kubectl krew install cautious
kubectl cautious
```

Install directly:
```sh
# choose a proper release
wget https://github.com/fedeztk/kubectl-cautious/releases/download/LATEST-TAG/cautious_YOUR_PLATFORM.tar.gz -O cautious.tar.gz
tar -xvf cautious.tar.gz
mv cautious /usr/local/bin/kubectl-cautious
kubectl cautious
```

## Configuration

The configuration file is located at `~/.kube/cautious.yaml` and has the following structure:
```yaml
# cautious.yaml
contexts:
    - name: prod # context name, supports regex
      actions:   # list of actions to be cautious about
        - name: apply
          dry-run: false # whether to run the command in dry-run mode
        - name: delete
          dry-run: true
```
