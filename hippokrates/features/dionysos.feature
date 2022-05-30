Feature: Dionysios
  In order to use grammar functions
  As a greek enthusiast
  We need to be able to validate the functioning of the Dionysios api

  @dionysios
  Scenario Outline: Feminine first declensions result in the correct declension rule
    Given the "<service>" is running
    When the grammar is checked for word "<word>"
    Then the declension "<declension>" should be included in the response
    Examples:
      | service  | declension                | word   |
      | dionysios | noun - sing - fem - nom   | μάχη   |
      | dionysios | noun - sing - fem - gen   | οἰκίας |
      | dionysios | noun - sing - fem - dat   | οἰκίᾳ  |
      | dionysios | noun - sing - fem - acc   | τιμήν  |
      | dionysios | noun - plural - fem - nom | μάχαι  |
      | dionysios | noun - plural - fem - gen | θεῶν   |
      | dionysios | noun - plural - fem - dat | δόξαις |
      | dionysios | noun - plural - fem - acc | χώρᾱς  |

  @dionysios
  Scenario Outline: Masculine first declensions result in the correct declension rule
    Given the "<service>" is running
    When the grammar is checked for word "<word>"
    Then the declension "<declension>" should be included in the response
    Examples:
      | service  | declension                 | word     |
      | dionysios | noun - sing - masc - nom   | πολίτης  |
      | dionysios | noun - sing - masc - gen   | κριτοῦ   |
      | dionysios | noun - sing - masc - dat   | νεανίᾳ   |
      | dionysios | noun - sing - masc - acc   | πολίτην  |
      | dionysios | noun - plural - masc - nom | κριταί   |
      | dionysios | noun - plural - masc - gen | πολίτῶν  |
      | dionysios | noun - plural - masc - dat | νεανίαις |
      | dionysios | noun - plural - masc - acc | κριτᾱ́ς  |

  @dionysios
  Scenario Outline: Masculine second declensions result in the correct declension rule
    Given the "<service>" is running
    When the grammar is checked for word "<word>"
    Then the declension "<declension>" should be included in the response
    Examples:
      | service  | declension                 | word     |
      | dionysios | noun - sing - masc - nom   | δοῦλος   |
      | dionysios | noun - sing - masc - gen   | πόλεμου  |
      | dionysios | noun - sing - masc - dat   | δοῦλῳ    |
      | dionysios | noun - sing - masc - acc   | πόλεμον  |
      | dionysios | noun - plural - masc - nom | θεοί     |
      | dionysios | noun - plural - masc - gen | νεανίῶν  |
      | dionysios | noun - plural - masc - dat | πόλεμοις |
      | dionysios | noun - plural - masc - acc | θεούς    |

  @dionysios
  Scenario Outline: Neuter second declensions result in the correct declension rule
    Given the "<service>" is running
    When the grammar is checked for word "<word>"
    Then the declension "<declension>" should be included in the response
    Examples:
      | service  | declension                 | word   |
      | dionysios | noun - sing - neut - nom   | μῆλον  |
      | dionysios | noun - sing - neut - gen   | δῶρου  |
      | dionysios | noun - sing - neut - dat   | δῶρῳ   |
      | dionysios | noun - sing - neut - acc   | μῆλον  |
      | dionysios | noun - plural - neut - nom | δῶρα   |
      | dionysios | noun - plural - neut - gen | δῶρων  |
      | dionysios | noun - plural - neut - dat | μήλοις |
      | dionysios | noun - plural - neut - acc | μῆλα   |

  @dionysios
  Scenario Outline: Queries with no results return an error
    Given the "<service>" is running
    When the grammar for word "<word>" is queried with an error
    Then an error containing "<message>" is returned
    Examples:
      | service  | word             | message         |
      | dionysios | ναυμαχίαναυμαχία | 200 but got 404 |

  @dionysios
  Scenario Outline: Some words have multiple dictionary entries
    Given the "<service>" is running
    When the grammar is checked for word "<word>"
    Then the number of results should be equal to or exceed "<results>"
    And the number of translations should be equal to er exceed "<translations>"
    And the number of declensions should be equal to or exceed "<declensions>"
    Examples:
      | service  | results | translations | declensions | word    |
      | dionysios | 2       | 2            | 1           | πόλεμου |
      | dionysios | 2       | 2            | 1           | μάχη    |

  @dionysios
  Scenario Outline: Some words have multiple declensions
    Given the "<service>" is running
    When the grammar is checked for word "<word>"
    Then the number of results should be equal to or exceed "<results>"
    And the number of translations should be equal to er exceed "<translations>"
    And the number of declensions should be equal to or exceed "<declensions>"
    Examples:
      | service  | results | translations | declensions | word |
      | dionysios | 2       | 1            | 2           | δῶρα |
      | dionysios | 2       | 2            | 2           | θεῶν |
