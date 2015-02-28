Feature: optional wrapping of normal git cmds in the shell
  In order to verify the shell git wrappers work correctly
  I want to make sure they intercept and wrap naked git commands properly

  Background:
    Given a mocked home directory

  Scenario Outline: Wrapped `git add` adds by number and echos status after
    Given I am in a git repository
      And a 4 byte file named "foo.bar"
      And a 4 byte file named "bar.foo"
    When I run `<shell>` interactively
      And I type `eval "$(scmpuff init -ws)"`
      And I type "scmpuff_status_shortcuts"
      And I type "git add 1"
      And I type "exit"
    Then the output should contain:
      """
      # On branch: master  |  [*] => $e*
      #
      ➤ Changes to be committed
      #
      #       new file:  [1] bar.foo
      #
      ➤ Untracked files
      #
      #      untracked:  [2] foo.bar
      #
      """
    Examples:
      | shell |
      | bash  |
      | zsh   |
