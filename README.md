# Vactl 

A simple ci/cd tool for Hashicorp Vault.

Not fully Functional yet..


#v0.1

Right now have simple and dumb functionallity.

## Configure
You configure your vactl with a .vactl.yaml file containing the following:
```
---
vaultToken: "<your vault token>"
vaultAddr: "<vault address"
```

When you run vactl it will look for this file in the same dir your running your binary.

## Functionality implemented

You can do the following verbs,
- get
- apply
- delete

### get

The get command can be used to get information in a format that vactl likes.
For example you can output users as list in a given path like so,

> vactl get users -p userpass

if you want it as yaml simply add -o flag
> vactl get users -o

you can the following resources
- users (as yml)
- ssh-roles (as yml)
- policies (as hcl)

TODO:
- should be possible to output policies as a collection of files

### apply

You can apply the following Vault resources,
- users (as yml)
- ssh-roles (as yml)
- policies (as .hcl)

apply can target directories or single files.

If they have .hcl suffix vactl will try to upload them as policies.

Example users yml should look like this,
```
---
type: users
users:
  - name: foo
    method: userpass  # Method is also path if you have multiple methods
	token_policies:
       - foo-policy
  - name: bar
    method: ldap
	token_policies:
	  - bar-policy
```

### delete

You can delete a resource with vactl by adding simple typing,
> vactl delete user sebbe
