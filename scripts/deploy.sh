#!/usr/bin/env bash

#deploy Web
chmod a+x /home/travis/build//Undermove/RockScissorPaperWebInterface/src/main
sshpass -p $FTP_PASSWORD scp -o StrictHostKeyChecking=no -r /home/travis/build//Undermove/RockScissorPaperWebInterface/src/* $FTP_USER@$FTP_HOST:$FTP_DIR
sshpass -p $FTP_PASSWORD ssh -o StrictHostKeyChecking=no $FTP_USER@$FTP_HOST 'sudo systemctl restart rsponline'