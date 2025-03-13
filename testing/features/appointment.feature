Feature: EPIC7-1 Appointment
    As a customer,
    I can create appointments based on sub package and selected time slots,
    so that my schedule is organized.

    Background: Server is running
        Given the server is running
        And a photographer has a package and sub package
        And a customer is logged in

    Scenario: Customer creates the appointment
        When a customer creates an appointment
        Then the appointment is created

    Scenario: Customer cannot create the appointment
        When a customer creates an appointment with wrong format    
        Then the appointment is not created