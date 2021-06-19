Feature: herodotos
  In order to test odysseia
  As a developer
  We need to be able to validate the functioning of the Herodotos api

  @herodotos
  Scenario Outline: A user can create a new sentence
    Given the "<service>" is running
    When a new sentence is requested for author "<author>"
    Then the responseCode should be "<response>"
    Examples:
      | service  | response | author    |
      | herodotos | 200     | herodotos |

  @herodotos
  Scenario Outline: A user cannot create sentences with bad input
    Given the "<service>" is running
    When a new sentence is requested for author "<author>"
    Then the responseCode should be "<response>"
    Examples:
      | service  | response | author    |
      | herodotos | 404     | notanautor |
      | herodotos | 400     |     |