# url-shortener

This is a simple URL shortener service that allows you to shorten long URLs into shorter, more manageable links. It's built using  react as frontend, golang's fiber as the backend server and redis as the database.

You must have commonly observed these shortened links when you share links accross social apps or posts on LinkedIn.

Features
1. Shorten long URLs into concise, easy-to-share links.
2. Redirect users to the original URL when they access the shortened link.
3. Track the number of clicks on each shortened link.
4. Prevents DDOS and bot attacks by allocating a QUOTA for each user.

Setup

Since the project is setup using docker compose. Once the keys are added, the project should ideally work with a simple 

    ```
        docker compose up --build
    ```

To Connect to redis database

    ```
    docker compose exec -it database redis-cli
    ```
    
    alternate between different dbs using `select`.

Quick Demo:-

https://github.com/rohitchauraisa1997/url-shortener/assets/67869038/8ae2c12f-cbaf-495f-baa6-ed97e380ae21
