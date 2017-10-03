#!/bin/bash

key=$(xxd -l 32 -p /dev/urandom | tr -d " \n")

$1 "encrypt" -k $key -i $2 -o $3
ct=$(cat $3)
echo "The Script thinks that the ciphertext for key " $key " is: "
echo $ct


$1 "decrypt" -k $key -i $3 -o $4
pt=$(cat $4)
echo "The Script thinks that the recovered plaintext for key " $key " is: "
echo $pt




