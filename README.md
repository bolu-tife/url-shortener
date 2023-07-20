## API Endpoints

longUrl - the url to be shortened
shortUrl - shortened url
POST `api/v1/shorten`
    param : {longUrl: string}
    return {shortUrl: string}


GET `api/v1/{shortUrl}`
    Query Parameter - shortUrl
    return: longUrl

GET `api/v1`
    return : returns all 


## Future works
1. Authentication
2. Customization
