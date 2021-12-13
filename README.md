# passwordCritic
a Go library for critiquing your selection of passwords

## the current plan

check the entropy of the password, stopping things like "12341234"

if it passes an entropy check, look it up in a ?markov-chain generated list of combinations of the most common passwords

if it's not anywhere in the list, allow  the user to continue,

else give an exact explanation to the user as to why they shold not choose that password
