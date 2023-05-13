export TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6IjFiYjI2MzY4YTNkMWExNDg1YmNhNTJiNGY4M2JkYjQ5YjY0ZWM2MmYiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiVHlsZXIgQXllcnMiLCJwaWN0dXJlIjoiaHR0cHM6Ly9saDMuZ29vZ2xldXNlcmNvbnRlbnQuY29tL2EvQUVkRlRwN0NvTFo1X0ZGczRNVF9tNXl4TkJaSTlGYlRUOEsyZE9jQ1JkUU5aUT1zOTYtYyIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS9jbG91ZDMyeCIsImF1ZCI6ImNsb3VkMzJ4IiwiYXV0aF90aW1lIjoxNjgyMjI3MTAyLCJ1c2VyX2lkIjoiQ1FkWFJmcFB4VWJ3TnJjdkw3Wm9XemtqamFjMiIsInN1YiI6IkNRZFhSZnBQeFVid05yY3ZMN1pvV3pramphYzIiLCJpYXQiOjE2ODM5NTMwMTEsImV4cCI6MTY4Mzk1NjYxMSwiZW1haWwiOiJ0eWxlci5heWVyc0BnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiZmlyZWJhc2UiOnsiaWRlbnRpdGllcyI6eyJnb29nbGUuY29tIjpbIjEwMTYwNDkzMDkzNjg4ODQ3NDQzMSJdLCJlbWFpbCI6WyJ0eWxlci5heWVyc0BnbWFpbC5jb20iXX0sInNpZ25faW5fcHJvdmlkZXIiOiJnb29nbGUuY29tIn19.uVJWytiGxNH52typpFpw0dOlXx8t1EsTVok9_FxFOKMDmuAa-gkcVxocO1eC9BH6mOPbrWJfSEVYlbgn6QpwRaf_x0UtJigPOYr0OU_rtJUWbmEpstENuaLgyksnP1IiuNUPmZ2ObtiItcd1GfPu1qaDg-PpyC1qmPprWNDsELXhyUcNTBjI8xiAbsAYmGLxeih4A2wo1awVCN8f_FhoJDKR8Rr16-TggACeRt2cmu6i3wJ7tTi4E9QNCYwCez_etQXUm1_5S3p0cUOSCmwejBQN07WsxcwDnnvIhpliU4am1wgnaTVjfPURHyTKr4opHExx1nBQWF3W3umKXQeBFA

export HOST=https://cms2-service-ghfontasua-ew.a.run.app/posts
# export HOST=http://localhost:8080/posts

export CONTENT="A snowflake is a single crystal that has achieved a sufficient size, and may have amalgamated with others, which falls through the  as . Each flake nucleates around a dust particle in  air masses by attracting  cloud water droplets, which  and accrete in crystal"

for i in {1..10000}
do
   echo "Loop number $i"

   curl -X POST -i $HOST \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/x-www-form-urlencoded" \
        -d "title=test$i&content=$CONTENT&summary=$CONTENT&authorDisplayName=LoadTest&authorProfilePic=https://lh3.googleusercontent.com/a/AEdFTp7CoLZ5_FFs4MT_m5yxNBZI9FbTT8K2dOcCRdQNZQ=s96-c"
   
   sleep .1
done