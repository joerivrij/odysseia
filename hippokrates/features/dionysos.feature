Feature: Dionysos
  In order to use grammar functions
  As a greek enthusiast
  We need to be able to validate the functioning of the Dionysos api

  @dionysos
  Scenario Outline: Feminine first declensions result in the correct declension rule
    Given the "<service>" is running
    When the grammar is checked for word "<word>"
    Then the declension "<declension>" should be included in the response
    Examples:
      | service  | declension                | word   |
      | dionysos | noun - sing - fem - nom   | μάχη   |
      | dionysos | noun - sing - fem - gen   | οἰκίας |
      | dionysos | noun - sing - fem - dat   | οἰκίᾳ  |
      | dionysos | noun - sing - fem - acc   | τιμήν  |
      | dionysos | noun - plural - fem - nom | μάχαι  |
      | dionysos | noun - plural - fem - gen | θεῶν   |
      | dionysos | noun - plural - fem - dat | δόξαις |
      | dionysos | noun - plural - fem - acc | χώρᾱς  |

  @dionysos
  Scenario Outline: Masculine first declensions result in the correct declension rule
    Given the "<service>" is running
    When the grammar is checked for word "<word>"
    Then the declension "<declension>" should be included in the response
    Examples:
      | service  | declension                 | word     |
      | dionysos | noun - sing - masc - nom   | πολίτης  |
      | dionysos | noun - sing - masc - gen   | κριτοῦ   |
      | dionysos | noun - sing - masc - dat   | νεανίᾳ   |
      | dionysos | noun - sing - masc - acc   | πολίτην  |
      | dionysos | noun - plural - masc - nom | κριταί   |
      | dionysos | noun - plural - masc - gen | πολίτῶν  |
      | dionysos | noun - plural - masc - dat | νεανίαις |
      | dionysos | noun - plural - masc - acc | κριτᾱ́ς  |

  @dionysos
  Scenario Outline: Masculine second declensions result in the correct declension rule
    Given the "<service>" is running
    When the grammar is checked for word "<word>"
    Then the declension "<declension>" should be included in the response
    Examples:
      | service  | declension                 | word     |
      | dionysos | noun - sing - masc - nom   | δοῦλος   |
      | dionysos | noun - sing - masc - gen   | πόλεμου  |
      | dionysos | noun - sing - masc - dat   | δοῦλῳ    |
      | dionysos | noun - sing - masc - acc   | πόλεμον  |
      | dionysos | noun - plural - masc - nom | θεοί     |
      | dionysos | noun - plural - masc - gen | νεανίῶν  |
      | dionysos | noun - plural - masc - dat | πόλεμοις |
      | dionysos | noun - plural - masc - acc | θεούς    |

  @dionysos
  Scenario Outline: Neuter second declensions result in the correct declension rule
    Given the "<service>" is running
    When the grammar is checked for word "<word>"
    Then the declension "<declension>" should be included in the response
    Examples:
      | service  | declension                 | word   |
      | dionysos | noun - sing - neut - nom   | μῆλον  |
      | dionysos | noun - sing - neut - gen   | δῶρου  |
      | dionysos | noun - sing - neut - dat   | δῶρῳ   |
      | dionysos | noun - sing - neut - acc   | μῆλον  |
      | dionysos | noun - plural - neut - nom | δῶρα   |
      | dionysos | noun - plural - neut - gen | δῶρων  |
      | dionysos | noun - plural - neut - dat | μήλοις |
      | dionysos | noun - plural - neut - acc | μῆλα   |

  @dionysos
  Scenario Outline: Queries with no results return an error
    Given the "<service>" is running
    When the grammar for word "<word>" is queried with an error
    Then an error containing "<message>" is returned
    Examples:
      | service  | word             | message         |
      | dionysos | ναυμαχίαναυμαχία | 200 but got 404 |

  @dionysos
  Scenario Outline: Some words have multiple dictionary entries
    Given the "<service>" is running
    When the grammar is checked for word "<word>"
    Then the number of results should be equal to or exceed "<results>"
    And the number of translations should be equal to er exceed "<translations>"
    And the number of declensions should be equal to or exceed "<declensions>"
    Examples:
      | service  | results | translations | declensions | word    |
      | dionysos | 2       | 2            | 1           | πόλεμου |
      | dionysos | 2       | 2            | 1           | μάχη    |

  @dionysos
  Scenario Outline: Some words have multiple declensions
    Given the "<service>" is running
    When the grammar is checked for word "<word>"
    Then the number of results should be equal to or exceed "<results>"
    And the number of translations should be equal to er exceed "<translations>"
    And the number of declensions should be equal to or exceed "<declensions>"
    Examples:
      | service  | results | translations | declensions | word |
      | dionysos | 2       | 1            | 2           | δῶρα |
      | dionysos | 2       | 2            | 2           | θεῶν |
