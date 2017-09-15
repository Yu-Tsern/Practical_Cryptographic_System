
====================== PROBLEM ONE TO THREE ======================

Please put the input files and main.go in the same directory.
Get into the directory that contains all the files. (main.go, plain.txt, ...)
After inserting "go run main.go" in the command line, the program will start executing. 
Type in one of the commands specified in the document hw1.pdf, for example:

    vigenere-encrypt ABCDEFGHIJKLMN plaintext.txt

    vigenere-decrypt ABCDEFGHIJKLMN ciphertext.txt

    vigenere-keylength ciphertext.txt

    vigenere-cryptanalyze ciphertext.txt 14

The program will end every time it outputs a result.

PS:

(1) For problem number one, my program only accept upper-case letters as an input of keys. 

(2) For problem number two, the program not only output the length of the key, but also recover the key I guess. 

(3) I choose the option: "make a best guess for the decryption key using cipher text only" instead of "write a tool that attempts to complete decrypt a Vigen`ere cipher text using frequency analysis" in problem number three, since the second function can be easily done by feeding my key to the command in problem number 1.

========================== PROBLEM FOUR ==========================

To answer this question, I would like to start with the ways to break a Vigen`ere cipher. As far as I know, there are at least three methods that can achieve this goal. 

(1) Exhaustive key search

(2) See whether there are repeated words in the cipher text. It is very likely that one of the common divisors of distances between them is the key length.

(3) Guess a key length "l", and partition the cipher text "c[N]" into "l" part in the following way. The first part contains (c[0], c[l], c[2l],...). The second part contains (c[1], c[1+l], c[1+2l],...) and so on and so forth. Then calculate every IC( Index of Coincidence) of the partitions, and see whether they are close to 0.066. After several trials we might get the possible key length. 

Maybe there are more ways to gain pieces of information from the cipher text only. However, the benefit of extending the key length can be sufficiently showed in these cases. First, the longer the keys are, the harder the adversary can do an exhaustive key search. Second, for a given plaintext, extending key lengths lowers the expectation of seeing repeated patterns in the cipher text. Third, when keys become longer, the size of each partition become smaller, which might lead to difference between the calculated IC and the theoretical value, since there are fewer samples. In sum, it is more secure to use a longer key in Vigen`ere encipherment. Therefore, a key that is at least as long as the message is definitely advantageous. However, just like the case in one time pad, the keys must be as random as possible as we can to prevent the adversary from predicting them. Besides, the fewer the repeated patterns in the key, the less likely the adversary can carry out the second method above successfully. Last but not least, the key should never be used twice! Otherwise it would face the same problem of using a key of length N/2 to encrypt a length N plaintext. 

