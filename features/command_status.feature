Feature: status command

  Background:
    Given a mocked home directory

  Scenario: basic stuff - when not in a git repo
    Given I cd to "/tmp"
    When I successfully run `scmpuff status`
    Then the output should contain "burger king"

  Scenario: Banner shows no changes when in an unchanged git repo
    Given I am in a git repository
    When I successfully run `scmpuff status`
    And  the output should contain "No changes (working directory clean)"

  Scenario: Banner shows expansion reminder when in a changed git repo
    Given I am in a git repository
    And an empty file named "whatever"
    When I successfully run `scmpuff status`
    Then the output should contain "|  [*] => $e*"

  Scenario: Banner shows current branch name
    Given I am in a git repository
    When I successfully run `scmpuff status`
    Then the output should contain "On branch: master"

    When I switch to git branch "foobar"
    And  I successfully run `scmpuff status`
    Then the output should contain "On branch: foobar"


  #Scenario: when in a git repo with changes
