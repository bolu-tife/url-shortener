## API Endpoints

longUrl - the url to be shortened
shortUrl - shortened url
POST `api/v1/shorten`
    param : {longUrl: string}
    return {shortUrl: string}


GET `api/v1/{shortUrl}`
    Query PArameter - shortUrl
    return: longUrl

GET `api/v1`
    return : Welcome message




## Todo: 
1. add postgres service - for data storage
2. add redis service for caching
3. add a rate limiter
4. add tests
5. input validations
