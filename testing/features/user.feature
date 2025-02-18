Feature: Login System US2-1
    As a new user,
    I can create an account using my email, password, or social media accounts
    so that I can access the platform.

    Background: Server is running
        Given the server is running

    Scenario: User logs in
        Given valid credentials are provided
        When the login is submitted
        Then access to the account is granted
        
    Scenario: User invalid credentials
        Given invalid credentials are provided
        When the user attempts to log in
        Then the system should reject the login and display an error message saying "Invalid email or password"

