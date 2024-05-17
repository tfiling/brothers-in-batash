## Test Implementation:

Implement tests for <> . Make sure the tests fully covers the code.
Make sure that each test tests a single use case.
Make sure you add "Arrange" "Act" "Assert" comments.
Make sure the test uses the mock type <> as a dependency injection for <>
The tests naming convention should be: Test<tested type>_<tested method>__<tested use case in snake case>
The project's source code:

## Implement CRUD API controller

# code 1:

I would like to implement a new controller type that would expose CRUD API for <> .
Generate a code block of const strings that would be used as URIs for the crude API.
The project's source code:

# code 2:

Implement a new controller type named <> that would expose CRUD API for <> .
The new controller should be using fiber package.
The new controller should implement the controllers.Controller interface.
The new controller should accept the following types ( marked with * ) to it's constructor function and store them as
properties:

* <>
  When registering it's API routes it should use the above string consts.
  A relevant example would be controllers.ShiftController

# tests 1:

Implement a mock type for interface <>

# tests 2:

Implement tests for controllers.SoldierController . Make sure the tests fully covers the code.
Make sure the tests make use of the following mocks ( marked with * ) as dependencies for
controllers.SoldierController :

* MockISoldierStore as a mock for ISoldierStore
  Make sure that each test tests a single use case.
  Make sure you add "Arrange" "Act" "Assert" comments.
  Make sure you assert expectations for the mocks being used.
  The tests naming convention should be: Test<tested type>_<tested method>__<tested use case in snake case>



