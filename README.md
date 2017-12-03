# 443-Assignment-2

The file "test.go" tries to generate a password with a given pattern. 
However, I still can not select a random word from the given dictionary. 

Update on Assignment2:
I checked the code and to my knowledge it seems like it works

Update on Assignment 3:
Got create to work. Check for any bugs
Also for my path variable you will have to change it to where you have
the file.

12/2/2017 23:41
1. I found it not so easy to type in a 32 byte password every time, so 
I accept the master password if it is not longer than 32 bytes. 

2. I check if the two passwords match right after the user re-enters the 
master password instead of checking it in the if statement.

3. When the user enters an invalid master password or if the two passwords 
do not match, I change "return &wallet" to "os.Exit(0)". Because previously,
I found that even if the user enters an invalid password, a new file will
still be created (it is wierd, ummmm). 

4. input[32] is out of index, and it seems like doing this only compares 
one byte of the passwords. 
