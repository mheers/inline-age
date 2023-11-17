#!/bin/bash

set -e

vault secrets enable -path=kv kv
vault kv put kv/foo bar=baz

echo done > /tmp/vault-init-done

# vault kv get kv/foo
