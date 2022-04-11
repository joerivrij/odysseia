Feature: Sokrates
  In order to work with multiple choice quizes
  As a greek enthusiast
  We need to be able to validate the functioning of the Sokrates api

  @sokrates
  Scenario Outline: Querying methods should return a list of methods
    Given the "<service>" is running
    When a query is made for all methods
    Then the method "<method>" should be included
    And the number of methods should exceed "<results>"
    Examples:
      | service  | method       | results |
      | sokrates | herodotos    | 4       |
      | sokrates | aristophanes | 4       |
      | sokrates | plato        | 4       |
      | sokrates | mouseion     | 4       |
      | sokrates | logos        | 4       |

  @sokrates
  Scenario Outline: Querying categories within a method should return categories
    Given the "<service>" is running
    When a query is made for all methods
    And a random method is queried for categories
    Then a category should be returned
    Examples:
      | service  |
      | sokrates |

  @sokrates
  Scenario Outline: Querying for a last chapter should return a last chapter
    Given the "<service>" is running
    When a query is made for all methods
    And a random method is queried for categories
    And a random category is queried for the last chapter
    Then that chapter should be a number above 0
    Examples:
      | service  |
      | sokrates |

  @sokrates
  Scenario Outline: The flow to create and answer a question should return a right or wrong answer
    Given the "<service>" is running
    When a new quiz question is requested
    And that question is answered with a "<answer>" answer
    Then the result should be "<answer>"
    Examples:
      | service  | answer |
      | sokrates | true   |
      | sokrates | false  |
