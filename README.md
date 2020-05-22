## Mountup

Hi, and welcome to the Mountup CLI repo. This project is in **alpha** and
very much explorative.

Mountup is a command-line utility that syncs code on remote
machines to your local filesystem. You edit locally, with all your
favorite tooling, and see your changes synced remotely.

### Installation

Currently **only OSX and Linux** are supported.

```bash
brew install mountup-io/mountup/mountup
```

### How does it work?
`mountup sync` will sync a remote directory to a local `~/mountup/servername`
folder. Any local file changes are synced right away to the remote ones.

#### Bring your own server

Existing remote code -> Local

```bash
mountup sync pull username@remote_host:directory_on_remote <ssh_key_path>
```

Local -> Existing remote code

```bash
mountup sync push username@remote_host:directory_on_remote <ssh_key_path>
```
---

#### No server? No problem!

First authenticate!
```bash
mountup signup
mountup login
```

Provision a new server
>Currently an [AWS EC2 Instance](https://aws.amazon.com/ec2/)

```bash
create <servername>
```

Off to the races!
```bash
mountup sync <push/pull> <servername>:<directory_on_remote>
```
