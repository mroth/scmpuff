Feature: status command

  Background:
    Given a mocked home directory

  Scenario: basic stuff - when not in a git repo
    #When I successfully run `scmpuff status`
    #Then the output should contain

  Scenario: basic output - when in an unchanged git repo
    Given I am in a git repository
    When I successfully run `scmpuff status`
    Then the output should contain "On branch: master"
    And  the output should contain "No changes (working directory clean)"

  Scenario: basic output - when in an unchanged git repo on a different branch
    Given I am in a git repository
    And I switch to git branch "foobar"
    When I successfully run `scmpuff status`
    Then the output should contain "On branch: foobar"

  #Scenario: when in a git repo with changes
