The user srevice

[] admin firebase custom claims

[x] create user profile

[x] Downlaod isomia, its last stuff in ur terminal just press up key 

[x] write test for profile service 

[] write make admin func in gatewat

PATH="${PATH}:${HOME}/go/bin"

to export it permanently export PATH="${PATH}:${HOME}/go/bin"

[] to figure out how to know if user has already signed in before from the front end, anytime user uses oauth to sign in use the get profile route with their token to know if the prodile already exists, if itr does skip sending post request to create new profile

[] Update migration to have a seperarte make command to run it and not run directly from main.go

[] Fix export path to be permanent so you dont always have to deal with that

[] Add context timeout in function calls and db calls. Use bff guy repo for refernce 

[x] use stream for get profiles and get scores in quiz service, for anything that returns multiple features to client

[x] write code for get profiles and its test, remember to make it a stream grpc

[x] write code test for profile_service

[] Add redis as background worker and for caching data



##Common error 

The most common error you will likely encounter at this point is something like:

sql: expected 2 destination arguments in Scan, not 1

This means that you are passing the incorrect number of arguments into the Scan() method. You need to pass in enough arguments for all of the columns retrieved by your SQL statement. If you are unsure of how many this is, I suggest calling the rows.Columns method which will return a slice of column names that are being returned by your SQL statement