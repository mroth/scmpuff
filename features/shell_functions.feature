Feature: scmpuff_status_shortcuts function
  The scmpuff_status_shortcuts shell function wraps the underlying
  `scmpuff status` command, passing along the `--filelist` option and then
  parsing the results to set environment variables in the current shell.

  Background:
    Given a mocked home directory

  @outside-repo
  Scenario Outline: Handle error conditions from wrapped binary command
    It is possible for the underlying `scmpuff status` command wrapped by our
    shell function to produce errors, for example, when not in a git repository.

    In keeping with the design theory of this program (handle as much as
    possible in the binary), we want to make sure those error messages are
    propogated to the user and not swallowed by the shell function, and that
    non-zero exit codes from the underlying process are preserved.

    When I run `<shell>` interactively
      And I type `eval "$(scmpuff init -ws)"`
      And I type "scmpuff_status_shortcuts"
      And I type "exit $?"
    Then the exit status should be 128
    And the output should contain:
      """
      Not a git repository (or any of the parent directories)
      """
    Examples:
      | shell |
      | bash  |
      | zsh   |

  Scenario Outline: Basic functionality works with shell wrapper.
    Given I am in a git repository
    When I run `<shell>` interactively
      And I type `eval "$(scmpuff init -ws)"`
      And I type "scmpuff_status_shortcuts"
      And I type "exit $?"
    Then the exit status should be 0
    And the output should contain "No changes (working directory clean)"
    Examples:
      | shell |
      | bash  |
      | zsh   |

  Scenario Outline: Sets proper environment variables in shell
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

  Scenario Outline: Clears extra environment variables from before
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
