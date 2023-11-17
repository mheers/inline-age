#!/bin/bash

set -e

# Set the path to the id_rsa.pub file
id_rsa_pub_path="/ssh/id_rsa.pub"

# Read the contents of the id_rsa.pub file
pub_key=$(cat "$id_rsa_pub_path")

# Extract the name from the public key
name=$(echo "$pub_key" | cut -d " " -f 3-)

# Calculate the fingerprint from the public key
ssh-keygen -lf "$id_rsa_pub_path" | awk '{print $2}' | sed 's/://g' | tr '[:upper:]' '[:lower:]' > fingerprint.txt
fingerprint=$(cat fingerprint.txt)

# Clean up the temporary fingerprint file
rm fingerprint.txt

# Set the missing values for the SQL INSERT statement
id=1
owner_id=1
mode=2
type=1

# Construct the SQL INSERT statement
insert_stmt="INSERT INTO public_key (id, owner_id, name, fingerprint, content, mode, type) VALUES ($id, $owner_id, \"$name\", \"$fingerprint\", \"$pub_key\", $mode, $type);"

# Execute the SQL INSERT statement
sqlite3 /data/gitea/gitea.db "$insert_stmt"

# check success
echo 'select * from public_key' | sqlite3 /data/gitea/gitea.db

# create gitea admin user
sudo -u git gitea admin user create --username gitea_admin --password admin --email gitea@example.com --admin --config /data/gitea/conf/app.ini

# create gitea admin user
sudo -u git gitea admin regenerate keys

# set git config
git config --global user.email "you@example.com"
git config --global user.name "Your Name"

# login as gitea admin
token=$(sudo -u git gitea admin user generate-access-token --username gitea_admin --scopes repo | awk '{print $6}')
tea login add --url http://localhost:3000 --name localtoken --token $token

# create a gitea demo repo
tea repo create --name demo --description "Demo repo" --private

# clone the gitea demo repo
git clone http://gitea_admin:admin@localhost:3000/gitea_admin/demo.git

# add files to the gitea demo repo
cd demo && cp -r /import/* . && git add . && git commit -m "add demo files" && git push

exit 0
