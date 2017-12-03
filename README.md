# 443-Assignment-2

The file "test.go" tries to generate a password with a given pattern. 
However, I still can not select a random word from the given dictionary. 

Update on Assignment2:
I checked the code and to my knowledge it seems like it works

Update on Assignment 3:
Got create to work. Check for any bugs
Also for my path variable you will have to change it to where you have
the file.

12/3/2017 10 am
1. I moved some code around. For example, I check if the file exists right
after the user enters the file name. And I check if the two master passwords
match right after the user enters the password for a second time. 

2. I found it not so easy to type in a 32 byte password every time, so 
I accept the master password if it is not longer than 32 bytes. 

3. When the user enters an invalid master password or if the two passwords 
do not match, I change "return &wallet" to "os.Exit(0)". Because previously,
I found that even if the user enters an invalid password, a new file will
still be created (it is wierd, ummmm). 
