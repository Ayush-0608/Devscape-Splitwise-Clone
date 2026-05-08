Splitwise Clone is the backend of an expense sharing platform that can hadle friends, groups, expenses and splitting of the expenses .

Features:

User Management: register, login using jwt auth, and profile(incredibly basic) management

Friends: Searching for users, send friend requests, and accept/reject them

Groups: Create groups, add friends to groups, check group info

Expense Logging: Named expenses with description, dates and payer info, deleting expneses and amount (checks>0)(amounts are stored as int where last two digits are decimal so divide by 100 to get actual amount and inverse while adding smount)

Splitting: equal split: Amount/number of people in group, 
           amount based split: Checks whether amounts are valid(>0 and sum not greater than expense amount) and splits remainder with other people equally, 
           percentage split: Checks whether all percentages are >0 and sum is leass than 100 percent and splits remainder with other people equally
           note: Every split Can be marked as paid or not
           
Balances: Two levels person to person: Gives who you owe and who owes you and details about what expense you owe for
                     group level: Same as person to person but only checks people in a group
                     note: Ignore splits that have alrwady been paid

Setup: Create an empty database and modify connection string in db.go to contain username, password, database name.
Table creation is handled by progarm.

Video drive link: https://drive.google.com/drive/folders/1yN1wyKo8me5pQ37uiXtBgwN5L2TD6rwm?usp=sharing
