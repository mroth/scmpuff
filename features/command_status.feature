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

  #@wip
  #Scenario: Banner shows position relative to upstream status

  @wip
  Scenario: Status properly reports all file changes
    Given I am in a git repository
    And I successfully run the following commands:
      | touch deleted_file          |
      | git add deleted_file        |
      | git commit -m "Test commit" |
      | touch new_file              |
      | touch untracked_file        |
      | git add new_file            |
      | echo "changed" > new_file   |
      | rm deleted_file             |
    When I run `scmpuff status`
    Then the exit status should be 0
      And the output should match / new file: *\[1\] *new_file/
      And the output should match /  deleted: *\[2\] *deleted_file/
      And the output should match / modified: *\[3\] *new_file/
      And the output should match /untracked: *\[4\] *untracked_file/


  #Scenario: status shows relative paths
  # TODO: port test_git_status_produces_relative_paths()  from scm_breeze

  #Scenario: status for a complex merge conflict
  # TODO: port test_git_status_shortcuts_merge_conflicts()  from scm_breeze
