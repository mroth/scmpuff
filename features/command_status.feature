Feature: status command

  Background:
    Given a mocked home directory

  @wip
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

  Scenario: Status properly reports all file changes
    # See: http://git.io/IR8qcg for scm_breeze version of test
    Given I am in a git repository
      And an empty file named "deleted_file"
      And I successfully run `git add deleted_file`
      And I successfully run `git commit -m "Test commit"`
      And an empty file named "new_file"
      And an empty file named "untracked_file"
      And I successfully run `git add new_file`
      And I overwrite "new_file" with:
        """
        changed contents lolol
        """
      And I remove the file "deleted_file"
    When I run `scmpuff status`
    Then the exit status should be 0
      And the output should match / new file: *\[1\] *new_file/
      And the output should match /  deleted: *\[2\] *deleted_file/
      And the output should match / modified: *\[3\] *new_file/
      And the output should match /untracked: *\[4\] *untracked_file/


  Scenario: status shows relative paths
    Given PENDING: port from scm_breeze
  # TODO: port test_git_status_produces_relative_paths()  from scm_breeze

  Scenario: status for a complex merge conflict
    Given PENDING: port from scm_breeze
  # TODO: port test_git_status_shortcuts_merge_conflicts()  from scm_breeze
