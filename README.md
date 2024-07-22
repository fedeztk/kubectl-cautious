# kubectl-cautious 

<!-- # TODO: add exhaustive description -->

Cautiously run kubectl commands, no more accidental deletions!

Supports regexes and is configured via a yaml file under ~/.kube.

This plugin is thought to be used in production environments where you want to be extra cautious when running kubectl commands. It allows pattern matching for contexts and actions, and it can be configured to run in dry-run mode.

## Quick Start

Install using krew:
```
kubectl krew install cautious
kubectl cautious
```

Install directly:
```
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