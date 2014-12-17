Feature: status command

  Background:
    Given a mocked home directory

  @outside-repo
  Scenario: Appropriate error status when not in a git repo
    When I run `scmpuff status`
    Then the exit status should be 128
    #                              ^^^ same as `git status`
    And the output should contain:
      """
      Not a git repository (or any of the parent directories)
      """

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

  Scenario: Banner shows position relative to upstream status
    Given PENDING: need to write this still

  Scenario: Status properly reports all file changes
    # See: http://git.io/IR8qcg for scm_breeze version of test
    Given I am in a complex working tree status matching scm_breeze tests
    When I run `scmpuff status`
    Then the exit status should be 0
      And the output should match / new file: *\[1\] *new_file/
      And the output should match /  deleted: *\[2\] *deleted_file/
      And the output should match / modified: *\[3\] *new_file/
      And the output should match /untracked: *\[4\] *untracked_file/


  Scenario: Status shows relative paths
    Given PENDING: port from scm_breeze
  # TODO: port test_git_status_produces_relative_paths()  from scm_breeze

  Scenario: Status for a complex merge conflict
    Given PENDING: port from scm_breeze
  # TODO: port test_git_status_shortcuts_merge_conflicts()  from scm_breeze

  @wip
  Scenario Outline: Status sets proper environment variables in shell
    Given I am in a complex working tree status matching scm_breeze tests
      And the scmpuff environment variables have been cleared
    When I run `<shell>` interactively
      And I type `eval "$(scmpuff init -s)"`
      And I type "scmpuff_status_shortcuts"
      And I type `echo -e "e1:$e1\ne2:$e2\ne3:$e3\ne4:$e4\ne5:$e5\nEND"`
      And I type "exit"
    Then the output should contain:
      """
      e1:new_file
      e2:deleted_file
      e3:new_file
      e4:untracked_file
      e5:
      END
      """
    Examples:
      | shell |
      | bash  |
      | zsh   |

  @wip
  Scenario Outline: Status clears extra environment variables from before
    Given I am in a complex working tree status matching scm_breeze tests
      And the scmpuff environment variables have been cleared
    When I run `<shell>` interactively
      And I type `eval "$(scmpuff init -s)"`
      And I type "scmpuff_status_shortcuts"
      And I type "git add new_file"
      And I type "git commit -m 'so be it'"
      And I type "scmpuff_status_shortcuts"
      And I type `echo -e "e1:$e1\ne2:$e2\ne3:$e3\ne4:$e4\ne5:$e5\nEND"`
      And I type "exit"
    Then the output should contain:
      """
      e1:deleted_file
      e2:untracked_file
      e3:
      e4:
      e5:
      END
      """
    Examples:
      | shell |
      | bash  |
      | zsh   |
