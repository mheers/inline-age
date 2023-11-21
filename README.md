# Inline-age

Knows JSON files and can encrypt/decrypt paths in them.

It works with SSH keys and leverages [age](https://github.com/FiloSottile/age) for security.

## Installation

```bash
go install github.com/mheers/inline-age@latest

# create a shortcut symlink for convenience
ln -s $HOME/go/bin/inline-age $HOME/go/bin/ia
```

## Inspiration

This project is inspired by a lot of existing technologies.

![inline-age-meme](docs/inline-age-meme.png)

# TODO

- [x] referencing
- [ ] tests
- [ ] cleanup
- [ ] implement YAML
- [ ] implement CUElang
- [ ] make ./ia init-file error when file is already initialized to not break secrets
- [ ] store vault token and git ssh key as secrets
- [ ] strip age header from encrypted strings
