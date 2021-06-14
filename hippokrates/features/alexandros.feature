Feature: Alexandros
  In order to test Lexiko
  As a developer
  We need to be able to validate the functioning of the Alexandros api

  @alexandros
  Scenario Outline: A user gets a response when searching for a word
    Given the "<service>" is running
    When the word "<word>" is queried
    Then the responseCode should be "<response>"
    Examples:
      | service  | response |  word |
      | alexandros | 200      |  ἀγαθός  |
