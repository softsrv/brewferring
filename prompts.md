# 3/29/2025 The first prompt

Hello, please channel your finest devloper expertise to craft a web application according to the below specifications.

# Summary

"Your Coffee is Brewferring" is a simple web app designed to be used as an interface to the www.terminal.shop API for viewing and ordering coffee.

# Technology

This website uses the following technology:

- latest version of golang as a back end
- vanilla html/javascript for the front end
- add tailwindcss and daisyUI with the "synthwave" theme for styling
- include HTMX for making dynamic AJAX style calls

# ordered tasks

please complete the following tasks in order.

1. set up a basic website back end. create a home page "/" and a dashboard page "/". Golang should handle the routing, and each route should return the rendered HTML. Use the "templ" library for creating html templates for each of these pages. The home page should have some nicely designed filler text about coffee.
2. set up authentication via www.terminal.shop oauth. All pages on the site should include a navbar at the top. On the right side of the navbar, there should be a "login" button. When the user clicks "login", it redirects to the termina.shop oauth workflow. The access_token from the oauth response should be stored as a user session in the back end. It should be returned as an http only secure cookie in the front end. Use the following data to complete this workflow:

- `{"issuer":"https://auth.terminal.shop","authorization_endpoint":"https://auth.terminal.shop/authorize","token_endpoint":"https://auth.terminal.shop/token","jwks_uri":"https://auth.terminal.shop/.well-known/jwks.json","response_types_supported":["code","token"]}`

3. If the access_token is found in the cookies, the "login" button should instead change to a "hamburger" icon that when clicked, expands with 2 options. The first option says "profile", and has no action on click. The second option says "logout" and will clear the access_token cookie and take the user back to the home page.
4. When a user successfully logs in, they should be taken to the "/products" page. The back end for this page should use the terminal.shop client SDK found here https://github.com/terminaldotshop/terminal-sdk-go combined with the access token in the user session to fetch the list of products. This list of products should be displayed on the products page in a nicely formatted table using "templ" to inject the data into the html page.
5. Create a "/profile" page. The back end for this page should use the terminal-sdk-go package to get the profile data for the current user and display it on the page. The data should be injected via "templ" template.
6. create a "/orders" page. Thge back end for this page should use the terminal-sdk-go package to get the list of orders for the current user and display it on the page in a nicely formatted table. Update the navbar to show "products" and "orders" as clickable options. Each one taking you to the respective page. "orders" should only show if the user is logged in, but "products" can show even if the user is not logged in.

# Note: I had to make some manual changes and some additional small prompts to get the oauth going. The AI could not figure out how to use the terminal-sdk-go until it offered to vendor the dependency and read the code. The AI

# 3/31/2025 - adding db and some objects.

Please implement the following design spec to the best of your ability. The spec is written in markdown to help you with the parsing.

# New Feature: Devices

## New Feature Summary

The Brewferring website allows customers to order coffee from terminal.shop. We'd like to add a new exciting way for customers to accomplish this simple task, and that is via "devices". A device can be anything that is capable of making an http call to report some data. Once a customer has a device successfully reporting data, they can create a "scheduled order", and they can specify a threshold based on the data. For example, if my "device" reports "15%", then place an order. This device could report a number that represents the weight of the coffee, the visual fullness of the container, or any other fun creative things a customer might come up with.

## Design

### database

- add a database layer. Use sqlite. all database tables should have "created_at" and "updated_at" columns that are automatically filled, as well as an auto-increment unique ID as the primary key.
- Create a "users" table with columns that are based on the data returned by the terminal-sdk "profile".
- Create a "devices" table. Devices should have a name and should belong to a user. Each user can have many devices.
- create a "device*tokens" table. A "device" can have one or more "tokens". Each token is a cryptographically securely generated random string of length 32 with a prefix of "dt*". Device tokens have the standard set of columns plus an additional date column called "last_used_at" which is an optional column.
- Create a "schedulers" table. A scheduler can have an associated devices, but it is optional. A scheduler can be configured to use a "date" or can instead use some threshold based on the associated device. A scheduler can have at most 1 device associated with it. It cannot have a device and a date at the same time.
- create a "device_data" table. One device can have many device_data table entries associated with it. The device_data table has only one field (in addition to the standard fields listed above) called "value" which is a float64 value. This table should also have a "user_id" field that links the device_data back to the original user who owns the device.
- create an ORM with golang types and CRUD methods for all of the above objects

### New Web Functionality

- when a user logs in via terminal.shop oauth, the callback route should create a new entry in the "users" table
- create a "schedulers" page "/schedulers". Include a "schedulers" tab in the navbar if the user is logged in. This page shoud list the user's existing schedulers in a table. Each row contains the details of the scheduler plus a delete button. This page should also have a "new" button that brings up a modal form for creating a new scheduler. Create the necessary back end route handlers for this functionality to work via the ORM created previously.
- create a "devices" page "/devices". Include a "devices" tab in the navbar if the user is logged in. This page should list the user's existing devices in a table. Each row contains the details of the device plus a delete button. Create the necessary back end route handlers for this functionality to work via the ORM created previously.

### Devices API

- Create an API layer in the application. Each API uri should start with "/api".
- Any route that starts with "/api" will be authenticated via device_token instead of the usual terminal.shop oauth token. The device_token should be in the "Authorization" header as a Bearer token. When authenticating, check that the token is valid (exists in the db) and then fetch the device from the database and add it to the request context
- create a POST "/api/data" route which takes in a POST body JSON payload that looks like the following

```
{
  "value": 5
}
```

Where the value should always be a JSON number. When parsing this data in the "/api/data" route handler, the number should be cast to a float64.

- update the sqlite ORM to be able to handle CRUD operations for the device_data table. The new "/api/data" POST route should create a new entry in this table via the ORM. Make sure the entry includes the unique ID of the device as well as the unique user id as foreign keys.
- the "/api/data" route should be rate limited per unique device. Each device key can make at most 1 call per hour. When a key is used to make a call, if the call is successful, update the last_used_at column to the appropriate date/time. Any subsequent calls will be rejected with the appropriate http status code until the time difference between the request time and the last_used_at colum exceeds 1 hour.
