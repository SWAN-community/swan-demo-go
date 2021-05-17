#!/bin/bash
if [[ $(comm -12 <(sort -u /etc/hosts) <(sort -u hosts-sample)) ]]; then 
    echo "No hosts written as your system hosts file already contains some SWAN entries."
else 
    cat hosts-sample >> /etc/hosts
fi

